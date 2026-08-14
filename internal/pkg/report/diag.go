// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package report

import (
	"fmt"
	"strings"
	"time"

	"gitea.com/gitea/runner/internal/pkg/ver"

	runnerv1 "gitea.dev/actionslib/runner/v1"
	log "github.com/sirupsen/logrus"
)

// diagFields are the entry fields worth keeping in diag.log. A whitelist, so a field
// holding a whole matrix or environment does not end up in a file nobody reviews.
var diagFields = []string{"stage", "stepNumber", "stepID", "raw_output", "jobResult", "stepResult"}

// diagContextFields name the task context values the header reports, in this order.
var diagContextFields = []string{"repository", "workflow", "job", "run_id", "event_name", "ref", "sha"}

// taskHeader describes the task in one line, for a reader holding a directory of logs.
func taskHeader(task *runnerv1.Task, started time.Time) string {
	parts := []string{
		"gitea-runner " + ver.Version(),
		fmt.Sprintf("task=%d", task.Id),
	}
	for _, field := range diagContextFields {
		if v := task.Context.GetFields()[field].GetStringValue(); v != "" {
			parts = append(parts, field+"="+v)
		}
	}
	return strings.Join(parts, " ") + " started=" + started.UTC().Format(time.RFC3339)
}

// diagEntry mirrors one logrus entry into diag.log, including the entries the upload path
// drops: everything below the configured level, everything logged between the steps, and a
// suppressed ::debug:: command. Fire defers it so an ::add-mask:: the entry itself carries
// is registered before its value is written here.
func (r *Reporter) diagEntry(entry *log.Entry) {
	if r.sink == nil {
		return
	}
	parts := []string{strings.ToUpper(entry.Level.String())}
	for _, field := range diagFields {
		if v, ok := entry.Data[field]; ok {
			parts = append(parts, fmt.Sprintf("%s=%v", field, v))
		}
	}
	r.sink.Diag(entry.Time, r.mask(strings.Join(parts, " ")+" | "+entry.Message))
}

// diagf records a runner-side note in diag.log only: it is about reporting the job, not
// part of it, so it never reaches Gitea or job.log.
func (r *Reporter) diagf(format string, a ...any) {
	if r.sink == nil {
		return
	}
	r.sink.Diag(time.Now(), r.mask("[report] "+fmt.Sprintf(format, a...)))
}

// mask applies the row masking to a line built outside the row path.
func (r *Reporter) mask(content string) string {
	r.stateMu.RLock()
	defer r.stateMu.RUnlock()
	return r.maskLocked(content)
}

// maskLocked is mask for a caller already holding stateMu. Masking is what makes a
// secret ***: an entry reaching the reporter is not masked yet, because act's masker is
// a formatter and a hook sees the entry before it is formatted.
func (r *Reporter) maskLocked(content string) string {
	return strings.ToValidUTF8(r.logReplacer.Replace(content), "?")
}
