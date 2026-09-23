// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/act/container/hooks"
)

func (rc *RunContext) buildContainerPrepareHookInput(ctx context.Context) *hooks.PrepareHookInput {
	input := &hooks.PrepareHookInput{}

	// --- job.container ---
	image := rc.platformImage(ctx)
	ext := container.LinuxContainerEnvironmentExtensions{}
	envList := make([]string, 0)
	envList = append(envList, rc.runnerEnv(ctx)...)
	envList = append(envList, "LANG=C.UTF-8")

	envMap := make(map[string]string, len(envList))
	for _, kv := range envList {
		if k, v, ok := strings.Cut(kv, "="); ok {
			envMap[k] = v
		}
	}

	jobCreds := rc.buildJobContainerCredentials(ctx)

	input.Job.Container = hooks.ContainerJSON{
		Name:           rc.jobContainerName(),
		Image:          image,
		Credentials:    jobCreds,
		Entrypoint:     []string{"/bin/sleep", fmt.Sprint(rc.Config.ContainerMaxLifetime.Round(time.Second).Seconds())},
		WorkingDir:     ext.ToContainerPath(rc.Config.Workdir),
		Env:            envMap,
		Privileged:     rc.Config.Privileged,
		UsernsMode:     rc.Config.UsernsMode,
		Platform:       rc.Config.ContainerArchitecture,
		Options:        rc.options(ctx),
		NetworkAliases: []string{rc.Name},
	}

	// --- job.services ---
	for serviceID, spec := range rc.Run.Job().Services {
		serviceImage := rc.ExprEval.Interpolate(ctx, spec.Image)
		if serviceImage == "" {
			continue
		}

		interpolatedEnvs := make(map[string]string, len(spec.Env)+len(rc.Config.ProxyEnv))
		maps.Copy(interpolatedEnvs, rc.Config.ProxyEnv)
		for k, v := range spec.Env {
			interpolatedEnvs[k] = rc.ExprEval.Interpolate(ctx, v)
		}

		interpolatedCmd := make([]string, 0, len(spec.Cmd))
		for _, v := range spec.Cmd {
			interpolatedCmd = append(interpolatedCmd, rc.ExprEval.Interpolate(ctx, v))
		}

		interpolatedPorts := make([]string, 0, len(spec.Ports))
		for _, p := range spec.Ports {
			interpolatedPorts = append(interpolatedPorts, rc.ExprEval.Interpolate(ctx, p))
		}

		svcCreds := rc.buildServiceCredentials(ctx, serviceID, spec.Credentials)

		input.Job.Services = append(input.Job.Services, hooks.ServiceJSON{
			Name:           serviceID,
			Image:          serviceImage,
			Credentials:    svcCreds,
			Cmd:            interpolatedCmd,
			WorkingDir:     ext.ToContainerPath(rc.Config.Workdir),
			Env:            interpolatedEnvs,
			Ports:          interpolatedPorts,
			Options:        rc.ExprEval.Interpolate(ctx, spec.Options),
			NetworkAliases: []string{serviceID},
		})
	}

	// --- job.steps ---
	stepNumbers := make([]int, 0, len(rc.containerHooks.StepSpecs))
	for n := range rc.containerHooks.StepSpecs {
		stepNumbers = append(stepNumbers, n)
	}
	slices.Sort(stepNumbers)
	for _, n := range stepNumbers {
		spec := rc.containerHooks.StepSpecs[n]
		input.Job.Steps = append(input.Job.Steps, hooks.StepJSON{
			Number:     spec.StepNumber,
			ID:         spec.StepID,
			Name:       spec.StepName,
			Image:      spec.Image,
			Entrypoint: spec.Entrypoint,
			Args:       spec.Args,
			WorkingDir: spec.WorkingDir,
			Env:        spec.Env,
		})
	}

	// --- job.mounts ---
	binds, mounts := rc.GetBindsAndMounts()
	input.Job.Mounts = hooks.MountsJSON{
		Volumes: mounts,
		Binds:   binds,
	}

	// --- job.network ---
	networkName, createAndDelete := rc.networkNameForGitea()
	input.Job.Network = hooks.NetworkJSON{
		Name:            networkName,
		CreateAndDelete: createAndDelete,
	}

	// --- runner ---
	input.Runner = hooks.PrepareHookInputRunner{
		RunnerJSON: hooks.RunnerJSON{
			Workdir:         rc.Config.Workdir,
			ToolCache:       rc.toolCache(container.DefaultToolCache),
			BindWorkdir:     rc.Config.BindWorkdir,
			SharedToolCache: rc.Config.SharedToolCache,
			EventJSON:       rc.EventJSON,
			ActionsDir:      rc.ActionCacheDir(),
			Actions:         rc.containerHooks.Actions,
		},
	}

	return input
}

