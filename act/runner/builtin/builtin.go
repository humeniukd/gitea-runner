// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// Package builtin hosts the runner's built-in actions and the small
// registry that maps reserved no-slash `uses:` names to native Go
// implementations.
//
// See docs/builtin-actions.md for the design overview.
package builtin

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"gitea.com/gitea/runner/act/container"
	"gitea.dev/actionslib/pkg/model"

	log "github.com/sirupsen/logrus"
)

// Action is implemented by every built-in. It deliberately does not depend on
// the `act/runner` package so the registry can live below it without an
// import cycle.
//
// Implementations should be cheap to construct (typically one zero-value
// struct per process) and should be safe for concurrent use across jobs.
type Action interface {
	// Name is the reserved no-slash token a workflow uses in `uses:`.
	// Must match `^[a-z0-9][a-z0-9-]*$`.
	Name() string

	// HasPre / HasPost let the runner skip dispatching to stages the
	// built-in does not implement. Pre / Post are called only when these
	// return true.
	HasPre() bool
	HasPost() bool

	Pre(*Context) error
	Main(*Context) error
	Post(*Context) error
}

// Stage tags Pre/Main/Post so a single helper can dispatch.
type Stage int

const (
	StagePre Stage = iota
	StageMain
	StagePost
)

func (s Stage) String() string {
	switch s {
	case StagePre:
		return "pre"
	case StageMain:
		return "main"
	case StagePost:
		return "post"
	}
	return "unknown"
}

// Context is everything a built-in needs to run a single stage. The runner
// allocates one per (step, stage) and discards it afterwards.
type Context struct {
	Ctx        context.Context
	RunContext RunContext

	// Step is the workflow step model. Built-ins should treat it as
	// read-only.
	Step *model.Step

	// Inputs holds the step's `with:` keys after expression interpolation.
	// It is a separate map so built-ins do not need to filter INPUT_*
	// out of Env.
	Inputs map[string]string

	// Env is the live step environment. The runner reads it back after
	// the stage returns to drive `GITHUB_OUTPUT`/`GITHUB_ENV`/etc., so
	// built-ins may add or replace keys here.
	Env map[string]string

	Logger log.FieldLogger
}

// RunContext is the narrowed view of `*runner.RunContext` exposed to
// built-ins. Keeping it as an interface means this package depends only on
// `act/model` and `act/container`, not on `act/runner`.
type RunContext interface {
	// JobContainer is the execution environment for the job (Docker
	// container, host environment, etc.). Built-ins drive the container
	// only through this interface.
	JobContainer() container.ExecutionsEnvironment

	// Workdir is the host-side workspace path (the directory the runner
	// considers the job's workspace on the runner machine).
	Workdir() string

	// ContainerWorkdir is the workspace path as seen from inside the job
	// container; equivalent to `JobContainer().ToContainerPath(Workdir())`.
	ContainerWorkdir() string

	// ActionCacheDir is a host-side directory where built-ins may stage
	// per-step working data (clones, downloads, ...). The runner does
	// not require it to be empty between steps; built-ins should
	// namespace inside it (e.g. by `Step.UsesHash()`).
	ActionCacheDir() string

	// GithubContext returns the live github context for the job (used to
	// resolve default repository / ref / token).
	GithubContext() *model.GithubContext

	// Config exposes a narrowed read-only view of the runner config.
	Config() Config
}

// Config is the subset of `*runner.Config` that built-ins may read.
// It is intentionally narrow: any new knob a built-in needs should be
// surfaced explicitly here so the contract remains visible.
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

// Registry is the package-private store for built-ins.
type registry struct {
	mu  sync.RWMutex
	all map[string]Action
}

var reg = &registry{all: map[string]Action{}}

// Register adds an Action to the package-wide registry. It panics on
// duplicate names because registration always happens at process start
// from `init()` and a duplicate is a programmer error rather than a
// runtime condition.
func Register(a Action) {
	if a == nil {
		panic("builtin: Register called with nil Action")
	}
	name := a.Name()
	if !validName(name) {
		panic(fmt.Sprintf("builtin: invalid Action name %q (must match ^[a-z0-9][a-z0-9-]*$)", name))
	}
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if _, ok := reg.all[name]; ok {
		panic(fmt.Sprintf("builtin: Action %q already registered", name))
	}
	reg.all[name] = a
}

// Lookup returns the Action registered under name, if any.
func Lookup(name string) (Action, bool) {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	a, ok := reg.all[name]
	return a, ok
}

// Names returns all registered names, sorted, for diagnostics.
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

// validName mirrors the regex used by `model.(*Step).Type()` to recognize
// a built-in. Kept in code (not regexp) to avoid a one-shot compile per
// registration.
func validName(s string) bool {
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

// resetForTest empties the registry. It is exported via builtin_test.go's
// access; in production code Register is the only mutator.
func resetForTest() {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.all = map[string]Action{}
}

// DispatchStage is a convenience used by `stepActionBuiltin` to invoke
// the right hook with a single switch.
func DispatchStage(a Action, c *Context, s Stage) error {
	switch s {
	case StagePre:
		return a.Pre(c)
	case StageMain:
		return a.Main(c)
	case StagePost:
		return a.Post(c)
	}
	return fmt.Errorf("builtin: unknown stage %d", int(s))
}
