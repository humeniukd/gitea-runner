// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"bytes"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

// BenchmarkCacheLookup measures warm HTTP lookups, excluding archive transfers.
func BenchmarkCacheLookup(b *testing.B) {
	for _, api := range []string{"v1", "v2"} {
		for _, match := range []string{"exact", "prefix", "miss"} {
			for _, clients := range []int{1, 8, 32} {
				b.Run(fmt.Sprintf("%s/%s/clients=%d", api, match, clients), func(b *testing.B) {
					benchmarkCacheLookup(b, api, match, clients)
				})
			}
		}
	}
}

func benchmarkCacheLookup(b *testing.B, api, match string, clients int) {
	const keys = 50
	dir := b.TempDir()
	storage, err := NewStorage(filepath.Join(dir, "cache"))
	if err != nil {
		b.Fatal(err)
	}
	db, err := bolthold.Open(filepath.Join(dir, "bolt.db"), 0o600, &bolthold.Options{
		Encoder: func(v any) ([]byte, error) { return json.Marshal(v) },
		Decoder: func(data []byte, v any) error { return json.Unmarshal(data, v) },
	})
	if err != nil {
		b.Fatal(err)
	}
	now := time.Now().Unix()
	err = db.Bolt().Update(func(tx *bbolt.Tx) error {
		for repo := range 20 {
			for key := range keys {
				entry := &Cache{ID: uint64(repo*keys + key + 1), Repo: fmt.Sprintf("owner/repo-%d", repo), Key: fmt.Sprintf("key-%03d-build", key), Version: "v1", Complete: true, Size: 1, CreatedAt: now, UsedAt: now}
				if err := db.TxInsert(tx, entry.ID, entry); err != nil {
					return err
				}
				if err := storage.Write(entry.ID, 0, strings.NewReader("x")); err != nil {
					return err
				}
				if _, err := storage.Commit(entry.ID, 1); err != nil {
					return err
				}
			}
		}
		return nil
	})
	closeErr := db.Close()
	if err != nil {
		b.Fatal(err)
	}
	if closeErr != nil {
		b.Fatal(closeErr)
	}
	h, err := StartHandler(Options{Dir: dir, OutboundIP: "127.0.0.1"})
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if err := h.Close(); err != nil {
			b.Error(err)
		}
	})
	h.RegisterJob("benchmark-token", JobCredential{Repo: "owner/repo-0"})
	transport := &http.Transport{MaxIdleConns: clients, MaxIdleConnsPerHost: clients, MaxConnsPerHost: clients}
	b.Cleanup(transport.CloseIdleConnections)
	client := &http.Client{Transport: transport, Timeout: 15 * time.Second}
	request := func(index int) error {
		key := fmt.Sprintf("key-%03d-build", index%keys)
		query := key
		if match == "prefix" {
			query = fmt.Sprintf("key-%03d-", index%keys)
		}
		if match == "miss" {
			query = "missing"
		}
		endpoint := h.ExternalURL() + apiPath + "/cache?keys=" + query + "&version=v1"
		method := http.MethodGet
		var body io.Reader
		if api == "v2" {
			endpoint = h.ExternalURL() + cacheServiceV2Path + "/GetCacheEntryDownloadURL"
			method = http.MethodPost
			payload, err := json.Marshal(map[string]any{"key": query, "version": "v1"})
			if err != nil {
				return err
			}
			body = bytes.NewReader(payload)
		}
		req, err := http.NewRequestWithContext(b.Context(), method, endpoint, body)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer benchmark-token")
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		payload, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		expected := http.StatusOK
		if api == "v1" && match == "miss" {
			expected = http.StatusNoContent
		}
		if resp.StatusCode != expected {
			return fmt.Errorf("status %d: %s", resp.StatusCode, payload)
		}
		if expected == http.StatusNoContent {
			return nil
		}
		var result map[string]any
		if err := json.Unmarshal(payload, &result); err != nil {
			return err
		}
		if api == "v2" {
			if result["ok"] != (match != "miss") {
				return fmt.Errorf("unexpected hit status: %s", payload)
			}
			if match != "miss" && result["matched_key"] != key {
				return fmt.Errorf("unexpected key: %s", payload)
			}
		} else if result["cacheKey"] != key {
			return fmt.Errorf("unexpected key: %s", payload)
		}
		return nil
	}
	for i := range keys {
		if err := request(i); err != nil {
			b.Fatal(err)
		}
	}
	latencies := make([]time.Duration, b.N)
	var next atomic.Int64
	var failures atomic.Int64
	var wg sync.WaitGroup
	b.ResetTimer()
	for range clients {
		wg.Go(func() {
			for {
				index := int(next.Add(1) - 1)
				if index >= b.N {
					return
				}
				started := time.Now()
				if err := request(index); err != nil {
					if failures.Add(1) == 1 {
						b.Logf("first failed request: %v", err)
					}
				}
				latencies[index] = time.Since(started)
			}
		})
	}
	wg.Wait()
	b.StopTimer()
	slices.Sort(latencies)
	b.ReportMetric(float64(latencies[(len(latencies)-1)*95/100].Nanoseconds())/1e6, "p95-ms")
	b.ReportMetric(float64(failures.Load()), "errors")
}

// BenchmarkCacheSaveRestore measures concurrent reserve/upload/finalize/download cycles.
func BenchmarkCacheSaveRestore(b *testing.B) {
	const size = 1 << 20
	for _, api := range []string{"v1", "v2"} {
		for _, clients := range []int{1, 8} {
			b.Run(fmt.Sprintf("%s/size=1MiB/clients=%d", api, clients), func(b *testing.B) {
				benchmarkCacheSaveRestore(b, api, size, clients)
			})
		}
	}
}

