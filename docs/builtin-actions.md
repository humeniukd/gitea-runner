# Built-in actions

Native Go actions selected with `uses: builtin:<name>`. They run in the runner
process. `actions/checkout@v4` and other remote actions are unchanged.

```yaml
- uses: builtin:checkout
  with:
    fetch-depth: 0
```

Unknown `builtin:` names fail at step construction. Names match
`^[a-z0-9][a-z0-9-]*$`.

## checkout

Subset of `actions/checkout` v4. Other inputs are ignored.

| Input | Default | Notes |
| --- | --- | --- |
| `repository` | `github.repository` | `org/repo` |
| `ref` | `github.ref` | |
| `token` | `github.token` | clone auth |
| `path` | workspace root | must stay under `GITHUB_WORKSPACE` |
| `fetch-depth` | `1` | honored on clone; `0` is full history |

Same-repo checkout copies the host workspace in (or no-ops if it is already
bind-mounted) unless skip is disabled. Otherwise the runner clones on the host
and copies the tree into the container.

## Adding one

1. Implement `builtin.Action` under `act/runner/builtin/<name>` and `Register` from `init()`.
2. Blank-import the package from `act/runner/builtins_register.go`.
