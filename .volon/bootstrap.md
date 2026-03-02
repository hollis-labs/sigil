---
type: bootstrap
iteration: 6
status: ready
updated_at: 2026-03-02
branch: main
---

# Bootstrap — Sigil

## Current State

**Iteration 6** — Sprints 0–5 complete. **MVP is complete.** Full CLI + MCP + code
generation + preview + export/import. Dogfood configs created for Sigil's own GUI.

## What exists

- Full architecture docs in `docs/` (8 files)
- Sprint guides in `artifacts/plan/sprints/` (Sprint 0–5, 34 tasks total)
- Go module with Cobra CLI and yaml.v3
- Core types: Page, Component, Action, ProjectConfig, etc.
- YAML parser, defaults, writer, templates
- Config validator with 14 structural checks + deep prop validation (type, required, enum)
- CLI commands: `init`, `new`, `list`, `validate`, `generate`, `preview`, `export`, `import`, `version`, `mcp serve`
- Component registry: 49 types, 6 categories, YAML schemas with full prop definitions
- Schema loader: go:embed for built-in, filesystem for custom (.sigil/components/)
- Renderer engine with Go/Templ implementation (page templates, shared components, theme CSS, handler stubs)
- **Preview**: Standalone HTML with Tailwind CDN, mock data, inline theme CSS, browser open
- **Export/Import**: JSON round-trip for page configs with validation
- **CLI UX**: Color output, --no-color, error suggestions (did you mean?), help examples, version command
- **MCP server**: JSON-RPC 2.0 over stdio (`sigil mcp serve`)
- **MCP tools**: 9 tools (list/get/create/update pages, validate, list/get components, list/create datasources)
- **MCP resources**: 7 URI types (pages, components, datasources, themes, project)
- **MCP prompts**: sigil_design_page, sigil_review_config
- **Dogfood**: 6 page configs for Sigil's own GUI (pages, editor, components, preview, datasources, themes)
- 160 tests passing (`go test ./...`)

## Sprint history

- Sprint 0: Project scaffold, core types, parser, validator, CLI skeleton, examples, registry foundation
- Sprint 1: CLI commands (init, new, list), YAML writer, templates
- Sprint 2: 49 YAML schema files, embedded loader, deep prop validation, component detail listing
- Sprint 3: Go/Templ renderer, shared components, theme CSS, handler stubs, generate CLI, E2E tests
- Sprint 4: MCP server, 9 tools, 7 resources, 2 prompts, integration tests
- Sprint 5: Preview command, export/import, CLI UX polish, dogfood configs

## Next action

**MVP complete.** All 6 sprints done. Next steps are post-MVP:
- GUI implementation using the dogfood Sigil configs
- React/shadcn renderer (priority 2)
- Nanite integration for config sharing
- Plugin system for custom renderers

## Sprint roadmap

| Sprint | Title | Status | Days |
|--------|-------|--------|------|
| 0 | Project Scaffold + Core Types + Config Parser | **done** | 3 |
| 1 | CLI Commands: init, new, list | **done** | 3 |
| 2 | Component Schemas + Full Registry | **done** | 3 |
| 3 | Go/Templ Renderer | **done** | 4 |
| 4 | MCP Server + Agent Integration | **done** | 3 |
| 5 | Preview + Polish + Dogfood | **done** | 3 |

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
- 2026-03-02: Sprint 5 complete. Preview, export/import, CLI polish, 6 dogfood page configs. MVP done.
