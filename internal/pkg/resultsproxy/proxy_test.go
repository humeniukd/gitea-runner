// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package resultsproxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startProxy fronts two recording upstreams and returns the proxy plus the requests each
// upstream saw.
func startProxy(t *testing.T) (proxy *Proxy, gitea, cache *[]*http.Request) {
	t.Helper()

	giteaSeen := &[]*http.Request{}
	cacheSeen := &[]*http.Request{}

	giteaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*giteaSeen = append(*giteaSeen, r.Clone(t.Context()))
		body, _ := io.ReadAll(r.Body)
		_, _ = w.Write([]byte("gitea:" + r.URL.Path + ":" + string(body)))
	}))
	t.Cleanup(giteaServer.Close)

	cacheServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*cacheSeen = append(*cacheSeen, r.Clone(t.Context()))
		_, _ = w.Write([]byte("cache:" + r.URL.Path))
	}))
	t.Cleanup(cacheServer.Close)

	proxy, err := Start("127.0.0.1", 0, giteaServer.URL, cacheServer.URL+"/", false, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = proxy.Close() })

	return proxy, giteaSeen, cacheSeen
}

func get(t *testing.T, url string) string {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(body)
}

func TestProxyRoutesCacheServiceToTheCacheServer(t *testing.T) {
	proxy, giteaSeen, cacheSeen := startProxy(t)

	body := get(t, proxy.ExternalURL()+cacheServicePrefix+"GetCacheEntryDownloadURL")

	assert.Equal(t, "cache:"+cacheServicePrefix+"GetCacheEntryDownloadURL", body)
	assert.Len(t, *cacheSeen, 1)
	assert.Empty(t, *giteaSeen)
}

// Everything that is not the cache service belongs to Gitea; the artifact service above
// all, which must keep working exactly as it does without the proxy.
func TestProxyForwardsEverythingElseToGitea(t *testing.T) {
	proxy, giteaSeen, cacheSeen := startProxy(t)

	for _, path := range []string{
		"/twirp/github.actions.results.api.v1.ArtifactService/CreateArtifact",
		"/api/actions_pipeline/_apis/pipelines/workflows/1/artifacts",
		"/",
	} {
		assert.Equal(t, "gitea:"+path+":", get(t, proxy.ExternalURL()+path))
	}

	assert.Len(t, *giteaSeen, 3)
	assert.Empty(t, *cacheSeen)
}

func TestProxyPreservesRequestDetails(t *testing.T) {
	proxy, giteaSeen, _ := startProxy(t)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
		proxy.ExternalURL()+"/twirp/github.actions.results.api.v1.ArtifactService/FinalizeArtifact?name=x",
		strings.NewReader(`{"ok":true}`))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer token")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, `gitea:/twirp/github.actions.results.api.v1.ArtifactService/FinalizeArtifact:{"ok":true}`, string(body))
	require.Len(t, *giteaSeen, 1)
	forwarded := (*giteaSeen)[0]
	assert.Equal(t, http.MethodPost, forwarded.Method)
	assert.Equal(t, "name=x", forwarded.URL.RawQuery)
	assert.Equal(t, "Bearer token", forwarded.Header.Get("Authorization"))
	// The Host header follows the target, so a virtual-hosted Gitea routes the request.
	assert.NotEqual(t, proxy.ExternalURL(), forwarded.Host)
}

func TestProxyReportsAnUnreachableUpstream(t *testing.T) {
	// A port nothing listens on: the proxy must answer, not hang or panic.
	proxy, err := Start("127.0.0.1", 0, "http://127.0.0.1:1", "http://127.0.0.1:1", false, nil)
	require.NoError(t, err)
	defer proxy.Close()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, proxy.ExternalURL()+"/anything", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
}

func TestStartRejectsAnUnusableTarget(t *testing.T) {
	_, err := Start("127.0.0.1", 0, "", "http://127.0.0.1:1", false, nil)
	require.ErrorContains(t, err, "gitea url")

	_, err = Start("127.0.0.1", 0, "http://127.0.0.1:1", "not-a-url", false, nil)
	require.ErrorContains(t, err, "cache url")
}

func TestCloseIsIdempotent(t *testing.T) {
	proxy, _, _ := startProxy(t)

	require.NoError(t, proxy.Close())
	require.NoError(t, proxy.Close())
	require.NoError(t, (*Proxy)(nil).Close())
	assert.Empty(t, (*Proxy)(nil).ExternalURL())
}
