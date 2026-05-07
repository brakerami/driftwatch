// Package digest provides deterministic fingerprinting of container runtime
// state for use by the drift detection pipeline.
//
// # Overview
//
// Calling Compute on a ContainerState produces a stable SHA-256 hex string
// that encodes the image reference, environment variables, labels, and
// exposed ports. Because Go maps have non-deterministic iteration order,
// all map keys are sorted before hashing so that the same logical state
// always yields the same digest regardless of insertion order.
//
// # Usage
//
// The digest can be persisted alongside a snapshot (see internal/snapshot)
// so that the scheduler can skip a full field-by-field comparison when the
// digest has not changed since the last poll cycle, reducing CPU overhead
// for containers that are stable.
//
//	state := digest.ContainerState{
//		Image:  info.Image,
//		Env:    info.EnvMap(),
//		Labels: info.Labels,
//		Ports:  info.Ports,
//	}
//	sum, err := digest.Compute(state)
package digest
