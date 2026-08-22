// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

// TestLogCommandEncoding round-trips workflow command lines through Gitea's log storage. The
// encoding rules in DEVELOPMENT.md are invisible to the rest of the suite — a wrong escape
// "never fails a test here, it only shows up in the browser" — so this is the one scenario that
// reads stored rows back and checks the payloads survived intact.
func TestLogCommandEncoding(t *testing.T) {
	requireFixture(t)
	t.Parallel() // own repo and runner, so it shares nothing with the other scenarios

	api, repo, _ := startScenario(t, "logformat.yml")

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)

	// Payloads must arrive whole, whether they came from a plain line or a command line.
	for _, want := range []string{
		"plain-100%-done;-[bracket]", // unescaped characters the ##[…] rules would otherwise decode
		"multiline-first",
		"multiline-second",
		"notice-payload-here",
		"warning-payload-here",
		"error-payload-here",
		"group-payload-here",
		"inside-the-group",
	} {
		if !strings.Contains(logs, want) {
			t.Errorf("stored logs are missing %q", want)
		}
	}

	// A command forwarded from step output keeps its payload as it arrived: re-escaping here
	// would store %250A, which the UI decodes to a literal %0A instead of a line break. The
	// comparison is per row because the group header for the step quotes the command verbatim,
	// so a substring search cannot tell the emitted header from the forwarded row.
	forwarded := commandRow(logs, "::notice::encoded-first")
	switch forwarded {
	case "":
		t.Errorf("no forwarded ::notice:: row carried the encoded newline")
	case "::notice::encoded-first%0Aencoded-second":
	default:
		t.Errorf("forwarded command row is %q, want its payload passed through unchanged", forwarded)
	}

	// The other half of the rule: data the runner emits itself is escaped, so a % in a step
	// name is stored as %25 and the UI decodes it back rather than eating the following chars.
	if !strings.Contains(logs, "plain-100%25-done") {
		t.Errorf("the emitted group header did not escape %% in the step name")
	}

	if t.Failed() {
		t.Logf("full stored logs:\n%s", logs)
	}
}

// commandRow returns the payload of the stored row that starts with prefix, stripped of the
// leading timestamp Gitea prepends. It returns "" when no row starts with prefix.
func commandRow(logs, prefix string) string {
	for line := range strings.SplitSeq(logs, "\n") {
		_, payload, found := strings.Cut(line, "Z ") // rows are "<RFC3339 nanos>Z <payload>"
		if !found {
			continue
		}
		if payload = strings.TrimRight(payload, "\r"); strings.HasPrefix(payload, prefix) {
			return payload
		}
	}
	return ""
}
