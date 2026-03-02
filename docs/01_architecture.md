---
type: architecture
version: 1
updated_at: 2026-03-01
---

# Sigil Architecture

## What Sigil Does

Sigil turns **declarative YAML configs** into **framework-specific UI code**. It is a
compiler for user interfaces: you describe *what* the UI should be (pages, components,
layouts, data bindings, themes), and Sigil generates *how* to render it in your target
framework.

## Design Principles

1. **Config is the product.** The YAML config fully describes the UI. All code is
   generated from it. If the config is valid, the output is correct.

2. **System-agnostic core.** The config format, component schemas, and theme tokens
   contain zero framework-specific code. Framework details live only in renderers.

3. **CLI/MCP/agent-first.** The primary interfaces are the CLI and MCP server. AI
   agents (Mentat) and human developers use the same tool. GUI comes last (built
   with Sigil itself).

4. **Code generation over runtime rendering.** Sigil produces source files that you
   compile/build normally. No runtime config loading, no dynamic rendering. The
   generated code is readable, editable, and version-controlled.

5. **Composable, not comprehensive.** Start with 25 core components. Add more as
   needed. Don't build 100+ components upfront.

## System Overview

```
                    ┌──────────────────┐
                    │   User / Agent   │
                    └────────┬─────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              ▼              ▼
        ┌──────────┐  ┌──────────┐  ┌──────────┐
        │ Sigil CLI │  │ MCP Srvr │  │  Mentat  │
        │ (Cobra)   │  │ (stdio)  │  │ (MRO →   │
        │           │  │          │  │  config)  │
        └─────┬─────┘  └────┬─────┘  └────┬─────┘
              │              │              │
              └──────────────┼──────────────┘
                             │
                             ▼
                 ┌───────────────────────┐
                 │     Sigil Core        │
                 │                       │
                 │  ┌─────────────────┐  │
                 │  │  Config Parser  │  │  YAML → Go structs
                 │  └────────┬────────┘  │
                 │           │           │
                 │  ┌────────▼────────┐  │
                 │  │  Schema Valid.  │  │  Validate against component schemas
                 │  └────────┬────────┘  │
                 │           │           │
                 │  ┌────────▼────────┐  │
                 │  │ Component Reg.  │  │  Resolve type → schema + defaults
                 │  └────────┬────────┘  │
                 │           │           │
                 │  ┌────────▼────────┐  │
                 │  │  Theme Resolver │  │  Tokens → CSS variables
                 │  └────────┬────────┘  │
                 │           │           │
                 │  ┌────────▼────────┐  │
                 │  │  DataSource Res │  │  Alias → endpoint contract
                 │  └────────┬────────┘  │
                 │           │           │
                 └───────────┼───────────┘
                             │
                             ▼
                 ┌───────────────────────┐
                 │    Renderer Engine    │
                 │                       │
                 │  ┌─────────────────┐  │
                 │  │   go-templ      │  │  → .templ + handlers + CSS
                 │  │   react-shadcn  │  │  → .tsx + hooks + tailwind
                 │  │   html-static   │  │  → .html + inline CSS
                 │  └─────────────────┘  │
                 │                       │
                 └───────────┬───────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  Generated Code │
                    │  (source files) │
                    └─────────────────┘
```

## Go Module Structure

