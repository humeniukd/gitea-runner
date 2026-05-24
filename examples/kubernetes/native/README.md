# Kubernetes-Native Pod Backend — Deployment Manifests

This directory contains the Kubernetes manifests for the Gitea Actions runner
operating in **native pod mode**: each Actions job runs as one Kubernetes Pod
(init container for workspace setup, long-lived job container for `act` step
execution via `remotecommand`, native sidecars for services).

> **Do not confuse with the adjacent DinD-backed manifests** in
> `examples/kubernetes/dind-docker.yaml` and `rootless-docker.yaml`, which run
> Docker-in-Docker inside the runner pod. This directory is for the
> Kubernetes-native backend that needs no Docker daemon in the runner itself.

---

## Files

| File | Purpose |
|---|---|
| `namespace.yaml` | `gitea-e2e` namespace |
| `serviceaccount.yaml` | `gitea-runner` ServiceAccount |
| `role.yaml` | Least-privilege Role (see RBAC section) |
| `rolebinding.yaml` | Binds Role to ServiceAccount |
| `configmap.yaml` | Runner config (includes `kubernetes:` section) |
| `service.yaml` | `gitea-runner-cache` ClusterIP Service for the embedded cache |
| `pvc.yaml` | PVC for `.runner` registration state |
| `runtimeclass.yaml` | Optional gVisor/Kata `RuntimeClass` objects for sandboxed jobs (see [Sandboxed runtimes](#sandboxed-runtimes-gvisor--kata)) |
| `deployment.yaml` | Runner Deployment |

---

## Quick Start

### 1. Create the registration token Secret

Obtain a runner registration token from Gitea (Site Admin → Runners → Add Runner),
or let the E2E Deno harness mint one (`tests/e2e/k8s/run.ts`). Then:

```bash
kubectl -n gitea-e2e create secret generic gitea-runner-token \
  --from-literal=token=<registration-token>
```

### 2. Apply manifests in order

```bash
kubectl apply -f examples/kubernetes/native/namespace.yaml
kubectl apply -f examples/kubernetes/native/
```

Applying the namespace first avoids ordering failures when all files are
applied together.

### 3. Verify the runner is online

```bash
kubectl -n gitea-e2e get pods -l app.kubernetes.io/name=gitea-runner
kubectl -n gitea-e2e logs -l app.kubernetes.io/name=gitea-runner --follow
```

The runner registers against `GITEA_INSTANCE_URL` on first start (when
`/data/.runner` does not yet exist on the PVC), then enters daemon mode.

---

## RBAC

The `role.yaml` grants the minimum permissions required by the native backend:

| Resource | Verbs | Purpose |
|---|---|---|
| `pods` | create, get, list, watch, delete | Job pod lifecycle |
| `pods/exec` | create | Step execution via `remotecommand` |
| `pods/log` | get | Log streaming |
| `secrets` | create, get, delete | Per-job credential Secrets |
| `configmaps` | create, get, delete | Per-job ConfigMaps |

`pods/attach` is **not granted** — there is no consumer in the native backend.
Step execution uses `pods/exec` (`remotecommand`) exclusively.

---

## Cache Service reachability

The runner's embedded artifact-cache server must be reachable from inside job
pods. The `gitea-runner-cache` Service exposes port 3000 and points at the
runner pod. The ConfigMap sets:

```yaml
cache:
  host: gitea-runner-cache.gitea-e2e.svc.cluster.local
  port: 3000
kubernetes:
  connect_back_address: "http://gitea-runner-cache.gitea-e2e.svc.cluster.local:3000"
```

The runner's Phase 6 startup guard rejects loopback (`127.0.0.1`/`localhost`)
values here when `kubernetes.enabled: true`, so misconfiguration fails fast.

---

## NetworkPolicy check

If your cluster enforces a **default-deny ingress** NetworkPolicy in
`gitea-e2e`, job pods will be unable to reach the cache Service and
`actions/cache`, `upload-artifact`, and `download-artifact` will fail.

**Always verify before running jobs:**

```bash
kubectl get networkpolicy -n gitea-e2e
```

If the output shows a default-deny policy (e.g. `default-deny-ingress` or
a policy with an empty `spec.ingress: []`), apply the allow policy below.

### Allow-ingress NetworkPolicy for the cache Service

Apply this manifest if a default-deny ingress policy exists in `gitea-e2e`:

```bash
kubectl apply -f - <<'EOF'
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-cache-ingress
  namespace: gitea-e2e
spec:
  # Select the runner pod(s) hosting the embedded cache server.
  podSelector:
    matchLabels:
      app.kubernetes.io/name: gitea-runner
      app.kubernetes.io/component: runner
  policyTypes:
  - Ingress
  ingress:
  # Allow all pods in gitea-e2e (job pods) to reach the cache on port 3000.
  - from:
    - namespaceSelector:
        matchLabels:
          kubernetes.io/metadata.name: gitea-e2e
    ports:
    - protocol: TCP
      port: 3000
EOF
```

After applying, confirm a job pod can reach the cache (the E2E harness does
this automatically via a `curl` probe before running jobs):

```bash
kubectl -n gitea-e2e run cache-probe --rm -it --restart=Never \
  --image=curlimages/curl -- \
  curl -sf http://gitea-runner-cache.gitea-e2e.svc.cluster.local:3000/ \
  && echo "cache reachable"
```

---

## DinD (optional, for `uses: docker://` steps)

To enable `uses: docker://` container steps, set in `configmap.yaml`:

```yaml
kubernetes:
  dind:
    enabled: true      # explicit opt-in; no auto-detect
    rootless: true     # default: lower privilege; set false for privileged
    image: docker:28-dind-rootless
```

When enabled, the backend injects an on-demand `docker:dind-rootless` native
sidecar (`restartPolicy: Always`) into each job pod. `using: docker` actions
(a `uses:` action whose `action.yml` sets `runs.using: docker`) are **not
supported in v1** and produce an explicit error — they are not silently
ignored.

---

## Customisation

| ConfigMap field | Default | Notes |
|---|---|---|
| `kubernetes.default_job_image` | `node:20-bookworm` | Image for bare `run:` jobs with no `container: image:` |
| `kubernetes.pod_ready_timeout` | `10m` | Increase for slow image pulls |
| `kubernetes.pod_cleanup` | `always` | When to delete completed job pods: `always`, `on-failure`, or `never` (see [Job pod retention](#job-pod-retention-auditing)) |
| `kubernetes.max_concurrent_jobs` | `0` (unlimited) | Concurrent job-pod cap in k8s mode; `0` = unlimited (see [Job concurrency](#job-concurrency-kube-scheduler-handles-the-backlog)) |
| `kubernetes.runtime_class_name` | _(unset)_ | `RuntimeClass` enforced on every job pod, e.g. `gvisor`/`kata` (see [Sandboxed runtimes](#sandboxed-runtimes-gvisor--kata)) |
| `kubernetes.resources.job` | `250m/256Mi` req, `2/2Gi` lim | Per-job container |
| `kubernetes.resources.service` | `100m/128Mi` req, `1/1Gi` lim | Per-service sidecar |
| `kubernetes.resources.dind` | `500m/1Gi` req, `2/4Gi` lim | DinD sidecar (≥1Gi required) |
| `runner.capacity` | `4` | Max concurrent jobs for **docker/host** backends; **bypassed** in k8s mode (use `kubernetes.max_concurrent_jobs`) |
| `GITEA_INSTANCE_URL` | `http://gitea.192.168.2.190.nip.io` | Set to your Gitea URL |

---

## Job concurrency (kube-scheduler handles the backlog)

In Kubernetes mode the runner only **creates** job pods and execs into them — the
**kube-scheduler** queues and places the real workloads. So the small
host-process `runner.capacity` (which bounds the docker/host backends) is an
artificial bottleneck: it makes the runner queue jobs internally instead of
letting Kubernetes pend them. The pod backend therefore uses
`kubernetes.max_concurrent_jobs` instead:

| Value | Behaviour |
|---|---|
| `0` (default) | **Unlimited.** The runner accepts every assigned job and creates a pod for it; the kube-scheduler pends whatever the cluster cannot place yet. |
| _N_ > 0 | Optional client-side cap: at most _N_ in-flight job pods from this runner. |

`runner.capacity` is **ignored** when `kubernetes.enabled: true`; docker/host
runners are unaffected.

> **Let the cluster be the backpressure, not the runner.** True-unlimited keeps
> one runner goroutine **and** one Kubernetes API watch alive per in-flight job,
> so a very large fan-out adds API-server load. The correct way to bound a build
> namespace is a **ResourceQuota** (and/or `LimitRange`), which makes the
> scheduler reject/pend excess pods — far better than throttling the runner to a
> tiny number and serialising the queue inside it.

Example: cap the `gitea-e2e` namespace at 50 pods and a CPU/memory budget so an
unbounded runner cannot exhaust the cluster:

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: gitea-e2e-quota
  namespace: gitea-e2e
spec:
  hard:
    pods: "50"
    requests.cpu: "20"
    requests.memory: 40Gi
    limits.cpu: "80"
    limits.memory: 160Gi
```

With the quota in place you can safely run `kubernetes.max_concurrent_jobs: 0`
and let pods pend when the budget is full.

---

## Job pod retention (auditing)

By default a job pod (and its per-job Secret/ConfigMap) is deleted as soon as the
job finishes. To keep completed pods around for post-mortem inspection, set
`kubernetes.pod_cleanup`:

| Value | Behaviour |
|---|---|
| `always` (default) | Delete the pod after every job. |
| `on-failure` | Delete only when the job **succeeded**; **keep** pods for **failed** jobs. |
| `never` | **Keep** every pod, regardless of result. |

When a pod is kept, the job container stays alive running its `sleep`
entrypoint, so you can inspect the workspace and rerun commands:

```bash
# find retained pods
kubectl -n gitea-e2e get pods -l actions.gitea.com/retain=true

# inspect the workspace of a retained job
kubectl -n gitea-e2e exec -it <pod-name> -c job -- sh
```

Retained pods are labelled `actions.gitea.com/retain=true`. The runner's
startup dangling-owner sweep **skips** pods carrying this marker, so a runner
restart will not reap them.

> **Cost of retention.** A retained pod keeps consuming its requested CPU and
> memory **and** counts against the namespace `ResourceQuota` (see [Job
> concurrency](#job-concurrency-kube-scheduler-handles-the-backlog)) until it is
> removed. Under `never` (or `on-failure` on a flaky pipeline) retained pods
> accumulate and will eventually exhaust the quota and pend new jobs. Prune them
> on a schedule:
>
> ```bash
> kubectl -n gitea-e2e delete pods -l actions.gitea.com/retain=true
> ```

---

## Sandboxed runtimes (gVisor / Kata)

Each job is one Pod, so a Kubernetes [`RuntimeClass`](https://kubernetes.io/docs/concepts/containers/runtime-class/)
is all that's needed to run untrusted job code under a sandboxed runtime —
**gVisor** (`runsc`, a userspace kernel) or **Kata Containers** (lightweight
VMs) — instead of the node default (`runc`). The runner sets
`pod.spec.runtimeClassName` from config; everything else (services as native
sidecars, exec-in steps, ephemeral workspace) is unchanged.

**Prerequisite — node side.** `RuntimeClass` is a contract with the node: the
runtime handler binary (runsc / kata) and a matching containerd/CRI-O runtime
must already be installed on every node that may schedule these pods. Install
the binaries, register the runtime name, then apply the `RuntimeClass` objects:

```bash
kubectl apply -f runtimeclass.yaml   # defines `gvisor` and `kata`
```

**Enable it** — one RuntimeClass is enforced on every job pod the runner creates:

```yaml
kubernetes:
  runtime_class_name: gvisor
```

### Why this is a per-runner policy, not a per-job label

The sandbox confines code you don't trust — the workflow author. So the runtime
is deliberately an **operator/cluster policy applied to every job this runner
accepts**, not something a workflow selects: if a job could pick its own
`runs-on: …-gvisor` vs `…-native`, a malicious or compromised workflow would
simply choose the least-isolated option and opt out of the sandbox. A boundary
the attacker can choose is not a boundary.

To offer **different isolation levels**, run **one runner per trust tier** — each
with its own `runtime_class_name` (and ideally its own namespace / node pool) —
and route jobs with `runs-on`. The label then selects the *runner* (its native
purpose), and each runner enforces its own runtime on everything it picks up:

```
runner-untrusted   runtime_class_name: gvisor      labels: [k8s, gvisor]
runner-trusted     runtime_class_name: <unset>     labels: [k8s, native]
```

```yaml
# .gitea/workflows/ci.yml — routed to the gVisor-enforcing runner
jobs:
  build:
    runs-on: [k8s, gvisor]
```

Notes:

- If only some nodes provide the sandbox, give each `RuntimeClass` a
  `scheduling.nodeSelector`/`tolerations` (commented examples in
  `runtimeclass.yaml`) so the scheduler never places a sandboxed pod on a node
  that cannot honor it — otherwise the pod stays `Pending`.
- A **privileged** DinD sidecar (`kubernetes.dind.rootless: false`) is generally
  incompatible with gVisor; prefer rootless DinD, or run `uses: docker://` jobs
  on a `runc`-backed runner.
- A misnamed or absent `RuntimeClass` surfaces at pod admission (creation error
  or a `Pending` pod with a clear event), so misconfiguration fails visibly
  rather than silently downgrading isolation.

---

## Security considerations

The native pod backend deliberately runs each job as an unprivileged Pod with a
least-privilege ServiceAccount. A few behaviours nonetheless warrant operator
attention:

- **DinD daemon serves plaintext TCP on pod loopback.** When
  `kubernetes.dind.enabled: true`, the docker-in-docker sidecar listens on
  `tcp://localhost:2375` with TLS disabled (`DOCKER_TLS_CERTDIR=""`). This is
  only reachable over the job Pod's loopback interface (the shared pod network),
  so it is safe **provided no untrusted container is co-scheduled into the job
  Pod**. Do not add arbitrary user-controlled sidecars to a job Pod that also
  runs DinD.
- **Translated `uses: docker://` steps run with `--network host`.** Inside the
  DinD daemon, the translated `docker run` uses `--network host`, so the step
  container shares the DinD network namespace rather than being isolated on its
  own bridge. Treat container-step images with the same trust level as the job.
- **`kubernetes.dind.rootless: false` makes every job's DinD privileged.**
  The default is rootless (`true`). Setting it to `false` injects a
  **privileged** `docker:dind` sidecar into *every* job Pod, which materially
  weakens isolation (a compromised job could escalate via the privileged
  daemon). Only disable rootless when a specific image genuinely cannot build
  under rootless DinD, and prefer per-runner scoping in that case.
