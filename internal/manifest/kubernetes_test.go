package manifest

import (
	"os"
	"testing"
)

const validDeploymentYAML = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: default
  labels:
    app: my-app
spec:
  template:
    spec:
      containers:
        - name: app
          image: my-app:1.2.3
          ports:
            - containerPort: 8080
              protocol: TCP
          env:
            - name: ENV
              value: production
`

func writeK8sTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "k8s-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseK8sManifest_Valid(t *testing.T) {
	path := writeK8sTempFile(t, validDeploymentYAML)

	m, err := ParseK8sManifest(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.Metadata.Name != "my-app" {
		t.Errorf("expected name %q, got %q", "my-app", m.Metadata.Name)
	}

	containers := m.Spec.Template.Spec.Containers
	if len(containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(containers))
	}

	if containers[0].Image != "my-app:1.2.3" {
		t.Errorf("expected image %q, got %q", "my-app:1.2.3", containers[0].Image)
	}
}

func TestParseK8sManifest_MissingFile(t *testing.T) {
	_, err := ParseK8sManifest("/nonexistent/path/deploy.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseK8sManifest_UnsupportedKind(t *testing.T) {
	yaml := `apiVersion: v1
kind: Service
metadata:
  name: svc
spec:
  template:
    spec:
      containers: []
`
	path := writeK8sTempFile(t, yaml)
	_, err := ParseK8sManifest(path)
	if err == nil {
		t.Fatal("expected error for unsupported kind, got nil")
	}
}

func TestParseK8sManifest_NoContainers(t *testing.T) {
	yaml := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: empty
spec:
  template:
    spec:
      containers: []
`
	path := writeK8sTempFile(t, yaml)
	_, err := ParseK8sManifest(path)
	if err == nil {
		t.Fatal("expected error for no containers, got nil")
	}
}
