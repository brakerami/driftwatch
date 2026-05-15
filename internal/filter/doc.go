// Package filter provides lightweight pre-flight filtering of container
// candidates before drift detection is performed.
//
// Two independent axes of filtering are supported:
//
//   - Name filtering: restrict detection to an explicit list of container
//     names supplied via CLI flags or configuration. When no names are
//     provided, all container names are considered a match.
//
//   - Label selector filtering: restrict detection to containers whose
//     Docker labels satisfy a comma-separated key=value selector string,
//     following the same syntax used by kubectl --selector. When no
//     selector is provided, all label sets are considered a match.
//
// Both axes are ANDed together: a container must satisfy both the name
// filter and the label selector filter to be considered a candidate for
// drift detection.
//
// Usage:
//
//	opts := filter.Options{
//	    Names:         []string{"api", "worker"},
//	    LabelSelector: "env=production,managed-by=driftwatch",
//	}
//	f, err := filter.New(opts)
//	if err != nil { ... }
//
//	if f.Match(name, labels) {
//	    // proceed with drift detection
//	}
//
// The convenience method Match combines MatchName and MatchLabels into a
// single call for the common case where both checks are needed together.
package filter
