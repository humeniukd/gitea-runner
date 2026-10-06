// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checkout

import (
	"cmp"
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
	runGit(t, "-C", filepath.Join(server, "other", "repo"), "symbolic-ref", "refs/remotes/origin/HEAD", "refs/heads/main")
	sha256Repo := filepath.Join(server, "org", "sha256")
	runGit(t, "init", "--object-format=sha256", "--initial-branch=main", sha256Repo)
	runGit(t, "-C", sha256Repo, "commit", "--allow-empty", "-m", "s1")
	runGit(t, "-C", sha256Repo, "checkout", "--detach")
	host := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(host, "marker"), nil, 0o644))

	for _, tc := range []struct {
		name           string
		bind, noSkip   bool
		eventRef       string
		before, inputs map[string]string
		dst, head      string
		log, refs      []string
		hostCopy, err  bool
	}{
		{name: "bound workspace is skipped", bind: true},
		{name: "workflow ref copies host workdir", inputs: map[string]string{"ref": "main", "path": "a/../src"}, dst: "src", hostCopy: true},
		{name: "bound path checks out the event branch", bind: true, inputs: map[string]string{"path": "src"}, dst: "src", head: "main", log: []string{"c1"}},
		{name: "no-skip checks out the event branch", noSkip: true, head: "main", log: []string{"c1"}, refs: []string{"refs/remotes/origin/main"}},
		{name: "event tag is kept", noSkip: true, eventRef: "refs/tags/v1", log: []string{"c1"}, refs: []string{"refs/tags/v1"}},
		{name: "moved event tag checks out the event commit", noSkip: true, eventRef: "refs/tags/trunk", log: []string{"c1"}, refs: []string{"refs/tags/trunk"}},
		{
			name: "sha input removes previously persisted credentials without clean", noSkip: true, before: map[string]string{"ref": sha},
			inputs: map[string]string{"ref": sha, "persist-credentials": "false", "clean": "false"},
			log:    []string{"c1"},
		},
		{name: "ref clones shallow", inputs: map[string]string{"ref": "refs/heads/feature"}, head: "feature", log: []string{"f1"}},
		{name: "tag ref wins over same-named branch and is kept", inputs: map[string]string{"ref": "refs/tags/trunk"}, log: []string{"c2"}, refs: []string{"refs/tags/trunk"}},
		{name: "branch ref wins over same-named tag", inputs: map[string]string{"ref": "refs/heads/trunk"}, head: "trunk", log: []string{"t2"}},
		{
			name: "short branch wins over same-named tag and fetch-tags fetches every tag", inputs: map[string]string{"ref": "trunk", "fetch-tags": "true"},
			head: "trunk", log: []string{"t2"}, refs: []string{"refs/tags/trunk", "refs/tags/v1"},
		},
		{name: "short tag", inputs: map[string]string{"ref": "v1"}, log: []string{"c1"}, refs: []string{"refs/tags/v1"}},
		{
			name: "fetch-depth 0 deepens a shallow checkout to every branch and tag", before: map[string]string{"ref": "feature"},
			inputs: map[string]string{"ref": "refs/heads/feature", "fetch-depth": "0"}, head: "feature", log: []string{"f1", "c2", "c1"},
			refs: []string{"refs/remotes/origin/feature", "refs/remotes/origin/main", "refs/remotes/origin/trunk", "refs/tags/trunk", "refs/tags/v1"},
		},
		{name: "other repository defaults to its default branch", inputs: map[string]string{"repository": "other/repo"}, head: "trunk", log: []string{"t2"}},
		{name: "explicit HEAD checks out the default branch", inputs: map[string]string{"repository": "other/repo", "ref": "HEAD"}, head: "trunk", log: []string{"t2"}},
		{name: "detached sha256 HEAD is checked out", inputs: map[string]string{"repository": "org/sha256"}, log: []string{"s1"}},
		{name: "unknown ref", inputs: map[string]string{"ref": "missing"}, err: true},
		{name: "unsupported input", inputs: map[string]string{"submodules": "true"}, err: true},
		{name: "relative path escaping workspace", inputs: map[string]string{"path": "../x"}, err: true},
		{name: "absolute path", inputs: map[string]string{"path": "/x"}, err: true},
		{name: "negative fetch-depth", inputs: map[string]string{"fetch-depth": "-1"}, err: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			workspace, runnerTemp := t.TempDir(), t.TempDir()
			output := filepath.Join(runnerTemp, "output")
			checkout := func(inputs map[string]string) error {
				return Main(t.Context(), &action.Context{
					Container:      &container.HostEnvironment{Path: workspace, TmpDir: runnerTemp, StdOut: io.Discard},
					Github:         &model.GithubContext{Repository: "org/repo", Ref: cmp.Or(tc.eventRef, "refs/heads/main"), Sha: sha, ServerURL: server, Token: "tok"},
					Inputs:         inputs,
					HostWorkdir:    host,
					Workspace:      workspace,
					BindWorkdir:    tc.bind,
					NoSkipCheckout: tc.noSkip,
					Env: map[string]string{
						"PATH": os.Getenv("PATH"), "GIT_CONFIG_GLOBAL": os.DevNull, "GIT_CONFIG_NOSYSTEM": "1", "GIT_DIR": host, "GITHUB_OUTPUT": output,
					},
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
			outputs, err := os.ReadFile(output)
			require.NoError(t, err)
			dst := filepath.Join(workspace, tc.dst)
			assert.Equal(t, tc.hostCopy, fileExists(filepath.Join(dst, "marker")))
			assert.Equal(t, tc.log, gitLog(dst))
			if tc.refs != nil {
				assert.Subset(t, strings.Fields(runGit(t, "-C", dst, "for-each-ref", "--format=%(refname)")), tc.refs)
			}
			if tc.log == nil {
				assert.Contains(t, string(outputs), "commit="+sha+"\n")
				assert.Contains(t, string(outputs), "ref="+cmp.Or(tc.eventRef, "refs/heads/main")+"\n")
				return
			}
			if tc.head == "" {
				assert.Empty(t, gitOutput(dst, "symbolic-ref", "-q", "HEAD"))
			} else {
				assert.Equal(t, "refs/heads/"+tc.head, gitOutput(dst, "symbolic-ref", "-q", "HEAD"))
				assert.Equal(t, "origin/"+tc.head, gitOutput(dst, "rev-parse", "--abbrev-ref", "@{u}"))
			}
			persisted := tc.inputs["persist-credentials"] != "false"
			headers, _ := exec.Command("git", "-C", dst, "config", "--get-all", "http."+server+"/.extraheader").Output()
			assert.Equal(t, persisted, string(headers) == "\nAuthorization: Basic dG9rZW46dG9r\n")
			assert.Contains(t, string(outputs), "commit="+gitOutput(dst, "rev-parse", "HEAD")+"\n")
			if inputRef := tc.inputs["ref"]; inputRef != "" && inputRef != sha {
				assert.Contains(t, string(outputs), "ref="+inputRef+"\n")
			}

			require.NoError(t, os.WriteFile(filepath.Join(dst, "untracked"), nil, 0o644))
			require.NoError(t, checkout(tc.inputs), "checking out into the same path again")
			assert.Equal(t, tc.log, gitLog(dst))
			assert.Equal(t, tc.inputs["clean"] == "false", fileExists(filepath.Join(dst, "untracked")))
		})
	}
}

func gitOutput(dir string, args ...string) string {
	out, _ := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	return strings.TrimSpace(string(out))
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
