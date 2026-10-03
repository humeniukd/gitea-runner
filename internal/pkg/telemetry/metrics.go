// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package telemetry

import (
	"cmp"
	"context"
	"errors"
	"math"
	"runtime"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/resource"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const cumulative = metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE

// pushMetrics also pushes on shutdown so short-lived runners report.
func pushMetrics(exporter *client, res *resource.Resource, gatherer prometheus.Gatherer) func(context.Context) error {
	attrs := make([]*commonpb.KeyValue, 0, res.Len())
	for _, kv := range res.Attributes() {
		attrs = append(attrs, keyValue(string(kv.Key), kv.Value.String()))
	}
	interval, timeout := millis("OTEL_METRIC_EXPORT_INTERVAL", time.Minute, time.Minute), millis("OTEL_METRIC_EXPORT_TIMEOUT", 30*time.Second, noLimit)
	push := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		families, err := gatherer.Gather()
		return errors.Join(err, exporter.upload(ctx, &metricpb.MetricsData{ResourceMetrics: []*metricpb.ResourceMetrics{{
			Resource:     &resourcepb.Resource{Attributes: attrs},
			ScopeMetrics: []*metricpb.ScopeMetrics{{Scope: &commonpb.InstrumentationScope{Name: scope}, Metrics: convert(families)}},
			SchemaUrl:    res.SchemaURL(),
		}}}))
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := push(ctx); err != nil && ctx.Err() == nil {
					log.WithError(err).Warn("OTLP metrics export failed, set OTEL_METRICS_EXPORTER=none to disable it")
				}
			}
		}
	}()
	return func(ctx context.Context) error {
		cancel()
		<-done
		return push(ctx)
	}
}

var ucum = map[string]string{
	"days": "d", "hours": "h", "minutes": "min", "seconds": "s", "milliseconds": "ms", "microseconds": "us", "nanoseconds": "ns",
	"bytes": "By", "kibibytes": "KiBy", "mebibytes": "MiBy", "gibibytes": "GiBy", "tebibytes": "TiBy",
	"kilobytes": "kBy", "megabytes": "MBy", "gigabytes": "GBy", "terabytes": "TBy",
	"meters": "m", "volts": "V", "amperes": "A", "joules": "J", "watts": "W", "grams": "g", "celsius": "Cel", "hertz": "Hz", "percent": "%",
}

// otlpMetric is how a Prometheus metric is sent over OTLP.
type otlpMetric struct {
	name, unit, description string
	upDown                  bool              // a gauge semconv defines as an UpDownCounter
	labels                  map[string]string // label → attribute
	values                  map[string]string // label value → attribute value
	extra                   []*commonpb.KeyValue
	unlimited               float64 // gauge value meaning no limit, not sent
}

var rpcMethods = map[string]string{
	"FetchTask":  "runner.v1.RunnerService/FetchTask",
	"UpdateLog":  "runner.v1.RunnerService/UpdateLog",
	"UpdateTask": "runner.v1.RunnerService/UpdateTask",
}

func rpcCall(method string) otlpMetric {
	return otlpMetric{
		name: "rpc.client.call.duration", unit: "s", description: "Duration of outgoing RPCs to Gitea.",
		extra: []*commonpb.KeyValue{keyValue("rpc.system.name", "connectrpc"), keyValue("rpc.method", rpcMethods[method])},
	}
}

func networkIO(direction string) otlpMetric {
	return otlpMetric{
		name: "process.network.io", unit: "By", description: "Network bytes transferred.",
		extra: []*commonpb.KeyValue{keyValue("network.io.direction", direction)},
	}
}

func goCPU(state string) otlpMetric {
	return otlpMetric{
		name: "go.cpu.time", unit: "s", description: "Estimated CPU time spent by the Go runtime.",
		extra: []*commonpb.KeyValue{keyValue("go.cpu.state", state)},
	}
}

// openFiles maps process_open_fds, which client_golang fills with the handle count on Windows.
func openFiles() otlpMetric {
	if runtime.GOOS == "windows" {
		return otlpMetric{name: "process.windows.handle.count", unit: "{handle}", upDown: true}
	}
	return otlpMetric{name: "process.unix.file_descriptor.count", unit: "{file_descriptor}", upDown: true}
}

