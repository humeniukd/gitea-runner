// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/timshannon/bolthold"
)

// The cache service v2 API. A client on this version talks twirp to
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

	// blobUploadPurpose keeps an upload URL from being replayed to read an entry.
	blobUploadPurpose = "upload:"

	// blobUploadURLTTL bounds how long a signed upload URL stays valid. Uploading a
	// large cache takes a while, so it is more generous than the download TTL.
	blobUploadURLTTL = time.Hour

	// twirpInternal is the only error code that is not the client's fault.
	twirpInternal = "internal"
)

func (h *Handler) registerV2Routes(router *httprouter.Router) {
	router.POST(cacheServiceV2Path+"/CreateCacheEntry", h.bearerAuth(h.v2CreateCacheEntry))
	router.POST(cacheServiceV2Path+"/FinalizeCacheEntryUpload", h.bearerAuth(h.v2FinalizeCacheEntryUpload))
	router.POST(cacheServiceV2Path+"/GetCacheEntryDownloadURL", h.bearerAuth(h.v2GetCacheEntryDownloadURL))
	router.PUT(blobPath+"/:id", h.signedAuth(blobUploadPurpose, h.v2UploadBlob))
	router.GET(blobPath+"/:id", h.signedAuth("", h.get))
}

// Reserves an entry and returns the URL to upload the archive to. An entry that already
// exists is reported as not ok, which is how the client learns to skip the upload.
func (h *Handler) v2CreateCacheEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cred := credFromContext(r.Context())
	req, err := decodeTwirpRequest[v2CreateRequest](r)
	if err != nil {
		h.twirpError(w, r, "malformed_request", err)
		return
	}
	if req.Key == "" || req.Version == "" {
		h.twirpError(w, r, "invalid_argument", errors.New("key and version are required"))
		return
	}

	db, err := h.openDB()
	if err != nil {
		h.twirpError(w, r, twirpInternal, err)
		return
	}
	defer db.Close()

	// An exact (key, version) match means the entry is already cached; the client then skips
	// the upload. A prefix match must not count here, or a shorter key would be reported as
	// existing and silently never saved.
	if existing, err := findExactCache(db, cred.Repo, req.Key, req.Version); err != nil {
		h.twirpError(w, r, twirpInternal, err)
		return
	} else if existing != nil {
		h.twirpMiss(w, r)
		return
	}

	now := time.Now().Unix()
	cache := &Cache{
		Repo:      cred.Repo,
		Key:       req.Key,
		Version:   req.Version,
		Size:      -1, // the size is only known at finalize time
		CreatedAt: now,
		UsedAt:    now,
	}
	if err := insertCache(db, cache); err != nil {
		h.twirpError(w, r, twirpInternal, err)
		return
	}

	h.responseJSON(w, r, http.StatusOK, map[string]any{
		"ok":              true,
		"signedUploadUrl": h.signedURL(blobPath, blobUploadPurpose, cache.ID, time.Now().Add(blobUploadURLTTL)),
	})
}

// Completes the entry once the archive is uploaded, which is what makes it visible to
// GetCacheEntryDownloadURL.
func (h *Handler) v2FinalizeCacheEntryUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cred := credFromContext(r.Context())
	req, err := decodeTwirpRequest[v2FinalizeRequest](r)
	if err != nil {
		h.twirpError(w, r, "malformed_request", err)
		return
	}

	db, err := h.openDB()
	if err != nil {
		h.twirpError(w, r, twirpInternal, err)
		return
	}
	defer db.Close()

	cache, err := findReservedCache(db, cred.Repo, req.Key, req.Version)
	if errors.Is(err, bolthold.ErrNotFound) {
		h.twirpMiss(w, r)
		return
	} else if err != nil {
		h.twirpError(w, r, twirpInternal, err)
		return
	}
	db.Close() // commitCache needs the store closed

	size, _ := cmp.Or(req.SizeBytes, req.SizeBytesCamel).Int64()
	if err := h.commitCache(cache, size); err != nil {
		h.logger.Errorf("finalize cache %d (%s): %v", cache.ID, cache.Key, err)
		h.twirpMiss(w, r)
		return
	}

	h.responseJSON(w, r, http.StatusOK, map[string]any{
		"ok": true,
		// int64 fields travel as strings in the proto JSON mapping.
		"entryId": strconv.FormatUint(cache.ID, 10),
	})
}

func (h *Handler) v2GetCacheEntryDownloadURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cred := credFromContext(r.Context())
	req, err := decodeTwirpRequest[v2DownloadRequest](r)
	if err != nil {
		h.twirpError(w, r, "malformed_request", err)
		return
	}

	db, err := h.openDB()
	if err != nil {
		h.twirpError(w, r, twirpInternal, err)
		return
	}
	defer db.Close()

	cache, err := findCache(db, cred.Repo, req.keys(), req.Version)
	if err != nil {
		h.twirpError(w, r, twirpInternal, err)
		return
	}
	if cache == nil {
		h.twirpMiss(w, r)
		return
	}
	if ok, err := h.storage.Exist(cache.ID); err != nil {
		h.twirpError(w, r, twirpInternal, err)
		return
	} else if !ok {
		// The blob is gone (pruned or lost); drop the dangling entry and report a miss.
		_ = db.Delete(cache.ID, cache)
		h.twirpMiss(w, r)
		return
	}

	h.responseJSON(w, r, http.StatusOK, map[string]any{
		"ok":                true,
		"signedDownloadUrl": h.signedURL(blobPath, "", cache.ID, time.Now().Add(artifactURLTTL)),
		"matchedKey":        cache.Key,
	})
}

