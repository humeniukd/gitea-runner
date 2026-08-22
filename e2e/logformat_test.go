// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"strings"
	"testing"
)

// Round-trips workflow command encoding through Gitea log storage (see DEVELOPMENT.md).
func TestLogCommandEncoding(t *testing.T) {
	requireFixture(t)
	t.Parallel()

	api, repo, _ := startScenario(t, "logformat.yml")

	wfRun := waitForRun(t, api, repo)
	requireSuccess(t, api, repo, wfRun.ID)

	logs := runLogs(t, api, repo, wfRun.ID)

	for _, want := range []string{
		"plain-100%-done;-[bracket]",
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

	// Forwarded commands must not be re-escaped (%0A stays %0A, not %250A).
	forwarded := commandRow(logs, "::notice::encoded-first")
	switch forwarded {
	case "":
		t.Errorf("no forwarded ::notice:: row carried the encoded newline")
	case "::notice::encoded-first%0Aencoded-second":
	default:
		t.Errorf("forwarded command row is %q, want its payload passed through unchanged", forwarded)
	}

	if !strings.Contains(logs, "plain-100%25-done") {
		t.Errorf("the emitted group header did not escape %% in the step name")
	}

	if t.Failed() {
		t.Logf("full stored logs:\n%s", logs)
	}
}

func commandRow(logs, prefix string) string {
	for line := range strings.SplitSeq(logs, "\n") {
		_, payload, found := strings.Cut(line, "Z ") // "<timestamp>Z <payload>"
		if !found {
			continue
		}
		if payload = strings.TrimRight(payload, "\r"); strings.HasPrefix(payload, prefix) {
			return payload
		}
	}
	return ""
}
