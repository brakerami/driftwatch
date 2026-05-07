// Package baseline captures and stores a known-good snapshot of a container's
// runtime configuration (environment variables and image) so that future runs
// can detect drift relative to that baseline rather than — or in addition to —
// a declarative manifest.
//
// # Usage
//
//	store, err := baseline.NewStore("/var/lib/driftwatch/baselines")
//	if err != nil { ... }
//
//	// Capture current state as the baseline.
//	entry := baseline.Entry{
//		ContainerID:   info.ID,
//		ContainerName: info.Name,
//		Image:         info.Image,
//		Env:           info.EnvMap(),
//	}
//	if err := store.Save(entry); err != nil { ... }
//
//	// Later, compare running state against the stored baseline.
//	base, err := store.Load(info.Name)
//	findings := baseline.DiffFromBaseline(base, info.EnvMap(), info.Image)
package baseline
