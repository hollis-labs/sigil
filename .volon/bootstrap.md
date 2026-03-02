---
type: bootstrap
iteration: 0
status: ready
updated_at: 2026-03-01
branch: main
---

# Bootstrap — Sigil

## Current State

**Iteration 0** — Project initialization. All architecture docs and sprint guides are
complete. No code has been written yet.

## What exists

- Full architecture docs in `docs/` (8 files)
- Detailed sprint guides in `artifacts/plan/sprints/` (Sprint 0–5, 34 tasks total)
- Volon boot infrastructure (`.volon/`, `volon.yaml`, `CLAUDE.md`, `AGENTS.md`)
- Empty Go project structure (`cmd/`, `internal/`, `schemas/`, `examples/`)

## What does NOT exist yet

- No Go code (`go.mod`, `main.go`, etc.)
- No compiled binary
- No tests
- No example config files (those are part of Sprint 0)

## Next action

**Start Sprint 0** — "Project Scaffold + Core Types + Config Parser"

Read `artifacts/plan/sprints/sprint-0.md` for the full sprint guide.

Tasks in priority order:
1. TASK-S0-001 | Priority A | Initialize Go module and project structure
2. TASK-S0-002 | Priority A | Define core config types
3. TASK-S0-003 | Priority A | Implement YAML config parser
4. TASK-S0-004 | Priority A | Implement config validator
5. TASK-S0-005 | Priority A | Implement Cobra CLI skeleton with validate command
6. TASK-S0-006 | Priority A | Create example configs
7. TASK-S0-007 | Priority B | Component registry foundation

## Sprint roadmap

| Sprint | Title | Status | Days |
|--------|-------|--------|------|
| 0 | Project Scaffold + Core Types + Config Parser | **next** | 3 |
| 1 | CLI Commands: init, new, list | todo | 3 |
| 2 | Component Schemas + Full Registry | todo | 3 |
| 3 | Go/Templ Renderer | todo | 4 |
| 4 | MCP Server + Agent Integration | todo | 3 |
| 5 | Preview + Polish + Dogfood | todo | 3 |

## Decision log

- 2026-03-01: Project created. CLI/MCP/agent-first approach. GUI deferred to post-MVP.
- 2026-03-01: Go/Templ+HTMX is priority 1 render target. React/shadcn priority 2.
- 2026-03-01: Keyboard shortcuts included in config spec from day 1.
- 2026-03-01: Dogfood strategy — Sigil will design its own GUI using Sigil configs.