// blockList is the body of an Azure "put block list" request.
type blockList struct {
	Latest []string
}

// The archive arrives over the subset of the Azure blob API the toolkit uses: a small
// cache is a single PUT, a large one is staged as blocks that a final block list puts
// in order.
func (h *Handler) v2UploadBlob(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.ParseUint(params.ByName("id"), 10, 64)
	if err != nil {
		h.responseJSON(w, r, http.StatusBadRequest, err)
		return
	}

	// The upload URL outlives the finalize call, so refuse to write over an entry that is
	// already complete: its blob is what other jobs restore, and its recorded size would
	// no longer describe it. This is the guard the v1 upload route applies.
	if err := h.touchIncompleteCache(id); err != nil {
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
		var list blockList
		if err := xml.NewDecoder(io.LimitReader(r.Body, 8<<20)).Decode(&list); err != nil {
			h.responseJSON(w, r, http.StatusBadRequest, fmt.Errorf("malformed block list: %w", err))
			return
		}
		if err := h.storage.OrderBlocks(id, list.Latest); err != nil {
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

	w.WriteHeader(http.StatusCreated)
}

// touchIncompleteCache rejects an already finalized entry and, in the same transaction, bumps
// UsedAt so gcCache does not reap the reservation mid-upload: one bolt open per part, not two.
// An entry missing from the store is accepted, since the signature proves the id was handed out.
func (h *Handler) touchIncompleteCache(id uint64) error {
	db, err := h.openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	cache := &Cache{}
	if err := db.Get(id, cache); err != nil {
		if errors.Is(err, bolthold.ErrNotFound) {
			return nil
		}
		return err
	}
	if cache.Complete {
		return fmt.Errorf("cache %d: already complete", id)
	}
	cache.UsedAt = time.Now().Unix()
	return db.Update(cache.ID, cache)
}

// findReservedCache returns the reservation written to most recently, since the upload route
// touches UsedAt on every part. The request carries no entry id, so two jobs reserving the same
// key and version cannot be told apart here; the size check in Commit rejects a mismatched pair.
func findReservedCache(db *bolthold.Store, repo, key, version string) (*Cache, error) {
	cache := &Cache{}
	if err := db.FindOne(cache, bolthold.Where("Repo").Eq(repo).
		And("Key").Eq(key).
		And("Version").Eq(version).
		And("Complete").Eq(false).
		SortBy("UsedAt").Reverse()); err != nil {
		return nil, err
	}
	return cache, nil
}

// twirpMiss is the cache service's "no such entry": not an error, just nothing to restore.
func (h *Handler) twirpMiss(w http.ResponseWriter, r *http.Request) {
	h.responseJSON(w, r, http.StatusOK, map[string]any{"ok": false})
}

// twirpError reports a failure in the shape a twirp client expects, so the toolkit
// surfaces the message instead of a parse error.
func (h *Handler) twirpError(w http.ResponseWriter, r *http.Request, code string, err error) {
	h.logger.Debugf("%s %s: %v", r.Method, r.URL.Path, err)
	status := http.StatusBadRequest
	if code == twirpInternal {
		status = http.StatusInternalServerError
	}
	h.responseJSON(w, r, status, map[string]any{"code": code, "msg": err.Error()})
}

// The twirp request bodies. The toolkit's client serialises with useProtoFieldName, so the
// proto names are what arrive; the camelCase spellings are accepted too. json.Number takes
// a 64-bit integer as either the string the proto JSON mapping prescribes or a bare number.
type (
	v2CreateRequest struct {
		Key     string `json:"key"`
		Version string `json:"version"`
	}

	v2FinalizeRequest struct {
		Key            string      `json:"key"`
		Version        string      `json:"version"`
		SizeBytes      json.Number `json:"size_bytes"`
		SizeBytesCamel json.Number `json:"sizeBytes"`
	}

	v2DownloadRequest struct {
		Key              string   `json:"key"`
		Version          string   `json:"version"`
		RestoreKeys      []string `json:"restore_keys"`
		RestoreKeysCamel []string `json:"restoreKeys"`
	}
)

// keys is the lookup order: the exact key first, then the restore-key prefixes.
func (d v2DownloadRequest) keys() []string {
	restoreKeys := d.RestoreKeys
	if len(restoreKeys) == 0 {
		restoreKeys = d.RestoreKeysCamel
	}
	return append([]string{d.Key}, restoreKeys...)
}

func decodeTwirpRequest[T any](r *http.Request) (T, error) {
	var req T
	err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req)
	return req, err
}
