// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const namespace = "act_runner"

// Label value constants for Prometheus metrics.
// Using constants prevents typos from silently creating new time-series.
const (
	LabelResultTask    = "task"
	LabelResultEmpty   = "empty"
	LabelResultError   = "error"
	LabelResultSuccess = "success"

	LabelMethodFetchTask  = "FetchTask"
	LabelMethodUpdateLog  = "UpdateLog"
	LabelMethodUpdateTask = "UpdateTask"
	LabelMethodDeclare    = "Declare"

	LabelStatusSuccess   = "success"
	LabelStatusFailure   = "failure"
	LabelStatusCancelled = "cancelled"
	LabelStatusSkipped   = "skipped"
	LabelStatusUnknown   = "unknown"
)

var (
	RunnerInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "info",
		Help:      "Runner metadata. Always 1. Labels carry version and name.",
	}, []string{"version", "name"})

	RunnerCapacity = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "capacity",
		Help:      "Configured maximum concurrent jobs.",
	})

	PollFetchTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "poll",
		Name:      "fetch_total",
		Help:      "Total number of FetchTask RPCs by result (task, empty, error).",
	}, []string{"result"})

	PollFetchDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "poll",
		Name:      "fetch_duration_seconds",
		Help:      "Latency of FetchTask RPCs, excluding expected long-poll timeouts.",
		Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10},
	})

	PollBackoffSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "poll",
		Name:      "backoff_seconds",
		Help:      "Current effective polling backoff interval.",
	})

	JobsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "job",
		Name:      "total",
		Help:      "Total jobs processed by status (success, failure, cancelled, skipped, unknown).",
	}, []string{"status"})

	JobDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "job",
		Name:      "duration_seconds",
		Help:      "Duration of job execution from start to finish.",
		Buckets:   prometheus.ExponentialBuckets(1, 2, 14), // 1s to ~4.5h
	})

	ReportLogTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "report",
		Name:      "log_total",
		Help:      "Total UpdateLog RPCs by result (success, error).",
	}, []string{"result"})

	ReportLogDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "report",
		Name:      "log_duration_seconds",
		Help:      "Latency of UpdateLog RPCs.",
		Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
	})

	ReportStateTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "report",
		Name:      "state_total",
		Help:      "Total UpdateTask (state) RPCs by result (success, error).",
	}, []string{"result"})

	ReportStateDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "report",
		Name:      "state_duration_seconds",
		Help:      "Latency of UpdateTask RPCs.",
		Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
	})

	ReportLogBufferRows = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "report",
		Name:      "log_buffer_rows",
		Help:      "Current number of buffered log rows awaiting send.",
	})

	ClientErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "client",
		Name:      "errors_total",
		Help:      "Total client RPC errors by method.",
	}, []string{"method"})
)

// Registry is the custom Prometheus registry used by the runner.
var Registry = prometheus.NewRegistry()

// Init registers all static metrics and the standard Go/process collectors.
// It must be called once before StartServer. Dynamic metrics (uptime, running
// jobs) are registered separately via prometheus.NewGaugeFunc in daemon.go
// because they require runtime values unavailable at package init time.
func Init() {
	Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		RunnerInfo, RunnerCapacity,
		PollFetchTotal, PollFetchDuration, PollBackoffSeconds,
		JobsTotal, JobDuration,
		ReportLogTotal, ReportLogDuration,
		ReportStateTotal, ReportStateDuration, ReportLogBufferRows,
		ClientErrors,
	)
}
