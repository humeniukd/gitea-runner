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
func (f *fakeExecEnv) Pull(_ bool) common.Executor  { return noopExecutor() }
func (f *fakeExecEnv) Start(_ bool) common.Executor { return noopExecutor() }
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
	step := NewKubernetesDockerStepContainer(job, &NewContainerInput{
		Image:      "alpine:3",
		WorkingDir: "/workspace",
		Cmd:        []string{"echo", "hi"},
	})

	ctx := context.Background()
	// forcePull=false: no docker pull (run pulls lazily).
	require.NoError(t, step.Pull(false)(ctx))
	assert.Empty(t, job.execCalls)

	// forcePull=true: docker pull.
	require.NoError(t, step.Pull(true)(ctx))
	require.Len(t, job.execCalls, 1)
	assert.Equal(t, []string{"docker", "pull", "alpine:3"}, job.execCalls[0])

	require.NoError(t, step.Create([]string{"SYS_PTRACE"}, nil)(ctx))
	require.NoError(t, step.Start(false)(ctx))
	require.Len(t, job.execCalls, 2)
	run := job.execCalls[1]
	assert.Equal(t, "docker", run[0])
	assert.Equal(t, "run", run[1])
	assert.Contains(t, run, "--cap-add")
	assert.Contains(t, run, "SYS_PTRACE")
	assert.Equal(t, "alpine:3", run[len(run)-3])
	assert.Equal(t, []string{"echo", "hi"}, run[len(run)-2:])

	// Remove/Close are no-ops (docker run --rm self-cleans).
	require.NoError(t, step.Remove()(ctx))
	require.NoError(t, step.Close()(ctx))
	assert.Len(t, job.execCalls, 2)
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
