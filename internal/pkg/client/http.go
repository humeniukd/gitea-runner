// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"gitea.com/gitea/runner/internal/pkg/ver"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"gitea.dev/actionslib/ping/v1/pingv1connect"
	"gitea.dev/actionslib/pkg/protocol"
	"gitea.dev/actionslib/runner/v1/runnerv1connect"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/propagation"
)

func getHTTPClient(endpoint string, insecure bool, timeout time.Duration) *http.Client {
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
	return &http.Client{Transport: transport, Timeout: timeout}
}

// New returns a new runner client. timeout bounds every RPC: without it a
// stalled connection parks the reporter for the whole job context, so logs and
// heartbeats stop together and the task is reaped as a zombie.
//
// rpcTraces enables client RPC spans and W3C trace-context propagation to Gitea.
// It is off by default because the control-plane RPCs (FetchTask polls, log and
// status updates) are very noisy; when off, the interceptor still records RPC
// metrics but creates no spans and sends no traceparent to Gitea.
func New(endpoint string, insecure bool, uuid, token string, timeout time.Duration, rpcTraces bool, opts ...connect.ClientOption) *HTTPClient {
	baseURL := strings.TrimRight(endpoint, "/") + "/api/actions"

	// Instruments every RPC with metrics; when rpcTraces is on, also spans nested
	// under the caller's ctx (job span for reporter calls) and traceparent
	// propagation. When off, an empty propagator suppresses trace headers too.
	otelOpts := []otelconnect.Option{}
	if !rpcTraces {
		otelOpts = append(otelOpts,
			otelconnect.WithoutTracing(),
			otelconnect.WithPropagator(propagation.NewCompositeTextMapPropagator()),
		)
	}
	otelInterceptor, err := otelconnect.NewInterceptor(otelOpts...)
	if err != nil {
		log.WithError(err).Warn("otel rpc instrumentation disabled")
	} else {
		opts = append(opts, connect.WithInterceptors(otelInterceptor))
	}
	opts = append(opts, connect.WithInterceptors(connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
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

	httpClient := getHTTPClient(endpoint, insecure, timeout)
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
