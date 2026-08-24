// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// Package action is the registry for `uses: builtin:<name>` actions.
package action

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"gitea.com/gitea/runner/act/container"

	"gitea.dev/actionslib/pkg/model"
	log "github.com/sirupsen/logrus"
)

const usesPrefix = "builtin:"

type Action interface {
	Name() string
	HasPre() bool
	HasPost() bool
	Pre(*Context) error
	Main(*Context) error
	Post(*Context) error
}

type Context struct {
	Ctx        context.Context
	RunContext RunContext
	Step       *model.Step
	Inputs     map[string]string
	Env        map[string]string // step env; written back after the stage
	Logger     log.FieldLogger
}

// RunContext is the *runner.RunContext surface built-ins may use (avoids an import cycle).
type RunContext interface {
	JobContainer() container.ExecutionsEnvironment
	Workdir() string
	ContainerWorkdir() string
	ActionCacheDir() string
	GithubContext() *model.GithubContext
	Config() Config
}

type Config interface {
	Token() string
	GitHubInstance() string
	DefaultActionInstance() string
	InsecureSkipTLS() bool
	BindWorkdir() bool
	NoSkipCheckout() bool
	UseGitIgnore() bool
	ActionOfflineMode() bool
}

type registry struct {
	mu  sync.RWMutex
	all map[string]Action
}

var reg = &registry{all: map[string]Action{}}

func Register(a Action) {
	if a == nil {
		panic("builtin: Register called with nil Action")
	}
	name := a.Name()
	if !ValidName(name) {
		panic(fmt.Sprintf("builtin: invalid Action name %q (must match ^[a-z0-9][a-z0-9-]*$)", name))
	}
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if _, ok := reg.all[name]; ok {
		panic(fmt.Sprintf("builtin: Action %q already registered", name))
	}
	reg.all[name] = a
}

func Lookup(name string) (Action, bool) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	a, ok := reg.all[name]
	return a, ok
}

func Names() []string {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	out := make([]string, 0, len(reg.all))
	for k := range reg.all {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func ValidName(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
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

func ParseUses(uses string) (string, bool) {
	name, ok := strings.CutPrefix(uses, usesPrefix)
	if !ok {
		return "", false
	}
	return name, true
}

func resetForTest() {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.all = map[string]Action{}
}
