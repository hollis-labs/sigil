package reactshadcn

import (
	"strings"
	"testing"

	"github.com/chrispian/sigil/internal/components"
	"github.com/chrispian/sigil/internal/config"
	"github.com/chrispian/sigil/internal/renderer"
)

func testPage(layout config.Component) *config.Page {
	return &config.Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "test-page",
		Title:   "Test Page",
		Overlay: "page",
		Layout:  layout,
	}
}

func testCtx(page *config.Page) *renderer.RenderContext {
	return &renderer.RenderContext{
		Page:     page,
		Registry: components.NewDefaultRegistry(),
		Theme: &renderer.ThemeConfig{
			Name: "default",
			Tokens: map[string]map[string]interface{}{
				"colors": {"background": "9 9 11", "text": "244 244 245"},
			},
		},
		DataSources: map[string]*renderer.DataSourceManifest{},
	}
}

func TestRendererName(t *testing.T) {
	r := &ReactShadcnRenderer{}
	if r.Name() != "react-shadcn" {
		t.Errorf("expected react-shadcn, got %s", r.Name())
	}
}

func TestRendererRegistered(t *testing.T) {
	r, ok := renderer.GetRenderer("react-shadcn")
	if !ok {
		t.Fatal("react-shadcn renderer not registered")
	}
	if r.Name() != "react-shadcn" {
		t.Errorf("expected react-shadcn, got %s", r.Name())
	}
}

func TestRenderBasicPage(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Props: map[string]interface{}{"gap": 4},
		Children: []config.Component{
			{
				ID:   "header",
				Type: "heading",
				Props: map[string]interface{}{
					"level": 2,
					"text":  "Hello World",
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Path != "pages/test-page.tsx" {
		t.Errorf("expected pages/test-page.tsx, got %s", files[0].Path)
	}

	content := string(files[0].Content)
	checks := []string{
		"export default function TestPage()",
		"Hello World",
		"flex flex-col",
		`"use client"`,
		"h2",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output:\n%s", check, content)
		}
	}
}

func TestRenderButton(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				ID:   "btn",
				Type: "button",
				Props: map[string]interface{}{
					"label":   "Click Me",
					"variant": "primary",
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}

	content := string(files[0].Content)
	checks := []string{
		"import { Button }",
		"@/components/ui/button",
		"Click Me",
		`variant="default"`,
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output:\n%s", check, content)
		}
	}
}

func TestRenderButtonWithIcon(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "button",
				Props: map[string]interface{}{
					"label":   "Add",
					"variant": "primary",
					"icon":    "plus",
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "lucide-react") {
		t.Error("expected lucide-react import")
	}
	if !strings.Contains(content, "<Plus") {
		t.Error("expected <Plus icon component")
	}
}

func TestRenderBadge(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "badge",
				Props: map[string]interface{}{
					"text":    "Active",
					"variant": "success",
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "<Badge") {
		t.Error("expected Badge component")
	}
	if !strings.Contains(content, "Active") {
		t.Error("expected badge value")
	}
}

func TestRenderDataTable(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "data-table",
				Props: map[string]interface{}{
					"datasource": "users",
					"columns": []interface{}{
						map[string]interface{}{"field": "name", "label": "Name"},
						map[string]interface{}{"field": "email", "label": "Email"},
					},
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	checks := []string{
		"<DataTable",
		"accessorKey",
		`"name"`,
		`"email"`,
		"data-table",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output:\n%s", check, content)
		}
	}
}

func TestRenderForm(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "form",
				Props: map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"name":  "email",
							"label": "Email",
							"type":  "email",
						},
						map[string]interface{}{
							"name":  "bio",
							"label": "Bio",
							"type":  "textarea",
						},
					},
					"submit": map[string]interface{}{"label": "Save"},
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	checks := []string{
		"<form",
		"<Input",
		"<Textarea",
		"<Label",
		"Email",
		"Bio",
		"Save",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output:\n%s", check, content)
		}
	}
}

func TestRenderSelect(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "select",
				Props: map[string]interface{}{
					"placeholder": "Choose...",
					"options": []interface{}{
						map[string]interface{}{"value": "a", "label": "Option A"},
						map[string]interface{}{"value": "b", "label": "Option B"},
					},
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	checks := []string{
		"<Select>",
		"<SelectTrigger>",
		"<SelectContent>",
		"<SelectItem",
		"Option A",
		"Option B",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output:\n%s", check, content)
		}
	}
}

