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
	"os"
	"time"

	"gitea.com/gitea/runner/act/common"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// podReadyPollInterval is how often Start polls a submitted Pod for readiness.
const podReadyPollInterval = 500 * time.Millisecond

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
	builder := newPodBuilder(cfg.PodName, cfg.Config.Namespace, cfg.Config, cfg.MaxLifetime, owner)

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
// already-absent objects so cleanup is idempotent (R2).
func (p *KubernetesPod) remove(ctx context.Context) error {
	var firstErr error
	if p.builder != nil && p.builder.name != "" {
		if err := p.clientset.CoreV1().Pods(p.namespace).Delete(ctx, p.builder.name, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			firstErr = fmt.Errorf("delete job pod %s/%s: %w", p.namespace, p.builder.name, err)
		}
	}
	for _, secret := range p.perJobSecrets {
		if err := p.clientset.CoreV1().Secrets(p.namespace).Delete(ctx, secret, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
			if firstErr == nil {
				firstErr = fmt.Errorf("delete per-job secret %s/%s: %w", p.namespace, secret, err)
			}
		}
	}
	return firstErr
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
// container; service Removes are no-ops since they share the Pod.
func (e *KubernetesExecutionEnvironment) Remove() common.Executor {
	return func(ctx context.Context) error {
		if e.role != roleJob {
			return nil
		}
		return e.pod.remove(ctx)
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

// serviceContainerName sanitises an act service container name into a DNS-label
// compatible container name. act already produces lowercase hyphenated names.
func serviceContainerName(name string) string {
	if name == "" {
		return "service"
	}
	return name
}
