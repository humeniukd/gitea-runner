// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checkout

import (
	"cmp"
	"context"
	"fmt"
	"maps"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"gitea.com/gitea/runner/act/common/git"
	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/internal/pkg/action"

	"github.com/kballard/go-shellquote"
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
	ref := c.Inputs["ref"]
	isWorkflowRepo := repository == ghc.Repository
	local := isWorkflowRepo && slices.Contains([]string{"", shortBranch(ghc.Ref), ghc.Sha}, shortBranch(ref)) && !c.NoSkipCheckout
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
	knownSHA := ref
	if isWorkflowRepo {
		knownSHA = ghc.Sha
	}
	url := strings.TrimSuffix(ghc.ServerURL, "/") + "/" + repository
	token := cmp.Or(c.Inputs["token"], ghc.Token)
	objectFormat := git.ObjectFormat(knownSHA)
	if objectFormat == "" {
		if objectFormat, err = git.RemoteObjectFormat(ctx, url, ref, token, c.InsecureSkipTLS); err != nil {
			return fmt.Errorf("checkout: %s@%s: %w", repository, ref, err)
		}
	}
	if err := clone(ctx, c, cloneCommands(url, ref, dst, objectFormat, depth), token); err != nil {
		return fmt.Errorf("checkout: clone %s@%s: %w", repository, ref, err)
	}
	return nil
}

func clone(ctx context.Context, c *action.Context, commands [][]string, token string) error {
	env := map[string]string{"GIT_TERMINAL_PROMPT": "0"}
	maps.Copy(env, c.Env)
	maps.DeleteFunc(env, func(name, _ string) bool { return slices.Contains(git.LocalEnvVars, name) })
	// safe.directory: a bound workdir belongs to the host user
	maps.Copy(env, git.ConfigEnv(append([]string{"protocol.version=2", "init.defaultBranch=main", "advice.detachedHead=false", "maintenance.auto=false", "safe.directory=*"},
		git.RemoteConfig(token, c.InsecureSkipTLS)...)...))
	if _, host := c.Container.(*container.HostEnvironment); !host || runtime.GOOS != "windows" {
		script := []string{"set -e", `command -v git >/dev/null || { echo "builtin:checkout needs git on PATH" >&2; exit 127; }`}
		for _, command := range commands {
			script = append(script, shellquote.Join(command...))
		}
		return c.Container.Exec([]string{"sh", "-c", strings.Join(script, "\n")}, env, "", "")(ctx)
	}
	for _, command := range commands {
		if err := c.Container.Exec(command, env, "", "")(ctx); err != nil {
			return err
		}
	}
	return nil
}

const fullDepth = "2147483647" // git's INFINITE_DEPTH, also deepens a shallow dst

func cloneCommands(url, ref, dst, objectFormat string, depth int) [][]string {
	repo := func(args ...string) []string { return append([]string{"git", "-C", dst}, args...) }
	commands := [][]string{
		{"git", "init", "--quiet", "--object-format=" + objectFormat, dst},
		repo("config", "remote.origin.url", url), // `remote add` fails on a repeated checkout
		repo("config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*"),
	}
	refSpec := ref // a short name or SHA lands in FETCH_HEAD only
	if branch, ok := strings.CutPrefix(ref, "refs/heads/"); ok {
		refSpec = "+" + ref + ":refs/remotes/origin/" + branch
	} else if strings.HasPrefix(ref, "refs/tags/") {
		refSpec = "+" + ref + ":" + ref
	}
	if depth == 0 {
		commands = append(commands, repo("fetch", "--force", "--tags", "--depth", fullDepth, "origin", refSpec, "+refs/heads/*:refs/remotes/origin/*"))
	} else {
		commands = append(commands, repo("fetch", "--no-tags", "--depth", strconv.Itoa(depth), "origin", refSpec))
	}
	return append(commands, repo("checkout", "--force", "FETCH_HEAD"))
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
