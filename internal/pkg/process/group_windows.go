// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package process

import (
	"os"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Group is a task-scoped Windows Job Object that owns every process a job's
// steps start, so that nothing a job spawned can outlive it.
//
// It closes the gap left by the per-step Killer, which only tears a tree down on
// cancellation: a step that runs to completion closes its job handle and leaves
// whatever it spawned still running. A GUI application that re-launches itself,
// or a tool that detaches a child and exits, therefore survives its own step.
// Such an orphan has no parent left, so no process-tree walk can find it, and
// its executable frequently lives outside the workspace (an installed
// application, e.g. Excel), so the workspace path scan at teardown cannot see it
// either.
//
// The job is created with JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE: when Close drops
// the last handle, the kernel terminates every process still assigned to the
// job, whatever its parentage or executable path. This also holds when the
// runner dies unexpectedly — the kernel closes the handle on process exit — so
// a crashed or force-killed runner cannot strand a job's processes on the
// machine.
//
// A process may belong to several jobs (nested jobs, Windows 8+). Assign a
// process here before handing it to NewKiller so the step's job nests inside the
// task's job: Killer.Kill then still terminates exactly the step's own tree on
// cancellation, while Close reaches everything the whole task left behind.
type Group struct {
	mu  sync.Mutex
	job windows.Handle
}

// NewGroup creates the job object. The caller must Close it when the task ends,
// which is what terminates the leftovers.
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

// Assign adds an already-started process to the job. Every process it spawns
// afterwards is assigned automatically, so one call per step process covers the
// whole tree the step goes on to create.
//
// It is safe to call on a nil Group, which is what a caller holds when the job
// object could not be created; leftover reaping then falls back to the
// best-effort workspace scan at teardown.
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

// Close drops the job handle, which terminates every process still assigned to
// the job. It is safe to call on a nil Group and more than once.
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
