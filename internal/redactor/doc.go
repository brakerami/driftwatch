// Package redactor provides utilities for masking sensitive environment
// variable values before they are included in drift findings, snapshots,
// or any external output.
//
// # Usage
//
//	r, err := redactor.New(nil) // use DefaultPatterns
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Mask a single value
//	safe := r.MaskValue("DB_PASSWORD", os.Getenv("DB_PASSWORD"))
//
//	// Mask a full env slice (e.g. from container inspect)
//	masked := r.MaskEnv(container.Env)
//
// Custom sensitive patterns can be supplied to New; they are matched
// case-insensitively as substrings of the environment variable key.
package redactor
