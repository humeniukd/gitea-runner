# Local job logs

A job's output normally exists in one place only: Gitea. When the path that uploads it is
what broke — `UpdateLog` failing, the instance unreachable, the runner killed mid-job — the
output is gone, and the job page shows a fraction of what ran.

Setting `log.job.dir` makes the runner keep its own copy of every task's log, on the machine
that ran it. It is the equivalent of GitHub's `_diag/Worker_*.log`. Nothing about the job
changes: the same rows are still uploaded, and Gitea shows what it always did.

## Configuration

```yaml
log:
  job:
    dir: /var/log/gitea-runner/jobs
    retention: 168h
    max_size: 0
```

| key | default | meaning |
| --- | --- | --- |
| `log.job.dir` | empty | Where to write the logs. Empty writes nothing, which is the default. |
| `log.job.retention` | `168h` | Delete a task's directory once it is older than this. Negative keeps every directory forever. |
| `log.job.max_size` | `0` | Cap for each of the two files, in bytes. `0` is unlimited. |

These are separate from `log.level`, which is the severity of the runner process's own output
on stderr and changes neither file.

Retention is applied whenever a task opens its own directory, so it also runs on a runner
that is never idle. Only directories the runner created are considered: anything else under
`log.job.dir` is left alone.

## Layout

Each task gets a directory named after its start time and id, holding two files:

```
/var/log/gitea-runner/jobs/20260814-091203-task-41/
├── job.log
└── diag.log
```

**`job.log`** is what was handed to the uploader, in order, and nothing else. Secrets are
masked exactly as they are for Gitea, and the rows the workflow-command handler drops
(`::add-mask::`, a `::debug::` without `ACTIONS_STEP_DEBUG`) are absent here too. Diffing it
against the job page in the web UI is the point: whatever is in the file and not on the page
was lost on the way.

```
# gitea-runner v0.2.13 task=41 repository=org/repo workflow=ci.yml run_id=7 started=2026-08-14T09:12:03Z
2026-08-14T09:12:03.114Z ##[group]Run actions/checkout@v4
2026-08-14T09:12:03.402Z ***
2026-08-14T09:12:04.001Z [runner] task 41 finished: success
```

**`diag.log`** holds every entry the reporter saw, which is more than it uploads: the entries
below the configured log level, the runner's own debug output during a step, and the commands
dropped from the job log. It also records each failed call to Gitea, prefixed `[report]`:

```
2026-08-14T09:12:03.115Z DEBUG stage=Pre stepID=[checkout] | evaluating 'always()'
2026-08-14T09:12:04.900Z ERROR | [report] UpdateLog(index=41 rows=12 no_more=false) failed: dial tcp: connection refused
```

Both files end with a `[runner] task <id> finished: <result>` line, so a file without one
belongs to a runner that died before it could finish the job.

## Operational notes

- **The files hold repository output.** The directory is created `0700` and the files `0600`,
  readable by the user the runner runs as. Masking is the same as for the uploaded rows, which
  means a secret the job never told the runner about is not masked in either place.
- **Writing never blocks or fails a job.** Lines are queued and written by one goroutine; if
  it cannot keep up, lines are dropped rather than delayed, and both files say how many. A
  directory that cannot be created disables the feature for that task with a warning in the
  runner log.
- **Only the daemon writes these files.** `gitea-runner exec` already prints to the terminal.
- Rotation is per task, not per file: there is nothing to signal or reopen, and log rotation
  tools should be pointed at `log.job.retention` instead.
