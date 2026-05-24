// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package config

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.yaml.in/yaml/v4"
)

// DefaultPostTaskScriptTimeout is the fallback cap on how long the post-task
// script may run when post_task_script is set without an explicit timeout. It is
// applied both at config load (for a configured script) and at the point of use
// (so a programmatically built config still gets a sane bound).
const DefaultPostTaskScriptTimeout = 5 * time.Minute

// Log represents the configuration for logging.
type Log struct {
	Level string `yaml:"level"` // Level indicates the logging level.
}

// Runner represents the configuration for the runner.
type Runner struct {
	File                  string            `yaml:"file"`                     // File specifies the file path for the runner.
	Capacity              int               `yaml:"capacity"`                 // Capacity specifies the capacity of the runner.
	Envs                  map[string]string `yaml:"envs"`                     // Envs stores environment variables for the runner.
	EnvFile               string            `yaml:"env_file"`                 // EnvFile specifies the path to the file containing environment variables for the runner.
	Timeout               time.Duration     `yaml:"timeout"`                  // Timeout specifies the duration for runner timeout.
	ShutdownTimeout       time.Duration     `yaml:"shutdown_timeout"`         // ShutdownTimeout specifies the duration to wait for running jobs to complete during a shutdown of the runner.
	Insecure              bool              `yaml:"insecure"`                 // Insecure indicates whether the runner operates in an insecure mode.
	FetchTimeout          time.Duration     `yaml:"fetch_timeout"`            // FetchTimeout specifies the timeout duration for fetching resources.
	FetchInterval         time.Duration     `yaml:"fetch_interval"`           // FetchInterval specifies the interval duration for fetching resources.
	FetchIntervalMax      time.Duration     `yaml:"fetch_interval_max"`       // FetchIntervalMax specifies the maximum backoff interval when idle.
	WorkdirCleanupAge     time.Duration     `yaml:"workdir_cleanup_age"`      // WorkdirCleanupAge removes stale bind-workdir task directories and orphaned host-mode scratch dirs older than this duration during idle cleanup.
	IdleCleanupInterval   time.Duration     `yaml:"idle_cleanup_interval"`    // IdleCleanupInterval runs stale-directory cleanup periodically while the runner is idle. Set to 0 to disable cleanup cadence.
	LogReportInterval     time.Duration     `yaml:"log_report_interval"`      // LogReportInterval specifies the base interval for periodic log flush.
	LogReportMaxLatency   time.Duration     `yaml:"log_report_max_latency"`   // LogReportMaxLatency specifies the max time a log row can wait before being sent.
	LogReportBatchSize    int               `yaml:"log_report_batch_size"`    // LogReportBatchSize triggers immediate log flush when buffer reaches this size.
	StateReportInterval   time.Duration     `yaml:"state_report_interval"`    // StateReportInterval specifies the interval for state reporting.
	ReportCloseTimeout    time.Duration     `yaml:"report_close_timeout"`     // ReportCloseTimeout caps each RPC attempt when flushing the final logs and task state at job completion, on a detached context so a server cancel can't block the acknowledgement.
	Labels                []string          `yaml:"labels"`                   // Labels specify the labels of the runner. Labels are declared on each startup
	GithubMirror          string            `yaml:"github_mirror"`            // GithubMirror defines what mirrors should be used when using github
	AllocatePTY           bool              `yaml:"allocate_pty"`             // AllocatePTY allocates a pseudo-TTY for each step's process. Default is false, matching GitHub's actions/runner. Enable only for jobs that need an interactive terminal; tools like docker build emit redrawing progress frames into the captured log when a TTY is present. Applies to both host and docker backends.
	PostTaskScript        string            `yaml:"post_task_script"`         // PostTaskScript is the path to an executable script run on the host after each task's cleanup completes. Empty disables the hook. On Windows use .exe/.bat/.cmd; PowerShell (.ps1) is not supported yet as the configured path.
	PostTaskScriptTimeout time.Duration     `yaml:"post_task_script_timeout"` // PostTaskScriptTimeout caps how long the post-task script may run. Default is 5m when post_task_script is set.
}

