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
`), 0600)

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
`), 0600)

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
`), 0600)
	os.WriteFile(filepath.Join(dir, "child.yaml"), []byte(`
name: child
extends: base
tokens:
  colors:
    accent: "99 102 241"
`), 0600)

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
`), 0600)

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
`), 0600)
	os.WriteFile(filepath.Join(dir, "other-page.yaml"), []byte(`
sigil: "1.0"
id: other-page
title: Other Page
overlay: modal
layout:
  id: root
  type: columns
`), 0600)

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

// TestGenerateTargetModeOverride drives the engine through a synthesized
// .sigil/ tree and asserts that GenerateConfig.TargetMode wins over
// app.yaml's target_mode, with unset falling through to app.yaml.
func TestGenerateTargetModeOverride(t *testing.T) {
	cases := []struct {
		name           string
		appYAMLMode    string
		cliOverride    string
		wantTargetMode string
	}{
		{"unset everywhere defaults to app-router", "", "", "app-router"},
		{"app.yaml spa, no override", "spa", "", "spa"},
		{"app.yaml app-router, override spa", "app-router", "spa", "spa"},
		{"app.yaml spa, override app-router", "spa", "app-router", "app-router"},
		{"app.yaml unset, override spa", "", "spa", "spa"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Register a capturing mock renderer.
			captured := &captureRenderer{name: "capture-render-target"}
			RegisterRenderer(captured)
			defer delete(renderers, captured.name)

			dir := t.TempDir()
			sigilDir := filepath.Join(dir, ".sigil")
			os.MkdirAll(filepath.Join(sigilDir, "pages"), 0750)
			os.MkdirAll(filepath.Join(sigilDir, "themes"), 0750)
			os.MkdirAll(filepath.Join(sigilDir, "datasources"), 0750)

			os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: capture-test
defaults:
  theme: default
  renderer: capture-render-target
`), 0600)

			appYAML := "sigil: \"1.0\"\nkind: app\nname: Test\nmodules:\n  - id: x\n    shell: x-shell\n    route_group: \"(x)\"\n    pages: [x-home]\n"
			if tc.appYAMLMode != "" {
				appYAML += "target_mode: " + tc.appYAMLMode + "\n"
			}
			os.WriteFile(filepath.Join(sigilDir, "app.yaml"), []byte(appYAML), 0600)

			os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
tokens:
  colors:
    background: "0 0 0"
`), 0600)

			os.WriteFile(filepath.Join(sigilDir, "pages", "x-home.yaml"), []byte(`
sigil: "1.0"
kind: page
id: x-home
title: Home
overlay: page
module: x
layout:
  id: root
  type: rows
`), 0600)
			os.WriteFile(filepath.Join(sigilDir, "pages", "x-shell.yaml"), []byte(`
sigil: "1.0"
kind: page
id: x-shell
title: Shell
overlay: page
layout:
  id: root
  type: rows
`), 0600)

			cfg := GenerateConfig{
				Target:     captured.name,
				OutputDir:  filepath.Join(dir, "out"),
				SigilDir:   sigilDir,
				DryRun:     true,
				TargetMode: tc.cliOverride,
			}
			if _, err := Generate(cfg); err != nil {
				t.Fatalf("Generate failed: %v", err)
			}
			if captured.lastTargetMode != tc.wantTargetMode {
				t.Errorf("TargetMode in RenderContext = %q, want %q", captured.lastTargetMode, tc.wantTargetMode)
			}
		})
	}
}

// TestGenerateTargetModeInvalidOverride ensures an invalid CLI override is
// rejected by validation (engine path; the CLI cmd has its own up-front
// rejection too).
func TestGenerateTargetModeInvalidOverride(t *testing.T) {
	captured := &captureRenderer{name: "capture-render-invalid"}
	RegisterRenderer(captured)
	defer delete(renderers, captured.name)

	dir := t.TempDir()
	sigilDir := filepath.Join(dir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0750)
	os.MkdirAll(filepath.Join(sigilDir, "themes"), 0750)
	os.MkdirAll(filepath.Join(sigilDir, "datasources"), 0750)

	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: invalid-test
defaults:
  theme: default
  renderer: capture-render-invalid
`), 0600)
	os.WriteFile(filepath.Join(sigilDir, "app.yaml"), []byte(`
sigil: "1.0"
kind: app
name: Invalid
modules:
  - id: x
    shell: x-shell
    route_group: "(x)"
    pages: [x-home]
`), 0600)
	os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
tokens:
  colors:
    background: "0 0 0"
`), 0600)
	os.WriteFile(filepath.Join(sigilDir, "pages", "x-home.yaml"), []byte(`
sigil: "1.0"
kind: page
id: x-home
title: Home
overlay: page
module: x
layout:
  id: root
  type: rows
`), 0600)
	os.WriteFile(filepath.Join(sigilDir, "pages", "x-shell.yaml"), []byte(`
sigil: "1.0"
kind: page
id: x-shell
title: Shell
overlay: page
layout:
  id: root
  type: rows
`), 0600)

	cfg := GenerateConfig{
		Target:     captured.name,
		OutputDir:  filepath.Join(dir, "out"),
		SigilDir:   sigilDir,
		DryRun:     true,
		TargetMode: "wails",
	}
	if _, err := Generate(cfg); err == nil {
		t.Fatal("expected validation failure for invalid target_mode override, got nil")
	}
}

type captureRenderer struct {
	name           string
	lastTargetMode string
}

func (m *captureRenderer) Name() string { return m.name }
func (m *captureRenderer) Render(ctx *RenderContext) ([]OutputFile, error) {
	m.lastTargetMode = ctx.TargetMode
	return nil, nil
}
func (m *captureRenderer) RenderTheme(_ *ThemeConfig) ([]OutputFile, error) {
	return nil, nil
}
func (m *captureRenderer) RenderDataSourceStubs(_ *DataSourceManifest) ([]OutputFile, error) {
	return nil, nil
}
func (m *captureRenderer) SharedComponents(_ []string) ([]OutputFile, error) {
	return nil, nil
}
func (m *captureRenderer) RenderLayout(_ *LayoutContext) ([]OutputFile, error) {
	return nil, nil
}
func (m *captureRenderer) RenderAPIClient(_ *config.AppConfig) ([]OutputFile, error) {
	return nil, nil
}
func (m *captureRenderer) RenderProviders(_ *config.AppConfig, _ string) ([]OutputFile, error) {
	return nil, nil
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
func (m *mockRenderer) RenderLayout(_ *LayoutContext) ([]OutputFile, error) {
	return nil, nil
}
func (m *mockRenderer) RenderAPIClient(_ *config.AppConfig) ([]OutputFile, error) {
	return nil, nil
}
func (m *mockRenderer) RenderProviders(_ *config.AppConfig, _ string) ([]OutputFile, error) {
	return nil, nil
}
