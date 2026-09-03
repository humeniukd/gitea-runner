// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"gitea.com/gitea/runner/internal/pkg/config"
	"gitea.com/gitea/runner/internal/pkg/tracing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	noop "go.opentelemetry.io/otel/trace/noop"
)

// TestTracing pins the propagation contract: the runner exports spans of a real
// job to the configured endpoint, and the job container sees TRACEPARENT (the
// job span) plus the runner's OTEL_* environment.
func TestTracing(t *testing.T) {
	if skipReason != "" {
		t.Skip(skipReason)
	}

	var hits atomic.Int64
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/traces" && r.Method == http.MethodPost {
			hits.Add(1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(collector.Close)

	// Setup picks the endpoint up from the env, exactly as a real deployment
	// configures the daemon.
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", collector.URL)
	shutdown, err := tracing.Setup(context.Background(), "e2e-runner", config.Otel{})
	if err != nil {
		t.Fatalf("setup tracing: %v", err)
	}
	// Flush the batch span processor before asserting on the collector.
	shutdownTracing := func() {
		sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = shutdown(sctx)
		otel.SetTracerProvider(noop.NewTracerProvider())
	}
	t.Cleanup(shutdownTracing)

	sharedPoller := startRunner(t, "", "ubuntu-latest", runnerOptions{})
	t.Cleanup(func() {
		select {
		case <-sharedPoller.Done():
			t.Error("shared runner stopped")
		default:
		}
	})

	api, repo := newScenario(t)
	pushWorkflow(t, api, repo, "tracing.yml")

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)

	traceparent := strings.TrimPrefix(commandRow(logs, "traceparent-env="), "traceparent-env=")
	cctx := propagation.TraceContext{}.Extract(context.Background(), propagation.MapCarrier{"traceparent": traceparent})
	if !trace.SpanContextFromContext(cctx).IsValid() {
		t.Fatalf("job echoed an invalid TRACEPARENT %q; logs:\n%s", traceparent, logs)
	}
	if endpoint := commandRow(logs, "otlp-endpoint="); !strings.HasPrefix(endpoint, "otlp-endpoint=http://") {
		t.Errorf("the job did not see OTEL_EXPORTER_OTLP_ENDPOINT, got %q", endpoint)
	}
	if vars := strings.TrimPrefix(commandRow(logs, "otel-vars="), "otel-vars="); vars == "" || vars == "0" {
		t.Errorf("the job saw no OTEL_* variables, got %q", vars)
	}

	// Flush before asserting: the batch processor buffers spans until shutdown.
	// A job container only sees a traceparent when the runner itself is inside
	// a real (sampled) span, whose export this check observes.
	sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdown(sctx)
	otel.SetTracerProvider(noop.NewTracerProvider())
	if n := hits.Load(); n == 0 {
		t.Error("collector received no trace export")
	}
}
