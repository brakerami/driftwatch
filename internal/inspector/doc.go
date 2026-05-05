// Package inspector provides utilities for querying the Docker daemon
// to retrieve runtime state of running containers.
//
// It wraps the official Docker SDK client and exposes a simplified
// ContainerInfo struct that captures the fields relevant for drift
// detection: image, environment variables, and labels.
//
// Typical usage:
//
//	insp, err := inspector.New()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer insp.Close()
//
//	containers, err := insp.ListRunning(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, c := range containers {
//		fmt.Println(c.Name, c.Image)
//	}
package inspector
