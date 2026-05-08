// Package watchdog monitors container lifecycle events from the Docker daemon
// and triggers drift detection whenever a container is restarted or recreated.
//
// # Overview
//
// A Watchdog polls an EventSource on a configurable interval. When events are
// received they are forwarded to a Handler function, typically one that
// schedules an immediate drift check for the affected container.
//
// # Usage
//
//	src, _ := watchdog.NewDockerSource(dockerClient)
//	w, _ := watchdog.New(src, myHandler, 15*time.Second)
//	w.Run(ctx)
//
// EventSource is an interface so the watchdog is easily testable without a
// live Docker daemon.
package watchdog
