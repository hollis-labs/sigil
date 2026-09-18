package mcp

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// setupTestProject creates a minimal .sigil project in a temp directory
// and returns a fully configured MCP server.
func setupTestProject(t *testing.T) (*Server, string) {
	t.Helper()

	root := t.TempDir()
	sigilDir := filepath.Join(root, ".sigil")

	for _, dir := range []string{
		filepath.Join(sigilDir, "pages"),
		filepath.Join(sigilDir, "datasources"),
		filepath.Join(sigilDir, "themes"),
	} {
		os.MkdirAll(dir, 0755)
	}

	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: test-project
defaults:
  theme: default
  renderer: go-templ
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
description: "Default dark theme"
tokens:
  colors:
    background: "9 9 21"
    text: "244 244 245"
    accent: "16 185 129"
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "pages", "dashboard.yaml"), []byte(`
sigil: "1.0"
id: dashboard
title: Dashboard
overlay: page
layout:
  id: root
  type: rows
  children:
    - id: title
      type: heading
      props:
        level: 2
        text: Dashboard
`), 0644)

	os.WriteFile(filepath.Join(sigilDir, "datasources", "sprint.yaml"), []byte(`
alias: Sprint
description: "Sprint management"
capabilities: [search, create, update, delete]
fields:
  - name: id
    type: integer
    primary: true
  - name: title
    type: string
    required: true
endpoints:
  list: "GET /api/sprints"
  create: "POST /api/sprints"
`), 0644)

	s := NewServer(sigilDir)
	RegisterAllTools(s)
	RegisterAllResources(s)
	RegisterAllPrompts(s)

	return s, root
}

// connectClient wires an in-memory MCP client to s over the real protocol
// and returns the connected client session. Both sides are closed
// automatically at test cleanup.
func connectClient(t *testing.T, s *Server) *mcpsdk.ClientSession {
	t.Helper()
	ctx := context.Background()

	serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()

	serverSession, err := s.SDKServer().Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	t.Cleanup(func() { serverSession.Close() })

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "test-client", Version: "1.0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { clientSession.Close() })

	return clientSession
}

func callTool(t *testing.T, cs *mcpsdk.ClientSession, name string, args map[string]any) *mcpsdk.CallToolResult {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcpsdk.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		t.Fatalf("call tool %s: %v", name, err)
	}
	return res
}

// extractToolText returns the text of a tool result's first content block.
func extractToolText(t *testing.T, res *mcpsdk.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		return ""
	}
	tc, ok := res.Content[0].(*mcpsdk.TextContent)
	if !ok {
		t.Fatalf("expected text content, got %T", res.Content[0])
	}
	return tc.Text
}

// promptText returns the text of a prompt result's first message.
func promptText(t *testing.T, res *mcpsdk.GetPromptResult) string {
	t.Helper()
	if len(res.Messages) == 0 {
		t.Fatal("expected at least one message")
	}
	tc, ok := res.Messages[0].Content.(*mcpsdk.TextContent)
	if !ok {
		t.Fatalf("expected text content, got %T", res.Messages[0].Content)
	}
	return tc.Text
}

func TestIntegrationToolsList(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res, err := cs.ListTools(context.Background(), &mcpsdk.ListToolsParams{})
	if err != nil {
		t.Fatalf("tools/list error: %v", err)
	}

	expectedTools := []string{
		"sigil_list_pages", "sigil_get_page", "sigil_create_page",
		"sigil_update_page", "sigil_validate", "sigil_list_components",
		"sigil_get_component_schema", "sigil_list_datasources",
		"sigil_create_datasource",
	}

	toolNames := map[string]bool{}
	for _, tl := range res.Tools {
		toolNames[tl.Name] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("expected tool %q in tools list", expected)
		}
	}

	if len(res.Tools) != 9 {
		t.Errorf("expected 9 tools, got %d", len(res.Tools))
	}
}

func TestIntegrationListPages(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res := callTool(t, cs, "sigil_list_pages", map[string]any{})
	if res.IsError {
		t.Fatalf("unexpected error: %s", extractToolText(t, res))
	}

	text := extractToolText(t, res)
	if !strings.Contains(text, "dashboard") {
		t.Error("expected dashboard in pages list")
	}
	if !strings.Contains(text, "Dashboard") {
		t.Error("expected Dashboard title")
	}
}

