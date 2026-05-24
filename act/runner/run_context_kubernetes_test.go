// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"testing"

	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/act/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsesKubernetes(t *testing.T) {
	t.Run("nil config disables the backend", func(t *testing.T) {
		rc := &RunContext{Config: &Config{}}
		assert.False(t, rc.usesKubernetes())
	})

	t.Run("present but disabled selects docker/host", func(t *testing.T) {
		rc := &RunContext{Config: &Config{Kubernetes: &container.KubernetesConfig{Enabled: false}}}
		assert.False(t, rc.usesKubernetes())
	})

	t.Run("enabled selects the pod backend", func(t *testing.T) {
		rc := &RunContext{Config: &Config{Kubernetes: &container.KubernetesConfig{Enabled: true}}}
		assert.True(t, rc.usesKubernetes())
	})
}

// TestGuardKubernetesActionUnsupported is the RC1 unit coverage: under the
// kubernetes backend a `using: docker` action is rejected with the explicit
// unsupported-action error (no socket dial / nil-deref), while JS, composite,
// and Go actions are allowed (the guard returns nil, so dispatch falls through
// to rc.execJobContainer / execAsComposite). The guard is also inert when the
// backend is disabled.
func TestGuardKubernetesActionUnsupported(t *testing.T) {
	k8sRC := &RunContext{Config: &Config{Kubernetes: &container.KubernetesConfig{Enabled: true}}}

	t.Run("rejects using: docker under kubernetes", func(t *testing.T) {
		err := k8sRC.guardKubernetesActionUnsupported(model.ActionRunsUsingDocker, "my-org/my-action")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "using: docker")
		assert.Contains(t, err.Error(), "not supported")
		assert.Contains(t, err.Error(), "my-org/my-action")
	})

	t.Run("allows JS, composite, and Go actions under kubernetes", func(t *testing.T) {
		for _, using := range []model.ActionRunsUsing{
			model.ActionRunsUsingNode16,
			model.ActionRunsUsingNode20,
			model.ActionRunsUsingNode24,
			model.ActionRunsUsingComposite,
			model.ActionRunsUsingGo,
		} {
			assert.NoError(t, k8sRC.guardKubernetesActionUnsupported(using, "a"),
				"using %q must run via exec-in, not be rejected", using)
		}
	})

	t.Run("inert when the kubernetes backend is disabled", func(t *testing.T) {
		for _, rc := range []*RunContext{
			{Config: &Config{}},
			{Config: &Config{Kubernetes: &container.KubernetesConfig{Enabled: false}}},
		} {
			assert.NoError(t, rc.guardKubernetesActionUnsupported(model.ActionRunsUsingDocker, "a"))
		}
	})
}

func TestJobHasDockerURLStep(t *testing.T) {
	newRC := func(steps []*model.Step) *RunContext {
		return &RunContext{Run: &model.Run{
			JobID: "j",
			Workflow: &model.Workflow{
				Jobs: map[string]*model.Job{"j": {Steps: steps}},
			},
		}}
	}

	t.Run("job with a uses: docker:// step", func(t *testing.T) {
		rc := newRC([]*model.Step{
			{Run: "echo hi"},
			{Uses: "docker://alpine:3"},
		})
		assert.True(t, rc.jobHasDockerURLStep())
	})

	t.Run("bare run / JS / composite job has none", func(t *testing.T) {
		rc := newRC([]*model.Step{
			{Run: "echo hi"},
			{Uses: "actions/checkout@v4"},
			{Uses: "./local/action"},
		})
		assert.False(t, rc.jobHasDockerURLStep())
	})

	t.Run("nil run is safe", func(t *testing.T) {
		assert.False(t, (&RunContext{}).jobHasDockerURLStep())
	})
}

func TestResolveKubernetesJobImage(t *testing.T) {
	const def = "node:20-bookworm"

	t.Run("explicit/real image wins over the default", func(t *testing.T) {
		got, err := resolveKubernetesJobImage("postgres:16-alpine", def, "build")
		require.NoError(t, err)
		assert.Equal(t, "postgres:16-alpine", got)
	})

	t.Run("bare run job (empty image) falls back to default", func(t *testing.T) {
		got, err := resolveKubernetesJobImage("", def, "build")
		require.NoError(t, err)
		assert.Equal(t, def, got)
	})

	t.Run("self-hosted sentinel falls back to default", func(t *testing.T) {
		got, err := resolveKubernetesJobImage("-self-hosted", def, "build")
		require.NoError(t, err)
		assert.Equal(t, def, got)
	})

	t.Run("empty image and no default is a clear error", func(t *testing.T) {
		_, err := resolveKubernetesJobImage("", "", "build")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "default_job_image")
		assert.Contains(t, err.Error(), "build")
	})

	t.Run("sentinel image and no default is a clear error", func(t *testing.T) {
		_, err := resolveKubernetesJobImage("-self-hosted", "", "build")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "default_job_image")
	})
}
