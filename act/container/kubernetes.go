// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"gitea.com/gitea/runner/act/common"

	"github.com/kballard/go-shellquote"
	"github.com/spf13/pflag"
)

// KubernetesOptions configures the pod a job runs in.
type KubernetesOptions struct {
	Kubeconfig  string
	Namespace   string
	PodTemplate map[string]any
	MaxLifetime time.Duration
	ForcePull   bool
}

// KubernetesPod runs a job container and its service containers as one pod.
type KubernetesPod struct {
	options          KubernetesOptions
	client           *kubeClient
	name, secretName string
	job              *kubernetesContainer
	services         []*kubernetesContainer
	arch, path       string
}

type kubernetesContainer struct {
	LinuxContainerEnvironmentExtensions
	pod   *KubernetesPod
	id    string // the workflow's service id, which name may differ from
	name  string
	input *NewContainerInput
	probe map[string]any
}

type containerStatus struct {
	Name  string
	Ready bool
	State struct {
		Waiting    *struct{ Reason, Message string }
		Running    *struct{}
		Terminated *struct {
			ExitCode int
			Message  string
		}
	}
}

type podStatus struct {
	Spec   struct{ Containers []struct{} }
	Status struct {
		Phase, Reason, Message                   string
		Conditions                               []struct{ Type, Status, Message string }
		InitContainerStatuses, ContainerStatuses []containerStatus
	}
}

func NewKubernetesPod(options KubernetesOptions) (*KubernetesPod, error) {
	client, err := newKubeClient(options.Kubeconfig, options.Namespace)
	if err != nil {
		return nil, fmt.Errorf("connect to Kubernetes: %w", err)
	}
	return &KubernetesPod{options: options, client: client}, nil
}

// JobContainer returns the container of the pod the job's steps run in.
func (p *KubernetesPod) JobContainer(input *NewContainerInput) ExecutionsEnvironment {
	p.job = &kubernetesContainer{pod: p, name: "job", input: input}
	return p.job
}

var invalidLabelChars = regexp.MustCompile(`[^a-z0-9-]+`)

// ServiceContainer returns a container of the pod the job reaches at hostname id, made a DNS label as container names and hostnames must be.
func (p *KubernetesPod) ServiceContainer(id string, input *NewContainerInput) ExecutionsEnvironment {
	name := strings.Trim(invalidLabelChars.ReplaceAllString(strings.ToLower(id), "-"), "-")
	service := &kubernetesContainer{pod: p, id: id, name: strings.TrimRight(name[:min(63, len(name))], "-"), input: input}
	p.services = append(p.services, service)
	return service
}

func (p *KubernetesPod) create(ctx context.Context) error {
	prefix := strings.Trim(strings.ToLower(p.job.input.Name), "-")
	prefix = cmp.Or(strings.TrimRight(prefix[:min(40, len(prefix))], "-"), "gitea-runner") + "-"
	name := prefix + strings.ToLower(rand.Text()[:8]) // chosen here, so cleanup knows it even when a create response is lost
	metadata := map[string]any{"name": name, "labels": map[string]any{"app.kubernetes.io/managed-by": "gitea-runner"}}
	if hostname, err := os.Hostname(); p.client.tokenFile != "" && err == nil {
		var runner struct{ Metadata struct{ Name, UID string } }
		if err := p.client.do(ctx, http.MethodGet, p.client.path("pods/"+hostname), nil, &runner); err == nil {
			metadata["ownerReferences"] = []any{map[string]any{"apiVersion": "v1", "kind": "Pod", "name": runner.Metadata.Name, "uid": runner.Metadata.UID}}
		} else {
			common.Logger(ctx).Debugf("The runner is not pod %s, job pods get no owner: %v", hostname, err)
		}
	}
	secretData := map[string]string{}
	manifest, err := p.manifest(name, secretData, metadata)
	if err != nil {
		return err
	}
	for _, service := range p.services {
		if !strings.EqualFold(service.id, service.name) {
			common.Logger(ctx).Warnf("Service %s is reachable as %s or localhost, as a Kubernetes hostname cannot be %q", service.id, service.name, service.id)
		}
	}
	if len(secretData) > 0 {
		p.secretName = name
		if err := p.client.do(ctx, http.MethodPost, p.client.path("secrets"), map[string]any{"metadata": metadata, "stringData": secretData}, nil); err != nil {
			return fmt.Errorf("create secret: %w", err)
		}
	}
	p.name = name
	if err := p.client.do(ctx, http.MethodPost, p.client.path("pods"), manifest, nil); err != nil {
		return fmt.Errorf("create pod: %w", err)
	}
	return nil
}

