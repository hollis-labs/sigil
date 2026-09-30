package reactshadcn

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hollis-labs/sigil/internal/config"
	"github.com/hollis-labs/sigil/internal/renderer"
)

// kit.go — sysop UI kit (@hollis-labs/sysop-ui) integration for the
// react-shadcn renderer. Everything here is gated on importTracker.kitMode()
// (the --ui-kit sysop flag); legacy per-app shadcn output is untouched.

// kitModule is the npm package generated pages import kit components from.
const kitModule = "@hollis-labs/sysop-ui"

// kitUIModules is the set of `@/components/ui/<suffix>` shadcn primitive
// modules that the sysop kit re-exports from its barrel. In kit mode,
// addShadcn redirects imports of these to kitModule. Primitives NOT listed
// here (select, card, label, tabs, separator, avatar, progress,
// dropdown-menu, sheet, accordion, alert, ...) have no kit equivalent and
// stay app-local under @/components/ui.
var kitUIModules = map[string]bool{
	"badge":       true,
	"button":      true,
	"command":     true,
	"dialog":      true,
	"input":       true,
	"input-group": true,
	"popover":     true,
	"scroll-area": true,
	"skeleton":    true,
	"sonner":      true,
	"table":       true,
	"textarea":    true,
	"tooltip":     true,
}

// kitRedirect maps a shadcn ui import module to kitModule when generation
// targets the sysop kit and the module is one the kit re-exports. Returns the
// module unchanged otherwise. The kit's domain components (DataTable,
// PageHeader, ...) have no @/components/ui path — those go through addKit.
func kitRedirect(module string) string {
	if suffix, ok := strings.CutPrefix(module, "@/components/ui/"); ok && kitUIModules[suffix] {
		return kitModule
	}
	return module
}

// addKit registers a named import from the sysop kit barrel (kitModule). Used
// for the kit's domain components — PageHeader, StatusBadge, DataTable,
// FilterBar, EmptyState, DetailDialog, NavRail, SummaryCards, etc. — which
// have no @/components/ui counterpart to redirect.
func (t *importTracker) addKit(name string) {
	t.addShadcn(name, kitModule)
}

// customKitComponents is the set of custom-component types that, in kit mode,
// are emitted as their @hollis-labs/sysop-ui equivalent instead of being
// copied in as hand-written app-local .tsx. This is the heart of FND-4 — the
// kit IS the interactive core, so Sigil wires it rather than regenerating it.
var customKitComponents = map[string]bool{
	"status-badge": true,
	"filter-bar":   true,
}

// renderKitCustomComponent emits the kit equivalent of a custom component.
// Returns false if the type has no kit mapping (caller falls back to the
// legacy copy-the-.tsx path).
func renderKitCustomComponent(buf *bytes.Buffer, c *config.Component, indent string, imports *importTracker, ctx *renderer.RenderContext) bool {
	switch c.Type {
	case "status-badge":
		imports.addKit("StatusBadge")
		// The kit's StatusBadge does its own theme-aware tone lookup, so the
		// legacy size/showLabel/dot props are dropped — only `status` carries.
		fmt.Fprintf(buf, "%s<StatusBadge status={%s} />\n", indent, kitStatusExpr(c.Props["status"], ctx))
		return true
	case "filter-bar":
		renderKitFilterBar(buf, c, indent, imports)
		return true
	}
	return false
}

// kitStatusExpr resolves a status prop value to a string JSX expression. A
// `{{...}}` template ref becomes the resolved datasource/item reference; a
// plain string becomes a quoted literal.
func kitStatusExpr(val interface{}, ctx *renderer.RenderContext) string {
	s, _ := val.(string)
	if strings.HasPrefix(s, "{{") && strings.HasSuffix(s, "}}") {
		ref := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, "{{"), "}}"))
		return "String(" + normalizeDatasourceRef(ref, ctx) + " ?? \"\")"
	}
	return fmt.Sprintf("%q", s)
}

