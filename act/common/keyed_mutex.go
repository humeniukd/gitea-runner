// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package common

import "sync"

// KeyedMutex serializes the work done under one key while letting different keys run in
// parallel. Its zero value is ready to use, and it keeps one mutex per key it has seen.
type KeyedMutex[K comparable] struct {
	mu     sync.Mutex
	locked map[K]*sync.Mutex
}

// Lock locks the mutex for key and returns its unlock function.
func (km *KeyedMutex[K]) Lock(key K) func() {
	km.mu.Lock()
	if km.locked == nil {
		km.locked = map[K]*sync.Mutex{}
	}
	lock, ok := km.locked[key]
	if !ok {
		lock = &sync.Mutex{}
		km.locked[key] = lock
	}
	km.mu.Unlock()

	lock.Lock()
	return lock.Unlock
}

// Delete drops the mutex for key, so a KeyedMutex whose keys come and go does not grow
// unbounded. Holders of the mutex keep working on it, they just no longer share it with a
// later Lock of the same key.
func (km *KeyedMutex[K]) Delete(key K) {
	km.mu.Lock()
	delete(km.locked, key)
	km.mu.Unlock()
}
