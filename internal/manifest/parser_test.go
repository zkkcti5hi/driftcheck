package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/driftcheck/internal/manifest"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestParseComposeFile_Valid(t *testing.T) {
	content := `
version: "3.8"
services:
  web:
    image: nginx:1.25
    ports:
      - "80:80"
    environment:
      ENV: production
    volumes:
      - ./html:/usr/share/nginx/html
`
	path := writeTempFile(t, content)

	services, err := manifest.ParseComposeFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(services))
	}

	svc := services[0]
	if svc.Name != "web" {
		t.Errorf("expected name %q, got %q", "web", svc.Name)
	}
	if svc.Image != "nginx:1.25" {
		t.Errorf("expected image %q, got %q", "nginx:1.25", svc.Image)
	}
	if svc.Environment["ENV"] != "production" {
		t.Errorf("expected ENV=production, got %q", svc.Environment["ENV"])
	}
}

func TestParseComposeFile_MissingFile(t *testing.T) {
	_, err := manifest.ParseComposeFile("/nonexistent/path/docker-compose.yml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParseComposeFile_EmptyServices(t *testing.T) {
	content := `version: "3.8"
services: {}
`
	path := writeTempFile(t, content)
	_, err := manifest.ParseComposeFile(path)
	if err == nil {
		t.Fatal("expected error for empty services, got nil")
	}
}
