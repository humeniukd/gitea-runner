// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

// Package e2e runs the runner against a real Gitea server (a gitea/gitea container) to catch
// compatibility regressions across Gitea versions.
//
// It is local-only and never runs as part of `make test`/`go test ./...`: every file here
// carries the `e2e` build tag, and the tests self-skip if no Docker daemon is reachable.
//
// Run it with `make test-e2e` (single version, GITEA_E2E_IMAGE overridable) or
// `make test-e2e-matrix` (loops the versions listed in e2e/.env). See DEVELOPMENT.md for
// details on what each scenario covers and the Gitea REST endpoints relied on.
package e2e
