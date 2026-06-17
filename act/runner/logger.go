// Copyright 2022 The Gitea Authors. All rights reserved.
// Copyright 2020 The nektos/act Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"

	"gitea.com/gitea/runner/act/common"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const (
	// nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	magenta = 35
	cyan    = 36
	gray    = 37
)

const (
	rawOutputField      = "raw_output"
	scriptLineCyanField = "script_line_cyan"
)

var (
	colors    []int
	nextColor int
	mux       sync.Mutex
)

func init() {
	nextColor = 0
	colors = []int{
		blue, yellow, green, magenta, red, gray, cyan,
	}
}

type masksContextKey string

const masksContextKeyVal = masksContextKey("logrus.FieldLogger")

// Logger returns the appropriate logger for current context
func Masks(ctx context.Context) *[]string {
	val := ctx.Value(masksContextKeyVal)
	if val != nil {
		if masks, ok := val.(*[]string); ok {
			return masks
		}
	}
	return &[]string{}
}

// WithMasks adds a value to the context for the logger
func WithMasks(ctx context.Context, masks *[]string) context.Context {
	return context.WithValue(ctx, masksContextKeyVal, masks)
}

type JobLoggerFactory interface {
	WithJobLogger() *logrus.Logger
}

type jobLoggerFactoryContextKey string

var jobLoggerFactoryContextKeyVal = (jobLoggerFactoryContextKey)("jobloggerkey")

func WithJobLoggerFactory(ctx context.Context, factory JobLoggerFactory) context.Context {
	return context.WithValue(ctx, jobLoggerFactoryContextKeyVal, factory)
}

// WithJobLogger attaches a new logger to context that is aware of steps
func WithJobLogger(ctx context.Context, jobID, jobName string, config *Config, masks *[]string, matrix map[string]any) context.Context {
	ctx = WithMasks(ctx, masks)

	var logger *logrus.Logger
	if jobLoggerFactory, ok := ctx.Value(jobLoggerFactoryContextKeyVal).(JobLoggerFactory); ok && jobLoggerFactory != nil {
		logger = jobLoggerFactory.WithJobLogger()
	} else {
		var formatter logrus.Formatter
		if config.JSONLogger {
			formatter = &logrus.JSONFormatter{}
		} else {
			mux.Lock()
			defer mux.Unlock()
			nextColor++
			formatter = &jobLogFormatter{
				color:          colors[nextColor%len(colors)],
				logPrefixJobID: config.LogPrefixJobID,
			}
		}

		logger = logrus.New()
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.GetLevel())
		logger.SetFormatter(formatter)
	}

	{ // Adapt to Gitea
		if hook := common.LoggerHook(ctx); hook != nil {
			logger.AddHook(hook)
		}
		if config.JobLoggerLevel != nil {
			logger.SetLevel(*config.JobLoggerLevel)
		} else {
			logger.SetLevel(logrus.TraceLevel)
		}
	}

	logger.SetFormatter(&maskedFormatter{
		Formatter: logger.Formatter,
		masker:    valueMasker(config.InsecureSecrets, config.Secrets),
	})
	rtn := logger.WithFields(logrus.Fields{
		"job":    jobName,
		"jobID":  jobID,
		"dryrun": common.Dryrun(ctx),
		"matrix": matrix,
	}).WithContext(ctx)

	return common.WithLogger(ctx, rtn)
}

func WithCompositeLogger(ctx context.Context, masks *[]string) context.Context {
	ctx = WithMasks(ctx, masks)
	return common.WithLogger(ctx, common.Logger(ctx).WithFields(logrus.Fields{}).WithContext(ctx))
}

func WithCompositeStepLogger(ctx context.Context, stepID string) context.Context {
	val := common.Logger(ctx)
	stepIDs := make([]string, 0)

	if logger, ok := val.(*logrus.Entry); ok {
		if oldStepIDs, ok := logger.Data["stepID"].([]string); ok {
			stepIDs = append(stepIDs, oldStepIDs...)
		}
	}

	stepIDs = append(stepIDs, stepID)

	return common.WithLogger(ctx, common.Logger(ctx).WithFields(logrus.Fields{
		"stepID": stepIDs,
	}).WithContext(ctx))
}

func withStepLogger(ctx context.Context, stepNumber int, stepID, stepName, stageName string) context.Context {
	rtn := common.Logger(ctx).WithFields(logrus.Fields{
		"stepNumber": stepNumber,
		"step":       stepName,
		"stepID":     []string{stepID},
		"stage":      stageName,
	})
	return common.WithLogger(ctx, rtn)
}

type entryProcessor func(entry *logrus.Entry) *logrus.Entry

