// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubDaemon counts the /info requests that reach it.
func stubDaemon(t *testing.T, arch string) (dockerHost string, infoCalls *atomic.Int32) {
	t.Helper()

	infoCalls = &atomic.Int32{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/_ping"):
			w.Header().Set("Api-Version", "1.41")
			w.WriteHeader(http.StatusOK)
		case strings.HasSuffix(r.URL.Path, "/info"):
			infoCalls.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"Architecture":"` + arch + `","OSType":"linux"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	return "tcp://" + strings.TrimPrefix(server.URL, "http://"), infoCalls
}

func resetHostInfoCache(t *testing.T) {
	t.Helper()

	hostInfo.Store(nil)
	t.Cleanup(func() {
		hostInfo.Store(nil)
	})
}

func TestGetHostInfoQueriesTheDaemonOnce(t *testing.T) {
	resetHostInfoCache(t)
	dockerHost, infoCalls := stubDaemon(t, "x86_64")
	t.Setenv("DOCKER_HOST", dockerHost)

	for range 20 {
		info, err := GetHostInfo(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "x86_64", info.Architecture)
		assert.Equal(t, "linux", info.OSType)
	}

	assert.Equal(t, int32(1), infoCalls.Load())
}

func TestRunnerArchDoesNotRedialPerCall(t *testing.T) {
	resetHostInfoCache(t)
	dockerHost, infoCalls := stubDaemon(t, "x86_64")
	t.Setenv("DOCKER_HOST", dockerHost)

	for range 20 {
		assert.Equal(t, "X64", RunnerArch(context.Background()))
	}

	assert.Equal(t, int32(1), infoCalls.Load())
}

func TestGetHostInfoDoesNotAnswerForAnotherDaemon(t *testing.T) {
	resetHostInfoCache(t)
	amdHost, amdCalls := stubDaemon(t, "x86_64")
	armHost, armCalls := stubDaemon(t, "aarch64")

	t.Setenv("DOCKER_HOST", amdHost)
	info, err := GetHostInfo(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "x86_64", info.Architecture)

	t.Setenv("DOCKER_HOST", armHost)
	info, err = GetHostInfo(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "aarch64", info.Architecture)

	t.Setenv("DOCKER_HOST", armHost)
	info, err = GetHostInfo(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "aarch64", info.Architecture)

	assert.Equal(t, int32(1), amdCalls.Load())
	assert.Equal(t, int32(1), armCalls.Load())
}

func TestGetHostInfoDoesNotCacheFailures(t *testing.T) {
	resetHostInfoCache(t)

	failing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/_ping") {
			w.Header().Set("Api-Version", "1.41")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failing.Close()

	t.Setenv("DOCKER_HOST", "tcp://"+strings.TrimPrefix(failing.URL, "http://"))
	_, err := GetHostInfo(context.Background())
	require.Error(t, err)

	assert.Nil(t, hostInfo.Load())
}
