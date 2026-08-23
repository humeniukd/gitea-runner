// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package checkout

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/container"
	"gitea.com/gitea/runner/act/runner/builtin"

	"gitea.dev/actionslib/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeContainer struct {
	container.LinuxContainerEnvironmentExtensions
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
func (f *fakeContainer) Inspect(_ context.Context) (*container.Info, error) {
	return &container.Info{}, nil
}
func (f *fakeContainer) Remove() common.Executor                                { return nopExec }
func (f *fakeContainer) Close() common.Executor                                 { return nopExec }
func (f *fakeContainer) DumpLogs(_ context.Context) error                       { return nil }
func (f *fakeContainer) ReplaceLogWriter(_, _ io.Writer) (io.Writer, io.Writer) { return nil, nil }

var nopExec common.Executor = func(_ context.Context) error { return nil }

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

func checkoutCtx(t *testing.T, fc *fakeContainer, cfg *fakeConfig, inputs map[string]string) *builtin.Context {
	t.Helper()
	if fc == nil {
		fc = &fakeContainer{}
	}
	if cfg == nil {
		cfg = &fakeConfig{}
	}
	if inputs == nil {
		inputs = map[string]string{}
	}
	return &builtin.Context{
		Ctx: context.Background(),
		RunContext: &fakeRunContext{
			jc:    fc,
			wd:    "/host/work",
			cwd:   "/workspace",
			cache: t.TempDir(),
			gh:    &model.GithubContext{Repository: "org/repo", Ref: "refs/heads/main"},
			cfg:   cfg,
		},
		Step:   &model.Step{Uses: "builtin:checkout"},
		Inputs: inputs,
		Env:    map[string]string{},
		Logger: log.New(),
	}
}

func TestCheckoutSkipsWhenWorkspaceBound(t *testing.T) {
	fc := &fakeContainer{}
	err := (&Action{}).Main(checkoutCtx(t, fc, &fakeConfig{bind: true}, nil))
	require.NoError(t, err)
	assert.False(t, fc.called)
}

func TestCheckoutCopiesLocalWorkspace(t *testing.T) {
	fc := &fakeContainer{}
	err := (&Action{}).Main(checkoutCtx(t, fc, &fakeConfig{gitIgnore: true}, map[string]string{"path": "src"}))
	require.NoError(t, err)
	assert.True(t, fc.called)
	assert.Equal(t, "/workspace/src", fc.copyDirDst)
	assert.True(t, strings.HasPrefix(fc.copyDirSrc, "/host/work"))
	assert.True(t, fc.copyDirGI)
}

func TestCheckoutRejectsPathOutsideWorkspace(t *testing.T) {
	fc := &fakeContainer{}
	err := (&Action{}).Main(checkoutCtx(t, fc, &fakeConfig{bind: true}, map[string]string{"path": "../outside"}))
	require.Error(t, err)
	assert.False(t, fc.called)
}

func TestWorkspaceDestination(t *testing.T) {
	for _, tc := range []struct {
		path string
		want string
		err  bool
	}{
		{path: "", want: "/workspace"},
		{path: ".", want: "/workspace"},
		{path: "src", want: "/workspace/src"},
		{path: "nested/../src", want: "/workspace/src"},
		{path: "../outside", err: true},
		{path: "/outside", err: true},
		{path: "nested/../../outside", err: true},
	} {
		t.Run(tc.path, func(t *testing.T) {
			got, err := workspaceDestination("/workspace", tc.path)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFetchDepth(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  int
		err   bool
	}{
		{input: "", want: 1},
		{input: " ", want: 1},
		{input: "0", want: 0},
		{input: "5", want: 5},
		{input: " 5", want: 5},
		{input: "-1", err: true},
		{input: "invalid", err: true},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, err := fetchDepth(tc.input)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCheckoutTakesClonePath(t *testing.T) {
	for _, tc := range []struct {
		name   string
		cfg    *fakeConfig
		inputs map[string]string
	}{
		{name: "no-skip", cfg: &fakeConfig{noSkip: true, offline: true}},
		{name: "other-repo", cfg: &fakeConfig{offline: true}, inputs: map[string]string{"repository": "other/repo"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Error(t, (&Action{}).Main(checkoutCtx(t, nil, tc.cfg, tc.inputs)))
		})
	}
}

func TestRepoURL(t *testing.T) {
	for _, tc := range []struct {
		name string
		gh   *model.GithubContext
		cfg  *fakeConfig
		want string
	}{
		{"server-url", &model.GithubContext{ServerURL: "https://gitea.example.com"}, &fakeConfig{githubInstance: "gitea.example.com", defaultActionInst: "https://github.com"}, "https://gitea.example.com/org/repo"},
		{"default-instance", &model.GithubContext{}, &fakeConfig{defaultActionInst: "github.com"}, "https://github.com/org/repo"},
		{"github-instance", &model.GithubContext{}, &fakeConfig{githubInstance: "gitea.local"}, "https://gitea.local/org/repo"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, repoURL(tc.gh, tc.cfg, "org/repo"))
		})
	}
}

func TestIsLocalCheckoutLike(t *testing.T) {
	gh := &model.GithubContext{Repository: "org/repo", Ref: "refs/heads/main"}
	assert.True(t, isLocalCheckoutLike(map[string]string{}, gh))
	assert.True(t, isLocalCheckoutLike(map[string]string{"repository": "org/repo"}, gh))
	assert.True(t, isLocalCheckoutLike(map[string]string{"repository": "org/repo", "ref": "refs/heads/main"}, gh))
	assert.False(t, isLocalCheckoutLike(map[string]string{"repository": "other/repo"}, gh))
	assert.False(t, isLocalCheckoutLike(map[string]string{"ref": "refs/heads/feature"}, gh))
}
