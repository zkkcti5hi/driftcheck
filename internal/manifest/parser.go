package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceConfig represents a normalized service configuration
// extracted from either a docker-compose or Kubernetes manifest.
type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
}

// ComposeFile represents the top-level structure of a docker-compose file.
type ComposeFile struct {
	Version  string                     `yaml:"version"`
	Services map[string]ComposeService  `yaml:"services"`
}

// ComposeService represents a single service in a docker-compose file.
type ComposeService struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Volumes     []string          `yaml:"volumes"`
}

// ParseComposeFile reads and parses a docker-compose YAML file,
// returning a slice of normalized ServiceConfig entries.
func ParseComposeFile(path string) ([]ServiceConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading compose file %q: %w", path, err)
	}

	var compose ComposeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, fmt.Errorf("parsing compose file %q: %w", path, err)
	}

	if len(compose.Services) == 0 {
		return nil, fmt.Errorf("no services found in %q", path)
	}

	services := make([]ServiceConfig, 0, len(compose.Services))
	for name, svc := range compose.Services {
		services = append(services, ServiceConfig{
			Name:        name,
			Image:       svc.Image,
			Environment: svc.Environment,
			Ports:       svc.Ports,
			Volumes:     svc.Volumes,
		})
	}
	return services, nil
}
