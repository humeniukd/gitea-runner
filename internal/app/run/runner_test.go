// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/runner"
	"gitea.com/gitea/runner/internal/app/poll"
	clientmocks "gitea.com/gitea/runner/internal/pkg/client/mocks"
	"gitea.com/gitea/runner/internal/pkg/config"
	"gitea.com/gitea/runner/internal/pkg/labels"
	"gitea.com/gitea/runner/internal/pkg/ver"

	"connectrpc.com/connect"
	runnerv1 "gitea.dev/actionslib/runner/v1"
	"github.com/stretchr/testify/assert"
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

func TestRunnerReclaimsVolumesAfterReporting(t *testing.T) {
	for _, mode := range []string{"regular", "ephemeral", "post-task script"} {
		t.Run(mode, func(t *testing.T) {
			cfg, err := config.LoadDefault("")
			require.NoError(t, err)
			cfg.Cache.Enabled = new(false)
			cfg.HealthCheck.Enabled = false
			cfg.Runner.Capacity = 1
			cfg.Container.Network = "host"
			cfg.Container.DockerHost = "-"
			cfg.Container.ForcePull = false
			cfg.Host.WorkdirParent = t.TempDir()
			if mode == "post-task script" {
				if runtime.GOOS == "windows" {
					t.Skip("uses a POSIX script")
				}
				cfg.Runner.PostTaskScript = filepath.Join(t.TempDir(), "post-task.sh")
				require.NoError(t, os.WriteFile(cfg.Runner.PostTaskScript, []byte("#!/bin/sh\n: > \"$0.done\"\n"), 0o700))
			}
			cli := clientmocks.NewClient(t)
			cli.AddressValue = "https://gitea.example/"
			r := NewRunner(cfg, &config.Registration{UUID: "runner-1", Labels: []string{"ubuntu:docker://node:20"}}, cli)
			t.Cleanup(func() { require.NoError(t, r.Close()) })
			task := &runnerv1.Task{
				Id: 123,
				Context: &structpb.Struct{Fields: map[string]*structpb.Value{
					"repository": structpb.NewStringValue("owner/repo"),
				}},
				WorkflowPayload: []byte("name: test\njobs:\n  job:\n    runs-on: ubuntu\n    steps:\n      - run: exit 0\n        if: false\n"),
				Secrets:         map[string]string{},
			}
			var reported atomic.Bool
			var removed atomic.Int64
			cli.On("UpdateLog", mock.Anything, mock.Anything).Return(func(_ context.Context, req *connect.Request[runnerv1.UpdateLogRequest]) (*connect.Response[runnerv1.UpdateLogResponse], error) {
				assert.False(t, reported.Load())
				return connect.NewResponse(&runnerv1.UpdateLogResponse{AckIndex: req.Msg.Index + int64(len(req.Msg.Rows))}), nil
			})
			cli.On("UpdateTask", mock.Anything, mock.Anything).Return(func(_ context.Context, req *connect.Request[runnerv1.UpdateTaskRequest]) (*connect.Response[runnerv1.UpdateTaskResponse], error) {
				if req.Msg.State.Result != runnerv1.Result_RESULT_UNSPECIFIED {
					assert.Equal(t, runnerv1.Result_RESULT_SUCCESS, req.Msg.State.Result)
					assert.Equal(t, int64(1), r.RunningCount())
					if mode == "post-task script" {
						assert.Equal(t, int64(2), removed.Load())
						assert.FileExists(t, cfg.Runner.PostTaskScript+".done")
					} else {
						assert.Zero(t, removed.Load())
					}
					reported.Store(true)
				}
				return connect.NewResponse(&runnerv1.UpdateTaskResponse{State: req.Msg.State}), nil
			})
			removalStarted := make(chan struct{})
			releaseRemoval := make(chan struct{})
			var releaseOnce sync.Once
			release := func() { releaseOnce.Do(func() { close(releaseRemoval) }) }
			t.Cleanup(release)
			var volumeNames sync.Map
			daemon := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json")
				path := strings.TrimPrefix(request.URL.Path, "/v1.47")
				switch {
				case path == "/_ping":
					writer.Header().Set("API-Version", "1.47")
					_, _ = io.WriteString(writer, "OK")
				case path == "/info":
					_, _ = io.WriteString(writer, `{"Architecture":"amd64","OSType":"linux"}`)
				case path == "/containers/json", path == "/networks":
					_, _ = io.WriteString(writer, "[]")
				case strings.HasPrefix(path, "/images/"):
					_, _ = io.WriteString(writer, `{"Id":"image-id","Config":{},"Os":"linux","Architecture":"amd64"}`)
				case path == "/volumes":
					_, _ = io.WriteString(writer, `{"Volumes":[]}`)
				case path == "/volumes/create":
					var volume struct {
						Name   string
						Labels map[string]string
					}
					assert.NoError(t, json.UnmarshalRead(request.Body, &volume))
					assert.Equal(t, map[string]string{"com.gitea.runner.uuid": "runner-1"}, volume.Labels)
					volumeNames.Store(volume.Name, struct{}{})
					_, _ = fmt.Fprintf(writer, `{"Name":%q}`, volume.Name)
				case request.Method == http.MethodDelete && strings.HasPrefix(path, "/volumes/"):
					if _, exists := volumeNames.LoadAndDelete(strings.TrimPrefix(path, "/volumes/")); exists {
						assert.Equal(t, mode != "post-task script", reported.Load())
						if mode == "post-task script" {
							assert.NoFileExists(t, cfg.Runner.PostTaskScript+".done")
						}
						assert.Equal(t, int64(1), r.RunningCount())
						if removed.Add(1) == 1 {
							close(removalStarted)
							<-releaseRemoval
						}
					}
					writer.WriteHeader(http.StatusNoContent)
				case path == "/containers/create":
					_, _ = io.WriteString(writer, `{"Id":"job-id"}`)
				case path == "/containers/job-id/json":
					_, _ = io.WriteString(writer, `{"Id":"job-id","Config":{},"State":{"Status":"running"}}`)
				case request.Method == http.MethodDelete || strings.HasSuffix(path, "/start") || strings.HasSuffix(path, "/kill") || strings.HasSuffix(path, "/archive"):
					writer.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected Docker request: %s %s", request.Method, request.URL)
					http.NotFound(writer, request)
				}
			}))
			t.Cleanup(func() { release(); daemon.Close() })
			t.Setenv("DOCKER_HOST", daemon.URL)
			var fetched atomic.Bool
			cli.On("FetchTask", mock.Anything, mock.Anything).Return(func(ctx context.Context, _ *connect.Request[runnerv1.FetchTaskRequest]) (*connect.Response[runnerv1.FetchTaskResponse], error) {
				if !fetched.Swap(true) {
					return connect.NewResponse(&runnerv1.FetchTaskResponse{Task: task}), nil
				}
				assert.Equal(t, int64(2), removed.Load())
				<-ctx.Done()
				return nil, ctx.Err()
			})
			poller := poll.New(cfg, cli, r)
			t.Cleanup(func() {
				release()
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_ = poller.Shutdown(ctx)
			})
			if mode == "ephemeral" {
				go poller.PollOnce()
			} else {
				go poller.Poll()
			}
			select {
			case <-removalStarted:
			case <-time.After(10 * time.Second):
				t.Fatal("volume removal was never reached")
			}
			assert.Equal(t, int64(1), r.RunningCount())
			select {
			case <-poller.Done():
				t.Fatal("poller finished before volume reclamation")
			default:
			}
			shutdownCtx, cancel := context.WithCancel(t.Context())
			cancel()
			shutdownEntered := make(chan struct{})
			shutdownDone := make(chan error, 1)
			go func() {
				shutdownDone <- poller.Shutdown(&shutdownWaitContext{Context: shutdownCtx, entered: shutdownEntered})
			}()
			<-shutdownEntered
			select {
			case err := <-shutdownDone:
				t.Fatalf("shutdown returned before volume reclamation: %v", err)
			default:
			}
			release()
			if err := <-shutdownDone; err != nil {
				require.ErrorIs(t, err, context.Canceled)
			}
			assert.True(t, reported.Load())
			assert.Equal(t, int64(2), removed.Load())
			assert.Zero(t, r.RunningCount())
		})
	}
}

