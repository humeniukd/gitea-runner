// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"crypto/sha256"
	"encoding/hex"
	"maps"
	"sort"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// kubernetesBackendLabel marks every Pod the kubernetes backend creates so
	// the startup sweep can enumerate candidate job pods. It is a class label;
	// the sweep NEVER deletes by this label alone — it additionally filters on a
	// resolvable ownerReference so live sibling runners' pods are spared (RC2).
	kubernetesBackendLabel = "actions.gitea.com/backend"
	kubernetesBackendValue = "kubernetes"

	// kubernetesRetainLabel marks a completed job Pod that the pod_cleanup policy
	// chose to KEEP for auditing. The startup dangling-owner sweep skips any pod
	// carrying this marker so a runner restart cannot reap a retained pod (which
	// would defeat the audit). The pod is retained until manually deleted.
	kubernetesRetainLabel = "actions.gitea.com/retain"
	kubernetesRetainValue = "true"

	// kubernetesJobContainerName is the long-lived container act execs steps
	// into. Service sidecars are named after their service id.
	kubernetesJobContainerName = "job"

	// kubernetesWorkspaceVolumeName backs the shared job workspace.
	kubernetesWorkspaceVolumeName = "gitea-workspace"

	// kubernetesActPathVolumeName backs the act runtime directory (/var/run/act/)
	kubernetesActPathVolumeName = "gitea-act-path"

	kubernetesWorkspaceTypePVC       = "pvc"
	kubernetesWorkspaceTypeEphemeral = "ephemeral"

	// dind sidecar wiring: the daemon listens on TCP 2375 (TLS disabled) and the
	// job container reaches it at localhost:2375 via the shared pod network.
	kubernetesDindContainerName   = "dind"
	kubernetesDindPort            = 2375
	kubernetesDindHost            = "tcp://localhost:2375"
	kubernetesDindImageRootless   = "docker:dind-rootless"
	kubernetesDindImagePrivileged = "docker:dind"

	// defaultWorkspaceAccessMode backs claim-based workspaces when none is set.
	defaultWorkspaceAccessMode = corev1.ReadWriteOnce
)

// podBuilder assembles the single corev1.Pod that backs one Actions job: a
// long-lived job container plus optional native service sidecars (init
// containers with restartPolicy=Always). It is pure — Build() never touches the
// API server — so it is fully unit-testable with a fake clientset.
type podBuilder struct {
	name             string
	namespace        string
	labels           map[string]string
	annotations      map[string]string
	serviceAccount   string
	imagePullSecrets []string
	workspaceVolume  KubernetesWorkspaceVolume

	activeDeadlineSeconds *int64
	ownerReference        *metav1.OwnerReference
	hostAliases           []corev1.HostAlias
	runtimeClassName      string
	nodeSelector          map[string]string
	tolerations           []corev1.Toleration
	affinity              *corev1.Affinity

	jobContainer   *corev1.Container
	sidecars       []corev1.Container
	initContainers []corev1.Container
}

// newPodBuilder seeds a builder with the pod identity and the orphan-reaping
// primitives (activeDeadlineSeconds bounded by maxLifetime, ownerReference at
// THIS runner's own Pod). A non-positive maxLifetime leaves the deadline unset.
func newPodBuilder(name, namespace string, cfg *KubernetesConfig, maxLifetime time.Duration, owner *metav1.OwnerReference) *podBuilder {
	labels := map[string]string{kubernetesBackendLabel: kubernetesBackendValue}
	maps.Copy(labels, cfg.PodLabels)

	var annotations map[string]string
	if len(cfg.PodAnnotations) > 0 {
		annotations = make(map[string]string, len(cfg.PodAnnotations))
		maps.Copy(annotations, cfg.PodAnnotations)
	}

	var deadline *int64
	if maxLifetime > 0 {
		secs := int64(maxLifetime.Round(time.Second).Seconds())
		if secs > 0 {
			deadline = &secs
		}
	}

	return &podBuilder{
		name:                  name,
		namespace:             namespace,
		labels:                labels,
		annotations:           annotations,
		serviceAccount:        cfg.ServiceAccount,
		imagePullSecrets:      cfg.ImagePullSecrets,
		workspaceVolume:       cfg.WorkspaceVolume,
		activeDeadlineSeconds: deadline,
		ownerReference:        owner,
		runtimeClassName:      cfg.RuntimeClassName,
		nodeSelector:          cfg.NodeSelector,
		tolerations:           cfg.Tolerations,
		affinity:              cfg.Affinity,
	}
}

