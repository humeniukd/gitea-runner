// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package process

// KillProcessesWithCWDUnder is a no-op outside Windows: on Unix an open handle
// does not block os.RemoveAll, so there is no workspace to unblock.
func KillProcessesWithCWDUnder(_ []string) (int, error) {
	return 0, nil
}
