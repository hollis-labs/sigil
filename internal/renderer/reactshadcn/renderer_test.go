package reactshadcn

import (
	"os"
	"path/filepath"
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
		ID:    "root",
		Type:  "rows",
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
		ID:    "root",
		Type:  "grid",
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
	if !strings.Contains(content, "grid grid-cols-1 md:grid-cols-3 gap-4") {
		t.Error("expected responsive grid classes")
	}
}

func TestRenderColumns(t *testing.T) {
	page := testPage(config.Component{
		ID:    "root",
		Type:  "columns",
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
	files, err := renderSharedComponents([]string{"data-table", "button"}, "")
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
	files, err := renderSharedComponents([]string{"button", "badge"}, "")
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

	files, err := renderTheme(theme, "")
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
			ID:    "root",
			Type:  "rows",
			Props: map[string]interface{}{"gap": 4},
			Children: []config.Component{
				{Type: "heading", Props: map[string]interface{}{"level": 1, "text": "Dashboard"}},
				{
					Type:  "columns",
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
				Type:  "button",
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

// pageWithNavigate builds a page whose row-click action navigates to another
// page. This exercises both the framework router import and the push call.
func pageWithNavigate() *config.Page {
	return testPage(config.Component{
		ID:   "root",
		Type: "rows",
		Children: []config.Component{
			{
				ID:   "table",
				Type: "data-table",
				Props: map[string]interface{}{
					"datasource": "Repo",
					"columns": []interface{}{
						map[string]interface{}{"field": "name", "label": "Name"},
					},
				},
				Actions: map[string]config.Action{
					"rowClick": {
						Type: "navigate",
						Page: "test-page-detail",
					},
				},
			},
		},
	})
}

func TestRenderPageTargetModeBranching(t *testing.T) {
	cases := []struct {
		name        string
		targetMode  string
		mustContain []string
		mustNotHave []string
	}{
		{
			name:       "default (empty) treated as app-router",
			targetMode: "",
			mustContain: []string{
				`"use client";`,
				`from "next/navigation"`,
				`useRouter`,
				`router.push`,
			},
			mustNotHave: []string{
				`react-router-dom`,
				`useNavigate`,
				`navigate(`,
			},
		},
		{
			name:       "explicit app-router",
			targetMode: "app-router",
			mustContain: []string{
				`"use client";`,
				`from "next/navigation"`,
				`useRouter`,
				`router.push`,
			},
			mustNotHave: []string{
				`react-router-dom`,
			},
		},
		{
			name:       "spa",
			targetMode: "spa",
			mustContain: []string{
				`from "react-router-dom"`,
				`useNavigate`,
				`navigate(`,
			},
			mustNotHave: []string{
				`"use client";`,
				`next/navigation`,
				`useRouter`,
				`router.push`,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			page := pageWithNavigate()
			ctx := testCtx(page)
			ctx.TargetMode = tc.targetMode

			files, err := renderPage(ctx)
			if err != nil {
				t.Fatalf("renderPage: %v", err)
			}
			content := string(files[0].Content)

			for _, want := range tc.mustContain {
				if !strings.Contains(content, want) {
					t.Errorf("expected %q in output:\n%s", want, content)
				}
			}
			for _, banned := range tc.mustNotHave {
				if strings.Contains(content, banned) {
					t.Errorf("did not expect %q in output:\n%s", banned, content)
				}
			}
		})
	}
}

func TestRenderAPIClientTargetMode(t *testing.T) {
	apiCfg := &config.APIConfig{
		BaseURLEnv:     "NEXT_PUBLIC_X_URL",
		BaseURLDefault: "http://localhost:8080",
		Prefix:         "/api",
	}

	cases := []struct {
		name        string
		targetMode  string
		mustContain []string
		mustNotHave []string
	}{
		{
			name:       "default → process.env + use-client",
			targetMode: "app-router",
			mustContain: []string{
				`"use client";`,
				`process.env.NEXT_PUBLIC_X_URL`,
			},
			mustNotHave: []string{
				`import.meta.env`,
			},
		},
		{
			name:       "spa → import.meta.env, no use-client",
			targetMode: "spa",
			mustContain: []string{
				`import.meta.env.NEXT_PUBLIC_X_URL`,
			},
			mustNotHave: []string{
				`"use client";`,
				`process.env.`,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files, err := renderAPIClient(apiCfg, tc.targetMode)
			if err != nil {
				t.Fatalf("renderAPIClient: %v", err)
			}
			if len(files) != 1 {
				t.Fatalf("expected 1 file, got %d", len(files))
			}
			content := string(files[0].Content)
			for _, want := range tc.mustContain {
				if !strings.Contains(content, want) {
					t.Errorf("expected %q in output:\n%s", want, content)
				}
			}
			for _, banned := range tc.mustNotHave {
				if strings.Contains(content, banned) {
					t.Errorf("did not expect %q in output:\n%s", banned, content)
				}
			}
		})
	}
}

func TestRenderRoutesEmitsLazyImportsAndRouteJSX(t *testing.T) {
	pages := []*config.Page{
		{
			Sigil:   "1.0",
			Kind:    "page",
			ID:      "se-repos",
			Title:   "Repos",
			Overlay: "page",
			Module:  "se",
			Layout:  config.Component{Type: "rows"},
		},
		{
			Sigil:   "1.0",
			Kind:    "page",
			ID:      "se-repo-detail",
			Title:   "Repo Detail",
			Overlay: "page",
			Module:  "se",
			Layout:  config.Component{Type: "rows"},
			DataSources: []config.DataSourceRef{{
				Alias:  "Repo",
				Params: map[string]interface{}{"id": "{{param.id}}"},
			}},
		},
	}

	files := renderRoutes(pages, nil)
	if len(files) != 1 || files[0].Path != "routes.tsx" {
		t.Fatalf("expected routes.tsx, got %+v", files)
	}

	content := string(files[0].Content)
	checks := []string{
		`import { lazy } from "react";`,
		`import { Route } from "react-router-dom";`,
		`const SeReposPage = lazy(() => import("@/pages/se-repos"));`,
		`const SeRepoDetailPage = lazy(() => import("@/pages/se-repo-detail"));`,
		`<Route path="/repos" element={<SeReposPage />} />`,
		`<Route path="/repo-detail/:id" element={<SeRepoDetailPage />} />`,
		`<Route path="/repo-detail" element={<SeRepoDetailPage />} />`,
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Errorf("expected %q in routes.tsx:\n%s", check, content)
		}
	}
}

// TestRenderRoutesMultiModulePrefixing covers the SPA multi-module case:
// when two modules each declare a route_group, their pages must mount under
// distinct URL prefixes so they don't collide silently. Single-module
// configs keep flat paths (unchanged from Phase 1).
func TestRenderRoutesMultiModulePrefixing(t *testing.T) {
	pages := []*config.Page{
		{
			Sigil: "1.0", Kind: "page", ID: "se-dashboard",
			Title: "SE Dashboard", Overlay: "page", Module: "se",
			Layout: config.Component{Type: "rows"},
		},
		{
			Sigil: "1.0", Kind: "page", ID: "clockwork-dashboard",
			Title: "CW Dashboard", Overlay: "page", Module: "clockwork",
			Layout: config.Component{Type: "rows"},
		},
		{
			Sigil: "1.0", Kind: "page", ID: "clockwork-task-detail",
			Title: "Task Detail", Overlay: "page", Module: "clockwork",
			Layout: config.Component{Type: "rows"},
			DataSources: []config.DataSourceRef{{
				Alias:  "Task",
				Params: map[string]interface{}{"id": "{{param.id}}"},
			}},
		},
	}

	modules := []*config.ModuleConfig{
		{ID: "se", RouteGroup: "(explorer)", Pages: []string{"se-dashboard"}},
		{ID: "clockwork", RouteGroup: "(clockwork)", Pages: []string{"clockwork-dashboard", "clockwork-task-detail"}},
	}

	files := renderRoutes(pages, modules)
	if len(files) != 1 || files[0].Path != "routes.tsx" {
		t.Fatalf("expected routes.tsx, got %+v", files)
	}
	content := string(files[0].Content)
	want := []string{
		`<Route path="/explorer/dashboard" element={<SeDashboardPage />} />`,
		`<Route path="/clockwork/dashboard" element={<ClockworkDashboardPage />} />`,
		`<Route path="/clockwork/task-detail/:id" element={<ClockworkTaskDetailPage />} />`,
		`<Route path="/clockwork/task-detail" element={<ClockworkTaskDetailPage />} />`,
	}
	for _, w := range want {
		if !strings.Contains(content, w) {
			t.Errorf("missing %q in routes.tsx:\n%s", w, content)
		}
	}
	// And the un-prefixed single-module paths must NOT appear; otherwise the
	// two modules would collide on /dashboard.
	bad := []string{
		`<Route path="/dashboard" element={<SeDashboardPage />} />`,
		`<Route path="/dashboard" element={<ClockworkDashboardPage />} />`,
	}
	for _, b := range bad {
		if strings.Contains(content, b) {
			t.Errorf("unexpected unprefixed route %q in multi-module output:\n%s", b, content)
		}
	}
}

// TestPageRouteSPASingleModuleBackCompat confirms a single-module SPA config
// keeps flat (un-prefixed) URLs — Phase 1 SPA output is preserved.
func TestPageRouteSPASingleModuleBackCompat(t *testing.T) {
	modules := []*config.ModuleConfig{
		{ID: "clockwork", RouteGroup: "(clockwork)", Pages: []string{"clockwork-board"}},
	}
	got := pageRouteSPA("clockwork-board", "clockwork", modules)
	if got != "/board" {
		t.Errorf("single-module SPA: pageRouteSPA = %q, want /board", got)
	}
}

// TestPageRouteSPACrossModule covers the cross-module navigation case: a
// page in module "se" linking to a page owned by module "clockwork" picks
// up clockwork's route prefix (not se's).
func TestPageRouteSPACrossModule(t *testing.T) {
	modules := []*config.ModuleConfig{
		{ID: "se", RouteGroup: "(explorer)", Pages: []string{"se-dashboard"}},
		{ID: "clockwork", RouteGroup: "(clockwork)", Pages: []string{"clockwork-board"}},
	}
	got := pageRouteSPA("clockwork-board", "se", modules)
	if got != "/clockwork/board" {
		t.Errorf("cross-module SPA: pageRouteSPA = %q, want /clockwork/board", got)
	}
}

func TestRenderLayoutTargetModeSelectsOutputShape(t *testing.T) {
	shell := &config.Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "x-shell",
		Title:   "Shell",
		Overlay: "page",
		Layout: config.Component{
			Type: "rows",
			Children: []config.Component{
				{
					Type: "nav-menu",
					Props: map[string]interface{}{
						"items": []interface{}{
							map[string]interface{}{"label": "Home", "page": "x-home", "icon": "home"},
						},
					},
				},
			},
		},
	}
	mod := &config.ModuleConfig{ID: "x", Shell: "x-shell", RouteGroup: "(x)", Pages: []string{"x-home"}}

	cases := []struct {
		name       string
		targetMode string
		wantPath   string
		wantExtra  string // additional file path expected (empty if none)
	}{
		{"app-router default", "", "app/(x)/layout.tsx", ""},
		{"app-router explicit", "app-router", "app/(x)/layout.tsx", ""},
		{"spa", "spa", "App.tsx", "routes.tsx"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &renderer.LayoutContext{
				Module:    mod,
				Shell:     shell,
				AppConfig: &config.AppConfig{Name: "X", TargetMode: tc.targetMode},
				Pages:     []*config.Page{},
			}
			files, err := renderLayout(ctx)
			if err != nil {
				t.Fatalf("renderLayout: %v", err)
			}
			paths := map[string]bool{}
			for _, f := range files {
				paths[f.Path] = true
			}
			if !paths[tc.wantPath] {
				t.Errorf("expected output path %q, got files %+v", tc.wantPath, paths)
			}
			if tc.wantExtra != "" && !paths[tc.wantExtra] {
				t.Errorf("expected extra path %q, got files %+v", tc.wantExtra, paths)
			}
			// SPA mode must NOT produce an app/.../layout.tsx; app-router must NOT produce App.tsx.
			if tc.targetMode == "spa" {
				for path := range paths {
					if strings.HasPrefix(path, "app/") {
						t.Errorf("spa mode produced unexpected app-router path %q", path)
					}
				}
			} else {
				if paths["App.tsx"] {
					t.Error("app-router mode produced unexpected SPA App.tsx")
				}
			}
		})
	}
}

// TestRenderLayoutSPAMultiModuleEmitsOnceAndMergesNav exercises the SPA
// multi-module short-circuit: the first module's call emits a unified
// App.tsx with both modules' nav items, prefix-namespaced; the second
// module's call returns no files (avoiding overwrite).
func TestRenderLayoutSPAMultiModuleEmitsOnceAndMergesNav(t *testing.T) {
	seShell := &config.Page{
		Sigil: "1.0", Kind: "page", ID: "se-shell", Title: "SE", Overlay: "page",
		Layout: config.Component{Type: "rows", Children: []config.Component{
			{Type: "nav-menu", Props: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"label": "Dashboard", "page": "se-dashboard", "icon": "home"},
				},
			}},
		}},
	}
	cwShell := &config.Page{
		Sigil: "1.0", Kind: "page", ID: "clockwork-shell", Title: "CW", Overlay: "page",
		Layout: config.Component{Type: "rows", Children: []config.Component{
			{Type: "nav-menu", Props: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"label": "Board", "page": "clockwork-board", "icon": "list"},
				},
			}},
		}},
	}
	seMod := &config.ModuleConfig{ID: "se", Shell: "se-shell", RouteGroup: "(explorer)", Pages: []string{"se-dashboard"}}
	cwMod := &config.ModuleConfig{ID: "clockwork", Shell: "clockwork-shell", RouteGroup: "(clockwork)", Pages: []string{"clockwork-board"}}
	modules := []*config.ModuleConfig{seMod, cwMod}
	shells := map[string]*config.Page{"se": seShell, "clockwork": cwShell}

	app := &config.AppConfig{
		Name:       "Multi",
		TargetMode: "spa",
		Modules:    []config.ModuleConfig{*seMod, *cwMod},
	}

	pages := []*config.Page{
		{Sigil: "1.0", Kind: "page", ID: "se-dashboard", Title: "SE Home", Overlay: "page", Module: "se", Layout: config.Component{Type: "rows"}},
		{Sigil: "1.0", Kind: "page", ID: "clockwork-board", Title: "Board", Overlay: "page", Module: "clockwork", Layout: config.Component{Type: "rows"}},
	}

	// First module: must emit App.tsx + routes.tsx with both modules' nav + routes.
	primary := &renderer.LayoutContext{
		Module: seMod, Shell: seShell, AppConfig: app,
		Pages: []*config.Page{pages[0]}, AllModules: modules, AllPages: pages, AllShells: shells,
	}
	files, err := renderLayout(primary)
	if err != nil {
		t.Fatalf("renderLayout(primary): %v", err)
	}
	paths := map[string]string{}
	for _, f := range files {
		paths[f.Path] = string(f.Content)
	}
	pathKeys := func(m map[string]string) []string {
		var out []string
		for k := range m {
			out = append(out, k)
		}
		return out
	}
	if _, ok := paths["App.tsx"]; !ok {
		t.Fatalf("primary call must emit App.tsx, got %v", pathKeys(paths))
	}
	if _, ok := paths["routes.tsx"]; !ok {
		t.Fatalf("primary call must emit routes.tsx, got %v", pathKeys(paths))
	}
	// Nav merges items from BOTH modules with their own prefixes.
	appTsx := paths["App.tsx"]
	for _, want := range []string{
		`"/explorer/dashboard"`,
		`"/clockwork/board"`,
		`"Dashboard"`,
		`"Board"`,
	} {
		if !strings.Contains(appTsx, want) {
			t.Errorf("App.tsx missing %q\n--- output ---\n%s", want, appTsx)
		}
	}

	// Second module: must NOT emit App.tsx again (would overwrite the unified one).
	secondary := &renderer.LayoutContext{
		Module: cwMod, Shell: cwShell, AppConfig: app,
		Pages: []*config.Page{pages[1]}, AllModules: modules, AllPages: pages, AllShells: shells,
	}
	files2, err := renderLayout(secondary)
	if err != nil {
		t.Fatalf("renderLayout(secondary): %v", err)
	}
	if len(files2) != 0 {
		var p []string
		for _, f := range files2 {
			p = append(p, f.Path)
		}
		t.Errorf("secondary SPA module must emit no layout files (App.tsx already emitted by primary); got %v", p)
	}
}

