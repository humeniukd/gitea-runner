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

// runnerLabel is the label every scenario registers with: an image already present on most
// dev/CI hosts (docs/scripts/test.yml pre-pulls it too), so job containers don't need a fresh
// pull on top of the Gitea image pull the fixture already does.
const runnerLabel = "ubuntu-latest:docker://node:24-bookworm-slim"

// runTimeout bounds how long a scenario waits for its workflow run to conclude.
const runTimeout = 3 * time.Minute

var nonRepoChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

// repoName derives a Gitea-safe repository name from the test name, so each scenario gets its
// own repo inside the shared fixture instance.
func repoName(t *testing.T) string {
	return strings.ToLower(nonRepoChars.ReplaceAllString(t.Name(), "-"))
}

// fixture is shared across every scenario in this binary run: starting one Gitea container
// (image pull + admin bootstrap) is expensive, and each scenario uses its own repo/runner
// name (see t.Name() usage in Harness/tests) so sharing costs no isolation between them.
var fixture *GiteaFixture

// skipReason, when non-empty, means no docker daemon was reachable at startup; every test
// checks it and calls t.Skip, so `go test -v` shows each test as SKIP rather than the binary
// exiting silently.
var skipReason string

func TestMain(m *testing.M) {
	os.Exit(runSuite(m)) // the only Exit, so every defer in runSuite still runs
}

// runSuite owns the shared fixture's lifetime and returns the exit code for TestMain.
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

	// Straight to stderr, not t.Log, so it shows without -v: it is the only place the run says
	// which Gitea it actually tested, as opposed to which tag it was asked for.
	fmt.Fprintf(os.Stderr, "gitea fixture: %s\n", fixture.Describe())

	return m.Run()
}

// startScenario creates this test's repository, pushes the named workflow fixture into it, and
// starts a runner registered against the shared fixture.
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

// pushWorkflow publishes a testdata workflow into repo's .gitea/workflows directory.
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

// requireSuccess waits for run to finish and fails with its logs unless it succeeded.
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

// runLogs concatenates every job log of a run, for assertions that do not care which job emitted
// a line.
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

// dumpRuns reports every run in repo, so a scenario waiting for a run that never appears says
// what did appear instead of only that it timed out.
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

// dumpRunLogs reports every job's outcome and log for run, so a scenario that fails on the
// server side says why instead of only reporting the conclusion.
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

// requireFixture skips t if TestMain could not start the shared fixture.
func requireFixture(t *testing.T) {
	t.Helper()
	if skipReason != "" {
		t.Skip(skipReason)
	}
}
