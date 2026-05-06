// Package filter provides lightweight pre-flight filtering of container
// candidates before drift detection is performed.
//
// Two independent axes of filtering are supported:
//
//   - Name filtering: restrict detection to an explicit list of container
//     names supplied via CLI flags or configuration.
//
//   - Label selector filtering: restrict detection to containers whose
//     Docker labels satisfy a comma-separated key=value selector string,
//     following the same syntax used by kubectl --selector.
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
//	if f.MatchName(name) && f.MatchLabels(labels) {
//	    // proceed with drift detection
//	}
package filter
