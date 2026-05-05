package inspector

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// ContainerInfo holds relevant runtime state for a container.
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Env     []string
	Labels  map[string]string
	Running bool
}

// Inspector wraps a Docker client to inspect running containers.
type Inspector struct {
	docker *client.Client
}

// New creates a new Inspector using the default Docker environment.
func New() (*Inspector, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("inspector: failed to create docker client: %w", err)
	}
	return &Inspector{docker: c}, nil
}

// ListRunning returns ContainerInfo for all currently running containers.
func (i *Inspector) ListRunning(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := i.docker.ContainerList(ctx, types.ContainerListOptions{All: false})
	if err != nil {
		return nil, fmt.Errorf("inspector: list containers: %w", err)
	}

	var result []ContainerInfo
	for _, c := range containers {
		info, err := i.Inspect(ctx, c.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}

// Inspect returns ContainerInfo for a single container by ID or name.
func (i *Inspector) Inspect(ctx context.Context, containerID string) (ContainerInfo, error) {
	data, err := i.docker.ContainerInspect(ctx, containerID)
	if err != nil {
		return ContainerInfo{}, fmt.Errorf("inspector: inspect %s: %w", containerID, err)
	}

	name := data.Name
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	return ContainerInfo{
		ID:      data.ID,
		Name:    name,
		Image:   data.Config.Image,
		Env:     data.Config.Env,
		Labels:  data.Config.Labels,
		Running: data.State.Running,
	}, nil
}

// Close releases the underlying Docker client resources.
func (i *Inspector) Close() error {
	return i.docker.Close()
}
