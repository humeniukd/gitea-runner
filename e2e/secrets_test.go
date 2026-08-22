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

// TestSecretsMasking confirms a repository secret interpolated into a step's log is masked
// with the runner's literal "***" marker (act/runner/logger.go's AppendSecretMasker), not
// leaked in plaintext.
func TestSecretsMasking(t *testing.T) {
	requireFixture(t)
	t.Parallel() // own repo and runner, so it shares nothing with the other scenarios
	ctx := context.Background()

	const secretValue = "s3cr3t-value-xyz"

	repo := repoName(t)
	api := fixture.API()
	if err := api.CreateRepo(ctx, repo, true); err != nil {
		t.Fatalf("create repo: %v", err)
	}
	if err := api.CreateSecret(ctx, giteaAdminUser, repo, "FOO", secretValue); err != nil {
		t.Fatalf("create secret: %v", err)
	}

	workflow, err := os.ReadFile("testdata/workflows/secrets.yml")
	if err != nil {
		t.Fatalf("read workflow fixture: %v", err)
	}
	if err := api.CreateFile(ctx, giteaAdminUser, repo, ".gitea/workflows/secrets.yml", string(workflow), "add workflow"); err != nil {
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
	if strings.Contains(logs, secretValue) {
		t.Fatalf("job logs leak the raw secret value: %s", logs)
	}
	if !strings.Contains(logs, "***") {
		t.Fatalf("job logs missing the mask marker %q: %s", "***", logs)
	}
}
