// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2022 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package git

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/internal/pkg/lock"

	log "github.com/sirupsen/logrus"
)

var (
	codeCommitHTTPRegex = regexp.MustCompile(`^https?://git-codecommit\.(.+)\.amazonaws.com/v1/repos/(.+)$`)
	codeCommitSSHRegex  = regexp.MustCompile(`ssh://git-codecommit\.(.+)\.amazonaws.com/v1/repos/(.+)$`)
	githubHTTPRegex     = regexp.MustCompile(`^https?://.*github.com.*/(.+)/(.+?)(?:.git)?$`)
	githubSSHRegex      = regexp.MustCompile(`github.com[:/](.+)/(.+?)(?:.git)?$`)
	hexRefRegex         = regexp.MustCompile(`^[0-9a-fA-F]+$`)

	cloneLocks lock.Keyed[string]

	ErrShortRef = errors.New("short SHA references are not supported")
)

// AcquireCloneLock returns an unlock function after locking the per-directory mutex for dir.
func AcquireCloneLock(dir string) func() { return cloneLocks.Lock(dir) }

type Error struct {
	err    error
	commit string
}

func (e *Error) Error() string  { return e.err.Error() }
func (e *Error) Unwrap() error  { return e.err }
func (e *Error) Commit() string { return e.commit }

// FindGitRevision gets the current git revision.
func FindGitRevision(ctx context.Context, file string) (shortSHA, sha string, err error) {
	logger := common.Logger(ctx)
	sha, err = gitOutput(ctx, file, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		logger.WithError(err).Error("path", file, "not located inside a git repository")
		return "", "", err
	}
	if len(sha) < 7 {
		return "", "", errors.New("HEAD could not be resolved")
	}
	logger.Debugf("Found revision: %s", sha)
	return sha[:7], sha, nil
}

// FindGitRef gets the current git ref.
func FindGitRef(ctx context.Context, file string) (string, error) {
	logger := common.Logger(ctx)
	logger.Debugf("Loading revision from git directory")
	_, revision, err := FindGitRevision(ctx, file)
	if err != nil {
		return "", err
	}
	logger.Debugf("HEAD points to '%s'", revision)

	refs, err := gitOutput(ctx, file, "for-each-ref", "--points-at=HEAD", "--format=%(refname)", "refs/tags", "refs/heads")
	if err != nil {
		return "", err
	}
	var tag, branch string
	for ref := range strings.FieldsSeq(refs) {
		switch {
		case strings.HasPrefix(ref, "refs/tags/"):
			tag = ref
		case strings.HasPrefix(ref, "refs/heads/"):
			branch = ref
		}
	}
	if tag != "" {
		return tag, nil
	}
	if branch != "" {
		return branch, nil
	}
	return "", fmt.Errorf("failed to identify reference (tag/branch) for the checked-out revision '%s'", revision)
}

// FindGithubRepo gets the repository slug from origin.
func FindGithubRepo(ctx context.Context, file, githubInstance string) (string, error) {
	url, err := findGitRemoteURL(ctx, file)
	if err != nil {
		return "", err
	}
	_, slug := findGitSlug(url, githubInstance)
	return slug, nil
}

func findGitRemoteURL(ctx context.Context, file string) (string, error) {
	url, err := gitOutput(ctx, file, "remote", "get-url", "origin")
	if err != nil {
		return "", err
	}
	if url == "" {
		return "", errors.New("remote 'origin' exists but has no URL")
	}
	return url, nil
}

func findGitSlug(url, githubInstance string) (string, string) {
	if matches := codeCommitHTTPRegex.FindStringSubmatch(url); matches != nil {
		return "CodeCommit", matches[2]
	} else if matches := codeCommitSSHRegex.FindStringSubmatch(url); matches != nil {
		return "CodeCommit", matches[2]
	} else if matches := githubHTTPRegex.FindStringSubmatch(url); matches != nil {
		return "GitHub", fmt.Sprintf("%s/%s", matches[1], matches[2])
	} else if matches := githubSSHRegex.FindStringSubmatch(url); matches != nil {
		return "GitHub", fmt.Sprintf("%s/%s", matches[1], matches[2])
	} else if githubInstance != "github.com" {
		gheHTTPRegex := regexp.MustCompile(fmt.Sprintf(`^https?://%s/(.+)/(.+?)(?:.git)?$`, githubInstance))
		gheSSHRegex := regexp.MustCompile(githubInstance + "[:/](.+)/(.+?)(?:.git)?$")
		if matches := gheHTTPRegex.FindStringSubmatch(url); matches != nil {
			return "GitHubEnterprise", fmt.Sprintf("%s/%s", matches[1], matches[2])
		} else if matches := gheSSHRegex.FindStringSubmatch(url); matches != nil {
			return "GitHubEnterprise", fmt.Sprintf("%s/%s", matches[1], matches[2])
		}
	}
	return "", url
}