// renderKitFilterBar emits the kit's <FilterBar> with <FilterCycleToggle>
// chips for a `filter-bar` custom component. Each `select`-type filter is
// registered against the datasource (registerFilter) so the page's data-table
// applies the matching `.filter()` chain — the filter bar and the table share
// one piece of state per field. Search is owned by the FilterBar (the kit's
// FilterBar has a built-in search input) and registered so the table consumes
// it too. `date-range` filters have no kit equivalent and are skipped.
func renderKitFilterBar(buf *bytes.Buffer, c *config.Component, indent string, imports *importTracker) {
	imports.addKit("FilterBar")
	datasource := getPropString(c.Props, "datasource", "")

	// Resolve select filters into (label, stateVar, setterVar, options).
	type cycleFilter struct {
		label, stateVar, setterVar string
		options                    []map[string]string // {value,label}
	}
	var filters []cycleFilter
	if raw, ok := c.Props["filters"].([]interface{}); ok {
		for _, f := range raw {
			fm, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			ftype, _ := fm["type"].(string)
			field := fmt.Sprintf("%v", fm["field"])
			if ftype != "select" || field == "" || datasource == "" {
				continue // skip date-range and malformed filters
			}
			stateVar, setterVar := registerFilter(imports, datasource, field)
			cf := cycleFilter{label: fmt.Sprintf("%v", fm["label"]), stateVar: stateVar, setterVar: setterVar}
			if opts, ok := fm["options"].([]interface{}); ok {
				for _, o := range opts {
					om, ok := o.(map[string]interface{})
					if !ok {
						continue
					}
					val := fmt.Sprintf("%v", om["value"])
					if val == "" {
						val = "all" // empty == "show all"; the table treats "all" as unfiltered
					}
					cf.options = append(cf.options, map[string]string{
						"value": val,
						"label": fmt.Sprintf("%v", om["label"]),
					})
				}
			}
			// The shared filter state defaults to "all" — guarantee it's a
			// valid cycle position so the toggle renders and can clear itself.
			hasAll := false
			for _, o := range cf.options {
				if o["value"] == "all" {
					hasAll = true
					break
				}
			}
			if !hasAll {
				cf.options = append([]map[string]string{{"value": "all", "label": "All"}}, cf.options...)
			}
			filters = append(filters, cf)
		}
	}

	// Search is owned by the kit FilterBar.
	searchVar, searchSetter := registerSearch(imports, datasource)

	// activeFilterCount = number of filters not on "all".
	activeExpr := "0"
	if len(filters) > 0 {
		var vars []string
		for _, f := range filters {
			vars = append(vars, f.stateVar)
		}
		activeExpr = fmt.Sprintf("[%s].filter((v) => v !== \"all\").length", strings.Join(vars, ", "))
	}

	// onClear resets search + every filter.
	var resets []string
	resets = append(resets, fmt.Sprintf("%s(\"\")", searchSetter))
	for _, f := range filters {
		resets = append(resets, fmt.Sprintf("%s(\"all\")", f.setterVar))
	}

	fmt.Fprintf(buf, "%s<FilterBar\n", indent)
	fmt.Fprintf(buf, "%s  searchQuery={%s}\n", indent, searchVar)
	fmt.Fprintf(buf, "%s  onSearchChange={%s}\n", indent, searchSetter)
	fmt.Fprintf(buf, "%s  searchPlaceholder=\"Search…\"\n", indent)
	fmt.Fprintf(buf, "%s  activeFilterCount={%s}\n", indent, activeExpr)
	fmt.Fprintf(buf, "%s  onClear={() => { %s; }}\n", indent, strings.Join(resets, "; "))
	fmt.Fprintf(buf, "%s>\n", indent)
	if len(filters) > 0 {
		imports.addKit("FilterCycleToggle")
	}
	for _, f := range filters {
		// Explicit <FilterCycleToggle<string>> so T=string — the shared filter
		// state var is a plain string, not the option-value literal union.
		fmt.Fprintf(buf, "%s  <FilterCycleToggle<string>\n", indent)
		fmt.Fprintf(buf, "%s    ariaLabel=%q\n", indent, f.label)
		fmt.Fprintf(buf, "%s    value={%s}\n", indent, f.stateVar)
		fmt.Fprintf(buf, "%s    onChange={%s}\n", indent, f.setterVar)
		fmt.Fprintf(buf, "%s    options={[\n", indent)
		for _, o := range f.options {
			fmt.Fprintf(buf, "%s      { value: %q, label: %q },\n", indent, o["value"], o["label"])
		}
		fmt.Fprintf(buf, "%s    ]}\n", indent)
		fmt.Fprintf(buf, "%s  />\n", indent)
	}
	fmt.Fprintf(buf, "%s</FilterBar>\n", indent)
}

