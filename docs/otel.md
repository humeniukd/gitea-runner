# OpenTelemetry

The runner can export telemetry about the jobs it executes through
[OpenTelemetry](https://opentelemetry.io) (OTel), a vendor-neutral standard
with two parts:

- **Traces** follow one job from start to finish as a tree of **spans** — each
  span is one unit of work with a name, timing, and success/failure status.
  One span covers the whole job, one each workflow step, and one each RPC or
  HTTP call the runner makes, so a slow or failed job can be decomposed into
  where the time went without reading logs.
- **Metrics** are aggregate numbers over time — job durations, active jobs,
  error counts — sent to the same destination.

Both are delivered with OTLP (OpenTelemetry's wire protocol) to a **collector**,
a small server you run that forwards to any backend (Jaeger, Tempo, Datadog,
Grafana Cloud, ...). Attribute names follow the OTel **semantic conventions**
("semconv") — standardized names such as `cicd.pipeline.run.id`, so backends
and dashboards understand them without runner-specific mapping.

**This feature is opt-in.** The runner exports nothing, and starts no OTel
SDK at all, until you configure an OTLP endpoint — in the config file or
through the standard `OTEL_*` environment variables. The Prometheus
`/metrics` endpoint is independent and unaffected.

## Configuration

A "signal" is one of the two exports — traces or metrics; each is off until its
OTLP endpoint is configured. Configuration comes from two sources, merged per
setting: the runner's own config file wins, and the standard `OTEL_*`
environment variables fill in everything it leaves empty. The top-level
`otel.{enabled,endpoint,protocol,headers}` are shared defaults for both
signals — mirroring the generic `OTEL_EXPORTER_OTLP_*` env vars — and each
per-signal section overrides them.

```yaml
otel:
  enabled: true                       # shared default; false turns both signals off (also honors OTEL_SDK_DISABLED)
  resource_attributes:                # extra resource attributes, merged over OTEL_RESOURCE_ATTRIBUTES (see below)
    deployment.environment.name: "production"
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

`otel.resource_attributes` sets [resource attributes](https://opentelemetry.io/docs/specs/semconv/resource/) that identify the runner on every span and metric, and are also passed to job containers (via `OTEL_RESOURCE_ATTRIBUTES`) so in-job telemetry carries them too. They resolve per key: the built-in default, then `OTEL_RESOURCE_ATTRIBUTES`, then `otel.resource_attributes` (most specific wins). `deployment.environment.name` defaults to `test`, and that default persists unless it is overridden explicitly — setting other attributes does not remove it.

### Environment variables

Every setting can also be given as a [standard `OTEL_*` variable](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/)
instead of the config file. The runner reads them directly, and
config-file value (when set) wins per setting.

| Config key | Environment variable |
| --- | --- |
| `otel.endpoint` | `OTEL_EXPORTER_OTLP_ENDPOINT` |
| `otel.protocol` | `OTEL_EXPORTER_OTLP_PROTOCOL` |
| `otel.headers` | `OTEL_EXPORTER_OTLP_HEADERS` |
| `otel.resource_attributes` | `OTEL_RESOURCE_ATTRIBUTES` |
| `otel.traces.endpoint` | `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` |
| `otel.traces.protocol` | `OTEL_EXPORTER_OTLP_TRACES_PROTOCOL` |
| `otel.traces.headers` | `OTEL_EXPORTER_OTLP_TRACES_HEADERS` |
| `otel.traces.sampler_ratio` | `OTEL_TRACES_SAMPLER=parentbased_traceidratio` + `OTEL_TRACES_SAMPLER_ARG` |
| `otel.metrics.endpoint` | `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT` |
| `otel.metrics.protocol` | `OTEL_EXPORTER_OTLP_METRICS_PROTOCOL` |
| `otel.metrics.headers` | `OTEL_EXPORTER_OTLP_METRICS_HEADERS` |

```sh
export OTEL_EXPORTER_OTLP_ENDPOINT="http://collector:4318"
export OTEL_EXPORTER_OTLP_PROTOCOL="http/protobuf"
export OTEL_EXPORTER_OTLP_HEADERS="authorization=Bearer secret"
export OTEL_RESOURCE_ATTRIBUTES="deployment.environment.name=production"
export OTEL_TRACES_SAMPLER="parentbased_traceidratio"
export OTEL_TRACES_SAMPLER_ARG="0.25"
gitea-runner daemon
```

## Job containers

The runner gives job containers the full OTLP configuration so applications
running inside a job continue the CI trace and export to the same collector
instead of starting disjoint traces. Each container receives:

- every `OTEL_*` variable from the runner's own environment;
- the `otel.*` config settings translated to their `OTEL_*` env vars (endpoint,
  protocol, headers, `resource_attributes`, and per-signal equivalents), so file
  config reaches jobs too — the yml config is otherwise never in the process
  environment. The config value wins over the runner env, per setting;
- `TRACEPARENT` and `TRACESTATE` (W3C trace context in environment variables)
  pointing at the job's span, enabling spans inside the job to be tied to the
  job itself.

A signal disabled with `otel.traces.enabled: false` or `otel.metrics.enabled:
false` contributes no OTLP env vars to the container and is turned off in-job
with `OTEL_TRACES_EXPORTER=none` / `OTEL_METRICS_EXPORTER=none`, so a shared
endpoint does not re-enable it.

> **Security limitation.** Exporter credentials are **not** withheld: job code —
> including jobs from untrusted pull requests — can read the collector headers
> (`OTEL_EXPORTER_OTLP*_HEADERS`) and any TLS client key/certificate passed to
> the SDK. Do not put credentials for a shared backend on the runner. The
> recommended setup is to point the runner at a local collector that accepts
> OTLP with no authentication and forwards to the real destination, keeping the
> backend credentials on the collector.

## Kubernetes: inject config with the OpenTelemetry Operator

If you run the [OpenTelemetry Operator](https://github.com/open-telemetry/opentelemetry-operator),
you can let it inject the OTLP configuration into the runner pod instead of
setting the variables by hand. The runner already ships its own SDK, so use the
`inject-sdk` annotation: it injects only the `OTEL_*` environment (endpoint,
resource attributes, propagators, sampler) from an `Instrumentation` resource,
without adding a language auto-instrumentation agent.

Define an `Instrumentation` resource:

```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: gitea-runner
spec:
  exporter:
    endpoint: http://collector.observability:4318
  propagators:
    - tracecontext
    - baggage
  resource:
    resourceAttributes:
      deployment.environment.name: production
  sampler:
    type: parentbased_traceidratio
    argument: "0.25"
```

Annotate the runner's pod template so the operator mutates the pod on creation:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: runner
spec:
  template:
    metadata:
      annotations:
        # "<namespace>/<name>" of the Instrumentation, or "true" for one named
        # "default" in the pod's namespace.
        instrumentation.opentelemetry.io/inject-sdk: "gitea-runner"
    spec:
      containers:
        - name: runner
          image: gitea/runner:nightly
```

The operator adds `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_RESOURCE_ATTRIBUTES`,
`OTEL_PROPAGATORS` and the sampler variables to the container. The runner reads
them like any other `OTEL_*` env and passes them on to job containers as
described above.

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
