// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package lookpath

type Error struct {
	Name string
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}
