// Package enricher provides post-detection metadata enrichment for drift
// findings.
//
// After the drift detector identifies divergence between a running container
// and its source manifest, findings pass through the enricher before reaching
// reporters, alerters, or the history store. The enricher stamps each finding
// with:
//
//   - The hostname of the node running driftwatch.
//   - Any static key/value tags supplied via configuration (e.g. environment,
//     team, region).
//   - Selected Docker container label values promoted to finding metadata.
//
// Enrichment is purely additive; original finding fields are never modified.
package enricher
