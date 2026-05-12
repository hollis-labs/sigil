# Backend Context — Sigil

> Project-specific backend conventions. Loaded by the backend agent role when working in this project.
> Lives at `sigil/.agentrc/agents/backend.md`.

## Stack

- **Go version:** 1.25.3
- **Module path:** `github.com/chrispian/sigil`
- **CLI framework:** Cobra (`github.com/spf13/cobra v1.10.2`)
- **Config format:** YAML (`gopkg.in/yaml.v3`)
- **Database:** None (file-based: `.sigil/` project directories)
- **Notable dependencies:** Minimal — only Cobra, pflag, mousetrap, yaml.v3. No HTTP routers, no database drivers, no auth libraries.

## Project Structure

```
cmd/
└── sigil/
    └── main.go                 # Bootstrap only — calls cli.NewRootCmd().Execute()
internal/
├── cli/
│   ├── root.go                 # Root Cobra command + global flags (--no-color, --verbose)
│   ├── validate.go             # sigil validate [path] --strict
│   ├── generate.go             # sigil generate --target --output --pages --clean --dry-run
│   ├── init.go                 # sigil init --name --theme --force
│   ├── new.go                  # sigil new {page,datasource,theme}
│   ├── serve.go                # sigil serve --port
│   ├── preview.go              # sigil preview <page> --theme --port
│   ├── export.go               # sigil export --output --pages --format
│   ├── import.go               # sigil import --input --from --dry-run
│   ├── diff.go                 # sigil diff <a> <b>
│   ├── migrate.go              # sigil migrate --pages --dry-run
│   ├── schema.go               # sigil schema export --output
│   ├── doctor.go               # sigil doctor (environment health check)
│   ├── mcp.go                  # sigil mcp serve (MCP server over stdio)
│   ├── list.go                 # sigil list {pages,components,themes}
│   ├── suggest.go              # sigil suggest (AI-assisted page suggestions)
│   ├── version.go              # sigil version
│   └── color.go                # ANSI color helpers, NO_COLOR support
├── config/
│   ├── types.go                # Page, Component, Action, DataSourceRef, ThemeRef, ProjectConfig structs
│   ├── parser.go               # YAML parsing, ParseFile(), ParseString()
│   ├── validator.go            # Schema validation (props, types, required fields)
│   ├── diff.go                 # Semantic config diffing
│   ├── migrate.go              # Schema migration and auto-fix
│   ├── defaults.go             # Default value application for optional fields
│   ├── writer.go               # YAML serialization back to files
│   └── templates.go            # Factory functions for new page/datasource/theme configs
├── components/
│   ├── registry.go             # Schema registry: Register, Get, Has, Types, GetSchema
│   ├── loader.go               # Load built-in YAML schemas from builtin/ + custom schemas
│   └── builtin/                # 49 component schemas (.yaml files)
│       ├── button.yaml         # category: primitives
│       ├── data-table.yaml     # category: data
│       ├── modal.yaml          # category: overlay
│       └── ...                 # 46 more across 6 categories
├── renderer/
│   ├── renderer.go             # Renderer interface + OutputFile, RenderContext, DataSourceManifest types
│   ├── engine.go               # Generation orchestration, renderer registry
│   ├── gotempl/                # Go/Templ + HTMX renderer
│   │   ├── renderer.go         # GoTemplRenderer implementation
│   │   ├── pages.go            # Page → .templ file generation
│   │   ├── shared.go           # Shared component .templ generation
│   │   ├── theme.go            # Theme → CSS generation
│   │   └── stubs.go            # Handler stub generation
│   ├── reactshadcn/            # React/shadcn + Tailwind renderer
│   │   ├── renderer.go         # ReactShadcnRenderer implementation
│   │   ├── pages.go            # Page → .tsx file generation
│   │   ├── shared.go           # Shared component generation
│   │   └── theme.go            # Theme → Tailwind config + CSS generation
│   └── preview/
│       └── preview.go          # HTML preview with mock data (standalone, no build step)
├── mcp/
│   ├── server.go               # JSON-RPC 2.0 MCP server over stdio
│   ├── tools.go                # 9 MCP tools (list/get/create/update pages, validate, components, datasources)
│   ├── resources.go            # MCP resources (pages, components, themes, datasources)
│   └── prompts.go              # MCP prompts (page design guidance, config review)
├── server/
│   ├── server.go               # Live dev server (HTTP) with routes for pages, theme, SSE
│   ├── watcher.go              # File watcher for auto-reload
│   └── index.go                # Index page generation
├── theme/
│   └── doc.go                  # Package doc (theme logic lives in renderer)
└── datasource/
    └── doc.go                  # Package doc (datasource logic lives in renderer)
schemas/                        # Empty — generated via `sigil schema export`
artifacts/plan/sprints/         # Sprint planning docs (sprint-0 through sprint-8)
examples/
├── sprint-dashboard.yaml       # Full dashboard example with data tables, search, modals
├── login-form.yaml             # Simple form example
├── invalid-example.yaml        # Invalid config for testing
└── task-board-demo/            # Complete example with generated output for both renderers
    ├── page.yaml
    ├── datasource.yaml
    ├── generated-react/        # Generated React/shadcn output
    └── generated-gotempl/      # Generated Go/Templ output
.sigil/                         # Sigil dogfoods itself
├── sigil.yaml                  # Project config
├── pages/*.yaml                # Sigil's own UI pages (editor, components, datasources, themes, preview)
├── themes/default.yaml         # Default dark theme
└── preview/*.html              # Generated preview HTML
```

