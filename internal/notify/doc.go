// Package notify implements a fan-out notification layer for drift findings.
//
// A Dispatcher holds one or more Sink implementations and delivers findings
// to all of them when Dispatch is called.  Built-in sinks:
//
//   - LogSink  – writes findings to any io.Writer (default: stderr)
//
// Custom sinks can be added by implementing the Sink interface:
//
//	type Sink interface {
//	    Name() string
//	    Send(ctx context.Context, findings []drift.Finding) error
//	}
//
// Example:
//
//	logSink := notify.NewLogSink("stderr", os.Stderr)
//	d, err := notify.New(logger, logSink)
//	if err != nil { ... }
//	d.Dispatch(ctx, findings)
package notify
