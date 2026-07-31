// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

const sweepNS = "gitea-e2e"

func backendLabels() map[string]string {
	return map[string]string{kubernetesBackendLabel: kubernetesBackendValue}
}

func jobPodOwnedBy(name, ownerName string, ownerUID types.UID) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: sweepNS,
			Labels:    backendLabels(),
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: "v1",
				Kind:       "Pod",
				Name:       ownerName,
				UID:        ownerUID,
			}},
		},
	}
}

func perJobSecret(podName string) *corev1.Secret {
	return &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: perJobSecretName(podName), Namespace: sweepNS}}
}

func podExists(t *testing.T, cs *fake.Clientset, name string) bool {
	t.Helper()
	_, err := cs.CoreV1().Pods(sweepNS).Get(context.Background(), name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return false
	}
	require.NoError(t, err)
	return true
}

func secretExists(t *testing.T, cs *fake.Clientset, podName string) bool {
	t.Helper()
	_, err := cs.CoreV1().Secrets(sweepNS).Get(context.Background(), perJobSecretName(podName), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return false
	}
	require.NoError(t, err)
	return true
}

// TestSweepDeletesDanglingOwner: a job pod whose owner Pod UID no longer
// resolves is reaped along with its per-job Secret (C2/RC2).
func TestSweepDeletesDanglingOwner(t *testing.T) {
	cs := fake.NewSimpleClientset(
		jobPodOwnedBy("job-dangling", "runner-dead", types.UID("dead-uid")),
		perJobSecret("job-dangling"),
	)

	require.NoError(t, sweepDanglingJobPods(context.Background(), cs, sweepNS))

	assert.False(t, podExists(t, cs, "job-dangling"), "dangling pod should be deleted")
	assert.False(t, secretExists(t, cs, "job-dangling"), "dangling pod's secret should be deleted")
}

// TestSweepSparesLiveSibling: a job pod owned by a runner Pod that is still
// alive (resolvable UID) must NOT be deleted — concurrency/rolling-update
// safety (RC2).
func TestSweepSparesLiveSibling(t *testing.T) {
	liveRunner := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "runner-live", Namespace: sweepNS, UID: types.UID("live-uid")}}
	cs := fake.NewSimpleClientset(
		liveRunner,
		jobPodOwnedBy("job-live", "runner-live", types.UID("live-uid")),
		perJobSecret("job-live"),
	)

	require.NoError(t, sweepDanglingJobPods(context.Background(), cs, sweepNS))

	assert.True(t, podExists(t, cs, "job-live"), "live sibling's pod must be spared")
	assert.True(t, secretExists(t, cs, "job-live"), "live sibling's secret must be spared")
}

// TestSweepNameReuseDifferentUIDIsDangling: an owner Pod name that exists but
// with a different UID (name reuse) counts as dangling.
func TestSweepNameReuseDifferentUIDIsDangling(t *testing.T) {
	reusedName := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "runner-x", Namespace: sweepNS, UID: types.UID("new-uid")}}
	cs := fake.NewSimpleClientset(
		reusedName,
		jobPodOwnedBy("job-old", "runner-x", types.UID("old-uid")),
		perJobSecret("job-old"),
	)

	require.NoError(t, sweepDanglingJobPods(context.Background(), cs, sweepNS))
	assert.False(t, podExists(t, cs, "job-old"), "pod owned by a reused-name/old-UID runner is dangling")
}

// TestSweepMixedFleet: dangling and live pods coexisting are handled correctly
// in one pass.
func TestSweepMixedFleet(t *testing.T) {
	liveRunner := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "runner-live", Namespace: sweepNS, UID: types.UID("live-uid")}}
	cs := fake.NewSimpleClientset(
		liveRunner,
		jobPodOwnedBy("job-live", "runner-live", types.UID("live-uid")),
		perJobSecret("job-live"),
		jobPodOwnedBy("job-dangling", "runner-dead", types.UID("dead-uid")),
		perJobSecret("job-dangling"),
	)

	require.NoError(t, sweepDanglingJobPods(context.Background(), cs, sweepNS))

	assert.True(t, podExists(t, cs, "job-live"))
	assert.False(t, podExists(t, cs, "job-dangling"))
}

// TestSweepSkipsRetainedPod: a pod the pod_cleanup policy chose to KEEP (carries
// the retain marker) must survive the sweep even when its owner is dangling,
// else a runner restart would reap exactly the pods kept for auditing.
func TestSweepSkipsRetainedPod(t *testing.T) {
	retained := jobPodOwnedBy("job-retained", "runner-dead", types.UID("dead-uid"))
	retained.Labels[kubernetesRetainLabel] = kubernetesRetainValue
	cs := fake.NewSimpleClientset(
		retained,
		perJobSecret("job-retained"),
	)

	require.NoError(t, sweepDanglingJobPods(context.Background(), cs, sweepNS))

	assert.True(t, podExists(t, cs, "job-retained"), "retained pod must survive the sweep")
	assert.True(t, secretExists(t, cs, "job-retained"), "retained pod's secret must survive the sweep")
}

