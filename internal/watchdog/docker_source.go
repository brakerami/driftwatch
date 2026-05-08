package watchdog

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// DockerSource fetches container lifecycle events from the Docker daemon.
type DockerSource struct {
	cli *client.Client
}

// NewDockerSource creates a DockerSource backed by the local Docker socket.
func NewDockerSource(cli *client.Client) (*DockerSource, error) {
	if cli == nil {
		return nil, errNilClient
	}
	return &DockerSource{cli: cli}, nil
}

var errNilClient = fmt.Errorf("watchdog: docker client must not be nil")

// Events returns container lifecycle events that occurred after since.
func (d *DockerSource) Events(ctx context.Context, since time.Time) ([]Event, error) {
	f := filters.NewArgs()
	f.Add("type", "container")
	f.Add("event", "die")
	f.Add("event", "start")

	msgCh, errCh := d.cli.Events(ctx, events.ListOptions{
		Since:   since.Format(time.RFC3339Nano),
		Until:   time.Now().Format(time.RFC3339Nano),
		Filters: f,
	})

	var out []Event
	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				return out, nil
			}
			if e, ok := toEvent(msg); ok {
				out = append(out, e)
			}
		case err := <-errCh:
			if err != nil && err != context.Canceled {
				return out, err
			}
			return out, nil
		}
	}
}

func toEvent(msg events.Message) (Event, bool) {
	var kind EventKind
	switch msg.Action {
	case "die":
		kind = EventStop
	case "start":
		kind = EventRestart
	default:
		return Event{}, false
	}
	name := strings.TrimPrefix(msg.Actor.Attributes["name"], "/")
	return Event{
		ContainerID:   msg.Actor.ID,
		ContainerName: name,
		Kind:          kind,
		OccurredAt:    time.Unix(msg.Time, 0),
	}, true
}
