package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ContainerInfo holds relevant runtime info for a running container.
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Env     []string
	Ports   []types.Port
	Labels  map[string]string
}

// Client wraps the Docker SDK client.
type Client struct {
	docker *client.Client
}

// NewClient creates a new Docker client using environment variables.
func NewClient() (*Client, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{docker: c}, nil
}

// ListContainers returns info about all running containers.
func (c *Client) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := c.docker.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("status", "running")),
	})
	if err != nil {
		return nil, err
	}

	var result []ContainerInfo
	for _, ctr := range containers {
		inspect, err := c.docker.ContainerInspect(ctx, ctr.ID)
		if err != nil {
			return nil, err
		}
		name := ""
		if len(inspect.Name) > 1 {
			name = inspect.Name[1:] // strip leading '/'
		}
		result = append(result, ContainerInfo{
			ID:     ctr.ID[:12],
			Name:   name,
			Image:  ctr.Image,
			Env:    inspect.Config.Env,
			Ports:  ctr.Ports,
			Labels: ctr.Labels,
		})
	}
	return result, nil
}

// Close releases the Docker client resources.
func (c *Client) Close() error {
	return c.docker.Close()
}
