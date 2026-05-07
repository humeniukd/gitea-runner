// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"gitea.com/gitea/runner/act/common/git"
	"gitea.com/gitea/runner/act/model"

	"github.com/stretchr/testify/require"
)

// Regression test for go-gitea/gitea#37483: a remote reusable workflow at a moving
// ref (branch/tag) must reflect the new tip on every invocation, not stay pinned
// to the cache populated on the first run.
func TestReusableWorkflowCachedBranchRefRefreshes(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	remoteDir := t.TempDir()
	gitMust(t, "", "init", "--bare", "--initial-branch=master", remoteDir)

	workDir := t.TempDir()
	gitMust(t, "", "clone", remoteDir, workDir)
	gitMust(t, workDir, "config", "user.email", "test@test")
	gitMust(t, workDir, "config", "user.name", "test")
	gitMust(t, workDir, "checkout", "-b", "master")

	const workflowPath = ".gitea/workflows/reusable.yml"
	tmpl := func(tag string) string {
		return "name: reusable\non:\n  workflow_call:\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo " + tag + "\n"
	}

	require.NoError(t, os.MkdirAll(filepath.Join(workDir, ".gitea/workflows"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(workDir, workflowPath), []byte(tmpl("v1")), 0o644))
	gitMust(t, workDir, "add", workflowPath)
	gitMust(t, workDir, "commit", "-m", "v1")
	gitMust(t, workDir, "push", "-u", "origin", "master")

	rc := &RunContext{
		Config: &Config{},
		Run: &model.Run{
			JobID: "j1",
			Workflow: &model.Workflow{
				Name: "wf",
				Jobs: map[string]*model.Job{"j1": {}},
			},
		},
	}
	cacheDir := t.TempDir()

	require.NoError(t, cloneRemoteReusableWorkflow(rc, remoteDir, "master", cacheDir, "")(context.Background()))
	got, err := os.ReadFile(filepath.Join(cacheDir, workflowPath))
	require.NoError(t, err)
	require.Equal(t, tmpl("v1"), string(got))

	// Branch tip moves; cache key (cacheDir) does not.
	require.NoError(t, os.WriteFile(filepath.Join(workDir, workflowPath), []byte(tmpl("v2")), 0o644))
	gitMust(t, workDir, "commit", "-am", "v2")
	gitMust(t, workDir, "push", "origin", "master")

	require.NoError(t, cloneRemoteReusableWorkflow(rc, remoteDir, "master", cacheDir, "")(context.Background()))
	got, err = os.ReadFile(filepath.Join(cacheDir, workflowPath))
	require.NoError(t, err)
	require.Equal(t, tmpl("v2"), string(got), "cached workflow file must reflect the updated branch tip")
}

// TestNewReusableWorkflowExecutorHoldsCloneLock guards the fix for the action-cache
// race in newReusableWorkflowExecutor: model.NewWorkflowPlanner reads the workflow
// yaml from <ActionCacheDir>/<safeFilename(uses)>, which a peer matrix job's
// cloneRemoteReusableWorkflow can `git checkout --force` mid-parse if the lock
// is not held during the read.
func TestNewReusableWorkflowExecutorHoldsCloneLock(t *testing.T) {
	workflowDir := t.TempDir()
	const workflowFile = "reusable.yml"
	require.NoError(t, os.WriteFile(filepath.Join(workflowDir, workflowFile), []byte(
		"name: reusable\non:\n  workflow_call:\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hi\n",
	), 0o644))

	unlockOnce := sync.OnceFunc(git.AcquireCloneLock(workflowDir))
	defer unlockOnce()

	rc := &RunContext{
		Config: &Config{},
		Run: &model.Run{
			JobID: "build",
			Workflow: &model.Workflow{
				Jobs: map[string]*model.Job{"build": {}},
			},
		},
	}
	exec := newReusableWorkflowExecutor(rc, workflowDir, workflowFile)

	done := make(chan error, 1)
	go func() { done <- exec(context.Background()) }()

	select {
	case err := <-done:
		t.Fatalf("executor returned without taking the clone lock: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	unlockOnce()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("executor did not return after clone lock was released")
	}
}

func gitMust(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, string(out))
}
