// Package suppression allows operators to declare known-acceptable drift
// findings that should be excluded from reports and webhook alerts.
//
// Rules are expressed as a JSON array and can be loaded from a file path
// supplied via the --suppress-file CLI flag. Each rule specifies a
// container name (or empty string to match all containers) and a drift
// type string such as "env", "image", or "label".
//
// Example rules file:
//
//	[
//	  { "container_name": "sidecar", "drift_type": "env",   "reason": "injected by platform" },
//	  { "container_name": "",        "drift_type": "label", "reason": "CI adds build labels" }
//	]
package suppression