func (rc *RunContext) buildContainerExecuteHookCommandInput(
	cmd []string,
	env map[string]string,
	workdir string,
) *hooks.ExecuteHookInput {
	if workdir == "" {
		workdir = rc.Config.Workdir
	}

	// In the Docker backend, action files are copied into the job container at
	// /var/run/act/actions. In container-hooks mode they stay in the runner's
	// action cache, so rewrite the prefix here to keep the hook unaware of that
	// detail.
	fixedCmd := make([]string, len(cmd))
	for i, c := range cmd {
		fixedCmd[i] = strings.ReplaceAll(c, "/var/run/act/actions", rc.containerHooks.ActPath+"/actions")
	}

	return &hooks.ExecuteHookInput{
		Kind:         "command",
		Command:      fixedCmd,
		WorkingDir:   workdir,
		Env:          env,
		FileCommands: rc.containerHooks.StepFileCmds,
	}
}

func (rc *RunContext) buildContainerExecuteHookScriptInput(
	script string,
	env map[string]string,
	workdir string,
) *hooks.ExecuteHookInput {
	if workdir == "" {
		workdir = rc.Config.Workdir
	}
	return &hooks.ExecuteHookInput{
		Kind:         "script",
		Script:       script,
		WorkingDir:   workdir,
		Env:          env,
		FileCommands: rc.containerHooks.StepFileCmds,
	}
}

func (rc *RunContext) buildContainerExecuteHookStepInput(
	spec *hooks.ContainerStepSpec,
) *hooks.ExecuteHookInput {
	return &hooks.ExecuteHookInput{
		Kind:         "container-step",
		StepNumber:   spec.StepNumber,
		Image:        spec.Image,
		Entrypoint:   spec.Entrypoint,
		Args:         spec.Args,
		WorkingDir:   rc.Config.Workdir,
		Env:          spec.Env,
		FileCommands: rc.containerHooks.StepFileCmds,
	}
}

// buildContainerHookCredentials interpolates a username/password pair, stores
// the real values in rc.containerHooks.Secrets, and returns CredentialsJSON that
// refers to them via $env. envPrefix is the common prefix for the generated
// environment variable names, without the trailing USERNAME/PASSWORD.
func (rc *RunContext) buildContainerHookCredentials(
	ctx context.Context,
	creds map[string]string,
	envPrefix string,
) *hooks.CredentialsJSON {
	if len(creds) != 2 {
		return nil
	}

	username := rc.ExprEval.Interpolate(ctx, creds["username"])
	password := rc.ExprEval.Interpolate(ctx, creds["password"])
	if username == "" || password == "" {
		return nil
	}

	userKey := envPrefix + "USERNAME"
	passKey := envPrefix + "PASSWORD"

	if rc.containerHooks.Secrets == nil {
		rc.containerHooks.Secrets = make(map[string]string)
	}
	rc.containerHooks.Secrets[userKey] = username
	rc.containerHooks.Secrets[passKey] = password

	return &hooks.CredentialsJSON{
		Username: hooks.EnvRef{Env: userKey},
		Password: hooks.EnvRef{Env: passKey},
	}
}

func (rc *RunContext) buildJobContainerCredentials(ctx context.Context) *hooks.CredentialsJSON {
	c := rc.Run.Job().Container()
	if c == nil {
		return nil
	}
	return rc.buildContainerHookCredentials(ctx, c.Credentials,
		"RUNNER_CONTAINER_HOOK_SECRET_JOB_CONTAINER_CREDENTIALS_")
}

