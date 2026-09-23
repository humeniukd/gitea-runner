// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package hooks

// Config configures the container-hooks mode. When Enabled is true, the runner
// stops calling the Docker API for job container operations and instead
// invokes Cmd (optionally through Entrypoint) for every hook event.
type Config struct {
	Enabled    bool   `yaml:"enabled"`    // Enabled turns container-hooks mode on. When false, the runner uses the built-in Docker backend and ignores Entrypoint and Cmd.
	Entrypoint string `yaml:"entrypoint"` // Entrypoint is an optional interpreter used to launch Cmd (for example, "/bin/sh" on Linux or "cmd.exe" on Windows). Empty means Cmd is executed directly.
	Cmd        string `yaml:"cmd"`        // Cmd is the path to the single hook entrypoint the runner invokes for every hook event.
}

// ExecuteHookInput is what the runner passes to the execute hook before
// running a step. It is written to <state_dir>/execute-input.json.
type ExecuteHookInput struct {
	// Kind selects what to run: "script" (a run: step), "command" (an action or
	// a direct exec), or "container-step" (uses: docker://...).
	Kind string `json:"kind"`
	// Populated when Kind is "script".
	Script string `json:"script,omitempty"`
	// Populated when Kind is "command".
	Command []string `json:"command,omitempty"`
	// Populated when Kind is "container-step".
	StepNumber int      `json:"stepNumber,omitempty"`
	Image      string   `json:"image,omitempty"`
	Entrypoint []string `json:"entrypoint,omitempty"`
	Args       []string `json:"args,omitempty"`
	// Common fields for all kinds.
	WorkingDir   string            `json:"workingDir,omitempty"`
	Env          map[string]string `json:"env,omitempty"`
	FileCommands []string          `json:"fileCommands,omitempty"`
}

type ExecuteHookResponse struct {
	ExitCode     int               `json:"exitCode"`
	FileCommands map[string]string `json:"fileCommands"`
}

type PrepareHookInput struct {
	Job    PrepareHookInputJob    `json:"job"`
	Runner PrepareHookInputRunner `json:"runner"`
}

type PrepareHookInputJob struct {
	Container ContainerJSON `json:"container"`
	Services  []ServiceJSON `json:"services,omitempty"`
	Steps     []StepJSON    `json:"steps,omitempty"`
	Mounts    MountsJSON    `json:"mounts"`
	Network   NetworkJSON   `json:"network"`
}

type PrepareHookInputRunner struct {
	RunnerJSON
}

// PrepareHookResponse is the shape of the response file written by the prepare
// hook. The runner reads it back from <state_dir>/prepare.json.
type PrepareHookResponse struct {
	ActPath      string            `json:"actPath,omitempty"`
	Environments map[string]string `json:"environments"`
	ID           string            `json:"id"`
}

// ContainerStepSpec describes a step that runs in its own container
// (uses: docker://...). Collected via prepareContainerStep() before
// execution starts and passed to the prepare hook inside prepare-input.json.
type ContainerStepSpec struct {
	StepNumber int
	StepID     string
	StepName   string
	Image      string
	Entrypoint []string
	Args       []string
	WorkingDir string
	Env        map[string]string
}

// EnvRef is a reference to an environment variable inside JSON. Instead of a
// literal value, the JSON carries {"$env": "VAR_NAME"}; the real value lives in
// the hook's process environment.
type EnvRef struct {
	Env string `json:"$env"`
}

// HookAction is one external action the runner downloaded for this job.
// The hook must stage LocalPath into <actPath>/actions/<Hash>/ before the
// job runs, so that node/exec inside the job find it at the expected path.
type HookAction struct {
	Hash      string `json:"hash"`
	LocalPath string `json:"localPath"`
}

// ContainerJSON describes a single container (the job container or a container
// step declared with uses: docker://...).
type ContainerJSON struct {
	Name           string            `json:"name"`
	Image          string            `json:"image"`
	Credentials    *CredentialsJSON  `json:"credentials,omitempty"`
	Entrypoint     []string          `json:"entrypoint,omitempty"`
	Cmd            []string          `json:"cmd,omitempty"`
	WorkingDir     string            `json:"workingDir,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	Privileged     bool              `json:"privileged,omitempty"`
	UsernsMode     string            `json:"usernsMode,omitempty"`
	Platform       string            `json:"platform,omitempty"`
	Options        string            `json:"options,omitempty"`
	NetworkAliases []string          `json:"networkAliases,omitempty"`
}

type CredentialsJSON struct {
	Username any `json:"username"`
	Password any `json:"password"`
}

type MountsJSON struct {
	Volumes map[string]string `json:"volumes,omitempty"`
	Binds   []string          `json:"binds,omitempty"`
}

type NetworkJSON struct {
	Name            string `json:"name"`
	CreateAndDelete bool   `json:"createAndDelete"`
}

// RunnerJSON carries runner-side paths and flags the hook may need to reproduce
// the execution environment (for example, where the workspace or action cache
// lives on the runner's filesystem).
type RunnerJSON struct {
	Workdir         string       `json:"workdir"`
	ToolCache       string       `json:"toolCache"`
	BindWorkdir     bool         `json:"bindWorkdir"`
	SharedToolCache bool         `json:"sharedToolCache"`
	EventJSON       string       `json:"eventJSON"`
	ActionsDir      string       `json:"actionsDir"`
	Actions         []HookAction `json:"actions,omitempty"`
}

// ServiceJSON describes a service container.
type ServiceJSON struct {
	Name           string            `json:"name"`
	Image          string            `json:"image"`
	Credentials    *CredentialsJSON  `json:"credentials,omitempty"`
	Cmd            []string          `json:"cmd,omitempty"`
	WorkingDir     string            `json:"workingDir,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	Ports          []string          `json:"ports,omitempty"`
	Options        string            `json:"options,omitempty"`
	NetworkAliases []string          `json:"networkAliases,omitempty"`
}

// StepJSON describes a step that runs in its own container (uses: docker://...).
type StepJSON struct {
	Number     int               `json:"number"`
	ID         string            `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Image      string            `json:"image"`
	Entrypoint []string          `json:"entrypoint,omitempty"`
	Args       []string          `json:"args,omitempty"`
	WorkingDir string            `json:"workingDir,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
}
