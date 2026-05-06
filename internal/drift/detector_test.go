package drift_test

import (
	"errors"
	"testing"

	"github.com/user/driftcheck/internal/docker"
	"github.com/user/driftcheck/internal/drift"
	"github.com/user/driftcheck/internal/manifest"
)

func composeSingleService(name, image string) *manifest.ComposeFile {
	return &manifest.ComposeFile{
		Services: map[string]manifest.Service{
			name: {Image: image},
		},
	}
}

func runningContainer(id, image, svcName string) docker.ContainerInfo {
	return docker.ContainerInfo{
		ID:    id,
		Image: image,
		Labels: map[string]string{
			"com.docker.compose.service": svcName,
		},
	}
}

func TestDetect_NoDrift(t *testing.T) {
	client := &docker.MockClient{
		Containers: []docker.ContainerInfo{
			runningContainer("abc123", "nginx:1.25", "web"),
		},
	}
	detector := drift.NewDetector(client)
	results, err := detector.Detect(composeSingleService("web", "nginx:1.25"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift, got issues: %v", results[0].Issues)
	}
}

func TestDetect_ImageMismatch(t *testing.T) {
	client := &docker.MockClient{
		Containers: []docker.ContainerInfo{
			runningContainer("abc123", "nginx:1.24", "web"),
		},
	}
	detector := drift.NewDetector(client)
	results, err := detector.Detect(composeSingleService("web", "nginx:1.25"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Drifted {
		t.Error("expected drift to be detected")
	}
	if len(results[0].Issues) == 0 {
		t.Error("expected at least one issue")
	}
}

func TestDetect_ServiceNotRunning(t *testing.T) {
	client := &docker.MockClient{Containers: []docker.ContainerInfo{}}
	detector := drift.NewDetector(client)
	results, err := detector.Detect(composeSingleService("web", "nginx:1.25"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Drifted {
		t.Error("expected drift for missing service")
	}
}

func TestDetect_ClientError(t *testing.T) {
	client := &docker.MockClient{Err: errors.New("daemon unavailable")}
	detector := drift.NewDetector(client)
	_, err := detector.Detect(composeSingleService("web", "nginx:1.25"))
	if err == nil {
		t.Error("expected error from client, got nil")
	}
}

func TestDetect_MultipleServices(t *testing.T) {
	client := &docker.MockClient{
		Containers: []docker.ContainerInfo{
			runningContainer("abc123", "nginx:1.25", "web"),
			runningContainer("def456", "redis:7.0", "cache"),
		},
	}
	compose := &manifest.ComposeFile{
		Services: map[string]manifest.Service{
			"web":   {Image: "nginx:1.25"},
			"cache": {Image: "redis:6.2"}, // intentional mismatch
		},
	}
	detector := drift.NewDetector(client)
	results, err := detector.Detect(compose)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	driftedCount := 0
	for _, r := range results {
		if r.Drifted {
			driftedCount++
		}
	}
	if driftedCount != 1 {
		t.Errorf("expected 1 drifted service, got %d", driftedCount)
	}
}
