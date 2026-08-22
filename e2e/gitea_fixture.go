// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"

	"gitea.com/gitea/runner/act/container"

	apicontainer "github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	mobyclient "github.com/moby/moby/client"
)

const (
	giteaAdminUser = "e2e-admin"
	giteaAdminMail = "e2e-admin@example.com"
)

// GiteaFixture is a running gitea/gitea server, started fresh for one test binary run.
type GiteaFixture struct {
	cli        mobyclient.APIClient
	id         string
	baseURL    string
	adminToken string
}

// dockerClient connects to the ambient docker daemon, or returns an error the caller turns
// into a skip. Shared by TestMain (to decide whether to start the fixture at all) and by
// StartGitea itself.
func dockerClient(ctx context.Context) (mobyclient.APIClient, error) {
	cli, err := container.GetDockerClient(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := cli.Ping(ctx, mobyclient.PingOptions{}); err != nil {
		return nil, fmt.Errorf("docker daemon unreachable: %w", err)
	}
	return cli, nil
}

// hostAddress returns an address for the fixture that resolves both on this host and from inside
// a job container. It cannot be 127.0.0.1: the runner hands its own instance address to jobs as
// ACTIONS_RUNTIME_URL (see internal/app/run/runner.go), where loopback means the job container
// itself, so artifact upload would get ECONNREFUSED.
//
// The docker bridge gateway is preferred because it is the same 172.17.0.1 on every Linux host,
// so no machine-identifying IP reaches the logs. Whether the host can actually use it is decided
// by a bind probe rather than by GOOS: with Docker Desktop the daemon runs inside a VM, so the
// gateway exists only there and binding it fails here. The routable LAN IP is the fallback that
// works in both places, and loopback the last resort when offline (callback-free scenarios only).
func hostAddress(ctx context.Context, cli mobyclient.APIClient) netip.Addr {
	if gateway, ok := bridgeGateway(ctx, cli); ok && bindable(gateway) {
		return gateway
	}

	loopback := netip.AddrFrom4([4]byte{127, 0, 0, 1})

	conn, err := net.Dial("udp", "192.0.2.1:80") // TEST-NET-1, only consults the routing table
	if err != nil {
		return loopback
	}
	defer conn.Close()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return loopback
	}
	return addr.AddrPort().Addr().Unmap()
}

// bridgeGateway reports the docker bridge network's gateway, which containers reach the host at.
func bridgeGateway(ctx context.Context, cli mobyclient.APIClient) (netip.Addr, bool) {
	inspected, err := cli.NetworkInspect(ctx, "bridge", mobyclient.NetworkInspectOptions{})
	if err != nil {
		return netip.Addr{}, false
	}
	for _, cfg := range inspected.Network.IPAM.Config {
		if cfg.Gateway.Is4() {
			return cfg.Gateway, true
		}
	}
	return netip.Addr{}, false
}

// bindable reports whether this host owns addr, i.e. whether the fixture can publish a port on it.
func bindable(addr netip.Addr) bool {
	l, err := net.Listen("tcp", net.JoinHostPort(addr.String(), "0"))
	if err != nil {
		return false
	}
	_ = l.Close()
	return true
}

// freeHostPort finds a free port on host so the container's published port and
// GITEA__server__ROOT_URL can agree on it before the container starts, avoiding a
// chicken-and-egg wait on docker's own port-0 assignment.
func freeHostPort(host string) (int, error) {
	l, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
	if err != nil {
		return 0, err
	}
	defer l.Close()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("unexpected listener address type %T", l.Addr())
	}
	return addr.Port, nil
}

// pullImage fetches image, tolerating a pull failure when a local copy already exists. This
// mirrors the runner's own force_pull rule: an unreachable registry should not fail a run whose
// image is already on the host, which is what keeps a CI pre-pull useful against a flaky registry.
func pullImage(ctx context.Context, cli mobyclient.APIClient, image string) error {
	pull, err := cli.ImagePull(ctx, image, mobyclient.ImagePullOptions{})
	if err == nil {
		err = pull.Wait(ctx)
	}
	if err == nil {
		return nil
	}
	if _, inspectErr := cli.ImageInspect(ctx, image); inspectErr == nil {
		fmt.Fprintf(os.Stderr, "pull %s failed (%v), using the local copy\n", image, err)
		return nil
	}
	return fmt.Errorf("pull %s: %w", image, err)
}

