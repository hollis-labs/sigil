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
}

// RenderContext provides everything a renderer needs to generate code for a page.
type RenderContext struct {
	Page          *config.Page
	Theme         *ThemeConfig
	DataSources   map[string]*DataSourceManifest
	Registry      *components.Registry
	ProjectConfig *config.ProjectConfig
	GoModule      string
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