// minMaskLen is the minimum substring length to mask.
// Shorter substrings would cause excessive false-positive masking.
const minMaskLen = 3

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func AppendSecretMasker(oldnew []string, v string) []string {
	ret := oldnew

	for l := range strings.SplitSeq(v, "\n") {
		tm := strings.TrimSpace(l)
		if len(tm) < minMaskLen {
			continue
		}
		// Mask the full line and its reverse
		ret = append(ret, tm, "***")
		if rev := reverseString(tm); rev != tm {
			ret = append(ret, rev, "***")
		}
		// Mask all contiguous substrings of length >= minMaskLen
		// to prevent leakage via fold, cut, indexing, etc.
		runes := []rune(tm)
		seen := make(map[string]struct{})
		seen[tm] = struct{}{}
		for size := minMaskLen; size < len(runes); size++ {
			for start := 0; start+size <= len(runes); start++ {
				sub := string(runes[start : start+size])
				if _, ok := seen[sub]; ok {
					continue
				}
				seen[sub] = struct{}{}
				ret = append(ret, sub, "***")
			}
		}
	}

	return ret
}

// secretChars builds a set of all individual characters present in the secrets.
func secretChars(secrets map[string]string, masks []string) map[rune]struct{} {
	chars := make(map[rune]struct{})
	for _, v := range secrets {
		for _, r := range v {
			chars[r] = struct{}{}
		}
	}
	for _, v := range masks {
		for _, r := range v {
			chars[r] = struct{}{}
		}
	}
	return chars
}

// valueMasker applies secrets and ::add-mask:: patterns to every log entry, including
// raw_output (command/stream) lines; there is no bypass by field.
func valueMasker(insecureSecrets bool, secrets map[string]string) entryProcessor {
	var oldnew []string
	for _, v := range secrets {
		oldnew = AppendSecretMasker(oldnew, v)
	}
	oldnew = slices.Clip(oldnew)
	defReplacer := strings.NewReplacer(oldnew...)
	defChars := secretChars(secrets, nil)
	return func(entry *logrus.Entry) *logrus.Entry {
		if insecureSecrets {
			return entry
		}

		masks := Masks(entry.Context)

		var replacer *strings.Replacer
		var chars map[rune]struct{}
		if len(*masks) == 0 {
			replacer = defReplacer
			chars = defChars
		} else {
			cmasker := oldnew
			for _, v := range *masks {
				cmasker = AppendSecretMasker(cmasker, v)
			}
			replacer = strings.NewReplacer(cmasker...)
			chars = secretChars(secrets, *masks)
		}

		entry.Message = replacer.Replace(entry.Message)

		// If the entire message after replacement is a single character that
		// belongs to a secret, mask it. This prevents char-by-char leakage
		// via fold -w1, array indexing, etc.
		if msg := strings.TrimSpace(entry.Message); len([]rune(msg)) == 1 {
			r := []rune(msg)[0]
			if _, ok := chars[r]; ok {
				entry.Message = strings.Replace(entry.Message, msg, "***", 1)
			}
		}

		return entry
	}
}

type maskedFormatter struct {
	logrus.Formatter
	masker entryProcessor
}

func (f *maskedFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return f.Formatter.Format(f.masker(entry))
}

type jobLogFormatter struct {
	color          int
	logPrefixJobID bool
}

func (f *jobLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	if f.isColored(entry) {
		f.printColored(b, entry)
	} else {
		f.print(b, entry)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *jobLogFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry) {
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	var job any
	if f.logPrefixJobID {
		job = entry.Data["jobID"]
	} else {
		job = entry.Data["job"]
	}

	debugFlag := ""
	if entry.Level == logrus.DebugLevel {
		debugFlag = "[DEBUG] "
	}

	if entry.Data[rawOutputField] == true {
		if entry.Data[scriptLineCyanField] == true {
			fmt.Fprintf(b, "\x1b[%dm|\x1b[0m \x1b[36;1m%s\x1b[0m", f.color, entry.Message)
		} else {
			fmt.Fprintf(b, "\x1b[%dm|\x1b[0m %s", f.color, entry.Message)
		}
	} else if entry.Data["dryrun"] == true {
		fmt.Fprintf(b, "\x1b[1m\x1b[%dm\x1b[7m*DRYRUN*\x1b[0m \x1b[%dm[%s] \x1b[0m%s%s", gray, f.color, job, debugFlag, entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%dm[%s] \x1b[0m%s%s", f.color, job, debugFlag, entry.Message)
	}
}

func (f *jobLogFormatter) print(b *bytes.Buffer, entry *logrus.Entry) {
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	var job any
	if f.logPrefixJobID {
		job = entry.Data["jobID"]
	} else {
		job = entry.Data["job"]
	}

	debugFlag := ""
	if entry.Level == logrus.DebugLevel {
		debugFlag = "[DEBUG] "
	}

	if entry.Data[rawOutputField] == true {
		fmt.Fprintf(b, "[%s]   | %s", job, entry.Message)
	} else if entry.Data["dryrun"] == true {
		fmt.Fprintf(b, "*DRYRUN* [%s] %s%s", job, debugFlag, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] %s%s", job, debugFlag, entry.Message)
	}
}

func (f *jobLogFormatter) isColored(entry *logrus.Entry) bool {
	isColored := checkIfTerminal(entry.Logger.Out)

	if force, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && force != "0" {
		isColored = true
	} else if ok && force == "0" {
		isColored = false
	} else if os.Getenv("CLICOLOR") == "0" {
		isColored = false
	}

	return isColored
}

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
