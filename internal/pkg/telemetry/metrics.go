// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package telemetry

import (
	"context"
	"errors"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/resource"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
)

const cumulative = metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE

// pushMetrics exports the Prometheus registry on an interval and once more on shutdown, so short-lived runners report too.
func pushMetrics(c *client, res *resource.Resource, gatherer prometheus.Gatherer) func(context.Context) error {
	attrs := make([]*commonpb.KeyValue, 0, res.Len())
	for _, kv := range res.Attributes() {
		attrs = append(attrs, keyValue(string(kv.Key), kv.Value.String()))
	}
	start := uint64(time.Now().UnixNano())
	push := func(ctx context.Context) error {
		families, err := gatherer.Gather()
		return errors.Join(err, c.upload(ctx, &metricpb.MetricsData{ResourceMetrics: []*metricpb.ResourceMetrics{{
			Resource:     &resourcepb.Resource{Attributes: attrs},
			ScopeMetrics: []*metricpb.ScopeMetrics{{Scope: &commonpb.InstrumentationScope{Name: scope}, Metrics: convert(families, start)}},
			SchemaUrl:    res.SchemaURL(),
		}}}))
	}
	interval, timeout := millis("OTEL_METRIC_EXPORT_INTERVAL", time.Minute), millis("OTEL_METRIC_EXPORT_TIMEOUT", 30*time.Second)
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
				pushCtx, cancelPush := context.WithTimeout(ctx, timeout)
				if err := push(pushCtx); err != nil {
					log.WithError(err).Warn("OTLP metrics export failed")
				}
				cancelPush()
			}
		}
	}()
	return func(ctx context.Context) error {
		cancel()
		<-done
		return push(ctx)
	}
}

func millis(key string, fallback time.Duration) time.Duration {
	if ms, err := strconv.Atoi(os.Getenv(key)); err == nil && ms > 0 {
		return time.Duration(ms) * time.Millisecond
	}
	return fallback
}

// convert keeps Prometheus names, as the OTel compatibility spec asks; cumulative series start at runner startup.
func convert(families []*dto.MetricFamily, start uint64) []*metricpb.Metric {
	now := uint64(time.Now().UnixNano())
	metrics := make([]*metricpb.Metric, 0, len(families))
	for _, family := range families {
		metric := &metricpb.Metric{Name: family.GetName(), Description: family.GetHelp()}
		switch family.GetType() {
		case dto.MetricType_COUNTER:
			sum := &metricpb.Sum{AggregationTemporality: cumulative, IsMonotonic: true}
			for _, m := range family.GetMetric() {
				sum.DataPoints = append(sum.DataPoints, number(m, start, now, m.GetCounter().GetValue()))
			}
			metric.Data = &metricpb.Metric_Sum{Sum: sum}
		case dto.MetricType_GAUGE, dto.MetricType_UNTYPED:
			gauge := &metricpb.Gauge{}
			for _, m := range family.GetMetric() {
				value := m.GetGauge().GetValue()
				if m.Untyped != nil {
					value = m.GetUntyped().GetValue()
				}
				gauge.DataPoints = append(gauge.DataPoints, number(m, 0, now, value))
			}
			metric.Data = &metricpb.Metric_Gauge{Gauge: gauge}
		case dto.MetricType_HISTOGRAM:
			histogram := &metricpb.Histogram{AggregationTemporality: cumulative}
			for _, m := range family.GetMetric() {
				histogram.DataPoints = append(histogram.DataPoints, histogramPoint(m, start, now))
			}
			metric.Data = &metricpb.Metric_Histogram{Histogram: histogram}
		case dto.MetricType_SUMMARY:
			summary := &metricpb.Summary{}
			for _, m := range family.GetMetric() {
				point := &metricpb.SummaryDataPoint{Attributes: labels(m), StartTimeUnixNano: start, TimeUnixNano: now, Count: m.GetSummary().GetSampleCount(), Sum: m.GetSummary().GetSampleSum()}
				for _, q := range m.GetSummary().GetQuantile() {
					point.QuantileValues = append(point.QuantileValues, &metricpb.SummaryDataPoint_ValueAtQuantile{Quantile: q.GetQuantile(), Value: q.GetValue()})
				}
				summary.DataPoints = append(summary.DataPoints, point)
			}
			metric.Data = &metricpb.Metric_Summary{Summary: summary}
		default:
			continue
		}
		metrics = append(metrics, metric)
	}
	return metrics
}

func number(m *dto.Metric, start, now uint64, value float64) *metricpb.NumberDataPoint {
	return &metricpb.NumberDataPoint{Attributes: labels(m), StartTimeUnixNano: start, TimeUnixNano: now, Value: &metricpb.NumberDataPoint_AsDouble{AsDouble: value}}
}

// histogramPoint turns Prometheus' cumulative buckets into per-bucket counts, with the +Inf bucket last.
func histogramPoint(m *dto.Metric, start, now uint64) *metricpb.HistogramDataPoint {
	h := m.GetHistogram()
	buckets := h.GetBucket()
	if n := len(buckets); n > 0 && math.IsInf(buckets[n-1].GetUpperBound(), 1) {
		buckets = buckets[:n-1]
	}
	point := &metricpb.HistogramDataPoint{
		Attributes: labels(m), StartTimeUnixNano: start, TimeUnixNano: now, Count: h.GetSampleCount(), Sum: new(h.GetSampleSum()),
		ExplicitBounds: make([]float64, len(buckets)), BucketCounts: make([]uint64, len(buckets)+1),
	}
	var below uint64
	for i, b := range buckets {
		point.ExplicitBounds[i], point.BucketCounts[i] = b.GetUpperBound(), b.GetCumulativeCount()-below
		below = b.GetCumulativeCount()
	}
	point.BucketCounts[len(buckets)] = h.GetSampleCount() - below
	return point
}

func labels(m *dto.Metric) []*commonpb.KeyValue {
	attrs := make([]*commonpb.KeyValue, 0, len(m.GetLabel()))
	for _, l := range m.GetLabel() {
		attrs = append(attrs, keyValue(l.GetName(), l.GetValue()))
	}
	return attrs
}

func keyValue(key, value string) *commonpb.KeyValue {
	return &commonpb.KeyValue{Key: key, Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: value}}}
}
