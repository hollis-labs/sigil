---
type: sprint-guide
sprint: 4
title: "MCP Server + Agent Integration"
status: todo
created_at: 2026-03-01
estimated_days: 3
depends_on: [2]
---

# Sprint 4 — MCP Server + Agent Integration

## Objective

Implement the MCP server so AI agents (Claude Code, Mentat) can create, validate,
and generate Sigil configs programmatically. After this sprint, an agent can scaffold
a complete UI page through MCP tool calls.

Note: Sprint 4 depends on Sprint 2 (component schemas) but NOT on Sprint 3 (renderer).
It can run in parallel with Sprint 3.

## Context

**Read `docs/08_mcp-integration.md`** for the MCP tool and resource definitions.

The MCP server runs via stdio transport (`sigil mcp serve`) and exposes tools for
all Sigil operations.

**MCP SDK:** Use a Go MCP library or implement the JSON-RPC protocol directly.
The protocol is straightforward: JSON-RPC 2.0 over stdin/stdout.

## Tasks

### TASK-S4-001 | Priority A | MCP server scaffolding

**Context:** Set up the MCP server process with JSON-RPC 2.0 transport over stdio.

**What to do:**

Create `internal/mcp/server.go`:
- Read JSON-RPC messages from stdin
- Route to handlers based on method name
- Write JSON-RPC responses to stdout
- Handle `initialize`, `tools/list`, `tools/call`, `resources/list`, `resources/read`

Create `internal/cli/mcp.go`:
```go
func NewMCPCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mcp",
        Short: "MCP server commands",
    }
    cmd.AddCommand(&cobra.Command{
        Use:   "serve",
        Short: "Start MCP server (stdio)",
        RunE:  runMCPServe,
    })
    return cmd
}
```

Register in root command.

**Acceptance criteria:**
- [ ] `sigil mcp serve` starts and waits for JSON-RPC input on stdin
- [ ] Responds to `initialize` with server capabilities
- [ ] Responds to `tools/list` with available tools
- [ ] Responds to `resources/list` with available resources
- [ ] Clean shutdown on EOF/SIGTERM

---

### TASK-S4-002 | Priority A | Implement core MCP tools

**Context:** The tools that agents need most: list pages, create pages, validate,
list components.

**What to do:**

Implement these MCP tools (see `docs/08_mcp-integration.md` for schemas):

1. **`sigil_list_pages`** — scan `.sigil/pages/`, return summaries
2. **`sigil_get_page`** — read a specific page YAML, return content
3. **`sigil_create_page`** — write YAML to `.sigil/pages/`, validate first
4. **`sigil_update_page`** — overwrite existing page
5. **`sigil_validate`** — validate a config string, return errors/warnings
6. **`sigil_list_components`** — return component types with schemas
7. **`sigil_get_component_schema`** — return full schema for a type
8. **`sigil_list_datasources`** — return datasource summaries
9. **`sigil_create_datasource`** — write datasource manifest

Each tool handler:
- Parse input parameters
- Call existing internal functions (same code as CLI commands)
- Return structured JSON response

**Acceptance criteria:**
- [ ] All 9 tools respond correctly
- [ ] `sigil_create_page` validates before saving (returns errors if invalid)
- [ ] `sigil_list_components` returns schemas that an agent can use as context
- [ ] Error responses use proper JSON-RPC error format
- [ ] Tools are discoverable via `tools/list`

---

### TASK-S4-003 | Priority A | Implement MCP resources

**Context:** Resources provide read-only access to project state.

**What to do:**

Implement these MCP resources:
1. `sigil://pages` — page listing
2. `sigil://pages/{id}` — page content
3. `sigil://components` — all component schemas
4. `sigil://components/{type}` — specific component schema
5. `sigil://datasources` — all datasource manifests
6. `sigil://themes` — all themes
7. `sigil://project` — project config

**Acceptance criteria:**
- [ ] All resources listed in `resources/list`
- [ ] `resources/read` returns content for each URI
- [ ] URI templating works for parameterized resources

---

### TASK-S4-004 | Priority B | MCP prompts

**Context:** Prompts guide agents through common workflows.

**What to do:**

Implement:
1. **`sigil_design_page`** — returns a system prompt with component schemas and config spec,
   tailored for the requested page description
2. **`sigil_review_config`** — validates and suggests improvements

**Acceptance criteria:**
- [ ] Prompts listed in `prompts/list`
- [ ] `sigil_design_page` returns useful context for UI generation
- [ ] `sigil_review_config` returns validation + improvement suggestions

---

### TASK-S4-005 | Priority B | Integration test with Claude Code config

**Context:** Verify the MCP server works with Claude Code's MCP configuration.

**What to do:**

1. Create a test script that:
   - Starts `sigil mcp serve` as a subprocess
   - Sends `initialize` → verifies response
   - Sends `tools/list` → verifies tools are listed
   - Sends `tools/call` for `sigil_list_components` → verifies response
   - Sends `tools/call` for `sigil_create_page` with a sample config → verifies file created
   - Sends `tools/call` for `sigil_validate` on the created file → verifies valid

2. Document MCP server registration for Claude Code:
   ```json
   {
     "mcpServers": {
       "sigil": {
         "command": "sigil",
         "args": ["mcp", "serve"],
         "cwd": "/path/to/project"
       }
     }
   }
   ```

**Acceptance criteria:**
- [ ] Integration test passes end-to-end
- [ ] MCP registration config documented
- [ ] Agent can create a valid page through MCP tools alone
