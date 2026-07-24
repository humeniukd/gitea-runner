// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// v2Call posts a twirp request to the cache service and returns the decoded response.
func v2Call(t *testing.T, handler *Handler, method string, request any) map[string]any {
	t.Helper()

	body, err := json.Marshal(request)
	require.NoError(t, err)

	resp, err := testClient.Post(handler.ExternalURL()+cacheServiceV2Path+"/"+method, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	got := map[string]any{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	return got
}

func putBlob(t *testing.T, url string, content []byte) {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPut, url, bytes.NewReader(content))
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func startTestHandler(t *testing.T) *Handler {
	t.Helper()

	handler, err := StartHandler(filepath.Join(t.TempDir(), "artifactcache"), "", 0, "", nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = handler.Close() })
	handler.RegisterJob(testToken, testRepo)
	return handler
}

// The whole round trip an actions/cache@v4.2 or later client makes: reserve, upload the
// archive to the returned URL, finalize, then restore it.
func TestCacheServiceV2RoundTrip(t *testing.T) {
	handler := startTestHandler(t)
	content := []byte("the cached archive")

	created := v2Call(t, handler, "CreateCacheEntry", map[string]any{
		"key":     "deps-v1",
		"version": "abc123",
	})
	require.Equal(t, true, created["ok"])
	uploadURL, _ := created["signedUploadUrl"].(string)
	require.NotEmpty(t, uploadURL)

	putBlob(t, uploadURL, content)

	finalized := v2Call(t, handler, "FinalizeCacheEntryUpload", map[string]any{
		"key":     "deps-v1",
		"version": "abc123",
		// int64 fields travel as strings in the proto JSON mapping.
		"sizeBytes": strconv.Itoa(len(content)),
	})
	require.Equal(t, true, finalized["ok"])
	assert.NotEmpty(t, finalized["entryId"])

	got := v2Call(t, handler, "GetCacheEntryDownloadURL", map[string]any{
		"key":     "deps-v1",
		"version": "abc123",
	})
	require.Equal(t, true, got["ok"])
	assert.Equal(t, "deps-v1", got["matchedKey"])

	downloadURL, _ := got["signedDownloadUrl"].(string)
	require.NotEmpty(t, downloadURL)
	resp, err := http.Get(downloadURL) //nolint:noctx // the URL is the server under test
	require.NoError(t, err)
	defer resp.Body.Close()
	downloaded, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, content, downloaded)
}

// A large archive is staged as blocks and only put in order by the final block list, so
// blocks that arrive out of order must still be assembled the way the client asked.
func TestCacheServiceV2BlockUpload(t *testing.T) {
	handler := startTestHandler(t)

	created := v2Call(t, handler, "CreateCacheEntry", map[string]any{"key": "blocks", "version": "v1"})
	uploadURL, _ := created["signedUploadUrl"].(string)
	require.NotEmpty(t, uploadURL)

	blocks := map[string][]byte{
		base64.StdEncoding.EncodeToString([]byte("block-0")): []byte("hello "),
		base64.StdEncoding.EncodeToString([]byte("block-1")): []byte("world"),
		base64.StdEncoding.EncodeToString([]byte("block-2")): []byte("!"),
	}
	order := []string{
		base64.StdEncoding.EncodeToString([]byte("block-0")),
		base64.StdEncoding.EncodeToString([]byte("block-1")),
		base64.StdEncoding.EncodeToString([]byte("block-2")),
	}
	// Upload in an order that is not the block list order.
	for _, blockID := range []string{order[2], order[0], order[1]} {
		putBlob(t, uploadURL+"&comp=block&blockid="+blockID, blocks[blockID])
	}

	var list bytes.Buffer
	list.WriteString(`<?xml version="1.0" encoding="utf-8"?><BlockList>`)
	for _, blockID := range order {
		fmt.Fprintf(&list, "<Latest>%s</Latest>", blockID)
	}
	list.WriteString(`</BlockList>`)
	putBlob(t, uploadURL+"&comp=blocklist", list.Bytes())

	finalized := v2Call(t, handler, "FinalizeCacheEntryUpload", map[string]any{
		"key":       "blocks",
		"version":   "v1",
		"sizeBytes": len("hello world!"),
	})
	require.Equal(t, true, finalized["ok"])

	got := v2Call(t, handler, "GetCacheEntryDownloadURL", map[string]any{"key": "blocks", "version": "v1"})
	require.Equal(t, true, got["ok"])
	downloadURL, _ := got["signedDownloadUrl"].(string)
	resp, err := http.Get(downloadURL) //nolint:noctx // the URL is the server under test
	require.NoError(t, err)
	defer resp.Body.Close()
	downloaded, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "hello world!", string(downloaded))
}

