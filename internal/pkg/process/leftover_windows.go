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
// Both paths are cleaned first; matching is case-insensitive, as Windows
// filesystems are. The boundary check appends a separator so a sibling that only
// shares a name prefix (job1 vs job10) is not considered contained.
func dirContainsPath(dir, target string) bool {
	return pathWithin(filepath.Clean(dir), filepath.Clean(target), true)
}

// KillProcessesWithCWDUnder terminates every running process whose current
// working directory is one of dirs or lives beneath it.
//
// It complements the workspace path scan (which matches on a process's
// executable path and command line): Win32_Process exposes neither the working
// directory, so a process that was merely *started in* a workspace directory —
// e.g. a compiler or language runtime resolved from PATH, launched with its
// working directory inside the checkout and only relative arguments — is
// invisible to that scan. Such a process survives teardown and, because its
// current directory pins an open handle on that folder, can also block the
// workspace deletion that follows. Children inherit their parent's working
// directory, so a whole spawned tree is caught as each member matches in turn.
//
// It is best-effort: a process we cannot open (higher integrity level or a
// different user) is skipped, mirroring the path scan's inability to read the
// ExecutablePath of such processes. It returns the number of processes
// terminated and the joined errors from any that matched but could not be
// killed.
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
		// Skip our own process, the System Idle Process (0) and System (4);
		// none of them can legitimately have a workspace as its cwd and none
		// can (or should) be terminated here.
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

// processCWD reads the current working directory of another process out of its
// PEB. It walks PEB -> ProcessParameters -> CurrentDirectory.DosPath, reading
// each pointer and the final UTF-16 string over the process boundary with
// ReadProcessMemory. Field offsets come from unsafe.Offsetof on the typed
// x/sys/windows structs so they stay correct for the build architecture; a
// 64-bit runner reads the 64-bit PEB of both 32- and 64-bit targets.
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

	// PEB.ProcessParameters (a remote pointer).
	paramsAddr, err := readRemotePtr(h, pebAddr+unsafe.Offsetof(windows.PEB{}.ProcessParameters))
	if err != nil || paramsAddr == 0 {
		return "", err
	}

	// RTL_USER_PROCESS_PARAMETERS.CurrentDirectory.DosPath is the leading
	// NTUnicodeString of the CURDIR, so its address is the CurrentDirectory
	// address. Read its Length (in bytes) and Buffer pointer.
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
