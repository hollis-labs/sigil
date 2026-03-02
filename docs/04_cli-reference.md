---
type: specification
version: 2
updated_at: 2026-03-02
---

# Sigil CLI Reference

## Overview

`sigil` is a Go CLI tool built with Cobra. It manages Sigil configs, validates
them, and generates framework-specific code.

## Commands

### `sigil init`

Initialize a Sigil project in the current directory.

```
sigil init [--name <project>]
```

Creates:
```
.sigil/
├── sigil.yaml             # Project config
├── pages/                 # Page definitions (empty)
├── components/            # Custom component schemas (empty)
├── datasources/           # DataSource manifests (empty)
└── themes/
    └── default.yaml       # Default theme with standard tokens
```

**Flags:**
- `--name <string>` — Project name (defaults to directory name)
- `--theme <dark|light>` — Initial theme preset

### `sigil new page <id>`

Create a new page config from a template.

```
sigil new page sprint-dashboard [--title "Sprint Dashboard"] [--overlay page] [--module project.sprints]
```

Creates `.sigil/pages/<id>.yaml` with a starter template.

**Flags:**
- `--title <string>` — Page title (inferred from ID if omitted)
- `--overlay <page|modal|sheet|drawer|fullscreen>` — Display mode (default: page)
- `--module <string>` — Module grouping
- `--layout <rows|columns|grid>` — Root layout type (default: rows)
- `--datasource <alias>` — Add a datasource declaration
- `--interactive` — Interactive mode: prompt for each option

### `sigil new datasource <alias>`

Create a new datasource manifest.

```
sigil new datasource Sprint [--capabilities search,filter,sort,paginate,create,update,delete]
```

Creates `.sigil/datasources/<alias>.yaml`.

**Flags:**
- `--capabilities <list>` — Comma-separated capabilities
- `--fields <json>` — Field definitions (JSON array)
- `--interactive` — Interactive mode

### `sigil new theme <id>`

Create a new theme.

```
sigil new theme dark-blue [--extends default]
```

Creates `.sigil/themes/<id>.yaml`.

**Flags:**
- `--extends <theme>` — Base theme to inherit from
- `--preset <dark|light>` — Start from a preset

### `sigil new component <type>`

Create a custom component schema.

```
sigil new component sprint-card [--category data] [--extends card]
```

Creates `.sigil/components/<type>.schema.yaml`.

### `sigil validate [<path>]`

Validate Sigil configs against component schemas.

```
sigil validate                          # Validate all configs in .sigil/
sigil validate .sigil/pages/dashboard.yaml  # Validate a specific file
sigil validate --strict                  # Fail on warnings too
```

**Checks performed:**
- YAML syntax valid
- Required fields present
- Component types exist in registry
- Props match component schema
- DataSource aliases declared
- Action types valid
- Template variables reference valid fields
- Keyboard shortcuts have valid key expressions
- No duplicate component IDs

**Output:**
```
✓ .sigil/pages/dashboard.yaml — valid
✗ .sigil/pages/users.yaml — 2 errors, 1 warning
  ERROR: line 34: unknown component type "data-grid" (did you mean "data-table"?)
  ERROR: line 45: prop "columns" is required for "data-table"
  WARN:  line 12: datasource "User" not declared in datasources section
```

**Exit codes:** 0 = valid, 1 = errors, 2 = invalid args

### `sigil generate`

Generate framework-specific code from Sigil configs.

```
sigil generate --target go-templ --output internal/ui/
sigil generate --target react-shadcn --output src/components/generated/
sigil generate --target html-static --output dist/
```

**Flags:**
- `--target <renderer>` — Target framework (required)
- `--output <dir>` — Output directory (required)
- `--pages <glob>` — Pages to generate (default: all)
- `--theme <id>` — Theme to use (default: project default)
- `--clean` — Remove existing files in output dir before generating
- `--dry-run` — Show what would be generated without writing

**Available targets:**
- `go-templ` — Go Templ components + HTMX + Tailwind + handler stubs
- `react-shadcn` — React/TypeScript + shadcn/ui + SWR hooks + Tailwind config

**Generated structure (go-templ):**
```
internal/ui/
├── pages/
│   └── sprint_dashboard.templ
├── components/
│   ├── data_table.templ
│   ├── button.templ
│   ├── modal.templ
│   └── ...
├── handlers/
│   └── sprint_handler.go
├── layouts/
│   └── page_layout.templ
└── theme.css
```

**Generated structure (react-shadcn):**
```
src/generated/
├── pages/
│   └── sprint-dashboard.tsx
├── components/
│   ├── data-table.tsx
│   └── index.ts
├── hooks/
│   └── use-sprint.ts
├── types/
│   └── sprint.ts
├── lib/
│   └── utils.ts
├── globals.css
└── tailwind.config.ts
```

### `sigil list <type>`

List registered items.

