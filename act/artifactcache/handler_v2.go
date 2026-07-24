// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/timshannon/bolthold"
)

// The cache service v2 API. From actions/cache@v4.2 on, the toolkit talks twirp to
// `github.actions.results.api.v1.CacheService` instead of the /_apis/artifactcache
// endpoints, and uploads the archive to the returned URL with the Azure blob protocol.
// Both API versions are served from the same store, so a repository keeps its cache
// when a workflow moves between action versions.
const (
	cacheServiceV2Path = "/twirp/github.actions.results.api.v1.CacheService"

	// blobPath is where the URLs handed out by the v2 API point. Like the v1 artifact
	// route it authenticates by signature, because the client uploads and downloads
	// without an Authorization header.
	blobPath = apiPath + "/blobs"

	// blobUploadURLTTL bounds how long a signed upload URL stays valid. Uploading a
	// large cache takes a while, so it is more generous than the download TTL.
	blobUploadURLTTL = time.Hour
)

// registerV2Routes wires the cache service v2 API onto the router.
func (h *Handler) registerV2Routes(router *httprouter.Router) {
	router.POST(cacheServiceV2Path+"/CreateCacheEntry", h.bearerAuth(h.v2CreateCacheEntry))
	router.POST(cacheServiceV2Path+"/FinalizeCacheEntryUpload", h.bearerAuth(h.v2FinalizeCacheEntryUpload))
	router.POST(cacheServiceV2Path+"/GetCacheEntryDownloadURL", h.bearerAuth(h.v2GetCacheEntryDownloadURL))
	router.PUT(blobPath+"/:id", h.signedUploadAuth(h.v2UploadBlob))
	router.GET(blobPath+"/:id", h.signedURLAuth(h.get))
}

// POST /twirp/github.actions.results.api.v1.CacheService/CreateCacheEntry
//
// Reserves an entry and returns the URL to upload the archive to. An entry that already
// exists is reported as not ok, which is how the client learns to skip the upload.
func (h *Handler) v2CreateCacheEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cred := credFromContext(r.Context())
	req, err := decodeTwirpRequest(r)
	if err != nil {
		h.twirpError(w, r, http.StatusBadRequest, "malformed_request", err)
		return
	}
	key := req.string("key")
	version := req.string("version")
	if key == "" || version == "" {
		h.twirpError(w, r, http.StatusBadRequest, "invalid_argument", errors.New("key and version are required"))
		return
	}

	db, err := h.openDB()
	if err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	defer db.Close()

	// An exact (key, version) match means the entry is already cached; the client then skips
	// the upload. A prefix match must not count here, or a shorter key would be reported as
	// existing and silently never saved.
	if existing, err := findExactCache(db, cred.Repo, key, version); err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	} else if existing != nil {
		h.responseJSON(w, r, http.StatusOK, map[string]any{"ok": false})
		return
	}

	now := time.Now().Unix()
	cache := &Cache{
		Repo:      cred.Repo,
		Key:       key,
		Version:   version,
		Size:      -1, // the size is only known at finalize time
		CreatedAt: now,
		UsedAt:    now,
	}
	if err := insertCache(db, cache); err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}

	h.responseJSON(w, r, http.StatusOK, map[string]any{
		"ok":              true,
		"signedUploadUrl": h.signedUploadURL(cache.ID, time.Now().Add(blobUploadURLTTL)),
		// Emitted under both spellings so a client using the proto field names finds it.
		"signed_upload_url": h.signedUploadURL(cache.ID, time.Now().Add(blobUploadURLTTL)),
	})
}

