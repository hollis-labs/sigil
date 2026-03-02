package config

import "testing"

// mockRegistry implements ComponentRegistry for testing.
type mockRegistry struct {
	types map[string]bool
}

func (r *mockRegistry) Has(t string) bool {
	return r.types[t]
}

func newMockRegistry(types ...string) *mockRegistry {
	m := &mockRegistry{types: make(map[string]bool)}
	for _, t := range types {
		m.types[t] = true
	}
	return m
}

func TestValidateValidConfig(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "test-page",
		Title:   "Test Page",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{ID: "heading", Type: "heading"},
			},
		},
	}
	result := Validate(page, newMockRegistry("rows", "heading"))
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidateMissingRequiredFields(t *testing.T) {
	page := &Page{}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for empty page")
	}
	// Should have errors for: sigil, id, title, overlay, layout.type
	if len(result.Errors) < 5 {
		t.Errorf("expected at least 5 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

func TestValidateInvalidOverlay(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "invalid-type",
		Layout:  Component{Type: "rows"},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for bad overlay")
	}
	found := false
	for _, e := range result.Errors {
		if e.Path == "overlay" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected overlay error")
	}
}

func TestValidateAllOverlayTypes(t *testing.T) {
	for _, overlay := range []string{"page", "modal", "sheet", "drawer", "fullscreen"} {
		page := &Page{
			Sigil:   "1.0",
			ID:      "test",
			Title:   "Test",
			Overlay: overlay,
			Layout:  Component{ID: "root", Type: "rows"},
		}
		result := Validate(page, nil)
		if !result.Valid {
			t.Errorf("overlay %q should be valid, got errors: %v", overlay, result.Errors)
		}
	}
}

func TestValidateUnknownComponentType(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{ID: "unknown", Type: "nonexistent-widget"},
			},
		},
	}
	reg := newMockRegistry("rows")
	result := Validate(page, reg)
	if result.Valid {
		t.Error("expected invalid for unknown component type")
	}
}

func TestValidateDuplicateIDs(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{ID: "dup", Type: "button"},
				{ID: "dup", Type: "button"},
			},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for duplicate IDs")
	}
	found := false
	for _, e := range result.Errors {
		if e.Message == `duplicate component id "dup"` {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected duplicate id error, got: %v", result.Errors)
	}
}

func TestValidateInvalidActionType(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{
					ID:   "btn",
					Type: "button",
					Actions: map[string]Action{
						"click": {Type: "teleport"},
					},
				},
			},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for bad action type")
	}
}

func TestValidateModalActionRequirements(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{
					ID:   "btn",
					Type: "button",
					Actions: map[string]Action{
						"click": {Type: "modal"}, // missing title and fields
					},
				},
			},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for modal without title/fields")
	}
	if len(result.Errors) < 2 {
		t.Errorf("expected at least 2 errors (title + fields), got %d", len(result.Errors))
	}
}

func TestValidateNavigateActionRequirements(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{
					ID:   "btn",
					Type: "button",
					Actions: map[string]Action{
						"click": {Type: "navigate"}, // missing page and url
					},
				},
			},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for navigate without page/url")
	}
}

func TestValidateHTTPActionRequirements(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{
					ID:   "btn",
					Type: "button",
					Actions: map[string]Action{
						"click": {Type: "http"}, // missing url and method
					},
				},
			},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for http without url/method")
	}
	if len(result.Errors) < 2 {
		t.Errorf("expected at least 2 errors (url + method), got %d", len(result.Errors))
	}
}

func TestValidateDataSourceRef(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout:  Component{ID: "root", Type: "rows"},
		DataSources: []DataSourceRef{
			{Alias: ""},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for datasource without alias")
	}
}

func TestValidateShortcutKeyRequired(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout:  Component{ID: "root", Type: "rows"},
		Shortcuts: []Shortcut{
			{Key: "", Action: Action{Type: "close"}},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for shortcut without key")
	}
}

func TestValidateFormFieldRequirements(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{
					ID:   "btn",
					Type: "button",
					Actions: map[string]Action{
						"click": {
							Type:  "modal",
							Title: "Test",
							Fields: []FormField{
								{}, // missing name, label, type
							},
						},
					},
				},
			},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for form field without name/label/type")
	}
	// Should have errors for name, label, type
	fieldErrors := 0
	for _, e := range result.Errors {
		if e.Path == "layout.children[0].actions.click.fields[0]" {
			fieldErrors++
		}
	}
	if fieldErrors < 3 {
		t.Errorf("expected at least 3 field errors, got %d", fieldErrors)
	}
}

func TestValidateSelectFieldWarning(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{
					ID:   "btn",
					Type: "button",
					Actions: map[string]Action{
						"click": {
							Type:  "modal",
							Title: "Test",
							Fields: []FormField{
								{Name: "role", Label: "Role", Type: "select"}, // no options
							},
						},
					},
				},
			},
		},
	}
	result := Validate(page, nil)
	if len(result.Warnings) == 0 {
		t.Error("expected warning for select without options")
	}
}

func TestValidateComponentWithoutIDWarning(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{Type: "button"}, // no id
			},
		},
	}
	result := Validate(page, nil)
	if len(result.Warnings) == 0 {
		t.Error("expected warning for component without id")
	}
}

func TestValidateNestedOnConfirm(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "rows",
			Children: []Component{
				{
					ID:   "btn",
					Type: "button",
					Actions: map[string]Action{
						"click": {
							Type:    "confirm",
							Title:   "Sure?",
							Message: "Really?",
							OnConfirm: &Action{
								Type: "http", // missing url and method
							},
						},
					},
				},
			},
		},
	}
	result := Validate(page, nil)
	if result.Valid {
		t.Error("expected invalid for nested onConfirm with bad http action")
	}
}

func TestValidateNilRegistry(t *testing.T) {
	page := &Page{
		Sigil:   "1.0",
		ID:      "test",
		Title:   "Test",
		Overlay: "page",
		Layout: Component{
			ID:   "root",
			Type: "anything-goes",
		},
	}
	result := Validate(page, nil)
	if !result.Valid {
		t.Errorf("expected valid with nil registry (type check skipped), got: %v", result.Errors)
	}
}
