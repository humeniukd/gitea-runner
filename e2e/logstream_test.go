// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogStreaming(t *testing.T) {
	requireFixture(t)
	ctx := context.Background()

	repo := repoName(t)
	api := fixture.API()
	if err := api.CreateRepo(ctx, repo, true); err != nil {
		t.Fatalf("create repo: %v", err)
	}

	workflow, err := os.ReadFile("testdata/workflows/logstream.yml")
	if err != nil {
		t.Fatalf("read workflow fixture: %v", err)
	}
	if err := api.CreateFile(ctx, giteaAdminUser, repo, ".gitea/workflows/logstream.yml", string(workflow), "add workflow"); err != nil {
		t.Fatalf("push workflow: %v", err)
	}

	token, err := fixture.RegistrationToken(ctx)
	if err != nil {
		t.Fatalf("get registration token: %v", err)
	}
	NewHarness(ctx, t, fixture.BaseURL(), token, []string{runnerLabel})

	run := waitForRun(t, api, repo)
	jobID := waitForJob(t, api, repo, run.ID)

	sawGrowth := false
	lastLen := -1
	pollCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
pollLoop:
	for {
		logs, err := api.JobLogs(pollCtx, giteaAdminUser, repo, jobID)
		if err == nil && len(logs) > 0 {
			if lastLen >= 0 && len(logs) > lastLen {
				sawGrowth = true
			}
			lastLen = len(logs)
			if sawGrowth {
				break pollLoop
			}
		}
		select {
		case <-pollCtx.Done():
			break pollLoop
		case <-time.After(time.Second):
		}
	}
	if !sawGrowth {
		t.Fatalf("job logs never grew while the job was running (last observed length %d)", lastLen)
	}

	completed, err := api.WaitForRunConclusion(ctx, giteaAdminUser, repo, run.ID, 2*time.Minute)
	if err != nil {
		t.Fatalf("wait for run: %v", err)
	}
	if completed.Conclusion != "success" {
		t.Fatalf("run concluded %q, want success", completed.Conclusion)
	}

	finalLogs, err := api.JobLogs(ctx, giteaAdminUser, repo, jobID)
	if err != nil {
		t.Fatalf("final job logs: %v", err)
	}
	if !strings.Contains(finalLogs, "line 5") {
		t.Fatalf("final job logs missing %q: %s", "line 5", finalLogs)
	}
}

func waitForJob(t *testing.T, api *GiteaAPI, repo string, runID int64) int64 {
	t.Helper()
	return waitForJobStatus(t, api, repo, runID, "")
}

func waitForJobStatus(t *testing.T, api *GiteaAPI, repo string, runID int64, status string) int64 {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		jobs, err := api.Jobs(ctx, giteaAdminUser, repo, runID)
		if err != nil {
			t.Fatalf("list jobs: %v", err)
		}
		if len(jobs) > 0 && (status == "" || jobs[0].Status == status) {
			return jobs[0].ID
		}
		select {
		case <-ctx.Done():
			if status == "" {
				t.Fatalf("no job appeared for run %d within timeout", runID)
			}
			t.Fatalf("job for run %d never reached %q within timeout", runID, status)
		case <-time.After(pollInterval):
		}
	}
}

func waitForJobLogContains(t *testing.T, api *GiteaAPI, repo string, jobID int64, substr string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for {
		logs, err := api.JobLogs(ctx, giteaAdminUser, repo, jobID)
		if err == nil && strings.Contains(logs, substr) {
			return
		}
		select {
		case <-ctx.Done():
			t.Fatalf("job %d logs never contained %q", jobID, substr)
		case <-time.After(pollInterval):
		}
	}
}
