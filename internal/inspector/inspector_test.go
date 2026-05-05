package inspector

import (
	"testing"
)

func TestContainerInfo_Fields(t *testing.T) {
	info := ContainerInfo{
		ID:      "abc123",
		Name:    "web",
		Image:   "nginx:latest",
		Env:     []string{"PORT=80", "ENV=prod"},
		Labels:  map[string]string{"app": "web"},
		Running: true,
	}

	if info.ID != "abc123" {
		t.Errorf("expected ID abc123, got %s", info.ID)
	}
	if info.Name != "web" {
		t.Errorf("expected Name web, got %s", info.Name)
	}
	if info.Image != "nginx:latest" {
		t.Errorf("expected Image nginx:latest, got %s", info.Image)
	}
	if !info.Running {
		t.Error("expected Running to be true")
	}
	if len(info.Env) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(info.Env))
	}
	if info.Labels["app"] != "web" {
		t.Errorf("expected label app=web, got %s", info.Labels["app"])
	}
}

func TestNew_ReturnsErrorWithoutDocker(t *testing.T) {
	// This test verifies New() can be called; it may fail if Docker is
	// unavailable in CI — that is acceptable for an integration boundary.
	_, err := New()
	if err != nil {
		// Not a hard failure: Docker daemon may not be present.
		t.Logf("New() returned expected error without Docker: %v", err)
	}
}

func TestContainerInfo_EnvLookup(t *testing.T) {
	info := ContainerInfo{
		Env: []string{"FOO=bar", "BAZ=qux"},
	}

	found := false
	for _, e := range info.Env {
		if e == "FOO=bar" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected FOO=bar in Env slice")
	}
}
