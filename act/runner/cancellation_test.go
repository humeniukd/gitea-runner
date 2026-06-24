// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"testing"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/exprparser"
	"gitea.com/gitea/runner/act/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCancelledJobStatusEnablesAlwaysAndCancelledSteps verifies that once a job is
// cancelled, getJobContext reports the "cancelled" status so the step `if` functions
// evaluate the way GitHub Actions does: cancelled()/always() are true, success()/failure()
// are false. A step that defaults to success() is therefore skipped while an always() step
// still runs. Before the fix the status could only ever be success/failure, so cancelled()
// was structurally impossible and cancel-only cleanup steps never ran.
func TestCancelledJobStatusEnablesAlwaysAndCancelledSteps(t *testing.T) {
	rc := createIfTestRunContext(map[string]*model.Job{
		"job1": createJob(t, `runs-on: ubuntu-latest`, ""),
	})
	rc.markCancelled()

	// The core fix: the job status context now reports "cancelled" instead of being
	// pinned to success/failure.
	jobCtx := rc.getJobContext()
	require.Equal(t, "cancelled", jobCtx.Status)

	// Feed that status through the step-context expression functions, which is what a
	// step `if` evaluates. On a cancelled job only always()/cancelled() are true.
	interp := exprparser.NewInterpeter(
		&exprparser.EvaluationEnvironment{Job: jobCtx},
		exprparser.Config{Context: "step"},
	)
	for expr, want := range map[string]bool{
		"cancelled()":  true,
		"always()":     true,
		"success()":    false,
		"failure()":    false,
		"!cancelled()": false,
	} {
		got, err := interp.Evaluate(expr, exprparser.DefaultStatusCheckNone)
		require.NoErrorf(t, err, "Evaluate(%q)", expr)
		assert.Equalf(t, want, got, "Evaluate(%q) on a cancelled job", expr)
	}

	// A step without an `if` defaults to success() and must be skipped on cancel,
	// while an `if: always()` step must still run.
	disabled, err := interp.Evaluate("", exprparser.DefaultStatusCheckSuccess)
	require.NoError(t, err)
	assert.Equal(t, false, disabled, "default-success step must be skipped on a cancelled job")

	enabled, err := interp.Evaluate("always()", exprparser.DefaultStatusCheckSuccess)
	require.NoError(t, err)
	assert.Equal(t, true, enabled, "`if: always()` step must run on a cancelled job")
}

// TestMainStepsExecutorRunsAlwaysStepsAfterCancel verifies that newMainStepsExecutor does
// not abandon the remaining steps when the run is cancelled mid-pipeline. The later step
// still runs (so a main-stage always() step is reached), it runs under a fresh,
// non-cancelled context, and the job is marked cancelled. The interrupt error is still
// propagated so callers up the chain see the cancellation.
func TestMainStepsExecutorRunsAlwaysStepsAfterCancel(t *testing.T) {
	rc := createIfTestRunContext(map[string]*model.Job{
		"job1": createJob(t, `runs-on: ubuntu-latest`, ""),
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var ran []string
	var laterStepCtxErr error
	steps := []common.Executor{
		func(_ context.Context) error {
			ran = append(ran, "step1")
			cancel() // server cancellation lands while step1 runs
			return nil
		},
		func(c context.Context) error {
			ran = append(ran, "always-step")
			laterStepCtxErr = c.Err()
			return nil
		},
	}

	err := newMainStepsExecutor(rc, steps)(ctx)

	require.ErrorIs(t, err, context.Canceled, "interrupt error is propagated")
	assert.Equal(t, []string{"step1", "always-step"}, ran, "the always() step still runs after cancel")
	require.NoError(t, laterStepCtxErr, "remaining steps run under a fresh, non-cancelled context")
	assert.True(t, rc.jobCancelled, "the job is marked cancelled")
}
