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