## Package Inventory

| Package | Location | Responsibility |
|---------|----------|----------------|
| cli | `internal/cli/` | 18 Cobra subcommands, global flags, color output, version injection via ldflags |
| config | `internal/config/` | YAML parsing, validation, diffing, migration, defaults, serialization. Central types: `Page`, `Component`, `Action`, `DataSourceRef`, `ProjectConfig` |
| components | `internal/components/` | Component schema registry. 49 built-in schemas loaded from `builtin/*.yaml`. Extensible with custom schemas from `.sigil/components/` |
| renderer | `internal/renderer/` | `Renderer` interface + orchestration engine. Two implementations: `gotempl`, `reactshadcn`. Plus `preview` for standalone HTML |
| mcp | `internal/mcp/` | JSON-RPC 2.0 MCP server: 9 tools, 7 resource types, 2 prompts. Standalone — no external dependencies |
| server | `internal/server/` | HTTP dev server with file watching and SSE for live reload |

## Core Types

### Page Config (`config/types.go`)
- `Page` — top-level: sigil version, kind, id, title, overlay, module, theme, datasources, layout, shortcuts, meta
- `Component` — tree node: id, type, props (map), actions (map), children, condition
- `Action` — event handler: type (navigate/modal/sheet/http/emit/confirm/close/focus), page, url, method, fields, submit, onSuccess, onConfirm, payload
- `DataSourceRef` — declaration: alias, capabilities (search/filter/sort/paginate/create/update/delete), params
- `ProjectConfig` — `.sigil/sigil.yaml`: version, name, defaults, components, datasources, themes, generation

### Renderer Interface (`renderer/renderer.go`)
- `Renderer` — `Name()`, `Render(*RenderContext)`, `RenderTheme(*ThemeConfig)`, `RenderDataSourceStubs(*DataSourceManifest)`, `SharedComponents([]string)`
- `RenderContext` — page, theme, datasources, registry, project config, go module
- `OutputFile` — path, content, mode
- `DataSourceManifest` — alias, capabilities, fields, relations, endpoints, defaults

### Component Schema (`components/registry.go`)
- `Schema` — type, category, description, props (map of PropDef), actions (map of ActionDef), slots, shortcuts
- `Registry` — `Register()`, `Get()`, `Has()`, `Types()`, `GetSchema()`
- 6 categories: primitives, layouts, data, forms, navigation, overlay

## Patterns to Follow

### CLI Commands
- All commands in `internal/cli/`, one file per command. Bootstrap in `cmd/sigil/main.go` is minimal.
- Version injected at build time: `-ldflags "-X github.com/chrispian/sigil/internal/cli.Version=$(VERSION)"`.
- Global flags (`--no-color`, `--verbose`) set on root command, inherited by all subcommands.

### Config Pipeline
- Parse → Defaults → Validate → Render. Each step is a separate function.
- `config.ParseFile(path)` returns `*Page`. Defaults applied in `defaults.go`. Validation in `validator.go`.
- Validation returns `[]ValidationError` (not a single error). Each error has field path, message, severity.

