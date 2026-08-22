// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gitea.com/gitea/runner/internal/app/poll"
	"gitea.com/gitea/runner/internal/app/run"
	"gitea.com/gitea/runner/internal/pkg/client"
	"gitea.com/gitea/runner/internal/pkg/config"
	"gitea.com/gitea/runner/internal/pkg/labels"

	"connectrpc.com/connect"
	pingv1 "gitea.dev/actionslib/ping/v1"
	runnerv1 "gitea.dev/actionslib/runner/v1"
)

// Harness is one in-process runner, registered and polling against a fixture Gitea. It
// reimplements the register+daemon call sequence from internal/app/cmd directly, skipping the
// disk/file-lock concerns that only matter for the CLI.
type Harness struct {
	runner *run.Runner
	poller *poll.Poller

	capabilities     string
	ephemeralWanted  bool
	ephemeralGranted bool
}

// harnessOptions collects the settings NewHarness applies before registering.
type harnessOptions struct {
	ephemeral bool
	cacheV2   *bool
}

// HarnessOption customises a Harness, so scenarios can vary registration and config without
// every existing call site having to pass defaults.
type HarnessOption func(*harnessOptions)

// WithEphemeral registers the runner as ephemeral. Servers too old to support it silently
// downgrade, which Harness.EphemeralGranted reports.
func WithEphemeral() HarnessOption {
	return func(o *harnessOptions) { o.ephemeral = true }
}

// WithCacheV2 forces the cache service version advertised to jobs (config cache.v2).
func WithCacheV2(enabled bool) HarnessOption {
	return func(o *harnessOptions) { o.cacheV2 = &enabled }
}

// NewHarness registers a runner named t.Name() against addr using token, then starts polling
// for tasks in the background.
func NewHarness(ctx context.Context, t *testing.T, addr, token string, lbls []string, opts ...HarnessOption) *Harness {
	t.Helper()

	var options harnessOptions
	for _, opt := range opts {
		opt(&options)
	}

	cfg, err := config.LoadDefault("")
	if err != nil {
		t.Fatalf("load default config: %v", err)
	}
	if options.cacheV2 != nil {
		cfg.Cache.V2 = options.cacheV2
	}
	cfg.Container.DockerHost = os.Getenv("DOCKER_HOST")
	if cfg.Container.DockerHost == "" {
		cfg.Container.DockerHost = "unix:///var/run/docker.sock"
	}
	cfg.Runner.File = t.TempDir() + "/.runner" // never written to; only config.LoadDefault's file-based defaults apply
	cfg.Runner.Insecure = true
	// Test-only tuning: every scenario registers a fresh runner and then waits for it to pick the
	// task up, so the 2s/5s production poll defaults are mostly idle sleep here. Capacity lets the
	// multi-job scenarios (matrix, artifacts, cache, outputs) run their jobs concurrently.
	cfg.Runner.FetchInterval = 250 * time.Millisecond
	cfg.Runner.FetchIntervalMax = 250 * time.Millisecond
	cfg.Runner.Capacity = 4
	// In container mode the fixture is only reachable by name on the network we share with it, so
	// job containers have to join that network too or they cannot resolve ACTIONS_RUNTIME_URL.
	if fixture != nil && fixture.Network() != "" {
		cfg.Container.Network = fixture.Network()
	}

	parsed := make(labels.Labels, 0, len(lbls))
	for _, l := range lbls {
		label, err := labels.Parse(l)
		if err != nil {
			t.Fatalf("parse label %q: %v", l, err)
		}
		parsed = append(parsed, label)
	}

	pingCli := client.New(addr, cfg.Runner.Insecure, "", "", config.RequestTimeout)
	if _, err := pingCli.Ping(ctx, connect.NewRequest(&pingv1.PingRequest{Data: t.Name()})); err != nil {
		t.Fatalf("ping %s: %v", addr, err)
	}

	regResp, err := pingCli.Register(ctx, connect.NewRequest(&runnerv1.RegisterRequest{
		Name:         t.Name(),
		Token:        token,
		Version:      "e2e",
		Labels:       parsed.Names(),
		Ephemeral:    options.ephemeral,
		Capabilities: run.RunnerCapabilities(),
	}))
	if err != nil {
		t.Fatalf("register runner: %v", err)
	}

	reg := &config.Registration{
		ID:        regResp.Msg.Runner.Id,
		UUID:      regResp.Msg.Runner.Uuid,
		Name:      regResp.Msg.Runner.Name,
		Token:     regResp.Msg.Runner.Token,
		Address:   addr,
		Labels:    lbls,
		Ephemeral: regResp.Msg.Runner.Ephemeral,
	}

	cli := client.New(addr, cfg.Runner.Insecure, reg.UUID, reg.Token, config.RequestTimeout)

	runner := run.NewRunner(cfg, reg, cli)
	declResp, err := runner.Declare(ctx, parsed.Names())
	if err != nil {
		_ = runner.Close()
		t.Fatalf("declare runner: %v", err)
	}
	runner.SetCapabilitiesFromDeclare(declResp)

	poller := poll.New(cfg, cli, runner)
	go poller.Poll()

	h := &Harness{
		runner: runner,
		poller: poller,
		// Negotiated via response header rather than the proto, see run.Runner.SetCapabilitiesFromDeclare.
		capabilities:     strings.TrimSpace(declResp.Header().Get("X-Gitea-Actions-Capabilities")),
		ephemeralWanted:  options.ephemeral,
		ephemeralGranted: regResp.Msg.Runner.Ephemeral,
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := h.Close(ctx); err != nil {
			t.Logf("harness shutdown: %v", err)
		}
	})
	return h
}

// Capabilities is the capability set the server negotiated back on Declare. It legitimately
// differs by Gitea version, which is why scenarios record it rather than assert one value.
func (h *Harness) Capabilities() string { return h.capabilities }

// EphemeralGranted reports whether the server honoured an ephemeral registration. Older servers
// accept the request but downgrade it, see internal/app/cmd/register.go.
func (h *Harness) EphemeralGranted() bool { return h.ephemeralGranted }

// EphemeralDowngraded reports an ephemeral registration the server silently refused.
func (h *Harness) EphemeralDowngraded() bool { return h.ephemeralWanted && !h.ephemeralGranted }

// Close stops polling and releases runner resources.
func (h *Harness) Close(ctx context.Context) error {
	if err := h.poller.Shutdown(ctx); err != nil {
		return fmt.Errorf("poller shutdown: %w", err)
	}
	return h.runner.Close()
}
