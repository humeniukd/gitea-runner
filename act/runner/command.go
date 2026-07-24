// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"gitea.com/gitea/runner/act/common"
)

var commandPatternGA *regexp.Regexp

var commandPatternADO *regexp.Regexp

func init() {
	commandPatternGA = regexp.MustCompile("^::([^ ]+)( (.+))?::([^\r\n]*)[\r\n]+$")
	commandPatternADO = regexp.MustCompile("^##\\[([^ ]+)( (.+))?]([^\r\n]*)[\r\n]+$")
}

func tryParseRawActionCommand(line string) (command string, kvPairs map[string]string, arg string, ok bool) {
	if m := commandPatternGA.FindStringSubmatch(line); m != nil {
		command = m[1]
		kvPairs = parseKeyValuePairs(m[3], ",")
		arg = m[4]
		ok = true
	} else if m := commandPatternADO.FindStringSubmatch(line); m != nil {
		command = m[1]
		kvPairs = parseKeyValuePairs(m[3], ";")
		arg = m[4]
		ok = true
	}
	return command, kvPairs, arg, ok
}

func (rc *RunContext) commandHandler(ctx context.Context) common.LineHandler {
	logger := common.Logger(ctx)
	resumeCommand := ""
	return func(line string) bool {
		command, kvPairs, arg, ok := tryParseRawActionCommand(line)
		if !ok {
			return true
		}

		if resumeCommand != "" && command != resumeCommand {
			// There should not be any emojis in the log output for Gitea.
			// The code in the switch statement is the same.
			// Return true (not false) so the line still reaches the raw_output
			// log handler; otherwise everything between ::stop-commands:: and
			// its end token is silently dropped from the step log.
			logger.Infof("%s", line)
			return true
		}
		arg = UnescapeCommandData(arg)
		kvPairs = unescapeKvPairs(kvPairs)
		switch command {
		case "set-env":
			if rc.refuseUnsecureCommand(ctx, command) {
				return true
			}
			rc.setEnv(ctx, kvPairs, arg)
		case "set-output":
			rc.setOutput(ctx, kvPairs, arg)
		case "add-path":
			if rc.refuseUnsecureCommand(ctx, command) {
				return true
			}
			rc.addPath(ctx, arg)
		case "debug":
			logger.Infof("%s", line)
		case "warning":
			logger.Infof("%s", line)
		case "error":
			logger.Infof("%s", line)
		case "add-mask":
			rc.AddMask(arg)
			logger.Infof("%s", "***")
		case "stop-commands":
			resumeCommand = arg
			logger.Infof("%s", line)
		case resumeCommand:
			resumeCommand = ""
			logger.Infof("%s", line)
		case "save-state":
			logger.Infof("%s", line)
			rc.saveState(ctx, kvPairs, arg)
		case "add-matcher":
			logger.Infof("%s", line)
		default:
			logger.Infof("%s", line)
		}

		// return true to let gitea's logger handle these special outputs also
		return true
	}
}

// refuseUnsecureCommand reports whether a deprecated ::set-env:: or ::add-path:: command
// must not run, recording the error that fails the step. GitHub disabled both because a
// step that echoes untrusted content can use them to set variables such as NODE_OPTIONS
// or PATH for every later step of the job; opting back in requires the step or job to set
// ACTIONS_ALLOW_UNSECURE_COMMANDS.
func (rc *RunContext) refuseUnsecureCommand(ctx context.Context, command string) bool {
	if rc.allowUnsecureCommands() {
		rc.warnUnsecureCommandOnce(ctx, command)
		return false
	}

	err := fmt.Errorf("the `%s` workflow command is disabled: it can set the environment of every later step from untrusted output. "+
		"Write to $GITHUB_ENV or $GITHUB_PATH instead, or set ACTIONS_ALLOW_UNSECURE_COMMANDS to allow it", command)
	common.Logger(ctx).WithField(rawOutputField, true).Errorf("##[error]%v", err)

	rc.unsecureCommandMu.Lock()
	defer rc.unsecureCommandMu.Unlock()
	if rc.unsecureCommandErr == nil {
		rc.unsecureCommandErr = err
	}
	return true
}

// allowUnsecureCommands reports whether ::set-env:: and ::add-path:: may still run. The
// step's own environment is checked first, since that is where GitHub's documentation
// puts ACTIONS_ALLOW_UNSECURE_COMMANDS.
func (rc *RunContext) allowUnsecureCommands() bool {
	const allowVar = "ACTIONS_ALLOW_UNSECURE_COMMANDS"

	if env := rc.currentStepEnv(); env != nil && isTruthyEnv((*env)[allowVar]) {
		return true
	}
	if isTruthyEnv(rc.Env[allowVar]) || isTruthyEnv(rc.GlobalEnv[allowVar]) {
		return true
	}
	// Not opted in: the runner setting decides, and still allows them by default.
	return rc.Config == nil || rc.Config.AllowUnsecureCommands == nil || *rc.Config.AllowUnsecureCommands
}

func isTruthyEnv(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "y", "yes", "on":
		return true
	default:
		return false
	}
}

