// Package history provides a persistent, file-backed store for recording
// and querying drift scan results over time.
//
// Each call to [Store.Record] appends a timestamped JSON entry for a
// named container. [Store.List] retrieves all entries for a container
// in chronological order, enabling callers to spot recurring drift,
// track when a container first diverged from its manifest, or surface
// trends in a dashboard.
//
// Typical usage:
//
//	store, err := history.NewStore("/var/lib/driftwatch/history")
//	if err != nil { ... }
//
//	// after each scan:
//	store.Record(containerName, findings)
//
//	// to review past results:
//	entries, err := store.List(containerName)
package history
