// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// The Kubernetes execution backend runs each Actions job as a single Pod and
// execs steps into a long-lived job container via the client-go remotecommand
// transport. Unlike the docker backend it carries NO build tag: client-go is
// pure Go and cross-platform, so the type is always compiled and only its
// selection is gated at runtime by the kubernetes.enabled config flag.
//
// One job maps to one KubernetesPod coordinator that owns the clientset and the
// shared podBuilder. The job container and each service sidecar are exposed as
// KubernetesExecutionEnvironment values that implement the Container interface;
// they register their spec into the shared builder on Create, and the job
// container's Start submits the whole Pod (service Starts are no-ops).

package container

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"gitea.com/gitea/runner/act/common"

	"github.com/avast/retry-go/v5"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// podReadyPollInterval is how often Start polls a submitted Pod for readiness.
const podReadyPollInterval = 500 * time.Millisecond

// kubernetesDeleteRetryAttempts bounds the retry loop around a pod/secret
// delete so transient control-plane instability cannot strand a finished job's
// Pod. 5 attempts with exponential backoff (500ms→8s, capped) lands ~16s worst
// case — short enough that a failed job still tears down promptly, long enough
// to ride out a typical API server hiccup.
const (
	kubernetesDeleteRetryAttempts = 5
	kubernetesDeleteRetryMaxDelay = 8 * time.Second
)

type kubernetesContainerRole int

const (
	roleJob kubernetesContainerRole = iota
	roleService
)

// KubernetesPodConfig carries everything the coordinator needs to assemble and
// submit the job Pod.
type KubernetesPodConfig struct {
	Config      *KubernetesConfig
	PodName     string
	MaxLifetime time.Duration
}

// KubernetesPod coordinates the single Pod backing one Actions job. It holds the
// clientset and the shared builder; the job/service environments it vends share
// this state so the job container's Start submits one Pod containing the service
// sidecars.
type KubernetesPod struct {
	config    *KubernetesConfig
	clientset kubernetes.Interface
	rest      *rest.Config
	namespace string
	builder   *podBuilder

	podReadyTimeout time.Duration

	submitted     *corev1.Pod
	perJobSecrets []string
}