// Cache represents the configuration for caching.
type Cache struct {
	Enabled        *bool  `yaml:"enabled"`         // Enabled indicates whether caching is enabled. It is a pointer to distinguish between false and not set. If not set, it will be true.
	Dir            string `yaml:"dir"`             // Dir specifies the directory path for caching.
	Host           string `yaml:"host"`            // Host specifies the caching host.
	Port           uint16 `yaml:"port"`            // Port specifies the caching port.
	ExternalServer string `yaml:"external_server"` // ExternalServer specifies the URL of external cache server
	ExternalSecret string `yaml:"external_secret"` // ExternalSecret is a shared secret between this runner and an external gitea-runner cache-server, enabling per-job ACTIONS_RUNTIME_TOKEN authentication and repo scoping over the network. Leave empty to keep the legacy unauthenticated behavior.
	OfflineMode    bool   `yaml:"offline_mode"`    // OfflineMode reuses a cached action without fetching from the remote; a moved tag or branch stays at the cached commit until the cache entry is removed.
}

// Container represents the configuration for the container.
type Container struct {
	Network              string                        `yaml:"network"`                // Network specifies the network for the container.
	NetworkCreateOptions ContainerNetworkCreateOptions `yaml:"network_create_options"` // Add options when the network need to be created by the runner
	NetworkMode          string                        `yaml:"network_mode"`           // Deprecated: use Network instead. Could be removed after Gitea 1.20
	Privileged           bool                          `yaml:"privileged"`             // Privileged indicates whether the container runs in privileged mode.
	Options              string                        `yaml:"options"`                // Options specifies additional options for the container.
	WorkdirParent        string                        `yaml:"workdir_parent"`         // WorkdirParent specifies the parent directory for the container's working directory.
	ValidVolumes         []string                      `yaml:"valid_volumes"`          // ValidVolumes specifies the volumes (including bind mounts) can be mounted to containers.
	DockerHost           string                        `yaml:"docker_host"`            // DockerHost specifies the Docker host. It overrides the value specified in environment variable DOCKER_HOST.
	ForcePull            bool                          `yaml:"force_pull"`             // Pull docker image(s) even if already present
	ForceRebuild         bool                          `yaml:"force_rebuild"`          // Rebuild docker image(s) even if already present
	RequireDocker        bool                          `yaml:"require_docker"`         // Always require a reachable docker daemon, even if not required by runner
	DockerTimeout        time.Duration                 `yaml:"docker_timeout"`         // Timeout to wait for the docker daemon to be reachable, if docker is required by require_docker or runner
	BindWorkdir          bool                          `yaml:"bind_workdir"`           // BindWorkdir binds the workspace to the host filesystem instead of using Docker volumes. Required for DinD when jobs use docker compose with bind mounts.
}

type ContainerNetworkCreateOptions struct {
	EnableIPv4 *bool `yaml:"enable_ipv4"` // Enable or disable IPv4 for the network (true for docker by default)
	EnableIPv6 *bool `yaml:"enable_ipv6"` // Enable or disable IPv6 for the network (false for docker by default)
}

// Host represents the configuration for the host.
type Host struct {
	WorkdirParent string `yaml:"workdir_parent"` // WorkdirParent specifies the parent directory for the host's working directory.
}

// Metrics represents the configuration for the Prometheus metrics endpoint.
type Metrics struct {
	Enabled bool   `yaml:"enabled"` // Enabled indicates whether the metrics endpoint is exposed.
	Addr    string `yaml:"addr"`    // Addr specifies the listen address for the metrics HTTP server (e.g., ":9101").
}

// KubernetesResourceList is the cpu/memory pair for a requests or limits entry.
// Values are Kubernetes resource.Quantity strings (e.g. "250m", "256Mi").
type KubernetesResourceList struct {
	CPU    string `yaml:"cpu"`    // CPU quantity (e.g. "250m", "2").
	Memory string `yaml:"memory"` // Memory quantity (e.g. "256Mi", "2Gi").
}