func (p *KubernetesPod) manifest(secret string, secretData map[string]string, metadata map[string]any) (any, error) {
	templateSpec, _ := p.options.PodTemplate["spec"].(map[string]any)
	templateContainers, _ := templateSpec["containers"].([]any)
	var containers, hostnames []any
	names := map[string]bool{}
	for _, c := range append([]*kubernetesContainer{p.job}, p.services...) {
		var template any
		if index := slices.IndexFunc(templateContainers, func(container any) bool { return manifestName(container) == c.name }); index >= 0 {
			template = templateContainers[index]
		}
		if c != p.job && (c.name == "" || names[c.name] || template != nil) {
			return nil, fmt.Errorf("service %s maps to container %q, which is empty or named like another container of the pod", c.id, c.name)
		}
		names[c.name] = true
		spec, err := c.spec(secret, secretData)
		if err != nil {
			return nil, err
		}
		containers = append(containers, mergeManifest(template, spec))
		if c != p.job {
			hostnames = append(hostnames, c.name)
		}
	}
	spec := map[string]any{
		"restartPolicy": "Never",
		"containers":    containers,
		"volumes": []any{
			map[string]any{"name": "workspace", "emptyDir": map[string]any{}},
			map[string]any{"name": "act", "emptyDir": map[string]any{}},
		},
	}
	if len(hostnames) > 0 {
		spec["hostAliases"] = []any{map[string]any{"ip": "127.0.0.1", "hostnames": hostnames}}
	}
	if p.options.MaxLifetime > 0 {
		spec["activeDeadlineSeconds"] = max(1, int64(p.options.MaxLifetime.Seconds()))
	}
	defaults := map[string]any{"spec": map[string]any{
		"automountServiceAccountToken":  false,
		"enableServiceLinks":            false,
		"terminationGracePeriodSeconds": 0, // the job container's tail ignores SIGTERM
	}}
	pod := map[string]any{"apiVersion": "v1", "kind": "Pod", "metadata": metadata, "spec": spec}
	return mergeManifest(mergeManifest(defaults, p.options.PodTemplate), pod), nil
}

// mergeManifest merges over onto base: maps key by key, list items of the same name replaced, other values replaced.
func mergeManifest(base, over any) any {
	switch over := over.(type) {
	case map[string]any:
		merged := map[string]any{}
		if base, ok := base.(map[string]any); ok {
			maps.Copy(merged, base)
		}
		for key, value := range over {
			merged[key] = mergeManifest(merged[key], value)
		}
		return merged
	case []any:
		baseList, _ := base.([]any)
		merged := slices.Clone(baseList)
		for _, item := range over {
			name := manifestName(item)
			index := slices.IndexFunc(merged, func(existing any) bool { return name != "" && manifestName(existing) == name })
			if index < 0 {
				merged = append(merged, item)
			} else {
				merged[index] = item
			}
		}
		return merged
	}
	return over
}

func manifestName(item any) string {
	fields, _ := item.(map[string]any)
	name, _ := fields["name"].(string)
	return name
}

