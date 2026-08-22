// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// TestRunCancellation covers REST run cancellation, which only nightly exposes: 1.26 and 1.27
// register no CancelWorkflowRun route. Rather than skip the older versions outright, it asserts
// they degrade gracefully — the unsupported call must not disturb the in-flight job or the
// runner, which keeps picking up work afterwards.
func TestRunCancellation(t *testing.T) {
	requireFixture(t)
	ctx := context.Background()

	api, repo, _ := startScenario(t, "cancel.yml")

	wfRun := waitForRun(t, api, repo)
	waitForJob(t, api, repo, wfRun.ID) // cancel only means something once the job is picked up

	err := api.CancelRun(ctx, giteaAdminUser, repo, wfRun.ID)
	switch {
	case errors.Is(err, ErrCancelUnsupported):
		t.Log("REST run cancellation is unavailable on this gitea version; asserting it degrades gracefully")
		assertRunnerStillHealthy(t, api, repo, wfRun.ID)
	case err != nil:
		dumpRunLogs(t, api, repo, wfRun.ID)
		t.Fatalf("cancel run: %v", err)
	default:
		assertRunCancelled(t, api, repo, wfRun.ID)
	}
}

// assertRunCancelled expects the cancelled run to finish promptly, well inside the workflow's
// own runtime, so a run that merely ran to completion cannot pass.
func assertRunCancelled(t *testing.T, api *GiteaAPI, repo string, runID int64) {
	t.Helper()

	completed, err := api.WaitForRunConclusion(context.Background(), giteaAdminUser, repo, runID, time.Minute)
	if err != nil {
		dumpRunLogs(t, api, repo, runID)
		t.Fatalf("cancelled run did not finish promptly: %v", err)
	}
	if completed.Conclusion == "success" {
		dumpRunLogs(t, api, repo, runID)
		t.Fatalf("cancelled run concluded %q; it should not have completed successfully", completed.Conclusion)
	}
	t.Logf("cancelled run concluded %q", completed.Conclusion)

	assertPicksUpNewWork(t, api, repo)
}

// assertRunnerStillHealthy checks that a rejected cancel left the original job alone and the
// runner able to take more work.
func assertRunnerStillHealthy(t *testing.T, api *GiteaAPI, repo string, runID int64) {
	t.Helper()

	completed, err := api.WaitForRunConclusion(context.Background(), giteaAdminUser, repo, runID, 5*time.Minute)
	if err != nil {
		dumpRunLogs(t, api, repo, runID)
		t.Fatalf("run did not finish after an unsupported cancel: %v", err)
	}
	if completed.Conclusion != "success" {
		dumpRunLogs(t, api, repo, runID)
		t.Fatalf("run concluded %q after an unsupported cancel, want success", completed.Conclusion)
	}

	assertPicksUpNewWork(t, api, repo)
}

// assertPicksUpNewWork pushes a second workflow and requires the same runner to run it, proving
// the runner survived whatever happened to the first run.
func assertPicksUpNewWork(t *testing.T, api *GiteaAPI, repo string) {
	t.Helper()

	pushWorkflow(t, api, repo, "quick.yml")

	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()
	for {
		runs, err := api.RunsForWorkflow(ctx, giteaAdminUser, repo, "quick.yml")
		if err != nil {
			t.Fatalf("list runs: %v", err)
		}
		if len(runs) > 0 {
			requireSuccess(t, api, repo, runs[0].ID)
			logs := runLogs(t, api, repo, runs[0].ID)
			if !strings.Contains(logs, "still-alive") {
				t.Fatalf("follow-up run did not produce its output: %s", logs)
			}
			return
		}
		select {
		case <-ctx.Done():
			dumpRuns(t, api, repo)
			t.Fatal("the runner never picked up the follow-up workflow")
		case <-time.After(pollInterval):
		}
	}
}
