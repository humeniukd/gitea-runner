// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

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
