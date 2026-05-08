// Package policy provides rule-based severity evaluation for drift findings.
//
// A Policy is loaded from a JSON file containing a list of Rule definitions.
// Each rule maps a drift type (e.g. "env", "image", "label") to a Severity
// level (info / warning / critical) and can optionally be scoped to a set of
// container names.
//
// Usage:
//
//	p, err := policy.LoadFile("policy.json")
//	if err != nil { ... }
//	store := policy.New(p)
//	sev := store.Evaluate("env", "my-container")
//
Only rules whose Enabled field is true are considered during evaluation.
// When multiple rules match, the highest severity wins.
package policy
