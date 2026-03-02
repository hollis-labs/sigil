// Package mcp implements the Model Context Protocol server for Sigil.
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// ResourceHandler reads a resource by URI and returns content.
type ResourceHandler func(uri string) (*ResourceContent, error)

// PromptHandler processes a prompt get request and returns messages.
type PromptHandler func(args map[string]string) ([]PromptMessage, error)

// Server handles MCP JSON-RPC 2.0 communication over stdio.
type Server struct {
	tools            map[string]Tool
	resources        []Resource
	resourceHandler  ResourceHandler
	prompts          []Prompt
	promptHandlers   map[string]PromptHandler
	sigilDir         string
	mu               sync.RWMutex
}

// Tool defines an MCP tool.
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
	Handler     ToolHandler     `json:"-"`
}

// ToolHandler processes a tool call and returns content.
type ToolHandler func(params json.RawMessage) (*ToolResult, error)

// ToolResult is the response from a tool call.
type ToolResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

// ContentBlock is a piece of content in a tool result.
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Resource defines an MCP resource.
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ResourceContent is the content of a resource read.
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
}

// Prompt defines an MCP prompt.
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument defines a prompt parameter.
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// PromptMessage is a message in a prompt response.
type PromptMessage struct {
	Role    string       `json:"role"`
	Content ContentBlock `json:"content"`
}

// JSON-RPC types
type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewServer creates a new MCP server.
func NewServer(sigilDir string) *Server {
	return &Server{
		tools:          make(map[string]Tool),
		promptHandlers: make(map[string]PromptHandler),
		sigilDir:       sigilDir,
	}
}

// RegisterTool adds a tool to the server.
func (s *Server) RegisterTool(tool Tool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[tool.Name] = tool
}

// SetResources sets the available resources.
func (s *Server) SetResources(resources []Resource) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resources = resources
}

// SetPrompts sets the available prompts.
func (s *Server) SetPrompts(prompts []Prompt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prompts = prompts
}

// SetResourceHandler sets the handler for resource reads.
func (s *Server) SetResourceHandler(handler ResourceHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resourceHandler = handler
}

// RegisterPromptHandler adds a prompt handler.
func (s *Server) RegisterPromptHandler(name string, handler PromptHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.promptHandlers[name] = handler
}

// SigilDir returns the configured .sigil directory path.
func (s *Server) SigilDir() string {
	return s.sigilDir
}

// Run starts the server, reading from reader and writing to writer.
func (s *Server) Run(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // 10MB max

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req jsonRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.writeError(writer, nil, -32700, "Parse error", err.Error())
			continue
		}

		resp := s.handleRequest(&req)
		if resp != nil {
			data, _ := json.Marshal(resp)
			fmt.Fprintf(writer, "%s\n", data)
		}
	}

	return scanner.Err()
}

// Serve starts the server on stdin/stdout.
func (s *Server) Serve() error {
	return s.Run(os.Stdin, os.Stdout)
}

func (s *Server) handleRequest(req *jsonRPCRequest) *jsonRPCResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "initialized":
		// Notification, no response
		return nil
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	case "resources/list":
		return s.handleResourcesList(req)
	case "resources/read":
		return s.handleResourcesRead(req)
	case "prompts/list":
		return s.handlePromptsList(req)
	case "prompts/get":
		return s.handlePromptsGet(req)
	case "ping":
		return &jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]interface{}{}}
	default:
		return s.errorResponse(req.ID, -32601, "Method not found", req.Method)
	}
}

func (s *Server) handleInitialize(req *jsonRPCRequest) *jsonRPCResponse {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools":     map[string]interface{}{},
			"resources": map[string]interface{}{},
			"prompts":   map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "sigil",
			"version": "0.1.0",
		},
	}
	return &jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: result}
}

func (s *Server) handleToolsList(req *jsonRPCRequest) *jsonRPCResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tools []map[string]interface{}
	for _, t := range s.tools {
		tool := map[string]interface{}{
			"name":        t.Name,
			"description": t.Description,
		}
		if t.InputSchema != nil {
			tool["inputSchema"] = json.RawMessage(t.InputSchema)
		}
		tools = append(tools, tool)
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"tools": tools},
	}
}

func (s *Server) handleToolsCall(req *jsonRPCRequest) *jsonRPCResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.errorResponse(req.ID, -32602, "Invalid params", err.Error())
	}

	s.mu.RLock()
	tool, ok := s.tools[params.Name]
	s.mu.RUnlock()

	if !ok {
		return s.errorResponse(req.ID, -32602, "Unknown tool", params.Name)
	}

	result, err := tool.Handler(params.Arguments)
	if err != nil {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: &ToolResult{
				Content: []ContentBlock{{Type: "text", Text: err.Error()}},
				IsError: true,
			},
		}
	}

	return &jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: result}
}

func (s *Server) handleResourcesList(req *jsonRPCRequest) *jsonRPCResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"resources": s.resources},
	}
}

func (s *Server) handleResourcesRead(req *jsonRPCRequest) *jsonRPCResponse {
	var params struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.errorResponse(req.ID, -32602, "Invalid params", err.Error())
	}

	s.mu.RLock()
	handler := s.resourceHandler
	s.mu.RUnlock()

	if handler == nil {
		return s.errorResponse(req.ID, -32602, "No resource handler", nil)
	}

	content, err := handler(params.URI)
	if err != nil {
		return s.errorResponse(req.ID, -32602, "Resource read error", err.Error())
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"contents": []interface{}{content}},
	}
}

func (s *Server) handlePromptsList(req *jsonRPCRequest) *jsonRPCResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"prompts": s.prompts},
	}
}

func (s *Server) handlePromptsGet(req *jsonRPCRequest) *jsonRPCResponse {
	var params struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.errorResponse(req.ID, -32602, "Invalid params", err.Error())
	}

	s.mu.RLock()
	handler, ok := s.promptHandlers[params.Name]
	s.mu.RUnlock()

	if !ok {
		return s.errorResponse(req.ID, -32602, "Unknown prompt", params.Name)
	}

	messages, err := handler(params.Arguments)
	if err != nil {
		return s.errorResponse(req.ID, -32602, "Prompt error", err.Error())
	}

	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"messages": messages},
	}
}

func (s *Server) writeError(writer io.Writer, id interface{}, code int, message string, data interface{}) {
	resp := s.errorResponse(id, code, message, data)
	jsonData, _ := json.Marshal(resp)
	fmt.Fprintf(writer, "%s\n", jsonData)
}

func (s *Server) errorResponse(id interface{}, code int, message string, data interface{}) *jsonRPCResponse {
	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &jsonRPCError{Code: code, Message: message, Data: data},
	}
}
