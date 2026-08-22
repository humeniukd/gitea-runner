// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestJobOutputs(t *testing.T) {
	requireFixture(t)
	t.Parallel()

	api, repo, _ := startScenario(t, "outputs.yml")

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)
	if !strings.Contains(logs, "received=produced-value-42") {
		t.Fatalf("the consuming job did not receive the produced output: %s", logs)
	}
}
