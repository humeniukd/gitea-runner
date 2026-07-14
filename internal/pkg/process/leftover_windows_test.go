// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package process

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

// startSleeperIn launches a long-running child whose cwd is dir but whose
// executable (ping.exe) lives outside it — what the path scan cannot see.
func startSleeperIn(t *testing.T, dir string) *exec.Cmd {
	t.Helper()
	ping := filepath.Join(os.Getenv("SystemRoot"), "System32", "ping.exe")
	cmd := exec.Command(ping, "-n", "60", "127.0.0.1")
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start sleeper: %v", err)
	}
	return cmd
}

func TestProcessCWDReadsWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	cmd := startSleeperIn(t, dir)
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()

	// Give the child a moment to initialize its PEB.
	time.Sleep(300 * time.Millisecond)

	cwd, err := processCWD(uint32(cmd.Process.Pid))
	if err != nil {
		t.Fatalf("processCWD returned error: %v", err)
	}
	// Windows reports the current directory with a trailing separator; compare
	// via the same containment check the reaper uses.
	if !dirContainsPath(dir, cwd) {
		t.Fatalf("processCWD = %q, expected it within %q", cwd, dir)
	}
}

func TestKillProcessesWithCWDUnderTerminatesMatch(t *testing.T) {
	dir := t.TempDir()
	// The executable is outside dir and the arguments carry no workspace path,
	// so only the working directory links this process to dir.
	cmd := startSleeperIn(t, dir)
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()

	time.Sleep(300 * time.Millisecond)

	killed, err := KillProcessesWithCWDUnder([]string{dir})
	if err != nil {
		t.Fatalf("KillProcessesWithCWDUnder returned error: %v", err)
	}
	if killed < 1 {
		t.Fatalf("expected at least one process killed, got %d", killed)
	}

	// The process should exit promptly once terminated.
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("process still running after KillProcessesWithCWDUnder")
	}
}

func TestKillProcessesWithCWDUnderSparesUnrelated(t *testing.T) {
	// A process whose cwd is outside the owned dirs must not be touched.
	workspace := t.TempDir()
	elsewhere := t.TempDir()
	cmd := startSleeperIn(t, elsewhere)
	defer func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}()

	time.Sleep(300 * time.Millisecond)

	if _, err := KillProcessesWithCWDUnder([]string{workspace}); err != nil {
		t.Fatalf("KillProcessesWithCWDUnder returned error: %v", err)
	}

	// It must still be alive.
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
		t.Fatal("unrelated process was killed")
	case <-time.After(500 * time.Millisecond):
		// still running, as expected
	}
}
