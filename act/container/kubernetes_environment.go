// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// The Kubernetes execution backend runs each Actions job as a single Pod and
// execs steps into a long-lived job container via the client-go remotecommand
// transport. Unlike the docker backend it carries NO build tag: client-go is
// pure Go and cross-platform, so the type is always compiled and only its
// selection is gated at runtime by the kubernetes.enabled config flag.

package container

import (
	"context"
	"errors"
	"io"

	"gitea.com/gitea/runner/act/common"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// errKubernetesNotImplemented is returned by the not-yet-wired methods of the
// Kubernetes backend skeleton. Subsequent phases replace these stubs with real
// pod-lifecycle and exec-in implementations.
var errKubernetesNotImplemented = errors.New("kubernetes execution environment: not implemented")

// KubernetesExecutionEnvironment is the ExecutionsEnvironment backend that runs
// a Gitea Actions job as one Kubernetes Pod. The embedded
// LinuxContainerEnvironmentExtensions supplies the Linux path/runner-context
// helpers (job containers are Linux-only in v1).
type KubernetesExecutionEnvironment struct {
	LinuxContainerEnvironmentExtensions

	input *NewContainerInput

	restConfig *rest.Config
	clientset  kubernetes.Interface
	namespace  string
}

// compile-time assertion that the backend satisfies the interface.
var _ ExecutionsEnvironment = (*KubernetesExecutionEnvironment)(nil)

// KubernetesEnvironmentConfig carries the cluster connection settings the
// backend needs to build a clientset. An empty Kubeconfig selects in-cluster
// config; a non-empty path is the out-of-cluster dev fallback.
type KubernetesEnvironmentConfig struct {
	Input      *NewContainerInput
	Namespace  string
	Kubeconfig string
}

// NewKubernetesExecutionEnvironment builds a Kubernetes-backed
// ExecutionsEnvironment. It resolves a rest.Config from the in-cluster
// ServiceAccount, falling back to an explicit kubeconfig path for
// out-of-cluster development.
func NewKubernetesExecutionEnvironment(cfg KubernetesEnvironmentConfig) (*KubernetesExecutionEnvironment, error) {
	restConfig, err := resolveKubernetesRestConfig(cfg.Kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &KubernetesExecutionEnvironment{
		input:      cfg.Input,
		restConfig: restConfig,
		clientset:  clientset,
		namespace:  cfg.Namespace,
	}, nil
}

// resolveKubernetesRestConfig returns the in-cluster config when kubeconfig is
// empty, otherwise loads the supplied kubeconfig file.
func resolveKubernetesRestConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func (e *KubernetesExecutionEnvironment) Create(_, _ []string) common.Executor {
	return func(_ context.Context) error {
		_ = e.input
		return errKubernetesNotImplemented
	}
}

func (e *KubernetesExecutionEnvironment) ConnectToNetwork(_ string) common.Executor {
	// Containers in a Pod share a single network namespace, so there is no
	// per-container network to join. This is intentionally a no-op.
	return func(_ context.Context) error { return nil }
}

func (e *KubernetesExecutionEnvironment) Copy(_ string, _ ...*FileEntry) common.Executor {
	return notImplemented()
}

func (e *KubernetesExecutionEnvironment) CopyTarStream(_ context.Context, _ string, _ io.Reader) error {
	return errKubernetesNotImplemented
}

func (e *KubernetesExecutionEnvironment) CopyDir(_, _ string, _ bool) common.Executor {
	return notImplemented()
}

func (e *KubernetesExecutionEnvironment) GetContainerArchive(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, errKubernetesNotImplemented
}

func (e *KubernetesExecutionEnvironment) Pull(_ bool) common.Executor {
	// The kubelet pulls images per the Pod's imagePullPolicy; the runner never
	// pulls directly. This is a no-op by design.
	return func(_ context.Context) error { return nil }
}

func (e *KubernetesExecutionEnvironment) Start(_ bool) common.Executor {
	return func(_ context.Context) error {
		if e.clientset == nil || e.namespace == "" {
			return errKubernetesNotImplemented
		}
		return errKubernetesNotImplemented
	}
}

func (e *KubernetesExecutionEnvironment) Exec(_ []string, _ map[string]string, _, _ string) common.Executor {
	return func(_ context.Context) error {
		_ = e.restConfig
		return errKubernetesNotImplemented
	}
}

func (e *KubernetesExecutionEnvironment) UpdateFromEnv(_ string, _ *map[string]string) common.Executor {
	return notImplemented()
}

func (e *KubernetesExecutionEnvironment) UpdateFromImageEnv(_ *map[string]string) common.Executor {
	return notImplemented()
}

func (e *KubernetesExecutionEnvironment) Remove() common.Executor {
	return notImplemented()
}

func (e *KubernetesExecutionEnvironment) Close() common.Executor {
	return func(_ context.Context) error { return nil }
}

func (e *KubernetesExecutionEnvironment) ReplaceLogWriter(stdout, stderr io.Writer) (io.Writer, io.Writer) {
	return stdout, stderr
}

func notImplemented() common.Executor {
	return func(_ context.Context) error {
		return errKubernetesNotImplemented
	}
}
