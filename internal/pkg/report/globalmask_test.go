// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"strings"
	"testing"

	runnerv1 "gitea.dev/actionslib/runner/v1"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalMasks(t *testing.T) {
	formatter := MaskingFormatter(&log.TextFormatter{DisableTimestamp: true})
	format := func() string {
		line, err := formatter.Format(log.WithField("job", "build s3cr3t and other"))
		require.NoError(t, err)
		return string(line)
	}
	register := func(oldnew ...string) uint64 {
		id := registerGlobalMasks(oldnew)
		t.Cleanup(func() { deregisterGlobalMasks(id) })
		return id
	}

	assert.Contains(t, format(), "s3cr3t") // nothing registered yet

	first := register("s3cr3t", "***")
	register("other", "***")

	line := format()
	assert.NotContains(t, line, "s3cr3t") // masked though it rode a field, not the message
	assert.NotContains(t, line, "other")

	deregisterGlobalMasks(first)
	line = format()
	assert.Contains(t, line, "s3cr3t")   // that task ended
	assert.NotContains(t, line, "other") // the one still running keeps its own

	register("otherlonger", "***")
	assert.Equal(t, "***", maskGlobal("otherlonger")) // not "***longer", whichever registered first

	r := &Reporter{globalMaskID: register()}
	r.AddMask("runtime-secret")
	assert.NotContains(t, maskGlobal("saw runtime-secret"), "runtime-secret")
}

func TestSetOutputsSkipsSecrets(t *testing.T) {
	r := &Reporter{logReplacer: strings.NewReplacer("s3cr3t", "***"), state: &runnerv1.TaskState{}}

	r.SetOutputs(map[string]string{"safe": "value", "leaky": "has s3cr3t in it"})

	assert.Contains(t, r.outputs, "safe")
	assert.NotContains(t, r.outputs, "leaky") // GitHub skips it rather than sending a masked value
}
