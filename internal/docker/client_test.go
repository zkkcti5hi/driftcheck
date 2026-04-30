package docker

import (
	"context"
	"errors"
	"testing"
)

func TestMockClient_ListContainers_Success(t *testing.T) {
	mc := &MockClient{Containers: SampleContainers()}

	containers, err := mc.ListContainers(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(containers))
	}

	if containers[0].Name != "web" {
		t.Errorf("expected first container name 'web', got %q", containers[0].Name)
	}
	if containers[0].Image != "nginx:1.25" {
		t.Errorf("expected image 'nginx:1.25', got %q", containers[0].Image)
	}
}

func TestMockClient_ListContainers_Error(t *testing.T) {
	expected := errors.New("docker daemon unavailable")
	mc := &MockClient{Err: expected}

	_, err := mc.ListContainers(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() != expected.Error() {
		t.Errorf("expected error %q, got %q", expected, err)
	}
}

func TestContainerInfo_Fields(t *testing.T) {
	containers := SampleContainers()
	db := containers[1]

	if db.ID != "def456abc123" {
		t.Errorf("unexpected ID: %s", db.ID)
	}
	if len(db.Env) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(db.Env))
	}
	if db.Labels["com.docker.compose.service"] != "db" {
		t.Errorf("unexpected compose service label: %v", db.Labels)
	}
}
