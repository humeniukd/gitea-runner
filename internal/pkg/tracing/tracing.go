// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// Package tracing initializes the OpenTelemetry SDK for the runner. The yml
// `otel:` section wins over the standard OTEL_* env vars per setting; a signal
// exports only when an endpoint is configured. Endpoint paths follow the OTLP
// spec: a signal-specific endpoint is used as-is, a generic one gets the
// signal path appended.
package tracing

import (
	"context"
	"maps"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"

	"gitea.com/gitea/runner/internal/pkg/config"
	"gitea.com/gitea/runner/internal/pkg/ver"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/otelconf"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	protocolHTTP = "http/protobuf"
	protocolGRPC = "grpc"

	tracesPath  = "/v1/traces"
	metricsPath = "/v1/metrics"
)

func Setup(ctx context.Context, serviceName string, otelCfg config.Otel) (func(context.Context) error, error) {
	noop := func(context.Context) error { return nil }

	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	// OTEL_CONFIG_FILE is the spec escape hatch: otelconf reads the file and
	// env itself, and the file supersedes anything passed programmatically.
	if os.Getenv("OTEL_CONFIG_FILE") != "" {
		sdk, err := otelconf.NewSDK()
		if err != nil {
			return noop, err
		}
		otel.SetTracerProvider(sdk.TracerProvider())
		if signalEnabled(otelCfg.Metrics.Enabled, otelCfg.Enabled) {
			otel.SetMeterProvider(sdk.MeterProvider())
		}
		return sdk.Shutdown, nil
	}

	tracesOn := signalEnabled(otelCfg.Traces.Enabled, otelCfg.Enabled) && tracesEndpointConfigured(otelCfg)
	metricsOn := signalEnabled(otelCfg.Metrics.Enabled, otelCfg.Enabled) && metricsEndpointConfigured(otelCfg)
	if !tracesOn && !metricsOn {
		return noop, nil
	}

	cfg := otelconf.OpenTelemetryConfiguration{
		FileFormat: "1.0",
		Resource: &otelconf.Resource{
			Attributes: resourceAttributes(serviceName, otelCfg.ResourceAttributes),
		},
	}
	if tracesOn {
		cfg.TracerProvider = traceProviderConfig(otelCfg)
	}
	if metricsOn {
		cfg.MeterProvider = meterProviderConfig(otelCfg)
	}
	sdk, err := otelconf.NewSDK(otelconf.WithOpenTelemetryConfiguration(cfg))
	if err != nil {
		return noop, err
	}
	otel.SetTracerProvider(sdk.TracerProvider())
	if metricsOn {
		otel.SetMeterProvider(sdk.MeterProvider())
	}
	return sdk.Shutdown, nil
}

func resourceAttributes(serviceName string, extra map[string]string) []otelconf.AttributeNameValue {
	attrs := []otelconf.AttributeNameValue{
		{Name: "service.name", Value: serviceName},
		{Name: "service.version", Value: ver.Version()},
	}
	merged := mergedResourceAttributes(os.Getenv("OTEL_RESOURCE_ATTRIBUTES"), extra)
	for _, k := range slices.Sorted(maps.Keys(merged)) {
		attrs = append(attrs, otelconf.AttributeNameValue{Name: k, Value: merged[k]})
	}
	return attrs
}

