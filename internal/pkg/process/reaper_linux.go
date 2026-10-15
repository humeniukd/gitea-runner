// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package process

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// RunAsInit runs the current command as a child when this process is PID 1, as
// in a container started without an init, and returns the child's exit code.
// Orphans are reparented to PID 1 and the Go runtime only waits for its own
// children, so without this they would linger as zombies. ok is false when the
// caller should run the command itself.
func RunAsInit() (code int, ok bool) {
	if os.Getpid() != 1 {
		return 0, false
	}
	exe, err := os.Executable()
	if err != nil {
		return 0, false
	}
	signals := make(chan os.Signal, 32)
	signal.Notify(signals) // before Start, so the child's SIGCHLD is not missed
	cmd := exec.Command(exe)
	cmd.Args = os.Args
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		signal.Reset()
		return 0, false
	}
	return supervise(cmd.Process, signals), true
}

// supervise forwards signals to child and reaps every exited process until
// child exits, returning its exit code.
func supervise(child *os.Process, signals <-chan os.Signal) int {
	for sig := range signals {
		switch sig {
		case syscall.SIGURG: // Go runtime preemption, not meant for the child
		case syscall.SIGCHLD:
			if status, exited := reap(child.Pid); exited {
				return exitCode(status)
			}
		default:
			_ = child.Signal(sig)
		}
	}
	return 1
}

// reap waits for all exited processes, reporting whether child was among them.
func reap(child int) (syscall.WaitStatus, bool) {
	for {
		var status syscall.WaitStatus
		pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)
		if errors.Is(err, syscall.EINTR) {
			continue
		}
		if err != nil || pid <= 0 {
			return 0, false
		}
		if pid == child {
			return status, true
		}
	}
}

func exitCode(status syscall.WaitStatus) int {
	if status.Signaled() {
		return 128 + int(status.Signal())
	}
	return status.ExitStatus()
}
