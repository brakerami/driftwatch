package drift

import (
	"fmt"

	"github.com/driftwatch/internal/inspector"
	"github.com/driftwatch/internal/manifest"
)

// DriftType categorizes the kind of drift detected.
type DriftType string

const (
	DriftTypeImage  DriftType = "image_mismatch"
	DriftTypeEnv    DriftType = "env_mismatch"
	DriftTypeMissing DriftType = "container_missing"
)

// Finding represents a single drift finding between a manifest spec
// and a running container.
type Finding struct {
	ContainerName string
	DriftType     DriftType
	Expected      string
	Actual        string
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: expected=%q actual=%q",
		f.DriftType, f.ContainerName, f.Expected, f.Actual)
}

// Detector compares manifest specs against live container state.
type Detector struct {
	inspector *inspector.Inspector
}

// New creates a Detector backed by the given Inspector.
func New(i *inspector.Inspector) *Detector {
	return &Detector{inspector: i}
}

// Detect returns all drift findings for the given manifest specs.
func (d *Detector) Detect(specs []manifest.ContainerSpec) ([]Finding, error) {
	var findings []Finding

	for _, spec := range specs {
		info, err := d.inspector.ContainerInfo(spec.Name)
		if err != nil {
			findings = append(findings, Finding{
				ContainerName: spec.Name,
				DriftType:     DriftTypeMissing,
				Expected:      spec.Name,
				Actual:        "not found",
			})
			continue
		}

		if info.Image != spec.Image {
			findings = append(findings, Finding{
				ContainerName: spec.Name,
				DriftType:     DriftTypeImage,
				Expected:      spec.Image,
				Actual:        info.Image,
			})
		}

		for k, expected := range spec.Env {
			actual, ok := info.EnvLookup(k)
			if !ok || actual != expected {
				findings = append(findings, Finding{
					ContainerName: spec.Name,
					DriftType:     DriftTypeEnv,
					Expected:      fmt.Sprintf("%s=%s", k, expected),
					Actual:        fmt.Sprintf("%s=%s", k, actual),
				})
			}
		}
	}

	return findings, nil
}