// otlpMetrics renames metrics to the OTel semantic conventions. Runner metrics without one use gitea.runner.*.
// Unlisted Go and process metrics keep their Prometheus names, as semconv reserves go.* and process.* and
// dashboards know the client_golang names. Info metrics are resource attributes, so not sent.
var otlpMetrics = map[string]otlpMetric{
	"gitea_runner_info":                           {},
	"go_info":                                     {},
	"gitea_runner_capacity":                       {name: "gitea.runner.job.limit", unit: "{job}", upDown: true},
	"gitea_runner_uptime_seconds":                 {name: "process.uptime", unit: "s"},
	"gitea_runner_job_running":                    {name: "gitea.runner.job.active", unit: "{job}", upDown: true},
	"gitea_runner_job_capacity_utilization_ratio": {name: "gitea.runner.job.utilization", unit: "1"},
	"gitea_runner_job_duration_seconds":           {name: "gitea.runner.job.duration", unit: "s"},
	"gitea_runner_job_total": {
		name: "gitea.runner.jobs", unit: "{job}",
		labels: map[string]string{"status": "cicd.pipeline.result"},
		values: map[string]string{"cancelled": "cancellation", "skipped": "skip", "unknown": "error"},
	},
	"gitea_runner_poll_fetch_total": {
		name: "gitea.runner.poll.fetches", unit: "{fetch}",
		labels: map[string]string{"result": "gitea.runner.poll.result"},
	},
	"gitea_runner_poll_fetch_duration_seconds":   rpcCall("FetchTask"),
	"gitea_runner_poll_backoff_seconds":          {name: "gitea.runner.poll.backoff", unit: "s"},
	"gitea_runner_report_log_duration_seconds":   rpcCall("UpdateLog"),
	"gitea_runner_report_state_duration_seconds": rpcCall("UpdateTask"),
	"gitea_runner_report_log_buffer_rows":        {name: "gitea.runner.report.log.buffer", unit: "{row}", upDown: true},
	"gitea_runner_report_log_total": {
		name: "gitea.runner.report.log.requests", unit: "{request}",
		labels: map[string]string{"result": "gitea.runner.report.result"},
	},
	"gitea_runner_report_state_total": {
		name: "gitea.runner.report.state.requests", unit: "{request}",
		labels: map[string]string{"result": "gitea.runner.report.result"},
	},
	"gitea_runner_client_errors_total": {
		name: "cicd.system.errors", unit: "{error}",
		labels: map[string]string{"method": "rpc.method", "code": "error.type"}, values: rpcMethods,
		extra: []*commonpb.KeyValue{keyValue("cicd.system.component", "runner")},
	},
	"gitea_runner_state": {
		name: "cicd.worker.count", unit: "{worker}", upDown: true,
		labels: map[string]string{"state": "cicd.worker.state"},
		values: map[string]string{"idle": "available", "unavailable": "offline"},
	},
	"go_goroutines":                 {name: "go.goroutine.count", unit: "{goroutine}", upDown: true},
	"go_memstats_alloc_bytes_total": {name: "go.memory.allocated", unit: "By"},
	"go_memstats_mallocs_total":     {name: "go.memory.allocations", unit: "{allocation}"},
	"go_memstats_next_gc_bytes":     {name: "go.memory.gc.goal", unit: "By", upDown: true},
	"go_gc_gomemlimit_bytes":        {name: "go.memory.limit", unit: "By", upDown: true, unlimited: math.MaxInt64},
	"go_gc_gogc_percent":            {name: "go.config.gogc", unit: "%", upDown: true},
	"go_sched_gomaxprocs_threads":   {name: "go.processor.limit", unit: "{thread}", upDown: true},
	// Exported only if metrics.Init enables them with collectors.WithGoCollectorRuntimeMetrics.
	"go_sched_pauses_total_gc_seconds":                {name: "go.memory.gc.pause.duration", unit: "s"},
	"go_sched_latencies_seconds":                      {name: "go.schedule.duration", unit: "s"},
	"go_cpu_classes_user_cpu_seconds_total":           goCPU("user"),
	"go_cpu_classes_gc_total_cpu_seconds_total":       goCPU("gc"),
	"go_cpu_classes_scavenge_total_cpu_seconds_total": goCPU("scavenge"),
	"go_cpu_classes_idle_cpu_seconds_total":           goCPU("idle"),
	"process_resident_memory_bytes":                   {name: "process.memory.usage", unit: "By", upDown: true},
	"process_virtual_memory_bytes":                    {name: "process.memory.virtual", unit: "By", upDown: true},
	"process_open_fds":                                openFiles(),
	"process_network_receive_bytes_total":             networkIO("receive"),
	"process_network_transmit_bytes_total":            networkIO("transmit"),
}

