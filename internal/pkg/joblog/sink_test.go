// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package joblog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testStart = time.Date(2026, 8, 14, 9, 12, 3, 0, time.UTC)

func read(t *testing.T, sink *Sink, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(sink.Dir(), name))
	require.NoError(t, err)
	return string(content)
}

func TestSink_Disabled(t *testing.T) {
	sink := Open(Options{TaskID: 1, StartedAt: testStart})
	require.Nil(t, sink)

	// Every method tolerates the disabled sink, so callers need no branches.
	assert.Empty(t, sink.Dir())
	sink.Job(testStart, "hello")
	sink.Diag(testStart, "hello")
	sink.Close("done")
}

func TestSink_WritesBothFiles(t *testing.T) {
	root := t.TempDir()
	sink := Open(Options{Dir: root, TaskID: 41, StartedAt: testStart, Header: "task=41 repo=org/repo"})
	require.NotNil(t, sink)

	assert.Equal(t, filepath.Join(root, "20260814-091203-task-41"), sink.Dir())

	sink.Job(testStart, "##[group]Run actions/checkout@v4")
	sink.Diag(testStart, "DEBUG stage=Pre | resolving the action")
	sink.Job(testStart.Add(time.Millisecond), "line\nwith a newline")
	sink.Close("task 41 finished: success")

	job := read(t, sink, JobFileName)
	assert.Equal(t, []string{
		"# task=41 repo=org/repo",
		"2026-08-14T09:12:03.000Z ##[group]Run actions/checkout@v4",
		`2026-08-14T09:12:03.001Z line\nwith a newline`,
	}, strings.Split(job, "\n")[:3])
	assert.Contains(t, job, "[runner] task 41 finished: success")
	assert.NotContains(t, job, "resolving the action", "job.log holds only what was uploaded")

	diag := read(t, sink, DiagFileName)
	assert.Contains(t, diag, "# task=41 repo=org/repo")
	assert.Contains(t, diag, "2026-08-14T09:12:03.000Z DEBUG stage=Pre | resolving the action")
	assert.Contains(t, diag, "[runner] task 41 finished: success")
	assert.NotContains(t, diag, "actions/checkout", "the uploaded rows stay in job.log")
}

func TestSink_CloseIsIdempotent(t *testing.T) {
	sink := Open(Options{Dir: t.TempDir(), TaskID: 1, StartedAt: testStart})
	require.NotNil(t, sink)

	sink.Job(testStart, "hello")
	sink.Close("first")
	sink.Close("second")

	job := read(t, sink, JobFileName)
	assert.Contains(t, job, "first")
	assert.NotContains(t, job, "second")
}

func TestSink_MaxSize(t *testing.T) {
	sink := Open(Options{Dir: t.TempDir(), TaskID: 1, StartedAt: testStart, MaxSize: 120})
	require.NotNil(t, sink)

	for range 10 {
		sink.Job(testStart, strings.Repeat("x", 40))
	}
	sink.Close("done")

	job := read(t, sink, JobFileName)
	assert.Equal(t, 1, strings.Count(job, "log.job_max_size"), "the cap is reported once")
	assert.NotContains(t, job, "done", "nothing is written past the cap")
	assert.Less(t, len(job), 300)

	// The cap is per file, so the other one is unaffected.
	assert.NotContains(t, read(t, sink, DiagFileName), "log.job_max_size")
}

func TestSink_DroppedLinesAreReported(t *testing.T) {
	sink := open(Options{Dir: t.TempDir(), TaskID: 1, StartedAt: testStart}, 1)
	require.NotNil(t, sink)

	// More than the queue holds, so the writer goroutine cannot keep up with all of them.
	for range 5000 {
		sink.Job(testStart, "flood")
	}
	sink.Close("done")

	job := read(t, sink, JobFileName)
	if dropped := 5000 - strings.Count(job, "flood"); dropped > 0 {
		assert.Contains(t, job, "dropped")
	}
	assert.Contains(t, job, "[runner] done", "the trailer is written even after drops")
}

func TestPrune(t *testing.T) {
	root := t.TempDir()
	expired := filepath.Join(root, "20200101-000000-task-1")
	fresh := filepath.Join(root, testStart.Format(dirNameLayout)+"-task-2")
	unrelated := filepath.Join(root, "operator-data")
	for _, dir := range []string{expired, fresh, unrelated} {
		require.NoError(t, os.MkdirAll(dir, 0o700))
	}
	loose := filepath.Join(root, "20200101-000000-task-3")
	require.NoError(t, os.WriteFile(loose, []byte("not a directory"), 0o600))

	sink := Open(Options{Dir: root, TaskID: 4, StartedAt: testStart, Retention: 24 * time.Hour})
	require.NotNil(t, sink)
	sink.Close("done")

	assert.NoDirExists(t, expired)
	assert.DirExists(t, fresh)
	assert.DirExists(t, unrelated)
	assert.FileExists(t, loose)
	assert.DirExists(t, sink.Dir(), "the task's own directory is never pruned")
}

func TestPrune_RetentionDisabled(t *testing.T) {
	root := t.TempDir()
	expired := filepath.Join(root, "20200101-000000-task-1")
	require.NoError(t, os.MkdirAll(expired, 0o700))

	sink := Open(Options{Dir: root, TaskID: 2, StartedAt: testStart})
	require.NotNil(t, sink)
	sink.Close("done")

	assert.DirExists(t, expired)
}
