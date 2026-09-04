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

func TestLogHook(t *testing.T) {
	tp := sdktrace.NewTracerProvider()
	ctx, span := tp.Tracer("test").Start(context.Background(), "op")
	defer span.End()
	sc := span.SpanContext()

	entry := logrus.NewEntry(logrus.New()).WithContext(ctx)
	require.NoError(t, logHook{}.Fire(entry))
	assert.Equal(t, sc.TraceID().String(), entry.Data["trace_id"])
	assert.Equal(t, sc.SpanID().String(), entry.Data["span_id"])

	for _, e := range []*logrus.Entry{
		logrus.NewEntry(logrus.New()),
		logrus.NewEntry(logrus.New()).WithContext(context.Background()),
	} {
		require.NoError(t, logHook{}.Fire(e))
		assert.NotContains(t, e.Data, "trace_id")
		assert.NotContains(t, e.Data, "span_id")
	}
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

func TestSetPipelineName(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr)))
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	_, span := otel.Tracer("test").Start(context.Background(), "RUN")
	SetPipelineName(span, "Tests", "build", "Build and test")
	SetPipelineResult(span, "success")
	span.End()

	_, fallback := otel.Tracer("test").Start(context.Background(), "RUN")
	SetPipelineName(fallback, "Tests", "deploy", "")
	SetPipelineResult(fallback, "skip")
	fallback.End()

	byName := make(map[string]sdktrace.ReadOnlySpan)
	for _, s := range sr.Ended() {
		byName[s.Name()] = s
	}
	require.Contains(t, byName, "RUN Tests Build and test", "combined workflow and job name")
	require.Contains(t, byName, "RUN Tests deploy", "job name falls back to the job id")

	id, ok := attrString(byName["RUN Tests Build and test"].Attributes(), "gitea.job.id")
	require.True(t, ok)
	assert.Equal(t, "build", id)
	assert.Equal(t, codes.Ok, byName["RUN Tests Build and test"].Status().Code, "success sets Ok status")
	assert.Equal(t, codes.Unset, byName["RUN Tests deploy"].Status().Code, "skip leaves status unset")
}

func attrString(attrs []attribute.KeyValue, key string) (string, bool) {
	for _, kv := range attrs {
		if string(kv.Key) == key {
			return kv.Value.AsString(), true
		}
	}
	return "", false
}

func TestVCSAttributes(t *testing.T) {
	byKey := func(attrs []attribute.KeyValue) map[string]string {
		m := make(map[string]string, len(attrs))
		for _, a := range attrs {
			m[string(a.Key)] = a.Value.AsString()
		}
		return m
	}

	got := byKey(vcsAttributes("octocat/hello", "octocat", "main", "abc123", "branch", "https://gitea.example.com/"))
	assert.Equal(t, "octocat", got["vcs.owner.name"])
	assert.Equal(t, "hello", got["vcs.repository.name"])
	assert.Equal(t, "https://gitea.example.com/octocat/hello", got["vcs.repository.url.full"])
	assert.Equal(t, "main", got["vcs.ref.head.name"])
	assert.Equal(t, "abc123", got["vcs.ref.head.revision"])
	assert.Equal(t, "branch", got["vcs.ref.head.type"])

	// owner falls back to the first path segment; empty inputs are omitted.
	got = byKey(vcsAttributes("octocat/hello", "", "", "", "tag", ""))
	assert.Equal(t, "octocat", got["vcs.owner.name"])
	assert.Equal(t, "hello", got["vcs.repository.name"])
	assert.Equal(t, "tag", got["vcs.ref.head.type"])
	assert.NotContains(t, got, "vcs.repository.url.full")
	assert.NotContains(t, got, "vcs.ref.head.name")
}
