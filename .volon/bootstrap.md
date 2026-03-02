---
type: bootstrap
iteration: 2
status: ready
updated_at: 2026-03-01
branch: main
---

# Bootstrap — Sigil

## Current State

**Iteration 2** — Sprints 0 and 1 complete. Full CLI with init, new (page/datasource/theme),
list, and validate commands. YAML writer for clean config generation.

## What exists

- Full architecture docs in `docs/` (8 files)
- Detailed sprint guides in `artifacts/plan/sprints/` (Sprint 0–5, 34 tasks total)
- Volon boot infrastructure (`.volon/`, `volon.yaml`, `CLAUDE.md`, `AGENTS.md`)
- Go module with Cobra CLI and yaml.v3 (`go.mod`, `go.sum`)
- Core types in `internal/config/types.go` (Page, Component, Action, ProjectConfig, etc.)
- YAML parser (`parser.go`), defaults (`defaults.go`), writer (`writer.go`)
- Config validator (`validator.go`) with 14 checks, structured results
- Page template generator (`templates.go`)
- CLI commands: `init`, `new page|datasource|theme`, `list pages|datasources|themes|components`, `validate`
- Component registry with 49 types across 6 categories (`internal/components/`)
- Example configs: `sprint-dashboard.yaml`, `login-form.yaml`, `invalid-example.yaml`
- 29 tests passing (`go test ./...`)
- `go build ./cmd/sigil` produces working binary

## Sprint 0 Tasks (all complete)

- TASK-S0-001 through TASK-S0-007: all done ✓

## Sprint 1 Tasks (all complete)

- TASK-S1-001: Implement sigil init command ✓
- TASK-S1-002: Implement sigil new page command ✓
- TASK-S1-003: Implement sigil new datasource command ✓
- TASK-S1-004: Implement sigil new theme command ✓
- TASK-S1-005: Implement sigil list command ✓
- TASK-S1-006: YAML serialization (write clean YAML) ✓

## Next action

**Start Sprint 2** — "Component Schemas + Full Registry"

Read `artifacts/plan/sprints/sprint-2.md` for the full sprint guide.

## Sprint roadmap

| Sprint | Title | Status | Days |
|--------|-------|--------|------|
| 0 | Project Scaffold + Core Types + Config Parser | **done** | 3 |
| 1 | CLI Commands: init, new, list | **done** | 3 |
| 2 | Component Schemas + Full Registry | **next** | 3 |
| 3 | Go/Templ Renderer | todo | 4 |
| 4 | MCP Server + Agent Integration | todo | 3 |
| 5 | Preview + Polish + Dogfood | todo | 3 |

## Decision log

- 2026-03-01: Project created. CLI/MCP/agent-first approach. GUI deferred to post-MVP.
- 2026-03-01: Go/Templ+HTMX is priority 1 render target. React/shadcn priority 2.
- 2026-03-01: Keyboard shortcuts included in config spec from day 1.
- 2026-03-01: Dogfood strategy — Sigil will design its own GUI using Sigil configs.
- 2026-03-01: Sprint 0 complete. 49 component types registered. Validator checks types against registry.
- 2026-03-01: Sprint 1 complete. Full CLI workflow: init → new → list → validate.
