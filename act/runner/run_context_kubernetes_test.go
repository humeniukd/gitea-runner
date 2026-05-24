// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"testing"

	"gitea.com/gitea/runner/act/container"

	"github.com/stretchr/testify/assert"
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
