package common

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWorkflow(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// empty
	emptyWorkflow := NewPipelineExecutor()
	require.NoError(t, emptyWorkflow(ctx))

	// error case
	errorWorkflow := NewErrorExecutor(errors.New("test error"))
	require.Error(t, errorWorkflow(ctx))

	// multiple success case
	runcount := 0
	successWorkflow := NewPipelineExecutor(
		func(_ context.Context) error {
			runcount++
			return nil
		},
		func(_ context.Context) error {
			runcount++
			return nil
		})
	require.NoError(t, successWorkflow(ctx))
	assert.Equal(2, runcount)
}

func TestNewConditionalExecutor(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	trueCount := 0
	falseCount := 0

	err := NewConditionalExecutor(func(_ context.Context) bool {
		return false
	}, func(_ context.Context) error {
		trueCount++
		return nil
	}, func(_ context.Context) error {
		falseCount++
		return nil
	})(ctx)

	require.NoError(t, err)
	assert.Equal(0, trueCount)
	assert.Equal(1, falseCount)

	err = NewConditionalExecutor(func(_ context.Context) bool {
		return true
	}, func(_ context.Context) error {
		trueCount++
		return nil
	}, func(_ context.Context) error {
		falseCount++
		return nil
	})(ctx)

	require.NoError(t, err)
	assert.Equal(1, trueCount)
	assert.Equal(1, falseCount)
}

func TestNewParallelExecutor(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	var count atomic.Int32
	var activeCount atomic.Int32
	var maxCount atomic.Int32
	emptyWorkflow := NewPipelineExecutor(func(_ context.Context) error {
		count.Add(1)

		cur := activeCount.Add(1)
		for {
			old := maxCount.Load()
			if cur <= old || maxCount.CompareAndSwap(old, cur) {
				break
			}
		}
		time.Sleep(2 * time.Second)
		activeCount.Add(-1)

		return nil
	})

	err := NewParallelExecutor(2, emptyWorkflow, emptyWorkflow, emptyWorkflow)(ctx)

	assert.Equal(int32(3), count.Load(), "should run all 3 executors")
	assert.Equal(int32(2), maxCount.Load(), "should run at most 2 executors in parallel")
	require.NoError(t, err)

	// Reset to test running the executor with 0 parallelism
	count.Store(0)
	activeCount.Store(0)
	maxCount.Store(0)

	errSingle := NewParallelExecutor(0, emptyWorkflow, emptyWorkflow, emptyWorkflow)(ctx)

	assert.Equal(int32(3), count.Load(), "should run all 3 executors")
	assert.Equal(int32(1), maxCount.Load(), "should run at most 1 executors in parallel")
	require.NoError(t, errSingle)
}

func TestNewParallelExecutorFailed(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	count := 0
	errorWorkflow := NewPipelineExecutor(func(_ context.Context) error {
		count++
		return errors.New("fake error")
	})
	err := NewParallelExecutor(1, errorWorkflow)(ctx)
	assert.Equal(1, count)
	assert.ErrorIs(err, context.Canceled)
}

func TestNewParallelExecutorCanceled(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var count atomic.Int32
	successWorkflow := NewPipelineExecutor(func(_ context.Context) error {
		count.Add(1)
		return nil
	})
	errorWorkflow := NewPipelineExecutor(func(_ context.Context) error {
		count.Add(1)
		return errors.New("fake error")
	})
	err := NewParallelExecutor(3, errorWorkflow, successWorkflow, successWorkflow)(ctx)
	assert.Equal(int32(3), count.Load())
	assert.ErrorIs(err, context.Canceled)
}
