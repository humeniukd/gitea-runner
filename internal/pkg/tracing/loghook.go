// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package tracing

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type logHook struct{}

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
