// Package config handles loading and validating driftcheck configuration
// from a YAML file or environment variables.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ManifestType identifies the kind of manifest to compare against.
type ManifestType string

const (
	ManifestCompose    ManifestType = "compose"
	ManifestKubernetes ManifestType = "kubernetes"
)

// Config holds the top-level driftcheck configuration.
type Config struct {
	Manifest ManifestConfig `yaml:"manifest"`
	Output   OutputConfig   `yaml:"output"`
	Docker   DockerConfig   `yaml:"docker"`
}

// ManifestConfig specifies which manifest file to use and its type.
type ManifestConfig struct {
	Path string       `yaml:"path"`
	Type ManifestType `yaml:"type"`
}

// OutputConfig controls how results are presented.
type OutputConfig struct {
	Format string `yaml:"format"` // text | json | table
	Quiet  bool   `yaml:"quiet"`
}

// DockerConfig holds Docker connection settings.
type DockerConfig struct {
	Host string `yaml:"host"` // e.g. unix:///var/run/docker.sock
}

// Load reads a Config from the YAML file at path.
// If path is empty the returned Config contains zero values.
func Load(path string) (*Config, error) {
	if path == "" {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that required fields are consistent.
func (c *Config) Validate() error {
	if c.Manifest.Path != "" {
		switch c.Manifest.Type {
		case ManifestCompose, ManifestKubernetes, "":
			// ok
		default:
			return fmt.Errorf("config: unsupported manifest type %q", c.Manifest.Type)
		}
	}
	return nil
}
