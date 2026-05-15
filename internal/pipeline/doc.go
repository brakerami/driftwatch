// Package pipeline provides a lightweight, ordered stage-based
// processing pipeline for drift findings.
//
// # Overview
//
// A Pipeline is constructed from one or more Stage implementations.
// When Run is called the findings slice flows through each stage in
// registration order. Any stage may filter, transform, or annotate
// the findings before passing them on.
//
// # Usage
//
//	p, err := pipeline.New(
//		myDeduplicator,
//		myRedactor,
//		myEnricher,
//	)
//	if err != nil { ... }
//
//	results, err := p.Run(rawFindings)
//
// # StageFunc
//
// For simple one-off transformations the StageFunc adapter allows
// an anonymous function to be used directly as a Stage without
// defining a concrete type.
package pipeline
