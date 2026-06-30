# Gitea Actions Runner — Kubernetes-Native Pod Backend (Helm Chart)

This directory is a **Helm chart** that deploys the Gitea Actions runner in
**native pod mode**: each Actions job runs as one Kubernetes Pod (init container
for workspace setup, long-lived job container for step execution via
`remotecommand`, native sidecars for services). No Docker daemon is required in
the runner pod.

> Not to be confused with the DinD-backed manifests in
> `examples/kubernetes/dind-docker.yaml` and `rootless-docker.yaml`, which run
> Docker-in-Docker inside the runner pod.

---

## Install

```bash
# 1. Create the registration token Secret (or pass --set runner.registration.create=true
#    --set runner.registration.token=<token> to have the chart create it).
kubectl -n gitea-e2e create secret generic gitea-runner-token \
  --from-literal=token=<registration-token>

# 2. Install the chart.
helm install gitea-runner examples/kubernetes/native \
  --namespace gitea-e2e \
  --set runner.giteaInstanceUrl=https://gitea.example.com
```

The chart creates its own namespace by default; pass `--set namespace.create=false`
to render into an existing one.

### Verify the runner is online

```bash
kubectl -n gitea-e2e get pods -l app.kubernetes.io/name=gitea-runner
kubectl -n gitea-e2e logs -l app.kubernetes.io/name=gitea-runner --follow
```

The runner registers against `runner.giteaInstanceUrl` on first start (when
`/data/.runner` does not yet exist on the PVC), then enters daemon mode.

---

## How the chart maps to the runner

