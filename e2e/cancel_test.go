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

// Cancel exists only on nightly; older versions must leave the job and runner alone.
func TestRunCancellation(t *testing.T) {
	requireFixture(t)
	ctx := context.Background()

	api, repo, _ := startScenario(t, "cancel.yml")

	wfRun := waitForRun(t, api, repo)
	waitForJob(t, api, repo, wfRun.ID)

	err := api.CancelRun(ctx, giteaAdminUser, repo, wfRun.ID)
	switch {
	case errors.Is(err, ErrCancelUnsupported):
		t.Log("REST run cancellation unavailable; asserting graceful degrade")
		assertRunnerStillHealthy(t, api, repo, wfRun.ID)
	case err != nil:
		dumpRunLogs(t, api, repo, wfRun.ID)
		t.Fatalf("cancel run: %v", err)
	default:
		assertRunCancelled(t, api, repo, wfRun.ID)
	}
}

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
