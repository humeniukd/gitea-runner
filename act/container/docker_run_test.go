// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"gitea.com/gitea/runner/act/common"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	mobyclient "github.com/moby/moby/client"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDocker(t *testing.T) {
	requireDocker(t)
	ctx := context.Background()
	client, err := GetDockerClient(ctx)
	require.NoError(t, err)
	defer client.Close()

	dockerBuild := NewDockerBuildExecutor(NewDockerBuildExecutorInput{
		ContextDir: "testdata",
		ImageTag:   "envmergetest",
	})

	err = dockerBuild(ctx)
	assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act

	cr := &containerReference{
		cli: client,
		input: &NewContainerInput{
			Image: "envmergetest",
		},
	}
	env := map[string]string{
		"PATH":         "/usr/local/bin:/usr/bin:/usr/sbin:/bin:/sbin",
		"RANDOM_VAR":   "WITH_VALUE",
		"ANOTHER_VAR":  "",
		"CONFLICT_VAR": "I_EXIST_IN_MULTIPLE_PLACES",
	}

	envExecutor := cr.extractFromImageEnv(&env)
	err = envExecutor(ctx)
	assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act
	assert.Equal(t, map[string]string{
		"PATH":            "/usr/local/bin:/usr/bin:/usr/sbin:/bin:/sbin:/this/path/does/not/exists/anywhere:/this/either",
		"RANDOM_VAR":      "WITH_VALUE",
		"ANOTHER_VAR":     "",
		"SOME_RANDOM_VAR": "",
		"ANOTHER_ONE":     "BUT_I_HAVE_VALUE",
		"CONFLICT_VAR":    "I_EXIST_IN_MULTIPLE_PLACES",
	}, env)
}

type mockDockerClient struct {
	mobyclient.APIClient
	mock.Mock
}

func (m *mockDockerClient) ServerVersion(ctx context.Context, opts mobyclient.ServerVersionOptions) (mobyclient.ServerVersionResult, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(mobyclient.ServerVersionResult), args.Error(1)
}

func (m *mockDockerClient) ExecCreate(ctx context.Context, id string, opts mobyclient.ExecCreateOptions) (mobyclient.ExecCreateResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.ExecCreateResult), args.Error(1)
}

func (m *mockDockerClient) ExecAttach(ctx context.Context, id string, opts mobyclient.ExecAttachOptions) (mobyclient.ExecAttachResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.ExecAttachResult), args.Error(1)
}

func (m *mockDockerClient) ExecInspect(ctx context.Context, execID string, opts mobyclient.ExecInspectOptions) (mobyclient.ExecInspectResult, error) {
	args := m.Called(ctx, execID, opts)
	return args.Get(0).(mobyclient.ExecInspectResult), args.Error(1)
}

func (m *mockDockerClient) ContainerAttach(ctx context.Context, containerID string, opts mobyclient.ContainerAttachOptions) (mobyclient.ContainerAttachResult, error) {
	args := m.Called(ctx, containerID, opts)
	return args.Get(0).(mobyclient.ContainerAttachResult), args.Error(1)
}

func (m *mockDockerClient) ContainerCreate(ctx context.Context, opts mobyclient.ContainerCreateOptions) (mobyclient.ContainerCreateResult, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(mobyclient.ContainerCreateResult), args.Error(1)
}

func (m *mockDockerClient) ContainerStart(ctx context.Context, id string, opts mobyclient.ContainerStartOptions) (mobyclient.ContainerStartResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.ContainerStartResult), args.Error(1)
}

func (m *mockDockerClient) ContainerWait(ctx context.Context, containerID string, opts mobyclient.ContainerWaitOptions) mobyclient.ContainerWaitResult {
	args := m.Called(ctx, containerID, opts)
	return args.Get(0).(mobyclient.ContainerWaitResult)
}

func (m *mockDockerClient) CopyToContainer(ctx context.Context, id string, options mobyclient.CopyToContainerOptions) (mobyclient.CopyToContainerResult, error) {
	args := m.Called(ctx, id, options)
	return args.Get(0).(mobyclient.CopyToContainerResult), args.Error(1)
}

func (m *mockDockerClient) CopyFromContainer(ctx context.Context, id string, opts mobyclient.CopyFromContainerOptions) (mobyclient.CopyFromContainerResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.CopyFromContainerResult), args.Error(1)
}

func (m *mockDockerClient) ContainerInspect(ctx context.Context, id string, opts mobyclient.ContainerInspectOptions) (mobyclient.ContainerInspectResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.ContainerInspectResult), args.Error(1)
}

func (m *mockDockerClient) ContainerList(ctx context.Context, opts mobyclient.ContainerListOptions) (mobyclient.ContainerListResult, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(mobyclient.ContainerListResult), args.Error(1)
}

func (m *mockDockerClient) ContainerRemove(ctx context.Context, id string, opts mobyclient.ContainerRemoveOptions) (mobyclient.ContainerRemoveResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.ContainerRemoveResult), args.Error(1)
}

func (m *mockDockerClient) ContainerKill(ctx context.Context, id string, opts mobyclient.ContainerKillOptions) (mobyclient.ContainerKillResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.ContainerKillResult), args.Error(1)
}

func (m *mockDockerClient) NetworkList(ctx context.Context, opts mobyclient.NetworkListOptions) (mobyclient.NetworkListResult, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(mobyclient.NetworkListResult), args.Error(1)
}

func (m *mockDockerClient) NetworkInspect(ctx context.Context, id string, opts mobyclient.NetworkInspectOptions) (mobyclient.NetworkInspectResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.NetworkInspectResult), args.Error(1)
}

