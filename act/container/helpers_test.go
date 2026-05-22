// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	"net"
	"testing"
	"time"

	mobyclient "github.com/moby/moby/client"
)

// requireNetwork skips the test unless the network is reachable. Tests that pull images need
// it; gating lets the suite run offline without failing.
func requireNetwork(t *testing.T) {
	t.Helper()
	conn, err := net.DialTimeout("tcp", "github.com:443", 3*time.Second)
	if err != nil {
		t.Skipf("skipping: network unavailable: %v", err)
	}
	_ = conn.Close()
}

// requireDocker skips the test unless a reachable docker daemon is available.
// GetDockerClient succeeds even without a running daemon (its ping is best-effort),
// so the daemon has to be pinged explicitly here to decide whether to skip.
func requireDocker(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	cli, err := GetDockerClient(ctx)
	if err != nil {
		t.Skipf("skipping: docker client unavailable: %v", err)
	}
	defer cli.Close()
	if _, err := cli.Ping(ctx, mobyclient.PingOptions{}); err != nil {
		t.Skipf("skipping: docker daemon unreachable: %v", err)
	}
}
