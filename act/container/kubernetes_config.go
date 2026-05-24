// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import "time"

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

// KubernetesWorkspaceVolume selects how the shared job workspace is backed.
// Type is "emptyDir" (default) or "pvc"; PVCSize is the requested size when
// Type is "pvc".
type KubernetesWorkspaceVolume struct {
	Type    string
	PVCSize string
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
	ImagePullSecrets   []string
	WorkspaceVolume    KubernetesWorkspaceVolume
	PodLabels          map[string]string
	PodAnnotations     map[string]string
	Resources          KubernetesRoleResources
	PodReadyTimeout    time.Duration
	ConnectBackAddress string
	Dind               KubernetesDindConfig
}
