// Package rollup provides aggregation utilities for drift findings produced
// across a set of containers during a single driftwatch scan cycle.
//
// Use [Aggregate] to convert a map of container-name → []drift.Finding into a
// [Summary] that groups counts by drift type and surfaces the containers with
// the most findings.  The resulting [Summary] can be printed directly or
// embedded in reports and notifications.
//
// Example:
//
//	results := map[string][]drift.Finding{
//		"api":    apiFindings,
//		"worker": workerFindings,
//	}
//	summary := rollup.Aggregate(results)
//	fmt.Print(summary)
package rollup
