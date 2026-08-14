// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"gitea.com/gitea/runner/internal/pkg/client/mocks"
	"gitea.com/gitea/runner/internal/pkg/config"
	"gitea.com/gitea/runner/internal/pkg/joblog"

	connect_go "connectrpc.com/connect"
	runnerv1 "gitea.dev/actionslib/runner/v1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

var testTime = time.Date(2026, 8, 14, 9, 12, 3, 0, time.UTC)

// newSinkReporter returns a reporter writing its task logs to a temporary directory, with
// updateLog deciding the fate of each UpdateLog call by its one-based number.
func newSinkReporter(t *testing.T, secrets map[string]string, updateLog func(call int) error) *Reporter {
	t.Helper()

	calls := &atomic.Int64{}
	client := mocks.NewClient(t)
	client.On("UpdateLog", mock.Anything, mock.Anything).Return(func(_ context.Context, req *connect_go.Request[runnerv1.UpdateLogRequest]) (*connect_go.Response[runnerv1.UpdateLogResponse], error) {
		if err := updateLog(int(calls.Add(1))); err != nil {
			return nil, err
		}
		return connect_go.NewResponse(&runnerv1.UpdateLogResponse{
			AckIndex: req.Msg.Index + int64(len(req.Msg.Rows)),
		}), nil
	})
	client.On("UpdateTask", mock.Anything, mock.Anything).Return(func(_ context.Context, _ *connect_go.Request[runnerv1.UpdateTaskRequest]) (*connect_go.Response[runnerv1.UpdateTaskResponse], error) {
		return connect_go.NewResponse(&runnerv1.UpdateTaskResponse{}), nil
	})

	taskCtx, err := structpb.NewStruct(map[string]any{"repository": "org/repo", "workflow": "ci.yml"})
	require.NoError(t, err)
	cfg, err := config.LoadDefault("")
	require.NoError(t, err)
	cfg.Log.JobDir = t.TempDir()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	reporter := NewReporter(ctx, cancel, client, &runnerv1.Task{Id: 41, Context: taskCtx, Secrets: secrets}, cfg)
	require.NotNil(t, reporter.sink)
	return reporter
}

func readSink(t *testing.T, reporter *Reporter, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(reporter.sink.Dir(), name))
	require.NoError(t, err)
	return string(content)
}

func fireRaw(t *testing.T, reporter *Reporter, message string) {
	t.Helper()
	require.NoError(t, reporter.Fire(&log.Entry{
		Message: message,
		Level:   log.InfoLevel,
		Data:    log.Fields{"stage": "Main", "stepNumber": 0, "raw_output": true},
	}))
}

func TestReporter_JobLogMirrorsUploadedRows(t *testing.T) {
	reporter := newSinkReporter(t, map[string]string{"TOKEN": "s3cret-value"}, func(int) error { return nil })
	reporter.RunDaemon()
	reporter.ResetSteps(1)

	fireRaw(t, reporter, "the token is s3cret-value")
	fireRaw(t, reporter, "::add-mask::dyn4mic-value")
	fireRaw(t, reporter, "and dyn4mic-value too")
	fireRaw(t, reporter, "::debug::suppressed unless ACTIONS_STEP_DEBUG")
	require.NoError(t, reporter.Fire(&log.Entry{
		Message: "Job succeeded",
		Level:   log.InfoLevel,
		Data:    log.Fields{"stage": "Post", "jobResult": "success"},
	}))
	require.NoError(t, reporter.Close(""))

	job := readSink(t, reporter, joblog.JobFileName)
	assert.Contains(t, job, "task=41 repository=org/repo workflow=ci.yml")
	assert.Contains(t, job, "the token is ***")
	assert.Contains(t, job, "and *** too")
	assert.NotContains(t, job, "s3cret-value")
	assert.NotContains(t, job, "dyn4mic-value")
	assert.NotContains(t, job, "add-mask", "the row carrying the secret never reaches the log")
	assert.NotContains(t, job, "suppressed unless", "job.log holds only the uploaded rows")
	assert.Contains(t, job, "[runner] task 41 finished: success")

	diag := readSink(t, reporter, joblog.DiagFileName)
	assert.Contains(t, diag, "INFO stage=Main stepNumber=0 raw_output=true | the token is ***")
	assert.Contains(t, diag, "::add-mask::***", "an add-mask value is masked in its own entry")
	assert.Contains(t, diag, "::debug::suppressed unless", "diag.log keeps what the upload dropped")
	assert.NotContains(t, diag, "dyn4mic-value")
	assert.NotContains(t, diag, "s3cret-value")
}

func TestReporter_JobLogKeepsRowsWhenUploadFails(t *testing.T) {
	// The final flush fails once and is retried, which is the case the file exists for.
	reporter := newSinkReporter(t, nil, func(call int) error {
		if call == 1 {
			return errors.New("dial tcp: connection refused")
		}
		return nil
	})
	reporter.RunDaemon()
	reporter.ResetSteps(1)

	fireRaw(t, reporter, "output Gitea nearly lost")
	require.NoError(t, reporter.Close(""))

	assert.Contains(t, readSink(t, reporter, joblog.JobFileName), "output Gitea nearly lost")

	diag := readSink(t, reporter, joblog.DiagFileName)
	assert.Contains(t, diag, "[report] UpdateLog(")
	assert.Contains(t, diag, "failed: dial tcp: connection refused")
}

func TestReporter_DiagEntryOutsideSteps(t *testing.T) {
	reporter := newSinkReporter(t, nil, func(int) error { return nil })
	reporter.RunDaemon()

	require.NoError(t, reporter.Fire(&log.Entry{
		Message: "evaluating 'always()'",
		Level:   log.DebugLevel,
		Data:    log.Fields{"stage": "Pre", "stepID": []string{"checkout"}},
	}))
	require.NoError(t, reporter.Close(""))

	// Debug entries are below the default level, so only diag.log keeps this one.
	assert.NotContains(t, readSink(t, reporter, joblog.JobFileName), "evaluating")
	assert.Contains(t, readSink(t, reporter, joblog.DiagFileName), "DEBUG stage=Pre stepID=[checkout] | evaluating 'always()'")
}

func TestTaskHeader(t *testing.T) {
	taskCtx, err := structpb.NewStruct(map[string]any{
		"repository": "org/repo",
		"workflow":   "ci.yml",
		"run_id":     "7",
		"ref":        "refs/heads/main",
		"unrelated":  "dropped",
	})
	require.NoError(t, err)

	header := taskHeader(&runnerv1.Task{Id: 41, Context: taskCtx}, testTime)
	assert.Equal(t,
		"gitea-runner dev task=41 repository=org/repo workflow=ci.yml run_id=7 ref=refs/heads/main started=2026-08-14T09:12:03Z",
		header)
}