func TestEffectiveTargetModeDefaults(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", "app-router"},
		{"spa", "spa"},
		{"app-router", "app-router"},
	}
	for _, c := range cases {
		got := (&config.AppConfig{TargetMode: c.in}).EffectiveTargetMode()
		if got != c.want {
			t.Errorf("EffectiveTargetMode(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	var nilCfg *config.AppConfig
	if got := nilCfg.EffectiveTargetMode(); got != "app-router" {
		t.Errorf("nil AppConfig.EffectiveTargetMode() = %q, want app-router", got)
	}
}

// --- Provider declaration tests (sprint 10 phase 2) -------------------------

// shellWithNav returns a minimal shell page with one nav item so layout
// rendering proceeds end-to-end.
func shellWithNav() *config.Page {
	return &config.Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "x-shell",
		Title:   "Shell",
		Overlay: "page",
		Layout: config.Component{
			Type: "rows",
			Children: []config.Component{
				{
					Type: "nav-menu",
					Props: map[string]interface{}{
						"items": []interface{}{
							map[string]interface{}{"label": "Home", "page": "x-home", "icon": "home"},
						},
					},
				},
			},
		},
	}
}

func TestRenderProvidersBuiltinShortFormUnchanged(t *testing.T) {
	// The pre-phase-2 `lens` short-form (just datasource + ui hints) must
	// regenerate byte-identically to the way it did before custom-provider
	// support landed: a single lib/<id>-context.tsx file with createContext +
	// useEffect + fetchList wiring.
	providers := []config.ProviderConfig{
		{
			ID:         "lens",
			Datasource: "Lens",
			TrackField: "id",
			LabelField: "name",
			Position:   "topbar",
			UIType:     "select",
		},
	}
	files, err := renderProvidersForTests(providers, "" /* no sigil dir needed for built-ins */)
	if err != nil {
		t.Fatalf("renderProviders: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file for built-in provider, got %d", len(files))
	}
	if files[0].Path != "lib/lens-context.tsx" {
		t.Errorf("expected lib/lens-context.tsx, got %q", files[0].Path)
	}
	content := string(files[0].Content)
	mustHave := []string{
		`"use client";`,
		`fetchList`,
		`createContext`,
		`export function useLensContext()`,
		`export function LensProvider({ children }`,
	}
	for _, want := range mustHave {
		if !strings.Contains(content, want) {
			t.Errorf("built-in provider missing %q in output:\n%s", want, content)
		}
	}
}

func TestRenderProvidersCustomCopiesSourceFiles(t *testing.T) {
	dir := t.TempDir()
	if err := writeTestFile(dir, "providers/active-runs-context.tsx",
		"// active runs context\nexport function ActiveRunsProvider() { return null; }\n"); err != nil {
		t.Fatal(err)
	}
	if err := writeTestFile(dir, "providers/sse-active-runs-bridge.tsx",
		"// bridge\nexport function SSEActiveRunsBridge() { return null; }\n"); err != nil {
		t.Fatal(err)
	}

	providers := []config.ProviderConfig{
		{
			ID: "active-runs",
			Source: &config.ProviderSource{
				Component: "providers/active-runs-context.tsx",
				Export:    "ActiveRunsProvider",
				Includes:  []string{"providers/sse-active-runs-bridge.tsx"},
			},
			Mounts: []string{"SSEActiveRunsBridge"},
		},
	}
	files, err := renderProvidersForTests(providers, dir)
	if err != nil {
		t.Fatalf("renderProviders: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files (main + include), got %d", len(files))
	}
	got := map[string]string{}
	for _, f := range files {
		got[f.Path] = string(f.Content)
	}
	if _, ok := got["lib/active-runs-context.tsx"]; !ok {
		t.Errorf("missing lib/active-runs-context.tsx; got files: %v", keys(got))
	}
	if _, ok := got["lib/sse-active-runs-bridge.tsx"]; !ok {
		t.Errorf("missing lib/sse-active-runs-bridge.tsx; got files: %v", keys(got))
	}
	// Verify source content was preserved verbatim (the renderer must not
	// rewrite hand-authored bodies).
	if !strings.Contains(got["lib/active-runs-context.tsx"], "ActiveRunsProvider") {
		t.Errorf("provider source content not preserved: %q", got["lib/active-runs-context.tsx"])
	}
}

func TestRenderProvidersCustomBackwardCompatMixed(t *testing.T) {
	// A providers: list that contains both a built-in (lens) AND a custom
	// (active-runs) provider should generate one file per built-in plus the
	// copied source files for each custom provider — in declared order.
	dir := t.TempDir()
	if err := writeTestFile(dir, "providers/active-runs-context.tsx", "export function ActiveRunsProvider() { return null; }\n"); err != nil {
		t.Fatal(err)
	}

	providers := []config.ProviderConfig{
		{ID: "lens", Datasource: "Lens", TrackField: "id", LabelField: "name"},
		{ID: "active-runs", Source: &config.ProviderSource{Component: "providers/active-runs-context.tsx"}},
	}
	files, err := renderProvidersForTests(providers, dir)
	if err != nil {
		t.Fatalf("renderProviders: %v", err)
	}
	paths := map[string]bool{}
	for _, f := range files {
		paths[f.Path] = true
	}
	if !paths["lib/lens-context.tsx"] {
		t.Errorf("expected lib/lens-context.tsx, got %v", paths)
	}
	if !paths["lib/active-runs-context.tsx"] {
		t.Errorf("expected lib/active-runs-context.tsx, got %v", paths)
	}
}

func TestRenderProvidersCustomMissingSourceErrors(t *testing.T) {
	providers := []config.ProviderConfig{
		{
			ID:     "active-runs",
			Source: &config.ProviderSource{Component: "providers/does-not-exist.tsx"},
		},
	}
	dir := t.TempDir()
	_, err := renderProvidersForTests(providers, dir)
	if err == nil {
		t.Fatal("expected error for missing source file, got nil")
	}
	if !strings.Contains(err.Error(), "does-not-exist") {
		t.Errorf("error message should reference missing source file; got: %v", err)
	}
}

func TestRenderLayoutCustomProviderWrapsWithMounts(t *testing.T) {
	shell := shellWithNav()
	// Phase 3.5: providers live on ModuleConfig.
	mod := &config.ModuleConfig{
		ID: "x", Shell: "x-shell", RouteGroup: "(x)", Pages: []string{"x-home"},
		Providers: []config.ProviderConfig{
			{
				ID: "active-runs",
				Source: &config.ProviderSource{
					Component: "providers/active-runs-context.tsx",
					Export:    "ActiveRunsProvider",
					Includes:  []string{"providers/sse-active-runs-bridge.tsx"},
				},
				Mounts: []string{"SSEActiveRunsBridge"},
			},
		},
	}

	cases := []struct {
		name        string
		targetMode  string
		mustContain []string
		mustNotHave []string
	}{
		{
			name:       "app-router custom provider",
			targetMode: "app-router",
			mustContain: []string{
				`import { ActiveRunsProvider } from "@/lib/active-runs-context";`,
				`import { SSEActiveRunsBridge } from "@/lib/sse-active-runs-bridge";`,
				`<ActiveRunsProvider>`,
				`<SSEActiveRunsBridge />`,
				`</ActiveRunsProvider>`,
			},
			mustNotHave: []string{
				`import { ActiveRunsProvider, use`, // no auto-hook for custom providers
			},
		},
		{
			name:       "spa custom provider",
			targetMode: "spa",
			mustContain: []string{
				`import { ActiveRunsProvider } from "@/lib/active-runs-context";`,
				`import { SSEActiveRunsBridge } from "@/lib/sse-active-runs-bridge";`,
				`<ActiveRunsProvider>`,
				`<SSEActiveRunsBridge />`,
				`<BrowserRouter>`,
				`</ActiveRunsProvider>`,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &renderer.LayoutContext{
				Module: mod,
				Shell:  shell,
				AppConfig: &config.AppConfig{
					Name:       "X",
					TargetMode: tc.targetMode,
					Modules:    []config.ModuleConfig{*mod},
				},
				Pages: []*config.Page{},
			}
			files, err := renderLayout(ctx)
			if err != nil {
				t.Fatalf("renderLayout: %v", err)
			}
			// Combine all emitted files so the assertions don't depend on
			// which file the wrap lives in (App.tsx for spa, layout.tsx for
			// app-router).
			var combined strings.Builder
			for _, f := range files {
				combined.WriteString(f.Path)
				combined.WriteByte('\n')
				combined.Write(f.Content)
				combined.WriteString("\n")
			}
			content := combined.String()
			for _, want := range tc.mustContain {
				if !strings.Contains(content, want) {
					t.Errorf("expected %q in output:\n%s", want, content)
				}
			}
			for _, banned := range tc.mustNotHave {
				if strings.Contains(content, banned) {
					t.Errorf("did not expect %q in output:\n%s", banned, content)
				}
			}
		})
	}
}

func TestRenderLayoutBuiltinProviderShortFormStillWorks(t *testing.T) {
	// Backward-compat: a providers: list with only the original lens
	// short-form must continue to import + wrap with LensProvider + the
	// useLensContext hook (the topbar select uses it).
	// Phase 3.5: providers live per-module.
	shell := shellWithNav()
	mod := &config.ModuleConfig{
		ID: "x", Shell: "x-shell", RouteGroup: "(x)", Pages: []string{"x-home"},
		Providers: []config.ProviderConfig{{
			ID:         "lens",
			Datasource: "Lens",
			TrackField: "id",
			LabelField: "name",
			Position:   "topbar",
			UIType:     "select",
		}},
	}

	ctx := &renderer.LayoutContext{
		Module: mod,
		Shell:  shell,
		AppConfig: &config.AppConfig{
			Name:    "X",
			Modules: []config.ModuleConfig{*mod},
		},
		Pages: []*config.Page{},
	}
	files, err := renderLayout(ctx)
	if err != nil {
		t.Fatalf("renderLayout: %v", err)
	}
	var combined strings.Builder
	for _, f := range files {
		combined.Write(f.Content)
	}
	content := combined.String()
	mustHave := []string{
		`import { LensProvider, useLensContext } from "@/lib/lens-context";`,
		`<LensProvider>`,
		`</LensProvider>`,
		`useLensContext()`, // consumed by the topbar select
	}
	for _, want := range mustHave {
		if !strings.Contains(content, want) {
			t.Errorf("expected %q in output:\n%s", want, content)
		}
	}
}

func TestProviderWrapNameAndImportPath(t *testing.T) {
	// Built-in: derived from ID.
	builtin := config.ProviderConfig{ID: "lens", Datasource: "Lens"}
	if got := providerWrapName(&builtin); got != "LensProvider" {
		t.Errorf("builtin providerWrapName(lens) = %q, want LensProvider", got)
	}
	if got := providerImportPath(&builtin); got != "@/lib/lens-context" {
		t.Errorf("builtin providerImportPath(lens) = %q, want @/lib/lens-context", got)
	}

	// Custom with explicit Export.
	custom := config.ProviderConfig{
		ID:     "active-runs",
		Source: &config.ProviderSource{Component: "providers/active-runs-context.tsx", Export: "ActiveRunsProvider"},
	}
	if got := providerWrapName(&custom); got != "ActiveRunsProvider" {
		t.Errorf("custom providerWrapName = %q, want ActiveRunsProvider", got)
	}
	if got := providerImportPath(&custom); got != "@/lib/active-runs-context" {
		t.Errorf("custom providerImportPath = %q, want @/lib/active-runs-context", got)
	}

	// Custom without explicit Export — defaults to PascalCase(ID)+"Provider".
	customNoExport := config.ProviderConfig{
		ID:     "active-runs",
		Source: &config.ProviderSource{Component: "providers/active-runs-context.tsx"},
	}
	if got := providerWrapName(&customNoExport); got != "ActiveRunsProvider" {
		t.Errorf("custom default wrap name = %q, want ActiveRunsProvider", got)
	}
}

// --- Phase 3.5 per-module providers tests -----------------------------------

// TestRenderProvidersDedupAcrossModules confirms that two modules declaring
// the same provider id produce only one set of `lib/` files. This is the
// core dedup contract: shared providers (e.g. `lens` in both `se` and
// `clockwork` modules) get one source-file copy per app.
func TestRenderProvidersDedupAcrossModules(t *testing.T) {
	dir := t.TempDir()
	if err := writeTestFile(dir, "providers/active-runs-context.tsx",
		"export function ActiveRunsProvider() { return null; }\n"); err != nil {
		t.Fatal(err)
	}

	// Two modules, each declaring the same `lens` builtin AND the same
	// custom `active-runs` provider. Expected output: one lens-context file
	// and one active-runs-context file — no duplicates.
	modules := []config.ModuleConfig{
		{
			ID: "se", Shell: "se-shell", Pages: []string{"p"},
			Providers: []config.ProviderConfig{
				{ID: "lens", Datasource: "Lens", TrackField: "id", LabelField: "name"},
			},
		},
		{
			ID: "cw", Shell: "cw-shell", Pages: []string{"p"},
			Providers: []config.ProviderConfig{
				{ID: "lens", Datasource: "Lens", TrackField: "id", LabelField: "name"},
				{ID: "active-runs", Source: &config.ProviderSource{
					Component: "providers/active-runs-context.tsx",
					Export:    "ActiveRunsProvider",
				}},
			},
		},
	}

	files, err := renderProviders(modules, dir)
	if err != nil {
		t.Fatalf("renderProviders: %v", err)
	}
	paths := map[string]int{}
	for _, f := range files {
		paths[f.Path]++
	}
	if paths["lib/lens-context.tsx"] != 1 {
		t.Errorf("expected lib/lens-context.tsx emitted exactly once across modules; got %d (all paths: %v)", paths["lib/lens-context.tsx"], paths)
	}
	if paths["lib/active-runs-context.tsx"] != 1 {
		t.Errorf("expected lib/active-runs-context.tsx emitted exactly once; got %d (all paths: %v)", paths["lib/active-runs-context.tsx"], paths)
	}
}

// TestRenderLayoutPerModuleProviders confirms that in App Router mode each
// module's layout.tsx wraps its OWN providers — and only its own. A second
// module's providers must not leak into the first module's layout.
func TestRenderLayoutPerModuleProviders(t *testing.T) {
	shell := shellWithNav()
	cwMod := config.ModuleConfig{
		ID: "cw", Shell: "cw-shell", RouteGroup: "(cw)", Pages: []string{"cw-home"},
		Providers: []config.ProviderConfig{
			{ID: "active-runs", Source: &config.ProviderSource{
				Component: "providers/active-runs-context.tsx",
				Export:    "ActiveRunsProvider",
			}},
		},
	}
	seMod := config.ModuleConfig{
		ID: "se", Shell: "se-shell", RouteGroup: "(se)", Pages: []string{"se-home"},
		// Intentionally NO providers — SE module should not get
		// active-runs even though clockwork has it.
	}

	app := &config.AppConfig{
		Name:    "Multi",
		Modules: []config.ModuleConfig{seMod, cwMod},
	}

	// SE module's layout: must NOT wrap with ActiveRunsProvider.
	seCtx := &renderer.LayoutContext{
		Module: &app.Modules[0], Shell: shell, AppConfig: app, Pages: []*config.Page{},
	}
	seFiles, err := renderLayout(seCtx)
	if err != nil {
		t.Fatalf("renderLayout(se): %v", err)
	}
	var seBuf strings.Builder
	for _, f := range seFiles {
		seBuf.Write(f.Content)
	}
	seOut := seBuf.String()
	if strings.Contains(seOut, "ActiveRunsProvider") {
		t.Errorf("SE module layout must NOT contain ActiveRunsProvider (per-module providers leak):\n%s", seOut)
	}

	// CW module's layout: must wrap with ActiveRunsProvider.
	cwCtx := &renderer.LayoutContext{
		Module: &app.Modules[1], Shell: shell, AppConfig: app, Pages: []*config.Page{},
	}
	cwFiles, err := renderLayout(cwCtx)
	if err != nil {
		t.Fatalf("renderLayout(cw): %v", err)
	}
	var cwBuf strings.Builder
	for _, f := range cwFiles {
		cwBuf.Write(f.Content)
	}
	cwOut := cwBuf.String()
	if !strings.Contains(cwOut, "<ActiveRunsProvider>") {
		t.Errorf("CW module layout must contain <ActiveRunsProvider>:\n%s", cwOut)
	}
}

// TestValidatorRejectsDuplicateProviderIDsWithinModule was added with the
// validator changes; this just confirms the renderer accepts a single ID
// declared per-module across modules.
func TestProviderConfigFromModuleSurface(t *testing.T) {
	// Per-module Providers on ModuleConfig should round-trip through
	// renderProviders without complaint.
	mods := []config.ModuleConfig{
		{ID: "m1", Shell: "m1-shell", Pages: []string{"p"},
			Providers: []config.ProviderConfig{
				{ID: "lens", Datasource: "Lens", TrackField: "id", LabelField: "name"},
			},
		},
	}
	files, err := renderProviders(mods, "")
	if err != nil {
		t.Fatalf("renderProviders: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file from a single builtin provider, got %d", len(files))
	}
}

func TestProviderMountsMapsKebabFileNames(t *testing.T) {
	p := config.ProviderConfig{
		ID: "active-runs",
		Source: &config.ProviderSource{
			Component: "providers/active-runs-context.tsx",
			Includes:  []string{"providers/sse-active-runs-bridge.tsx"},
		},
		Mounts: []string{"SSEActiveRunsBridge"},
	}
	mounts := providerMounts(&p)
	if len(mounts) != 1 {
		t.Fatalf("expected 1 mount, got %d", len(mounts))
	}
	if mounts[0].Name != "SSEActiveRunsBridge" {
		t.Errorf("mount.Name = %q, want SSEActiveRunsBridge", mounts[0].Name)
	}
	if mounts[0].Path != "@/lib/sse-active-runs-bridge" {
		t.Errorf("mount.Path = %q, want @/lib/sse-active-runs-bridge", mounts[0].Path)
	}
}

// keys returns the keys of m sorted by Go's range order (sufficient for
// debugging assertions; we don't need stable ordering here).
func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func writeTestFile(rootDir, relPath, content string) error {
	full := filepath.Join(rootDir, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o750); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), 0o600)
}

