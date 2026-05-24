// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"context"
	goruntime "runtime"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func testKubernetesConfig() *KubernetesConfig {
	return &KubernetesConfig{
		Enabled:   true,
		Namespace: "gitea-e2e",
		Resources: KubernetesRoleResources{
			Job: KubernetesResources{
				Requests: KubernetesResourceList{CPU: "250m", Memory: "256Mi"},
				Limits:   KubernetesResourceList{CPU: "2", Memory: "2Gi"},
			},
			Service: KubernetesResources{
				Requests: KubernetesResourceList{CPU: "100m", Memory: "128Mi"},
				Limits:   KubernetesResourceList{CPU: "1", Memory: "1Gi"},
			},
		},
		PodReadyTimeout: time.Second,
	}
}

// readyOnCreateReactor mutates every created Pod to Running+Ready so Start's
// readiness wait completes immediately under the fake clientset.
func readyOnCreateReactor(action k8stesting.Action) (bool, runtime.Object, error) {
	pod := action.(k8stesting.CreateAction).GetObject().(*corev1.Pod)
	pod.Status.Phase = corev1.PodRunning
	pod.Status.Conditions = []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}
	return false, pod, nil
}

func TestKubernetesPodBuilder(t *testing.T) {
	cfg := testKubernetesConfig()
	cfg.PodLabels = map[string]string{"team": "ci"}
	cfg.ImagePullSecrets = []string{"regcred"}
	cfg.ServiceAccount = "gitea-runner"

	owner := ownerReferenceFromEnv("runner-pod", "uid-123")
	require.NotNil(t, owner)

	b := newPodBuilder("gitea-job-1", cfg.Namespace, cfg, 3*time.Hour, owner)
	jobSpec := kubernetesContainerSpec(kubernetesJobContainerName,
		&NewContainerInput{Image: "alpine:3", WorkingDir: "/workspace/owner/repo", Env: []string{"FOO=bar", "EQ=a=b"}},
		cfg.Resources.Job, []string{"SYS_PTRACE"}, []string{"MKNOD"}, false)
	b.setJobContainer(jobSpec)

	pod := b.build()

	// identity + labels + annotations
	assert.Equal(t, "gitea-job-1", pod.Name)
	assert.Equal(t, "gitea-e2e", pod.Namespace)
	assert.Equal(t, kubernetesBackendValue, pod.Labels[kubernetesBackendLabel])
	assert.Equal(t, "ci", pod.Labels["team"])
	assert.Equal(t, "gitea-runner", pod.Spec.ServiceAccountName)
	require.Len(t, pod.Spec.ImagePullSecrets, 1)
	assert.Equal(t, "regcred", pod.Spec.ImagePullSecrets[0].Name)

	// never-restart pod with a single job container
	assert.Equal(t, corev1.RestartPolicyNever, pod.Spec.RestartPolicy)
	require.Len(t, pod.Spec.Containers, 1)
	job := pod.Spec.Containers[0]
	assert.Equal(t, kubernetesJobContainerName, job.Name)
	assert.Equal(t, "alpine:3", job.Image)
	assert.Equal(t, "/workspace/owner/repo", job.WorkingDir)

	// env preserves order and values containing '='
	require.Len(t, job.Env, 2)
	assert.Equal(t, corev1.EnvVar{Name: "FOO", Value: "bar"}, job.Env[0])
	assert.Equal(t, corev1.EnvVar{Name: "EQ", Value: "a=b"}, job.Env[1])

	// per-role resources
	assert.True(t, job.Resources.Requests.Cpu().Equal(resource.MustParse("250m")))
	assert.True(t, job.Resources.Requests.Memory().Equal(resource.MustParse("256Mi")))
	assert.True(t, job.Resources.Limits.Cpu().Equal(resource.MustParse("2")))
	assert.True(t, job.Resources.Limits.Memory().Equal(resource.MustParse("2Gi")))

	// workspace volume + mount
	require.Len(t, pod.Spec.Volumes, 1)
	assert.Equal(t, kubernetesWorkspaceVolumeName, pod.Spec.Volumes[0].Name)
	require.NotNil(t, pod.Spec.Volumes[0].EmptyDir)
	require.Len(t, job.VolumeMounts, 1)
	assert.Equal(t, "/workspace/owner/repo", job.VolumeMounts[0].MountPath)

	// security context from cap add/drop
	require.NotNil(t, job.SecurityContext)
	require.NotNil(t, job.SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"SYS_PTRACE"}, job.SecurityContext.Capabilities.Add)
	assert.Equal(t, []corev1.Capability{"MKNOD"}, job.SecurityContext.Capabilities.Drop)

	// orphan-reaping primitives (C2/RC2)
	require.NotNil(t, pod.Spec.ActiveDeadlineSeconds)
	assert.Equal(t, int64(3*60*60), *pod.Spec.ActiveDeadlineSeconds)
	require.Len(t, pod.OwnerReferences, 1)
	assert.Equal(t, "Pod", pod.OwnerReferences[0].Kind)
	assert.Equal(t, "runner-pod", pod.OwnerReferences[0].Name)
	assert.Equal(t, "uid-123", string(pod.OwnerReferences[0].UID))
}

