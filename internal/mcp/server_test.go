package mcp

import (
	"context"
	"errors"
	"testing"

	gomcpserver "github.com/hollis-labs/libs/mcp/plugin-mcp/go-mcp/server"
	libauth "github.com/hollis-labs/libs/security/auth"
)

// newTestServer wraps a standalone go-mcp server, as the Substrate mcp module's would be.
func newTestServer(sigilDir string) *Server {
	return NewServer(gomcpserver.NewServer("sigil", "test"), sigilDir)
}

func TestNewServer(t *testing.T) {
	s := newTestServer("/tmp/example/.sigil")
	if got := s.SigilDir(); got != "/tmp/example/.sigil" {
		t.Errorf("SigilDir() = %q, want %q", got, "/tmp/example/.sigil")
	}
}

func TestRegisterAllToolsAnnotations(t *testing.T) {
	s := newTestServer(".")
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
	s := newTestServer(".")
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
	s := newTestServer(".")
	RegisterAllTools(s)

	if _, err := s.CallTool(context.Background(), "sigil_list_components", map[string]any{}); err != nil {
		t.Fatalf("CallTool error: %v", err)
	}
}

func TestCallToolUnknown(t *testing.T) {
	s := newTestServer(".")
	if _, err := s.CallTool(context.Background(), "nonexistent", map[string]any{}); err == nil {
		t.Fatal("expected error for unknown tool")
	}
}

// TestEveryToolHasAScope: the tool table is the only registration path and each entry
// carries a scope, mapped onto the auth module's table.
func TestEveryToolHasAScope(t *testing.T) {
	table := libauth.NewScopeTable()
	RegisterScopes(table)
	s := newTestServer(".")
	RegisterAllTools(s)
	for _, d := range s.ToolDefinitions() {
		scope, known := table.ForTool(d.Name)
		if !known || (scope != "sigil.read" && scope != "sigil.write") {
			t.Errorf("tool %s: scope %q (known %v)", d.Name, scope, known)
		}
		if wantWrite := !d.Annotations.ReadOnlyHint; wantWrite != (scope == "sigil.write") {
			t.Errorf("tool %s: scope %q does not match its read-only hint %v", d.Name, scope, d.Annotations.ReadOnlyHint)
		}
	}
}

func TestToolsAreAuthorizedBeforeTheyRun(t *testing.T) {
	s := newTestServer(".")
	RegisterAllTools(s)
	var asked []string
	s.UseAuthorizer(func(_ context.Context, tool string) error {
		asked = append(asked, tool)
		return errors.New("denied")
	})
	if _, err := s.CallTool(context.Background(), "sigil_list_components", map[string]any{}); err == nil {
		t.Fatal("a denied call must fail")
	}
	if len(asked) != 1 || asked[0] != "sigil_list_components" {
		t.Errorf("authorizer saw %v", asked)
	}
}
