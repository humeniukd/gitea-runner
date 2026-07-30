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

// startProxy fronts two upstreams that echo which one was reached, and records the
// requests Gitea saw.
func startProxy(t *testing.T) (*Proxy, *[]*http.Request) {
	t.Helper()

	giteaSeen := &[]*http.Request{}
	upstream := func(name string, seen *[]*http.Request) *httptest.Server {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if seen != nil {
				*seen = append(*seen, r.Clone(t.Context()))
			}
			body, _ := io.ReadAll(r.Body)
			_, _ = w.Write([]byte(name + ":" + r.URL.Path + ":" + string(body)))
		}))
		t.Cleanup(server.Close)
		return server
	}

	proxy, err := Start("127.0.0.1", 0, upstream("gitea", giteaSeen).URL, upstream("cache", nil).URL+"/", false, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = proxy.Close() })

	return proxy, giteaSeen
}

// The cache service is answered by the runner; everything else belongs to Gitea, the
// artifact service above all, which must keep working exactly as it does without the proxy.
func TestProxyRoutesByPath(t *testing.T) {
	proxy, _ := startProxy(t)

	for _, tc := range []struct{ path, want string }{
		{cacheServicePrefix + "GetCacheEntryDownloadURL", "cache"},
		{"/twirp/github.actions.results.api.v1.ArtifactService/CreateArtifact", "gitea"},
		{"/api/actions_pipeline/_apis/pipelines/workflows/1/artifacts", "gitea"},
		{"/", "gitea"},
	} {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, proxy.ExternalURL()+tc.path, nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		_ = resp.Body.Close()

		assert.Equal(t, tc.want+":"+tc.path+":", string(body))
	}
}

func TestProxyPreservesRequestDetails(t *testing.T) {
	proxy, giteaSeen := startProxy(t)

	const path = "/twirp/github.actions.results.api.v1.ArtifactService/FinalizeArtifact"
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
		proxy.ExternalURL()+path+"?name=x", strings.NewReader(`{"ok":true}`))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, `gitea:`+path+`:{"ok":true}`, string(body))
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
	proxy, _ := startProxy(t)
	url := proxy.ExternalURL()

	require.NoError(t, proxy.Close())
	require.NoError(t, proxy.Close())
	require.NoError(t, (*Proxy)(nil).Close())

	// Reported after shutdown too, so a caller that logs it cannot panic on a closed proxy.
	assert.Equal(t, url, proxy.ExternalURL())
	assert.Empty(t, (*Proxy)(nil).ExternalURL())
}
