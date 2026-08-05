// Copyright 2021 The Gitea Authors. All rights reserved.
// Copyright 2021 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package model

type JobContext struct {
	Status    string                `json:"status"`
	Container JobContainerContext   `json:"container"`
	Services  map[string]JobService `json:"services"`
}

// JobContainerContext describes the container a job's steps run in. Both fields are
// empty when the job runs on the host instead of in a container.
type JobContainerContext struct {
	ID      string `json:"id"`
	Network string `json:"network"`
}

// JobService describes one of a job's service containers. Ports maps a container port
// to the host port it was published on, which is how a workflow finds it.
type JobService struct {
	ID      string            `json:"id"`
	Network string            `json:"network"`
	Ports   map[string]string `json:"ports"`
}