// POST /twirp/github.actions.results.api.v1.CacheService/FinalizeCacheEntryUpload
//
// Completes the entry once the archive is uploaded, which is what makes it visible to
// GetCacheEntryDownloadURL.
func (h *Handler) v2FinalizeCacheEntryUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cred := credFromContext(r.Context())
	req, err := decodeTwirpRequest(r)
	if err != nil {
		h.twirpError(w, r, http.StatusBadRequest, "malformed_request", err)
		return
	}
	key := req.string("key")
	version := req.string("version")
	size := req.int64("sizeBytes", "size_bytes")

	db, err := h.openDB()
	if err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	defer db.Close()

	cache, err := findReservedCache(db, cred.Repo, key, version)
	if errors.Is(err, bolthold.ErrNotFound) {
		h.responseJSON(w, r, http.StatusOK, map[string]any{"ok": false})
		return
	} else if err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	db.Close()

	written, err := h.storage.Commit(cache.ID, size)
	if err != nil {
		h.logger.Errorf("finalize cache %d (%s): %v", cache.ID, cache.Key, err)
		h.responseJSON(w, r, http.StatusOK, map[string]any{"ok": false})
		return
	}
	cache.Size = written
	cache.Complete = true

	db, err = h.openDB()
	if err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	defer db.Close()
	if err := db.Update(cache.ID, cache); err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}

	h.responseJSON(w, r, http.StatusOK, map[string]any{
		"ok": true,
		// int64 fields travel as strings in the proto JSON mapping.
		"entryId":  strconv.FormatUint(cache.ID, 10),
		"entry_id": strconv.FormatUint(cache.ID, 10),
	})
}

// POST /twirp/github.actions.results.api.v1.CacheService/GetCacheEntryDownloadURL
func (h *Handler) v2GetCacheEntryDownloadURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cred := credFromContext(r.Context())
	req, err := decodeTwirpRequest(r)
	if err != nil {
		h.twirpError(w, r, http.StatusBadRequest, "malformed_request", err)
		return
	}
	keys := append([]string{req.string("key")}, req.strings("restoreKeys", "restore_keys")...)
	version := req.string("version")

	db, err := h.openDB()
	if err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	defer db.Close()

	cache, err := findCache(db, cred.Repo, keys, version)
	if err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	}
	if cache == nil {
		h.responseJSON(w, r, http.StatusOK, map[string]any{"ok": false})
		return
	}
	if ok, err := h.storage.Exist(cache.ID); err != nil {
		h.twirpError(w, r, http.StatusInternalServerError, "internal", err)
		return
	} else if !ok {
		// The blob is gone (pruned or lost); drop the dangling entry and report a miss.
		_ = db.Delete(cache.ID, cache)
		h.responseJSON(w, r, http.StatusOK, map[string]any{"ok": false})
		return
	}

	downloadURL := h.signedBlobURL(cache.ID, time.Now().Add(artifactURLTTL))
	h.responseJSON(w, r, http.StatusOK, map[string]any{
		"ok":                  true,
		"signedDownloadUrl":   downloadURL,
		"signed_download_url": downloadURL,
		"matchedKey":          cache.Key,
		"matched_key":         cache.Key,
	})
}

// blockList is the body of an Azure "put block list" request.
type blockList struct {
	Latest []string `xml:"Latest"`
}