// NewGitCloneExecutorInput is the input for NewGitCloneExecutor.
type NewGitCloneExecutorInput struct {
	URL             string
	Ref             string
	Dir             string
	Token           string
	OfflineMode     bool
	Depth           int
	Quiet           bool
	InsecureSkipTLS bool
}

// CloneIfRequired ensures Dir is a clone of URL and reports whether it reused a cached clone.
func CloneIfRequired(ctx context.Context, input NewGitCloneExecutorInput, logger log.FieldLogger) (bool, error) {
	origin, err := findGitRemoteURL(ctx, input.Dir)
	if err == nil && origin == input.URL {
		return true, nil
	}
	if err == nil {
		logger.Debugf("Removing cached clone at %s because origin URL changed from %s to %s", input.Dir, origin, input.URL)
	} else if _, statErr := os.Stat(input.Dir); statErr == nil {
		logger.Debugf("Removing cached clone at %s because origin cannot be read: %v", input.Dir, err)
	}
	if err := os.RemoveAll(input.Dir); err != nil {
		return false, fmt.Errorf("remove cached clone %s: %w", input.Dir, err)
	}
	if err := clone(ctx, input); err != nil {
		return false, err
	}
	if err := os.Chmod(input.Dir, 0o755); err != nil {
		return false, err
	}
	return false, nil
}

// NewGitCloneExecutor creates an executor to clone git repositories with the native Git CLI.
func NewGitCloneExecutor(input NewGitCloneExecutorInput) common.Executor {
	return func(ctx context.Context) error {
		logger := common.Logger(ctx)
		if input.Quiet {
			logger.Debugf("git clone '%s' # ref=%s", input.URL, input.Ref)
		} else {
			logger.Infof("git clone '%s' # ref=%s", input.URL, input.Ref)
		}
		logger.Debugf("  cloning %s to %s", input.URL, input.Dir)
		defer AcquireCloneLock(input.Dir)()

		reused, err := CloneIfRequired(ctx, input, logger)
		if err != nil {
			return err
		}
		resolved, err := resolveRef(ctx, input.Dir, input.Ref)
		if err != nil && input.OfflineMode {
			return err
		}
		if !input.OfflineMode && (err != nil || (reused && !isFullHash(ctx, input.Dir, input.Ref))) {
			if err := fetch(ctx, input); err != nil {
				return err
			}
			resolved, err = resolveRef(ctx, input.Dir, input.Ref)
		}
		if err != nil {
			logger.Errorf("Unable to resolve %s: %v", input.Ref, err)
			return err
		}
		if isShortHash(ctx, input.Dir, input.Ref, resolved) {
			return &Error{err: ErrShortRef, commit: resolved}
		}
		if _, err := runGit(ctx, input.Dir, nil, "checkout", "--force", "--detach", resolved); err != nil {
			return fmt.Errorf("checkout %s: %w", resolved, err)
		}
		reusedMsg := ""
		if input.OfflineMode && reused {
			reusedMsg = " (reused in offline mode)"
		}
		logger.Debugf("Cloned %s to %s%s", input.URL, input.Dir, reusedMsg)
		logger.Debugf("Checked out %s", input.Ref)
		return nil
	}
}

func clone(ctx context.Context, input NewGitCloneExecutorInput) error {
	if input.Depth > 0 && isHashRef(input.Ref) {
		if err := shallowPinnedClone(ctx, input); err == nil {
			return nil
		} else if err := removePartialClone(input.Dir); err != nil {
			return err
		}
	}
	args := []string{"clone", "--no-checkout"}
	if input.Depth > 0 {
		args = append(args, "--no-local", "--depth", strconv.Itoa(input.Depth), "--no-tags")
		if input.Ref != "HEAD" && !strings.HasPrefix(input.Ref, "refs/") {
			args = append(args, "--branch", input.Ref)
		}
	}
	args = append(args, input.URL, input.Dir)
	if _, err := runGit(ctx, "", &input, args...); err == nil {
		return nil
	} else if input.Depth == 0 {
		return err
	}
	if err := removePartialClone(input.Dir); err != nil {
		return err
	}
	_, err := runGit(ctx, "", &input, "clone", "--no-checkout", input.URL, input.Dir)
	return err
}

