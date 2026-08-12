// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2022 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"fmt"

	"gitea.com/gitea/runner/act/runner/builtin"

	"gitea.dev/actionslib/pkg/model"
)

type stepFactory interface {
	newStep(step *model.Step, rc *RunContext) (step, error)
}

type stepFactoryImpl struct{}

func (sf *stepFactoryImpl) newStep(stepModel *model.Step, rc *RunContext) (step, error) {
	switch stepModel.Type() {
	case model.StepTypeInvalid:
		return nil, fmt.Errorf("Invalid run/uses syntax for job:%s step:%+v", rc.Run, stepModel)
	case model.StepTypeRun:
		return &stepRun{
			Step:       stepModel,
			RunContext: rc,
		}, nil
	case model.StepTypeUsesActionLocal:
		return &stepActionLocal{
			Step:       stepModel,
			RunContext: rc,
			readAction: readActionImpl,
			runAction:  runActionImpl,
		}, nil
	case model.StepTypeUsesActionRemote:
		// For Gitea: reserved no-slash names dispatch to a built-in registered
		// in the act/runner/builtin tree. See docs/builtin-actions.md.
		// Classification used to be model.StepTypeUsesActionBuiltin in act/model;
		// that package now lives in actionslib without the builtin type, so we
		// detect the same pattern here.
		if isBuiltinUses(stepModel.Uses) {
			a, ok := builtin.Lookup(stepModel.Uses)
			if !ok {
				return nil, fmt.Errorf("unknown built-in action %q (known: %v) for job:%s step:%+v",
					stepModel.Uses, builtin.Names(), rc.Run, stepModel)
			}
			return &stepActionBuiltin{
				Step:       stepModel,
				RunContext: rc,
				action:     a,
			}, nil
		}
		return &stepActionRemote{
			Step:       stepModel,
			RunContext: rc,
			readAction: readActionImpl,
			runAction:  runActionImpl,
		}, nil
	case model.StepTypeUsesDockerURL:
		return &stepDocker{
			Step:       stepModel,
			RunContext: rc,
		}, nil
	}

	return nil, fmt.Errorf("Unable to determine how to run job:%s step:%+v", rc.Run, stepModel)
}

// isBuiltinUses reports whether uses is a reserved no-slash name the runner
// resolves via its built-in action registry (`^[a-z0-9][a-z0-9-]*$`).
func isBuiltinUses(uses string) bool {
	if uses == "" {
		return false
	}
	for i, r := range uses {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '-':
			if i == 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}