func TestIntegrationGetPage(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res := callTool(t, cs, "sigil_get_page", map[string]any{"id": "dashboard"})
	if res.IsError {
		t.Fatalf("unexpected error: %s", extractToolText(t, res))
	}

	text := extractToolText(t, res)
	if !strings.Contains(text, "sigil:") {
		t.Error("expected YAML content")
	}
	if !strings.Contains(text, "dashboard") {
		t.Error("expected dashboard content")
	}
}

func TestIntegrationCreatePage(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	pageConfig := `sigil: "1.0"
id: new-page
title: New Page
overlay: page
layout:
  id: root
  type: rows
  children:
    - id: heading
      type: heading
      props:
        level: 2
        text: New Page`

	res := callTool(t, cs, "sigil_create_page", map[string]any{
		"id":     "new-page",
		"config": pageConfig,
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", extractToolText(t, res))
	}

	text := extractToolText(t, res)
	if !strings.Contains(text, `"valid":true`) {
		t.Errorf("expected valid: true, got: %s", text)
	}
	if !strings.Contains(text, "new-page.yaml") {
		t.Error("expected file path in response")
	}

	path := filepath.Join(s.SigilDir(), "pages", "new-page.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to be created at %s", path)
	}
}

func TestIntegrationCreatePageInvalid(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res := callTool(t, cs, "sigil_create_page", map[string]any{
		"id":     "bad-page",
		"config": "id: bad\nlayout:\n  type: rows\n",
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", extractToolText(t, res))
	}

	text := extractToolText(t, res)
	if !strings.Contains(text, `"valid":false`) {
		t.Errorf("expected valid: false for invalid config, got: %s", text)
	}

	path := filepath.Join(s.SigilDir(), "pages", "bad-page.yaml")
	if _, err := os.Stat(path); err == nil {
		t.Error("expected invalid page to NOT be saved")
	}
}

func TestIntegrationValidate(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	validConfig := `sigil: "1.0"
id: test
title: Test
overlay: page
layout:
  id: root
  type: rows`

	res := callTool(t, cs, "sigil_validate", map[string]any{"config": validConfig})
	if res.IsError {
		t.Fatalf("unexpected error: %s", extractToolText(t, res))
	}

	text := extractToolText(t, res)
	if !strings.Contains(text, `"valid":true`) {
		t.Errorf("expected valid: true, got: %s", text)
	}
}

func TestIntegrationValidateInvalid(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res := callTool(t, cs, "sigil_validate", map[string]any{
		"config": "id: \ntitle: \noverlay: bad\nlayout:\n  type: \n",
	})

	text := extractToolText(t, res)
	if !strings.Contains(text, `"valid":false`) {
		t.Errorf("expected valid: false, got: %s", text)
	}
}

func TestIntegrationListComponents(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res := callTool(t, cs, "sigil_list_components", map[string]any{})
	text := extractToolText(t, res)
	if !strings.Contains(text, "button") {
		t.Error("expected button component")
	}
	if !strings.Contains(text, "heading") {
		t.Error("expected heading component")
	}
	if !strings.Contains(text, "data-table") {
		t.Error("expected data-table component")
	}
}

func TestIntegrationListComponentsByCategory(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res := callTool(t, cs, "sigil_list_components", map[string]any{"category": "layouts"})
	text := extractToolText(t, res)
	if !strings.Contains(text, "rows") {
		t.Error("expected rows in layouts")
	}
	if !strings.Contains(text, "columns") {
		t.Error("expected columns in layouts")
	}
	if strings.Contains(text, `"category": "primitives"`) {
		t.Error("unexpected primitives in layouts filter")
	}
}

func TestIntegrationGetComponentSchema(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res := callTool(t, cs, "sigil_get_component_schema", map[string]any{"type": "button"})
	text := extractToolText(t, res)
	if !strings.Contains(text, "button") {
		t.Error("expected button type")
	}
	if !strings.Contains(text, "props") {
		t.Error("expected props in schema")
	}
}

func TestIntegrationListDataSources(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res := callTool(t, cs, "sigil_list_datasources", map[string]any{})
	text := extractToolText(t, res)
	if !strings.Contains(text, "Sprint") {
		t.Error("expected Sprint datasource")
	}
}

func TestIntegrationCreateDataSource(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	dsConfig := `alias: Task
description: "Task management"
capabilities: [search, create, update, delete]
fields:
  - name: id
    type: integer
    primary: true
  - name: title
    type: string
    required: true`

	res := callTool(t, cs, "sigil_create_datasource", map[string]any{
		"alias":  "Task",
		"config": dsConfig,
	})
	text := extractToolText(t, res)
	if !strings.Contains(text, "task.yaml") {
		t.Errorf("expected file path, got: %s", text)
	}

	path := filepath.Join(s.SigilDir(), "datasources", "task.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected datasource file at %s", path)
	}
}

