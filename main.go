// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitea.com/gitea/runner/internal/app/cmd"
	"gitea.com/gitea/runner/internal/pkg/process"
)

func main() {
	if code, ok := process.RunAsInit(); ok {
		os.Exit(code)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// run the command
	cmd.Execute(ctx)
}
