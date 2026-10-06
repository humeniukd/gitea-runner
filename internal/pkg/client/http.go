// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package client

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitea.com/gitea/runner/internal/pkg/ver"

	"connectrpc.com/connect"
	"gitea.dev/actionslib/ping/v1/pingv1connect"
	"gitea.dev/actionslib/pkg/protocol"
	"gitea.dev/actionslib/runner/v1/runnerv1connect"
)

func getHTTPClient(endpoint string, insecure bool, timeout time.Duration, extraHeaders http.Header) *http.Client {
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10, // All requests go to one host; default is 2 which causes frequent reconnects.
		IdleConnTimeout:     90 * time.Second,
	}
	if strings.HasPrefix(endpoint, "https://") && insecure {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 { // a custom CheckRedirect replaces net/http's own limit
				return errors.New("stopped after 10 redirects")
			}
			if sameOrigin(req.URL, via[0].URL) {
				return nil
			}

			for header := range extraHeaders {
				req.Header.Del(header)
			}
			req.Header.Del(protocol.UUIDHeader)
			req.Header.Del(protocol.TokenHeader)
			return nil
		},
	}
}

func sameOrigin(a, b *url.URL) bool {
	return a.Scheme == b.Scheme && a.Host == b.Host
}

// New returns a new runner client. timeout bounds every RPC: without it a
// stalled connection parks the reporter for the whole job context, so logs and
// heartbeats stop together and the task is reaped as a zombie.
func New(endpoint string, insecure bool, uuid, token string, timeout time.Duration, extraHeaders map[string]string, opts ...connect.ClientOption) *HTTPClient {
	baseURL := strings.TrimRight(endpoint, "/") + "/api/actions"
	headers := make(http.Header, len(extraHeaders))
	for name, value := range extraHeaders {
		headers.Set(name, value)
	}
	headers.Del("User-Agent")
	headers.Del(protocol.UUIDHeader)
	headers.Del(protocol.TokenHeader)

	opts = append(opts, connect.WithInterceptors(connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			for name, values := range headers {
				req.Header().Del(name)
				for _, value := range values {
					req.Header().Add(name, value)
				}
			}
			req.Header().Set("User-Agent", "gitea-runner/"+ver.Version())
			if uuid != "" {
				req.Header().Set(protocol.UUIDHeader, uuid)
			}
			if token != "" {
				req.Header().Set(protocol.TokenHeader, token)
			}
			return next(ctx, req)
		}
	})))

	httpClient := getHTTPClient(endpoint, insecure, timeout, headers)
	return &HTTPClient{
		PingServiceClient: pingv1connect.NewPingServiceClient(
			httpClient,
			baseURL,
			opts...,
		),
		RunnerServiceClient: runnerv1connect.NewRunnerServiceClient(
			httpClient,
			baseURL,
			opts...,
		),
		endpoint: endpoint,
	}
}

func (c *HTTPClient) Address() string {
	return c.endpoint
}

var _ Client = (*HTTPClient)(nil)

// An HTTPClient manages communication with the runner API.
type HTTPClient struct {
	pingv1connect.PingServiceClient
	runnerv1connect.RunnerServiceClient
	endpoint string
}
