// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
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
