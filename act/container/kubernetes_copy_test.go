// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"archive/tar"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildFileEntryTar(t *testing.T) {
	const eventBody = "event-payload-contents"
	r, err := buildFileEntryTar([]*FileEntry{
		{Name: "workflow/event.json", Mode: 0o644, Body: eventBody},
		{Name: "workflow/envs.txt", Mode: 0o666, Body: ""},
	})
	require.NoError(t, err)

	tr := tar.NewReader(r)

	hdr, err := tr.Next()
	require.NoError(t, err)
	assert.Equal(t, "workflow/event.json", hdr.Name)
	assert.Equal(t, int64(0o644), hdr.Mode)
	body, err := io.ReadAll(tr)
	require.NoError(t, err)
	assert.Equal(t, eventBody, string(body))

	hdr, err = tr.Next()
	require.NoError(t, err)
	assert.Equal(t, "workflow/envs.txt", hdr.Name)
	assert.Equal(t, int64(0o666), hdr.Mode)
	assert.Equal(t, int64(0), hdr.Size)

	_, err = tr.Next()
	assert.ErrorIs(t, err, io.EOF)
}

func TestShellQuoteSingle(t *testing.T) {
	assert.Equal(t, "'/var/run/act/'", shellQuoteSingle("/var/run/act/"))
	// embedded single quotes are escaped so the script stays well-formed
	assert.Equal(t, `'a'\''b'`, shellQuoteSingle("a'b"))
}
