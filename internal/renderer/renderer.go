// Package renderer defines the renderer interface and implementations.
package renderer

import (
	"os"

	"github.com/chrispian/sigil/internal/components"
	"github.com/chrispian/sigil/internal/config"
)

// Renderer generates framework-specific code from Sigil configs.
type Renderer interface {
	// Name returns the renderer identifier (e.g., "go-templ").
	Name() string

	// Render generates code for a single page config.
	Render(ctx *RenderContext) ([]OutputFile, error)

	// RenderTheme generates theme-specific files (CSS, config).
	RenderTheme(theme *ThemeConfig) ([]OutputFile, error)

	// RenderDataSourceStubs generates handler stubs for datasources.
	RenderDataSourceStubs(ds *DataSourceManifest) ([]OutputFile, error)

	// SharedComponents returns shared/reusable component files.
	SharedComponents(usedTypes []string) ([]OutputFile, error)

	// RenderLayout generates the layout wrapper for a module (shell + nav + chrome).
	RenderLayout(ctx *LayoutContext) ([]OutputFile, error)

	// RenderAPIClient generates a typed API client from app config.
	// Receives the full *AppConfig so renderers can branch on target mode
	// (e.g., env var access pattern). app.API may still be nil — callers
	// already short-circuit in that case, but implementations should handle it.
	RenderAPIClient(app *config.AppConfig) ([]OutputFile, error)

	// RenderProviders generates React Context providers from app config.
	// Receives the full *AppConfig (so renderers can branch on target mode if
	// needed) and the path to the .sigil/ directory (so custom providers can
	// have their source files copied into the output).
	//
	// Since sprint 10 phase 3.5, providers live per-module. The engine
	// aggregates each module's providers and dedups custom-provider source
	// file copies across modules (two modules declaring the same custom
	// provider produce a single set of `lib/<basename>` files). The signature
	// stays *AppConfig so renderers see every module's providers in one call.
	RenderProviders(app *config.AppConfig, sigilDir string) ([]OutputFile, error)
}

// RenderContext provides everything a renderer needs to generate code for a page.
type RenderContext struct {
	Page          *config.Page
	Theme         *ThemeConfig
	DataSources   map[string]*DataSourceManifest
	Registry      *components.Registry
	ProjectConfig *config.ProjectConfig
	GoModule      string
	SigilDir      string // path to .sigil directory (for locating custom component sources)
	// TargetMode is the AppConfig.EffectiveTargetMode() at generation time
	// (e.g., "app-router" or "spa"). Renderers that emit framework-specific
	// imports/calls branch on this. Empty means "app-router" for safety.
	TargetMode string
	// Modules is the full app's module list, used by SPA renderers to
	// resolve cross-module route prefixes when emitting navigation calls.
	// Nil/empty is fine — renderers fall back to the legacy single-module
	// path layout (no prefix). App Router renderers ignore this field.
	Modules []*config.ModuleConfig
}

// LayoutContext provides everything a renderer needs to generate a module layout.
type LayoutContext struct {
	Module    *config.ModuleConfig
	Shell     *config.Page
	AppConfig *config.AppConfig
	Theme     *ThemeConfig
	Pages     []*config.Page // All pages in this module (for nav extraction)
	// AllModules and AllPages let SPA renderers emit a single App.tsx that
	// roots every module at its own URL prefix (derived from
	// ModuleConfig.RouteGroup). For App Router renderers these are ignored
	// (each module already gets its own app/(group)/layout.tsx).
	AllModules []*config.ModuleConfig
	AllPages   []*config.Page          // shells included for parity, but pages are filtered to non-shells by the engine
	AllShells  map[string]*config.Page // moduleID → shell page
}

// OutputFile represents a file to be written to disk.
type OutputFile struct {
	Path    string      // Relative to output directory
	Content []byte      // File content
	Mode    os.FileMode // File permissions (default: 0644)
}

// ThemeConfig represents a parsed theme file.
type ThemeConfig struct {
	Name        string                            `yaml:"name"`
	Description string                            `yaml:"description"`
	Extends     string                            `yaml:"extends,omitempty"`
	Tokens      map[string]map[string]interface{} `yaml:"tokens"`
	Variants    map[string]map[string]interface{} `yaml:"variants,omitempty"`
}

// DataSourceManifest represents a parsed datasource manifest.
type DataSourceManifest struct {
	Alias        string                 `yaml:"alias"`
	Description  string                 `yaml:"description"`
	Capabilities []string               `yaml:"capabilities"`
	Fields       []DataSourceField      `yaml:"fields"`
	Relations    []DataSourceRelation   `yaml:"relations,omitempty"`
	Endpoints    map[string]string      `yaml:"endpoints"`
	Defaults     map[string]interface{} `yaml:"defaults,omitempty"`
}

// DataSourceField describes a field in a datasource manifest.
type DataSourceField struct {
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Label      string      `yaml:"label,omitempty"`
	Primary    bool        `yaml:"primary,omitempty"`
	Required   bool        `yaml:"required,omitempty"`
	Readonly   bool        `yaml:"readonly,omitempty"`
	Hidden     bool        `yaml:"hidden,omitempty"`
	Searchable bool        `yaml:"searchable,omitempty"`
	Sortable   bool        `yaml:"sortable,omitempty"`
	Filterable bool        `yaml:"filterable,omitempty"`
	Default    interface{} `yaml:"default,omitempty"`
	Values     []string    `yaml:"values,omitempty"`
}

// DataSourceRelation describes a relation to another datasource.
type DataSourceRelation struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Target     string `yaml:"target"`
	ForeignKey string `yaml:"foreignKey"`
	Label      string `yaml:"label,omitempty"`
}

// HasCapability checks if the datasource supports a given capability.
func (ds *DataSourceManifest) HasCapability(cap string) bool {
	for _, c := range ds.Capabilities {
		if c == cap {
			return true
		}
	}
	return false
}
