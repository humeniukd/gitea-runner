// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"testing"

	"gitea.com/gitea/runner/act/container"

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

func TestStartKubernetesEnvironmentNotImplemented(t *testing.T) {
	rc := &RunContext{Config: &Config{Kubernetes: &container.KubernetesConfig{Enabled: true}}}
	err := rc.startKubernetesEnvironment()(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}
