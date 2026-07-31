// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"testing"

	"gitea.com/gitea/runner/internal/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
	corev1 "k8s.io/api/core/v1"
)

// TestKubernetesRunnerConfigAffinityRoundTrip guards against a regression where
// the kubernetes.affinity / kubernetes.tolerations config blocks were silently
// dropped: corev1 types carry json tags (e.g. nodeAffinity), so the generic
// yaml.Unmarshal produced a nil Affinity while swallowing the parse result.
// sigs.k8s.io/yaml (YAML -> JSON -> struct) is required for a correct round trip.
func TestKubernetesRunnerConfigAffinityRoundTrip(t *testing.T) {
	cfgYaml := `
kubernetes:
  enabled: true
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
              - key: node-role.kubernetes.io/control-plane
                operator: NotIn
                values:
                  - "true"
  tolerations:
    - key: node-role.kubernetes.io/control-plane
      operator: Exists
      effect: NoSchedule
`
	cfg := &config.Config{}
	require.NoError(t, yaml.Unmarshal([]byte(cfgYaml), cfg))

	k8sConfig, err := KubernetesRunnerConfig(&cfg.Kubernetes)
	require.NoError(t, err)
	require.NotNil(t, k8sConfig)
	require.NotNil(t, k8sConfig.Affinity, "affinity must survive the config round trip")
	require.NotNil(t, k8sConfig.Affinity.NodeAffinity, "nodeAffinity must survive the config round trip")

	terms := k8sConfig.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	require.Len(t, terms, 1)
	expr := terms[0].MatchExpressions[0]
	assert.Equal(t, "node-role.kubernetes.io/control-plane", expr.Key)
	assert.Equal(t, corev1.NodeSelectorOpNotIn, expr.Operator)
	assert.Equal(t, []string{"true"}, expr.Values)

	require.Len(t, k8sConfig.Tolerations, 1)
	assert.Equal(t, "node-role.kubernetes.io/control-plane", k8sConfig.Tolerations[0].Key)
	assert.Equal(t, corev1.TolerationOpExists, k8sConfig.Tolerations[0].Operator)
	assert.Equal(t, corev1.TaintEffectNoSchedule, k8sConfig.Tolerations[0].Effect)
}

// TestKubernetesRunnerConfigDisabled verifies a non-enabled kubernetes section
// returns a nil config rather than erroring.
func TestKubernetesRunnerConfigDisabled(t *testing.T) {
	cfg := &config.Config{}
	k8sConfig, err := KubernetesRunnerConfig(&cfg.Kubernetes)
	require.NoError(t, err)
	assert.Nil(t, k8sConfig)
}
