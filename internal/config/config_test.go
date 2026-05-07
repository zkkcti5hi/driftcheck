package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftcheck/internal/config"
)

func writeCfgFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "driftcheck.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeCfgFile: %v", err)
	}
	return p
}

func TestLoad_EmptyPath(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := `
manifest:
  path: docker-compose.yml
  type: compose
output:
  format: json
  quiet: false
docker:
  host: unix:///var/run/docker.sock
`
	p := writeCfgFile(t, raw)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Manifest.Path != "docker-compose.yml" {
		t.Errorf("manifest path = %q, want docker-compose.yml", cfg.Manifest.Path)
	}
	if cfg.Manifest.Type != config.ManifestCompose {
		t.Errorf("manifest type = %q, want compose", cfg.Manifest.Type)
	}
	if cfg.Output.Format != "json" {
		t.Errorf("output format = %q, want json", cfg.Output.Format)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/driftcheck.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	p := writeCfgFile(t, ": this is not valid yaml: [")
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestValidate_UnsupportedType(t *testing.T) {
	cfg := &config.Config{
		Manifest: config.ManifestConfig{
			Path: "some-file.yml",
			Type: "helm",
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for unsupported manifest type")
	}
}

func TestValidate_KubernetesType(t *testing.T) {
	cfg := &config.Config{
		Manifest: config.ManifestConfig{
			Path: "deployment.yaml",
			Type: config.ManifestKubernetes,
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
