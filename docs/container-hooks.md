# Container Hooks: runner ↔ hooks contract

Status: prototype, works end-to-end in Docker mode without ever touching
the Docker API.

## 1. Overview

Container Hooks is a mode in which the runner delegates all container
lifecycle operations to an external entrypoint. When the mode is enabled,
the runner does not call the Docker API at all — not for create, start,
exec, copy, remove, or network setup. Instead, the runner invokes a single
external command once per hook event and exchanges structured data with it
through files and environment variables.

The external command is referred to below as "the hook". The hook decides
everything about the runtime: Kubernetes pods, Nomad jobs, plain Docker
commands run out-of-process, or no runtime at all in a test. The runner
does not care which backend the hook talks to, and the contract does not
assume anything about where the hook or the job physically run.

Communication channels:

- Environment variables, runner → hook.
- JSON files in a per-job state directory, both directions.

The runner never writes into the job's filesystem and never reads from it.
All file exchange goes through the state directory, whose location the
runner controls.

## 2. Activation

Mode is enabled via the runner configuration:

```yaml
container:
  hooks:
    enabled: true
    # Optional interpreter used to launch cmd. Empty means cmd is executed
    # directly. Useful on Windows, where shebangs in cmd are not honored.
    entrypoint: /bin/sh
    # Required when enabled=true. Path to the hook entrypoint. The runner
    # invokes it for every hook event.
    cmd: /opt/container-hooks/main.sh
```

When `enabled` is false, the runner ignores `entrypoint` and `cmd` and uses
the built-in Docker backend. When `enabled` is true but `cmd` is empty, the
runner fails the job with a clear error before any step runs.

## 3. Lifecycle

The order of invocations is fixed:

```text
prepare  ->  execute  ->  execute  ->  ...  ->  execute  ->  cleanup
             (step 1)     (step 2)              (step N)
```

- `prepare` is called exactly once, before any step runs.
- `execute` is called once per step: for `run:` steps, for JS/Go actions,
  and for `uses: docker://...` steps.
- `cleanup` is called exactly once, after every step has finished
  (including post steps of actions).

The hook learns which event it is being invoked for from the
`RUNNER_CONTAINER_HOOK_TYPE` environment variable. A single `cmd` is
expected to dispatch on this value and route to its own `prepare`,
`execute`, and `cleanup` implementations.

## 4. State directory

The runner creates a per-job state directory:

```text
/tmp/gitea-runner-container-hooks-<random>/
```

Its absolute path is passed to the hook via
`RUNNER_CONTAINER_HOOK_STATE_DIR`. The runner and the hook exchange files
in this directory:

```text
<event>-input.json    # runner -> hook   (written by the runner)
<event>.json          # hook -> runner   (written by the hook)
```

For events:

```text
prepare-input.json   / prepare.json
execute-input.json   / execute.json
```

The directory is created before the first hook call and removed after
`cleanup` returns.

`cleanup` has no input or response file. The environment variables
`RUNNER_CONTAINER_HOOK_INPUT_FILE` and `RUNNER_CONTAINER_HOOK_RESPONSE_FILE`
are still set for it, pointing to paths that do not exist.

## 5. Environment

### 5.1. Runner-owned variables

Set for every hook invocation.

| Variable | Description |
| -------- | ----------- |
| `RUNNER_CONTAINER_HOOK_TYPE` | Event: `prepare`, `execute`, or `cleanup`. |
| `RUNNER_CONTAINER_HOOK_STATE_DIR` | Absolute path to the per-job state directory. |
| `RUNNER_CONTAINER_HOOK_INPUT_FILE` | Path to `<state_dir>/<event>-input.json`. |
| `RUNNER_CONTAINER_HOOK_RESPONSE_FILE` | Path to `<state_dir>/<event>.json`. |
| `RUNNER_CONTAINER_HOOK_ID` | The `id` returned by prepare. Empty in prepare, populated in execute and cleanup. |
| `RUNNER_CONTAINER_HOOK_NAME` | Job container name, as it appears in the job log. |
| `RUNNER_CONTAINER_HOOK_JOB_NAME` | The job's name. |
| `RUNNER_CONTAINER_HOOK_WORKFLOW_NAME` | The workflow's name. |
| `RUNNER_CONTAINER_HOOK_LOG_LEVEL` | Log level the runner uses (`info`, `debug`, ...). |

### 5.2. Passthrough context

All variables from the job environment are copied verbatim into the hook
environment, except the ones the runner already owns. This includes:

- `GITHUB_*` — the full set of GitHub Actions variables.
- `RUNNER_*` — OS, ARCH, TEMP, TOOL_CACHE, WORKSPACE, NAME, ENVIRONMENT.
- `CI`, `ACT`.
- Any workflow-level `env:`.

