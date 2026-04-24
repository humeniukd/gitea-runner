// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"gitea.com/gitea/act_runner/act/artifactcache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRunner_registerCacheForTask pins the runner → cache server contract
// that the PR relies on: while a task is running, cache requests carrying
// the task's ACTIONS_RUNTIME_TOKEN authenticate successfully; once the
// revoker returned by registerCacheForTask runs, the token stops working.
func TestRunner_registerCacheForTask(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "artifactcache")
	handler, err := artifactcache.StartHandler(dir, "127.0.0.1", 0, nil)
	require.NoError(t, err)
	defer handler.Close()

	r := &Runner{cacheHandler: handler}
	token := "run-token-123"
	unregister := r.registerCacheForTask(token, "owner/repo", "7")

	base := handler.ExternalURL() + "/_apis/artifactcache"
	probe := func() int {
		req, err := http.NewRequest(http.MethodGet, base+"/cache?keys=x&version=v", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		return resp.StatusCode
	}

	assert.NotEqual(t, http.StatusUnauthorized, probe(),
		"token should be accepted while task is registered")

	unregister()
	assert.Equal(t, http.StatusUnauthorized, probe(),
		"token must be rejected after the revoker runs")
}

// TestRunner_registerCacheForTask_NoOps covers the two paths where the
// runner should not call RegisterJob at all: no embedded cache server (the
// operator configured cfg.Cache.ExternalServer) and no runtime token (older
// Gitea server). Both must return a safe no-op revoker rather than panic.
func TestRunner_registerCacheForTask_NoOps(t *testing.T) {
	t.Run("nil cacheHandler", func(t *testing.T) {
		r := &Runner{}
		unregister := r.registerCacheForTask("tok", "owner/repo", "1")
		require.NotNil(t, unregister)
		unregister()
	})

	t.Run("empty token", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "artifactcache")
		handler, err := artifactcache.StartHandler(dir, "127.0.0.1", 0, nil)
		require.NoError(t, err)
		defer handler.Close()

		r := &Runner{cacheHandler: handler}
		unregister := r.registerCacheForTask("", "owner/repo", "1")
		require.NotNil(t, unregister)
		unregister()
	})
}

// TestRunner_CacheFullFlow_MatchesToolkit walks through the exact sequence
// @actions/cache performs against the cache server, to lock in toolkit
// compatibility:
//
//  1. reserve with Bearer ACTIONS_RUNTIME_TOKEN (POST /caches)
//  2. upload chunks with Bearer + Content-Range (PATCH /caches/:id)
//  3. commit with Bearer (POST /caches/:id)
//  4. find with Bearer and comma-separated keys (GET /cache)
//  5. download archiveLocation *without* Authorization (GET /artifacts/:id?...)
//
// If any of these behaviors changes and the toolkit stops interoperating,
// this test fails.
func TestRunner_CacheFullFlow_MatchesToolkit(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "artifactcache")
	handler, err := artifactcache.StartHandler(dir, "127.0.0.1", 0, nil)
	require.NoError(t, err)
	defer handler.Close()

	r := &Runner{cacheHandler: handler}
	token := "full-flow-token"
	unregister := r.registerCacheForTask(token, "owner/repo", "1")
	defer unregister()

	base := handler.ExternalURL() + "/_apis/artifactcache"
	do := func(method, url, contentType, contentRange, body string) *http.Response {
		req, err := http.NewRequest(method, url, strings.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		if contentRange != "" {
			req.Header.Set("Content-Range", contentRange)
		}
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		return resp
	}

	key := "toolkit-flow"
	version := "c19da02a2bd7e77277f1ac29ab45c09b7d46a4ee758284e26bb3045ad11d9d20"
	body := `hello-cache-body`

	// reserve
	resp := do(http.MethodPost, base+"/caches", "application/json", "",
		fmt.Sprintf(`{"key":"%s","version":"%s","cacheSize":%d}`, key, version, len(body)))
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var reserved struct {
		CacheID uint64 `json:"cacheId"`
	}
	require.NoError(t, decodeJSON(resp, &reserved))
	require.NotZero(t, reserved.CacheID)

	// upload
	resp = do(http.MethodPatch, fmt.Sprintf("%s/caches/%d", base, reserved.CacheID),
		"application/octet-stream", fmt.Sprintf("bytes 0-%d/*", len(body)-1), body)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// commit
	resp = do(http.MethodPost, fmt.Sprintf("%s/caches/%d", base, reserved.CacheID), "", "", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// find — @actions/cache always sends comma-separated keys here
	resp = do(http.MethodGet,
		fmt.Sprintf("%s/cache?keys=%s,fallback&version=%s", base, key, version), "", "", "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var hit struct {
		ArchiveLocation string `json:"archiveLocation"`
		CacheKey        string `json:"cacheKey"`
	}
	require.NoError(t, decodeJSON(resp, &hit))
	require.Equal(t, key, hit.CacheKey)
	require.NotEmpty(t, hit.ArchiveLocation)

	// download — toolkit does NOT attach Authorization here; the signature
	// in the URL must be enough.
	dl, err := http.Get(hit.ArchiveLocation)
	require.NoError(t, err)
	defer dl.Body.Close()
	require.Equal(t, http.StatusOK, dl.StatusCode)
	got := make([]byte, 64)
	n, _ := dl.Body.Read(got)
	assert.Equal(t, body, string(got[:n]))
}

func decodeJSON(resp *http.Response, v any) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
