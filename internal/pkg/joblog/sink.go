// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// Package joblog writes per-task diagnostic logs to the runner's own disk, so a job's
// output survives a failure of the path that uploads it to Gitea.
package joblog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// JobFileName holds the rows as they were sent to Gitea, DiagFileName every entry the
	// reporter saw, including the ones never uploaded.
	JobFileName  = "job.log"
	DiagFileName = "diag.log"

	dirNameLayout   = "20060102-150405"
	timestampLayout = "2006-01-02T15:04:05.000Z"
	queueSize       = 4096
)

// Options configures one task's sink.
type Options struct {
	Dir       string        // log.job.dir; empty disables the sink
	Retention time.Duration // log.job.retention; <=0 keeps every task directory
	MaxSize   int64         // log.job.max_size per file; <=0 is unlimited
	TaskID    int64
	StartedAt time.Time
	Header    string // metadata line written to both files, already masked
}

type record struct {
	diag bool
	line string
}

// Sink is one task's pair of log files. A nil *Sink is a no-op, so a caller that has
// the feature disabled needs no branches: every method tolerates it.
type Sink struct {
	dir       string
	ch        chan record
	done      chan struct{}
	dropped   atomic.Int64
	closeOnce sync.Once
}

// Open prepares the task's log directory. It returns nil when logging is disabled or
// when the files cannot be created: diagnostics must never fail a job.
func Open(opts Options) *Sink {
	return open(opts, queueSize)
}

func open(opts Options, queue int) *Sink {
	if opts.Dir == "" {
		return nil
	}
	if opts.StartedAt.IsZero() {
		opts.StartedAt = time.Now()
	}

	dir := filepath.Join(opts.Dir, fmt.Sprintf("%s-task-%d", opts.StartedAt.UTC().Format(dirNameLayout), opts.TaskID))
	// The files carry repository output, so they stay readable by this user only.
	if err := os.MkdirAll(dir, 0o700); err != nil {
		log.Warnf("cannot create job log directory %s: %v", dir, err)
		return nil
	}

	jobFile, err := openFile(dir, JobFileName, opts.MaxSize)
	if err != nil {
		log.Warnf("cannot create job log: %v", err)
		return nil
	}
	diagFile, err := openFile(dir, DiagFileName, opts.MaxSize)
	if err != nil {
		log.Warnf("cannot create job log: %v", err)
		jobFile.close()
		return nil
	}

	s := &Sink{
		dir:  dir,
		ch:   make(chan record, queue),
		done: make(chan struct{}),
	}
	go s.run(jobFile, diagFile, opts)
	return s
}

// Dir returns the task's log directory, empty when logging is disabled.
func (s *Sink) Dir() string {
	if s == nil {
		return ""
	}
	return s.dir
}

// Job records one row as it was handed to the uploader.
func (s *Sink) Job(t time.Time, content string) {
	s.send(record{line: stamp(t, content)})
}

// Diag records one runner-side diagnostic line. The caller masks it.
func (s *Sink) Diag(t time.Time, content string) {
	s.send(record{diag: true, line: stamp(t, content)})
}

// send never blocks: it is called from the reporter's hot path, under its lock. A line
// that finds the queue full is dropped and counted, and Close reports the count.
func (s *Sink) send(rec record) {
	if s == nil {
		return
	}
	select {
	case s.ch <- rec:
	default:
		s.dropped.Add(1)
	}
}

// Close writes trailer to both files and waits for everything queued to reach disk.
func (s *Sink) Close(trailer string) {
	if s == nil {
		return
	}
	s.closeOnce.Do(func() {
		if trailer != "" {
			// Blocking sends: the writer goroutine drains until the channel is closed.
			s.ch <- record{line: runnerLine(trailer)}
			s.ch <- record{diag: true, line: runnerLine(trailer)}
		}
		if dropped := s.dropped.Load(); dropped > 0 {
			note := runnerLine(fmt.Sprintf("dropped %d line(s): the log writer could not keep up", dropped))
			s.ch <- record{line: note}
			s.ch <- record{diag: true, line: note}
			log.Warnf("dropped %d line(s) from the job log in %s", dropped, s.dir)
		}
		close(s.ch)
		<-s.done
	})
}

