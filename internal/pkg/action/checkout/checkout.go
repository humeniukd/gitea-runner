// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checkout

import (
	"cmp"
	"context"
	"crypto/sha256"
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

var inputNames = []string{"repository", "ref", "token", "path", "fetch-depth", "fetch-tags", "clean", "persist-credentials", "show-progress"}

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
	isWorkflowRepo := strings.EqualFold(repository, ghc.Repository)
	local := isWorkflowRepo && slices.Contains([]string{"", shortBranch(ghc.Ref), ghc.Sha}, shortBranch(ref)) && !c.NoSkipCheckout
	if local && (!c.BindWorkdir || dst == c.Workspace) {
		if !c.BindWorkdir {
			if err := c.Container.CopyDir(dst, c.HostWorkdir+string(filepath.Separator)+".", c.UseGitIgnore, false)(ctx); err != nil {
				return err
			}
		}
		if output := c.Env["GITHUB_OUTPUT"]; output != "" {
			return c.Container.Copy(path.Dir(output), &container.FileEntry{Name: path.Base(output), Mode: 0o666, Body: fmt.Sprintf("ref=%s\ncommit=%s\n", ghc.Ref, ghc.Sha)})(ctx)
		}
		return nil
	}

	co := &checkout{
		url:       strings.TrimSuffix(ghc.ServerURL, "/") + "/" + repository,
		dst:       dst,
		ref:       ref,
		depth:     depth,
		fetchTags: boolInput(c.Inputs, "fetch-tags", false),
		clean:     boolInput(c.Inputs, "clean", true),
		output:    c.Env["GITHUB_OUTPUT"],
	}
	co.objectFormat = git.ObjectFormat(ref)
	switch {
	case co.objectFormat != "":
		co.sha, co.ref = ref, ""
	case ref == "" && isWorkflowRepo && strings.HasPrefix(ghc.Ref, "refs/"):
		co.sha, co.ref = ghc.Sha, ghc.Ref
	case ref == "" && isWorkflowRepo: // a merged pull request's ref is its merge commit
		co.sha = ghc.Sha
	case ref == "HEAD":
		co.ref = ""
	}
	if isWorkflowRepo {
		co.objectFormat = cmp.Or(co.objectFormat, git.ObjectFormat(ghc.Sha))
	}
	token := cmp.Or(c.Inputs["token"], ghc.Token)
	config, remoteConfig := serverConfig(ghc.ServerURL, token, c.InsecureSkipTLS)
	if co.objectFormat == "" || co.sha == "" && !strings.HasPrefix(co.ref, "refs/") {
		advertised, err := lsRemote(ctx, c, gitEnv(c, remoteConfig), co.url, co.ref)
		if err != nil {
			return fmt.Errorf("checkout: %s@%s: %w", repository, ref, err)
		}
		if co.resolve(advertised); co.ref == "" && co.sha == "" {
			return fmt.Errorf("checkout: %s has no ref %q", repository, ref)
		}
	}

	co.outputRef = co.ref
	if ref != co.sha {
		co.outputRef = cmp.Or(ref, co.ref) // as given, like actions/checkout
	}

	// named like actions/checkout's, which create-pull-request knows to hide, and gone with the job
	runnerTemp, _ := c.Container.GetRunnerContext(ctx)["temp"].(string)
	digest := sha256.Sum256([]byte(dst))
	credentials := fmt.Sprintf("git-credentials-%x.config", digest[:8])
	co.credentials = path.Join(runnerTemp, credentials)
	if boolInput(c.Inputs, "persist-credentials", true) {
		remoteConfig = nil
	} else {
		config = ""
	}
	if err := c.Container.Copy(runnerTemp, &container.FileEntry{Name: credentials, Mode: 0o600, Body: config})(ctx); err != nil {
		return err
	}
	if err := run(ctx, c, gitEnv(c, remoteConfig), co.commands()); err != nil {
		return fmt.Errorf("checkout: clone %s@%s: %w", repository, cmp.Or(co.ref, co.sha), err)
	}
	return nil
}

func gitEnv(c *action.Context, remoteConfig []string) map[string]string {
	env := map[string]string{"GIT_TERMINAL_PROMPT": "0", "GIT_LFS_SKIP_SMUDGE": "1"}
	maps.Copy(env, c.Env)
	maps.DeleteFunc(env, func(name, _ string) bool { return slices.Contains(git.LocalEnvVars, name) })
	// safe.directory: a bound workdir belongs to the host user
	maps.Copy(env, git.ConfigEnv(append([]string{"protocol.version=2", "init.defaultBranch=main", "advice.detachedHead=false", "maintenance.auto=false", "safe.directory=*"},
		remoteConfig...)...))
	return env
}

