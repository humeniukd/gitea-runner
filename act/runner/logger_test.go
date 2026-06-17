// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestValueMasker(t *testing.T) {
	table := []struct {
		name       string
		lines      string
		secrets    map[string]string
		masks      []string
		disallowed []string
	}{
		{
			name:  "Multiline Private Key",
			lines: "cat << EOF > private.key\nPRIVATE_KEY_BEGIN\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\nPRIVATE_KEY_END\nEOF",
			secrets: map[string]string{
				"PRIVATE_KEY": "PRIVATE_KEY_BEGIN\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\nPRIVATE_KEY_END",
			},
			disallowed: []string{"KEY", "dsdfseffefsefes", "PRIVATE_KEY_END"},
		},
		{
			name:       "Multiline Private Key in masks",
			lines:      "cat << EOF > private.key\nPRIVATE_KEY_BEGIN\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\nPRIVATE_KEY_END\nEOF",
			masks:      []string{"PRIVATE_KEY_BEGIN\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\ndsdfseffefsefes\nPRIVATE_KEY_END"},
			disallowed: []string{"KEY", "dsdfseffefsefes", "PRIVATE_KEY_END"},
		},
		{
			name:  "Reversed secret via rev command",
			lines: "terces_ym",
			secrets: map[string]string{
				"MY_SECRET": "my_secret",
			},
			disallowed: []string{"terces_ym"},
		},
		{
			name:  "Substring leak via cut or fold",
			lines: "my_\nsec\nret",
			secrets: map[string]string{
				"MY_SECRET": "my_secret",
			},
			disallowed: []string{"my_", "sec", "ret"},
		},
		{
			name:  "Single char leak via fold -w1 or array index",
			lines: "m\ny\n_\ns\ne\nc\nr\ne\nt",
			secrets: map[string]string{
				"MY_SECRET": "my_secret",
			},
			disallowed: []string{"m", "y", "_", "s", "e", "c", "r", "t"},
		},
	}
	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			ctx := WithMasks(t.Context(), &entry.masks)
			masker := valueMasker(false, entry.secrets)
			for line := range strings.SplitSeq(entry.lines, "\n") {
				lentry := masker(&logrus.Entry{
					Context: ctx,
					Message: line,
				})
				for _, line := range entry.disallowed {
					assert.NotContains(t, lentry.Message, line)
				}
			}
		})
	}
}
