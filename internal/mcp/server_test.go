package mcp

import (
	"context"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer("/tmp/example/.sigil")
	if got := s.SigilDir(); got != "/tmp/example/.sigil" {
		t.Errorf("SigilDir() = %q, want %q", got, "/tmp/example/.sigil")
	}
}

func TestRegisterAllToolsAnnotations(t *testing.T) {
	s := NewServer(".")
	RegisterAllTools(s)

	defs := s.ToolDefinitions()
	if len(defs) != 9 {
		t.Fatalf("expected 9 tools, got %d", len(defs))
	}

	wantReadOnly := map[string]bool{
		"sigil_list_pages":           true,
		"sigil_get_page":             true,
		"sigil_create_page":          false,
		"sigil_update_page":          false,
		"sigil_validate":             true,
		"sigil_list_components":      true,
		"sigil_get_component_schema": true,
		"sigil_list_datasources":     true,
		"sigil_create_datasource":    false,
	}

	seen := map[string]bool{}
	for _, d := range defs {
		seen[d.Name] = true
		want, ok := wantReadOnly[d.Name]
		if !ok {
			t.Errorf("unexpected tool %q", d.Name)
			continue
		}
		if d.Annotations.ReadOnlyHint != want {
			t.Errorf("%s: ReadOnlyHint = %v, want %v", d.Name, d.Annotations.ReadOnlyHint, want)
		}
		if !d.Annotations.IdempotentHint {
			t.Errorf("%s: expected IdempotentHint true", d.Name)
		}
		if d.InputSchema == nil {
			t.Errorf("%s: expected a non-nil input schema", d.Name)
		}
	}
	for name := range wantReadOnly {
		if !seen[name] {
			t.Errorf("missing tool %q", name)
		}
	}
}

func TestRegisterAllToolsUpdatePageIsDestructive(t *testing.T) {
	s := NewServer(".")
	RegisterAllTools(s)

	for _, d := range s.ToolDefinitions() {
		if d.Name != "sigil_update_page" {
			continue
		}
		if !d.Annotations.DestructiveHint {
			t.Error("sigil_update_page: expected DestructiveHint true (it overwrites file content)")
		}
		return
	}
	t.Fatal("sigil_update_page not registered")
}

func TestCallToolDirect(t *testing.T) {
	s := NewServer(".")
	RegisterAllTools(s)

	if _, err := s.CallTool(context.Background(), "sigil_list_components", map[string]any{}); err != nil {
		t.Fatalf("CallTool error: %v", err)
	}
}

func TestCallToolUnknown(t *testing.T) {
	s := NewServer(".")
	if _, err := s.CallTool(context.Background(), "nonexistent", map[string]any{}); err == nil {
		t.Fatal("expected error for unknown tool")
	}
}
