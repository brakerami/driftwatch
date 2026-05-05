package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultManifestNames lists filenames searched when no explicit path is given.
var DefaultManifestNames = []string{
	"driftwatch.yaml",
	"driftwatch.yml",
	".driftwatch.yaml",
	".driftwatch.yml",
}

// Loader discovers and loads manifests from the filesystem.
type Loader struct {
	// SearchDirs is the ordered list of directories to search for manifests.
	SearchDirs []string
}

// NewLoader creates a Loader that searches the given directories.
func NewLoader(dirs ...string) *Loader {
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	return &Loader{SearchDirs: dirs}
}

// Discover finds all manifest files across SearchDirs using DefaultManifestNames.
func (l *Loader) Discover() ([]string, error) {
	var found []string
	for _, dir := range l.SearchDirs {
		for _, name := range DefaultManifestNames {
			candidate := filepath.Join(dir, name)
			if _, err := os.Stat(candidate); err == nil {
				found = append(found, candidate)
			}
		}
	}
	return found, nil
}

// LoadAll discovers and parses all manifests, returning a combined slice of specs.
func (l *Loader) LoadAll() ([]*ContainerSpec, error) {
	paths, err := l.Discover()
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no manifest files found in: %s", strings.Join(l.SearchDirs, ", "))
	}

	var specs []*ContainerSpec
	for _, p := range paths {
		m, err := LoadFromFile(p)
		if err != nil {
			return nil, fmt.Errorf("loading %q: %w", p, err)
		}
		for i := range m.Containers {
			specs = append(specs, &m.Containers[i])
		}
	}
	return specs, nil
}
