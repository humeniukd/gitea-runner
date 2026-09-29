# Jobs as Kubernetes pods

A label with the `kubernetes` schema runs each of its jobs as a pod, without a Docker daemon:

```text
ubuntu-latest:kubernetes://docker.gitea.com/runner-images:ubuntu-latest
```

The pod holds a `job` container from the label's image or the job's `container.image`, plus one container per service. It is deleted when the job ends, or with the runner's pod.

## Setup

In a cluster, the runner uses its pod's service account and namespace. [`examples/kubernetes/kubernetes-pods.yaml`](../examples/kubernetes/kubernetes-pods.yaml) deploys a runner with these permissions:

| Resource | Verbs |
| --- | --- |
| `pods` | `create`, `get`, `delete` |
| `pods/exec` | `create`, `get` |
| `pods/log` | `get` |
| `secrets` | `create`, `delete` |

Outside a cluster, it uses the current context of `$KUBECONFIG` or `~/.kube/config`. The cluster must run Kubernetes 1.30 or later, and job images need `sh`, `tar`, `sleep` and `env`.

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

Job pods get no service account token or service links unless the template sets `automountServiceAccountToken` or `enableServiceLinks`.

## Services

Services share the pod's network, so a step reaches one by its id or `localhost` on the port it listens on, and two services cannot use the same port. Service `env` goes into a per-job secret, so pass secrets that way, as `command` and `options` are visible in the pod spec. Of `options`, only `--health-*` is accepted, and `--health-cmd` becomes a readiness probe the job waits for.

## Unsupported

`docker://` steps, `runs.using: docker` actions, the job container's `options`, `volumes` and `credentials`, and services' `volumes` and `credentials` fail the job. Pull private images with `imagePullSecrets` in `pod_template`. Docker settings such as `container.options` and `container.network` do not apply.