No filtering is applied. If a value ends up in the job environment, it
ends up in the hook environment as well.

### 5.3. Variables exported by the prepare hook

The prepare hook may return an `environments` map in `prepare.json`. Its
entries are merged into the environment of every subsequent hook call
(`execute` and `cleanup`), under their given names. They are **not**
injected into the job's step environment.

This is the intended channel for hook-to-hook communication: for example,
prepare can stash a pod name or an API address that execute and cleanup
will need.

### 5.4. Credentials

Registry credentials are never written to JSON. Instead, JSON contains
references of the form `{"$env": "VAR_NAME"}` and the real values travel
in the hook's process environment. The naming convention is:

```text
RUNNER_CONTAINER_HOOK_SECRET_<PATH_IN_UPPER_SNAKE>
```

Examples:

- Job container credentials:
  - `RUNNER_CONTAINER_HOOK_SECRET_JOB_CONTAINER_CREDENTIALS_USERNAME`
  - `RUNNER_CONTAINER_HOOK_SECRET_JOB_CONTAINER_CREDENTIALS_PASSWORD`
- Service container credentials for a service named `postgres`:
  - `RUNNER_CONTAINER_HOOK_SECRET_JOB_SERVICES_POSTGRES_CREDENTIALS_USERNAME`
  - `RUNNER_CONTAINER_HOOK_SECRET_JOB_SERVICES_POSTGRES_CREDENTIALS_PASSWORD`

`uses: docker://...` steps inherit credentials from `job.container`;
there is no separate source of credentials for container steps.

## 6. JSON files

### 6.1. `prepare-input.json` (runner → prepare)

Describes the intended runtime: the job container, its service containers,
every `uses: docker://...` step, and the actions the runner has already
downloaded. The hook is free to ignore any field it does not need.

```json
{
  "job": {
    "container": {
      "name": "GITEA-ACTIONS-TASK-...-JOB-test-...",
      "image": "docker.io/owner/image:tag",
      "credentials": {
        "username": { "$env": "RUNNER_CONTAINER_HOOK_SECRET_JOB_CONTAINER_CREDENTIALS_USERNAME" },
        "password": { "$env": "RUNNER_CONTAINER_HOOK_SECRET_JOB_CONTAINER_CREDENTIALS_PASSWORD" }
      },
      "entrypoint": ["/bin/sleep", "10800"],
      "cmd": null,
      "workingDir": "/workspace/owner/repo",
      "env": { "LANG": "C.UTF-8", "RUNNER_TEMP": "/tmp" },
      "privileged": false,
      "usernsMode": "",
      "platform": "linux/amd64",
      "options": "",
      "networkAliases": ["test"]
    },
    "services": [
      {
        "name": "postgres",
        "image": "postgres:16",
        "cmd": [],
        "workingDir": "/workspace/owner/repo",
        "env": { "POSTGRES_DB": "testdb" },
        "ports": ["5432:5432"],
        "options": "--health-cmd pg_isready",
        "networkAliases": ["postgres"]
      }
    ],
    "steps": [
      {
        "number": 1,
        "id": "1",
        "image": "alpine:3.20",
        "entrypoint": ["/bin/sh"],
        "args": ["-c", "echo hi"],
        "workingDir": "/workspace/owner/repo",
        "env": {}
      }
    ],
    "mounts": {
      "volumes": {},
      "binds": []
    },
    "network": {
      "name": "container-hook",
      "createAndDelete": true
    }
  },
  "runner": {
    "workdir": "/workspace/owner/repo",
    "toolCache": "/opt/hostedtoolcache",
    "bindWorkdir": false,
    "sharedToolCache": false,
    "eventJSON": "{\"ref\":\"refs/heads/main\",...}",
    "actionsDir": "/root/.cache/act",
    "actions": [
      {
        "hash": "c3fe249fe73091a17d6638fe1341e7bd0bcc3466ce52323c0688e83e2463a4ab",
        "localPath": "/root/.cache/act/c3fe249fe73091a17d6638fe1341e7bd0bcc3466ce52323c0688e83e2463a4ab"
      }
    ]
  }
}
```

Notes:

- `job.container` describes the main job container: the one that normally
  runs the steps and holds the workspace.
- `job.services` describes service containers declared in the workflow.
  Each entry mirrors the corresponding Docker container configuration.
- `job.steps` describes every `uses: docker://...` step in the workflow.
  The runner collects these before prepare runs, so the hook can create
  all required containers up front, in one pod if it wants.
- `job.mounts` and `job.network` carry the Docker-model view for hooks
  that choose to reproduce it.
- `runner.workdir`, `runner.toolCache`, `runner.bindWorkdir`,
  `runner.sharedToolCache` describe the runtime the runner would have
  built itself.
- `runner.eventJSON` is the raw event payload; the hook writes it to its
  own `<actPath>/workflow/event.json`.