func (rc *RunContext) buildServiceCredentials(
	ctx context.Context,
	serviceID string,
	creds map[string]string,
) *hooks.CredentialsJSON {
	prefix := fmt.Sprintf("RUNNER_CONTAINER_HOOK_SECRET_JOB_SERVICES_%s_CREDENTIALS_",
		strings.ToUpper(serviceID))
	return rc.buildContainerHookCredentials(ctx, creds, prefix)
}

func (rc *RunContext) getContainerHookEnvironments(event string) []string {
	envMap := make(map[string]string)
	maps.Copy(envMap, rc.Env)
	maps.Copy(envMap, rc.containerHooks.Vars)
	maps.Copy(envMap, rc.containerHooks.Secrets)

	envMap["RUNNER_CONTAINER_HOOK_ID"] = rc.jobContainerID
	envMap["RUNNER_CONTAINER_HOOK_INPUT_FILE"] = hooks.GetContainerHookInputFile(event, rc.containerHooks.StateDir)
	envMap["RUNNER_CONTAINER_HOOK_JOB_NAME"] = rc.JobName
	envMap["RUNNER_CONTAINER_HOOK_LOG_LEVEL"] = rc.getContainerHookLogLevel()
	envMap["RUNNER_CONTAINER_HOOK_NAME"] = rc.jobContainerName()
	envMap["RUNNER_CONTAINER_HOOK_RESPONSE_FILE"] = hooks.GetContainerHookResponseFile(event, rc.containerHooks.StateDir)
	envMap["RUNNER_CONTAINER_HOOK_STATE_DIR"] = rc.containerHooks.StateDir
	envMap["RUNNER_CONTAINER_HOOK_TYPE"] = event
	envMap["RUNNER_CONTAINER_HOOK_WORKFLOW_NAME"] = rc.Run.Workflow.Name

	result := make([]string, 0, len(envMap))
	for k, v := range envMap {
		result = append(result, k+"="+v)
	}
	slices.Sort(result)
	return result
}

func (rc *RunContext) getContainerHookLogLevel() string {
	if rc.Config.JobLoggerLevel != nil {
		return rc.Config.JobLoggerLevel.String()
	}
	if rc.Config.RunnerDebug() {
		return "debug"
	}
	return "info"
}

// getContainerHookLogWriter wraps the hook's output in the same log format used
// by regular steps, so it shows up in the job UI alongside everything else.
func (rc *RunContext) getContainerHookLogWriter(ctx context.Context) io.Writer {
	logger := common.Logger(ctx)
	rawLogger := logger.WithField(rawOutputField, true)
	return common.NewLineWriter(rc.commandHandler(ctx), func(s string) bool {
		rawLogger.Infof("%s", s)
		return true
	})
}

func (rc *RunContext) readContainerPrepareHookResponse() error {
	path := hooks.GetContainerHookResponseFile("prepare", rc.containerHooks.StateDir)
	resp, err := hooks.ReadContainerHookResponse[hooks.PrepareHookResponse]("prepare", path)
	if err != nil {
		return fmt.Errorf("container hook prepare failed: %w", err)
	}
	if resp == nil {
		return fmt.Errorf("prepare hook did not return a response in %s", path)
	}
	if resp.ID == "" {
		return fmt.Errorf("prepare hook must return a non-empty \"id\" in %s", path)
	}
	rc.jobContainerID = resp.ID

	if resp.ActPath != "" {
		rc.containerHooks.ActPath = resp.ActPath
	} else {
		rc.containerHooks.ActPath = rc.JobContainer.GetActPath()
	}

	for k, v := range resp.Environments {
		if k == "" {
			return fmt.Errorf("prepare hook returned an empty env name in %s", path)
		}
		if rc.containerHooks.Vars == nil {
			rc.containerHooks.Vars = make(map[string]string)
		}
		rc.containerHooks.Vars[k] = v
	}
	return nil
}

func (rc *RunContext) readContainerExecuteHookResponse(hookErr error) error {
	path := hooks.GetContainerHookResponseFile("execute", rc.containerHooks.StateDir)
	resp, readErr := hooks.ReadContainerHookResponse[hooks.ExecuteHookResponse]("execute", path)
	if readErr != nil {
		if hookErr != nil {
			return hookErr
		}
		return readErr
	}
	if resp != nil && resp.FileCommands != nil {
		rc.containerHooks.StepFileCmdContents = resp.FileCommands
	}
	if resp != nil && resp.ExitCode != 0 {
		return container.ExitCodeError(resp.ExitCode)
	}
	return hookErr
}

