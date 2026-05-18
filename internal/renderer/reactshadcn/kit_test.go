package reactshadcn

import (
	"strings"
	"testing"

	"github.com/chrispian/sigil/internal/config"
	"github.com/chrispian/sigil/internal/renderer"
)

// kitCtx returns a RenderContext for the given page with sysop kit mode on.
func kitCtx(page *config.Page) *renderer.RenderContext {
	ctx := testCtx(page)
	ctx.UIKit = "sysop"
	return ctx
}

// TestEmitActionButtonOnClick guards the FND-4 emit-action bug: a button whose
// emit/filter action carries a datasource+field must emit an onClick handler
// that sets the shared filter state — never an onChange reading e.target.value
// (which React silently drops on a <button>).
func TestEmitActionButtonOnClick(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type:  "button",
				Props: map[string]interface{}{"label": "Running"},
				Actions: map[string]config.Action{
					"click": {Type: "emit", Datasource: "Run", Field: "status", Value: "running"},
				},
			},
			{
				Type:  "button",
				Props: map[string]interface{}{"label": "All"},
				Actions: map[string]config.Action{
					"click": {Type: "emit", Datasource: "Run", Field: "status"},
				},
			},
		},
	})

	files, err := renderPage(testCtx(page))
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)

	if strings.Contains(content, "onChange") {
		t.Errorf("emit action on a button must not emit onChange:\n%s", content)
	}
	if strings.Contains(content, "e.target.value") {
		t.Errorf("emit action on a button must not read e.target.value:\n%s", content)
	}
	for _, want := range []string{
		`onClick={() => setFilterRunStatus("running")}`,
		`onClick={() => setFilterRunStatus("all")}`, // empty Value resets to "all"
		`const [filterRunStatus, setFilterRunStatus] = useState("all");`,
	} {
		if !strings.Contains(content, want) {
			t.Errorf("expected %q in output:\n%s", want, content)
		}
	}
	// Two buttons share one piece of filter state — declared exactly once.
	if n := strings.Count(content, "setFilterRunStatus] = useState"); n != 1 {
		t.Errorf("expected filter state declared once, got %d", n)
	}
}

// TestKitModeDataTable verifies the data-table emits the sysop kit's DataTable
// (items/getRowId/columns/emptyState) imported from @hollis-labs/sysop-ui.
func TestKitModeDataTable(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "data-table",
				Props: map[string]interface{}{
					"datasource": "Project",
					"columns": []interface{}{
						map[string]interface{}{"field": "name", "label": "Name"},
						map[string]interface{}{"field": "status", "label": "Status", "render": "badge"},
					},
				},
			},
		},
	})

	files, err := renderPage(kitCtx(page))
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)

	for _, want := range []string{
		`from "@hollis-labs/sysop-ui"`,
		"<DataTable",
		"items={project ?? []}",
		"getRowId={(item) => String(item.id)}",
		"cell: (item) => String(item.name ?? \"\")",
		"<StatusBadge status={String(item.status ?? \"\")} />",
		"emptyState={<EmptyState",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("expected %q in output:\n%s", want, content)
		}
	}
	// Kit mode must not pull the app-local data-table component.
	if strings.Contains(content, `from "@/components/data-table"`) {
		t.Errorf("kit mode must not import the app-local data-table:\n%s", content)
	}
	if strings.Contains(content, "accessorKey") {
		t.Errorf("kit DataTable uses key/cell, not tanstack accessorKey:\n%s", content)
	}
}

// TestKitModeImportRedirect verifies shadcn ui primitives the kit re-exports
// are redirected to the kit barrel, while primitives it does not (select) stay
// app-local.
func TestKitModeImportRedirect(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{Type: "button", Props: map[string]interface{}{"label": "Go"}},
			{Type: "badge", Props: map[string]interface{}{"text": "New"}},
			{Type: "select", Props: map[string]interface{}{"placeholder": "Pick"}},
		},
	})

	files, err := renderPage(kitCtx(page))
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)

	if !strings.Contains(content, `from "@hollis-labs/sysop-ui"`) {
		t.Errorf("expected kit import for Button/Badge:\n%s", content)
	}
	if strings.Contains(content, `from "@/components/ui/button"`) ||
		strings.Contains(content, `from "@/components/ui/badge"`) {
		t.Errorf("Button/Badge must redirect to the kit barrel:\n%s", content)
	}
	if !strings.Contains(content, `from "@/components/ui/select"`) {
		t.Errorf("Select has no kit equivalent — must stay app-local:\n%s", content)
	}
}

