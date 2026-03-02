---
type: specification
version: 1
updated_at: 2026-03-01
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
- `react-shadcn` — React/TypeScript + shadcn/ui (future)
- `html-static` — Standalone HTML + Tailwind CDN (future)

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

Generate a temporary preview of a page. Opens in browser.

```
sigil preview sprint-dashboard
sigil preview sprint-dashboard --port 3456
```

Generates static HTML with mock data and opens in browser. Useful for
reviewing layout and styling without a backend.

**Flags:**
- `--port <int>` — Preview server port (default: 3210)
- `--mock` — Use auto-generated mock data (default: true)

### `sigil export`

Export configs for external use.

```
sigil export --format nanite            # Export as Nanite items
sigil export --format json --output configs/  # Export as JSON files
```

### `sigil import`

Import configs from external sources.

```
sigil import --from nanite --tag "sigil/page/*"
sigil import --from json --input configs/dashboard.json
```

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
