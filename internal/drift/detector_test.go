package drift_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/manifest"
)

// fakeInspector satisfies the interface used by Detector without Docker.
type fakeContainerInfo struct {
	image string
	env   map[string]string
}

func (f fakeContainerInfo) EnvLookup(key string) (string, bool) {
	v, ok := f.env[key]
	return v, ok
}

func TestFinding_String(t *testing.T) {
	f := drift.Finding{
		ContainerName: "web",
		DriftType:     drift.DriftTypeImage,
		Expected:      "nginx:1.25",
		Actual:        "nginx:1.24",
	}
	got := f.String()
	if got == "" {
		t.Fatal("expected non-empty string from Finding.String()")
	}
}

func TestFinding_MissingType(t *testing.T) {
	f := drift.Finding{
		ContainerName: "db",
		DriftType:     drift.DriftTypeMissing,
		Expected:      "db",
		Actual:        "not found",
	}
	if f.DriftType != drift.DriftTypeMissing {
		t.Errorf("expected DriftTypeMissing, got %s", f.DriftType)
	}
}

func TestDriftTypes_AreDistinct(t *testing.T) {
	types := []drift.DriftType{
		drift.DriftTypeImage,
		drift.DriftTypeEnv,
		drift.DriftTypeMissing,
	}
	seen := make(map[drift.DriftType]bool)
	for _, dt := range types {
		if seen[dt] {
			t.Errorf("duplicate drift type: %s", dt)
		}
		seen[dt] = true
	}
}

func TestManifestSpec_EnvMap(t *testing.T) {
	spec := manifest.ContainerSpec{
		Name:  "api",
		Image: "myapp:v2",
		Env:   map[string]string{"PORT": "8080", "DEBUG": "false"},
	}
	if spec.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", spec.Env["PORT"])
	}
	if spec.Env["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %s", spec.Env["DEBUG"])
	}
}
