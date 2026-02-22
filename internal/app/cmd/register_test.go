// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterNonInteractiveReturnsLabelValidationError(t *testing.T) {
	err := registerNoInteractive(t.Context(), "", &registerArgs{
		Labels:       "label:invalid",
		Token:        "token",
		InstanceAddr: "http://localhost:3000",
	})
	require.Error(t, err, "unsupported schema: invalid")
}
