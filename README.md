# Sigil

Sigil is a build-time UI compiler. Pages, themes, and datasources are
authored once as YAML, validated against 57 built-in component schemas, and
compiled into framework-native source for Go/Templ+HTMX or React/shadcn+
Tailwind. It has no runtime component: generated code never depends on Sigil
once it's written.

> **Pre-release.** Sigil is unreleased, not deployed, and has no outside
> consumers. It's being built in the open: the code, the docs, and this
> README describe what exists today, not a pitch for what's planned.
> Interfaces and behavior change without notice, and there are no
> compatibility guarantees yet.

## What it is today

- **57 built-in component types** across layout, primitive, input, form,
  data, composite, and navigation categories.
- **Two render targets** — `go-templ` (Go/Templ+HTMX) and `react-shadcn`
  (React/TypeScript+shadcn/ui), including an opt-in mode that generates
  against the shared `@hollis-labs/sysop-ui` design-system kit.
- **Live dev server and preview** — `sigil serve` with auto-reload, and
  instant browser preview against mock data with no backend needed.
- **Validation** — deep prop validation against component schemas with "did
  you mean?" suggestions.
- **An MCP server** — `sigil mcp serve` exposes 9 tools, 7 resource types,
  and 2 prompts over stdio, so an agent can list, create, update, and
  validate page configs directly.
- **DataSource model, export/import, config diffing, schema migration, and
  JSON Schema export** for IDE autocomplete.

## Where it sits in the stack

```
   YAML: pages, themes, datasources        (authored by hand or by an agent)
              │
        ┌──────────┐
        │  Sigil   │   validate → resolve → compile → generate
        └──────────┘
              │
   generated source: Go/Templ+HTMX, React/shadcn+Tailwind —
   optionally against the shared @hollis-labs/sysop-ui component kit
```

Sigil is a compiler, not a service — it doesn't run alongside the app it
generates for, and the other Hollis Labs tools don't call into it at runtime.
Its two seams are build-time (CLI) and agent-time (its MCP server).

## Examples

**Daily use.** Define a page once in YAML, then run
`sigil generate --target go-templ` for the Go backend and
`--target react-shadcn` for a throwaway prototype, without hand-porting the
UI between them.

**Agent-driven UI.** A coding agent (running in Nanite, Claude Code, or
elsewhere) connects to `sigil mcp serve`, creates and validates a page config
directly through MCP tool calls, then triggers generation — building a UI
from a description without a human hand-authoring YAML first.

**Composition.** With `--ui-kit sysop`, generated React output imports
shared primitives, data tables, and shell components from
`@hollis-labs/sysop-ui` instead of per-app copies, so a page generated for
one internal tool matches the look of the others.

**In this repo.** `examples/` holds small page configs plus a task-board demo
with its generated Go/Templ and React output checked in. Three demo apps are
generated from the `.sigil/` project: `demo/` (the "Forge" infra dashboard,
`make demo`), `stack-explorer-demo/` (`make se-demo`), and `demo-clockwork/`
(`make cw-demo`). `sysop/` is a Sysop UI admin app over a Sigil project.

## Roadmap

- **Sysop UI kit coverage.** `--ui-kit sysop` currently re-exports 13 shadcn
  primitives from the shared kit; the rest still fall back to per-app
  components. Closing that gap is active work; the Sysop admin screens
  aren't at parity yet.
- **Portfolio compiler direction.** `docs/ui-compiler-direction.md` sets out
  a direction for Sigil resolving shared design-system, datasource, and
  route contracts across Hollis Labs apps, rather than per-project component
  schemas. This is a directional draft, not a scheduled change.

## License

MIT — see [LICENSE](LICENSE).

## Install

Requires Go 1.26.8+ (toolchain 1.26.9). There are no published binaries yet. Install with Go:

```bash
go install github.com/hollis-labs/sigil/cmd/sigil@latest
```

or build from source:

```bash
git clone https://github.com/hollis-labs/sigil.git
cd sigil
make build     # bin/sigil
make install   # or install sigil to $GOBIN
```

## Quick Start

```bash

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
2. **Validate** — Sigil checks types, required props, and structure against 57 component schemas
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
| `sigil new theme <id>` | Create a theme |
| `sigil new datasource <alias>` | Create a DataSource manifest |
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

57 built-in types across 7 categories:

| Category | Types |
|----------|-------|
| **Layouts** | rows, columns, grid, split, spacer, card, sidebar, scroll-area, tabs, tab, accordion, accordion-item |
| **Primitives** | heading, text, button, icon, icon-button, badge, label, separator, avatar, progress, alert, input, textarea, select, combobox, checkbox, switch |
| **Inputs** | date-picker, radio-group, slider, toggle-group |
| **Forms** | form, field-group |
| **Data** | data-table, list, detail-view, stat-card, chart, timeline, search-bar |
| **Composites** | modal, dialog, confirm-dialog, sheet, popover, tooltip, toast, dropdown-menu, context-menu, collapsible, command-group |
| **Navigation** | nav-menu, breadcrumb, pagination, command-palette |

Custom component schemas can be added per project under `.sigil/components/`.

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

Provides 9 tools (list/get/create/update pages, validate, component info, datasource management), 7 resource types, and 2 prompts for page design and config review. See `docs/08_mcp-integration.md`.

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
  09_sysop-ui-kit.md  --ui-kit sysop integration
  ui-compiler-direction.md  Directional architecture draft
  adr/                Architecture decision records
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
