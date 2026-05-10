// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVolumeSpec(t *testing.T) {
	tests := []struct {
		name string
		spec string
		want volumeSpec
	}{
		{
			name: "single-rune target is anonymous volume",
			spec: "/",
			want: volumeSpec{Target: "/", Type: "volume"},
		},
		{
			name: "two-rune target is anonymous volume",
			spec: "ab",
			want: volumeSpec{Target: "ab", Type: "volume"},
		},
		{
			name: "absolute path target is anonymous volume",
			spec: "/data",
			want: volumeSpec{Target: "/data", Type: "volume"},
		},
		{
			name: "named volume with target",
			spec: "myvol:/container",
			want: volumeSpec{Source: "myvol", Target: "/container", Type: "volume"},
		},
		{
			name: "absolute host path is bind mount",
			spec: "/host:/container",
			want: volumeSpec{Source: "/host", Target: "/container", Type: "bind"},
		},
		{
			name: "relative path source is bind mount",
			spec: "./host:/container",
			want: volumeSpec{Source: "./host", Target: "/container", Type: "bind"},
		},
		{
			name: "home-prefixed source is bind mount",
			spec: "~/host:/container",
			want: volumeSpec{Source: "~/host", Target: "/container", Type: "bind"},
		},
		{
			name: "ro option sets ReadOnly",
			spec: "/host:/container:ro",
			want: volumeSpec{Source: "/host", Target: "/container", Type: "bind", ReadOnly: true},
		},
		{
			name: "rw option clears ReadOnly",
			spec: "/host:/container:rw",
			want: volumeSpec{Source: "/host", Target: "/container", Type: "bind", ReadOnly: false},
		},
		{
			name: "unknown option is silently dropped",
			spec: "/host:/container:Z",
			want: volumeSpec{Source: "/host", Target: "/container", Type: "bind"},
		},
		{
			name: "ro with extra unknown option",
			spec: "/host:/container:ro,Z",
			want: volumeSpec{Source: "/host", Target: "/container", Type: "bind", ReadOnly: true},
		},
		{
			name: "later option overrides earlier",
			spec: "/host:/container:rw,ro",
			want: volumeSpec{Source: "/host", Target: "/container", Type: "bind", ReadOnly: true},
		},
		{
			name: "windows drive letter source",
			spec: `C:\data:/container`,
			want: volumeSpec{Source: `C:\data`, Target: "/container", Type: "bind"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseVolumeSpec(tc.spec)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}

	errorCases := []struct {
		name string
		spec string
	}{
		{name: "empty spec", spec: ""},
		{name: "too many colons", spec: "/a:/b:/c:/d"},
		{name: "empty section between colons", spec: "/a::/b"},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseVolumeSpec(tc.spec)
			require.Error(t, err)
		})
	}
}