func benchmarkCacheSaveRestore(b *testing.B, api string, size, clients int) {
	h, err := StartHandler(Options{Dir: b.TempDir(), OutboundIP: "127.0.0.1"})
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if err := h.Close(); err != nil {
			b.Error(err)
		}
	})
	h.RegisterJob("benchmark-token", JobCredential{Repo: "owner/repo-0"})
	transport := &http.Transport{MaxIdleConns: clients, MaxIdleConnsPerHost: clients, MaxConnsPerHost: clients}
	b.Cleanup(transport.CloseIdleConnections)
	client := &http.Client{Transport: transport, Timeout: 2 * time.Minute}
	archive := bytes.Repeat([]byte("a"), size)
	request := func(index int) error {
		key := fmt.Sprintf("k-%d", index)
		if api == "v2" {
			return saveRestoreV2(b, h, client, key, archive)
		}
		return saveRestoreV1(b, h, client, key, archive)
	}
	if err := request(0); err != nil {
		b.Fatal(err)
	}
	latencies := make([]time.Duration, b.N)
	var next atomic.Int64
	var failures atomic.Int64
	var wg sync.WaitGroup
	b.SetBytes(int64(size))
	b.ResetTimer()
	for range clients {
		wg.Go(func() {
			for {
				index := int(next.Add(1) - 1)
				if index >= b.N {
					return
				}
				started := time.Now()
				if err := request(index + 1); err != nil {
					if failures.Add(1) == 1 {
						b.Logf("first failed request: %v", err)
					}
				}
				latencies[index] = time.Since(started)
			}
		})
	}
	wg.Wait()
	b.StopTimer()
	slices.Sort(latencies)
	b.ReportMetric(float64(latencies[(len(latencies)-1)*95/100].Nanoseconds())/1e6, "p95-ms")
	b.ReportMetric(float64(failures.Load()), "errors")
}

func saveRestoreV1(b *testing.B, h *Handler, client *http.Client, key string, archive []byte) error {
	base := h.ExternalURL() + apiPath
	body, err := json.Marshal(&Request{Key: key, Version: "v1", Size: int64(len(archive))})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(b.Context(), http.MethodPost, base+"/caches", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer benchmark-token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	var reserved struct {
		CacheID uint64 `json:"cacheId"`
	}
	err = json.UnmarshalRead(resp.Body, &reserved)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("reserve status %d", resp.StatusCode)
	}
	req, err = http.NewRequestWithContext(b.Context(), http.MethodPatch, fmt.Sprintf("%s/caches/%d", base, reserved.CacheID), bytes.NewReader(archive))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer benchmark-token")
	req.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/*", len(archive)-1))
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload status %d", resp.StatusCode)
	}
	req, err = http.NewRequestWithContext(b.Context(), http.MethodPost, fmt.Sprintf("%s/caches/%d", base, reserved.CacheID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer benchmark-token")
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("commit status %d", resp.StatusCode)
	}
	req, err = http.NewRequestWithContext(b.Context(), http.MethodGet, base+"/cache?keys="+key+"&version=v1", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer benchmark-token")
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	var found struct {
		ArchiveLocation string `json:"archiveLocation"`
	}
	err = json.UnmarshalRead(resp.Body, &found)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK || found.ArchiveLocation == "" {
		return fmt.Errorf("lookup status %d", resp.StatusCode)
	}
	return downloadArchive(b, client, found.ArchiveLocation, len(archive))
}

func saveRestoreV2(b *testing.B, h *Handler, client *http.Client, key string, archive []byte) error {
	created, err := benchV2Call(b, h, client, "CreateCacheEntry", map[string]any{"key": key, "version": "v1"})
	if err != nil {
		return err
	}
	if created["ok"] != true {
		return fmt.Errorf("create not ok: %v", created)
	}
	uploadURL, _ := created["signed_upload_url"].(string)
	if err := benchPutBlob(b, client, uploadURL, archive); err != nil {
		return err
	}
	finalized, err := benchV2Call(b, h, client, "FinalizeCacheEntryUpload", map[string]any{
		"key": key, "version": "v1", "size_bytes": strconv.Itoa(len(archive)),
	})
	if err != nil {
		return err
	}
	if finalized["ok"] != true {
		return fmt.Errorf("finalize not ok: %v", finalized)
	}
	got, err := benchV2Call(b, h, client, "GetCacheEntryDownloadURL", map[string]any{"key": key, "version": "v1"})
	if err != nil {
		return err
	}
	downloadURL, _ := got["signed_download_url"].(string)
	if got["ok"] != true || downloadURL == "" {
		return fmt.Errorf("download lookup not ok: %v", got)
	}
	return downloadArchive(b, client, downloadURL, len(archive))
}

func benchV2Call(b *testing.B, h *Handler, client *http.Client, method string, request any) (map[string]any, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(b.Context(), http.MethodPost, h.ExternalURL()+cacheServiceV2Path+"/"+method, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer benchmark-token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s status %d: %s", method, resp.StatusCode, payload)
	}
	var result map[string]any
	if err := json.Unmarshal(payload, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func benchPutBlob(b *testing.B, client *http.Client, url string, archive []byte) error {
	req, err := http.NewRequestWithContext(b.Context(), http.MethodPut, url, bytes.NewReader(archive))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload status %d", resp.StatusCode)
	}
	return nil
}

func downloadArchive(b *testing.B, client *http.Client, url string, size int) error {
	req, err := http.NewRequestWithContext(b.Context(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download status %d", resp.StatusCode)
	}
	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}
	if n != int64(size) {
		return fmt.Errorf("downloaded %d bytes, want %d", n, size)
	}
	return nil
}