func TestIntegrationResourcesRead(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	tests := []struct {
		uri      string
		contains string
	}{
		{"sigil://pages", "dashboard"},
		{"sigil://pages/dashboard", "sigil:"},
		{"sigil://components", "button"},
		{"sigil://components/heading", "heading"},
		{"sigil://datasources", "Sprint"},
		{"sigil://themes", "default"},
		{"sigil://project", "test-project"},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			res, err := cs.ReadResource(context.Background(), &mcpsdk.ReadResourceParams{URI: tt.uri})
			if err != nil {
				t.Fatalf("error reading %s: %v", tt.uri, err)
			}
			if len(res.Contents) == 0 {
				t.Fatalf("expected contents for %s", tt.uri)
			}
			if !strings.Contains(res.Contents[0].Text, tt.contains) {
				t.Errorf("expected %q in result for %s, got: %s", tt.contains, tt.uri, res.Contents[0].Text)
			}
		})
	}
}

func TestIntegrationPromptsDesignPage(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	res, err := cs.GetPrompt(context.Background(), &mcpsdk.GetPromptParams{
		Name:      "sigil_design_page",
		Arguments: map[string]string{"description": "A user management page"},
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	text := promptText(t, res)
	if !strings.Contains(text, "user management") {
		t.Error("expected description in prompt")
	}
	if !strings.Contains(text, "button") {
		t.Error("expected component catalog in prompt")
	}
	if !strings.Contains(text, "sigil:") {
		t.Error("expected config template in prompt")
	}
}

func TestIntegrationPromptsReviewConfig(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	validConfig := `sigil: "1.0"
id: test
title: Test
overlay: page
layout:
  id: root
  type: rows`

	res, err := cs.GetPrompt(context.Background(), &mcpsdk.GetPromptParams{
		Name:      "sigil_review_config",
		Arguments: map[string]string{"config": validConfig},
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	text := promptText(t, res)
	if !strings.Contains(text, "PASS") {
		t.Error("expected PASS for valid config")
	}
}

func TestIntegrationEndToEnd(t *testing.T) {
	s, _ := setupTestProject(t)
	cs := connectClient(t, s)

	// 1. List components to understand what's available
	res := callTool(t, cs, "sigil_list_components", map[string]any{"category": "primitives"})
	text := extractToolText(t, res)
	if !strings.Contains(text, "button") {
		t.Fatal("expected button in components")
	}

	// 2. Create a page
	pageConfig := `sigil: "1.0"
id: user-list
title: User List
overlay: page
layout:
  id: root
  type: rows
  props:
    gap: 4
  children:
    - id: header
      type: heading
      props:
        level: 2
        text: Users
    - id: add-btn
      type: button
      props:
        label: Add User
        variant: primary`

	res = callTool(t, cs, "sigil_create_page", map[string]any{
		"id":     "user-list",
		"config": pageConfig,
	})
	text = extractToolText(t, res)
	if !strings.Contains(text, `"valid":true`) {
		t.Fatalf("expected valid page, got: %s", text)
	}

	// 3. Validate the created page
	res = callTool(t, cs, "sigil_validate", map[string]any{
		"path": filepath.Join(s.SigilDir(), "pages", "user-list.yaml"),
	})
	text = extractToolText(t, res)
	if !strings.Contains(text, `"valid":true`) {
		t.Fatalf("expected validation pass, got: %s", text)
	}

	// 4. List pages — should include both dashboard and user-list
	res = callTool(t, cs, "sigil_list_pages", map[string]any{})
	text = extractToolText(t, res)
	if !strings.Contains(text, "dashboard") {
		t.Error("expected dashboard in pages list")
	}
	if !strings.Contains(text, "user-list") {
		t.Error("expected user-list in pages list")
	}

	// 5. Get the created page
	res = callTool(t, cs, "sigil_get_page", map[string]any{"id": "user-list"})
	text = extractToolText(t, res)
	if !strings.Contains(text, "User List") {
		t.Error("expected page title in content")
	}
}
