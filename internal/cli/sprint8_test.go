package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiffCommand(t *testing.T) {
	dir := t.TempDir()

	// Write two page configs
	pageA := `sigil: "1.0"
kind: page
id: test
title: Old Title
overlay: page
layout:
  id: root
  type: rows
  children:
    - id: h1
      type: heading
      props:
        level: 2
        text: Hello
`
	pageB := `sigil: "1.0"
kind: page
id: test
title: New Title
overlay: page
layout:
  id: root
  type: rows
  children:
    - id: h1
      type: heading
      props:
        level: 2
        text: World
`
	pathA := filepath.Join(dir, "a.yaml")
	pathB := filepath.Join(dir, "b.yaml")
	os.WriteFile(pathA, []byte(pageA), 0644)
	os.WriteFile(pathB, []byte(pageB), 0644)

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"diff", pathA, pathB, "--no-color"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("diff command: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "title") {
		t.Error("expected title change in diff output")
	}
	if !strings.Contains(output, "text") {
		t.Error("expected text prop change in diff output")
	}
}

func TestDiffCommandNoChanges(t *testing.T) {
	dir := t.TempDir()

	page := `sigil: "1.0"
kind: page
id: test
title: Same
overlay: page
layout:
  id: root
  type: rows
`
	path := filepath.Join(dir, "page.yaml")
	os.WriteFile(path, []byte(page), 0644)

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"diff", path, path, "--no-color"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("diff command: %v", err)
	}

	if !strings.Contains(out.String(), "No differences") {
		t.Error("expected 'No differences' message")
	}
}

func TestMigrateCommandDryRun(t *testing.T) {
	dir := t.TempDir()
	sigilDir := filepath.Join(dir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0755)

	// Page without component IDs
	page := `sigil: "1.0"
kind: page
id: test
title: Test
overlay: page
layout:
  id: root
  type: rows
  children:
    - type: heading
      props:
        level: 2
        text: Hello
`
	os.WriteFile(filepath.Join(sigilDir, "pages", "test.yaml"), []byte(page), 0644)

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"migrate", "--sigil-dir", sigilDir, "--dry-run", "--no-color"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("migrate command: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "DRY RUN") {
		t.Error("expected DRY RUN in output")
	}
	if !strings.Contains(output, "auto-generated ID") {
		t.Error("expected auto-generated ID message")
	}
}

func TestMigrateCommandNoChanges(t *testing.T) {
	dir := t.TempDir()
	sigilDir := filepath.Join(dir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0755)

	page := `sigil: "1.0"
kind: page
id: test
title: Test
overlay: page
layout:
  id: root
  type: rows
  children:
    - id: h1
      type: heading
      props:
        level: 2
        text: Hello
`
	os.WriteFile(filepath.Join(sigilDir, "pages", "test.yaml"), []byte(page), 0644)

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"migrate", "--sigil-dir", sigilDir, "--no-color"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("migrate command: %v", err)
	}

	if !strings.Contains(out.String(), "up to date") {
		t.Error("expected 'up to date' message")
	}
}

func TestSchemaExportCommand(t *testing.T) {
	dir := t.TempDir()
	outputDir := filepath.Join(dir, "schemas")

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"schema", "export", "--output", outputDir, "--no-color"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("schema export command: %v", err)
	}

	// Verify schema file exists
	schemaPath := filepath.Join(outputDir, "page.schema.json")
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("reading schema: %v", err)
	}

	// Verify it's valid JSON
	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Check key fields
	if schema["$schema"] != "https://json-schema.org/draft/2020-12/schema" {
		t.Error("expected JSON Schema draft 2020-12")
	}
	if schema["title"] != "Sigil Page Config" {
		t.Error("expected title")
	}

	// Check $defs
	defs, ok := schema["$defs"].(map[string]interface{})
	if !ok {
		t.Fatal("expected $defs")
	}
	if _, ok := defs["component"]; !ok {
		t.Error("expected component definition")
	}
	if _, ok := defs["action"]; !ok {
		t.Error("expected action definition")
	}

	// Check component type enum includes known types
	compDef := defs["component"].(map[string]interface{})
	props := compDef["properties"].(map[string]interface{})
	typeProp := props["type"].(map[string]interface{})
	typeEnum := typeProp["enum"].([]interface{})
	if len(typeEnum) < 10 {
		t.Errorf("expected at least 10 component types, got %d", len(typeEnum))
	}
}

func TestDoctorCommand(t *testing.T) {
	dir := t.TempDir()
	sigilDir := filepath.Join(dir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0755)
	os.MkdirAll(filepath.Join(sigilDir, "themes"), 0755)
	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte("name: test\nversion: \"1.0\""), 0644)

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"doctor", "--sigil-dir", sigilDir, "--no-color"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("doctor command: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Go") {
		t.Error("expected Go check")
	}
	if !strings.Contains(output, ".sigil dir") {
		t.Error("expected .sigil dir check")
	}
	if !strings.Contains(output, "sigil.yaml") {
		t.Error("expected sigil.yaml check")
	}
}

func TestDoctorCommandNoSigilDir(t *testing.T) {
	dir := t.TempDir()

	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"doctor", "--sigil-dir", filepath.Join(dir, "nonexistent"), "--no-color"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("doctor command: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "not found") {
		t.Error("expected 'not found' for missing .sigil dir")
	}
}
