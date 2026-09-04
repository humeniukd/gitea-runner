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

const scopeName = "gitea.com/gitea/runner"

func StartPipeline(ctx context.Context, task *runnerv1.Task, uuid, name, instanceURL string) (context.Context, trace.Span) {
	f := task.Context.Fields
	attrs := []attribute.KeyValue{
		semconv.CICDPipelineRunID(f["run_id"].GetStringValue()),
		semconv.CICDPipelineTaskRunID(strconv.FormatInt(task.Id, 10)),
		semconv.CICDPipelineActionNameRun,
		semconv.CICDWorkerID(uuid),
		semconv.CICDWorkerName(name),
	}
	attrs = append(attrs, vcsAttributes(
		f["repository"].GetStringValue(),
		f["repository_owner"].GetStringValue(),
		f["ref_name"].GetStringValue(),
		f["sha"].GetStringValue(),
		f["ref_type"].GetStringValue(),
		instanceURL,
	)...)
	return otel.Tracer(scopeName).Start(ctx, "RUN", trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attrs...))
}

func vcsAttributes(repository, owner, refName, sha, refType, instanceURL string) []attribute.KeyValue {
	name := repository
	if i := strings.LastIndex(repository, "/"); i >= 0 {
		if owner == "" {
			owner = repository[:i]
		}
		name = repository[i+1:]
	}
	var attrs []attribute.KeyValue
	if owner != "" {
		attrs = append(attrs, semconv.VCSOwnerName(owner))
	}
	if name != "" {
		attrs = append(attrs, semconv.VCSRepositoryName(name))
	}
	if repository != "" && instanceURL != "" {
		attrs = append(attrs, semconv.VCSRepositoryURLFull(strings.TrimRight(instanceURL, "/")+"/"+repository))
	}
	if refName != "" {
		attrs = append(attrs, semconv.VCSRefHeadName(refName))
	}
	if sha != "" {
		attrs = append(attrs, semconv.VCSRefHeadRevision(sha))
	}
	switch refType {
	case "branch":
		attrs = append(attrs, semconv.VCSRefHeadTypeBranch)
	case "tag":
		attrs = append(attrs, semconv.VCSRefHeadTypeTag)
	}
	return attrs
}

func SetPipelineName(span trace.Span, workflowName, jobID, jobName string) {
	job := jobName
	if job == "" {
		job = jobID
	}
	parts := []string{"RUN"}
	if workflowName != "" {
		parts = append(parts, workflowName)
	}
	if job != "" {
		parts = append(parts, job)
	}
	span.SetName(strings.Join(parts, " "))
	span.SetAttributes(
		semconv.CICDPipelineName(workflowName),
		attribute.String("gitea.job.id", jobID),
		attribute.String("gitea.job.name", jobName),
	)
}

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

func SetPipelineResult(span trace.Span, result string) {
	setResult(span, semconv.CICDPipelineResultKey, result)
}

func SetTaskResult(span trace.Span, result string) {
	setResult(span, semconv.CICDPipelineTaskRunResultKey, result)
}

func SetTaskFailed(span trace.Span, err error, lastOutput string) {
	span.SetAttributes(semconv.CICDPipelineTaskRunResultKey.String("failure"))
	span.RecordError(err)
	span.SetStatus(codes.Error, stepFailureMessage(err, lastOutput))
}

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
	switch result {
	case "success":
		span.SetStatus(codes.Ok, "")
	case "skip":
		span.SetStatus(codes.Unset, "skipped")
	default:
		span.SetStatus(codes.Error, result)
	}
}

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
