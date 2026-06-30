{{/*
Copyright 2026 The Gitea Authors. All rights reserved.
SPDX-License-Identifier: MIT
*/}}

{{/* Common labels shared by every rendered object. */}}
{{- define "gitea-runner.labels" -}}
app.kubernetes.io/name: gitea-runner
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: gitea-runner
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{/*
Selector labels for the runner pod. Only name+component — these MUST be
immutable for the lifetime of the Deployment.
*/}}
{{- define "gitea-runner.selectorLabels" -}}
app.kubernetes.io/name: gitea-runner
app.kubernetes.io/component: runner
{{- end -}}

{{/*
Namespace the chart renders into = the release namespace. The chart does NOT
create a Namespace object; it is the deployer's responsibility (e.g. Argo's
CreateNamespace=true, or an existing namespace). The runner also creates its
job pods HERE, so jobs run in the same namespace as the runner.
*/}}
{{- define "gitea-runner.namespace" -}}
{{- .Release.Namespace -}}
{{- end -}}

{{/*
Cache Service FQDN used as cache.host / connect_back_address. Derived from the
fixed Service name so the Phase 6 startup loopback guard always passes when
the embedded cache is on. Never override to a loopback value.
*/}}
{{- define "gitea-runner.cacheHost" -}}
gitea-runner-cache.{{ include "gitea-runner.namespace" . }}.svc.cluster.local
{{- end -}}

{{- define "gitea-runner.cacheURL" -}}
{{ printf "http://%s:%v" (include "gitea-runner.cacheHost" .) .Values.cache.port }}
{{- end -}}

{{/*
True when the embedded cache is active (not disabled, not replaced by an
external cache server). Gates the cache Service, the Deployment cache port,
and the cache config block.
*/}}
{{- define "gitea-runner.embeddedCache" -}}
{{- and .Values.cache.enabled (not .Values.cache.externalServer) -}}
{{- end -}}

{{/*
Effective metrics bind address. The code default (127.0.0.1:9101) is
unreachable from a Service; when metrics.service.enabled we force 0.0.0.0:9101.
*/}}
{{- define "gitea-runner.metricsAddr" -}}
{{- if .Values.metrics.addr -}}
{{- .Values.metrics.addr -}}
{{- else if .Values.metrics.service.enabled -}}
0.0.0.0:9101
{{- else -}}
127.0.0.1:9101
{{- end -}}
{{- end -}}

{{/*
labels env value (comma-joined) for GITEA_RUNNER_LABELS. The chart adds the
runner's own `k8s-native` routing label by default so workflows can target it
unless the operator overrides runner.labels.
*/}}
{{- define "gitea-runner.labelsEnv" -}}
{{- $labels := .Values.runner.labels -}}
{{- if not $labels -}}{{- $labels = list "k8s-native" -}}{{- end -}}
{{- join "," $labels -}}
{{- end -}}
