// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"gitea.com/gitea/runner/act/common"

	"gitea.dev/actionslib/pkg/model"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// spanRecorder collects ended spans via the sdktrace global provider.
type spanRecorder struct {
	mu    sync.Mutex
	spans []sdktrace.ReadOnlySpan
}

func (r *spanRecorder) OnStart(context.Context, sdktrace.ReadWriteSpan) {}
func (r *spanRecorder) OnEnd(s sdktrace.ReadOnlySpan) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.spans = append(r.spans, s)
}
func (r *spanRecorder) Shutdown(context.Context) error   { return nil }
func (r *spanRecorder) ForceFlush(context.Context) error { return nil }

func (r *spanRecorder) recorded() []sdktrace.ReadOnlySpan {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]sdktrace.ReadOnlySpan(nil), r.spans...)
}

func attrValue(attrs []attribute.KeyValue, key string) (string, bool) {
	for _, kv := range attrs {
		if string(kv.Key) == key {
			return kv.Value.String(), true
		}
	}
	return "", false
}

func TestNewJobExecutorCreatesStepSpans(t *testing.T) {
	rec := &spanRecorder{}
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	ctx := common.WithJobErrorContainer(context.Background())
	jim := &jobInfoMock{}
	sfm := &stepFactoryMock{}
	rc := &RunContext{
		JobContainer: &jobContainerMock{},
		Run: &model.Run{
			JobID: "test",
			Workflow: &model.Workflow{
				Jobs: map[string]*model.Job{"test": {}},
			},
		},
		Config: &Config{TaskID: 42},
	}
	rc.ExprEval = rc.NewExpressionEvaluator(ctx)

	steps := []*model.Step{{ID: "1", Name: "build"}, {ID: "2", Name: "test"}}
	executed := make([]string, 0)

	jim.On("steps").Return(steps)
	jim.On("matrix").Return(map[string]any{})
	jim.On("startContainer").Return(func(context.Context) error { return nil })
	jim.On("stopContainer").Return(func(context.Context) error { return nil })
	jim.On("closeContainer").Return(func(context.Context) error { return nil })
	jim.On("interpolateOutputs").Return(func(context.Context) error { return nil })
	jim.On("result", "success")

	for _, stepModel := range steps {
		sm := &stepMock{}
		sfm.On("newStep", stepModel, rc).Return(sm, nil)
		sm.On("pre").Return(func(context.Context) error { return nil })
		sm.On("main").Return(func(ctx context.Context) error {
			executed = append(executed, stepModel.Name)
			return nil
		})
		sm.On("post").Return(func(context.Context) error { return nil })
		defer sm.AssertExpectations(t)
	}

	logger, _ := logrustest.NewNullLogger()
	require.NoError(t, newJobExecutor(jim, sfm, rc)(common.WithLogger(ctx, logger.WithField("job", "test"))))

	assert.Equal(t, []string{"build", "test"}, executed)

	// The docker client emits its own spans (moby instruments its transport),
	// so pick out the step spans by name and preserve their execution order.
	spans := rec.recorded()
	stepSpans := make([]sdktrace.ReadOnlySpan, 0, len(steps))
	for _, span := range spans {
		if span.Name() == "build" || span.Name() == "test" {
			stepSpans = append(stepSpans, span)
		}
	}
	require.Len(t, stepSpans, 2)
	for i, span := range stepSpans {
		assert.Equal(t, steps[i].Name, span.Name())
		assert.Equal(t, trace.SpanKindInternal, span.SpanKind())
		name, ok := attrValue(span.Attributes(), "cicd.pipeline.task.name")
		require.True(t, ok)
		assert.Equal(t, steps[i].Name, name)
		runID, ok := attrValue(span.Attributes(), "cicd.pipeline.task.run.id")
		require.True(t, ok)
		assert.Equal(t, fmt.Sprintf("42.%d", steps[i].Number), runID)
		result, ok := attrValue(span.Attributes(), "cicd.pipeline.task.run.result")
		require.True(t, ok)
		assert.Equal(t, "success", result)
	}
}