// TestRenderModalRefetchAfterCreate is a regression test for the codegen bug
// surfaced by the FND-4 sysop dogfood loop: a modal that creates a row never
// refreshed the table, and the `refetch<Alias>` binding was emitted only when
// the page YAML carried an explicit `refresh:` field — leaving it destructured
// but unused (a `noUnusedLocals` failure) on the common case. refetch is now
// driven by the datasource's create/update/delete capability.
func TestRenderModalRefetchAfterCreate(t *testing.T) {
	page := &config.Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      "things",
		Title:   "Things",
		Overlay: "page",
		DataSources: []config.DataSourceRef{
			{Alias: "Thing", Capabilities: []string{"list", "create"}},
		},
		Layout: config.Component{
			ID:   "root",
			Type: "rows",
			Children: []config.Component{
				{
					Type:  "button",
					Props: map[string]interface{}{"label": "New Thing"},
					Actions: map[string]config.Action{
						"click": {
							Type:   "modal",
							Title:  "New Thing",
							Fields: []config.FormField{{Name: "name", Label: "Name", Type: "text"}},
							Submit: &config.SubmitConfig{DataSource: "Thing", Method: "POST"},
						},
					},
				},
				{Type: "data-table", Props: map[string]interface{}{"datasource": "Thing"}},
			},
		},
	}
	ctx := testCtx(page)
	ctx.DataSources = map[string]*renderer.DataSourceManifest{
		"Thing": {
			Alias:        "Thing",
			Capabilities: []string{"list", "create"},
			Endpoints:    map[string]string{"list": "/api/things", "create": "/api/things"},
		},
	}

	files, err := renderPage(ctx)
	if err != nil {
		t.Fatalf("renderPage: %v", err)
	}
	content := string(files[0].Content)
	if !strings.Contains(content, "refetch: refetchThing") {
		t.Errorf("expected the list hook to be destructured with a refetch binding:\n%s", content)
	}
	if !strings.Contains(content, "refetchThing();") {
		t.Errorf("expected the create handler to call refetchThing():\n%s", content)
	}
}