func (c *kubernetesContainer) spec(secret string, secretData map[string]string) (map[string]any, error) {
	env := []any{}
	for index, entry := range c.input.Env {
		key, value, _ := strings.Cut(entry, "=")
		if c == c.pod.job {
			env = append(env, map[string]any{"name": key, "value": value})
			continue
		}
		secretKey := fmt.Sprintf("%s.%d", c.name, index)
		secretData[secretKey] = value
		env = append(env, map[string]any{"name": key, "valueFrom": map[string]any{"secretKeyRef": map[string]any{"name": secret, "key": secretKey}}})
	}
	spec := map[string]any{
		"name": c.name, "image": c.input.Image, "command": escapeExpansion(c.input.Entrypoint), "args": escapeExpansion(c.input.Cmd),
		"workingDir": c.input.WorkingDir, "env": env, "imagePullPolicy": "IfNotPresent",
	}
	if c.pod.options.ForcePull {
		spec["imagePullPolicy"] = "Always"
	}
	if c == c.pod.job {
		spec["volumeMounts"] = []any{
			map[string]any{"name": "workspace", "mountPath": c.input.WorkingDir},
			map[string]any{"name": "act", "mountPath": c.GetActPath()},
		}
	}
	if err := c.parseOptions(); err != nil {
		return nil, fmt.Errorf("service %s: %w", c.name, err)
	}
	if c.probe != nil {
		spec["readinessProbe"] = c.probe
	}
	return spec, nil
}

// escapeExpansion keeps kubelet from expanding $(VAR) and $$ in command and args, which docker passes as is.
func escapeExpansion(args []string) []string {
	args = slices.Clone(args)
	for index := range args {
		args[index] = strings.ReplaceAll(args[index], "$", "$$")
	}
	return args
}

// parseOptions turns docker health options, the only ones a pod supports, into a readiness probe.
func (c *kubernetesContainer) parseOptions() error {
	args, err := shellquote.Split(c.input.WorkflowOptions)
	if err != nil {
		return err
	}
	flags := pflag.NewFlagSet("options", pflag.ContinueOnError)
	flags.SetOutput(io.Discard)
	command := flags.String("health-cmd", "", "")
	interval := flags.Duration("health-interval", 0, "")
	timeout := flags.Duration("health-timeout", 0, "")
	flags.Int("health-retries", 0, "")
	flags.Duration("health-start-period", 0, "")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("options %q: %w, only --health-* options are supported on Kubernetes", c.input.WorkflowOptions, err)
	}
	if *command == "" {
		return nil
	}
	c.probe = map[string]any{"exec": map[string]any{"command": []string{"/bin/sh", "-c", *command}}}
	if *interval > 0 {
		c.probe["periodSeconds"] = max(1, int(interval.Seconds()))
	}
	if *timeout > 0 {
		c.probe["timeoutSeconds"] = max(1, int(timeout.Seconds()))
	}
	return nil
}

func (p *KubernetesPod) get(ctx context.Context) (*podStatus, error) {
	var pod podStatus
	return &pod, p.client.do(ctx, http.MethodGet, p.client.path("pods/"+p.name), nil, &pod)
}

