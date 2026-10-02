// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checkout

import (
	"cmp"
	"context"
	"encoding/base64"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/internal/pkg/action"
)

var inputNames = []string{"repository", "ref", "token", "path", "fetch-depth"}

func Main(ctx context.Context, c *action.Context) error {
	for name := range c.Inputs {
		if !slices.Contains(inputNames, name) {
			return fmt.Errorf("checkout: unsupported input %q", name)
		}
	}
	dst, err := workspaceDestination(c.Workspace, c.Inputs["path"])
	if err != nil {
		return err
	}
	depth, err := fetchDepth(c.Inputs["fetch-depth"])
	if err != nil {
		return err
	}

	ghc := c.Github
	repository := cmp.Or(c.Inputs["repository"], ghc.Repository)
	ref := shortBranch(c.Inputs["ref"])
	isWorkflowRepo := repository == ghc.Repository
	local := isWorkflowRepo && slices.Contains([]string{"", shortBranch(ghc.Ref), ghc.Sha}, ref) && !c.NoSkipCheckout
	if local && !c.BindWorkdir {
		return c.Container.CopyDir(dst, c.HostWorkdir+string(filepath.Separator)+".", c.UseGitIgnore, false)(ctx)
	}
	if local && dst == c.Workspace {
		return nil
	}

	if ref == "" {
		ref = "HEAD"
		if isWorkflowRepo {
			ref = ghc.Sha
		}
	}
	url := strings.TrimSuffix(ghc.ServerURL, "/") + "/" + repository
	if err := clone(ctx, c.Container, url, ref, dst, cmp.Or(c.Inputs["token"], ghc.Token), depth, c.InsecureSkipTLS); err != nil {
		return fmt.Errorf("checkout: clone %s@%s: %w", repository, ref, err)
	}
	return nil
}

func clone(ctx context.Context, env container.ExecutionsEnvironment, url, ref, dst, token string, depth int, insecureSkipTLS bool) error {
	run := func(args ...string) error {
		return env.Exec(args, nil, "", "")(ctx)
	}
	git := func(args ...string) error {
		cmd := append([]string{"git"}, gitConfigArgs(token, insecureSkipTLS)...)
		return run(append(cmd, args...)...)
	}
	if err := run("git", "--version"); err != nil {
		return fmt.Errorf("git is required in the job image: %w", err)
	}
	if err := run("mkdir", "-p", dst); err != nil {
		return err
	}
	if err := git("init", dst); err != nil {
		return err
	}
	if err := git("-C", dst, "remote", "add", "origin", url); err != nil {
		return err
	}
	fetch := []string{"-C", dst, "fetch", "--no-tags"}
	if depth > 0 {
		fetch = append(fetch, "--depth", strconv.Itoa(depth))
	}
	if err := git(append(fetch, "origin", ref)...); err != nil {
		return err
	}
	return git("-C", dst, "checkout", "--force", "FETCH_HEAD")
}

func gitConfigArgs(token string, insecureSkipTLS bool) []string {
	var args []string
	if insecureSkipTLS {
		args = append(args, "-c", "http.sslVerify=false")
	}
	if token != "" {
		auth := base64.StdEncoding.EncodeToString([]byte("token:" + token))
		args = append(args, "-c", "http.extraHeader=AUTHORIZATION: basic "+auth)
	}
	return args
}

func workspaceDestination(workspace, checkoutPath string) (string, error) {
	cleanPath := path.Clean(strings.TrimSpace(checkoutPath))
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

func shortBranch(ref string) string { // full tag refs stay exact
	return strings.TrimPrefix(ref, "refs/heads/")
}
