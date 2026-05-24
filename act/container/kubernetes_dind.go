// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	"io"

	"gitea.com/gitea/runner/act/common"
)

// KubernetesDindDockerHost is the DOCKER_HOST the job container uses to reach
// the on-demand dind sidecar over the shared pod network.
const KubernetesDindDockerHost = kubernetesDindHost

// DindEnabled reports whether the on-demand docker-in-docker sidecar is enabled.
func (p *KubernetesPod) DindEnabled() bool {
	return p.config != nil && p.config.Dind.Enabled
}

// AddDindSidecar registers the docker-in-docker native sidecar into the pod,
// mounting the shared workspace at workdir so `docker run -v <workdir>:<workdir>`
// resolves against the same files the job container sees.
func (p *KubernetesPod) AddDindSidecar(workdir string) {
	p.builder.addDindSidecar(p.config, workdir)
}

// NewKubernetesDockerStepContainer returns a Container that translates a
// `uses: docker://` step into docker CLI commands exec'd inside the dind
// sidecar (which has docker CLI in PATH). Pull/Start become `docker pull` /
// `docker run`; file and path operations delegate to the job container.
func NewKubernetesDockerStepContainer(job *KubernetesExecutionEnvironment, input *NewContainerInput) Container {
	return newKubernetesDockerStepContainer(job, input, job.pod.execInDindContainer)
}

// newKubernetesDockerStepContainer is the testable constructor that accepts an
// explicit dindExec function so unit tests can stub the sidecar exec without a
// live cluster.
func newKubernetesDockerStepContainer(job ExecutionsEnvironment, input *NewContainerInput, dindExec func(context.Context, []string, io.Writer, io.Writer) error) Container {
	return &kubernetesDockerStepContainer{ExecutionsEnvironment: job, input: input, dindExec: dindExec}
}

// kubernetesDockerStepContainer embeds the job container's ExecutionsEnvironment
// so file/path operations delegate to it, and overrides the lifecycle methods to
// drive the dind daemon via docker CLI exec'd in the dind sidecar.
type kubernetesDockerStepContainer struct {
	ExecutionsEnvironment

	input   *NewContainerInput
	capAdd  []string
	capDrop []string
	// dindExec runs a shell-wrapped command in the dind sidecar, streaming
	// stdout/stderr to the supplied writers.
	dindExec func(context.Context, []string, io.Writer, io.Writer) error
}

// execInDind wraps command with buildExecCommand and runs it in the dind
// sidecar, streaming to the step's log writers.
func (c *kubernetesDockerStepContainer) execInDind(ctx context.Context, command []string) error {
	wrapped := buildExecCommand(command, nil, "", "")
	return c.dindExec(ctx, wrapped, c.input.Stdout, c.input.Stderr)
}

// Pull runs `docker pull` in the dind sidecar only when forcePull is set;
// otherwise the image is pulled lazily by `docker run`.
func (c *kubernetesDockerStepContainer) Pull(forcePull bool) common.Executor {
	return func(ctx context.Context) error {
		if !forcePull {
			return nil
		}
		return c.execInDind(ctx, []string{"docker", "pull", c.input.Image})
	}
}

// Create records the cap add/drop for the eventual `docker run`; the container
// is not created until Start.
func (c *kubernetesDockerStepContainer) Create(capAdd, capDrop []string) common.Executor {
	return func(_ context.Context) error {
		c.capAdd = capAdd
		c.capDrop = capDrop
		return nil
	}
}

// Start runs the step image via `docker run` against the dind daemon, exec'd
// inside the dind sidecar and streamed to the job's log writers.
func (c *kubernetesDockerStepContainer) Start(_ bool) common.Executor {
	return func(ctx context.Context) error {
		return c.execInDind(ctx, dockerRunArgs(c.input, c.capAdd, c.capDrop))
	}
}

// Remove is a no-op: `docker run --rm` removes the container on exit.
func (c *kubernetesDockerStepContainer) Remove() common.Executor {
	return func(_ context.Context) error { return nil }
}

func (c *kubernetesDockerStepContainer) Close() common.Executor {
	return func(_ context.Context) error { return nil }
}

// dockerRunArgs assembles the `docker run` command for a translated step: the
// shared workspace is bind-mounted at its own path and made the working dir, env
// is forwarded, optional privileged/cap flags and entrypoint are applied, then
// the image and its args. --rm cleans up on exit; --network host shares the
// dind daemon's network.
func dockerRunArgs(input *NewContainerInput, capAdd, capDrop []string) []string {
	args := []string{"docker", "run", "--rm", "--network", "host"}
	if input.WorkingDir != "" {
		args = append(args, "-v", input.WorkingDir+":"+input.WorkingDir, "-w", input.WorkingDir)
	}
	for _, env := range input.Env {
		args = append(args, "-e", env)
	}
	if input.Privileged {
		args = append(args, "--privileged")
	}
	for _, capability := range capAdd {
		args = append(args, "--cap-add", capability)
	}
	for _, capability := range capDrop {
		args = append(args, "--cap-drop", capability)
	}
	if len(input.Entrypoint) > 0 {
		args = append(args, "--entrypoint", input.Entrypoint[0])
	}
	args = append(args, input.Image)
	args = append(args, input.Cmd...)
	return args
}
