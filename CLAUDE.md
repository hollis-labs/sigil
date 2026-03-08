# Sigil — agent boot

## Agent Auto-Boot Override

This repo boots as `worker`. **Skip profile selection** — go directly to:

4. Read `.agentrc/agent-boot.md`
5. Read `.agentrc/boot/worker.md`
6. Read `.agentrc/bootstrap.md`
7. Follow the worker profile instructions — emit boot confirmation and begin work

## Project Overview

Sigil is a system-agnostic UI configuration and code generation tool. It defines UIs declaratively in YAML, validates them against component schemas, and generates framework-specific code (Go/Templ+HTMX, React/shadcn, static HTML).

## Build & Test

```bash
go build ./cmd/sigil/
go test ./...
```

## Architecture

- `cmd/sigil/` — Entry point, Cobra CLI
- `internal/cli/` — CLI command implementations
- `internal/components/` — Component type system and schemas
- `internal/config/` — Configuration loading and validation
- `internal/datasource/` — DataSource abstraction
- `internal/mcp/` — MCP server for AI/agent integration
- `internal/renderer/` — Code generation renderers
- `internal/server/` — Live dev server
- `internal/theme/` — Theme tokens and CSS generation

## Key docs

| File | Purpose |
|---|---|
| `agentrc.yaml` | System configuration |
| `.agentrc/bootstrap.md` | Current iteration state — start here |
| `.agentrc/agent-boot.md` | Full boot reference (rules, reference map) |
| `.agentrc/boot/` | Role-specific boot profiles |
| `docs/01_architecture.md` | System architecture |
| `docs/02_config-spec.md` | Sigil config YAML specification |
| `docs/03_component-model.md` | Component type system |
| `docs/04_cli-reference.md` | CLI commands reference |

## Profiles

| Profile | Use when |
|---|---|
| `orchestrator` | Driving implementation, managing tasks, writing to the repo |
| `architect` | Planning, ADRs, doc-only work — limited write scope |
| `worker` | Bounded read-only analysis (scans, audits, reports) |
| `reviewer` | Structured reviews, verification, targeted investigation |
