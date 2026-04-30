package drift

import (
	"fmt"

	"github.com/user/driftcheck/internal/docker"
	"github.com/user/driftcheck/internal/manifest"
)

// DriftResult holds the result of a drift check for a single service.
type DriftResult struct {
	ServiceName string
	ContainerID string
	Drifted     bool
	Issues      []string
}

// Detector compares running containers against compose manifests.
type Detector struct {
	client docker.Client
}

// NewDetector creates a new Detector with the provided docker client.
func NewDetector(client docker.Client) *Detector {
	return &Detector{client: client}
}

// Detect compares containers listed by the docker client against the
// services defined in the given ComposeFile and returns drift results.
func (d *Detector) Detect(compose *manifest.ComposeFile) ([]DriftResult, error) {
	containers, err := d.client.ListContainers()
	if err != nil {
		return nil, fmt.Errorf("listing containers: %w", err)
	}

	// Build a lookup map: service name -> container info
	running := make(map[string]docker.ContainerInfo)
	for _, c := range containers {
		if name, ok := c.Labels["com.docker.compose.service"]; ok {
			running[name] = c
		}
	}

	var results []DriftResult
	for svcName, svc := range compose.Services {
		result := DriftResult{ServiceName: svcName}

		c, found := running[svcName]
		if !found {
			result.Drifted = true
			result.Issues = append(result.Issues, "service not running")
			results = append(results, result)
			continue
		}

		result.ContainerID = c.ID
		checkImage(&result, svc.Image, c.Image)
		results = append(results, result)
	}

	return results, nil
}

// checkImage appends an issue if the running image differs from the manifest.
func checkImage(r *DriftResult, expected, actual string) {
	if expected == "" {
		return
	}
	if expected != actual {
		r.Drifted = true
		r.Issues = append(r.Issues,
			fmt.Sprintf("image mismatch: want %q, got %q", expected, actual))
	}
}
