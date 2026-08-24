// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"bytes"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// A job's plan runs before its job logger exists and logs through the process-wide one, which
// several tasks share, so that one holds the union of every live task's secrets.
var (
	globalMu       sync.Mutex
	globalMasks    = map[uint64][]string{}
	globalMaskSeq  uint64
	globalReplacer atomic.Pointer[strings.Replacer] // nil while no task is live
)

func registerGlobalMasks(oldnew []string) uint64 {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalMaskSeq++
	globalMasks[globalMaskSeq] = oldnew
	rebuildGlobalReplacer()

	return globalMaskSeq
}

func updateGlobalMasks(id uint64, oldnew []string) { // an unregistered id, such as zero, is a no-op
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
	total := 0
	for _, masks := range globalMasks {
		total += len(masks) / 2
	}
	if total == 0 {
		globalReplacer.Store(nil)
		return
	}

	pairs := make([][2]string, 0, total)
	for _, masks := range globalMasks {
		for i := 0; i+1 < len(masks); i += 2 {
			pairs = append(pairs, [2]string{masks[i], masks[i+1]})
		}
	}
	// Replacer matches in argument order, so a secret prefixing another's would mask only that.
	slices.SortFunc(pairs, func(a, b [2]string) int { return len(b[0]) - len(a[0]) })

	oldnew := make([]string, 0, len(pairs)*2)
	for _, pair := range pairs {
		oldnew = append(oldnew, pair[0], pair[1])
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

// MaskingFormatter wraps f so no live task's secret reaches the process-wide log, fields included.
func MaskingFormatter(f log.Formatter) log.Formatter {
	return &maskingFormatter{inner: f}
}

type maskingFormatter struct{ inner log.Formatter }

func (m *maskingFormatter) Format(entry *log.Entry) ([]byte, error) {
	line, err := m.inner.Format(entry)
	if err != nil {
		return nil, err
	}
	replacer := globalReplacer.Load()
	if replacer == nil {
		return line, nil
	}

	var masked bytes.Buffer // streams, rather than a []byte to string and back round trip
	masked.Grow(len(line))
	if _, err := replacer.WriteString(&masked, string(line)); err != nil {
		return nil, err
	}
	return masked.Bytes(), nil
}
