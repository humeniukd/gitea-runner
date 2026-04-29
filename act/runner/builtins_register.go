// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

// This file blank-imports every built-in action that ships with the
// runner so each built-in's `init()` registers it as soon as the runner
// package is linked into the binary. See docs/builtin-actions.md.
//
// To add a new built-in: drop a sub-package under `act/runner/builtin/`,
// register from its `init()`, and add a blank import here.

import (
	// Importing the built-in sub-packages for their `init()`-time
	// `builtin.Register` calls.
	_ "gitea.com/gitea/act_runner/act/runner/builtin/checkout"
)
