// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/common/git"
	"gitea.com/gitea/runner/act/model"
)

// Two lines of the @actions/cache client keep it from using this runner's cache service v2 API,
// and both are edited out of the action's own bundle.
//
// isGhes() decides that any host which is not github.com, *.ghe.com or *.localhost is GitHub
// Enterprise, and getCacheServiceVersion() forces v1 for it. The edit turns the last of those
// three tests, `hostname.endsWith('.LOCALHOST')`, into `endsWith("")`, which every hostname
// satisfies. That is the whole gate: one string literal, no call sites to resolve, and the same
// answer the toolkit's own proposed ACTIONS_VENDOR switch would give. Gitea already makes this
// edit by hand in its fork of upload-artifact.
//
// getCacheServiceURL() then resolves the service from ACTIONS_RESULTS_URL alone, where v1 reads
// ACTIONS_CACHE_URL first. Both reads in that function are given the same ACTIONS_CACHE_URL
// preference, which is what keeps the runner out of the artifact path: the results URL still
// points at Gitea and only the cache goes to the cache server. That function is found by its v1
// branch, the only place the two variables are read together.
//
// Neither edit is applied without the other. A bundle where one shape is unrecognised is left
// alone and keeps to v1, which the runner also serves, rather than being sent to a results URL
// that implements no cache service.
//
// Either of these landing upstream makes this file deletable:
//
//	https://github.com/actions/toolkit/pull/2123   — an ACTIONS_VENDOR switch, naming Gitea
//	https://github.com/actions/toolkit/issues/2439 — treat ACTIONS_RESULTS_URL as the signal
const (
	CacheServiceV2Env = "ACTIONS_CACHE_SERVICE_V2"
	cacheURLEnv       = "ACTIONS_CACHE_URL"
	resultsURLEnv     = "ACTIONS_RESULTS_URL"

	// artifactRefusal is the error @actions/artifact throws behind the same gate, and the only
	// thing that gate guards there. Its presence is what makes an artifact bundle safe to open:
	// a bundle carrying neither toolkit uses isGhes for something this runner has not looked at,
	// and is left alone.
	artifactRefusal = "GHESNotSupportedError"

	// urlBranchDistance separates the two branches of getCacheServiceURL: the v1 branch reads the
	// cache URL and falls back to the results URL, the v2 branch reads the results URL alone. Both
	// live in one function body, so the gap is small and stable, 2 bytes minified and 20 not,
	// across every bundle measured. The nearest read of the results URL belonging to anything else
	// sits some 75 kB away, so the window has room to spare for a larger function body.
	urlBranchDistance = 1 << 10

	// sidecarSuffix names the directory holding each patched bundle's untouched copy, a sibling
	// of the action directory because that directory is copied wholesale into job containers.
	sidecarSuffix = ".cache-gate"

	// skipMarker in the sidecar means a patched bundle already failed once here.
	skipMarker = "skip"

	maxBundleSize = 64 << 20
)

var (
	// localhostTest matches the third hostname test of isGhes, in any quoting. The match is case
	// sensitive on purpose, and that is load-bearing: isGhes uppercases the hostname before
	// testing it, while undici, bundled into all of these actions, tests a lowercase ".localhost"
	// in isURLPotentiallyTrustworthy. Opening that one would tell its HTTP client that every URL
	// is trustworthy. Uppercase, the literal occurs nowhere but this test, across sixteen bundles
	// from ten actions.
	localhostTest = regexp.MustCompile(`endsWith\s*\(\s*` + quoted(`\.LOCALHOST`) + `\s*\)`)

	// cacheURLRead and resultsURLRead match one read of either variable, in either property syntax.
	cacheURLRead   = regexp.MustCompile(envRead(cacheURLEnv))
	resultsURLRead = regexp.MustCompile(envRead(resultsURLEnv))
)

// envRead builds a pattern matching one read of an environment variable, as a property or as any
// flavour of quoted index.
func envRead(name string) string {
	return `process\s*\.\s*env\s*(?:\.\s*` + name + `\b|\[\s*` + quoted(name) + `\s*\])`
}

// quoted matches a JavaScript string literal holding exactly pattern, in any of the three quote
// characters. RE2 has no backreferences, so the pairs are spelled out.
func quoted(pattern string) string {
	return "(?:'" + pattern + "'|\"" + pattern + "\"|`" + pattern + "`)"
}

// actionScriptPaths returns the entrypoints of a node action, the only kind with a bundle. Only
// remote actions get here: a local one lives in the user's checkout, which the runner does not
// rewrite, so a repo-vendored action bundling @actions/cache stays on the v1 API.
func actionScriptPaths(dir string, action *model.Action) []string {
	if action == nil || !action.Runs.Using.IsNode() {
		return nil
	}
	var paths []string
	for _, script := range []string{action.Runs.Pre, action.Runs.Main, action.Runs.Post} {
		if script != "" {
			paths = append(paths, filepath.Join(dir, script))
		}
	}
	return paths
}

// patchCacheGate lets an action reach the v2 cache API by patching the version gate out of its
// bundles, keeping each original in a sidecar directory. Every failure is silent and leaves the
// bundle untouched: the client then picks v1, which the runner also serves.
func patchCacheGate(ctx context.Context, actionDir string, scripts []string) {
	if len(scripts) == 0 {
		return
	}
	if _, err := os.Stat(filepath.Join(sidecarDir(actionDir), skipMarker)); err == nil {
		return
	}
	defer git.AcquireCloneLock(actionDir)()

	for _, script := range scripts {
		if err := patchBundle(script, originalFor(actionDir, script)); err != nil {
			common.Logger(ctx).Debugf("cache service v2: %s left unpatched: %v", filepath.Base(script), err)
		}
	}
}

