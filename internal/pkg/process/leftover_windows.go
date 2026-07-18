// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package process

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

// dirContainsPath reports whether target is dir itself or lives beneath it.
// Matching is case-insensitive, and a sibling sharing a name prefix (job1 vs
// job10) is not contained.
func dirContainsPath(dir, target string) bool {
	return pathWithin(filepath.Clean(dir), filepath.Clean(target), true)
}

// KillProcessesWithCWDUnder terminates every process whose current working
// directory is one of dirs or lives beneath it, and returns how many it killed.
// This complements the workspace path scan, which cannot see a process that only
// *runs in* a workspace dir: Win32_Process exposes no working directory. Such a
// process pins an open handle that blocks the workspace deletion.
//
// Only processes in the runner's own session are considered: leftovers we
// spawned inherit it, and skipping the rest avoids reading the memory of
// system and other-user processes we could never open anyway.
//
// Best-effort: a process we cannot open (higher integrity level, another user)
// is skipped.
func KillProcessesWithCWDUnder(dirs []string) (int, error) {
	if len(dirs) == 0 {
		return 0, nil
	}

	pids, err := listProcessIDs()
	if err != nil {
		return 0, err
	}

	self := uint32(os.Getpid())
	var ownSession uint32
	if err := windows.ProcessIdToSessionId(self, &ownSession); err != nil {
		return 0, fmt.Errorf("look up own session: %w", err)
	}

	var killed int
	var errs []error
	for _, pid := range pids {
		// Never touch ourselves, System Idle (0) or System (4).
		if pid == self || pid == 0 || pid == 4 {
			continue
		}

		// Restrict the memory-reading scan to our own session.
		var session uint32
		if err := windows.ProcessIdToSessionId(pid, &session); err != nil || session != ownSession {
			continue
		}

		ok, err := terminateIfCWDUnder(pid, dirs)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if ok {
			killed++
		}
	}

	return killed, errors.Join(errs...)
}

// terminateIfCWDUnder opens pid once and both reads the cwd and terminates
// through that one handle. The shared handle binds to this exact process, so a
// PID recycled onto a new process between reading and killing cannot be hit.
// Returns true when it terminated the process.
func terminateIfCWDUnder(pid uint32, dirs []string) (bool, error) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ|windows.PROCESS_TERMINATE, false, pid)
	if err != nil {
		// Best-effort: cannot open (gone, higher integrity, another user).
		return false, nil
	}
	defer func() { _ = windows.CloseHandle(h) }()

	cwd, err := processCWDHandle(h)
	if err != nil || cwd == "" {
		return false, nil
	}

	matched := false
	for _, d := range dirs {
		if dirContainsPath(d, cwd) {
			matched = true
			break
		}
	}
	if !matched {
		return false, nil
	}

	if err := windows.TerminateProcess(h, 1); err != nil {
		return false, fmt.Errorf("terminate pid %d: %w", pid, err)
	}
	return true, nil
}

// listProcessIDs returns the PIDs of every process in a Toolhelp snapshot.
func listProcessIDs() ([]uint32, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer func() { _ = windows.CloseHandle(snap) }()

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	pids := make([]uint32, 0, 256)
	err = windows.Process32First(snap, &entry)
	for err == nil {
		pids = append(pids, entry.ProcessID)
		err = windows.Process32Next(snap, &entry)
	}
	if !errors.Is(err, windows.ERROR_NO_MORE_FILES) {
		return pids, err
	}
	return pids, nil
}

// processCWD opens pid and reads its working directory. It exists for tests;
// the reaper uses processCWDHandle so it can share one handle with the kill.
func processCWD(pid uint32) (string, error) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pid)
	if err != nil {
		return "", err
	}
	defer func() { _ = windows.CloseHandle(h) }()
	return processCWDHandle(h)
}

// processCWDHandle reads a process's working directory out of its PEB, walking
// PEB -> ProcessParameters -> CurrentDirectory.DosPath. Offsets come from
// unsafe.Offsetof so they stay correct for the build architecture. h must carry
// PROCESS_QUERY_INFORMATION|PROCESS_VM_READ.
func processCWDHandle(h windows.Handle) (string, error) {
	var pbi windows.PROCESS_BASIC_INFORMATION
	var retLen uint32
	if err := windows.NtQueryInformationProcess(h, windows.ProcessBasicInformation, unsafe.Pointer(&pbi), uint32(unsafe.Sizeof(pbi)), &retLen); err != nil {
		return "", err
	}
	if pbi.PebBaseAddress == nil {
		return "", errors.New("nil PEB base address")
	}
	pebAddr := uintptr(unsafe.Pointer(pbi.PebBaseAddress))

	paramsAddr, err := readRemotePtr(h, pebAddr+unsafe.Offsetof(windows.PEB{}.ProcessParameters))
	if err != nil || paramsAddr == 0 {
		return "", err
	}

	// DosPath is the leading NTUnicodeString of CURDIR, so it shares the
	// CurrentDirectory address.
	dosPathAddr := paramsAddr + unsafe.Offsetof(windows.RTL_USER_PROCESS_PARAMETERS{}.CurrentDirectory)
	length, err := readRemoteU16(h, dosPathAddr+unsafe.Offsetof(windows.NTUnicodeString{}.Length))
	if err != nil || length == 0 {
		return "", err
	}
	bufAddr, err := readRemotePtr(h, dosPathAddr+unsafe.Offsetof(windows.NTUnicodeString{}.Buffer))
	if err != nil || bufAddr == 0 {
		return "", err
	}

	u16 := make([]uint16, length/2)
	var n uintptr
	if err := windows.ReadProcessMemory(h, bufAddr, (*byte)(unsafe.Pointer(&u16[0])), uintptr(length), &n); err != nil {
		return "", err
	}
	if n != uintptr(length) {
		return "", errors.New("short read of working directory")
	}
	return windows.UTF16ToString(u16), nil
}

// readRemotePtr reads a single pointer-sized value from another process.
func readRemotePtr(h windows.Handle, addr uintptr) (uintptr, error) {
	var v uintptr
	var n uintptr
	if err := windows.ReadProcessMemory(h, addr, (*byte)(unsafe.Pointer(&v)), unsafe.Sizeof(v), &n); err != nil {
		return 0, err
	}
	if n != unsafe.Sizeof(v) {
		return 0, errors.New("short pointer read")
	}
	return v, nil
}

// readRemoteU16 reads a single uint16 from another process.
func readRemoteU16(h windows.Handle, addr uintptr) (uint16, error) {
	var v uint16
	var n uintptr
	if err := windows.ReadProcessMemory(h, addr, (*byte)(unsafe.Pointer(&v)), unsafe.Sizeof(v), &n); err != nil {
		return 0, err
	}
	if n != unsafe.Sizeof(v) {
		return 0, errors.New("short uint16 read")
	}
	return v, nil
}