func TestRenderModal(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "modal",
				Props: map[string]interface{}{
					"title": "Create Item",
				},
				Children: []config.Component{
					{Type: "text", Props: map[string]interface{}{"text": "Modal body"}},
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	checks := []string{
		"<Dialog>",
		"<DialogContent>",
		"<DialogHeader>",
		"<DialogTitle>Create Item</DialogTitle>",
		"Modal body",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output:\n%s", check, content)
		}
	}
}

func TestRenderSheet(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "sheet",
				Props: map[string]interface{}{
					"title": "Details",
					"side":  "right",
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	checks := []string{
		"<Sheet>",
		"<SheetContent",
		`side="right"`,
		"<SheetTitle>Details</SheetTitle>",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output:\n%s", check, content)
		}
	}
}

func TestRenderTabs(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "tabs",
				Children: []config.Component{
					{
						ID:    "tab1",
						Props: map[string]interface{}{"label": "First"},
						Children: []config.Component{
							{Type: "text", Props: map[string]interface{}{"text": "Tab 1 content"}},
						},
					},
					{
						ID:    "tab2",
						Props: map[string]interface{}{"label": "Second"},
					},
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	checks := []string{
		"<Tabs",
		"<TabsList>",
		"<TabsTrigger",
		"<TabsContent",
		"First",
		"Second",
		"Tab 1 content",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output:\n%s", check, content)
		}
	}
}

func TestRenderGrid(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "grid",
		Props: map[string]interface{}{"columns": 3, "gap": 4},
		Children: []config.Component{
			{Type: "text", Props: map[string]interface{}{"text": "A"}},
			{Type: "text", Props: map[string]interface{}{"text": "B"}},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "grid grid-cols-3 gap-4") {
		t.Error("expected grid classes")
	}
}

func TestRenderColumns(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "columns",
		Props: map[string]interface{}{"gap": 6, "justify": "between"},
		Children: []config.Component{
			{Type: "text", Props: map[string]interface{}{"text": "Left"}},
			{Type: "text", Props: map[string]interface{}{"text": "Right"}},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "flex flex-row gap-6 justify-between") {
		t.Error("expected columns classes")
	}
}

func TestRenderSeparator(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{Type: "separator"},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "<Separator") {
		t.Error("expected Separator component")
	}
}

func TestRenderAlert(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "alert",
				Props: map[string]interface{}{
					"message": "Something happened",
					"variant": "danger",
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "<Alert") {
		t.Error("expected Alert component")
	}
	if !strings.Contains(content, "destructive") {
		t.Error("expected destructive variant mapping")
	}
}

func TestRenderAvatar(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "avatar",
				Props: map[string]interface{}{
					"src": "/avatar.png",
					"alt": "John",
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	checks := []string{"<Avatar>", "<AvatarImage", "<AvatarFallback>", "JO"}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in output", check)
		}
	}
}

func TestRenderProgress(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{Type: "progress", Props: map[string]interface{}{"value": 75}},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "<Progress") {
		t.Error("expected Progress component")
	}
}

func TestRenderSpacer(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{Type: "spacer", Props: map[string]interface{}{"size": "8"}},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, `className="h-8"`) {
		t.Error("expected spacer div")
	}
}

func TestRenderLabel(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{Type: "label", Props: map[string]interface{}{"text": "Name"}},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "<Label>Name</Label>") {
		t.Error("expected Label component")
	}
}

func TestRenderIcon(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{Type: "icon", Props: map[string]interface{}{"name": "settings"}},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "<Settings") {
		t.Error("expected Settings icon")
	}
	if !strings.Contains(content, "lucide-react") {
		t.Error("expected lucide-react import")
	}
}

func TestRenderSearchBar(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{Type: "search-bar", Props: map[string]interface{}{"placeholder": "Search users..."}},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, `type="search"`) {
		t.Error("expected search input type")
	}
	if !strings.Contains(content, "Search users...") {
		t.Error("expected placeholder")
	}
}

func TestRenderInput(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				ID:   "email-input",
				Type: "input",
				Props: map[string]interface{}{
					"placeholder": "Enter email",
					"type":        "email",
					"required":    true,
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, `type="email"`) {
		t.Error("expected email type")
	}
	if !strings.Contains(content, "required") {
		t.Error("expected required attribute")
	}
}

func TestSharedComponentsDataTable(t *testing.T) {
	files, err := renderSharedComponents([]string{"data-table", "button"})
	if err != nil {
		t.Fatalf("renderSharedComponents: %v", err)
	}

	// Should include data-table.tsx, lib/utils.ts, components/index.ts
	paths := map[string]bool{}
	for _, f := range files {
		paths[f.Path] = true
	}

	if !paths["components/data-table.tsx"] {
		t.Error("expected data-table.tsx")
	}
	if !paths["lib/utils.ts"] {
		t.Error("expected lib/utils.ts")
	}
	if !paths["components/index.ts"] {
		t.Error("expected components/index.ts")
	}
}