func TestCacheServiceV2Misses(t *testing.T) {
	handler := startTestHandler(t)

	t.Run("reports a miss for an unknown key", func(t *testing.T) {
		got := v2Call(t, handler, "GetCacheEntryDownloadURL", map[string]any{"key": "nothing", "version": "v1"})
		assert.Equal(t, false, got["ok"])
	})

	t.Run("restore keys match by prefix", func(t *testing.T) {
		created := v2Call(t, handler, "CreateCacheEntry", map[string]any{"key": "deps-abc", "version": "v1"})
		putBlob(t, created["signedUploadUrl"].(string), []byte("x"))
		v2Call(t, handler, "FinalizeCacheEntryUpload", map[string]any{"key": "deps-abc", "version": "v1", "sizeBytes": 1})

		got := v2Call(t, handler, "GetCacheEntryDownloadURL", map[string]any{
			"key":         "deps-zzz",
			"restoreKeys": []string{"deps-"},
			"version":     "v1",
		})
		require.Equal(t, true, got["ok"])
		assert.Equal(t, "deps-abc", got["matchedKey"])
	})

	t.Run("an existing entry is not reserved twice", func(t *testing.T) {
		again := v2Call(t, handler, "CreateCacheEntry", map[string]any{"key": "deps-abc", "version": "v1"})
		assert.Equal(t, false, again["ok"])
	})

	t.Run("finalizing without a reservation is not ok", func(t *testing.T) {
		got := v2Call(t, handler, "FinalizeCacheEntryUpload", map[string]any{"key": "never-reserved", "version": "v1", "sizeBytes": 1})
		assert.Equal(t, false, got["ok"])
	})
}

// Reserving a key that is only a prefix of an existing completed key is a different entry
// and must succeed: the reservation is an exact (key, version) check, not a restore-key
// prefix match, or the shorter key would be reported as existing and silently never saved.
func TestCacheServiceV2ReservesPrefixOfExistingKey(t *testing.T) {
	handler := startTestHandler(t)

	created := v2Call(t, handler, "CreateCacheEntry", map[string]any{"key": "deps-abc", "version": "v1"})
	putBlob(t, created["signedUploadUrl"].(string), []byte("x"))
	v2Call(t, handler, "FinalizeCacheEntryUpload", map[string]any{"key": "deps-abc", "version": "v1", "sizeBytes": 1})

	reserved := v2Call(t, handler, "CreateCacheEntry", map[string]any{"key": "deps", "version": "v1"})
	require.Equal(t, true, reserved["ok"])
	require.NotEmpty(t, reserved["signedUploadUrl"])
}

// The cache of one repository must stay invisible to another, as it does for the v1 API.
func TestCacheServiceV2IsScopedPerRepository(t *testing.T) {
	handler := startTestHandler(t)

	created := v2Call(t, handler, "CreateCacheEntry", map[string]any{"key": "scoped", "version": "v1"})
	putBlob(t, created["signedUploadUrl"].(string), []byte("secret"))
	v2Call(t, handler, "FinalizeCacheEntryUpload", map[string]any{"key": "scoped", "version": "v1", "sizeBytes": 6})

	otherToken := "other-runtime-token"
	handler.RegisterJob(otherToken, "other/repo")
	otherClient := &http.Client{Transport: &bearerTransport{token: otherToken}}

	body, err := json.Marshal(map[string]any{"key": "scoped", "version": "v1"})
	require.NoError(t, err)
	resp, err := otherClient.Post(handler.ExternalURL()+cacheServiceV2Path+"/GetCacheEntryDownloadURL", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	got := map[string]any{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	assert.Equal(t, false, got["ok"])
}

func TestCacheServiceV2RejectsUnsignedBlobRequests(t *testing.T) {
	handler := startTestHandler(t)

	created := v2Call(t, handler, "CreateCacheEntry", map[string]any{"key": "signed", "version": "v1"})
	uploadURL, _ := created["signedUploadUrl"].(string)
	require.NotEmpty(t, uploadURL)

	t.Run("an unsigned upload is refused", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodPut,
			fmt.Sprintf("%s%s/1", handler.ExternalURL(), blobPath), bytes.NewReader([]byte("x")))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("an upload URL cannot be replayed as a download", func(t *testing.T) {
		resp, err := http.Get(uploadURL) //nolint:noctx // the URL is the server under test
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestTwirpRequestFieldNames(t *testing.T) {
	req := twirpRequest{
		"key":         json.RawMessage(`"k"`),
		"size_bytes":  json.RawMessage(`"4096"`),
		"restoreKeys": json.RawMessage(`["a","b"]`),
	}

	assert.Equal(t, "k", req.string("key"))
	// Either spelling, and a 64-bit integer sent as a string or as a number.
	assert.Equal(t, int64(4096), req.int64("sizeBytes", "size_bytes"))
	assert.Equal(t, []string{"a", "b"}, req.strings("restoreKeys", "restore_keys"))
	assert.Empty(t, req.string("missing"))
	assert.Zero(t, req.int64("missing"))
}
