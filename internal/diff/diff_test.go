package diff_test

import (
	"strings"
	"testing"

	"github.com/yourorg/driftcheck/internal/diff"
)

func TestField_String(t *testing.T) {
	f := diff.Field{Name: "image", Expected: "nginx:1.25", Actual: "nginx:1.24"}
	got := f.String()
	if !strings.Contains(got, "image") || !strings.Contains(got, "nginx:1.25") || !strings.Contains(got, "nginx:1.24") {
		t.Errorf("unexpected Field.String(): %q", got)
	}
}

func TestResult_HasDrift_False(t *testing.T) {
	r := diff.Result{Service: "web"}
	if r.HasDrift() {
		t.Error("expected no drift for empty Fields slice")
	}
}

func TestResult_HasDrift_True(t *testing.T) {
	r := diff.Result{
		Service: "web",
		Fields:  []diff.Field{{Name: "image", Expected: "a", Actual: "b"}},
	}
	if !r.HasDrift() {
		t.Error("expected drift to be detected")
	}
}

func TestResult_Summary_NoDrift(t *testing.T) {
	r := diff.Result{Service: "db"}
	s := r.Summary()
	if !strings.Contains(s, "no drift") {
		t.Errorf("expected 'no drift' in summary, got: %q", s)
	}
}

func TestResult_Summary_WithDrift(t *testing.T) {
	r := diff.Result{
		Service: "api",
		Fields: []diff.Field{
			{Name: "image", Expected: "app:v2", Actual: "app:v1"},
			{Name: "env.PORT", Expected: "8080", Actual: "9090"},
		},
	}
	s := r.Summary()
	if !strings.Contains(s, "2 field(s)") {
		t.Errorf("expected field count in summary, got: %q", s)
	}
	if !strings.Contains(s, "image") || !strings.Contains(s, "env.PORT") {
		t.Errorf("expected field names in summary, got: %q", s)
	}
}

func TestCompare_NoDiff(t *testing.T) {
	expected := map[string]string{"image": "nginx:1.25", "restart": "always"}
	actual := map[string]string{"image": "nginx:1.25", "restart": "always"}
	r := diff.Compare("web", expected, actual)
	if r.HasDrift() {
		t.Errorf("expected no drift, got: %+v", r.Fields)
	}
}

func TestCompare_MissingActualField(t *testing.T) {
	expected := map[string]string{"image": "nginx:1.25"}
	actual := map[string]string{}
	r := diff.Compare("web", expected, actual)
	if !r.HasDrift() {
		t.Error("expected drift when actual field is missing")
	}
	if r.Fields[0].Actual != "<missing>" {
		t.Errorf("expected '<missing>', got %q", r.Fields[0].Actual)
	}
}

func TestCompare_ValueMismatch(t *testing.T) {
	expected := map[string]string{"image": "redis:7", "port": "6379"}
	actual := map[string]string{"image": "redis:6", "port": "6379"}
	r := diff.Compare("cache", expected, actual)
	if len(r.Fields) != 1 {
		t.Errorf("expected 1 differing field, got %d", len(r.Fields))
	}
	if r.Fields[0].Name != "image" {
		t.Errorf("expected field name 'image', got %q", r.Fields[0].Name)
	}
}
