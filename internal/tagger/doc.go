// Package tagger enriches drift findings with user-defined tags based on
// configurable matching rules.
//
// A Tagger holds a set of Rule values, each specifying an optional drift-type
// filter, an optional container-name filter, and one or more string tags to
// attach when both filters match a finding.
//
// # Basic usage
//
//	rules := []tagger.Rule{
//		{DriftType: "env", Tags: []string{"config", "sensitive"}},
//		{Container: "api-server", Tags: []string{"critical"}},
//	}
//
//	t, err := tagger.New(rules)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	tagged := t.Tag(findings)
//
// When DriftType or Container is left empty the corresponding filter is treated
// as a wildcard and matches every finding.
//
// Tags are deduplicated and appended to the Finding.Tags slice; the original
// slice is never mutated — Tag returns a new slice of findings.
package tagger
