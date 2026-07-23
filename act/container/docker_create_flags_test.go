// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFlagsFromOptions(t *testing.T) {
	for _, tc := range []struct {
		options  string
		platform string
		pull     string
	}{
		{"", "", pullPolicyMissing},
		{"--privileged", "", pullPolicyMissing},
		{"-v /a:/b --platform=linux/arm64 -u 1000", "linux/arm64", pullPolicyMissing},
		{"--platform linux/arm/v7 --pull always", "linux/arm/v7", pullPolicyAlways},
		{"--pull never", "", pullPolicyNever},
		{`--platform "linux/amd64`, "", pullPolicyMissing}, // malformed, defaults kept
	} {
		t.Run(tc.options, func(t *testing.T) {
			cf := createFlagsFromOptions(tc.options)
			assert.Equal(t, tc.platform, cf.platform)
			assert.Equal(t, tc.pull, cf.pull)
		})
	}
}

func TestCreateFlagsValidate(t *testing.T) {
	for _, tc := range []struct {
		options string
		wantErr string
	}{
		{"--pull always", ""},
		{"--quiet --disable-content-trust --name mine", ""},
		{"--pull sometimes", `invalid --pull option "sometimes"`},
		{"--use-api-socket", "--use-api-socket is not supported"},
	} {
		t.Run(tc.options, func(t *testing.T) {
			err := createFlagsFromOptions(tc.options).validate()
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestNewContainerAppliesCreateFlags(t *testing.T) {
	t.Run("platform overrides the runner-wide architecture", func(t *testing.T) {
		input := &NewContainerInput{Platform: "linux/amd64", Options: "--platform linux/arm64"}
		NewContainer(input)
		assert.Equal(t, "linux/arm64", input.Platform)
	})

	t.Run("architecture is kept without a platform in the options", func(t *testing.T) {
		input := &NewContainerInput{Platform: "linux/amd64", Options: "--privileged"}
		NewContainer(input)
		assert.Equal(t, "linux/amd64", input.Platform)
	})

	t.Run("pull policy is taken from the options", func(t *testing.T) {
		cr := NewContainer(&NewContainerInput{Options: "--pull never"}).(*containerReference)
		assert.Equal(t, pullPolicyNever, cr.pullPolicy)
	})
}
