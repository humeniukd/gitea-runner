// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"regexp"
	"strings"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/internal/pkg/action"

	"gitea.dev/actionslib/pkg/model"
)

type stepActionBuiltin struct {
	Step       *model.Step
	RunContext *RunContext
	action     action.Action
	env        map[string]string
}

func (sab *stepActionBuiltin) pre() common.Executor {
	sab.env = map[string]string{}
	return func(ctx context.Context) error {
		if !sab.action.HasPre() {
			return nil
		}
		return runStepExecutor(sab, stepStagePre, sab.run(sab.action.Pre))(ctx)
	}
}

func (sab *stepActionBuiltin) main() common.Executor {
	if sab.env == nil {
		sab.env = map[string]string{}
	}
	return runStepExecutor(sab, stepStageMain, sab.run(sab.action.Main))
}

func (sab *stepActionBuiltin) post() common.Executor {
	return func(ctx context.Context) error {
		if !sab.action.HasPost() {
			return nil
		}
		if sab.RunContext != nil {
			if r := sab.RunContext.getStepsContext()[sab.Step.ID]; r == nil || r.Conclusion == model.StepStatusSkipped {
				return nil
			}
		}
		return runStepExecutor(sab, stepStagePost, sab.run(sab.action.Post))(ctx)
	}
}

func (sab *stepActionBuiltin) run(fn func(*action.Context) error) common.Executor {
	return func(ctx context.Context) error {
		rc := sab.RunContext
		printRunActionHeader(ctx, sab.Step, sab.env, rc)
		rawLogger := common.Logger(ctx).WithField(rawOutputField, true)
		defer rawLogger.Infof("::endgroup::")

		inputs := sab.evaluatedInputs(ctx)
		populateEnvsFromBuiltinInputs(sab.env, inputs)
		return fn(&action.Context{
			Ctx:        ctx,
			RunContext: &builtinRunContextAdapter{ctx: ctx, rc: rc},
			Step:       sab.Step,
			Inputs:     inputs,
			Env:        sab.env,
			Logger:     common.Logger(ctx),
		})
	}
}

func (sab *stepActionBuiltin) evaluatedInputs(ctx context.Context) map[string]string {
	out := make(map[string]string, len(sab.Step.With))
	if len(sab.Step.With) == 0 {
		return out
	}
	eval := sab.RunContext.NewExpressionEvaluator(ctx)
	for k, v := range sab.Step.With {
		out[k] = eval.Interpolate(ctx, v)
	}
	return out
}

func populateEnvsFromBuiltinInputs(env, inputs map[string]string) {
	for k, v := range inputs {
		envKey := "INPUT_" + builtinInputEnvKey(k)
		if _, ok := env[envKey]; !ok {
			env[envKey] = v
		}
	}
}

var builtinInputEnvKeyRe = regexp.MustCompile(`[^A-Z0-9-]`)

func builtinInputEnvKey(name string) string {
	return builtinInputEnvKeyRe.ReplaceAllString(strings.ToUpper(name), "_")
}

func (sab *stepActionBuiltin) getRunContext() *RunContext { return sab.RunContext }
func (sab *stepActionBuiltin) getGithubContext(ctx context.Context) *model.GithubContext {
	return sab.RunContext.getGithubContext(ctx)
}
func (sab *stepActionBuiltin) getStepModel() *model.Step  { return sab.Step }
func (sab *stepActionBuiltin) getEnv() *map[string]string { return &sab.env }

func (sab *stepActionBuiltin) getIfExpression(_ context.Context, stage stepStage) string {
	switch stage {
	case stepStagePre, stepStagePost:
		return "always()"
	case stepStageMain:
		return sab.Step.If.Value
	}
	return ""
}

type builtinRunContextAdapter struct {
	ctx context.Context
	rc  *RunContext
}

func (a *builtinRunContextAdapter) JobContainer() container.ExecutionsEnvironment {
	return a.rc.JobContainer
}
func (a *builtinRunContextAdapter) Workdir() string { return a.rc.Config.Workdir }
func (a *builtinRunContextAdapter) ContainerWorkdir() string {
	if a.rc.JobContainer == nil {
		return a.rc.Config.Workdir
	}
	return a.rc.JobContainer.ToContainerPath(a.rc.Config.Workdir)
}
func (a *builtinRunContextAdapter) ActionCacheDir() string { return a.rc.ActionCacheDir() }
func (a *builtinRunContextAdapter) GithubContext() *model.GithubContext {
	return a.rc.getGithubContext(a.ctx)
}

func (a *builtinRunContextAdapter) Config() action.Config {
	return &builtinConfigAdapter{cfg: a.rc.Config}
}

type builtinConfigAdapter struct{ cfg *Config }

func (a *builtinConfigAdapter) Token() string                 { return a.cfg.Token }
func (a *builtinConfigAdapter) GitHubInstance() string        { return a.cfg.GitHubInstance }
func (a *builtinConfigAdapter) DefaultActionInstance() string { return a.cfg.DefaultActionInstance }
func (a *builtinConfigAdapter) InsecureSkipTLS() bool         { return a.cfg.InsecureSkipTLS }
func (a *builtinConfigAdapter) BindWorkdir() bool             { return a.cfg.BindWorkdir }
func (a *builtinConfigAdapter) NoSkipCheckout() bool          { return a.cfg.NoSkipCheckout }
func (a *builtinConfigAdapter) UseGitIgnore() bool            { return a.cfg.UseGitIgnore }
func (a *builtinConfigAdapter) ActionOfflineMode() bool       { return a.cfg.ActionOfflineMode }