func (b *podBuilder) setJobContainer(c corev1.Container) {
	b.jobContainer = &c
}

// addSidecar registers a service container as a native sidecar (an init
// container with restartPolicy=Always, guaranteed on k0s v1.35).
func (b *podBuilder) addSidecar(c corev1.Container) {
	always := corev1.ContainerRestartPolicyAlways
	c.RestartPolicy = &always
	b.sidecars = append(b.sidecars, c)
}

// addHostAlias maps hostnames to an IP in the Pod's /etc/hosts so a job step can
// reach a service by name (servicename:port) in addition to localhost:port.
func (b *podBuilder) addHostAlias(ip string, hostnames ...string) {
	b.hostAliases = append(b.hostAliases, corev1.HostAlias{IP: ip, Hostnames: hostnames})
}

// addDindSidecar injects an on-demand docker-in-docker native sidecar so a
// `uses: docker://` step can build/run images. The job container reaches it at
// localhost:2375 (DOCKER_HOST). The same workspace volume is mounted so a
// `docker run -v <workdir>:<workdir>` resolves against shared files. rootless
// selects docker:dind-rootless; otherwise a privileged docker:dind is used.
func (b *podBuilder) addDindSidecar(cfg *KubernetesConfig, workdir string) {
	rootless := cfg.Dind.Rootless
	image := cfg.Dind.Image
	if image == "" {
		if rootless {
			image = kubernetesDindImageRootless
		} else {
			image = kubernetesDindImagePrivileged
		}
	}

	probe := tcpPortProbe(kubernetesDindPort)
	c := corev1.Container{
		Name:  kubernetesDindContainerName,
		Image: image,
		// Disable TLS so the daemon serves plain TCP on 2375.
		Env:       []corev1.EnvVar{{Name: "DOCKER_TLS_CERTDIR", Value: ""}},
		Resources: resourceRequirements(cfg.Resources.Dind),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      kubernetesWorkspaceVolumeName,
			MountPath: workspaceMountPath(workdir),
		}},
		Ports:          []corev1.ContainerPort{{ContainerPort: int32(kubernetesDindPort)}},
		StartupProbe:   probe,
		ReadinessProbe: probe,
	}
	if !rootless {
		privileged := true
		c.SecurityContext = &corev1.SecurityContext{Privileged: &privileged}
	}
	b.addSidecar(c)
}

// sidecarContainerNames returns the names of the registered service sidecars,
// used by the readiness gate.
func (b *podBuilder) sidecarContainerNames() []string {
	names := make([]string, 0, len(b.sidecars))
	for _, s := range b.sidecars {
		names = append(names, s.Name)
	}
	return names
}

