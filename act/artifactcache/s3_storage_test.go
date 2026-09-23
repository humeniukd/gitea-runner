// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewS3StorageRejectsInvalidEndpoint(t *testing.T) {
	_, err := NewS3Storage(t.TempDir(), S3Options{Endpoint: "minio", Bucket: "cache"})
	if err == nil || !strings.Contains(err.Error(), "invalid S3 endpoint") {
		t.Fatalf("expected invalid endpoint error, got %v", err)
	}
}

// TestS3StorageIntegration exercises MinIO/RustFS-compatible servers when a
// test bucket is supplied by CI or an operator. It is intentionally opt-in:
// unit tests must not require Docker or a reachable object store.
func TestS3StorageIntegration(t *testing.T) {
	opts := S3Options{
		Endpoint:  os.Getenv("RUNNER_TEST_S3_ENDPOINT"),
		Bucket:    os.Getenv("RUNNER_TEST_S3_BUCKET"),
		AccessKey: os.Getenv("RUNNER_TEST_S3_ACCESS_KEY"),
		SecretKey: os.Getenv("RUNNER_TEST_S3_SECRET_KEY"),
		Region:    os.Getenv("RUNNER_TEST_S3_REGION"),
		PathStyle: true,
		Prefix:    "runner-test-" + strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	if opts.Endpoint == "" || opts.Bucket == "" || opts.AccessKey == "" || opts.SecretKey == "" {
		t.Skip("set RUNNER_TEST_S3_ENDPOINT, RUNNER_TEST_S3_BUCKET, RUNNER_TEST_S3_ACCESS_KEY, and RUNNER_TEST_S3_SECRET_KEY to run against MinIO or RustFS")
	}
	if opts.Region == "" {
		opts.Region = "us-east-1"
	}
	storage, err := NewS3Storage(t.TempDir(), opts)
	if err != nil {
		t.Fatal(err)
	}
	const id = 42
	if err := storage.Write(id, 0, strings.NewReader("cache-data")); err != nil {
		t.Fatal(err)
	}
	if _, err := storage.Commit(id, int64(len("cache-data"))); err != nil {
		t.Fatal(err)
	}
	if ok, err := storage.Exist(id); err != nil || !ok {
		t.Fatalf("cache object does not exist: ok=%v err=%v", ok, err)
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	storage.Serve(response, request, id)
	if response.Code != 200 || response.Body.String() != "cache-data" {
		t.Fatalf("unexpected download: status=%d body=%q", response.Code, response.Body.String())
	}
	if err := storage.Remove(id); err != nil {
		t.Fatal(err)
	}
}
