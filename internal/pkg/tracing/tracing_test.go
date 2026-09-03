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
	"go.opentelemetry.io/contrib/otelconf"
	"go.opentelemetry.io/otel"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

// unsetConfigFile clears OTEL_CONFIG_FILE, which otelconf treats as a file to read even when empty.
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

	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "")
	envs := JobEnv(ctx, config.Otel{})
	assert.Equal(t, "1", envs["OTEL_TEST_PROP"])
	assert.NotContains(t, envs, "OTHER_VAR")
	assert.Equal(t, "deployment.environment.name=test", envs["OTEL_RESOURCE_ATTRIBUTES"], "default persists")

	cctx := propagation.TraceContext{}.Extract(context.Background(), propagation.MapCarrier{"traceparent": envs["TRACEPARENT"]})
	assert.True(t, trace.SpanContextFromContext(cctx).IsValid())

	assert.Equal(t, "authorization=Bearer secret", envs["OTEL_EXPORTER_OTLP_HEADERS"])
	assert.NotContains(t, JobEnv(context.Background(), config.Otel{}), "TRACEPARENT")
}

func TestJobEnvConfigToEnv(t *testing.T) {
	unsetConfigFile(t)
	clearEnv(t)
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "team=ci")
	ratio := 0.25
	cfg := config.Otel{
		Endpoint:           "http://collector:4318",
		Protocol:           "grpc",
		Headers:            map[string]string{"authorization": "Bearer shared"},
		ResourceAttributes: map[string]string{"deployment.environment.name": "production"},
		Traces:             config.OtelTraces{Endpoint: "http://collector:4318/v1/traces", SamplerRatio: &ratio},
		Metrics:            config.OtelMetrics{Headers: map[string]string{"x-scope-orgid": "ci"}},
	}
	envs := JobEnv(context.Background(), cfg)
	assert.Equal(t, "http://collector:4318", envs["OTEL_EXPORTER_OTLP_ENDPOINT"])
	assert.Equal(t, "grpc", envs["OTEL_EXPORTER_OTLP_PROTOCOL"])
	assert.Equal(t, "authorization=Bearer shared", envs["OTEL_EXPORTER_OTLP_HEADERS"])
	assert.Equal(t, "http://collector:4318/v1/traces", envs["OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"])
	assert.Equal(t, "parentbased_traceidratio", envs["OTEL_TRACES_SAMPLER"])
	assert.Equal(t, "0.25", envs["OTEL_TRACES_SAMPLER_ARG"])
	assert.Equal(t, "x-scope-orgid=ci", envs["OTEL_EXPORTER_OTLP_METRICS_HEADERS"])
	assert.Equal(t, "deployment.environment.name=production,team=ci", envs["OTEL_RESOURCE_ATTRIBUTES"])
}

func TestJobEnvDisabledSignal(t *testing.T) {
	unsetConfigFile(t)
	clearEnv(t)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://collector:4318")
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "http://collector:4318/v1/traces")

	envs := JobEnv(context.Background(), config.Otel{Traces: config.OtelTraces{Enabled: new(false)}})
	assert.Equal(t, "none", envs["OTEL_TRACES_EXPORTER"], "disabled traces turned off in-job")
	assert.NotContains(t, envs, "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "disabled signal's env dropped")
	assert.Equal(t, "http://collector:4318", envs["OTEL_EXPORTER_OTLP_ENDPOINT"], "enabled metrics keep the shared endpoint")
	assert.NotContains(t, envs, "OTEL_METRICS_EXPORTER")
}

func TestResourceAttributeResolution(t *testing.T) {
	assert.Equal(t, map[string]string{"deployment.environment.name": "test"},
		mergedResourceAttributes("", nil))
	assert.Equal(t, map[string]string{"deployment.environment.name": "test", "team": "ci"},
		mergedResourceAttributes("team=ci", nil))
	assert.Equal(t, map[string]string{"deployment.environment.name": "prod", "team": "ci"},
		mergedResourceAttributes("deployment.environment.name=staging,team=ci", map[string]string{"deployment.environment.name": "prod"}))

	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "")
	deploymentEnv := func(attrs []otelconf.AttributeNameValue) any {
		for _, a := range attrs {
			if a.Name == "deployment.environment.name" {
				return a.Value
			}
		}
		return nil
	}
	assert.Equal(t, "test", deploymentEnv(resourceAttributes("runner", nil)))
	assert.Equal(t, "prod", deploymentEnv(resourceAttributes("runner", map[string]string{"deployment.environment.name": "prod"})))
}

func TestSetupNoopWithoutEndpoint(t *testing.T) {
	unsetConfigFile(t)
	clearEnv(t)

	shutdown, err := Setup(context.Background(), "runner", config.Otel{})
	require.NoError(t, err)
	require.NoError(t, shutdown(context.Background()))
}

func TestSetupSignalGating(t *testing.T) {
	cases := []struct {
		name  string
		env   map[string]string
		cfg   func(url string) config.Otel
		check func(t *testing.T, s *otlpSink)
	}{
		{
			name: "shared endpoint enables both signals",
			cfg:  func(url string) config.Otel { return config.Otel{Endpoint: url} },
			check: func(t *testing.T, s *otlpSink) {
				assert.Positive(t, s.count("/v1/traces"))
				assert.Positive(t, s.count("/v1/metrics"))
			},
		},
		{
			name: "signal endpoint used as-is, no path appended",
			cfg: func(url string) config.Otel {
				return config.Otel{Traces: config.OtelTraces{Endpoint: url + "/custom/traces"}}
			},
			check: func(t *testing.T, s *otlpSink) {
				assert.Positive(t, s.count("/custom/traces"))
				assert.Equal(t, 0, s.count("/custom/traces/v1/traces"))
			},
		},
		{
			name:  "otel.enabled=false disables both",
			cfg:   func(url string) config.Otel { return config.Otel{Enabled: new(false), Endpoint: url} },
			check: func(t *testing.T, s *otlpSink) { assert.Equal(t, 0, s.total()) },
		},
		{
			name: "traces.enabled=true overrides otel.enabled=false",
			cfg: func(url string) config.Otel {
				return config.Otel{Enabled: new(false), Endpoint: url, Traces: config.OtelTraces{Enabled: new(true)}}
			},
			check: func(t *testing.T, s *otlpSink) {
				assert.Positive(t, s.count("/v1/traces"))
				assert.Equal(t, 0, s.count("/v1/metrics"))
			},
		},
		{
			name:  "OTEL_SDK_DISABLED=true disables both",
			env:   map[string]string{"OTEL_SDK_DISABLED": "true"},
			cfg:   func(url string) config.Otel { return config.Otel{Endpoint: url} },
			check: func(t *testing.T, s *otlpSink) { assert.Equal(t, 0, s.total()) },
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetGlobals(t)
			unsetConfigFile(t)
			clearEnv(t)
			sink := newSink(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			shutdown, err := Setup(context.Background(), "runner", tc.cfg(sink.URL))
			require.NoError(t, err)
			emit(t, shutdown)
			tc.check(t, sink)
		})
	}
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
