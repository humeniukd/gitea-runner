// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build integration

// These tests exercise the Kubernetes backend against a real API server
// (envtest or kind). They are gated behind the `integration` build tag because
// they require a cluster runtime that is not available in plain unit-test runs;
// the E2E lane runs them against the live cluster. Provide a kubeconfig via the
// GITEA_RUNNER_TEST_KUBECONFIG env var (empty selects in-cluster config) and a
// pre-pullable busybox-style image via GITEA_RUNNER_TEST_IMAGE.
package container

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func integrationConfig(t *testing.T) *KubernetesConfig {
	t.Helper()
	ns := os.Getenv("GITEA_RUNNER_TEST_NAMESPACE")
	if ns == "" {
		ns = "gitea-e2e"
	}
	cfg := testKubernetesConfig()
	cfg.Namespace = ns
	cfg.Kubeconfig = os.Getenv("GITEA_RUNNER_TEST_KUBECONFIG")
	cfg.PodReadyTimeout = 2 * time.Minute
	return cfg
}

func integrationImage() string {
	if img := os.Getenv("GITEA_RUNNER_TEST_IMAGE"); img != "" {
		return img
	}
	return "busybox:stable"
}

// startIntegrationJob submits a real job pod and returns the job environment +
// a cleanup func.
func startIntegrationJob(t *testing.T, name string) (*KubernetesExecutionEnvironment, func()) {
	t.Helper()
	pod, err := NewKubernetesPod(KubernetesPodConfig{
		Config:      integrationConfig(t),
		PodName:     name,
		MaxLifetime: 10 * time.Minute,
	})
	require.NoError(t, err)

	job := pod.NewJobContainer(&NewContainerInput{
		Image:      integrationImage(),
		Entrypoint: []string{"/bin/sleep", "600"},
		WorkingDir: "/workspace",
		Name:       name,
	})
	ctx := context.Background()
	require.NoError(t, job.Pull(false)(ctx))
	require.NoError(t, job.Create(nil, nil)(ctx))
	require.NoError(t, job.Start(false)(ctx))

	return job, func() { _ = job.Remove()(context.Background()) }
}

func TestKubernetesExecRoundTrip(t *testing.T) {
	job, cleanup := startIntegrationJob(t, "gitea-it-exec")
	defer cleanup()

	ctx := context.Background()

	var stdout bytes.Buffer
	job.input.Stdout = &stdout
	job.input.Stderr = &stdout
	require.NoError(t, job.Exec([]string{"echo", "hello-k8s"}, nil, "", "")(ctx))
	assert.Contains(t, stdout.String(), "hello-k8s")

	// non-zero exit propagates as ExitCodeError
	err := job.Exec([]string{"sh", "-c", "exit 7"}, nil, "", "")(ctx)
	var exit ExitCodeError
	require.ErrorAs(t, err, &exit)
	assert.Equal(t, 7, int(exit))

	// `false` is exit code 1
	err = job.Exec([]string{"false"}, nil, "", "")(ctx)
	require.ErrorAs(t, err, &exit)
	assert.Equal(t, 1, int(exit))
}

func TestKubernetesExecHonorsEnvAndWorkdir(t *testing.T) {
	job, cleanup := startIntegrationJob(t, "gitea-it-env")
	defer cleanup()

	ctx := context.Background()
	var stdout bytes.Buffer
	job.input.Stdout = &stdout
	job.input.Stderr = &stdout
	require.NoError(t, job.Exec([]string{"sh", "-c", "echo $GREETING in $(pwd)"},
		map[string]string{"GREETING": "hi"}, "", "/tmp")(ctx))
	assert.Contains(t, stdout.String(), "hi in /tmp")
}
