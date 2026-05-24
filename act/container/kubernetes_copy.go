// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"gitea.com/gitea/runner/act/common"
	"gitea.com/gitea/runner/act/filecollector"

	"github.com/go-git/go-billy/v5/helper/polyfill"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

// Copy writes in-memory files into the job container by piping a tar archive
// into `tar -xf - -C <destPath>` over the exec subresource.
func (e *KubernetesExecutionEnvironment) Copy(destPath string, files ...*FileEntry) common.Executor {
	return func(ctx context.Context) error {
		archive, err := buildFileEntryTar(files)
		if err != nil {
			return err
		}
		return e.pod.copyTarStream(ctx, destPath, archive)
	}
}

// CopyTarStream extracts an existing tar stream into destPath in the job
// container.
func (e *KubernetesExecutionEnvironment) CopyTarStream(ctx context.Context, destPath string, tarStream io.Reader) error {
	return e.pod.copyTarStream(ctx, destPath, tarStream)
}

// CopyDir tars a host directory (honoring .gitignore when requested) and
// extracts it into the job container. Entries carry their absolute destination
// path so extraction targets the container root, mirroring the docker backend.
func (e *KubernetesExecutionEnvironment) CopyDir(destPath, srcPath string, useGitIgnore bool) common.Executor {
	return func(ctx context.Context) error {
		archive, err := buildDirTar(ctx, destPath, srcPath, useGitIgnore)
		if err != nil {
			return err
		}
		return e.pod.copyTarStream(ctx, "/", archive)
	}
}

// GetContainerArchive streams a path out of the job container as a tar archive
// via `tar -cf - -C <dir> <base>`, mirroring docker's CopyFromContainer.
func (e *KubernetesExecutionEnvironment) GetContainerArchive(ctx context.Context, srcPath string) (io.ReadCloser, error) {
	clean := strings.TrimRight(srcPath, "/")
	dir, base := path.Dir(clean), path.Base(clean)

	var buf bytes.Buffer
	if err := e.pod.streamExec(ctx, []string{"tar", "-cf", "-", "-C", dir, base}, nil, &buf, io.Discard); err != nil {
		return nil, err
	}
	return io.NopCloser(&buf), nil
}

// UpdateFromEnv parses an env file inside the job container (read via
// GetContainerArchive) into the env map, reusing the shared parser.
func (e *KubernetesExecutionEnvironment) UpdateFromEnv(srcPath string, env *map[string]string) common.Executor {
	return parseEnvFile(e, srcPath, env)
}

// UpdateFromImageEnv resolves the image's baked environment using C1: it execs
// `env` in the (already-Ready) job container and merges the result. PATH is
// appended; other keys fill only when unset. Resolution is best-effort —
// ApplyExtraPath falls back to DefaultPathVariable — so a failure logs and
// returns nil rather than aborting the job. (C2, a go-containerregistry image
// config fetch, is the documented fallback for the residual case where exec env
// is insufficient; it is intentionally not wired in v1 to avoid the dependency.)
func (e *KubernetesExecutionEnvironment) UpdateFromImageEnv(env *map[string]string) common.Executor {
	envMap := *env
	return func(ctx context.Context) error {
		var buf bytes.Buffer
		if err := e.pod.execInJobContainer(ctx, []string{"env"}, &buf, io.Discard); err != nil {
			common.Logger(ctx).Warnf("kubernetes: reading image env via `exec env` failed, falling back to defaults: %v", err)
			return nil
		}

		scanner := bufio.NewScanner(&buf)
		scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
		for scanner.Scan() {
			key, value, ok := strings.Cut(scanner.Text(), "=")
			if !ok {
				continue
			}
			if key == "PATH" {
				if envMap[key] == "" {
					envMap[key] = value
				} else {
					envMap[key] += ":" + value
				}
			} else if envMap[key] == "" {
				envMap[key] = value
			}
		}
		return nil
	}
}

// copyTarStream creates destPath then extracts the tar stream into it inside the
// job container.
func (p *KubernetesPod) copyTarStream(ctx context.Context, destPath string, tarStream io.Reader) error {
	script := fmt.Sprintf("mkdir -p %s && tar -xf - -C %s", shellQuoteSingle(destPath), shellQuoteSingle(destPath))
	return p.streamExec(ctx, []string{"sh", "-c", script}, tarStream, io.Discard, io.Discard)
}

// buildFileEntryTar serialises in-memory FileEntries into a tar archive.
func buildFileEntryTar(files []*FileEntry) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: file.Mode,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

// buildDirTar tars srcPath with each entry rewritten under destPath, so the
// archive can be extracted at the container root.
func buildDirTar(ctx context.Context, destPath, srcPath string, useGitIgnore bool) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	srcPrefix := filepath.Dir(srcPath)
	if !strings.HasSuffix(srcPrefix, string(filepath.Separator)) {
		srcPrefix += string(filepath.Separator)
	}

	var ignorer gitignore.Matcher
	if useGitIgnore {
		ps, err := gitignore.ReadPatterns(polyfill.New(osfs.New(srcPath)), nil)
		if err != nil {
			common.Logger(ctx).Debugf("Error loading .gitignore: %v", err)
		}
		ignorer = gitignore.NewMatcher(ps)
	}

	fc := &filecollector.FileCollector{
		Fs:        &filecollector.DefaultFs{},
		Ignorer:   ignorer,
		SrcPath:   srcPath,
		SrcPrefix: srcPrefix,
		Handler: &filecollector.TarCollector{
			TarWriter: tw,
			DstDir:    strings.TrimPrefix(destPath, "/"),
		},
	}

	if err := filepath.Walk(srcPath, fc.CollectFiles(ctx, []string{})); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

// shellQuoteSingle wraps a value in single quotes, escaping embedded quotes, so
// it can be embedded literally in an `sh -c` script.
func shellQuoteSingle(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}
