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
	"time"

	"gitea.com/gitea/runner/act/common"

	"github.com/sirupsen/logrus"
)

// cacheServicePrefix is the twirp path of the cache service, the only part of the results
// API this runner serves itself.
const cacheServicePrefix = "/twirp/github.actions.results.api.v1.CacheService/"

type Proxy struct {
	listener net.Listener
	server   *http.Server
	host     string
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
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
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

	return &Proxy{listener: listener, server: server, host: host}, nil
}

// ExternalURL is the value to hand to jobs as ACTIONS_RESULTS_URL. It carries no trailing
// slash, matching what the runner reports without the proxy.
func (p *Proxy) ExternalURL() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("http://%s:%d", p.host, p.listener.Addr().(*net.TCPAddr).Port)
}

func (p *Proxy) Close() error {
	if p == nil {
		return nil
	}
	var errs []error
	if p.server != nil {
		errs = append(errs, p.server.Close())
		p.server = nil
	}
	if p.listener != nil {
		if err := p.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			errs = append(errs, err)
		}
		p.listener = nil
	}
	return errors.Join(errs...)
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
	if !insecureTLS {
		return http.DefaultTransport
	}
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // the operator asked for it with runner.insecure
	return tr
}