func (m *mockDockerClient) NetworkRemove(ctx context.Context, id string, opts mobyclient.NetworkRemoveOptions) (mobyclient.NetworkRemoveResult, error) {
	args := m.Called(ctx, id, opts)
	return args.Get(0).(mobyclient.NetworkRemoveResult), args.Error(1)
}

type interruptReader struct {
	started     chan struct{}
	interrupted chan struct{}
	stopped     chan struct{}
}

func (r *interruptReader) Read(_ []byte) (int, error) {
	close(r.started)
	<-r.interrupted
	close(r.stopped)
	return 0, io.EOF
}

type mockConn struct {
	net.Conn
	mock.Mock
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *mockConn) Close() (err error) {
	return nil
}

func TestRunnerArchCachesSuccessfulResultsByEndpoint(t *testing.T) {
	t.Setenv("DOCKER_TLS_VERIFY", "")
	t.Setenv("DOCKER_CERT_PATH", "")
	t.Setenv("DOCKER_API_VERSION", "1.55")
	for _, arch := range []string{"aarch64", "x86_64", "custom"} {
		t.Run(arch, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/_ping" {
					return
				}
				assert.Equal(t, "/v1.55/info", r.URL.Path)
				switch calls.Add(1) {
				case 1:
					w.WriteHeader(http.StatusInternalServerError)
				case 2:
					_, _ = io.WriteString(w, `{}`)
				default:
					_, _ = io.WriteString(w, `{"Architecture":"`+arch+`"}`)
				}
			}))
			defer server.Close()
			t.Setenv("DOCKER_HOST", server.URL)
			assert.Empty(t, RunnerArch(t.Context()))
			assert.Empty(t, RunnerArch(t.Context()))
			var workers sync.WaitGroup
			for range 8 {
				workers.Go(func() {
					assert.Equal(t, map[string]string{"aarch64": "ARM64", "x86_64": "X64", "custom": "custom"}[arch], RunnerArch(t.Context()))
				})
			}
			workers.Wait()
			assert.EqualValues(t, 3, calls.Load())
		})
	}
}

func TestDockerReadUIDGID(t *testing.T) {
	for _, tc := range []struct {
		name, user, stdout string
		skipExec           bool
		err                error
		uid, gid           int
	}{
		{name: "empty user", skipExec: true},
		{name: "numeric user and group", user: "1000:1001", skipExec: true, uid: 1000, gid: 1001},
		{name: "numeric root and group", user: "0:0", skipExec: true},
		{name: "named user ignores stderr", user: "runner", stdout: "1000\n1001\n", uid: 1000, gid: 1001},
		{name: "named root", user: "root", stdout: "0\n1001\n", gid: 1001},
		{name: "bare zero", user: "0", stdout: "0\n1001\n", gid: 1001},
		{name: "bare uid", user: "1000", stdout: "1000\n1001\n", uid: 1000, gid: 1001},
		{name: "named group", user: "1000:runner", stdout: "1000\n1001\n", uid: 1000, gid: 1001},
		{name: "truncated output preserves ids", user: "runner", stdout: "1000\n100", uid: 7, gid: 8},
		{name: "missing shell preserves ids", user: "runner", stdout: "sh: id: not found\n", uid: 7, gid: 8},
		{name: "overflow preserves ids", user: "runner", stdout: "2147483648\n1001\n", uid: 7, gid: 8},
		{name: "inspect and exec failure preserve ids", err: errors.New("daemon unreachable"), uid: 7, gid: 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockDockerClient{}
			client.On("ContainerInspect", t.Context(), "123", mobyclient.ContainerInspectOptions{}).Return(mobyclient.ContainerInspectResult{Container: container.InspectResponse{Config: &container.Config{User: tc.user}}}, tc.err).Once()
			if !tc.skipExec {
				client.On("ExecCreate", t.Context(), "123", mobyclient.ExecCreateOptions{Cmd: []string{"sh", "-c", "id -u; id -g"}, AttachStdout: true, AttachStderr: true}).
					Return(mobyclient.ExecCreateResult{ID: "id"}, tc.err).Once()
				if tc.err == nil {
					client.On("ExecAttach", t.Context(), "id", mobyclient.ExecAttachOptions{}).Return(mobyclient.ExecAttachResult{
						Conn: &mockConn{}, Reader: bufio.NewReader(bytes.NewReader(append(stdcopyFrame(stdcopy.Stderr, "9999\n9999\n"), stdcopyFrame(stdcopy.Stdout, tc.stdout)...))),
					}, nil).Once()
				}
			}
			cr := &containerReference{id: "123", cli: client, UID: 7, GID: 8}
			require.NoError(t, cr.tryReadUIDGID()(t.Context()))
			assert.Equal(t, [2]int{tc.uid, tc.gid}, [2]int{cr.UID, cr.GID})
			client.AssertExpectations(t)
		})
	}
}