// build materialises the corev1.Pod. Native sidecars precede any plain init
// containers so they are started (and kept running) before the job container.
func (b *podBuilder) build() *corev1.Pod {
	meta := metav1.ObjectMeta{
		Name:        b.name,
		Namespace:   b.namespace,
		Labels:      b.labels,
		Annotations: b.annotations,
	}
	if b.ownerReference != nil {
		meta.OwnerReferences = []metav1.OwnerReference{*b.ownerReference}
	}

	initContainers := make([]corev1.Container, 0, len(b.sidecars)+len(b.initContainers))
	initContainers = append(initContainers, b.sidecars...)
	initContainers = append(initContainers, b.initContainers...)

	var containers []corev1.Container
	if b.jobContainer != nil {
		containers = []corev1.Container{*b.jobContainer}
	}

	spec := corev1.PodSpec{
		RestartPolicy:         corev1.RestartPolicyNever,
		Containers:            containers,
		InitContainers:        initContainers,
		Volumes: []corev1.Volume{
			b.workspaceVolumeSpec(),
			{Name: kubernetesActPathVolumeName, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
		},
		ActiveDeadlineSeconds: b.activeDeadlineSeconds,
		HostAliases:           b.hostAliases,
	}
	// runtimeClassName selects a sandboxed runtime (gVisor/Kata) for the whole
	// pod. Left unset when empty so the cluster's default runtime applies.
	if b.runtimeClassName != "" {
		runtimeClass := b.runtimeClassName
		spec.RuntimeClassName = &runtimeClass
	}
	if b.serviceAccount != "" {
		spec.ServiceAccountName = b.serviceAccount
	}
	if len(b.nodeSelector) > 0 {
		spec.NodeSelector = b.nodeSelector
	}
	if len(b.tolerations) > 0 {
		spec.Tolerations = b.tolerations
	}
	if b.affinity != nil {
		spec.Affinity = b.affinity
	}
	for _, name := range b.imagePullSecrets {
		spec.ImagePullSecrets = append(spec.ImagePullSecrets, corev1.LocalObjectReference{Name: name})
	}

	return &corev1.Pod{ObjectMeta: meta, Spec: spec}
}

// workspaceVolumeSpec returns the shared workspace volume. "ephemeral" yields a
// pod-owned generic ephemeral volume (a PVC auto-created and GC'd with the pod);
// "pvc" references a static claim named after the pod; otherwise an emptyDir
// (default).
func (b *podBuilder) workspaceVolumeSpec() corev1.Volume {
	switch b.workspaceVolume.Type {
	case kubernetesWorkspaceTypeEphemeral:
		return corev1.Volume{
			Name: kubernetesWorkspaceVolumeName,
			VolumeSource: corev1.VolumeSource{
				Ephemeral: &corev1.EphemeralVolumeSource{
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimTemplate{
						Spec: b.workspaceClaimSpec(),
					},
				},
			},
		}
	case kubernetesWorkspaceTypePVC:
		return corev1.Volume{
			Name: kubernetesWorkspaceVolumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: b.name,
				},
			},
		}
	default:
		return corev1.Volume{
			Name:         kubernetesWorkspaceVolumeName,
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		}
	}
}

// workspaceClaimSpec builds the PVC spec for the generic ephemeral volume from
// the configured size, storage class, and access modes (defaulting to
// ReadWriteOnce).
func (b *podBuilder) workspaceClaimSpec() corev1.PersistentVolumeClaimSpec {
	accessModes := make([]corev1.PersistentVolumeAccessMode, 0, len(b.workspaceVolume.AccessModes))
	for _, mode := range b.workspaceVolume.AccessModes {
		accessModes = append(accessModes, corev1.PersistentVolumeAccessMode(mode))
	}
	if len(accessModes) == 0 {
		accessModes = []corev1.PersistentVolumeAccessMode{defaultWorkspaceAccessMode}
	}

	spec := corev1.PersistentVolumeClaimSpec{
		AccessModes: accessModes,
		Resources: corev1.VolumeResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(b.workspaceVolume.Size),
			},
		},
	}
	if b.workspaceVolume.StorageClassName != "" {
		storageClass := b.workspaceVolume.StorageClassName
		spec.StorageClassName = &storageClass
	}
	return spec
}

// kubernetesContainerSpec builds a corev1.Container from an act container input,
// applying the workspace mount, per-role resources, security context, env, and
// image-pull policy. workdir is the mount path for the shared workspace.
func kubernetesContainerSpec(name string, input *NewContainerInput, res KubernetesResources, capAdd, capDrop []string, forcePull bool) corev1.Container {
	c := corev1.Container{
		Name:       name,
		Image:      input.Image,
		WorkingDir: input.WorkingDir,
		Env:        envVarsFromList(input.Env),
		Resources:  resourceRequirements(res),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      kubernetesActPathVolumeName,
			MountPath: "/var/run/act",
		}, {
			Name:      kubernetesWorkspaceVolumeName,
			MountPath: workspaceMountPath(input.WorkingDir),
		}},
		ImagePullPolicy: imagePullPolicy(forcePull),
	}
	if len(input.Entrypoint) > 0 {
		c.Command = input.Entrypoint
	}
	if len(input.Cmd) > 0 {
		c.Args = input.Cmd
	}
	if sc := securityContext(input.Privileged, capAdd, capDrop); sc != nil {
		c.SecurityContext = sc
	}
	return c
}

