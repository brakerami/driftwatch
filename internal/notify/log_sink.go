package notify

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yourorg/driftwatch/internal/drift"
)

// LogSink writes a human-readable summary of findings to an io.Writer.
type LogSink struct {
	name   string
	writer io.Writer
	logger *log.Logger
}

// NewLogSink constructs a LogSink that writes to w.
// If w is nil, os.Stderr is used.
func NewLogSink(name string, w io.Writer) *LogSink {
	if w == nil {
		w = os.Stderr
	}
	return &LogSink{
		name:   name,
		writer: w,
		logger: log.New(w, "[notify/log] ", log.LstdFlags),
	}
}

// Name implements Sink.
func (l *LogSink) Name() string { return l.name }

// Send writes each finding as a single log line.
func (l *LogSink) Send(_ context.Context, findings []drift.Finding) error {
	for _, f := range findings {
		fmt.Fprintf(l.writer, "[notify/log] drift detected: %s\n", f)
	}
	return nil
}
