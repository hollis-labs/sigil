package config

// Page represents a complete Sigil page config.
type Page struct {
	Sigil       string          `yaml:"sigil" json:"sigil"`
	Kind        string          `yaml:"kind" json:"kind"`
	ID          string          `yaml:"id" json:"id"`
	Title       string          `yaml:"title" json:"title"`
	Description string          `yaml:"description,omitempty" json:"description,omitempty"`
	Overlay     string          `yaml:"overlay" json:"overlay"`
	Module      string          `yaml:"module,omitempty" json:"module,omitempty"`
	Theme       *ThemeRef       `yaml:"theme,omitempty" json:"theme,omitempty"`
	DataSources []DataSourceRef `yaml:"datasources,omitempty" json:"datasources,omitempty"`
	Layout      Component       `yaml:"layout" json:"layout"`
	Shortcuts   []Shortcut      `yaml:"shortcuts,omitempty" json:"shortcuts,omitempty"`
	Meta        *PageMeta       `yaml:"meta,omitempty" json:"meta,omitempty"`
}

// Component represents a UI component in the config tree.
type Component struct {
	ID        string                 `yaml:"id,omitempty" json:"id,omitempty"`
	Type      string                 `yaml:"type" json:"type"`
	Props     map[string]interface{} `yaml:"props,omitempty" json:"props,omitempty"`
	Actions   map[string]Action      `yaml:"actions,omitempty" json:"actions,omitempty"`
	Children  []Component            `yaml:"children,omitempty" json:"children,omitempty"`
	Shortcuts []Shortcut             `yaml:"shortcuts,omitempty" json:"shortcuts,omitempty"`
	Condition *Condition             `yaml:"condition,omitempty" json:"condition,omitempty"`
}

// Action defines an event handler.
type Action struct {
	Type       string                 `yaml:"type" json:"type"`
	Page       string                 `yaml:"page,omitempty" json:"page,omitempty"`
	URL        string                 `yaml:"url,omitempty" json:"url,omitempty"`
	Method     string                 `yaml:"method,omitempty" json:"method,omitempty"`
	Title      string                 `yaml:"title,omitempty" json:"title,omitempty"`
	Size       string                 `yaml:"size,omitempty" json:"size,omitempty"`
	Side       string                 `yaml:"side,omitempty" json:"side,omitempty"`
	Fields     []FormField            `yaml:"fields,omitempty" json:"fields,omitempty"`
	Submit     *SubmitConfig          `yaml:"submit,omitempty" json:"submit,omitempty"`
	Refresh    string                 `yaml:"refresh,omitempty" json:"refresh,omitempty"`
	Params     map[string]string      `yaml:"params,omitempty" json:"params,omitempty"`
	Event      string                 `yaml:"event,omitempty" json:"event,omitempty"`
	Payload    map[string]interface{} `yaml:"payload,omitempty" json:"payload,omitempty"`
	Message    string                 `yaml:"message,omitempty" json:"message,omitempty"`
	OnConfirm  *Action                `yaml:"onConfirm,omitempty" json:"onConfirm,omitempty"`
	OnSuccess  *Action                `yaml:"onSuccess,omitempty" json:"onSuccess,omitempty"`
	Target     string                 `yaml:"target,omitempty" json:"target,omitempty"`
	Datasource string                 `yaml:"datasource,omitempty" json:"datasource,omitempty"`
	Field      string                 `yaml:"field,omitempty" json:"field,omitempty"`
	Value      string                 `yaml:"value,omitempty" json:"value,omitempty"`
}

// FormField defines a field in a modal/form.
type FormField struct {
	Name        string         `yaml:"name" json:"name"`
	Label       string         `yaml:"label" json:"label"`
	Type        string         `yaml:"type" json:"type"`
	Required    bool           `yaml:"required,omitempty" json:"required,omitempty"`
	Placeholder string         `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
	Default     interface{}    `yaml:"default,omitempty" json:"default,omitempty"`
	Options     []SelectOption `yaml:"options,omitempty" json:"options,omitempty"`
	Validation  *Validation    `yaml:"validation,omitempty" json:"validation,omitempty"`
	Rows        int            `yaml:"rows,omitempty" json:"rows,omitempty"`
	MaxLength   int            `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Min         interface{}    `yaml:"min,omitempty" json:"min,omitempty"`
	Max         interface{}    `yaml:"max,omitempty" json:"max,omitempty"`
}

// SelectOption represents a value/label pair for select fields.
type SelectOption struct {
	Value string `yaml:"value" json:"value"`
	Label string `yaml:"label" json:"label"`
}

// Validation defines field-level validation rules.
type Validation struct {
	Pattern string `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Message string `yaml:"message,omitempty" json:"message,omitempty"`
}

// SubmitConfig defines how a form submits data.
type SubmitConfig struct {
	DataSource string `yaml:"datasource" json:"datasource"`
	Method     string `yaml:"method" json:"method"`
}

// Shortcut defines a keyboard shortcut binding.
type Shortcut struct {
	Key         string `yaml:"key" json:"key"`
	Action      Action `yaml:"action" json:"action"`
	When        string `yaml:"when,omitempty" json:"when,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Global      bool   `yaml:"global,omitempty" json:"global,omitempty"`
}

