// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunnerCloneEnvsReturnsTaskLocalCopy(t *testing.T) {
	r := &Runner{
		envs: map[string]string{
			"ACTIONS_CACHE_URL":   "http://cache.example",
			"ACTIONS_RUNTIME_URL": "http://runner.example",
		},
	}

	cloned := r.cloneEnvs()
	require.Equal(t, r.envs, cloned)

	cloned["ACTIONS_RUNTIME_TOKEN"] = "task-token"
	cloned["ACTIONS_ID_TOKEN_REQUEST_URL"] = "http://oidc.example"

	assert.NotContains(t, r.envs, "ACTIONS_RUNTIME_TOKEN")
	assert.NotContains(t, r.envs, "ACTIONS_ID_TOKEN_REQUEST_URL")
	assert.Equal(t, "http://cache.example", r.envs["ACTIONS_CACHE_URL"])
}
