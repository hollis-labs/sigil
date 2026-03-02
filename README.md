# Sigil

**Declarative UI configs to framework-specific code.**

Sigil turns YAML page definitions into working UI code for Go/Templ+HTMX or React/shadcn+Tailwind. Define your pages once, generate for any target.

## Features

- **49 component types** — layouts, forms, data tables, modals, buttons, badges, and more
- **Two render targets** — `go-templ` (Go/Templ+HTMX) and `react-shadcn` (React/TypeScript+shadcn/ui)
- **Live dev server** — `sigil serve` with auto-reload on config changes
- **Preview** — instant browser preview with mock data, no backend needed
- **Validation** — deep prop validation against component schemas with "did you mean?" suggestions
- **MCP server** — AI/agent integration via JSON-RPC 2.0 over stdio
- **Theme system** — design tokens as CSS variables with Tailwind integration
- **DataSource model** — declare data shapes, get typed hooks/handlers generated
- **Export/Import** — JSON round-trip for page configs
- **Config diffing** — semantic comparison of page configs
- **Schema migration** — auto-fix missing IDs, set defaults
- **JSON Schema** — export for IDE autocomplete (VS Code YAML extension)

## Quick Start

```bash
# Install
go install github.com/chrispian/sigil/cmd/sigil@latest

# Or build from source
make build

# Initialize a project
sigil init --name my-app

# Create a page
sigil new page dashboard --title "Dashboard"

# Validate configs
sigil validate

# Preview in browser
sigil preview dashboard

# Start live dev server
sigil serve

# Generate Go/Templ code
sigil generate --target go-templ --output internal/ui

# Generate React/shadcn code
sigil generate --target react-shadcn --output src/generated
```

## How It Works

1. **Define pages in YAML** — describe your UI as a component tree with props, actions, and data bindings
2. **Validate** — Sigil checks types, required props, and structure against 49 component schemas
3. **Generate** — choose a render target and get framework-specific source files
4. **Preview** — see your pages instantly with mock data before writing any backend code

### Example Page Config

```yaml
sigil: "1.0"
kind: page
id: user-list
title: User Management
overlay: page
datasources:
  - alias: users
layout:
  id: root
  type: rows
  props:
    gap: 4
    padding: 6
  children:
    - id: header
      type: columns
      props:
        justify: between
        align: center
      children:
        - id: title
          type: heading
          props:
            level: 2
            text: Users
        - id: add-btn
          type: button
          props:
            label: Add User
            variant: primary
            icon: plus
          actions:
            click:
              type: modal
              title: New User
    - id: search
      type: search-bar
      props:
        placeholder: Search users...
        datasource: users
    - id: table
      type: data-table
      props:
        datasource: users
        columns:
          - field: name
            label: Name
          - field: email
            label: Email
            sortable: true
          - field: role
            label: Role
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `sigil init` | Initialize a new project |
| `sigil new page <id>` | Create a page config |
| `sigil validate` | Validate all configs |
| `sigil generate` | Generate framework code |
| `sigil preview <page>` | Preview a page in browser |
| `sigil serve` | Start live dev server |
| `sigil list <type>` | List pages, components, themes |
| `sigil export` | Export pages as JSON |
| `sigil import` | Import pages from JSON |
| `sigil diff <a> <b>` | Compare configs semantically |
| `sigil migrate` | Migrate configs to latest schema |
| `sigil schema export` | Generate JSON Schema for IDE support |
| `sigil doctor` | Check environment health |
| `sigil version` | Show version |
| `sigil mcp serve` | Start MCP server for AI agents |

## Render Targets

### Go/Templ (`go-templ`)

Generates `.templ` files with HTMX interactions, shared component templates, handler stubs for datasources, and theme CSS.

```bash
sigil generate --target go-templ --output internal/ui
```

Output:
```
internal/ui/
  pages/          .templ page components
  components/     shared .templ components
  handlers/       Go handler stubs per datasource
  theme.css       CSS custom properties
  tailwind.sigil.cjs  Tailwind preset
```

### React/shadcn (`react-shadcn`)

Generates `.tsx` files with shadcn/ui imports, SWR data hooks, TypeScript types, and Tailwind config.

```bash
sigil generate --target react-shadcn --output src/generated
```

Output:
```
src/generated/
  pages/          .tsx page components
  components/     DataTable, barrel exports
  hooks/          SWR hooks per datasource
  types/          TypeScript interfaces
  lib/            cn() utility
  globals.css     CSS custom properties
  tailwind.config.ts  Tailwind config with theme tokens
```

## Component Types

49 built-in types across 6 categories:

| Category | Types |
|----------|-------|
| **Layout** | rows, columns, grid, card, sidebar, split |
| **Primitive** | heading, text, button, icon, badge, label, separator, spacer, avatar, progress, alert |
| **Form** | input, textarea, select, checkbox, radio, toggle, slider, date-picker, form |
| **Data** | data-table, list, detail-view, stat-card, chart, tree, timeline |
| **Overlay** | modal, sheet, drawer, popover, tooltip, dropdown, context-menu, command-palette |
| **Navigation** | tabs, breadcrumb, pagination, stepper, nav-menu, search-bar |

## Theme System

Themes define design tokens as CSS custom properties:

```yaml
name: dark
tokens:
  colors:
    background: "9 9 11"
    surface: "24 24 27"
    text: "244 244 245"
    accent: "79 70 229"
    border: "63 63 70"
  radius:
    md: "0.375rem"
    lg: "0.5rem"
```

Tokens become `--sigil-*` CSS variables and are available via Tailwind (`sigil-*` color palette).

## MCP Integration

Sigil includes an MCP server for AI agent workflows:

```bash
sigil mcp serve
```

Provides 9 tools (list/get/create/update pages, validate, component info, datasource management), 7 resource types, and 2 prompts for page design and config review.

## Project Structure

```
.sigil/
  sigil.yaml          Project config
  pages/              Page definitions (*.yaml)
  themes/             Theme files (*.yaml)
  datasources/        DataSource manifests (*.yaml)
  components/         Custom component schemas (*.schema.yaml)
docs/
  01_architecture.md  System architecture
  02_config-spec.md   Config YAML specification
  03_component-model.md  Component type system
  04_cli-reference.md CLI commands reference
  05_renderer-contract.md  Renderer plugin contract
  06_datasource-model.md  DataSource abstraction
  07_theme-system.md  Theme tokens and CSS
  08_mcp-integration.md  MCP server docs
```

## Development

```bash
# Run tests
make test

# Run vet
make vet

# Build
make build

# Install locally
make install
```

## License

MIT
