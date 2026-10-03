// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package envcheck

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckGit(t *testing.T) {
	require.NoError(t, CheckGit(t.Context()))

	minimum := parseVersion(MinGitVersion)
	for text, atLeastMinimum := range map[string]bool{
		"git version 2.34.1":                 true,
		"git version 2.39.5 (Apple Git-154)": true,
		"git version 2.45.1.windows.1":       true,
		"git version 3.0":                    true,
		"git version 2.34.0":                 false,
		"git version 2.25.1":                 false,
		"garbled":                            false,
	} {
		assert.Equal(t, atLeastMinimum, slices.Compare(parseVersion(text), minimum) >= 0, text)
	}
}
