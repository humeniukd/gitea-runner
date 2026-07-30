// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitea.com/gitea/runner/act/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The three shapes real bundlers emit, reduced to the bytes that matter. tsc keeps the name,
// webpack prefixes it, esbuild mangles it and writes a ternary, recording the real name in the
// export assignment.
const (
	gateTSC     = `function isGhes(){const h=u.hostname.toUpperCase();return h!=='GITHUB.COM'&&!h.endsWith('.GHE.COM')&&!h.endsWith('.LOCALHOST')}` + "\n" + `function getCacheServiceVersion() {` + "\n" + `    if (isGhes())` + "\n" + `        return 'v1';` + "\n" + `    return process.env['ACTIONS_CACHE_SERVICE_V2'] ? 'v2' : 'v1';` + "\n" + `}`
	gateWebpack = `function config_isGhes(){const h=u.hostname.toUpperCase();return h!=='GITHUB.COM'&&!h.endsWith('.GHE.COM')&&!h.endsWith('.LOCALHOST')}` + "\n" + `function config_getCacheServiceVersion() {` + "\n" + `    if (config_isGhes())` + "\n" + `        return 'v1';` + "\n" + `    return process.env['ACTIONS_CACHE_SERVICE_V2'] ? 'v2' : 'v1';` + "\n" + `}`
	gateEsbuild = `vu.isGhes=$K;vu.getCacheServiceVersion=XK;function $K(){let e=new URL(process.env.GITHUB_SERVER_URL||"https://github.com").hostname.toUpperCase(),r=e==="GITHUB.COM",n=e.endsWith(".GHE.COM"),i=e.endsWith(".LOCALHOST");return!r&&!n&&!i}function XK(){return $K()?"v1":process.env.ACTIONS_CACHE_SERVICE_V2?"v2":"v1"}`
)

func TestPatchedBundle(t *testing.T) {
	for _, tc := range []struct {
		name, body, wantContains string
		wantCount                int
	}{
		{"tsc keeps the name", gateTSC, "if (false)", 1},
		{"webpack prefixes it", gateWebpack, "if (false)", 1},
		{"esbuild mangles it into a ternary", gateEsbuild, `return false?"v1"`, 1},
		{"no gate at all", `console.log("hello")`, "", 0},
		// sccache-action sets the variable itself; there is nothing to patch.
		{"mentions the variable without the gate", `core.exportVariable("ACTIONS_CACHE_SERVICE_V2","on")`, "", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, n := patchedBundle([]byte(tc.body))
			assert.Equal(t, tc.wantCount, n)
			if tc.wantCount == 0 {
				assert.Equal(t, tc.body, string(out), "an unpatched bundle must come back byte for byte")
				return
			}
			assert.Contains(t, string(out), tc.wantContains)
			// Nothing but the guarding call may change.
			assert.Equal(t, strings.Count(tc.body, "ACTIONS_CACHE_SERVICE_V2"), strings.Count(string(out), "ACTIONS_CACHE_SERVICE_V2"))
			assert.Contains(t, string(out), ".LOCALHOST", "the hostname test itself is left alone")
		})
	}
}

func TestPatchBundleKeepsTheOriginal(t *testing.T) {
	dir, script := bundleFile(t, gateTSC)
	original := originalFor(dir, script)

	require.NoError(t, patchBundle(script, original))
	patched, err := os.ReadFile(script)
	require.NoError(t, err)
	assert.Contains(t, string(patched), "if (false)")

	kept, err := os.ReadFile(original)
	require.NoError(t, err)
	assert.Equal(t, gateTSC, string(kept), "the untouched bundle is kept outside the action tree")
	assert.NotContains(t, original, dir+string(filepath.Separator), "originals must not ship into job containers")

	// Patching again must not stack, and must not overwrite the kept original.
	require.NoError(t, patchBundle(script, original))
	again, err := os.ReadFile(script)
	require.NoError(t, err)
	assert.Equal(t, string(patched), string(again))
	kept, err = os.ReadFile(original)
	require.NoError(t, err)
	assert.Equal(t, gateTSC, string(kept))
}

// A bundle with nothing to patch is left exactly as it was, with no original kept beside it.
func TestPatchBundleLeavesOtherActionsAlone(t *testing.T) {
	dir, script := bundleFile(t, `console.log("checkout")`)
	original := originalFor(dir, script)

	require.NoError(t, patchBundle(script, original))
	body, err := os.ReadFile(script)
	require.NoError(t, err)
	assert.Equal(t, `console.log("checkout")`, string(body))
	_, err = os.Stat(original)
	assert.True(t, os.IsNotExist(err), "no original is kept for a bundle that was not patched")
}

// bundleFile writes one entrypoint into a fresh action directory.
func bundleFile(t *testing.T, body string) (dir, script string) {
	t.Helper()

	dir = t.TempDir()
	script = filepath.Join(dir, "index.js")
	require.NoError(t, os.WriteFile(script, []byte(body), 0o600))
	return dir, script
}

func TestActionScriptPaths(t *testing.T) {
	node := &model.Action{Runs: model.ActionRuns{Using: "node20", Main: "dist/restore/index.js", Post: "dist/save/index.js"}}
	assert.Equal(t, []string{"/a/dist/restore/index.js", "/a/dist/save/index.js"}, actionScriptPaths("/a", node))

	// Only a node action has a bundle to patch.
	assert.Nil(t, actionScriptPaths("/a", &model.Action{Runs: model.ActionRuns{Using: "docker", Image: "alpine"}}))
	assert.Nil(t, actionScriptPaths("/a", nil))
}

// A step that fails with a patched bundle gets the untouched bundle back, and the action is not
// patched again, so later jobs run it exactly as its author shipped it.
func TestRevertCacheGate(t *testing.T) {
	dir, script := bundleFile(t, gateTSC)
	scripts := []string{script}

	patchCacheGate(t.Context(), dir, scripts)
	body, err := os.ReadFile(script)
	require.NoError(t, err)
	require.Contains(t, string(body), "if (false)", "precondition: the bundle is patched")

	revertCacheGate(t.Context(), dir, scripts)
	body, err = os.ReadFile(script)
	require.NoError(t, err)
	assert.Equal(t, gateTSC, string(body), "the original bundle is back")

	// The skip marker survives, so the action stays unpatched from now on.
	patchCacheGate(t.Context(), dir, scripts)
	body, err = os.ReadFile(script)
	require.NoError(t, err)
	assert.Equal(t, gateTSC, string(body), "a reverted action stays unpatched")
}
