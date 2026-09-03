// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package tracing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"gitea.com/gitea/runner/internal/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
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

// clearEnv drops the OTLP endpoint env so a test starts from a known state.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		"OTEL_EXPORTER_OTLP_ENDPOINT",
		"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT",
		"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT",
		"OTEL_EXPORTER_OTLP_METRIC_ENDPOINT",
		"OTEL_EXPORTER_OTLP_HEADERS",
		"OTEL_SDK_DISABLED",
	} {
		t.Setenv(name, "")
	}
}

// otlpSink records the request paths and Authorization headers of OTLP posts.
type otlpSink struct {
	mu    sync.Mutex
	paths map[string]int
	auth  map[string]string
	*httptest.Server
}

func newSink(t *testing.T) *otlpSink {
	t.Helper()
	s := &otlpSink{paths: map[string]int{}, auth: map[string]string{}}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.mu.Lock()
			s.paths[r.URL.Path]++
			if a := r.Header.Get("Authorization"); a != "" {
				s.auth[r.URL.Path] = a
			}
			s.mu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *otlpSink) count(path string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.paths[path]
}

func (s *otlpSink) total() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, c := range s.paths {
		n += c
	}
	return n
}

func resetGlobals(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		otel.SetTracerProvider(tracenoop.NewTracerProvider())
		otel.SetMeterProvider(metricnoop.NewMeterProvider())
	})
}

// emit produces one span and one metric point, then flushes via shutdown.
func emit(t *testing.T, shutdown func(context.Context) error) {
	t.Helper()
	_, span := otel.Tracer("test").Start(context.Background(), "op")
	span.End()
	counter, err := otel.Meter("test").Int64Counter("probe")
	require.NoError(t, err)
	counter.Add(context.Background(), 1)
	sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, shutdown(sctx))
}

func TestJobEnvPassthroughAndTraceparent(t *testing.T) {
	t.Setenv("OTEL_TEST_PROP", "1")
	t.Setenv("OTHER_VAR", "skip")
	t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "authorization=Bearer secret")
	unsetConfigFile(t)

	tracer := sdktrace.NewTracerProvider().Tracer("test")
	ctx, span := tracer.Start(context.Background(), "op")
	defer span.End()

	envs := JobEnv(ctx)
	assert.Equal(t, "1", envs["OTEL_TEST_PROP"])
	assert.NotContains(t, envs, "OTHER_VAR")

	cctx := propagation.TraceContext{}.Extract(context.Background(), propagation.MapCarrier{"traceparent": envs["TRACEPARENT"]})
	assert.True(t, trace.SpanContextFromContext(cctx).IsValid())

	// Collector credentials stay on the runner and out of job containers.
	assert.NotContains(t, envs, "OTEL_EXPORTER_OTLP_HEADERS")
	assert.NotContains(t, JobEnv(context.Background()), "TRACEPARENT")
}

func TestSetupNoopWithoutEndpoint(t *testing.T) {
	unsetConfigFile(t)
	clearEnv(t)

	shutdown, err := Setup(context.Background(), "runner", config.Otel{})
	require.NoError(t, err)
	require.NoError(t, shutdown(context.Background()))
}

func TestSetupSharedEndpointEnablesBothSignals(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	clearEnv(t)
	sink := newSink(t)

	shutdown, err := Setup(context.Background(), "runner", config.Otel{Endpoint: sink.URL})
	require.NoError(t, err)
	emit(t, shutdown)

	// Generic endpoint: the signal path is appended for each signal.
	assert.Positive(t, sink.count("/v1/traces"), "traces exported to appended path")
	assert.Positive(t, sink.count("/v1/metrics"), "metrics exported to appended path")
}

func TestSetupSignalEndpointUsedAsIs(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	clearEnv(t)
	sink := newSink(t)

	shutdown, err := Setup(context.Background(), "runner", config.Otel{
		Traces: config.OtelTraces{Endpoint: sink.URL + "/custom/traces"},
	})
	require.NoError(t, err)
	emit(t, shutdown)

	// Signal endpoint used as-is: no signal path appended.
	assert.Positive(t, sink.count("/custom/traces"))
	assert.Equal(t, 0, sink.count("/custom/traces/v1/traces"))
}

func TestSetupSharedEnabledFalseDisablesBoth(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	clearEnv(t)
	sink := newSink(t)

	shutdown, err := Setup(context.Background(), "runner", config.Otel{
		Enabled:  new(false),
		Endpoint: sink.URL,
	})
	require.NoError(t, err)
	emit(t, shutdown)
	assert.Equal(t, 0, sink.total(), "otel.enabled=false disables both signals")
}

func TestSetupSignalEnabledOverridesShared(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	clearEnv(t)
	sink := newSink(t)

	// Shared disabled, but traces re-enabled per signal.
	shutdown, err := Setup(context.Background(), "runner", config.Otel{
		Enabled:  new(false),
		Endpoint: sink.URL,
		Traces:   config.OtelTraces{Enabled: new(true)},
	})
	require.NoError(t, err)
	emit(t, shutdown)
	assert.Positive(t, sink.count("/v1/traces"), "traces.enabled=true overrides otel.enabled=false")
	assert.Equal(t, 0, sink.count("/v1/metrics"), "metrics stays disabled by otel.enabled")
}

func TestSetupYmlEndpointWinsOverEnv(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	clearEnv(t)
	envSink := newSink(t)
	ymlSink := newSink(t)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", envSink.URL)

	shutdown, err := Setup(context.Background(), "runner", config.Otel{Endpoint: ymlSink.URL})
	require.NoError(t, err)
	emit(t, shutdown)

	assert.Positive(t, ymlSink.count("/v1/traces"), "yml shared endpoint wins over env")
	assert.Equal(t, 0, envSink.total(), "env endpoint must be ignored when yml is set")
}

func TestSetupYmlHeadersOverrideEnv(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	clearEnv(t)
	sink := newSink(t)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", sink.URL)
	t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "authorization=Bearer env-token")

	// Shared yml header overrides the env header for both signals.
	shutdown, err := Setup(context.Background(), "runner", config.Otel{
		Headers: map[string]string{"authorization": "Bearer yml-token"},
	})
	require.NoError(t, err)
	emit(t, shutdown)

	require.Positive(t, sink.count("/v1/traces"))
	require.Positive(t, sink.count("/v1/metrics"))
	sink.mu.Lock()
	defer sink.mu.Unlock()
	assert.Equal(t, "Bearer yml-token", sink.auth["/v1/traces"], "yml header wins over env for traces")
	assert.Equal(t, "Bearer yml-token", sink.auth["/v1/metrics"], "shared yml header reaches metrics")
}

func TestSetupSdkDisabledEnv(t *testing.T) {
	resetGlobals(t)
	unsetConfigFile(t)
	clearEnv(t)
	sink := newSink(t)
	t.Setenv("OTEL_SDK_DISABLED", "true")

	shutdown, err := Setup(context.Background(), "runner", config.Otel{Endpoint: sink.URL})
	require.NoError(t, err)
	emit(t, shutdown)
	assert.Equal(t, 0, sink.total(), "OTEL_SDK_DISABLED=true disables both signals")
}
