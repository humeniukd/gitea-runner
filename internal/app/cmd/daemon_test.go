// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitea.com/gitea/runner/internal/pkg/config"
	"gitea.com/gitea/runner/internal/pkg/report"

	"connectrpc.com/connect"
	runnerv1 "gitea.dev/actionslib/runner/v1"
	"gitea.dev/actionslib/runner/v1/runnerv1connect"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestResolveLabels(t *testing.T) {
	var (
		cfgLabels = []string{"cfg:host"}
		regLabels = []string{"reg:host"}
	)

	tests := []struct {
		name string
		arg  string
		cfg  []string
		reg  []string
		want []string
	}{
		{"flag wins", "flag:host,other", cfgLabels, regLabels, []string{"flag:host", "other"}},
		{"config wins over registration", "", cfgLabels, regLabels, cfgLabels},
		{"registration is the fallback", "", nil, regLabels, regLabels},
		{"blank flag is ignored", " , ", cfgLabels, regLabels, cfgLabels},
		{"nothing configured", "", nil, nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, resolveLabels(tt.arg, tt.cfg, tt.reg))
		})
	}
}

func TestGetDockerSocketPathUsesConfigAndEnvironment(t *testing.T) {
	got, err := getDockerSocketPath("tcp://docker.example:2376")
	require.NoError(t, err)
	require.Equal(t, "tcp://docker.example:2376", got)

	t.Setenv("DOCKER_HOST", "unix:///tmp/docker.sock")
	got, err = getDockerSocketPath("-")
	require.NoError(t, err)
	require.Equal(t, "unix:///tmp/docker.sock", got)
}

func TestInitLoggingSetsLevelAndCaller(t *testing.T) {
	oldLevel := log.GetLevel()
	oldReportCaller := log.StandardLogger().ReportCaller
	t.Cleanup(func() {
		log.SetLevel(oldLevel)
		log.SetReportCaller(oldReportCaller)
	})

	oldFormatter := log.StandardLogger().Formatter
	t.Cleanup(func() { log.SetFormatter(oldFormatter) })

	cfg := &config.Config{}
	cfg.Log.Level = "debug"
	initLogging(cfg)

	require.Equal(t, log.DebugLevel, log.GetLevel())
	require.True(t, log.StandardLogger().ReportCaller)
	// act plans a job on this logger, so a live task's secrets have to be masked out of it
	require.IsType(t, report.MaskingFormatter(nil), log.StandardLogger().Formatter)
}

type daemonTestService struct {
	runnerv1connect.UnimplementedRunnerServiceHandler
	beforeTask   bool
	unregistered bool
	cancel       context.CancelFunc
}

func (*daemonTestService) Declare(context.Context, *connect.Request[runnerv1.DeclareRequest]) (*connect.Response[runnerv1.DeclareResponse], error) {
	return connect.NewResponse(&runnerv1.DeclareResponse{Runner: &runnerv1.Runner{Name: "test"}}), nil
}

func (s *daemonTestService) FetchTask(context.Context, *connect.Request[runnerv1.FetchTaskRequest]) (*connect.Response[runnerv1.FetchTaskResponse], error) {
	if s.unregistered {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unregistered runner"))
	}
	if s.beforeTask {
		s.cancel()
		return nil, connect.NewError(connect.CodeCanceled, errors.New("interrupted before task"))
	}
	return connect.NewResponse(&runnerv1.FetchTaskResponse{Task: &runnerv1.Task{Id: 1, Context: &structpb.Struct{}, WorkflowPayload: []byte("invalid: [")}}), nil
}

func (*daemonTestService) UpdateTask(_ context.Context, req *connect.Request[runnerv1.UpdateTaskRequest]) (*connect.Response[runnerv1.UpdateTaskResponse], error) {
	return connect.NewResponse(&runnerv1.UpdateTaskResponse{State: req.Msg.State}), nil
}

func (*daemonTestService) UpdateLog(_ context.Context, req *connect.Request[runnerv1.UpdateLogRequest]) (*connect.Response[runnerv1.UpdateLogResponse], error) {
	return connect.NewResponse(&runnerv1.UpdateLogResponse{AckIndex: req.Msg.Index + int64(len(req.Msg.Rows))}), nil
}

func TestDaemonRemovesOnlyConsumedEphemeralRegistration(t *testing.T) {
	t.Setenv("OTEL_SDK_DISABLED", "true")
	for _, tc := range []struct {
		name                                string
		ephemeral, beforeTask, unregistered bool
	}{
		{"ephemeral failed task", true, false, false},
		{"ephemeral interrupted before task", true, true, false},
		{"persistent once", false, false, false},
		{"ephemeral unregistered", true, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()
			_, handler := runnerv1connect.NewRunnerServiceHandler(&daemonTestService{beforeTask: tc.beforeTask, unregistered: tc.unregistered, cancel: cancel})
			server := httptest.NewServer(http.StripPrefix("/api/actions", handler))
			defer server.Close()
			cfg, err := config.LoadDefault("")
			require.NoError(t, err)
			cfg.Runner.File = filepath.Join(t.TempDir(), "registration.json")
			cfg.Runner.Labels = []string{"host:host"}
			cfg.Cache.Enabled = new(false)
			cfg.Runner.ReportCloseTimeout = time.Second
			file := filepath.Join(t.TempDir(), "config.yaml")
			content, err := yaml.Marshal(cfg)
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(file, content, 0o600))
			require.NoError(t, config.SaveRegistration(cfg.Runner.File, &config.Registration{Address: server.URL, UUID: "test", Name: "test", Token: "test", Labels: cfg.Runner.Labels, Ephemeral: tc.ephemeral}))
			err = runDaemon(ctx, &daemonArgs{Once: true}, &file)(nil, nil)
			if tc.unregistered {
				require.ErrorContains(t, err, "no longer registered")
			} else {
				require.NoError(t, err)
			}
			if tc.ephemeral && !tc.beforeTask {
				require.NoFileExists(t, cfg.Runner.File)
			} else {
				require.FileExists(t, cfg.Runner.File)
			}
		})
	}
}
