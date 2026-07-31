// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	utilexec "k8s.io/client-go/util/exec"
)

func TestBuildExecCommandPlain(t *testing.T) {
	got := buildExecCommand([]string{"echo", "hello"}, nil, "", "")
	require.Len(t, got, 3)
	assert.Equal(t, "sh", got[0])
	assert.Equal(t, "-c", got[1])
	assert.Equal(t, "exec echo hello", got[2])
}

func TestBuildExecCommandEnvSortedAndWorkdir(t *testing.T) {
	got := buildExecCommand(
		[]string{"node", "index.js"},
		map[string]string{"BBB": "2", "AAA": "1"},
		"",
		"/workspace/owner/repo",
	)
	require.Len(t, got, 6)
	assert.Equal(t, "env", got[0])
	// env vars applied in sorted key order via env(1), then sh -c with cd+exec.
	assert.Equal(t, "AAA=1", got[1])
	assert.Equal(t, "BBB=2", got[2])
	assert.Equal(t, "sh", got[3])
	assert.Equal(t, "-c", got[4])
	expected := "cd /workspace/owner/repo || exit 1\nexec node index.js"
	assert.Equal(t, expected, got[5])
}

func TestBuildExecCommandPreservesDashedEnvKeys(t *testing.T) {
	got := buildExecCommand(
		[]string{"node", "dist/index.js"},
		map[string]string{"INPUT_FETCH-DEPTH": "0", "INPUT_PERSIST-CREDENTIALS": "true"},
		"", "",
	)
	require.GreaterOrEqual(t, len(got), 3)
	assert.Equal(t, "env", got[0])
	assert.Contains(t, got, "INPUT_FETCH-DEPTH=0")
	assert.Contains(t, got, "INPUT_PERSIST-CREDENTIALS=true")
}

func TestBuildExecCommandQuotesUnsafeValues(t *testing.T) {
	got := buildExecCommand([]string{"sh", "-c", "echo $FOO"}, map[string]string{"FOO": "a b; rm -rf /"}, "", "")
	// env(1) receives the value as a single argv element; the inner sh -c
	// script and command arg must still be shell-quoted.
	assert.Contains(t, got, "FOO=a b; rm -rf /")
	assert.Contains(t, got, "sh")
	assert.Contains(t, got, "exec sh -c 'echo $FOO'")
}

func TestBuildExecCommandPreservesUnusualEnvKeys(t *testing.T) {
	got := buildExecCommand(
		[]string{"echo", "hi"},
		map[string]string{
			"GOOD":        "1",
			"x; id #":     "evil", // shell metacharacters in the key
			"$(reboot)":   "evil", // command substitution in the key
			"1BAD":        "evil", // starts with a digit
			"with-dash":   "evil", // hyphen is allowed in env(1) names
			"INPUT_DEPTH": "0",
		},
		"", "",
	)
	// env(1) receives each NAME=VALUE as literal argv — no shell parses the
	// key — so metacharacter keys are preserved rather than dropped.
	assert.Contains(t, got, "GOOD=1")
	assert.Contains(t, got, "with-dash=evil")
	assert.Contains(t, got, "INPUT_DEPTH=0")
	assert.Contains(t, got, "1BAD=evil")
	assert.Contains(t, got, "x; id #=evil")
	assert.Contains(t, got, "$(reboot)=evil")
}

func TestBuildExecCommandSkipsUnrepresentableEnvKeys(t *testing.T) {
	got := buildExecCommand(
		[]string{"echo", "hi"},
		map[string]string{
			"GOOD":      "1",
			"A=B":       "evil", // '=' breaks env(1) NAME=VALUE parsing
			"-FLAG":     "evil", // leading '-' would be parsed as an env(1) option
			"x\x00y":    "evil", // NUL cannot round-trip through argv
			"with-dash": "evil", // hyphen is allowed in env(1) names
		},
		"", "",
	)
	assert.Contains(t, got, "GOOD=1")
	assert.Contains(t, got, "with-dash=evil")
	assert.NotContains(t, got, "A=B=evil")
	assert.NotContains(t, got, "-FLAG=evil")
	assert.NotContains(t, got, "x\x00y=evil")
}

func TestBuildExecCommandUser(t *testing.T) {
	got := buildExecCommand([]string{"whoami"}, nil, "runner", "")
	require.Len(t, got, 4)
	assert.Equal(t, "su", got[0])
	assert.Equal(t, "runner", got[1])
	assert.Equal(t, "-c", got[2])
	assert.Equal(t, "exec whoami", got[3])
}

func TestMapExecError(t *testing.T) {
	require.NoError(t, mapExecError(nil))

	plain := errors.New("connection reset")
	assert.Equal(t, plain, mapExecError(plain))

	mapped := mapExecError(utilexec.CodeExitError{Err: errors.New("command failed"), Code: 7})
	var exit ExitCodeError
	require.ErrorAs(t, mapped, &exit)
	assert.Equal(t, 7, int(exit))
	assert.Equal(t, "Process completed with exit code 7.", mapped.Error())
}
