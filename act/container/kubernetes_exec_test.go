// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"errors"
	"strings"
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
	require.Len(t, got, 3)
	assert.Equal(t, "sh", got[0])
	// env exported in sorted key order, then cd, then exec.
	expected := "export AAA=1\nexport BBB=2\ncd /workspace/owner/repo || exit 1\nexec node index.js"
	assert.Equal(t, expected, got[2])
}

func TestBuildExecCommandQuotesUnsafeValues(t *testing.T) {
	got := buildExecCommand([]string{"sh", "-c", "echo $FOO"}, map[string]string{"FOO": "a b; rm -rf /"}, "", "")
	// The value and the command arg must be shell-quoted so metacharacters are literal.
	assert.Contains(t, got[2], "export FOO='a b; rm -rf /'")
	assert.Contains(t, got[2], "exec sh -c 'echo $FOO'")
}

func TestBuildExecCommandSkipsMaliciousEnvKeys(t *testing.T) {
	got := buildExecCommand(
		[]string{"echo", "hi"},
		map[string]string{
			"GOOD":      "1",
			"x; id #":   "evil", // shell metacharacters in the key
			"$(reboot)": "evil", // command substitution in the key
			"1BAD":      "evil", // starts with a digit
			"with-dash": "evil", // hyphen not allowed in shell name
		},
		"", "",
	)
	script := got[2]
	// Only the valid key is exported; malicious keys are skipped entirely.
	assert.Contains(t, script, "export GOOD=1")
	assert.NotContains(t, script, "id #")
	assert.NotContains(t, script, "reboot")
	assert.NotContains(t, script, "1BAD")
	assert.NotContains(t, script, "with-dash")
	// exactly one export line
	assert.Equal(t, 1, strings.Count(script, "export "))
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
