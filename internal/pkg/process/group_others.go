// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package process

import "os"

// Group is a no-op outside Windows, where the task-scoped Job Object has no
// direct equivalent: a step's tree is already reclaimed through its process
// group (see SysProcAttr and Killer), and a surviving orphan does not block the
// workspace deletion at teardown the way an open Windows file handle does.
type Group struct{}

// NewGroup returns a Group that owns nothing.
func NewGroup() (*Group, error) {
	return &Group{}, nil
}

// Assign does nothing.
func (g *Group) Assign(_ *os.Process) error {
	return nil
}

// Close does nothing.
func (g *Group) Close() error {
	return nil
}
