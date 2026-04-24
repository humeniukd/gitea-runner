// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2022 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"fmt"
	"maps"
	"runtime"
	"strings"

	"gitea.com/gitea/act_runner/act/common"
	"gitea.com/gitea/act_runner/act/container"
	"gitea.com/gitea/act_runner/act/lookpath"
	"gitea.com/gitea/act_runner/act/model"

	"github.com/kballard/go-shellquote"
	yaml "go.yaml.in/yaml/v4"
)

type stepRun struct {
	Step               *model.Step
	RunContext         *RunContext
	cmd                []string
	cmdline            string
	env                map[string]string
	WorkingDirectory   string
	interpolatedScript string
}

func (sr *stepRun) pre() common.Executor {
	return func(ctx context.Context) error {
		return nil
	}
}

func (sr *stepRun) main() common.Executor {
	sr.env = map[string]string{}
	return runStepExecutor(sr, stepStageMain, common.NewPipelineExecutor(
		sr.setupShellCommandExecutor(),
		func(ctx context.Context) error {
			rc := sr.getRunContext()
			// Apply ::add-path:: effects before printing so PATH is accurate in the env: block.
			rc.ApplyExtraPath(ctx, &sr.env)
			sr.printRunScriptActionDetails(ctx)
			rawLogger := common.Logger(ctx).WithField(rawOutputField, true)
			defer rawLogger.Infof("::endgroup::")
			if he, ok := rc.JobContainer.(*container.HostEnvironment); ok && he != nil {
				return he.ExecWithCmdLine(sr.cmd, sr.cmdline, sr.env, "", sr.WorkingDirectory)(ctx)
			}
			return rc.JobContainer.Exec(sr.cmd, sr.env, "", sr.WorkingDirectory)(ctx)
		},
	))
}

// printRunScriptActionDetails mirrors actions/runner ScriptHandler.PrintActionDetails
// for script steps. The caller is responsible for emitting ::endgroup:: after step
// execution so command output is folded into the same collapsible block.
func (sr *stepRun) printRunScriptActionDetails(ctx context.Context) {
	rawLogger := common.Logger(ctx).WithField(rawOutputField, true)
	scriptLineLogger := rawLogger.WithField(scriptLineCyanField, true)

	normalized := strings.TrimRight(strings.ReplaceAll(sr.interpolatedScript, "\r\n", "\n"), "\n")

	rawLogger.Infof("::group::Run %s", sr.runScriptGroupTitle(normalized))

	if normalized != "" {
		for line := range strings.SplitSeq(normalized, "\n") {
			scriptLineLogger.Info(line)
		}
	}

	rawLogger.Infof("shell: %s", sr.cmdline)

	stepEnvKeys := stepDeclaredEnvKeysInOrder(sr.Step)
	if len(stepEnvKeys) == 0 {
		return
	}
	rawLogger.Infof("env:")
	envLookup := sr.env
	rc := sr.getRunContext()
	caseInsensitive := rc != nil && rc.JobContainer != nil && rc.JobContainer.IsEnvironmentCaseInsensitive()
	if caseInsensitive {
		envLookup = make(map[string]string, len(sr.env))
		for k, v := range sr.env {
			envLookup[strings.ToUpper(k)] = v
		}
	}
	for _, k := range stepEnvKeys {
		lookupKey := k
		if caseInsensitive {
			lookupKey = strings.ToUpper(k)
		}
		rawLogger.Infof("  %s: %s", k, envLookup[lookupKey])
	}
}

func (sr *stepRun) runScriptGroupTitle(normalizedScript string) string {
	trimmed := strings.TrimLeft(normalizedScript, " \t\n")
	if idx := strings.IndexByte(trimmed, '\n'); idx >= 0 {
		trimmed = trimmed[:idx]
	}
	if trimmed != "" {
		return trimmed
	}
	if sr.Step != nil {
		if sr.Step.Name != "" {
			return sr.Step.Name
		}
		return sr.Step.ID
	}
	return ""
}

// stepDeclaredEnvKeysInOrder walks the raw YAML Env mapping so keys are emitted in
// the order the workflow author wrote them; step.Environment() decodes into a Go map
// and loses ordering.
func stepDeclaredEnvKeysInOrder(step *model.Step) []string {
	if step == nil || step.Env.Kind != yaml.MappingNode {
		return nil
	}
	content := step.Env.Content
	keys := make([]string, 0, len(content)/2)
	for i := 0; i+1 < len(content); i += 2 {
		k := content[i]
		if k.Kind != yaml.ScalarNode || k.Tag == "!!merge" || k.Value == "<<" {
			continue
		}
		keys = append(keys, k.Value)
	}
	return keys
}

func (sr *stepRun) post() common.Executor {
	return func(ctx context.Context) error {
		return nil
	}
}

func (sr *stepRun) getRunContext() *RunContext {
	return sr.RunContext
}

func (sr *stepRun) getGithubContext(ctx context.Context) *model.GithubContext {
	return sr.getRunContext().getGithubContext(ctx)
}

func (sr *stepRun) getStepModel() *model.Step {
	return sr.Step
}

func (sr *stepRun) getEnv() *map[string]string {
	return &sr.env
}

func (sr *stepRun) getIfExpression(_ context.Context, _ stepStage) string {
	return sr.Step.If.Value
}

