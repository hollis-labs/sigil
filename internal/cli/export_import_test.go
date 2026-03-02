package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportImportRoundTrip(t *testing.T) {
	// Create a temp sigil project
	tmpDir := t.TempDir()
	sigilDir := filepath.Join(tmpDir, ".sigil")
	pagesDir := filepath.Join(sigilDir, "pages")
	os.MkdirAll(pagesDir, 0755)

	// Write a test page YAML
	pageYAML := `sigil: "1.0"
kind: page
id: test-page
title: Test Page
overlay: page
module: test
layout:
  id: root
  type: rows
  props:
    gap: 4
  children:
    - id: header
      type: heading
      props:
        level: 2
        text: Hello
    - id: btn
      type: button
      props:
        label: Click
        variant: primary
shortcuts:
  - key: Escape
    action:
      type: close
`
	os.WriteFile(filepath.Join(pagesDir, "test-page.yaml"), []byte(pageYAML), 0644)

	// Export to JSON
	exportDir := filepath.Join(tmpDir, "exported")
	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"export", "--sigil-dir", sigilDir, "--output", exportDir})
	var exportOut bytes.Buffer
	rootCmd.SetOut(&exportOut)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(exportOut.String(), "Exported 1 page") {
		t.Errorf("expected 'Exported 1 page' in output, got: %s", exportOut.String())
	}

	// Verify JSON file exists and is valid
	jsonPath := filepath.Join(exportDir, "test-page.json")
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("reading exported JSON: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["id"] != "test-page" {
		t.Errorf("expected id 'test-page', got %v", parsed["id"])
	}
	if parsed["title"] != "Test Page" {
		t.Errorf("expected title 'Test Page', got %v", parsed["title"])
	}

	// Import back into a new sigil dir
	newSigilDir := filepath.Join(tmpDir, ".sigil2")
	os.MkdirAll(filepath.Join(newSigilDir, "pages"), 0755)

	rootCmd2 := NewRootCmd()
	rootCmd2.SetArgs([]string{"import", "--from", "json", "--input", exportDir, "--sigil-dir", newSigilDir})
	var importOut bytes.Buffer
	rootCmd2.SetOut(&importOut)
	if err := rootCmd2.Execute(); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	if !strings.Contains(importOut.String(), "Imported 1 page") {
		t.Errorf("expected 'Imported 1 page' in output, got: %s", importOut.String())
	}

	// Verify YAML was written
	yamlPath := filepath.Join(newSigilDir, "pages", "test-page.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		t.Fatalf("imported YAML file not found: %v", err)
	}

	// Re-export from the imported dir and compare
	reexportDir := filepath.Join(tmpDir, "reexported")
	rootCmd3 := NewRootCmd()
	rootCmd3.SetArgs([]string{"export", "--sigil-dir", newSigilDir, "--output", reexportDir})
	rootCmd3.SetOut(&bytes.Buffer{})
	if err := rootCmd3.Execute(); err != nil {
		t.Fatalf("re-export failed: %v", err)
	}

	reexportJSON, err := os.ReadFile(filepath.Join(reexportDir, "test-page.json"))
	if err != nil {
		t.Fatalf("reading re-exported JSON: %v", err)
	}

	// Parse both JSON objects and compare key fields
	var orig, reimported map[string]interface{}
	json.Unmarshal(jsonData, &orig)
	json.Unmarshal(reexportJSON, &reimported)

	if orig["id"] != reimported["id"] {
		t.Errorf("round-trip ID mismatch: %v vs %v", orig["id"], reimported["id"])
	}
	if orig["title"] != reimported["title"] {
		t.Errorf("round-trip title mismatch: %v vs %v", orig["title"], reimported["title"])
	}
	if orig["overlay"] != reimported["overlay"] {
		t.Errorf("round-trip overlay mismatch: %v vs %v", orig["overlay"], reimported["overlay"])
	}
}

func TestExportFilterPages(t *testing.T) {
	tmpDir := t.TempDir()
	sigilDir := filepath.Join(tmpDir, ".sigil")
	pagesDir := filepath.Join(sigilDir, "pages")
	os.MkdirAll(pagesDir, 0755)

	page1 := `sigil: "1.0"
kind: page
id: page-one
title: Page One
overlay: page
layout:
  type: rows
  children:
    - type: heading
      props:
        text: One
`
	page2 := `sigil: "1.0"
kind: page
id: page-two
title: Page Two
overlay: page
layout:
  type: rows
  children:
    - type: heading
      props:
        text: Two
`
	os.WriteFile(filepath.Join(pagesDir, "page-one.yaml"), []byte(page1), 0644)
	os.WriteFile(filepath.Join(pagesDir, "page-two.yaml"), []byte(page2), 0644)

	exportDir := filepath.Join(tmpDir, "exported")
	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"export", "--sigil-dir", sigilDir, "--output", exportDir, "--pages", "page-one"})
	rootCmd.SetOut(&bytes.Buffer{})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Only page-one should be exported
	if _, err := os.Stat(filepath.Join(exportDir, "page-one.json")); err != nil {
		t.Error("expected page-one.json to exist")
	}
	if _, err := os.Stat(filepath.Join(exportDir, "page-two.json")); err == nil {
		t.Error("expected page-two.json NOT to exist")
	}
}

func TestImportInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	sigilDir := filepath.Join(tmpDir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0755)

	// Write invalid JSON
	inputDir := filepath.Join(tmpDir, "input")
	os.MkdirAll(inputDir, 0755)
	os.WriteFile(filepath.Join(inputDir, "bad.json"), []byte("{invalid json"), 0644)

	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"import", "--from", "json", "--input", inputDir, "--sigil-dir", sigilDir})
	rootCmd.SetOut(&bytes.Buffer{})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for invalid JSON import")
	}
}

func TestImportDryRun(t *testing.T) {
	tmpDir := t.TempDir()
	sigilDir := filepath.Join(tmpDir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0755)

	inputDir := filepath.Join(tmpDir, "input")
	os.MkdirAll(inputDir, 0755)

	pageJSON := `{
  "sigil": "1.0",
  "kind": "page",
  "id": "dry-test",
  "title": "Dry Run Test",
  "overlay": "page",
  "layout": {
    "type": "rows",
    "children": [{"type": "heading", "props": {"text": "Hi", "level": 2}}]
  }
}`
	os.WriteFile(filepath.Join(inputDir, "dry-test.json"), []byte(pageJSON), 0644)

	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"import", "--from", "json", "--input", inputDir, "--sigil-dir", sigilDir, "--dry-run"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("dry-run import failed: %v", err)
	}

	if !strings.Contains(out.String(), "VALID") {
		t.Error("expected VALID in dry-run output")
	}
	if !strings.Contains(out.String(), "Validated 1 page") {
		t.Error("expected 'Validated 1 page' in output")
	}

	// Verify no file was written
	if _, err := os.Stat(filepath.Join(sigilDir, "pages", "dry-test.yaml")); err == nil {
		t.Error("expected no YAML file in dry-run mode")
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"export", "--format", "xml"})
	rootCmd.SetOut(&bytes.Buffer{})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestImportUnsupportedFormat(t *testing.T) {
	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"import", "--from", "xml", "--input", "."})
	rootCmd.SetOut(&bytes.Buffer{})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
