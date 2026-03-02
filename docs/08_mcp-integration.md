---
type: specification
version: 1
updated_at: 2026-03-01
---

# Sigil MCP Integration

## Overview

Sigil exposes an **MCP (Model Context Protocol) server** so AI agents and tools like
Claude Code can interact with Sigil configs programmatically. The MCP server provides
tools for creating, validating, listing, and generating Sigil configs.

## MCP Server

The MCP server runs as a subprocess via stdio transport:

```
sigil mcp serve
```

This is registered as an MCP server in Claude Code's config:

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

## Available Tools

### `sigil_list_pages`

List all page configs in the project.

**Input:** `{}`
**Output:** Array of `{ id, title, overlay, module, path }`

### `sigil_get_page`

Get the full config for a page.

**Input:** `{ id: string }`
**Output:** Full page YAML content

### `sigil_create_page`

Create a new page config.

**Input:** `{ id: string, config: string (YAML) }`
**Output:** `{ path: string, valid: boolean, errors: string[] }`

### `sigil_update_page`

Update an existing page config.

**Input:** `{ id: string, config: string (YAML) }`
**Output:** `{ path: string, valid: boolean, errors: string[] }`

### `sigil_validate`

Validate a config string or file.

**Input:** `{ config: string (YAML) }` or `{ path: string }`
**Output:** `{ valid: boolean, errors: ValidationError[], warnings: ValidationWarning[] }`

### `sigil_list_components`

List available component types with their props.

**Input:** `{ category?: string }` (optional: primitives, layouts, data, forms, etc.)
**Output:** Array of `{ type, category, description, props: PropSchema[] }`

### `sigil_get_component_schema`

Get the full schema for a component type.

**Input:** `{ type: string }`
**Output:** Full component schema YAML

### `sigil_list_datasources`

List registered datasource manifests.

**Input:** `{}`
**Output:** Array of `{ alias, description, capabilities, fieldCount }`

### `sigil_get_datasource`

Get a datasource manifest.

**Input:** `{ alias: string }`
**Output:** Full datasource YAML

### `sigil_create_datasource`

Create a new datasource manifest.

**Input:** `{ alias: string, config: string (YAML) }`
**Output:** `{ path: string }`

### `sigil_list_themes`

List available themes.

**Input:** `{}`
**Output:** Array of `{ name, description, extends }`

### `sigil_generate`

Generate code from configs.

**Input:** `{ target: string, pages?: string[], output?: string, dryRun?: boolean }`
**Output:** `{ files: { path, size }[], errors: string[] }`

### `sigil_preview`

Generate a preview URL for a page.

**Input:** `{ id: string }`
**Output:** `{ url: string, html: string }`

## Available Resources

MCP resources provide read-only access to Sigil project state:

### `sigil://pages`
List of all page configs as a summary.

### `sigil://pages/{id}`
Full content of a specific page config.

### `sigil://components`
All available component types and their schemas.

### `sigil://components/{type}`
Schema for a specific component type.

### `sigil://datasources`
All datasource manifests.

### `sigil://themes`
All theme definitions.

### `sigil://project`
Project config (`.sigil/sigil.yaml`).

## AI Agent Usage

When an AI agent (like Mentat) generates UI configs, it should:

1. **Read component schemas** via `sigil_list_components` to know what's available
2. **Read datasource manifests** to know what data is available
3. **Generate a valid YAML config** following the spec in `docs/02_config-spec.md`
4. **Validate** via `sigil_validate` before saving
5. **Save** via `sigil_create_page`
6. **Generate code** via `sigil_generate` (with `dryRun: true` first)

## Prompts

The MCP server provides prompts for common workflows:

### `sigil_design_page`

Guide for designing a page config.

**Arguments:** `{ description: string }`
**Returns:** System prompt with component schemas, datasource info, and config spec.

### `sigil_review_config`

Review a Sigil config for issues.

**Arguments:** `{ config: string (YAML) }`
**Returns:** Validation results + improvement suggestions.
