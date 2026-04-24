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
	Step             *model.Step
	RunContext       *RunContext
	cmd              []string
	cmdline          string
	env              map[string]string
	WorkingDirectory string
	// interpolatedScript caches the result of interpolating Step.Run during
	// setupShellCommand so the group-header printer doesn't re-evaluate it.
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
			// ApplyExtraPath must run before the header is printed so that any PATH
			// modifications from prior ::add-path:: commands are reflected in the
			// rendered env: block.
			rc.ApplyExtraPath(ctx, &sr.env)
			// Match actions/runner ScriptHandler.PrintActionDetails: emit the collapsible
			// block before execution so command output stays outside the fold.
			if err := sr.printRunScriptActionDetails(ctx); err != nil {
				return err
			}
			if he, ok := rc.JobContainer.(*container.HostEnvironment); ok && he != nil {
				return he.ExecWithCmdLine(sr.cmd, sr.cmdline, sr.env, "", sr.WorkingDirectory)(ctx)
			}
			return rc.JobContainer.Exec(sr.cmd, sr.env, "", sr.WorkingDirectory)(ctx)
		},
	))
}

// printRunScriptActionDetails mirrors actions/runner ScriptHandler.PrintActionDetails for
// script steps. The group title is the first line of the interpolated script, script lines
// are bright cyan, and env entries are emitted in YAML declaration order.
func (sr *stepRun) printRunScriptActionDetails(ctx context.Context) error {
	rawLogger := common.Logger(ctx).WithField(rawOutputField, true)
	scriptLineLogger := rawLogger.WithField(scriptLineCyanField, true)
	logger := common.Logger(ctx)
	script := sr.interpolatedScript

	rawLogger.Infof("::group::Run %s", sr.runScriptGroupTitle(script))

	normalized := strings.ReplaceAll(script, "\r\n", "\n")
	normalized = strings.TrimRight(normalized, "\n")
	if normalized != "" {
		for line := range strings.SplitSeq(normalized, "\n") {
			scriptLineLogger.Info(line)
		}
	}

	rawLogger.Infof("shell: %s", sr.cmdline)

	stepEnvKeys := stepDeclaredEnvKeysInOrder(sr.Step)
	if len(stepEnvKeys) > 0 {
		rawLogger.Infof("env:")
		caseInsensitiveEnv := sr.isEnvCaseInsensitive()
		for _, k := range stepEnvKeys {
			val, ok := lookupDeclaredEnvValue(k, sr.env, caseInsensitiveEnv)
			if !ok {
				logger.Debugf("step env key %q declared but not present in merged env", k)
			}
			rawLogger.Infof(" %s: %s", k, val)
		}
	}

	rawLogger.Infof("::endgroup::")
	return nil
}

// runScriptGroupTitle returns the group-header suffix following "Run ". It matches
// actions/runner ScriptHandler: the first non-whitespace line of the script. When the
// script is empty, fall back to the step name or id so the header stays useful.
func (sr *stepRun) runScriptGroupTitle(script string) string {
	if title := firstScriptLine(script); title != "" {
		return title
	}
	if sr.Step != nil {
		if sr.Step.Name != "" {
			return sr.Step.Name
		}
		if sr.Step.ID != "" {
			return sr.Step.ID
		}
	}
	return ""
}

func firstScriptLine(script string) string {
	first := strings.TrimLeft(strings.ReplaceAll(script, "\r\n", "\n"), " \t\r\n")
	if idx := strings.IndexAny(first, "\r\n"); idx >= 0 {
		first = first[:idx]
	}
	return first
}

// stepDeclaredEnvKeysInOrder walks the raw YAML Env mapping so keys are emitted in
// the order the workflow author wrote them. step.Environment() loses ordering because
// it decodes into a Go map.
func stepDeclaredEnvKeysInOrder(step *model.Step) []string {
	if step == nil || step.Env.Kind != yaml.MappingNode {
		return nil
	}
	keys := make([]string, 0)
	seen := make(map[string]struct{})
	collectMappingKeys(&step.Env, &keys, seen)
	return keys
}

func collectMappingKeys(node *yaml.Node, keys *[]string, seen map[string]struct{}) {
	if node == nil || node.Kind != yaml.MappingNode {
		return
	}
	content := node.Content
	for i := 0; i+1 < len(content); i += 2 {
		k := content[i]
		if k.Kind != yaml.ScalarNode {
			continue
		}
		if k.Tag == "!!merge" || k.Value == "<<" {
			continue
		}
		if _, dup := seen[k.Value]; dup {
			continue
		}
		seen[k.Value] = struct{}{}
		*keys = append(*keys, k.Value)
	}
}

func (sr *stepRun) isEnvCaseInsensitive() bool {
	rc := sr.getRunContext()
	return rc != nil && rc.JobContainer != nil && rc.JobContainer.IsEnvironmentCaseInsensitive()
}

func lookupDeclaredEnvValue(declKey string, env map[string]string, caseInsensitive bool) (string, bool) {
	if caseInsensitive {
		for k, v := range env {
			if strings.EqualFold(k, declKey) {
				return v, true
			}
		}
		return "", false
	}
	v, ok := env[declKey]
	return v, ok
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