```
sigil list pages                # List all page configs
sigil list components           # List available component types
sigil list datasources          # List datasource manifests
sigil list themes               # List themes
```

**Output:**
```
PAGES (3 found)
  sprint-dashboard    Sprint Dashboard           page      project.sprints
  user-management     User Management            page      admin.users
  login               Login                      modal     auth

COMPONENTS (45 registered)
  Primitives:  heading, text, button, icon-button, input, textarea, ...
  Layouts:     rows, columns, grid, card, tabs, split, sidebar, ...
  Data:        data-table, list, detail-view, stat-card, chart, ...
  ...
```

### `sigil preview <page>`

Generate a standalone HTML preview. Opens in browser.

```
sigil preview dashboard
sigil preview dashboard --output preview/
sigil preview dashboard --no-open
```

Generates static HTML with Tailwind CDN, mock data from datasource field definitions,
and inline theme CSS. No backend needed.

**Flags:**
- `--sigil-dir <string>` — Path to .sigil directory (default: ".sigil")
- `--theme <string>` — Theme to use
- `--output <dir>` — Output directory (default: temp dir)
- `--no-open` — Don't open in browser

### `sigil serve`

Start a live development server with auto-reload.

```
sigil serve
sigil serve --port 8080
sigil serve --sigil-dir path/to/.sigil
```

Starts an HTTP server that renders pages on-the-fly. Watches `.sigil/` for
file changes and reloads the browser via SSE (Server-Sent Events).

**Features:**
- Index page listing all pages with validation status
- On-the-fly page rendering with theme CSS
- Auto-reload on file changes (pages, themes, datasources)
- Navigation bar injected into page previews

**Flags:**
- `--port <int>` — HTTP server port (default: 3210)
- `--sigil-dir <string>` — Path to .sigil directory (default: ".sigil")

### `sigil export`

Export page configs as JSON.

```
sigil export --format json
sigil export --format json --output configs/ --pages dashboard,users
```

**Flags:**
- `--format <string>` — Export format (default: "json")
- `--output <dir>` — Output directory
- `--pages <list>` — Specific page IDs to export (default: all)

### `sigil import`

Import page configs from JSON.

```
sigil import --from json --input configs/dashboard.json
sigil import --from json --input configs/ --dry-run
```

**Flags:**
- `--from <string>` — Import format (default: "json")
- `--input <path>` — Input file or directory
- `--dry-run` — Preview changes without writing

### `sigil diff <file-a> <file-b>`

Compare two page configs semantically.

```
sigil diff old.yaml new.yaml
sigil diff .sigil/pages/v1.yaml .sigil/pages/v2.yaml
```

Shows added/removed/modified components, prop changes, and datasource
differences at the config level (not a text diff).

### `sigil migrate`

Migrate page configs to the latest schema version.

```
sigil migrate
sigil migrate --dry-run
sigil migrate --pages dashboard
```

**Current migrations:**
- Auto-generate missing component IDs
- Set default `sigil` version, `kind`, and `overlay`

**Flags:**
- `--dry-run` — Preview changes without writing
- `--sigil-dir <string>` — Path to .sigil directory
- `--pages <list>` — Specific page IDs to migrate

### `sigil schema export`

Generate JSON Schema for page configs.

```
sigil schema export
sigil schema export --output my-schemas/
```

Generates `page.schema.json` from the component registry. Use with the
VS Code YAML extension for autocomplete and inline validation.

**Flags:**
- `--output <dir>` — Output directory (default: "schemas")

### `sigil doctor`

Check environment and project health.

```
sigil doctor
sigil doctor --sigil-dir path/to/.sigil
```

**Checks:**
- Go version
- Platform (OS/arch)
- `templ` binary (for go-templ renderer)
- Node.js (for react-shadcn renderer)
- `.sigil/` directory structure
- `sigil.yaml` config
- Page count

### `sigil version`

Show the Sigil version.

```
sigil version
```

Version is embedded at build time via `-ldflags`.

### `sigil mcp serve`

Start the MCP server for AI/agent integration.

```
sigil mcp serve
```

Provides JSON-RPC 2.0 over stdio with 9 tools, 7 resource types, and
2 prompts. See `docs/08_mcp-integration.md` for details.

## Global Flags

These flags are available on all commands:

- `--no-color` — Disable colored output (also respects `NO_COLOR` env var)
- `--verbose` — Enable verbose/debug output

## Project Config (.sigil/sigil.yaml)

```yaml
version: "1.0"
name: my-project

defaults:
  theme: default
  renderer: go-templ
  output: internal/ui/

components:
  builtin: true                  # Load built-in component schemas
  custom_dir: .sigil/components  # Custom component schemas

datasources:
  dir: .sigil/datasources

themes:
  dir: .sigil/themes

generation:
  clean: false                   # Don't auto-clean output dir
  go_package: "myproject/internal/ui"  # Go package path for generated code
  go_module: "myproject"         # Go module name
```