func TestKubernetesPodBuilderPVCWorkspace(t *testing.T) {
	cfg := testKubernetesConfig()
	cfg.WorkspaceVolume = KubernetesWorkspaceVolume{Type: "pvc", Size: "10Gi"}
	b := newPodBuilder("gitea-job-2", cfg.Namespace, cfg, 0, nil)
	b.setJobContainer(kubernetesContainerSpec(kubernetesJobContainerName, &NewContainerInput{Image: "alpine:3"}, cfg.Resources.Job, nil, nil, false))
	pod := b.build()

	require.Len(t, pod.Spec.Volumes, 1)
	require.NotNil(t, pod.Spec.Volumes[0].PersistentVolumeClaim)
	assert.Equal(t, "gitea-job-2", pod.Spec.Volumes[0].PersistentVolumeClaim.ClaimName)
	// no maxLifetime -> no deadline; no env owner -> no ownerReference
	assert.Nil(t, pod.Spec.ActiveDeadlineSeconds)
	assert.Empty(t, pod.OwnerReferences)
}

func TestKubernetesPodBuilderEphemeralWorkspace(t *testing.T) {
	cfg := testKubernetesConfig()
	cfg.WorkspaceVolume = KubernetesWorkspaceVolume{
		Type:             "ephemeral",
		Size:             "20Gi",
		StorageClassName: "fast-ssd",
		AccessModes:      []string{"ReadWriteOnce"},
	}
	b := newPodBuilder("gitea-job-eph", cfg.Namespace, cfg, 0, nil)
	b.setJobContainer(kubernetesContainerSpec(kubernetesJobContainerName, &NewContainerInput{Image: "alpine:3"}, cfg.Resources.Job, nil, nil, false))
	pod := b.build()

	require.Len(t, pod.Spec.Volumes, 1)
	eph := pod.Spec.Volumes[0].Ephemeral
	require.NotNil(t, eph)
	require.NotNil(t, eph.VolumeClaimTemplate)
	spec := eph.VolumeClaimTemplate.Spec
	assert.Equal(t, []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}, spec.AccessModes)
	require.NotNil(t, spec.StorageClassName)
	assert.Equal(t, "fast-ssd", *spec.StorageClassName)
	assert.True(t, spec.Resources.Requests.Storage().Equal(resource.MustParse("20Gi")))
	// emptyDir and static pvc must not be set for the ephemeral type.
	assert.Nil(t, pod.Spec.Volumes[0].EmptyDir)
	assert.Nil(t, pod.Spec.Volumes[0].PersistentVolumeClaim)
}

func TestKubernetesPodBuilderEphemeralDefaultsAccessMode(t *testing.T) {
	cfg := testKubernetesConfig()
	cfg.WorkspaceVolume = KubernetesWorkspaceVolume{Type: "ephemeral", Size: "5Gi"}
	b := newPodBuilder("gitea-job-eph2", cfg.Namespace, cfg, 0, nil)
	b.setJobContainer(kubernetesContainerSpec(kubernetesJobContainerName, &NewContainerInput{Image: "alpine:3"}, cfg.Resources.Job, nil, nil, false))
	pod := b.build()

	spec := pod.Spec.Volumes[0].Ephemeral.VolumeClaimTemplate.Spec
	assert.Equal(t, []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}, spec.AccessModes)
	assert.Nil(t, spec.StorageClassName)
}

func TestKubernetesSidecarRestartPolicyAlways(t *testing.T) {
	cfg := testKubernetesConfig()
	b := newPodBuilder("gitea-job-3", cfg.Namespace, cfg, 0, nil)
	b.addSidecar(kubernetesContainerSpec("postgres", &NewContainerInput{Image: "postgres:16", Name: "postgres"}, cfg.Resources.Service, nil, nil, false))
	pod := b.build()

	require.Len(t, pod.Spec.InitContainers, 1)
	sidecar := pod.Spec.InitContainers[0]
	assert.Equal(t, "postgres", sidecar.Name)
	require.NotNil(t, sidecar.RestartPolicy)
	assert.Equal(t, corev1.ContainerRestartPolicyAlways, *sidecar.RestartPolicy)
}

