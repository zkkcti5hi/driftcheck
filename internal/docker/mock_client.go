package docker

import (
	"context"
	"github.com/docker/docker/api/types"
)

// MockClient is a test double for the Docker client.
type MockClient struct {
	Containers []ContainerInfo
	Err        error
}

// ListContainers returns the pre-configured containers or error.
func (m *MockClient) ListContainers(_ context.Context) ([]ContainerInfo, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Containers, nil
}

// SampleContainers returns a slice of ContainerInfo useful in tests.
func SampleContainers() []ContainerInfo {
	return []ContainerInfo{
		{
			ID:    "abc123def456",
			Name:  "web",
			Image: "nginx:1.25",
			Env:   []string{"PORT=80", "ENV=production"},
			Ports: []types.Port{
				{PublicPort: 8080, PrivatePort: 80, Type: "tcp"},
			},
			Labels: map[string]string{
				"com.docker.compose.service": "web",
			},
		},
		{
			ID:    "def456abc123",
			Name:  "db",
			Image: "postgres:15",
			Env:   []string{"POSTGRES_USER=admin", "POSTGRES_DB=app"},
			Ports: []types.Port{
				{PublicPort: 5432, PrivatePort: 5432, Type: "tcp"},
			},
			Labels: map[string]string{
				"com.docker.compose.service": "db",
			},
		},
	}
}
