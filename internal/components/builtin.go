package components

// NewDefaultRegistry creates a registry pre-populated with all built-in component types.
// Full prop/action/slot schemas will be added in Sprint 2.
func NewDefaultRegistry() *Registry {
	r := NewRegistry()

	// Primitives
	for _, s := range []Schema{
		{Type: "heading", Category: "primitives", Description: "Section heading"},
		{Type: "text", Category: "primitives", Description: "Text paragraph"},
		{Type: "button", Category: "primitives", Description: "Clickable button"},
		{Type: "icon-button", Category: "primitives", Description: "Icon-only button"},
		{Type: "input", Category: "primitives", Description: "Text input"},
		{Type: "textarea", Category: "primitives", Description: "Multi-line input"},
		{Type: "select", Category: "primitives", Description: "Dropdown select"},
		{Type: "checkbox", Category: "primitives", Description: "Checkbox input"},
		{Type: "switch", Category: "primitives", Description: "Toggle switch"},
		{Type: "badge", Category: "primitives", Description: "Status badge"},
		{Type: "avatar", Category: "primitives", Description: "User avatar"},
		{Type: "separator", Category: "primitives", Description: "Visual divider"},
		{Type: "progress", Category: "primitives", Description: "Progress bar"},
		{Type: "alert", Category: "primitives", Description: "Alert message"},
		{Type: "label", Category: "primitives", Description: "Form label"},
		{Type: "icon", Category: "primitives", Description: "Icon display"},
	} {
		s := s
		r.Register(&s)
	}

	// Layouts
	for _, s := range []Schema{
		{Type: "rows", Category: "layouts", Description: "Vertical stack"},
		{Type: "columns", Category: "layouts", Description: "Horizontal columns"},
		{Type: "grid", Category: "layouts", Description: "CSS grid"},
		{Type: "card", Category: "layouts", Description: "Container with border"},
		{Type: "tabs", Category: "layouts", Description: "Tabbed panels"},
		{Type: "tab", Category: "layouts", Description: "Single tab content"},
		{Type: "split", Category: "layouts", Description: "Resizable split panes"},
		{Type: "sidebar", Category: "layouts", Description: "Sidebar + content"},
		{Type: "accordion", Category: "layouts", Description: "Collapsible sections"},
		{Type: "accordion-item", Category: "layouts", Description: "Accordion section"},
		{Type: "scroll-area", Category: "layouts", Description: "Scrollable region"},
		{Type: "spacer", Category: "layouts", Description: "Empty space"},
	} {
		s := s
		r.Register(&s)
	}

	// Navigation
	for _, s := range []Schema{
		{Type: "breadcrumb", Category: "navigation", Description: "Navigation breadcrumbs"},
		{Type: "pagination", Category: "navigation", Description: "Page navigation"},
		{Type: "nav-menu", Category: "navigation", Description: "Navigation menu"},
		{Type: "command-palette", Category: "navigation", Description: "Searchable command list"},
	} {
		s := s
		r.Register(&s)
	}

	// Composites
	for _, s := range []Schema{
		{Type: "modal", Category: "composites", Description: "Dialog overlay"},
		{Type: "sheet", Category: "composites", Description: "Side panel"},
		{Type: "dropdown-menu", Category: "composites", Description: "Dropdown with items"},
		{Type: "context-menu", Category: "composites", Description: "Right-click menu"},
		{Type: "tooltip", Category: "composites", Description: "Hover tooltip"},
		{Type: "popover", Category: "composites", Description: "Click popover"},
		{Type: "confirm-dialog", Category: "composites", Description: "Confirmation modal"},
		{Type: "toast", Category: "composites", Description: "Notification toast"},
	} {
		s := s
		r.Register(&s)
	}

	// Data
	for _, s := range []Schema{
		{Type: "data-table", Category: "data", Description: "Sortable/filterable data table"},
		{Type: "list", Category: "data", Description: "Simple data list"},
		{Type: "detail-view", Category: "data", Description: "Single record display"},
		{Type: "stat-card", Category: "data", Description: "Metric display"},
		{Type: "chart", Category: "data", Description: "Data chart"},
		{Type: "timeline", Category: "data", Description: "Event timeline"},
		{Type: "search-bar", Category: "data", Description: "Search input with targeting"},
	} {
		s := s
		r.Register(&s)
	}

	// Forms
	for _, s := range []Schema{
		{Type: "form", Category: "forms", Description: "Form container"},
		{Type: "field-group", Category: "forms", Description: "Grouped fields"},
	} {
		s := s
		r.Register(&s)
	}

	return r
}