// KubernetesResources holds the requests and limits for one container role.
type KubernetesResources struct {
	Requests KubernetesResourceList `yaml:"requests"` // Requests is the guaranteed reservation.
	Limits   KubernetesResourceList `yaml:"limits"`   // Limits is the hard cap.
}

// KubernetesContainerResources carries per-container-role resource requirements.
// Each role is tuned independently because their footprints differ markedly.
type KubernetesContainerResources struct {
	Job     KubernetesResources `yaml:"job"`     // Job container running the steps.
	Service KubernetesResources `yaml:"service"` // Service sidecar containers.
	Dind    KubernetesResources `yaml:"dind"`    // Optional docker-in-docker sidecar (needs >=1Gi).
}

// KubernetesWorkspaceVolume selects how the shared job workspace is backed.
//
//   - "emptyDir" (default): node-local scratch, freed with the pod.
//   - "ephemeral": a pod-owned generic ephemeral volume (a PVC auto-created and
//     auto-GC'd with the pod) — the recommended persistent-backed per-job
//     workspace; composes with the own-Pod ownerReference cleanup.
//   - "pvc": a pre-existing/static PersistentVolumeClaim named after the pod
//     (shared/advanced use).
type KubernetesWorkspaceVolume struct {
	Type             string   `yaml:"type"`               // Type is "emptyDir" (default), "ephemeral", or "pvc".
	Size             string   `yaml:"size"`               // Size is the requested storage when Type is "ephemeral" or "pvc" (e.g. "10Gi").
	StorageClassName string   `yaml:"storage_class_name"` // StorageClassName optionally selects a storage class for "ephemeral"/"pvc"; empty uses the cluster default.
	AccessModes      []string `yaml:"access_modes"`       // AccessModes for the backing claim; defaults to ["ReadWriteOnce"].
}

// KubernetesDind configures the optional on-demand docker-in-docker sidecar that
// serves `uses: docker://` container steps.
type KubernetesDind struct {
	Enabled  bool   `yaml:"enabled"`  // Enabled injects the dind sidecar for container steps. Explicit; no auto-detect.
	Rootless *bool  `yaml:"rootless"` // Rootless selects docker:dind-rootless. It is a pointer to distinguish unset (defaults true) from an explicit false (privileged sidecar).
	Image    string `yaml:"image"`    // Image overrides the dind image; empty selects a sensible default.
}

// Kubernetes represents the configuration for the Kubernetes-native pod backend.
// When Enabled, every job runs as one Pod (init + long-lived job container +
// native service sidecars) instead of via the docker/host backends.
type Kubernetes struct {
	Enabled            bool                         `yaml:"enabled"`              // Enabled routes every job to the Kubernetes pod backend.
	Namespace          string                       `yaml:"namespace"`            // Namespace for all runner-created pods/secrets/configmaps.
	Kubeconfig         string                       `yaml:"kubeconfig"`           // Kubeconfig path for out-of-cluster dev; empty selects in-cluster config.
	ServiceAccount     string                       `yaml:"service_account"`      // ServiceAccount assigned to job pods.
	DefaultJobImage    string                       `yaml:"default_job_image"`    // DefaultJobImage backs jobs with no explicit container:image (a bare run: job whose runs-on label maps to no image). Defaults to node:20-bookworm.
	ImagePullSecrets   []string                     `yaml:"image_pull_secrets"`   // ImagePullSecrets attached to job pods for private registries.
	WorkspaceVolume    KubernetesWorkspaceVolume    `yaml:"workspace_volume"`     // WorkspaceVolume backs the shared job workspace.
	PodLabels          map[string]string            `yaml:"pod_labels"`           // PodLabels are merged onto every job pod.
	PodAnnotations     map[string]string            `yaml:"pod_annotations"`      // PodAnnotations are merged onto every job pod.
	Resources          KubernetesContainerResources `yaml:"resources"`            // Resources are the per-container-role requests/limits.
	PodReadyTimeout    time.Duration                `yaml:"pod_ready_timeout"`    // PodReadyTimeout bounds the wait for a job pod to become Ready.
	ConnectBackAddress string                       `yaml:"connect_back_address"` // ConnectBackAddress is the cluster-reachable DNS job pods use to reach the runner's cache/results server.
	Dind               KubernetesDind               `yaml:"dind"`                 // Dind configures the optional docker-in-docker sidecar.
	PodCleanup         string                       `yaml:"pod_cleanup"`          // PodCleanup selects when completed job pods (and their per-job Secret/ConfigMap) are deleted: "always" (default), "on-failure" (delete only on success, keep failed pods), or "never" (keep all for auditing). Retained pods consume cluster resources and a capacity slot until manually deleted.
	MaxConcurrentJobs  int                          `yaml:"max_concurrent_jobs"`  // MaxConcurrentJobs caps how many job pods this runner will run concurrently in Kubernetes mode. 0 (default) means UNLIMITED: the runner accepts every assigned job and creates a pod for it, letting the kube-scheduler and a namespace ResourceQuota be the real backpressure. A positive value is an optional client-side cap. Ignored unless kubernetes.enabled; docker/host modes always use runner.capacity.
}

