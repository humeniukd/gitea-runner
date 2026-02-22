package main

import (
	_ "embed"

	"gitea.com/gitea/act_runner/cmd"
	"gitea.com/gitea/act_runner/pkg/common"
)

//go:embed VERSION
var version string

func main() {
	ctx, cancel := common.CreateGracefulJobCancellationContext()
	defer cancel()

	// run the command
	cmd.Execute(ctx, version)
}