func TestDockerExecAbort(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	reader := &interruptReader{started: make(chan struct{}), interrupted: make(chan struct{}), stopped: make(chan struct{})}
	conn := &mockConn{}
	conn.On("Write", []byte{3}).
		Run(func(mock.Arguments) { close(reader.interrupted) }).
		Return(1, nil)

	client := &mockDockerClient{}
	client.On("ExecCreate", ctx, "123", mock.AnythingOfType("client.ExecCreateOptions")).Return(mobyclient.ExecCreateResult{ID: "id"}, nil)
	client.On("ExecAttach", ctx, "id", mock.AnythingOfType("client.ExecAttachOptions")).Return(mobyclient.ExecAttachResult{
		Conn:   conn,
		Reader: bufio.NewReader(reader),
	}, nil)

	cr := &containerReference{
		id:  "123",
		cli: client,
		input: &NewContainerInput{
			Image: "image",
		},
	}

	channel := make(chan error)

	go func() {
		channel <- cr.exec([]string{""}, map[string]string{}, "user", "workdir")(ctx)
	}()

	<-reader.started
	cancel()

	err := <-channel
	<-reader.stopped
	assert.ErrorIs(t, err, context.Canceled) //nolint:testifylint // pre-existing issue from nektos/act

	conn.AssertExpectations(t)
	client.AssertExpectations(t)
}

func TestDockerExecFailure(t *testing.T) {
	notFound := cerrdefs.ErrNotFound.WithMessage("No such container: 123")
	for _, tc := range []struct {
		name                           string
		attachErr, inspectErr, wantErr error
	}{
		{name: "exit code", wantErr: ExitCodeError(1)},
		{name: "attach failure never recreates exec", attachErr: notFound, wantErr: notFound},
		{name: "inspect failure never replays exec", inspectErr: notFound, wantErr: notFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			client := &mockDockerClient{}
			client.On("ExecCreate", ctx, "123", mock.AnythingOfType("client.ExecCreateOptions")).Return(mobyclient.ExecCreateResult{ID: "id"}, nil).Once()
			client.On("ExecAttach", ctx, "id", mock.AnythingOfType("client.ExecAttachOptions")).Return(mobyclient.ExecAttachResult{Conn: &mockConn{}, Reader: bufio.NewReader(strings.NewReader(""))}, tc.attachErr).Once()
			if tc.attachErr == nil {
				client.On("ExecInspect", ctx, "id", mobyclient.ExecInspectOptions{}).Return(mobyclient.ExecInspectResult{ExitCode: 1}, tc.inspectErr).Once()
			}
			cr := &containerReference{id: "123", cli: client, input: &NewContainerInput{Name: "job-1"}}
			err := cr.Exec([]string{"echo"}, nil, "user", "workdir")(ctx)
			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == ExitCodeError(1) {
				assert.Equal(t, "Process completed with exit code 1.", err.Error())
			}
			assert.Equal(t, "123", cr.id)
			client.AssertExpectations(t)
		})
	}
}

func stdcopyFrame(stream stdcopy.StdType, payload string) []byte {
	b := make([]byte, 8+len(payload))
	b[0] = byte(stream)
	binary.BigEndian.PutUint32(b[4:8], uint32(len(payload)))
	copy(b[8:], payload)
	return b
}

func TestDockerStartWait(t *testing.T) {
	for _, tc := range []struct {
		name              string
		startErr, wantErr error
		exitCode          int64
	}{
		{name: "flushes trailing output before returning"},
		{name: "wait returns exit code", exitCode: 2, wantErr: ExitCodeError(2)},
		{name: "missing container after attach is never retried", startErr: cerrdefs.ErrNotFound.WithMessage("No such container: 123"), wantErr: ErrContainerNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var lines []string
			client := &mockDockerClient{}
			client.On("ContainerAttach", t.Context(), "123", mock.Anything).Return(mobyclient.ContainerAttachResult{
				Conn: &mockConn{}, Reader: bufio.NewReader(bytes.NewReader(stdcopyFrame(stdcopy.Stdout, "line one\nlast line without newline"))),
			}, nil).Once()
			client.On("ContainerStart", t.Context(), "123", mobyclient.ContainerStartOptions{}).Return(mobyclient.ContainerStartResult{}, tc.startErr).Once()
			if tc.startErr == nil {
				status := make(chan container.WaitResponse, 1)
				status <- container.WaitResponse{StatusCode: tc.exitCode}
				client.On("ContainerWait", t.Context(), "123", mobyclient.ContainerWaitOptions{Condition: container.WaitConditionNotRunning}).Return(mobyclient.ContainerWaitResult{Result: status}).Once()
				if tc.exitCode == 0 {
					client.On("ContainerInspect", t.Context(), "123", mobyclient.ContainerInspectOptions{}).Return(mobyclient.ContainerInspectResult{Container: container.InspectResponse{Config: &container.Config{}}}, nil).Once()
				}
			}
			cr := &containerReference{id: "123", cli: client, input: &NewContainerInput{Name: "job-1", Stdout: common.NewLineWriter(func(line string) bool {
				lines = append(lines, line)
				return true
			})}}
			err := cr.Start(true)(t.Context())
			require.ErrorIs(t, err, tc.wantErr)
			if tc.startErr != nil {
				assert.Empty(t, cr.id)
				assert.Equal(t, `container "job-1" does not exist; cannot start it`, err.Error())
				<-cr.attachDone
			} else {
				assert.Equal(t, []string{"line one\n", "last line without newline"}, lines)
			}
			client.AssertExpectations(t)
		})
	}
}

