# Changelog

## 3.3.2 (2026-08-31)

### ⚠️ Behavior changes

#### `NODE_OPTIONS` can no longer be set from a workflow (https://gitea.com/gitea/runner/pulls/1194)

Writing `NODE_OPTIONS` to `$GITHUB_ENV` or via `::set-env::` is refused with an error annotation, as on GitHub.

### Features

- `exec` now supports `--input name=value` and `--input-file <path>` for the inputs a workflow declares under `workflow_dispatch` or `workflow_call` (https://gitea.com/gitea/runner/pulls/1173)

### Bug fixes

- Register every encoded spelling of a secret, so a value like `base64("user:$TOKEN")` no longer prints in the clear (https://gitea.com/gitea/runner/pulls/1194)
- Report the status a job actually reached: a `continue-on-error` step logs success, a cancelled job reports cancelled, a failing `if:` reports failure (https://gitea.com/gitea/runner/pulls/1194)
- Resolve `${{ matrix.* }}` and `${{ strategy.* }}` inside composite actions, and stop composite inputs leaking into nested actions as `INPUT_*` (https://gitea.com/gitea/runner/pulls/1194)
- Stop `container.env` overriding job env and `$GITHUB_ENV` writes (https://gitea.com/gitea/runner/pulls/1194)
- Pass `--preserve-symlinks-main` to Node actions, so an action comparing `process.argv[1]` against its realpath runs instead of skipping itself (https://gitea.com/gitea/runner/pulls/1202)

**Full changelog**: https://gitea.com/gitea/runner/compare/v3.3.1...v3.3.2

## 3.3.1 (2026-08-26)

### ⚠️ Behavior changes

- A workflow's `container.options` may no longer set `--env-file` or `--label-file`, which name files read on the runner, nor carry a volume driver option, which turns a name `valid_volumes` allows into a bind of any host path. Both still serve the runner's own options (https://gitea.com/gitea/runner/pulls/1151)
- A bare `--env NAME` no longer resolves the value from the runner's own environment, for every source. Use `runner.envs` or `runner.env_file` to pass a variable on (https://gitea.com/gitea/runner/pulls/1151)
- A job output whose value carries a secret is skipped with a warning instead of sent, matching GitHub, so a downstream `needs.<job>.outputs.<name>` reading it is empty (https://gitea.com/gitea/runner/pulls/1188)

### Bug fixes

- Keep the runner's own `container.options` when privileged mode is off. The host-escape filter dropped the administrator's options along with the workflow's, so a setup needing `--device` or `--security-opt` from the config file had no way to get them (https://gitea.com/gitea/runner/pulls/1151)
- Mask secrets on every path they leave a job: uploaded log rows, the on-disk `job.log`, the runner's own log, debug stdout, job summaries, job outputs, and the job name that becomes a container name. `ACTIONS_STEP_DEBUG` and `ACTIONS_RUNNER_DEBUG` are never masked, as on GitHub (https://gitea.com/gitea/runner/pulls/1188)
- Keep a step's own `with:` values out of its `inputs` context, so a colliding `with:` key or an `INPUT_`-shaped variable from `env:` no longer flips `if:` conditions or forges an input (https://gitea.com/gitea/runner/pulls/1192)
- Honor volumes declared on service containers under the configured `valid_volumes` policy instead of dropping them silently (https://gitea.com/gitea/runner/pulls/1186)
- Report step log ranges with the log flush, so rows the server takes between two state reports no longer land under no step or under "Complete job" (https://gitea.com/gitea/runner/pulls/1189)
- Fail the run when matrix expansion fails, instead of reporting success without running anything (https://gitea.com/gitea/runner/pulls/1187)

### Dependencies

- Update the Go toolchain to 1.27.0, including the 1.26.6 security fixes. Building from source now needs Go 1.27 (https://gitea.com/gitea/runner/pulls/1183, https://gitea.com/gitea/runner/pulls/1185)
- Update dependencies, docker to 29.7.2 (https://gitea.com/gitea/runner/pulls/1185)

**Full changelog**: https://gitea.com/gitea/runner/compare/v3.3.0...v3.3.1

## 3.3.0 (2026-08-21)

### Features

- `log.job.dir` writes a copy of every task's log to that directory, the rows exactly as Gitea received them, bounded by `log.job.retention` and `log.job.max_size`. Off by default (https://gitea.com/gitea/runner/pulls/1165)

### Enhancements

- `runner.default_image` sets the image for a job whose `runs-on` matches none of the runner's labels, which includes every job that sets no `runs-on` at all. Such a job now runs on the host when no docker daemon is reachable, which matters most to a host-mode runner (https://gitea.com/gitea/runner/pulls/1164)

### Bug fixes

- Stop a failed step disabling the toolkit patch, which made every later `actions/upload-artifact` run fail with `GHESNotSupportedError`. `cache.v2` no longer decides whether the patch is applied, and the new `runner.patch_actions` turns it off (https://gitea.com/gitea/runner/pulls/1177)

**Full changelog**: https://gitea.com/gitea/runner/compare/v3.2.0...v3.3.0

## 3.2.0 (2026-08-19)

### ⚠️ Behavior changes

#### The tool cache is no longer shared between jobs (https://gitea.com/gitea/runner/pulls/1171)

The shared tool cache volume is not concurrency-safe, concurrent jobs can write it at the same time and corrupt it, for example `setup-go` with an explicit Go version. The new `runner.tool_cache_mode` (and `--tool-cache-mode` for `exec`) defaults to `none`, where a job uses what its image ships in the tool cache and discards what it installs, the same as GitHub Actions with its preinstalled images. Set it to `shared` to keep the single volume every job reuses.

### Features

- Size-based cache eviction via `cache.repo_size_limit` and `cache.size_limit`, and `cache.retention` now counts from last access instead of creation (https://gitea.com/gitea/runner/pulls/1170)
- Support GitHub's `$/` self-repository prefix in action `uses:` (https://gitea.com/gitea/runner/pulls/1150)

### Enhancements

- Show `workflow_dispatch` inputs in the "Set up job" section (https://gitea.com/gitea/runner/pulls/1167)

### Bug fixes

- Bound every RPC to Gitea and stop failing silently, jobs no longer go quiet until Gitea reaps them (https://gitea.com/gitea/runner/pulls/1174)
- Build cache job URLs on the address its runner registered, so jobs reaching the cache server through a reverse proxy are no longer sent to its internal address (https://gitea.com/gitea/runner/pulls/1153)

### Refactors

- Move `act/model` and `act/exprparser` to `gitea.dev/actionslib`, both packages are gone from this module (https://gitea.com/gitea/runner/pulls/1143)

### Dependencies

- Update `github.com/go-git/go-git/v5` to v5.19.2 (security) (https://gitea.com/gitea/runner/pulls/1156)
- Update dependencies, docker to 29.7.1 (https://gitea.com/gitea/runner/pulls/1160)

**Full changelog**: https://gitea.com/gitea/runner/compare/v3.1.0...v3.2.0

## 3.1.0 (2026-08-06)

### ⚠️ Behavior changes

#### `::set-env::` and `::add-path::` need an opt-in (https://gitea.com/gitea/runner/pulls/1109)

The runner refuses both commands, as GitHub has done since 2020. Set `ACTIONS_ALLOW_UNSECURE_COMMANDS` on the step or the job to allow them again.

### Features

- Wait for healthy services and fill the `job` context (https://gitea.com/gitea/runner/pulls/1107)
- `config` command to edit config files (https://gitea.com/gitea/runner/pulls/1140)
- Render annotation locations, so `::error file=main.go,line=12::msg` shows as `::error::main.go:12: msg` (https://gitea.com/gitea/runner/pulls/1109)

### Enhancements

- Improve config, comment out values in example file, add `config init` (https://gitea.com/gitea/runner/pulls/1145)

### Bug fixes

- Evaluate each `${{ }}` part on its own (https://gitea.com/gitea/runner/pulls/1146)
- Trim whitespace from register inputs (https://gitea.com/gitea/runner/pulls/1147)

**Full changelog**: https://gitea.com/gitea/runner/compare/v3.0.2...v3.1.0

## 3.0.2 (2026-08-03)

### Bug fixes

- Serve the whole results service from the cache server (https://gitea.com/gitea/runner/pulls/1141)

### Dependencies

- Update dependencies (https://gitea.com/gitea/runner/pulls/1138)

**Full changelog**: https://gitea.com/gitea/runner/compare/v3.0.1...v3.0.2

## 3.0.1 (2026-08-02)

### Bug fixes

- Coerce every expression value kind to a string and drop the `format()` panic — https://gitea.com/gitea/runner/pulls/1135

### Dependencies

- Drop `Masterminds/semver` and `pkg/errors`, re-sync the `docker/cli` copies — https://gitea.com/gitea/runner/pulls/1132

**Full changelog**: https://gitea.com/gitea/runner/compare/v3.0.0...v3.0.1

## 3.0.0 (2026-07-31)

### ⚠️ Breaking changes

#### Cache service v2 (https://gitea.com/gitea/runner/pulls/1110)

The runner serves the cache service v2 API alongside v1, on by default and working with `external_server`. To reach it, the official cache actions are patched before they run which fixes the official artifact actions, so `actions/upload-artifact` and `actions/download-artifact` v4.4.0 and newer work as published and the `gitea-upload-artifact` fork can be dropped. The v2 API can be disabled via `cache.v2: false` in the configuration.

#### Only one runner process per registration file (https://gitea.com/gitea/runner/pulls/1099)

Two daemons started from the same `.runner` file looked like one runner to Gitea and cancelled each other's jobs. The file is now locked at startup, so every runner needs its own `runner.file` and its own registration.

#### Container options that reach the host now require privileged mode (https://gitea.com/gitea/runner/pulls/1058)

A workflow could use `container.options` to reach the runner host, for example by joining its process namespace, even with privileged mode off. Those options are now dropped with a warning, turn on `container.privileged` to keep them working.

### Features

- Add cache service v2 and toolkit patches (https://gitea.com/gitea/runner/pulls/1110)
- Add job hooks `runner.hooks.job_started` and `runner.hooks.job_completed` (https://gitea.com/gitea/runner/pulls/1111)
- Propagate proxy variables to jobs, services and builds (https://gitea.com/gitea/runner/pulls/1112)
- Mask secrets that reach the log in an encoded form (https://gitea.com/gitea/runner/pulls/1108)

### Bug fixes

- Resolve symlinked container paths before building tar entries (https://gitea.com/gitea/runner/pulls/1130)
- Stop leaking per-job docker networks (https://gitea.com/gitea/runner/pulls/1124)
- Fix panics, enable the `forcetypeassert` lint (https://gitea.com/gitea/runner/pulls/1123)
- Allow relocating the tool cache and mounting over runner paths (https://gitea.com/gitea/runner/pulls/1122)
- Strip leading slash from mkdir tarball in `CopyTarStream` (https://gitea.com/gitea/runner/pulls/1129)

**Full changelog**: https://gitea.com/gitea/runner/compare/v2.3.0...v3.0.0

## 2.3.0 (2026-07-28)

### Features

- Run pre-entrypoint and post-entrypoint of docker actions — https://gitea.com/gitea/runner/pulls/1106
- Report runner name, environment, workspace and debug to jobs — https://gitea.com/gitea/runner/pulls/1105
- Support `--platform` and `--pull` in `container.options` — https://gitea.com/gitea/runner/pulls/1104
- Support reading `cache.external_secret` from a file — https://gitea.com/gitea/runner/pulls/1100

### Bug fixes

- Escape command data the runner writes itself — https://gitea.com/gitea/runner/pulls/1120
- Clean service containers after failed job setup — https://gitea.com/gitea/runner/pulls/1066

### Dependencies

- Update `go.opentelemetry.io/otel` to v1.44.0 (security) — https://gitea.com/gitea/runner/pulls/1115
- Update `golang.org/x/net` to v0.56.0 (security) — https://gitea.com/gitea/runner/pulls/1116
- Update `github.com/prometheus/client_golang` to v1.24.0 — https://gitea.com/gitea/runner/pulls/1113
- Update docker docker tag to v29.6.2 — https://gitea.com/gitea/runner/pulls/1101
- Update `actions/setup-go` action to v7 — https://gitea.com/gitea/runner/pulls/1102
- Pin dependencies — https://gitea.com/gitea/runner/pulls/1117

**Full changelog**: https://gitea.com/gitea/runner/compare/v2.2.0...v2.3.0

## 2.2.0 (2026-07-22)

### Features

- Add runner health admission checks — https://gitea.com/gitea/runner/pulls/1090
- Report a GitHub-style "Set up job" section — https://gitea.com/gitea/runner/pulls/1089

### Bug fixes

- Repair the free-disk-space build on FreeBSD — https://gitea.com/gitea/runner/pulls/1098
- Stop host-mode jobs from leaking processes on Windows — https://gitea.com/gitea/runner/pulls/1080
- Classify a cancelled step as an interruption, not a failure — https://gitea.com/gitea/runner/pulls/1095
- Stop racing the daemon when removing containers — https://gitea.com/gitea/runner/pulls/1093
- Guard status-check functions against a nil job context — https://gitea.com/gitea/runner/pulls/1092

### Dependencies

- Update `github.com/mattn/go-isatty` to v0.0.23 — https://gitea.com/gitea/runner/pulls/1097
- Update `github.com/docker/cli` to v29.6.2 — https://gitea.com/gitea/runner/pulls/1096
- Update `actions/setup-node` action to v7 — https://gitea.com/gitea/runner/pulls/1094

**Full changelog**: https://gitea.com/gitea/runner/compare/v2.1.0...v2.2.0

## 2.1.0 (2026-07-16)

### Features

- Honour `GITEA_RUNNER_LABELS` on daemon start and accept labels containing a colon — https://gitea.com/gitea/runner/pulls/1085
- Add `--token-file` flag to the `register` command — https://gitea.com/gitea/runner/pulls/1076
- Add a `bug-report` subcommand, `exec --eventpath`, and `runner.set_act_env` — https://gitea.com/gitea/runner/pulls/1075

### Bug fixes

- Ignore blank lines and decode UTF-16 in the runner env files — https://gitea.com/gitea/runner/pulls/1084
- Stop service containers from clobbering the job container's credentials — https://gitea.com/gitea/runner/pulls/1083
- Accept natively typed boolean workflow inputs — https://gitea.com/gitea/runner/pulls/1087
- Drop action outputs whose value exceeds the size limit — https://gitea.com/gitea/runner/pulls/1070

### Docs

- Explain labels, the docker image cache volume, and the dind-rootless UID — https://gitea.com/gitea/runner/pulls/1086

### Dependencies

- Update `golang.org/x/text` to v0.40.0 — https://gitea.com/gitea/runner/pulls/1091
- Update `golang.org/x/term` to v0.45.0 — https://gitea.com/gitea/runner/pulls/1081

**Full changelog**: https://gitea.com/gitea/runner/compare/v2.0.1...v2.1.0

## 2.0.1 (2026-07-09)

### Bug fixes

- Skip service containers with an empty image — https://gitea.com/gitea/runner/pulls/1074
- Prevent exponential growth of RunContext masks in composite actions — https://gitea.com/gitea/runner/pulls/1059
- Attach task token when cloning actions from self-hosted instance on a different host — https://gitea.com/gitea/runner/pulls/1056
- Install `nftables` in dind images to silence `nft` cleanup errors — https://gitea.com/gitea/runner/pulls/1064

### Tests & docs

- Enhance test coverage + CI — https://gitea.com/gitea/runner/pulls/1055
- Clarify cache reachability for dockerized runners — https://gitea.com/gitea/runner/pulls/1069

### Dependencies

- Update `github.com/docker/cli` to v29.6.1 — https://gitea.com/gitea/runner/pulls/1060
- Update docker docker tag to v29.6.1 — https://gitea.com/gitea/runner/pulls/1063
- Update `golang.org/x/sys` to v0.47.0 — https://gitea.com/gitea/runner/pulls/1073

**Full changelog**: https://gitea.com/gitea/runner/compare/v2.0.0...v2.0.1

## 2.0.0 (2026-06-30)

### ⚠️ Breaking changes

#### `DOCKER_USERNAME` / `DOCKER_PASSWORD` are no longer used implicitly for image pulls

See https://gitea.com/gitea/runner/pulls/1007

The runner previously treated secrets named `DOCKER_USERNAME` and `DOCKER_PASSWORD` as implicit Docker pull credentials, attaching them to every job- and step-container image pull unconditionally. Since these credentials carry no registry information, they were even sent to Docker Hub when pulling public images, causing those pulls to fail with authentication errors. This behavior has been removed.

Workflows or runner setups that relied on `DOCKER_USERNAME` / `DOCKER_PASSWORD` being implicitly used for Docker image pulls must migrate to an explicit authentication mechanism:

- For private **job container** images, use `container.credentials`:

```yaml
jobs:
  build:
    container:
      image: registry.example.com/image:tag
      credentials:
        username: ${{ secrets.REGISTRY_USERNAME }}
        password: ${{ secrets.REGISTRY_PASSWORD }}
```

- For private **service container** images, use service `credentials`.
- For private `uses: docker://...` or private Docker actions, configure Docker authentication in the runner environment before the job starts (e.g. run `docker login` on the runner host).

`DOCKER_USERNAME` and `DOCKER_PASSWORD` can still be used as ordinary workflow secrets, for example with `docker/login-action` before pushing images.

### Features

- Shallow clone action repositories — https://gitea.com/gitea/runner/pulls/1053
- Enable `jobs.<job_id>.timeout-minutes` and `jobs.<job_id>.continue-on-error` — https://gitea.com/gitea/runner/pulls/1032
- Add optional `runner.post_task_script` hook after task cleanup — https://gitea.com/gitea/runner/pulls/1026
- Support `ssh://` action URLs — https://gitea.com/gitea/runner/pulls/1035
- IPv6 options for network container creation — https://gitea.com/gitea/runner/pulls/1029
- Complete runner-side cancellation handling — https://gitea.com/gitea/runner/pulls/1016
- Upload job summary when supported — https://gitea.com/gitea/runner/pulls/917

### Bug fixes

- Composite nested `uses` clone token — https://gitea.com/gitea/runner/pulls/1041
- Namespace local docker action image tags per repository — https://gitea.com/gitea/runner/pulls/1051
- Run `always()`/`cancelled()` and post steps correctly on cancellation — https://gitea.com/gitea/runner/pulls/1043
- Guard `SetJobError` against missing job-error container — https://gitea.com/gitea/runner/pulls/1050
- Interpolate job `container.volume` — https://gitea.com/gitea/runner/pulls/1036
- Kill Unix step process group on cancel to avoid hang — https://gitea.com/gitea/runner/pulls/1025
- Prevent loss of step log output at end of step — https://gitea.com/gitea/runner/pulls/1028
- Bound host-environment cleanup and reclaim leaked scratch dirs — https://gitea.com/gitea/runner/pulls/1024
- Do not update cached actions with stale origin URL — https://gitea.com/gitea/runner/pulls/1014

### Build & docs

- Cover Windows Go files in build — https://gitea.com/gitea/runner/pulls/1052
- Improve cache documentation — https://gitea.com/gitea/runner/pulls/1034
- Update docker-images docs — https://gitea.com/gitea/runner/pulls/1020

### Dependencies

- Update `github.com/moby/moby/client` to v0.5.0 — https://gitea.com/gitea/runner/pulls/1049
- Update `github.com/docker/cli` to v29.6.0 — https://gitea.com/gitea/runner/pulls/1045
- Update docker docker tag to v29.6.0 — https://gitea.com/gitea/runner/pulls/1044
- Update `actions/checkout` action to v7 — https://gitea.com/gitea/runner/pulls/1042
- Update `go.etcd.io/bbolt` to v1.5.0 — https://gitea.com/gitea/runner/pulls/1040
- Update `golang.org/x/term` to v0.44.0 — https://gitea.com/gitea/runner/pulls/1031
- Update alpine docker tag to v3.24 — https://gitea.com/gitea/runner/pulls/1030
- Update `golang.org/x/crypto` to v0.52.0 — https://gitea.com/gitea/runner/pulls/1027
- Update docker docker tag to v29.5.3 — https://gitea.com/gitea/runner/pulls/1021
- Update `golang.org/x/sys` to v0.46.0 — https://gitea.com/gitea/runner/pulls/1019
- Update `github.com/docker/cli` to v29.5.3 — https://gitea.com/gitea/runner/pulls/1018
- Update `github.com/opencontainers/selinux` to v1.15.1 — https://gitea.com/gitea/runner/pulls/1017
- Update `golang.org/x/sys` to v0.45.0 — https://gitea.com/gitea/runner/pulls/1012
- Update `code.gitea.io/actions-proto-go` to v0.5.0 — https://gitea.com/gitea/runner/pulls/1009

**Full changelog**: https://gitea.com/gitea/runner/compare/v1.0.8...v2.0.0

🤖 Release notes compiled by Opus 4.8.

## 1.0.8 (2026-06-02)

- fix(cleanup): kill Windows step process tree on cancel to avoid hang (#1011)

## 1.0.7 (2026-05-29)

### Changelog

- fix(container): re-validate cached container id before reuse (#1003)
- fix: support multiline secret masking (#1001)
- fix: restore global docker config dir and socket env in tests (#1004)
- fix: matrix-job data races + outputs, leaner offline test suite (#994)
- fix: deliver cancel ack and reap leftover Windows job processes (#996)
- fix(reporter): respect configured log level for job log forwarding (#989)

## 1.0.6 (2026-05-22)

### Changelog

- fix(deps): update module github.com/opencontainers/selinux to v1.15.0 (#990)
- chore: pin Docker base images to explicit versions (#992)
- chore(deps): update actions/setup-node action to v6 (#991)
- test: make TestRunEvent integration suite runnable locally (#987)
- ci: add PR title linting against Conventional Commits (#988)
- fix: clean up job network and container when container start fails (#986)

## 1.0.5 (2026-05-21)

### Changelog

- fix(deps): update module connectrpc.com/connect to v1.20.0 (#985)
- update docker cli to v29.5.2 (#984)
- feat: add cache.offline_mode to reuse cached actions (#966)
- Remove stale Gitea 1.20 compatibility shims (#978)
- fix(deps): update module github.com/docker/cli to v29.5.1+incompatible (#979)
- fix(deps): update module github.com/go-git/go-git/v5 to v5.19.1 (#980)

## 1.0.4 (2026-05-18)

### Changelog

- Fix token use with schemaless Gitea instance (#977)
- Add OCI `source` and `version` labels to images (#975)
- fix(parse_env_file): support env-file lines larger than 64 KiB (#974)
- Fix host cleanup, volume allowlist, cache upload, and action host edge cases (#970)
- Remove dead code from `act/` (#971)
- fix: Return if executors length is zero in ParallelExecutor (#960)
- feat: make pseudo-TTY allocation opt-in (#961)
- fix(deps): update module github.com/docker/cli to v29.5.0+incompatible (#969)
- Simplify kubernetes dind example allowing for default docker config in workflows (#709)
- chore(deps): bump `retry-go`, `golangci-lint`, `govulncheck` (#965)
- fix(deps): bump `docker` deps, switch to `moby/moby` (#943)
- fix: respect proxy env vars in runner client (#962)

## 1.0.3 (2026-05-12)

### Changelog

- fix: isolate per-task runner envs (#959)
- fix(deps): update module github.com/opencontainers/selinux to v1.14.1 (#955)
- fix(deps): update mergo to v1.0.2 (now `dario.cat/mergo`) (#954)
- fix: ensure `dbfs_data` is cleaned up after task completion (#952)
- fix(artifactcache): preserve cache key case to stop redundant uploads (#947)
- fix(deps): update module golang.org/x/term to v0.43.0 (#948)

## 1.0.2 (2026-05-08)

### Changelog

- fix: overwrite read-only files when copying action directories (#942)
- docs: add apparmor=rootlesskit in security_opt (#937)
- fix: serialize action-cache reads to prevent worktree race (#938)
- feat: remove emojis from runner logging, add `Starting job container` group (#940)

## 1.0.1 (2026-05-07)

### Changelog   
* dff63b3ecc4320d340d3b5511f91a65eb2e2e6c2 fix(deps): update module github.com/go-git/go-git/v5 to v5.19.0 (#934)   
* a5d9fe965120c0b398a6349005dd9a73eb6e511c fix(deps): update module github.com/opencontainers/selinux to v1.14.0 (#928)   
* d607f3b342e351197691a08795189972071d6374 test: clean up dead/stale fixtures and bump test container images (#932)   
* 5e59402fb2ad128c3f3008d4c41c7d3d03ed3f19 fix: re-fetch cached reusable workflow on every run (#930)   
* dfeb463904e768bb33086563d620a1925c278618 chore(deps): update docker docker tag to v29 (#924)

## 1.0.0 (2026-05-05)

This is the first release under the new name **`gitea-runner`**. The project was previously published as `act_runner` and has been renamed end-to-end — including the canonical repository, which now lives at [`gitea/runner`](https://gitea.com/gitea/runner) (the old `gitea/act_runner` URL still redirects). Operators upgrading from `act_runner` must update wrapper scripts, systemd units, container images, monitoring, and any code importing the Go module.

### Breaking changes

The full rename in #850 changed many user-visible names:

- **Repository**: `gitea.com/gitea/act_runner` → `gitea.com/gitea/runner`. The old URL redirects, but bookmarks, CI checkouts, and `go get` paths should be updated.
- **Binary**: `act_runner` → `gitea-runner`. Wrapper scripts, systemd units, init scripts, and docs referencing `act_runner` will break.
- **Docker image**: `gitea/act_runner` → `gitea/runner` (incl. `*-dind-rootless` variants). Pulls of `gitea/act_runner:nightly` etc. will return stale images.
- **Go module path**: `gitea.com/gitea/act_runner` → `gitea.com/gitea/runner`. Anything importing `act/...` or `internal/...` packages (notably Gitea itself) must update imports.
- **Release artifact paths**: S3 directory `act_runner/{{.Version}}` → `gitea-runner/{{.Version}}`, and artifact filenames change accordingly. Existing download URLs break.
- **Prometheus metrics namespace**: `act_runner_*` → `gitea_runner_*` (e.g. `act_runner_jobs_total` → `gitea_runner_jobs_total`). Existing dashboards and alerts must be updated.
- **ldflags version path**: `gitea.com/gitea/act_runner/internal/pkg/ver.version` → `gitea.com/gitea/runner/internal/pkg/ver.version`. Affects custom builds that stamp the version.
- **Kubernetes example resources**: `act-runner` / `act-runner-vol` → `runner` / `runner-vol`. Users who copied the manifests verbatim will see resource churn on apply.
- **s6 service name**: `scripts/s6/act_runner/` → `scripts/s6/gitea-runner/`. Only matters for downstream image overrides.

Unchanged: YAML config field names, env vars (`GITEA_*`), CLI flags and subcommands, registration file format.

### User-facing changes

- **Idle-time workspace janitor** (#870). Stale `bind-workdir` task workspaces are now cleaned up in the background instead of on the task execution path. Two new config keys: `runner.startup_cleanup_age` (stale-age threshold, default `24h`) and `runner.idle_cleanup_interval` (cleanup cadence, default `10m`). Cleanup is scoped to numeric task directory names only, so operator-managed folders are left alone.
- **Step failure log output now matches GitHub Actions** (#927). The error block printed when a step fails has been realigned with the GitHub Actions format.
- **Host-mode workspace cleanup fixed on Windows** (#883). Host-mode runs now remove the job workspace after the run instead of leaving checkouts behind. On Windows, remaining step process trees are terminated during teardown before workspace deletion to avoid file-lock failures. Cleanup is skipped when `bind_workdir` is enabled, to avoid conflicting with the runner-level janitor above.

### Changelog

- Align step failure log output with GitHub Actions (#927)
- feat: add startup janitor for stale bind-workdir task workspaces (#870)
- fix(host): correct host workspace cleanup on Windows (#883)
- chore(deps): update crazy-max/ghaction-import-gpg action to v7 (#923)
- fix(deps): update module golang.org/x/term to v0.42.0 (#920)
- chore(deps): update actions/hello-world-docker-action action to v2 (#921)
- ci: Fix triggers (#882)
- chore(deps): update dependency @vercel/ncc to ^0.38.0 (#881)
- chore(deps): update dependency @actions/core to v1.11.1 (#880)
- fix(deps): update module github.com/rhysd/actionlint to v1.7.12 (#873)
- Rename act_runner to runner (#850)
- docs: Update docs (#872)

## 0.6.1 (2026-04-30)

### Changelog

- perf(config): lower default fetch_interval_max from 60s to 5s (#875)
- chore: clean up `nolint` directives in `act` package (#864)
- fix(deps): update module github.com/moby/patternmatcher to v0.6.1 (#868)
- fix(deps): update module github.com/docker/docker to v25.0.15+incompatible (#867)

## 0.6.0 (2026-04-28)

This release contains a security fix for CVE-2026-28708 which introduces authentication between the runner and the cache server. Thanks to [@haroutp](https://github.com/haroutp) for the report.

### Breaking changes

- **`cache.external_server` now requires `cache.external_secret`** (#849). Operators running a standalone `act_runner cache-server` must set the same `cache.external_secret` in both the runner config and the cache-server config — neither process starts without it. Runners using the default embedded cache server are unaffected and need no config changes. As part of the same security fix, pre-existing entries in `bolt.db` have no repo stamp and will be evicted by GC, causing a one-time cache miss per key after upgrade.

### Changelog

- Clone different git repos in parallel via per-directory locks (#866)
- fix: avoid 'filename too long' on matrix jobs by hashing container names (#853)
- fix(deps): update module github.com/mattn/go-isatty to v0.0.22 (#863)
- fix(deps): update module github.com/docker/docker to v25.0.14+incompatible (#862)
- Authenticate cache requests via `ACTIONS_RUNTIME_TOKEN` and scope by repo (#849)
- feat: show run command, shell and env in collapsible group before step output (#847)
- run docker step in host mode (#857)
- chore(deps): drop unused `distribution/reference` replace directive (#858)

## 0.5.0 (2026-04-27)

### Notice

This repository now contains the merged [`gitea/act`](https://gitea.com/gitea/act) repository under the `act/` folder. Future development of `act` happens here in `act_runner`.

### Changelog

- fix(deps): update module github.com/docker/cli to v25.0.7+incompatible (#855)
- fix(deps): update module connectrpc.com/connect to v1.19.2 (#854)
- fix: wait for docker supervise dir before `s6-svwait` (#851)
- fix: Heartbeat `ReportState` for long-running silent jobs (#852)
- fix: fetch when other refs get force-pushed (#846)
- Align root files with `gitea` (#844)
- Upgrade go-git and yaml (#842)
- Merge `gitea/act` into `act/` folder (#827)

---
This release was prepared with the help of Claude Opus 4.7

## 0.4.1 (2026-04-19)

### Changelog   
* 9aafec169b7ebf3d71df2b06780868fb552b435e perf: use single poller with semaphore-based capacity control (#822)   
* 48944e136c0c3fc0dd04c0ed6cc4401bb00c7cf4 Use `golangci-lint fmt` to format code (#823)

## 0.4.0 (2026-04-15)

### Changelog   
* 40dcee0991c3bd33b657bb77aa1f2f46d69cc0e2 chore(deps): upgrade golangci-lint from v2.10.1 to v2.11.4 (#821)   
* f33e5a6245696af3dbd3a2c6004da35c720abfed feat: add Prometheus metrics endpoint for runner observability (#820)   
* f2d545565fb93df2c06cb7aeb791d9b0a24b67e4 perf: reduce runner-to-server connection load with adaptive reporting and polling (#819)   
* 90c1275f0ec64754feb73e5270547779e97b9c2d Upgrade yaml (#816)

## 0.3.1 (2026-03-26)

### Changelog

- Add `run_attempt` to context (#632)
- feat: add configurable bind_workdir option with workspace cleanup for DinD setups (#810)
- Semver tags for Docker images (#720)
- chore(lint): add golangci-lint v2 and fix all lint issues (#803)
- fix: composite action log result reported as step result (#801)
- fix: race condition in reporter between RunDaemon and Close (#796)
- chore(deps): update go dependencies (#802)
- chore(deps): update workflow action versions (#792)

## 0.3.0 (2026-02-18)

### Changelog

- chore(deps): update act to v0.261.8 (#791)
- examples: improve DIND rootless network performance (#786)
- chore(deps): bump Go version to 1.26 in Dockerfile and Makefile (#788)
- chore(deps): bump Go version to 1.26.0 (#787)
- Remove duplicate assignment of DIST (#777)
- chore(deps): upgrade dependencies (#775)
- Bump gitea/act (#770)
- fix(dockerfile): Pin docker dind images to version 28 (#769)
- Fix some typos in config.example.yaml (#764)
- Fix URL to documentation of runner images (#765)
- chore: update config.example.yaml (#763)
- ci: release binary for linux/loong64 (#756)
- Bump gitea/act (#753)
- Do not implicitly mount /var/run/docker.sock (#751)

## 0.2.13 (2025-08-29)

### Changelog   
* 8920c4a170603dc8f23e9e0d4ba0a9de13cd5f94 Timeout to wait for and optionally require docker always (#741)   
* bbf9d7e90f3fee9e9968fc49c1ee6cfdb11f55e9 feat: support github mirror (#716)   
* 46f471a9007a18d965191c6ab07711ec360107b5 refac(workflow): use matrix to compile different docker images (#740)   
* aa28f8d99c88fe7d662e49efcb8caa6215213058 feat(docker): add TZ env variable working (#738)   
* 6a7e18b124a64221d2fda1d508390aedba270b46 Allow node24 actions (#737)   
* 53329c46ff29922cea0d02a0d52a039308f4bd99 Add ubuntu-24.04 label to defaults. (#724)

## 0.2.12 (2025-06-25)

### Changelog   
* 6b1aea9c04d2d7d5fa59ab97da6712b3b2d9b502 Upgrade github.com/go-git/go-git/v5 to v5.16.2 (#706)   
* edec9a8f91c8dafeb834215002ec9f4928205d6a Add a way to specify vars in act_runner exec (#704)   
* 6a9a447f86a3951b0b409e09aabfd40aae6dbe86 Report errors by setting raw_output when it's error level (#645)   
* 5302c25feb6d64b682dd789d41984295db919050 Add environment variables for OIDC token service (#674)   
* a616ed1a10efcb91f8ff3f67b35e16686cfa28b9 feat: register interactive with values from cli (#682)   
* f0b5aff3bbeda469a682fe21144c16485ea2189a fix: invalid label NoInteractive exit code (#683)   
* 44b473670330342d66300d6ac96bccd0e39c37e7 feat: docker env vars for ephemeral and once (#685)   
* b1ae30dda8292b255ba8e6d230817e4ca949090d ephemeral act runner (#649)   
* 0d687268c77e4e391f04d1685431734993600636 act_runner requires go 1.24 now   
* 425a570261817ce86a646334367afaf71a92158e use new docker image URLs (#661)   
* 4c8179ee12aff2862b7bf43d62ed9de28b070926 upgrade to go1.24, act to 0.261.4 and actions-proto-go to 0.4.1 (#662)   
* 5ae13f0bd73f381b92496f3c60352a0ce549663a Update xgo version to 1.24 (#651)   
* 3510152e3671880d483b111a2a66c38f924ac147 Fix Makefile make docker (#641)   
* 8dfb805c627633125e0812d05998ba76d71630a2 Update examples/kubernetes/dind-docker.yaml to reflect recent changes to Dockerfile (#633)   
* a7080f5457afbfbd4931347f2f9b39c871fdb56d Update examples for `GITEA_RUNNER_REGISTRATION_TOKEN` env (#630)   
* 8b72d1c7ae3a53059476ca600951e69efaf40913 add s390x and riscv64 as an arch for binaries   
* 8bc0275e74af3516cb9d952e144e48405c084aa5 feat: add once flag to daemon command (#19) (#598)   
* 0348aaac599d7c8f9bdbc213d5d80c5786b57b5b Wait for the Docker daemon to be ready before starting act runner (#620)   
* 9712481bed83b6ce9937833efcac9cca73381872 Support `basic`, `dind` and `dind-rootless` as multiple kinds of images (#619)   
* b5f901b2d96a1d310bb4efaee28574af8c0e97f8 Upgrade act from v0.261.2 -> v0.261.3 (#607)   
* 0e2a3e00f5b32d5407f4f929a1ae61d3df825362 examples/vm/rootless-docker.md aktualisiert (#487)   
* b282356e9e5247da16ebe6bf6a6804ad22d47cf4 update example for docker-compose to allow fix 502 errors in case Gitea not yet ready on runner startup (#605)

## 0.2.11 (2024-09-24)

### Changelog   
* b075e3a1d5a0ca4f1cfc5dc3a7e6dc187aa9d23a Bump goreleaser-action to use v2 by default (#604)   
* e27189ea321666f63ea70e86622572b0ae9c2c17 Fix releasing (#603)   
* 59e478464efde43f715849f83f81a743da182292 Bump act to v0.261.2 (#602)   
* e1c7b208985b59eedbe8bc4df1174ea1c455cae1 Bump Go to 1.22 and upgrade dependencies (#580)   
* 2f78411c3d871f0046832e4d082ba67c0655baf7 ci: standardize code style and update version extraction (#566)   
* d1d3cad4b0bfd0bfc5df344306f304043ff63223 feat: allow graceful shutdowns (#546)   
* 1735b26e66f8e81c526204c128f1bcfd7f578906 Don't log job output when debug logging is not enabled (#543)   
* 65ed62d2f5fd9d104e4dc4a941a54d9eccae2d64 Upgrade dependencies (#537)

## 0.2.10 (2024-04-10)

### Notice
The default value of env.GITHUB_TOKEN follows GitHub Actions conventions. However, you can still retrieve the value from secrets.GITHUB_TOKEN if needed.

### Changelog   
* ec03f19 Bump act to v0.261.1 (#535)   
* 8567324 Release nightly versions for tags (#536)

## 0.2.9 (2024-04-02)

### Changelog   
* be2df36 Ensure declare to use new labels (#530)   
* a5085dd Use gitea/runner-images by default (#529)   
* cef86d1 Bump act to v0.260.0 (#528)   
* 94c45ac Use node 20 by default (#526)

## 0.2.8 (2024-03-27)

### Changelog   
* 1760899 Set gitea token to release (#524)   
* a7eca81 Bump act to v0.260.2 (#523)

## 0.2.7 (2024-03-27)

### Changelog   
* 23ec12b Bump act to v0.260.0 (#522)   
* a1fc2b3 Improve the explanation of `workdir_parent` config (#519)   
* 5977042 Bump act to v0.259.2 (#515)   
* 6bc19cb Improve workflows (#516)   
* 75006a5 Support cloning remote actions from insecure Gitea instances (#508)   
* 4da97f5 Bump act to v0.259.0 (#491)   
* 4527065 Set the status of steps to `skipped` if job is skipped (#500)   
* e14f42c Add ACTIONS_RESULTS_URL to env (#473)   
* e6630e2 Use artifacts v4 jwt if available (#471)   
* f1f9142 fix: Exit docker container on any process exit (#462)   
* f17cad1 Update the docker image rebuild flag for config file and example gitea/act_runner#390 (#391)   
* daf52d0 Sanitize UFT-8 content in logs (#453)   
* 8c8a8ce fix(deps): update module github.com/avast/retry-go/v4 to v4.5.1 (#411) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* 08c681b fix(deps): update module golang.org/x/time to v0.5.0 (#429) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* 91bfe4c fix(deps): update module golang.org/x/time to v0.4.0 (#424) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* 825c6f9 fix(deps): update module github.com/mattn/go-isatty to v0.0.20 (#414) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* 2f3e5c7 fix(deps): update module github.com/docker/docker to v24.0.7+incompatible (#413) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* 4d9de6c fix(deps): update module github.com/spf13/cobra to v1.8.0 (#416) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* feb3966 fix(deps): update module code.gitea.io/gitea-vet to v0.2.3 (#410) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* 61cd71c bump to actions/checkout v4   
* 0adfc1c chore(deps): update actions/setup-go action to v4 (#418) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* e3c6866 chore(deps): update docker/build-push-action action to v5 (#419) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* f1b27d5 chore(deps): update docker/login-action action to v3 (#420) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* 655a39f chore(deps): update docker/setup-buildx-action action to v3 (#421) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* cca7d54 chore(deps): update docker/setup-qemu-action action to v3 (#422) Co-authored-by: Renovate Bot <renovate-bot@gitea.com> Co-committed-by: Renovate Bot <renovate-bot@gitea.com>   
* 9344718 Fix #404: nil map error when reading env file (#405)   
* 1e940f0 Add renovate config (#408)   
* 2020ce7 Remove ACTIONS_RUNTIME_TOKEN workaround (#396)   
* 00e9b3d Bump act (#394)   
* 130b9f1 Added environment variable for run.sh to specify runner state file. (#377)   
* 4c35288 Add DIND docker-compose example (#365)   
* 990db1b Fix k8s rootless Docker filesystem permissions (#366)

## 0.2.6 (2023-09-24)

Node20 Actions support

## 0.2.5 (2023-08-02)

### Changelog   
* 9fc823e Fix "http: invalid Host header" issue (#319)

## 0.2.4 (2023-07-28)

### Changelog   
* 12999b6 Reduce unnecessary DB queries for Actions tasks (#219)   
* 49a2fcc fix endless loop (#306)   
* 8f88e4f Update README.md (#302)   
* a1bb3b5 Catch the panic and print the error (#305)   
* 1a7ec5f Upgrade gitea/act (#298)   
* 5d01cb8 Add tips in config file (#297)   
* dcf84d8 Avoid https://github.com/nektos/act/issues/1908 by specifying golang 1.20.5. (#295)   
* 73adae0 Upgrade act (#291)   
* db662b3 ci(lint): refactor code for clarity and linting compliance (#289)   
* cf92a97 refactor(register): refactor registration functions and improve logging (#288)   
* 8705871 fix(register): refactor context usage in registration functions (#286)   
* c701ba4 Add a quick start runner method in README (#282)   
* 57ff1df config: default container workspace set to host path (#279)   
* 3dcfd6e Run as cache server (#275)

## 0.2.3 (2023-07-03)

### Changelog   
* c6006ee Upgrade act (#269)

## 0.2.2 (2023-07-01)

### Changelog   
* f2629f2 Add support for finding docker daemon from common socket paths (#263)

## 0.2.1 (2023-06-30)

### Changelog   
* cf48ed8 Revert supporting multiple default actions URLs and use github for exec by default (#262)   
* ccc2732 Improve the usage of `docker_host` configuration (#260)   
* b0bd503 add token support for exec (#247)   
* 8c14933 Upgrade act (#248)   
* 34d15f2 Add option to configure workspace on host (#238)   
* 32d29f0 add ACT_EXEC (#246)   
* 2e2c040 add --gitea-instance (#245)   
* 054c8d9 Move `docker.host` to `container.docker_host` (#244)   
* 9e4a5f7 feat: improve Docker configuration and detection handling (#242)   
* ec38401 Add `ValidVolumes` config (#226)

## 0.2.0 (2023-06-15)

### Changelog   
* 45bfe0a Accept empty file as config (#241)   
* 3165349 Build docker image `gitea/act_runner/x.y.z-dind-rootless` (#239)   
* 67b1363 Support changing labels (#201)   
* 946c41c Improve `run.sh` to handle empty labels and log to stdout (#237)   
* 341d49a implement act_runner rootless image (#208)   
* b21d476 Exit with Code 1 if registering a runner fails (#228)   
* a29307a Remove hadolint and improve Dockerfile (#234)   
* 4bfbfec fix defaultActionsUrls config for exec (#233)   
* fed01c9 Parse multiple default actions URLs (#200)   
* a83f29d docs: improve examples README and organization (#230)   
* 69c55ee refactor: daemon, config, and logging for better clarity (#225)   
* 01ef57c fix `artifactServerPath` and `artifactServerAddr` config for exec (#221)   
* a384adb Documentation enhancements (#207)   
* e3271d8 Remove trailing slash from instance address (#197)   
* 84386c1 Add exec command flag of `network` (#192)   
* fd7c858 Prevent exposing GITEA_RUNNER_REGISTRATION_TOKEN to act (#188)   
* 35596a1 Add configuration item of `container.network` (#184)   
* c9d3f67 Add .editorconfig and .gitattributes (#186)   
* 94031fc Fix README.md typo on daemon (#183)

## 0.1.8 (2023-05-10)

### Changelog   
* d5e4bae Fix missing runner version when building images from a Dockerfile (#181)   
* d4caa7e Dockerfile: Improve signal handling by adding a runtime init (#180)

## 0.1.7 (2023-05-06)

### Changelog   
* de4160b Skip counting log length when parseLogRow return nil (#176)   
* 609c0a0 fix `--event` option logic for `exec` (#175)   
* 0c029f7 Upgrade act and use new artifactcache (#174)   
* eef3c32 ci: improve release process and test robustness (#173)   
* c40b651 chore: improve Dockerfile, README, and testing settings (#172)

## 0.1.6 (2023-05-02)

### Changelog   
* b498341 build: improve compression and update GitHub actions (#168)   
* 0d727eb build: optimize Dockerfile and update dependencies (#162)

## 0.1.5 (2023-04-29)

### Changelog   
* 7c71c94 Document persisting `/data` for docker container (#160)   
* 49d2cb0 ci: improve API usage and test robustness across platforms (#159)   
* 85626b6 Support `configuration variables` (#157)   
* 35400f7 Add parent directory for working directory (#154)

## 0.1.4 (2023-04-27)

### Changelog   
* 0cf31b2 Update docker image tag (#153)

## 0.1.3 (2023-04-27)

### Changelog   
* c8cc7b2 Workflow commands (#149)   
* 3be962c Rename the download folder from main -> nightly (#152)   
* a5edbc9 Release docker images (#151)   
* 66bab3d ci(actions): add build docker image workflow (#118)   
* 293926f bump modernc.org/sqlite (#141)   
* 43c5ba9 make: skip --disable-content-trust at docker buildx (#139)   
* acc5afc allow building act_runner with cgo (#137)   
* 27a1a90 Upgrade act (#135)   
* 83ec0ba Support upload `outputs` and use `needs` context (#133)   
* ed86e2f Update workflow files (#131)   
* d4bebcc Update dependencies (#130)   
* c75b67e Upgrade act (#128)   
* bc6031e Fix reporting log in `Reporter.Close` (#126)   
* c69c353 Upgrade act (#124)

## 0.1.2 (2023-04-13)

### Changelog   
* fcc016e Special the release tag signing sub key (#121)   
* d5caee3 Add prefix to use ghaction-import-gpg (#120)   
* 9e26208 Add signature for tag releases, upload to gitea releases itself (#116)   
* a05c5ba Add make docker (#115)   
* c248520 Set specific environments to distinguish between Gitea and GitHub (#113)

## 0.1.1 (2023-04-12)

_No release notes._

## 0.1.0 (2023-04-05)

_No release notes._

## 0.0.1 (2023-02-26)

this tag exists for goreleaser to be able to build nightly releases. please don't use it as a release to target.
