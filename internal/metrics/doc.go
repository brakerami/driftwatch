// Package metrics provides a simple, thread-safe counter registry used
// throughout driftwatch to track operational statistics such as the number
// of containers scanned, drift findings raised, alerts suppressed, and
// webhook deliveries attempted.
//
// Usage:
//
//	reg := metrics.New()
//
//	// Increment a named counter anywhere in the application.
//	reg.Counter("scans_total").Inc()
//	reg.Counter("findings_total").Add(uint64(len(findings)))
//
//	// Obtain a point-in-time snapshot for reporting.
//	snap := reg.Snapshot()
//	fmt.Println(snap)
//
// Counters are created lazily on first access and are safe for concurrent
// use without additional locking by the caller.
package metrics
