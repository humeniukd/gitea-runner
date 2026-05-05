// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitea.com/gitea/runner/internal/pkg/config"

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
	alphaNumericTask := filepath.Join(workdirRoot, "123abc")
	for _, path := range []string{oldTask, freshTask, nonTask, alphaNumericTask} {
		require.NoError(t, os.MkdirAll(path, 0o700))
	}

	require.NoError(t, os.Chtimes(oldTask, now.Add(-3*time.Hour), now.Add(-3*time.Hour)))
	require.NoError(t, os.Chtimes(freshTask, now.Add(-30*time.Minute), now.Add(-30*time.Minute)))
	require.NoError(t, os.Chtimes(nonTask, now.Add(-5*time.Hour), now.Add(-5*time.Hour)))
	require.NoError(t, os.Chtimes(alphaNumericTask, now.Add(-5*time.Hour), now.Add(-5*time.Hour)))

	r := &Runner{
		cfg: &config.Config{
			Runner: config.Runner{
				WorkdirCleanupAge: 2 * time.Hour,
			},
		},
		now: func() time.Time { return now },
	}

	r.cleanupStaleTaskDirs(workdirRoot)

	assert.NoDirExists(t, oldTask)
	assert.DirExists(t, freshTask)
	assert.DirExists(t, nonTask)
	assert.DirExists(t, alphaNumericTask)
}

func TestRunnerShouldRunIdleCleanupThrottles(t *testing.T) {
	now := time.Date(2026, time.April, 29, 20, 0, 0, 0, time.UTC)
	r := &Runner{
		cfg: &config.Config{
			Container: config.Container{
				BindWorkdir: true,
			},
			Runner: config.Runner{
				WorkdirCleanupAge:   24 * time.Hour,
				IdleCleanupInterval: time.Hour,
			},
		},
		now: func() time.Time { return now },
	}

	assert.True(t, r.shouldRunIdleCleanup())

	now = now.Add(30 * time.Minute)
	assert.False(t, r.shouldRunIdleCleanup())

	now = now.Add(31 * time.Minute)
	assert.True(t, r.shouldRunIdleCleanup())
}

func TestRunnerShouldRunIdleCleanupSkipsWhenJobRunning(t *testing.T) {
	r := &Runner{
		cfg: &config.Config{
			Container: config.Container{
				BindWorkdir: true,
			},
			Runner: config.Runner{
				WorkdirCleanupAge:   24 * time.Hour,
				IdleCleanupInterval: time.Minute,
			},
		},
		now: time.Now,
	}
	r.runningCount.Store(1)

	assert.False(t, r.shouldRunIdleCleanup())
}

func TestRunnerShouldRunIdleCleanupSkipsWhenBindWorkdirDisabled(t *testing.T) {
	r := &Runner{
		cfg: &config.Config{
			Runner: config.Runner{
				WorkdirCleanupAge:   24 * time.Hour,
				IdleCleanupInterval: time.Minute,
			},
		},
		now: time.Now,
	}

	assert.False(t, r.shouldRunIdleCleanup())
}

func TestRunnerShouldRunIdleCleanupSkipsWhenDisabled(t *testing.T) {
	now := time.Date(2026, time.April, 29, 20, 0, 0, 0, time.UTC)

	t.Run("cleanup age disabled", func(t *testing.T) {
		r := &Runner{
			cfg: &config.Config{
				Container: config.Container{
					BindWorkdir: true,
				},
				Runner: config.Runner{
					WorkdirCleanupAge:   -1,
					IdleCleanupInterval: time.Minute,
				},
			},
			now: func() time.Time { return now },
		}

		assert.False(t, r.shouldRunIdleCleanup())
	})

	t.Run("idle interval disabled", func(t *testing.T) {
		r := &Runner{
			cfg: &config.Config{
				Container: config.Container{
					BindWorkdir: true,
				},
				Runner: config.Runner{
					WorkdirCleanupAge:   24 * time.Hour,
					IdleCleanupInterval: -1,
				},
			},
			now: func() time.Time { return now },
		}

		assert.False(t, r.shouldRunIdleCleanup())
	})
}
