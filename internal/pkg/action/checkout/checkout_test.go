// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checkout

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/internal/pkg/action"

	"gitea.dev/actionslib/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runGit(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@test.com",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null",
	)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	return strings.TrimSpace(string(out))
}

func gitLog(dir string) []string {
	out, err := exec.Command("git", "--git-dir", filepath.Join(dir, ".git"), "log", "--format=%s").Output()
	if err != nil {
		return nil
	}
	return strings.Fields(string(out))
}

func TestCheckout(t *testing.T) {
	work := t.TempDir()
	runGit(t, "init", "--initial-branch=main", work)
	runGit(t, "-C", work, "commit", "--allow-empty", "-m", "c1")
	sha := runGit(t, "-C", work, "rev-parse", "HEAD")
	runGit(t, "-C", work, "tag", "v1")
	runGit(t, "-C", work, "commit", "--allow-empty", "-m", "c2")
	runGit(t, "-C", work, "checkout", "-b", "feature")
	runGit(t, "-C", work, "commit", "--allow-empty", "-m", "f1")
	runGit(t, "-C", work, "checkout", "--orphan", "trunk")
	runGit(t, "-C", work, "commit", "--allow-empty", "-m", "t1")
	runGit(t, "-C", work, "commit", "--allow-empty", "-m", "t2")
	runGit(t, "-C", work, "tag", "trunk", "main")
	server := t.TempDir()
	runGit(t, "clone", "--bare", "--branch", "main", work, filepath.Join(server, "org", "repo"))
	runGit(t, "clone", "--bare", "--branch", "trunk", work, filepath.Join(server, "other", "repo"))
	sha256Repo := filepath.Join(server, "org", "sha256")
	runGit(t, "init", "--object-format=sha256", "--initial-branch=main", sha256Repo)
	runGit(t, "-C", sha256Repo, "commit", "--allow-empty", "-m", "s1")
	sha256 := runGit(t, "-C", sha256Repo, "rev-parse", "HEAD")
	host := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(host, "marker"), nil, 0o644))

	for _, tc := range []struct {
		name          string
		bind, noSkip  bool
		before        map[string]string
		inputs        map[string]string
		dst           string
		log, refs     []string
		hostCopy, err bool
	}{
		{name: "bound workspace is skipped", bind: true},
		{name: "workflow ref copies host workdir", inputs: map[string]string{"ref": "main", "path": "a/../src"}, dst: "src", hostCopy: true},
		{name: "bound path clones event sha", bind: true, inputs: map[string]string{"path": "src"}, dst: "src", log: []string{"c1"}},
		{name: "no-skip clones event sha", noSkip: true, log: []string{"c1"}},
		{name: "ref clones shallow", inputs: map[string]string{"ref": "refs/heads/feature"}, log: []string{"f1"}},
		{name: "tag ref wins over same-named branch", inputs: map[string]string{"ref": "refs/tags/trunk"}, log: []string{"c2"}},
		{
			name: "fetch-depth 0 deepens a shallow checkout to every branch and tag", before: map[string]string{"ref": "feature"},
			inputs: map[string]string{"ref": "feature", "fetch-depth": "0"}, log: []string{"f1", "c2", "c1"},
			refs: []string{"refs/remotes/origin/feature", "refs/remotes/origin/main", "refs/remotes/origin/trunk", "refs/tags/trunk", "refs/tags/v1"},
		},
		{name: "other repository defaults to its default branch", inputs: map[string]string{"repository": "other/repo"}, log: []string{"t2"}},
		{name: "sha256 repository", inputs: map[string]string{"repository": "org/sha256", "ref": sha256}, log: []string{"s1"}},
		{name: "unsupported input", inputs: map[string]string{"submodules": "true"}, err: true},
		{name: "relative path escaping workspace", inputs: map[string]string{"path": "../x"}, err: true},
		{name: "absolute path", inputs: map[string]string{"path": "/x"}, err: true},
		{name: "negative fetch-depth", inputs: map[string]string{"fetch-depth": "-1"}, err: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			workspace := t.TempDir()
			checkout := func(inputs map[string]string) error {
				return Main(t.Context(), &action.Context{
					Container:      &container.HostEnvironment{Path: workspace, StdOut: io.Discard},
					Github:         &model.GithubContext{Repository: "org/repo", Ref: "refs/heads/main", Sha: sha, ServerURL: server},
					Inputs:         inputs,
					HostWorkdir:    host,
					Workspace:      workspace,
					BindWorkdir:    tc.bind,
					NoSkipCheckout: tc.noSkip,
					Env:            map[string]string{"PATH": os.Getenv("PATH"), "GIT_CONFIG_GLOBAL": os.DevNull, "GIT_CONFIG_NOSYSTEM": "1"},
				})
			}
			if tc.before != nil {
				require.NoError(t, checkout(tc.before))
			}
			err := checkout(tc.inputs)
			if tc.err {
				require.Error(t, err)
				entries, _ := os.ReadDir(workspace)
				assert.Empty(t, entries)
				return
			}
			require.NoError(t, err)
			dst := filepath.Join(workspace, tc.dst)
			assert.Equal(t, tc.hostCopy, fileExists(filepath.Join(dst, "marker")))
			assert.Equal(t, tc.log, gitLog(dst))
			if tc.refs != nil {
				assert.Subset(t, strings.Fields(runGit(t, "-C", dst, "for-each-ref", "--format=%(refname)")), tc.refs)
			}
			if tc.log != nil {
				require.NoError(t, checkout(tc.inputs), "checking out into the same path again")
				assert.Equal(t, tc.log, gitLog(dst))
			}
		})
	}
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
