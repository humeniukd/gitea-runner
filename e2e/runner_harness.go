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

// Harness is an in-process runner registered against the fixture (register+poll without CLI disk lock).
type Harness struct {
	runner *run.Runner
	poller *poll.Poller

	capabilities     string
	ephemeralWanted  bool
	ephemeralGranted bool
}

type harnessOptions struct {
	ephemeral bool
	cacheV2   *bool
}

type HarnessOption func(*harnessOptions)

func WithEphemeral() HarnessOption {
	return func(o *harnessOptions) { o.ephemeral = true }
}

func WithCacheV2(enabled bool) HarnessOption {
	return func(o *harnessOptions) { o.cacheV2 = &enabled }
}

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
	cfg.Runner.File = t.TempDir() + "/.runner"
	cfg.Runner.Insecure = true
	cfg.Runner.FetchInterval = 250 * time.Millisecond // faster than prod defaults for local fixture
	cfg.Runner.FetchIntervalMax = 250 * time.Millisecond
	cfg.Runner.StateReportInterval = 500 * time.Millisecond // so cancel reaches the job quickly
	cfg.Runner.LogReportInterval = 500 * time.Millisecond
	cfg.Runner.Capacity = 4
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
		runner:           runner,
		poller:           poller,
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

func (h *Harness) Capabilities() string { return h.capabilities }

func (h *Harness) EphemeralGranted() bool { return h.ephemeralGranted }

func (h *Harness) EphemeralDowngraded() bool { return h.ephemeralWanted && !h.ephemeralGranted }

func (h *Harness) Close(ctx context.Context) error {
	if err := h.poller.Shutdown(ctx); err != nil {
		return fmt.Errorf("poller shutdown: %w", err)
	}
	return h.runner.Close()
}
