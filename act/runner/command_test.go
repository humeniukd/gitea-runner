// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/model"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetEnv(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	rc := new(RunContext)
	handler := rc.commandHandler(ctx)

	handler("::set-env name=x::valz\n")
	a.Equal("valz", rc.Env["x"])
}

func TestStopCommandsKeepsSuppressedLinesInLog(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	rc := new(RunContext)
	handler := rc.commandHandler(ctx)

	// Stop command processing until the matching end token is seen.
	a.True(handler("::stop-commands::my-end-token\n"))

	// A command-shaped line while stopped must not be executed (env unchanged),
	// but must still return true so it reaches the raw_output log handler and is
	// not dropped from the step log.
	a.True(handler("::set-env name=x::valz\n"))
	a.NotContains(rc.Env, "x")

	// The matching end token resumes command processing.
	a.True(handler("::my-end-token::\n"))

	// Commands are processed again after resuming.
	a.True(handler("::set-env name=y::valy\n"))
	a.Equal("valy", rc.Env["y"])
}

func TestSetOutput(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	rc := new(RunContext)
	rc.StepResults = make(map[string]*model.StepResult)
	handler := rc.commandHandler(ctx)

	rc.CurrentStep = "my-step"
	rc.StepResults[rc.CurrentStep] = &model.StepResult{
		Outputs: make(map[string]string),
	}
	handler("::set-output name=x::valz\n")
	a.Equal("valz", rc.StepResults["my-step"].Outputs["x"])

	handler("::set-output name=x::percent2%25\n")
	a.Equal("percent2%", rc.StepResults["my-step"].Outputs["x"])

	handler("::set-output name=x::percent2%25%0Atest\n")
	a.Equal("percent2%\ntest", rc.StepResults["my-step"].Outputs["x"])

	handler("::set-output name=x::percent2%25%0Atest another3%25test\n")
	a.Equal("percent2%\ntest another3%test", rc.StepResults["my-step"].Outputs["x"])

	handler("::set-output name=x%3A::percent2%25%0Atest\n")
	a.Equal("percent2%\ntest", rc.StepResults["my-step"].Outputs["x:"])

	handler("::set-output name=x%3A%2C%0A%25%0D%3A::percent2%25%0Atest\n")
	a.Equal("percent2%\ntest", rc.StepResults["my-step"].Outputs["x:,\n%\r:"])
}

func TestAddpath(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	rc := new(RunContext)
	handler := rc.commandHandler(ctx)

	handler("::add-path::/zoo\n")
	a.Equal("/zoo", rc.ExtraPath[0])

	handler("::add-path::/boo\n")
	a.Equal("/boo", rc.ExtraPath[0])
}

func TestStopCommands(t *testing.T) {
	logger, hook := test.NewNullLogger()

	a := assert.New(t)
	ctx := common.WithLogger(context.Background(), logger)
	rc := new(RunContext)
	handler := rc.commandHandler(ctx)

	handler("::set-env name=x::valz\n")
	a.Equal("valz", rc.Env["x"])
	handler("::stop-commands::my-end-token\n")
	handler("::set-env name=x::abcd\n")
	a.Equal("valz", rc.Env["x"])
	handler("::my-end-token::\n")
	handler("::set-env name=x::abcd\n")
	a.Equal("abcd", rc.Env["x"])

	messages := make([]string, 0)
	for _, entry := range hook.AllEntries() {
		messages = append(messages, entry.Message)
	}

	a.Contains(messages, "::set-env name=x::abcd\n")
}

func TestAddpathADO(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	rc := new(RunContext)
	handler := rc.commandHandler(ctx)

	handler("##[add-path]/zoo\n")
	a.Equal("/zoo", rc.ExtraPath[0])

	handler("##[add-path]/boo\n")
	a.Equal("/boo", rc.ExtraPath[0])
}

func TestAddmask(t *testing.T) {
	logger, hook := test.NewNullLogger()

	a := assert.New(t)
	ctx := context.Background()
	loggerCtx := common.WithLogger(ctx, logger)

	rc := new(RunContext)
	handler := rc.commandHandler(loggerCtx)
	handler("::add-mask::my-secret-value\n")

	a.Equal("***", hook.LastEntry().Message)
	a.NotEqual("*my-secret-value", hook.LastEntry().Message)
}

