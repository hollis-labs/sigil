package preview

import (
	"strings"
	"testing"

	"github.com/hollis-labs/sigil/internal/config"
	"github.com/hollis-labs/sigil/internal/renderer"
)

func TestGenerateHTML_BasicPage(t *testing.T) {
	page := &config.Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "test-page",
		Title:   "Test Page",
		Overlay: "page",
		Layout: config.Component{
			ID:   "root",
			Type: "rows",
			Props: map[string]interface{}{
				"gap":     "4",
				"padding": "6",
			},
			Children: []config.Component{
				{
					ID:   "title",
					Type: "heading",
					Props: map[string]interface{}{
						"level": "2",
						"text":  "Hello World",
					},
				},
				{
					ID:   "desc",
					Type: "text",
					Props: map[string]interface{}{
						"text":    "A test page",
						"variant": "muted",
					},
				},
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)

	checks := []string{
		"<!DOCTYPE html>",
		"<title>Test Page — Sigil Preview</title>",
		"cdn.tailwindcss.com",
		"Hello World",
		"A test page",
		"flex flex-col gap-4",
		":root {",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output", check)
		}
	}
}

func TestGenerateHTML_ModalOverlay(t *testing.T) {
	page := &config.Page{
		Title:   "Modal Test",
		Overlay: "modal",
		Layout: config.Component{
			Type: "rows",
			Children: []config.Component{
				{Type: "heading", Props: map[string]interface{}{"text": "Modal Content"}},
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	if !strings.Contains(content, "bg-black/50") {
		t.Error("expected modal backdrop overlay")
	}
	if !strings.Contains(content, "Modal Content") {
		t.Error("expected modal content")
	}
}

func TestGenerateHTML_SheetOverlay(t *testing.T) {
	page := &config.Page{
		Title:   "Sheet Test",
		Overlay: "sheet",
		Layout: config.Component{
			Type:     "rows",
			Children: []config.Component{{Type: "heading", Props: map[string]interface{}{"text": "Sheet"}}},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	if !strings.Contains(content, "w-[480px]") {
		t.Error("expected sheet side panel width")
	}
}

func TestGenerateHTML_WithTheme(t *testing.T) {
	page := &config.Page{
		Title:   "Themed",
		Overlay: "page",
		Layout: config.Component{
			Type:     "rows",
			Children: []config.Component{{Type: "heading", Props: map[string]interface{}{"text": "Themed"}}},
		},
	}

	theme := &renderer.ThemeConfig{
		Name: "test-theme",
		Tokens: map[string]map[string]interface{}{
			"colors": {
				"accent":     "99 102 241",
				"background": "15 23 42",
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page, Theme: theme})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	if !strings.Contains(content, "--sigil-accent: 99 102 241") {
		t.Error("expected theme accent token in CSS")
	}
	if !strings.Contains(content, "--sigil-background: 15 23 42") {
		t.Error("expected theme background token in CSS")
	}
}

func TestGenerateHTML_DataTable(t *testing.T) {
	page := &config.Page{
		Title:   "Table Test",
		Overlay: "page",
		Layout: config.Component{
			Type: "rows",
			Children: []config.Component{
				{
					Type: "data-table",
					Props: map[string]interface{}{
						"datasource": "Users",
						"columns": []interface{}{
							map[string]interface{}{"label": "Name", "field": "name"},
							map[string]interface{}{"label": "Email", "field": "email"},
							map[string]interface{}{"label": "Status", "field": "status"},
						},
					},
				},
			},
		},
	}

	ds := &renderer.DataSourceManifest{
		Alias: "Users",
		Fields: []renderer.DataSourceField{
			{Name: "name", Type: "string"},
			{Name: "email", Type: "email"},
			{Name: "status", Type: "string", Values: []string{"active", "inactive", "suspended"}},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{
		Page:        page,
		DataSources: map[string]*renderer.DataSourceManifest{"Users": ds},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	checks := []string{
		"<th>Name</th>",
		"<th>Email</th>",
		"<th>Status</th>",
		"<tbody>",
		"<tr>",
		"@example.com", // from mock email generation
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in data table output", check)
		}
	}
}

func TestGenerateHTML_Form(t *testing.T) {
	page := &config.Page{
		Title:   "Form Test",
		Overlay: "page",
		Layout: config.Component{
			Type: "rows",
			Children: []config.Component{
				{
					Type: "form",
					Props: map[string]interface{}{
						"fields": []interface{}{
							map[string]interface{}{
								"name":  "title",
								"label": "Title",
								"type":  "text",
							},
							map[string]interface{}{
								"name":  "body",
								"label": "Body",
								"type":  "textarea",
								"rows":  "5",
							},
						},
						"submit": map[string]interface{}{
							"label": "Save",
						},
					},
				},
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	checks := []string{"<form", "Title", "Body", "<textarea", `rows="5"`, "Save"}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in form output", check)
		}
	}
}

func TestGenerateHTML_AllComponents(t *testing.T) {
	page := &config.Page{
		Title:   "Components Test",
		Overlay: "page",
		Layout: config.Component{
			Type: "rows",
			Children: []config.Component{
				{Type: "heading", Props: map[string]interface{}{"text": "Title", "level": "1"}},
				{Type: "text", Props: map[string]interface{}{"text": "Body text"}},
				{Type: "button", Props: map[string]interface{}{"label": "Click Me", "variant": "primary"}},
				{Type: "badge", Props: map[string]interface{}{"text": "New", "variant": "success"}},
				{Type: "search-bar", Props: map[string]interface{}{"placeholder": "Search..."}},
				{Type: "separator"},
				{Type: "progress", Props: map[string]interface{}{"value": "75", "max": "100"}},
				{Type: "alert", Props: map[string]interface{}{"message": "Warning!", "variant": "warning"}},
				{Type: "label", Props: map[string]interface{}{"text": "Field Label"}},
				{Type: "spacer", Props: map[string]interface{}{"size": "8"}},
				{Type: "input", Props: map[string]interface{}{"name": "field1", "label": "Input Label", "type": "text"}},
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	checks := []string{
		"<h1", "Title",
		"Body text",
		"Click Me",
		"sigil-badge-success", "New",
		`type="search"`,
		"<hr",
		"sigil-progress-bar",
		"sigil-alert-warning", "Warning!",
		"Field Label",
		"h-8",
		"Input Label",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in component output", check)
		}
	}
}

func TestGenerateHTML_Grid(t *testing.T) {
	page := &config.Page{
		Title:   "Grid Test",
		Overlay: "page",
		Layout: config.Component{
			Type: "grid",
			Props: map[string]interface{}{
				"columns": "3",
				"gap":     "6",
			},
			Children: []config.Component{
				{Type: "text", Props: map[string]interface{}{"text": "Cell 1"}},
				{Type: "text", Props: map[string]interface{}{"text": "Cell 2"}},
				{Type: "text", Props: map[string]interface{}{"text": "Cell 3"}},
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	if !strings.Contains(content, "grid grid-cols-3 gap-6") {
		t.Error("expected grid classes")
	}
	if !strings.Contains(content, "Cell 1") || !strings.Contains(content, "Cell 2") || !strings.Contains(content, "Cell 3") {
		t.Error("expected all grid cells")
	}
}

func TestGenerateHTML_Columns(t *testing.T) {
	page := &config.Page{
		Title:   "Columns Test",
		Overlay: "page",
		Layout: config.Component{
			Type: "columns",
			Props: map[string]interface{}{
				"gap":     "4",
				"justify": "between",
				"align":   "center",
			},
			Children: []config.Component{
				{Type: "text", Props: map[string]interface{}{"text": "Left"}},
				{Type: "text", Props: map[string]interface{}{"text": "Right"}},
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	if !strings.Contains(content, "flex flex-row gap-4 justify-between items-center") {
		t.Error("expected columns classes with justify-between and items-center")
	}
}

func TestGenerateHTML_DefaultTokens(t *testing.T) {
	page := &config.Page{
		Title:   "No Theme",
		Overlay: "page",
		Layout: config.Component{
			Type:     "rows",
			Children: []config.Component{{Type: "text", Props: map[string]interface{}{"text": "Test"}}},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	// Should have default tokens
	if !strings.Contains(content, "--sigil-background: 255 255 255") {
		t.Error("expected default background token")
	}
	if !strings.Contains(content, "--sigil-accent:") {
		t.Error("expected default accent token")
	}
}

func TestGenerateHTML_UnknownComponent(t *testing.T) {
	page := &config.Page{
		Title:   "Unknown Test",
		Overlay: "page",
		Layout: config.Component{
			Type: "rows",
			Children: []config.Component{
				{ID: "custom", Type: "custom-widget"},
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	if !strings.Contains(content, "data-component=\"custom-widget\"") {
		t.Error("expected unknown component placeholder")
	}
	if !strings.Contains(content, "[custom-widget]") {
		t.Error("expected component type label")
	}
}

func TestGenerateHTML_ButtonWithIcon(t *testing.T) {
	page := &config.Page{
		Title:   "Icon Test",
		Overlay: "page",
		Layout: config.Component{
			Type: "rows",
			Children: []config.Component{
				{Type: "button", Props: map[string]interface{}{"label": "Add", "icon": "plus", "variant": "primary"}},
			},
		},
	}

	html, err := GenerateHTML(&PreviewConfig{Page: page})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	if !strings.Contains(content, "+") {
		t.Error("expected icon symbol for 'plus'")
	}
	if !strings.Contains(content, "Add") {
		t.Error("expected button label")
	}
}

func TestMockValue_FieldInference(t *testing.T) {
	tests := []struct {
		field    string
		contains string
	}{
		{"id", "1"},
		{"name", "Dashboard"},
		{"email", "@example.com"},
		{"status", "active"},
		{"created_at", "2026-"},
		{"total", "142"},
	}
	for _, tt := range tests {
		val := mockValue(tt.field, tt.field, 0, nil)
		if !strings.Contains(val, tt.contains) {
			t.Errorf("mockValue(%q) = %q, expected to contain %q", tt.field, val, tt.contains)
		}
	}
}

func TestMockValueForType(t *testing.T) {
	// With explicit values
	val := mockValueForType("string", "status", []string{"on", "off"}, 0)
	if val != "on" {
		t.Errorf("expected 'on', got %q", val)
	}

	// Boolean
	val = mockValueForType("boolean", "active", nil, 0)
	if val != "true" {
		t.Errorf("expected 'true', got %q", val)
	}
	val = mockValueForType("boolean", "active", nil, 1)
	if val != "false" {
		t.Errorf("expected 'false', got %q", val)
	}

	// Email
	val = mockValueForType("email", "contact", nil, 0)
	if !strings.Contains(val, "@") {
		t.Errorf("expected email, got %q", val)
	}
}