func shallowPinnedClone(ctx context.Context, input NewGitCloneExecutorInput) error {
	if _, err := runGit(ctx, "", &input, "init", input.Dir); err != nil {
		return err
	}
	if _, err := runGit(ctx, input.Dir, &input, "remote", "add", "origin", input.URL); err != nil {
		return err
	}
	_, err := runGit(ctx, input.Dir, &input, "fetch", "--depth", strconv.Itoa(input.Depth), "--no-tags", "origin", "+"+input.Ref+":refs/pinned/"+input.Ref)
	return err
}

func fetch(ctx context.Context, input NewGitCloneExecutorInput) error {
	args := []string{"fetch", "--force", "--prune"}
	if shallow, _ := gitOutput(ctx, input.Dir, "rev-parse", "--is-shallow-repository"); shallow == "true" {
		args = append(args, "--depth", "1", "--no-tags", "origin", fetchRefSpec(ctx, input.Dir, input.Ref))
	} else {
		args = append(args, "origin", fetchRefSpec(ctx, input.Dir, input.Ref))
	}
	_, err := runGit(ctx, input.Dir, &input, args...)
	return err
}

func fetchRefSpec(ctx context.Context, dir, ref string) string {
	if isHashRef(ref) {
		return "+" + ref + ":refs/pinned/" + ref
	}
	if strings.HasPrefix(ref, "refs/") {
		return "+" + ref + ":" + ref
	}
	if _, err := runGit(ctx, dir, nil, "show-ref", "--verify", "--quiet", "refs/tags/"+ref); err == nil {
		return "+refs/tags/" + ref + ":refs/tags/" + ref
	}
	return "+refs/heads/" + ref + ":refs/remotes/origin/" + ref
}

func resolveRef(ctx context.Context, dir, ref string) (string, error) {
	if isHashRef(ref) {
		return gitOutput(ctx, dir, "rev-parse", "--verify", ref+"^{commit}")
	}
	if sha, err := gitOutput(ctx, dir, "rev-parse", "--verify", "refs/tags/"+ref+"^{commit}"); err == nil {
		return sha, nil
	}
	if sha, err := gitOutput(ctx, dir, "rev-parse", "--verify", "refs/remotes/origin/"+ref+"^{commit}"); err == nil {
		return sha, nil
	}
	return gitOutput(ctx, dir, "rev-parse", "--verify", ref+"^{commit}")
}

func isHashRef(ref string) bool { return hexRefRegex.MatchString(ref) }

func isFullHash(ctx context.Context, dir, ref string) bool {
	if !isHashRef(ref) {
		return false
	}
	format, err := gitOutput(ctx, dir, "rev-parse", "--show-object-format")
	if err != nil {
		return false
	}
	return (format == "sha1" && len(ref) == 40) || (format == "sha256" && len(ref) == 64)
}

func isShortHash(ctx context.Context, dir, ref, resolved string) bool {
	return isHashRef(ref) && !isFullHash(ctx, dir, ref) && strings.HasPrefix(resolved, ref)
}

func removePartialClone(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("remove partial clone %s: %w", dir, err)
	}
	return nil
}

func gitOutput(ctx context.Context, dir string, args ...string) (string, error) {
	output, err := runGit(ctx, dir, nil, args...)
	return strings.TrimSpace(string(output)), err
}

func runGit(ctx context.Context, dir string, input *NewGitCloneExecutorInput, args ...string) ([]byte, error) {
	commandArgs := make([]string, 0, len(args)+4)
	if input != nil {
		if input.InsecureSkipTLS {
			commandArgs = append(commandArgs, "-c", "http.sslVerify=false")
		}
		if input.Token != "" {
			auth := base64.StdEncoding.EncodeToString([]byte("token:" + input.Token))
			commandArgs = append(commandArgs, "-c", "http.extraHeader=Authorization: Basic "+auth)
		}
	}
	if dir != "" {
		commandArgs = append(commandArgs, "-C", dir)
	}
	commandArgs = append(commandArgs, args...)
	cmd := exec.CommandContext(ctx, "git", commandArgs...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	cmd.WaitDelay = time.Second
	var stderr strings.Builder
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return output, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return output, nil
}
