package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrispian/sigil/internal/components"
	"github.com/chrispian/sigil/internal/config"
	"gopkg.in/yaml.v3"
)

// GenerateConfig holds options for the generate command.
type GenerateConfig struct {
	Target    string   // Renderer name (e.g., "go-templ")
	OutputDir string   // Output directory
	SigilDir  string   // .sigil directory (default: ".sigil")
	Pages     []string // Specific pages to generate (empty = all)
	Theme     string   // Theme name override (empty = project default)
	Clean     bool     // Remove output dir before generating
	DryRun    bool     // List files without writing
	GoModule  string   // Go module path for imports
}

// GenerateResult contains the results of a generation run.
type GenerateResult struct {
	Files     []OutputFile
	PageCount int
}

// renderers is the registry of available renderers.
var renderers = map[string]Renderer{}

// RegisterRenderer adds a renderer to the registry.
func RegisterRenderer(r Renderer) {
	renderers[r.Name()] = r
}

// GetRenderer returns a registered renderer by name.
func GetRenderer(name string) (Renderer, bool) {
	r, ok := renderers[name]
	return r, ok
}

// AvailableRenderers returns the names of all registered renderers.
func AvailableRenderers() []string {
	var names []string
	for name := range renderers {
		names = append(names, name)
	}
	return names
}

// Generate processes all pages and produces output files.
func Generate(cfg GenerateConfig) (*GenerateResult, error) {
	r, ok := GetRenderer(cfg.Target)
	if !ok {
		return nil, fmt.Errorf("unknown renderer %q (available: %s)", cfg.Target, strings.Join(AvailableRenderers(), ", "))
	}

	sigilDir := cfg.SigilDir
	if sigilDir == "" {
		sigilDir = ".sigil"
	}

	// 1. Load project config
	projConfig, err := loadProjectConfig(filepath.Join(sigilDir, "sigil.yaml"))
	if err != nil {
		return nil, fmt.Errorf("loading project config: %w", err)
	}

	// 2. Load component registry
	registry := components.NewDefaultRegistry()

	// 3. Load and resolve theme
	themeName := cfg.Theme
	if themeName == "" {
		themeName = projConfig.Defaults.Theme
	}
	if themeName == "" {
		themeName = "default"
	}
	theme, err := LoadTheme(filepath.Join(sigilDir, "themes"), themeName)
	if err != nil {
		return nil, fmt.Errorf("loading theme %q: %w", themeName, err)
	}

	// 4. Load datasource manifests
	dataSources, err := loadDataSources(filepath.Join(sigilDir, "datasources"))
	if err != nil {
		return nil, fmt.Errorf("loading datasources: %w", err)
	}

	// 5. Load and validate pages
	pages, err := loadPages(filepath.Join(sigilDir, "pages"), cfg.Pages)
	if err != nil {
		return nil, fmt.Errorf("loading pages: %w", err)
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("no pages found in %s", filepath.Join(sigilDir, "pages"))
	}

	// Validate pages
	for _, page := range pages {
		result := config.Validate(page, registry)
		if !result.Valid {
			var msgs []string
			for _, e := range result.Errors {
				msgs = append(msgs, e.String())
			}
			return nil, fmt.Errorf("validation failed for page %q:\n  %s", page.ID, strings.Join(msgs, "\n  "))
		}
	}

	// Resolve Go module path
	goModule := cfg.GoModule
	if goModule == "" {
		goModule = projConfig.Generation.GoModule
	}

	// 6. Render each page
	var allFiles []OutputFile
	usedTypes := map[string]bool{}

	for _, page := range pages {
		ctx := &RenderContext{
			Page:          page,
			Theme:         theme,
			DataSources:   dataSources,
			Registry:      registry,
			ProjectConfig: projConfig,
			GoModule:      goModule,
		}
		files, err := r.Render(ctx)
		if err != nil {
			return nil, fmt.Errorf("rendering page %q: %w", page.ID, err)
		}
		allFiles = append(allFiles, files...)

		// Track used component types
		collectUsedTypes(&page.Layout, usedTypes)
	}

	// 7. Generate shared components (once)
	var usedList []string
	for t := range usedTypes {
		usedList = append(usedList, t)
	}
	sharedFiles, err := r.SharedComponents(usedList)
	if err != nil {
		return nil, fmt.Errorf("generating shared components: %w", err)
	}
	allFiles = append(allFiles, sharedFiles...)

	// 8. Generate theme CSS (once)
	themeFiles, err := r.RenderTheme(theme)
	if err != nil {
		return nil, fmt.Errorf("generating theme: %w", err)
	}
	allFiles = append(allFiles, themeFiles...)

	// 9. Generate handler stubs (once per datasource)
	for _, ds := range dataSources {
		stubFiles, err := r.RenderDataSourceStubs(ds)
		if err != nil {
			return nil, fmt.Errorf("generating datasource stubs for %q: %w", ds.Alias, err)
		}
		allFiles = append(allFiles, stubFiles...)
	}

	result := &GenerateResult{
		Files:     allFiles,
		PageCount: len(pages),
	}

	// Dry run — just return the list
	if cfg.DryRun {
		return result, nil
	}

	// 10. Clean output dir if requested
	if cfg.Clean {
		if err := os.RemoveAll(cfg.OutputDir); err != nil {
			return nil, fmt.Errorf("cleaning output dir: %w", err)
		}
	}

	// 11. Write all files
	for _, f := range allFiles {
		outPath := filepath.Join(cfg.OutputDir, f.Path)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return nil, fmt.Errorf("creating directory for %s: %w", f.Path, err)
		}
		mode := f.Mode
		if mode == 0 {
			mode = 0644
		}
		if err := os.WriteFile(outPath, f.Content, mode); err != nil {
			return nil, fmt.Errorf("writing %s: %w", f.Path, err)
		}
	}

	return result, nil
}

