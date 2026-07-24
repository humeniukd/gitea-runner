// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2022 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"errors"
	"testing"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/model"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	yaml "go.yaml.in/yaml/v4"
)

func TestMergeIntoMap(t *testing.T) {
	table := []struct {
		name     string
		target   map[string]string
		maps     []map[string]string
		expected map[string]string
	}{
		{
			name:     "testEmptyMap",
			target:   map[string]string{},
			maps:     []map[string]string{},
			expected: map[string]string{},
		},
		{
			name:   "testMergeIntoEmptyMap",
			target: map[string]string{},
			maps: []map[string]string{
				{
					"key1": "value1",
					"key2": "value2",
				}, {
					"key2": "overridden",
					"key3": "value3",
				},
			},
			expected: map[string]string{
				"key1": "value1",
				"key2": "overridden",
				"key3": "value3",
			},
		},
		{
			name: "testMergeIntoExistingMap",
			target: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			maps: []map[string]string{
				{
					"key1": "overridden",
				},
			},
			expected: map[string]string{
				"key1": "overridden",
				"key2": "value2",
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			mergeIntoMapCaseSensitive(tt.target, tt.maps...)
			assert.Equal(t, tt.expected, tt.target)
			mergeIntoMapCaseInsensitive(tt.target, tt.maps...)
			assert.Equal(t, tt.expected, tt.target)
		})
	}
}

type stepMock struct {
	mock.Mock
	step
}

func (sm *stepMock) pre() common.Executor {
	args := sm.Called()
	return args.Get(0).(func(context.Context) error)
}

func (sm *stepMock) main() common.Executor {
	args := sm.Called()
	return args.Get(0).(func(context.Context) error)
}

func (sm *stepMock) post() common.Executor {
	args := sm.Called()
	return args.Get(0).(func(context.Context) error)
}

func (sm *stepMock) getRunContext() *RunContext {
	args := sm.Called()
	return args.Get(0).(*RunContext)
}

func (sm *stepMock) getGithubContext(ctx context.Context) *model.GithubContext {
	args := sm.Called()
	return args.Get(0).(*RunContext).getGithubContext(ctx)
}

func (sm *stepMock) getStepModel() *model.Step {
	args := sm.Called()
	return args.Get(0).(*model.Step)
}

func (sm *stepMock) getEnv() *map[string]string {
	args := sm.Called()
	return args.Get(0).(*map[string]string)
}

func TestSetupEnv(t *testing.T) {
	cm := &containerMock{}
	sm := &stepMock{}

	rc := &RunContext{
		Config: &Config{
			Env: map[string]string{
				"GITHUB_RUN_ID": "runId",
			},
		},
		Run: &model.Run{
			JobID: "1",
			Workflow: &model.Workflow{
				Jobs: map[string]*model.Job{
					"1": {
						Env: yaml.Node{
							Value: "JOB_KEY: jobvalue",
						},
					},
				},
			},
		},
		Env: map[string]string{
			"RC_KEY": "rcvalue",
		},
		JobContainer: cm,
	}
	step := &model.Step{
		With: map[string]string{
			"STEP_WITH": "with-value",
		},
	}
	env := map[string]string{}

	sm.On("getRunContext").Return(rc)
	sm.On("getGithubContext").Return(rc)
	sm.On("getStepModel").Return(step)
	sm.On("getEnv").Return(&env)

	err := setupEnv(context.Background(), sm)
	assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act

	// These are commit or system specific
	delete((env), "GITHUB_REF")
	delete((env), "GITHUB_REF_NAME")
	delete((env), "GITHUB_REF_TYPE")
	delete((env), "GITHUB_SHA")
	delete((env), "GITHUB_WORKSPACE")
	delete((env), "GITHUB_REPOSITORY")
	delete((env), "GITHUB_REPOSITORY_OWNER")
	delete((env), "GITHUB_ACTOR")

	assert.Equal(t, map[string]string{
		"ACT":                      "true",
		"ACT_SKIP_CHECKOUT":        "true",
		"CI":                       "true",
		"GITHUB_ACTION":            "",
		"GITHUB_ACTIONS":           "true",
		"GITHUB_ACTION_PATH":       "",
		"GITHUB_ACTION_REF":        "",
		"GITHUB_ACTION_REPOSITORY": "",
		"GITHUB_API_URL":           "https:///api/v1", // Gitea uses api/v1 (upstream GitHub: api/v3)
		"GITHUB_BASE_REF":          "",
		"GITHUB_EVENT_NAME":        "",
		"GITHUB_EVENT_PATH":        "/var/run/act/workflow/event.json",
		"GITHUB_GRAPHQL_URL":       "",
		"GITHUB_HEAD_REF":          "",
		"GITHUB_JOB":               "1",
		"GITHUB_RETENTION_DAYS":    "0",
		"GITHUB_RUN_ATTEMPT":       "",
		"GITHUB_RUN_ID":            "runId",
		"GITHUB_RUN_NUMBER":        "1",
		"GITHUB_SERVER_URL":        "https://",
		"GITHUB_WORKFLOW":          "",
		"INPUT_STEP_WITH":          "with-value",
		"RC_KEY":                   "rcvalue",
		"RUNNER_PERFLOG":           "/dev/null",
		"RUNNER_TRACKING_ID":       "",
	}, env)

	cm.AssertExpectations(t)
}