func TestRemoveHandlesAutoRemovalAndStaleIDs(t *testing.T) {
	removeFailure := errors.New("driver failed to remove root filesystem")
	for _, tc := range []struct {
		name, id, removeID             string
		err, killErr, waitErr, wantErr error
		wantWait, useName              bool
	}{
		{name: "removal in progress", id: "abc", removeID: "abc", err: cerrdefs.ErrConflict.WithMessage("removal of container abc is already in progress"), wantWait: true},
		{name: "wait canceled", id: "abc", removeID: "abc", err: cerrdefs.ErrConflict, wantWait: true, waitErr: context.Canceled, wantErr: context.Canceled},
		{name: "removed during wait", id: "abc", removeID: "abc", err: cerrdefs.ErrConflict, wantWait: true, waitErr: cerrdefs.ErrNotFound},
		{name: "already removed", id: "abc", removeID: "abc", err: cerrdefs.ErrNotFound.WithMessage("No such container: abc")},
		{name: "removed cleanly", id: "abc", removeID: "abc"},
		{name: "real failure", id: "abc", removeID: "abc", err: removeFailure, wantErr: removeFailure},
		{name: "unknown id removes by name", removeID: "job-1", useName: true},
		{name: "stale id retries by name", id: "stale", removeID: "stale", err: cerrdefs.ErrNotFound, killErr: cerrdefs.ErrNotFound, useName: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			ctx := common.WithLogger(t.Context(), logger)
			client := &mockDockerClient{}
			cr := &containerReference{id: tc.id, cli: client}
			if tc.id != "" {
				client.On("ContainerKill", ctx, tc.id, mobyclient.ContainerKillOptions{Signal: "SIGKILL"}).Return(mobyclient.ContainerKillResult{}, tc.killErr).Once()
			}
			client.On("ContainerRemove", ctx, tc.removeID, mobyclient.ContainerRemoveOptions{RemoveVolumes: true, Force: true}).Return(mobyclient.ContainerRemoveResult{}, tc.err).Once()
			if tc.useName {
				cr.input = &NewContainerInput{Name: "job-1"}
				if tc.id != "" {
					client.On("ContainerRemove", ctx, "job-1", mobyclient.ContainerRemoveOptions{RemoveVolumes: true, Force: true}).Return(mobyclient.ContainerRemoveResult{}, nil).Once()
				}
			}
			if tc.wantWait {
				removed, waitErrors := make(chan container.WaitResponse, 1), make(chan error, 1)
				if tc.waitErr == nil {
					removed <- container.WaitResponse{}
				} else {
					waitErrors <- tc.waitErr
				}
				client.On("ContainerWait", mock.Anything, tc.id, mobyclient.ContainerWaitOptions{Condition: container.WaitConditionRemoved}).Return(mobyclient.ContainerWaitResult{Result: removed, Error: waitErrors}).Once()
			}
			require.ErrorIs(t, cr.remove()(ctx), tc.wantErr)
			if tc.wantErr != nil {
				assert.Equal(t, tc.id, cr.id)
			} else {
				assert.Empty(t, cr.id)
			}
			assert.Empty(t, hook.AllEntries())
			client.AssertExpectations(t)
		})
	}
}

func TestCreateRevalidatesStaleID(t *testing.T) {
	notFound, transient := cerrdefs.ErrNotFound, errors.New("daemon unreachable")
	for _, tc := range []struct {
		name, id, wantID    string
		inspectErr, listErr error
	}{
		{name: "stale id recreated", id: "stale", inspectErr: notFound, wantID: "new"},
		{name: "stale id resolved by name", id: "stale", inspectErr: notFound, wantID: "fresh"},
		{name: "live id kept", id: "live", wantID: "live"},
		{name: "transient inspect error trusts cache", id: "live", inspectErr: transient, wantID: "live"},
		{name: "list error propagates", listErr: transient},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockDockerClient{}
			if tc.id != "" {
				client.On("ContainerInspect", t.Context(), tc.id, mobyclient.ContainerInspectOptions{}).Return(mobyclient.ContainerInspectResult{}, tc.inspectErr).Once()
			}
			if tc.id == "" || tc.inspectErr == notFound {
				var containers []container.Summary
				if tc.wantID == "fresh" {
					containers = []container.Summary{{ID: "other", Names: []string{"/somebody-else"}}, {ID: "fresh", Names: []string{"/job-1"}}}
				}
				client.On("ContainerList", t.Context(), mobyclient.ContainerListOptions{All: true}).Return(mobyclient.ContainerListResult{Items: containers}, tc.listErr).Once()
			}
			if tc.wantID == "new" {
				client.On("ContainerCreate", t.Context(), mock.Anything).Return(mobyclient.ContainerCreateResult{ID: "new"}, nil).Once()
			}
			cr := &containerReference{id: tc.id, cli: client, input: &NewContainerInput{Name: "job-1"}}
			require.ErrorIs(t, cr.Create(nil, nil)(t.Context()), tc.listErr)
			assert.Equal(t, tc.wantID, cr.id)
			client.AssertExpectations(t)
		})
	}
}

func TestRejectsUnknownAndRemovedContainers(t *testing.T) {
	ctx := context.Background()
	client := &mockDockerClient{}
	client.On("ContainerList", ctx, mobyclient.ContainerListOptions{All: true}).Return(mobyclient.ContainerListResult{}, nil)
	cr := &containerReference{cli: client, input: &NewContainerInput{Name: "job-1"}}
	check := func(op string, err error) {
		t.Helper()
		require.ErrorIs(t, err, ErrContainerNotFound, op)
		assert.Contains(t, err.Error(), `container "job-1" does not exist`, op)
	}
	check("copyContent", cr.copyContent("/var/run/act", &FileEntry{Name: "x", Mode: 0o644})(ctx))
	check("copyDir", cr.copyDir("/var/run/act", "/src", false, false)(ctx))
	check("exec", cr.exec([]string{"echo"}, nil, "", "")(ctx))
	_, err := cr.GetContainerArchive(ctx, "/var/run/act/x")
	check("GetContainerArchive", err)
	_, err = cr.Inspect(ctx)
	check("Inspect", err)
	client.On("ContainerInspect", ctx, "gone", mobyclient.ContainerInspectOptions{}).
		Return(mobyclient.ContainerInspectResult{}, cerrdefs.ErrNotFound)
	removed := &containerReference{id: "gone", cli: client, input: &NewContainerInput{Name: "job-1"}}
	_, err = removed.Inspect(ctx)
	check("Inspect after removal", err)
}

