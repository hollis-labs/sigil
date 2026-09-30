package components

import (
	"sort"

	"github.com/hollis-labs/sigil/internal/config"
)

// Schema defines a component type's props, actions, and slots.
type Schema struct {
	Type        string               `yaml:"type"`
	Category    string               `yaml:"category"`
	Description string               `yaml:"description"`
	Props       map[string]PropDef   `yaml:"props,omitempty"`
	Actions     map[string]ActionDef `yaml:"actions,omitempty"`
	Slots       map[string]SlotDef   `yaml:"slots,omitempty"`
	Shortcuts   []ShortcutDef        `yaml:"shortcuts,omitempty"`
	Source      *ComponentSource     `yaml:"source,omitempty"`
}

// ComponentSource defines the source files for a custom component.
type ComponentSource struct {
	Component string   `yaml:"component"`          // main .tsx file (required for custom)
	Includes  []string `yaml:"includes,omitempty"` // additional files to copy (hooks, utils)
}

// IsCustom returns true if this schema has source files (custom component).
func (s *Schema) IsCustom() bool {
	return s.Source != nil && s.Source.Component != ""
}

// PropDef defines a single component prop.
type PropDef struct {
	Type        string      `yaml:"type"`
	Required    bool        `yaml:"required,omitempty"`
	Default     interface{} `yaml:"default,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Enum        []string    `yaml:"enum,omitempty"`
}

// ActionDef defines a component action/event.
type ActionDef struct {
	Description string `yaml:"description"`
}

// SlotDef defines a named slot that accepts child components.
type SlotDef struct {
	Description string   `yaml:"description"`
	Accepts     []string `yaml:"accepts,omitempty"`
}

// ShortcutDef defines a built-in keyboard shortcut for a component.
type ShortcutDef struct {
	Key         string `yaml:"key"`
	Description string `yaml:"description"`
	When        string `yaml:"when,omitempty"`
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

// GetSchema returns a ComponentSchemaInfo for deep validation.
// Implements config.ComponentSchemaProvider.
func (r *Registry) GetSchema(componentType string) (*config.ComponentSchemaInfo, bool) {
	s, ok := r.schemas[componentType]
	if !ok {
		return nil, false
	}
	info := &config.ComponentSchemaInfo{
		Props: make(map[string]config.PropInfo, len(s.Props)),
	}
	for name, p := range s.Props {
		info.Props[name] = config.PropInfo{
			Type:     p.Type,
			Required: p.Required,
			Enum:     p.Enum,
		}
	}
	return info, true
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
