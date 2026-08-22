// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

// TestServiceContainer runs a job with a service container and reaches it by its alias over the
// job network. It covers the service lifecycle the runner drives itself (creating the network,
// starting the service, waiting on its readiness per container.service_ready_timeout) against a
// task definition the server produced.
func TestServiceContainer(t *testing.T) {
	requireFixture(t)
	t.Parallel() // own repo and runner, so it shares nothing with the other scenarios

	api, repo, _ := startScenario(t, "services.yml")

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)
	if !strings.Contains(logs, "service-reachable-ok") {
		t.Fatalf("the job could not reach its service container: %s", logs)
	}
}
