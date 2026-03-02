package gotempl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chrispian/sigil/internal/renderer"
)

// TestEndToEndGeneration verifies the complete pipeline:
// init project → add page + datasource + theme → generate → verify output.
func TestEndToEndGeneration(t *testing.T) {
	// Setup temp directory as a project
	root := t.TempDir()
	sigilDir := filepath.Join(root, ".sigil")
	outDir := filepath.Join(root, "out")

	// Create .sigil structure
	for _, dir := range []string{
		filepath.Join(sigilDir, "pages"),
		filepath.Join(sigilDir, "datasources"),
		filepath.Join(sigilDir, "themes"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	// Write project config
	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: test-project
defaults:
  theme: default
  renderer: go-templ
  output: out/
generation:
  go_module: github.com/test/project
`), 0644)

	// Write theme
	os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
description: "Test theme"
tokens:
  colors:
    background: "9 9 11"
    surface: "24 24 27"
    border: "63 63 70"
    text: "244 244 245"
    text-muted: "161 161 170"
    accent: "16 185 129"
    accent-hover: "5 150 105"
    danger: "239 68 68"
  typography:
    font-sans: "ui-sans-serif, system-ui, sans-serif"
    font-mono: "ui-monospace, monospace"
  radius:
    sm: "0.125rem"
    md: "0.375rem"
    lg: "0.5rem"
variants:
  button:
    primary:
      bg: accent
      text: background
      hover: accent-hover
      radius: md
    destructive:
      bg: danger
      text: background
      hover: danger
  badge:
    default:
      bg: surface
      text: text
`), 0644)

	// Write datasource
	os.WriteFile(filepath.Join(sigilDir, "datasources", "sprint.yaml"), []byte(`
alias: Sprint
description: "Sprint management"
capabilities:
  - search
  - filter
  - create
  - update
  - delete
fields:
  - name: id
    type: integer
    primary: true
    hidden: true
  - name: title
    type: string
    required: true
  - name: status
    type: enum
    values: [planning, active, done]
  - name: created_at
    type: datetime
    readonly: true
endpoints:
  list: "GET /api/sprints"
  create: "POST /api/sprints"
  read: "GET /api/sprints/{id}"
  update: "PUT /api/sprints/{id}"
  delete: "DELETE /api/sprints/{id}"
`), 0644)

	// Write page config
	os.WriteFile(filepath.Join(sigilDir, "pages", "dashboard.yaml"), []byte(`
sigil: "1.0"
id: dashboard
title: Dashboard
overlay: page

datasources:
  - alias: Sprint
    capabilities: [search, filter, create]

layout:
  id: root
  type: rows
  props:
    gap: 4
    padding: 6
  children:
    - id: header
      type: columns
      props:
        justify: between
        align: center
      children:
        - id: title
          type: heading
          props:
            level: 2
            text: Dashboard
        - id: btn-new
          type: button
          props:
            label: New Sprint
            variant: primary
            icon: plus
          actions:
            click:
              type: modal
              title: Create Sprint
              fields:
                - name: title
                  label: Title
                  type: text
                  required: true
              submit:
                datasource: Sprint
                method: POST
    - id: search
      type: search-bar
      props:
        placeholder: "Search sprints..."
        datasource: Sprint
        debounce: 300
    - id: table
      type: data-table
      props:
        datasource: Sprint
        columns:
          - { field: title, label: Title, sortable: true }
          - { field: status, label: Status }
`), 0644)

	// Run generation
	cfg := renderer.GenerateConfig{
		Target:    "go-templ",
		OutputDir: outDir,
		SigilDir:  sigilDir,
		DryRun:    false,
	}

	result, err := renderer.Generate(cfg)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	// Verify results
	if result.PageCount != 1 {
		t.Errorf("expected 1 page, got %d", result.PageCount)
	}
	if len(result.Files) == 0 {
		t.Fatal("expected generated files")
	}

	// Check all files were written
	for _, f := range result.Files {
		fullPath := filepath.Join(outDir, f.Path)
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Errorf("expected file %s to exist: %v", f.Path, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("expected non-empty file %s", f.Path)
		}
	}

	// Verify specific files exist
	expectedFiles := []string{
		"pages/dashboard.templ",
		"theme.css",
		"tailwind.sigil.cjs",
		"handlers/sprint_handler.go",
	}
	for _, expected := range expectedFiles {
		fullPath := filepath.Join(outDir, expected)
		if _, err := os.Stat(fullPath); err != nil {
			t.Errorf("expected file %s: %v", expected, err)
		}
	}

	// Verify shared component files for used types
	usedComponents := []string{
		"components/layout_helpers.templ",
		"components/button.templ",
		"components/heading.templ",
		"components/search_bar.templ",
		"components/data_table.templ",
	}
	for _, comp := range usedComponents {
		fullPath := filepath.Join(outDir, comp)
		if _, err := os.Stat(fullPath); err != nil {
			t.Errorf("expected shared component %s: %v", comp, err)
		}
	}

	// Verify page templ content
	pageContent, _ := os.ReadFile(filepath.Join(outDir, "pages/dashboard.templ"))
	pc := string(pageContent)
	if !strings.Contains(pc, "templ Dashboard()") {
		t.Error("expected Dashboard templ function")
	}
	if !strings.Contains(pc, "flex flex-col") {
		t.Error("expected rows layout class")
	}
	if !strings.Contains(pc, "Dashboard") {
		t.Error("expected page title in heading")
	}

	// Verify theme.css content
	themeContent, _ := os.ReadFile(filepath.Join(outDir, "theme.css"))
	tc := string(themeContent)
	if !strings.Contains(tc, ":root {") {
		t.Error("expected :root in theme.css")
	}
	if !strings.Contains(tc, "--sigil-accent") {
		t.Error("expected --sigil-accent token")
	}
	if !strings.Contains(tc, ".sigil-btn-primary") {
		t.Error("expected .sigil-btn-primary class")
	}

	// Verify handler stub content
	handlerContent, _ := os.ReadFile(filepath.Join(outDir, "handlers/sprint_handler.go"))
	hc := string(handlerContent)
	if !strings.Contains(hc, "type SprintHandler struct") {
		t.Error("expected SprintHandler struct")
	}
	if !strings.Contains(hc, "func (h *SprintHandler) List(") {
		t.Error("expected List method")
	}
	if !strings.Contains(hc, "func (h *SprintHandler) Create(") {
		t.Error("expected Create method")
	}
	if !strings.Contains(hc, "RegisterRoutes") {
		t.Error("expected RegisterRoutes method")
	}
}

