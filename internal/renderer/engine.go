package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hollis-labs/sigil/internal/components"
	"github.com/hollis-labs/sigil/internal/config"
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
	// TargetMode overrides AppConfig.TargetMode at generate-time. Valid
	// values: "", "spa", "app-router". Empty falls through to whatever the
	// app.yaml declares (or its default). Invalid values fail validation.
	TargetMode string
	// UIKit selects the shared component kit generated pages import from.
	// Valid values: "", "sysop". Empty (default) keeps the legacy behavior —
	// per-app hand-written shadcn re-exports under @/components/ui. "sysop"
	// makes the react-shadcn renderer import kit components from
	// @hollis-labs/sysop-ui. Ignored by renderers that don't support a kit.
	UIKit string
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

	// Hand the UI-kit selection to renderers that support one. RenderContext
	// and LayoutContext also carry UIKit, but the contextless interface methods
	// (RenderTheme/SharedComponents/RenderAPIClient/RenderProviders) need it set
	// on the renderer itself.
	if ka, ok := r.(UIKitAware); ok {
		ka.SetUIKit(cfg.UIKit)
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

	// 1.5. Load app config (optional — backward compatible)
	appConfig, err := loadAppConfig(filepath.Join(sigilDir, "app.yaml"))
	if err != nil {
		return nil, fmt.Errorf("loading app config: %w", err)
	}
	// CLI --target-mode override wins over app.yaml's target_mode. Applied
	// before validation so an invalid CLI value is caught by the same
	// validator that checks app.yaml values.
	if cfg.TargetMode != "" {
		if appConfig == nil {
			// Without an app.yaml the renderer has no modules to layout, so
			// target_mode is moot — but accept the flag for forward symmetry
			// (a CLI-only invocation that hits the renderer without app.yaml
			// is rare).
			appConfig = &config.AppConfig{TargetMode: cfg.TargetMode}
		} else {
			appConfig.TargetMode = cfg.TargetMode
		}
	}
	if appConfig != nil {
		result := config.ValidateAppConfig(appConfig)
		if !result.Valid {
			var msgs []string
			for _, e := range result.Errors {
				msgs = append(msgs, e.String())
			}
			return nil, fmt.Errorf("app config validation failed:\n  %s", strings.Join(msgs, "\n  "))
		}
	}

	// 2. Load component registry (built-in + custom)
	registry := components.NewDefaultRegistry()
	if projConfig.Components.CustomDir != "" {
		if loadErr := components.LoadCustomSchemas(registry, projConfig.Components.CustomDir); loadErr != nil {
			return nil, fmt.Errorf("loading custom components: %w", loadErr)
		}
	}

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

	targetMode := appConfig.EffectiveTargetMode() // safe on nil

	// Pre-build the module pointer list once so each page's RenderContext can
	// resolve cross-module routes in SPA mode. Filter to modules that have
	// at least one loaded (not --pages-filtered) page — modules with no
	// surviving pages have no routes emitted, so links into them would
	// 404; treating them as absent matches the layout filter below.
	var modulesForCtx []*config.ModuleConfig
	if appConfig != nil {
		loaded := map[string]bool{}
		for _, p := range pages {
			loaded[p.ID] = true
		}
		for i := range appConfig.Modules {
			mod := &appConfig.Modules[i]
			for _, mp := range mod.Pages {
				if loaded[mp] {
					modulesForCtx = append(modulesForCtx, mod)
					break
				}
			}
		}
	}

	for _, page := range pages {
		ctx := &RenderContext{
			Page:          page,
			Theme:         theme,
			DataSources:   dataSources,
			Registry:      registry,
			ProjectConfig: projConfig,
			GoModule:      goModule,
			SigilDir:      sigilDir,
			TargetMode:    targetMode,
			Modules:       modulesForCtx,
			UIKit:         cfg.UIKit,
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

	// 7.5. Copy custom component source files
	for typeName := range usedTypes {
		schema, ok := registry.Get(typeName)
		if !ok || !schema.IsCustom() {
			continue
		}
		src := schema.Source
		// Copy main component file
		srcPath := filepath.Join(sigilDir, src.Component)
		srcData, srcErr := os.ReadFile(srcPath) //nolint:gosec // Path is from validated schema config
		if srcErr != nil {
			return nil, fmt.Errorf("reading custom component source %q: %w", src.Component, srcErr)
		}
		allFiles = append(allFiles, OutputFile{Path: src.Component, Content: srcData})

		// Copy additional included files
		for _, inc := range src.Includes {
			incPath := filepath.Join(sigilDir, inc)
			incData, incErr := os.ReadFile(incPath) //nolint:gosec // Path is from validated schema config
			if incErr != nil {
				return nil, fmt.Errorf("reading custom component include %q: %w", inc, incErr)
			}
			allFiles = append(allFiles, OutputFile{Path: inc, Content: incData})
		}
	}

	// 8. Generate theme CSS (once)
	themeFiles, err := r.RenderTheme(theme)
	if err != nil {
		return nil, fmt.Errorf("generating theme: %w", err)
	}
	allFiles = append(allFiles, themeFiles...)

	// 9. Generate handler stubs.
	//
	// Phase 3.5: emit hooks only for datasources actually referenced by the
	// loaded pages (or by ANY page if there are no pages — preserves the
	// `sigil generate` ergonomic of "give me everything"). This prevents a
	// `--pages X,Y,Z` filter from leaking unrelated swr-using hook files
	// into the output, which then break `tsc` because swr isn't a
	// dependency of the target app.
	referencedDS := map[string]bool{}
	for _, page := range pages {
		for _, ref := range page.DataSources {
			referencedDS[ref.Alias] = true
		}
	}
	for _, ds := range dataSources {
		if len(referencedDS) > 0 && !referencedDS[ds.Alias] {
			continue
		}
		stubFiles, err := r.RenderDataSourceStubs(ds)
		if err != nil {
			return nil, fmt.Errorf("generating datasource stubs for %q: %w", ds.Alias, err)
		}
		allFiles = append(allFiles, stubFiles...)
	}

	// 10. Generate app layouts.
	//
	// App Router mode emits one layout.tsx per module (each in its own
	// app/(group)/ directory). SPA mode emits a single App.tsx for the whole
	// app — so we still loop here, but pass AllModules + AllShells into the
	// context and the renderer is expected to short-circuit on subsequent
	// modules. The first module's RenderLayout call sees the full picture
	// and emits the unified App.tsx + routes.tsx; subsequent module calls
	// return no files in SPA mode.
	if appConfig != nil {
		// Pre-load all shell pages so SPA mode has them at hand. Also drop
		// modules whose pages are all filtered out by --pages — leaving them
		// in would emit nav links to routes that don't exist in the routes
		// table, breaking navigation.
		loadedPages := map[string]*config.Page{}
		for _, p := range pages {
			loadedPages[p.ID] = p
		}
		allShells := map[string]*config.Page{}
		var allModulePages []*config.Page
		var activeModuleIdx []int
		for i := range appConfig.Modules {
			mod := &appConfig.Modules[i]
			hasActive := false
			for _, mp := range mod.Pages {
				if _, ok := loadedPages[mp]; ok {
					hasActive = true
					if p := loadedPages[mp]; p != nil {
						allModulePages = append(allModulePages, p)
					}
				}
			}
			if !hasActive {
				continue
			}
			activeModuleIdx = append(activeModuleIdx, i)
			shellPage, err := loadShellPage(filepath.Join(sigilDir, "pages"), mod.Shell)
			if err != nil {
				return nil, fmt.Errorf("loading shell page %q for module %q: %w", mod.Shell, mod.ID, err)
			}
			allShells[mod.ID] = shellPage
		}
		allModulePtrs := make([]*config.ModuleConfig, len(activeModuleIdx))
		for idx, i := range activeModuleIdx {
			allModulePtrs[idx] = &appConfig.Modules[i]
		}

		for _, i := range activeModuleIdx {
			mod := &appConfig.Modules[i]
			shellPage := allShells[mod.ID]
			// Collect module pages for nav extraction.
			var modulePages []*config.Page
			for _, mp := range mod.Pages {
				if p, ok := loadedPages[mp]; ok {
					modulePages = append(modulePages, p)
				}
			}
			layoutCtx := &LayoutContext{
				Module:     mod,
				Shell:      shellPage,
				AppConfig:  appConfig,
				Theme:      theme,
				Pages:      modulePages,
				AllModules: allModulePtrs,
				AllPages:   allModulePages,
				AllShells:  allShells,
				UIKit:      cfg.UIKit,
			}
			layoutFiles, err := r.RenderLayout(layoutCtx)
			if err != nil {
				return nil, fmt.Errorf("rendering layout for module %q: %w", mod.ID, err)
			}
			allFiles = append(allFiles, layoutFiles...)
		}

		// 11. Generate API client
		if appConfig.API != nil {
			apiFiles, err := r.RenderAPIClient(appConfig)
			if err != nil {
				return nil, fmt.Errorf("generating API client: %w", err)
			}
			allFiles = append(allFiles, apiFiles...)
		}

		// 12. Generate context providers.
		// Phase 3.5: providers live per-module. Trigger the renderer call if
		// any module declares at least one provider; the renderer handles
		// per-module rendering + cross-module dedup of custom-provider files.
		hasAnyProvider := false
		for _, m := range appConfig.Modules {
			if len(m.Providers) > 0 {
				hasAnyProvider = true
				break
			}
		}
		if hasAnyProvider {
			providerFiles, err := r.RenderProviders(appConfig, sigilDir)
			if err != nil {
				return nil, fmt.Errorf("generating providers: %w", err)
			}
			allFiles = append(allFiles, providerFiles...)
		}
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

func loadAppConfig(path string) (*config.AppConfig, error) {
	data, err := os.ReadFile(path) //nolint:gosec // Path is from CLI --config flag
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // app.yaml is optional
		}
		return nil, err
	}
	var cfg config.AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadShellPage(pagesDir, shellID string) (*config.Page, error) {
	path := filepath.Join(pagesDir, shellID+".yaml")
	return config.ParseFile(path)
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
		// Alias mappings: some YAML types render via a differently-named
		// custom component file. Phase 3.5 chart wiring: `type: chart` in
		// YAML emits `<SigilChart>` and reads from `.sigil/components/
		// sigil-chart.tsx`, so we need the engine's file-copy step to
		// know about it.
		for _, alias := range customComponentAliases[c.Type] {
			types[alias] = true
		}
	}
	for i := range c.Children {
		collectUsedTypes(&c.Children[i], types)
	}
}

// customComponentAliases maps a YAML component `type` to additional
// custom-component schemas the engine should treat as "used" so their
// source files are copied into the output.
var customComponentAliases = map[string][]string{
	"chart": {"sigil-chart"},
}
