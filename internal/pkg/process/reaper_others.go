// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !linux

package process

// RunAsInit does nothing outside Linux, where the runner does not run as a container's PID 1.
func RunAsInit() (code int, ok bool) {
	return 0, false
}