func (s *Sink) run(jobFile, diagFile *sinkFile, opts Options) {
	defer close(s.done)
	defer jobFile.close()
	defer diagFile.close()

	prune(opts.Dir, s.dir, opts.Retention, opts.StartedAt)

	if opts.Header != "" {
		header := "# " + sanitize(opts.Header)
		jobFile.write(header)
		diagFile.write(header)
	}

	for rec := range s.ch {
		file := jobFile
		if rec.diag {
			file = diagFile
		}
		file.write(rec.line)
		// Flushed whenever the queue runs dry, so a runner killed mid-job still leaves
		// everything it had written behind.
		if len(s.ch) == 0 {
			jobFile.flush()
			diagFile.flush()
		}
	}
}

func stamp(t time.Time, content string) string {
	return t.UTC().Format(timestampLayout) + " " + content
}

func runnerLine(content string) string {
	return stamp(time.Now(), "[runner] "+content)
}

// sanitize keeps one record on one line. A row reaching Gitea cannot contain a real
// newline either (see DEVELOPMENT.md), so this only guards against a caller that
// formats its own multi-line diagnostic.
var lineEscaper = strings.NewReplacer("\r\n", `\n`, "\n", `\n`, "\r", `\r`)

func sanitize(s string) string {
	return lineEscaper.Replace(s)
}

type sinkFile struct {
	f      *os.File
	w      *bufio.Writer
	max    int64
	size   int64
	capped bool
}

func openFile(dir, name string, maxSize int64) (*sinkFile, error) {
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, err
	}
	return &sinkFile{f: f, w: bufio.NewWriter(f), max: maxSize}, nil
}

func (sf *sinkFile) write(line string) {
	if sf.capped {
		return
	}
	line = sanitize(line) + "\n"
	if sf.max > 0 && sf.size+int64(len(line)) > sf.max {
		sf.capped = true
		line = runnerLine(fmt.Sprintf("truncated: log.job.max_size of %d bytes reached", sf.max)) + "\n"
	}
	n, err := sf.w.WriteString(line)
	sf.size += int64(n)
	if err != nil {
		// Reported once: a failing write here is almost always a full disk, and retrying
		// every line would flood the runner's own log.
		sf.capped = true
		log.Warnf("cannot write %s: %v", sf.f.Name(), err)
	}
}

func (sf *sinkFile) flush() {
	if err := sf.w.Flush(); err != nil && !sf.capped {
		sf.capped = true
		log.Warnf("cannot write %s: %v", sf.f.Name(), err)
	}
}

func (sf *sinkFile) close() {
	sf.flush()
	if err := sf.f.Close(); err != nil {
		log.Warnf("cannot close %s: %v", sf.f.Name(), err)
	}
}

// taskDirPattern matches the directories Open creates, so pruning leaves anything else
// under log.job.dir alone.
var taskDirPattern = regexp.MustCompile(`^(\d{8}-\d{6})-task-\d+$`)

// prune removes the task directories older than retention. The age comes from the name
// rather than the mtime, which a reader or a backup tool can move.
func prune(root, keep string, retention time.Duration, now time.Time) {
	if retention <= 0 {
		return
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Warnf("cannot list job log directory %s: %v", root, err)
		return
	}

	cutoff := now.Add(-retention)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		match := taskDirPattern.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}
		started, err := time.ParseInLocation(dirNameLayout, match[1], time.UTC)
		if err != nil || !started.Before(cutoff) {
			continue
		}
		dir := filepath.Join(root, entry.Name())
		if dir == keep {
			continue
		}
		if err := os.RemoveAll(dir); err != nil {
			log.Warnf("cannot remove expired job log directory %s: %v", dir, err)
		}
	}
}
