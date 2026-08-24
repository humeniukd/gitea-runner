// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaskingFormatter(t *testing.T) {
	formatter := MaskingFormatter(&log.TextFormatter{DisableTimestamp: true})
	format := func() string {
		line, err := formatter.Format(log.WithField("job", "build s3cr3t and other"))
		require.NoError(t, err)
		return string(line)
	}

	assert.Contains(t, format(), "s3cr3t") // nothing registered yet

	first := registerGlobalMasks([]string{"s3cr3t", "***"})
	t.Cleanup(func() { deregisterGlobalMasks(first) })
	second := registerGlobalMasks([]string{"other", "***"})
	t.Cleanup(func() { deregisterGlobalMasks(second) })

	line := format()
	assert.NotContains(t, line, "s3cr3t") // masked though it rode a field, not the message
	assert.NotContains(t, line, "other")

	deregisterGlobalMasks(first)
	line = format()
	assert.Contains(t, line, "s3cr3t")   // that task ended
	assert.NotContains(t, line, "other") // the task still running keeps its own
}

func TestAddMaskReachesGlobalMasker(t *testing.T) {
	r := &Reporter{globalMaskID: registerGlobalMasks(nil)}
	t.Cleanup(func() { deregisterGlobalMasks(r.globalMaskID) })

	r.addMask("runtime-secret")

	assert.NotContains(t, maskGlobal("saw runtime-secret"), "runtime-secret")
}

func TestGlobalMasksDoNotTruncateEachOther(t *testing.T) {
	short := registerGlobalMasks([]string{"abc", "***"})
	t.Cleanup(func() { deregisterGlobalMasks(short) })
	long := registerGlobalMasks([]string{"abcdef", "***"})
	t.Cleanup(func() { deregisterGlobalMasks(long) })

	assert.Equal(t, "***", maskGlobal("abcdef")) // not "***def", whichever task registered first
}
