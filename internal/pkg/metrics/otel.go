// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// OpenTelemetry instruments for the runner's own metrics, named to the CI/CD
// semantic conventions. The Prometheus endpoint is unaffected; these export
// through the OTLP pipeline configured under otel.metrics.
package metrics

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

type Telemetry struct {
	meter metric.Meter

	PipelineRunDuration metric.Float64Histogram   // cicd.pipeline.run.duration, unit s
	PipelineRunActive   metric.Int64UpDownCounter // cicd.pipeline.run.active, {run}
	WorkerCount         metric.Int64UpDownCounter // cicd.worker.count, {worker}
	SystemErrors        metric.Int64Counter       // cicd.system.errors, {error}
	PollFetch           metric.Int64Counter       // gitea.runner.poll.fetch, {fetch}
	PollFetchDuration   metric.Float64Histogram   // gitea.runner.poll.fetch.duration, unit s
	PollBackoff         metric.Float64Gauge       // gitea.runner.poll.backoff, unit s
}

func NewTelemetry() *Telemetry {
	meter := otel.Meter("gitea.com/gitea/runner/internal/pkg/metrics")
	// Boundaries span the runner.job timeout range, 5m default to hours.
	runDuration, err := meter.Float64Histogram(
		"cicd.pipeline.run.duration", metric.WithUnit("s"),
		metric.WithDescription("Duration of a pipeline run grouped by pipeline, state and result."),
		metric.WithExplicitBucketBoundaries(0.5, 1, 2.5, 5, 10, 30, 60, 300, 600, 1800, 3600),
	)
	handleErr(err)

	runActive, err := meter.Int64UpDownCounter(
		"cicd.pipeline.run.active", metric.WithUnit("{run}"),
		metric.WithDescription("The number of pipeline runs currently active in the system by state."),
	)
	handleErr(err)

	workerCount, err := meter.Int64UpDownCounter(
		"cicd.worker.count", metric.WithUnit("{worker}"),
		metric.WithDescription("The number of workers on the CI/CD system by state."),
	)
	handleErr(err)

	systemErrors, err := meter.Int64Counter(
		"cicd.system.errors", metric.WithUnit("{error}"),
		metric.WithDescription("The number of errors in a component of the CI/CD system."),
	)
	handleErr(err)

	pollFetch, err := meter.Int64Counter(
		"gitea.runner.poll.fetch", metric.WithUnit("{fetch}"),
		metric.WithDescription("Total FetchTask RPCs by result (task, empty, error)."),
	)
	handleErr(err)

	pollFetchDuration, err := meter.Float64Histogram(
		"gitea.runner.poll.fetch.duration", metric.WithUnit("s"),
		metric.WithDescription("Latency of FetchTask RPCs, excluding expected long-poll timeouts."),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	handleErr(err)

	pollBackoff, err := meter.Float64Gauge(
		"gitea.runner.poll.backoff", metric.WithUnit("s"),
		metric.WithDescription("Last observed polling backoff interval in seconds."),
	)
	handleErr(err)

	return &Telemetry{
		meter:               meter,
		PipelineRunDuration: runDuration,
		PipelineRunActive:   runActive,
		WorkerCount:         workerCount,
		SystemErrors:        systemErrors,
		PollFetch:           pollFetch,
		PollFetchDuration:   pollFetchDuration,
		PollBackoff:         pollBackoff,
	}
}

func handleErr(err error) {
	if err != nil {
		log.WithError(err).Warn("otel metric instrument unavailable")
	}
}

func (t *Telemetry) RecordJobStart(ctx context.Context, pipeline string) {
	attrs := metric.WithAttributes(
		semconv.CICDPipelineName(pipeline),
		semconv.CICDPipelineRunStateExecuting,
	)
	t.PipelineRunActive.Add(ctx, 1, attrs)
	t.WorkerCount.Add(ctx, 1, metric.WithAttributes(semconv.CICDWorkerStateBusy))
}

func (t *Telemetry) RecordJobEnd(ctx context.Context, pipeline, result, errType string) {
	attrs := metric.WithAttributes(
		semconv.CICDPipelineName(pipeline),
		semconv.CICDPipelineRunStateExecuting,
		semconv.CICDPipelineResultKey.String(result),
	)
	t.PipelineRunActive.Add(ctx, -1, attrs)
	t.WorkerCount.Add(ctx, -1, metric.WithAttributes(semconv.CICDWorkerStateAvailable))
	if result != "success" && errType != "" {
		t.SystemErrors.Add(ctx, 1, metric.WithAttributes(
			semconv.CICDSystemComponent("runner"),
			attribute.String("error.type", errType),
		))
	}
}

func (t *Telemetry) RecordPipelineRunDuration(ctx context.Context, pipeline, result string, seconds float64) {
	t.PipelineRunDuration.Record(ctx, seconds, metric.WithAttributes(
		semconv.CICDPipelineName(pipeline),
		semconv.CICDPipelineRunStateExecuting,
		semconv.CICDPipelineResultKey.String(result),
	))
}

func (t *Telemetry) RecordPollFetch(ctx context.Context, result string) {
	t.PollFetch.Add(ctx, 1, metric.WithAttributes(
		attribute.String("gitea.poll.result", result),
	))
}

func (t *Telemetry) RecordPollFetchDuration(ctx context.Context, seconds float64) {
	t.PollFetchDuration.Record(ctx, seconds)
}

func (t *Telemetry) RecordPollBackoff(ctx context.Context, seconds float64) {
	t.PollBackoff.Record(ctx, seconds)
}
