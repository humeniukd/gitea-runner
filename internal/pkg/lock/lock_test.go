// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !plan9

package lock

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestTryLock(t *testing.T) {
	runnerFile := filepath.Join(t.TempDir(), ".runner")

	release, err := TryLock(runnerFile)
	if err != nil {
		t.Fatalf("first TryLock failed: %v", err)
	}

	// A second lock on the same file must be refused while the first is held.
	if _, err := TryLock(runnerFile); !errors.Is(err, ErrLocked) {
		t.Fatalf("second TryLock: want ErrLocked, got %v", err)
	}

	if err := release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}

	// After release the lock is available again.
	release2, err := TryLock(runnerFile)
	if err != nil {
		t.Fatalf("TryLock after release failed: %v", err)
	}
	if err := release2(); err != nil {
		t.Fatalf("second release failed: %v", err)
	}
}
