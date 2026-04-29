// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checkout

import (
	"context"
	"errors"
	"io"
	"path"
	"strings"
	"testing"

	"gitea.com/gitea/act_runner/act/common"
	"gitea.com/gitea/act_runner/act/container"
	"gitea.com/gitea/act_runner/act/model"
	"gitea.com/gitea/act_runner/act/runner/builtin"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeContainer records the CopyDir invocation we want to assert and is a
// no-op for everything else. It satisfies container.ExecutionsEnvironment.
type fakeContainer struct {
	copyDirDst string
	copyDirSrc string
	copyDirGI  bool
	called     bool
	copyErr    error
}

func (f *fakeContainer) CopyDir(destPath, srcPath string, useGitIgnore bool) common.Executor {
	return func(_ context.Context) error {
		f.called = true
		f.copyDirDst = destPath
		f.copyDirSrc = srcPath
		f.copyDirGI = useGitIgnore
		return f.copyErr
	}
}

// container.Container surface (the rest are unused in tests).
func (f *fakeContainer) Create(_, _ []string) common.Executor      { return nopExec }
func (f *fakeContainer) ConnectToNetwork(_ string) common.Executor { return nopExec }
func (f *fakeContainer) Copy(_ string, _ ...*container.FileEntry) common.Executor {
	return nopExec
}
func (f *fakeContainer) CopyTarStream(_ context.Context, _ string, _ io.Reader) error { return nil }
func (f *fakeContainer) GetContainerArchive(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeContainer) Pull(_ bool) common.Executor  { return nopExec }
func (f *fakeContainer) Start(_ bool) common.Executor { return nopExec }
func (f *fakeContainer) Exec(_ []string, _ map[string]string, _, _ string) common.Executor {
	return nopExec
}

func (f *fakeContainer) UpdateFromEnv(_ string, _ *map[string]string) common.Executor {
	return nopExec
}
func (f *fakeContainer) UpdateFromImageEnv(_ *map[string]string) common.Executor { return nopExec }
func (f *fakeContainer) Remove() common.Executor                                 { return nopExec }
func (f *fakeContainer) Close() common.Executor                                  { return nopExec }
func (f *fakeContainer) ReplaceLogWriter(_, _ io.Writer) (io.Writer, io.Writer)  { return nil, nil }

// container.ExecutionsEnvironment surface.
func (f *fakeContainer) ToContainerPath(p string) string {
	return path.Join("/workspace", strings.TrimPrefix(p, "/"))
}
func (f *fakeContainer) GetActPath() string          { return "/var/run/act" }
func (f *fakeContainer) GetPathVariableName() string { return "PATH" }
func (f *fakeContainer) DefaultPathVariable() string { return "/usr/bin" }
func (f *fakeContainer) JoinPathVariable(parts ...string) string {
	return strings.Join(parts, ":")
}

func (f *fakeContainer) GetRunnerContext(_ context.Context) map[string]any {
	return map[string]any{}
}
func (f *fakeContainer) IsEnvironmentCaseInsensitive() bool { return false }

var nopExec common.Executor = func(_ context.Context) error { return nil }

// fakeRunContext satisfies builtin.RunContext for tests.
type fakeRunContext struct {
	jc    container.ExecutionsEnvironment
	wd    string
	cwd   string
	cache string
	gh    *model.GithubContext
	cfg   builtin.Config
}

func (f *fakeRunContext) JobContainer() container.ExecutionsEnvironment { return f.jc }
func (f *fakeRunContext) Workdir() string                               { return f.wd }
func (f *fakeRunContext) ContainerWorkdir() string                      { return f.cwd }
func (f *fakeRunContext) ActionCacheDir() string                        { return f.cache }
func (f *fakeRunContext) GithubContext() *model.GithubContext           { return f.gh }
func (f *fakeRunContext) Config() builtin.Config                        { return f.cfg }

// fakeConfig satisfies builtin.Config for tests.
type fakeConfig struct {
	token             string
	githubInstance    string
	defaultActionInst string
	insecureTLS       bool
	bind              bool
	noSkip            bool
	gitIgnore         bool
	offline           bool
}

func (f *fakeConfig) Token() string                 { return f.token }
func (f *fakeConfig) GitHubInstance() string        { return f.githubInstance }
func (f *fakeConfig) DefaultActionInstance() string { return f.defaultActionInst }
func (f *fakeConfig) InsecureSkipTLS() bool         { return f.insecureTLS }
func (f *fakeConfig) BindWorkdir() bool             { return f.bind }
func (f *fakeConfig) NoSkipCheckout() bool          { return f.noSkip }
func (f *fakeConfig) UseGitIgnore() bool            { return f.gitIgnore }
func (f *fakeConfig) ActionOfflineMode() bool       { return f.offline }

func TestCheckoutSkipsWhenWorkspaceBound(t *testing.T) {
	fc := &fakeContainer{}
	rc := &fakeRunContext{
		jc:    fc,
		wd:    "/host/work",
		cwd:   "/workspace",
		cache: t.TempDir(),
		gh:    &model.GithubContext{Repository: "org/repo", Ref: "refs/heads/main"},
		cfg:   &fakeConfig{bind: true},
	}

	a := &Action{}
	ctx := &builtin.Context{
		Ctx:        context.Background(),
		RunContext: rc,
		Step:       &model.Step{Uses: "checkout"},
		Inputs:     map[string]string{},
		Env:        map[string]string{},
		Logger:     log.New(),
	}

	err := a.Main(ctx)
	require.NoError(t, err)
	assert.False(t, fc.called, "CopyDir must not be called when BindWorkdir is true and workspace is local")
}

func TestCheckoutCopiesLocalWorkspace(t *testing.T) {
	fc := &fakeContainer{}
	rc := &fakeRunContext{
		jc:    fc,
		wd:    "/host/work",
		cwd:   "/workspace",
		cache: t.TempDir(),
		gh:    &model.GithubContext{Repository: "org/repo", Ref: "refs/heads/main"},
		cfg:   &fakeConfig{gitIgnore: true},
	}

	a := &Action{}
	ctx := &builtin.Context{
		Ctx:        context.Background(),
		RunContext: rc,
		Step:       &model.Step{Uses: "checkout"},
		Inputs:     map[string]string{"path": "src"},
		Env:        map[string]string{},
		Logger:     log.New(),
	}

	err := a.Main(ctx)
	require.NoError(t, err)
	assert.True(t, fc.called, "CopyDir must be called for local checkout when not bound")
	assert.Equal(t, "/workspace/src", fc.copyDirDst)
	assert.True(t, strings.HasPrefix(fc.copyDirSrc, "/host/work"))
	assert.True(t, fc.copyDirGI)
}

func TestCheckoutNoSkipForcesClonePath(t *testing.T) {
	rc := &fakeRunContext{
		jc:    &fakeContainer{},
		wd:    "/host/work",
		cwd:   "/workspace",
		cache: t.TempDir(),
		gh:    &model.GithubContext{Repository: "org/repo", Ref: "refs/heads/main"},
		cfg:   &fakeConfig{noSkip: true, offline: true},
	}

	a := &Action{}
	ctx := &builtin.Context{
		Ctx:        context.Background(),
		RunContext: rc,
		Step:       &model.Step{Uses: "checkout"},
		Inputs:     map[string]string{},
		Env:        map[string]string{},
		Logger:     log.New(),
	}

	err := a.Main(ctx)
	// We expect an error here because the offline clone of "org/repo" cannot
	// resolve, but reaching that path proves the skip-if-bound branch was
	// bypassed when NoSkipCheckout is set.
	assert.Error(t, err)
}

func TestCheckoutRedirectsToOtherRepoTakesClonePath(t *testing.T) {
	rc := &fakeRunContext{
		jc:    &fakeContainer{},
		wd:    "/host/work",
		cwd:   "/workspace",
		cache: t.TempDir(),
		gh:    &model.GithubContext{Repository: "org/repo", Ref: "refs/heads/main"},
		cfg:   &fakeConfig{offline: true},
	}

	a := &Action{}
	ctx := &builtin.Context{
		Ctx:        context.Background(),
		RunContext: rc,
		Step:       &model.Step{Uses: "checkout"},
		Inputs:     map[string]string{"repository": "other/repo"},
		Env:        map[string]string{},
		Logger:     log.New(),
	}

	err := a.Main(ctx)
	// Same as above: offline clone of an arbitrary repo errors out, but the
	// fact that we get past isLocalCheckoutLike is the assertion.
	assert.Error(t, err)
}

func TestRepoURLPrefersServerURL(t *testing.T) {
	gh := &model.GithubContext{ServerURL: "https://gitea.example.com"}
	cfg := &fakeConfig{githubInstance: "gitea.example.com", defaultActionInst: "https://github.com"}
	got := repoURL(gh, cfg, "org/repo")
	assert.Equal(t, "https://gitea.example.com/org/repo", got)
}

func TestRepoURLFallsBackToDefaultActionInstance(t *testing.T) {
	gh := &model.GithubContext{}
	cfg := &fakeConfig{defaultActionInst: "github.com"}
	got := repoURL(gh, cfg, "org/repo")
	assert.Equal(t, "https://github.com/org/repo", got)
}

func TestRepoURLFallsBackToGithubInstance(t *testing.T) {
	gh := &model.GithubContext{}
	cfg := &fakeConfig{githubInstance: "gitea.local"}
	got := repoURL(gh, cfg, "org/repo")
	assert.Equal(t, "https://gitea.local/org/repo", got)
}

func TestIsLocalCheckoutLike(t *testing.T) {
	gh := &model.GithubContext{Repository: "org/repo", Ref: "refs/heads/main"}

	assert.True(t, isLocalCheckoutLike(map[string]string{}, gh))
	assert.True(t, isLocalCheckoutLike(map[string]string{"repository": "org/repo"}, gh))
	assert.True(t, isLocalCheckoutLike(map[string]string{"repository": "org/repo", "ref": "refs/heads/main"}, gh))

	assert.False(t, isLocalCheckoutLike(map[string]string{"repository": "other/repo"}, gh))
	assert.False(t, isLocalCheckoutLike(map[string]string{"ref": "refs/heads/feature"}, gh))
}

func TestActionMetadata(t *testing.T) {
	a := &Action{}
	assert.Equal(t, "checkout", a.Name())
	assert.False(t, a.HasPre())
	assert.False(t, a.HasPost())
	assert.NoError(t, a.Pre(&builtin.Context{}))
	assert.NoError(t, a.Post(&builtin.Context{}))
}
