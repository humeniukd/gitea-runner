// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package client

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect"
	pingv1 "gitea.dev/actionslib/ping/v1"
	"gitea.dev/actionslib/pkg/protocol"
	"github.com/stretchr/testify/require"
)

func TestGetHTTPClientUsesProxyFromEnvironment(t *testing.T) {
	t.Setenv("HTTP_PROXY", "http://proxy.example.com:8080")

	client := getHTTPClient("http://gitea.example.com", nil, time.Minute, nil)
	require.Equal(t, time.Minute, client.Timeout)
	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)

	req, err := http.NewRequest(http.MethodGet, "http://gitea.example.com/api/actions/ping", nil)
	require.NoError(t, err)

	proxyURL, err := transport.Proxy(req)
	require.NoError(t, err)
	require.NotNil(t, proxyURL)
	require.Equal(t, "http://proxy.example.com:8080", proxyURL.String())
}

func TestGetHTTPClientAppliesTLSConfigOnlyToHTTPS(t *testing.T) {
	tlsConfig := &tls.Config{}
	require.Same(t, tlsConfig, getHTTPClient("https://gitea.example.com", tlsConfig, time.Minute, nil).Transport.(*http.Transport).TLSClientConfig)
	require.Nil(t, getHTTPClient("http://gitea.example.com", tlsConfig, time.Minute, nil).Transport.(*http.Transport).TLSClientConfig)
}

func TestNewSetsBaseURLAndHeaders(t *testing.T) {
	var gotPath string
	gotHeaders := make(http.Header)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// trailing slash must be trimmed before "/api/actions" is appended
	c := New(server.URL+"/", nil, "the-uuid", "the-token", time.Minute, map[string]string{
		"X-Proxy-Token":      "proxy-token",
		protocol.UUIDHeader:  "other-uuid",
		protocol.TokenHeader: "other-token",
		"Content-Type":       "text/plain",
	})
	// Address returns the endpoint as supplied (untrimmed)
	require.Equal(t, server.URL+"/", c.Address())

	// the call is expected to fail (server returns 500), we only assert what was sent
	_, _ = c.Ping(t.Context(), connect.NewRequest(&pingv1.PingRequest{Data: "hi"}))

	require.True(t, strings.HasPrefix(gotPath, "/api/actions/"), "unexpected path %q", gotPath)
	require.Equal(t, "the-uuid", gotHeaders.Get(protocol.UUIDHeader))
	require.Equal(t, "the-token", gotHeaders.Get(protocol.TokenHeader))
	require.Equal(t, "proxy-token", gotHeaders.Get("X-Proxy-Token"))
	require.Equal(t, "application/proto", gotHeaders.Get("Content-Type"))
}

func TestNewOmitsEmptyHeaders(t *testing.T) {
	gotHeaders := make(http.Header)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := New(server.URL, nil, "", "", time.Minute, nil)
	_, _ = c.Ping(t.Context(), connect.NewRequest(&pingv1.PingRequest{Data: "hi"}))

	require.Empty(t, gotHeaders.Get(protocol.UUIDHeader))
	require.Empty(t, gotHeaders.Get(protocol.TokenHeader))
}

func TestNewStopsRedirectLoop(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.URL.Path, http.StatusTemporaryRedirect)
	}))
	defer server.Close()

	c := New(server.URL, nil, "", "", time.Minute, nil)
	_, err := c.Ping(t.Context(), connect.NewRequest(&pingv1.PingRequest{Data: "hi"}))
	require.ErrorContains(t, err, "stopped after 10 redirects")
}

func TestNewRedirectKeepsExtraHeadersOnlyOnSameHost(t *testing.T) {
	var gotHeaders http.Header
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer target.Close()
	_, port, err := net.SplitHostPort(target.Listener.Addr().String())
	require.NoError(t, err)

	for _, tc := range []struct {
		newOrigin func(http.Handler) *httptest.Server
		host      string
		kept      bool
	}{
		{httptest.NewServer, "127.0.0.1", true},
		{httptest.NewServer, "localhost", false},
		{httptest.NewTLSServer, "127.0.0.1", false},
	} {
		gotHeaders = nil
		origin := tc.newOrigin(http.RedirectHandler("http://"+net.JoinHostPort(tc.host, port), http.StatusPermanentRedirect))
		c := New(origin.URL, &tls.Config{InsecureSkipVerify: true}, "the-uuid", "the-token", time.Minute, map[string]string{"X-Proxy-Token": "proxy-token", "Content-Type": "text/plain"})
		_, _ = c.Ping(t.Context(), connect.NewRequest(&pingv1.PingRequest{Data: "hi"}))
		origin.Close()

		require.NotNil(t, gotHeaders, origin.URL)
		require.Equal(t, tc.kept, gotHeaders.Get("X-Proxy-Token") != "", "%s to %s", origin.URL, tc.host)
		require.Equal(t, "the-token", gotHeaders.Get(protocol.TokenHeader))
		require.Equal(t, "application/proto", gotHeaders.Get("Content-Type"))
	}
}

func TestTransportReplacesSilentHTTP2Connection(t *testing.T) {
	sendPing, ping := http2SendPingTimeout, http2PingTimeout
	http2SendPingTimeout, http2PingTimeout = 100*time.Millisecond, 100*time.Millisecond
	t.Cleanup(func() { http2SendPingTimeout, http2PingTimeout = sendPing, ping })

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Proto))
	}))
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()
	addr, silence := startSilenceableProxy(t, server.Listener.Addr().String())

	const timeout = 5 * time.Second
	client := &http.Client{Transport: NewTransport("https://"+addr, server.Client().Transport.(*http.Transport).TLSClientConfig), Timeout: timeout}
	get := func() (string, error) {
		resp, err := client.Get("https://" + addr)
		if err != nil {
			return "", err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)
		return resp.Proto, nil
	}

	proto, err := get()
	require.NoError(t, err)
	require.Equal(t, "HTTP/2.0", proto)

	silence()
	start := time.Now()
	_, err = get()
	require.Error(t, err)
	require.Less(t, time.Since(start), timeout, "the silent connection was not closed by the ping health check")

	proto, err = get()
	require.NoError(t, err, "the next request must open a new connection")
	require.Equal(t, "HTTP/2.0", proto)
}

// startSilenceableProxy forwards TCP connections to target. silence makes the connections it
// already carries drop all traffic without closing, like a peer that vanished; new ones still pass.
func startSilenceableProxy(t *testing.T, target string) (addr string, silence func()) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	var mu sync.Mutex
	var quiet []*atomic.Bool
	var conns []net.Conn
	t.Cleanup(func() {
		_ = listener.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, conn := range conns {
			_ = conn.Close()
		}
	})
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			upstream, err := net.Dial("tcp", target)
			if err != nil {
				_ = conn.Close()
				continue
			}
			silent := new(atomic.Bool)
			mu.Lock()
			quiet = append(quiet, silent)
			conns = append(conns, conn, upstream)
			mu.Unlock()
			go forward(upstream, conn, silent)
			go forward(conn, upstream, silent)
		}
	}()
	return listener.Addr().String(), func() {
		mu.Lock()
		defer mu.Unlock()
		for _, silent := range quiet {
			silent.Store(true)
		}
	}
}

func forward(dst, src net.Conn, silent *atomic.Bool) {
	buf := make([]byte, 32<<10)
	for {
		n, err := src.Read(buf)
		if n > 0 && !silent.Load() {
			if _, err := dst.Write(buf[:n]); err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}
