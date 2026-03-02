package gotempl

import (
	"strings"
	"testing"
)

func TestRenderSharedComponentsAll(t *testing.T) {
	usedTypes := []string{
		"button", "heading", "text", "badge", "input", "select",
		"search-bar", "data-table", "modal", "form",
	}

	files, err := renderSharedComponents(usedTypes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have layout_helpers + all 10 component files
	if len(files) != 11 {
		t.Errorf("expected 11 files, got %d", len(files))
		for _, f := range files {
			t.Logf("  %s", f.Path)
		}
	}

	// Verify all expected files
	expectedPaths := map[string]bool{
		"components/layout_helpers.templ": false,
		"components/button.templ":         false,
		"components/heading.templ":        false,
		"components/text.templ":           false,
		"components/badge.templ":          false,
		"components/input.templ":          false,
		"components/select_field.templ":   false,
		"components/search_bar.templ":     false,
		"components/data_table.templ":     false,
		"components/modal.templ":          false,
		"components/form.templ":           false,
	}

	for _, f := range files {
		if _, ok := expectedPaths[f.Path]; ok {
			expectedPaths[f.Path] = true
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf("expected file %q not generated", path)
		}
	}
}

func TestRenderSharedComponentsSubset(t *testing.T) {
	usedTypes := []string{"button", "heading"}

	files, err := renderSharedComponents(usedTypes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have layout_helpers + button + heading = 3
	if len(files) != 3 {
		t.Errorf("expected 3 files, got %d", len(files))
	}

	paths := map[string]bool{}
	for _, f := range files {
		paths[f.Path] = true
	}

	if !paths["components/layout_helpers.templ"] {
		t.Error("expected layout_helpers.templ")
	}
	if !paths["components/button.templ"] {
		t.Error("expected button.templ")
	}
	if !paths["components/heading.templ"] {
		t.Error("expected heading.templ")
	}
	if paths["components/data_table.templ"] {
		t.Error("unexpected data_table.templ (not used)")
	}
}

func TestSharedComponentsContent(t *testing.T) {
	tests := []struct {
		name     string
		types    []string
		path     string
		contains []string
	}{
		{
			name:  "button has props struct",
			types: []string{"button"},
			path:  "components/button.templ",
			contains: []string{
				"ButtonProps",
				"Label",
				"Variant",
				"sigil-btn",
				"hx-get",
				"package components",
			},
		},
		{
			name:  "heading has level switch",
			types: []string{"heading"},
			path:  "components/heading.templ",
			contains: []string{
				"HeadingProps",
				"Level",
				"--sigil-text",
				"package components",
			},
		},
		{
			name:  "data-table has HTMX",
			types: []string{"data-table"},
			path:  "components/data_table.templ",
			contains: []string{
				"DataTableProps",
				"Column",
				"hx-get",
				"hx-trigger",
				"package components",
			},
		},
		{
			name:  "modal has overlay",
			types: []string{"modal"},
			path:  "components/modal.templ",
			contains: []string{
				"ModalProps",
				"data-modal",
				"z-50",
				"package components",
			},
		},
		{
			name:  "form has HTMX",
			types: []string{"form"},
			path:  "components/form.templ",
			contains: []string{
				"FormProps",
				"hx-post",
				"sigil-form",
				"package components",
			},
		},
		{
			name:  "layout helpers",
			types: []string{},
			path:  "components/layout_helpers.templ",
			contains: []string{
				"templ Rows",
				"templ Columns",
				"flex flex-col",
				"flex flex-row",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := renderSharedComponents(tt.types)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var content string
			for _, f := range files {
				if f.Path == tt.path {
					content = string(f.Content)
					break
				}
			}
			if content == "" {
				t.Fatalf("file %q not found", tt.path)
			}

			for _, c := range tt.contains {
				if !strings.Contains(content, c) {
					t.Errorf("expected %q in %s content", c, tt.path)
				}
			}
		})
	}
}

func TestSharedComponentsEmpty(t *testing.T) {
	// Even with no used types, should still get layout helpers
	files, err := renderSharedComponents([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file (layout_helpers), got %d", len(files))
	}
}