type shutdownWaitContext struct {
	context.Context
	entered chan struct{}
	once    sync.Once
}

func (ctx *shutdownWaitContext) Done() <-chan struct{} {
	ctx.once.Do(func() { close(ctx.entered) })
	return ctx.Context.Done()
}

func TestCleanupJobVolumesJoinsErrorsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	firstError, secondError := errors.New("workspace removal"), errors.New("environment removal")
	err := cleanupJobVolumes(ctx, []common.Executor{
		func(ctx context.Context) error {
			require.NoError(t, ctx.Err())
			deadline, ok := ctx.Deadline()
			assert.True(t, ok)
			assert.InDelta(t, time.Minute.Seconds(), time.Until(deadline).Seconds(), 1)
			return firstError
		},
		func(context.Context) error { return secondError },
	})
	require.ErrorIs(t, err, firstError)
	require.ErrorIs(t, err, secondError)
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
	cli.AddressValue = "https://gitea.example/"

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
	require.Empty(t, r.envs[runner.CacheServiceV2Env], "no cache server, nothing to serve v2 from")
}

func TestRunnerFallbackPlatform(t *testing.T) {
	tests := []struct {
		name          string
		label         string
		dockerRunning bool
		want          string
	}{
		{"a docker label needs no daemon probe", "ubuntu:docker://node:18", false, "mirror.example/ci:noble"},
		{"host labels keep the image where docker runs", "ubuntu:host", true, "mirror.example/ci:noble"},
		{"host labels without docker run on the host", "ubuntu:host", false, labels.SelfHostedPlatform},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reachable := dockerReachable
			dockerReachable = func(context.Context) bool { return tt.dockerRunning }
			t.Cleanup(func() { dockerReachable = reachable })

			label, err := labels.Parse(tt.label)
			require.NoError(t, err)
			cfg := &config.Config{}
			cfg.Runner.DefaultImage = "mirror.example/ci:noble"
			r := &Runner{cfg: cfg, labels: labels.Labels{label}}

			require.Equal(t, tt.want, r.fallbackPlatform(t.Context()))
		})
	}
}

