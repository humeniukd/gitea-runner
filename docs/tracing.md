# Tracing

The runner can export OpenTelemetry traces for the jobs it executes: one span
per job (`RUN {workflow name}`) and one span per workflow step, connected with
[W3C Trace Context](https://www.w3.org/TR/trace-context/). RPCs to the Gitea
server and the job cache server are instrumented as well.

## Configuration

Tracing is off unless an OTLP endpoint is configured. Configuration comes from
two sources, merged per setting: the runner's own config file wins, and the
standard `OTEL_*` environment variables fill in everything it leaves empty.

```yaml
otel:
  endpoint: "http://collector:4318" # wins over OTEL_EXPORTER_OTLP(_TRACES)_ENDPOINT
  protocol: "http/protobuf"         # or "grpc"; wins over OTEL_EXPORTER_OTLP_TRACES_PROTOCOL
  headers:                          # merged over OTEL_EXPORTER_OTLP_HEADERS, per key
    authorization: "Bearer secret"
  sampler_ratio: 0.25               # parent-based traceidratio; wins over OTEL_TRACES_SAMPLER
```

Every other standard variable is honored, including
`OTEL_EXPORTER_OTLP_COMPRESSION`, `OTEL_EXPORTER_OTLP_TIMEOUT`,
`OTEL_SERVICE_NAME`, `OTEL_RESOURCE_ATTRIBUTES`, and (as a spec escape hatch)
`OTEL_CONFIG_FILE`, which hands the SDK a full OpenTelemetry configuration file
and supersedes the settings above. Invalid values are ignored with a warning;
the runner keeps working untraced if the SDK cannot start.

## Job containers

The runner passes its own `OTEL_*` environment variables into the containers
where jobs run, plus `TRACEPARENT` and `TRACESTATE` (W3C trace context in
environment variables) pointing at the job's span. Applications running inside
a job that use OpenTelemetry therefore continue the CI trace and export to the
same collector, instead of starting disjoint traces.

The runner's own collector credentials are withheld: variables carrying
exporter headers (`*HEADERS`) or a TLS client key/certificate are not handed to
job containers, because they are readable by workflow code, including jobs from
untrusted pull requests. Configure those in the runner's environment only; a
job application that needs its own collector credentials should receive them
through a workflow secret.

## Spans

Every task produces one job span with one child span per workflow step, using
[CI/CD semantic-convention](https://opentelemetry.io/docs/specs/semconv/cicd/cicd-spans/)
attributes.

### Job span

Kind `SERVER`, named `RUN` and refined to `RUN <workflow>` once the workflow is
parsed. Reporter RPCs and step spans are children of it.

| Attribute | Value |
| --- | --- |
| `cicd.pipeline.action.name` | `RUN` |
| `cicd.pipeline.name` | workflow name |
| `cicd.pipeline.run.id` | Gitea run id |
| `cicd.pipeline.task.run.id` | Gitea task id |
| `cicd.pipeline.result` | `success`, `failure`, `cancellation`, `skip`, or `error` |
| `cicd.worker.id` | runner UUID |
| `cicd.worker.name` | runner name |

The span status is set to error for any non-success result, and a run error is
recorded on the span.

### Step span

One per workflow step (main stage), kind `INTERNAL`, named after the step
(`<step name>`, or `step <n>` when unnamed).

| Attribute | Value |
| --- | --- |
| `cicd.pipeline.task.name` | step name |
| `cicd.pipeline.task.run.id` | `<task id>.<step index>` |
| `cicd.pipeline.task.type` | `build` |
| `cicd.pipeline.task.run.result` | `success`, `failure`, or `cancellation` |
| `gitea.step.uses` | action reference, for `uses:` steps |
| `gitea.step.run` | script, for `run:` steps |

### Instrumented clients

RPCs to Gitea (connect) and the job cache server plus job-summary uploads (HTTP)
carry the standard `rpc.*` / `http.*` span attributes from the OpenTelemetry
contrib instrumentation, and propagate W3C context to the peer. Spans emitted by
applications inside a job attach to the job span through the injected
`TRACEPARENT`.
