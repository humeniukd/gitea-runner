// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	"fmt"

	"gitea.com/gitea/runner/act/common"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SweepDanglingJobPods runs the startup backstop sweep (C2/RC2): it builds a
// clientset from the config and deletes any leaked job pods whose owning runner
// Pod no longer exists. It is best-effort; callers should log and continue on
// error rather than block startup.
func SweepDanglingJobPods(ctx context.Context, cfg *KubernetesConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	restConfig, err := resolveKubernetesRestConfig(cfg.Kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	return sweepDanglingJobPods(ctx, clientset, cfg.Namespace)
}

// sweepDanglingJobPods deletes ONLY job pods whose Pod-kind ownerReference no
// longer resolves to a live runner Pod (a "dangling owner"), along with their
// per-job Secret and ConfigMap. This is the narrow backstop for leaks that
// ownerReference GC has not yet reaped (e.g. a hard runner crash) — it is NOT a
// blanket label delete: a pod owned by a LIVE sibling/rolling-update runner has
// a resolvable owner and is spared, keeping the sweep concurrency-safe under
// replicas>1 (RC2). Pods carrying the backend label but no Pod owner are left
// untouched (their provenance cannot be determined).
func sweepDanglingJobPods(ctx context.Context, clientset kubernetes.Interface, namespace string) error {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: kubernetesBackendLabel + "=" + kubernetesBackendValue,
	})
	if err != nil {
		return fmt.Errorf("list job pods in %s: %w", namespace, err)
	}

	logger := common.Logger(ctx)
	var firstErr error
	for i := range pods.Items {
		pod := &pods.Items[i]
		owner := podOwnerReference(pod)
		if owner == nil {
			// No Pod owner to resolve; leave it alone (conservative).
			continue
		}
		if ownerResolvable(ctx, clientset, namespace, owner) {
			// Owning runner is alive — a live sibling's job pod. Spare it.
			continue
		}

		logger.Warnf("Sweeping dangling job pod %s/%s (owner %s/%s no longer resolvable)", namespace, pod.Name, owner.Name, owner.UID)
		if err := deleteJobPodResources(ctx, clientset, namespace, pod.Name); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// podOwnerReference returns the Pod-kind ownerReference of a job pod, or nil.
func podOwnerReference(pod *corev1.Pod) *metav1.OwnerReference {
	for i := range pod.OwnerReferences {
		if pod.OwnerReferences[i].Kind == "Pod" {
			return &pod.OwnerReferences[i]
		}
	}
	return nil
}

// ownerResolvable reports whether the owning runner Pod still exists with the
// same UID. A NotFound (or a name reused with a different UID) means the owner
// is gone — the job pod is dangling.
func ownerResolvable(ctx context.Context, clientset kubernetes.Interface, namespace string, owner *metav1.OwnerReference) bool {
	ownerPod, err := clientset.CoreV1().Pods(namespace).Get(ctx, owner.Name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return ownerPod.UID == owner.UID
}

// deleteJobPodResources deletes a job pod and its deterministic per-job Secret
// and ConfigMap, tolerating already-absent objects.
func deleteJobPodResources(ctx context.Context, clientset kubernetes.Interface, namespace, podName string) error {
	var firstErr error
	record := func(err error) {
		if err != nil && !apierrors.IsNotFound(err) && firstErr == nil {
			firstErr = err
		}
	}

	record(clientset.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{}))
	record(clientset.CoreV1().Secrets(namespace).Delete(ctx, perJobSecretName(podName), metav1.DeleteOptions{}))
	record(clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, perJobConfigMapName(podName), metav1.DeleteOptions{}))
	return firstErr
}
