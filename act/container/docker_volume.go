// Copyright 2023 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !(WITHOUT_DOCKER || !(linux || darwin || windows || netbsd))

package container

import (
	"context"

	"gitea.com/gitea/runner/act/common"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/client"
)

func NewDockerVolumeRemoveExecutor(volumeName string, force bool) common.Executor {
	return func(ctx context.Context) error {
		common.Logger(ctx).Debugf("docker volume rm %s", volumeName)

		if common.Dryrun(ctx) {
			return nil
		}

		cli, err := GetDockerClient(ctx)
		if err != nil {
			return err
		}
		defer cli.Close()

		_, err = cli.VolumeRemove(ctx, volumeName, client.VolumeRemoveOptions{Force: force})
		if cerrdefs.IsNotFound(err) { // already gone is the outcome we wanted
			return nil
		}
		return err
	}
}