func (rc *RunContext) recordContainerHookAction(hash, localPath string) {
	if !rc.containerHooks.Enabled {
		return
	}
	rc.containerHooks.actionsMutex.Lock()
	defer rc.containerHooks.actionsMutex.Unlock()
	rc.containerHooks.Actions = append(rc.containerHooks.Actions,
		hooks.HookAction{Hash: hash, LocalPath: localPath})
}

func (rc *RunContext) runContainerPrepareHook(ctx context.Context) error {
	if !rc.containerHooks.Enabled {
		return nil
	}

	if err := rc.writeContainerHookInput("prepare", rc.buildContainerPrepareHookInput(ctx)); err != nil {
		return err
	}

	env := rc.getContainerHookEnvironments("prepare")
	if err := hooks.RunContainerHook(ctx, rc.Config.ContainerHooks, "prepare", env, rc.getContainerHookLogWriter(ctx)); err != nil {
		return fmt.Errorf("container hook prepare failed: %w", err)
	}

	return rc.readContainerPrepareHookResponse()
}

func (rc *RunContext) runContainerExecuteHook(
	ctx context.Context,
	input *hooks.ExecuteHookInput,
) error {
	if !rc.containerHooks.Enabled {
		return nil
	}

	if err := rc.writeContainerHookInput("execute", input); err != nil {
		return err
	}

	// Clear previous response file
	_ = os.Remove(hooks.GetContainerHookResponseFile("execute", rc.containerHooks.StateDir))

	env := rc.getContainerHookEnvironments("execute")
	hookErr := hooks.RunContainerHook(ctx, rc.Config.ContainerHooks, "execute", env, rc.getContainerHookLogWriter(ctx))

	return rc.readContainerExecuteHookResponse(hookErr)
}

func (rc *RunContext) runContainerCleanupHook(ctx context.Context) error {
	if !rc.containerHooks.Enabled {
		return nil
	}

	env := rc.getContainerHookEnvironments("cleanup")
	err := hooks.RunContainerHook(ctx, rc.Config.ContainerHooks, "cleanup", env, rc.getContainerHookLogWriter(ctx))

	_ = os.RemoveAll(rc.containerHooks.StateDir)
	rc.containerHooks.StateDir = ""

	return err
}

func (rc *RunContext) runScriptStepInContainerHook(
	script string,
	env map[string]string,
	workdir string,
) common.Executor {
	return func(ctx context.Context) error {
		if !rc.Config.ContainerHooks.Enabled {
			return errors.New("container hooks directory is not set")
		}
		input := rc.buildContainerExecuteHookScriptInput(script, env, workdir)
		return rc.runContainerExecuteHook(ctx, input)
	}
}

func (rc *RunContext) writeContainerHookInput(event string, v any) error {
	path := hooks.GetContainerHookInputFile(event, rc.containerHooks.StateDir)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

func (rc *RunContext) parseFileCommandsFromContainerHook(ctx context.Context) error {
	if rc.containerHooks.StepFileCmdContents == nil {
		return nil
	}
	for _, p := range rc.containerHooks.StepFileCmds {
		content, ok := rc.containerHooks.StepFileCmdContents[p]
		if !ok {
			continue
		}
		switch filepath.Base(p) {
		case "envs.txt":
			if err := hooks.ParseEnvFileContent(ctx, content, rc.setEnv); err != nil {
				return err
			}
		case "statecmd.txt":
			if err := hooks.ParseEnvFileContent(ctx, content, rc.saveState); err != nil {
				return err
			}
		case "outputcmd.txt":
			if err := hooks.ParseEnvFileContent(ctx, content, rc.setOutput); err != nil {
				return err
			}
		case "pathcmd.txt":
			for line := range strings.SplitSeq(content, "\n") {
				if line != "" {
					rc.addPath(ctx, line)
				}
			}
		}
	}
	rc.containerHooks.StepFileCmds = nil
	rc.containerHooks.StepFileCmdContents = nil
	return nil
}