// Config represents the overall configuration.
type Config struct {
	Log        Log        `yaml:"log"`        // Log represents the configuration for logging.
	Runner     Runner     `yaml:"runner"`     // Runner represents the configuration for the runner.
	Cache      Cache      `yaml:"cache"`      // Cache represents the configuration for caching.
	Container  Container  `yaml:"container"`  // Container represents the configuration for the container.
	Host       Host       `yaml:"host"`       // Host represents the configuration for the host.
	Metrics    Metrics    `yaml:"metrics"`    // Metrics represents the configuration for the Prometheus metrics endpoint.
	Kubernetes Kubernetes `yaml:"kubernetes"` // Kubernetes represents the configuration for the Kubernetes-native pod backend.
}

// LoadDefault returns the default configuration.
// If file is not empty, it will be used to load the configuration.
func LoadDefault(file string) (*Config, error) {
	cfg := &Config{}
	definedRunnerKeys := map[string]bool{}
	if file != "" {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("open config file %q: %w", file, err)
		}
		if err := yaml.Unmarshal(content, cfg); err != nil {
			return nil, fmt.Errorf("parse config file %q: %w", file, err)
		}
		definedRunnerKeys, err = definedRunnerConfigKeys(content)
		if err != nil {
			return nil, fmt.Errorf("parse config file %q for defaults metadata: %w", file, err)
		}
	}

	if cfg.Runner.EnvFile != "" {
		if stat, err := os.Stat(cfg.Runner.EnvFile); err == nil && !stat.IsDir() {
			envs, err := godotenv.Read(cfg.Runner.EnvFile)
			if err != nil {
				return nil, fmt.Errorf("read env file %q: %w", cfg.Runner.EnvFile, err)
			}
			if cfg.Runner.Envs == nil {
				cfg.Runner.Envs = map[string]string{}
			}
			maps.Copy(cfg.Runner.Envs, envs)
		}
	}

	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Runner.File == "" {
		cfg.Runner.File = ".runner"
	}
	if cfg.Runner.Capacity <= 0 {
		cfg.Runner.Capacity = 1
	}
	if cfg.Runner.Timeout <= 0 {
		cfg.Runner.Timeout = 3 * time.Hour
	}
	if cfg.Cache.Enabled == nil {
		b := true
		cfg.Cache.Enabled = &b
	}
	if *cfg.Cache.Enabled {
		if cfg.Cache.Dir == "" {
			home, _ := os.UserHomeDir()
			cfg.Cache.Dir = filepath.Join(home, ".cache", "actcache")
		}
		if cfg.Cache.ExternalServer != "" && cfg.Cache.ExternalSecret == "" {
			return nil, errors.New("cache.external_server is set but cache.external_secret is empty; configure the same external_secret on this runner and the gitea-runner cache-server")
		}
	}
	if cfg.Container.WorkdirParent == "" {
		cfg.Container.WorkdirParent = "workspace"
	}
	if cfg.Host.WorkdirParent == "" {
		home, _ := os.UserHomeDir()
		cfg.Host.WorkdirParent = filepath.Join(home, ".cache", "act")
	}
	if cfg.Runner.FetchTimeout <= 0 {
		cfg.Runner.FetchTimeout = 5 * time.Second
	}
	if cfg.Runner.FetchInterval <= 0 {
		cfg.Runner.FetchInterval = 2 * time.Second
	}
	if cfg.Runner.FetchIntervalMax <= 0 {
		cfg.Runner.FetchIntervalMax = 5 * time.Second
	}
	if cfg.Runner.WorkdirCleanupAge == 0 && !definedRunnerKeys["workdir_cleanup_age"] {
		cfg.Runner.WorkdirCleanupAge = 24 * time.Hour
	}
	if cfg.Runner.IdleCleanupInterval == 0 && !definedRunnerKeys["idle_cleanup_interval"] {
		cfg.Runner.IdleCleanupInterval = 10 * time.Minute
	}
	if cfg.Runner.LogReportInterval <= 0 {
		cfg.Runner.LogReportInterval = 5 * time.Second
	}
	if cfg.Runner.LogReportMaxLatency <= 0 {
		cfg.Runner.LogReportMaxLatency = 3 * time.Second
	}
	if cfg.Runner.LogReportBatchSize <= 0 {
		cfg.Runner.LogReportBatchSize = 100
	}
	if cfg.Runner.StateReportInterval <= 0 {
		cfg.Runner.StateReportInterval = 5 * time.Second
	}
	if cfg.Runner.ReportCloseTimeout <= 0 {
		cfg.Runner.ReportCloseTimeout = 10 * time.Second
	}
	if cfg.Runner.PostTaskScript != "" && cfg.Runner.PostTaskScriptTimeout <= 0 {
		cfg.Runner.PostTaskScriptTimeout = DefaultPostTaskScriptTimeout
	}
	if cfg.Metrics.Addr == "" {
		cfg.Metrics.Addr = "127.0.0.1:9101"
	}

	// Validate and fix invalid config combinations to prevent confusing behavior.
	if cfg.Runner.FetchIntervalMax < cfg.Runner.FetchInterval {
		log.Warnf("fetch_interval_max (%v) is less than fetch_interval (%v), setting fetch_interval_max to fetch_interval",
			cfg.Runner.FetchIntervalMax, cfg.Runner.FetchInterval)
		cfg.Runner.FetchIntervalMax = cfg.Runner.FetchInterval
	}
	if cfg.Runner.LogReportMaxLatency >= cfg.Runner.LogReportInterval {
		log.Warnf("log_report_max_latency (%v) >= log_report_interval (%v), the max-latency timer will never fire before the periodic ticker; consider lowering log_report_max_latency",
			cfg.Runner.LogReportMaxLatency, cfg.Runner.LogReportInterval)
	}

	applyKubernetesDefaults(&cfg.Kubernetes)
	if err := validateKubernetes(&cfg.Kubernetes); err != nil {
		return nil, err
	}
	if err := validateKubernetesCacheReachability(cfg); err != nil {
		return nil, err
	}

	// although `container.network_mode` will be deprecated, but we have to be compatible with it for now.
	if cfg.Container.NetworkMode != "" && cfg.Container.Network == "" {
		log.Warn("You are trying to use deprecated configuration item of `container.network_mode`, please use `container.network` instead.")
		if cfg.Container.NetworkMode == "bridge" {
			// Previously, if the value of `container.network_mode` is `bridge`, we will create a new network for job.
			// But “bridge” is easily confused with the bridge network created by Docker by default.
			// So we set the value of `container.network` to empty string to make `runner` automatically create a new network for job.
			cfg.Container.Network = ""
		} else {
			cfg.Container.Network = cfg.Container.NetworkMode
		}
	}

	return cfg, nil
}

