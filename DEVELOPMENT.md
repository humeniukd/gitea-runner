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
  covers both forms: it escapes `%` first, which turns a literal `%3B` into `%253B` that the
  extra `##[…]` rules cannot match, and a raw `;` or `]` is never decoded. That is also what
  makes multi-line work, `\n` becomes `%0A` and the UI turns it back into a line break.
- **Forwarding a command from step output?** Leave the payload alone, it arrived escaped and is
  decoded once. Decoding here double-decodes, and destroys multi-line by putting a real newline
  into the row.
- **No prefix?** Do not escape, and split multi-line values into one row each.
- **Interpolating a secret?** Masking runs after escaping, so `AppendSecretMasker` registers the
  encoded forms too.
- Command *properties* use a superset that also escapes `%3A` and `%2C`, which the UI never
  decodes, so the reporter decodes exactly those two when folding a location into an annotation.

Not worth fighting: a lone `\r` renders as a line break, one trailing newline is stripped, and
leading whitespace is trimmed from labelled messages.

### Why annotations are formatted here

Gitea has no annotation store and its UI strips command properties, so `::error file=…::` would
lose its location. The reporter folds it into the message instead, in the runner rather than the
UI because raw log downloads are served straight from storage.

`TestEscapeCommandData` pins the escaping. To check the decoder, run a scratch test in a Gitea
checkout against `parseLogLineCommand` and `createLogLineMessage`.