func (sr *stepRun) setupShellCommandExecutor() common.Executor {
	return func(ctx context.Context) error {
		scriptName, script, err := sr.setupShellCommand(ctx)
		if err != nil {
			return err
		}

		rc := sr.getRunContext()
		return rc.JobContainer.Copy(rc.JobContainer.GetActPath(), &container.FileEntry{
			Name: scriptName,
			Mode: 0o755,
			Body: script,
		})(ctx)
	}
}

func getScriptName(rc *RunContext, step *model.Step) string {
	scriptName := step.ID
	for rcs := rc; rcs.Parent != nil; rcs = rcs.Parent {
		scriptName = fmt.Sprintf("%s-composite-%s", rcs.Parent.CurrentStep, scriptName)
	}
	return "workflow/" + scriptName
}

// TODO: Currently we just ignore top level keys, BUT we should return proper error on them
// BUTx2 I leave this for when we rewrite act to use actionlint for workflow validation
// so we return proper errors before any execution or spawning containers
// it will error anyway with:
// OCI runtime exec failed: exec failed: container_linux.go:380: starting container process caused: exec: "${{": executable file not found in $PATH: unknown
func (sr *stepRun) setupShellCommand(ctx context.Context) (name, script string, err error) {
	logger := common.Logger(ctx)
	sr.setupShell(ctx)
	sr.setupWorkingDirectory(ctx)

	step := sr.Step

	script = sr.RunContext.NewStepExpressionEvaluator(ctx, sr).Interpolate(ctx, step.Run)
	sr.interpolatedScript = script

	scCmd := step.ShellCommand()

	name = getScriptName(sr.RunContext, step)

	// Reference: https://github.com/actions/runner/blob/8109c962f09d9acc473d92c595ff43afceddb347/src/Runner.Worker/Handlers/ScriptHandlerHelpers.cs#L47-L64
	// Reference: https://github.com/actions/runner/blob/8109c962f09d9acc473d92c595ff43afceddb347/src/Runner.Worker/Handlers/ScriptHandlerHelpers.cs#L19-L27
	runPrepend := ""
	runAppend := ""
	switch step.Shell {
	case "bash", "sh":
		name += ".sh"
	case "pwsh", "powershell":
		name += ".ps1"
		runPrepend = "$ErrorActionPreference = 'stop'"
		runAppend = "if ((Test-Path -LiteralPath variable:/LASTEXITCODE)) { exit $LASTEXITCODE }"
	case "cmd":
		name += ".cmd"
		runPrepend = "@echo off"
	case "python":
		name += ".py"
	}

	script = fmt.Sprintf("%s\n%s\n%s", runPrepend, script, runAppend)

	if !strings.Contains(script, "::add-mask::") && !sr.RunContext.Config.InsecureSecrets {
		logger.Debugf("Wrote command \n%s\n to '%s'", script, name)
	} else {
		logger.Debugf("Wrote add-mask command to '%s'", name)
	}

	rc := sr.getRunContext()
	scriptPath := fmt.Sprintf("%s/%s", rc.JobContainer.GetActPath(), name)
	sr.cmdline = strings.Replace(scCmd, `{0}`, scriptPath, 1)
	sr.cmd, err = shellquote.Split(sr.cmdline)

	return name, script, err
}

type localEnv struct {
	env map[string]string
}

func (l *localEnv) Getenv(name string) string {
	if runtime.GOOS == "windows" {
		for k, v := range l.env {
			if strings.EqualFold(name, k) {
				return v
			}
		}
		return ""
	}
	return l.env[name]
}

func (sr *stepRun) setupShell(ctx context.Context) {
	rc := sr.RunContext
	step := sr.Step

	if step.Shell == "" {
		step.Shell = rc.Run.Job().Defaults.Run.Shell
	}

	step.Shell = rc.NewExpressionEvaluator(ctx).Interpolate(ctx, step.Shell)

	if step.Shell == "" {
		step.Shell = rc.Run.Workflow.Defaults.Run.Shell
	}

	if step.Shell == "" {
		if _, ok := rc.JobContainer.(*container.HostEnvironment); ok {
			shellWithFallback := []string{"bash", "sh"}
			// Don't use bash on windows by default, if not using a docker container
			if runtime.GOOS == "windows" {
				shellWithFallback = []string{"pwsh", "powershell"}
			}
			step.Shell = shellWithFallback[0]
			lenv := &localEnv{env: map[string]string{}}
			maps.Copy(lenv.env, sr.env)
			sr.getRunContext().ApplyExtraPath(ctx, &lenv.env)
			_, err := lookpath.LookPath2(shellWithFallback[0], lenv)
			if err != nil {
				step.Shell = shellWithFallback[1]
			}
		} else if containerImage := rc.containerImage(ctx); containerImage != "" {
			// Currently only linux containers are supported, use sh by default like actions/runner
			step.Shell = "sh"
		}
	}
}

func (sr *stepRun) setupWorkingDirectory(ctx context.Context) {
	rc := sr.RunContext
	step := sr.Step
	var workingdirectory string

	if step.WorkingDirectory == "" {
		workingdirectory = rc.Run.Job().Defaults.Run.WorkingDirectory
	} else {
		workingdirectory = step.WorkingDirectory
	}

	// jobs can receive context values, so we interpolate
	workingdirectory = rc.NewExpressionEvaluator(ctx).Interpolate(ctx, workingdirectory)

	// but top level keys in workflow file like `defaults` or `env` can't
	if workingdirectory == "" {
		workingdirectory = rc.Run.Workflow.Defaults.Run.WorkingDirectory
	}
	sr.WorkingDirectory = workingdirectory
}