// columnRendersBadge reports whether a data-table column's `render` value asks
// for a badge. Specs use two forms: a bare string (`render: badge`) and a
// nested map (`render: { type: badge, variants: {...} }`).
func columnRendersBadge(render interface{}) bool {
	switch v := render.(type) {
	case string:
		return v == "badge"
	case map[string]interface{}:
		t, _ := v["type"].(string)
		return t == "badge"
	default:
		return false
	}
}

// primaryField returns the datasource's primary-key field name, or "id" when
// the datasource is unknown or declares no primary field.
func primaryField(ctx *renderer.RenderContext, datasource string) string {
	if ctx != nil {
		if ds, ok := ctx.DataSources[datasource]; ok {
			for _, f := range ds.Fields {
				if f.Primary {
					return f.Name
				}
			}
		}
	}
	return "id"
}

// renderKitLayoutSPA emits the SPA App.tsx as a thin composition of the sysop
// kit's shell — NavRail + PageHeader + ThemeSwitcher — matching the shape the
// folio `sysop-ui` preset scaffolds. No next-themes, no app-local
// @/components/ui/* chrome: the kit IS the shell, Sigil just composes it.
// react-router routing and per-module provider wrapping are retained (Sigil
// apps are multi-page/multi-module; the folio starter is single-page).
func renderKitLayoutSPA(ctx *renderer.LayoutContext) ([]renderer.OutputFile, error) {
	app := ctx.AppConfig
	mod := ctx.Module

	// Only the first module generates App.tsx (the unified shell covers all).
	if len(ctx.AllModules) > 0 && ctx.AllModules[0] != nil && ctx.AllModules[0].ID != mod.ID {
		return nil, nil
	}

	modules := ctx.AllModules
	if len(modules) == 0 {
		modules = []*config.ModuleConfig{mod}
	}
	shellByID := ctx.AllShells
	if shellByID == nil {
		shellByID = map[string]*config.Page{mod.ID: ctx.Shell}
	}

	// Nav items: walk every module's shell with a SPA-aware route resolver.
	var groups []navGroup
	for _, m := range modules {
		if m == nil {
			continue
		}
		shellPage := shellByID[m.ID]
		if shellPage == nil {
			continue
		}
		groups = append(groups, extractNavGroups(&shellPage.Layout, spaRouteResolver(m.ID, modules))...)
	}
	var navItems []navItem
	for _, g := range groups {
		navItems = append(navItems, g.Items...)
	}

	// Brand from the primary module's shell.
	primaryShell := shellByID[mod.ID]
	if primaryShell == nil {
		primaryShell = ctx.Shell
	}
	brandIcon, brandTitle := "", ""
	if primaryShell != nil {
		brandIcon, brandTitle = extractBrand(&primaryShell.Layout)
	}
	if brandTitle == "" {
		brandTitle = app.Name
	}
	brandIconComp := lucideComponentName(brandIcon)
	if brandIcon == "" {
		brandIconComp = "Hexagon"
	}

	// Icon set: nav icons + brand.
	iconSet := map[string]bool{brandIconComp: true}
	for _, item := range navItems {
		iconSet[lucideComponentName(item.Icon)] = true
	}
	var iconImports []string
	for name := range iconSet {
		iconImports = append(iconImports, name)
	}
	sortStrings(iconImports)

	uniqueProviders := dedupProvidersAcrossModules(modules)
	anyProviders := len(uniqueProviders) > 0

	var buf bytes.Buffer
	buf.WriteString("// App.tsx\n")
	buf.WriteString("// Generated by Sigil (sysop kit mode) — do not edit manually\n\n")
	buf.WriteString("import { Suspense } from \"react\";\n")
	buf.WriteString("import {\n  BrowserRouter,\n  Routes,\n  Route,\n  Navigate,\n  useLocation,\n  useNavigate,\n} from \"react-router-dom\";\n")
	if len(iconImports) > 0 {
		buf.WriteString("import {\n")
		for _, name := range iconImports {
			fmt.Fprintf(&buf, "  %s,\n", name)
		}
		buf.WriteString("} from \"lucide-react\";\n")
	}
	buf.WriteString("import {\n  NavRail,\n  PageHeader,\n  ThemeSwitcher,\n  Toaster,\n  applyTheme,\n  getInitialTheme,\n  type NavRailItem,\n} from \"@hollis-labs/sysop-ui\";\n")

	// Route helpers — per-module when providers wrap groups, else flat getRoutes.
	if anyProviders {
		buf.WriteString("import {")
		for _, m := range modules {
			if m == nil || len(pagesForModule(allModulePagesFromCtx(ctx), m.ID)) == 0 {
				continue
			}
			fmt.Fprintf(&buf, " getModuleRoutes_%s,", toPascalCase(m.ID))
		}
		buf.WriteString(" } from \"@/routes\";\n")
	} else {
		buf.WriteString("import { getRoutes } from \"@/routes\";\n")
	}

	// Provider imports (custom app providers — orthogonal to the kit).
	for _, p := range uniqueProviders {
		wrapName := providerWrapName(p)
		importPath := providerImportPath(p)
		hookName, hookPath := providerHookImport(p)
		if hookName != "" && hookPath == importPath {
			fmt.Fprintf(&buf, "import { %s, %s } from %q;\n", wrapName, hookName, importPath)
		} else {
			fmt.Fprintf(&buf, "import { %s } from %q;\n", wrapName, importPath)
			if hookName != "" {
				fmt.Fprintf(&buf, "import { %s } from %q;\n", hookName, hookPath)
			}
		}
		for _, m := range providerMounts(p) {
			fmt.Fprintf(&buf, "import { %s } from %q;\n", m.Name, m.Path)
		}
	}
	if anyProviders {
		buf.WriteString("import { Outlet } from \"react-router-dom\";\n")
	}
	buf.WriteString("\n")

	buf.WriteString("// Apply the persisted Sysop UI palette before first paint.\n")
	buf.WriteString("applyTheme(getInitialTheme());\n\n")

	// Per-module provider wrappers (each renders <Outlet/>).
	if anyProviders {
		for _, m := range modules {
			if m == nil || len(m.Providers) == 0 {
				continue
			}
			fmt.Fprintf(&buf, "function %sModuleProviders() {\n  return (\n", pascalCase(m.ID))
			indent := "    "
			for i := range m.Providers {
				p := &m.Providers[i]
				fmt.Fprintf(&buf, "%s<%s>\n", indent, providerWrapName(p))
				indent += "  "
				for _, mt := range providerMounts(p) {
					fmt.Fprintf(&buf, "%s<%s />\n", indent, mt.Name)
				}
			}
			fmt.Fprintf(&buf, "%s<Outlet />\n", indent)
			for i := len(m.Providers) - 1; i >= 0; i-- {
				indent = indent[:len(indent)-2]
				fmt.Fprintf(&buf, "%s</%s>\n", indent, providerWrapName(&m.Providers[i]))
			}
			buf.WriteString("  );\n}\n\n")
		}
	}

	// ── Navigation table ──
	buf.WriteString("const navItems: { key: string; label: string; href: string; icon: React.ReactNode }[] = [\n")
	for _, item := range navItems {
		fmt.Fprintf(&buf, "  { key: %q, label: %q, href: %q, icon: <%s className=\"h-4 w-4\" /> },\n",
			item.Href, item.Label, item.Href, lucideComponentName(item.Icon))
	}
	buf.WriteString("];\n\n")
	buf.WriteString("const pageTitles: Record<string, string> = {\n")
	for _, item := range navItems {
		fmt.Fprintf(&buf, "  %q: %q,\n", item.Href, item.Label)
	}
	buf.WriteString("};\n\n")

	// ── AppShell ──
	buf.WriteString("function AppShell({ children }: { children: React.ReactNode }) {\n")
	buf.WriteString("  const location = useLocation();\n")
	buf.WriteString("  const navigate = useNavigate();\n")
	buf.WriteString("  const pathname = location.pathname;\n")
	buf.WriteString("  const nav: NavRailItem[] = navItems.map((item) => ({\n")
	buf.WriteString("    key: item.key,\n")
	buf.WriteString("    label: item.label,\n")
	buf.WriteString("    icon: item.icon,\n")
	buf.WriteString("    active: pathname === item.href || pathname.startsWith(item.href + \"/\"),\n")
	buf.WriteString("    onSelect: () => navigate(item.href),\n")
	buf.WriteString("  }));\n")
	fmt.Fprintf(&buf, "  const title = pageTitles[pathname] ?? %q;\n", brandTitle)
	buf.WriteString("  return (\n")
	buf.WriteString("    <div className=\"flex h-screen bg-bg text-text\">\n")
	fmt.Fprintf(&buf, "      <NavRail items={nav} logo={<%s className=\"h-4 w-4\" />} logoLabel=%q />\n", brandIconComp, brandTitle)
	buf.WriteString("      <div className=\"flex min-w-0 flex-1 flex-col\">\n")
	buf.WriteString("        <PageHeader title={title}>\n")
	buf.WriteString("          <ThemeSwitcher />\n")
	buf.WriteString("        </PageHeader>\n")
	buf.WriteString("        <main className=\"min-h-0 flex-1 overflow-auto\">\n")
	buf.WriteString("          <Suspense fallback={<div className=\"flex items-center justify-center py-12 text-text-muted\">Loading…</div>}>\n")
	buf.WriteString("            {children}\n")
	buf.WriteString("          </Suspense>\n")
	buf.WriteString("        </main>\n")
	buf.WriteString("      </div>\n")
	buf.WriteString("      <Toaster />\n")
	buf.WriteString("    </div>\n")
	buf.WriteString("  );\n}\n\n")

	// ── App export ──
	buf.WriteString("export default function App() {\n  return (\n")
	buf.WriteString("    <BrowserRouter>\n      <AppShell>\n        <Routes>\n")
	if len(navItems) > 0 {
		fmt.Fprintf(&buf, "          <Route path=\"/\" element={<Navigate to=%q replace />} />\n", navItems[0].Href)
	}
	if anyProviders {
		for _, m := range modules {
			if m == nil || len(pagesForModule(allModulePagesFromCtx(ctx), m.ID)) == 0 {
				continue
			}
			if len(m.Providers) > 0 {
				fmt.Fprintf(&buf, "          <Route element={<%sModuleProviders />}>\n", pascalCase(m.ID))
				fmt.Fprintf(&buf, "            {getModuleRoutes_%s()}\n", toPascalCase(m.ID))
				buf.WriteString("          </Route>\n")
			} else {
				fmt.Fprintf(&buf, "          {getModuleRoutes_%s()}\n", toPascalCase(m.ID))
			}
		}
	} else {
		buf.WriteString("          {getRoutes()}\n")
	}
	buf.WriteString("        </Routes>\n      </AppShell>\n    </BrowserRouter>\n")
	buf.WriteString("  );\n}\n")

	routesPages := ctx.AllPages
	if len(routesPages) == 0 {
		routesPages = ctx.Pages
	}
	files := []renderer.OutputFile{{Path: "App.tsx", Content: buf.Bytes()}}
	files = append(files, renderRoutes(routesPages, modules)...)
	return files, nil
}

