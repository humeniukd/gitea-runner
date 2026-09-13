// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

func TestHandlerStartupReleasesStore(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, ".secret"), 0o700))
	_, err := StartHandler(Options{Dir: dir, OutboundIP: "127.0.0.1"})
	require.ErrorContains(t, err, "read cache secret")
	db, err := openDB(dir)
	require.NoError(t, err)
	require.NoError(t, db.Close())
}

func TestHandlerCloseReleasesStore(t *testing.T) {
	h := newTestHandler(t, Policy{})
	h.gcAt = time.Time{}
	require.NoError(t, h.Close())
	h.gcCache() // a sweep queued by the last request may start after shutdown
	require.NoError(t, h.Close())
	db, err := openDB(h.dir)
	require.NoError(t, err)
	require.NoError(t, db.Close())
}

func TestHandlerTouchCoalescesWrites(t *testing.T) {
	cache := &Cache{Repo: testRepo, Key: "touch", Version: "v1", Complete: true, UsedAt: time.Now().Add(-2 * touchStale).Unix()}
	h := newTestHandler(t, Policy{}, cache)
	require.NoError(t, h.touchCache(cache.ID, false))
	var touched Cache
	require.NoError(t, h.db.Get(cache.ID, &touched))
	require.Greater(t, touched.UsedAt, cache.UsedAt)
	txID := func() int {
		var id int
		require.NoError(t, h.db.Bolt().View(func(tx *bbolt.Tx) error { id = tx.ID(); return nil }))
		return id
	}
	before := txID()
	require.NoError(t, h.touchCache(cache.ID, false))
	assert.Equal(t, before, txID(), "fresh accesses must not start a write transaction")
	require.ErrorContains(t, h.touchCache(cache.ID, true), "already complete")
	assert.GreaterOrEqual(t, inUseGrace, artifactURLTTL+touchStale)
}

func TestCacheServiceV2ConcurrentReservations(t *testing.T) {
	h := newTestHandler(t, Policy{})
	const jobs = 32
	start := make(chan struct{})
	results := make(chan map[string]any, jobs)
	var wg sync.WaitGroup
	for i := range jobs {
		token := "job-" + strconv.Itoa(i)
		h.RegisterJob(token, JobCredential{Repo: testRepo})
		wg.Go(func() {
			<-start
			client := &http.Client{Transport: &bearerTransport{token: token}}
			results <- v2Call(t, h, client, "CreateCacheEntry", map[string]any{"key": "same-key", "version": "v1"})
		})
	}
	close(start)
	wg.Wait()
	close(results)
	winners := 0
	for result := range results {
		if result["ok"] == true {
			winners++
		}
	}
	require.Equal(t, 1, winners)
	var entries []Cache
	require.NoError(t, h.db.Find(&entries, nil))
	require.Len(t, entries, 1)
}

func TestCacheServiceV2ConcurrentFinalize(t *testing.T) {
	h := newTestHandler(t, Policy{})
	content := []byte("the cached archive")
	created := v2Call(t, h, testClient, "CreateCacheEntry", map[string]any{"key": "shared", "version": "v1"})
	require.Equal(t, true, created["ok"])
	require.Equal(t, http.StatusCreated, putBlob(t, created["signed_upload_url"].(string), content))

	const retries = 8
	start := make(chan struct{})
	results := make(chan map[string]any, retries)
	var wg sync.WaitGroup
	for range retries {
		wg.Go(func() {
			<-start
			results <- v2Call(t, h, testClient, "FinalizeCacheEntryUpload", map[string]any{
				"key": "shared", "version": "v1", "size_bytes": strconv.Itoa(len(content)),
			})
		})
	}
	close(start)
	wg.Wait()
	close(results)
	var entryID string
	for result := range results {
		require.Equal(t, true, result["ok"])
		if entryID == "" {
			entryID = result["entry_id"].(string)
		}
		assert.Equal(t, entryID, result["entry_id"])
	}
	var entries []Cache
	require.NoError(t, h.db.Find(&entries, nil))
	completed := 0
	for _, e := range entries {
		if e.Complete {
			completed++
		}
	}
	require.Equal(t, 1, completed)
}

func TestFindExactCacheIgnoresOtherRepositories(t *testing.T) {
	h := newTestHandler(t, Policy{})
	now := time.Now().Unix()
	for i := range 32 {
		require.NoError(t, insertCache(h.db, &Cache{
			Repo: fmt.Sprintf("owner/repo-%d", i), Key: "deps-node", Version: "v1",
			Complete: true, Size: 1, CreatedAt: now, UsedAt: now,
		}))
	}
	got, err := findExactCache(h.db, "owner/repo-7", "deps-node", "v1", true)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "owner/repo-7", got.Repo)
}

func TestLookupCacheDropsTruncatedBlob(t *testing.T) {
	cache := &Cache{Repo: testRepo, Key: "broken", Version: "v", Complete: true, Size: 4, UsedAt: time.Now().Unix()}
	h := newTestHandler(t, Policy{}, cache)
	require.NoError(t, h.storage.Write(cache.ID, 0, bytes.NewReader([]byte("x"))))
	_, err := h.storage.Commit(cache.ID, -1)
	require.NoError(t, err)

	resp, err := testClient.Get(fmt.Sprintf("%s%s/cache?keys=broken&version=v", h.ExternalURL(), apiPath))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.ErrorIs(t, h.db.Get(cache.ID, &Cache{}), bolthold.ErrNotFound)
}
