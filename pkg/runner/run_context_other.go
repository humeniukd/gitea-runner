//go:build !darwin

package runner

import (
	"context"
	"errors"

	"gitea.com/gitea/act_runner/pkg/common"
)

func (rc *RunContext) startTartEnvironment() common.Executor {
	return func(_ context.Context) error {
		return errors.New("you need macOS for tart")
	}
}
