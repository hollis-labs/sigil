---
type: bootstrap
iteration: 1
status: ready
updated_at: 2026-03-01
branch: main
---

# Bootstrap — Sigil

## Current State

**Iteration 1** — Sprint 0 complete. Core types, parser, validator, CLI, component
registry, and example configs all implemented and tested.

## What exists

- Full architecture docs in `docs/` (8 files)
- Detailed sprint guides in `artifacts/plan/sprints/` (Sprint 0–5, 34 tasks total)
- Volon boot infrastructure (`.volon/`, `volon.yaml`, `CLAUDE.md`, `AGENTS.md`)
- Go module with Cobra CLI and yaml.v3 (`go.mod`, `go.sum`)
- Core types in `internal/config/types.go` (Page, Component, Action, etc.)
- YAML parser in `internal/config/parser.go` with defaults in `defaults.go`
- Config validator in `internal/config/validator.go` (14 checks, structured results)
- Cobra CLI with `sigil validate` command (`internal/cli/`)
- Component registry with 49 types across 6 categories (`internal/components/`)
- Example configs: `sprint-dashboard.yaml`, `login-form.yaml`, `invalid-example.yaml`
- 26 tests passing (`go test ./...`)
- `go build ./cmd/sigil` produces working binary

## Sprint 0 Tasks (all complete)

- TASK-S0-001: Initialize Go module and project structure ✓
- TASK-S0-002: Define core config types ✓
- TASK-S0-003: Implement YAML config parser ✓
- TASK-S0-004: Implement config validator ✓
- TASK-S0-005: Implement Cobra CLI skeleton with validate command ✓
- TASK-S0-006: Create example configs ✓
- TASK-S0-007: Component registry foundation ✓

## Next action

**Start Sprint 1** — "CLI Commands: init, new, list"

Read `artifacts/plan/sprints/sprint-1.md` for the full sprint guide.

## Sprint roadmap

| Sprint | Title | Status | Days |
|--------|-------|--------|------|
| 0 | Project Scaffold + Core Types + Config Parser | **done** | 3 |
| 1 | CLI Commands: init, new, list | **next** | 3 |
| 2 | Component Schemas + Full Registry | todo | 3 |
| 3 | Go/Templ Renderer | todo | 4 |
| 4 | MCP Server + Agent Integration | todo | 3 |
| 5 | Preview + Polish + Dogfood | todo | 3 |

## Decision log

- 2026-03-01: Project created. CLI/MCP/agent-first approach. GUI deferred to post-MVP.
- 2026-03-01: Go/Templ+HTMX is priority 1 render target. React/shadcn priority 2.
- 2026-03-01: Keyboard shortcuts included in config spec from day 1.
- 2026-03-01: Dogfood strategy — Sigil will design its own GUI using Sigil configs.
- 2026-03-01: Sprint 0 complete. 49 component types registered. Validator checks types against registry.
