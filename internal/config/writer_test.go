package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMarshalYAML(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "test",
		Title:   "Test Page",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Props: map[string]interface{}{
				"gap": 4,
			},
			Children: []Component{
				{ID: "heading", Type: "heading", Props: map[string]interface{}{"text": "Hello"}},
			},
		},
	}

	data, err := MarshalYAML(page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	yaml := string(data)
	if yaml == "" {
		t.Fatal("expected non-empty YAML")
	}

	// Verify it's parseable
	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("round-trip parse failed: %v", err)
	}
	if parsed.ID != "test" {
		t.Errorf("round-trip: expected id 'test', got %q", parsed.ID)
	}
	if parsed.Layout.Type != "rows" {
		t.Errorf("round-trip: expected layout type 'rows', got %q", parsed.Layout.Type)
	}
	if len(parsed.Layout.Children) != 1 {
		t.Fatalf("round-trip: expected 1 child, got %d", len(parsed.Layout.Children))
	}
	if parsed.Layout.Children[0].ID != "heading" {
		t.Errorf("round-trip: expected child id 'heading', got %q", parsed.Layout.Children[0].ID)
	}
}

func TestMarshalOmitsEmptyFields(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "minimal",
		Title:   "Minimal",
		Overlay: "page",
		Layout:  Component{ID: "root", Type: "rows"},
	}

	data, err := MarshalYAML(page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	yaml := string(data)
	// These optional fields should not appear
	for _, absent := range []string{"description:", "module:", "theme:", "datasources:", "shortcuts:", "meta:"} {
		if contains(yaml, absent) {
			t.Errorf("expected %q to be omitted from YAML output", absent)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestWriteFile(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "write-test",
		Title:   "Write Test",
		Overlay: "page",
		Layout:  Component{ID: "root", Type: "rows"},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")

	err := WriteFile(path, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	// Verify round-trip
	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("failed to parse written file: %v", err)
	}
	if parsed.ID != "write-test" {
		t.Errorf("expected id 'write-test', got %q", parsed.ID)
	}
}
