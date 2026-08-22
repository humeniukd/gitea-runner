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

func TestJobMatrix(t *testing.T) {
	requireFixture(t)
	t.Parallel()
	ctx := context.Background()

	repo := repoName(t)
	api := fixture.API()
	if err := api.CreateRepo(ctx, repo, true); err != nil {
		t.Fatalf("create repo: %v", err)
	}

	workflow, err := os.ReadFile("testdata/workflows/matrix.yml")
	if err != nil {
		t.Fatalf("read workflow fixture: %v", err)
	}
	if err := api.CreateFile(ctx, giteaAdminUser, repo, ".gitea/workflows/matrix.yml", string(workflow), "add workflow"); err != nil {
		t.Fatalf("push workflow: %v", err)
	}

	token, err := fixture.RegistrationToken(ctx)
	if err != nil {
		t.Fatalf("get registration token: %v", err)
	}
	NewHarness(ctx, t, fixture.BaseURL(), token, []string{runnerLabel})

	run := waitForRun(t, api, repo)
	completed, err := api.WaitForRunConclusion(ctx, giteaAdminUser, repo, run.ID, 4*time.Minute)
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
	if len(jobs) != 4 {
		t.Fatalf("got %d jobs, want 4 (2x2 matrix)", len(jobs))
	}

	want := map[string]bool{
		"cell-a-1": false, "cell-a-2": false, "cell-b-1": false, "cell-b-2": false,
	}
	for _, job := range jobs {
		logs, err := api.JobLogs(ctx, giteaAdminUser, repo, job.ID)
		if err != nil {
			t.Fatalf("job logs for %s: %v", job.Name, err)
		}
		for cell := range want {
			if strings.Contains(logs, cell) {
				want[cell] = true
			}
		}
	}
	for cell, seen := range want {
		if !seen {
			t.Errorf("no job logged %q", cell)
		}
	}
}
