# Sigil

## agentrc
- If `.agentrc/boot-prompt.md` exists, read it first for session context.
- If the user says "Boot <agent>", look up the agent in `.agentrc/config.yaml` under `agents:`. Load each role file from `~/.agentrc/roles/` (using the `file:` path from `~/.agentrc/config.yaml` role definitions), load the listed skills, and read the project context file from `.agentrc/` if specified.
- If the user says "Boot <role>" and no agent matches, fall back to loading that single role from `~/.agentrc/roles/` by type directory (domain/, stack/, meta/).
- After context compaction, re-read the active role and project context files.
- Do not guess when uncertain. Stop and ask.
- Prefer focused, minimal output. No trailing summaries.
- Sub-agent output stays in the sub-agent. Main context gets one-line confirmations.

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

## Demo App

Next.js 16 + Tailwind v4 + shadcn/ui v4 demo at `demo/`. Renders real shadcn components from Sigil-generated output.

```bash
make demo-generate   # build sigil, generate react output, sync to demo
make demo            # generate + start dev server on port 3333
```

**Important:** Demo hooks in `demo/src/hooks/` use mock data (SWR incompatible with React 19). The `make demo-generate` target uses `--ignore-existing` to avoid overwriting them. If hooks get overwritten, restore mock versions manually.

## Key docs

| File | Purpose |
|---|---|
| `.agentrc/config.yaml` | Agent definitions for this project |
| `docs/01_architecture.md` | System architecture |
| `docs/02_config-spec.md` | Sigil config YAML specification |
| `docs/03_component-model.md` | Component type system |
| `docs/04_cli-reference.md` | CLI commands reference |
