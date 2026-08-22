# Development

## Job log line format

Gitea stores one log row per line and its web UI decodes the payload, so getting the encoding
wrong never fails a test here, it only shows up in the browser.

**A row cannot contain a real newline.** `FormatLog` rewrites `\n` to a literal backslash-n and
truncates at 64 KiB on a byte boundary.

**The payload of a line starting with a recognised prefix is decoded**, with the escape set
depending on the prefix:

| prefix | decodes |
| --- | --- |
| `##[error]` `##[warning]` `##[notice]` `##[debug]` `##[group]` `##[endgroup]` `##[add-matcher]` | `%25` `%0D` `%0A` `%3B` `%5D` |
| `::error::` `::warning::` `::notice::` `::debug::` (with or without ` key=value` properties), `::group::` `::endgroup::` `::add-matcher::` | `%25` `%0D` `%0A` |
| `##[command]` `[command]`, or no recognised prefix | nothing |

### Rules

- **Emitting a command line?** Escape the payload with `runner.EscapeCommandData`. One escaper
  covers both forms: it escapes `%` first, so a literal `%3B` becomes `%253B` that the extra
  `##[…]` rules cannot match, and a raw `;` or `]` is never decoded. It is also what makes
  multi-line work, `\n` becomes `%0A` and the UI turns it back into a line break.
- **Forwarding a command from step output?** Leave the payload alone, it arrived escaped and is
  decoded once. Decoding here double-decodes and destroys multi-line.
- **No prefix?** Do not escape, and split multi-line values into one row each.
- **Interpolating a secret?** Masking runs after escaping, so `AppendSecretMasker` registers the
  encoded forms too.
- Command *properties* also escape `%3A` and `%2C`, which the UI never decodes, so the reporter
  decodes exactly those two when folding a location into an annotation.

## End-to-end compatibility tests

`e2e/` runs the runner against a real `gitea/gitea` server container to catch compatibility
regressions in Gitea's runner/Actions API across versions. Unlike the mocked-client unit tests,
these exercise the actual RPC and REST surface a real Gitea exposes.

Scenarios: registration, artifacts, job matrices, secret masking, log streaming, capability
negotiation, run cancellation, the Actions cache (both advertised service versions), workflow
command encoding, ephemeral runners, job outputs, repository variables, workflow dispatch, and
service containers.

Some of these differ by server version by design, so they report rather than hard-assert:

- **Run cancellation** exists only on the development branch. Stable releases register no
  `CancelWorkflowRun` route, so `GiteaAPI.CancelRun` reports `ErrCancelUnsupported` and the
  scenario instead asserts the unsupported call leaves the job and the runner unharmed.
- **Capabilities** are negotiated through the `X-Gitea-Actions-Capabilities` response header on
  Declare, and the answer varies; the scenario logs what each version negotiates. Note the two
  directions are independent: the runner advertises its own flags (`cancelling`) in Register and
  Declare, while this header is the *server* listing what it supports, so the two never match up.
- **Ephemeral** registration is accepted but silently downgraded by older servers, which
  `Harness.EphemeralDowngraded` reports.

**Not run in CI yet.** The suite is local-only and opt-in: every file under `e2e/` carries the
`e2e` build tag, so it is invisible to `go build ./...`, `go vet ./...`, `make test`, and
`.gitea/workflows/test.yml`.

Run it with:

- `make test-e2e` — one version, `GITEA_E2E_IMAGE` (default `gitea/gitea:main-nightly`) overridable.
- `make test-e2e-matrix` — loops over the versions listed in `e2e/.env`'s `GITEA_E2E_VERSIONS`,
  currently the newest stable (`latest`) and the development branch (`main-nightly`). Gitea only
  patches the newest release, so older minors are deliberately out of scope.

Most scenarios call `t.Parallel()`; the four that cannot (log streaming, cancellation, the cache
subtests, and ephemeral registration) depend on timing or on their own runner taking the work.

Both need a local Docker daemon and network access to pull `gitea/gitea` images; tests skip
cleanly if no daemon is reachable. `e2e/main_test.go`'s `TestMain` starts one Gitea container
for the whole test binary run and every scenario creates its own repository (named after the
test), so add a new scenario as a `testdata/workflows/*.yml` file plus a `*_test.go` using the
shared `fixture` (`*GiteaFixture`), a fresh `Harness` (in-process runner registration + polling,
`e2e/runner_harness.go`), and `fixture.API()` (`e2e/gitea_api.go`) to drive repos/workflows/runs.

The Gitea REST endpoints relied on — `POST /api/v1/admin/actions/runners/registration-token`, the
Contents API, `.../actions/runs`, `.../actions/runs/{run}/jobs`, `.../actions/jobs/{job_id}/logs`,
`.../actions/artifacts`, `.../actions/secrets/{name}` (PUT), `.../actions/variables/{name}`
(POST creates, PUT updates), and `.../actions/workflows/{workflow_id}/dispatches` — are stable
across the tested versions; `.../actions/runs/{run}/cancel` exists only on the development branch.
Re-verify against Gitea's source if a scenario starts failing only on one version.

Two shapes are worth knowing when adding scenarios: a run's `path` is reported as
`<workflow file>@<ref>` (not a repository path), and a job log row is `<timestamp>Z <payload>`.
