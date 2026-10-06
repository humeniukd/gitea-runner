# Built-in actions

Built-in actions ship with the runner, so a job needs no action download and no Node in its image.

```yaml
- uses: builtin:checkout
```

## checkout

Works like `actions/checkout` with these inputs. Any other input fails the step. It runs `git` where the job runs, so the job image needs `git` 2.34.1 or newer.

| Input | Default |
| --- | --- |
| `repository` | `github.repository` |
| `ref` | the triggering ref and commit, or the default branch of another repository |
| `token` | `github.token` |
| `path` | workspace root |
| `fetch-depth` | `1`, `0` for the full history of every branch and tag |
| `fetch-tags` | `false`, `true` fetches every tag |
| `clean` | `true`, removes untracked files |
| `persist-credentials` | `true`, later steps of the job can fetch and push with `token` |
| `show-progress` | ignored |

A branch is checked out as a local branch tracking `origin`, a tag or commit as a detached `HEAD`. The step outputs `ref` and `commit` like `actions/checkout`.

With `gitea-runner exec`, checking out the workflow's own repository copies your local working directory, unless `--no-skip-checkout` is set.
