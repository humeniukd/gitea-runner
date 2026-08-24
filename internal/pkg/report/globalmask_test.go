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
	second := register("other", "***", "otherlonger", "***")

	line := format("build s3cr3t and other")
	assert.NotContains(t, line, "s3cr3t") // masked though it rode a field, not the message
	assert.NotContains(t, line, "other")
	assert.NotContains(t, format("otherlonger"), "longer") // longest first, so not "***longer"

	deregisterGlobalMasks(first)
	line = format("build s3cr3t and other")
	assert.Contains(t, line, "s3cr3t")   // that task ended
	assert.NotContains(t, line, "other") // the one still running keeps its own

	deregisterGlobalMasks(second)
	assert.Nil(t, globalReplacer.Load()) // no task live, so the process log pays nothing

	// A workflow could otherwise mask "error" here and rewrite every other task's log.
	r := &Reporter{globalMaskID: register("s3cr3t", "***"), oldnew: []string{"s3cr3t", "***"}}
	r.addMask("runtime-secret")
	assert.Contains(t, format("saw runtime-secret"), "runtime-secret")
	assert.NotContains(t, r.mask("saw runtime-secret"), "runtime-secret") // the job's own log still masks it
}