func (p *KubernetesPod) waitStarted(ctx context.Context) error {
	var unscheduled string
	for delay := 100 * time.Millisecond; ; delay = min(2*delay, time.Second) {
		pod, err := p.get(ctx)
		if err != nil {
			return err
		}
		if pod.Status.Phase == "Failed" {
			return fmt.Errorf("pod %s failed: %s", p.name, cmp.Or(pod.Status.Message, pod.Status.Reason))
		}
		started := len(pod.Status.ContainerStatuses) == len(pod.Spec.Containers)
		for _, status := range append(pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses...) {
			if waiting := status.State.Waiting; waiting != nil {
				switch waiting.Reason {
				case "ErrImagePull", "ImagePullBackOff", "InvalidImageName", "ErrImageNeverPull", "CreateContainerConfigError", "CreateContainerError", "RunContainerError":
					return fmt.Errorf("container %s of pod %s: %s: %s", status.Name, p.name, waiting.Reason, waiting.Message)
				}
				started = false
			}
		}
		if started {
			return nil
		}
		for _, condition := range pod.Status.Conditions {
			if condition.Type == "PodScheduled" && condition.Status == "False" && condition.Message != unscheduled {
				unscheduled = condition.Message
				common.Logger(ctx).Infof("Waiting for pod %s to be scheduled: %s", p.name, unscheduled)
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

func (c *kubernetesContainer) Create(_, _ []string) common.Executor {
	return func(ctx context.Context) error {
		if c != c.pod.job {
			return nil
		}
		return c.pod.create(ctx)
	}
}

func (c *kubernetesContainer) Start(_ bool) common.Executor {
	return func(ctx context.Context) error {
		if c != c.pod.job {
			return nil
		}
		if err := c.pod.waitStarted(ctx); err != nil {
			return err
		}
		var stdout bytes.Buffer
		if err := c.pod.client.exec(ctx, c.pod.name, c.name, []string{"sh", "-c", `printf '%s\n%s' "$(uname -m)" "$PATH"`}, nil, &stdout, io.Discard); err != nil {
			return err
		}
		machine, path, _ := strings.Cut(stdout.String(), "\n")
		c.pod.arch, c.pod.path = goArchToActionArch(machine), path
		return nil
	}
}

func (*kubernetesContainer) Pull(bool) common.Executor {
	return func(context.Context) error { return nil }
}

// Exec feeds the command to sh on stdin, keeping its env out of the request URL and the process list.
func (c *kubernetesContainer) Exec(command []string, env map[string]string, _, workdir string) common.Executor {
	return func(ctx context.Context) error {
		script := "exec " + shellquote.Join(append(append([]string{"env", "--"}, getEnvListFromMap(env)...), command...)...)
		if workdir != "" {
			script = "cd " + shellquote.Join(workdir) + " || exit\n" + script
		}
		defer common.FlushWriter(c.input.Stdout)
		defer common.FlushWriter(c.input.Stderr)
		return c.pod.client.exec(ctx, c.pod.name, c.name, []string{"sh"}, strings.NewReader(script), c.input.Stdout, c.input.Stderr)
	}
}

func (c *kubernetesContainer) extract(ctx context.Context, destPath string, write func(io.Writer) error) error {
	reader, writer := io.Pipe()
	defer reader.Close()
	go func() { writer.CloseWithError(write(writer)) }()
	var stderr bytes.Buffer
	if err := c.pod.client.exec(ctx, c.pod.name, c.name, []string{"sh", "-c", `mkdir -p "$1" && tar -xf - -C "$1"`, "sh", destPath}, reader, io.Discard, &stderr); err != nil {
		return fmt.Errorf("extract to %s: %w: %s", destPath, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func (c *kubernetesContainer) Copy(destPath string, files ...*FileEntry) common.Executor {
	return func(ctx context.Context) error {
		return c.extract(ctx, destPath, func(writer io.Writer) error { return writeFilesTar(ctx, writer, 0, 0, files...) })
	}
}

func (c *kubernetesContainer) CopyDir(destPath, srcPath string, useGitIgnore, skipGitDir bool) common.Executor {
	return func(ctx context.Context) error {
		return c.extract(ctx, "/", func(writer io.Writer) error {
			return writeDirTar(ctx, writer, destPath, srcPath, useGitIgnore, skipGitDir, 0, 0)
		})
	}
}

// GetContainerArchive streams the archive, checking the path first as tar writes an empty archive for a missing one.
func (c *kubernetesContainer) GetContainerArchive(ctx context.Context, srcPath string) (io.ReadCloser, error) {
	reader, writer := io.Pipe()
	go func() {
		var stderr bytes.Buffer
		err := c.pod.client.exec(ctx, c.pod.name, c.name, []string{"sh", "-c", `test -e "$1" -o -L "$1" && exec tar -cf - -C "$2" "$3"`, "sh", srcPath, path.Dir(srcPath), path.Base(srcPath)}, nil, writer, &stderr)
		if err != nil {
			err = fmt.Errorf("archive %s: %w: %s", srcPath, err, strings.TrimSpace(stderr.String()))
		}
		writer.CloseWithError(err)
	}()
	buffered := bufio.NewReader(reader)
	if _, err := buffered.Peek(1); err != nil {
		reader.Close()
		return nil, err
	}
	return struct {
		io.Reader
		io.Closer
	}{buffered, reader}, nil
}

func (c *kubernetesContainer) UpdateFromEnv(srcPath string, env *map[string]string) common.Executor {
	return parseEnvFile(c, srcPath, env)
}

func (c *kubernetesContainer) UpdateFromImageEnv(env *map[string]string) common.Executor {
	return func(context.Context) error {
		(*env)["PATH"] = cmp.Or((*env)["PATH"], c.pod.path)
		return nil
	}
}

func (c *kubernetesContainer) Inspect(ctx context.Context) (*Info, error) {
	if c.pod.name == "" {
		return nil, fmt.Errorf("container %s %w", c.name, ErrContainerNotFound)
	}
	pod, err := c.pod.get(ctx)
	if isNotFound(err) {
		return nil, fmt.Errorf("pod %s %w", c.pod.name, ErrContainerNotFound)
	} else if err != nil {
		return nil, err
	}
	info := &Info{ID: c.pod.name + "/" + c.name, State: "created", Health: HealthNone, Ports: map[string]string{}}
	for port := range c.input.ExposedPorts {
		info.Ports[port.Port()] = port.Port() // services share the job's network, so a port is reached as is
	}
	for _, status := range pod.Status.ContainerStatuses {
		if status.Name != c.name {
			continue
		}
		if status.State.Running != nil {
			info.State = StateRunning
		} else if terminated := status.State.Terminated; terminated != nil {
			info.State, info.ExitCode, info.HealthOutput = "exited", terminated.ExitCode, terminated.Message
		}
		if c.probe != nil {
			switch {
			case status.Ready:
				info.Health = HealthHealthy
			case info.State == "exited":
				info.Health = HealthUnhealthy
			default:
				info.Health = HealthStarting
			}
		}
	}
	return info, nil
}

func (c *kubernetesContainer) DumpLogs(ctx context.Context) error {
	if c.pod.name == "" {
		return fmt.Errorf("container %s %w", c.name, ErrContainerNotFound)
	}
	defer common.FlushWriter(c.input.Stdout)
	return c.pod.client.do(ctx, http.MethodGet, c.pod.client.path("pods/"+c.pod.name+"/log?container="+c.name), nil, c.input.Stdout)
}

func (c *kubernetesContainer) Remove() common.Executor {
	return func(ctx context.Context) error {
		if c != c.pod.job {
			return nil
		}
		return errors.Join(c.pod.delete(ctx, "pods", c.pod.name), c.pod.delete(ctx, "secrets", c.pod.secretName))
	}
}

func (p *KubernetesPod) delete(ctx context.Context, resource, name string) error {
	if name == "" {
		return nil
	}
	if err := p.client.do(ctx, http.MethodDelete, p.client.path(resource+"/"+name), nil, nil); err != nil && !isNotFound(err) {
		return fmt.Errorf("delete %s %s: %w", resource, name, err)
	}
	return nil
}

func (c *kubernetesContainer) Close() common.Executor {
	return func(context.Context) error {
		c.pod.client.http.CloseIdleConnections()
		return nil
	}
}

func (c *kubernetesContainer) ReplaceLogWriter(stdout, stderr io.Writer) (io.Writer, io.Writer) {
	oldStdout, oldStderr := c.input.Stdout, c.input.Stderr
	c.input.Stdout, c.input.Stderr = stdout, stderr
	return oldStdout, oldStderr
}

func (c *kubernetesContainer) GetRunnerContext(_ context.Context) map[string]any {
	return map[string]any{"os": "Linux", "arch": c.pod.arch, "temp": "/tmp", "tool_cache": DefaultToolCache}
}
