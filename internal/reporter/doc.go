// Package reporter provides formatting and output capabilities for drift
// findings produced by the drift detector.
//
// Two output formats are supported:
//
//	- text: human-readable tabular output suitable for terminal use.
//	- json: machine-readable JSON, useful for piping into other tools
//	        or log aggregation systems.
//
// Example usage:
//
//	r := reporter.New(os.Stdout, reporter.FormatText)
//	if err := r.Write(findings); err != nil {
//	    log.Fatal(err)
//	}
package reporter
