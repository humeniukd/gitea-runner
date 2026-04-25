// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"gitea.com/gitea/act_runner/act/artifactcache"
	"gitea.com/gitea/act_runner/internal/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyCfg returns a Config with no ExternalServer/ExternalSecret so the
// cache-for-task helper hits the embedded-handler path rather than the
// network-control-plane path.
func emptyCfg() *config.Config { return &config.Config{} }

// TestRunner_registerCacheForTask pins the runner → cache server contract
// that the PR relies on: while a task is running, cache requests carrying
// the task's ACTIONS_RUNTIME_TOKEN authenticate successfully; once the
// revoker returned by registerCacheForTask runs, the token stops working.
func TestRunner_registerCacheForTask(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "artifactcache")
	handler, err := artifactcache.StartHandler(dir, "127.0.0.1", 0, "", nil)
	require.NoError(t, err)
	defer handler.Close()

	r := &Runner{cfg: emptyCfg(), cacheHandler: handler}
	token := "run-token-123"
	unregister := r.registerCacheForTask(token, "owner/repo")

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
		r := &Runner{cfg: emptyCfg()}
		unregister := r.registerCacheForTask("tok", "owner/repo")
		require.NotNil(t, unregister)
		unregister()
	})

	t.Run("empty token", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "artifactcache")
		handler, err := artifactcache.StartHandler(dir, "127.0.0.1", 0, "", nil)
		require.NoError(t, err)
		defer handler.Close()

		r := &Runner{cfg: emptyCfg(), cacheHandler: handler}
		unregister := r.registerCacheForTask("", "owner/repo")
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
	handler, err := artifactcache.StartHandler(dir, "127.0.0.1", 0, "", nil)
	require.NoError(t, err)
	defer handler.Close()

	r := &Runner{cfg: emptyCfg(), cacheHandler: handler}
	token := "full-flow-token"
	unregister := r.registerCacheForTask(token, "owner/repo")
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

// TestRunner_ExternalCacheServer_RegisterRevoke is the high-level
// proof for cfg.Cache.ExternalServer with cfg.Cache.ExternalSecret set.
//
// Setup mirrors a production deployment: one process runs `act_runner
// cache-server` (here: a Handler with the same shared secret), another runs
// the runner daemon configured with external_server pointing at it. The
// runner has no embedded handler. Around each task, registerCacheForTask
// must POST to /_internal/register with the runtime token, then revoke on
// completion.
//
// What this test pins, end to end:
//
//   - Without the side-channel registration, the runtime token is unknown
//     to the cache server and management requests get 401 (proves auth is
//     active).
//   - After registerCacheForTask, the same token is accepted and the full
//     reserve→upload→commit→find→signed-download cycle works against the
//     remote server.
//   - After the deferred revoker runs, the token is rejected again (proves
//     revoke is wired and the server side honors it).
//
// If any of those break, this test fails.
func TestRunner_ExternalCacheServer_RegisterRevoke(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "remote-cache")
	const secret = "shared-secret-for-tests"
	remote, err := artifactcache.StartHandler(dir, "127.0.0.1", 0, secret, nil)
	require.NoError(t, err)
	defer remote.Close()

	r := &Runner{cfg: &config.Config{Cache: config.Cache{
		ExternalServer: remote.ExternalURL(),
		ExternalSecret: secret,
	}}}

	token := "external-task-token"
	repo := "owner/repoX"
	base := remote.ExternalURL() + "/_apis/artifactcache"
	probe := func() int {
		req, _ := http.NewRequest(http.MethodGet, base+"/cache?keys=k&version=v", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		return resp.StatusCode
	}

	require.Equal(t, http.StatusUnauthorized, probe(),
		"token must be unknown to the remote server before registration")

	unregister := r.registerCacheForTask(token, repo)
	require.NotEqual(t, http.StatusUnauthorized, probe(),
		"token must be accepted after registerCacheForTask")

	// Full reserve→upload→commit→find→download cycle, identical to what
	// @actions/cache does, against the remote (external) server.
	body := []byte("payload-from-task")
	reserveBody, _ := json.Marshal(&artifactcache.Request{Key: "ext-key", Version: "v", Size: int64(len(body))})
	req, _ := http.NewRequest(http.MethodPost, base+"/caches", bytes.NewReader(reserveBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var reserved struct {
		CacheID uint64 `json:"cacheId"`
	}
	require.NoError(t, decodeJSON(resp, &reserved))
	require.NotZero(t, reserved.CacheID)

	req, _ = http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/caches/%d", base, reserved.CacheID), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/*", len(body)-1))
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/caches/%d", base, reserved.CacheID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, base+"/cache?keys=ext-key&version=v", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var hit struct {
		ArchiveLocation string `json:"archiveLocation"`
	}
	require.NoError(t, decodeJSON(resp, &hit))
	require.NotEmpty(t, hit.ArchiveLocation)

	dl, err := http.Get(hit.ArchiveLocation)
	require.NoError(t, err)
	defer dl.Body.Close()
	require.Equal(t, http.StatusOK, dl.StatusCode)

	unregister()
	assert.Equal(t, http.StatusUnauthorized, probe(),
		"token must be rejected after the revoker runs")
}

// TestRunner_ExternalCacheServer_NoSecret pins the legacy compatibility path:
// when external_server is set without external_secret, the runner does not
// attempt registration. This guards against a regression where the runner
// would suddenly start firing register requests at a server that does not
// implement the control-plane (and would 404 them, polluting logs).
func TestRunner_ExternalCacheServer_NoSecret(t *testing.T) {
	registerCalls := 0
	srv := newTrackingExternalServer(&registerCalls)
	defer srv.Close()

	r := &Runner{cfg: &config.Config{Cache: config.Cache{
		ExternalServer: srv.URL,
		// ExternalSecret intentionally empty
	}}}

	unregister := r.registerCacheForTask("any-token", "owner/repo")
	unregister()
	assert.Zero(t, registerCalls,
		"runner must not call /_internal/register when external_secret is unset")
}

func newTrackingExternalServer(calls *int) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/_internal/", func(_ http.ResponseWriter, _ *http.Request) {
		*calls++
	})
	return httptest.NewServer(mux)
}
