// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePlatform(t *testing.T) {
	t.Run("empty input returns nil platform without error", func(t *testing.T) {
		got, err := parsePlatform("")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("os/arch", func(t *testing.T) {
		got, err := parsePlatform("linux/amd64")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "linux", got.OS)
		assert.Equal(t, "amd64", got.Architecture)
		assert.Empty(t, got.Variant)
	})

	t.Run("os/arch/variant", func(t *testing.T) {
		got, err := parsePlatform("linux/arm/v7")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "linux", got.OS)
		assert.Equal(t, "arm", got.Architecture)
		assert.Equal(t, "v7", got.Variant)
	})

	t.Run("input is lowercased", func(t *testing.T) {
		got, err := parsePlatform("Linux/AMD64/V8")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "linux", got.OS)
		assert.Equal(t, "amd64", got.Architecture)
		assert.Equal(t, "v8", got.Variant)
	})

	for _, bad := range []string{
		"amd64",
		"linux",
		"linux/",
		"/amd64",
		"/",
		"//",
		"linux/arm/",
		"linux/arm/v7/extra",
	} {
		t.Run("rejects "+bad, func(t *testing.T) {
			got, err := parsePlatform(bad)
			require.Error(t, err)
			assert.Nil(t, got)
		})
	}
}

func TestPlatformMatches(t *testing.T) {
	for _, tt := range []struct {
		name              string
		want              string
		os, arch, variant string
		matches           bool
	}{
		{name: "exact os and arch", want: "linux/amd64", os: "linux", arch: "amd64", matches: true},
		{name: "requested variant matches image variant", want: "linux/arm64/v8", os: "linux", arch: "arm64", variant: "v8", matches: true},
		{name: "unrequested variant accepts any image variant", want: "linux/arm64", os: "linux", arch: "arm64", variant: "v8", matches: true},
		{name: "requested variant rejects a bare image", want: "linux/arm64/v8", os: "linux", arch: "arm64", matches: false},
		{name: "requested variant rejects a different variant", want: "linux/arm64/v8", os: "linux", arch: "arm64", variant: "v7", matches: false},
		{name: "different architecture", want: "linux/amd64", os: "linux", arch: "arm64", matches: false},
		{name: "different os", want: "windows/amd64", os: "linux", arch: "amd64", matches: false},
		{name: "image casing is ignored", want: "linux/arm64/v8", os: "Linux", arch: "ARM64", variant: "V8", matches: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			want, err := parsePlatform(tt.want)
			require.NoError(t, err)
			assert.Equal(t, tt.matches, platformMatches(want, tt.os, tt.arch, tt.variant))
		})
	}
}
