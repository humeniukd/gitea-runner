// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitea.com/gitea/act_runner/internal/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunnerCleanupStaleTaskDirs(t *testing.T) {
	now := time.Date(2026, time.April, 29, 20, 0, 0, 0, time.UTC)
	workdirRoot := filepath.Join(t.TempDir(), "workspace")
	require.NoError(t, os.MkdirAll(workdirRoot, 0o700))

	oldTask := filepath.Join(workdirRoot, "1001")
	freshTask := filepath.Join(workdirRoot, "1002")
	nonTask := filepath.Join(workdirRoot, "shared")
	for _, path := range []string{oldTask, freshTask, nonTask} {
		require.NoError(t, os.MkdirAll(path, 0o700))
	}

	require.NoError(t, os.Chtimes(oldTask, now.Add(-3*time.Hour), now.Add(-3*time.Hour)))
	require.NoError(t, os.Chtimes(freshTask, now.Add(-30*time.Minute), now.Add(-30*time.Minute)))
	require.NoError(t, os.Chtimes(nonTask, now.Add(-5*time.Hour), now.Add(-5*time.Hour)))

	r := &Runner{
		cfg: &config.Config{
			Runner: config.Runner{
				StartupCleanupAge: 2 * time.Hour,
			},
		},
		now: func() time.Time { return now },
	}

	r.cleanupStaleTaskDirs(workdirRoot)

	assert.NoDirExists(t, oldTask)
	assert.DirExists(t, freshTask)
	assert.DirExists(t, nonTask)
}
