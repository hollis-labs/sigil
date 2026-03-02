package config

import (
	"testing"
)

func TestMigrateNoChanges(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{ID: "h1", Type: "heading"},
			},
		},
	}

	result := Migrate(page)
	if result.HasChanges() {
		t.Errorf("expected no changes, got %d: %v", len(result.Changes), result.Changes)
	}
}

func TestMigrateAutoGenerateIDs(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{Type: "heading"},
				{Type: "button"},
				{Type: "data-table"},
			},
		},
	}

	result := Migrate(page)
	if !result.HasChanges() {
		t.Fatal("expected changes")
	}

	// Verify IDs were generated
	for i, child := range page.Layout.Children {
		if child.ID == "" {
			t.Errorf("child %d still has no ID", i)
		}
	}

	// Check specific IDs
	if page.Layout.Children[0].ID != "root_heading_0" {
		t.Errorf("expected root_heading_0, got %s", page.Layout.Children[0].ID)
	}
	if page.Layout.Children[1].ID != "root_button_1" {
		t.Errorf("expected root_button_1, got %s", page.Layout.Children[1].ID)
	}
	if page.Layout.Children[2].ID != "root_data_table_2" {
		t.Errorf("expected root_data_table_2, got %s", page.Layout.Children[2].ID)
	}
}

func TestMigrateSetDefaults(t *testing.T) {
	page := &Page{
		ID:     "test",
		Layout: Component{ID: "root", Type: "rows"},
	}

	result := Migrate(page)
	if !result.HasChanges() {
		t.Fatal("expected changes")
	}

	if page.Sigil != "1.0" {
		t.Errorf("expected sigil=1.0, got %s", page.Sigil)
	}
	if page.Kind != "page" {
		t.Errorf("expected kind=page, got %s", page.Kind)
	}
	if page.Overlay != "page" {
		t.Errorf("expected overlay=page, got %s", page.Overlay)
	}
}

func TestMigrateNestedIDs(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		Kind:    "page",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{
					Type: "columns",
					Children: []Component{
						{Type: "button"},
						{Type: "badge"},
					},
				},
			},
		},
	}

	Migrate(page)

	// columns should get ID
	cols := page.Layout.Children[0]
	if cols.ID == "" {
		t.Error("columns should have an ID")
	}

	// Nested children should have IDs based on parent
	for _, child := range cols.Children {
		if child.ID == "" {
			t.Error("nested child should have an ID")
		}
	}
}

func TestFormatMigrateResultNoChanges(t *testing.T) {
	result := &MigrateResult{FromVersion: "1.0", ToVersion: "1.0"}
	output := FormatMigrateResult(result, false)
	if output != "No migrations needed" {
		t.Errorf("expected 'No migrations needed', got %q", output)
	}
}

func TestFormatMigrateResultWithChanges(t *testing.T) {
	result := &MigrateResult{
		FromVersion: "1.0",
		ToVersion:   "1.0",
		Changes: []MigrateChange{
			{Path: "component[root_heading_0]", Message: "auto-generated ID"},
		},
	}
	output := FormatMigrateResult(result, false)
	if output == "No migrations needed" {
		t.Error("expected changes in output")
	}
}
