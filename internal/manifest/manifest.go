// Package manifest handles loading and parsing of container source manifests.
package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ContainerSpec describes the desired state of a container as defined in a manifest file.
type ContainerSpec struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Env         map[string]string `yaml:"env"`
	Ports       []string          `yaml:"ports"`
	Labels      map[string]string `yaml:"labels"`
	Command     []string          `yaml:"command"`
	RestartPolicy string          `yaml:"restart_policy"`
}

// Manifest represents a driftwatch manifest file containing one or more container specs.
type Manifest struct {
	Version    string          `yaml:"version"`
	Containers []ContainerSpec `yaml:"containers"`
}

// LoadFromFile reads and parses a manifest YAML file from the given path.
func LoadFromFile(path string) (*Manifest, error) {
	path = filepath.Clean(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest %q: %w", path, err)
	}
	return parse(data)
}

// parse decodes raw YAML bytes into a Manifest struct.
func parse(data []byte) (*Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest yaml: %w", err)
	}
	if m.Version == "" {
		m.Version = "1"
	}
	return &m, nil
}

// SpecByName returns the ContainerSpec with the given name, or an error if not found.
func (m *Manifest) SpecByName(name string) (*ContainerSpec, error) {
	for i := range m.Containers {
		if m.Containers[i].Name == name {
			return &m.Containers[i], nil
		}
	}
	return nil, fmt.Errorf("no container spec named %q in manifest", name)
}