// serviceSidecarSpec builds a service container with its exposed ports declared
// and TCP startup/readiness probes on the first port (C3). The probes let the
// readiness gate (and native-sidecar ordering) hold the job container's first
// step until the service is actually accepting connections, eliminating the
// connection-refused flake. A service with no declared ports gets no probe and
// is considered ready once running.
func serviceSidecarSpec(input *NewContainerInput, res KubernetesResources, capAdd, capDrop []string, forcePull bool) corev1.Container {
	// Use the workflow service key (NetworkAliases[0]) as the sidecar container
	// name so it stays human-readable within the pod scope.  The compound
	// Docker-style name stored in input.Name is too long for Kubernetes and
	// becomes an opaque hash after RFC 1123 sanitization.
	sidecarName := input.Name
	if len(input.NetworkAliases) > 0 {
		sidecarName = input.NetworkAliases[0]
	}
	c := kubernetesContainerSpec(serviceContainerName(sidecarName), input, res, capAdd, capDrop, forcePull)

	ports := sortedExposedPorts(input.ExposedPorts)
	for _, port := range ports {
		c.Ports = append(c.Ports, corev1.ContainerPort{ContainerPort: int32(port)})
	}
	if len(ports) > 0 {
		probe := tcpPortProbe(ports[0])
		c.StartupProbe = probe
		c.ReadinessProbe = probe
	}
	return c
}

// tcpPortProbe returns a TCP socket probe that polls every 2s and tolerates a
// slow-starting service for up to ~2 minutes before giving up.
func tcpPortProbe(port int) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{Port: intstr.FromInt(port)},
		},
		PeriodSeconds:    2,
		TimeoutSeconds:   2,
		FailureThreshold: 60,
	}
}

// sortedExposedPorts returns the numeric container ports from an act PortSet in
// ascending order for deterministic output.
func sortedExposedPorts(exposed nat.PortSet) []int {
	ports := make([]int, 0, len(exposed))
	for p := range exposed {
		ports = append(ports, p.Int())
	}
	sort.Ints(ports)
	return ports
}

// workspaceMountPath returns where the shared workspace volume is mounted. The
// job's configured working directory is used when present so step output lands
// on the shared volume; otherwise a sensible default root is used.
func workspaceMountPath(workdir string) string {
	if workdir != "" {
		return workdir
	}
	return "/workspace"
}

func imagePullPolicy(forcePull bool) corev1.PullPolicy {
	if forcePull {
		return corev1.PullAlways
	}
	return corev1.PullIfNotPresent
}

// securityContext maps act's privileged flag and cap add/drop lists onto a
// container SecurityContext, returning nil when there is nothing to set.
func securityContext(privileged bool, capAdd, capDrop []string) *corev1.SecurityContext {
	if !privileged && len(capAdd) == 0 && len(capDrop) == 0 {
		return nil
	}
	sc := &corev1.SecurityContext{}
	if privileged {
		p := true
		sc.Privileged = &p
	}
	if len(capAdd) > 0 || len(capDrop) > 0 {
		caps := &corev1.Capabilities{}
		for _, a := range capAdd {
			caps.Add = append(caps.Add, corev1.Capability(a))
		}
		for _, d := range capDrop {
			caps.Drop = append(caps.Drop, corev1.Capability(d))
		}
		sc.Capabilities = caps
	}
	return sc
}

// envVarsFromList converts act's "KEY=VALUE" env slice into corev1.EnvVars,
// preserving order and tolerating values that themselves contain '='.
func envVarsFromList(env []string) []corev1.EnvVar {
	if len(env) == 0 {
		return nil
	}
	vars := make([]corev1.EnvVar, 0, len(env))
	for _, e := range env {
		key, value := splitEnv(e)
		if key == "" {
			continue
		}
		vars = append(vars, corev1.EnvVar{Name: key, Value: value})
	}
	return vars
}

func splitEnv(entry string) (key, value string) {
	for i := 0; i < len(entry); i++ {
		if entry[i] == '=' {
			return entry[:i], entry[i+1:]
		}
	}
	return entry, ""
}