func run(ctx context.Context, c *action.Context, env map[string]string, commands [][]string) error {
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

// lsRemote asks the remote where the job runs, through the same network and git config as the fetch.
func lsRemote(ctx context.Context, c *action.Context, env map[string]string, url, ref string) (string, error) {
	patterns := []string{ref}
	switch {
	case ref == "":
		patterns = []string{"HEAD"}
	case !strings.HasPrefix(ref, "refs/"):
		patterns = []string{"refs/heads/" + ref, "refs/tags/" + ref}
	}
	var out strings.Builder
	stdout, stderr := c.Container.ReplaceLogWriter(&out, &out)
	defer c.Container.ReplaceLogWriter(stdout, stderr)
	if err := c.Container.Exec(append([]string{"git", "ls-remote", "--symref", url}, patterns...), env, "", "")(ctx); err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(out.String()))
	}
	return out.String(), nil
}

// serverConfig scopes the token to the server, so no other host like a repository's LFS endpoint gets it.
func serverConfig(serverURL, token string, insecureSkipTLS bool) (file string, env []string) {
	scope := strings.TrimSuffix(serverURL, "/") + "/"
	file = fmt.Sprintf("[http %q]\n", scope)
	for _, entry := range git.RemoteConfig(token, insecureSkipTLS) {
		setting := strings.TrimPrefix(entry, "http.")
		key, value, _ := strings.Cut(setting, "=")
		file += "\t" + key + " = " + value + "\n"
		env = append(env, "http."+scope+"."+setting)
	}
	return file, env
}

type checkout struct {
	url, dst, objectFormat, credentials, output, outputRef string
	ref, sha                                               string // ref is fully qualified or empty, sha pins it when known
	depth                                                  int
	fetchTags, clean                                       bool
}

const fullDepth = "2147483647" // git's INFINITE_DEPTH, also deepens a shallow dst

func (co *checkout) commands() [][]string {
	repo := func(args ...string) []string { return append([]string{"git", "-C", co.dst}, args...) }
	depth := strconv.Itoa(co.depth)
	if co.depth == 0 {
		depth = fullDepth
	}
	commands := [][]string{
		{"git", "init", "--quiet", "--object-format=" + co.objectFormat, co.dst},
		repo("config", "remote.origin.url", co.url), // `remote add` fails on a repeated checkout
		repo("config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*"),
		repo("config", "--replace-all", "--fixed-value", "include.path", co.credentials, co.credentials),
	}
	var refSpecs []string
	if co.depth == 0 {
		refSpecs = append(refSpecs, "+refs/heads/*:refs/remotes/origin/*")
	}
	if co.fetchTags || co.depth == 0 {
		refSpecs = append(refSpecs, "+refs/tags/*:refs/tags/*")
	}
	pinned := co.sha != "" && co.ref != "" && !strings.HasPrefix(co.ref, "refs/tags/") // a tag is fetched itself to keep its annotation
	if co.ref != "" && !pinned {
		refSpecs = append(refSpecs, "+"+co.ref+":"+trackingRef(co.ref))
	}
	if co.sha != "" {
		refSpecs = append(refSpecs, co.sha)
	}
	commands = append(commands, repo(append([]string{"fetch", "--force", "--no-tags", "--depth", depth, "origin"}, refSpecs...)...))
	if pinned {
		commands = append(commands, repo("update-ref", trackingRef(co.ref), co.sha))
	}
	if branch, ok := strings.CutPrefix(co.ref, "refs/heads/"); ok {
		commands = append(commands, repo("checkout", "--force", "--ignore-other-worktrees", "-B", branch, "--track", trackingRef(co.ref)))
	} else {
		commands = append(commands, repo("checkout", "--force", cmp.Or(co.sha, trackingRef(co.ref))))
	}
	if co.clean {
		commands = append(commands, repo("clean", "-ffdx"))
	}
	if co.output != "" {
		commands = append(commands, repo("log", "-1", "--format=commit=%H%nref="+strings.ReplaceAll(co.outputRef, "%", "%%"), "--output="+co.output))
	}
	return commands
}

func (co *checkout) resolve(advertised string) {
	ref := co.ref
	if !strings.HasPrefix(ref, "refs/") {
		co.ref = ""
	}
	for line := range strings.Lines(advertised) {
		value, name, _ := strings.Cut(strings.TrimSpace(line), "\t")
		if target, ok := strings.CutPrefix(value, "ref: "); ok {
			if name == "HEAD" {
				co.ref = target
			}
			continue
		}
		co.objectFormat = cmp.Or(co.objectFormat, git.ObjectFormat(value))
		switch name {
		case "HEAD":
			if co.ref == "" {
				co.sha = value
			}
		case "refs/heads/" + ref:
			co.ref = name
		case "refs/tags/" + ref:
			co.ref = cmp.Or(co.ref, name)
		}
	}
	co.objectFormat = cmp.Or(co.objectFormat, "sha1")
}

func trackingRef(ref string) string {
	if branch, ok := strings.CutPrefix(ref, "refs/heads/"); ok {
		return "refs/remotes/origin/" + branch
	}
	if pull, ok := strings.CutPrefix(ref, "refs/pull/"); ok {
		return "refs/remotes/pull/" + pull
	}
	return ref
}

func boolInput(inputs map[string]string, name string, fallback bool) bool {
	value := strings.TrimSpace(inputs[name])
	return strings.EqualFold(value, "true") || value == "" && fallback
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
