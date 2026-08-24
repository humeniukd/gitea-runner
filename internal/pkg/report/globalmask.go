// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"strings"
	"sync"
	"sync/atomic"

	"gitea.com/gitea/runner/act/runner"

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
	var oldnew []string
	for _, masks := range globalMasks {
		oldnew = append(oldnew, masks...)
	}
	if len(oldnew) == 0 {
		globalReplacer.Store(nil)
		return
	}
	globalReplacer.Store(runner.NewSecretReplacer(oldnew))
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
	if replacer == nil { // no task is live, so this costs what it did before the wrapper
		return line, nil
	}
	return []byte(replacer.Replace(string(line))), nil
}
