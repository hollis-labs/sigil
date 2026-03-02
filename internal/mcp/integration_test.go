package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
    background: "9 9 11"
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

func callTool(t *testing.T, s *Server, name string, args interface{}) *jsonRPCResponse {
	t.Helper()
	return sendRequest(t, s, "tools/call", 1, map[string]interface{}{
		"name":      name,
		"arguments": args,
	})
}

func TestIntegrationInitialize(t *testing.T) {
	s, _ := setupTestProject(t)
	resp := sendRequest(t, s, "initialize", 1, map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]interface{}{"name": "test-client", "version": "1.0"},
	})

	if resp.Error != nil {
		t.Fatalf("initialize error: %v", resp.Error)
	}

	result := resp.Result.(map[string]interface{})
	caps := result["capabilities"].(map[string]interface{})
	if _, ok := caps["tools"]; !ok {
		t.Error("expected tools capability")
	}
	if _, ok := caps["resources"]; !ok {
		t.Error("expected resources capability")
	}
	if _, ok := caps["prompts"]; !ok {
		t.Error("expected prompts capability")
	}
}

func TestIntegrationToolsList(t *testing.T) {
	s, _ := setupTestProject(t)
	resp := sendRequest(t, s, "tools/list", 1, nil)

	if resp.Error != nil {
		t.Fatalf("tools/list error: %v", resp.Error)
	}

	result := resp.Result.(map[string]interface{})
	tools := result["tools"].([]interface{})

	expectedTools := []string{
		"sigil_list_pages", "sigil_get_page", "sigil_create_page",
		"sigil_update_page", "sigil_validate", "sigil_list_components",
		"sigil_get_component_schema", "sigil_list_datasources",
		"sigil_create_datasource",
	}

	toolNames := map[string]bool{}
	for _, t := range tools {
		tm := t.(map[string]interface{})
		toolNames[tm["name"].(string)] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("expected tool %q in tools list", expected)
		}
	}

	if len(tools) != 9 {
		t.Errorf("expected 9 tools, got %d", len(tools))
	}
}

func TestIntegrationListPages(t *testing.T) {
	s, _ := setupTestProject(t)
	resp := callTool(t, s, "sigil_list_pages", map[string]interface{}{})

	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}

	text := extractToolText(t, resp)
	if !strings.Contains(text, "dashboard") {
		t.Error("expected dashboard in pages list")
	}
	if !strings.Contains(text, "Dashboard") {
		t.Error("expected Dashboard title")
	}
}

func TestIntegrationGetPage(t *testing.T) {
	s, _ := setupTestProject(t)
	resp := callTool(t, s, "sigil_get_page", map[string]interface{}{
		"id": "dashboard",
	})

	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}

	text := extractToolText(t, resp)
	if !strings.Contains(text, "sigil:") {
		t.Error("expected YAML content")
	}
	if !strings.Contains(text, "dashboard") {
		t.Error("expected dashboard content")
	}
}

func TestIntegrationCreatePage(t *testing.T) {
	s, _ := setupTestProject(t)

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

	resp := callTool(t, s, "sigil_create_page", map[string]interface{}{
		"id":     "new-page",
		"config": pageConfig,
	})

	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}

	text := extractToolText(t, resp)
	if !strings.Contains(text, `"valid": true`) {
		t.Errorf("expected valid: true, got: %s", text)
	}
	if !strings.Contains(text, "new-page.yaml") {
		t.Error("expected file path in response")
	}

	// Verify file was created
	path := filepath.Join(s.SigilDir(), "pages", "new-page.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to be created at %s", path)
	}
}

func TestIntegrationCreatePageInvalid(t *testing.T) {
	s, _ := setupTestProject(t)

	// Missing required fields
	resp := callTool(t, s, "sigil_create_page", map[string]interface{}{
		"id":     "bad-page",
		"config": "id: bad\nlayout:\n  type: rows\n",
	})

	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}

	text := extractToolText(t, resp)
	if !strings.Contains(text, `"valid": false`) {
		t.Errorf("expected valid: false for invalid config, got: %s", text)
	}

	// File should NOT be created
	path := filepath.Join(s.SigilDir(), "pages", "bad-page.yaml")
	if _, err := os.Stat(path); err == nil {
		t.Error("expected invalid page to NOT be saved")
	}
}

func TestIntegrationValidate(t *testing.T) {
	s, _ := setupTestProject(t)

	validConfig := `sigil: "1.0"
id: test
title: Test
overlay: page
layout:
  id: root
  type: rows`

	resp := callTool(t, s, "sigil_validate", map[string]interface{}{
		"config": validConfig,
	})

	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}

	text := extractToolText(t, resp)
	if !strings.Contains(text, `"valid": true`) {
		t.Errorf("expected valid: true, got: %s", text)
	}
}

func TestIntegrationValidateInvalid(t *testing.T) {
	s, _ := setupTestProject(t)

	resp := callTool(t, s, "sigil_validate", map[string]interface{}{
		"config": "id: \ntitle: \noverlay: bad\nlayout:\n  type: \n",
	})

	text := extractToolText(t, resp)
	if !strings.Contains(text, `"valid": false`) {
		t.Errorf("expected valid: false, got: %s", text)
	}
}

