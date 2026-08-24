// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalMasks(t *testing.T) {
	formatter := MaskingFormatter(&log.TextFormatter{DisableTimestamp: true})
	format := func(job string) string {
		line, err := formatter.Format(log.WithField("job", job))
		require.NoError(t, err)
		return string(line)
	}
	register := func(oldnew ...string) uint64 {
		id := registerGlobalMasks(oldnew)
		t.Cleanup(func() { deregisterGlobalMasks(id) })
		return id
	}

	assert.Contains(t, format("build s3cr3t"), "s3cr3t") // nothing registered yet

	first := register("s3cr3t", "***")
	register("other", "***")

	line := format("build s3cr3t and other")
	assert.NotContains(t, line, "s3cr3t") // masked though it rode a field, not the message
	assert.NotContains(t, line, "other")

	deregisterGlobalMasks(first)
	line = format("build s3cr3t and other")
	assert.Contains(t, line, "s3cr3t")   // that task ended
	assert.NotContains(t, line, "other") // the one still running keeps its own

	register("otherlonger", "***")
	assert.NotContains(t, format("otherlonger"), "longer") // not "***longer", whichever registered first

	r := &Reporter{globalMaskID: register()}
	r.addMask("runtime-secret")
	assert.NotContains(t, format("saw runtime-secret"), "runtime-secret")
}
