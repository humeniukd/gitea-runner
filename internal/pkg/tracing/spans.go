// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package tracing

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

// StartPipeline opens the pipeline-run span for one fetched task. The workflow
// name is unknown yet, so the caller refines the name to "RUN {workflow}" later.
func StartPipeline(ctx context.Context, task *runnerv1.Task, uuid, name string) (context.Context, trace.Span) {
	return otel.Tracer(scopeName).Start(ctx, "RUN", trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			semconv.CICDPipelineRunID(task.Context.Fields["run_id"].GetStringValue()),
			semconv.CICDPipelineTaskRunID(strconv.FormatInt(task.Id, 10)),
			semconv.CICDPipelineActionNameRun,
			semconv.CICDWorkerID(uuid),
			semconv.CICDWorkerName(name),
		))
}

// StartJob opens the span for the single job the task runs, as a child of the
// pipeline-run span, so the job name is visible and the step spans nest under
// it. jobName is the job's name:, falling back to its workflow id when unset.
func StartJob(ctx context.Context, jobID, jobName string) (context.Context, trace.Span) {
	name := jobName
	if name == "" {
		name = jobID
	}
	return otel.Tracer(scopeName).Start(ctx, name, trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("gitea.job.id", jobID),
			attribute.String("gitea.job.name", jobName),
		))
}

// StartStep opens the task-run span for one step's main stage, named as the
// runner logs it (model.Step.String()), else "step N".
func StartStep(ctx context.Context, step *model.Step, taskID int64) (context.Context, trace.Span) {
	name := step.String()
	if name == "" {
		name = fmt.Sprintf("step %d", step.Number)
	}
	attrs := []attribute.KeyValue{
		semconv.CICDPipelineTaskName(name),
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

// SetPipelineResult marks the pipeline-run span with the run's outcome.
func SetPipelineResult(span trace.Span, result string) {
	setResult(span, semconv.CICDPipelineResultKey, result)
}

// SetJobResult marks the job span with the run's outcome.
func SetJobResult(span trace.Span, result string) {
	span.SetAttributes(attribute.String("gitea.job.result", result))
	if result != "success" {
		span.SetStatus(codes.Error, result)
	}
}

// SetTaskResult marks a task-run span with the step's outcome.
func SetTaskResult(span trace.Span, result string) {
	setResult(span, semconv.CICDPipelineTaskRunResultKey, result)
}

// SetTaskFailed marks a task-run span failed, recording the error and using the
// last command output as the status when present (clearer than the exit error).
func SetTaskFailed(span trace.Span, err error, lastOutput string) {
	span.SetAttributes(semconv.CICDPipelineTaskRunResultKey.String("failure"))
	span.RecordError(err)
	span.SetStatus(codes.Error, stepFailureMessage(err, lastOutput))
}

// stepFailureMessage combines the step error with the last command output line,
// falling back to the error alone when no output was captured.
func stepFailureMessage(err error, lastOutput string) string {
	out := strings.TrimSpace(lastOutput)
	switch {
	case err == nil:
		return out
	case out == "":
		return err.Error()
	default:
		return err.Error() + ": " + out
	}
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
