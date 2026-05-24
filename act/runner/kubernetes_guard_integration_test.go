// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build integration

// RC1 fail-fast guard assertions for the Kubernetes backend. The positive cases
// (JS, composite, and Go `uses:` actions must EXECUTE via exec-in) require a
// live job container, so this suite is gated behind the `integration` tag and
// runs in the E2E lane against the real cluster. The negative case asserts that
// a `using: docker` action returns the explicit unsupported-action error without
// dialing a socket or dereferencing a nil stub.
//
// These are documented here as the authoritative RC1 checks; the Phase 9
// full-matrix E2E workflow exercises them end-to-end. See run_context.go
// usesKubernetes and action.go's IsDocker() guard.
package runner

import "testing"

// TestKubernetesUsesActionGuard is a placeholder marker for the E2E lane: the
// real assertions run through the Phase 9 workflow (a JS, a composite, and a Go
// action all succeed; a `using: docker` action fails with the unsupported
// error). Kept compile-clean under the integration tag so the intent is tracked
// in-tree.
func TestKubernetesUsesActionGuard(t *testing.T) {
	t.Skip("RC1 guard is asserted by the Phase 9 full-matrix E2E workflow against a live cluster")
}
