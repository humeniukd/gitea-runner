// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"fmt"
	"path"
	"strings"

	"gitea.com/gitea/runner/act/common"
)

// Job scripts are operator-provided scripts that run inside the job environment, before its
// first step and after its last one. They are the equivalent of GitHub's
// ACTIONS_RUNNER_HOOK_JOB_STARTED and ACTIONS_RUNNER_HOOK_JOB_COMPLETED, and the place to
// put per-job setup that no workflow should have to carry: registry logins, mirror
// configuration, credential files.
//
// Unlike runner.host_post_script, which runs on the host after the job is torn down,
// these run where the steps run and see the job's environment.
const (
	jobPreScriptEnv  = "ACTIONS_RUNNER_HOOK_JOB_STARTED"
	jobPostScriptEnv = "ACTIONS_RUNNER_HOOK_JOB_COMPLETED"
)

// jobPreScript returns the script to run before the job's first step: the runner
// setting, or the environment variable GitHub uses for it.
func (rc *RunContext) jobPreScript() string {
	if rc.Config.JobPreScript != "" {
		return rc.Config.JobPreScript
	}
	return rc.Config.Env[jobPreScriptEnv]
}

func (rc *RunContext) jobPostScript() string {
	if rc.Config.JobPostScript != "" {
		return rc.Config.JobPostScript
	}
	return rc.Config.Env[jobPostScriptEnv]
}

// runJobScript runs one script in the job environment. A failing pre script fails the
// job, because whatever it was setting up is missing; a failing post script is only
// reported, since the job's result is already decided by then.
func (rc *RunContext) runJobScript(scriptPath, name string, failJob bool) common.Executor {
	return func(ctx context.Context) error {
		if scriptPath == "" || rc.JobContainer == nil {
			return nil
		}

		logger := common.Logger(ctx)
		rawLogger := logger.WithField(rawOutputField, true)
		rawLogger.Infof("::group::Run %s script", name)
		defer rawLogger.Infof("::endgroup::")
		rawLogger.Infof("script: %s", scriptPath)

		env := mergeMaps(rc.GetEnv())
		rc.withGithubEnv(ctx, rc.getGithubContext(ctx), env)
		rc.ApplyExtraPath(ctx, &env)

		err := rc.execJobContainer(scriptCommand(scriptPath), env, "", "")(ctx)
		if err == nil {
			return nil
		}
		err = fmt.Errorf("the %s script %q failed: %w", name, scriptPath, err)
		if failJob {
			// Flip the job status the way a failing pre step does, so success()-default main
			// steps skip and the task is reported failed, then short-circuit the pre pipeline.
			reportStepError(ctx, rc, err)
			return err
		}
		logger.Errorf("##[error]%v", err)
		return nil
	}
}

// scriptCommand runs the script the way its extension says it wants to be run, and lets
// an executable without a known extension speak for itself through its shebang.
func scriptCommand(scriptPath string) []string {
	switch strings.ToLower(path.Ext(scriptPath)) {
	case ".sh":
		return []string{"bash", "-e", scriptPath}
	case ".ps1":
		return []string{"pwsh", "-command", ". '" + scriptPath + "'"}
	default:
		return []string{scriptPath}
	}
}