// convert follows https://opentelemetry.io/docs/specs/otel/compatibility/prometheus_and_openmetrics/#prometheus-metric-points-to-otlp,
// renaming per otlpMetrics and merging metrics that share a name.
func convert(families []*dto.MetricFamily) []*metricpb.Metric {
	now := uint64(time.Now().UnixNano())
	metrics := make([]*metricpb.Metric, 0, len(families))
	byName := make(map[string]*metricpb.Metric, len(families))
	byFamily := make(map[string]*dto.MetricFamily, len(families))
	for _, family := range families {
		byFamily[family.GetName()] = family
		mapping, ok := otlpMetrics[family.GetName()]
		if !ok {
			mapping.name = family.GetName()
		} else if mapping.name == "" {
			continue
		}
		metric := &metricpb.Metric{
			Name: mapping.name, Description: cmp.Or(mapping.description, family.GetHelp()), Unit: cmp.Or(mapping.unit, ucum[family.GetUnit()], family.GetUnit()),
			Metadata: []*commonpb.KeyValue{keyValue("prometheus.type", strings.ToLower(family.GetType().String()))},
		}
		switch family.GetType() {
		case dto.MetricType_COUNTER:
			sum := &metricpb.Sum{AggregationTemporality: cumulative, IsMonotonic: true}
			for _, series := range family.GetMetric() {
				sum.DataPoints = append(sum.DataPoints, number(labels(series, mapping), unixNano(series.GetCounter().GetCreatedTimestamp()), now, series.GetCounter().GetValue()))
			}
			metric.Data = &metricpb.Metric_Sum{Sum: sum}
		case dto.MetricType_GAUGE:
			var points []*metricpb.NumberDataPoint
			for _, series := range family.GetMetric() {
				if value := series.GetGauge().GetValue(); mapping.unlimited == 0 || value != mapping.unlimited {
					points = append(points, number(labels(series, mapping), 0, now, value))
				}
			}
			if len(points) == 0 {
				continue
			}
			metric.Data = &metricpb.Metric_Gauge{Gauge: &metricpb.Gauge{DataPoints: points}}
			if mapping.upDown {
				metric.Data = &metricpb.Metric_Sum{Sum: &metricpb.Sum{AggregationTemporality: cumulative, DataPoints: points}}
			}
		case dto.MetricType_HISTOGRAM:
			histogram := &metricpb.Histogram{AggregationTemporality: cumulative}
			for _, series := range family.GetMetric() {
				histogram.DataPoints = append(histogram.DataPoints, histogramPoint(series, labels(series, mapping), now))
			}
			metric.Data = &metricpb.Metric_Histogram{Histogram: histogram}
		case dto.MetricType_SUMMARY:
			summary := &metricpb.Summary{}
			for _, series := range family.GetMetric() {
				point := &metricpb.SummaryDataPoint{
					Attributes: labels(series, mapping), StartTimeUnixNano: unixNano(series.GetSummary().GetCreatedTimestamp()), TimeUnixNano: now,
					Count: series.GetSummary().GetSampleCount(), Sum: series.GetSummary().GetSampleSum(),
				}
				for _, quantile := range series.GetSummary().GetQuantile() {
					point.QuantileValues = append(point.QuantileValues, &metricpb.SummaryDataPoint_ValueAtQuantile{Quantile: quantile.GetQuantile(), Value: quantile.GetValue()})
				}
				summary.DataPoints = append(summary.DataPoints, point)
			}
			metric.Data = &metricpb.Metric_Summary{Summary: summary}
		default:
			continue
		}
		switch existing := byName[metric.GetName()]; {
		case existing == nil:
			byName[metric.GetName()] = metric
			metrics = append(metrics, metric)
		case existing.GetSum() != nil:
			existing.GetSum().DataPoints = append(existing.GetSum().DataPoints, metric.GetSum().GetDataPoints()...)
		case existing.GetGauge() != nil:
			existing.GetGauge().DataPoints = append(existing.GetGauge().DataPoints, metric.GetGauge().GetDataPoints()...)
		case existing.GetHistogram() != nil:
			existing.GetHistogram().DataPoints = append(existing.GetHistogram().DataPoints, metric.GetHistogram().GetDataPoints()...)
		}
	}
	return append(metrics, derived(byFamily, now)...)
}

