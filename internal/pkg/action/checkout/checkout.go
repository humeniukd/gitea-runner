// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checkout

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"gitea.com/gitea/runner/act/common/git"
	"gitea.com/gitea/runner/internal/pkg/action"

	"gitea.dev/actionslib/pkg/model"
)

type Action struct{}

func init() {
	action.Register(&Action{})
}

func (*Action) Name() string  { return "checkout" }
func (*Action) HasPre() bool  { return false }
func (*Action) HasPost() bool { return false }

func (*Action) Pre(_ *action.Context) error  { return nil }
func (*Action) Post(_ *action.Context) error { return nil }

func (a *Action) Main(c *action.Context) error {
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
	dst, err := workspaceDestination(c.RunContext.ContainerWorkdir(), subPath)
	if err != nil {
		return err
	}
	depth, err := fetchDepth(c.Inputs["fetch-depth"])
	if err != nil {
		return err
	}

	c.Logger.Infof("checkout: repository=%s ref=%s path=%s", repository, ref, subPath)

	if isLocalCheckoutLike(c.Inputs, ghc) && !cfg.NoSkipCheckout() {
		if cfg.BindWorkdir() {
			c.Logger.Debugf("checkout: workspace already bind-mounted, skipping")
			return nil
		}
		jobContainer := c.RunContext.JobContainer()
		if jobContainer == nil {
			return errors.New("checkout: no job container available")
		}
		src := c.RunContext.Workdir() + string(filepath.Separator) + "."
		c.Logger.Debugf("checkout: copying workspace %s -> %s", src, dst)
		return jobContainer.CopyDir(dst, src, cfg.UseGitIgnore())(c.Ctx)
	}

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
		Depth:           depth,
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
	src := cloneDir + string(filepath.Separator) + "."
	c.Logger.Debugf("checkout: copying clone %s -> %s", src, dst)
	return jobContainer.CopyDir(dst, src, cfg.UseGitIgnore())(c.Ctx)
}

func workspaceDestination(workspace, checkoutPath string) (string, error) {
	cleanPath := path.Clean(checkoutPath)
	if path.IsAbs(cleanPath) || cleanPath == ".." || strings.HasPrefix(cleanPath, "../") {
		return "", fmt.Errorf("checkout: path %q must be relative to the workspace", checkoutPath)
	}
	return path.Join(workspace, cleanPath), nil
}

func fetchDepth(input string) (int, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 1, nil
	}
	depth, err := strconv.Atoi(input)
	if err != nil || depth < 0 {
		return 0, fmt.Errorf("checkout: fetch-depth %q must be a non-negative integer", input)
	}
	return depth, nil
}

func isLocalCheckoutLike(inputs map[string]string, ghc *model.GithubContext) bool {
	if repo, ok := inputs["repository"]; ok && repo != "" && repo != ghc.Repository {
		return false
	}
	if ref, ok := inputs["ref"]; ok && ref != "" && ref != ghc.Ref {
		return false
	}
	return true
}

func repoURL(ghc *model.GithubContext, cfg action.Config, repository string) string {
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

func stageDir(actionCacheDir, usesHash, repository string) string {
	suffix := sha256.Sum256([]byte(repository))
	return filepath.Join(actionCacheDir, "builtin-checkout", usesHash, hex.EncodeToString(suffix[:8]))
}

// NewGitCloneExecutor re-prefixes the ref as refs/heads/<ref>.
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
