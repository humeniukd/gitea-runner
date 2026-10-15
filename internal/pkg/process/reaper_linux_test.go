// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package process

import (
	"bufio"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// startSupervised starts script under supervise, the way RunAsInit does, with
// the test process as subreaper standing in for PID 1.
func startSupervised(t *testing.T, script string) (*bufio.Reader, chan os.Signal, <-chan int) {
	require.NoError(t, unix.Prctl(unix.PR_SET_CHILD_SUBREAPER, 1, 0, 0, 0))
	t.Cleanup(func() { _ = unix.Prctl(unix.PR_SET_CHILD_SUBREAPER, 0, 0, 0, 0) })

	signals := make(chan os.Signal, 32)
	signal.Notify(signals, syscall.SIGCHLD)
	t.Cleanup(func() { signal.Stop(signals) })

	cmd := exec.Command("sh", "-c", script)
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	require.NoError(t, cmd.Start())
	code := make(chan int, 1)
	go func() { code <- supervise(cmd.Process, signals) }()
	return bufio.NewReader(stdout), signals, code
}

func TestSuperviseReapsOrphansAndForwardsSignals(t *testing.T) {
	// The inner shell exits at once, orphaning the backgrounded true.
	out, signals, code := startSupervised(t, `trap 'kill $!; exit 7' TERM; sh -c 'true & echo $!'; sleep 60 & wait`)

	line, err := out.ReadString('\n')
	require.NoError(t, err)
	orphan, err := strconv.Atoi(strings.TrimSpace(line))
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		_, err := os.Stat("/proc/" + strconv.Itoa(orphan))
		return os.IsNotExist(err)
	}, 10*time.Second, 10*time.Millisecond, "orphan %d was not reaped", orphan)

	signals <- syscall.SIGTERM
	select {
	case c := <-code:
		require.Equal(t, 7, c)
	case <-time.After(10 * time.Second):
		t.Fatal("child did not exit on the forwarded SIGTERM")
	}
}

func TestSuperviseReturnsSignalExitCode(t *testing.T) {
	_, _, code := startSupervised(t, `kill -KILL $$`)
	require.Equal(t, 128+int(syscall.SIGKILL), <-code)
}