// derived adds the semconv Go metrics that are computed from other metrics rather than renamed.
func derived(families map[string]*dto.MetricFamily, now uint64) []*metricpb.Metric {
	first := func(name string) *dto.Metric {
		if series := families[name].GetMetric(); len(series) > 0 {
			return series[0]
		}
		return nil
	}
	gauge := func(name string) (float64, bool) {
		series := first(name)
		return series.GetGauge().GetValue(), series != nil
	}
	sum := func(name, unit, description string, monotonic bool, points ...*metricpb.NumberDataPoint) *metricpb.Metric {
		return &metricpb.Metric{Name: name, Unit: unit, Description: description, Data: &metricpb.Metric_Sum{Sum: &metricpb.Sum{
			AggregationTemporality: cumulative, IsMonotonic: monotonic, DataPoints: points,
		}}}
	}
	var metrics []*metricpb.Metric
	total, ok1 := gauge("go_memstats_sys_bytes")
	released, ok2 := gauge("go_memstats_heap_released_bytes")
	stack, ok3 := gauge("go_memstats_stack_inuse_bytes") // /memory/classes/heap/stacks:bytes, as semconv defines stack
	if ok1 && ok2 && ok3 {
		metrics = append(metrics, sum("go.memory.used", "By", "Memory used by the Go runtime.", false,
			number([]*commonpb.KeyValue{keyValue("go.memory.type", "stack")}, 0, now, stack),
			number([]*commonpb.KeyValue{keyValue("go.memory.type", "other")}, 0, now, total-released-stack)))
	}
	if gc := first("go_gc_duration_seconds"); gc != nil {
		metrics = append(metrics, sum("go.memory.gc.cycles", "{gc_cycle}", "Number of completed GC cycles.", true,
			number(nil, unixNano(gc.GetSummary().GetCreatedTimestamp()), now, float64(gc.GetSummary().GetSampleCount()))))
	}
	return metrics
}

func number(attrs []*commonpb.KeyValue, start, now uint64, value float64) *metricpb.NumberDataPoint {
	return &metricpb.NumberDataPoint{Attributes: attrs, StartTimeUnixNano: start, TimeUnixNano: now, Value: &metricpb.NumberDataPoint_AsDouble{AsDouble: value}}
}

// unixNano maps a nil timestamp to 0, leaving the start time unset as the spec asks.
func unixNano(timestamp *timestamppb.Timestamp) uint64 {
	return uint64(timestamp.AsTime().UnixNano())
}

func histogramPoint(series *dto.Metric, attrs []*commonpb.KeyValue, now uint64) *metricpb.HistogramDataPoint {
	histogram := series.GetHistogram()
	point := &metricpb.HistogramDataPoint{
		Attributes: attrs, StartTimeUnixNano: unixNano(histogram.GetCreatedTimestamp()), TimeUnixNano: now,
		Count: histogram.GetSampleCount(), Sum: histogram.SampleSum,
		BucketCounts: make([]uint64, len(histogram.GetBucket())+1),
	}
	var below uint64
	for i, bucket := range histogram.GetBucket() {
		point.ExplicitBounds = append(point.ExplicitBounds, bucket.GetUpperBound())
		point.BucketCounts[i], below = bucket.GetCumulativeCount()-below, bucket.GetCumulativeCount()
	}
	point.BucketCounts[len(histogram.GetBucket())] = histogram.GetSampleCount() - below
	return point
}

func labels(series *dto.Metric, mapping otlpMetric) []*commonpb.KeyValue {
	attrs := make([]*commonpb.KeyValue, 0, len(series.GetLabel())+len(mapping.extra))
	for _, label := range series.GetLabel() {
		attrs = append(attrs, keyValue(cmp.Or(mapping.labels[label.GetName()], label.GetName()), cmp.Or(mapping.values[label.GetValue()], label.GetValue())))
	}
	return append(attrs, mapping.extra...)
}

func keyValue(key, value string) *commonpb.KeyValue {
	return &commonpb.KeyValue{Key: key, Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: value}}}
}
