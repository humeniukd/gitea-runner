// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"bytes"
	"context"
	"io"
	"testing"

	"gitea.com/gitea/runner/act/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

// fakeExecEnv is a minimal ExecutionsEnvironment that records Exec command lines
// and no-ops everything else, so the docker-step translation can be unit-tested
// without a cluster. Path helpers come from the embedded Linux extensions.
type fakeExecEnv struct {
	LinuxContainerEnvironmentExtensions
	execCalls [][]string
}

func noopExecutor() common.Executor { return func(_ context.Context) error { return nil } }

func (f *fakeExecEnv) Exec(command []string, _ map[string]string, _, _ string) common.Executor {
	return func(_ context.Context) error {
		f.execCalls = append(f.execCalls, command)
		return nil
	}
}

func (f *fakeExecEnv) Create(_, _ []string) common.Executor           { return noopExecutor() }
func (f *fakeExecEnv) ConnectToNetwork(_ string) common.Executor      { return noopExecutor() }
func (f *fakeExecEnv) Copy(_ string, _ ...*FileEntry) common.Executor { return noopExecutor() }
func (f *fakeExecEnv) CopyTarStream(_ context.Context, _ string, _ io.Reader) error {
	return nil
}
func (f *fakeExecEnv) CopyDir(_, _ string, _ bool) common.Executor { return noopExecutor() }
func (f *fakeExecEnv) GetContainerArchive(_ context.Context, _ string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(nil)), nil
}

func (f *fakeExecEnv) Inspect(_ context.Context) (*Info, error) {
	return &Info{Health: HealthNone, Ports: map[string]string{}}, nil
}
func (f *fakeExecEnv) DumpLogs(_ context.Context) error { return nil }
func (f *fakeExecEnv) Pull(_ bool) common.Executor      { return noopExecutor() }
func (f *fakeExecEnv) Start(_ bool) common.Executor     { return noopExecutor() }
func (f *fakeExecEnv) UpdateFromEnv(_ string, _ *map[string]string) common.Executor {
	return noopExecutor()
}
func (f *fakeExecEnv) UpdateFromImageEnv(_ *map[string]string) common.Executor { return noopExecutor() }
func (f *fakeExecEnv) Remove() common.Executor                                 { return noopExecutor() }
func (f *fakeExecEnv) Close() common.Executor                                  { return noopExecutor() }
func (f *fakeExecEnv) ReplaceLogWriter(stdout, stderr io.Writer) (io.Writer, io.Writer) {
	return stdout, stderr
}

var _ ExecutionsEnvironment = (*fakeExecEnv)(nil)

func TestDockerRunArgs(t *testing.T) {
	args := dockerRunArgs(&NewContainerInput{
		Image:      "alpine:3",
		WorkingDir: "/workspace/owner/repo",
		Env:        []string{"FOO=bar"},
		Entrypoint: []string{"/entry.sh"},
		Cmd:        []string{"echo", "hi"},
		Privileged: true,
	}, []string{"SYS_PTRACE"}, []string{"MKNOD"})

	assert.Equal(t, []string{
		"docker", "run", "--rm", "--network", "host",
		"-v", "/workspace/owner/repo:/workspace/owner/repo", "-w", "/workspace/owner/repo",
		"-e", "FOO=bar",
		"--privileged",
		"--cap-add", "SYS_PTRACE",
		"--cap-drop", "MKNOD",
		"--entrypoint", "/entry.sh",
		"alpine:3", "echo", "hi",
	}, args)
}

