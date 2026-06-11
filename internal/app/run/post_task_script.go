// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"gitea.com/gitea/runner/internal/pkg/report"

	runnerv1 "gitea.dev/actions-proto-go/runner/v1"
	log "github.com/sirupsen/logrus"
)

func (r *Runner) runPostTaskScript(ctx context.Context, reporter *report.Reporter, task *runnerv1.Task, workdir string) {
	script := r.cfg.Runner.PostTaskScript
	if script == "" {
		return
	}

	timeout := r.cfg.Runner.PostTaskScriptTimeout
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}

	scriptCtx, cancel := postTaskScriptContext(ctx, timeout)
	defer cancel()

	env := r.postTaskScriptEnv(reporter, task, workdir)
	log.Infof("running post-task script %q for task %d", script, task.Id)

	cmd := exec.CommandContext(scriptCtx, script)
	cmd.Env = envListFromMap(env)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Warnf("post-task script %q for task %d: stdout pipe: %v", script, task.Id, err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Warnf("post-task script %q for task %d: stderr pipe: %v", script, task.Id, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		writePostTaskScriptOutput("stdout", stdout)
	}()
	go func() {
		defer wg.Done()
		writePostTaskScriptOutput("stderr", stderr)
	}()

	if err := cmd.Start(); err != nil {
		log.Warnf("post-task script %q for task %d: %v", script, task.Id, err)
		return
	}

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			log.Warnf("post-task script %q for task %d: %v", script, task.Id, err)
			return
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			log.Warnf("post-task script %q for task %d exited with code %d", script, task.Id, exitErr.ExitCode())
			return
		}
		log.Warnf("post-task script %q for task %d: %v", script, task.Id, err)
	}
}

func postTaskScriptContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx.Err() != nil {
		return context.WithTimeout(context.Background(), timeout)
	}
	return context.WithTimeout(ctx, timeout)
}

func (r *Runner) postTaskScriptEnv(reporter *report.Reporter, task *runnerv1.Task, workdir string) map[string]string {
	env := r.cloneEnvs()
	maps.Copy(env, map[string]string{
		"GITEA_TASK_ID":    strconv.FormatInt(task.Id, 10),
		"GITEA_WORKSPACE":  workdir,
		"GITEA_JOB_RESULT": jobResultString(reporter.Result()),
	})
	if v := task.Context.Fields["run_id"].GetStringValue(); v != "" {
		env["GITEA_RUN_ID"] = v
	}
	if v := task.Context.Fields["repository"].GetStringValue(); v != "" {
		env["GITEA_REPOSITORY"] = v
	}
	return env
}

func jobResultString(result runnerv1.Result) string {
	switch result {
	case runnerv1.Result_RESULT_SUCCESS:
		return "success"
	case runnerv1.Result_RESULT_FAILURE:
		return "failure"
	case runnerv1.Result_RESULT_CANCELLED:
		return "cancelled"
	case runnerv1.Result_RESULT_SKIPPED:
		return "skipped"
	default:
		return "unknown"
	}
}

func envListFromMap(env map[string]string) []string {
	envList := make([]string, 0, len(env))
	for k, v := range env {
		envList = append(envList, fmt.Sprintf("%s=%s", k, v))
	}
	return envList
}

func writePostTaskScriptOutput(stream string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log.Infof("post-task script %s: %s", stream, scanner.Text())
	}
	if err := scanner.Err(); err != nil && !errors.Is(err, os.ErrClosed) {
		log.Warnf("post-task script %s: read output: %v", stream, err)
	}
}
