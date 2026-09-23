# OpenTelemetry

The runner can send traces of the jobs it runs to an OpenTelemetry collector over OTLP/HTTP. It stays off until you set an endpoint in the runner's environment:

```sh
OTEL_EXPORTER_OTLP_ENDPOINT=http://collector:4318
```

The exporter also reads `OTEL_EXPORTER_OTLP_HEADERS` and `OTEL_EXPORTER_OTLP_TRACES_HEADERS`. The standard [SDK variables](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/) for resource attributes, sampling and batching apply too, and `OTEL_TRACES_EXPORTER=none` or `OTEL_SDK_DISABLED=true` turns tracing off. Only the `http/protobuf` protocol is supported.

## What is sent

A trace per job, named `RUN {workflow} {job}`, with a span for every step named as in Gitea, including `Pre` and `Post` stages and the steps of composite actions. Spans carry the repository, ref, commit and result.

## Tracing inside jobs

Jobs receive `TRACEPARENT`, so tools in a job that use OpenTelemetry add their spans to the job's trace. Set the `OTEL_*` variables they export with in `runner.envs`.

## Metrics

Metrics stay on the Prometheus endpoint enabled by `metrics.enabled`. A collector picks them up by scraping `/metrics` with its `prometheus` receiver.
