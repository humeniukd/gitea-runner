// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package tracing

import (
	"context"
	"errors"
	"testing"

	"gitea.dev/actionslib/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestStepFailureMessage(t *testing.T) {
	const exitMsg = "Process completed with exit code 1."
	err := errors.New(exitMsg)
	assert.Equal(t, "Process completed with exit code 1.", stepFailureMessage(err, ""))
	assert.Equal(t, "Process completed with exit code 1.", stepFailureMessage(err, "   "))
	assert.Equal(t, "Process completed with exit code 1.: npm ERR! missing script: build",
		stepFailureMessage(err, "npm ERR! missing script: build\n"))
	assert.Equal(t, "boom", stepFailureMessage(nil, "boom"))
}

func TestSetTaskFailedUsesLastOutput(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	_, span := tp.Tracer("test").Start(context.Background(), "step")
	SetTaskFailed(span, errors.New("process completed with exit code 1"), "go: build failed")
	span.End()

	ended := sr.Ended()
	require.Len(t, ended, 1)
	got := ended[0]
	assert.Equal(t, codes.Error, got.Status().Code)
	assert.Equal(t, "process completed with exit code 1: go: build failed", got.Status().Description)

	var found bool
	for _, a := range got.Attributes() {
		if string(a.Key) == "cicd.pipeline.task.run.result" {
			found = true
			assert.Equal(t, "failure", a.Value.AsString())
		}
	}
	assert.True(t, found, "result attribute set")
	require.Len(t, got.Events(), 1, "error recorded as an event")
	assert.Equal(t, "exception", got.Events()[0].Name)
}

func TestLogHookAddsTraceAndSpanID(t *testing.T) {
	tp := sdktrace.NewTracerProvider()
	ctx, span := tp.Tracer("test").Start(context.Background(), "op")
	defer span.End()
	sc := span.SpanContext()

	entry := logrus.NewEntry(logrus.New()).WithContext(ctx)
	require.NoError(t, logHook{}.Fire(entry))
	assert.Equal(t, sc.TraceID().String(), entry.Data["trace_id"])
	assert.Equal(t, sc.SpanID().String(), entry.Data["span_id"])
}

func TestLogHookNoopWithoutSpan(t *testing.T) {
	entry := logrus.NewEntry(logrus.New())
	require.NoError(t, logHook{}.Fire(entry))
	assert.NotContains(t, entry.Data, "trace_id")

	entry = logrus.NewEntry(logrus.New()).WithContext(context.Background())
	require.NoError(t, logHook{}.Fire(entry))
	assert.NotContains(t, entry.Data, "trace_id")
	assert.NotContains(t, entry.Data, "span_id")
}

func TestStartStepName(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	cases := []struct {
		label string
		step  model.Step
		want  string
	}{
		{"named", model.Step{Name: "Build", Number: 0}, "Build"},
		{"uses falls back", model.Step{Uses: "actions/checkout@v4", Number: 1}, "actions/checkout@v4"},
		{"run falls back", model.Step{Run: "make test", Number: 2}, "make test"},
		{"empty is step N", model.Step{Number: 3}, "step 3"},
	}
	for _, tc := range cases {
		step := tc.step
		_, span := StartStep(context.Background(), &step, 42)
		span.End()
	}

	got := make(map[string]bool)
	for _, s := range sr.Ended() {
		got[s.Name()] = true
	}
	for _, tc := range cases {
		assert.Truef(t, got[tc.want], "%s: expected a span named %q", tc.label, tc.want)
	}
}

func TestStartJobName(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	_, named := StartJob(context.Background(), "build", "Build and test")
	named.End()
	_, unnamed := StartJob(context.Background(), "deploy", "")
	unnamed.End()

	byName := make(map[string]sdktrace.ReadOnlySpan)
	for _, s := range sr.Ended() {
		byName[s.Name()] = s
	}
	require.Contains(t, byName, "Build and test", "job span uses the job name")
	require.Contains(t, byName, "deploy", "job span falls back to the job id")

	id, ok := attrString(byName["Build and test"].Attributes(), "gitea.job.id")
	require.True(t, ok)
	assert.Equal(t, "build", id)
}

func attrString(attrs []attribute.KeyValue, key string) (string, bool) {
	for _, kv := range attrs {
		if string(kv.Key) == key {
			return kv.Value.AsString(), true
		}
	}
	return "", false
}
