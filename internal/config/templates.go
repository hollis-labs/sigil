package config

import "strings"

// NewPageTemplate creates a starter Page config from parameters.
func NewPageTemplate(id, title, overlay, module, layout, datasource string) *Page {
	if title == "" {
		title = formatID(id)
	}
	if overlay == "" {
		overlay = "page"
	}
	if layout == "" {
		layout = "rows"
	}

	page := &Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      id,
		Title:   title,
		Overlay: overlay,
		Layout: Component{
			ID:   "root",
			Type: layout,
			Props: map[string]interface{}{
				"gap":     4,
				"padding": 6,
			},
			Children: []Component{
				{ID: "page-title", Type: "heading", Props: map[string]interface{}{"level": 2, "text": title}},
			},
		},
		Shortcuts: []Shortcut{
			{Key: "Escape", Action: Action{Type: "close"}},
		},
	}

	if module != "" {
		page.Module = module
	}

	if datasource != "" {
		page.DataSources = []DataSourceRef{
			{
				Alias:        datasource,
				Capabilities: []string{"search", "filter", "sort", "paginate"},
			},
		}
		page.Layout.Children = append(page.Layout.Children, Component{
			ID:   "table-" + strings.ToLower(datasource),
			Type: "data-table",
			Props: map[string]interface{}{
				"datasource": datasource,
				"columns":    []interface{}{},
				"pagination":  map[string]interface{}{"enabled": true, "pageSize": 25},
			},
		})
	}

	return page
}

// formatID converts kebab-case to Title Case.
func formatID(id string) string {
	parts := strings.Split(id, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
