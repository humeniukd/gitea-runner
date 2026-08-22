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

func TestCapabilityNegotiation(t *testing.T) {
	requireFixture(t)
	t.Parallel()
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

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)
}