func TestSharedComponentsNoDT(t *testing.T) {
	files, err := renderSharedComponents([]string{"button", "badge"})
	if err != nil {
		t.Fatalf("renderSharedComponents: %v", err)
	}

	for _, f := range files {
		if f.Path == "components/data-table.tsx" {
			t.Error("data-table.tsx should not be generated when not used")
		}
	}
}

func TestDataSourceHooks(t *testing.T) {
	ds := &renderer.DataSourceManifest{
		Alias:        "sprint",
		Description:  "Sprint items",
		Capabilities: []string{"list", "create", "update", "delete"},
		Fields: []renderer.DataSourceField{
			{Name: "id", Type: "integer", Primary: true},
			{Name: "title", Type: "string", Required: true},
			{Name: "status", Type: "enum", Values: []string{"todo", "doing", "done"}},
			{Name: "created_at", Type: "datetime", Readonly: true},
		},
		Endpoints: map[string]string{
			"list":   "GET /api/sprints",
			"create": "POST /api/sprints",
			"read":   "GET /api/sprints/{id}",
			"update": "PUT /api/sprints/{id}",
			"delete": "DELETE /api/sprints/{id}",
		},
	}

	files, err := renderDataSourceHooks(ds)
	if err != nil {
		t.Fatalf("renderDataSourceHooks: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	// Check types file
	var typesContent, hookContent string
	for _, f := range files {
		if strings.HasPrefix(f.Path, "types/") {
			typesContent = string(f.Content)
		}
		if strings.HasPrefix(f.Path, "hooks/") {
			hookContent = string(f.Content)
		}
	}

	// Types checks
	typesChecks := []string{
		"export interface Sprint {",
		"id: number;",
		"title: string;",
		"status?: string;",
		"export interface SprintInput {",
	}
	for _, check := range typesChecks {
		if !strings.Contains(typesContent, check) {
			t.Errorf("expected %q in types:\n%s", check, typesContent)
		}
	}
	// Primary and readonly should not be in input
	if strings.Contains(typesContent, "SprintInput {\n  id") {
		t.Error("id should not be in SprintInput")
	}

	// Hook checks
	hookChecks := []string{
		"export function useSprint()",
		"useSWR",
		"/api/sprints",
		"export async function createSprint",
		"export async function updateSprint",
		"export async function deleteSprint",
		`method: "POST"`,
		`method: "PUT"`,
		`method: "DELETE"`,
	}
	for _, check := range hookChecks {
		if !strings.Contains(hookContent, check) {
			t.Errorf("expected %q in hook:\n%s", check, hookContent)
		}
	}
}

func TestThemeGeneration(t *testing.T) {
	theme := &renderer.ThemeConfig{
		Name: "dark",
		Tokens: map[string]map[string]interface{}{
			"colors": {
				"background": "9 9 11",
				"text":       "244 244 245",
				"accent":     "79 70 229",
			},
			"radius": {
				"md": "0.375rem",
				"lg": "0.5rem",
			},
		},
	}

	files, err := renderTheme(theme)
	if err != nil {
		t.Fatalf("renderTheme: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	var cssContent, twContent string
	for _, f := range files {
		switch f.Path {
		case "globals.css":
			cssContent = string(f.Content)
		case "tailwind.config.ts":
			twContent = string(f.Content)
		}
	}

	// CSS checks
	cssChecks := []string{
		"@tailwind base",
		"@tailwind components",
		"@tailwind utilities",
		":root {",
		"--sigil-background: 9 9 11",
		"--sigil-text: 244 244 245",
		"--sigil-accent: 79 70 229",
		"--sigil-md: 0.375rem",
	}
	for _, check := range cssChecks {
		if !strings.Contains(cssContent, check) {
			t.Errorf("expected %q in CSS:\n%s", check, cssContent)
		}
	}

	// Tailwind config checks
	twChecks := []string{
		`import type { Config }`,
		"sigil: {",
		`"rgb(var(--sigil-background)`,
		`"rgb(var(--sigil-text)`,
		"borderRadius:",
		`"sigil-md"`,
		"export default config",
	}
	for _, check := range twChecks {
		if !strings.Contains(twContent, check) {
			t.Errorf("expected %q in Tailwind config:\n%s", check, twContent)
		}
	}
}

func TestE2EGeneration(t *testing.T) {
	r := &ReactShadcnRenderer{}

	page := &config.Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "dashboard",
		Title:   "Dashboard",
		Overlay: "page",
		Layout: config.Component{
			ID:   "root",
			Type: "rows",
			Props: map[string]interface{}{"gap": 4},
			Children: []config.Component{
				{Type: "heading", Props: map[string]interface{}{"level": 1, "text": "Dashboard"}},
				{
					Type: "columns",
					Props: map[string]interface{}{"gap": 4},
					Children: []config.Component{
						{Type: "button", Props: map[string]interface{}{"label": "New", "variant": "primary"}},
						{Type: "badge", Props: map[string]interface{}{"text": "3", "variant": "default"}},
					},
				},
				{
					Type: "data-table",
					Props: map[string]interface{}{
						"datasource": "items",
						"columns": []interface{}{
							map[string]interface{}{"field": "name", "label": "Name"},
						},
					},
				},
				{Type: "separator"},
				{Type: "alert", Props: map[string]interface{}{"message": "Done!", "variant": "info"}},
			},
		},
	}

	ctx := &renderer.RenderContext{
		Page:     page,
		Registry: components.NewDefaultRegistry(),
		Theme: &renderer.ThemeConfig{
			Name:   "test",
			Tokens: map[string]map[string]interface{}{"colors": {"bg": "0 0 0"}},
		},
		DataSources: map[string]*renderer.DataSourceManifest{},
	}

	// Render page
	pageFiles, err := r.Render(ctx)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if len(pageFiles) == 0 {
		t.Fatal("expected page files")
	}

	// Shared components
	sharedFiles, err := r.SharedComponents([]string{"heading", "button", "badge", "data-table", "separator", "alert"})
	if err != nil {
		t.Fatalf("SharedComponents: %v", err)
	}
	if len(sharedFiles) < 2 {
		t.Errorf("expected at least 2 shared files, got %d", len(sharedFiles))
	}

	// Theme
	themeFiles, err := r.RenderTheme(ctx.Theme)
	if err != nil {
		t.Fatalf("RenderTheme: %v", err)
	}
	if len(themeFiles) != 2 {
		t.Errorf("expected 2 theme files, got %d", len(themeFiles))
	}

	// Verify all files have content
	allFiles := append(pageFiles, sharedFiles...)
	allFiles = append(allFiles, themeFiles...)
	for _, f := range allFiles {
		if len(f.Content) == 0 {
			t.Errorf("empty file: %s", f.Path)
		}
	}

	// Verify page TSX contains all component types
	content := string(pageFiles[0].Content)
	expected := []string{"Button", "Badge", "DataTable", "Separator", "Alert", "h1"}
	for _, e := range expected {
		if !strings.Contains(content, e) {
			t.Errorf("expected %q in E2E page output", e)
		}
	}
}

func TestPageWithDataSource(t *testing.T) {
	page := &config.Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "user-list",
		Title:   "User List",
		Overlay: "page",
		DataSources: []config.DataSourceRef{
			{Alias: "users"},
		},
		Layout: config.Component{
			ID:   "root",
			Type: "rows",
			Children: []config.Component{
				{Type: "heading", Props: map[string]interface{}{"level": 2, "text": "Users"}},
			},
		},
	}
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "useUsers()") {
		t.Error("expected datasource hook call")
	}
	if !strings.Contains(content, "usersLoading") {
		t.Error("expected loading state variable")
	}
}

func TestNavigateAction(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				Type: "button",
				Props: map[string]interface{}{"label": "Go"},
				Actions: map[string]config.Action{
					"click": {Type: "navigate", Page: "dashboard"},
				},
			},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "router.push") {
		t.Error("expected router.push for navigate action")
	}
}

func TestUnknownComponent(t *testing.T) {
	page := testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{ID: "custom", Type: "my-widget"},
		},
	})
	ctx := testCtx(page)

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, `data-component="my-widget"`) {
		t.Error("expected fallback div for unknown component")
	}
}

func TestHelperFunctions(t *testing.T) {
	if toPascalCase("user-profile") != "UserProfile" {
		t.Errorf("toPascalCase: %s", toPascalCase("user-profile"))
	}
	if toCamelCase("user-profile") != "userProfile" {
		t.Errorf("toCamelCase: %s", toCamelCase("user-profile"))
	}
	if toKebabCase("user_profile") != "user-profile" {
		t.Errorf("toKebabCase: %s", toKebabCase("user_profile"))
	}
	if mapButtonVariant("primary") != "default" {
		t.Error("primary should map to default")
	}
	if mapButtonVariant("destructive") != "destructive" {
		t.Error("destructive should stay destructive")
	}
}
