// Package filter provides container name and label filtering
// for scoping drift detection to a subset of running containers.
package filter

import "strings"

// Options holds the filtering criteria applied before drift detection.
type Options struct {
	// Names is an explicit list of container names to include.
	// If empty, all containers are considered.
	Names []string

	// LabelSelector is a comma-separated list of key=value pairs.
	// A container must match ALL supplied labels to be included.
	LabelSelector string
}

// Filter evaluates containers against a set of Options.
type Filter struct {
	opts Options
	labels map[string]string
}

// New creates a Filter from the given Options.
func New(opts Options) (*Filter, error) {
	labels, err := parseLabels(opts.LabelSelector)
	if err != nil {
		return nil, err
	}
	return &Filter{opts: opts, labels: labels}, nil
}

// MatchName reports whether the given container name passes the name filter.
func (f *Filter) MatchName(name string) bool {
	if len(f.opts.Names) == 0 {
		return true
	}
	for _, n := range f.opts.Names {
		if n == name {
			return true
		}
	}
	return false
}

// MatchLabels reports whether the supplied container labels satisfy
// every selector requirement in the filter.
func (f *Filter) MatchLabels(containerLabels map[string]string) bool {
	for k, v := range f.labels {
		if containerLabels[k] != v {
			return false
		}
	}
	return true
}

// parseLabels converts a comma-separated key=value string into a map.
// An empty selector string returns a nil map and no error.
func parseLabels(selector string) (map[string]string, error) {
	if selector == "" {
		return nil, nil
	}
	result := make(map[string]string)
	for _, pair := range strings.Split(selector, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, &ParseError{Raw: pair}
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

// ParseError is returned when a label selector cannot be parsed.
type ParseError struct {
	Raw string
}

func (e *ParseError) Error() string {
	return "filter: invalid label selector segment: " + e.Raw
}
