// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package telemetry

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitea.com/gitea/runner/internal/pkg/ver"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"gitea.dev/actionslib/pkg/model"
	runnerv1 "gitea.dev/actionslib/runner/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/semconv/v1.43.0/cicdconv"
	"go.opentelemetry.io/otel/trace"
)

const scope = "gitea.com/gitea/runner"

type taskRunKey struct{}

type taskRun struct {
	id, url string
}

func Setup(ctx context.Context, uuid, name string) (func(context.Context) error, error) {
	var shutdowns []func(context.Context) error
	shutdown := func(ctx context.Context) (err error) {
		for _, fn := range shutdowns {
			err = errors.Join(err, fn(ctx))
		}
		return err
	}
	traces, metrics := exporting("TRACES"), exporting("METRICS")
	if !traces && !metrics {
		return shutdown, nil
	}
	for _, key := range []string{"OTEL_EXPORTER_OTLP_PROTOCOL", "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL", "OTEL_EXPORTER_OTLP_METRICS_PROTOCOL"} {
		if protocol := os.Getenv(key); protocol != "" && !strings.EqualFold(protocol, "http/protobuf") {
			return shutdown, fmt.Errorf("%s=%s is not supported, only http/protobuf is", key, protocol)
		}
	}
	res, err := resource.New(ctx, resource.WithTelemetrySDK(), resource.WithAttributes(
		semconv.ServiceName("gitea-runner"),
		semconv.ServiceVersion(ver.Version()),
		semconv.ServiceInstanceID(uuid),
		semconv.CICDWorkerID(uuid),
		semconv.CICDWorkerName(name),
	), resource.WithFromEnv())
	if err != nil {
		return shutdown, err
	}
	if traces {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			return shutdown, err
		}
		provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
		otel.SetTracerProvider(provider)
		shutdowns = append(shutdowns, provider.Shutdown)
	}
	if metrics {
		exporter, err := otlpmetrichttp.New(ctx)
		if err != nil {
			return shutdown, err
		}
		provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)), sdkmetric.WithResource(res))
		otel.SetMeterProvider(provider)
		shutdowns = append(shutdowns, provider.Shutdown)
	}
	return shutdown, nil
}

