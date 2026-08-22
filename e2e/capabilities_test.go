// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"strings"
	"testing"

	"gitea.com/gitea/runner/internal/app/run"
)

// TestCapabilityNegotiation records what the server negotiates back on Declare. The runner
// advertises its capabilities in the request and reads the server's answer from the
// X-Gitea-Actions-Capabilities response header (run.Runner.SetCapabilitiesFromDeclare), a
// side channel that exists precisely because the answer varies by Gitea version. The value
// here is reporting that per-version truth in the matrix output, so the assertion only covers
// what must hold everywhere: Declare succeeds and the runner keeps working.
func TestCapabilityNegotiation(t *testing.T) {
	requireFixture(t)
	t.Parallel() // own repo and runner, so it shares nothing with the other scenarios
	ctx := context.Background()

	repo := repoName(t)
	api := fixture.API()
	if err := api.CreateRepo(ctx, repo, true); err != nil {
		t.Fatalf("create repo: %v", err)
	}
	pushWorkflow(t, api, repo, "basic.yml")

	token, err := fixture.RegistrationToken(ctx)
	if err != nil {
		t.Fatalf("get registration token: %v", err)
	}
	h := NewHarness(ctx, t, fixture.BaseURL(), token, []string{runnerLabel})

	advertised := run.RunnerCapabilities()
	negotiated := h.Capabilities()
	t.Logf("runner advertised %v; server negotiated %q", advertised, negotiated)
	for _, c := range advertised {
		t.Logf("capability %q honoured by this gitea: %t", c, strings.Contains(negotiated, c))
	}

	// A server that negotiates nothing must still run jobs, so the declare path is proven end
	// to end rather than only by the header being readable.
	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)
}
