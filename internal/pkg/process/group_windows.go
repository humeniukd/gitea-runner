// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package process

import (
	"os"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Group is a task-scoped Windows Job Object holding every process a job's steps
// start. Closing it terminates whatever is still assigned, catching orphans the
// per-step Killer misses: a step that completes leaves its spawn running, with
// no parent to walk from and often an executable outside the workspace.
type Group struct {
	mu  sync.Mutex
	job windows.Handle
}

// NewGroup creates the job object. Closing it is what terminates the leftovers,
// so the caller must Close it when the task ends.
func NewGroup() (*Group, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	); err != nil {
		_ = windows.CloseHandle(job)
		return nil, err
	}

	return &Group{job: job}, nil
}

// Assign adds a started process to the job; anything it spawns afterwards joins
// automatically. Call it before NewKiller so the step's job nests inside this
// one. Safe on a nil Group, which is what a caller holds when NewGroup failed.
func (g *Group) Assign(p *os.Process) error {
	if g == nil || p == nil {
		return nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	if g.job == 0 {
		return nil
	}

	h, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(p.Pid))
	if err != nil {
		return err
	}
	defer func() { _ = windows.CloseHandle(h) }()

	return windows.AssignProcessToJobObject(g.job, h)
}

// Close drops the job handle, terminating every process still assigned to it.
// Safe on a nil Group and on repeat calls.
func (g *Group) Close() error {
	if g == nil {
		return nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	if g.job == 0 {
		return nil
	}

	h := g.job
	g.job = 0
	return windows.CloseHandle(h)
}