// NewKubernetesPod builds the coordinator: it resolves a clientset (in-cluster
// with a kubeconfig dev fallback) and seeds the pod builder with identity,
// orphan-reaping primitives, and the own-Pod ownerReference (RC2).
func NewKubernetesPod(cfg KubernetesPodConfig) (*KubernetesPod, error) {
	if cfg.Config == nil {
		return nil, errors.New("kubernetes pod: nil config")
	}

	restConfig, err := resolveKubernetesRestConfig(cfg.Config.Kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return newKubernetesPodWithClientset(clientset, restConfig, cfg), nil
}

// newKubernetesPodWithClientset is the testable seam: it builds the coordinator
// from an already-constructed clientset (a fake clientset in unit tests).
func newKubernetesPodWithClientset(clientset kubernetes.Interface, restConfig *rest.Config, cfg KubernetesPodConfig) *KubernetesPod {
	owner := ownerReferenceFromEnv(os.Getenv("POD_NAME"), os.Getenv("POD_UID"))
	// act's jobContainerName() is mixed-case and ~128 chars; the API server
	// requires the Pod name to be an RFC 1123 name. Sanitize before it reaches
	// the builder (and thus the per-job Secret/ConfigMap names derived from it).
	podName := sanitizeRFC1123Name(cfg.PodName)
	builder := newPodBuilder(podName, cfg.Config.Namespace, cfg.Config, cfg.MaxLifetime, owner)

	return &KubernetesPod{
		config:          cfg.Config,
		clientset:       clientset,
		rest:            restConfig,
		namespace:       cfg.Config.Namespace,
		builder:         builder,
		podReadyTimeout: cfg.Config.PodReadyTimeout,
	}
}

// NewJobContainer returns the environment representing the long-lived job
// container act execs steps into.
func (p *KubernetesPod) NewJobContainer(input *NewContainerInput) *KubernetesExecutionEnvironment {
	return &KubernetesExecutionEnvironment{pod: p, input: input, role: roleJob, resources: p.config.Resources.Job}
}

// NewServiceContainer returns the environment representing one service sidecar.
func (p *KubernetesPod) NewServiceContainer(input *NewContainerInput) *KubernetesExecutionEnvironment {
	return &KubernetesExecutionEnvironment{pod: p, input: input, role: roleService, resources: p.config.Resources.Service}
}

// registerPerJobSecret records a Secret name the coordinator must delete on
// cleanup alongside the Pod.
func (p *KubernetesPod) registerPerJobSecret(name string) {
	p.perJobSecrets = append(p.perJobSecrets, name)
}

// submit creates the Pod and blocks until it is Ready or the timeout elapses.
func (p *KubernetesPod) submit(ctx context.Context) error {
	pod := p.builder.build()
	created, err := p.clientset.CoreV1().Pods(p.namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create job pod %s/%s: %w", p.namespace, pod.Name, err)
	}
	p.submitted = created
	common.Logger(ctx).Infof("Created job pod %s/%s", p.namespace, created.Name)
	return p.waitReady(ctx, created.Name)
}

// waitReady polls the Pod until its Ready condition is true. A Failed pod or an
// elapsed pod_ready_timeout is a hard error so a stuck start surfaces promptly.
func (p *KubernetesPod) waitReady(ctx context.Context, name string) error {
	timeout := p.podReadyTimeout
	if timeout <= 0 {
		timeout = defaultPodReadyTimeout
	}
	deadline := time.Now().Add(timeout)

	for {
		pod, err := p.clientset.CoreV1().Pods(p.namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get job pod %s/%s: %w", p.namespace, name, err)
		}
		if pod.Status.Phase == corev1.PodFailed {
			return fmt.Errorf("job pod %s/%s entered phase %s", p.namespace, name, pod.Status.Phase)
		}
		if isPodReady(pod) {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("job pod %s/%s not ready within %s", p.namespace, name, timeout)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(podReadyPollInterval):
		}
	}
}

// remove deletes the Pod and any registered per-job Secrets, tolerating
// already-absent objects so cleanup is idempotent (R2). Transient API-server
// errors (timeouts, throttling, connection resets during a control-plane blip)
// are retried with bounded backoff so a failed job does not strand its Pod
// forever; permanent errors (NotFound, Forbidden) fail fast.
func (p *KubernetesPod) remove(ctx context.Context) error {
	var firstErr error
	if p.builder != nil && p.builder.name != "" {
		name := p.builder.name
		if err := deleteWithRetry(ctx, p.namespace, "job pod", name, func(ctx context.Context) error {
			return p.clientset.CoreV1().Pods(p.namespace).Delete(ctx, name, metav1.DeleteOptions{})
		}); err != nil && !apierrors.IsNotFound(err) {
			firstErr = err
		}
	}
	for _, secret := range p.perJobSecrets {
		if err := deleteWithRetry(ctx, p.namespace, "per-job secret", secret, func(ctx context.Context) error {
			return p.clientset.CoreV1().Secrets(p.namespace).Delete(ctx, secret, metav1.DeleteOptions{})
		}); err != nil && !apierrors.IsNotFound(err) && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// deleteWithRetry deletes a single Kubernetes object, retrying transient
// API-server errors (throttling, server timeouts, connection resets) with
// bounded exponential backoff so a control-plane blip cannot strand a finished
// job's Pod. Permanent errors (NotFound, Forbidden) abort immediately — retrying
// cannot fix them. The returned error is already namespaced/qualified.
func deleteWithRetry(ctx context.Context, namespace, kind, name string, doDelete func(context.Context) error) error {
	err := retry.New(
		retry.Context(ctx),
		retry.Attempts(kubernetesDeleteRetryAttempts),
		retry.DelayType(retry.BackOffDelay),
		retry.MaxDelay(kubernetesDeleteRetryMaxDelay),
		retry.LastErrorOnly(true),
		retry.RetryIf(isTransientDeleteError),
	).Do(func() error { return doDelete(ctx) })
	if err == nil {
		return nil
	}
	if isTransientDeleteError(err) {
		return fmt.Errorf("delete %s %s/%s failed after %d attempts: %w", kind, namespace, name, kubernetesDeleteRetryAttempts, err)
	}
	return fmt.Errorf("delete %s %s/%s: %w", kind, namespace, name, err)
}

// isTransientDeleteError reports whether a Delete error is worth retrying: API
// server throttling (429), server timeouts, and transport-level timeouts.
// NotFound and Forbidden are deliberately excluded — the object is already
// absent (success), or retrying cannot fix RBAC.
func isTransientDeleteError(err error) bool {
	if err == nil {
		return false
	}
	if apierrors.IsTooManyRequests(err) || apierrors.IsServerTimeout(err) || apierrors.IsTimeout(err) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// removeWithPolicy applies the configured pod_cleanup policy to a completed job:
// it either deletes the Pod (and per-job Secrets) or retains it for auditing,
// marking a retained Pod so the dangling-owner sweep spares it. jobSucceeded
// carries the job conclusion the cleanup decision depends on.
func (p *KubernetesPod) removeWithPolicy(ctx context.Context, jobSucceeded bool) error {
	if !shouldRetainPod(p.cleanupPolicy(), jobSucceeded) {
		return p.remove(ctx)
	}

	result := "failure"
	if jobSucceeded {
		result = "success"
	}
	name := ""
	if p.builder != nil {
		name = p.builder.name
	}
	if err := p.markRetained(ctx, name); err != nil {
		// Marking is best-effort: even if the label patch fails the operator
		// asked us to KEEP the pod, so we must not fall through to deletion.
		common.Logger(ctx).Warnf("failed to mark retained job pod %s/%s (sweep may later reap it): %v", p.namespace, name, err)
	}
	common.Logger(ctx).Infof("retaining job pod %s/%s per kubernetes.pod_cleanup=%s (job result=%s)", p.namespace, name, p.cleanupPolicy(), result)
	return nil
}

// cleanupPolicy returns the effective pod_cleanup policy, defaulting to "always"
// so a zero-value config (e.g. a directly-constructed coordinator in tests)
// preserves the historical delete-on-completion behavior.
func (p *KubernetesPod) cleanupPolicy() string {
	if p.config == nil || p.config.PodCleanup == "" {
		return KubernetesPodCleanupAlways
	}
	return p.config.PodCleanup
}

// markRetained patches the retain marker onto the live Pod so the startup
// dangling-owner sweep skips it. A merge patch is idempotent and tolerates an
// already-absent pod (NotFound).
func (p *KubernetesPod) markRetained(ctx context.Context, name string) error {
	if name == "" {
		return nil
	}
	patch := fmt.Appendf(nil, `{"metadata":{"labels":{%q:%q}}}`, kubernetesRetainLabel, kubernetesRetainValue)
	_, err := p.clientset.CoreV1().Pods(p.namespace).Patch(ctx, name, k8stypes.MergePatchType, patch, metav1.PatchOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

// shouldRetainPod decides whether a completed job Pod is kept per the policy:
// "never" keeps all, "on-failure" keeps only failed jobs, anything else
// ("always", the default) keeps none.
func shouldRetainPod(policy string, jobSucceeded bool) bool {
	switch policy {
	case KubernetesPodCleanupNever:
		return true
	case KubernetesPodCleanupOnFailure:
		return !jobSucceeded
	default:
		return false
	}
}

// AddServiceAlias maps a service id to localhost in the Pod's hostAliases so a
// job step can reach the service at both servicename:port and localhost:port.
func (p *KubernetesPod) AddServiceAlias(serviceID string) {
	p.builder.addHostAlias("127.0.0.1", serviceID)
}

// WaitForServicesReady blocks until every service sidecar reports Ready, so the
// first job step never races a service that is still starting (C3). Native
// sidecars surface in the Pod's initContainerStatuses.
func (p *KubernetesPod) WaitForServicesReady(ctx context.Context) common.Executor {
	return func(ctx context.Context) error {
		names := p.builder.sidecarContainerNames()
		if len(names) == 0 {
			return nil
		}
		want := make(map[string]struct{}, len(names))
		for _, n := range names {
			want[n] = struct{}{}
		}

		timeout := p.podReadyTimeout
		if timeout <= 0 {
			timeout = defaultPodReadyTimeout
		}
		deadline := time.Now().Add(timeout)

		for {
			pod, err := p.clientset.CoreV1().Pods(p.namespace).Get(ctx, p.builder.name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("get job pod %s/%s: %w", p.namespace, p.builder.name, err)
			}
			if sidecarsReady(pod, want) {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("service sidecars not ready within %s", timeout)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(podReadyPollInterval):
			}
		}
	}
}

// sidecarsReady reports whether every named sidecar has Ready=true.
func sidecarsReady(pod *corev1.Pod, want map[string]struct{}) bool {
	ready := 0
	for _, st := range pod.Status.InitContainerStatuses {
		if _, ok := want[st.Name]; ok && st.Ready {
			ready++
		}
	}
	return ready == len(want)
}

// isPodReady reports whether the Pod's Ready condition is true.
func isPodReady(pod *corev1.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady {
			return cond.Status == corev1.ConditionTrue
		}
	}
	return false
}

// defaultPodReadyTimeout mirrors the config default so a zero timeout (e.g. a
// directly-constructed coordinator in tests) still bounds the wait.
const defaultPodReadyTimeout = 5 * time.Minute

// KubernetesExecutionEnvironment is one container's view of the job Pod (the job
// container or a service sidecar). The embedded LinuxContainerEnvironmentExtensions
// supplies the Linux path/runner-context helpers (job containers are Linux-only).
type KubernetesExecutionEnvironment struct {
	LinuxContainerEnvironmentExtensions

	pod       *KubernetesPod
	input     *NewContainerInput
	role      kubernetesContainerRole
	resources KubernetesResources
	forcePull bool
}

// compile-time assertion that the backend satisfies the interface.
var _ ExecutionsEnvironment = (*KubernetesExecutionEnvironment)(nil)

// GetRunnerContext overrides the embedded LinuxContainerEnvironmentExtensions
// implementation so the runner.arch value is resolved locally (RunnerArchLocal)
// instead of by pinging a host docker daemon (RunnerArch). The pod backend has
// no host docker socket, and the expression evaluator calls GetRunnerContext
// many times per run — deferring to the docker resolver would emit a flood of
// non-fatal "docker daemon ping ... /var/run/docker.sock: no such file" lines.
func (e *KubernetesExecutionEnvironment) GetRunnerContext(_ context.Context) map[string]any {
	return map[string]any{
		"os":         "Linux",
		"arch":       RunnerArchLocal(),
		"temp":       "/tmp",
		"tool_cache": "/opt/hostedtoolcache",
	}
}

func resolveKubernetesRestConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// Create registers this container's spec into the shared pod builder. It does
// not submit the Pod — that happens in the job container's Start.
func (e *KubernetesExecutionEnvironment) Create(capAdd, capDrop []string) common.Executor {
	return func(_ context.Context) error {
		switch e.role {
		case roleJob:
			spec := kubernetesContainerSpec(kubernetesJobContainerName, e.input, e.resources, capAdd, capDrop, e.forcePull)
			e.pod.builder.setJobContainer(spec)
		case roleService:
			e.pod.builder.addSidecar(serviceSidecarSpec(e.input, e.resources, capAdd, capDrop, e.forcePull))
		}
		return nil
	}
}

// ConnectToNetwork is a no-op: containers in a Pod share one network namespace.
func (e *KubernetesExecutionEnvironment) ConnectToNetwork(_ string) common.Executor {
	return func(_ context.Context) error { return nil }
}

// Pull is a no-op against the API server: the kubelet pulls images per the
// Pod's imagePullPolicy. It records forcePull so Create can set imagePullPolicy
// to Always (Pull runs before Create in the start pipeline).
func (e *KubernetesExecutionEnvironment) Pull(forcePull bool) common.Executor {
	return func(_ context.Context) error {
		e.forcePull = forcePull
		return nil
	}
}

// Start submits the Pod and waits for readiness when called on the job
// container. Service-sidecar Starts are no-ops: the sidecars are submitted as
// part of the same Pod.
func (e *KubernetesExecutionEnvironment) Start(_ bool) common.Executor {
	return func(ctx context.Context) error {
		if e.role != roleJob {
			return nil
		}
		return e.pod.submit(ctx)
	}
}

// Exec runs a step command in the job container via remotecommand, streaming
// output to this container's log writers. Only the job container is exec'd into.
func (e *KubernetesExecutionEnvironment) Exec(command []string, env map[string]string, user, workdir string) common.Executor {
	return func(ctx context.Context) error {
		wrapped := buildExecCommand(command, env, user, workdir)
		return e.pod.execInJobContainer(ctx, wrapped, e.input.Stdout, e.input.Stderr)
	}
}

// Remove deletes the whole Pod (and per-job Secrets) when called on the job
// container; service Removes are no-ops since they share the Pod. It always
// deletes regardless of pod_cleanup — it is the failed-start backstop and the
// interface-level teardown; the policy-aware path is RemoveWithRetention.
func (e *KubernetesExecutionEnvironment) Remove() common.Executor {
	return func(ctx context.Context) error {
		if e.role != roleJob {
			return nil
		}
		return e.pod.remove(ctx)
	}
}

// RemoveWithRetention deletes or retains the job Pod per the configured
// pod_cleanup policy and the job conclusion (jobSucceeded). Service-sidecar
// calls are no-ops since they share the Pod. Callers that know the job result
// (the run-context teardown) use this instead of Remove so completed pods can
// be kept for auditing.
func (e *KubernetesExecutionEnvironment) RemoveWithRetention(jobSucceeded bool) common.Executor {
	return func(ctx context.Context) error {
		if e.role != roleJob {
			return nil
		}
		return e.pod.removeWithPolicy(ctx, jobSucceeded)
	}
}

func (e *KubernetesExecutionEnvironment) Close() common.Executor {
	return func(_ context.Context) error { return nil }
}

// ReplaceLogWriter swaps the writers Exec streams to, returning the previous
// pair so callers can restore them — mirroring the docker backend.
func (e *KubernetesExecutionEnvironment) ReplaceLogWriter(stdout, stderr io.Writer) (io.Writer, io.Writer) {
	prevOut, prevErr := e.input.Stdout, e.input.Stderr
	e.input.Stdout = stdout
	e.input.Stderr = stderr
	return prevOut, prevErr
}

// Inspect snapshots this container's status from the shared job Pod. A
// not-yet-submitted or already-deleted pod is reported as ErrContainerNotFound
// so the docker wait-for-healthy path can treat it as gone.
func (e *KubernetesExecutionEnvironment) Inspect(ctx context.Context) (*Info, error) {
	if e.pod == nil || e.pod.builder == nil || e.pod.builder.name == "" {
		return nil, fmt.Errorf("inspect kubernetes container: %w", ErrContainerNotFound)
	}
	pod, err := e.pod.clientset.CoreV1().Pods(e.pod.namespace).Get(ctx, e.pod.builder.name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil, fmt.Errorf("inspect kubernetes container: %w", ErrContainerNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("inspect kubernetes container %s/%s: %w", e.pod.namespace, e.pod.builder.name, err)
	}
	return kubernetesContainerInfo(pod, e.k8sContainerName(), e.input), nil
}

// DumpLogs copies this container's log so far onto its output writers. Used
// when a service is unhealthy so the job log still shows why the sidecar died.
func (e *KubernetesExecutionEnvironment) DumpLogs(ctx context.Context) error {
	if e.pod == nil || e.pod.builder == nil || e.pod.builder.name == "" {
		return fmt.Errorf("read kubernetes container logs: %w", ErrContainerNotFound)
	}
	req := e.pod.clientset.CoreV1().Pods(e.pod.namespace).GetLogs(e.pod.builder.name, &corev1.PodLogOptions{
		Container: e.k8sContainerName(),
	})
	stream, err := req.Stream(ctx)
	if apierrors.IsNotFound(err) {
		return fmt.Errorf("read kubernetes container logs: %w", ErrContainerNotFound)
	}
	if err != nil {
		return fmt.Errorf("read kubernetes container logs %s/%s: %w", e.pod.namespace, e.pod.builder.name, err)
	}
	defer stream.Close()

	out := io.Discard
	if e.input != nil && e.input.Stdout != nil {
		out = e.input.Stdout
	}
	_, err = io.Copy(out, stream)
	return err
}

// k8sContainerName is the Pod container this environment maps to: the
// long-lived job container, or the RFC 1123-sanitized service sidecar name.
func (e *KubernetesExecutionEnvironment) k8sContainerName() string {
	if e.role == roleJob {
		return kubernetesJobContainerName
	}
	sidecarName := ""
	if e.input != nil {
		sidecarName = e.input.Name
		if len(e.input.NetworkAliases) > 0 {
			sidecarName = e.input.NetworkAliases[0]
		}
	}
	return serviceContainerName(sidecarName)
}

// kubernetesContainerInfo maps a Pod + named container status onto the
// backend-agnostic Info the act engine consumes for job.container / job.services.
func kubernetesContainerInfo(pod *corev1.Pod, containerName string, input *NewContainerInput) *Info {
	info := &Info{
		ID:     string(pod.UID),
		Health: HealthNone,
		Ports:  map[string]string{},
	}
	if input != nil {
		for _, port := range sortedExposedPorts(input.ExposedPorts) {
			// Sidecars share the job Pod's network namespace, so the container
			// port is the address a step uses (via hostAliases → 127.0.0.1).
			key := strconv.Itoa(port)
			info.Ports[key] = key
		}
	}

	status := findContainerStatus(pod, containerName)
	if status == nil {
		info.State = string(pod.Status.Phase)
		return info
	}
	switch {
	case status.State.Running != nil:
		info.State = StateRunning
		if status.Ready {
			info.Health = HealthHealthy
		} else {
			info.Health = HealthStarting
		}
	case status.State.Waiting != nil:
		info.State = "created"
		info.Health = HealthStarting
		info.HealthOutput = status.State.Waiting.Message
	case status.State.Terminated != nil:
		info.State = "exited"
		info.ExitCode = int(status.State.Terminated.ExitCode)
		if status.State.Terminated.ExitCode != 0 {
			info.Health = HealthUnhealthy
		}
		info.HealthOutput = status.State.Terminated.Message
	}
	return info
}

// findContainerStatus looks up a container in both regular and init statuses:
// service sidecars are native sidecars and surface under InitContainerStatuses.
func findContainerStatus(pod *corev1.Pod, name string) *corev1.ContainerStatus {
	for i := range pod.Status.ContainerStatuses {
		if pod.Status.ContainerStatuses[i].Name == name {
			return &pod.Status.ContainerStatuses[i]
		}
	}
	for i := range pod.Status.InitContainerStatuses {
		if pod.Status.InitContainerStatuses[i].Name == name {
			return &pod.Status.InitContainerStatuses[i]
		}
	}
	return nil
}

// serviceContainerName converts an act service container name (mixed-case and
// ~128 chars, from createContainerName) into a valid RFC 1123 label (<=63,
// lowercase, alphanumeric start/end) for use as a Pod container name. A short
// hash keeps distinct services distinct after truncation.
func serviceContainerName(name string) string {
	if name == "" {
		return "service"
	}
	return sanitizeRFC1123Name(name)
}