| Chart value | What it controls | Runner side |
|---|---|---|
| `runner.image` | Runner container image | — |
| `runner.giteaInstanceUrl` | Gitea server URL | `GITEA_INSTANCE_URL` env |
| `runner.labels` | Job-routing labels (`runs-on`) | `GITEA_RUNNER_LABELS` env (default `k8s-native`) |
| `runner.name` / `ephemeral` / `once` | Registration identity/mode | `GITEA_RUNNER_NAME` / `_EPHEMERAL` / `_ONE` |
| `runner.registration.secretName`/`secretKey` | Token Secret reference | `GITEA_RUNNER_REGISTRATION_TOKEN` env |
| `runner.registration.create`/`token` | Optional Secret creation | creates the Secret above |
| `runner.registration.maxAttempts` | Reg. retry count | `GITEA_MAX_REG_ATTEMPTS` env |
| `runner.persistence.*` | `/data/.runner` PVC | persists registration across restarts |
| `serviceAccount.name` | Runner SA | config `kubernetes.service_account` |
| `rbac.forwardLooking` | Strict vs documented RBAC | see [RBAC](#rbac) |
| `cache.enabled`/`port`/`dir` | Embedded artifact cache | config `cache.*` |
| `cache.externalServer`/`externalSecret` | External cache-server mode | config `cache.external_*` |
| `cache.service.enabled` | ClusterIP Service for cache | the `gitea-runner-cache` Service |
| `metrics.service.enabled` | Prometheus scrape Service | config `metrics.addr` → `0.0.0.0:9101` |
| `config.kubernetes.*` | Backend knobs | config `kubernetes.*` |
| `config.kubernetes.runtimeClassName` | Per-runner sandbox policy | config `kubernetes.runtime_class_name` |
| `runtimeClasses.create` | Ship gVisor/Kata objects | the RuntimeClass YAML |

`POD_NAME`/`POD_UID` are injected via the downward API so the runner sets an
own-Pod `ownerReference` on every job pod (cluster GC + the startup
dangling-owner sweep rely on this).

---

## RBAC

The Role (`rbac.forwardLooking: true`, default) grants the documented surface:

| Resource | Verbs | Purpose |
|---|---|---|
| `pods` | create, get, list, watch, delete, **patch** | Job-pod lifecycle + retain-marker merge-patch |
| `pods/exec` | create | Step execution via `remotecommand` |
| `pods/log` | get | _(documented; currently unused — logs stream over exec)_ |
| `secrets`, `configmaps` | create, get, delete | Per-job cleanup; create/get forward-looking |

> **`patch` on `pods` is required** even though the previous static manifest
> omitted it: the runner merge-patches `actions.gitea.com/retain=true` onto a
> job pod when `pod_cleanup` keeps it. Without `patch`, retained-pod marking
> fails.

Set `rbac.forwardLooking: false` for strictest least-privilege — the Role then
grants only what the runtime exercises today (`pods/exec create`, `secrets`/
`configmaps delete`, and the `pods` set above minus the redundant `pods/log`).

`pods/attach` is never granted — the backend uses `pods/exec` exclusively.

---

## Cache reachability

In Kubernetes mode job pods live in a different network namespace than the
runner, so the embedded artifact-cache must be reachable at a cluster-routable
address. The chart wires this three ways:

1. **Embedded cache (default).** `cache.enabled: true` with no
   `externalServer`. The chart renders the `gitea-runner-cache` ClusterIP
   Service, exposes port `cache.port` on the Deployment, and sets both
   `cache.host` and `kubernetes.connect_back_address` to that Service's DNS.
   The runner's Phase 6 startup guard rejects loopback values here, so a
   misconfiguration fails fast.

2. **External cache server.** Set `cache.externalServer` (and a matching
   `cache.externalSecret`). The embedded cache is not started, no cache
   Service/port is rendered, and the Deployment drops its cache probes. Point
   this at a standalone `gitea-runner cache-server` shared across runners.

3. **Disabled.** `cache.enabled: false`. No cache of any kind. `actions/cache`,
   `upload-artifact`, and `download-artifact` will not work.

### NetworkPolicy

If the namespace enforces a **default-deny ingress** policy, job pods cannot
reach the cache Service. Enable the allow policy:

```yaml
networkPolicy:
  create: true
```

It selects the runner pod for ingress on `cache.port` from same-namespace pods.

---

## DinD (optional, for `uses: docker://` steps)

```yaml
config:
  kubernetes:
    dind:
      enabled: true      # explicit opt-in; no auto-detect
      rootless: true     # default: lower privilege; false = privileged sidecar
      image: docker:28-dind-rootless
```

When enabled, the backend injects an on-demand DinD sidecar into each job pod
that actually has a `uses: docker://` step. `using: docker` actions
(`action.yml` `runs.using: docker`) are **not supported in v1** and produce an
explicit error — they are not silently ignored.

---

## Customisation reference

| Field | Default | Notes |
|---|---|---|
| `config.kubernetes.defaultJobImage` | `node:20-bookworm` | Image for bare `run:` jobs |
| `config.kubernetes.podReadyTimeout` | `10m` | Code default is `5m` |
| `config.kubernetes.podCleanup` | `always` | `always`/`on-failure`/`never` |
| `config.kubernetes.maxConcurrentJobs` | `0` (unlimited) | Bypasses `runner.capacity` in k8s mode |
| `config.kubernetes.runtimeClassName` | _(unset)_ | Enforced on every job pod |
| `config.kubernetes.workspaceVolume.type` | `emptyDir` | `emptyDir`/`ephemeral`/`pvc` |
| `config.kubernetes.resources.{job,service,dind}` | M3 defaults | dind needs ≥1Gi |
| `runner.capacity` | `4` | Cosmetic in k8s mode (bypassed) |

---

## Job concurrency (kube-scheduler handles the backlog)

In Kubernetes mode the runner only **creates** job pods and execs into them —
the **kube-scheduler** queues and places the real workloads. So `runner.capacity`
(which bounds the docker/host backends) is bypassed; use
`kubernetes.max_concurrent_jobs` instead. `0` (default) means unlimited: the
runner accepts every assigned job and lets the scheduler pend the backlog.

The correct way to bound a build namespace is a **ResourceQuota** (the scheduler
then rejects/pend excess pods) — far better than throttling the runner:

```yaml
resourceQuota:
  create: true
  hard:
    pods: "50"
    requests.cpu: "20"
    requests.memory: 40Gi
    limits.cpu: "80"
    limits.memory: 160Gi
```

---

## Job pod retention (auditing)

`config.kubernetes.podCleanup` controls when a completed job pod (and its
per-job Secret/ConfigMap) is deleted:

| Value | Behaviour |
|---|---|
| `always` (default) | Delete after every job. |
| `on-failure` | Keep pods for **failed** jobs; delete on success. |
| `never` | Keep every pod. |

Retained pods are merge-patched with `actions.gitea.com/retain=true` (this is
why the Role needs `patch` on `pods`) and survive runner-restart sweeps:

```bash
kubectl -n gitea-e2e get pods -l actions.gitea.com/retain=true
kubectl -n gitea-e2e exec -it <pod-name> -c job -- sh
# prune:
kubectl -n gitea-e2e delete pods -l actions.gitea.com/retain=true
```

> Retained pods keep consuming their requested CPU/memory **and** count against
> the namespace ResourceQuota until deleted.

---

## Sandboxed runtimes (gVisor / Kata)

A Kubernetes `RuntimeClass` is all that's needed to run untrusted job code under
gVisor (`runsc`) or Kata (`kata`) instead of the node default (`runc`). The
runner sets `pod.spec.runtimeClassName` from config on every job pod.

**Node prerequisite:** the runtime handler binary + a matching containerd/CRI-O
runtime must be installed on every scheduling node. Then:

```yaml
runtimeClasses:
  create: true
  gvisor:
    handler: runsc
    nodeSelector: { sandbox.gitea.com/runtime: gvisor }

config:
  kubernetes:
    runtimeClassName: gvisor
```

The runtime is a **per-runner operator policy**, not a per-job label: a workflow
cannot opt out of the sandbox. To offer different isolation levels, run **one
runner per trust tier** (each with its own `runtimeClassName`, namespace, and
node pool) and route jobs with `runs-on` — the label selects the runner, and
each runner enforces its own runtime.

Notes:

- A misnamed/absent `RuntimeClass` surfaces at pod admission (creation error or
  a `Pending` pod with a clear event).
- A **privileged** DinD sidecar (`dind.rootless: false`) is generally
  incompatible with gVisor; prefer rootless DinD or a `runc`-backed runner.

---

## Security considerations

- **DinD daemon serves plaintext TCP on pod loopback.** When
  `dind.enabled: true` the sidecar listens on `tcp://localhost:2375` with TLS
  disabled, reachable only over the job Pod's shared network namespace. Safe
  **provided no untrusted container is co-scheduled into the job Pod**.
- **Translated `uses: docker://` steps run with `--network host`** inside DinD,
  sharing the DinD network namespace. Treat container-step images with the same
  trust level as the job.
- **`dind.rootless: false` makes every job's DinD privileged.** Only disable
  rootless when a specific image genuinely cannot build under rootless DinD.
- **`replicas > 1` is unsupported today:** the registration-state PVC is
  `ReadWriteOnce` and the cache Service reaches only one pod's embedded cache.
  Run `replicas: 1` unless you understand the caveats.