// based on https://stackoverflow.com/a/10476304
func captureOutput(t *testing.T, f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			a := assert.New(t)
			a.Fail("io.Copy failed")
		}
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	out := <-outC

	return out
}

func TestAddmaskUsemask(t *testing.T) {
	rc := new(RunContext)
	rc.StepResults = make(map[string]*model.StepResult)
	rc.CurrentStep = "my-step"
	rc.StepResults[rc.CurrentStep] = &model.StepResult{
		Outputs: make(map[string]string),
	}

	a := assert.New(t)

	config := &Config{
		Secrets:         map[string]string{},
		InsecureSecrets: false,
	}

	re := captureOutput(t, func() {
		ctx := context.Background()
		ctx = WithJobLogger(ctx, "0", "testjob", config, &rc.Masks, map[string]any{})

		handler := rc.commandHandler(ctx)
		handler("::add-mask::secret\n")
		handler("::set-output:: token=secret\n")
	})

	a.Equal("[testjob] ***\n[testjob] ::set-output:: = token=***\n", re)
}

func TestSaveState(t *testing.T) {
	rc := &RunContext{
		CurrentStep: "step",
		StepResults: map[string]*model.StepResult{},
	}

	ctx := context.Background()

	handler := rc.commandHandler(ctx)
	handler("::save-state name=state-name::state-value\n")

	assert.Equal(t, "state-value", rc.IntraActionState["step"]["state-name"])
}

func TestUnsecureCommandsAreRefusedWhenDisabled(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	disallow := false
	rc := &RunContext{Config: &Config{AllowUnsecureCommands: &disallow}}
	handler := rc.commandHandler(ctx)

	handler("::set-env name=x::valz\n")
	handler("::add-path::/opt/bin\n")

	a.Empty(rc.Env["x"])
	a.Empty(rc.ExtraPath)

	// The refusal fails the step that produced it, once.
	err := rc.takeUnsecureCommandError()
	require.ErrorContains(t, err, "set-env")
	a.NoError(rc.takeUnsecureCommandError())
}

func TestUnsecureCommandsOptIn(t *testing.T) {
	disallow := false

	t.Run("through the step environment", func(t *testing.T) {
		a := assert.New(t)
		rc := &RunContext{Config: &Config{AllowUnsecureCommands: &disallow}}
		env := map[string]string{"ACTIONS_ALLOW_UNSECURE_COMMANDS": "true"}
		rc.setCurrentStepEnv(&env)

		rc.commandHandler(context.Background())("::set-env name=x::valz\n")

		a.Equal("valz", rc.Env["x"])
		a.NoError(rc.takeUnsecureCommandError())
	})

	t.Run("through the job environment", func(t *testing.T) {
		a := assert.New(t)
		rc := &RunContext{
			Config: &Config{AllowUnsecureCommands: &disallow},
			Env:    map[string]string{"ACTIONS_ALLOW_UNSECURE_COMMANDS": "1"},
		}

		rc.commandHandler(context.Background())("::add-path::/opt/bin\n")

		a.Equal([]string{"/opt/bin"}, rc.ExtraPath)
		a.NoError(rc.takeUnsecureCommandError())
	})
}

// Until the default flips, the commands still run, but the job log says so exactly once.
func TestUnsecureCommandsWarnOnceWhenAllowed(t *testing.T) {
	a := assert.New(t)
	logger, hook := test.NewNullLogger()
	ctx := common.WithLogger(context.Background(), logger.WithField("test", true))

	rc := &RunContext{Config: &Config{}}
	handler := rc.commandHandler(ctx)
	handler("::set-env name=x::valz\n")
	handler("::add-path::/opt/bin\n")

	a.Equal("valz", rc.Env["x"])

	warnings := 0
	for _, entry := range hook.AllEntries() {
		if strings.Contains(entry.Message, "##[warning]") {
			warnings++
		}
	}
	a.Equal(1, warnings)
}
