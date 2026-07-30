// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCacheV2TriState(t *testing.T) {
	for _, tc := range []struct {
		name, yaml string
		want       bool
	}{
		{"absent defaults to true", "cache:\n  dir: /tmp/x\n", true},
		{"explicit true", "cache:\n  dir: /tmp/x\n  v2: true\n", true},
		{"explicit false stays false", "cache:\n  dir: /tmp/x\n  v2: false\n", false},
		{"empty file defaults to true", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := filepath.Join(t.TempDir(), "c.yaml")
			require.NoError(t, os.WriteFile(f, []byte(tc.yaml), 0o600))
			cfg, err := LoadDefault(f)
			require.NoError(t, err)
			require.Equal(t, tc.want, cfg.Cache.V2 != nil && *cfg.Cache.V2)
		})
	}
}

// The shipped example must parse and must not warn about unknown keys.
func TestExampleConfigParses(t *testing.T) {
	cfg, err := LoadDefault("config.example.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg.Cache.V2)
	require.True(t, *cfg.Cache.V2)
}
