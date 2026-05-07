// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2022 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/common/git"
	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/act/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type closerMock struct {
	mock.Mock
}

func (m *closerMock) Close() error {
	m.Called()
	return nil
}

func TestActionReader(t *testing.T) {
	yaml := strings.ReplaceAll(`
name: 'name'
runs:
	using: 'node16'
	main: 'main.js'
`, "\t", "  ")

	table := []struct {
		name        string
		step        *model.Step
		filename    string
		fileContent string
		expected    *model.Action
	}{
		{
			name:        "readActionYml",
			step:        &model.Step{},
			filename:    "action.yml",
			fileContent: yaml,
			expected: &model.Action{
				Name: "name",
				Runs: model.ActionRuns{
					Using:  "node16",
					Main:   "main.js",
					PreIf:  "always()",
					PostIf: "always()",
				},
			},
		},
		{
			name:        "readActionYaml",
			step:        &model.Step{},
			filename:    "action.yaml",
			fileContent: yaml,
			expected: &model.Action{
				Name: "name",
				Runs: model.ActionRuns{
					Using:  "node16",
					Main:   "main.js",
					PreIf:  "always()",
					PostIf: "always()",
				},
			},
		},
		{
			name:        "readDockerfile",
			step:        &model.Step{},
			filename:    "Dockerfile",
			fileContent: "FROM ubuntu:20.04",
			expected: &model.Action{
				Name: "(Synthetic)",
				Runs: model.ActionRuns{
					Using: "docker",
					Image: "Dockerfile",
				},
			},
		},
		{
			name: "readWithArgs",
			step: &model.Step{
				With: map[string]string{
					"args": "cmd",
				},
			},
			expected: &model.Action{
				Name: "(Synthetic)",
				Inputs: map[string]model.Input{
					"cwd": {
						Description: "(Actual working directory)",
						Required:    false,
						Default:     "actionDir/actionPath",
					},
					"command": {
						Description: "(Actual program)",
						Required:    false,
						Default:     "cmd",
					},
				},
				Runs: model.ActionRuns{
					Using: "node12",
					Main:  "trampoline.js",
				},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			closerMock := &closerMock{}

			readFile := func(filename string) (io.Reader, io.Closer, error) {
				if tt.filename != filename {
					return nil, nil, fs.ErrNotExist
				}

				return strings.NewReader(tt.fileContent), closerMock, nil
			}

			writeFile := func(filename string, data []byte, perm fs.FileMode) error {
				assert.Equal(t, "actionDir/actionPath/trampoline.js", filename)
				assert.Equal(t, fs.FileMode(0o400), perm)
				return nil
			}

			if tt.filename != "" {
				closerMock.On("Close")
			}

			action, err := readActionImpl(context.Background(), tt.step, "actionDir", "actionPath", readFile, writeFile)

			assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act
			assert.Equal(t, tt.expected, action)

			closerMock.AssertExpectations(t)
		})
	}
}

func TestActionRunner(t *testing.T) {
	table := []struct {
		name        string
		step        actionStep
		expectedEnv map[string]string
	}{
		{
			name: "with-input",
			step: &stepActionRemote{
				Step: &model.Step{
					Uses: "org/repo/path@ref",
				},
				RunContext: &RunContext{
					Config: &Config{},
					Run: &model.Run{
						JobID: "job",
						Workflow: &model.Workflow{
							Jobs: map[string]*model.Job{
								"job": {
									Name: "job",
								},
							},
						},
					},
				},
				action: &model.Action{
					Inputs: map[string]model.Input{
						"key": {
							Default: "default value",
						},
					},
					Runs: model.ActionRuns{
						Using: "node16",
					},
				},
				env: map[string]string{},
			},
			expectedEnv: map[string]string{"INPUT_KEY": "default value"},
		},
		{
			name: "restore-saved-state",
			step: &stepActionRemote{
				Step: &model.Step{
					ID:   "step",
					Uses: "org/repo/path@ref",
				},
				RunContext: &RunContext{
					ActionPath: "path",
					Config:     &Config{},
					Run: &model.Run{
						JobID: "job",
						Workflow: &model.Workflow{
							Jobs: map[string]*model.Job{
								"job": {
									Name: "job",
								},
							},
						},
					},
					CurrentStep: "post-step",
					StepResults: map[string]*model.StepResult{
						"step": {},
					},
					IntraActionState: map[string]map[string]string{
						"step": {
							"name": "state value",
						},
					},
				},
				action: &model.Action{
					Runs: model.ActionRuns{
						Using: "node16",
					},
				},
				env: map[string]string{},
			},
			expectedEnv: map[string]string{"STATE_name": "state value"},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			cm := &containerMock{}
			cm.On("CopyDir", "/var/run/act/actions/dir/", "dir/", false).Return(func(ctx context.Context) error { return nil })

			envMatcher := mock.MatchedBy(func(env map[string]string) bool {
				for k, v := range tt.expectedEnv {
					if env[k] != v {
						return false
					}
				}
				return true
			})

			cm.On("Exec", []string{"node", "/var/run/act/actions/dir/path"}, envMatcher, "", "").Return(func(ctx context.Context) error { return nil })

			tt.step.getRunContext().JobContainer = cm

			err := runActionImpl(tt.step, "dir", newRemoteAction("org/repo/path@ref"))(ctx)

			assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act
			cm.AssertExpectations(t)
		})
	}
}

