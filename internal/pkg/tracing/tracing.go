// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// Package tracing initializes the OpenTelemetry tracing SDK and the spans the
// runner itself emits. Configuration merges the runner yml `otel:` section
// (higher precedence) with the standard OTEL_* environment variables; with
// neither an endpoint nor OTEL_CONFIG_FILE the SDK stays off entirely.
package tracing

import (
	"context"
	"maps"
	"os"
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
)

// Setup builds the tracing SDK from the merged yml and env configuration and
// installs it as the global tracer provider. It returns the shutdown function
// (always non-nil) and an error only when SDK creation failed, so the caller
// can log and keep running untraced.
func Setup(ctx context.Context, serviceName string, otelCfg config.Otel) (func(context.Context) error, error) {
	noop := func(context.Context) error { return nil }

	// The otel default propagator is a no-op composite, so the RPC and HTTP
	// instrumentation would inject nothing; the runner always propagates W3C,
	// operator-extendable through OTEL_PROPAGATORS.
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	// OTEL_CONFIG_FILE is the spec escape hatch: otelconf reads the file and
	// env itself, and the file supersedes anything passed programmatically.
	configFile := os.Getenv("OTEL_CONFIG_FILE")
	if configFile != "" {
		sdk, err := otelconf.NewSDK()
		if err != nil {
			return noop, err
		}
		otel.SetTracerProvider(sdk.TracerProvider())
		return sdk.Shutdown, nil
	}

	endpoint := effectiveEndpoint(otelCfg)
	if endpoint == "" {
		return noop, nil
	}

	sdk, err := otelconf.NewSDK(otelconf.WithOpenTelemetryConfiguration(buildConfiguration(serviceName, otelCfg, endpoint)))
	if err != nil {
		return noop, err
	}
	otel.SetTracerProvider(sdk.TracerProvider())
	return sdk.Shutdown, nil
}

// effectiveEndpoint resolves the endpoint with the config file winning over env.
func effectiveEndpoint(otelCfg config.Otel) string {
	if otelCfg.Endpoint != "" {
		return otelCfg.Endpoint
	}
	if v := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"); v != "" {
		return v
	}
	return os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
}

// effectiveProtocol resolves the protocol with the config file winning over env;
// an invalid file value is ignored in favor of the env/default.
func effectiveProtocol(otelCfg config.Otel) string {
	if otelCfg.Protocol != "" {
		if otelCfg.Protocol == protocolHTTP || otelCfg.Protocol == protocolGRPC {
			return otelCfg.Protocol
		}
		log.Warnf("invalid otel.protocol %q, falling back to env or %s", otelCfg.Protocol, protocolHTTP)
	}
	for _, envName := range []string{"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL", "OTEL_EXPORTER_OTLP_PROTOCOL"} {
		if v := os.Getenv(envName); v == protocolGRPC || v == protocolHTTP {
			return v
		}
	}
	return protocolHTTP
}

func buildConfiguration(serviceName string, otelCfg config.Otel, endpoint string) otelconf.OpenTelemetryConfiguration {
	// When the endpoint came from env, leave the exporter's Endpoint unset so
	// the otlptrace exporters apply the remaining OTEL_EXPORTER_OTLP_* vars
	// (compression, timeout) with their own spec precedence.
	fromYml := otelCfg.Endpoint != ""
	exporter := otelconf.SpanExporter{}
	switch effectiveProtocol(otelCfg) {
	case protocolGRPC:
		exp := &otelconf.OTLPGrpcExporter{}
		if fromYml {
			ep := endpoint
			exp.Endpoint = &ep
		}
		exp.Headers = mergedHeaders(otelCfg)
		exporter.OTLPGrpc = exp
	default:
		exp := &otelconf.OTLPHttpExporter{}
		if fromYml {
			ep := endpoint
			exp.Endpoint = &ep
		}
		exp.Headers = mergedHeaders(otelCfg)
		exporter.OTLPHttp = exp
	}

	sampler := envSampler()
	if otelCfg.SamplerRatio != nil {
		ratio := *otelCfg.SamplerRatio
		sampler = &otelconf.Sampler{
			ParentBased: &otelconf.ParentBasedSampler{
				Root: &otelconf.Sampler{
					TraceIDRatioBased: &otelconf.TraceIDRatioBasedSampler{Ratio: &ratio},
				},
			},
		}
	}

	return otelconf.OpenTelemetryConfiguration{
		FileFormat: "1.0",
		Resource: &otelconf.Resource{
			Attributes: []otelconf.AttributeNameValue{
				{Name: "service.name", Value: serviceName},
				{Name: "service.version", Value: ver.Version()},
			},
		},
		TracerProvider: &otelconf.TracerProvider{
			Processors: []otelconf.SpanProcessor{{
				Batch: &otelconf.BatchSpanProcessor{Exporter: exporter},
			}},
			Sampler: sampler,
		},
	}
}

// exporterHeaders merges the env OTLP headers with the yml ones, file wins per
// key. Needed whenever the yml section is non-empty because otlptrace's
// WithHeaders replaces the env-derived map wholesale.
func mergedHeaders(otelCfg config.Otel) []otelconf.NameStringValuePair {
	if len(otelCfg.Headers) == 0 {
		return nil
	}
	merged := make(map[string]string)
	for _, envName := range []string{"OTEL_EXPORTER_OTLP_HEADERS", "OTEL_EXPORTER_OTLP_TRACES_HEADERS"} {
		for pair := range strings.SplitSeq(os.Getenv(envName), ",") {
			k, v, found := strings.Cut(strings.TrimSpace(pair), "=")
			if found {
				merged[k] = v
			}
		}
	}
	maps.Copy(merged, otelCfg.Headers)
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
// variable, plus TRACEPARENT/TRACESTATE continuing the trace of ctx's span, so
// applications running in the job join the CI trace and export to the same
// collector.
//
// Collector credentials stay on the runner: header and TLS client-key variables
// are withheld from job containers, which are reachable by untrusted workflow
// code.
func JobEnv(ctx context.Context) map[string]string {
	envs := make(map[string]string)
	for _, kv := range os.Environ() {
		name, value, found := strings.Cut(kv, "=")
		if found && strings.HasPrefix(name, "OTEL_") && !isJobSecretVar(name) {
			envs[name] = value
		}
	}

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

// isJobSecretVar reports whether an OTEL_* variable carries collector
// credentials that must not reach job containers.
func isJobSecretVar(name string) bool {
	return strings.Contains(name, "HEADERS") ||
		strings.Contains(name, "CLIENT_KEY") ||
		strings.Contains(name, "CLIENT_CERTIFICATE")
}
