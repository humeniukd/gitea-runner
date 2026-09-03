# OpenTelemetry

The runner can send [OpenTelemetry](https://opentelemetry.io) traces and metrics
about the jobs it runs:

- **Traces** show each job as a tree of timed **spans** — the job and each step —
  so you can see where a slow or failed job spent its time without reading the
  logs.
- **Metrics** are counts and durations over time: job durations, active jobs,
  and errors.

Both are sent over OTLP to a **collector** you run, which forwards them to a
backend such as Jaeger, Tempo, or Grafana. Span and metric names follow the
OpenTelemetry [semantic conventions](https://opentelemetry.io/docs/specs/semconv/),
so backends understand them without extra mapping.

This is **opt-in**: nothing is exported, and no SDK starts, until you set an OTLP
endpoint. The Prometheus `/metrics` endpoint is separate and unchanged.

## Configuration

Configure it in the runner config file or with the standard `OTEL_*` environment
variables; the file wins per setting, env fills the rest. Traces and metrics are
configured separately, and the top-level `otel.*` keys are shared defaults each
section can override.

```yaml
otel:
  enabled: true                       # false turns everything off (also honors OTEL_SDK_DISABLED)
  endpoint: "http://collector:4318"   # the signal path (/v1/traces, /v1/metrics) is appended
  protocol: "http/protobuf"           # or "grpc"
  headers:                            # sent on every export
    authorization: "Bearer secret"
  resource_attributes:                # attached to every span and metric (see below)
    deployment.environment.name: "production"
  traces:
    enabled: true                     # defaults to otel.enabled
    endpoint: "http://collector:4318/v1/traces" # used as-is; overrides otel.endpoint
    protocol: "grpc"
    headers:
      x-scope-orgid: "ci"
    sampler_ratio: 0.25               # sample this fraction of traces
    rpc: false                        # also trace the runner's own calls to Gitea; off (noisy)
  metrics:
    enabled: true
    endpoint: "http://collector:4318/v1/metrics"
    protocol: "http/protobuf"
    headers:
      x-scope-orgid: "ci"
```

Setting just `otel.endpoint` turns on both traces and metrics. A shared endpoint
has `/v1/traces` or `/v1/metrics` appended; an endpoint set under `traces:` or
`metrics:` is used exactly as given (include the path).

`otel.resource_attributes` labels every span and metric, and is also passed to
job containers. Values resolve per key: a built-in default, then
`OTEL_RESOURCE_ATTRIBUTES`, then `otel.resource_attributes`.
`deployment.environment.name` defaults to `test` and stays unless you set it.

### Environment variables

Every setting also has a [standard `OTEL_*` variable](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/);
the config file wins over it per setting.

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

So tools inside a job can add to the same trace, each job container gets:

- every `OTEL_*` variable from the runner's environment;
- the `otel.*` config settings as their `OTEL_*` equivalents (config-file
  settings do not otherwise reach the container);
- `TRACEPARENT`/`TRACESTATE` pointing at the job's span, so spans created inside
  the job attach to it.

Disabling traces or metrics (`otel.traces.enabled: false` /
`otel.metrics.enabled: false`) drops that section's variables and sets
`OTEL_TRACES_EXPORTER=none` / `OTEL_METRICS_EXPORTER=none`, so a shared endpoint
cannot turn it back on.

> **Security.** Exporter credentials are not hidden from jobs: workflow code —
> including from untrusted pull requests — can read the collector headers and any
> TLS client key/cert. Do not put shared-backend credentials on the runner; point
> it at a local collector that accepts OTLP without auth and forwards on, keeping
> the credentials on the collector.

## Kubernetes: inject config with the OpenTelemetry Operator

With the [OpenTelemetry Operator](https://github.com/open-telemetry/opentelemetry-operator)
you can inject the `OTEL_*` configuration into the runner pod instead of setting
it by hand. The runner already has its own SDK, so use the `inject-sdk`
annotation, which only injects the config from an `Instrumentation` resource and
adds no language agent.

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

The operator adds the `OTEL_*` variables to the container; the runner reads them
and passes them to job containers as above.

## Spans

Each task produces a **job** span with one **step** span per step, using
[CI/CD](https://opentelemetry.io/docs/specs/semconv/cicd/cicd-spans/) attributes.

A task runs one job of a workflow, so today one span covers both; a separate
per-workflow span that ties a run's jobs together is left to Gitea.

### Job span

Kind `SERVER`, named `RUN <workflow> <job>` (the job name, else its id) once the
workflow is parsed. The reporter's calls to Gitea and the step spans are children
of it.

| Attribute | Value |
| --- | --- |
| `cicd.pipeline.action.name` | `RUN` |
| `cicd.pipeline.name` | workflow name |
| `cicd.pipeline.run.id` | Gitea run id |
| `cicd.pipeline.task.run.id` | Gitea task id |
| `cicd.pipeline.result` | `success`, `failure`, `cancellation`, `skip`, or `error` |
| `cicd.worker.id` | runner UUID |
| `cicd.worker.name` | runner name |
| `gitea.job.id` | job id (the `jobs.<id>` key) |
| `gitea.job.name` | job `name:` (may be empty) |
| `vcs.owner.name` | repository owner |
| `vcs.repository.name` | repository name |
| `vcs.repository.url.full` | repository URL |
| `vcs.ref.head.name` | branch or tag name |
| `vcs.ref.head.revision` | commit SHA |
| `vcs.ref.head.type` | `branch` or `tag` |

The status is set to error for any non-success result.

### Step span

Kind `INTERNAL`, named as the runner logs the step: its `name:`, else the
`uses:` action, else the `run:` script, and only `step <n>` (0-based) when it has
none of those.

| Attribute | Value |
| --- | --- |
| `cicd.pipeline.task.name` | step name (as above) |
| `cicd.pipeline.task.run.id` | `<task id>.<step index>` |
| `cicd.pipeline.task.type` | `build` |
| `cicd.pipeline.task.run.result` | `success`, `failure`, or `cancellation` |
| `gitea.step.uses` | action reference, for `uses:` steps |
| `gitea.step.run` | script, for `run:` steps |

On failure the status is set to error and its message is the step's last output
line (e.g. `npm ERR! missing script: build`), which is clearer than the bare exit
code; the error is also recorded as an event.

### Instrumented clients

The job cache server and job-summary uploads (HTTP) get `http.*` spans, and
spans from applications inside a job attach through the injected `TRACEPARENT`.

The runner's own calls to Gitea are **not** traced by default — the polling and
status-update calls are constant and rarely useful. Set `otel.traces.rpc: true`
to record them (and propagate trace context to Gitea) when debugging. RPC
metrics (`rpc.client.*`) are recorded either way.

## Logs

When tracing is on, the runner's log lines for a job carry `trace_id` and
`span_id`, so you can jump from a log line to the trace. Lines outside a job
(startup, polling) are unchanged.

## Metrics

When metrics are on, the runner reports its own instruments — named to the CI/CD
conventions where one fits, otherwise under `gitea.runner.*` — plus library and
Go runtime metrics. This is separate from the Prometheus `/metrics` endpoint.

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
start, and the repository keeps the label low-cardinality.

### Library and runtime metrics

- Connect RPC client metrics (`rpc.client.duration`, `rpc.client.request.size`,
  `rpc.client.response.size`), and HTTP client metrics (`http.client.*`) from the
  cache and job-summary clients.
- Go runtime metrics (`go.memory.*`, `go.gc.*`, `go.goroutine.count`,
  `go.schedule.duration`, …).
