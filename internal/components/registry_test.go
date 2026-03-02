package components

import "testing"

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if len(r.Types()) != 0 {
		t.Errorf("new registry should be empty, got %d types", len(r.Types()))
	}
}

func TestRegisterAndGet(t *testing.T) {
	r := NewRegistry()
	s := &Schema{Type: "button", Category: "primitives", Description: "A button"}
	r.Register(s)

	got, ok := r.Get("button")
	if !ok {
		t.Fatal("expected to find 'button'")
	}
	if got.Description != "A button" {
		t.Errorf("expected description 'A button', got %q", got.Description)
	}
}

func TestHas(t *testing.T) {
	r := NewRegistry()
	r.Register(&Schema{Type: "heading"})

	if !r.Has("heading") {
		t.Error("expected Has('heading') to be true")
	}
	if r.Has("nonexistent") {
		t.Error("expected Has('nonexistent') to be false")
	}
}

func TestTypes(t *testing.T) {
	r := NewRegistry()
	r.Register(&Schema{Type: "button"})
	r.Register(&Schema{Type: "alert"})
	r.Register(&Schema{Type: "heading"})

	types := r.Types()
	if len(types) != 3 {
		t.Fatalf("expected 3 types, got %d", len(types))
	}
	// Should be sorted
	if types[0] != "alert" || types[1] != "button" || types[2] != "heading" {
		t.Errorf("expected sorted [alert, button, heading], got %v", types)
	}
}

func TestByCategory(t *testing.T) {
	r := NewRegistry()
	r.Register(&Schema{Type: "button", Category: "primitives"})
	r.Register(&Schema{Type: "heading", Category: "primitives"})
	r.Register(&Schema{Type: "rows", Category: "layouts"})

	prims := r.ByCategory("primitives")
	if len(prims) != 2 {
		t.Errorf("expected 2 primitives, got %d", len(prims))
	}

	layouts := r.ByCategory("layouts")
	if len(layouts) != 1 {
		t.Errorf("expected 1 layout, got %d", len(layouts))
	}

	empty := r.ByCategory("nonexistent")
	if len(empty) != 0 {
		t.Errorf("expected 0 for unknown category, got %d", len(empty))
	}
}

func TestDefaultRegistryHasExpectedTypes(t *testing.T) {
	r := NewDefaultRegistry()

	types := r.Types()
	if len(types) < 40 {
		t.Errorf("expected at least 40 registered types, got %d", len(types))
	}

	// Spot-check key types from each category
	expected := []string{
		"heading", "button", "badge",        // primitives
		"rows", "columns", "grid", "card",   // layouts
		"data-table", "search-bar", "form",  // data + forms
		"breadcrumb", "nav-menu",            // navigation
		"modal", "sheet", "tooltip",         // composites
	}
	for _, typ := range expected {
		if !r.Has(typ) {
			t.Errorf("expected default registry to have %q", typ)
		}
	}
}

func TestDefaultRegistryCategories(t *testing.T) {
	r := NewDefaultRegistry()

	categories := map[string]int{
		"primitives": 16,
		"layouts":    12,
		"navigation": 4,
		"composites": 8,
		"data":       7,
		"forms":      2,
	}

	for cat, expectedCount := range categories {
		got := r.ByCategory(cat)
		if len(got) != expectedCount {
			t.Errorf("category %q: expected %d types, got %d", cat, expectedCount, len(got))
		}
	}
}

func TestRegistryOverwrite(t *testing.T) {
	r := NewRegistry()
	r.Register(&Schema{Type: "button", Description: "v1"})
	r.Register(&Schema{Type: "button", Description: "v2"})

	got, _ := r.Get("button")
	if got.Description != "v2" {
		t.Errorf("expected overwritten description 'v2', got %q", got.Description)
	}
	if len(r.Types()) != 1 {
		t.Errorf("expected 1 type after overwrite, got %d", len(r.Types()))
	}
}