func TestPublicOperationsHandleStaleID(t *testing.T) {
	notFound := cerrdefs.ErrNotFound.WithMessage("Error response from daemon: No such container: stale")
	transient := errors.New("daemon unreachable")
	missingPath := cerrdefs.ErrNotFound.WithMessage("Could not find the file /No such container: x in container stale")
	for _, operation := range []struct {
		name, method string
		result       any
	}{
		{"copy", "CopyToContainer", mobyclient.CopyToContainerResult{}},
		{"copy directory", "CopyToContainer", mobyclient.CopyToContainerResult{}},
		{"get archive", "CopyFromContainer", mobyclient.CopyFromContainerResult{}},
		{"exec", "ExecCreate", mobyclient.ExecCreateResult{ID: "exec-id"}},
		{"start", "ContainerStart", mobyclient.ContainerStartResult{}},
		{"attach", "ContainerAttach", mobyclient.ContainerAttachResult{Conn: &mockConn{}, Reader: bufio.NewReader(strings.NewReader(""))}},
	} {
		for _, tc := range []struct {
			name, wantID                           string
			firstErr, retryErr, lookupErr, wantErr error
			resolve                                bool
		}{
			{name: "known id used directly", wantID: "stale"},
			{name: "stale id cleared", firstErr: notFound, wantErr: ErrContainerNotFound},
			{name: "replacement found by name", firstErr: notFound, resolve: true, wantID: "fresh"},
			{name: "replacement also removed, no second retry", firstErr: notFound, resolve: true, retryErr: notFound, wantErr: ErrContainerNotFound},
			{name: "transient error retains id", firstErr: transient, wantErr: transient, wantID: "stale"},
			{name: "lookup error propagates", firstErr: notFound, lookupErr: transient, wantErr: transient},
			{name: "missing path retains id", firstErr: missingPath, wantErr: missingPath, wantID: "stale"},
			{name: "unclassified 404 retains id", firstErr: cerrdefs.ErrNotFound, wantErr: cerrdefs.ErrNotFound, wantID: "stale"},
		} {
			t.Run(operation.name+"/"+tc.name, func(t *testing.T) {
				ctx := t.Context()
				client := &mockDockerClient{}
				var copies [][]byte
				expectOperation := func(id string, operationErr error) {
					options, result := any(mock.Anything), operation.result
					if operation.name == "get archive" {
						options, result = mobyclient.CopyFromContainerOptions{SourcePath: "/var/run/act/x"}, mobyclient.CopyFromContainerResult{Content: io.NopCloser(strings.NewReader("archive"))}
					}
					client.On(operation.method, ctx, id, options).Run(func(args mock.Arguments) {
						if operation.method == "CopyToContainer" {
							content, err := io.ReadAll(args.Get(2).(mobyclient.CopyToContainerOptions).Content)
							require.NoError(t, err)
							require.NotEmpty(t, content)
							copies = append(copies, content)
						}
					}).Return(result, operationErr).Once()
				}
				expectOperation("stale", tc.firstErr)
				if tc.firstErr == notFound {
					var containers []container.Summary
					if tc.resolve {
						containers = []container.Summary{{ID: "fresh", Names: []string{"/job-1"}}}
						expectOperation("fresh", tc.retryErr)
					}
					client.On("ContainerList", ctx, mobyclient.ContainerListOptions{All: true}).Return(mobyclient.ContainerListResult{Items: containers}, tc.lookupErr).Once()
				}
				if tc.wantErr == nil {
					switch operation.name {
					case "exec":
						client.On("ExecAttach", ctx, "exec-id", mobyclient.ExecAttachOptions{}).Return(mobyclient.ExecAttachResult{Conn: &mockConn{}, Reader: bufio.NewReader(strings.NewReader(""))}, nil).Once()
						client.On("ExecInspect", ctx, "exec-id", mobyclient.ExecInspectOptions{}).Return(mobyclient.ExecInspectResult{}, nil).Once()
					case "attach":
						status := make(chan container.WaitResponse, 1)
						status <- container.WaitResponse{}
						client.On("ContainerWait", ctx, tc.wantID, mobyclient.ContainerWaitOptions{Condition: container.WaitConditionNotRunning}).Return(mobyclient.ContainerWaitResult{Result: status}).Once()
						client.On("ContainerStart", ctx, tc.wantID, mobyclient.ContainerStartOptions{}).Return(mobyclient.ContainerStartResult{}, nil).Once()
						fallthrough
					case "start":
						client.On("ContainerInspect", ctx, tc.wantID, mobyclient.ContainerInspectOptions{}).Return(mobyclient.ContainerInspectResult{Container: container.InspectResponse{Config: &container.Config{}}}, nil).Once()
					}
				}
				cr := &containerReference{id: "stale", cli: client, input: &NewContainerInput{Name: "job-1"}}
				var err error
				switch operation.name {
				case "copy":
					err = cr.Copy("/var/run/act", &FileEntry{Name: "x", Mode: 0o644, Body: "content"})(ctx)
				case "copy directory":
					source := t.TempDir()
					require.NoError(t, os.WriteFile(filepath.Join(source, "x"), []byte("content"), 0o644))
					err = cr.CopyDir("/var/run/act", source, false, false)(ctx)
				case "get archive":
					var archive io.ReadCloser
					archive, err = cr.GetContainerArchive(ctx, "/var/run/act/x")
					if err == nil {
						content, readErr := io.ReadAll(archive)
						require.NoError(t, readErr)
						assert.Equal(t, "archive", string(content))
						require.NoError(t, archive.Close())
					}
				case "exec":
					err = cr.Exec([]string{"echo"}, nil, "", "")(ctx)
				case "start", "attach":
					err = cr.Start(operation.name == "attach")(ctx)
				}
				require.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, tc.wantID, cr.id)
				if len(copies) == 2 {
					assert.Equal(t, copies[0], copies[1])
				}
				client.AssertExpectations(t)
			})
		}
	}
}

