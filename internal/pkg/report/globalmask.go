// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"strings"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// A job logger masks one task's secrets, but the plan runs before any job logger exists and
// logs through the process-wide logger, which several tasks share. So that one holds the union
// of every live task's secrets: masking one task's output with another's secret is harmless,
// missing one is not.
var (
	globalMu       sync.Mutex
	globalMasks    = map[*Reporter][]string{}
	globalReplacer atomic.Pointer[strings.Replacer]
)

func registerGlobalMasks(key *Reporter, oldnew []string) func() {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalMasks[key] = oldnew
	rebuildGlobalReplacer()

	return func() {
		globalMu.Lock()
		defer globalMu.Unlock()
		delete(globalMasks, key)
		rebuildGlobalReplacer()
	}
}

func rebuildGlobalReplacer() { // caller holds globalMu
	oldnew := make([]string, 0, len(globalMasks)*2)
	for _, masks := range globalMasks {
		oldnew = append(oldnew, masks...)
	}
	globalReplacer.Store(strings.NewReplacer(oldnew...))
}

func maskGlobal(content string) string {
	replacer := globalReplacer.Load()
	if replacer == nil {
		return content
	}
	return replacer.Replace(content)
}

// MaskingFormatter wraps f so no live task's secret reaches the process-wide log. It masks the
// rendered line, so fields are covered as well as the message.
func MaskingFormatter(f log.Formatter) log.Formatter {
	return &maskingFormatter{inner: f}
}

type maskingFormatter struct{ inner log.Formatter }

func (m *maskingFormatter) Format(entry *log.Entry) ([]byte, error) {
	line, err := m.inner.Format(entry)
	if err != nil {
		return line, err
	}
	return []byte(maskGlobal(string(line))), nil
}