func TestKubernetesDockerStepContainerTranslatesToDockerCLI(t *testing.T) {
	job := &fakeExecEnv{}
	var dindCalls [][]string
	fakeDindExec := func(_ context.Context, command []string, _, _ io.Writer) error {
		dindCalls = append(dindCalls, command)
		return nil
	}
	step := newKubernetesDockerStepContainer(job, &NewContainerInput{
		Image:      "alpine:3",
		WorkingDir: "/workspace",
		Cmd:        []string{"echo", "hi"},
	}, fakeDindExec)

	ctx := context.Background()
	// forcePull=false: no docker pull (run pulls lazily).
	require.NoError(t, step.Pull(false)(ctx))
	assert.Empty(t, dindCalls)

	// forcePull=true: docker pull exec'd in dind sidecar.
	require.NoError(t, step.Pull(true)(ctx))
	require.Len(t, dindCalls, 1)
	// buildExecCommand wraps with sh -c so dind sidecar gets: ["sh", "-c", "exec docker pull alpine:3"]
	assert.Equal(t, []string{"sh", "-c", "exec docker pull alpine:3"}, dindCalls[0])

	require.NoError(t, step.Create([]string{"SYS_PTRACE"}, nil)(ctx))
	require.NoError(t, step.Start(false)(ctx))
	require.Len(t, dindCalls, 2)
	// Start sends the full docker run command wrapped in sh -c to the dind sidecar.
	run := dindCalls[1]
	assert.Equal(t, []string{"sh", "-c"}, run[:2])
	assert.Contains(t, run[2], "docker run")
	assert.Contains(t, run[2], "--cap-add SYS_PTRACE")
	assert.Contains(t, run[2], "alpine:3")
	assert.Contains(t, run[2], "echo hi")

	// job.Exec is NOT called: docker commands go to dind, not job container.
	assert.Empty(t, job.execCalls)

	// Remove/Close are no-ops (docker run --rm self-cleans).
	require.NoError(t, step.Remove()(ctx))
	require.NoError(t, step.Close()(ctx))
	assert.Len(t, dindCalls, 2)
}

func TestAddDindSidecarRootless(t *testing.T) {
	cfg := testKubernetesConfig()
	cfg.Dind = KubernetesDindConfig{Enabled: true, Rootless: true}
	cfg.Resources.Dind = KubernetesResources{Requests: KubernetesResourceList{CPU: "500m", Memory: "1Gi"}}

	b := newPodBuilder("gitea-job-dind", cfg.Namespace, cfg, 0, nil)
	b.addDindSidecar(cfg, "/workspace")
	pod := b.build()

	require.Len(t, pod.Spec.InitContainers, 1)
	dind := pod.Spec.InitContainers[0]
	assert.Equal(t, kubernetesDindContainerName, dind.Name)
	assert.Equal(t, kubernetesDindImageRootless, dind.Image)
	require.NotNil(t, dind.RestartPolicy)
	assert.Equal(t, corev1.ContainerRestartPolicyAlways, *dind.RestartPolicy)
	assert.Nil(t, dind.SecurityContext, "rootless dind must not be privileged")
	require.Len(t, dind.Ports, 1)
	assert.Equal(t, int32(kubernetesDindPort), dind.Ports[0].ContainerPort)
	require.NotNil(t, dind.StartupProbe)
	assert.Equal(t, kubernetesDindPort, dind.StartupProbe.TCPSocket.Port.IntValue())
	// TLS disabled so the daemon serves plain TCP.
	assert.Contains(t, dind.Env, corev1.EnvVar{Name: "DOCKER_TLS_CERTDIR", Value: ""})
	// workspace shared so docker run -v resolves.
	require.Len(t, dind.VolumeMounts, 1)
	assert.Equal(t, "/workspace", dind.VolumeMounts[0].MountPath)
}

func TestAddDindSidecarPrivileged(t *testing.T) {
	cfg := testKubernetesConfig()
	cfg.Dind = KubernetesDindConfig{Enabled: true, Rootless: false, Image: "docker:27-dind"}

	b := newPodBuilder("gitea-job-dind2", cfg.Namespace, cfg, 0, nil)
	b.addDindSidecar(cfg, "/workspace")
	pod := b.build()

	dind := pod.Spec.InitContainers[0]
	assert.Equal(t, "docker:27-dind", dind.Image)
	require.NotNil(t, dind.SecurityContext)
	require.NotNil(t, dind.SecurityContext.Privileged)
	assert.True(t, *dind.SecurityContext.Privileged, "non-rootless dind must be privileged")
}