// Type assert containerReference implements ExecutionsEnvironment
var _ ExecutionsEnvironment = &containerReference{}

func TestCheckVolumes(t *testing.T) {
	testCases := []struct {
		desc          string
		validVolumes  []string
		binds         []string
		expectedBinds []string
	}{
		{
			desc:         "match all volumes",
			validVolumes: []string{"**"},
			binds: []string{
				"shared_volume:/shared_volume",
				"/home/test/data:/test_data",
				"/etc/conf.d/base.json:/config/base.json",
				"sql_data:/sql_data",
				"/secrets/keys:/keys",
			},
			expectedBinds: []string{
				"shared_volume:/shared_volume",
				"/home/test/data:/test_data",
				"/etc/conf.d/base.json:/config/base.json",
				"sql_data:/sql_data",
				"/secrets/keys:/keys",
			},
		},
		{
			desc:         "no volumes can be matched",
			validVolumes: []string{},
			binds: []string{
				"shared_volume:/shared_volume",
				"/home/test/data:/test_data",
				"/etc/conf.d/base.json:/config/base.json",
				"sql_data:/sql_data",
				"/secrets/keys:/keys",
			},
			expectedBinds: []string{},
		},
		{
			desc: "only allowed volumes can be matched",
			validVolumes: []string{
				"shared_volume",
				"/home/test/data",
				"/etc/conf.d/*.json",
			},
			binds: []string{
				"shared_volume:/shared_volume",
				"/home/test/data:/test_data",
				"/etc/conf.d/base.json:/config/base.json",
				"sql_data:/sql_data",
				"/secrets/keys:/keys",
			},
			expectedBinds: []string{
				"shared_volume:/shared_volume",
				"/home/test/data:/test_data",
				"/etc/conf.d/base.json:/config/base.json",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			ctx := common.WithLogger(context.Background(), logger)
			cr := &containerReference{
				input: &NewContainerInput{
					ValidVolumes: tc.validVolumes,
				},
			}
			_, hostConf := cr.sanitizeConfig(ctx, &container.Config{}, &container.HostConfig{Binds: tc.binds})
			assert.Equal(t, tc.expectedBinds, hostConf.Binds)
			assert.Len(t, hook.AllEntries(), len(tc.binds)-len(tc.expectedBinds)) // every drop is warned about
		})
	}
}

// A volume driver decides for itself what it mounts, e.g. the local driver with device= binds
// any host path, which valid_volumes never gets to see.
func TestMergeContainerConfigsDropsVolumeDriversFromWorkflows(t *testing.T) {
	const escape = "--mount type=volume,src=job-escape,dst=/host,volume-driver=local,volume-opt=type=none,volume-opt=o=bind,volume-opt=device=/"

	hostConfig, _ := mergeOptions(t, "", escape+" --mount type=volume,src=job-plain,dst=/cache", false)
	require.Len(t, hostConfig.Mounts, 1)
	assert.Equal(t, "job-plain", hostConfig.Mounts[0].Source)

	// the same mount from the runner's own options is the administrator's to make
	hostConfig, _ = mergeOptions(t, escape, "", false)
	require.Len(t, hostConfig.Mounts, 1)
	assert.Equal(t, "job-escape", hostConfig.Mounts[0].Source)
}

// Both of these are read here, on the runner, so a workflow could read the runner's files
// and environment with them.
func TestMergeContainerConfigsKeepsTheRunnersFilesAndEnvToItself(t *testing.T) {
	hostFile := filepath.Join(t.TempDir(), "host.env")
	require.NoError(t, os.WriteFile(hostFile, []byte("STOLEN=from-the-host\n"), 0o600))
	t.Setenv("RUNNER_SECRET", "s3cr3t")

	for _, option := range []string{"--env-file " + hostFile, "--label-file " + hostFile} {
		logger, _ := test.NewNullLogger()
		cr := &containerReference{input: &NewContainerInput{NetworkMode: "bridge", WorkflowOptions: option}}

		_, _, err := cr.mergeContainerConfigs(common.WithLogger(context.Background(), logger), &container.Config{}, &container.HostConfig{})
		require.ErrorContains(t, err, "not allowed in a workflow")

		// the runner reading its own files is what those options are for
		cr = &containerReference{input: &NewContainerInput{NetworkMode: "bridge", RunnerOptions: option}}
		_, _, err = cr.mergeContainerConfigs(common.WithLogger(context.Background(), logger), &container.Config{}, &container.HostConfig{})
		require.NoError(t, err)
	}

	// a bare name is no longer resolved from the runner's environment, for either source
	logger, _ := test.NewNullLogger()
	cr := &containerReference{input: &NewContainerInput{
		NetworkMode:     "bridge",
		RunnerOptions:   "--env RUNNER_SECRET",
		WorkflowOptions: "--env RUNNER_SECRET --env GIVEN=value",
	}}

	config, _, err := cr.mergeContainerConfigs(common.WithLogger(context.Background(), logger), &container.Config{}, &container.HostConfig{})
	require.NoError(t, err)
	assert.Equal(t, []string{"RUNNER_SECRET", "RUNNER_SECRET", "GIVEN=value"}, config.Env)
}

