// Package sampler implements probabilistic sampling for drift findings.
//
// In high-throughput environments, every drift event may not need to be
// forwarded to alerting or history backends. Sampler allows the caller to
// retain a statistically representative fraction of findings while
// discarding the rest, reducing downstream load without losing signal.
//
// Usage:
//
//	s, err := sampler.New(0.25) // forward ~25 % of findings
//	if err != nil {
//		log.Fatal(err)
//	}
//	filtered := s.Sample(findings)
//
// A rate of 1.0 is a no-op (all findings pass through).
// A rate of 0.0 drops all findings.
package sampler
