// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gitea.com/gitea/runner/internal/pkg/config"
	"gitea.com/gitea/runner/internal/pkg/tracing"

	contribruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	metricsv1 "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	"google.golang.org/protobuf/proto"
)

// otelCollector captures OTLP/HTTP exports: it counts trace posts and decodes
// the exported metric names and units.
type otelCollector struct {
	mu         sync.Mutex
	names      map[string]bool
	units      map[string]string
	traceHits  atomic.Int64
	metricHits atomic.Int64
	server     *httptest.Server
}

func startOtelCollector(t *testing.T) *otelCollector {
	t.Helper()
	c := &otelCollector{names: map[string]bool{}, units: map[string]string{}}
	c.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/v1/traces"):
			c.traceHits.Add(1)
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/v1/metrics"):
			if body, err := io.ReadAll(r.Body); err == nil {
				req := &metricsv1.ExportMetricsServiceRequest{}
				if proto.Unmarshal(body, req) == nil {
					c.mu.Lock()
					for _, rm := range req.ResourceMetrics {
						for _, sm := range rm.ScopeMetrics {
							for _, m := range sm.Metrics {
								c.names[m.Name] = true
								c.units[m.Name] = m.Unit
							}
						}
					}
					c.mu.Unlock()
				}
			}
			c.metricHits.Add(1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(c.server.Close)
	return c
}

func (c *otelCollector) hasMetric(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.names[name]
}

func (c *otelCollector) metricUnit(name string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.units[name]
}

func (c *otelCollector) metricNames() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	names := make([]string, 0, len(c.names))
	for name := range c.names {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// runOtelJob installs the SDK from otelCfg, starts Go runtime metrics, runs one
// real job, flushes, and returns the job log.
func runOtelJob(t *testing.T, otelCfg config.Otel) string {
	t.Helper()

	shutdown, err := tracing.Setup(context.Background(), "e2e-runner", otelCfg)
	if err != nil {
		t.Fatalf("setup otel: %v", err)
	}
	t.Cleanup(func() {
		otel.SetTracerProvider(tracenoop.NewTracerProvider())
		otel.SetMeterProvider(metricnoop.NewMeterProvider())
	})

	// Go runtime metrics ride on the same global meter provider, exactly as the
	// daemon starts them.
	if err := contribruntime.Start(contribruntime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		t.Fatalf("start runtime metrics: %v", err)
	}

	api, repo, _ := startIsolatedScenario(t, "tracing.yml", "ubuntu-latest", runnerOptions{})
	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)
	logs := runLogs(t, api, repo, wfRun.ID)

	// Flush both pipelines before asserting: batch/periodic processors buffer
	// until shutdown.
	sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := shutdown(sctx); err != nil {
		t.Fatalf("flush otel: %v", err)
	}
	return logs
}

// assertOtelExports checks the job propagated a valid TRACEPARENT and that the
// collector received the trace and the full metric set with correct units.
func assertOtelExports(t *testing.T, collector *otelCollector, logs string) {
	t.Helper()

	traceparent := strings.TrimPrefix(commandRow(logs, "traceparent-env="), "traceparent-env=")
	cctx := propagation.TraceContext{}.Extract(context.Background(), propagation.MapCarrier{"traceparent": traceparent})
	if !trace.SpanContextFromContext(cctx).IsValid() {
		t.Fatalf("job echoed an invalid TRACEPARENT %q; logs:\n%s", traceparent, logs)
	}
	if vars := strings.TrimPrefix(commandRow(logs, "otel-vars="), "otel-vars="); vars == "" || vars == "0" {
		t.Errorf("the job saw no OTEL_* variables, got %q", vars)
	}
	if collector.traceHits.Load() == 0 {
		t.Error("collector received no trace export")
	}
	for _, name := range []string{
		"cicd.pipeline.run.duration",
		"cicd.pipeline.run.active",
		"cicd.worker.count",
		"gitea.runner.poll.fetch",
		"gitea.runner.poll.fetch.duration",
		"gitea.runner.poll.backoff",
		"rpc.client.duration",
		"go.goroutine.count",
		"go.memory.used",
	} {
		if !collector.hasMetric(name) {
			t.Errorf("collector did not receive %q; got:\n%v", name, collector.metricNames())
		}
	}
	if u := collector.metricUnit("cicd.pipeline.run.duration"); u != "s" {
		t.Errorf("cicd.pipeline.run.duration unit is %q, want s", u)
	}
	if u := collector.metricUnit("cicd.pipeline.run.active"); u != "{run}" {
		t.Errorf("cicd.pipeline.run.active unit is %q, want {run}", u)
	}
}

// TestOtel pins the full OTel contract in one job run configured through the
// standard per-signal env vars (complete URLs including the signal path): the
// runner exports its job span, the job container sees TRACEPARENT plus the
// runner's OTEL_* environment, and the runner exports the CI/CD semconv
// metrics, poll and RPC instruments, and Go runtime metrics.
func TestOtel(t *testing.T) {
	if skipReason != "" {
		t.Skip(skipReason)
	}

	collector := startOtelCollector(t)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", collector.server.URL+"/v1/traces")
	t.Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", collector.server.URL+"/v1/metrics")

	logs := runOtelJob(t, config.Otel{})
	assertOtelExports(t, collector, logs)
}

// TestOtelSharedEndpoint proves the top-level otel.endpoint enables and points
// both signals from a single setting, with the OTLP signal path appended to the
// shared base.
func TestOtelSharedEndpoint(t *testing.T) {
	if skipReason != "" {
		t.Skip(skipReason)
	}

	collector := startOtelCollector(t)
	// No per-signal endpoints anywhere; only the shared base URL.
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "")

	logs := runOtelJob(t, config.Otel{Endpoint: collector.server.URL})
	assertOtelExports(t, collector, logs)
}
