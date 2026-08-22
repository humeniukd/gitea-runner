// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

// TestJobOutputs confirms a value written to GITHUB_OUTPUT in one job reaches a dependent job
// through needs.<job>.outputs, which crosses the runner/server boundary: the producing job
// reports its outputs to Gitea, and Gitea feeds them into the consuming job's task.
func TestJobOutputs(t *testing.T) {
	requireFixture(t)
	t.Parallel() // own repo and runner, so it shares nothing with the other scenarios

	api, repo, _ := startScenario(t, "outputs.yml")

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)
	if !strings.Contains(logs, "received=produced-value-42") {
		t.Fatalf("the consuming job did not receive the produced output: %s", logs)
	}
}
