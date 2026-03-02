package config

import (
	"strings"
	"testing"
)

func TestDiffNoChanges(t *testing.T) {
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
				{ID: "h1", Type: "heading", Props: map[string]interface{}{"level": 2, "text": "Hello"}},
			},
		},
	}

	result := Diff(page, page)
	if result.HasChanges() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestDiffTitleChange(t *testing.T) {
	a := &Page{Title: "Old Title", Overlay: "page", Layout: Component{Type: "rows"}}
	b := &Page{Title: "New Title", Overlay: "page", Layout: Component{Type: "rows"}}

	result := Diff(a, b)
	if !result.HasChanges() {
		t.Fatal("expected changes")
	}

	found := false
	for _, c := range result.Changes {
		if c.Path == "title" && c.OldVal == "Old Title" && c.NewVal == "New Title" {
			found = true
		}
	}
	if !found {
		t.Error("expected title change")
	}
}

func TestDiffOverlayChange(t *testing.T) {
	a := &Page{Overlay: "page", Layout: Component{Type: "rows"}}
	b := &Page{Overlay: "modal", Layout: Component{Type: "rows"}}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if c.Path == "overlay" {
			found = true
		}
	}
	if !found {
		t.Error("expected overlay change")
	}
}

func TestDiffPropChange(t *testing.T) {
	a := &Page{
		Layout: Component{
			Type:  "heading",
			Props: map[string]interface{}{"text": "Old"},
		},
	}
	b := &Page{
		Layout: Component{
			Type:  "heading",
			Props: map[string]interface{}{"text": "New"},
		},
	}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if strings.Contains(c.Path, "props.text") && c.Type == "modified" {
			found = true
		}
	}
	if !found {
		t.Error("expected prop change")
	}
}

func TestDiffAddedProp(t *testing.T) {
	a := &Page{
		Layout: Component{
			Type:  "button",
			Props: map[string]interface{}{"label": "Click"},
		},
	}
	b := &Page{
		Layout: Component{
			Type:  "button",
			Props: map[string]interface{}{"label": "Click", "variant": "primary"},
		},
	}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if c.Type == "added" && strings.Contains(c.Path, "variant") {
			found = true
		}
	}
	if !found {
		t.Error("expected added prop")
	}
}

func TestDiffRemovedProp(t *testing.T) {
	a := &Page{
		Layout: Component{
			Type:  "button",
			Props: map[string]interface{}{"label": "Click", "icon": "plus"},
		},
	}
	b := &Page{
		Layout: Component{
			Type:  "button",
			Props: map[string]interface{}{"label": "Click"},
		},
	}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if c.Type == "removed" && strings.Contains(c.Path, "icon") {
			found = true
		}
	}
	if !found {
		t.Error("expected removed prop")
	}
}

func TestDiffAddedChild(t *testing.T) {
	a := &Page{
		Layout: Component{
			Type: "rows",
			Children: []Component{
				{Type: "heading", Props: map[string]interface{}{"text": "Hello"}},
			},
		},
	}
	b := &Page{
		Layout: Component{
			Type: "rows",
			Children: []Component{
				{Type: "heading", Props: map[string]interface{}{"text": "Hello"}},
				{Type: "button", Props: map[string]interface{}{"label": "New"}},
			},
		},
	}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if c.Type == "added" && strings.Contains(c.Path, "children[1]") {
			found = true
		}
	}
	if !found {
		t.Error("expected added child")
	}
}

func TestDiffRemovedChild(t *testing.T) {
	a := &Page{
		Layout: Component{
			Type: "rows",
			Children: []Component{
				{Type: "heading"},
				{Type: "button"},
			},
		},
	}
	b := &Page{
		Layout: Component{
			Type:     "rows",
			Children: []Component{{Type: "heading"}},
		},
	}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if c.Type == "removed" && strings.Contains(c.Path, "children[1]") {
			found = true
		}
	}
	if !found {
		t.Error("expected removed child")
	}
}

func TestDiffTypeChange(t *testing.T) {
	a := &Page{
		Layout: Component{Type: "rows"},
	}
	b := &Page{
		Layout: Component{Type: "columns"},
	}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if c.Type == "modified" && strings.Contains(c.Path, "type") {
			found = true
		}
	}
	if !found {
		t.Error("expected type change")
	}
}

func TestDiffDataSourceAdded(t *testing.T) {
	a := &Page{Layout: Component{Type: "rows"}}
	b := &Page{
		Layout:      Component{Type: "rows"},
		DataSources: []DataSourceRef{{Alias: "users"}},
	}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if c.Type == "added" && c.NewVal == "users" {
			found = true
		}
	}
	if !found {
		t.Error("expected added datasource")
	}
}

func TestDiffDataSourceRemoved(t *testing.T) {
	a := &Page{
		Layout:      Component{Type: "rows"},
		DataSources: []DataSourceRef{{Alias: "users"}},
	}
	b := &Page{Layout: Component{Type: "rows"}}

	result := Diff(a, b)
	found := false
	for _, c := range result.Changes {
		if c.Type == "removed" && c.OldVal == "users" {
			found = true
		}
	}
	if !found {
		t.Error("expected removed datasource")
	}
}

func TestFormatDiffNoColor(t *testing.T) {
	result := &DiffResult{
		Changes: []DiffChange{
			{Type: "added", Path: "layout.children[1]", NewVal: "button"},
			{Type: "removed", Path: "layout.props.icon", OldVal: "plus"},
			{Type: "modified", Path: "title", OldVal: "Old", NewVal: "New"},
		},
	}

	output := FormatDiff(result, false)
	if !strings.Contains(output, "+ layout.children[1]") {
		t.Error("expected + prefix for added")
	}
	if !strings.Contains(output, "- layout.props.icon") {
		t.Error("expected - prefix for removed")
	}
	if !strings.Contains(output, "~ title") {
		t.Error("expected ~ prefix for modified")
	}
	if !strings.Contains(output, "Old → New") {
		t.Error("expected old → new format")
	}
}

func TestFormatDiffNoChanges(t *testing.T) {
	result := &DiffResult{}
	output := FormatDiff(result, false)
	if output != "No changes" {
		t.Errorf("expected 'No changes', got %q", output)
	}
}