func loadProjectConfig(path string) (*config.ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg config.ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadTheme loads a theme by name, resolving inheritance via the extends field.
func LoadTheme(themesDir, name string) (*ThemeConfig, error) {
	path := filepath.Join(themesDir, name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var theme ThemeConfig
	if err := yaml.Unmarshal(data, &theme); err != nil {
		return nil, err
	}

	// Resolve inheritance
	if theme.Extends != "" {
		base, err := LoadTheme(themesDir, theme.Extends)
		if err != nil {
			return nil, fmt.Errorf("loading base theme %q: %w", theme.Extends, err)
		}
		mergeTheme(base, &theme)
		return base, nil
	}

	return &theme, nil
}

func mergeTheme(base, child *ThemeConfig) {
	if base.Tokens == nil {
		base.Tokens = map[string]map[string]interface{}{}
	}
	for category, tokens := range child.Tokens {
		if base.Tokens[category] == nil {
			base.Tokens[category] = map[string]interface{}{}
		}
		for k, v := range tokens {
			base.Tokens[category][k] = v
		}
	}
	if base.Variants == nil {
		base.Variants = map[string]map[string]interface{}{}
	}
	for comp, variants := range child.Variants {
		if base.Variants[comp] == nil {
			base.Variants[comp] = map[string]interface{}{}
		}
		for k, v := range variants {
			base.Variants[comp][k] = v
		}
	}
}

func loadDataSources(dir string) (map[string]*DataSourceManifest, error) {
	result := map[string]*DataSourceManifest{}
	pattern := filepath.Join(dir, "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		var ds DataSourceManifest
		if err := yaml.Unmarshal(data, &ds); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}
		if ds.Alias == "" {
			ds.Alias = strings.TrimSuffix(filepath.Base(path), ".yaml")
		}
		result[ds.Alias] = &ds
	}
	return result, nil
}

func loadPages(pagesDir string, filter []string) ([]*config.Page, error) {
	pattern := filepath.Join(pagesDir, "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	filterSet := map[string]bool{}
	for _, f := range filter {
		filterSet[f] = true
	}

	var pages []*config.Page
	for _, path := range matches {
		page, err := config.ParseFile(path)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}
		if len(filterSet) > 0 && !filterSet[page.ID] {
			continue
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func collectUsedTypes(c *config.Component, types map[string]bool) {
	if c.Type != "" {
		types[c.Type] = true
	}
	for i := range c.Children {
		collectUsedTypes(&c.Children[i], types)
	}
}
