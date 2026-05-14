package drift

// Finding represents a single detected divergence between the desired state
// described by a manifest and the observed state of a running container.
//
// ContainerLabels holds the raw Docker labels of the inspected container and
// is populated by the inspector before findings reach the enricher.
// Metadata is an open-ended map that downstream components (enricher, tagger,
// policy) may populate with additional context.
type Finding struct {
	// Container is the name or ID of the container where drift was detected.
	Container string
	// Type categorises the kind of drift (see TypeEnv, TypeImage, etc.).
	Type DriftType
	// Field is the specific attribute that drifted (e.g. an env var name).
	Field string
	// Expected is the value prescribed by the manifest.
	Expected string
	// Actual is the value observed on the running container.
	Actual string
	// ContainerLabels are the Docker labels attached to the container.
	ContainerLabels map[string]string
	// Metadata holds enrichment data added after detection.
	Metadata map[string]string
}
