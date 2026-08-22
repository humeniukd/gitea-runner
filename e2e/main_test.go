// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

const runnerLabel = "ubuntu-latest:docker://node:24-bookworm-slim" // pre-pulled in CI/docs scripts

const runTimeout = 3 * time.Minute

var nonRepoChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func repoName(t *testing.T) string {
	return strings.ToLower(nonRepoChars.ReplaceAllString(t.Name(), "-"))
}

var fixture *GiteaFixture

var skipReason string // set when docker is unavailable; tests call t.Skip

func TestMain(m *testing.M) {
	os.Exit(runSuite(m))
}

func runSuite(m *testing.M) int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if _, err := dockerClient(ctx); err != nil {
		skipReason = fmt.Sprintf("docker unavailable: %v", err)
		return m.Run()
	}

	f, err := StartGitea(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start gitea fixture: %v\n", err)
		return 1
	}
	fixture = f
	defer func() { _ = fixture.Close(context.Background()) }()

	fmt.Fprintf(os.Stderr, "gitea fixture: %s\n", fixture.Describe())

	return m.Run()
}

func startScenario(t *testing.T, workflow string, opts ...HarnessOption) (*GiteaAPI, string, *Harness) {
	t.Helper()
	ctx := context.Background()

	repo := repoName(t)
	api := fixture.API()
	if err := api.CreateRepo(ctx, repo, true); err != nil {
		t.Fatalf("create repo: %v", err)
	}
	pushWorkflow(t, api, repo, workflow)

	token, err := fixture.RegistrationToken(ctx)
	if err != nil {
		t.Fatalf("get registration token: %v", err)
	}
	return api, repo, NewHarness(ctx, t, fixture.BaseURL(), token, []string{runnerLabel}, opts...)
}

func pushWorkflow(t *testing.T, api *GiteaAPI, repo, workflow string) {
	t.Helper()

	content, err := os.ReadFile("testdata/workflows/" + workflow)
	if err != nil {
		t.Fatalf("read workflow fixture %s: %v", workflow, err)
	}
	if err := api.CreateFile(context.Background(), giteaAdminUser, repo, ".gitea/workflows/"+workflow, string(content), "add "+workflow); err != nil {
		t.Fatalf("push workflow %s: %v", workflow, err)
	}
}

func requireSuccess(t *testing.T, api *GiteaAPI, repo string, runID int64) {
	t.Helper()

	completed, err := api.WaitForRunConclusion(context.Background(), giteaAdminUser, repo, runID, runTimeout)
	if err != nil {
		dumpRunLogs(t, api, repo, runID)
		t.Fatalf("wait for run: %v", err)
	}
	if completed.Conclusion != "success" {
		dumpRunLogs(t, api, repo, runID)
		t.Fatalf("run concluded %q, want success", completed.Conclusion)
	}
}

func runLogs(t *testing.T, api *GiteaAPI, repo string, runID int64) string {
	t.Helper()
	ctx := context.Background()

	jobs, err := api.Jobs(ctx, giteaAdminUser, repo, runID)
	if err != nil {
		t.Fatalf("list jobs: %v", err)
	}
	var all strings.Builder
	for _, job := range jobs {
		logs, err := api.JobLogs(ctx, giteaAdminUser, repo, job.ID)
		if err != nil {
			t.Fatalf("job logs for %s: %v", job.Name, err)
		}
		all.WriteString(logs)
	}
	return all.String()
}

func dumpRuns(t *testing.T, api *GiteaAPI, repo string) {
	t.Helper()

	runs, err := api.Runs(context.Background(), giteaAdminUser, repo)
	if err != nil {
		t.Logf("dump runs: %v", err)
		return
	}
	if len(runs) == 0 {
		t.Log("the repository has no runs at all")
		return
	}
	for _, run := range runs {
		t.Logf("run %d: path=%q event=%q status=%q conclusion=%q", run.ID, run.Path, run.Event, run.Status, run.Conclusion)
	}
}

func dumpRunLogs(t *testing.T, api *GiteaAPI, repo string, runID int64) {
	t.Helper()
	ctx := context.Background()

	jobs, err := api.Jobs(ctx, giteaAdminUser, repo, runID)
	if err != nil {
		t.Logf("dump run %d: list jobs: %v", runID, err)
		return
	}
	for _, job := range jobs {
		logs, err := api.JobLogs(ctx, giteaAdminUser, repo, job.ID)
		if err != nil {
			t.Logf("job %q (%s): logs unavailable: %v", job.Name, job.Conclusion, err)
			continue
		}
		t.Logf("job %q concluded %q:\n%s", job.Name, job.Conclusion, logs)
	}
}

func requireFixture(t *testing.T) {
	t.Helper()
	if skipReason != "" {
		t.Skip(skipReason)
	}
}
