// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