// TestSweepLeavesOwnerlessPods: a labeled pod with no Pod owner is left
// untouched (provenance unknown).
func TestSweepLeavesOwnerlessPods(t *testing.T) {
	ownerless := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "job-ownerless", Namespace: sweepNS, Labels: backendLabels()}}
	cs := fake.NewSimpleClientset(ownerless)

	require.NoError(t, sweepDanglingJobPods(context.Background(), cs, sweepNS))
	assert.True(t, podExists(t, cs, "job-ownerless"), "ownerless pod must be left untouched")
}

func TestIsTransientDeleteError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"too many requests", apierrors.NewTooManyRequests("throttled", 10), true},
		{"server timeout", apierrors.NewTimeoutError("timed out", 0), true},
		{"not found", apierrors.NewNotFound(corev1.Resource("pods"), "job-x"), false},
		{"forbidden", apierrors.NewForbidden(corev1.Resource("pods"), "job-x", errors.New("rbac")), false},
		{"net timeout", &netTimeoutErr{}, true},
		{"plain error", errors.New("boom"), false},
		{"wrapped server timeout", fmt.Errorf("delete: %w", apierrors.NewTimeoutError("t", 0)), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, isTransientDeleteError(tc.err))
		})
	}
}

func TestDeleteWithRetry(t *testing.T) {
	// Succeeds on the first attempt.
	t.Run("immediate success", func(t *testing.T) {
		var calls int
		err := deleteWithRetry(context.Background(), nil, sweepNS, "job pod", "job-x", func(context.Context) error {
			calls++
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, 1, calls)
	})

	// A transient failure followed by success is retried.
	t.Run("retries then succeeds", func(t *testing.T) {
		var calls int
		err := deleteWithRetry(context.Background(), nil, sweepNS, "job pod", "job-x", func(context.Context) error {
			calls++
			if calls < 3 {
				return apierrors.NewTimeoutError("transient", 0)
			}
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, 3, calls)
	})

	// A permanent error aborts immediately without retrying.
	t.Run("permanent error fails fast", func(t *testing.T) {
		var calls int
		err := deleteWithRetry(context.Background(), nil, sweepNS, "job pod", "job-x", func(context.Context) error {
			calls++
			return apierrors.NewForbidden(corev1.Resource("pods"), "job-x", errors.New("rbac"))
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete job pod")
		assert.Equal(t, 1, calls)
	})

	// Exhausting all attempts on a transient error surfaces a qualified error.
	t.Run("exhausts attempts on persistent transient error", func(t *testing.T) {
		var calls int
		err := deleteWithRetry(context.Background(), nil, sweepNS, "job pod", "job-x", func(context.Context) error {
			calls++
			return apierrors.NewTimeoutError("always", 0)
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed after 5 attempts")
		assert.Equal(t, kubernetesDeleteRetryAttempts, calls)
	})
}

// netTimeoutErr implements net.Error for the isTransientDeleteError transport
// timeout branch.
type netTimeoutErr struct{}

func (e *netTimeoutErr) Error() string   { return "read tcp: i/o timeout" }
func (e *netTimeoutErr) Timeout() bool   { return true }
func (e *netTimeoutErr) Temporary() bool { return true }

func TestRemoveRetriesPodDelete(t *testing.T) {
	cs := fake.NewSimpleClientset(
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "job-x", Namespace: sweepNS}},
		perJobSecret("job-x"),
	)

	// First Delete of the pod fails transiently, subsequent calls succeed.
	var podDeletes int
	cs.PrependReactor("delete", "pods", func(action k8stesting.Action) (bool, runtime.Object, error) {
		podDeletes++
		if podDeletes == 1 {
			return true, nil, apierrors.NewTimeoutError("transient", 0)
		}
		return false, nil, nil
	})

	p := &KubernetesPod{
		namespace:     sweepNS,
		builder:       &podBuilder{name: "job-x"},
		clientset:     cs,
		perJobSecrets: []string{perJobSecretName("job-x")},
	}
	require.NoError(t, p.remove(context.Background()))
	assert.False(t, podExists(t, cs, "job-x"))
	assert.False(t, secretExists(t, cs, "job-x"))
}

// io timeout wrapper for the transient-error transport branch.
var _ net.Error = (*netTimeoutErr)(nil)