func TestSanitizeOptionsHostConfig(t *testing.T) {
	logger, _ := test.NewNullLogger()

	// every field the sanitizer resets, so a reset dropped in a refactor fails here
	hostConfig := &container.HostConfig{
		PidMode:           "host",
		IpcMode:           "host",
		UTSMode:           "host",
		CgroupnsMode:      "host",
		UsernsMode:        "host",
		CapAdd:            []string{"ALL"},
		SecurityOpt:       []string{"seccomp=unconfined", "apparmor=unconfined"},
		VolumesFrom:       []string{"other"},
		Runtime:           "runc",
		Isolation:         "process",
		VolumeDriver:      "rogue",
		MaskedPaths:       []string{},
		ReadonlyPaths:     []string{},
		CgroupParent:      "/custom",
		Devices:           []container.DeviceMapping{{PathOnHost: "/dev/sda", PathInContainer: "/dev/sda", CgroupPermissions: "rwm"}},
		DeviceCgroupRules: []string{"a *:* rwm"},
		DeviceRequests:    []container.DeviceRequest{{Count: -1, Capabilities: [][]string{{"gpu"}}}},
		Sysctls:           map[string]string{"net.ipv4.ip_forward": "1"},
	}

	sanitizeOptionsHostConfig(logger, hostConfig, &container.HostConfig{})

	assert.Equal(t, &container.HostConfig{}, hostConfig)
}

// mergeOptions merges both option sources into a bare container, returning the result and its log.
func mergeOptions(t *testing.T, runnerOptions, workflowOptions string, privileged bool) (*container.HostConfig, *test.Hook) {
	t.Helper()
	logger, hook := test.NewNullLogger()
	cr := &containerReference{input: &NewContainerInput{
		RunnerOptions:   runnerOptions,
		WorkflowOptions: workflowOptions,
		NetworkMode:     "bridge",
		UsernsMode:      "private",
	}}

	_, hostConfig, err := cr.mergeContainerConfigs(common.WithLogger(context.Background(), logger), &container.Config{}, &container.HostConfig{
		Privileged:  privileged,
		UsernsMode:  container.UsernsMode("private"),
		NetworkMode: container.NetworkMode("bridge"),
	})
	require.NoError(t, err)
	return hostConfig, hook
}

func TestMergeContainerConfigsStripsDangerousOptionsWhenUnprivileged(t *testing.T) {
	// OS-independent options only, --device and --gpus need a linux/windows server OS
	const dangerousOptions = "--pid=host --ipc=host --uts=host --cgroupns=host " +
		"--userns=host --cap-add=ALL --security-opt seccomp=unconfined " +
		"--security-opt apparmor=unconfined --volumes-from other --isolation process " +
		"--runtime runc --cgroup-parent /custom --sysctl net.ipv4.ip_forward=1"

	// whatever the workflow adds, an unprivileged container comes out exactly as the runner's
	// own options alone describe it, field for field
	for _, runnerOptions := range []string{"--shm-size 1g", dangerousOptions, "--cap-add SYS_ADMIN --security-opt seccomp=unconfined"} {
		runnerOnly, _ := mergeOptions(t, runnerOptions, "", false)
		withWorkflow, _ := mergeOptions(t, runnerOptions, dangerousOptions, false)

		assert.Equal(t, runnerOnly, withWorkflow, "runner options: %q", runnerOptions)
	}

	// the same options from the runner reach the daemon, even --userns, which no workflow may set
	kept, _ := mergeOptions(t, dangerousOptions, "", false)
	assert.Equal(t, "host", string(kept.PidMode))
	assert.Equal(t, []string{"ALL"}, kept.CapAdd)
	assert.Equal(t, "runc", kept.Runtime)
	assert.Equal(t, "host", string(kept.UsernsMode))
	assert.False(t, kept.Privileged)

	// privileged is the administrator opting in, so the workflow's options are honored
	privileged, _ := mergeOptions(t, "", dangerousOptions, true)
	assert.Equal(t, "host", string(privileged.PidMode))
	assert.Equal(t, []string{"ALL"}, privileged.CapAdd)
	assert.Equal(t, []string{"seccomp=unconfined", "apparmor=unconfined"}, privileged.SecurityOpt)
}

func TestCheckVolumesRejectsEscapingHostPaths(t *testing.T) {
	logger, _ := test.NewNullLogger()
	ctx := common.WithLogger(context.Background(), logger)

	base := t.TempDir()
	allowed := filepath.Join(base, "allowed")
	denied := filepath.Join(base, "denied")
	require.NoError(t, os.MkdirAll(allowed, 0o700))
	require.NoError(t, os.MkdirAll(denied, 0o700))

	cr := &containerReference{
		input: &NewContainerInput{
			ValidVolumes: []string{filepath.Join(allowed, "**")},
		},
	}

	escapingPath := allowed + string(filepath.Separator) + ".." + string(filepath.Separator) + "denied"
	_, hostConf := cr.sanitizeConfig(ctx, &container.Config{}, &container.HostConfig{
		Binds: []string{escapingPath + ":/mnt"},
	})
	assert.Empty(t, hostConf.Binds)

	linkPath := filepath.Join(allowed, "link")
	if err := os.Symlink(denied, linkPath); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}
	_, hostConf = cr.sanitizeConfig(ctx, &container.Config{}, &container.HostConfig{
		Binds: []string{linkPath + ":/mnt"},
	})
	assert.Empty(t, hostConf.Binds)

	_, hostConf = cr.sanitizeConfig(ctx, &container.Config{}, &container.HostConfig{
		Binds: []string{filepath.Join(linkPath, "missing") + ":/mnt"},
	})
	assert.Empty(t, hostConf.Binds)
}

