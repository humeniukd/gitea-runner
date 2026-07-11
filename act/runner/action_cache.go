// Copyright 2023 The Gitea Authors. All rights reserved.
// Copyright 2023 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"archive/tar"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"path"
	"strings"

	git "github.com/go-git/go-git/v6"
	config "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/client"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

type ActionCache interface {
	Fetch(ctx context.Context, cacheDir, url, ref, token string) (string, error)
	GetTarArchive(ctx context.Context, cacheDir, sha, includePrefix string) (io.ReadCloser, error)
}

type GoGitActionCache struct {
	Path string
}

func (c GoGitActionCache) Fetch(ctx context.Context, cacheDir, url, ref, token string) (string, error) {
	gitPath := path.Join(c.Path, safeFilename(cacheDir)+".git")
	gogitrepo, err := git.PlainInit(gitPath, true)
	if errors.Is(err, git.ErrTargetDirNotEmpty) {
		gogitrepo, err = git.PlainOpen(gitPath)
	}
	if err != nil {
		return "", err
	}
	tmpBranch := make([]byte, 12)
	if _, err := rand.Read(tmpBranch); err != nil {
		return "", err
	}
	branchName := hex.EncodeToString(tmpBranch)

	var clientOptions []client.Option
	if token != "" {
		clientOptions = append(clientOptions, client.WithHTTPAuth(&http.BasicAuth{
			Username: "token",
			Password: token,
		}))
	}
	remote, err := gogitrepo.CreateRemoteAnonymous(&config.RemoteConfig{
		Name: "anonymous",
		URLs: []string{
			url,
		},
	})
	if err != nil {
		return "", err
	}
	defer func() {
		_ = gogitrepo.DeleteBranch(branchName)
	}()

	resolveRef := branchName
	err = remote.FetchContext(ctx, &git.FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec(ref + ":" + branchName),
		},
		ClientOptions: clientOptions,
		Force:         true,
	})
	if errors.Is(err, git.ErrExactSHA1NotSupported) {
		// ref is an exact commit SHA and the server does not advertise
		// allow-*-sha1-in-want, so it cannot be fetched directly. go-git's own
		// upload-pack server never advertises it, and not every Gitea instance
		// enables it either. Fetch the advertised branches and tags instead and
		// resolve the commit from the local object store; this succeeds whenever
		// the commit is reachable from a ref, as SHA-pinned actions are.
		err = remote.FetchContext(ctx, &git.FetchOptions{
			RefSpecs: []config.RefSpec{
				"+refs/heads/*:refs/heads/*",
				"+refs/tags/*:refs/tags/*",
			},
			ClientOptions: clientOptions,
			Force:         true,
		})
		resolveRef = ref
	}
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return "", err
	}
	hash, err := gogitrepo.ResolveRevision(plumbing.Revision(resolveRef))
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

func (c GoGitActionCache) GetTarArchive(ctx context.Context, cacheDir, sha, includePrefix string) (io.ReadCloser, error) {
	gitPath := path.Join(c.Path, safeFilename(cacheDir)+".git")
	gogitrepo, err := git.PlainOpen(gitPath)
	if err != nil {
		return nil, err
	}
	commit, err := gogitrepo.CommitObject(plumbing.NewHash(sha))
	if err != nil {
		return nil, err
	}
	files, err := commit.Files()
	if err != nil {
		return nil, err
	}
	rpipe, wpipe := io.Pipe()
	// Interrupt io.Copy using ctx
	ch := make(chan int, 1)
	go func() {
		select {
		case <-ctx.Done():
			wpipe.CloseWithError(ctx.Err())
		case <-ch:
		}
	}()
	go func() {
		defer wpipe.Close()
		defer close(ch)
		tw := tar.NewWriter(wpipe)
		cleanIncludePrefix := path.Clean(includePrefix)
		wpipe.CloseWithError(files.ForEach(func(f *object.File) error {
			if err := ctx.Err(); err != nil {
				return err
			}
			name := f.Name
			if strings.HasPrefix(name, cleanIncludePrefix+"/") {
				name = name[len(cleanIncludePrefix)+1:]
			} else if cleanIncludePrefix != "." && name != cleanIncludePrefix {
				return nil
			}
			fmode, err := f.Mode.ToOSFileMode()
			if err != nil {
				return err
			}
			if fmode&fs.ModeSymlink == fs.ModeSymlink {
				content, err := f.Contents()
				if err != nil {
					return err
				}
				return tw.WriteHeader(&tar.Header{
					Name:     name,
					Mode:     int64(fmode),
					Linkname: content,
				})
			}
			err = tw.WriteHeader(&tar.Header{
				Name: name,
				Mode: int64(fmode),
				Size: f.Size,
			})
			if err != nil {
				return err
			}
			reader, err := f.Reader()
			if err != nil {
				return err
			}
			_, err = io.Copy(tw, reader)
			return err
		}))
	}()
	return rpipe, err
}