// revertCacheGate puts the originals back and stops this action being patched again, so the next
// job runs it exactly as shipped. Called when a step failed with a patched bundle; it does not
// re-run the step, because a step's outputs and env-file writes are already recorded by then.
func revertCacheGate(ctx context.Context, actionDir string, scripts []string) {
	if len(scripts) == 0 {
		return
	}
	if _, err := os.Stat(sidecarDir(actionDir)); err != nil {
		return
	}
	defer git.AcquireCloneLock(actionDir)()

	reverted := false
	for _, script := range scripts {
		original := originalFor(actionDir, script)
		if !isPatchOf(original, script) {
			continue
		}
		if err := os.Rename(original, script); err == nil {
			reverted = true
		}
	}
	_ = os.WriteFile(filepath.Join(sidecarDir(actionDir), skipMarker), nil, 0o600)
	if reverted {
		common.Logger(ctx).Warnf("cache service v2: restored the original %s, it will not be patched again", filepath.Base(actionDir))
	}
}

// sidecarDir holds an action's untouched bundles, and the marker that stops it being patched.
func sidecarDir(actionDir string) string {
	return actionDir + sidecarSuffix
}

func originalFor(actionDir, script string) string {
	return filepath.Join(sidecarDir(actionDir), strings.TrimPrefix(script, actionDir))
}

// patchBundle rewrites one entrypoint in place. The original is moved aside, not copied, and it
// is what marks the bundle as already patched.
func patchBundle(script, original string) error {
	if _, err := os.Stat(original); err == nil {
		if isPatchOf(original, script) {
			return nil
		}
		// The action's ref moved and git checked the new bundle out over the patched one, so
		// the pair no longer belongs together. Patch afresh rather than keep an original that
		// would restore an older version of the action.
		if err := os.Remove(original); err != nil {
			return err
		}
	}
	info, err := os.Stat(script)
	if err != nil {
		return err
	}
	if info.Size() > maxBundleSize {
		return nil
	}
	data, err := os.ReadFile(script)
	if err != nil {
		return err
	}
	patched, ok := patchedBundle(data)
	if !ok {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(original), 0o755); err != nil {
		return err
	}
	if err := os.Rename(script, original); err != nil {
		return err
	}
	return os.WriteFile(script, patched, info.Mode().Perm())
}

// isPatchOf reports whether script is exactly what patching original produced. It is what proves
// the two still belong together: an action whose ref moved is checked out over the patched bundle,
// leaving an original that would restore the version before the move.
func isPatchOf(original, script string) bool {
	data, err := os.ReadFile(original)
	if err != nil {
		return false
	}
	current, err := os.ReadFile(script)
	if err != nil {
		return false
	}
	patched, ok := patchedBundle(data)
	return ok && bytes.Equal(patched, current)
}

// edit is a byte range of the bundle and what replaces it.
type edit struct {
	start, end int
	text       string
}

// patchedBundle opens the GHES gate, and where the cache toolkit is present, points the cache
// service at the cache server. A bundle this runner cannot account for comes back untouched.
func patchedBundle(data []byte) ([]byte, bool) {
	gates := localhostTest.FindAllIndex(data, -1)
	if len(gates) == 0 {
		return data, false
	}
	urls := cacheURLPreference(data)

	switch {
	case bytes.Contains(data, []byte(CacheServiceV2Env)):
		// The cache toolkit: both edits or neither, because choosing v2 without redirecting the
		// URL would send the client to a results URL that serves no cache service.
		if len(urls) == 0 {
			return data, false
		}
	case bytes.Contains(data, []byte(artifactRefusal)):
		// The artifact toolkit, where the gate is a plain refusal and there is no URL to move:
		// artifacts already go to Gitea, which implements that service.
	default:
		return data, false
	}

	edits := make([]edit, 0, len(gates)+len(urls))
	for _, g := range gates {
		edits = append(edits, edit{g[0], g[1], `endsWith("")`})
	}
	return applyEdits(data, append(edits, urls...)), true
}

// cacheURLPreference rewrites the reads of the results URL that sit in getCacheServiceURL so they
// prefer the cache URL, which is where this runner serves the cache service. That function is
// found by its v1 branch, the only place the two variables are read together: the v2 branch is
// what follows. Giving both branches the same preference leaves v1 as it already behaves.
func cacheURLPreference(data []byte) []edit {
	var edits []edit
	for _, anchor := range cacheURLRead.FindAllIndex(data, -1) {
		hi := min(anchor[1]+urlBranchDistance, len(data))
		for _, m := range resultsURLRead.FindAllIndex(data[anchor[1]:hi], -1) {
			start, end := anchor[1]+m[0], anchor[1]+m[1]
			edits = append(edits, edit{start, end, "(process.env." + cacheURLEnv + "||" + string(data[start:end]) + ")"})
		}
	}
	return edits
}

// applyEdits splices the edits into data. They are sorted because a later edit may sit before an
// earlier one's anchor, and any that overlap a kept edit is dropped rather than applied twice.
func applyEdits(data []byte, edits []edit) []byte {
	slices.SortFunc(edits, func(a, b edit) int { return a.start - b.start })

	grown := len(data)
	for _, e := range edits {
		grown += len(e.text)
	}
	out := make([]byte, 0, grown)
	last := 0
	for _, e := range edits {
		if e.start < last {
			continue
		}
		out = append(out, data[last:e.start]...)
		out = append(out, e.text...)
		last = e.end
	}
	return append(out, data[last:]...)
}
