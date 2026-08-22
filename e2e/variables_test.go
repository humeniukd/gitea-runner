// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"strings"
	"testing"
)

func TestRepoVariables(t *testing.T) {
	requireFixture(t)
	t.Parallel()
	ctx := context.Background()

	const variableValue = "hello-from-a-variable"

	repo := repoName(t)
	api := fixture.API()
	if err := api.CreateRepo(ctx, repo, true); err != nil {
		t.Fatalf("create repo: %v", err)
	}
	if err := api.CreateVariable(ctx, giteaAdminUser, repo, "GREETING", variableValue); err != nil {
		t.Fatalf("create variable: %v", err)
	}
	pushWorkflow(t, api, repo, "variables.yml")

	token, err := fixture.RegistrationToken(ctx)
	if err != nil {
		t.Fatalf("get registration token: %v", err)
	}
	NewHarness(ctx, t, fixture.BaseURL(), token, []string{runnerLabel})

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)
	if !strings.Contains(logs, variableValue) {
		t.Fatalf("job logs do not contain the variable value %q: %s", variableValue, logs)
	}
}
