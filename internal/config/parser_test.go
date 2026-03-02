package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseValidYAML(t *testing.T) {
	yaml := `
sigil: "1.0"
kind: page
id: test-page
title: Test Page
overlay: page

layout:
  type: rows
  props:
    gap: 4
  children:
    - id: heading
      type: heading
      props:
        level: 2
        text: Hello
`
	page, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.ID != "test-page" {
		t.Errorf("expected id 'test-page', got %q", page.ID)
	}
	if page.Title != "Test Page" {
		t.Errorf("expected title 'Test Page', got %q", page.Title)
	}
	if page.Layout.Type != "rows" {
		t.Errorf("expected layout type 'rows', got %q", page.Layout.Type)
	}
}

func TestParseInvalidYAML(t *testing.T) {
	yaml := `
sigil: "1.0"
  bad indent: [
`
	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestDefaultsApplied(t *testing.T) {
	yaml := `
id: minimal-page
title: Minimal
layout:
  type: rows
`
	page, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Sigil != "1.0" {
		t.Errorf("expected default sigil '1.0', got %q", page.Sigil)
	}
	if page.Kind != "page" {
		t.Errorf("expected default kind 'page', got %q", page.Kind)
	}
	if page.Overlay != "page" {
		t.Errorf("expected default overlay 'page', got %q", page.Overlay)
	}
}

func TestNestedComponents(t *testing.T) {
	yaml := `
sigil: "1.0"
kind: page
id: nested
title: Nested
overlay: page

layout:
  type: rows
  children:
    - type: columns
      children:
        - id: btn
          type: button
          props:
            label: Click
`
	page, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Layout.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(page.Layout.Children))
	}
	cols := page.Layout.Children[0]
	if cols.Type != "columns" {
		t.Errorf("expected columns, got %q", cols.Type)
	}
	if len(cols.Children) != 1 {
		t.Fatalf("expected 1 nested child, got %d", len(cols.Children))
	}
	btn := cols.Children[0]
	if btn.ID != "btn" {
		t.Errorf("expected id 'btn', got %q", btn.ID)
	}
	if btn.Props["label"] != "Click" {
		t.Errorf("expected label 'Click', got %v", btn.Props["label"])
	}
}

func TestActionsAndShortcuts(t *testing.T) {
	yaml := `
sigil: "1.0"
kind: page
id: actions-test
title: Actions Test
overlay: page

layout:
  type: rows
  children:
    - id: btn
      type: button
      actions:
        click:
          type: navigate
          page: target-page
          params:
            id: "123"

shortcuts:
  - key: Escape
    action:
      type: close
    description: Close overlay
  - key: Ctrl+N
    action:
      type: modal
      title: New Item
    global: true
`
	page, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	btn := page.Layout.Children[0]
	click, ok := btn.Actions["click"]
	if !ok {
		t.Fatal("expected click action")
	}
	if click.Type != "navigate" {
		t.Errorf("expected navigate, got %q", click.Type)
	}
	if click.Page != "target-page" {
		t.Errorf("expected page 'target-page', got %q", click.Page)
	}
	if click.Params["id"] != "123" {
		t.Errorf("expected param id '123', got %q", click.Params["id"])
	}

	if len(page.Shortcuts) != 2 {
		t.Fatalf("expected 2 shortcuts, got %d", len(page.Shortcuts))
	}
	if page.Shortcuts[0].Key != "Escape" {
		t.Errorf("expected key 'Escape', got %q", page.Shortcuts[0].Key)
	}
	if page.Shortcuts[0].Action.Type != "close" {
		t.Errorf("expected action type 'close', got %q", page.Shortcuts[0].Action.Type)
	}
	if page.Shortcuts[1].Global != true {
		t.Error("expected second shortcut to be global")
	}
}

func TestParseFile(t *testing.T) {
	content := `
sigil: "1.0"
kind: page
id: file-test
title: File Test
overlay: page
layout:
  type: rows
`
	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	page, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.ID != "file-test" {
		t.Errorf("expected id 'file-test', got %q", page.ID)
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestDataSourcesParsed(t *testing.T) {
	yaml := `
sigil: "1.0"
kind: page
id: ds-test
title: DataSource Test
overlay: page
datasources:
  - alias: Sprint
    capabilities: [search, filter, sort]
  - alias: Task
    capabilities: [search]
    params:
      filter:
        sprint_id: "123"
layout:
  type: rows
`
	page, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.DataSources) != 2 {
		t.Fatalf("expected 2 datasources, got %d", len(page.DataSources))
	}
	if page.DataSources[0].Alias != "Sprint" {
		t.Errorf("expected alias 'Sprint', got %q", page.DataSources[0].Alias)
	}
	if len(page.DataSources[0].Capabilities) != 3 {
		t.Errorf("expected 3 capabilities, got %d", len(page.DataSources[0].Capabilities))
	}
}

func TestMetaParsed(t *testing.T) {
	yaml := `
sigil: "1.0"
kind: page
id: meta-test
title: Meta Test
overlay: page
layout:
  type: rows
meta:
  guards: [auth, admin]
  tags: [dashboard]
  version: 2
`
	page, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Meta == nil {
		t.Fatal("expected meta to be non-nil")
	}
	if len(page.Meta.Guards) != 2 {
		t.Errorf("expected 2 guards, got %d", len(page.Meta.Guards))
	}
	if page.Meta.Version != 2 {
		t.Errorf("expected version 2, got %d", page.Meta.Version)
	}
}