func signalEnabled(signal, shared *bool) bool {
	if signal != nil {
		return *signal
	}
	if shared != nil {
		return *shared
	}
	disabled, _ := strconv.ParseBool(os.Getenv("OTEL_SDK_DISABLED"))
	return !disabled
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func tracesEndpointConfigured(cfg config.Otel) bool {
	return firstNonEmpty(cfg.Traces.Endpoint, cfg.Endpoint,
		os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")) != ""
}

func metricsEndpointConfigured(cfg config.Otel) bool {
	return firstNonEmpty(cfg.Metrics.Endpoint, cfg.Endpoint,
		os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")) != ""
}

func ymlEndpoint(signal, generic, protocol, signalPath string) (string, bool) {
	switch {
	case signal != "":
		if protocol == protocolGRPC {
			return signal, true
		}
		return asIsHTTP(signal), true
	case generic != "":
		if protocol == protocolGRPC {
			return generic, true
		}
		return strings.TrimRight(generic, "/") + signalPath, true
	}
	return "", false
}

func asIsHTTP(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Path != "" {
		return raw
	}
	u.Path = "/"
	return u.String()
}

func endpointPtr(endpoint string, fromYml bool) *string {
	if !fromYml {
		return nil
	}
	return &endpoint
}

func signalProtocol(signalYml, sharedYml, signalEnv, genericEnv string) string {
	for _, v := range []string{signalYml, sharedYml} {
		if v == protocolHTTP || v == protocolGRPC {
			return v
		}
		if v != "" {
			log.Warnf("invalid otel protocol %q, ignoring", v)
		}
	}
	for _, envName := range []string{signalEnv, genericEnv} {
		if v := os.Getenv(envName); v == protocolGRPC || v == protocolHTTP {
			return v
		}
	}
	return protocolHTTP
}

func traceProviderConfig(cfg config.Otel) *otelconf.TracerProvider {
	protocol := signalProtocol(cfg.Traces.Protocol, cfg.Protocol,
		"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL", "OTEL_EXPORTER_OTLP_PROTOCOL")
	headers := mergedHeaders(cfg.Traces.Headers, cfg.Headers, "OTEL_EXPORTER_OTLP_TRACES_HEADERS")
	endpoint, fromYml := ymlEndpoint(cfg.Traces.Endpoint, cfg.Endpoint, protocol, tracesPath)

	ep := endpointPtr(endpoint, fromYml)
	exporter := otelconf.SpanExporter{}
	switch protocol {
	case protocolGRPC:
		exporter.OTLPGrpc = &otelconf.OTLPGrpcExporter{Headers: headers, Endpoint: ep}
	default:
		exporter.OTLPHttp = &otelconf.OTLPHttpExporter{Headers: headers, Endpoint: ep}
	}

	sampler := envSampler()
	if cfg.Traces.SamplerRatio != nil {
		ratio := *cfg.Traces.SamplerRatio
		sampler = &otelconf.Sampler{
			ParentBased: &otelconf.ParentBasedSampler{
				Root: &otelconf.Sampler{
					TraceIDRatioBased: &otelconf.TraceIDRatioBasedSampler{Ratio: &ratio},
				},
			},
		}
	}

	return &otelconf.TracerProvider{
		Processors: []otelconf.SpanProcessor{{
			Batch: &otelconf.BatchSpanProcessor{Exporter: exporter},
		}},
		Sampler: sampler,
	}
}

func meterProviderConfig(cfg config.Otel) *otelconf.MeterProvider {
	protocol := signalProtocol(cfg.Metrics.Protocol, cfg.Protocol,
		"OTEL_EXPORTER_OTLP_METRICS_PROTOCOL", "OTEL_EXPORTER_OTLP_PROTOCOL")
	endpoint, fromYml := ymlEndpoint(cfg.Metrics.Endpoint, cfg.Endpoint, protocol, metricsPath)
	headers := mergedHeaders(cfg.Metrics.Headers, cfg.Headers, "OTEL_EXPORTER_OTLP_METRICS_HEADERS")

	ep := endpointPtr(endpoint, fromYml)
	exporter := otelconf.PushMetricExporter{}
	switch protocol {
	case protocolGRPC:
		exporter.OTLPGrpc = &otelconf.OTLPGrpcMetricExporter{Headers: headers, Endpoint: ep}
	default:
		exporter.OTLPHttp = &otelconf.OTLPHttpMetricExporter{Headers: headers, Endpoint: ep}
	}

	return &otelconf.MeterProvider{
		Readers: []otelconf.MetricReader{{
			Periodic: &otelconf.PeriodicMetricReader{Exporter: exporter},
		}},
	}
}

func mergedHeaders(signal, shared map[string]string, signalEnv string) []otelconf.NameStringValuePair {
	if len(signal) == 0 && len(shared) == 0 {
		return nil
	}
	merged := make(map[string]string)
	for _, envName := range []string{"OTEL_EXPORTER_OTLP_HEADERS", signalEnv} {
		for pair := range strings.SplitSeq(os.Getenv(envName), ",") {
			k, v, found := strings.Cut(strings.TrimSpace(pair), "=")
			if found {
				merged[k] = v
			}
		}
	}
	maps.Copy(merged, shared)
	maps.Copy(merged, signal)
	pairs := make([]otelconf.NameStringValuePair, 0, len(merged))
	for k, v := range merged {
		pairs = append(pairs, otelconf.NameStringValuePair{Name: k, Value: &v})
	}
	return pairs
}

func envSampler() *otelconf.Sampler {
	switch strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER")) {
	case "", "always_on", "parentbased_always_on":
		return nil
	case "always_off", "parentbased_always_off":
		return &otelconf.Sampler{AlwaysOff: otelconf.AlwaysOffSampler{}}
	case "traceidratio":
		ratio := samplerArg()
		return &otelconf.Sampler{TraceIDRatioBased: &otelconf.TraceIDRatioBasedSampler{Ratio: &ratio}}
	case "parentbased_traceidratio":
		ratio := samplerArg()
		return &otelconf.Sampler{
			ParentBased: &otelconf.ParentBasedSampler{
				Root: &otelconf.Sampler{TraceIDRatioBased: &otelconf.TraceIDRatioBasedSampler{Ratio: &ratio}},
			},
		}
	default:
		log.Warnf("unsupported OTEL_TRACES_SAMPLER %q, using default sampling", os.Getenv("OTEL_TRACES_SAMPLER"))
		return nil
	}
}

func samplerArg() float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER_ARG")), 64)
	if err != nil || v < 0 || v > 1 {
		return 1.0
	}
	return v
}

// JobEnv translates the merged otel.* configuration back to OTEL_* env vars so
// job processes export with the same settings, plus TRACEPARENT/TRACESTATE from
// ctx. Job code can read the collector credentials it carries; see docs/otel.md.
func JobEnv(ctx context.Context, otelCfg config.Otel) map[string]string {
	envs := make(map[string]string)
	for _, kv := range os.Environ() {
		name, value, found := strings.Cut(kv, "=")
		if found && strings.HasPrefix(name, "OTEL_") {
			envs[name] = value
		}
	}
	setJobOtelEnv(envs, otelCfg)
	delete(envs, "OTEL_SERVICE_NAME")
	merged := mergedResourceAttributes(envs["OTEL_RESOURCE_ATTRIBUTES"], otelCfg.ResourceAttributes)
	maps.DeleteFunc(merged, func(k, _ string) bool { return strings.HasPrefix(k, "service.") })
	envs["OTEL_RESOURCE_ATTRIBUTES"] = serializeKeyValues(merged)

	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return envs
	}
	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(ctx, carrier)
	if tp := carrier.Get("traceparent"); tp != "" {
		envs["TRACEPARENT"] = tp
	}
	if ts := carrier.Get("tracestate"); ts != "" {
		envs["TRACESTATE"] = ts
	}
	return envs
}

