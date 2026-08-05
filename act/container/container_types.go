// Copyright 2023 The Gitea Authors. All rights reserved.
// Copyright 2023 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	"errors"
	"fmt"
	"io"

	"gitea.com/gitea/runner/act/common"

	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/api/types/container"
)

// ExitCodeError reports a non-zero process exit code from a container command.
type ExitCodeError int

func (e ExitCodeError) Error() string {
	return fmt.Sprintf("Process completed with exit code %d.", int(e))
}

// NewContainerInput the input for the New function
type NewContainerInput struct {
	Image          string
	Username       string
	Password       string
	Entrypoint     []string
	Cmd            []string
	WorkingDir     string
	Env            []string
	Binds          []string
	Mounts         map[string]string
	Name           string
	Stdout         io.Writer
	Stderr         io.Writer
	NetworkMode    string
	Privileged     bool
	UsernsMode     string
	Platform       string
	Options        string
	NetworkAliases []string
	ExposedPorts   nat.PortSet
	PortBindings   nat.PortMap

	// Gitea specific
	AutoRemove   bool
	ValidVolumes []string
	AllocatePTY  bool // allocate a pseudo-TTY for the container's exec processes
}

// FileEntry is a file to copy to a container
type FileEntry struct {
	Name string
	Mode int64
	Body string
}

// Health states of a container, as reported by its healthcheck. They stay plain strings,
// so a caller of Info needs no docker SDK of its own.
const (
	HealthNone      = string(container.NoHealthcheck)
	HealthStarting  = string(container.Starting)
	HealthHealthy   = string(container.Healthy)
	HealthUnhealthy = string(container.Unhealthy)
)

// ErrContainerNotFound is what an AutoRemove container becomes moments after it exits.
// Its text is a fragment, missingContainerError composes it into the shared message.
var ErrContainerNotFound = errors.New("does not exist")

// Info describes a created container: what the `job` context reports about it, plus the
// state the runner waits on before it lets a job's steps talk to a service.
type Info struct {
	ID string
	// State is the docker container state: "created", "running", "exited", ...
	State    string
	ExitCode int
	// Health is one of the Health* constants, HealthNone when nothing is declared.
	Health string
	// HealthOutput is the last healthcheck probe's output, for reporting failures.
	HealthOutput string
	// Ports maps a container port ("5432/tcp") to the host port it is published on.
	Ports map[string]string
}

// Stopped reports a container that will not run again, including one the daemon is
// already dropping because it exited under AutoRemove.
func (i *Info) Stopped() bool {
	switch container.ContainerState(i.State) {
	case container.StateExited, container.StateDead, container.StateRemoving:
		return true
	default:
		return false
	}
}

// Container for managing docker run containers
type Container interface {
	Create(capAdd, capDrop []string) common.Executor
	ConnectToNetwork(name string) common.Executor
	Copy(destPath string, files ...*FileEntry) common.Executor
	CopyTarStream(ctx context.Context, destPath string, tarStream io.Reader) error
	CopyDir(destPath, srcPath string, useGitIgnore bool) common.Executor
	GetContainerArchive(ctx context.Context, srcPath string) (io.ReadCloser, error)
	Inspect(ctx context.Context) (*Info, error)
	Pull(forcePull bool) common.Executor
	Start(attach bool) common.Executor
	Exec(command []string, env map[string]string, user, workdir string) common.Executor
	UpdateFromEnv(srcPath string, env *map[string]string) common.Executor
	UpdateFromImageEnv(env *map[string]string) common.Executor
	Remove() common.Executor
	Close() common.Executor
	ReplaceLogWriter(io.Writer, io.Writer) (io.Writer, io.Writer)
}

// NewDockerBuildExecutorInput the input for the NewDockerBuildExecutor function
type NewDockerBuildExecutorInput struct {
	ContextDir   string
	Dockerfile   string
	BuildContext io.Reader
	ImageTag     string
	Platform     string
	BuildArgs    map[string]*string
}

// NewDockerNetworkCreateExecutorInput the input for the NewDockerNetworkCreateExecutor function
type NewDockerNetworkCreateExecutorInput struct {
	EnableIPv4 *bool
	EnableIPv6 *bool
	RunnerUUID string
}

// NewDockerPullExecutorInput the input for the NewDockerPullExecutor function
type NewDockerPullExecutorInput struct {
	Image     string
	ForcePull bool
	Platform  string
	Username  string
	Password  string
}
