// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// Pod cleanup policies select when a completed job Pod (and its per-job Secret
// and ConfigMap) is deleted. Retained pods keep the backend label and survive
// the startup dangling-owner sweep via the retain marker, so they remain
// available for kubectl-exec auditing until manually removed.
const (
	// KubernetesPodCleanupAlways deletes the pod after every job (default).
	KubernetesPodCleanupAlways = "always"
	// KubernetesPodCleanupOnFailure deletes the pod only when the job succeeded,
	// keeping failed pods for debugging.
	KubernetesPodCleanupOnFailure = "on-failure"
	// KubernetesPodCleanupNever keeps every pod for full audit retention.
	KubernetesPodCleanupNever = "never"
)

// KubernetesResourceList is the cpu/memory pair for a single requests or limits
// entry. Values are Kubernetes resource.Quantity strings (e.g. "250m", "256Mi").
type KubernetesResourceList struct {
	CPU    string
	Memory string
}

// KubernetesResources holds the requests and limits for one container role.
type KubernetesResources struct {
	Requests KubernetesResourceList
	Limits   KubernetesResourceList
}

// KubernetesRoleResources carries per-container-role resource requirements. Each
// role (job container, service sidecar, dind sidecar) is tuned independently
// because their footprints differ markedly (dind needs >=1Gi).
type KubernetesRoleResources struct {
	Job     KubernetesResources
	Service KubernetesResources
	Dind    KubernetesResources
}

// KubernetesWorkspaceVolume selects how the shared job workspace is backed:
// "emptyDir" (default), "ephemeral" (a pod-owned generic ephemeral volume), or
// "pvc" (a static claim named after the pod). Size/StorageClassName/AccessModes
// apply to the claim-backed types.
type KubernetesWorkspaceVolume struct {
	Type             string
	Size             string
	StorageClassName string
	AccessModes      []string
}

// KubernetesDindConfig configures the optional on-demand docker-in-docker
// sidecar used to serve `uses: docker://` container steps.
type KubernetesDindConfig struct {
	Enabled  bool
	Rootless bool
	Image    string
}

// KubernetesConfig is the runtime view of the Kubernetes backend settings that
// the act engine and the pod builder consume. It is the framework-agnostic
// counterpart to the YAML-facing config.Kubernetes struct, kept separate so the
// vendored act packages do not import the runner's config package.
type KubernetesConfig struct {
	Enabled            bool
	Namespace          string
	Kubeconfig         string
	ServiceAccount     string
	DefaultJobImage    string
	ImagePullSecrets   []string
	WorkspaceVolume    KubernetesWorkspaceVolume
	PodLabels          map[string]string
	PodAnnotations     map[string]string
	Resources          KubernetesRoleResources
	PodReadyTimeout    time.Duration
	ConnectBackAddress string
	Dind               KubernetesDindConfig
	PodCleanup         string

	// RuntimeClassName is the cluster RuntimeClass applied to every job pod
	// (e.g. "gvisor" or "kata"), selecting a sandboxed container runtime in place
	// of the node default (runc). It is an operator/cluster policy enforced on
	// every job this runner accepts, so a workflow cannot opt out of isolation;
	// route different trust tiers to different runners via runs-on. Empty leaves
	// spec.runtimeClassName unset, so the cluster's default runtime is used.
	RuntimeClassName string

	// NodeSelector constrains job pods to nodes matching all key-value pairs.
	NodeSelector map[string]string
	// Tolerations allow job pods to schedule on tainted nodes.
	Tolerations []corev1.Toleration
	// Affinity is a corev1.Affinity object applied to the job pod.
	Affinity *corev1.Affinity
}
