// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeConfig(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
	return path
}

func TestLoadDefault_KubernetesDefaults(t *testing.T) {
	cfg, err := LoadDefault("")
	require.NoError(t, err)

	k := cfg.Kubernetes
	assert.False(t, k.Enabled)
	assert.Equal(t, "gitea-e2e", k.Namespace)
	assert.Equal(t, 5*time.Minute, k.PodReadyTimeout)
	assert.Equal(t, "emptyDir", k.WorkspaceVolume.Type)

	// dind.enabled defaults to false; rootless defaults to true.
	assert.False(t, k.Dind.Enabled)
	require.NotNil(t, k.Dind.Rootless)
	assert.True(t, *k.Dind.Rootless)

	// Per-role resource defaults (M3).
	assert.Equal(t, "250m", k.Resources.Job.Requests.CPU)
	assert.Equal(t, "256Mi", k.Resources.Job.Requests.Memory)
	assert.Equal(t, "2", k.Resources.Job.Limits.CPU)
	assert.Equal(t, "2Gi", k.Resources.Job.Limits.Memory)

	assert.Equal(t, "100m", k.Resources.Service.Requests.CPU)
	assert.Equal(t, "128Mi", k.Resources.Service.Requests.Memory)
	assert.Equal(t, "1", k.Resources.Service.Limits.CPU)
	assert.Equal(t, "1Gi", k.Resources.Service.Limits.Memory)

	assert.Equal(t, "500m", k.Resources.Dind.Requests.CPU)
	assert.Equal(t, "1Gi", k.Resources.Dind.Requests.Memory)
	assert.Equal(t, "2", k.Resources.Dind.Limits.CPU)
	assert.Equal(t, "4Gi", k.Resources.Dind.Limits.Memory)
}

func TestLoadDefault_KubernetesEnabledValidBlock(t *testing.T) {
	path := writeConfig(t, `
kubernetes:
  enabled: true
  namespace: ci
  service_account: gitea-runner
  kubeconfig: /etc/kube/config
  image_pull_secrets:
    - regcred
  workspace_volume:
    type: pvc
    pvc_size: 10Gi
  pod_labels:
    team: ci
  pod_annotations:
    note: hello
  pod_ready_timeout: 90s
  connect_back_address: http://gitea-runner-cache.gitea-e2e.svc:8088
  resources:
    job:
      requests:
        cpu: 500m
        memory: 512Mi
      limits:
        cpu: 4
        memory: 4Gi
  dind:
    enabled: true
    rootless: false
    image: docker:27-dind
`)

	cfg, err := LoadDefault(path)
	require.NoError(t, err)

	k := cfg.Kubernetes
	assert.True(t, k.Enabled)
	assert.Equal(t, "ci", k.Namespace)
	assert.Equal(t, "gitea-runner", k.ServiceAccount)
	assert.Equal(t, "/etc/kube/config", k.Kubeconfig)
	assert.Equal(t, []string{"regcred"}, k.ImagePullSecrets)
	assert.Equal(t, "pvc", k.WorkspaceVolume.Type)
	assert.Equal(t, "10Gi", k.WorkspaceVolume.PVCSize)
	assert.Equal(t, "ci", k.PodLabels["team"])
	assert.Equal(t, "hello", k.PodAnnotations["note"])
	assert.Equal(t, 90*time.Second, k.PodReadyTimeout)
	assert.Equal(t, "http://gitea-runner-cache.gitea-e2e.svc:8088", k.ConnectBackAddress)

	// Operator overrides preserved.
	assert.Equal(t, "500m", k.Resources.Job.Requests.CPU)
	assert.Equal(t, "4", k.Resources.Job.Limits.CPU)
	// Unspecified roles still receive their M3 defaults.
	assert.Equal(t, "500m", k.Resources.Dind.Requests.CPU)

	assert.True(t, k.Dind.Enabled)
	require.NotNil(t, k.Dind.Rootless)
	assert.False(t, *k.Dind.Rootless)
	assert.Equal(t, "docker:27-dind", k.Dind.Image)
}

func TestLoadDefault_KubernetesPVCRequiresSize(t *testing.T) {
	path := writeConfig(t, `
kubernetes:
  enabled: true
  workspace_volume:
    type: pvc
`)

	_, err := LoadDefault(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pvc_size")
}

func TestLoadDefault_KubernetesInvalidWorkspaceType(t *testing.T) {
	path := writeConfig(t, `
kubernetes:
  enabled: true
  workspace_volume:
    type: nfs
`)

	_, err := LoadDefault(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "workspace_volume.type")
}

func TestLoadDefault_KubernetesDisabledSkipsValidation(t *testing.T) {
	// An invalid block must not break a runner that does not use the backend.
	path := writeConfig(t, `
kubernetes:
  enabled: false
  workspace_volume:
    type: nfs
`)

	_, err := LoadDefault(path)
	require.NoError(t, err)
}

func TestLoadDefault_KubernetesRootlessExplicitFalse(t *testing.T) {
	path := writeConfig(t, `
kubernetes:
  enabled: true
  dind:
    enabled: true
    rootless: false
`)

	cfg, err := LoadDefault(path)
	require.NoError(t, err)
	require.NotNil(t, cfg.Kubernetes.Dind.Rootless)
	assert.False(t, *cfg.Kubernetes.Dind.Rootless)
}

func TestLoadDefault_DisabledKubernetesLeavesContainerDefaults(t *testing.T) {
	// With the backend off, unrelated defaults remain untouched.
	cfg, err := LoadDefault("")
	require.NoError(t, err)
	assert.Equal(t, "workspace", cfg.Container.WorkdirParent)
	assert.False(t, cfg.Kubernetes.Enabled)
}
