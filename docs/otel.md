# OpenTelemetry

The runner exports OpenTelemetry traces and metrics for the jobs it executes:
one span per job (`RUN {workflow name}`) and one per workflow step, connected
with [W3C Trace Context](https://www.w3.org/TR/trace-context/), plus its own
runner metrics over OTLP. RPCs to the Gitea server and the job cache server are
instrumented as well.

## Configuration

A signal is off unless its OTLP endpoint is configured. Configuration comes
from two sources, merged per setting: the runner's own config file wins, and
the standard `OTEL_*` environment variables fill in everything it leaves empty.
The top-level `otel.{enabled,endpoint,protocol,headers}` are shared defaults for
both signals — mirroring the generic `OTEL_EXPORTER_OTLP_*` env vars — and each
per-signal section overrides them.

```yaml
otel:
  enabled: true                       # shared default; false turns both signals off (also honors OTEL_SDK_DISABLED)
  endpoint: "http://collector:4318"   # shared base; the signal path is appended (mirrors OTEL_EXPORTER_OTLP_ENDPOINT)
  protocol: "http/protobuf"           # shared; "http/protobuf" or "grpc" (mirrors OTEL_EXPORTER_OTLP_PROTOCOL)
  headers:                            # shared; merged over OTEL_EXPORTER_OTLP_HEADERS, per key
    authorization: "Bearer secret"
  traces:
    enabled: true                     # defaults to otel.enabled
    endpoint: "http://collector:4318/v1/traces" # used as-is; wins over otel.endpoint and the env vars
    protocol: "grpc"                  # wins over otel.protocol
    headers:                          # merged over otel.headers, per key
      x-scope-orgid: "ci"
    sampler_ratio: 0.25               # parent-based traceidratio; wins over OTEL_TRACES_SAMPLER
  metrics:
    enabled: true                     # defaults to otel.enabled
    endpoint: "http://collector:4318/v1/metrics" # used as-is
    protocol: "http/protobuf"
    headers:                          # merged over otel.headers, per key
      x-scope-orgid: "ci"
```

Setting only `otel.endpoint` (or the standard `OTEL_EXPORTER_OTLP_ENDPOINT`)
enables and points both signals. Endpoint path handling matches the OTLP spec:
a shared/generic endpoint has the signal path (`/v1/traces`, `/v1/metrics`)
appended, while a per-signal endpoint is used as-is (include the path). The
Prometheus metrics endpoint (`metrics.*`) is independent of `otel.metrics`.
When OTel metrics export is on, the runner emits the CI/CD semconv metrics
(`cicd.pipeline.run.duration`, `cicd.pipeline.run.active`, `cicd.worker.count`,
`cicd.system.errors`), poll instruments under `gitea.runner.poll.*`, the Go
runtime metrics (`go.*`) from the contrib runtime instrumentation, and the
connect RPC client metrics.

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

## Metrics

When OTel metrics export is on, the runner emits its own instruments — named to
the CI/CD semantic conventions where one applies, otherwise under the
`gitea.runner.*` namespace — plus instrumentation-library and Go runtime
metrics. This is independent of the Prometheus `/metrics` endpoint, which is
unchanged.

### Runner instruments

| Metric | Instrument | Unit | Attributes |
| --- | --- | --- | --- |
| `cicd.pipeline.run.duration` | histogram | `s` | `cicd.pipeline.name`, `cicd.pipeline.run.state=executing`, `cicd.pipeline.result` |
| `cicd.pipeline.run.active` | up/down counter | `{run}` | `cicd.pipeline.name`, `cicd.pipeline.run.state=executing` |
| `cicd.worker.count` | up/down counter | `{worker}` | `cicd.worker.state` (`busy`, `available`) |
| `cicd.system.errors` | counter | `{error}` | `cicd.system.component=runner`, `error.type` |
| `gitea.runner.poll.fetch` | counter | `{fetch}` | `gitea.poll.result` (`task`, `empty`, `error`) |
| `gitea.runner.poll.fetch.duration` | histogram | `s` | — |
| `gitea.runner.poll.backoff` | gauge | `s` | — |

`cicd.pipeline.name` is the repository: the workflow name is not known at job
start, and the repository is the low-cardinality grouping.

### Library and runtime metrics

- Connect RPC client metrics (`rpc.client.duration`, `rpc.client.request.size`,
  `rpc.client.response.size`) from the connect instrumentation, and HTTP client
  metrics (`http.client.*`) from the cache and job-summary clients.
- Go runtime metrics (`go.memory.*`, `go.gc.*`, `go.goroutine.count`,
  `go.schedule.duration`, …) from the contrib runtime instrumentation.
