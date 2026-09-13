# Cache server

The runner provides a disk-backed cache for `actions/cache` and actions that use
its cache API. It can run inside each runner or as a standalone service shared by
multiple runners. Both cache API versions use the same stored entries.

## Embedded cache

Caching is enabled by default. The runner starts its own server and registers each
job automatically. A minimal configuration with a fixed address is:

```yaml
cache:
  enabled: true
  dir: /var/lib/gitea-runner/cache
  host: runner.example.internal
  port: 8088
```

`host` is the address advertised to job containers, not a bind address. The server
listens on all interfaces. With no host configured, the runner detects its outbound
IP; port `0` selects a free port. If the runner runs in Docker, publish a fixed
cache port and use an address reachable from the job containers, including jobs
on separate Docker networks. `localhost` inside a job container refers to that
container.

The default directory is `$HOME/.cache/actcache`. Mount it on persistent storage
if caches should survive replacement of the runner container.

## Shared cache server

Generate a shared secret and provision the same value to the cache server and
each trusted runner. For example, create a file with restrictive permissions:

```sh
umask 077
openssl rand -hex 32 > cache-secret
```

Configure the standalone server in `cache-server.yaml`:

```yaml
cache:
  dir: /var/lib/gitea-runner/cache
  host: cache.example.internal
  port: 8088
  external_secret_file: /run/secrets/cache-secret
  retention: 168h
  repo_size_limit: 10GB
  size_limit: 100GB
  sweep_interval: 1h
```

Start it with:

```sh
gitea-runner --config cache-server.yaml cache-server
```

The command requires no runner registration. Non-empty `--dir` and `--host`, and
non-zero `--port`, override the corresponding configuration values.

On each runner, configure:

```yaml
cache:
  enabled: true
  external_server: https://cache.example.internal/
  external_secret_file: /run/secrets/cache-secret
```

This HTTPS example assumes a reverse proxy forwarding to the cache server on
port `8088`. The service itself speaks HTTP. The external URL must be reachable
from both runner processes and job environments. Preserve request paths, query
strings, authorization headers, and upload bodies at the proxy; allow large
uploads and suitable transfer timeouts.

`external_secret` can hold the secret directly instead of using a file. Setting
both secret options is an error. Local directory, host, port, and eviction settings
on a runner are ignored when it uses an external server; configure eviction on
the shared server. `cache.enabled: false` disables the embedded server and external
job registration.

## Authentication and scope

Runners use the shared secret to register and revoke per-job credentials through
`/_internal/register` and `/_internal/revoke`. Cache management requests use the
job's bearer token, and entries are scoped to its repository. Keep the shared
secret available only to trusted runner operators: it authorizes credential
registration. Repository names are the scope, so use separate servers for
independent Gitea instances whose repository names can overlap.

Blob transfers use signed URLs. Download URLs expire after ten minutes; v2 upload
URLs expire after one hour. Treat these URLs as credentials. Use HTTPS through a
reverse proxy when traffic crosses an untrusted network.

The server persists its signing key in `.secret` inside the cache directory. This
is separate from the shared runner secret. Credentials registered by jobs live in
memory and must be registered again after a server restart.

## Retention and disk usage

| Setting | Default | Meaning |
| --- | --- | --- |
| `cache.retention` | `168h` | Remove entries that have not been accessed within this window. `0s` disables age expiry. |
| `cache.repo_size_limit` | `10GB` | Limit completed entries per repository. `0` disables the limit. |
| `cache.size_limit` | `0` | Limit completed entries across repositories. `0` disables the limit. |
| `cache.sweep_interval` | `1h` | Minimum interval between sweeps; non-positive values use the default. |

Size units are binary and case-insensitive: `GB` and `GiB` both mean 1024³ bytes.
Size eviction removes least recently used entries first. Entries larger than a
limit are eligible for removal themselves. Recent accesses are protected for
eleven minutes, so limits can be exceeded temporarily. Access timestamps are
written at most once per minute per entry to reduce database writes.

Sweeps run at startup and are triggered by authenticated cache traffic, subject
to the sweep interval. They remove stalled incomplete uploads, expired entries,
superseded entries, and excess data. Upload reservations become stale after five
minutes without a recorded access. Completing an upload also checks repository,
total-size, and free-space limits.

The cache keeps a default free-space floor of 1024 MiB on its volume. When health
checks are enabled, `health_check.min_free_disk_space_mb` sets that floor instead.
This safeguard remains active when retention and size limits are disabled.
Recently used entries and unrelated files can prevent it from reclaiming enough
space; it is not a hard filesystem quota.

## Storage and troubleshooting

The cache directory contains `bolt.db` metadata, blobs under `cache/`, and the
`.secret` signing key. One server owns the directory for its entire lifetime.
Do not point multiple running processes at the same directory; share the service
through `external_server` instead. Stop the service before copying or manually
removing its storage, and preserve all three parts together when backing it up.

For connection failures, check reachability from inside a job container and
verify Docker port publishing and the advertised host. For authentication
failures, check that the runner and server use the same secret and that the proxy
preserves headers. A database-open timeout can mean another process still owns
the cache directory. A full volume can require removing unrelated data even when
cache eviction is enabled.

`cache.v2` defaults to `true`; `runner.patch_actions` controls the compatibility
patches applied to action bundles. The action repository cache controlled by
`cache.offline_mode` is separate from the workflow cache stored by this service.