// TestKitModeFilterBar verifies the filter-bar custom component is emitted as
// the kit FilterBar with FilterCycleToggle chips, wired to shared filter state.
func TestKitModeFilterBar(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "filter-bar",
				Props: map[string]interface{}{
					"datasource": "Run",
					"filters": []interface{}{
						map[string]interface{}{
							"type":  "select",
							"field": "status",
							"label": "Status",
							"options": []interface{}{
								map[string]interface{}{"label": "All", "value": ""},
								map[string]interface{}{"label": "Running", "value": "running"},
							},
						},
						map[string]interface{}{"type": "date-range", "field": "created_at", "label": "Date"},
					},
				},
			},
		},
	})

	files, err := renderPage(kitCtx(page))
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)

	for _, want := range []string{
		"<FilterBar",
		"<FilterCycleToggle<string>",
		"value={filterRunStatus}",
		"onChange={setFilterRunStatus}",
		`searchQuery={searchRun}`,
		`activeFilterCount={[filterRunStatus].filter((v) => v !== "all").length}`,
	} {
		if !strings.Contains(content, want) {
			t.Errorf("expected %q in output:\n%s", want, content)
		}
	}
	// The empty-valued "All" option is normalized to "all".
	if !strings.Contains(content, `{ value: "all", label: "All" }`) {
		t.Errorf("empty filter value should normalize to \"all\":\n%s", content)
	}
	// date-range has no kit equivalent — skipped, not crashed.
	if strings.Contains(content, "date-range") {
		t.Errorf("date-range filter should be skipped:\n%s", content)
	}
}

// TestKitLayoutSPA verifies kit mode emits a thin App.tsx composed from the
// kit shell — NavRail/PageHeader/ThemeSwitcher, kit theme tokens — with no
// next-themes and no app-local @/components/ui/* chrome.
func TestKitLayoutSPA(t *testing.T) {
	shell := &config.Page{
		Sigil: "1.0", Kind: "page", ID: "app-shell", Title: "App", Overlay: "page",
		Layout: config.Component{Type: "rows", Children: []config.Component{
			{Type: "nav-menu", Props: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"label": "Board", "page": "board", "icon": "list"},
				},
			}},
		}},
	}
	mod := &config.ModuleConfig{ID: "app", Shell: "app-shell", Pages: []string{"board"}}
	app := &config.AppConfig{
		Name:       "Kit App",
		TargetMode: "spa",
		Features:   &config.AppFeatures{ThemeToggle: true, CommandPalette: true},
		Modules:    []config.ModuleConfig{*mod},
	}
	pages := []*config.Page{
		{Sigil: "1.0", Kind: "page", ID: "board", Title: "Board", Overlay: "page", Module: "app", Layout: config.Component{Type: "rows"}},
	}
	ctx := &renderer.LayoutContext{
		Module: mod, Shell: shell, AppConfig: app,
		Pages: pages, AllModules: []*config.ModuleConfig{mod}, AllPages: pages,
		AllShells: map[string]*config.Page{"app": shell},
		UIKit:     "sysop",
	}

	files, err := renderLayout(ctx)
	if err != nil {
		t.Fatalf("renderLayout: %v", err)
	}
	var appTsx string
	for _, f := range files {
		if f.Path == "App.tsx" {
			appTsx = string(f.Content)
		}
	}
	if appTsx == "" {
		t.Fatalf("kit layout must emit App.tsx")
	}

	for _, want := range []string{
		`from "@hollis-labs/sysop-ui"`,
		"NavRail",
		"PageHeader",
		"ThemeSwitcher",
		"bg-bg text-text",
		"applyTheme(getInitialTheme())",
	} {
		if !strings.Contains(appTsx, want) {
			t.Errorf("kit App.tsx missing %q:\n%s", want, appTsx)
		}
	}
	for _, forbidden := range []string{"next-themes", `@/components/ui/`} {
		if strings.Contains(appTsx, forbidden) {
			t.Errorf("kit App.tsx must not reference %q:\n%s", forbidden, appTsx)
		}
	}
}

// TestKitTheme verifies kit mode emits a thin entry stylesheet importing the
// kit theme and skips tailwind.config.ts.
func TestKitTheme(t *testing.T) {
	theme := &renderer.ThemeConfig{Name: "default", Tokens: map[string]map[string]interface{}{}}

	files, err := renderTheme(theme, "sysop")
	if err != nil {
		t.Fatalf("renderTheme: %v", err)
	}
	if len(files) != 1 || files[0].Path != "globals.css" {
		t.Fatalf("kit mode should emit only globals.css, got %d files", len(files))
	}
	css := string(files[0].Content)
	if !strings.Contains(css, `@import "@hollis-labs/sysop-ui/theme.css"`) {
		t.Errorf("globals.css should import the kit theme:\n%s", css)
	}
}