// TestGenerateDSHookRefetch verifies the generated SWR list hook exposes a
// `refetch` alias (SWR itself ships `mutate`, not `refetch`) and that the
// `<Type>Input` type import is emitted only when a mutation helper consumes it.
func TestGenerateDSHookRefetch(t *testing.T) {
	crud := &renderer.DataSourceManifest{
		Alias:        "Thing",
		Capabilities: []string{"list", "create"},
		Endpoints:    map[string]string{"list": "/api/things", "create": "/api/things"},
	}
	crudOut := string(generateDSHook(crud).Content)
	if !strings.Contains(crudOut, "refetch: () => swr.mutate()") {
		t.Errorf("expected list hook to expose a refetch alias:\n%s", crudOut)
	}
	if !strings.Contains(crudOut, "ThingInput") {
		t.Errorf("expected ThingInput import for a create-capable datasource:\n%s", crudOut)
	}

	listOnly := &renderer.DataSourceManifest{
		Alias:        "Note",
		Capabilities: []string{"list"},
		Endpoints:    map[string]string{"list": "/api/notes"},
	}
	listOut := string(generateDSHook(listOnly).Content)
	if !strings.Contains(listOut, "refetch: () => swr.mutate()") {
		t.Errorf("expected list-only hook to still expose a refetch alias:\n%s", listOut)
	}
	if strings.Contains(listOut, "NoteInput") {
		t.Errorf("list-only datasource should not import the unused NoteInput type:\n%s", listOut)
	}
}
