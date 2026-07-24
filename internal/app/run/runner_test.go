// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"context"
	"testing"

	clientmocks "gitea.com/gitea/runner/internal/pkg/client/mocks"
	"gitea.com/gitea/runner/internal/pkg/config"
	"gitea.com/gitea/runner/internal/pkg/ver"

	"connectrpc.com/connect"
	runnerv1 "gitea.dev/actions-proto-go/runner/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestRunnerCapabilitiesAndDeclare(t *testing.T) {
	require.Equal(t, []string{CapabilityCancelling}, RunnerCapabilities())

	cli := clientmocks.NewClient(t)
	cli.On("Declare", mock.Anything, mock.MatchedBy(func(req *connect.Request[runnerv1.DeclareRequest]) bool {
		return req.Msg.Version == ver.Version() &&
			len(req.Msg.Labels) == 1 &&
			req.Msg.Labels[0] == "ubuntu" &&
			len(req.Msg.Capabilities) == 1 &&
			req.Msg.Capabilities[0] == CapabilityCancelling
	})).Return(connect.NewResponse(&runnerv1.DeclareResponse{}), nil)

	r := &Runner{client: cli}
	_, err := r.Declare(context.Background(), []string{"ubuntu"})
	require.NoError(t, err)
}

func TestRunnerSetCapabilitiesFromDeclare(t *testing.T) {
	r := &Runner{}
	r.SetCapabilitiesFromDeclare(nil)
	require.Empty(t, r.capabilities)

	resp := connect.NewResponse(&runnerv1.DeclareResponse{})
	resp.Header().Set("X-Gitea-Actions-Capabilities", " cancelling,cache-v2 ")
	r.SetCapabilitiesFromDeclare(resp)
	require.Equal(t, "cancelling,cache-v2", r.capabilities)
}

func TestRunnerDefaultActionsURLUsesMirrorOnlyForGithub(t *testing.T) {
	r := &Runner{cfg: &config.Config{}}
	r.cfg.Runner.GithubMirror = "https://mirror.example"

	task := taskWithDefaultActionsURL("https://github.com")
	require.Equal(t, "https://mirror.example", r.getDefaultActionsURL(task))

	task = taskWithDefaultActionsURL("https://gitea.example")
	require.Equal(t, "https://gitea.example", r.getDefaultActionsURL(task))
}

func TestRunnerRunningCountAndNullLogger(t *testing.T) {
	r := &Runner{}
	require.Equal(t, int64(0), r.RunningCount())
	r.runningCount.Add(2)
	require.Equal(t, int64(2), r.RunningCount())

	logger := NullLogger{}.WithJobLogger()
	require.NotNil(t, logger)
	require.NotNil(t, logger.Out)
}

func TestNewRunnerInitializesLabelsAndEnvironment(t *testing.T) {
	cacheEnabled := false
	cfg := &config.Config{}
	cfg.Cache.Enabled = &cacheEnabled
	cfg.Runner.Envs = map[string]string{"EXISTING": "value"}
	reg := &config.Registration{
		Name:   "runner",
		Labels: []string{"ubuntu:host", "", "pool:e57e18d4"},
	}
	cli := clientmocks.NewClient(t)
	cli.On("Address").Return("https://gitea.example/").Maybe()

	r := NewRunner(cfg, reg, cli)

	require.Equal(t, "runner", r.name)
	require.Len(t, r.labels, 2)
	require.Equal(t, []string{"ubuntu", "pool:e57e18d4"}, r.labels.Names())
	require.Equal(t, "value", r.envs["EXISTING"])
	require.Equal(t, "https://gitea.example/api/actions_pipeline/", r.envs["ACTIONS_RUNTIME_URL"])
	require.Equal(t, "https://gitea.example", r.envs["ACTIONS_RESULTS_URL"])
	require.Equal(t, "true", r.envs["GITEA_ACTIONS"])
	require.NotEmpty(t, r.envs["GITEA_ACTIONS_RUNNER_VERSION"])
	require.Nil(t, r.cacheHandler)
}

// A proxied runner must force only its own cache server into the jobs' no_proxy. Gitea
// (ACTIONS_RESULTS_URL) is reached through the proxy like everything else, so job steps
// that upload artifacts to it keep working when Gitea is only reachable via the proxy.
func TestNewRunnerDoesNotForceGiteaDirectForJobs(t *testing.T) {
	cfg := &config.Config{}
	cfg.Cache.ExternalServer = "http://cache.local:8088/"
	cfg.Proxy = config.Proxy{HTTPProxy: "http://proxy:3128", NoProxy: "internal.example"}
	reg := &config.Registration{Name: "runner"}
	cli := clientmocks.NewClient(t)
	cli.On("Address").Return("https://gitea.example/").Maybe()

	r := NewRunner(cfg, reg, cli)

	// The runner-local cache host is bypassed and the operator's own entry is kept.
	require.Contains(t, r.envs["no_proxy"], "cache.local")
	require.Contains(t, r.envs["no_proxy"], "internal.example")
	// Gitea must not be forced direct; the operator's no_proxy governs it instead.
	require.NotContains(t, r.envs["no_proxy"], "gitea.example")
}

func taskWithDefaultActionsURL(url string) *runnerv1.Task {
	return &runnerv1.Task{
		Context: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"gitea_default_actions_url": structpb.NewStringValue(url),
			},
		},
	}
}