func TestKubernetesServiceSidecarPortsAndProbes(t *testing.T) {
	cfg := testKubernetesConfig()
	input := &NewContainerInput{
		Image:          "postgres:16",
		Name:           "gitea-job-postgres",
		ExposedPorts:   nat.PortSet{"5432/tcp": struct{}{}},
		NetworkAliases: []string{"postgres"},
	}
	spec := serviceSidecarSpec(input, cfg.Resources.Service, nil, nil, false)

	require.Len(t, spec.Ports, 1)
	assert.Equal(t, int32(5432), spec.Ports[0].ContainerPort)

	require.NotNil(t, spec.StartupProbe)
	require.NotNil(t, spec.StartupProbe.TCPSocket)
	assert.Equal(t, 5432, spec.StartupProbe.TCPSocket.Port.IntValue())
	require.NotNil(t, spec.ReadinessProbe)
	require.NotNil(t, spec.ReadinessProbe.TCPSocket)
	assert.Equal(t, 5432, spec.ReadinessProbe.TCPSocket.Port.IntValue())
}

func TestKubernetesServiceSidecarNoPortsNoProbe(t *testing.T) {
	cfg := testKubernetesConfig()
	spec := serviceSidecarSpec(&NewContainerInput{Image: "redis:7", Name: "redis"}, cfg.Resources.Service, nil, nil, false)
	assert.Empty(t, spec.Ports)
	assert.Nil(t, spec.StartupProbe)
	assert.Nil(t, spec.ReadinessProbe)
}

func TestKubernetesHostAliases(t *testing.T) {
	cfg := testKubernetesConfig()
	b := newPodBuilder("gitea-job-aliases", cfg.Namespace, cfg, 0, nil)
	b.addHostAlias("127.0.0.1", "postgres")
	b.addHostAlias("127.0.0.1", "redis")
	pod := b.build()

	require.Len(t, pod.Spec.HostAliases, 2)
	assert.Equal(t, "127.0.0.1", pod.Spec.HostAliases[0].IP)
	assert.Equal(t, []string{"postgres"}, pod.Spec.HostAliases[0].Hostnames)
	assert.Equal(t, []string{"redis"}, pod.Spec.HostAliases[1].Hostnames)
}

func TestKubernetesPullPolicy(t *testing.T) {
	assert.Equal(t, corev1.PullIfNotPresent, imagePullPolicy(false))
	assert.Equal(t, corev1.PullAlways, imagePullPolicy(true))
}

func TestOwnerReferenceFromEnv(t *testing.T) {
	assert.Nil(t, ownerReferenceFromEnv("", "uid"))
	assert.Nil(t, ownerReferenceFromEnv("name", ""))
	owner := ownerReferenceFromEnv("runner-pod", "uid-9")
	require.NotNil(t, owner)
	assert.Equal(t, "v1", owner.APIVersion)
	assert.Equal(t, "Pod", owner.Kind)
}

// TestRunnerArchLocal asserts the docker-free resolver maps the local
// architecture to GitHub runner.arch names without touching a docker daemon.
func TestRunnerArchLocal(t *testing.T) {
	want := map[string]string{
		"amd64":   "X64",
		"x86_64":  "X64",
		"386":     "X86",
		"arm64":   "ARM64",
		"aarch64": "ARM64",
	}
	expected, ok := want[goruntime.GOARCH]
	if !ok {
		expected = goruntime.GOARCH
	}
	assert.Equal(t, expected, RunnerArchLocal())
}

// TestKubernetesGetRunnerContext asserts the pod backend resolves runner.arch
// locally (RunnerArchLocal) — never via the docker-pinging RunnerArch — so the
// expression evaluator never constructs a docker client. The environment is
// built without a pod/clientset precisely to prove GetRunnerContext needs no
// docker (or kubernetes) connection.
func TestKubernetesGetRunnerContext(t *testing.T) {
	env := &KubernetesExecutionEnvironment{}
	ctx := context.Background()

	runnerCtx := env.GetRunnerContext(ctx)

	assert.Equal(t, "Linux", runnerCtx["os"])
	assert.Equal(t, RunnerArchLocal(), runnerCtx["arch"])
	assert.NotEmpty(t, runnerCtx["arch"])
	assert.Equal(t, "/tmp", runnerCtx["temp"])
	assert.Equal(t, "/opt/hostedtoolcache", runnerCtx["tool_cache"])
}

