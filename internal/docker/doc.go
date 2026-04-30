// Package docker provides utilities for interacting with the Docker daemon
// to retrieve runtime state of running containers.
//
// It exposes a Client type backed by the official Docker SDK, as well as a
// MockClient suitable for unit testing without a live Docker daemon.
//
// Typical usage:
//
//	c, err := docker.NewClient()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer c.Close()
//
//	containers, err := c.ListContainers(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, ctr := range containers {
//		fmt.Printf("%s -> %s\n", ctr.Name, ctr.Image)
//	}
package docker
