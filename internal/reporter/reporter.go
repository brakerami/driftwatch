// Package reporter formats and outputs drift findings.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Format controls how findings are rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a snapshot of findings at a point in time.
type Report struct {
	Timestamp  time.Time        `json:"timestamp"`
	Findings   []drift.Finding  `json:"findings"`
	TotalDrift int              `json:"total_drift"`
}

// Reporter writes drift reports to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to out in the given format.
// If out is nil, os.Stdout is used.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Write renders findings to the configured output.
func (r *Reporter) Write(findings []drift.Finding) error {
	rep := Report{
		Timestamp:  time.Now().UTC(),
		Findings:   findings,
		TotalDrift: len(findings),
	}
	switch r.format {
	case FormatJSON:
		return r.writeJSON(rep)
	default:
		return r.writeText(rep)
	}
}

func (r *Reporter) writeJSON(rep Report) error {
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}

func (r *Reporter) writeText(rep Report) error {
	fmt.Fprintf(r.out, "Drift Report — %s\n", rep.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(r.out, "Total findings: %d\n\n", rep.TotalDrift)
	if rep.TotalDrift == 0 {
		fmt.Fprintln(r.out, "No drift detected.")
		return nil
	}
	tw := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "CONTAINER\tTYPE\tDETAIL")
	for _, f := range rep.Findings {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", f.ContainerID, f.Type, f.Detail)
	}
	return tw.Flush()
}
