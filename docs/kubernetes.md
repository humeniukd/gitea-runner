# Jobs as Kubernetes pods

A `kubernetes` label runs each job as a pod, without a Docker daemon:

```text
ubuntu-latest:kubernetes://docker.gitea.com/runner-images:ubuntu-latest
```

The pod has a `job` container from the label's image or the job's `container.image`, plus one container per service.

## Setup

In a cluster, the runner uses its pod's service account and namespace. [`examples/kubernetes/kubernetes-pods.yaml`](../examples/kubernetes/kubernetes-pods.yaml) deploys a runner with these permissions:

| Resource | Verbs |
| --- | --- |
| `pods` | `create`, `get`, `delete` |
| `pods/exec` | `create`, `get` |
| `pods/log` | `get` |
| `secrets` | `create`, `delete` |

Outside a cluster, it uses the current context of `$KUBECONFIG` or `~/.kube/config`. The cluster must run Kubernetes 1.30 or later, and job images need `sh`, `tar`, `tail` and `env`.

## Configuration

```yaml
kubernetes:
  kubeconfig: "" # instead of the service account
  namespace: "" # defaults to the service account's, the kubeconfig context's, or default
  pod_template: # merged under each job pod, a container named job merges into the job container
    spec:
      nodeSelector:
        kubernetes.io/arch: amd64
      containers:
        - name: job
          resources:
            limits:
              memory: 4Gi
```

Job pods mount no service account token unless the template sets `automountServiceAccountToken: true`.

## Services

Services share the pod's network, so a step reaches one by its id or `localhost`, and two services cannot listen on the same port. An id is lowercased and has characters other than letters, digits and `-` replaced by `-`, so the service `my_db` is reached as `my-db`. Pass secrets to services through `env`, as `command` is visible in the pod spec. Of `options`, only `--health-*` is supported, and `--health-cmd` becomes a readiness probe.

## Unsupported

Docker actions, `docker://` steps, and the `options`, `volumes` and `credentials` of the job container fail the job, as do `volumes` and `credentials` of services. Run such workflows on a [Docker in Docker](../examples/kubernetes/README.md) runner. Pull private images with `imagePullSecrets` in `pod_template`.