func resourceRequirements(r KubernetesResources) corev1.ResourceRequirements {
	req := corev1.ResourceRequirements{}
	if list := resourceList(r.Requests); len(list) > 0 {
		req.Requests = list
	}
	if list := resourceList(r.Limits); len(list) > 0 {
		req.Limits = list
	}
	return req
}

// resourceList parses cpu/memory quantity strings, silently skipping empty or
// unparseable entries (config validation owns surfacing bad quantities).
func resourceList(rl KubernetesResourceList) corev1.ResourceList {
	list := corev1.ResourceList{}
	if rl.CPU != "" {
		if q, err := resource.ParseQuantity(rl.CPU); err == nil {
			list[corev1.ResourceCPU] = q
		}
	}
	if rl.Memory != "" {
		if q, err := resource.ParseQuantity(rl.Memory); err == nil {
			list[corev1.ResourceMemory] = q
		}
	}
	return list
}

// ownerReferenceFromEnv builds an ownerReference pointing at THIS runner's own
// Pod, resolved from the downward-API POD_NAME/POD_UID env the Deployment sets
// (RC2). It returns nil out-of-cluster (env unset) so dev runs are unaffected.
func ownerReferenceFromEnv(podName, podUID string) *metav1.OwnerReference {
	if podName == "" || podUID == "" {
		return nil
	}
	controller := true
	blockDeletion := true
	return &metav1.OwnerReference{
		APIVersion:         "v1",
		Kind:               "Pod",
		Name:               podName,
		UID:                types.UID(podUID),
		Controller:         &controller,
		BlockOwnerDeletion: &blockDeletion,
	}
}

// perJobSecretName derives a deterministic per-job Secret name from the pod name
// so cleanup can target it even after a restart.
func perJobSecretName(podName string) string {
	return podName + "-secrets"
}

// perJobConfigMapName derives a deterministic per-job ConfigMap name from the
// pod name so the dangling-owner sweep can target it even after a restart.
func perJobConfigMapName(podName string) string {
	return podName + "-config"
}

// kubernetesNameMaxLen is the RFC 1123 label limit (63). Pod names may be up to
// 253, but constraining them to a single valid label is a strict subset that is
// always acceptable and leaves ample room for the per-job Secret/ConfigMap
// suffixes; container names MUST be <=63.
const kubernetesNameMaxLen = 63

// sanitizeRFC1123Name converts an arbitrary act container/pod name (which may be
// mixed-case and ~128 chars, e.g. "GITEA-ACTIONS-TASK-1-WORKFLOW-CI-JOB-build-...-<sha>")
// into a valid RFC 1123 label: lowercase [a-z0-9-], starting and ending
// alphanumeric, at most kubernetesNameMaxLen (63) chars. When truncation is
// needed a short hash of the ORIGINAL name is appended so distinct inputs stay
// distinct (preserving per-service uniqueness). The API server rejects
// Pods()/containers whose names violate this, so every name fed to the cluster
// must pass through here.
func sanitizeRFC1123Name(name string) string {
	lowered := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			return r
		case r >= 'A' && r <= 'Z':
			return r + ('a' - 'A')
		default:
			return '-'
		}
	}, name)

	sanitized := strings.Trim(collapseDashes(lowered), "-")
	if sanitized == "" {
		sanitized = "x"
	}

	if len(sanitized) > kubernetesNameMaxLen {
		suffix := "-" + shortNameHash(name)
		keep := max(kubernetesNameMaxLen-len(suffix), 1)
		sanitized = strings.TrimRight(sanitized[:keep], "-") + suffix
	}
	return sanitized
}

// collapseDashes replaces runs of '-' with a single '-'.
func collapseDashes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevDash := false
	for _, r := range s {
		if r == '-' {
			if prevDash {
				continue
			}
			prevDash = true
		} else {
			prevDash = false
		}
		b.WriteRune(r)
	}
	return b.String()
}

// shortNameHash returns the first 8 hex chars of the SHA-256 of name — a stable,
// collision-resistant suffix for disambiguating truncated names.
func shortNameHash(name string) string {
	sum := sha256.Sum256([]byte(name))
	return hex.EncodeToString(sum[:])[:8]
}
