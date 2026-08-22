// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

// TestActionsCacheRoundTrip saves a cache entry in one job and restores it in a dependent job.
// It runs against both advertised cache service versions because the runner patches the actions
// toolkit to resolve the service URL differently per version (act/runner/patch_actions.go), and
// the toolkit's own resolution has changed across releases.
func TestActionsCacheRoundTrip(t *testing.T) {
	requireFixture(t)

	for _, tc := range []struct {
		name string
		v2   bool
	}{
		{name: "cache_v2", v2: true},
		{name: "cache_v1", v2: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Each subtest needs its own repo and its own runner: cache.v2 is runner config, so
			// it cannot be varied within a single harness.
			api, repo, _ := startScenario(t, "cache.yml", WithCacheV2(tc.v2))

			wfRun := waitForRun(t, api, repo)
			requireSuccess(t, api, repo, wfRun.ID)

			logs := runLogs(t, api, repo, wfRun.ID)
			if !strings.Contains(logs, "cache-restored-ok") {
				t.Fatalf("restore job did not confirm a cache hit: %s", logs)
			}
		})
	}
}
