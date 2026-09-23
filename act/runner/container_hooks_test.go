// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/container/hooks"

	"gitea.dev/actionslib/pkg/model"
	log "github.com/sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func hooksTestConfig(t *testing.T) hooks.Config {
	t.Helper()
	cmd, err := filepath.Abs("testdata/container-hooks/hooks/main.sh")
	assert.NoError(t, err)
	return hooks.Config{
		Enabled:    true,
		Entrypoint: "/bin/sh",
		Cmd:        cmd,
	}
}

type containerHooksTestInfo struct {
	workflowPath string
	errorMessage string
}

type hooksJobLoggerFactory struct {
	Output bytes.Buffer
}

func (f *hooksJobLoggerFactory) WithJobLogger() *log.Logger {
	logger := log.New()
	logger.SetOutput(io.MultiWriter(&f.Output, os.Stdout))
	logger.SetLevel(log.DebugLevel)
	return logger
}

func (j *containerHooksTestInfo) run(t *testing.T, loggerFactory *hooksJobLoggerFactory) {
	t.Helper()

	log.SetLevel(logLevel)

	ctx := WithJobLoggerFactory(
		common.WithLogger(context.Background(), loggerFactory.WithJobLogger()),
		loggerFactory,
	)

	workdir, err := filepath.Abs("testdata/container-hooks")
	require.NoError(t, err)

	runnerConfig := &Config{
		Workdir:              workdir,
		BindWorkdir:          false,
		EventName:            "push",
		Platforms:            map[string]string{"docker": "node:24-bookworm-slim"},
		ContainerNamePrefix:  "container-hooks-test",
		AutoRemove:           true,
		MaxParallel:          1,
		ContainerMaxLifetime: time.Hour,
		ContainerHooks:       hooksTestConfig(t),
	}

	runner, err := New(runnerConfig)
	require.NoError(t, err, j.workflowPath)

	planner, err := model.NewWorkflowPlanner(filepath.Join(workdir, j.workflowPath), true)
	require.NoError(t, err, j.workflowPath)

	plan, err := planner.PlanEvent("push")
	require.NoError(t, err, j.workflowPath)
	assert.NotNil(t, plan)

	err = runner.NewPlanExecutor(plan)(ctx)
	if j.errorMessage == "" {
		assert.NoError(t, err, j.workflowPath)
	} else {
		require.Error(t, err, j.workflowPath)
		assert.Contains(t, err.Error(), j.errorMessage, j.workflowPath)
	}
}

func TestContainerHooksRun(t *testing.T) {
	// No requireDocker: in container-hooks mode the runner must not touch the
	// Docker API at all.
	t.Parallel()

	tables := []containerHooksTestInfo{
		{"push.yml", ""},
	}

	for _, table := range tables {
		t.Run(table.workflowPath, func(t *testing.T) {
			t.Parallel()
			factory := &hooksJobLoggerFactory{}
			table.run(t, factory)
			assert.Contains(t, factory.Output.String(), "prepare hook: start")
			assert.Contains(t, factory.Output.String(), "execute hook: start")
			assert.Contains(t, factory.Output.String(), "cleanup hook: start")
		})
	}
}
