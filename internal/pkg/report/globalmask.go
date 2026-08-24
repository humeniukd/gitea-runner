// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"slices"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// A job logger masks one task's secrets, but a job's plan runs before that logger exists and
// logs through the process-wide one, which several tasks share. So that one holds the union of
// every live task's secrets: masking one task's output with another's is harmless, missing one
// is not.
var (
	globalMu       sync.RWMutex
	globalMasks    = map[uint64][]string{}
	globalMaskSeq  uint64
	globalReplacer = strings.NewReplacer()
)

func registerGlobalMasks(oldnew []string) uint64 {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalMaskSeq++
	globalMasks[globalMaskSeq] = oldnew
	rebuildGlobalReplacer()
	return globalMaskSeq
}

// updateGlobalMasks follows a task whose secrets grew, such as by ::add-mask::. An id that is
// not registered, including the zero one a Reporter built outside NewReporter carries, is a no-op.
func updateGlobalMasks(id uint64, oldnew []string) {
	globalMu.Lock()
	defer globalMu.Unlock()
	if _, ok := globalMasks[id]; !ok {
		return
	}
	globalMasks[id] = oldnew
	rebuildGlobalReplacer()
}

func deregisterGlobalMasks(id uint64) {
	globalMu.Lock()
	defer globalMu.Unlock()
	delete(globalMasks, id)
	rebuildGlobalReplacer()
}

func rebuildGlobalReplacer() { // caller holds globalMu
	var pairs [][2]string
	for _, masks := range globalMasks {
		for i := 0; i+1 < len(masks); i += 2 {
			pairs = append(pairs, [2]string{masks[i], masks[i+1]})
		}
	}
	// Replacer matches in argument order, so one task's secret that prefixes another's would
	// mask only that prefix and print the rest. Longest first, and map order stops mattering.
	slices.SortFunc(pairs, func(a, b [2]string) int { return len(b[0]) - len(a[0]) })
	_ = slices.Sort[[]int]

	oldnew := make([]string, 0, len(pairs)*2)
	for _, pair := range pairs {
		oldnew = append(oldnew, pair[0], pair[1])
	}
	globalReplacer = strings.NewReplacer(oldnew...)
}

func maskGlobal(content string) string {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return strings.ToValidUTF8(globalReplacer.Replace(content), "?")
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
		return nil, err
	}
	return []byte(maskGlobal(string(line))), nil
}
