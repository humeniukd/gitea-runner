// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package git

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// initSHA256Repo creates a SHA256 object-format repository with a single commit,
// mirroring the extensions Gitea records on such repositories: objectformat=sha256
// and compatObjectFormat=sha1. The compat extension is configured after the commit
// so the test does not depend on the host git being built with full sha1<->sha256
// interop. It skips the test when the host git cannot create SHA256 repositories.
func initSHA256Repo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	if err := gitCmd("init", "--object-format=sha256", "--initial-branch=main", dir); err != nil {
		t.Skipf("git does not support SHA256 repositories: %v", err)
	}
	require.NoError(t, cleanGitHooks(dir))
	require.NoError(t, gitCmd("-C", dir, "commit", "--allow-empty", "-m", "init"))
	require.NoError(t, gitCmd("-C", dir, "config", "extensions.compatObjectFormat", "sha1"))

	return dir
}

// TestFindGitRevisionSHA256 is a regression test for reading SHA256 repositories.
// go-git v5 could not open them at all (it rejected the objectformat and
// compatobjectformat extensions), and its fixed 20-byte plumbing.Hash truncated
// 64-hex SHA256 object names to 40 characters. The runner must instead resolve the
// full, correct revision.
func TestFindGitRevisionSHA256(t *testing.T) {
	dir := initSHA256Repo(t)

	want := gitRevParse(t, dir, "HEAD")
	require.Len(t, want, 64, "sanity: SHA256 object names are 64 hex characters")

	shortSha, sha, err := FindGitRevision(context.Background(), dir)
	require.NoError(t, err)
	require.Equal(t, want, sha)
	require.Equal(t, want[:7], shortSha)
}

// TestFindGitRefSHA256 ensures ref resolution works on SHA256 repositories.
func TestFindGitRefSHA256(t *testing.T) {
	dir := initSHA256Repo(t)

	ref, err := FindGitRef(context.Background(), dir)
	require.NoError(t, err)
	require.Equal(t, "refs/heads/main", ref)
}

// TestFindGithubRepoSHA256 ensures the remote URL (and derived slug) can be read
// from SHA256 repositories.
func TestFindGithubRepoSHA256(t *testing.T) {
	dir := initSHA256Repo(t)
	require.NoError(t, gitCmd("-C", dir, "remote", "add", "origin", "https://github.com/owner/repo.git"))

	slug, err := FindGithubRepo(context.Background(), dir, "github.com", "")
	require.NoError(t, err)
	require.Equal(t, "owner/repo", slug)
}
