package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/manifest"
)

const sampleYAML = `
version: "1"
containers:
  - name: web
    image: nginx:1.25
    env:
      PORT: "8080"
      ENV: production
    ports:
      - "80:80"
    labels:
      app: web
    restart_policy: always
  - name: worker
    image: myapp/worker:latest
    command: ["/app/worker", "--verbose"]
`

func TestParseManifest(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "driftwatch.yaml")
	if err := os.WriteFile(path, []byte(sampleYAML), 0o644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}

	m, err := manifest.LoadFromFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(m.Containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(m.Containers))
	}

	spec, err := m.SpecByName("web")
	if err != nil {
		t.Fatalf("SpecByName: %v", err)
	}
	if spec.Image != "nginx:1.25" {
		t.Errorf("expected image nginx:1.25, got %q", spec.Image)
	}
	if spec.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", spec.Env["PORT"])
	}
	if spec.RestartPolicy != "always" {
		t.Errorf("expected restart_policy=always, got %q", spec.RestartPolicy)
	}
}

func TestSpecByName_NotFound(t *testing.T) {
	m := &manifest.Manifest{
		Containers: []manifest.ContainerSpec{{Name: "api", Image: "api:v1"}},
	}
	_, err := m.SpecByName("missing")
	if err == nil {
		t.Error("expected error for missing container, got nil")
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	_, err := manifest.LoadFromFile("/nonexistent/path/manifest.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