func newTestPod(t *testing.T, clientset *fake.Clientset, podName string) *KubernetesPod {
	t.Helper()
	cfg := testKubernetesConfig()
	return newKubernetesPodWithClientset(clientset, nil, KubernetesPodConfig{
		Config:      cfg,
		PodName:     podName,
		MaxLifetime: time.Hour,
	})
}

func newTestPodWithCleanup(t *testing.T, clientset *fake.Clientset, podName, policy string) *KubernetesPod {
	t.Helper()
	cfg := testKubernetesConfig()
	cfg.PodCleanup = policy
	return newKubernetesPodWithClientset(clientset, nil, KubernetesPodConfig{
		Config:      cfg,
		PodName:     podName,
		MaxLifetime: time.Hour,
	})
}

func TestShouldRetainPod(t *testing.T) {
	cases := []struct {
		policy       string
		jobSucceeded bool
		want         bool
	}{
		{KubernetesPodCleanupAlways, true, false},
		{KubernetesPodCleanupAlways, false, false},
		{KubernetesPodCleanupOnFailure, true, false},
		{KubernetesPodCleanupOnFailure, false, true},
		{KubernetesPodCleanupNever, true, true},
		{KubernetesPodCleanupNever, false, true},
		{"", true, false}, // unset behaves as "always"
	}
	for _, c := range cases {
		assert.Equal(t, c.want, shouldRetainPod(c.policy, c.jobSucceeded),
			"policy=%q jobSucceeded=%v", c.policy, c.jobSucceeded)
	}
}

// submitJobPod creates+starts a ready job pod under a cleanup policy, returning
// the job environment so a test can drive RemoveWithRetention.
func submitJobPod(t *testing.T, clientset *fake.Clientset, podName, policy string) *KubernetesExecutionEnvironment {
	t.Helper()
	clientset.PrependReactor("create", "pods", readyOnCreateReactor)
	pod := newTestPodWithCleanup(t, clientset, podName, policy)
	job := pod.NewJobContainer(&NewContainerInput{Image: "alpine:3", WorkingDir: "/workspace"})
	ctx := context.Background()
	require.NoError(t, job.Create(nil, nil)(ctx))
	require.NoError(t, job.Start(false)(ctx))
	return job
}

func TestKubernetesPodCleanupAlwaysDeletes(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	job := submitJobPod(t, clientset, "gitea-job-always", KubernetesPodCleanupAlways)
	ctx := context.Background()

	// Even a successful job is deleted under "always".
	require.NoError(t, job.RemoveWithRetention(true)(ctx))
	_, err := clientset.CoreV1().Pods("gitea-e2e").Get(ctx, "gitea-job-always", metav1.GetOptions{})
	assert.True(t, apierrors.IsNotFound(err), "always must delete the pod")
}

func TestKubernetesPodCleanupOnFailureDeletesOnSuccess(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	job := submitJobPod(t, clientset, "gitea-job-onfail-ok", KubernetesPodCleanupOnFailure)
	ctx := context.Background()

	require.NoError(t, job.RemoveWithRetention(true)(ctx))
	_, err := clientset.CoreV1().Pods("gitea-e2e").Get(ctx, "gitea-job-onfail-ok", metav1.GetOptions{})
	assert.True(t, apierrors.IsNotFound(err), "on-failure must delete a SUCCEEDED job's pod")
}

func TestKubernetesPodCleanupOnFailureKeepsOnFailure(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	job := submitJobPod(t, clientset, "gitea-job-onfail-bad", KubernetesPodCleanupOnFailure)
	ctx := context.Background()

	require.NoError(t, job.RemoveWithRetention(false)(ctx))
	got, err := clientset.CoreV1().Pods("gitea-e2e").Get(ctx, "gitea-job-onfail-bad", metav1.GetOptions{})
	require.NoError(t, err, "on-failure must KEEP a FAILED job's pod")
	assert.Equal(t, kubernetesRetainValue, got.Labels[kubernetesRetainLabel], "kept pod must carry the retain marker")
}

func TestKubernetesPodCleanupNeverKeeps(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	job := submitJobPod(t, clientset, "gitea-job-never", KubernetesPodCleanupNever)
	ctx := context.Background()

	// Even a successful job is kept under "never".
	require.NoError(t, job.RemoveWithRetention(true)(ctx))
	got, err := clientset.CoreV1().Pods("gitea-e2e").Get(ctx, "gitea-job-never", metav1.GetOptions{})
	require.NoError(t, err, "never must KEEP the pod")
	assert.Equal(t, kubernetesRetainValue, got.Labels[kubernetesRetainLabel])
	// The backend class label is preserved alongside the retain marker.
	assert.Equal(t, kubernetesBackendValue, got.Labels[kubernetesBackendLabel])
}

