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
variables; the file wins per setting, env fills the rest. See the `otel:` section
of [config.example.yaml](../internal/pkg/config/config.example.yaml) for every
key. Each config key also has a
[standard `OTEL_*` variable](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/).

Traces and metrics are configured separately, and the top-level `otel.*` keys are
shared defaults each section can override. Setting just `otel.endpoint` turns on
both. A shared endpoint has `/v1/traces` or `/v1/metrics` appended; an endpoint
set under `traces:` or `metrics:` is used exactly as given (include the path).

`otel.resource_attributes` labels every span and metric and is also passed to job
containers. Values resolve per key: a built-in default, then
`OTEL_RESOURCE_ATTRIBUTES`, then `otel.resource_attributes`.
`deployment.environment.name` defaults to `test` and stays unless you set it.

## Job containers

So tools inside a job can add to the same trace, each job container gets:

- every `OTEL_*` variable from the runner's environment;
- the `otel.*` config settings as their `OTEL_*` equivalents (config-file
  settings do not otherwise reach the container);
- `TRACEPARENT`/`TRACESTATE` pointing at the job's span, so spans created inside
  the job attach to it.

`service.*` resource attributes are not forwarded, so each job is its own
service rather than inheriting the runner's identity.

Disabling traces or metrics (`otel.traces.enabled: false` /
`otel.metrics.enabled: false`) drops that section's variables and sets
`OTEL_TRACES_EXPORTER=none` / `OTEL_METRICS_EXPORTER=none`, so a shared endpoint
cannot turn it back on.

> **Security.** Exporter credentials are not hidden from jobs: workflow code —
> including from untrusted pull requests — can read the collector headers and any
> TLS client key/cert. Do not put shared-backend credentials on the runner; point
> it at a local collector that accepts OTLP without auth and forwards on, keeping
> the credentials on the collector.

On Kubernetes, the [OpenTelemetry Operator](https://github.com/open-telemetry/opentelemetry-operator)
can inject this `OTEL_*` configuration into the runner pod via its `inject-sdk`
annotation instead of you setting it by hand.

## Spans

Each task records:

- a **job span** (`SERVER`, named `RUN <workflow> <job>`) covering the whole task,
  carrying the pipeline, worker, and job identity, the repository and ref it runs
  against (VCS attributes), and the run result;
- a **step span** (`INTERNAL`) per workflow step, named as the runner logs the
  step (its `name:`, else the `uses:` action, else the `run:` script). On failure
  its status message is the step's last output line, which is clearer than the
  bare exit code.

Reporter calls to Gitea and the step spans are children of the job span. A task is
one job of a workflow; a span tying a whole run's jobs together is left to Gitea.

The job cache server and job-summary uploads get `http.*` client spans.
Applications inside a job attach to the trace through the injected `TRACEPARENT`.
The runner's own RPCs to Gitea are not traced by default — they are constant and
rarely useful; set `otel.traces.rpc: true` to trace them when debugging.

## Logs

When tracing is on, the runner's log lines for a job carry `trace_id` and
`span_id`, so you can jump from a log line to the trace. Lines outside a job
(startup, polling) are unchanged.

## Metrics

When metrics are on, the runner reports, separately from the Prometheus
`/metrics` endpoint:

- **pipeline metrics** — `cicd.pipeline.run.duration` and `cicd.pipeline.run.active`
  (job run timing and the number in progress);
- **worker and error metrics** — `cicd.worker.count` (busy vs available runners)
  and `cicd.system.errors`;
- **poll metrics** under `gitea.runner.poll.*` — FetchTask count, latency, and
  backoff;
- **library and Go runtime metrics** — connect `rpc.client.*`, `http.client.*`,
  and `go.*`.

`cicd.pipeline.name` is the repository, not the workflow name (which is not known
at job start), which keeps the label low-cardinality.
