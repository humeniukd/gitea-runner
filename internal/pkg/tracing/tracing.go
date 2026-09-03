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

// Setup builds the tracing and metrics SDKs from the merged yml and env
// configuration and installs them as the global providers. It returns the
// shutdown function (always non-nil) and an error only when SDK creation
// failed, so the caller can log and keep running untraced.
func Setup(ctx context.Context, serviceName string, otelCfg config.Otel) (func(context.Context) error, error) {
	noop := func(context.Context) error { return nil }

	// The otel default propagator is a no-op composite, so the RPC and HTTP
	// instrumentation would inject nothing; the runner always propagates W3C,
	// operator-extendable through OTEL_PROPAGATORS.
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

// resourceAttributes builds the OTel resource identifying the runner. Attributes
// resolve as the built-in default, then OTEL_RESOURCE_ATTRIBUTES, then the yml
// attributes — most specific winning per key; keys are sorted for determinism.
// They are applied after otelconf's env detector so they take effect as resolved.
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

// signalEnabled resolves a signal's enabled flag: signal, then the shared
// otel.enabled, then the OTEL_SDK_DISABLED env, defaulting to true.
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

// tracesEndpointConfigured reports whether any yml or env source provides a
// traces endpoint, which is what gates the signal on.
func tracesEndpointConfigured(cfg config.Otel) bool {
	return firstNonEmpty(cfg.Traces.Endpoint, cfg.Endpoint,
		os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")) != ""
}

func metricsEndpointConfigured(cfg config.Otel) bool {
	return firstNonEmpty(cfg.Metrics.Endpoint, cfg.Endpoint,
		os.Getenv("OTEL_EXPORTER_OTLP_METRIC_ENDPOINT"), os.Getenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"),
		os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")) != ""
}

// ymlEndpoint resolves the endpoint URL to configure on the exporter from the
// yml sources only, applying OTLP path semantics: a signal endpoint is used
// as-is (a bare host posts to the root path), a generic endpoint has the
// signal path appended. Returns ("", false) when only env provides the
// endpoint, so the otlp exporter applies its own env handling.
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

// asIsHTTP keeps a signal-specific HTTP endpoint verbatim; a bare host with no
// path is given the root path so otelconf does not append a default signal
// path to it.
func asIsHTTP(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Path != "" {
		return raw
	}
	u.Path = "/"
	return u.String()
}

// signalProtocol resolves the protocol: yml signal, yml shared, env signal, env
// generic, then http/protobuf. An invalid yml value is ignored with a warning.
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

	exporter := otelconf.SpanExporter{}
	switch protocol {
	case protocolGRPC:
		exp := &otelconf.OTLPGrpcExporter{Headers: headers}
		if fromYml {
			ep := endpoint
			exp.Endpoint = &ep
		}
		exporter.OTLPGrpc = exp
	default:
		exp := &otelconf.OTLPHttpExporter{Headers: headers}
		if fromYml {
			ep := endpoint
			exp.Endpoint = &ep
		}
		exporter.OTLPHttp = exp
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

	exporter := otelconf.PushMetricExporter{}
	switch protocol {
	case protocolGRPC:
		exp := &otelconf.OTLPGrpcMetricExporter{Headers: headers}
		if fromYml {
			ep := endpoint
			exp.Endpoint = &ep
		}
		exporter.OTLPGrpc = exp
	default:
		exp := &otelconf.OTLPHttpMetricExporter{Headers: headers}
		if fromYml {
			ep := endpoint
			exp.Endpoint = &ep
		}
		exporter.OTLPHttp = exp
	}

	return &otelconf.MeterProvider{
		Readers: []otelconf.MetricReader{{
			Periodic: &otelconf.PeriodicMetricReader{Exporter: exporter},
		}},
	}
}

// mergedHeaders merges the env OTLP headers with the shared and signal yml
// ones, most specific winning per key. Needed whenever any yml header is set
// because otlptrace's WithHeaders replaces the env-derived map wholesale.
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
		name := k
		value := v
		pairs = append(pairs, otelconf.NameStringValuePair{Name: name, Value: &value})
	}
	return pairs
}

// envSampler honors OTEL_TRACES_SAMPLER; unknown values keep the SDK default.
func envSampler() *otelconf.Sampler {
	switch strings.TrimSpace(os.Getenv("OTEL_TRACES_SAMPLER")) {
	case "", "always_on", "parentbased_always_on":
		return nil // SDK default is parent-based always-on
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

// JobEnv returns the env to inject into job containers: every process OTEL_*
// variable, the yml otel.* settings translated to their OTEL_* env vars (so
// file config reaches jobs too, since yml is never in the process env), and
// TRACEPARENT/TRACESTATE continuing the trace of ctx's span. Applications in the
// job then join the CI trace and export with the same settings. yml wins over
// the process env, per setting.
//
// A disabled signal contributes no OTLP env vars and is turned off in-job with
// OTEL_{TRACES,METRICS}_EXPORTER=none, so a shared endpoint does not re-enable it.
//
// Security limitation: OTEL_* vars are passed through unfiltered, so job code
// (including untrusted pull requests) can read the collector credentials in
// OTEL_EXPORTER_OTLP*_HEADERS and TLS client key/cert vars. Point the runner at
// a local collector that accepts OTLP without auth and forwards to the real
// backend, rather than putting backend credentials on the runner.
func JobEnv(ctx context.Context, otelCfg config.Otel) map[string]string {
	envs := make(map[string]string)
	for _, kv := range os.Environ() {
		name, value, found := strings.Cut(kv, "=")
		if found && strings.HasPrefix(name, "OTEL_") {
			envs[name] = value
		}
	}
	setJobOtelEnv(envs, otelCfg)
	envs["OTEL_RESOURCE_ATTRIBUTES"] = serializeKeyValues(
		mergedResourceAttributes(envs["OTEL_RESOURCE_ATTRIBUTES"], otelCfg.ResourceAttributes))

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

// setJobOtelEnv overlays the yml otel.* settings onto the job env as their
// OTEL_* env vars, gating each signal on its enabled flag. yml wins per setting.
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

// defaultDeploymentEnvironment is the deployment.environment.name applied to the
// runner and to job containers unless overridden by env or yml config.
const defaultDeploymentEnvironment = "test"

// mergedResourceAttributes resolves the resource attributes: the built-in
// default, then OTEL_RESOURCE_ATTRIBUTES (envValue, comma-separated key=value),
// then the yml attributes — most specific winning per key, so a caller-set
// attribute never displaces the default of a different key.
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

// serializeKeyValues renders a map as a sorted comma-separated key=value list,
// the format of OTEL_RESOURCE_ATTRIBUTES and OTEL_EXPORTER_OTLP_HEADERS.
func serializeKeyValues(attrs map[string]string) string {
	pairs := make([]string, 0, len(attrs))
	for _, k := range slices.Sorted(maps.Keys(attrs)) {
		pairs = append(pairs, k+"="+attrs[k])
	}
	return strings.Join(pairs, ",")
}
