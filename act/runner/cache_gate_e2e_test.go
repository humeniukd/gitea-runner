// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"gitea.com/gitea/runner/act/artifactcache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// actionsCacheRef pins the actions/cache release this is verified against. Bump it
// deliberately: a new release is exactly what can stop the patch matching.
const actionsCacheRef = "v6.1.0"

// actionsCacheBundle downloads one of the pinned action's entrypoints, keeping it in the user
// cache dir so repeated runs cost nothing. The bundles are megabytes, too large to vendor.
func actionsCacheBundle(t *testing.T, stage string) string {
	t.Helper()

	cacheDir, err := os.UserCacheDir()
	require.NoError(t, err)
	dir := filepath.Join(cacheDir, "gitea-runner-test", "actions-cache-"+actionsCacheRef)
	bundle := filepath.Join(dir, stage+".js")
	if _, err := os.Stat(bundle); err == nil {
		return bundle
	}
	require.NoError(t, os.MkdirAll(dir, 0o755))

	url := "https://raw.githubusercontent.com/actions/cache/" + actionsCacheRef + "/dist/" + stage + "/index.js"
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("cannot reach %s: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Skipf("GET %s: %s", url, resp.Status)
	}
	file, err := os.Create(bundle)
	require.NoError(t, err)
	_, err = io.Copy(file, resp.Body)
	require.NoError(t, file.Close())
	require.NoError(t, err)
	return bundle
}

// runCacheAction runs one entrypoint the way a job would. The server URL is the real Gitea host:
// nothing about the environment is rewritten, so only the patch can make the client choose v2.
func runCacheAction(t *testing.T, script, workspace, runnerTemp, resultsURL, token, key string) string {
	t.Helper()

	state := filepath.Join(runnerTemp, "state")
	output := filepath.Join(runnerTemp, "output")
	for _, name := range []string{state, output} {
		require.NoError(t, os.WriteFile(name, nil, 0o600))
	}

	cmd := exec.CommandContext(t.Context(), "node", script)
	cmd.Dir = workspace
	cmd.Env = append(os.Environ(),
		"INPUT_PATH=to-cache",
		"INPUT_KEY="+key,
		"ACTIONS_RUNTIME_TOKEN="+token,
		"ACTIONS_RESULTS_URL="+resultsURL+"/",
		"ACTIONS_CACHE_SERVICE_V2=true",
		"GITHUB_SERVER_URL=https://gitea.example.com",
		"GITHUB_REF=refs/heads/main",
		"GITHUB_EVENT_NAME=push",
		"GITHUB_WORKSPACE="+workspace,
		"RUNNER_TEMP="+runnerTemp,
		"GITHUB_STATE="+state,
		"GITHUB_OUTPUT="+output,
	)
	out, err := cmd.CombinedOutput()
	t.Logf("%s:\n%s", filepath.Base(filepath.Dir(script)), out)
	require.NoError(t, err, "%s failed", script)
	return string(out)
}

// tempDirPath is TempDir with symlinks resolved, because macOS hands out /var paths that resolve
// to /private/var and the client derives archive paths relative to the workspace.
func tempDirPath(t *testing.T) string {
	t.Helper()

	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	return dir
}

// The whole chain against the pinned release: patch the real bundles, run the real client with an
// ordinary Gitea server URL, and have it save and restore through this runner's cache server. If
// a release stops matching the patch the client falls back to v1 and this fails on the version
// line, which is the signal to look at the new bundle.
func TestCacheServiceV2EndToEnd(t *testing.T) {
	requireHostTools(t, "node")

	// A stand-in action directory, patched exactly as a downloaded one would be.
	actionDir := tempDirPath(t)
	scripts := map[string]string{}
	for _, stage := range []string{"restore", "save"} {
		body, err := os.ReadFile(actionsCacheBundle(t, stage))
		require.NoError(t, err)
		scripts[stage] = filepath.Join(actionDir, stage+".js")
		require.NoError(t, os.WriteFile(scripts[stage], body, 0o600))
	}
	patchCacheGate(t.Context(), actionDir, []string{scripts["restore"], scripts["save"]})

	handler, err := artifactcache.StartHandler(filepath.Join(t.TempDir(), "cache"), "127.0.0.1", 0, "", nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = handler.Close() })
	const token, repo = "e2e-runtime-token", "testuser/testrepo"
	handler.RegisterJob(token, repo)

	workspace, runnerTemp := tempDirPath(t), tempDirPath(t)
	require.NoError(t, os.MkdirAll(filepath.Join(workspace, "to-cache"), 0o755))
	content := []byte("cached through the patched gate")
	require.NoError(t, os.WriteFile(filepath.Join(workspace, "to-cache", "data.txt"), content, 0o600))

	const key = "patched-gate-key"
	missed := runCacheAction(t, scripts["restore"], workspace, runnerTemp, handler.ExternalURL(), token, key)
	require.Contains(t, missed, "Cache service version: v2", "the patch did not take, the client stayed on v1")
	require.Contains(t, missed, "Cache not found for input keys: "+key)

	saved := runCacheAction(t, scripts["save"], workspace, runnerTemp, handler.ExternalURL(), token, key)
	require.Contains(t, saved, "Cache saved with key: "+key)

	restored := tempDirPath(t)
	hit := runCacheAction(t, scripts["restore"], restored, runnerTemp, handler.ExternalURL(), token, key)
	require.Contains(t, hit, "Cache restored from key: "+key)

	got, err := os.ReadFile(filepath.Join(restored, "to-cache", "data.txt"))
	require.NoError(t, err)
	assert.Equal(t, content, got)
}
