// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package process

// KillProcessesWithCWDUnder is a no-op on non-Windows platforms. Leftover
// process reaping at teardown is only needed on Windows, where an open handle
// (including a process's current directory) blocks deletion of the workspace;
// on Unix os.RemoveAll succeeds regardless, and per-step process groups already
// tear the step's tree down on cancellation.
func KillProcessesWithCWDUnder(_ []string) (int, error) {
	return 0, nil
}
