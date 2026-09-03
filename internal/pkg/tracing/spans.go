// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package tracing

import (
	"context"
	"fmt"
	"strconv"

	"gitea.dev/actionslib/pkg/model"
	runnerv1 "gitea.dev/actionslib/runner/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
)

// scopeName is the instrumentation scope for the runner's own spans. The
// tracer is resolved per call, not cached, so it always follows the current
// global provider (e.g. after a re-Setup in tests).
const scopeName = "gitea.com/gitea/runner"

// StartJob opens the pipeline-run span for one fetched task. The workflow name
// is unknown yet, so the caller refines the name to "RUN {workflow}" later.
func StartJob(ctx context.Context, task *runnerv1.Task, uuid, name string) (context.Context, trace.Span) {
	return otel.Tracer(scopeName).Start(ctx, "RUN", trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			semconv.CICDPipelineRunID(task.Context.Fields["run_id"].GetStringValue()),
			semconv.CICDPipelineTaskRunID(strconv.FormatInt(task.Id, 10)),
			semconv.CICDPipelineActionNameRun,
			semconv.CICDWorkerID(uuid),
			semconv.CICDWorkerName(name),
		))
}

// StartStep opens the task-run span for one workflow step's main stage.
func StartStep(ctx context.Context, step *model.Step, taskID int64) (context.Context, trace.Span) {
	name := step.Name
	if name == "" {
		name = fmt.Sprintf("step %d", step.Number)
	}
	attrs := []attribute.KeyValue{
		semconv.CICDPipelineTaskName(step.Name),
		semconv.CICDPipelineTaskRunID(fmt.Sprintf("%d.%d", taskID, step.Number)),
		semconv.CICDPipelineTaskTypeBuild,
	}
	if step.Uses != "" {
		attrs = append(attrs, attribute.String("gitea.step.uses", step.Uses))
	}
	if step.Run != "" {
		attrs = append(attrs, attribute.String("gitea.step.run", step.Run))
	}
	return otel.Tracer(scopeName).Start(ctx, name, trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(attrs...))
}

// SetJobResult marks the pipeline-run span with the run's outcome.
func SetJobResult(span trace.Span, result string) {
	setResult(span, semconv.CICDPipelineResultKey, result)
}

// SetTaskResult marks a task-run span with the step's outcome.
func SetTaskResult(span trace.Span, result string) {
	setResult(span, semconv.CICDPipelineTaskRunResultKey, result)
}

func setResult(span trace.Span, key attribute.Key, result string) {
	span.SetAttributes(key.String(result))
	if result != "success" {
		span.SetStatus(codes.Error, result)
	}
}

// ResultToSemconv maps the runnerv1.Result enum onto the semconv result values.
func ResultToSemconv(result runnerv1.Result) string {
	switch result {
	case runnerv1.Result_RESULT_SUCCESS:
		return "success"
	case runnerv1.Result_RESULT_FAILURE:
		return "failure"
	case runnerv1.Result_RESULT_CANCELLED:
		return "cancellation"
	case runnerv1.Result_RESULT_SKIPPED:
		return "skip"
	default:
		return "error"
	}
}
