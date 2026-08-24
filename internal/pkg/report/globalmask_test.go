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
	entry := func() *log.Entry {
		return log.WithField("job", "build s3cr3t").WithField("msg", "x")
	}

	line, err := formatter.Format(entry())
	require.NoError(t, err)
	assert.Contains(t, string(line), "s3cr3t") // nothing registered yet

	first := &Reporter{}
	stopFirst := registerGlobalMasks(first, []string{"s3cr3t", "***"})
	second := &Reporter{}
	stopSecond := registerGlobalMasks(second, []string{"other", "***"})

	line, err = formatter.Format(entry())
	require.NoError(t, err)
	assert.NotContains(t, string(line), "s3cr3t") // masked even though it is a field, not the message

	stopFirst()
	line, err = formatter.Format(entry())
	require.NoError(t, err)
	assert.Contains(t, string(line), "s3cr3t") // that task is gone, the other task's masks remain

	stopSecond()
	require.Empty(t, globalMasks)
}