- `runner.actionsDir` is the root of the runner's action cache.
- `runner.actions[]` lists every action the runner has downloaded for this
  job, as `{hash, localPath}`. The hook must stage each entry into
  `<hookActPath>/actions/<hash>/` before any step runs, because the job's
  step environment and the tooling inside action code expect the action
  at that path.

### 6.2. `prepare.json` (prepare → runner)

```json
{
  "id": "test-pod-123",
  "actPath": "/tmp/gitea-runner-container-hooks-XXX/act",
  "environments": {
    "CUSTOM_VAR": "value"
  }
}
```

- `id` (required, non-empty). Passed back to the hook in `execute` and
  `cleanup` as `RUNNER_CONTAINER_HOOK_ID`. What it means is up to the
  hook: pod name, container name, composite ID, anything.
- `actPath` (optional). Base path for workflow-scoped files — `event.json`,
  file-command files, generated scripts, staged actions. The runner uses
  it to fill the step environment (`GITHUB_EVENT_PATH`, `GITHUB_OUTPUT`,
  ...) and the file-command list. When omitted, the runner falls back to
  its own default (`/var/run/act` on Linux). The hook is responsible for
  ensuring the path exists and is visible to the job.
- `environments` (optional). Merged into the environment of every
  subsequent hook call. See section 5.3.

If the file is missing, empty, or `id` is empty, the job fails before the
first step.

### 6.3. `execute-input.json` (runner → execute)

The `kind` field selects what to run.

`kind = "script"` (a `run:` step):

```json
{
  "kind": "script",
  "script": "echo hello\nls -la",
  "workingDir": "/workspace/owner/repo",
  "env": { "VAR": "value" },
  "fileCommands": [
    "<actPath>/workflow/outputcmd.txt",
    "<actPath>/workflow/statecmd.txt",
    "<actPath>/workflow/envs.txt",
    "<actPath>/workflow/pathcmd.txt",
    "<actPath>/workflow/step-summary-0.md"
  ]
}
```

`kind = "command"` (a JS or Go action):

```json
{
  "kind": "command",
  "command": ["node", "<actPath>/actions/<hash>/dist/index.js"],
  "workingDir": "/workspace/owner/repo",
  "env": { "INPUT_TOKEN": "...", "INPUT_REPOSITORY": "..." },
  "fileCommands": [ ... ]
}
```

The runner rewrites `/var/run/act/actions` in the command to
`<actPath>/actions`, so the hook always sees paths under its own actPath.

`kind = "container-step"` (a `uses: docker://...` step):

```json
{
  "kind": "container-step",
  "stepNumber": 1,
  "image": "alpine:3.20",
  "entrypoint": ["/bin/sh"],
  "args": ["-c", "echo 'hello'; ls"],
  "workingDir": "/workspace/owner/repo",
  "env": {},
  "fileCommands": [ ... ]
}
```

### 6.4. `execute.json` (execute → runner)

```json
{
  "exitCode": 0,
  "fileCommands": {
    "<actPath>/workflow/outputcmd.txt": "commit=abc123\nref=refs/heads/main\n",
    "<actPath>/workflow/statecmd.txt": "isPost=true\n",
    "<actPath>/workflow/envs.txt": "",
    "<actPath>/workflow/pathcmd.txt": "",
    "<actPath>/workflow/step-summary-0.md": "## Summary\n"
  }
}
```

- `exitCode`: the exit code of the command that was run. A non-zero value
  becomes `container.ExitCodeError` and marks the step as failed.
- `fileCommands`: contents of the file-command files after the command
  ran, keyed by the same paths that came in `execute-input.json`. The
  runner applies them to its own state via `setEnv` / `saveState` /
  `setOutput` / `addPath`.

The hook must always write a valid `execute.json`, even when the command
fails. Missing `fileCommands` is treated as "no file commands were
produced".

## 7. Hook responsibilities

### 7.1. prepare

**Receives:**

- env: `RUNNER_CONTAINER_HOOK_TYPE=prepare`, `RUNNER_CONTAINER_HOOK_STATE_DIR`,
  `RUNNER_CONTAINER_HOOK_INPUT_FILE`, `RUNNER_CONTAINER_HOOK_RESPONSE_FILE`,
  `RUNNER_CONTAINER_HOOK_LOG_LEVEL`, passthrough context.
- File `prepare-input.json`.

**Does:**

- Ensures the runtime is ready: creates a pod, prepares a filesystem,
  mounts volumes, writes `event.json` and `envs.txt` — whatever the
  chosen backend requires.
- Stages every action listed in `runner.actions[]` into
  `<actPath>/actions/<hash>/`.
- Writes `prepare.json` with a non-empty `id`, the `actPath` it chose,
  and (optionally) `environments`.

