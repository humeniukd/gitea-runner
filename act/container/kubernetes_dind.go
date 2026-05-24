// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"

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
// `uses: docker://` step into docker CLI commands exec'd inside the job
// container (which has DOCKER_HOST pointed at the dind sidecar). Pull/Start
// become `docker pull` / `docker run`; file and path operations delegate to the
// job container. It only needs the job environment, not the pod coordinator.
func NewKubernetesDockerStepContainer(job ExecutionsEnvironment, input *NewContainerInput) Container {
	return &kubernetesDockerStepContainer{ExecutionsEnvironment: job, input: input}
}

// kubernetesDockerStepContainer embeds the job container's ExecutionsEnvironment
// so file/path operations delegate to it, and overrides the lifecycle methods to
// drive the dind daemon via docker CLI exec'd in the job container.
type kubernetesDockerStepContainer struct {
	ExecutionsEnvironment

	input   *NewContainerInput
	capAdd  []string
	capDrop []string
}

// Pull runs `docker pull` only when forcePull is set; otherwise the image is
// pulled lazily by `docker run`.
func (c *kubernetesDockerStepContainer) Pull(forcePull bool) common.Executor {
	return func(ctx context.Context) error {
		if !forcePull {
			return nil
		}
		return c.Exec([]string{"docker", "pull", c.input.Image}, nil, "", "")(ctx)
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
// inside the job container and streamed to its log writers.
func (c *kubernetesDockerStepContainer) Start(_ bool) common.Executor {
	return func(ctx context.Context) error {
		return c.Exec(dockerRunArgs(c.input, c.capAdd, c.capDrop), nil, "", "")(ctx)
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
