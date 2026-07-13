// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package process

import (
	"path/filepath"
	"strings"
)

// pathWithin is the separator- and OS-agnostic core of dirContainsPath. dir and
// target must already be cleaned. When fold is true the comparison is
// case-insensitive (Windows filesystems). It is a pure function so the boundary
// and case rules can be unit-tested without touching the process table.
func pathWithin(dir, target string, fold bool) bool {
	if dir == "" || target == "" || dir == "." || target == "." {
		return false
	}
	if fold {
		dir = strings.ToLower(dir)
		target = strings.ToLower(target)
	}
	if dir == target {
		return true
	}
	sep := string(filepath.Separator)
	if !strings.HasSuffix(dir, sep) {
		dir += sep
	}
	return strings.HasPrefix(target, dir)
}
