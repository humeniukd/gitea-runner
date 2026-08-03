// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"gitea.com/gitea/runner/internal/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigCmdDefinesExpectedSubcommands(t *testing.T) {
	configFile := ""
	cmd := loadConfigCmd(&configFile)

	require.Equal(t, "config", cmd.Use)

	names := map[string]bool{}
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	for _, name := range []string{"generate", "get", "set", "add", "remove"} {
		assert.True(t, names[name], "expected subcommand %q to be registered", name)
	}
}

func TestConfigCmdGeneratePrintsTheExample(t *testing.T) {
	configFile := ""
	cmd := loadConfigCmd(&configFile)
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"generate"})

	require.NoError(t, cmd.Execute())
	assert.Equal(t, string(config.Example), out.String())
}

func TestConfigCmdRequiresConfigFile(t *testing.T) {
	t.Chdir(t.TempDir())

	configFile := ""
	cmd := loadConfigCmd(&configFile)
	cmd.SetArgs([]string{"set", "runner.capacity", "4"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--config")
}

func TestConfigCmdFallsBackToWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("runner:\n  capacity: 2\n"), 0o600))
	t.Chdir(dir)

	configFile := ""
	out, errOut := &bytes.Buffer{}, &bytes.Buffer{}
	cmd := loadConfigCmd(&configFile)
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.SetArgs([]string{"get", "runner.capacity"})

	require.NoError(t, cmd.Execute())
	assert.Equal(t, "2\n", out.String())
	assert.Contains(t, errOut.String(), "using config file")
}

func TestConfigCmdAddAndGet(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(configFile, []byte("runner:\n  labels:\n    - self-hosted\n"), 0o600))

	cmd := loadConfigCmd(&configFile)
	cmd.SetArgs([]string{"add", "runner.labels", "ubuntu:docker://node:22"})
	require.NoError(t, cmd.Execute())

	out := &bytes.Buffer{}
	cmd = loadConfigCmd(&configFile)
	cmd.SetOut(out)
	cmd.SetArgs([]string{"get", "runner.labels"})
	require.NoError(t, cmd.Execute())

	assert.Equal(t, "self-hosted\nubuntu:docker://node:22\n", out.String())
}
