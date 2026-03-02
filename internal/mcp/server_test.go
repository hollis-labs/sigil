package mcp

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sendRequest(t *testing.T, s *Server, method string, id interface{}, params interface{}) *jsonRPCResponse {
	t.Helper()

	var paramsJSON json.RawMessage
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("marshaling params: %v", err)
		}
		paramsJSON = data
	}

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  paramsJSON,
	}

	reqData, _ := json.Marshal(req)
	reader := bytes.NewReader(append(reqData, '\n'))
	var writer bytes.Buffer

	s.Run(reader, &writer)

	var resp jsonRPCResponse
	if err := json.Unmarshal(writer.Bytes(), &resp); err != nil {
		t.Fatalf("parsing response: %v (raw: %s)", err, writer.String())
	}
	return &resp
}

func TestInitialize(t *testing.T) {
	s := NewServer(".")
	resp := sendRequest(t, s, "initialize", 1, map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]interface{}{"name": "test", "version": "1.0"},
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}
	if result["protocolVersion"] != "2024-11-05" {
		t.Errorf("expected protocol version 2024-11-05, got %v", result["protocolVersion"])
	}

	serverInfo, ok := result["serverInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("expected serverInfo map")
	}
	if serverInfo["name"] != "sigil" {
		t.Errorf("expected server name 'sigil', got %v", serverInfo["name"])
	}
}

func TestToolsList(t *testing.T) {
	s := NewServer(".")
	s.RegisterTool(Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: json.RawMessage(`{"type":"object"}`),
		Handler:     func(params json.RawMessage) (*ToolResult, error) { return textResult("ok"), nil },
	})

	resp := sendRequest(t, s, "tools/list", 2, nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}
	tools, ok := result["tools"].([]interface{})
	if !ok {
		t.Fatal("expected tools array")
	}
	if len(tools) != 1 {
		t.Errorf("expected 1 tool, got %d", len(tools))
	}
}

func TestToolsCall(t *testing.T) {
	s := NewServer(".")
	s.RegisterTool(Tool{
		Name:        "echo_tool",
		Description: "Echo back the input",
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			return textResult(string(params)), nil
		},
	})

	resp := sendRequest(t, s, "tools/call", 3, map[string]interface{}{
		"name":      "echo_tool",
		"arguments": map[string]interface{}{"msg": "hello"},
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	// The result should contain our tool result
	resultData, _ := json.Marshal(resp.Result)
	if !strings.Contains(string(resultData), "hello") {
		t.Errorf("expected 'hello' in result, got %s", string(resultData))
	}
}

func TestToolsCallUnknown(t *testing.T) {
	s := NewServer(".")
	resp := sendRequest(t, s, "tools/call", 4, map[string]interface{}{
		"name":      "nonexistent",
		"arguments": map[string]interface{}{},
	})

	if resp.Error == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestResourcesList(t *testing.T) {
	s := NewServer(".")
	s.SetResources([]Resource{
		{URI: "sigil://pages", Name: "Pages", MimeType: "application/json"},
		{URI: "sigil://project", Name: "Project", MimeType: "text/yaml"},
	})

	resp := sendRequest(t, s, "resources/list", 5, nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}
	resources, ok := result["resources"].([]interface{})
	if !ok {
		t.Fatal("expected resources array")
	}
	if len(resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(resources))
	}
}

func TestResourcesRead(t *testing.T) {
	s := NewServer(".")
	s.SetResourceHandler(func(uri string) (*ResourceContent, error) {
		return &ResourceContent{
			URI:      uri,
			MimeType: "text/plain",
			Text:     "test content for " + uri,
		}, nil
	})

	resp := sendRequest(t, s, "resources/read", 6, map[string]interface{}{
		"uri": "sigil://test",
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	resultData, _ := json.Marshal(resp.Result)
	if !strings.Contains(string(resultData), "test content") {
		t.Errorf("expected 'test content' in result, got %s", string(resultData))
	}
}

func TestPromptsList(t *testing.T) {
	s := NewServer(".")
	s.SetPrompts([]Prompt{
		{Name: "test_prompt", Description: "A test prompt"},
	})

	resp := sendRequest(t, s, "prompts/list", 7, nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("expected map result")
	}
	prompts, ok := result["prompts"].([]interface{})
	if !ok {
		t.Fatal("expected prompts array")
	}
	if len(prompts) != 1 {
		t.Errorf("expected 1 prompt, got %d", len(prompts))
	}
}

func TestPromptsGet(t *testing.T) {
	s := NewServer(".")
	s.RegisterPromptHandler("test_prompt", func(args map[string]string) ([]PromptMessage, error) {
		return []PromptMessage{
			{Role: "user", Content: ContentBlock{Type: "text", Text: "Hello " + args["name"]}},
		}, nil
	})

	resp := sendRequest(t, s, "prompts/get", 8, map[string]interface{}{
		"name":      "test_prompt",
		"arguments": map[string]string{"name": "World"},
	})

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	resultData, _ := json.Marshal(resp.Result)
	if !strings.Contains(string(resultData), "Hello World") {
		t.Errorf("expected 'Hello World' in result, got %s", string(resultData))
	}
}

func TestPing(t *testing.T) {
	s := NewServer(".")
	resp := sendRequest(t, s, "ping", 9, nil)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
}

func TestUnknownMethod(t *testing.T) {
	s := NewServer(".")
	resp := sendRequest(t, s, "nonexistent/method", 10, nil)
	if resp.Error == nil {
		t.Fatal("expected error for unknown method")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("expected code -32601, got %d", resp.Error.Code)
	}
}

func TestMultipleRequests(t *testing.T) {
	s := NewServer(".")
	s.RegisterTool(Tool{
		Name: "counter",
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			return textResult("counted"), nil
		},
	})

	// Send multiple requests in one stream
	var input bytes.Buffer
	for i := 0; i < 3; i++ {
		req := jsonRPCRequest{
			JSONRPC: "2.0",
			ID:      i + 1,
			Method:  "tools/call",
			Params:  json.RawMessage(`{"name":"counter","arguments":{}}`),
		}
		data, _ := json.Marshal(req)
		input.Write(data)
		input.WriteByte('\n')
	}

	var output bytes.Buffer
	s.Run(&input, &output)

	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 responses, got %d", len(lines))
	}
}