// StartGitea starts a gitea/gitea container (image from GITEA_E2E_IMAGE, default main-nightly),
// waits for it to become healthy, and bootstraps an admin user plus API token.
//
// It takes no *testing.T because TestMain (which owns the fixture's lifetime, starting it once
// for the whole binary run) has no T to give it; callers in test functions surface a returned
// error via t.Fatal themselves.
func StartGitea(ctx context.Context) (*GiteaFixture, error) {
	cli, err := dockerClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("docker unavailable: %w", err)
	}

	image := os.Getenv("GITEA_E2E_IMAGE")
	if image == "" {
		image = "gitea/gitea:main-nightly"
	}

	host := hostAddress(ctx, cli)
	port, err := freeHostPort(host.String())
	if err != nil {
		return nil, fmt.Errorf("find a free host port: %w", err)
	}
	baseURL := fmt.Sprintf("http://%s:%d", host, port)

	if err := pullImage(ctx, cli, image); err != nil {
		return nil, err
	}

	containerPort := network.MustParsePort("3000/tcp")
	resp, err := cli.ContainerCreate(ctx, mobyclient.ContainerCreateOptions{
		Config: &apicontainer.Config{
			Image: image,
			Env: []string{
				"GITEA__security__INSTALL_LOCK=true",
				"GITEA__database__DB_TYPE=sqlite3",
				"GITEA__actions__ENABLED=true",
				"GITEA__server__ROOT_URL=" + baseURL + "/",
			},
			ExposedPorts: network.PortSet{containerPort: struct{}{}},
		},
		HostConfig: &apicontainer.HostConfig{
			PortBindings: network.PortMap{
				containerPort: []network.PortBinding{{HostIP: host, HostPort: strconv.Itoa(port)}},
			},
			AutoRemove: false,
		},
		Name: fmt.Sprintf("gitea-runner-e2e-%d", time.Now().UnixNano()),
	})
	if err != nil {
		return nil, fmt.Errorf("create gitea container: %w", err)
	}

	f := &GiteaFixture{cli: cli, id: resp.ID, baseURL: baseURL}

	if _, err := cli.ContainerStart(ctx, f.id, mobyclient.ContainerStartOptions{}); err != nil {
		_ = f.Close(ctx)
		return nil, fmt.Errorf("start gitea container: %w", err)
	}

	if err := f.waitHealthy(ctx); err != nil {
		_ = f.Close(ctx)
		return nil, fmt.Errorf("gitea did not become healthy: %w", err)
	}

	if err := f.bootstrapAdmin(ctx); err != nil {
		_ = f.Close(ctx)
		return nil, fmt.Errorf("bootstrap gitea admin: %w", err)
	}

	return f, nil
}

// waitHealthy polls /api/healthz until it returns 200, bounded by ctx.
func (f *GiteaFixture) waitHealthy(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	url := f.baseURL + "/api/healthz"
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err == nil {
			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				_ = resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return nil
				}
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

// bootstrapAdmin creates the admin user and mints its API token via `gitea` CLI execs, so no
// REST login/CSRF flow is needed for setup itself.
func (f *GiteaFixture) bootstrapAdmin(ctx context.Context) error {
	password := randomToken(16)
	if _, err := f.exec(ctx, []string{
		"gitea", "admin", "user", "create",
		"--username", giteaAdminUser,
		"--password", password,
		"--email", giteaAdminMail,
		"--admin",
		"--must-change-password=false",
	}); err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	out, err := f.exec(ctx, []string{
		"gitea", "admin", "user", "generate-access-token",
		"--username", giteaAdminUser,
		"--scopes", "all",
		"-t", "e2e-admin-token",
	})
	if err != nil {
		return fmt.Errorf("generate admin token: %w", err)
	}
	token := extractToken(out)
	if token == "" {
		return fmt.Errorf("could not parse access token from CLI output: %q", out)
	}
	f.adminToken = token
	return nil
}

// extractToken pulls the token value out of `gitea admin user generate-access-token`'s output,
// which prints a sentence ending in the token (format has varied across Gitea releases, so this
// takes the last whitespace-separated field rather than matching an exact prefix).
func extractToken(cliOutput string) string {
	fields := strings.Fields(cliOutput)
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}

func randomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// exec runs cmd inside the gitea container and returns combined stdout+stderr.
func (f *GiteaFixture) exec(ctx context.Context, cmd []string) (string, error) {
	created, err := f.cli.ExecCreate(ctx, f.id, mobyclient.ExecCreateOptions{
		Cmd:          cmd,
		User:         "git", // gitea refuses to run as root, the image's own user
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return "", err
	}

	attached, err := f.cli.ExecAttach(ctx, created.ID, mobyclient.ExecAttachOptions{})
	if err != nil {
		return "", err
	}
	defer attached.Close()

	var out bytes.Buffer
	if _, err := io.Copy(&out, attached.Reader); err != nil {
		return "", err
	}

	inspected, err := f.cli.ExecInspect(ctx, created.ID, mobyclient.ExecInspectOptions{})
	if err != nil {
		return "", err
	}
	if inspected.ExitCode != 0 {
		return "", fmt.Errorf("exec %v exited %d: %s", cmd, inspected.ExitCode, out.String())
	}
	return out.String(), nil
}

// RegistrationToken mints a fresh global runner registration token via Gitea's admin API.
func (f *GiteaFixture) RegistrationToken(ctx context.Context) (string, error) {
	var body struct {
		Token string `json:"token"`
	}
	if err := f.doJSON(ctx, http.MethodPost, "/api/v1/admin/actions/runners/registration-token", nil, &body); err != nil {
		return "", err
	}
	return body.Token, nil
}

func (f *GiteaFixture) doJSON(ctx context.Context, method, path string, reqBody, respBody any) error {
	api := &GiteaAPI{baseURL: f.baseURL, token: f.adminToken}
	return api.doJSON(ctx, method, path, reqBody, respBody)
}

// BaseURL is the address the runner registers against.
func (f *GiteaFixture) BaseURL() string { return f.baseURL }

// API returns a client authenticated as the bootstrap admin, for driving repos/workflows/runs.
func (f *GiteaFixture) API() *GiteaAPI {
	return &GiteaAPI{baseURL: f.baseURL, token: f.adminToken}
}

// Close removes the gitea container.
func (f *GiteaFixture) Close(ctx context.Context) error {
	if f.id == "" {
		return nil
	}
	_, err := f.cli.ContainerRemove(ctx, f.id, mobyclient.ContainerRemoveOptions{Force: true})
	return err
}
