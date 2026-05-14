// Package enricher attaches additional metadata to drift findings before
// they are forwarded to sinks or stored in history. Enrichment sources
// include container labels, host information, and user-supplied static tags.
package enricher

import (
	"fmt"
	"os"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Enricher adds metadata fields to a slice of findings.
type Enricher struct {
	hostname    string
	staticTags  map[string]string
	labelKeys   []string // container label keys to promote to metadata
}

// Config holds options for constructing an Enricher.
type Config struct {
	// StaticTags are key/value pairs unconditionally added to every finding.
	StaticTags map[string]string
	// LabelKeys is a list of Docker label keys whose values should be
	// promoted into finding metadata when present on the container.
	LabelKeys []string
}

// New creates an Enricher. It automatically resolves the current hostname;
// a non-fatal error is returned if hostname resolution fails (the field is
// left as "unknown").
func New(cfg Config) (*Enricher, error) {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	if cfg.StaticTags == nil {
		cfg.StaticTags = map[string]string{}
	}
	return &Enricher{
		hostname:   host,
		staticTags: cfg.StaticTags,
		labelKeys:  cfg.LabelKeys,
	}, nil
}

// Enrich returns a new slice of findings with metadata populated.
// The original findings are not mutated.
func (e *Enricher) Enrich(findings []drift.Finding) []drift.Finding {
	out := make([]drift.Finding, len(findings))
	for i, f := range findings {
		f.Metadata = mergeMetadata(f.Metadata, e.buildMeta(f))
		out[i] = f
	}
	return out
}

func (e *Enricher) buildMeta(f drift.Finding) map[string]string {
	m := map[string]string{
		"host": e.hostname,
	}
	for k, v := range e.staticTags {
		m[k] = v
	}
	for _, key := range e.labelKeys {
		if val, ok := f.ContainerLabels[key]; ok {
			m[fmt.Sprintf("label.%s", key)] = val
		}
	}
	return m
}

func mergeMetadata(base, extra map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(extra))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
