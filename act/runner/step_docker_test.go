// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2022 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/act/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// runStepDockerMain reports how many times the step container was actually removed.
func runStepDockerMain(t *testing.T, autoRemove bool) (int, *container.NewContainerInput) {
	cm := &containerMock{}

	var input *container.NewContainerInput

	// mock the new container call
	origContainerNewContainer := ContainerNewContainer
	ContainerNewContainer = func(containerInput *container.NewContainerInput) container.ExecutionsEnvironment {
		input = containerInput
		return cm
	}
	defer (func() {
		ContainerNewContainer = origContainerNewContainer
	})()

	ctx := context.Background()

	sd := &stepDocker{
		RunContext: &RunContext{
			StepResults: map[string]*model.StepResult{},
			Config: &Config{
				AutoRemove: autoRemove,
				Secrets: map[string]string{
					"DOCKER_USERNAME": "docker-user",
					"DOCKER_PASSWORD": "docker-password",
				},
			},
			Run: &model.Run{
				JobID: "1",
				Workflow: &model.Workflow{
					Jobs: map[string]*model.Job{
						"1": {
							Defaults: model.Defaults{
								Run: model.RunDefaults{
									Shell: "bash",
								},
							},
						},
					},
				},
			},
			JobContainer: cm,
		},
		Step: &model.Step{
			ID:               "1",
			Uses:             "docker://node:14",
			WorkingDirectory: "workdir",
		},
	}
	sd.RunContext.ExprEval = sd.RunContext.NewExpressionEvaluator(ctx)

	cm.On("Pull", false).Return(func(ctx context.Context) error {
		return nil
	})

	// Count executions, not calls: IfBool gates the executor, not the Remove() that builds it.
	removes := 0
	cm.On("Remove").Return(func(ctx context.Context) error {
		removes++
		return nil
	})

	cm.On("Create", []string(nil), []string(nil)).Return(func(ctx context.Context) error {
		return nil
	})

	cm.On("Start", true).Return(func(ctx context.Context) error {
		return nil
	})

	cm.On("Close").Return(func(ctx context.Context) error {
		return nil
	})

	cm.On("Copy", "/var/run/act", mock.AnythingOfType("[]*container.FileEntry")).Return(func(ctx context.Context) error {
		return nil
	})

	cm.On("UpdateFromEnv", "/var/run/act/workflow/envs.txt", mock.AnythingOfType("*map[string]string")).Return(func(ctx context.Context) error {
		return nil
	})

	cm.On("UpdateFromEnv", "/var/run/act/workflow/statecmd.txt", mock.AnythingOfType("*map[string]string")).Return(func(ctx context.Context) error {
		return nil
	})

	cm.On("UpdateFromEnv", "/var/run/act/workflow/outputcmd.txt", mock.AnythingOfType("*map[string]string")).Return(func(ctx context.Context) error {
		return nil
	})

	cm.On("GetContainerArchive", ctx, "/var/run/act/workflow/pathcmd.txt").Return(io.NopCloser(&bytes.Buffer{}), nil)

	err := sd.main()(ctx)
	assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act

	cm.AssertExpectations(t)
	return removes, input
}

func TestStepDockerMain(t *testing.T) {
	removes, input := runStepDockerMain(t, false)

	assert.Equal(t, "node:14", input.Image)

	// DOCKER_USERNAME/DOCKER_PASSWORD secrets should not be used as implicit pull credentials for docker:// action containers.
	assert.Empty(t, input.Username)
	assert.Empty(t, input.Password)

	// act owns removal: stale container before create, then after the run.
	assert.Equal(t, 2, removes)
}

// With AutoRemove the daemon reaps the container on exit, so act must not remove it
// afterwards and race that teardown.
func TestStepDockerMainAutoRemove(t *testing.T) {
	removes, _ := runStepDockerMain(t, true)

	assert.Equal(t, 1, removes)
}

func TestStepDockerNewStepContainerAllocatePTY(t *testing.T) {
	for _, tc := range []struct {
		name     string
		allocPTY bool
	}{
		{name: "off", allocPTY: false},
		{name: "on", allocPTY: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cm := &containerMock{}

			var captured *container.NewContainerInput
			origContainerNewContainer := ContainerNewContainer
			ContainerNewContainer = func(input *container.NewContainerInput) container.ExecutionsEnvironment {
				captured = input
				return cm
			}
			defer func() {
				ContainerNewContainer = origContainerNewContainer
			}()

			ctx := context.Background()
			sd := &stepDocker{
				RunContext: &RunContext{
					StepResults: map[string]*model.StepResult{},
					Config: &Config{
						AllocatePTY: tc.allocPTY,
						PlatformPicker: func(_ []string) string {
							return "node:14"
						},
					},
					Run: &model.Run{
						JobID: "1",
						Workflow: &model.Workflow{
							Jobs: map[string]*model.Job{"1": {}},
						},
					},
					JobContainer: cm,
				},
				Step: &model.Step{ID: "1", Uses: "docker://node:14"},
			}
			sd.RunContext.ExprEval = sd.RunContext.NewExpressionEvaluator(ctx)

			_ = sd.newStepContainer(ctx, "node:14", []string{"echo", "hi"}, nil)
			assert.Equal(t, tc.allocPTY, captured.AllocatePTY)
		})
	}
}

func TestStepDockerPrePost(t *testing.T) {
	ctx := context.Background()
	sd := &stepDocker{}

	err := sd.pre()(ctx)
	assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act

	err = sd.post()(ctx)
	assert.NoError(t, err)
}

func TestStepDockerNewStepContainerNetworkMode(t *testing.T) {
	cases := []struct {
		name          string
		platform      string
		expectDefault bool
	}{
		{
			name:          "docker mode attaches to job container network",
			platform:      "node:14",
			expectDefault: false,
		},
		{
			name:          "host mode uses default network",
			platform:      "-self-hosted",
			expectDefault: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cm := &containerMock{}

			var captured *container.NewContainerInput
			origContainerNewContainer := ContainerNewContainer
			ContainerNewContainer = func(input *container.NewContainerInput) container.ExecutionsEnvironment {
				captured = input
				return cm
			}
			defer func() {
				ContainerNewContainer = origContainerNewContainer
			}()

			ctx := context.Background()

			platform := tc.platform
			sd := &stepDocker{
				RunContext: &RunContext{
					StepResults: map[string]*model.StepResult{},
					Config: &Config{
						PlatformPicker: func(_ []string) string {
							return platform
						},
					},
					Run: &model.Run{
						JobID: "1",
						Workflow: &model.Workflow{
							Jobs: map[string]*model.Job{
								"1": {},
							},
						},
					},
					JobContainer: cm,
				},
				Step: &model.Step{
					ID:   "1",
					Uses: "docker://alpine:3.20",
				},
			}
			sd.RunContext.ExprEval = sd.RunContext.NewExpressionEvaluator(ctx)

			assert.Equal(t, tc.expectDefault, sd.RunContext.IsHostEnv(ctx),
				"IsHostEnv mismatch for platform %q", tc.platform)

			_ = sd.newStepContainer(ctx, "alpine:3.20", []string{"echo", "hello"}, nil)

			if tc.expectDefault {
				assert.Equal(t, "default", captured.NetworkMode,
					"host-mode step container must use 'default' network, got %q",
					captured.NetworkMode)
			} else {
				assert.True(t, strings.HasPrefix(captured.NetworkMode, "container:"),
					"docker-mode step container must attach to job container network, got %q",
					captured.NetworkMode)
			}
		})
	}
}