func setJobOtelEnv(envs map[string]string, cfg config.Otel) {
	setEnvIf(envs, "OTEL_EXPORTER_OTLP_ENDPOINT", cfg.Endpoint)
	setEnvIf(envs, "OTEL_EXPORTER_OTLP_PROTOCOL", cfg.Protocol)
	setEnvIf(envs, "OTEL_EXPORTER_OTLP_HEADERS", serializeKeyValues(cfg.Headers))

	if signalEnabled(cfg.Traces.Enabled, cfg.Enabled) {
		setEnvIf(envs, "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", cfg.Traces.Endpoint)
		setEnvIf(envs, "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL", cfg.Traces.Protocol)
		setEnvIf(envs, "OTEL_EXPORTER_OTLP_TRACES_HEADERS", serializeKeyValues(cfg.Traces.Headers))
		if cfg.Traces.SamplerRatio != nil {
			envs["OTEL_TRACES_SAMPLER"] = "parentbased_traceidratio"
			envs["OTEL_TRACES_SAMPLER_ARG"] = strconv.FormatFloat(*cfg.Traces.SamplerRatio, 'g', -1, 64)
		}
	} else {
		for _, k := range []string{
			"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "OTEL_EXPORTER_OTLP_TRACES_PROTOCOL",
			"OTEL_EXPORTER_OTLP_TRACES_HEADERS", "OTEL_TRACES_SAMPLER", "OTEL_TRACES_SAMPLER_ARG",
		} {
			delete(envs, k)
		}
		envs["OTEL_TRACES_EXPORTER"] = "none"
	}

	if signalEnabled(cfg.Metrics.Enabled, cfg.Enabled) {
		setEnvIf(envs, "OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", cfg.Metrics.Endpoint)
		setEnvIf(envs, "OTEL_EXPORTER_OTLP_METRICS_PROTOCOL", cfg.Metrics.Protocol)
		setEnvIf(envs, "OTEL_EXPORTER_OTLP_METRICS_HEADERS", serializeKeyValues(cfg.Metrics.Headers))
	} else {
		for _, k := range []string{
			"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT", "OTEL_EXPORTER_OTLP_METRICS_PROTOCOL",
			"OTEL_EXPORTER_OTLP_METRICS_HEADERS",
		} {
			delete(envs, k)
		}
		envs["OTEL_METRICS_EXPORTER"] = "none"
	}
}

func setEnvIf(envs map[string]string, key, value string) {
	if value != "" {
		envs[key] = value
	}
}

const defaultDeploymentEnvironment = "test"

func mergedResourceAttributes(envValue string, yml map[string]string) map[string]string {
	merged := map[string]string{"deployment.environment.name": defaultDeploymentEnvironment}
	for pair := range strings.SplitSeq(envValue, ",") {
		k, v, found := strings.Cut(strings.TrimSpace(pair), "=")
		if found && k != "" {
			merged[k] = v
		}
	}
	maps.Copy(merged, yml)
	return merged
}

func serializeKeyValues(attrs map[string]string) string {
	pairs := make([]string, 0, len(attrs))
	for _, k := range slices.Sorted(maps.Keys(attrs)) {
		pairs = append(pairs, k+"="+attrs[k])
	}
	return strings.Join(pairs, ",")
}