func TestContainerInfoFromInspect(t *testing.T) {
	t.Run("reports no healthcheck when the image declares none", func(t *testing.T) {
		info := containerInfoFromInspect(container.InspectResponse{
			ID:    "abc123",
			State: &container.State{Status: "running", Running: true},
		})

		assert.Equal(t, "abc123", info.ID)
		assert.Equal(t, "running", info.State)
		assert.Equal(t, HealthNone, info.Health)
		assert.Empty(t, info.Ports)
	})

	t.Run("reports the health status and the last probe output", func(t *testing.T) {
		info := containerInfoFromInspect(container.InspectResponse{
			State: &container.State{
				Status: "running",
				Health: &container.Health{
					Status: container.Unhealthy,
					Log: []*container.HealthcheckResult{
						{Output: "first\n"},
						{Output: "connection refused\n"},
					},
				},
			},
		})

		assert.Equal(t, HealthUnhealthy, info.Health)
		assert.Equal(t, "connection refused", info.HealthOutput)
	})

	t.Run("reports the published ports", func(t *testing.T) {
		info := containerInfoFromInspect(container.InspectResponse{
			State: &container.State{Status: "running"},
			NetworkSettings: &container.NetworkSettings{
				Ports: network.PortMap{
					network.MustParsePort("5432/tcp"): []network.PortBinding{{HostPort: "49153"}},
					network.MustParsePort("6379/tcp"): nil,
				},
			},
		})

		assert.Equal(t, map[string]string{"5432": "49153"}, info.Ports)
	})

	t.Run("tolerates a container without state", func(t *testing.T) {
		info := containerInfoFromInspect(container.InspectResponse{ID: "abc123"})

		assert.Equal(t, "abc123", info.ID)
		assert.Equal(t, HealthNone, info.Health)
	})

	t.Run("reports the mount sources by container path", func(t *testing.T) {
		info := containerInfoFromInspect(container.InspectResponse{
			Mounts: []container.MountPoint{
				{Type: mount.TypeVolume, Name: "job", Source: "/var/lib/docker/volumes/job/_data", Destination: "/workspace/owner/repo"},
				{Type: mount.TypeBind, Source: "/var/run/docker.sock", Destination: "/var/run/docker.sock"},
			},
		})

		assert.Equal(t, map[string]string{
			"/workspace/owner/repo": "/var/lib/docker/volumes/job/_data",
			"/var/run/docker.sock":  "/var/run/docker.sock",
		}, info.Mounts)
	})
}

func TestMergeContainerConfigsVolumesReplaceRunnerMounts(t *testing.T) {
	logger, _ := test.NewNullLogger()
	ctx := common.WithLogger(context.Background(), logger)
	cr := &containerReference{
		input: &NewContainerInput{
			NetworkMode:   "bridge",
			RunnerOptions: "--volume /host/tools:/opt/hostedtoolcache",
		},
	}

	_, hostConf, err := cr.mergeContainerConfigs(ctx, &container.Config{}, &container.HostConfig{
		Binds:  []string{"/var/run/docker.sock:/var/run/docker.sock"},
		Mounts: []mount.Mount{{Type: mount.TypeVolume, Source: "act-toolcache", Target: "/opt/hostedtoolcache"}},
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"/var/run/docker.sock:/var/run/docker.sock", "/host/tools:/opt/hostedtoolcache"}, hostConf.Binds)
	assert.Empty(t, hostConf.Mounts)
}

func TestMergeContainerConfigsWarnsOnlyAboutOptionsThatWereGiven(t *testing.T) {
	warnings := func(runnerOptions, workflowOptions string) int {
		_, hook := mergeOptions(t, runnerOptions, workflowOptions, false)
		return len(hook.AllEntries())
	}

	assert.Zero(t, warnings("--volume /host/tools:/opt/hostedtoolcache", ""))
	assert.Zero(t, warnings("", "--shm-size 1g"))
	assert.Equal(t, 1, warnings("--network host", ""))
}

func TestMergeContainerConfigsKeepsNetworkAliasesFromOptions(t *testing.T) {
	logger, _ := test.NewNullLogger()
	cr := &containerReference{input: &NewContainerInput{
		NetworkMode: "job-network", NetworkAliases: []string{"redis"}, WorkflowOptions: "--network-alias redis-primary",
	}}

	_, _, err := cr.mergeContainerConfigs(common.WithLogger(t.Context(), logger), &container.Config{}, &container.HostConfig{})
	require.NoError(t, err)
	assert.Equal(t, []string{"redis", "redis-primary"}, cr.input.NetworkAliases)
}

// A dead daemon must fail the job, not panic through logrus and not silently
// drop the requested platform.
func TestSupportsContainerImagePlatformDaemonError(t *testing.T) {
	cli := &mockDockerClient{}
	cli.On("ServerVersion", mock.Anything, mock.Anything).
		Return(mobyclient.ServerVersionResult{}, errors.New("cannot connect to the Docker daemon"))

	supported, err := supportsContainerImagePlatform(t.Context(), cli)
	require.ErrorContains(t, err, "cannot connect to the Docker daemon")
	assert.False(t, supported)
}
