// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package git

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func initSHA256Repo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := gitCmd("init", "--object-format=sha256", "--initial-branch=main", dir); err != nil {
		t.Skipf("git does not support SHA256 repositories: %v", err)
	}
	require.NoError(t, cleanGitHooks(dir))
	require.NoError(t, gitCmd("-C", dir, "commit", "--allow-empty", "-m", "init"))
	return dir
}

func TestGitCLIReadsSHA256Repository(t *testing.T) {
	dir := initSHA256Repo(t)
	want := gitRevParse(t, dir, "HEAD")
	require.Len(t, want, 64)

	shortSHA, sha, err := FindGitRevision(context.Background(), dir)
	require.NoError(t, err)
	require.Equal(t, want, sha)
	require.Equal(t, want[:7], shortSHA)

	ref, err := FindGitRef(context.Background(), dir)
	require.NoError(t, err)
	require.Equal(t, "refs/heads/main", ref)

	require.NoError(t, gitCmd("-C", dir, "remote", "add", "origin", "https://github.com/owner/repo.git"))
	slug, err := FindGithubRepo(context.Background(), dir, "github.com")
	require.NoError(t, err)
	require.Equal(t, "owner/repo", slug)
}