// PUT /_apis/artifactcache/blobs/:id
//
// The archive arrives over the subset of the Azure blob API the toolkit uses: a small
// cache is a single PUT, a large one is staged as blocks that a final block list puts
// in order.
func (h *Handler) v2UploadBlob(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.ParseUint(params.ByName("id"), 10, 64)
	if err != nil {
		h.responseJSON(w, r, http.StatusBadRequest, err)
		return
	}

	switch strings.ToLower(r.URL.Query().Get("comp")) {
	case "block":
		blockID := r.URL.Query().Get("blockid")
		if blockID == "" {
			h.responseJSON(w, r, http.StatusBadRequest, errors.New("missing blockid"))
			return
		}
		if err := h.storage.WriteBlock(id, blockID, r.Body); err != nil {
			h.responseJSON(w, r, http.StatusInternalServerError, err)
			return
		}
	case "blocklist":
		body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
		if err != nil {
			h.responseJSON(w, r, http.StatusBadRequest, err)
			return
		}
		var list blockList
		if err := xml.Unmarshal(body, &list); err != nil {
			h.responseJSON(w, r, http.StatusBadRequest, fmt.Errorf("malformed block list: %w", err))
			return
		}
		if err := h.storage.OrderBlocks(id, list.Latest); err != nil {
			h.responseJSON(w, r, http.StatusInternalServerError, err)
			return
		}
	case "appendblock":
		if err := h.storage.AppendBlock(id, r.Body); err != nil {
			h.responseJSON(w, r, http.StatusInternalServerError, err)
			return
		}
	default:
		// A single-shot block blob upload: the whole archive in one request.
		if err := h.storage.Write(id, 0, r.Body); err != nil {
			h.responseJSON(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	h.useCache(int64(id))
	w.WriteHeader(http.StatusCreated)
}

// findReservedCache returns the entry a finalize call refers to: the newest reservation
// for this key and version that is not complete yet. A missing reservation is reported as
// bolthold.ErrNotFound.
func findReservedCache(db *bolthold.Store, repo, key, version string) (*Cache, error) {
	cache := &Cache{}
	if err := db.FindOne(cache, bolthold.Where("Repo").Eq(repo).
		And("Key").Eq(key).
		And("Version").Eq(version).
		And("Complete").Eq(false).
		SortBy("CreatedAt").Reverse()); err != nil {
		return nil, err
	}
	return cache, nil
}

func (h *Handler) signedUploadURL(cacheID uint64, exp time.Time) string {
	return h.signedBlobURLWith(cacheID, exp, h.computeUploadSignature)
}

func (h *Handler) signedBlobURL(cacheID uint64, exp time.Time) string {
	return h.signedBlobURLWith(cacheID, exp, h.computeSignature)
}

func (h *Handler) signedBlobURLWith(cacheID uint64, exp time.Time, sign func(cacheID, exp int64) string) string {
	expUnix := exp.Unix()
	q := url.Values{}
	q.Set("exp", strconv.FormatInt(expUnix, 10))
	q.Set("sig", sign(int64(cacheID), expUnix))
	return fmt.Sprintf("%s%s/%d?%s", h.ExternalURL(), blobPath, cacheID, q.Encode())
}

// twirpError reports a failure in the shape a twirp client expects, so the toolkit
// surfaces the message instead of a parse error.
func (h *Handler) twirpError(w http.ResponseWriter, r *http.Request, status int, code string, err error) {
	h.logger.Debugf("%s %s: %v", r.Method, r.URL.Path, err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code": code,
		"msg":  err.Error(),
	})
}

// twirpRequest is a decoded twirp JSON body. Field names are read leniently because a
// client may send either the proto field names or their camelCase form, and the proto
// JSON mapping writes 64-bit integers as strings.
type twirpRequest map[string]json.RawMessage

func decodeTwirpRequest(r *http.Request) (twirpRequest, error) {
	req := twirpRequest{}
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func (t twirpRequest) raw(names ...string) json.RawMessage {
	for _, name := range names {
		if raw, ok := t[name]; ok {
			return raw
		}
	}
	return nil
}

func (t twirpRequest) string(names ...string) string {
	var value string
	if raw := t.raw(names...); raw != nil {
		_ = json.Unmarshal(raw, &value)
	}
	return value
}

func (t twirpRequest) strings(names ...string) []string {
	var value []string
	if raw := t.raw(names...); raw != nil {
		_ = json.Unmarshal(raw, &value)
	}
	return value
}

func (t twirpRequest) int64(names ...string) int64 {
	raw := t.raw(names...)
	if raw == nil {
		return 0
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		if value, err := number.Int64(); err == nil {
			return value
		}
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		if value, err := strconv.ParseInt(text, 10, 64); err == nil {
			return value
		}
	}
	return 0
}
