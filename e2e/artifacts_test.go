// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestArtifactUploadDownload exercises upload-artifact in one job and download-artifact in a
// dependent job, verifying the artifact also shows up via the artifacts API.
func TestArtifactUploadDownload(t *testing.T) {
	requireFixture(t)
	t.Parallel() // own repo and runner, so it shares nothing with the other scenarios
	ctx := context.Background()

	repo := repoName(t)
	api := fixture.API()
	if err := api.CreateRepo(ctx, repo, true); err != nil {
		t.Fatalf("create repo: %v", err)
	}

	workflow, err := os.ReadFile("testdata/workflows/artifacts.yml")
	if err != nil {
		t.Fatalf("read workflow fixture: %v", err)
	}
	if err := api.CreateFile(ctx, giteaAdminUser, repo, ".gitea/workflows/artifacts.yml", string(workflow), "add workflow"); err != nil {
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
		dumpRunLogs(t, api, repo, run.ID)
		t.Fatalf("run concluded %q, want success", completed.Conclusion)
	}

	artifacts, err := api.Artifacts(ctx, giteaAdminUser, repo, run.ID)
	if err != nil {
		t.Fatalf("list artifacts: %v", err)
	}
	if len(artifacts) != 1 || artifacts[0].Name != "payload" {
		t.Fatalf("got artifacts %+v, want one named payload", artifacts)
	}
}
