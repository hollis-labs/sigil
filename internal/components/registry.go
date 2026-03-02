package components

import "sort"

// Schema defines a component type's props, actions, and slots.
type Schema struct {
	Type        string
	Category    string // primitives, layouts, data, forms, navigation, composites
	Description string
	Props       map[string]PropDef
	Actions     map[string]ActionDef
	Slots       map[string]SlotDef
	Shortcuts   []ShortcutDef
}

// PropDef defines a single component prop.
type PropDef struct {
	Type        string      // string, integer, boolean, array, object, enum
	Required    bool
	Default     interface{}
	Description string
	Enum        []string
}

// ActionDef defines a component action/event.
type ActionDef struct {
	Description string
}

// SlotDef defines a named slot that accepts child components.
type SlotDef struct {
	Description string
	Accepts     []string
}

// ShortcutDef defines a built-in keyboard shortcut for a component.
type ShortcutDef struct {
	Key         string
	Description string
	When        string
}

// Registry holds all known component schemas.
type Registry struct {
	schemas map[string]*Schema
}

// NewRegistry creates an empty component registry.
func NewRegistry() *Registry {
	return &Registry{schemas: make(map[string]*Schema)}
}

// Register adds a component schema to the registry.
func (r *Registry) Register(schema *Schema) {
	r.schemas[schema.Type] = schema
}

// Get returns the schema for a component type.
func (r *Registry) Get(componentType string) (*Schema, bool) {
	s, ok := r.schemas[componentType]
	return s, ok
}

// Has returns true if the component type is registered.
func (r *Registry) Has(componentType string) bool {
	_, ok := r.schemas[componentType]
	return ok
}

// Types returns all registered component type names, sorted.
func (r *Registry) Types() []string {
	types := make([]string, 0, len(r.schemas))
	for t := range r.schemas {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// ByCategory returns all schemas in the given category.
func (r *Registry) ByCategory(category string) []*Schema {
	var result []*Schema
	for _, s := range r.schemas {
		if s.Category == category {
			result = append(result, s)
		}
	}
	return result
}
