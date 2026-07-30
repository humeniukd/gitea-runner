// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/common/git"
	"gitea.com/gitea/runner/act/model"
)

// The @actions/cache client picks its API version from GITHUB_SERVER_URL: getCacheServiceVersion
// returns v1 unless the host is github.com, *.ghe.com or *.localhost, and no server-side setting
// overrides that. Rather than lie to the whole action about the server URL, which breaks anything
// that builds a real URL out of it, the one call that gates the choice is replaced with false in
// the action's own bundle. Gitea already does the same edit by hand in its fork of
// upload-artifact; this applies it to the one call site, so other uses of isGhes keep working.
//
// Either of these landing upstream makes this file deletable:
//
//	https://github.com/actions/toolkit/pull/2123   — an ACTIONS_VENDOR switch, naming Gitea
//	https://github.com/actions/toolkit/issues/2439 — treat ACTIONS_RESULTS_URL as the signal
const (
	cacheServiceV2Env = "ACTIONS_CACHE_SERVICE_V2"
	gateName          = "isGhes"

	// cacheGateMarker must occur within maxMarkerDistance of the variable for the bundle to
	// hold the gate rather than an unrelated mention: they are adjacent statements of
	// getCacheServiceVersion and bundling keeps them so, measured 310-335 bytes apart across
	// actions/cache, the setup-* family, Swatinem/rust-cache and ruby/setup-ruby.
	cacheGateMarker   = ".LOCALHOST"
	maxMarkerDistance = 2 << 10

	// gateCallDistance is how far ahead of the variable the guarding call sits, in every shape
	// seen: an if statement, or the ternary esbuild emits.
	gateCallDistance = 400

	// sidecarSuffix names the directory holding each patched bundle's untouched copy. It is a
	// sibling of the action directory, never inside it, because the action directory is copied
	// wholesale into every job container.
	sidecarSuffix = ".cache-gate"

	// skipMarker in the sidecar means a patched bundle already failed once here, so the action
	// is left as its author shipped it from then on.
	skipMarker = "skip"

	maxBundleSize = 64 << 20
)

// gateAlias recovers the gate's real name after a bundler mangled it, as in `vu.isGhes=$K`.
var gateAlias = regexp.MustCompile(`\.` + gateName + `\s*=\s*([A-Za-z_$][A-Za-z0-9_$]*)`)

// actionScriptPaths returns the entrypoints of a node action, the only kind with a bundle.
func actionScriptPaths(dir string, action *model.Action) []string {
	if action == nil || !strings.HasPrefix(string(action.Runs.Using), "node") {
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

// patchCacheGateIn lets an action reach the v2 cache API by patching the version gate out of its
// bundles, keeping each original in a sidecar directory. Every failure is silent and leaves the
// bundle untouched: the client then picks v1, which the runner also serves.
func patchCacheGate(ctx context.Context, actionDir string, scripts []string) {
	if len(scripts) == 0 {
		return
	}
	sidecar := actionDir + sidecarSuffix
	if _, err := os.Stat(filepath.Join(sidecar, skipMarker)); err == nil {
		return
	}
	defer git.AcquireCloneLock(actionDir)()

	for _, script := range scripts {
		if err := patchBundle(script, originalFor(actionDir, script)); err != nil {
			common.Logger(ctx).Debugf("cache service v2: %s left unpatched: %v", filepath.Base(script), err)
		}
	}
}

// revertCacheGateIn puts the originals back and stops this action being patched again, so the
// next job runs it exactly as shipped. Called when a step failed with a patched bundle; it does
// not re-run the step, because a step's outputs and env-file writes are already recorded by then.
func revertCacheGate(ctx context.Context, actionDir string, scripts []string) {
	if len(scripts) == 0 {
		return
	}
	sidecar := actionDir + sidecarSuffix
	if _, err := os.Stat(sidecar); err != nil {
		return
	}
	defer git.AcquireCloneLock(actionDir)()

	reverted := false
	for _, script := range scripts {
		if err := os.Rename(originalFor(actionDir, script), script); err == nil {
			reverted = true
		}
	}
	if err := os.MkdirAll(sidecar, 0o755); err == nil {
		_ = os.WriteFile(filepath.Join(sidecar, skipMarker), nil, 0o600)
	}
	if reverted {
		common.Logger(ctx).Warnf("cache service v2: restored the original %s, it will not be patched again", filepath.Base(actionDir))
	}
}

func originalFor(actionDir, script string) string {
	rel, err := filepath.Rel(actionDir, script)
	if err != nil {
		rel = filepath.Base(script)
	}
	return filepath.Join(actionDir+sidecarSuffix, rel)
}

// patchBundle rewrites one entrypoint in place. The original is moved aside, not copied, and its
// presence is what marks the bundle as already patched.
func patchBundle(script, original string) error {
	if _, err := os.Stat(original); err == nil {
		return nil
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
	patched, count := patchedBundle(data)
	if count == 0 {
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

// patchCacheGate replaces the call guarding each read of the variable with false. A bundle that
// only mentions the variable comes back untouched.
func patchedBundle(data []byte) ([]byte, int) {
	if !hasVersionGate(data) {
		return data, 0
	}
	calls := gateCalls(data)

	out := make([]byte, 0, len(data))
	last, count := 0, 0
	marker := []byte(cacheServiceV2Env)
	for i := 0; ; {
		j := bytes.Index(data[i:], marker)
		if j < 0 {
			break
		}
		at := i + j
		i = at + len(marker)

		// The nearest call before the read is the guard; earlier ones include the gate's own
		// definition, which reads as an empty call too.
		k := len(calls) - 1
		for ; k >= 0 && calls[k][0] >= at; k-- {
		}
		if k < 0 || calls[k][0] < max(last, at-gateCallDistance) {
			continue
		}
		out = append(out, data[last:calls[k][0]]...)
		out = append(out, "false"...)
		last, count = calls[k][1], count+1
	}
	if count == 0 {
		return data, 0
	}
	return append(out, data[last:]...), count
}

// gateCalls returns the byte range of every empty call to the gate, under the plain name, any
// prefix a bundler added, and the short name it was mangled to. A definition is not a call.
func gateCalls(data []byte) [][2]int {
	names := []string{`[A-Za-z0-9_$]*` + gateName}
	for _, m := range gateAlias.FindAllSubmatch(data, -1) {
		names = append(names, regexp.QuoteMeta(string(m[1])))
	}
	call := regexp.MustCompile(`(?:\A|[^A-Za-z0-9_$])((?:` + strings.Join(names, "|") + `)\s*\(\s*\))`)

	var calls [][2]int
	for _, m := range call.FindAllSubmatchIndex(data, -1) {
		if bytes.HasSuffix(bytes.TrimRight(data[max(m[2]-10, 0):m[2]], " \t"), []byte("function")) {
			continue
		}
		calls = append(calls, [2]int{m[2], m[3]})
	}
	return calls
}

func hasVersionGate(data []byte) bool {
	marker, gate := []byte(cacheServiceV2Env), []byte(cacheGateMarker)
	for i := 0; ; {
		j := bytes.Index(data[i:], marker)
		if j < 0 {
			return false
		}
		at := i + j
		lo := max(at-maxMarkerDistance, 0)
		hi := min(at+maxMarkerDistance+len(gate), len(data))
		if bytes.Contains(data[lo:hi], gate) {
			return true
		}
		i = at + 1
	}
}
