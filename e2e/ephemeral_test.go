// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

// TestEphemeralRunner registers a one-shot runner and runs a job on it. Whether the server
// honours the request is version-dependent: older Gitea accepts the registration but silently
// clears the flag (internal/app/cmd/register.go warns about exactly this), so the downgrade is
// reported rather than failed, and the job must run either way.
func TestEphemeralRunner(t *testing.T) {
	requireFixture(t)

	api, repo, h := startScenario(t, "basic.yml", WithEphemeral())

	if h.EphemeralDowngraded() {
		t.Log("this gitea version downgraded the ephemeral registration to a normal runner")
	} else {
		t.Log("ephemeral registration honoured")
	}

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)
	if !strings.Contains(logs, "hello") {
		t.Fatalf("ephemeral runner did not produce the job output: %s", logs)
	}
}
