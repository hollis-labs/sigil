package renderer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chrispian/sigil/internal/config"
)

func TestLoadProjectConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "sigil.yaml")
	os.WriteFile(configPath, []byte(`
version: "1.0"
name: test-project
defaults:
  theme: default
  renderer: go-templ
  output: internal/ui/
`), 0644)

	cfg, err := loadProjectConfig(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "test-project" {
		t.Errorf("expected name 'test-project', got %q", cfg.Name)
	}
	if cfg.Defaults.Theme != "default" {
		t.Errorf("expected default theme 'default', got %q", cfg.Defaults.Theme)
	}
}

func TestLoadTheme(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "default.yaml"), []byte(`
name: default
description: "Default dark theme"
tokens:
  colors:
    background: "9 9 11"
    text: "244 244 245"
  radius:
    md: "0.375rem"
`), 0644)

	theme, err := LoadTheme(dir, "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if theme.Name != "default" {
		t.Errorf("expected name 'default', got %q", theme.Name)
	}
	if theme.Tokens["colors"]["background"] != "9 9 11" {
		t.Errorf("expected background '9 9 11', got %v", theme.Tokens["colors"]["background"])
	}
}

func TestLoadThemeInheritance(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "base.yaml"), []byte(`
name: base
tokens:
  colors:
    background: "0 0 0"
    text: "255 255 255"
    accent: "16 185 129"
`), 0644)
	os.WriteFile(filepath.Join(dir, "child.yaml"), []byte(`
name: child
extends: base
tokens:
  colors:
    accent: "99 102 241"
`), 0644)

	theme, err := LoadTheme(dir, "child")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should inherit base colors
	if theme.Tokens["colors"]["background"] != "0 0 0" {
		t.Errorf("expected inherited background, got %v", theme.Tokens["colors"]["background"])
	}
	// Should have child's override
	if theme.Tokens["colors"]["accent"] != "99 102 241" {
		t.Errorf("expected overridden accent, got %v", theme.Tokens["colors"]["accent"])
	}
}

func TestLoadDataSources(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "sprint.yaml"), []byte(`
alias: Sprint
description: "Sprint management"
capabilities:
  - search
  - create
  - delete
fields:
  - name: id
    type: integer
    primary: true
  - name: title
    type: string
    required: true
endpoints:
  list: "GET /api/sprints"
  create: "POST /api/sprints"
`), 0644)

	ds, err := loadDataSources(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ds) != 1 {
		t.Fatalf("expected 1 datasource, got %d", len(ds))
	}
	sprint := ds["Sprint"]
	if sprint == nil {
		t.Fatal("expected Sprint datasource")
	}
	if !sprint.HasCapability("search") {
		t.Error("expected search capability")
	}
	if sprint.HasCapability("update") {
		t.Error("unexpected update capability")
	}
	if len(sprint.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(sprint.Fields))
	}
}

func TestLoadPages(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test-page.yaml"), []byte(`
sigil: "1.0"
id: test-page
title: Test Page
overlay: page
layout:
  id: root
  type: rows
`), 0644)
	os.WriteFile(filepath.Join(dir, "other-page.yaml"), []byte(`
sigil: "1.0"
id: other-page
title: Other Page
overlay: modal
layout:
  id: root
  type: columns
`), 0644)

	// Load all
	pages, err := loadPages(dir, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}

	// Load with filter
	pages, err = loadPages(dir, []string{"test-page"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 1 {
		t.Fatalf("expected 1 filtered page, got %d", len(pages))
	}
	if pages[0].ID != "test-page" {
		t.Errorf("expected test-page, got %q", pages[0].ID)
	}
}

func TestCollectUsedTypes(t *testing.T) {
	layout := &config.Component{
		Type: "rows",
		Children: []config.Component{
			{Type: "columns", Children: []config.Component{
				{Type: "heading"},
				{Type: "button"},
			}},
			{Type: "search-bar"},
			{Type: "data-table"},
		},
	}
	types := map[string]bool{}
	collectUsedTypes(layout, types)

	expected := []string{"rows", "columns", "heading", "button", "search-bar", "data-table"}
	for _, e := range expected {
		if !types[e] {
			t.Errorf("expected type %q in used types", e)
		}
	}
	if len(types) != 6 {
		t.Errorf("expected 6 types, got %d", len(types))
	}
}

func TestRendererRegistry(t *testing.T) {
	// Create a mock renderer
	mock := &mockRenderer{name: "test-renderer"}
	RegisterRenderer(mock)
	defer func() {
		delete(renderers, "test-renderer")
	}()

	r, ok := GetRenderer("test-renderer")
	if !ok {
		t.Fatal("expected to find test-renderer")
	}
	if r.Name() != "test-renderer" {
		t.Errorf("expected name 'test-renderer', got %q", r.Name())
	}

	_, ok = GetRenderer("nonexistent")
	if ok {
		t.Error("expected not to find nonexistent renderer")
	}

	names := AvailableRenderers()
	found := false
	for _, n := range names {
		if n == "test-renderer" {
			found = true
		}
	}
	if !found {
		t.Error("expected test-renderer in available renderers")
	}
}

func TestGenerateUnknownRenderer(t *testing.T) {
	cfg := GenerateConfig{
		Target:    "nonexistent-renderer",
		OutputDir: t.TempDir(),
	}
	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error for unknown renderer")
	}
}

func TestHasCapability(t *testing.T) {
	ds := &DataSourceManifest{
		Capabilities: []string{"search", "create", "delete"},
	}
	if !ds.HasCapability("search") {
		t.Error("expected search capability")
	}
	if ds.HasCapability("update") {
		t.Error("unexpected update capability")
	}
}

// Helpers

type mockRenderer struct {
	name string
}

func (m *mockRenderer) Name() string { return m.name }
func (m *mockRenderer) Render(ctx *RenderContext) ([]OutputFile, error) {
	return []OutputFile{{Path: "test.txt", Content: []byte("test")}}, nil
}
func (m *mockRenderer) RenderTheme(theme *ThemeConfig) ([]OutputFile, error) {
	return nil, nil
}
func (m *mockRenderer) RenderDataSourceStubs(ds *DataSourceManifest) ([]OutputFile, error) {
	return nil, nil
}
func (m *mockRenderer) SharedComponents(usedTypes []string) ([]OutputFile, error) {
	return nil, nil
}