// TestExecAsDockerHoldsCloneLockForRemoteUncached guards the fix for the
// docker-action build context race: when a remote uses: docker action runs
// without ActionCache configured, container.NewDockerBuildExecutor walks
// contextDir directly to send the build context to the daemon. contextDir
// lives under <ActionCacheDir>/<UsesHash> where a peer matrix job's
// NewGitCloneExecutor can `git checkout --force` the worktree at any time;
// the build executor must hold the per-dir clone lock for its duration.
func TestExecAsDockerHoldsCloneLockForRemoteUncached(t *testing.T) {
	basedir := t.TempDir()
	contextDir, _ := filepath.Split(filepath.Join(basedir, "Dockerfile"))

	unlockOnce := sync.OnceFunc(git.AcquireCloneLock(contextDir))
	defer unlockOnce()

	innerEntered := make(chan struct{})
	releaseInner := make(chan struct{})
	releaseOnce := sync.OnceFunc(func() { close(releaseInner) })
	defer releaseOnce()

	origImageExists := ContainerImageExistsLocally
	ContainerImageExistsLocally = func(_ context.Context, _, _ string) (bool, error) {
		return false, nil
	}
	defer func() { ContainerImageExistsLocally = origImageExists }()

	origBuildExec := ContainerNewDockerBuildExecutor
	ContainerNewDockerBuildExecutor = func(_ container.NewDockerBuildExecutorInput) common.Executor {
		return func(_ context.Context) error {
			close(innerEntered)
			<-releaseInner
			return nil
		}
	}
	defer func() { ContainerNewDockerBuildExecutor = origBuildExec }()

	cm := &containerMock{}
	cm.On("Pull", false).Return(func(ctx context.Context) error { return ctx.Err() }).Maybe()
	cm.On("Remove").Return(func(ctx context.Context) error { return nil }).Maybe()
	cm.On("Create", []string(nil), []string(nil)).Return(func(ctx context.Context) error { return nil }).Maybe()
	cm.On("Start", true).Return(func(ctx context.Context) error { return nil }).Maybe()
	cm.On("Close").Return(func(ctx context.Context) error { return nil }).Maybe()

	step := &stepActionRemote{
		Step: &model.Step{ID: "1", Uses: "remote/action@v1", With: map[string]string{}},
		RunContext: &RunContext{
			Config: &Config{},
			Run: &model.Run{
				JobID: "1",
				Workflow: &model.Workflow{
					Name: "wf",
					Jobs: map[string]*model.Job{"1": {}},
				},
			},
			JobContainer: cm,
		},
		action: &model.Action{Runs: model.ActionRuns{Using: "docker", Image: "Dockerfile"}},
		env:    map[string]string{},
	}

	origNewContainer := ContainerNewContainer
	ContainerNewContainer = func(_ *container.NewContainerInput) container.ExecutionsEnvironment {
		return cm
	}
	defer func() { ContainerNewContainer = origNewContainer }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- execAsDocker(ctx, step, "test-action", basedir, false) }()

	select {
	case <-innerEntered:
		t.Fatal("inner build executor ran before clone lock was released; wrap not held")
	case err := <-done:
		t.Fatalf("execAsDocker returned before inner was entered: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	unlockOnce()

	select {
	case <-innerEntered:
	case err := <-done:
		t.Fatalf("execAsDocker returned without entering inner: %v", err)
	case <-time.After(time.Second):
		t.Fatal("inner build executor not entered within 1s after clone lock released")
	}

	cancel()
	releaseOnce()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("execAsDocker did not return after inner was released and ctx was canceled")
	}
}