```
sigil/
├── cmd/sigil/main.go              # Entry point: cobra root command
├── internal/
│   ├── config/
│   │   ├── types.go               # Core types: Page, Component, Action, etc.
│   │   ├── parser.go              # YAML/JSON parsing
│   │   ├── validator.go           # Schema validation
│   │   └── defaults.go            # Default values for optional fields
│   ├── cli/
│   │   ├── root.go                # Root cobra command
│   │   ├── init.go                # sigil init
│   │   ├── new.go                 # sigil new page|component|datasource|theme
│   │   ├── validate.go            # sigil validate
│   │   ├── generate.go            # sigil generate
│   │   ├── list.go                # sigil list pages|components|themes|datasources
│   │   └── preview.go             # sigil preview (open browser)
│   ├── renderer/
│   │   ├── renderer.go            # Renderer interface
│   │   ├── context.go             # Render context (resolved config + theme)
│   │   └── gotempl/
│   │       ├── renderer.go        # Go/Templ renderer implementation
│   │       ├── templates/         # Go template files for code generation
│   │       └── helpers.go         # Templ-specific helpers
│   ├── components/
│   │   ├── registry.go            # Component type registry
│   │   ├── schema.go              # Component schema loading
│   │   └── builtin/               # Built-in component schemas (YAML files)
│   ├── theme/
│   │   ├── types.go               # Token types
│   │   ├── resolver.go            # Token resolution + inheritance
│   │   └── css.go                 # CSS variable generation
│   └── datasource/
│       ├── types.go               # DataSource manifest types
│       ├── parser.go              # Manifest parsing
│       └── stubgen.go             # Handler stub generation
├── schemas/                        # JSON Schema files
│   ├── page.schema.json
│   ├── component.schema.json
│   ├── datasource.schema.json
│   └── theme.schema.json
└── examples/                       # Example configs + expected output
    ├── sprint-dashboard.yaml
    ├── user-management.yaml
    └── login-form.yaml
```

## Data Flow: Config → Generated Code

```
1. Developer writes .sigil/pages/dashboard.yaml (or Mentat generates it)
2. `sigil validate .sigil/pages/dashboard.yaml`
   → Parse YAML → Load component schemas → Validate props → Report errors
3. `sigil generate --target go-templ --output internal/ui/`
   → Parse YAML
   → Resolve theme tokens → CSS variables
   → Resolve data sources → endpoint contracts
   → Walk component tree → generate Templ code for each component
   → Write .templ files + handler stubs + CSS file
4. Developer builds: `go generate ./... && go build`
5. App serves the generated UI
```

## Project File Layout (for Sigil-managed projects)

When a developer runs `sigil init` in their project, Sigil creates:

```
myproject/
├── .sigil/
│   ├── sigil.yaml            # Project-level Sigil config
│   ├── pages/                # Page definitions
│   │   └── dashboard.yaml
│   ├── components/           # Custom component schemas (optional)
│   ├── datasources/          # DataSource manifests
│   │   └── sprint.yaml
│   └── themes/               # Theme definitions
│       └── default.yaml
└── (generated code goes wherever --output points)
```

## Key Design Decisions

### 1. Code generation, not runtime rendering

The ui-builder v3 loads configs from the database at runtime and renders them
dynamically. Sigil generates source files. Trade-offs:

| | Runtime rendering | Code generation |
|---|---|---|
| **Deploy changes** | Database update only | Rebuild + redeploy |
| **Generated code quality** | N/A — no code | Readable, editable, reviewable |
| **Performance** | Config parse on every request | Native compiled code |
| **Debugging** | Inspect JSON at runtime | Read generated source files |
| **AI authoring** | Write JSON to DB | Write YAML to file → generate |
| **Version control** | DB migrations | Git-tracked YAML files |

For the ecosystem (Volon, Nanite, Mentat), code generation is better: configs are
git-tracked artifacts, generated code is reviewable, and the output is native.

### 2. YAML as primary format

YAML is more readable than JSON for humans and more reliable for LLM generation
(fewer syntax errors with brackets/commas). JSON is supported via conversion for
programmatic use.

### 3. Component schemas as LLM context

Each component type has a schema (YAML file) that describes its props, slots, actions,
and constraints. These schemas serve as documentation for humans AND as context for
AI agents generating configs. When Mentat generates a Sigil config, it receives the
component schemas as part of its context packet.

### 4. DataSource as abstract contract

DataSources are described by alias + capabilities + field definitions + endpoint
patterns. The config doesn't know or care about the backend implementation. The
renderer generates handler stubs; the developer fills in the query logic.

### 5. Theme tokens → CSS variables

Themes define design tokens (colors, typography, spacing, radius) that map to CSS
custom properties. This works with any CSS framework: Tailwind, plain CSS, or
framework-specific styling. Renderers consume tokens through CSS variables.
