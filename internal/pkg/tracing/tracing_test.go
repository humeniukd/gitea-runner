// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package tracing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"gitea.com/gitea/runner/internal/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// unsetConfigFile removes OTEL_CONFIG_FILE and restores it at test end;
// otelconf reacts to a present-but-empty value by trying to read a "" file.
func unsetConfigFile(t *testing.T) {
	t.Helper()
	old, had := os.LookupEnv("OTEL_CONFIG_FILE")
	require.NoError(t, os.Unsetenv("OTEL_CONFIG_FILE"))
	t.Cleanup(func() {
		if had {
			t.Setenv("OTEL_CONFIG_FILE", old)
		}
	})
}

func TestJobEnvPassthroughAndTraceparent(t *testing.T) {
	t.Setenv("OTEL_TEST_PROP", "1")
	t.Setenv("OTHER_VAR", "skip")
	t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "authorization=Bearer secret")
	// Ensure a stale OTEL_CONFIG_FILE cannot influence anything here.
	unsetConfigFile(t)

	tracer := sdktrace.NewTracerProvider().Tracer("test")
	ctx, span := tracer.Start(context.Background(), "op")
	defer span.End()

	envs := JobEnv(ctx)
	assert.Equal(t, "1", envs["OTEL_TEST_PROP"])
	assert.NotContains(t, envs, "OTHER_VAR")

	cctx := propagation.TraceContext{}.Extract(context.Background(), propagation.MapCarrier{"traceparent": envs["TRACEPARENT"]})
	assert.True(t, trace.SpanContextFromContext(cctx).IsValid())
	assert.NotContains(t, envs, "TRACESTATE")

	// Collector credentials stay on the runner and out of job containers.
	assert.NotContains(t, envs, "OTEL_EXPORTER_OTLP_HEADERS")
	assert.NotContains(t, JobEnv(context.Background()), "TRACEPARENT")
}

// captureServer counts POSTs to /v1/traces and records the Authorization and
// X-Extra headers it received.
func captureServer(t *testing.T, hits *atomic.Int64, authHeader, extraHeader *atomic.Pointer[string]) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/traces" && r.Method == http.MethodPost {
			hits.Add(1)
			if authHeader != nil {
				v := r.Header.Get("Authorization")
				authHeader.Store(&v)
			}
			if extraHeader != nil {
				x := r.Header.Get("X-Extra")
				extraHeader.Store(&x)
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func resetGlobals(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		otel.SetTracerProvider(noop.NewTracerProvider())
	})
}

func waitFor(t *testing.T, hits *atomic.Int64) {
	t.Helper()
	require.Eventually(t, func() bool { return hits.Load() > 0 }, 10*time.Second, 10*time.Millisecond)
}

func TestSetupNoopWithoutEndpoint(t *testing.T) {
	unsetConfigFile(t)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")

	shutdown, err := Setup(context.Background(), "runner", config.Otel{})
	require.NoError(t, err)
	require.NoError(t, shutdown(context.Background()))
}

func TestSetupExportsToEndpoint(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")
	var hits atomic.Int64
	srv := captureServer(t, &hits, nil, nil)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", srv.URL)

	shutdown, err := Setup(context.Background(), "runner", config.Otel{})
	require.NoError(t, err)

	_, span := otel.Tracer("test").Start(context.Background(), "op")
	span.End()
	require.NoError(t, shutdown(context.Background()))
	waitFor(t, &hits)
}

func TestSetupYmlEndpointWinsOverEnv(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")
	var envHits, ymlHits atomic.Int64
	envSrv := captureServer(t, &envHits, nil, nil)
	ymlSrv := captureServer(t, &ymlHits, nil, nil)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", envSrv.URL)

	shutdown, err := Setup(context.Background(), "runner", config.Otel{Endpoint: ymlSrv.URL})
	require.NoError(t, err)

	_, span := otel.Tracer("test").Start(context.Background(), "op")
	span.End()
	require.NoError(t, shutdown(context.Background()))
	waitFor(t, &ymlHits)
	assert.Equal(t, int64(0), envHits.Load(), "yml endpoint must take precedence over env")
}

func TestSetupYmlHeadersOverrideEnv(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")
	var hits atomic.Int64
	var auth atomic.Pointer[string]
	var extra atomic.Pointer[string]
	srv := captureServer(t, &hits, &auth, &extra)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", srv.URL)
	t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "authorization=Bearer env-token,x-extra=from-env")

	shutdown, err := Setup(context.Background(), "runner", config.Otel{
		Endpoint: srv.URL,
		Headers:  map[string]string{"authorization": "Bearer yml-token"},
	})
	require.NoError(t, err)

	_, span := otel.Tracer("test").Start(context.Background(), "op")
	span.End()
	require.NoError(t, shutdown(context.Background()))

	require.Eventually(t, func() bool {
		return auth.Load() != nil && extra.Load() != nil
	}, 10*time.Second, 10*time.Millisecond)
	assert.Equal(t, "Bearer yml-token", *auth.Load(), "yml header wins per key")
	assert.Equal(t, "from-env", *extra.Load(), "env-only headers survive the merge")
}