**Must:**

- Produce a valid `prepare.json` with a non-empty `id`.
- Ensure the returned `actPath` exists and is visible to the job (the
  same path must be seen by the job's actions and scripts).
- Exit with code 0 on success.

**May:**

- Ignore any field of `prepare-input.json` it does not need.

### 7.2. execute

**Receives:**

- env: `RUNNER_CONTAINER_HOOK_TYPE=execute`, `RUNNER_CONTAINER_HOOK_STATE_DIR`,
  `RUNNER_CONTAINER_HOOK_INPUT_FILE`, `RUNNER_CONTAINER_HOOK_RESPONSE_FILE`,
  `RUNNER_CONTAINER_HOOK_ID`, passthrough context, variables from prepare.
- File `execute-input.json` with one of the three kinds.

**Does:**

- Creates the file-command files listed in `fileCommands` before running
  the command.
- Exports the `env` map from the input to the command process. Must use
  the `env` command (or `execve` semantics), not shell `export`, because
  variable names can contain hyphens (e.g. `INPUT_ALLOW-UNSAFE-PR-CHECKOUT`).
- Changes the working directory to `workingDir` (if it exists).
- Runs the command: script via stdin, command and container-step as argv.
- Reads the file-command files after execution and writes their contents
  into `execute.json`.

**Must:**

- Always write a valid `execute.json`.
- Record `exitCode`.
- Record the contents of the file-command files in `fileCommands`.
- Always exit with code 0. The step's exit code is transported in the
  JSON, not in the hook's own exit status.

**Contract with actions:**

Actions toolkit writes to `GITHUB_OUTPUT`, `GITHUB_STATE`, `GITHUB_ENV`,
and `GITHUB_PATH`. These are the files whose paths arrive in the step
environment as `<actPath>/workflow/*.txt`. The hook creates them before
the command, reads them after, and hands their contents to the runner
through `execute.json.fileCommands`. The runner applies them using its
existing mechanisms.

### 7.3. cleanup

**Receives:**

- env: `RUNNER_CONTAINER_HOOK_TYPE=cleanup`, `RUNNER_CONTAINER_HOOK_STATE_DIR`,
  `RUNNER_CONTAINER_HOOK_ID`, passthrough context, variables from prepare.

**Does:**

- Removes everything prepare created: pod, temporary files, secrets,
  staged actions, actPath, etc.

**Must:**

- Exit with code 0.
- Be idempotent: repeated invocation must not fail.

**Notes:**

- Called after the job's execution phase, with its own timeout.
- Does not affect the job result: a failing cleanup is logged but does
  not turn a successful job into a failed one.

## 8. Additional contracts

### 8.1. File-command paths

The runner and hook agree that file-command files live under
`<actPath>/workflow/`:

```text
<actPath>/workflow/outputcmd.txt
<actPath>/workflow/statecmd.txt
<actPath>/workflow/pathcmd.txt
<actPath>/workflow/envs.txt
<actPath>/workflow/step-summary-<N>.md
```

`<actPath>` is whatever the prepare hook returned in `prepare.json`, or
the runner's default when the hook omitted it. These paths are produced
by the runner and passed to the hook in `execute-input.json.fileCommands`.

The file format matches GitHub Actions file commands and is parsed by
`container.ParseEnvFileContent`, which is shared between Docker and hook
modes.

### 8.2. Action staging

In the Docker backend, action files are copied into the job container at
`/var/run/act/actions/<hash>/`. In hook mode there is no container the
runner owns, so:

- The runner lists every downloaded action under `runner.actions[]`.
- The hook stages each entry into `<actPath>/actions/<hash>/`.
- The runner rewrites the `/var/run/act/actions` prefix in every action
  command to `<actPath>/actions`, so the command the hook receives
  already points where the hook staged the files.

The runner does not stage actions itself and does not fall back to its
own cache.

### 8.3. Working directory

`workingDir` in `execute-input.json`:

- for `script` and `container-step`: `rc.Config.Workdir`, or the step's
  explicit `working-directory` when set.
- for `command`: the workdir provided by the caller.

The hook must `cd` into this directory before running the command, if it
exists.

## 9. Current limitations

- Job hooks (`JobStartedHook` / `JobCompletedHook` from the runner
  configuration) are skipped in hook mode. They are a separate mechanism
  from container hooks and are not yet integrated.
- Cleanup has no input file. Hook-to-hook state must travel through
  `environments` in `prepare.json`.
- `Config.ActionCache` (the git-object-store backend used for offline
  fixtures) is incompatible with hook mode: the hook cannot reach the
  object store. The runner refuses to start an action step in this
  combination with a clear error.
- No timeout is enforced on individual hook invocations beyond the job's
  own timeout.
- No way for the hook to signal "retry this step" back to the runner.