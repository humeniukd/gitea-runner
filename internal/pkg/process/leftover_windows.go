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
	var killed int
	var errs []error
	for _, pid := range pids {
		// Never touch ourselves, System Idle (0) or System (4).
		if pid == self || pid == 0 || pid == 4 {
			continue
		}

		cwd, err := processCWD(pid)
		if err != nil || cwd == "" {
			continue
		}

		matched := false
		for _, d := range dirs {
			if dirContainsPath(d, cwd) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		if err := terminateProcess(pid); err != nil {
			errs = append(errs, fmt.Errorf("terminate pid %d: %w", pid, err))
			continue
		}
		killed++
	}

	return killed, errors.Join(errs...)
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

// terminateProcess force-terminates a single process by PID.
func terminateProcess(pid uint32) error {
	h, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, pid)
	if err != nil {
		return err
	}
	defer func() { _ = windows.CloseHandle(h) }()
	return windows.TerminateProcess(h, 1)
}

// processCWD reads another process's working directory out of its PEB, walking
// PEB -> ProcessParameters -> CurrentDirectory.DosPath. Offsets come from
// unsafe.Offsetof so they stay correct for the build architecture.
func processCWD(pid uint32) (string, error) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pid)
	if err != nil {
		return "", err
	}
	defer func() { _ = windows.CloseHandle(h) }()

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
