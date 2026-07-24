// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"errors"
	"strings"
	"testing"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/model"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// scriptContainer records the command a script was run with and answers with a fixed result.
type scriptContainer struct {
	fakeContainer
	cmd []string
	env map[string]string
	err error
}

func (c *scriptContainer) ToContainerPath(path string) string { return path }

func (c *scriptContainer) GetRunnerContext(context.Context) map[string]any {
	return map[string]any{"os": "Linux"}
}

func (c *scriptContainer) Exec(command []string, env map[string]string, _, _ string) common.Executor {
	return func(context.Context) error {
		c.cmd = command
		c.env = env
		return c.err
	}
}

func newScriptRunContext(t *testing.T, container *scriptContainer, config *Config) *RunContext {
	t.Helper()

	workflow, err := model.ReadWorkflow(strings.NewReader(`
name: test
on: push
jobs:
  job:
    runs-on: ubuntu-latest
    steps: []
`))
	require.NoError(t, err)

	// Env is left nil so that it is built from the config, as it is for a real job.
	rc := &RunContext{
		Config:       config,
		Run:          &model.Run{JobID: "job", Workflow: workflow},
		JobContainer: container,
	}
	rc.ExprEval = rc.NewExpressionEvaluator(context.Background())
	return rc
}

func TestRunJobScript(t *testing.T) {
	t.Run("runs the script with the job environment", func(t *testing.T) {
		container := &scriptContainer{}
		rc := newScriptRunContext(t, container, &Config{Env: map[string]string{"A_VAR": "value"}})

		require.NoError(t, rc.runJobScript("/hooks/setup.sh", "job pre", true)(context.Background()))

		assert.Equal(t, []string{"bash", "-e", "/hooks/setup.sh"}, container.cmd)
		assert.Equal(t, "value", container.env["A_VAR"])
		// The github environment is there too, so a hook can tell which job it runs for.
		assert.Equal(t, "job", container.env["GITHUB_JOB"])
	})

	t.Run("a failing started hook fails the job", func(t *testing.T) {
		container := &scriptContainer{err: errors.New("boom")}
		rc := newScriptRunContext(t, container, &Config{})

		err := rc.runJobScript("/hooks/setup.sh", "job pre", true)(context.Background())
		require.ErrorContains(t, err, "the job pre script")
		require.ErrorContains(t, err, "boom")
	})

	t.Run("a failing completed hook is only reported", func(t *testing.T) {
		container := &scriptContainer{err: errors.New("boom")}
		rc := newScriptRunContext(t, container, &Config{})
		logger, hook := test.NewNullLogger()
		ctx := common.WithLogger(context.Background(), logger.WithField("test", true))

		require.NoError(t, rc.runJobScript("/hooks/teardown.sh", "job post", false)(ctx))

		var reported bool
		for _, entry := range hook.AllEntries() {
			if assert.ObjectsAreEqual(entry.Level.String(), "error") {
				reported = true
			}
		}
		assert.True(t, reported, "the failure must reach the job log")
	})

	t.Run("is a no-op without a script", func(t *testing.T) {
		container := &scriptContainer{}
		rc := newScriptRunContext(t, container, &Config{})

		require.NoError(t, rc.runJobScript("", "job pre", true)(context.Background()))
		assert.Nil(t, container.cmd)
	})
}

func TestJobScriptConfiguration(t *testing.T) {
	t.Run("the setting wins over the environment variable", func(t *testing.T) {
		rc := &RunContext{Config: &Config{
			JobPreScript: "/from/config.sh",
			Env:          map[string]string{jobPreScriptEnv: "/from/env.sh"},
		}}
		assert.Equal(t, "/from/config.sh", rc.jobPreScript())
	})

	t.Run("the GitHub environment variables are honored", func(t *testing.T) {
		rc := &RunContext{Config: &Config{Env: map[string]string{
			jobPreScriptEnv:  "/from/env.sh",
			jobPostScriptEnv: "/from/env-done.sh",
		}}}
		assert.Equal(t, "/from/env.sh", rc.jobPreScript())
		assert.Equal(t, "/from/env-done.sh", rc.jobPostScript())
	})

	t.Run("nothing configured", func(t *testing.T) {
		rc := &RunContext{Config: &Config{}}
		assert.Empty(t, rc.jobPreScript())
		assert.Empty(t, rc.jobPostScript())
	})
}

func TestScriptCommand(t *testing.T) {
	assert.Equal(t, []string{"bash", "-e", "/hooks/setup.sh"}, scriptCommand("/hooks/setup.sh"))
	assert.Equal(t, []string{"pwsh", "-command", ". '/hooks/setup.ps1'"}, scriptCommand("/hooks/setup.ps1"))
	// An executable without a known extension speaks for itself through its shebang.
	assert.Equal(t, []string{"/hooks/setup"}, scriptCommand("/hooks/setup"))
}