// kitItemsExpr builds the JSX expression for the kit DataTable `items` prop:
// the datasource hook variable, with any select-filter and search `.filter()`
// chains registered by sibling controls applied. Mirrors the legacy
// data-table's `data=` chain so button/select/search wiring "actually feeds"
// the table.
func kitItemsExpr(varName, datasource string, imports *importTracker) string {
	searchVar, hasSearch := imports.searchDatasources[datasource]
	filters := imports.filterBindings[datasource]
	if !hasSearch && len(filters) == 0 {
		return varName + " ?? []"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "(%s ?? [])", varName)
	for _, fb := range filters {
		fmt.Fprintf(&b, "\n    .filter((item) => %s === \"all\" || String(item.%s) === %s)", fb.stateVar, fb.field, fb.stateVar)
	}
	if hasSearch {
		fmt.Fprintf(&b, "\n    .filter((item) => {")
		fmt.Fprintf(&b, "\n      if (!%s) return true;", searchVar)
		fmt.Fprintf(&b, "\n      const q = %s.toLowerCase();", searchVar)
		fmt.Fprintf(&b, "\n      return Object.values(item).some((v) => String(v).toLowerCase().includes(q));")
		fmt.Fprintf(&b, "\n    })")
	}
	return b.String()
}

// renderKitDataTableComponent emits the sysop kit's <DataTable> for a
// `data-table` component. The kit's DataTable API differs from the legacy
// local one:
//
//	legacy: <DataTable data={rows} columns={[{accessorKey,header,cell}]} onRowClick=.../>
//	kit:    <DataTable items={rows} columns={[{key,header,cell:(item)=>...}]}
//	          getRowId={(item)=>...} onRowOpen={(id,item)=>...}
//	          emptyState={<EmptyState .../>} />
//
// Columns carry function-valued `cell` renderers, emitted as arrow closures.
// `render: badge` columns become the kit's <StatusBadge>, which does its own
// theme-aware tone lookup — replacing the legacy hand-rolled variants map.
func renderKitDataTableComponent(buf *bytes.Buffer, c *config.Component, indent string, imports *importTracker, ctx *renderer.RenderContext) {
	imports.addKit("DataTable")
	datasource := getPropString(c.Props, "datasource", "")
	dataExpr := "[]"
	idField := "id"
	if datasource != "" {
		dataExpr = kitItemsExpr(toCamelCase(datasource), datasource, imports)
		idField = primaryField(ctx, datasource)
	}

	fmt.Fprintf(buf, "%s<DataTable\n", indent)
	fmt.Fprintf(buf, "%s  items={%s}\n", indent, dataExpr)
	fmt.Fprintf(buf, "%s  getRowId={(item) => String(item.%s)}\n", indent, idField)

	// Columns.
	fmt.Fprintf(buf, "%s  columns={[\n", indent)
	if cols, ok := c.Props["columns"].([]interface{}); ok {
		for _, col := range cols {
			colMap, ok := col.(map[string]interface{})
			if !ok {
				continue
			}
			field := fmt.Sprintf("%v", colMap["field"])
			label := fmt.Sprintf("%v", colMap["label"])

			var cellExpr string
			if columnRendersBadge(colMap["render"]) {
				imports.addKit("StatusBadge")
				cellExpr = fmt.Sprintf("<StatusBadge status={String(item.%s ?? \"\")} />", field)
			} else {
				cellExpr = fmt.Sprintf("String(item.%s ?? \"\")", field)
			}
			fmt.Fprintf(buf, "%s    {\n", indent)
			fmt.Fprintf(buf, "%s      key: %q,\n", indent, field)
			fmt.Fprintf(buf, "%s      header: %q,\n", indent, label)
			fmt.Fprintf(buf, "%s      cell: (item) => %s,\n", indent, cellExpr)
			fmt.Fprintf(buf, "%s      sortValue: (item) => String(item.%s ?? \"\"),\n", indent, field)
			fmt.Fprintf(buf, "%s    },\n", indent)
		}
	}
	fmt.Fprintf(buf, "%s  ]}\n", indent)

	// Row-open: kit's onRowOpen replaces the legacy onRowClick.
	if rowClick, ok := c.Actions["rowClick"]; ok && rowClick.Type == "navigate" && rowClick.Page != "" {
		imports.addRouter()
		route := imports.resolveRoute(rowClick.Page)
		if _, hasID := rowClick.Params["id"]; hasID {
			fmt.Fprintf(buf, "%s  onRowOpen={(id) => %s}\n", indent, imports.routerPushCall(fmt.Sprintf("`%s/${id}`", route)))
		} else {
			fmt.Fprintf(buf, "%s  onRowOpen={() => %s}\n", indent, imports.routerPushCall(fmt.Sprintf("%q", route)))
		}
	}

	// Empty state: kit's <EmptyState> replaces the legacy emptyMessage string.
	imports.addKit("EmptyState")
	_, hasSearch := imports.searchDatasources[datasource]
	if hasSearch || len(imports.filterBindings[datasource]) > 0 {
		fmt.Fprintf(buf, "%s  emptyState={<EmptyState variant=\"no-results\" title=\"No matches\" description=\"No rows match the current filters or search.\" />}\n", indent)
	} else {
		fmt.Fprintf(buf, "%s  emptyState={<EmptyState variant=\"empty\" title=\"Nothing here yet\" description=\"There are no rows to show.\" />}\n", indent)
	}
	fmt.Fprintf(buf, "%s/>\n", indent)
}
