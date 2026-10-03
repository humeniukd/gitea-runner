// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package envcheck

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const MinGitVersion = "2.34.1" // Ubuntu 22.04

var gitVersionRegex = regexp.MustCompile(`(\d+)\.(\d+)(?:\.(\d+))?`)

func CheckGit(ctx context.Context) error {
	out, err := exec.CommandContext(ctx, "git", "version").Output()
	if err != nil {
		return fmt.Errorf("git %s or newer must be on PATH: %w", MinGitVersion, err)
	}
	if slices.Compare(parseVersion(string(out)), parseVersion(MinGitVersion)) < 0 {
		return fmt.Errorf("git %s or newer must be on PATH, found %q", MinGitVersion, strings.TrimSpace(string(out)))
	}
	return nil
}

func parseVersion(text string) []int {
	match := gitVersionRegex.FindStringSubmatch(text)
	if match == nil {
		return nil
	}
	version := make([]int, 0, 3)
	for _, part := range match[1:] {
		number, _ := strconv.Atoi(part) // missing patch is 0
		version = append(version, number)
	}
	return version
}
