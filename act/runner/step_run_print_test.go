// Copyright 2026 The Gitea Authors. All rights reserved.
// Copyright 2026 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"gitea.com/gitea/act_runner/act/common"
	"gitea.com/gitea/act_runner/act/model"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "go.yaml.in/yaml/v4"
)

func TestFirstScriptLine(t *testing.T) {
	tests := []struct {
		name   string
		script string
		want   string
	}{
		{"multiline", "echo one\necho two", "echo one"},
		{"crlf", "echo one\r\necho two", "echo one"},
		{"leading_ws", "\n  \techo hi", "echo hi"},
		{"empty", "", ""},
		{"single_line", "npm test", "npm test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, firstScriptLine(tt.script))
		})
	}
}

func TestRunScriptGroupTitleFallsBackToStepName(t *testing.T) {
	sr := &stepRun{Step: &model.Step{Name: "Build"}}
	assert.Equal(t, "make build", sr.runScriptGroupTitle("make build"))
	assert.Equal(t, "Build", sr.runScriptGroupTitle(""))

	sr = &stepRun{Step: &model.Step{ID: "s1"}}
	assert.Equal(t, "s1", sr.runScriptGroupTitle("\n  \n"))
}

func TestStepDeclaredEnvOrderPreservesYAML(t *testing.T) {
	raw := `id: s1
run: "echo 1"
env:
  GITHUB_TOKEN: tok
  PATH: /custom/bin
  MY_VAR: hello
`
	var step model.Step
	require.NoError(t, yaml.Unmarshal([]byte(raw), &step))

	keys := stepDeclaredEnvKeysInOrder(&step)
	assert.Equal(t, []string{"GITHUB_TOKEN", "PATH", "MY_VAR"}, keys)

	sr := &stepRun{
		Step:       &step,
		RunContext: &RunContext{},
		env: map[string]string{
			"GITHUB_TOKEN": "tok",
			"PATH":         "/custom/bin",
			"MY_VAR":       "hello",
			"TERM":         "xterm",
		},
	}
	v, ok := lookupDeclaredEnvValue("GITHUB_TOKEN", sr.env, sr.foldedEnvIfCaseInsensitive())
	assert.True(t, ok)
	assert.Equal(t, "tok", v)
	v, ok = lookupDeclaredEnvValue("PATH", sr.env, sr.foldedEnvIfCaseInsensitive())
	assert.True(t, ok)
	assert.Equal(t, "/custom/bin", v)
	_, ok = lookupDeclaredEnvValue("MISSING", sr.env, sr.foldedEnvIfCaseInsensitive())
	assert.False(t, ok)
}

func TestStepDeclaredEnvKeysInOrderEmpty(t *testing.T) {
	assert.Nil(t, stepDeclaredEnvKeysInOrder(nil))
	keys := stepDeclaredEnvKeysInOrder(&model.Step{})
	assert.Empty(t, keys)
}

func TestStepDeclaredEnvKeysExpandYAMLMergeKeys(t *testing.T) {
	doc := `
common: &common
  COMMON_A: a
  COMMON_B: b
step:
  env:
    LOCAL_BEFORE: before
    <<: *common
    COMMON_B: overridden
    LOCAL_AFTER: after
`
	var root struct {
		Step model.Step `yaml:"step"`
	}
	require.NoError(t, yaml.Unmarshal([]byte(doc), &root))

	keys := stepDeclaredEnvKeysInOrder(&root.Step)
	assert.Equal(t, []string{"LOCAL_BEFORE", "COMMON_A", "COMMON_B", "LOCAL_AFTER"}, keys)
	assert.NotContains(t, keys, "<<")

	resolved := root.Step.Environment()
	assert.Equal(t, "a", resolved["COMMON_A"])
	assert.Equal(t, "overridden", resolved["COMMON_B"])
	assert.Equal(t, "before", resolved["LOCAL_BEFORE"])
	assert.Equal(t, "after", resolved["LOCAL_AFTER"])
}

func TestApplyExtraPathBeforeHeader(t *testing.T) {
	raw := `id: s1
run: "echo 1"
env:
  PATH: /step/bin
`
	var step model.Step
	require.NoError(t, yaml.Unmarshal([]byte(raw), &step))

	rc := &RunContext{
		ExtraPath:    []string{"/extra/bin"},
		JobContainer: &containerMock{},
	}
	env := map[string]string{"PATH": "/step/bin"}
	rc.ApplyExtraPath(context.Background(), &env)
	assert.Equal(t, "/extra/bin:/step/bin", env["PATH"])

	sr := &stepRun{
		Step:               &step,
		RunContext:         rc,
		cmdline:            "bash ...",
		interpolatedScript: "echo 1\n",
		env:                env,
	}

	buf := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetOutput(buf)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&jobLogFormatter{color: cyan})
	entry := logger.WithFields(logrus.Fields{"job": "j1"})
	ctx := common.WithLogger(context.Background(), entry)

	require.NoError(t, sr.printRunScriptActionDetails(ctx))
	assert.Contains(t, buf.String(), "PATH: /extra/bin:/step/bin")
	assert.NotContains(t, buf.String(), "PATH: /step/bin\n")
}

func TestPrintRunScriptActionDetailsGolden(t *testing.T) {
	raw := `id: s1
name: Build
run: |
  echo one
  echo two
shell: pwsh
env:
  PATH_PREFIX: /custom/bin
  GREETING: hello
`
	var step model.Step
	require.NoError(t, yaml.Unmarshal([]byte(raw), &step))

	buf := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetOutput(buf)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&jobLogFormatter{color: cyan})
	entry := logger.WithFields(logrus.Fields{"job": "j1"})
	ctx := common.WithLogger(context.Background(), entry)

	sr := &stepRun{
		Step:               &step,
		RunContext:         &RunContext{},
		cmdline:            "pwsh -command . '/var/run/act/workflow/s1.ps1'",
		interpolatedScript: "echo one\necho two\n",
		env: map[string]string{
			"PATH_PREFIX": "/custom/bin",
			"GREETING":    "hello",
		},
	}

	require.NoError(t, sr.printRunScriptActionDetails(ctx))

	want := strings.Join([]string{
		"[j1]   | ::group::Run echo one",
		"[j1]   | echo one",
		"[j1]   | echo two",
		"[j1]   | shell: pwsh -command . '/var/run/act/workflow/s1.ps1'",
		"[j1]   | env:",
		"[j1]   |   PATH_PREFIX: /custom/bin",
		"[j1]   |   GREETING: hello",
		"[j1]   | ::endgroup::",
		"",
	}, "\n")
	assert.Equal(t, want, buf.String())
}

func TestPrintColoredScriptLineCyan(t *testing.T) {
	f := &jobLogFormatter{color: cyan}
	entry := &logrus.Entry{
		Level:   logrus.InfoLevel,
		Message: "echo one",
		Data: logrus.Fields{
			"job":               "j1",
			rawOutputField:      true,
			scriptLineCyanField: true,
		},
	}
	buf := &bytes.Buffer{}
	f.printColored(buf, entry)
	assert.Equal(t, "\x1b[36m|\x1b[0m \x1b[36;1mecho one\x1b[0m", buf.String())
}