### Component Registry
- Built-in schemas loaded from `internal/components/builtin/*.yaml` at startup.
- Custom schemas loaded from `.sigil/components/*.schema.yaml` if the project config enables it.
- Components are referenced by `type` string in page configs. Registry validates that all referenced types exist.

### Renderer Architecture
- `Renderer` interface allows swappable backends. Two implementations: `gotempl`, `reactshadcn`.
- Each renderer produces `[]OutputFile` — the engine writes them to disk.
- Renderers are stateless — all context passed via `RenderContext`.
- `--dry-run` flag prevents file writes and prints what would be generated.

### MCP Server
- Standalone JSON-RPC 2.0 over stdio. No external dependencies beyond the sigil CLI itself.
- Tools return structured JSON (never raw errors). Pattern: `{success: bool, data: ..., error: string}`.
- Resources expose pages, component schemas, datasource manifests, and themes as URIs.
- Two prompts provide AI agent guidance for page design and config review.

### File Organization
- Project state lives in `.sigil/` directory: `sigil.yaml` (config), `pages/`, `themes/`, `datasources/`, `components/`.
- Generated code goes to the output directory specified by `--output` or `defaults.output` in project config.
- Preview HTML goes to `.sigil/preview/`.

### Testing
- 20 test files. Standard `testing` package — no mocking frameworks.
- Config tests use inline YAML strings via `config.ParseString()`.
- Renderer tests validate generated output contains expected patterns.
- MCP integration tests send JSON-RPC messages and validate responses.
- Server tests use `httptest.NewServer`.

## Anti-Patterns to Avoid

- **Don't add runtime dependencies.** Sigil is a code generator with a minimal footprint (4 dependencies). Don't add HTTP clients, database drivers, or heavy frameworks.
- **Don't mix config parsing and validation.** These are separate pipeline stages. Parsing produces structs, validation checks them.
- **Don't hardcode renderer-specific logic in the engine.** All framework-specific code belongs in the renderer implementation, not in `engine.go`.
- **Don't generate code without `OutputFile`.** All generated content must flow through the `[]OutputFile` return value so `--dry-run` and `--clean` work correctly.
- **Don't add component types without schemas.** Every component must have a `.yaml` schema in `builtin/` (or custom dir). The registry enforces this.

## Reference Implementations

| Pattern | Reference File | Why it's good |
|---------|---------------|---------------|
| CLI command | `internal/cli/generate.go` | Clean Cobra command with flags, validation, error handling |
| Config parsing + validation | `internal/config/parser.go` + `validator.go` | Separated concerns, structured error returns |
| Renderer implementation | `internal/renderer/reactshadcn/renderer.go` | Clean interface implementation, stateless, produces OutputFile |
| Component schema | `internal/components/builtin/button.yaml` | Concise schema format with props, actions, slots |
| MCP tool | `internal/mcp/tools.go` | Structured JSON responses, input schema validation |

## Build & Run

- **Build:** `make build` → `bin/sigil`
- **Test:** `make test` → `go test ./...`
- **Lint:** `make vet` → `go vet ./...`
- **Install:** `make install` → `go install ./cmd/sigil`
- **Preview:** `sigil preview <page-id>` (opens browser)
- **Dev server:** `sigil serve` (HTTP on port 3000 with live reload)
- **Generate:** `sigil generate --target react-shadcn --output src/generated`
- **MCP server:** `sigil mcp serve` (JSON-RPC 2.0 over stdio)

## Notes

- Sigil is a **code generator**, not a runtime framework. It produces source files that developers integrate into their own projects.
- The project dogfoods itself — `.sigil/pages/` contains page definitions for Sigil's own UI (editor, component browser, etc.).
- Two render targets: `go-templ` (Go/Templ + HTMX) and `react-shadcn` (React/TypeScript + shadcn/ui + Tailwind). The `react-shadcn` renderer supports both App Router (Next.js) and SPA (Vite + react-router-dom) output via `target_mode` in `app.yaml` (or `--target-mode` flag).
- Preview mode generates standalone HTML with mock data — no build step required.
- Pre-commit hooks via lefthook: `gofmt`, `go vet`.
- Version injected at build time via ldflags from `git describe --tags`.