func TestIsStepEnabled(t *testing.T) {
	createTestStep := func(t *testing.T, input string) step {
		var step *model.Step
		err := yaml.Unmarshal([]byte(input), &step)
		assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act

		return &stepRun{
			RunContext: &RunContext{
				Config: &Config{
					Workdir: ".",
					Platforms: map[string]string{
						"ubuntu-latest": "ubuntu-latest",
					},
				},
				StepResults: map[string]*model.StepResult{},
				Env:         map[string]string{},
				Run: &model.Run{
					JobID: "job1",
					Workflow: &model.Workflow{
						Name: "workflow1",
						Jobs: map[string]*model.Job{
							"job1": createJob(t, `runs-on: ubuntu-latest`, ""),
						},
					},
				},
			},
			Step: step,
		}
	}

	log.SetLevel(log.DebugLevel)
	assertObject := assert.New(t)

	// success()
	step := createTestStep(t, "if: success()")
	assertObject.True(isStepEnabled(context.Background(), step.getIfExpression(context.Background(), stepStageMain), step, stepStageMain))

	step = createTestStep(t, "if: success()")
	step.getRunContext().StepResults["a"] = &model.StepResult{
		Conclusion: model.StepStatusSuccess,
	}
	assertObject.True(isStepEnabled(context.Background(), step.getStepModel().If.Value, step, stepStageMain))

	step = createTestStep(t, "if: success()")
	step.getRunContext().StepResults["a"] = &model.StepResult{
		Conclusion: model.StepStatusFailure,
	}
	assertObject.False(isStepEnabled(context.Background(), step.getStepModel().If.Value, step, stepStageMain))

	// failure()
	step = createTestStep(t, "if: failure()")
	assertObject.False(isStepEnabled(context.Background(), step.getStepModel().If.Value, step, stepStageMain))

	step = createTestStep(t, "if: failure()")
	step.getRunContext().StepResults["a"] = &model.StepResult{
		Conclusion: model.StepStatusSuccess,
	}
	assertObject.False(isStepEnabled(context.Background(), step.getStepModel().If.Value, step, stepStageMain))

	step = createTestStep(t, "if: failure()")
	step.getRunContext().StepResults["a"] = &model.StepResult{
		Conclusion: model.StepStatusFailure,
	}
	assertObject.True(isStepEnabled(context.Background(), step.getStepModel().If.Value, step, stepStageMain))

	// always()
	step = createTestStep(t, "if: always()")
	assertObject.True(isStepEnabled(context.Background(), step.getStepModel().If.Value, step, stepStageMain))

	step = createTestStep(t, "if: always()")
	step.getRunContext().StepResults["a"] = &model.StepResult{
		Conclusion: model.StepStatusSuccess,
	}
	assertObject.True(isStepEnabled(context.Background(), step.getStepModel().If.Value, step, stepStageMain))

	step = createTestStep(t, "if: always()")
	step.getRunContext().StepResults["a"] = &model.StepResult{
		Conclusion: model.StepStatusFailure,
	}
	assertObject.True(isStepEnabled(context.Background(), step.getStepModel().If.Value, step, stepStageMain))
}

