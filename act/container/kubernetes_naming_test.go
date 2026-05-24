// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rfc1123Label mirrors the API server's container/label name validation.
var rfc1123Label = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

// rfc1123Subdomain mirrors the API server's Pod/Secret/ConfigMap name validation.
var rfc1123Subdomain = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

func assertValidRFC1123(t *testing.T, name string) {
	t.Helper()
	assert.LessOrEqual(t, len(name), kubernetesNameMaxLen, "name %q exceeds %d", name, kubernetesNameMaxLen)
	assert.NotEmpty(t, name)
	assert.Regexp(t, rfc1123Label, name, "name %q is not a valid RFC 1123 label", name)
}

func TestSanitizeRFC1123Name(t *testing.T) {
	t.Run("lowercases and replaces invalid chars", func(t *testing.T) {
		got := sanitizeRFC1123Name("GITEA-ACTIONS-TASK-1_WORKFLOW.CI")
		assert.Equal(t, "gitea-actions-task-1-workflow-ci", got)
		assertValidRFC1123(t, got)
	})

	t.Run("collapses repeated dashes and trims edges", func(t *testing.T) {
		got := sanitizeRFC1123Name("--Foo___Bar--")
		assert.Equal(t, "foo-bar", got)
	})

	t.Run("truncates long names with a stable hash suffix", func(t *testing.T) {
		long := "GITEA-ACTIONS-TASK-1-WORKFLOW-Continuous-Integration-JOB-build-and-test-CALLED-BY-reusable-aaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		got := sanitizeRFC1123Name(long)
		assertValidRFC1123(t, got)
		// deterministic
		assert.Equal(t, got, sanitizeRFC1123Name(long))
		// distinct inputs stay distinct even after truncation
		other := long + "-different"
		assert.NotEqual(t, got, sanitizeRFC1123Name(other))
	})

	t.Run("empty becomes a safe placeholder", func(t *testing.T) {
		assert.Equal(t, "x", sanitizeRFC1123Name(""))
		assert.Equal(t, "x", sanitizeRFC1123Name("---"))
	})
}

// TestKubernetesNamesAreRFC1123 drives realistic act-style mixed-case, >63-char
// names (the shape createContainerName produces: <name>-<64 hex>) through the
// pod coordinator and builder, asserting the resulting Pod name and service
// container names are valid RFC 1123 and unique per service.
func TestKubernetesNamesAreRFC1123(t *testing.T) {
	cfg := testKubernetesConfig()
	rawPodName := "GITEA-ACTIONS-TASK-7-WORKFLOW-Continuous-Integration-JOB-build-1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

	pod := newKubernetesPodWithClientset(nil, nil, KubernetesPodConfig{
		Config:  cfg,
		PodName: rawPodName,
	})

	// Pod name (and the Secret/ConfigMap names derived from it) must be valid.
	assertValidRFC1123(t, pod.builder.name)
	secretName := perJobSecretName(pod.builder.name)
	assert.Regexp(t, rfc1123Subdomain, secretName, "per-job secret name must be a valid RFC 1123 subdomain")
	assert.LessOrEqual(t, len(secretName), 253)

	// Two services with distinct mixed-case, >63-char names must yield distinct,
	// valid container names.
	postgresRaw := createStyleName("postgres")
	redisRaw := createStyleName("redis")
	pgName := serviceContainerName(postgresRaw)
	redisName := serviceContainerName(redisRaw)

	assertValidRFC1123(t, pgName)
	assertValidRFC1123(t, redisName)
	require.NotEqual(t, pgName, redisName, "distinct services must get distinct container names")

	// And they must actually compose into a buildable Pod whose every container
	// name is RFC 1123 valid.
	b := newPodBuilder(pod.builder.name, cfg.Namespace, cfg, 0, nil)
	b.setJobContainer(kubernetesContainerSpec(kubernetesJobContainerName, &NewContainerInput{Image: "alpine:3"}, cfg.Resources.Job, nil, nil, false))
	b.addSidecar(serviceSidecarSpec(&NewContainerInput{Image: "postgres:16", Name: postgresRaw}, cfg.Resources.Service, nil, nil, false))
	built := b.build()

	assertValidRFC1123(t, built.Name)
	for _, c := range built.Spec.Containers {
		assertValidRFC1123(t, c.Name)
	}
	for _, c := range built.Spec.InitContainers {
		assertValidRFC1123(t, c.Name)
	}
}

// createStyleName mimics createContainerName's output: an uppercase prefix and a
// 64-char hex hash, ~128 chars total.
func createStyleName(service string) string {
	return "GITEA-ACTIONS-TASK-7-WORKFLOW-CI-JOB-build-" + service +
		"-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
}
