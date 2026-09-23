// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gitea.com/gitea/runner/act/container"
)

func GetContainerHookInputFile(event, stateDir string) string {
	return filepath.Join(stateDir, event+"-input.json")
}

func GetContainerHookResponseFile(event, stateDir string) string {
	return filepath.Join(stateDir, event+".json")
}

func ParseEnvFileContent(
	ctx context.Context,
	content string,
	setter func(context.Context, map[string]string, string),
) error {
	parsed, err := container.ParseEnvFileContent(strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("container-hook: parse env file content: %w", err)
	}
	for k, v := range parsed {
		setter(ctx, map[string]string{"name": k}, v)
	}
	return nil
}

func ReadContainerHookResponse[T any](event, path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("container hook %s: read %s: %w", event, path, err)
	}
	var resp T
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("container hook %s returned invalid JSON in %s: %w", event, path, err)
	}
	return &resp, nil
}

// RunContainerHook launches the container-hooks entrypoint for a single event
// and streams its stdout and stderr into logWriter. When hooks are disabled or
// Cmd is empty, it is a no-op. The hook's exit code is returned as an error.
func RunContainerHook(
	ctx context.Context,
	cfg Config,
	event string,
	env []string,
	logWriter io.Writer,
) error {
	if !cfg.Enabled || cfg.Cmd == "" {
		return nil
	}
	var cmd *exec.Cmd
	if cfg.Entrypoint != "" {
		cmd = exec.CommandContext(ctx, cfg.Entrypoint, cfg.Cmd)
	} else {
		cmd = exec.CommandContext(ctx, cfg.Cmd)
	}
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("container hook %s failed: %w", event, err)
	}
	return nil
}