func exporting(signal string) bool {
	return !strings.EqualFold(os.Getenv("OTEL_SDK_DISABLED"), "true") && !strings.EqualFold(os.Getenv("OTEL_"+signal+"_EXPORTER"), "none") &&
		cmp.Or(os.Getenv("OTEL_EXPORTER_OTLP_"+signal+"_ENDPOINT"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")) != ""
}

func ClientOptions() []connect.ClientOption {
	if _, ok := otel.GetMeterProvider().(*sdkmetric.MeterProvider); !ok {
		return nil
	}
	interceptor, err := otelconnect.NewInterceptor(otelconnect.WithoutTracing())
	if err != nil {
		otel.Handle(err)
		return nil
	}
	return []connect.ClientOption{connect.WithInterceptors(interceptor)}
}

func StartJob(ctx context.Context, task *runnerv1.Task) (context.Context, func(runnerv1.Result)) {
	_, tracing := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	_, metering := otel.GetMeterProvider().(*sdkmetric.MeterProvider)
	if !tracing && !metering {
		return ctx, func(runnerv1.Result) {}
	}
	field := func(key string) string { return task.GetContext().GetFields()[key].GetStringValue() }
	pipeline, owner, repo := field("workflow"), field("repository_owner"), field("repository")
	repoURL := strings.TrimSuffix(field("server_url"), "/") + "/" + repo
	runURL := repoURL + "/actions/runs/" + field("run_id")
	repository := []attribute.KeyValue{semconv.VCSOwnerName(owner), semconv.VCSRepositoryName(strings.TrimPrefix(repo, owner+"/"))}
	ctx, span := otel.Tracer(scope).Start(ctx, "RUN "+pipeline, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(repository...), trace.WithAttributes(
		semconv.CICDPipelineName(pipeline),
		semconv.CICDPipelineActionNameRun,
		semconv.CICDPipelineRunID(field("run_id")),
		semconv.CICDPipelineRunURLFull(runURL),
		semconv.VCSRepositoryURLFull(repoURL),
		semconv.VCSRefHeadName(field("ref_name")),
		semconv.VCSRefHeadRevision(field("sha")),
		attribute.String("gitea.job.id", field("job")),
	))
	meter := otel.Meter(scope)
	active, _ := cicdconv.NewPipelineRunActive(meter)
	duration, _ := cicdconv.NewPipelineRunDuration(meter)
	active.Add(ctx, 1, pipeline, cicdconv.PipelineRunStateExecuting, repository...)
	start := time.Now()
	return context.WithValue(ctx, taskRunKey{}, taskRun{strconv.FormatInt(task.Id, 10), runURL}), func(result runnerv1.Result) {
		attrs := end(span, semconv.CICDPipelineResultKey, jobResult(ctx, result))
		active.Add(ctx, -1, pipeline, cicdconv.PipelineRunStateExecuting, repository...)
		duration.Record(ctx, time.Since(start).Seconds(), pipeline, cicdconv.PipelineRunStateExecuting, append(attrs, repository...)...)
	}
}

func AddJobEnv(ctx context.Context, envs map[string]string) {
	if ctx.Value(taskRunKey{}) == nil {
		return
	}
	add := func(key, value string) {
		if _, ok := envs[key]; !ok {
			envs[key] = value
		}
	}
	for _, kv := range os.Environ() {
		if key, value, _ := strings.Cut(kv, "="); strings.HasPrefix(key, "OTEL_") && key != "OTEL_SERVICE_NAME" {
			add(key, value)
		}
	}
	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, carrier)
	for key, value := range carrier {
		add(strings.ToUpper(key), value)
	}
}

// StartStep names spans like Gitea's step list does, prefixing Pre and Post stages.
func StartStep(ctx context.Context, step *model.Step, stage string) (context.Context, func(error, *model.StepResult)) {
	parent, ok := ctx.Value(taskRunKey{}).(taskRun)
	if !ok || !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, func(error, *model.StepResult) {}
	}
	name, id := step.Name, parent.id+"."+strconv.Itoa(step.Number)
	if name == "" {
		run, _, _ := strings.Cut(strings.TrimSpace(step.Run), "\n")
		name = "Run " + cmp.Or(strings.TrimSpace(run), step.String())
	}
	if stage != "Main" {
		name, id = stage+" "+name, id+"."+strings.ToLower(stage)
	}
	ctx, span := otel.Tracer(scope).Start(ctx, name, trace.WithAttributes(
		semconv.CICDPipelineTaskName(name),
		semconv.CICDPipelineTaskRunID(id),
		semconv.CICDPipelineTaskRunURLFull(parent.url),
	))
	return context.WithValue(ctx, taskRunKey{}, taskRun{id, parent.url}), func(err error, result *model.StepResult) {
		end(span, semconv.CICDPipelineTaskRunResultKey, stepResult(ctx, err, result))
	}
}

func jobResult(ctx context.Context, result runnerv1.Result) string {
	switch {
	case result == runnerv1.Result_RESULT_SUCCESS:
		return "success"
	case result == runnerv1.Result_RESULT_SKIPPED:
		return "skip"
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		return "timeout"
	case result == runnerv1.Result_RESULT_CANCELLED:
		return "cancellation"
	case result == runnerv1.Result_RESULT_FAILURE:
		return "failure"
	}
	return "error"
}

func stepResult(ctx context.Context, err error, result *model.StepResult) string {
	switch {
	case errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded):
		return "timeout"
	case errors.Is(err, context.Canceled) || ctx.Err() != nil:
		return "cancellation"
	case err != nil || result.Outcome == model.StepStatusFailure:
		return "failure"
	}
	return "success"
}

func end(span trace.Span, key attribute.Key, result string) []attribute.KeyValue {
	attrs := []attribute.KeyValue{key.String(result)}
	if result == "failure" || result == "timeout" || result == "error" {
		attrs = append(attrs, semconv.ErrorTypeKey.String(result))
		span.SetStatus(codes.Error, "")
	}
	span.SetAttributes(attrs...)
	span.End()
	return attrs
}