func TestIsContinueOnError(t *testing.T) {
	createTestStep := func(t *testing.T, input string) step {
		var step *model.Step
		err := yaml.Unmarshal([]byte(input), &step)
		assert.NoError(t, err) //nolint:testifylint // pre-existing issue from nektos/act

		return &stepRun{
			RunContext: &RunContext{
				Config: &Config{
					Workdir: ".",
					Platforms: map[string]string{
						"ubuntu-latest": "ubuntu-latest",
					},
				},
				StepResults: map[string]*model.StepResult{},
				Env:         map[string]string{},
				Run: &model.Run{
					JobID: "job1",
					Workflow: &model.Workflow{
						Name: "workflow1",
						Jobs: map[string]*model.Job{
							"job1": createJob(t, `runs-on: ubuntu-latest`, ""),
						},
					},
				},
			},
			Step: step,
		}
	}

	log.SetLevel(log.DebugLevel)
	assertObject := assert.New(t)

	// absent
	step := createTestStep(t, "name: test")
	continueOnError, err := isContinueOnError(context.Background(), step.getStepModel().RawContinueOnError, step, stepStageMain)
	assertObject.False(continueOnError)
	assertObject.NoError(err) //nolint:testifylint // pre-existing issue from nektos/act

	// explcit true
	step = createTestStep(t, "continue-on-error: true")
	continueOnError, err = isContinueOnError(context.Background(), step.getStepModel().RawContinueOnError, step, stepStageMain)
	assertObject.True(continueOnError)
	assertObject.NoError(err) //nolint:testifylint // pre-existing issue from nektos/act

	// explicit false
	step = createTestStep(t, "continue-on-error: false")
	continueOnError, err = isContinueOnError(context.Background(), step.getStepModel().RawContinueOnError, step, stepStageMain)
	assertObject.False(continueOnError)
	assertObject.NoError(err) //nolint:testifylint // pre-existing issue from nektos/act

	// expression true
	step = createTestStep(t, "continue-on-error: ${{ 'test' == 'test' }}")
	continueOnError, err = isContinueOnError(context.Background(), step.getStepModel().RawContinueOnError, step, stepStageMain)
	assertObject.True(continueOnError)
	assertObject.NoError(err) //nolint:testifylint // pre-existing issue from nektos/act

	// expression false
	step = createTestStep(t, "continue-on-error: ${{ 'test' != 'test' }}")
	continueOnError, err = isContinueOnError(context.Background(), step.getStepModel().RawContinueOnError, step, stepStageMain)
	assertObject.False(continueOnError)
	assertObject.NoError(err) //nolint:testifylint // pre-existing issue from nektos/act

	// expression parse error
	step = createTestStep(t, "continue-on-error: ${{ 'test' != test }}")
	continueOnError, err = isContinueOnError(context.Background(), step.getStepModel().RawContinueOnError, step, stepStageMain)
	assertObject.False(continueOnError)
	assertObject.Error(err)
}

// fakeStep is a minimal step for driving runStepExecutor without the shell/action machinery.
type fakeStep struct {
	rc     *RunContext
	model  *model.Step
	env    map[string]string
	ifExpr string
}

func (s *fakeStep) pre() common.Executor       { return func(context.Context) error { return nil } }
func (s *fakeStep) main() common.Executor      { return func(context.Context) error { return nil } }
func (s *fakeStep) post() common.Executor      { return func(context.Context) error { return nil } }
func (s *fakeStep) getRunContext() *RunContext { return s.rc }
func (s *fakeStep) getGithubContext(ctx context.Context) *model.GithubContext {
	return s.rc.getGithubContext(ctx)
}
func (s *fakeStep) getStepModel() *model.Step                         { return s.model }
func (s *fakeStep) getEnv() *map[string]string                        { return &s.env }
func (s *fakeStep) getIfExpression(context.Context, stepStage) string { return s.ifExpr }

// A refused ::set-env::/::add-path:: records a job-scoped error. When the step that
// produced it also fails on its own, the refusal must be cleared at the step boundary, so
// it fails only that step and never leaks onto a later step that runs anyway (if: always()).
func TestRunStepExecutorDoesNotLeakRefusalToNextStep(t *testing.T) {
	cm := &containerMock{}
	noop := func(context.Context) error { return nil }
	cm.On("Copy", mock.Anything, mock.Anything).Return(noop)
	cm.On("UpdateFromEnv", mock.Anything, mock.Anything).Return(noop)

	disallow := false
	rc := &RunContext{
		Config: &Config{AllowUnsecureCommands: &disallow, Env: map[string]string{}},
		Run: &model.Run{
			JobID:    "1",
			Workflow: &model.Workflow{Jobs: map[string]*model.Job{"1": {}}},
		},
		Env:          map[string]string{},
		StepResults:  map[string]*model.StepResult{},
		JobContainer: cm,
	}
	rc.ExprEval = rc.NewExpressionEvaluator(context.Background())
	// Dryrun skips reading the path file back from the (mocked) container.
	ctx := common.WithDryrun(context.Background(), true)

	// Step A refuses a ::set-env:: and then fails on its own.
	stepA := &fakeStep{rc: rc, model: &model.Step{ID: "a"}, env: map[string]string{}}
	errA := runStepExecutor(stepA, stepStageMain, func(context.Context) error {
		rc.commandHandler(ctx)("::set-env name=x::y\n")
		return errors.New("boom")
	})(ctx)
	// The step fails with its own error, not the refusal.
	require.ErrorContains(t, errA, "boom")

	// Step B runs despite step A's failure (if: always()) and issues no unsecure command;
	// it must not inherit step A's refusal.
	stepB := &fakeStep{rc: rc, model: &model.Step{ID: "b"}, env: map[string]string{}, ifExpr: "always()"}
	errB := runStepExecutor(stepB, stepStageMain, func(context.Context) error { return nil })(ctx)
	require.NoError(t, errB)
}