func TestEndToEndDryRun(t *testing.T) {
	root := t.TempDir()
	sigilDir := filepath.Join(root, ".sigil")
	outDir := filepath.Join(root, "out")

	// Minimal setup
	for _, dir := range []string{
		filepath.Join(sigilDir, "pages"),
		filepath.Join(sigilDir, "datasources"),
		filepath.Join(sigilDir, "themes"),
	} {
		os.MkdirAll(dir, 0755)
	}

	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: test
defaults:
  theme: default
  renderer: go-templ
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
tokens:
  colors:
    text: "255 255 255"
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "pages", "test.yaml"), []byte(`
sigil: "1.0"
id: test-page
title: Test
overlay: page
layout:
  id: root
  type: rows
`), 0644)

	cfg := renderer.GenerateConfig{
		Target:    "go-templ",
		OutputDir: outDir,
		SigilDir:  sigilDir,
		DryRun:    true,
	}

	result, err := renderer.Generate(cfg)
	if err != nil {
		t.Fatalf("dry run failed: %v", err)
	}

	if result.PageCount != 1 {
		t.Errorf("expected 1 page, got %d", result.PageCount)
	}
	if len(result.Files) == 0 {
		t.Error("expected files in dry run result")
	}

	// Output dir should NOT exist (dry run)
	if _, err := os.Stat(outDir); err == nil {
		t.Error("expected output dir to not exist in dry run")
	}
}

func TestEndToEndClean(t *testing.T) {
	root := t.TempDir()
	sigilDir := filepath.Join(root, ".sigil")
	outDir := filepath.Join(root, "out")

	// Create pre-existing output
	os.MkdirAll(filepath.Join(outDir, "old"), 0755)
	os.WriteFile(filepath.Join(outDir, "old", "stale.txt"), []byte("stale"), 0644)

	// Minimal setup
	for _, dir := range []string{
		filepath.Join(sigilDir, "pages"),
		filepath.Join(sigilDir, "datasources"),
		filepath.Join(sigilDir, "themes"),
	} {
		os.MkdirAll(dir, 0755)
	}

	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: test
defaults:
  theme: default
  renderer: go-templ
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
tokens:
  colors:
    text: "255 255 255"
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "pages", "test.yaml"), []byte(`
sigil: "1.0"
id: test
title: Test
overlay: page
layout:
  id: root
  type: rows
`), 0644)

	cfg := renderer.GenerateConfig{
		Target:    "go-templ",
		OutputDir: outDir,
		SigilDir:  sigilDir,
		Clean:     true,
	}

	_, err := renderer.Generate(cfg)
	if err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	// Old file should be gone
	if _, err := os.Stat(filepath.Join(outDir, "old", "stale.txt")); err == nil {
		t.Error("expected stale file to be cleaned")
	}

	// New files should exist
	if _, err := os.Stat(filepath.Join(outDir, "pages", "test.templ")); err != nil {
		t.Error("expected generated page file")
	}
}

func TestEndToEndValidationFailure(t *testing.T) {
	root := t.TempDir()
	sigilDir := filepath.Join(root, ".sigil")

	for _, dir := range []string{
		filepath.Join(sigilDir, "pages"),
		filepath.Join(sigilDir, "datasources"),
		filepath.Join(sigilDir, "themes"),
	} {
		os.MkdirAll(dir, 0755)
	}

	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: test
defaults:
  theme: default
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
tokens:
  colors:
    text: "255"
`), 0644)

	// Invalid page — missing required fields
	os.WriteFile(filepath.Join(sigilDir, "pages", "bad.yaml"), []byte(`
id: ""
title: ""
layout:
  type: ""
`), 0644)

	cfg := renderer.GenerateConfig{
		Target:    "go-templ",
		OutputDir: filepath.Join(root, "out"),
		SigilDir:  sigilDir,
	}

	_, err := renderer.Generate(cfg)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("expected validation failure message, got: %v", err)
	}
}

func TestEndToEndNoPages(t *testing.T) {
	root := t.TempDir()
	sigilDir := filepath.Join(root, ".sigil")

	for _, dir := range []string{
		filepath.Join(sigilDir, "pages"),
		filepath.Join(sigilDir, "datasources"),
		filepath.Join(sigilDir, "themes"),
	} {
		os.MkdirAll(dir, 0755)
	}

	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: test
defaults:
  theme: default
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
tokens:
  colors:
    text: "255"
`), 0644)

	cfg := renderer.GenerateConfig{
		Target:    "go-templ",
		OutputDir: filepath.Join(root, "out"),
		SigilDir:  sigilDir,
	}

	_, err := renderer.Generate(cfg)
	if err == nil {
		t.Fatal("expected error for no pages")
	}
	if !strings.Contains(err.Error(), "no pages found") {
		t.Errorf("expected 'no pages found' error, got: %v", err)
	}
}
