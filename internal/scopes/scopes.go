// Package scopes names the scopes Sigil's HTTP routes and MCP tools require. Sigil runs
// with auth off on loopback, so none is enforced today; they are declared so the tables
// are correct the day the mode changes.
package scopes

const (
	// Read covers every GET route and every tool that does not write.
	Read = "sigil.read"
	// Write covers every route and tool that changes the project's files.
	Write = "sigil.write"
	// MCP is the scope the Substrate mcp module requires on its /mcp endpoint.
	MCP = "sigil.mcp"
)
