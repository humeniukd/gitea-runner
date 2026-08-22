// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

func TestServiceContainer(t *testing.T) {
	requireFixture(t)
	t.Parallel()

	api, repo, _ := startScenario(t, "services.yml")

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)
	if !strings.Contains(logs, "service-reachable-ok") {
		t.Fatalf("the job could not reach its service container: %s", logs)
	}
}
