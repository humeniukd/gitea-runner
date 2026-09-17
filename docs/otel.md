# OpenTelemetry

The runner exports OpenTelemetry traces and metrics over OTLP/HTTP, configured with the standard [`OTEL_*` environment variables](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/). A signal is only exported when an endpoint is set for it. Without one, OpenTelemetry stays off.

| Variable | Effect |
| --- | --- |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Exports traces to `<endpoint>/v1/traces` and metrics to `<endpoint>/v1/metrics` |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`, `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` | Exports one signal to the URL as given |
| `OTEL_TRACES_EXPORTER=none`, `OTEL_METRICS_EXPORTER=none` | Turns one signal off |
| `OTEL_SDK_DISABLED=true` | Turns both off |

The other exporter, resource, sampler and batching variables apply as specified. Only the `http/protobuf` protocol is supported. The runner reports itself as `service.name` `gitea-runner`, with its UUID as `service.instance.id` and `cicd.worker.id`.

## Traces

Each job is a `RUN <workflow>` span of kind `SERVER`, named after its workflow file, and each step is a child span named after the step. They follow the [CI/CD semantic conventions](https://opentelemetry.io/docs/specs/semconv/cicd/cicd-spans/), with the repository, ref and commit as `vcs.*` attributes and the outcome as `cicd.pipeline.result` or `cicd.pipeline.task.run.result`. A `failure`, `timeout` or `error` outcome sets the span status to error and `error.type` to the outcome. The runner sees one job of a workflow run at a time, so every job is a trace of its own. Calls to the Docker API appear as HTTP client spans.

## Metrics

- `cicd.pipeline.run.duration` and `cicd.pipeline.run.active`, with the workflow file as `cicd.pipeline.name`
- `rpc.client.*` for the calls to Gitea

The Prometheus endpoint is unaffected.

## Jobs

Jobs receive the runner's `OTEL_*` variables except `OTEL_SERVICE_NAME`, and the job span as `TRACEPARENT` and `TRACESTATE`, so instrumented tools in a job export into its trace. Values set in `runner.envs` take precedence. Workflows can read these variables, so point the runner at a collector that needs no credentials instead of passing them in `OTEL_EXPORTER_OTLP_HEADERS`.
