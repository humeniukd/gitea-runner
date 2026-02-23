package workflowpattern

import (
	"fmt"
	"os"
)

type TraceWriter interface {
	Info(string, ...any)
}

type EmptyTraceWriter struct{}

func (*EmptyTraceWriter) Info(string, ...any) {
}

type StdOutTraceWriter struct{}

func (*StdOutTraceWriter) Info(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}
