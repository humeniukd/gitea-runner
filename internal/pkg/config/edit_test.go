// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const editFixture = `# A leading comment.
log:
  # The logging level.
  level: info

runner:
  capacity: 1
  envs:
    EXISTING: value
  timeout: 3h
  labels:
    - ubuntu-latest:docker://node:20
    - self-hosted
`

func writeEditFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(editFixture), 0o600))
	return path
}

func TestEditValues(t *testing.T) {
	tests := []struct {
		name   string
		edit   func(file string) error
		assert func(t *testing.T, cfg *Config, content string)
	}{
		{
			name: "set scalar",
			edit: func(file string) error { return SetValue(file, "runner.capacity", "4") },
			assert: func(t *testing.T, cfg *Config, _ string) {
				assert.Equal(t, 4, cfg.Runner.Capacity)
			},
		},
		{
			name: "set duration",
			edit: func(file string) error { return SetValue(file, "runner.timeout", "90m") },
			assert: func(t *testing.T, cfg *Config, content string) {
				assert.Equal(t, 90*time.Minute, cfg.Runner.Timeout)
				assert.Contains(t, content, "timeout: 1h30m0s")
			},
		},
		{
			name: "set bool in a missing section",
			edit: func(file string) error { return SetValue(file, "metrics.enabled", "true") },
			assert: func(t *testing.T, cfg *Config, _ string) {
				assert.True(t, cfg.Metrics.Enabled)
			},
		},
		{
			name: "set map entry",
			edit: func(file string) error { return SetValue(file, "runner.envs.ADDED", "yes") },
			assert: func(t *testing.T, cfg *Config, _ string) {
				assert.Equal(t, map[string]string{"EXISTING": "value", "ADDED": "yes"}, cfg.Runner.Envs)
			},
		},
		{
			name: "set replaces a list",
			edit: func(file string) error { return SetValue(file, "runner.labels", "one", "two") },
			assert: func(t *testing.T, cfg *Config, _ string) {
				assert.Equal(t, []string{"one", "two"}, cfg.Runner.Labels)
			},
		},
		{
			name: "add appends to a list",
			edit: func(file string) error { return AddValue(file, "runner.labels", "ubuntu:docker://node:22") },
			assert: func(t *testing.T, cfg *Config, _ string) {
				assert.Equal(t, []string{"ubuntu-latest:docker://node:20", "self-hosted", "ubuntu:docker://node:22"}, cfg.Runner.Labels)
			},
		},
		{
			name: "remove drops a list entry",
			edit: func(file string) error { return RemoveValue(file, "runner.labels", "self-hosted") },
			assert: func(t *testing.T, cfg *Config, _ string) {
				assert.Equal(t, []string{"ubuntu-latest:docker://node:20"}, cfg.Runner.Labels)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := writeEditFixture(t)
			require.NoError(t, tt.edit(file))

			content, err := os.ReadFile(file)
			require.NoError(t, err)
			cfg, err := LoadDefault(file)
			require.NoError(t, err)

			tt.assert(t, cfg, string(content))

			// Every edit keeps the comments, the section spacing and the untouched keys.
			assert.Contains(t, string(content), "# A leading comment.")
			assert.Contains(t, string(content), "  # The logging level.")
			assert.Contains(t, string(content), "\n\nrunner:")
			assert.Equal(t, "info", cfg.Log.Level)
			assert.Equal(t, "value", cfg.Runner.Envs["EXISTING"])
		})
	}
}

func TestEditValuesRejectsBadInput(t *testing.T) {
	tests := []struct {
		name    string
		edit    func(file string) error
		wantErr string
	}{
		{
			name:    "unknown key",
			edit:    func(file string) error { return SetValue(file, "runner.labl", "x") },
			wantErr: `unknown config key "runner.labl"`,
		},
		{
			name:    "unknown section",
			edit:    func(file string) error { return SetValue(file, "nope.key", "x") },
			wantErr: `unknown config key "nope"`,
		},
		{
			name:    "value is not a number",
			edit:    func(file string) error { return SetValue(file, "runner.capacity", "many") },
			wantErr: `"many" is not a valid int`,
		},
		{
			name:    "value is not a duration",
			edit:    func(file string) error { return SetValue(file, "runner.timeout", "soon") },
			wantErr: `"soon" is not a duration`,
		},
		{
			name:    "value is not a boolean",
			edit:    func(file string) error { return SetValue(file, "runner.insecure", "maybe") },
			wantErr: `"maybe" is not a boolean`,
		},
		{
			name:    "set needs a single value",
			edit:    func(file string) error { return SetValue(file, "runner.capacity", "1", "2") },
			wantErr: "takes exactly one value",
		},
		{
			name:    "set on a section",
			edit:    func(file string) error { return SetValue(file, "runner", "x") },
			wantErr: "is a section",
		},
		{
			name:    "add on a scalar",
			edit:    func(file string) error { return AddValue(file, "runner.capacity", "4") },
			wantErr: "is not a list",
		},
		{
			name:    "add a duplicate",
			edit:    func(file string) error { return AddValue(file, "runner.labels", "self-hosted") },
			wantErr: `already contains "self-hosted"`,
		},
		{
			name:    "remove a missing entry",
			edit:    func(file string) error { return RemoveValue(file, "runner.labels", "absent") },
			wantErr: `does not contain "absent"`,
		},
		{
			name:    "sub-key of a free-form map entry",
			edit:    func(file string) error { return SetValue(file, "runner.envs.A.B", "x") },
			wantErr: "has no sub-keys",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := writeEditFixture(t)
			err := tt.edit(file)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)

			content, err := os.ReadFile(file)
			require.NoError(t, err)
			assert.Equal(t, editFixture, string(content), "a rejected edit must leave the file untouched")
		})
	}
}

func TestEditValuesMissingFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "absent.yaml")
	err := SetValue(file, "runner.capacity", "4")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetValue(t *testing.T) {
	file := writeEditFixture(t)

	value, err := GetValue(file, "runner.capacity")
	require.NoError(t, err)
	assert.Equal(t, "1", value)

	value, err = GetValue(file, "runner.labels")
	require.NoError(t, err)
	assert.Equal(t, "ubuntu-latest:docker://node:20\nself-hosted", value)

	value, err = GetValue(file, "runner.envs")
	require.NoError(t, err)
	assert.Equal(t, "EXISTING=value", value)

	_, err = GetValue(file, "metrics.addr")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not set")
}

// The example config is the file users edit, so it has to survive a round-trip.
func TestEditValuesKeepsExampleConfigIntact(t *testing.T) {
	file := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(file, Example, 0o600))

	require.NoError(t, AddValue(file, "runner.labels", "ubuntu:docker://node:22"))

	content, err := os.ReadFile(file)
	require.NoError(t, err)

	const added = "    - ubuntu:docker://node:22"
	var (
		withoutAdded []string
		found        bool
	)
	for line := range strings.SplitSeq(string(content), "\n") {
		if !found && line == added {
			found = true
			continue
		}
		withoutAdded = append(withoutAdded, line)
	}

	require.True(t, found, "the label was not appended")
	assert.Equal(t, strings.Split(string(Example), "\n"), withoutAdded, "only the appended label may differ")
}
