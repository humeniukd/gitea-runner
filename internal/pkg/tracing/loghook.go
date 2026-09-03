// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package tracing

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// logHook stamps trace_id and span_id onto log entries whose context carries a
// valid span, so runner log lines can be correlated with the trace when OTel
// tracing is enabled. Entries without a span context — tracing off, or a log
// call made without WithContext — are left unchanged, so the hook is safe to
// install unconditionally.
type logHook struct{}

// NewLogHook returns a logrus hook that adds trace_id/span_id from the entry's
// context.
func NewLogHook() logrus.Hook { return logHook{} }

func (logHook) Levels() []logrus.Level { return logrus.AllLevels }

func (logHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}
	sc := trace.SpanContextFromContext(entry.Context)
	if !sc.IsValid() {
		return nil
	}
	entry.Data["trace_id"] = sc.TraceID().String()
	entry.Data["span_id"] = sc.SpanID().String()
	return nil
}
