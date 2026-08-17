// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package disk

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFreeBytes(t *testing.T) {
	free, err := FreeBytes(t.TempDir())
	switch runtime.GOOS {
	case "aix", "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd", "solaris", "windows":
		require.NoError(t, err)
		assert.Positive(t, free)
	default:
		require.Error(t, err, "unsupported platforms report unavailable rather than a full disk")
	}

	_, err = FreeBytes(filepath.Join(t.TempDir(), "does-not-exist"))
	require.Error(t, err)
}