// warnUnsecureCommandOnce warns the first time a job runs a deprecated command without
// opting in, so the log says what will stop working once the setting flips.
func (rc *RunContext) warnUnsecureCommandOnce(ctx context.Context, command string) {
	if rc.allowUnsecureCommandsOptIn() {
		return
	}

	rc.unsecureCommandMu.Lock()
	defer rc.unsecureCommandMu.Unlock()
	if rc.unsecureCommandWarned {
		return
	}
	rc.unsecureCommandWarned = true
	common.Logger(ctx).WithField(rawOutputField, true).Warnf(
		"##[warning]The `%s` workflow command is deprecated and will stop working. "+
			"Write to $GITHUB_ENV or $GITHUB_PATH instead, or set ACTIONS_ALLOW_UNSECURE_COMMANDS to keep using it", command)
}

// allowUnsecureCommandsOptIn reports whether the workflow itself asked for the deprecated
// commands, in which case there is nothing to warn about.
func (rc *RunContext) allowUnsecureCommandsOptIn() bool {
	const allowVar = "ACTIONS_ALLOW_UNSECURE_COMMANDS"

	if env := rc.currentStepEnv(); env != nil && isTruthyEnv((*env)[allowVar]) {
		return true
	}
	return isTruthyEnv(rc.Env[allowVar]) || isTruthyEnv(rc.GlobalEnv[allowVar])
}

// takeUnsecureCommandError returns and clears the error left by a refused command, so the
// step that produced it fails the way GitHub fails it.
func (rc *RunContext) takeUnsecureCommandError() error {
	rc.unsecureCommandMu.Lock()
	defer rc.unsecureCommandMu.Unlock()
	err := rc.unsecureCommandErr
	rc.unsecureCommandErr = nil
	return err
}

func (rc *RunContext) setEnv(ctx context.Context, kvPairs map[string]string, arg string) {
	name := kvPairs["name"]
	common.Logger(ctx).Infof("::set-env:: %s=%s", name, arg)
	if rc.Env == nil {
		rc.Env = make(map[string]string)
	}
	if rc.GlobalEnv == nil {
		rc.GlobalEnv = map[string]string{}
	}
	newenv := map[string]string{
		name: arg,
	}
	mergeIntoMap := mergeIntoMapCaseSensitive
	if rc.JobContainer != nil && rc.JobContainer.IsEnvironmentCaseInsensitive() {
		mergeIntoMap = mergeIntoMapCaseInsensitive
	}
	mergeIntoMap(rc.Env, newenv)
	mergeIntoMap(rc.GlobalEnv, newenv)
}

func (rc *RunContext) setOutput(ctx context.Context, kvPairs map[string]string, arg string) {
	logger := common.Logger(ctx)
	stepID := rc.CurrentStep
	outputName := kvPairs["name"]
	if outputMapping, ok := rc.OutputMappings[MappableOutput{StepID: stepID, OutputName: outputName}]; ok {
		stepID = outputMapping.StepID
		outputName = outputMapping.OutputName
	}

	result, ok := rc.StepResults[stepID]
	if !ok {
		logger.Infof("No outputs registered for step '%s'", stepID)
		return
	}

	logger.Infof("::set-output:: %s=%s", outputName, arg)
	result.Outputs[outputName] = arg
}

func (rc *RunContext) addPath(ctx context.Context, arg string) {
	common.Logger(ctx).Infof("::add-path:: %s", arg)
	extraPath := []string{arg}
	for _, v := range rc.ExtraPath {
		if v != arg {
			extraPath = append(extraPath, v)
		}
	}
	rc.ExtraPath = extraPath
}

func parseKeyValuePairs(kvPairs, separator string) map[string]string {
	rtn := make(map[string]string)
	kvPairList := strings.SplitSeq(kvPairs, separator)
	for kvPair := range kvPairList {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			rtn[kv[0]] = kv[1]
		}
	}
	return rtn
}

func UnescapeCommandData(arg string) string {
	escapeMap := map[string]string{
		"%25": "%",
		"%0D": "\r",
		"%0A": "\n",
	}
	for k, v := range escapeMap {
		arg = strings.ReplaceAll(arg, k, v)
	}
	return arg
}

// UnescapeCommandProperty decodes a workflow command property value, whose escaping
// covers the separators of the property list as well.
func UnescapeCommandProperty(arg string) string {
	escapeMap := map[string]string{
		"%25": "%",
		"%0D": "\r",
		"%0A": "\n",
		"%3A": ":",
		"%2C": ",",
	}
	for k, v := range escapeMap {
		arg = strings.ReplaceAll(arg, k, v)
	}
	return arg
}

func unescapeKvPairs(kvPairs map[string]string) map[string]string {
	for k, v := range kvPairs {
		kvPairs[k] = UnescapeCommandProperty(v)
	}
	return kvPairs
}

func (rc *RunContext) saveState(_ context.Context, kvPairs map[string]string, arg string) {
	stepID := rc.CurrentStep
	if stepID != "" {
		if rc.IntraActionState == nil {
			rc.IntraActionState = map[string]map[string]string{}
		}
		state, ok := rc.IntraActionState[stepID]
		if !ok {
			state = map[string]string{}
			rc.IntraActionState[stepID] = state
		}
		state[kvPairs["name"]] = arg
	}
}
