// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// Package checkout implements the `checkout` built-in action: a strict
// subset of `actions/checkout` v4 that runs natively in the runner
// process and copies the working tree into the job container instead of
// cloning the action and shelling into Node.
//
// See docs/builtin-actions.md.
package checkout

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gitea.com/gitea/runner/act/common/git"
	"gitea.com/gitea/runner/act/runner/builtin"

	"gitea.dev/actionslib/pkg/model"
)

// Action is the `checkout` built-in.
type Action struct{}

func init() {
	builtin.Register(&Action{})
}

func (*Action) Name() string  { return "checkout" }
func (*Action) HasPre() bool  { return false }
func (*Action) HasPost() bool { return false }

func (*Action) Pre(_ *builtin.Context) error  { return nil }
func (*Action) Post(_ *builtin.Context) error { return nil }

func (a *Action) Main(c *builtin.Context) error {
	if c == nil || c.RunContext == nil {
		return errors.New("checkout: nil context")
	}

	ghc := c.RunContext.GithubContext()
	if ghc == nil {
		return errors.New("checkout: no github context available")
	}
	cfg := c.RunContext.Config()

	repository := firstNonEmpty(c.Inputs["repository"], ghc.Repository)
	ref := shortRef(firstNonEmpty(c.Inputs["ref"], ghc.Ref))
	token := firstNonEmpty(c.Inputs["token"], ghc.Token, cfg.Token())
	subPath := strings.TrimSpace(c.Inputs["path"])

	c.Logger.Infof("checkout: repository=%s ref=%s path=%s", repository, ref, subPath)

	// Skip-if-bound parity: when the workflow's `with:` does not redirect
	// to a different repo or ref, the runner already has the workspace on
	// the host (the same code path that `step_action_remote.go` takes for
	// `actions/checkout`). For `Config.BindWorkdir == true` the workspace
	// is already mounted into the container; otherwise copy it in.
	if isLocalCheckoutLike(c.Inputs, ghc) && !cfg.NoSkipCheckout() {
		if cfg.BindWorkdir() {
			c.Logger.Debugf("checkout: workspace already bind-mounted, skipping")
			return nil
		}
		jobContainer := c.RunContext.JobContainer()
		if jobContainer == nil {
			return errors.New("checkout: no job container available")
		}
		dst := path.Join(c.RunContext.ContainerWorkdir(), subPath)
		src := c.RunContext.Workdir() + string(filepath.Separator) + "."
		c.Logger.Debugf("checkout: copying workspace %s -> %s", src, dst)
		return jobContainer.CopyDir(dst, src, cfg.UseGitIgnore())(c.Ctx)
	}

	// Otherwise clone on the runner host into ActionCacheDir/<step-hash>,
	// then copy the working tree into the container workspace at `path`.
	if repository == "" {
		return errors.New("checkout: empty `repository` and no github.repository in context")
	}
	if ref == "" {
		return errors.New("checkout: empty `ref` and no github.ref in context")
	}

	cloneURL := repoURL(ghc, cfg, repository)
	cloneDir := stageDir(c.RunContext.ActionCacheDir(), c.Step.UsesHash(), repository)

	if err := os.MkdirAll(filepath.Dir(cloneDir), 0o755); err != nil {
		return fmt.Errorf("checkout: prepare stage dir: %w", err)
	}

	c.Logger.Infof("checkout: clone %s ref=%s -> %s", cloneURL, ref, cloneDir)
	clone := git.NewGitCloneExecutor(git.NewGitCloneExecutorInput{
		URL:             cloneURL,
		Ref:             ref,
		Dir:             cloneDir,
		Token:           token,
		OfflineMode:     cfg.ActionOfflineMode(),
		InsecureSkipTLS: cfg.InsecureSkipTLS(),
	})
	if err := clone(c.Ctx); err != nil {
		return fmt.Errorf("checkout: clone %s @ %s: %w", cloneURL, ref, err)
	}

	jobContainer := c.RunContext.JobContainer()
	if jobContainer == nil {
		return errors.New("checkout: no job container available")
	}
	dst := path.Join(c.RunContext.ContainerWorkdir(), subPath)
	src := cloneDir + string(filepath.Separator) + "."
	c.Logger.Debugf("checkout: copying clone %s -> %s", src, dst)
	return jobContainer.CopyDir(dst, src, cfg.UseGitIgnore())(c.Ctx)
}

// isLocalCheckoutLike mirrors the precondition used by
// `runner.isLocalCheckout` for `actions/checkout`: the step refers to the
// same repository / ref as the job's github context (or omits them
// entirely), so the working tree the runner already has on disk is fine.
func isLocalCheckoutLike(inputs map[string]string, ghc *model.GithubContext) bool {
	if repo, ok := inputs["repository"]; ok && repo != "" && repo != ghc.Repository {
		return false
	}
	if ref, ok := inputs["ref"]; ok && ref != "" && ref != ghc.Ref {
		return false
	}
	return true
}

// repoURL builds the clone URL. Mirrors what the existing remote-action
// path does: prefer the github context's ServerURL; fall back to the
// configured DefaultActionInstance, then GitHubInstance.
func repoURL(ghc *model.GithubContext, cfg builtin.Config, repository string) string {
	base := strings.TrimRight(ghc.ServerURL, "/")
	if base == "" {
		base = strings.TrimRight(cfg.DefaultActionInstance(), "/")
		if base != "" && !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
			base = "https://" + base
		}
	}
	if base == "" {
		inst := strings.TrimRight(cfg.GitHubInstance(), "/")
		if inst != "" && !strings.HasPrefix(inst, "http://") && !strings.HasPrefix(inst, "https://") {
			inst = "https://" + inst
		}
		base = inst
	}
	return base + "/" + repository
}

// stageDir returns a host-side directory unique to (step, repository) where
// the clone is staged before being copied into the container.
func stageDir(actionCacheDir, usesHash, repository string) string {
	suffix := sha256.Sum256([]byte(repository))
	return filepath.Join(actionCacheDir, "builtin-checkout", usesHash, hex.EncodeToString(suffix[:8]))
}

// shortRef strips a fully-qualified ref (as found in `github.ref`, e.g.
// "refs/heads/main") down to the short name that git.NewGitCloneExecutor
// expects, since it unconditionally re-qualifies its Ref input as
// "refs/heads/<ref>".
func shortRef(ref string) string {
	for _, prefix := range []string{"refs/heads/", "refs/tags/"} {
		if s, ok := strings.CutPrefix(ref, prefix); ok {
			return s
		}
	}
	return ref
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
