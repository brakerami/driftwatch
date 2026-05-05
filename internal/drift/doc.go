// Package drift provides the core drift detection logic for driftwatch.
//
// It compares a set of [manifest.ContainerSpec] definitions against the live
// state of running containers retrieved via [inspector.Inspector].
//
// Basic usage:
//
//	// Load the desired state from a manifest file.
//	// m, _ := manifest.LoadFromFile("containers.yaml")
//
//	// Connect to the Docker daemon.
//	// insp, _ := inspector.New(ctx)
//
//	// Run drift detection.
//	// detector := drift.New(insp)
//	// findings, _ := detector.Detect(m.Specs)
//	// for _, f := range findings {
//	//     fmt.Println(f)
//	// }
//
// A [Finding] describes a single discrepancy and includes the container name,
// the type of drift ([DriftType]), and the expected vs. actual values.
package drift