// Kubernetes backend defaults (M3). Each container role is tuned independently;
// the dind sidecar needs >=1Gi to build images reliably.
const (
	defaultKubernetesNamespace           = "gitea-e2e"
	defaultKubernetesPodReadyTimeout     = 5 * time.Minute
	defaultKubernetesWorkspaceType       = "emptyDir"
	defaultKubernetesWorkspaceAccessMode = "ReadWriteOnce"
	defaultKubernetesDindRootless        = true
	defaultKubernetesPodCleanup          = "always"
	// defaultKubernetesJobImage backs bare run: jobs (no container:image). It
	// carries node (for JS actions) plus common build tools.
	defaultKubernetesJobImage = "node:20-bookworm"

	defaultKubernetesJobCPURequest    = "250m"
	defaultKubernetesJobMemoryRequest = "256Mi"
	defaultKubernetesJobCPULimit      = "2"
	defaultKubernetesJobMemoryLimit   = "2Gi"

	defaultKubernetesServiceCPURequest    = "100m"
	defaultKubernetesServiceMemoryRequest = "128Mi"
	defaultKubernetesServiceCPULimit      = "1"
	defaultKubernetesServiceMemoryLimit   = "1Gi"

	defaultKubernetesDindCPURequest    = "500m"
	defaultKubernetesDindMemoryRequest = "1Gi"
	defaultKubernetesDindCPULimit      = "2"
	defaultKubernetesDindMemoryLimit   = "4Gi"
)

