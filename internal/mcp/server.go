// Package mcp implements Sigil's Model Context Protocol server on top of
// the portfolio's go-mcp library (github.com/hollis-labs/go-mcp), targeting
// the 2026-07-28 MCP specification.
package mcp

import (
	gomcpserver "github.com/hollis-labs/go-mcp/server"
)

// Tool is go-mcp's tool registration type, re-exported so callers of this
// package don't need to import go-mcp/server directly.
type Tool = gomcpserver.Tool

// Server wraps a go-mcp server with the Sigil project directory it operates
// against. Resources, prompts, and stdio serving are inherited from the
// embedded go-mcp Server; see RegisterAllTools, RegisterAllResources, and
// RegisterAllPrompts.
type Server struct {
	*gomcpserver.Server
	sigilDir string
}

// NewServer creates a new Sigil MCP server rooted at sigilDir, a .sigil
// project directory.
func NewServer(sigilDir string) *Server {
	return &Server{
		Server:   gomcpserver.NewServer("sigil", "0.1.0"),
		sigilDir: sigilDir,
	}
}

// SigilDir returns the configured .sigil directory path.
func (s *Server) SigilDir() string {
	return s.sigilDir
}