func TestIntegrationListComponents(t *testing.T) {
	s, _ := setupTestProject(t)
	resp := callTool(t, s, "sigil_list_components", map[string]interface{}{})

	text := extractToolText(t, resp)
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
	resp := callTool(t, s, "sigil_list_components", map[string]interface{}{
		"category": "layouts",
	})

	text := extractToolText(t, resp)
	if !strings.Contains(text, "rows") {
		t.Error("expected rows in layouts")
	}
	if !strings.Contains(text, "columns") {
		t.Error("expected columns in layouts")
	}
	// Should NOT include primitives
	if strings.Contains(text, `"category": "primitives"`) {
		t.Error("unexpected primitives in layouts filter")
	}
}

func TestIntegrationGetComponentSchema(t *testing.T) {
	s, _ := setupTestProject(t)
	resp := callTool(t, s, "sigil_get_component_schema", map[string]interface{}{
		"type": "button",
	})

	text := extractToolText(t, resp)
	if !strings.Contains(text, "button") {
		t.Error("expected button type")
	}
	if !strings.Contains(text, "props") {
		t.Error("expected props in schema")
	}
}

func TestIntegrationListDataSources(t *testing.T) {
	s, _ := setupTestProject(t)
	resp := callTool(t, s, "sigil_list_datasources", map[string]interface{}{})

	text := extractToolText(t, resp)
	if !strings.Contains(text, "Sprint") {
		t.Error("expected Sprint datasource")
	}
}

func TestIntegrationCreateDataSource(t *testing.T) {
	s, _ := setupTestProject(t)

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

	resp := callTool(t, s, "sigil_create_datasource", map[string]interface{}{
		"alias":  "Task",
		"config": dsConfig,
	})

	text := extractToolText(t, resp)
	if !strings.Contains(text, "task.yaml") {
		t.Errorf("expected file path, got: %s", text)
	}

	// Verify file was created
	path := filepath.Join(s.SigilDir(), "datasources", "task.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected datasource file at %s", path)
	}
}

func TestIntegrationResourcesRead(t *testing.T) {
	s, _ := setupTestProject(t)

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
			resp := sendRequest(t, s, "resources/read", 1, map[string]interface{}{
				"uri": tt.uri,
			})
			if resp.Error != nil {
				t.Fatalf("error reading %s: %v", tt.uri, resp.Error)
			}

			resultData, _ := json.Marshal(resp.Result)
			if !strings.Contains(string(resultData), tt.contains) {
				t.Errorf("expected %q in result for %s, got: %s", tt.contains, tt.uri, string(resultData))
			}
		})
	}
}

func TestIntegrationPromptsDesignPage(t *testing.T) {
	s, _ := setupTestProject(t)

	resp := sendRequest(t, s, "prompts/get", 1, map[string]interface{}{
		"name":      "sigil_design_page",
		"arguments": map[string]string{"description": "A user management page"},
	})

	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}

	resultData, _ := json.Marshal(resp.Result)
	text := string(resultData)
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

	validConfig := `sigil: "1.0"
id: test
title: Test
overlay: page
layout:
  id: root
  type: rows`

	resp := sendRequest(t, s, "prompts/get", 1, map[string]interface{}{
		"name":      "sigil_review_config",
		"arguments": map[string]string{"config": validConfig},
	})

	if resp.Error != nil {
		t.Fatalf("error: %v", resp.Error)
	}

	resultData, _ := json.Marshal(resp.Result)
	if !strings.Contains(string(resultData), "PASS") {
		t.Error("expected PASS for valid config")
	}
}

func TestIntegrationEndToEnd(t *testing.T) {
	s, _ := setupTestProject(t)

	// Simulate an agent workflow:
	// 1. List components to understand what's available
	resp := callTool(t, s, "sigil_list_components", map[string]interface{}{
		"category": "primitives",
	})
	if resp.Error != nil {
		t.Fatalf("list components error: %v", resp.Error)
	}
	text := extractToolText(t, resp)
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

	resp = callTool(t, s, "sigil_create_page", map[string]interface{}{
		"id":     "user-list",
		"config": pageConfig,
	})
	text = extractToolText(t, resp)
	if !strings.Contains(text, `"valid": true`) {
		t.Fatalf("expected valid page, got: %s", text)
	}

	// 3. Validate the created page
	resp = callTool(t, s, "sigil_validate", map[string]interface{}{
		"path": filepath.Join(s.SigilDir(), "pages", "user-list.yaml"),
	})
	text = extractToolText(t, resp)
	if !strings.Contains(text, `"valid": true`) {
		t.Fatalf("expected validation pass, got: %s", text)
	}

	// 4. List pages — should include both dashboard and user-list
	resp = callTool(t, s, "sigil_list_pages", map[string]interface{}{})
	text = extractToolText(t, resp)
	if !strings.Contains(text, "dashboard") {
		t.Error("expected dashboard in pages list")
	}
	if !strings.Contains(text, "user-list") {
		t.Error("expected user-list in pages list")
	}

	// 5. Get the created page
	resp = callTool(t, s, "sigil_get_page", map[string]interface{}{"id": "user-list"})
	text = extractToolText(t, resp)
	if !strings.Contains(text, "User List") {
		t.Error("expected page title in content")
	}
}

// extractToolText extracts the text content from a tool result response.
func extractToolText(t *testing.T, resp *jsonRPCResponse) string {
	t.Helper()

	resultData, _ := json.Marshal(resp.Result)
	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(resultData, &result); err != nil {
		// Try to find text in the raw result
		return string(resultData)
	}
	if len(result.Content) > 0 {
		return result.Content[0].Text
	}
	return string(resultData)
}
