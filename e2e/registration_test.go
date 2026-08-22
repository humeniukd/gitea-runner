// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRegisterAndRunWorkflow(t *testing.T) {
	requireFixture(t)
	t.Parallel()
	ctx := context.Background()

	repo := repoName(t)
	api := fixture.API()
	if err := api.CreateRepo(ctx, repo, true); err != nil {
		t.Fatalf("create repo: %v", err)
	}

	workflow, err := os.ReadFile("testdata/workflows/basic.yml")
	if err != nil {
		t.Fatalf("read workflow fixture: %v", err)
	}
	if err := api.CreateFile(ctx, giteaAdminUser, repo, ".gitea/workflows/basic.yml", string(workflow), "add workflow"); err != nil {
		t.Fatalf("push workflow: %v", err)
	}

	token, err := fixture.RegistrationToken(ctx)
	if err != nil {
		t.Fatalf("get registration token: %v", err)
	}
	NewHarness(ctx, t, fixture.BaseURL(), token, []string{runnerLabel})

	run := waitForRun(t, api, repo)
	completed, err := api.WaitForRunConclusion(ctx, giteaAdminUser, repo, run.ID, 3*time.Minute)
	if err != nil {
		t.Fatalf("wait for run: %v", err)
	}
	if completed.Conclusion != "success" {
		t.Fatalf("run concluded %q, want success", completed.Conclusion)
	}

	jobs, err := api.Jobs(ctx, giteaAdminUser, repo, run.ID)
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}
	logs, err := api.JobLogs(ctx, giteaAdminUser, repo, jobs[0].ID)
	if err != nil {
		t.Fatalf("job logs: %v", err)
	}
	if !strings.Contains(logs, "hello") {
		t.Fatalf("job logs missing %q: %s", "hello", logs)
	}
}

func waitForRun(t *testing.T, api *GiteaAPI, repo string) *ActionRun {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for {
		run, err := api.LatestRun(ctx, giteaAdminUser, repo)
		if err != nil && !errors.Is(err, ErrNoRun) {
			t.Fatalf("get latest run: %v", err)
		}
		if run != nil {
			return run
		}
		select {
		case <-ctx.Done():
			t.Fatalf("no run appeared for %s within timeout", repo)
		case <-time.After(pollInterval):
		}
	}
}
