package drift

import (
	"fmt"
	"strings"

	"github.com/driftcheck/internal/docker"
	"github.com/driftcheck/internal/manifest"
)

// K8sResult holds the drift result for a single Kubernetes container spec.
type K8sResult struct {
	DeploymentName string
	ContainerName  string
	ExpectedImage  string
	ActualImage    string
	ContainerID    string
	Drifted        bool
	Reason         string
}

// DetectK8sDrift compares a parsed Kubernetes Deployment manifest against
// the list of running containers and returns per-container drift results.
func DetectK8sDrift(m *manifest.K8sManifest, running []docker.ContainerInfo) []K8sResult {
	var results []K8sResult

	for _, specContainer := range m.Spec.Template.Spec.Containers {
		result := K8sResult{
			DeploymentName: m.Metadata.Name,
			ContainerName:  specContainer.Name,
			ExpectedImage:  specContainer.Image,
		}

		match, found := findK8sContainer(specContainer.Name, m.Metadata.Name, running)
		if !found {
			result.Drifted = true
			result.Reason = fmt.Sprintf("no running container found for deployment %q container %q",
				m.Metadata.Name, specContainer.Name)
			results = append(results, result)
			continue
		}

		result.ContainerID = match.ID
		result.ActualImage = match.Image

		if !imagesMatch(specContainer.Image, match.Image) {
			result.Drifted = true
			result.Reason = fmt.Sprintf("image mismatch: manifest=%q running=%q",
				specContainer.Image, match.Image)
		}

		results = append(results, result)
	}

	return results
}

// findK8sContainer looks for a running container whose name contains both
// the deployment name and the container name (common k8s naming convention).
func findK8sContainer(containerName, deploymentName string, running []docker.ContainerInfo) (docker.ContainerInfo, bool) {
	for _, c := range running {
		lower := strings.ToLower(c.Name)
		if strings.Contains(lower, strings.ToLower(deploymentName)) &&
			strings.Contains(lower, strings.ToLower(containerName)) {
			return c, true
		}
	}
	return docker.ContainerInfo{}, false
}

// imagesMatch returns true when two image references are considered equivalent.
func imagesMatch(expected, actual string) bool {
	return normaliseImage(expected) == normaliseImage(actual)
}

// normaliseImage strips the docker.io/library/ prefix for comparison purposes.
func normaliseImage(image string) string {
	image = strings.TrimPrefix(image, "docker.io/library/")
	image = strings.TrimPrefix(image, "docker.io/")
	return image
}
