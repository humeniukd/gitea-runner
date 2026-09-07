// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !(WITHOUT_DOCKER || !(linux || darwin || windows || netbsd))

package container

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"gitea.com/gitea/runner/act/common"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
)

const (
	jobLabel      = "com.gitea.runner.job"
	maxCreateBody = 8 << 20

	probeDir   = "probe.d"
	probeFile  = "probe"
	probeMount = "/gitea-runner-probe"
	probePath  = probeMount + "/" + probeFile
)

var (
	createPath    = regexp.MustCompile(`^(/v[0-9.]+)?/(containers|networks|volumes)/create$`)
	rawStreamPath = regexp.MustCompile(`^(/v[0-9.]+)?/(containers/[^/]+/attach|exec/[^/]+/start)$`)
	proxyProbe    struct {
		sync.Mutex
		decided bool
		dir     string
	}
)

// DockerProxyDir returns where job proxy sockets live, "" while undecided or when the daemon cannot open the runner's files.
func DockerProxyDir(ctx context.Context) string {
	proxyProbe.Lock()
	defer proxyProbe.Unlock()
	if proxyProbe.decided {
		return proxyProbe.dir
	}
	dir := filepath.Join(os.TempDir(), "gitea-runner-docker")
	cli, err := GetDockerClient(ctx)
	if err != nil {
		common.Logger(ctx).Debugf("docker proxy probe postponed: %v", err)
		return ""
	}
	defer cli.Close()
	ok, err := daemonSeesDir(ctx, cli, dir)
	if err != nil {
		common.Logger(ctx).Debugf("docker proxy probe postponed: %v", err)
		return ""
	}
	proxyProbe.decided = true
	if ok {
		proxyProbe.dir = dir
	} else {
		common.Logger(ctx).Infof("the docker daemon cannot reach the runner's filesystem, jobs get the daemon socket directly")
	}
	return proxyProbe.dir
}

// daemonSeesDir reports whether the daemon opens the files the runner writes in dir,
// which is what a job's proxy socket mounted from there needs.
func daemonSeesDir(ctx context.Context, cli client.APIClient, dir string) (bool, error) {
	source := filepath.Join(dir, probeDir)
	if err := os.MkdirAll(source, 0o700); err != nil {
		return false, err
	}
	want := rand.Text() // fresh per probe
	if err := os.WriteFile(filepath.Join(source, probeFile), []byte(want), 0o600); err != nil {
		return false, err
	}
	images, err := cli.ImageList(ctx, client.ImageListOptions{})
	if err != nil {
		return false, err
	}
	if len(images.Items) == 0 {
		return false, errors.New("no image to probe with yet")
	}
	// nothing is started, the daemon reads the file out of the mount for us
	created, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{Image: images.Items[0].ID, Cmd: []string{"true"}},
		HostConfig: &container.HostConfig{Mounts: []mount.Mount{
			{Type: mount.TypeBind, Source: source, Target: probeMount, ReadOnly: true},
		}},
	})
	if err != nil {
		return false, err
	}
	defer func() {
		if _, err := cli.ContainerRemove(ctx, created.ID, client.ContainerRemoveOptions{Force: true}); err != nil {
			common.Logger(ctx).Debugf("removing the docker proxy probe container failed: %v", err)
		}
	}()
	// a daemon that cannot reach dir creates an empty one of its own instead of refusing the mount,
	// so only the content coming back says whose file the job would get
	res, err := cli.CopyFromContainer(ctx, created.ID, client.CopyFromContainerOptions{SourcePath: probePath})
	if cerrdefs.IsNotFound(err) {
		return false, nil // the daemon mounted an empty probe.d it created itself, not the runner's
	}
	if err != nil {
		return false, err
	}
	defer res.Content.Close()
	archive := tar.NewReader(res.Content)
	header, err := archive.Next()
	if errors.Is(err, io.EOF) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if header.Typeflag != tar.TypeReg {
		return false, nil // something other than the file the runner wrote is at that path
	}
	got, err := io.ReadAll(io.LimitReader(archive, int64(len(want))+1))
	if err != nil {
		return false, err
	}
	return string(got) == want, nil
}

