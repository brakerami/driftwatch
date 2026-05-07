// Package suppression allows operators to declare known-acceptable drift
// findings that should be excluded from reports and webhook alerts.
//
// Rules are expressed as a JSON array and can be loaded from a file path
// supplied via the --suppress-file CLI flag. Each rule specifies a
// container name (or empty string to match all containers) and a drift
// type string such as "env", "image", or "label".
//
// Matching behaviour:
//
//   - If container_name is non-empty, the rule applies only to the named
//     container. If it is an empty string, the rule applies to every
//     container in the pod.
//   - drift_type is matched case-insensitively against the finding type
//     produced by the drift-detection engine.
//   - reason is a free-text field for audit purposes and has no effect on
//     matching logic.
//
// Example rules file:
//
//	[
//	  { "container_name": "sidecar", "drift_type": "env",   "reason": "injected by platform" },
//	  { "container_name": "",        "drift_type": "label", "reason": "CI adds build labels" }
//	]
package suppression