func TestKubernetesCreateStartSubmitsReadyPod(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("create", "pods", readyOnCreateReactor)

	pod := newTestPod(t, clientset, "gitea-job-start")
	job := pod.NewJobContainer(&NewContainerInput{Image: "alpine:3", WorkingDir: "/workspace"})

	ctx := context.Background()
	require.NoError(t, job.Create(nil, nil)(ctx))
	require.NoError(t, job.Start(false)(ctx))

	got, err := clientset.CoreV1().Pods("gitea-e2e").Get(ctx, "gitea-job-start", metav1.GetOptions{})
	require.NoError(t, err)
	assert.Equal(t, "gitea-job-start", got.Name)
	require.Len(t, got.Spec.Containers, 1)

	// Remove deletes the Pod.
	require.NoError(t, job.Remove()(ctx))
	_, err = clientset.CoreV1().Pods("gitea-e2e").Get(ctx, "gitea-job-start", metav1.GetOptions{})
	assert.True(t, apierrors.IsNotFound(err))
}

func TestKubernetesServiceStartIsNoop(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	pod := newTestPod(t, clientset, "gitea-job-svc")
	svc := pod.NewServiceContainer(&NewContainerInput{Image: "postgres:16", Name: "postgres"})

	ctx := context.Background()
	require.NoError(t, svc.Create(nil, nil)(ctx))
	// A service Start must not submit a Pod.
	require.NoError(t, svc.Start(false)(ctx))
	list, err := clientset.CoreV1().Pods("gitea-e2e").List(ctx, metav1.ListOptions{})
	require.NoError(t, err)
	assert.Empty(t, list.Items)
}

// TestKubernetesFailedStartCleanup forces a mid-pipeline Start failure (Ready
// timeout) and asserts the explicit Remove path deletes the Pod AND its per-job
// Secret (R2). It does NOT assert real ownerReference GC or activeDeadlineSeconds
// — the fake clientset has no GC controller or kubelet.
func TestKubernetesFailedStartCleanup(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	// No ready reactor: the Pod never becomes Ready, so Start times out.

	pod := newTestPod(t, clientset, "gitea-job-fail")
	secretName := perJobSecretName("gitea-job-fail")
	_, err := clientset.CoreV1().Secrets("gitea-e2e").Create(context.Background(),
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: "gitea-e2e"}}, metav1.CreateOptions{})
	require.NoError(t, err)
	pod.registerPerJobSecret(secretName)

	job := pod.NewJobContainer(&NewContainerInput{Image: "alpine:3", WorkingDir: "/workspace"})
	ctx := context.Background()
	require.NoError(t, job.Create(nil, nil)(ctx))

	// Start submits the Pod then blocks on readiness; with the 1s timeout it fails.
	startErr := job.Start(false)(ctx)
	require.Error(t, startErr)
	assert.Contains(t, startErr.Error(), "not ready")

	// The created Pod exists at this point (submit succeeded, readiness failed).
	_, err = clientset.CoreV1().Pods("gitea-e2e").Get(ctx, "gitea-job-fail", metav1.GetOptions{})
	require.NoError(t, err)

	// cleanupFailedStart routes through Remove: Pod and per-job Secret deleted.
	require.NoError(t, job.Remove()(ctx))
	_, err = clientset.CoreV1().Pods("gitea-e2e").Get(ctx, "gitea-job-fail", metav1.GetOptions{})
	assert.True(t, apierrors.IsNotFound(err), "pod should be deleted")
	_, err = clientset.CoreV1().Secrets("gitea-e2e").Get(ctx, secretName, metav1.GetOptions{})
	assert.True(t, apierrors.IsNotFound(err), "per-job secret should be deleted")
}

func TestKubernetesPodBuilderRuntimeClass(t *testing.T) {
	// Unset by default: spec.runtimeClassName stays nil (cluster default runtime).
	cfg := testKubernetesConfig()
	assert.Nil(t, newPodBuilder("gitea-job-rc0", cfg.Namespace, cfg, 0, nil).build().Spec.RuntimeClassName)

	// A configured RuntimeClass is applied to every job pod the runner builds.
	cfg.RuntimeClassName = "gvisor"
	rc := newPodBuilder("gitea-job-rc1", cfg.Namespace, cfg, 0, nil).build().Spec.RuntimeClassName
	require.NotNil(t, rc)
	assert.Equal(t, "gvisor", *rc)
}
