# Development

Notes for working on the runner itself. See `README.md` for installation and configuration, and
`AGENTS.md` for build and lint commands.

## Job log line format

Job log lines are stored by Gitea one row per line and decoded by its web UI. Getting the
encoding wrong never fails a test here, it only shows up in the browser.

### Two constraints

**A row cannot contain a real newline.** Gitea's `FormatLog` rewrites `\n` to a literal
backslash-n and truncates at 64 KiB on a byte boundary.

**The web UI decodes the payload of a line starting with a recognised prefix**, with the escape
set depending on the prefix:

| prefix | decodes |
| --- | --- |
| `##[error]` `##[warning]` `##[notice]` `##[debug]` `##[group]` `##[endgroup]` `##[add-matcher]` | `%25` `%0D` `%0A` `%3B` `%5D` |
| `::error::` `::warning::` `::notice::` `::debug::` (with or without ` key=value` properties), `::group::` `::endgroup::` `::add-matcher::` | `%25` `%0D` `%0A` |
| `##[command]` `[command]`, or no recognised prefix | nothing |

`##[…]` also decodes `%3B` and `%5D` because `;` and `]` delimit its header, so moving a message
between the two forms silently changes how it is decoded.

### Rules

- **Emitting a command line?** Escape the payload with `runner.EscapeCommandData`. One escaper
  covers both prefix forms: it escapes `%` first, which turns a literal `%3B` into `%253B` that
  the extra `##[…]` rules cannot match, and a raw `;` or `]` is never decoded. That is also what
  makes multi-line work: `\n` becomes `%0A` and the UI turns it back into a line break.
- **Passing a workflow command through from container output?** Leave the payload alone, it
  arrived escaped and the UI decodes it once. Decoding here double-decodes, and destroys
  multi-line by putting a real newline into the row.
- **No prefix?** Do not escape, nothing will decode it.
- **Interpolating a secret?** Masking runs after escaping, so `AppendSecretMasker` registers the
  encoded forms as well.
- Command *properties* use a superset that also escapes `%3A` and `%2C`, which the UI never
  decodes, so the reporter decodes exactly those two when folding a location into an annotation.

Not worth fighting: a lone `\r` renders as a line break, one trailing newline is stripped, and
leading whitespace is trimmed from labelled messages.

### Why annotations are formatted here

Gitea has no annotation store and its UI strips command properties, so `::error file=…::` would
lose its location. The reporter folds it into the message, keeping the `::error::` prefix so the
payload still gets the escape set it was written in. It happens in the runner because Gitea
serves raw log downloads straight from storage, with the frontend not in the path.

`TestCommandLineEscaping` pins both escape tables side by side. To check the decoder itself, run
a scratch test in a Gitea checkout against `parseLogLineCommand` and `createLogLineMessage`.
