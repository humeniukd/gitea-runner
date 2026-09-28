// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntrypointUsesConfiguredRegistrationFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Docker entrypoint requires bash")
	}
	bash, err := exec.LookPath("bash")
	require.NoError(t, err)
	script, err := os.ReadFile("../../../scripts/run.sh")
	require.NoError(t, err)
	for _, tc := range []struct {
		name                 string
		registered, override bool
	}{
		{"configured existing", true, false},
		{"configured missing with stale default", false, false},
		{"environment override", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			data := filepath.Join(dir, "data dir")
			require.NoError(t, os.MkdirAll(data, 0o700))
			configured := filepath.Join(data, "configured state.json")
			selected := configured
			if tc.override {
				selected = filepath.Join(data, "override state.json")
			}
			require.NoError(t, os.WriteFile(filepath.Join(data, ".runner"), []byte("stale default"), 0o600))
			if tc.registered {
				require.NoError(t, os.WriteFile(selected, []byte("registered"), 0o600))
			}
			t.Setenv("CONFIG_FILE", filepath.Join(dir, "config with spaces.yaml"))
			t.Setenv("RUNNER_STATE_FILE", "")
			if tc.override {
				t.Setenv("RUNNER_STATE_FILE", selected)
			}
			t.Setenv("TEST_CONFIG_STATE", configured)
			argsFile := filepath.Join(dir, "args")
			t.Setenv("TEST_ARGS_FILE", argsFile)
			t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
			stub := `#!/usr/bin/env bash
mode=""
state=""
previous=""
for arg in "$@"; do
  if [[ "$previous" == "--runner-file" ]]; then state="$arg"; fi
  case "$arg" in config|register|daemon) mode="$arg";; esac
  previous="$arg"
done
if [[ "$mode" == config ]]; then printf '%s\n' "$TEST_CONFIG_STATE"; exit 0; fi
printf '%s\n' "$@" >> "$TEST_ARGS_FILE"
if [[ "$mode" == register ]]; then
  printf registered > "$state"
  echo 'Runner registered successfully'
fi
`
			require.NoError(t, os.WriteFile(filepath.Join(dir, "gitea-runner"), []byte(stub), 0o700))
			entrypoint := filepath.Join(dir, "run.sh")
			require.NoError(t, os.WriteFile(entrypoint, []byte(strings.ReplaceAll(string(script), "/data", strconv.Quote(data))), 0o600))
			cmd := exec.CommandContext(t.Context(), bash, entrypoint)
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, string(output))
			args, err := os.ReadFile(argsFile)
			require.NoError(t, err)
			assert.Contains(t, string(args), "--runner-file\n"+selected+"\n")
			assert.Contains(t, string(args), "--config\n"+os.Getenv("CONFIG_FILE")+"\n")
			if tc.registered {
				assert.NotContains(t, string(args), "register\n")
			} else {
				assert.Contains(t, string(args), "register\n")
			}
			require.FileExists(t, selected)
		})
	}
}
