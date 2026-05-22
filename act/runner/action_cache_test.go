// Copyright 2023 The Gitea Authors. All rights reserved.
// Copyright 2023 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	if dir != "" {
		args = append([]string{"-C", dir}, args...)
	}
	out, err := exec.Command("git", args...).CombinedOutput()
	require.NoError(t, err, string(out))
}

func TestActionCache(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()

	// Build a local bare repo with a `js` action dir so this runs offline (formerly cloned
	// github.com/nektos/act-test-actions over the network). allowAnySHA1InWant lets the
	// "Fetch Sha" case fetch a commit hash directly.
	remoteDir := t.TempDir()
	runGit(t, "", "init", "--bare", "--initial-branch=main", remoteDir)
	runGit(t, remoteDir, "config", "uploadpack.allowAnySHA1InWant", "true")

	workDir := t.TempDir()
	runGit(t, "", "clone", remoteDir, workDir)
	runGit(t, workDir, "config", "user.email", "test@example.com")
	runGit(t, workDir, "config", "user.name", "test")
	require.NoError(t, os.MkdirAll(filepath.Join(workDir, "js"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(workDir, "js", "action.yml"),
		[]byte("name: js\nruns:\n  using: node24\n  main: index.js\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(workDir, "js", "index.js"),
		[]byte("console.log('hello');\n"), 0o644))
	runGit(t, workDir, "add", ".")
	runGit(t, workDir, "commit", "-m", "initial")
	runGit(t, workDir, "push", "-u", "origin", "main")

	out, err := exec.Command("git", "-C", workDir, "rev-parse", "main").Output()
	require.NoError(t, err)
	fullSha := strings.TrimSpace(string(out))

	cache := &GoGitActionCache{
		Path: t.TempDir(),
	}
	cacheDir := "local/act-test-actions"
	refs := []struct {
		Name string
		Ref  string
	}{
		{Name: "Fetch Branch Name", Ref: "main"},
		{Name: "Fetch Branch Name Absolutely", Ref: "refs/heads/main"},
		{Name: "Fetch HEAD", Ref: "HEAD"},
		{Name: "Fetch Sha", Ref: fullSha},
	}
	for _, c := range refs {
		t.Run(c.Name, func(t *testing.T) {
			sha, err := cache.Fetch(ctx, cacheDir, remoteDir, c.Ref, "")
			if !a.NoError(err) || !a.NotEmpty(sha) { //nolint:testifylint // pre-existing issue from nektos/act
				return
			}
			atar, err := cache.GetTarArchive(ctx, cacheDir, sha, "js")
			// NotNil, not NotEmpty: atar is a live io.PipeReader whose producer goroutine is
			// writing concurrently; NotEmpty deep-reflects over its internals and races.
			if !a.NoError(err) || !a.NotNil(atar) { //nolint:testifylint // pre-existing issue from nektos/act
				return
			}
			// GetTarArchive streams from a background goroutine walking the shared repo.
			// Drain and close so it finishes before the next subtest fetches into the same
			// repo; otherwise the lingering walk races with that fetch.
			defer func() {
				_, _ = io.Copy(io.Discard, atar)
				_ = atar.Close()
			}()
			mytar := tar.NewReader(atar)
			th, err := mytar.Next()
			if !a.NoError(err) || !a.NotEqual(0, th.Size) { //nolint:testifylint // pre-existing issue from nektos/act
				return
			}
			buf := &bytes.Buffer{}
			// G110: Potential DoS vulnerability via decompression bomb (gosec)
			_, err = io.Copy(buf, mytar)
			a.NoError(err) //nolint:testifylint // pre-existing issue from nektos/act
			str := buf.String()
			a.NotEmpty(str)
		})
	}
}
