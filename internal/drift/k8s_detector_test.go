package drift_test

import (
	"testing"

	"github.com/yourorg/driftcheck/internal/drift"
	"github.com/yourorg/driftcheck/internal/manifest"
)

func k8sDeployment(name, image string) manifest.K8sService {
	return manifest.K8sService{
		Name:  name,
		Image: image,
	}
}

func k8sContainer(name, image, id string) drift.ContainerInfo {
	return drift.ContainerInfo{
		Name:  "/" + name,
		Image: image,
		ID:    id,
	}
}

func TestDetectK8sDrift_NoDrift(t *testing.T) {
	services := []manifest.K8sService{
		k8sDeployment("api", "nginx:1.25"),
	}
	containers := []drift.ContainerInfo{
		k8sContainer("api", "nginx:1.25", "abc123"),
	}

	results := drift.DetectK8sDrift(services, containers)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift for matching image")
	}
}

func TestDetectK8sDrift_ImageMismatch(t *testing.T) {
	services := []manifest.K8sService{
		k8sDeployment("api", "nginx:1.25"),
	}
	containers := []drift.ContainerInfo{
		k8sContainer("api", "nginx:1.24", "def456"),
	}

	results := drift.DetectK8sDrift(services, containers)
	if !results[0].Drifted {
		t.Errorf("expected drift for image mismatch")
	}
	if results[0].ExpectedImage != "nginx:1.25" {
		t.Errorf("unexpected ExpectedImage: %s", results[0].ExpectedImage)
	}
	if results[0].ActualImage != "nginx:1.24" {
		t.Errorf("unexpected ActualImage: %s", results[0].ActualImage)
	}
}

func TestDetectK8sDrift_ServiceNotRunning(t *testing.T) {
	services := []manifest.K8sService{
		k8sDeployment("worker", "myapp:latest"),
	}
	containers := []drift.ContainerInfo{}

	results := drift.DetectK8sDrift(services, containers)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Drifted {
		t.Errorf("expected drift when service has no running container")
	}
	if results[0].ContainerID != "" {
		t.Errorf("expected empty ContainerID for missing service")
	}
}

func TestDetectK8sDrift_NormalisedImageMatch(t *testing.T) {
	services := []manifest.K8sService{
		k8sDeployment("cache", "redis"),
	}
	containers := []drift.ContainerInfo{
		k8sContainer("cache", "redis:latest", "ccc789"),
	}

	results := drift.DetectK8sDrift(services, containers)
	if results[0].Drifted {
		t.Errorf("expected no drift when 'redis' matches 'redis:latest'")
	}
}