// applyKubernetesDefaults fills empty Kubernetes settings with their defaults.
// It is safe to call regardless of Enabled: the values are inert for the
// docker/host backends and only consulted when the pod backend is selected.
func applyKubernetesDefaults(k *Kubernetes) {
	if k.Namespace == "" {
		k.Namespace = defaultKubernetesNamespace
	}
	if k.DefaultJobImage == "" {
		k.DefaultJobImage = defaultKubernetesJobImage
	}
	if k.PodReadyTimeout <= 0 {
		k.PodReadyTimeout = defaultKubernetesPodReadyTimeout
	}
	if k.WorkspaceVolume.Type == "" {
		k.WorkspaceVolume.Type = defaultKubernetesWorkspaceType
	}
	// Claim-backed workspaces default to ReadWriteOnce when access modes are
	// unset; emptyDir ignores the field.
	if (k.WorkspaceVolume.Type == "ephemeral" || k.WorkspaceVolume.Type == "pvc") && len(k.WorkspaceVolume.AccessModes) == 0 {
		k.WorkspaceVolume.AccessModes = []string{defaultKubernetesWorkspaceAccessMode}
	}
	// Rootless defaults to true; an explicit `rootless: false` opts into a
	// privileged dind sidecar. The pointer distinguishes the two.
	if k.Dind.Rootless == nil {
		rootless := defaultKubernetesDindRootless
		k.Dind.Rootless = &rootless
	}
	if k.PodCleanup == "" {
		k.PodCleanup = defaultKubernetesPodCleanup
	}
	defaultRoleResources(&k.Resources)
}

// defaultRoleResources fills any empty per-role requests/limits with the M3
// defaults, preserving operator overrides.
func defaultRoleResources(r *KubernetesContainerResources) {
	defaultResource(&r.Job,
		defaultKubernetesJobCPURequest, defaultKubernetesJobMemoryRequest,
		defaultKubernetesJobCPULimit, defaultKubernetesJobMemoryLimit)
	defaultResource(&r.Service,
		defaultKubernetesServiceCPURequest, defaultKubernetesServiceMemoryRequest,
		defaultKubernetesServiceCPULimit, defaultKubernetesServiceMemoryLimit)
	defaultResource(&r.Dind,
		defaultKubernetesDindCPURequest, defaultKubernetesDindMemoryRequest,
		defaultKubernetesDindCPULimit, defaultKubernetesDindMemoryLimit)
}