// StartDockerProxy serves a job's docker socket in dir, labelling what the job creates through it.
func StartDockerProxy(daemonSocket, dir, job string) (*DockerProxy, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	digest := sha256.Sum256([]byte(job))
	socket := filepath.Join(dir, hex.EncodeToString(digest[:8])+".sock")
	_ = os.Remove(socket)
	listener, err := net.Listen("unix", socket)
	if err != nil {
		return nil, err
	}
	if info, err := os.Stat(daemonSocket); err == nil {
		_ = os.Chmod(socket, info.Mode().Perm())
	}
	dial := func(ctx context.Context, _, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "unix", daemonSocket)
	}
	transport := &http.Transport{DialContext: dial}
	forward := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.URL.Scheme = "http"
			r.Out.URL.Host = "docker"
		},
		Transport: transport,
	}
	server := &http.Server{ReadHeaderTimeout: 30 * time.Second, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method != http.MethodPost:
		case createPath.MatchString(r.URL.Path):
			if err := addLabel(r, job); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case rawStreamPath.MatchString(r.URL.Path):
			tunnel(w, r, dial)
			return
		}
		forward.ServeHTTP(w, r)
	})}
	go func() { _ = server.Serve(listener) }()
	return &DockerProxy{Socket: socket, close: func(ctx context.Context) error {
		err := removeJobResources(ctx, job)
		_ = server.Close()
		transport.CloseIdleConnections()
		_ = os.Remove(socket)
		return err
	}}, nil
}

func addLabel(r *http.Request, job string) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxCreateBody+1))
	if err != nil {
		return err
	}
	if len(body) > maxCreateBody {
		return errors.New("create request too large")
	}
	if len(bytes.TrimSpace(body)) == 0 {
		r.Body = io.NopCloser(bytes.NewReader(body))
		return nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return fmt.Errorf("invalid create request: %w", err)
	}
	key := "Labels"
	for name := range fields {
		if strings.EqualFold(name, key) {
			key = name
			break
		}
	}
	labels := map[string]string{}
	if raw := fields[key]; len(raw) > 0 && string(raw) != "null" {
		if err := json.Unmarshal(raw, &labels); err != nil {
			return fmt.Errorf("invalid create request: %w", err)
		}
	}
	labels[jobLabel] = job
	if fields[key], err = json.Marshal(labels); err != nil {
		return err
	}
	if body, err = json.Marshal(fields); err != nil {
		return err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))
	r.TransferEncoding = nil
	return nil
}

// tunnel splices attach and exec streams, which the daemon hijacks with or without an HTTP upgrade
func tunnel(w http.ResponseWriter, r *http.Request, dial func(context.Context, string, string) (net.Conn, error)) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "connection cannot be hijacked", http.StatusInternalServerError)
		return
	}
	upstream, err := dial(r.Context(), "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer upstream.Close()
	if err := r.Write(upstream); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	downstream, buffered, err := hijacker.Hijack()
	if err != nil {
		return
	}
	defer downstream.Close()
	if _, err := io.CopyN(upstream, buffered, int64(buffered.Reader.Buffered())); err != nil {
		return
	}
	done := make(chan struct{}, 2)
	go func() {
		_, _ = io.Copy(upstream, downstream)
		done <- struct{}{}
	}()
	go func() {
		_, _ = io.Copy(downstream, upstream)
		done <- struct{}{}
	}()
	<-done
}

func removeJobResources(ctx context.Context, job string) error {
	cli, err := GetDockerClient(ctx)
	if err != nil {
		return err
	}
	defer cli.Close()
	return removeLabelled(ctx, cli, job)
}

func removeLabelled(ctx context.Context, cli client.APIClient, job string) error {
	logger := common.Logger(ctx)
	filters := make(client.Filters).Add("label", jobLabel+"="+job)
	containers, err := cli.ContainerList(ctx, client.ContainerListOptions{All: true, Filters: filters})
	if err != nil {
		return err
	}
	var errs []error
	for _, c := range containers.Items {
		logger.Infof("removing container %s the job left behind", strings.TrimPrefix(strings.Join(c.Names, ","), "/"))
		errs = append(errs, (&containerReference{cli: cli, id: c.ID}).remove()(ctx))
	}
	networks, err := cli.NetworkList(ctx, client.NetworkListOptions{Filters: filters})
	if err != nil {
		return errors.Join(append(errs, err)...)
	}
	for _, n := range networks.Items {
		if _, err := cli.NetworkRemove(ctx, n.ID, client.NetworkRemoveOptions{}); err != nil && !cerrdefs.IsNotFound(err) {
			errs = append(errs, fmt.Errorf("failed to remove network %s: %w", n.Name, err))
		}
	}
	volumes, err := cli.VolumeList(ctx, client.VolumeListOptions{Filters: filters})
	if err != nil {
		return errors.Join(append(errs, err)...)
	}
	for _, v := range volumes.Items {
		if _, err := cli.VolumeRemove(ctx, v.Name, client.VolumeRemoveOptions{}); err != nil && !cerrdefs.IsNotFound(err) {
			errs = append(errs, fmt.Errorf("failed to remove volume %s: %w", v.Name, err))
		}
	}
	return errors.Join(errs...)
}
