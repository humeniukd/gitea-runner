//go:build !darwin

package runner

import (
	"context"
	"fmt"

	"gitea.com/gitea/act_runner/pkg/common"
)

func (rc *RunContext) startTartEnvironment() common.Executor {
	return func(_ context.Context) error {
		return fmt.Errorf("you need macOS for tart")
	}
}