// Condition for conditional rendering.
type Condition struct {
	Field string      `yaml:"field" json:"field"`
	Op    string      `yaml:"op" json:"op"`
	Value interface{} `yaml:"value" json:"value"`
}

// DataSourceRef declares a datasource dependency in a page.
type DataSourceRef struct {
	Alias        string                 `yaml:"alias" json:"alias"`
	Capabilities []string               `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`
	Params       map[string]interface{} `yaml:"params,omitempty" json:"params,omitempty"`
}

// ThemeRef references a theme with optional overrides.
type ThemeRef struct {
	Extends string                 `yaml:"extends,omitempty" json:"extends,omitempty"`
	Tokens  map[string]interface{} `yaml:"tokens,omitempty" json:"tokens,omitempty"`
}

// PageMeta contains page metadata.
type PageMeta struct {
	Guards  []string `yaml:"guards,omitempty" json:"guards,omitempty"`
	Tags    []string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Version int      `yaml:"version,omitempty" json:"version,omitempty"`
}

// ComponentRegistry is an interface for checking component types during validation.
type ComponentRegistry interface {
	Has(componentType string) bool
}

// PropInfo describes a component prop for deep validation.
type PropInfo struct {
	Type     string
	Required bool
	Enum     []string
}

// ComponentSchemaInfo provides detailed schema info for deep validation.
type ComponentSchemaInfo struct {
	Props map[string]PropInfo
}

// ComponentSchemaProvider extends ComponentRegistry with schema lookup.
type ComponentSchemaProvider interface {
	ComponentRegistry
	GetSchema(componentType string) (*ComponentSchemaInfo, bool)
}

// AppConfig represents the .sigil/app.yaml application configuration.
//
// Providers live on each ModuleConfig (since sprint 10 phase 3.5) — there is
// no app-level fallback. Each module declares the providers wrapping its own
// layout/routes. See ModuleConfig.Providers.
type AppConfig struct {
	Sigil       string `yaml:"sigil" json:"sigil"`
	Kind        string `yaml:"kind" json:"kind"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// TargetMode selects the output shape for renderers that support multiple
	// frontend targets (currently react-shadcn). Valid values: "spa", "app-router".
	// Empty / unset is treated as "app-router" (the original shape) so existing
	// configs regenerate byte-identically.
	TargetMode string         `yaml:"target_mode,omitempty" json:"target_mode,omitempty"`
	API        *APIConfig     `yaml:"api,omitempty" json:"api,omitempty"`
	Modules    []ModuleConfig `yaml:"modules" json:"modules"`
	Features   *AppFeatures   `yaml:"features,omitempty" json:"features,omitempty"`
	Actions    []AppAction    `yaml:"actions,omitempty" json:"actions,omitempty"`
}

// EffectiveTargetMode returns the AppConfig's target mode, defaulting to
// "app-router" when unset. Callers should use this rather than reading
// TargetMode directly so an absent flag stays equivalent to today's output.
func (a *AppConfig) EffectiveTargetMode() string {
	if a == nil || a.TargetMode == "" {
		return "app-router"
	}
	return a.TargetMode
}

// APIConfig defines the API client generation settings.
type APIConfig struct {
	BaseURLEnv     string `yaml:"base_url_env" json:"base_url_env"`
	BaseURLDefault string `yaml:"base_url_default" json:"base_url_default"`
	Prefix         string `yaml:"prefix" json:"prefix"`
}

// ModuleConfig defines a module within an app.
//
// Providers are declared per-module (since sprint 10 phase 3.5): each module
// wraps its own layout/routes with the providers it lists, in declared order
// (first = outermost). Two modules may share the same custom provider
// declaration — the renderer copies each provider's source files once per
// app, deduplicated by source path.
type ModuleConfig struct {
	ID         string           `yaml:"id" json:"id"`
	Shell      string           `yaml:"shell" json:"shell"`
	RouteGroup string           `yaml:"route_group" json:"route_group"`
	Pages      []string         `yaml:"pages" json:"pages"`
	Providers  []ProviderConfig `yaml:"providers,omitempty" json:"providers,omitempty"`
}

// AppFeatures controls app-wide feature flags.
type AppFeatures struct {
	ThemeToggle    bool `yaml:"theme_toggle" json:"theme_toggle"`
	CommandPalette bool `yaml:"command_palette" json:"command_palette"`
	MobileSidebar  bool `yaml:"mobile_sidebar" json:"mobile_sidebar"`
}

// ProviderConfig defines a context provider to generate and wire into the layout.
//
// Two shapes are supported:
//
//  1. Built-in (the historical form): set Datasource/TrackField/LabelField to
//     have Sigil generate a data-driven context provider against the named
//     datasource. Optional UI bits (Position/UIType/Placeholder) wire the
//     provider into the layout topbar.
//
//  2. Custom (since sprint 10 phase 2): set Source.Component to point at a
//     hand-written provider .tsx file under .sigil/providers/. Sigil copies
//     the source files into the output's `lib/` directory and wraps the
//     layout with <Source.Export>. Additional sidecar files (e.g. an SSE
//     bridge) can be listed in Source.Includes, and any of those whose
//     exported component name appears in Mounts will be rendered as a
//     sibling of `{children}` inside the provider scope.
//
// A ProviderConfig is "custom" when Source.Component is non-empty.
// Datasource/TrackField/LabelField are ignored for custom providers; Mounts
// is ignored for built-in providers.
type ProviderConfig struct {
	ID          string          `yaml:"id" json:"id"`
	Datasource  string          `yaml:"datasource,omitempty" json:"datasource,omitempty"`
	TrackField  string          `yaml:"track_field,omitempty" json:"track_field,omitempty"`
	LabelField  string          `yaml:"label_field,omitempty" json:"label_field,omitempty"`
	Default     string          `yaml:"default,omitempty" json:"default,omitempty"`
	Placeholder string          `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
	Position    string          `yaml:"position,omitempty" json:"position,omitempty"` // "topbar" or ""
	UIType      string          `yaml:"ui_type,omitempty" json:"ui_type,omitempty"`   // "select" or ""
	Source      *ProviderSource `yaml:"source,omitempty" json:"source,omitempty"`
	Mounts      []string        `yaml:"mounts,omitempty" json:"mounts,omitempty"`
}

