// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestWorkflowDispatch(t *testing.T) {
	requireFixture(t)
	t.Parallel()
	ctx := context.Background()

	api, repo, _ := startScenario(t, "dispatch.yml")

	branch := defaultBranch(t, api, repo)
	inputs := map[string]string{"subject": "dispatch-input-value"}
	if err := waitForDispatch(ctx, api, repo, branch, inputs); err != nil {
		t.Fatalf("dispatch workflow: %v", err)
	}

	dispatched := waitForDispatchedRun(t, api, repo)
	requireSuccess(t, api, repo, dispatched.ID)

	logs := runLogs(t, api, repo, dispatched.ID)
	if !strings.Contains(logs, "dispatched-for=dispatch-input-value") {
		t.Fatalf("the dispatch input did not reach the job: %s", logs)
	}
}

// Gitea indexes workflows asynchronously after the Contents API commit.
func waitForDispatch(ctx context.Context, api *GiteaAPI, repo, branch string, inputs map[string]string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	var lastErr error
	for {
		lastErr = api.DispatchWorkflow(ctx, giteaAdminUser, repo, "dispatch.yml", branch, inputs)
		if lastErr == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return lastErr
		case <-time.After(pollInterval):
		}
	}
}

func waitForDispatchedRun(t *testing.T, api *GiteaAPI, repo string) *ActionRun {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()

	for {
		runs, err := api.RunsForWorkflow(ctx, giteaAdminUser, repo, "dispatch.yml")
		if err != nil {
			t.Fatalf("list runs: %v", err)
		}
		for i := range runs {
			if runs[i].Event == "workflow_dispatch" {
				return &runs[i]
			}
		}
		select {
		case <-ctx.Done():
			dumpRuns(t, api, repo)
			t.Fatal("no workflow_dispatch run appeared")
		case <-time.After(pollInterval):
		}
	}
}

func defaultBranch(t *testing.T, api *GiteaAPI, repo string) string {
	t.Helper()

	branch, err := api.DefaultBranch(context.Background(), giteaAdminUser, repo)
	if err != nil {
		t.Fatalf("get default branch: %v", err)
	}
	return branch
}
