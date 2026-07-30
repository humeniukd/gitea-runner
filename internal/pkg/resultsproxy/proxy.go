// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// Package resultsproxy fronts ACTIONS_RESULTS_URL for job containers.
//
// The results API hosts several services under one base URL. Gitea implements the
// artifact service; the cache service is implemented by this runner's cache server. A
// client only knows the one base URL, so this proxy answers for both: cache service calls
// go to the cache server, everything else is forwarded to Gitea untouched.
package resultsproxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"gitea.com/gitea/runner/act/common"

	"github.com/sirupsen/logrus"
)

// cacheServicePrefix is the twirp path of the cache service, the only part of the results
// API this runner serves itself.
const cacheServicePrefix = "/twirp/github.actions.results.api.v1.CacheService/"

type Proxy struct {
	server *http.Server
	url    string
}

// Start listens on port (0 picks a free one) and proxies to gitea and cache. host is the
// address job containers use to reach the proxy; it is detected when empty, the same way
// the cache server detects it.
func Start(host string, port uint16, gitea, cache string, insecureTLS bool, logger logrus.FieldLogger) (*Proxy, error) {
	giteaURL, err := parseTarget(gitea)
	if err != nil {
		return nil, fmt.Errorf("results proxy: gitea url: %w", err)
	}
	cacheURL, err := parseTarget(cache)
	if err != nil {
		return nil, fmt.Errorf("results proxy: cache url: %w", err)
	}

	if host == "" {
		ip := common.GetOutboundIP()
		if ip == nil {
			return nil, errors.New("results proxy: unable to determine outbound IP address")
		}
		host = ip.String()
	}

	if logger == nil {
		logger = logrus.StandardLogger()
	}
	logger = logger.WithField("module", "resultsproxy")

	reverse := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			target := giteaURL
			if strings.HasPrefix(r.In.URL.Path, cacheServicePrefix) {
				target = cacheURL
			}
			r.SetURL(target)
			// SetURL keeps the inbound path appended to the target path, which is what
			// both services expect, but the Host header has to follow the target so a
			// virtual-hosted Gitea routes the request.
			r.Out.Host = target.Host
			r.SetXForwarded()
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Warnf("proxying %s %s failed: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusBadGateway)
		},
		Transport: transport(insecureTLS),
		// Artifact uploads are many parallel chunked requests; the default is a fresh 32 KiB
		// buffer per request.
		BufferPool: &bufferPool{},
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("results proxy: %w", err)
	}
	_, boundPort, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("results proxy: %w", err)
	}
	server := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           reverse,
	}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, net.ErrClosed) {
			logger.Errorf("http serve: %v", err)
		}
	}()

	return &Proxy{server: server, url: "http://" + net.JoinHostPort(host, boundPort)}, nil
}

// ExternalURL is the value to hand to jobs as ACTIONS_RESULTS_URL. It carries no trailing
// slash, matching what the runner reports without the proxy.
func (p *Proxy) ExternalURL() string {
	if p == nil {
		return ""
	}
	return p.url
}

func (p *Proxy) Close() error {
	if p == nil {
		return nil
	}
	return p.server.Close()
}

func parseTarget(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, errors.New("no address configured")
	}
	target, err := url.Parse(strings.TrimSuffix(raw, "/"))
	if err != nil {
		return nil, err
	}
	if target.Host == "" {
		return nil, fmt.Errorf("%q has no host", raw)
	}
	return target, nil
}

// transport clones the default transport so that a Gitea instance with a certificate the
// runner is configured to skip verifying stays reachable through the proxy as well.
func transport(insecureTLS bool) http.RoundTripper {
	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return http.DefaultTransport
	}
	tr := defaultTransport.Clone()
	// Artifact uploads run many chunks in parallel, and the default of 2 idle connections
	// per host would make most of them reconnect.
	tr.MaxIdleConnsPerHost = 32
	if insecureTLS {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // the operator asked for it with runner.insecure
	}
	return tr
}

// bufferPool reuses the copy buffers the proxy would otherwise allocate per request.
type bufferPool struct{ pool sync.Pool }

func (b *bufferPool) Get() []byte {
	if buf, ok := b.pool.Get().(*[]byte); ok {
		return *buf
	}
	return make([]byte, 256<<10)
}

func (b *bufferPool) Put(buf []byte) { b.pool.Put(&buf) }