// ProviderSource defines the source files for a custom (hand-written) provider.
// Paths are relative to the .sigil/ directory (typically under
// `.sigil/providers/`). Component is the file that exports the provider
// component referenced by Export; Includes lists additional files copied to
// the output as sidecars (e.g. an SSE bridge component, a shared types file).
type ProviderSource struct {
	Component string   `yaml:"component" json:"component"`                   // main .tsx file (required for custom)
	Export    string   `yaml:"export,omitempty" json:"export,omitempty"`     // exported provider component name (default: pascalCase(ID)+"Provider")
	Includes  []string `yaml:"includes,omitempty" json:"includes,omitempty"` // additional files to copy
}

// IsCustom returns true if this provider has a source file and should be
// copied into the output rather than generated from a datasource template.
func (p *ProviderConfig) IsCustom() bool {
	return p != nil && p.Source != nil && p.Source.Component != ""
}

// ExportName returns the React component name to import from the provider's
// source file. Defaults to PascalCase(ID) + "Provider" when Source.Export is
// unset.
func (p *ProviderConfig) ExportName() string {
	if p == nil {
		return ""
	}
	if p.Source != nil && p.Source.Export != "" {
		return p.Source.Export
	}
	return pascalCaseID(p.ID) + "Provider"
}

// pascalCaseID converts a kebab/snake-case id to PascalCase. Duplicated from
// the renderer's helper so config can compute the default export name without
// importing the renderer package.
func pascalCaseID(s string) string {
	parts := splitOnDelims(s)
	var b []byte
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		b = append(b, upperFirst(p)...)
	}
	return string(b)
}

func splitOnDelims(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '-' || s[i] == '_' {
			if i > start {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 'a' - 'A'
	}
	return string(b)
}

// AppAction defines a quick action available in command palette and menus.
type AppAction struct {
	Label string `yaml:"label" json:"label"`
	Icon  string `yaml:"icon" json:"icon"`
	Page  string `yaml:"page" json:"page"`
	Query string `yaml:"query,omitempty" json:"query,omitempty"`
}

// ProjectConfig represents the .sigil/sigil.yaml project configuration.
type ProjectConfig struct {
	Version  string `yaml:"version" json:"version"`
	Name     string `yaml:"name" json:"name"`
	Defaults struct {
		Theme    string `yaml:"theme" json:"theme"`
		Renderer string `yaml:"renderer" json:"renderer"`
		Output   string `yaml:"output" json:"output"`
	} `yaml:"defaults" json:"defaults"`
	Components struct {
		Builtin   bool   `yaml:"builtin" json:"builtin"`
		CustomDir string `yaml:"custom_dir" json:"custom_dir"`
	} `yaml:"components" json:"components"`
	DataSources struct {
		Dir string `yaml:"dir" json:"dir"`
	} `yaml:"datasources" json:"datasources"`
	Themes struct {
		Dir string `yaml:"dir" json:"dir"`
	} `yaml:"themes" json:"themes"`
	Generation struct {
		Clean     bool   `yaml:"clean" json:"clean"`
		GoPackage string `yaml:"go_package,omitempty" json:"go_package,omitempty"`
		GoModule  string `yaml:"go_module,omitempty" json:"go_module,omitempty"`
	} `yaml:"generation" json:"generation"`
}