func defaultResource(res *KubernetesResources, cpuReq, memReq, cpuLim, memLim string) {
	if res.Requests.CPU == "" {
		res.Requests.CPU = cpuReq
	}
	if res.Requests.Memory == "" {
		res.Requests.Memory = memReq
	}
	if res.Limits.CPU == "" {
		res.Limits.CPU = cpuLim
	}
	if res.Limits.Memory == "" {
		res.Limits.Memory = memLim
	}
}

// validateKubernetes rejects incoherent Kubernetes settings. Validation only
// fires when the backend is enabled so a stray disabled block cannot break a
// docker/host runner.
func validateKubernetes(k *Kubernetes) error {
	if !k.Enabled {
		return nil
	}
	switch k.WorkspaceVolume.Type {
	case "emptyDir":
	case "ephemeral", "pvc":
		if k.WorkspaceVolume.Size == "" {
			return fmt.Errorf("kubernetes.workspace_volume.size is required when workspace_volume.type is %q", k.WorkspaceVolume.Type)
		}
	default:
		return fmt.Errorf("kubernetes.workspace_volume.type %q is invalid; must be \"emptyDir\", \"ephemeral\", or \"pvc\"", k.WorkspaceVolume.Type)
	}
	switch k.PodCleanup {
	case "", "always", "on-failure", "never":
	default:
		return fmt.Errorf("kubernetes.pod_cleanup %q is invalid; must be \"always\", \"on-failure\", or \"never\"", k.PodCleanup)
	}
	if k.MaxConcurrentJobs < 0 {
		return fmt.Errorf("kubernetes.max_concurrent_jobs must be >= 0 (0 = unlimited), got %d", k.MaxConcurrentJobs)
	}
	return nil
}

// validateKubernetesCacheReachability fails fast when the Kubernetes backend is
// enabled but the embedded cache server would advertise a URL job pods cannot
// reach (pre-mortem #2). The embedded cache binds to cache.host (loopback/empty
// in-cluster); job pods can only reach it via a cluster DNS name, so either
// kubernetes.connect_back_address (the cache Service DNS) or an external cache
// (cache.external_server) must be set. The existing external_server/external_secret
// rule is enforced separately and is unaffected (D6).
func validateKubernetesCacheReachability(cfg *Config) error {
	if !cfg.Kubernetes.Enabled {
		return nil
	}
	cacheEnabled := cfg.Cache.Enabled == nil || *cfg.Cache.Enabled
	if !cacheEnabled || cfg.Cache.ExternalServer != "" {
		// Disabled cache or an explicit external server: nothing loopback to guard.
		return nil
	}
	if cfg.Kubernetes.ConnectBackAddress == "" && isLoopbackOrEmptyHost(cfg.Cache.Host) {
		return errors.New("kubernetes.enabled with the embedded cache, but cache.host is loopback/empty and kubernetes.connect_back_address is unset; job pods cannot reach 127.0.0.1 — set kubernetes.connect_back_address to a cluster-reachable cache Service DNS (or set cache.external_server)")
	}
	return nil
}

// isLoopbackOrEmptyHost reports whether host is empty or a loopback/wildcard
// address that a job pod in another network namespace cannot route to.
func isLoopbackOrEmptyHost(host string) bool {
	switch host {
	case "", "127.0.0.1", "localhost", "::1", "0.0.0.0":
		return true
	default:
		return false
	}
}

func definedRunnerConfigKeys(content []byte) (map[string]bool, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(content, &root); err != nil {
		return nil, err
	}

	defined := map[string]bool{}
	if len(root.Content) == 0 {
		return defined, nil
	}

	doc := root.Content[0]
	for i := 0; i+1 < len(doc.Content); i += 2 {
		key := doc.Content[i]
		value := doc.Content[i+1]
		if key.Value != "runner" || value.Kind != yaml.MappingNode {
			continue
		}
		for j := 0; j+1 < len(value.Content); j += 2 {
			defined[value.Content[j].Value] = true
		}
		break
	}

	return defined, nil
}
