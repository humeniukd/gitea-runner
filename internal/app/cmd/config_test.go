// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

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
	for _, name := range []string{"get", "set", "add", "remove"} {
		assert.True(t, names[name], "expected subcommand %q to be registered", name)
	}
}

func TestConfigCmdRequiresConfigFile(t *testing.T) {
	configFile := ""
	cmd := loadConfigCmd(&configFile)
	cmd.SetArgs([]string{"set", "runner.capacity", "4"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--config")
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