// Proxy variables are assembled per task, because a job's service containers have to be
// reached directly and they are only known once the workflow is parsed.
func TestNewRunnerLeavesProxyToTheTask(t *testing.T) {
	clearProxyEnv(t)
	t.Setenv("http_proxy", "http://proxy:3128")

	cfg := &config.Config{}
	cfg.Cache.ExternalServer = "http://cache.local:8088/"
	reg := &config.Registration{Name: "runner"}
	cli := clientmocks.NewClient(t)
	cli.AddressValue = "https://gitea.example/"

	r := NewRunner(cfg, reg, cli)

	require.NotContains(t, r.envs, "http_proxy")
	require.NotContains(t, r.envs, "no_proxy")
	assert.Empty(t, r.builtInCacheURL(), "an external cache server is the operator's to exempt, not ours")
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

// The results service is decided per task, because that is where the job's token and its
// instance are both known. NewRunner only leaves the address Gitea serves.
func TestNewRunnerCacheServiceV2(t *testing.T) {
	cfg := &config.Config{}
	cfg.Cache.Dir, cfg.Cache.Host = t.TempDir(), "127.0.0.1"
	cli := clientmocks.NewClient(t)
	cli.AddressValue = "https://gitea.example/"

	r := NewRunner(cfg, &config.Registration{Name: "runner"}, cli)
	t.Cleanup(func() { _ = r.Close() })
	const token = "task-token"

	assert.Equal(t, "https://gitea.example", r.envs["ACTIONS_RESULTS_URL"])
	assert.Empty(t, r.envs[runner.CacheServiceV2Env], "a promise the runner has not made yet")

	// The registration is what makes it true: the cache server takes the results service over,
	// having been told which instance to forward the artifact half to.
	revoke, resultsURL := r.registerCacheForTask(token, "owner/repo", nil)
	defer revoke()
	require.Equal(t, r.cacheHandler.ExternalURL(), resultsURL)

	// And what is advertised has to answer the cache service. A client without the GHES escape
	// hatch, docker buildx among them, posts its cache calls at exactly this URL and nowhere else,
	// so a 404 here is the failure this whole change is about.
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost,
		resultsURL+"/twirp/github.actions.results.api.v1.CacheService/GetCacheEntryDownloadURL",
		strings.NewReader(`{"key":"k","version":"v"}`))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "the advertised results service serves no cache service")
	assert.Equal(t, r.envs["ACTIONS_CACHE_URL"], r.builtInCacheURL(), "the address only the runner knows is bypassed for the operator")

	envs := r.cloneEnvs()
	r.setResultsService(envs, resultsURL)
	assert.Equal(t, resultsURL, envs["ACTIONS_RESULTS_URL"])
	assert.Equal(t, "true", envs[runner.CacheServiceV2Env])

	cfg.Cache.V2 = new(bool)
	envs = r.cloneEnvs()
	envs[runner.CacheServiceV2Env] = "true"
	r.setResultsService(envs, resultsURL)
	assert.Equal(t, "https://gitea.example", envs["ACTIONS_RESULTS_URL"])
	assert.Empty(t, envs[runner.CacheServiceV2Env], "a runner.envs entry would promise v2 at an origin not serving it")

	envs = r.cloneEnvs()
	envs[runner.CacheServiceV2Env] = "true"
	envs["ACTIONS_RESULTS_URL"] = "https://gitea.example/sub"
	r.setResultsService(envs, "")
	assert.Equal(t, "https://gitea.example/sub", envs["ACTIONS_RESULTS_URL"], "with no cache server there is nothing to front with")
	assert.Empty(t, envs[runner.CacheServiceV2Env])

	for instance, insecure := range map[string]bool{
		"https://gitea.example/sub":   false,
		"https://self-signed.example": true,
	} {
		cfg.Runner.Insecure = insecure
		envs = r.cloneEnvs()
		envs["ACTIONS_RESULTS_URL"] = instance
		r.setResultsService(envs, resultsURL)
		assert.Equal(t, resultsURL, envs["ACTIONS_RESULTS_URL"], instance)
		assert.Empty(t, envs[runner.CacheServiceV2Env])
	}
}

// The v1 cache client appends its path to ACTIONS_CACHE_URL without a separator, so a configured
// server that is missing the slash would send it to a URL whose port swallows the path.
func TestNewRunnerNormalizesTheExternalCacheServer(t *testing.T) {
	cfg := &config.Config{}
	cfg.Cache.ExternalServer = "http://cache.local:8088//"
	cli := clientmocks.NewClient(t)
	cli.AddressValue = "https://gitea.example/"

	r := NewRunner(cfg, &config.Registration{Name: "runner"}, cli)

	assert.Equal(t, "http://cache.local:8088/", r.envs["ACTIONS_CACHE_URL"])
	// Nothing to front the results service with, so the variable stays unset and the client keeps
	// to v1, which reads the cache URL first.
	assert.Equal(t, "https://gitea.example", r.envs["ACTIONS_RESULTS_URL"])
	assert.Empty(t, r.envs[runner.CacheServiceV2Env])
}
