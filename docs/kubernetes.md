# Jobs as Kubernetes pods

A label with the `kubernetes` schema runs each of its jobs as a pod, without a Docker daemon:

```text
ubuntu-latest:kubernetes://docker.gitea.com/runner-images:ubuntu-latest
```

The pod holds a `job` container, created from the label's image or the job's `container.image`, and one container per service. The pod is deleted when the job ends, and with the runner's pod when the runner runs as a pod in the same namespace.

## Setup

The runner talks to the Kubernetes API as the service account of its pod, and creates job pods in that account's namespace. [`examples/kubernetes/kubernetes-pods.yaml`](../examples/kubernetes/kubernetes-pods.yaml) deploys a runner with the permissions it needs:

| Resource | Verbs |
| --- | --- |
| `pods` | `create`, `get`, `delete` |
| `pods/exec` | `create`, `get` |
| `pods/log` | `get` |
| `secrets` | `create`, `delete` |

The `env` of a job's services goes into a secret the runner deletes with the pod.

Outside a cluster, the runner uses the current context of the first file in `$KUBECONFIG`, or of `~/.kube/config`, with a token or a client certificate. `gitea-runner exec --image kubernetes://<image>` runs a workflow the same way.

The cluster must run Kubernetes 1.30 or later. Job images need `sh`, `tar` and coreutils such as `sleep` and `env`, which the [runner images](https://gitea.com/gitea/runner-images) have.

## Configuration

```yaml
kubernetes:
  # Kubeconfig to use instead of the service account of the runner's pod.
  kubeconfig: ""
  # Namespace of job pods. Empty uses the service account's, else the kubeconfig context's, else default.
  namespace: ""
  pod_template:
    spec:
      nodeSelector:
        kubernetes.io/arch: amd64
      containers:
        - name: job
          resources:
            limits:
              memory: 4Gi
```

The runner merges its pod onto `pod_template` map key by map key, replacing list items of the same `name`, except that a container named `job` is merged into the job container. Job pods mount no service account token and get no environment variables for the namespace's services unless the template sets `automountServiceAccountToken` or `enableServiceLinks`.

## Services

A step reaches a service by its id or `localhost`, on the port the service listens on, so two services cannot listen on the same port. A service id must be a valid DNS label other than `job`. `job.services.<id>.ports` maps each declared port to itself.

A service's `options` may only hold `--health-cmd`, `--health-interval`, `--health-timeout`, `--health-retries` and `--health-start-period`, of which the last two are ignored. `--health-cmd` becomes a readiness probe the job waits for. A `HEALTHCHECK` in the service's image is not used. A service's `command`, `entrypoint` and `--health-cmd` are readable in the pod spec, so pass secrets to services through `env`.

## Unsupported

These fail the job:

- `docker://` steps and actions with `runs.using: docker`
- `options`, `volumes` and `credentials` of the job container, `volumes` and `credentials` of services. Pull private images with `imagePullSecrets` in `pod_template`.

Docker settings of the runner config, such as `container.options`, `container.network` and `runner.tool_cache_mode`, do not apply to pods.
