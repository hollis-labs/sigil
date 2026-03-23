---
type: agent-boot
version: 1
updated_at: 2026-03-01
---

# Volon Agent Boot — Sigil

## What is Sigil

Sigil is a system-agnostic UI configuration and code generation tool. It defines UIs
declaratively in YAML, validates them against component schemas, and generates
framework-specific code (Go/Templ+HTMX priority 1, React/shadcn priority 2, static HTML priority 3).

Sigil is CLI/MCP/agent-first. The GUI comes later (and will be built using Sigil itself).

## What is Volon

Volon is a multi-session orchestration system for repository automation. It coordinates
human-directed workflows, task execution, and knowledge artifact generation through a
single-writer Orchestrator that delegates bounded read-only work to sub-agents.

## Ground truth (always read these first)

- `volon.yaml` — system configuration
- `.volon/bootstrap.md` — current iteration state and next actions
- `docs/01_architecture.md` — system architecture and design decisions
- `docs/02_config-spec.md` — Sigil YAML config specification (the core schema)
- `docs/03_component-model.md` — component type system and schemas
- `docs/04_cli-reference.md` — CLI commands
- `docs/05_renderer-contract.md` — renderer plugin interface
- `docs/06_datasource-model.md` — DataSource abstraction
- `docs/07_theme-system.md` — theme tokens and CSS generation
- `docs/08_mcp-integration.md` — MCP server for AI/agent tooling
- `artifacts/plan/sprints/` — sprint guides with detailed task specs

## Core rules

- **Single writer**: Only the Orchestrator may modify tasks, logs, PCC, or bootstrap. Sub-agents are read-only.
- **Ground truth in files, not chat**: Always re-ground from repo artifacts, never rely on conversation context.
- **Minimal diffs**: Small, verifiable steps. Each task should produce incremental changes.
- **Bootstrap boundaries**: Finalize each iteration with `/bootstrap-update` to enable clean restarts.
- **Config-first**: The Sigil YAML config is the single source of truth for UI definitions. All code is generated from it.
- **System-agnostic**: Configs must not contain framework-specific code. Framework details live in renderers only.

## Tech stack

- **Language**: Go 1.22+
- **CLI framework**: Cobra
- **Config format**: YAML (primary), JSON (supported via conversion)
- **Schema validation**: JSON Schema (for validating Sigil configs)
- **Template engine**: Go `text/template` for renderer output
- **Code generation target (priority 1)**: Go Templ + HTMX + Tailwind CSS
- **Testing**: Go standard `testing` package
- **Build**: `go build ./cmd/sigil`

## Project structure

```
sigil/
├── cmd/sigil/            # CLI entry point (main.go)
├── internal/
│   ├── config/           # YAML parsing, validation, config types
│   ├── cli/              # Cobra commands
│   ├── renderer/         # Renderer interface + implementations
│   │   ├── gotempl/      # Go/Templ+HTMX renderer
│   │   ├── react/        # React/shadcn renderer (later)
│   │   └── html/         # Static HTML renderer (later)
│   ├── components/       # Component registry + schemas
│   ├── theme/            # Theme token system + CSS generation
│   └── datasource/       # DataSource manifest parsing + stub generation
├── schemas/              # JSON Schema files for Sigil configs
├── docs/                 # Architecture and reference docs
├── examples/             # Example Sigil configs and generated output
├── artifacts/plan/       # Sprint guides and task specs
└── .volon/               # Volon orchestration state
```

## How to start

1. Read `.volon/bootstrap.md` (if present) for current state and next action.
2. Identify your role: check `.volon/boot/` for your role addendum.
3. For Orchestrator: read sprint guides in `artifacts/plan/sprints/` and begin the canonical loop.

## Your role

Load the addendum for your role:
- **Orchestrator**: `.volon/boot/orchestrator.md` — you drive loops, write state, finalize iterations
- **Architect**: `.volon/boot/architect.md` — planning, ADRs, doc-only work, limited write scope
- **Worker**: `.volon/boot/worker.md` — you execute scoped tasks, return results, read-only
- **Reviewer**: `.volon/boot/reviewer.md` — you scan and summarize, read-only

## Reference map

| Topic | Doc |
|---|---|
| Architecture | `docs/01_architecture.md` |
| Config spec | `docs/02_config-spec.md` |
| Component model | `docs/03_component-model.md` |
| CLI reference | `docs/04_cli-reference.md` |
| Renderer contract | `docs/05_renderer-contract.md` |
| DataSource model | `docs/06_datasource-model.md` |
| Theme system | `docs/07_theme-system.md` |
| MCP integration | `docs/08_mcp-integration.md` |
| Sprint guides | `artifacts/plan/sprints/` |
