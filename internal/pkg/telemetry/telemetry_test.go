// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package telemetry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitea.dev/actionslib/pkg/model"
	runnerv1 "gitea.dev/actionslib/runner/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/protobuf/types/known/structpb"
)

func resetProviders(t *testing.T) {
	t.Cleanup(func() {
		otel.SetTracerProvider(tracenoop.NewTracerProvider())
		otel.SetMeterProvider(metricnoop.NewMeterProvider())
	})
}

func TestSetup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()
	for _, tc := range []struct {
		env             map[string]string
		traces, metrics bool
		err             string
	}{
		{},
		{env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": server.URL}, traces: true, metrics: true},
		{env: map[string]string{"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": server.URL}, metrics: true},
		{env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": server.URL, "OTEL_TRACES_EXPORTER": "NONE"}, metrics: true},
		{env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": server.URL, "OTEL_SDK_DISABLED": "TRUE"}},
		{env: map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": server.URL, "OTEL_EXPORTER_OTLP_PROTOCOL": "grpc"}, err: "OTEL_EXPORTER_OTLP_PROTOCOL=grpc is not supported, only http/protobuf is"},
	} {
		t.Run(fmt.Sprint(tc.env), func(t *testing.T) {
			resetProviders(t)
			for _, key := range []string{
				"OTEL_SDK_DISABLED", "OTEL_TRACES_EXPORTER", "OTEL_METRICS_EXPORTER",
				"OTEL_EXPORTER_OTLP_ENDPOINT", "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "OTEL_EXPORTER_OTLP_METRICS_ENDPOINT",
				"OTEL_EXPORTER_OTLP_PROTOCOL", "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL", "OTEL_EXPORTER_OTLP_METRICS_PROTOCOL",
			} {
				t.Setenv(key, tc.env[key])
			}
			shutdown, err := Setup(t.Context(), "uuid", "runner")
			if tc.err != "" {
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
			_, traces := otel.GetTracerProvider().(*sdktrace.TracerProvider)
			_, metrics := otel.GetMeterProvider().(*sdkmetric.MeterProvider)
			assert.Equal(t, tc.traces, traces)
			assert.Equal(t, tc.metrics, metrics)
			require.NoError(t, shutdown(t.Context()))
		})
	}
}

func TestJob(t *testing.T) {
	fields, err := structpb.NewStruct(map[string]any{
		"workflow":         "ci.yml",
		"repository":       "owner/repo",
		"repository_owner": "owner",
		"server_url":       "https://gitea.example/",
		"run_id":           "42",
	})
	require.NoError(t, err)
	task := &runnerv1.Task{Id: 7, Context: fields}
	t.Setenv("OTEL_SERVICE_NAME", "runner")
	t.Setenv("OTEL_EXPORTER_OTLP_HEADERS", "authorization=runner")

	envs := map[string]string{}
	ctx, endJob := StartJob(t.Context(), task)
	AddJobEnv(ctx, envs)
	stepCtx, endStep := StartStep(ctx, &model.Step{}, "Main")
	endStep(nil, nil)
	endJob(runnerv1.Result_RESULT_SUCCESS)
	assert.Equal(t, ctx, stepCtx)
	assert.Empty(t, envs)

	resetProviders(t)
	spans := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans)))
	reader := sdkmetric.NewManualReader()
	otel.SetMeterProvider(sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader)))

	envs = map[string]string{"OTEL_EXPORTER_OTLP_HEADERS": "authorization=job"}
	ctx, endJob = StartJob(t.Context(), task)
	AddJobEnv(ctx, envs)
	_, endStep = StartStep(ctx, &model.Step{Run: "make\nmake test", Number: 1}, "Main")
	endStep(errors.New("exit code 1"), &model.StepResult{Outcome: model.StepStatusFailure})
	endJob(runnerv1.Result_RESULT_FAILURE)

	assert.Equal(t, "authorization=job", envs["OTEL_EXPORTER_OTLP_HEADERS"])
	assert.NotContains(t, envs, "OTEL_SERVICE_NAME")
	jobSpan := trace.SpanContextFromContext(ctx)
	assert.Equal(t, "00-"+jobSpan.TraceID().String()+"-"+jobSpan.SpanID().String()+"-01", envs["TRACEPARENT"])

	ended := spans.Ended()
	require.Len(t, ended, 2)
	step, job := ended[0], ended[1]
	assert.Equal(t, "Run make", step.Name())
	assert.Equal(t, jobSpan.SpanID(), step.Parent().SpanID())
	assert.Equal(t, codes.Error, step.Status().Code)
	assert.Subset(t, step.Attributes(), []any{
		semconv.CICDPipelineTaskRunID("7.1"),
		semconv.CICDPipelineTaskRunURLFull("https://gitea.example/owner/repo/actions/runs/42"),
		semconv.CICDPipelineTaskRunResultKey.String("failure"),
		semconv.ErrorTypeKey.String("failure"),
	})
	assert.Equal(t, "RUN ci.yml", job.Name())
	assert.Equal(t, trace.SpanKindServer, job.SpanKind())
	assert.Subset(t, job.Attributes(), []any{
		semconv.VCSOwnerName("owner"),
		semconv.VCSRepositoryName("repo"),
		semconv.CICDPipelineResultKey.String("failure"),
	})

	var data metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &data))
	aggregations := map[string]metricdata.Aggregation{}
	for _, m := range data.ScopeMetrics[0].Metrics {
		aggregations[m.Name] = m.Data
	}
	assert.Zero(t, aggregations["cicd.pipeline.run.active"].(metricdata.Sum[int64]).DataPoints[0].Value)
	assert.Subset(t, aggregations["cicd.pipeline.run.duration"].(metricdata.Histogram[float64]).DataPoints[0].Attributes.ToSlice(), []any{
		semconv.ErrorTypeKey.String("failure"),
		semconv.VCSOwnerName("owner"),
		semconv.VCSRepositoryName("repo"),
	})
}

func TestResults(t *testing.T) {
	cancelled, cancel := context.WithCancel(t.Context())
	cancel()
	expired, expire := context.WithDeadline(t.Context(), time.Time{})
	defer expire()
	for want, tc := range map[string]struct {
		ctx    context.Context
		err    error
		result *model.StepResult
	}{
		"timeout":      {ctx: t.Context(), err: fmt.Errorf("exec: %w", context.DeadlineExceeded)},
		"cancellation": {ctx: cancelled},
		"failure":      {ctx: t.Context(), result: &model.StepResult{Outcome: model.StepStatusFailure, Conclusion: model.StepStatusSuccess}},
		"success":      {ctx: t.Context(), result: &model.StepResult{}},
	} {
		assert.Equal(t, want, stepResult(tc.ctx, tc.err, tc.result))
	}
	assert.Equal(t, "timeout", jobResult(expired, runnerv1.Result_RESULT_FAILURE))
	assert.Equal(t, "cancellation", jobResult(t.Context(), runnerv1.Result_RESULT_CANCELLED))
}
