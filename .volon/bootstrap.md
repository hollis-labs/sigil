---
type: bootstrap
iteration: 5
status: ready
updated_at: 2026-03-02
branch: main
---

# Bootstrap — Sigil

## Current State

**Iteration 5** — Sprints 0–4 complete. Full MCP server with 9 tools, 7 resources,
2 prompts. AI agents can create, validate, and inspect Sigil configs via JSON-RPC.

## What exists

- Full architecture docs in `docs/` (8 files)
- Sprint guides in `artifacts/plan/sprints/` (Sprint 0–5, 34 tasks total)
- Go module with Cobra CLI and yaml.v3
- Core types: Page, Component, Action, ProjectConfig, etc.
- YAML parser, defaults, writer, templates
- Config validator with 14 structural checks + deep prop validation (type, required, enum)
- CLI commands: `init`, `new page|datasource|theme`, `list`, `validate`, `generate`, `mcp serve`
- Component registry: 49 types, 6 categories, YAML schemas with full prop definitions
- Schema loader: go:embed for built-in, filesystem for custom (.sigil/components/)
- Renderer engine with Go/Templ implementation (page templates, shared components, theme CSS, handler stubs)
- **MCP server**: JSON-RPC 2.0 over stdio (`sigil mcp serve`)
- **MCP tools**: 9 tools (list/get/create/update pages, validate, list/get components, list/create datasources)
- **MCP resources**: 7 URI types (pages, components, datasources, themes, project)
- **MCP prompts**: sigil_design_page (context for UI generation), sigil_review_config (validation + suggestions)
- 114 tests passing (`go test ./...`)

## Sprint history

- Sprint 0: Project scaffold, core types, parser, validator, CLI skeleton, examples, registry foundation
- Sprint 1: CLI commands (init, new, list), YAML writer, templates
- Sprint 2: 49 YAML schema files, embedded loader, deep prop validation, component detail listing
- Sprint 3: Go/Templ renderer, shared components, theme CSS, handler stubs, generate CLI, E2E tests
- Sprint 4: MCP server, 9 tools, 7 resources, 2 prompts, integration tests

## Next action

**Start Sprint 5** — "Preview + Polish + Dogfood"

Read `artifacts/plan/sprints/sprint-5.md` for the full sprint guide.

## Sprint roadmap

| Sprint | Title | Status | Days |
|--------|-------|--------|------|
| 0 | Project Scaffold + Core Types + Config Parser | **done** | 3 |
| 1 | CLI Commands: init, new, list | **done** | 3 |
| 2 | Component Schemas + Full Registry | **done** | 3 |
| 3 | Go/Templ Renderer | **done** | 4 |
| 4 | MCP Server + Agent Integration | **done** | 3 |
| 5 | Preview + Polish + Dogfood | **next** | 3 |

## Decision log

- 2026-03-01: Project created. CLI/MCP/agent-first approach. GUI deferred to post-MVP.
- 2026-03-01: Go/Templ+HTMX is priority 1 render target. React/shadcn priority 2.
- 2026-03-01: Keyboard shortcuts included in config spec from day 1.
- 2026-03-01: Dogfood strategy — Sigil will design its own GUI using Sigil configs.
- 2026-03-01: Sprint 0 complete. 49 component types registered.
- 2026-03-01: Sprint 1 complete. Full CLI workflow: init → new → list → validate.
- 2026-03-02: Sprint 2 complete. YAML schema files with go:embed. Deep prop validation. ComponentSchemaProvider interface.
- 2026-03-02: Sprint 3 complete. Go/Templ renderer with page generation, shared components, theme CSS, handler stubs.
- 2026-03-02: Sprint 4 complete. MCP server with JSON-RPC 2.0. 9 tools, 7 resources, 2 prompts. Full agent workflow tested.
