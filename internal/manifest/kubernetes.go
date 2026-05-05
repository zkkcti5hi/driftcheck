package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// K8sManifest represents a minimal Kubernetes Deployment manifest.
type K8sManifest struct {
	APIVersion string  `yaml:"apiVersion"`
	Kind       string  `yaml:"kind"`
	Metadata   K8sMeta `yaml:"metadata"`
	Spec       K8sSpec `yaml:"spec"`
}

// K8sMeta holds the metadata section of a Kubernetes resource.
type K8sMeta struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}

// K8sSpec holds the spec section of a Kubernetes Deployment.
type K8sSpec struct {
	Template K8sPodTemplate `yaml:"template"`
}

// K8sPodTemplate holds the pod template spec.
type K8sPodTemplate struct {
	Spec K8sPodSpec `yaml:"spec"`
}

// K8sPodSpec holds the list of containers in a pod.
type K8sPodSpec struct {
	Containers []K8sContainer `yaml:"containers"`
}

// K8sContainer represents a container definition in a Kubernetes manifest.
type K8sContainer struct {
	Name  string            `yaml:"name"`
	Image string            `yaml:"image"`
	Env   []K8sEnvVar       `yaml:"env"`
	Ports []K8sContainerPort `yaml:"ports"`
}

// K8sEnvVar is a name/value environment variable pair.
type K8sEnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// K8sContainerPort exposes a port from a container.
type K8sContainerPort struct {
	ContainerPort int    `yaml:"containerPort"`
	Protocol      string `yaml:"protocol"`
}

// ParseK8sManifest reads and parses a Kubernetes Deployment YAML file.
// Only Deployment kind is supported; other kinds return an error.
func ParseK8sManifest(path string) (*K8sManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("kubernetes: read file %q: %w", path, err)
	}

	var m K8sManifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("kubernetes: parse yaml %q: %w", path, err)
	}

	if m.Kind != "Deployment" {
		return nil, fmt.Errorf("kubernetes: unsupported kind %q (only Deployment is supported)", m.Kind)
	}

	if len(m.Spec.Template.Spec.Containers) == 0 {
		return nil, fmt.Errorf("kubernetes: manifest %q defines no containers", path)
	}

	return &m, nil
}
