// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"maps"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// kubernetesBackendLabel marks every Pod the kubernetes backend creates so
	// the startup sweep can enumerate candidate job pods. It is a class label;
	// the sweep NEVER deletes by this label alone — it additionally filters on a
	// resolvable ownerReference so live sibling runners' pods are spared (RC2).
	kubernetesBackendLabel = "actions.gitea.com/backend"
	kubernetesBackendValue = "kubernetes"

	// kubernetesJobContainerName is the long-lived container act execs steps
	// into. Service sidecars are named after their service id.
	kubernetesJobContainerName = "job"

	// kubernetesWorkspaceVolumeName backs the shared job workspace.
	kubernetesWorkspaceVolumeName = "gitea-workspace"

	kubernetesWorkspaceTypePVC = "pvc"
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
		Volumes:               []corev1.Volume{b.workspaceVolumeSpec()},
		ActiveDeadlineSeconds: b.activeDeadlineSeconds,
		HostAliases:           b.hostAliases,
	}
	if b.serviceAccount != "" {
		spec.ServiceAccountName = b.serviceAccount
	}
	for _, name := range b.imagePullSecrets {
		spec.ImagePullSecrets = append(spec.ImagePullSecrets, corev1.LocalObjectReference{Name: name})
	}

	return &corev1.Pod{ObjectMeta: meta, Spec: spec}
}

// workspaceVolumeSpec returns the shared workspace volume: a PVC-backed volume
// when configured, otherwise an emptyDir (default).
func (b *podBuilder) workspaceVolumeSpec() corev1.Volume {
	if b.workspaceVolume.Type == kubernetesWorkspaceTypePVC {
		return corev1.Volume{
			Name: kubernetesWorkspaceVolumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: b.name,
				},
			},
		}
	}
	return corev1.Volume{
		Name:         kubernetesWorkspaceVolumeName,
		VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
	}
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
