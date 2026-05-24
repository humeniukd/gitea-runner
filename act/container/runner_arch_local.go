// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import "runtime"

// RunnerArchLocal resolves the GitHub-actions-compatible runner.arch value from
// the local process architecture (runtime.GOARCH) WITHOUT contacting a docker
// daemon. Unlike RunnerArch — which pings the docker server to report its
// architecture — this resolver is safe on backends that have no host docker
// socket (e.g. the Kubernetes pod backend), where probing /var/run/docker.sock
// only produces noisy, non-fatal "docker daemon ping failed" log lines.
//
// https://github.com/github/docs/blob/main/data/reusables/actions/runner-arch-description.md
func RunnerArchLocal() string {
	archMapper := map[string]string{
		"amd64":   "X64",
		"x86_64":  "X64",
		"386":     "X86",
		"arm64":   "ARM64",
		"aarch64": "ARM64",
	}
	if arch, ok := archMapper[runtime.GOARCH]; ok {
		return arch
	}
	return runtime.GOARCH
}
