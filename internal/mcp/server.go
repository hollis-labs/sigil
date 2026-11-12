// Package mcp implements Sigil's Model Context Protocol server on top of
// the portfolio's go-mcp library (github.com/hollis-labs/libs/mcp/plugin-mcp/go-mcp), targeting
// the 2026-07-28 MCP specification. The server instance, its transports and the /mcp scope
// belong to the Substrate mcp module; this package registers Sigil's tools, resources and
// prompts on it.
package mcp

import (
	"context"

	gomcpserver "github.com/hollis-labs/libs/mcp/plugin-mcp/go-mcp/server"
)

// Tool is go-mcp's tool registration type, re-exported so callers of this
// package don't need to import go-mcp/server directly.
type Tool = gomcpserver.Tool

// Server wraps a go-mcp server with the Sigil project directory it operates
// against. Resources and prompts are registered on the embedded server's SDK server;
// see RegisterAllTools, RegisterAllResources, and RegisterAllPrompts.
type Server struct {
	*gomcpserver.Server
	sigilDir  string
	authorize func(ctx context.Context, tool string) error
}

// NewServer wraps base (the Substrate mcp module's server, or a standalone go-mcp server
// in tests) for the .sigil project directory sigilDir.
func NewServer(base *gomcpserver.Server, sigilDir string) *Server {
	return &Server{Server: base, sigilDir: sigilDir}
}

// SigilDir returns the configured .sigil directory path.
func (s *Server) SigilDir() string {
	return s.sigilDir
}

// UseAuthorizer sets the per-tool scope check run before every tool call (the mcp module's
// Authorize). Without one every call is allowed.
func (s *Server) UseAuthorizer(a func(ctx context.Context, tool string) error) {
	s.authorize = a
}
