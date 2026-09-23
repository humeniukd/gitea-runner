// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/timshannon/bolthold"
)

// S3Options describes the static-credential S3-compatible endpoint used for
// actions/cache data. It deliberately has no artifact settings: Gitea still
// owns workflow artifacts.
type S3Options struct {
	Endpoint  string
	Region    string
	Bucket    string
	Prefix    string
	PathStyle bool
	AccessKey string
	SecretKey string
}

type cacheStorage interface {
	Exist(uint64) (bool, error)
	Write(uint64, int64, io.Reader) error
	WriteBlock(uint64, string, io.Reader) error
	OrderBlocks(uint64, []string) error
	Commit(uint64, int64) (int64, error)
	Serve(http.ResponseWriter, *http.Request, uint64)
	Remove(uint64) error
	filename(uint64) string
}

type cacheMetadata interface {
	syncMetadata(*bolthold.Store) error
	persistMetadata(*Cache) error
	removeMetadata(uint64) error
}

// S3Storage keeps temporary, in-progress ranges on the runner, and writes
// only completed archives to S3. This preserves v1 range and v2 block upload
// semantics while keeping durable cache data out of the runner filesystem.
type S3Storage struct {
	staging *Storage
	client  *minio.Client
	bucket  string
	prefix  string
}

func NewS3Storage(stagingDir string, opts S3Options) (*S3Storage, error) {
	u, err := url.Parse(opts.Endpoint)
	if err != nil || u.Scheme == "" || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return nil, fmt.Errorf("invalid S3 endpoint %q", opts.Endpoint)
	}
	lookup := minio.BucketLookupDNS
	if opts.PathStyle {
		lookup = minio.BucketLookupPath
	}
	client, err := minio.New(u.Host, &minio.Options{
		Creds:        credentials.NewStaticV4(opts.AccessKey, opts.SecretKey, ""),
		Secure:       u.Scheme == "https",
		Region:       opts.Region,
		BucketLookup: lookup,
	})
	if err != nil {
		return nil, fmt.Errorf("create S3 client: %w", err)
	}
	client.SetAppInfo("gitea-runner", "cache")
	storage, err := NewStorage(stagingDir)
	if err != nil {
		return nil, err
	}
	return &S3Storage{staging: storage, client: client, bucket: opts.Bucket, prefix: opts.Prefix}, nil
}

func (s *S3Storage) key(id uint64) string {
	return path.Join(s.prefix, "blobs", strconv.FormatUint(id, 10))
}

func (s *S3Storage) filename(id uint64) string {
	return s.staging.filename(id)
}

func (s *S3Storage) metadataKey(id uint64) string {
	return path.Join(s.prefix, "metadata", strconv.FormatUint(id, 10)+".json")
}

// syncMetadata imports cache records that other runners committed. The local
// Bolt database remains the handler's fast working index; S3 is the durable
// shared source of completed cache records.
func (s *S3Storage) syncMetadata(db *bolthold.Store) error {
	objects := s.client.ListObjects(context.Background(), s.bucket, minio.ListObjectsOptions{Prefix: path.Join(s.prefix, "metadata") + "/", Recursive: true})
	for object := range objects {
		if object.Err != nil {
			return fmt.Errorf("list S3 cache metadata: %w", object.Err)
		}
		reader, err := s.client.GetObject(context.Background(), s.bucket, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return fmt.Errorf("get S3 cache metadata: %w", err)
		}
		var remote Cache
		err = json.UnmarshalRead(reader, &remote)
		_ = reader.Close()
		if err != nil {
			return fmt.Errorf("decode S3 cache metadata %q: %w", object.Key, err)
		}
		if !remote.Complete {
			continue
		}
		local := &Cache{}
		err = db.Get(remote.ID, local)
		if errors.Is(err, bolthold.ErrNotFound) {
			if err := db.Insert(remote.ID, &remote); err != nil {
				return fmt.Errorf("import S3 cache metadata: %w", err)
			}
			continue
		}
		if err != nil {
			return err
		}
		if remote.UpdatedAt > local.UpdatedAt {
			if err := db.Update(remote.ID, &remote); err != nil {
				return fmt.Errorf("update S3 cache metadata: %w", err)
			}
		}
	}
	return nil
}

func (s *S3Storage) persistMetadata(cache *Cache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(context.Background(), s.bucket, s.metadataKey(cache.ID), bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return fmt.Errorf("write S3 cache metadata: %w", err)
	}
	return nil
}

func (s *S3Storage) removeMetadata(id uint64) error {
	if err := s.client.RemoveObject(context.Background(), s.bucket, s.metadataKey(id), minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("remove S3 cache metadata: %w", err)
	}
	return nil
}

func (s *S3Storage) Exist(id uint64) (bool, error) {
	_, err := s.client.StatObject(context.Background(), s.bucket, s.key(id), minio.StatObjectOptions{})
	if minio.ToErrorResponse(err).Code == "NoSuchKey" || minio.ToErrorResponse(err).Code == "NoSuchObject" {
		return false, nil
	}
	return err == nil, err
}

func (s *S3Storage) Write(id uint64, offset int64, reader io.Reader) error {
	return s.staging.Write(id, offset, reader)
}

func (s *S3Storage) WriteBlock(id uint64, blockID string, reader io.Reader) error {
	return s.staging.WriteBlock(id, blockID, reader)
}

func (s *S3Storage) OrderBlocks(id uint64, blockIDs []string) error {
	return s.staging.OrderBlocks(id, blockIDs)
}

func (s *S3Storage) Commit(id uint64, size int64) (int64, error) {
	written, err := s.staging.Commit(id, size)
	if err != nil {
		return 0, err
	}
	name := s.staging.filename(id)
	defer func() { _ = s.staging.Remove(id) }()
	if _, err := s.client.FPutObject(context.Background(), s.bucket, s.key(id), name, minio.PutObjectOptions{ContentType: "application/octet-stream"}); err != nil {
		return 0, fmt.Errorf("upload cache to S3: %w", err)
	}
	return written, nil
}

func (s *S3Storage) Serve(w http.ResponseWriter, r *http.Request, id uint64) {
	object, err := s.client.GetObject(r.Context(), s.bucket, s.key(id), minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, "read cache from S3", http.StatusBadGateway)
		return
	}
	defer object.Close()
	info, err := object.Stat()
	if err != nil {
		http.Error(w, "cache not found", http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, strconv.FormatUint(id, 10), info.LastModified, object)
}

func (s *S3Storage) Remove(id uint64) error {
	if err := s.client.RemoveObject(context.Background(), s.bucket, s.key(id), minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("remove cache from S3: %w", err)
	}
	return s.staging.Remove(id)
}
