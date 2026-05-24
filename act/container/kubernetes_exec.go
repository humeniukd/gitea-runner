// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	"errors"
	"io"
	"regexp"
	"slices"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	utilexec "k8s.io/client-go/util/exec"
)

// execInJobContainer runs command in the job container over the pods/exec
// subresource, streaming stdout/stderr to the provided writers. A non-zero
// command exit is surfaced as ExitCodeError so it propagates exactly like the
// Docker backend's exec.
func (p *KubernetesPod) execInJobContainer(ctx context.Context, command []string, stdout, stderr io.Writer) error {
	return p.streamExec(ctx, command, nil, stdout, stderr)
}

// streamExec is the shared pods/exec transport. A non-nil stdin is piped into
// the command (used to deliver tar archives for Copy); nil stdout/stderr
// streams are omitted from the request. A non-zero exit maps to ExitCodeError.
func (p *KubernetesPod) streamExec(ctx context.Context, command []string, stdin io.Reader, stdout, stderr io.Writer) error {
	req := p.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(p.builder.name).
		Namespace(p.namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: kubernetesJobContainerName,
			Command:   command,
			Stdin:     stdin != nil,
			Stdout:    stdout != nil,
			Stderr:    stderr != nil,
			TTY:       false,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(p.rest, "POST", req.URL())
	if err != nil {
		return err
	}

	return mapExecError(executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}))
}

// mapExecError maps a non-zero process exit onto ExitCodeError
// (container_types.go) so the runner treats it identically to a failed docker
// exec. Other errors pass through unchanged; nil stays nil.
func mapExecError(err error) error {
	if err == nil {
		return nil
	}
	var exitErr utilexec.ExitError
	if errors.As(err, &exitErr) && exitErr.Exited() {
		return ExitCodeError(exitErr.ExitStatus())
	}
	return err
}

// validEnvKey matches a POSIX shell name. Env keys are emitted UNQUOTED into the
// `export <key>=` script, so a key must be validated to a safe shell identifier
// to avoid command injection through the key (the value is always shell-quoted).
var validEnvKey = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// buildExecCommand wraps an act step command so the pods/exec subresource — which
// has no native env/workdir/user fields — reproduces the Docker exec semantics:
// it exports each env var, changes into the working directory, then execs the
// command. When user is set the script runs via `su <user> -c`, otherwise `sh`.
// Env keys are emitted in sorted order for deterministic output; keys that are
// not valid shell identifiers are skipped (the unquoted key would otherwise be a
// shell-injection vector).
func buildExecCommand(command []string, env map[string]string, user, workdir string) []string {
	var script strings.Builder
	for _, key := range sortedKeys(env) {
		if !validEnvKey.MatchString(key) {
			continue
		}
		script.WriteString("export ")
		script.WriteString(key)
		script.WriteString("=")
		script.WriteString(shellquote.Join(env[key]))
		script.WriteByte('\n')
	}
	if workdir != "" {
		script.WriteString("cd ")
		script.WriteString(shellquote.Join(workdir))
		script.WriteString(" || exit 1\n")
	}
	script.WriteString("exec ")
	script.WriteString(shellquote.Join(command...))

	body := script.String()
	if user != "" {
		return []string{"su", user, "-c", body}
	}
	return []string{"sh", "-c", body}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
