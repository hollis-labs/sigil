---
intent: pcc_global
project: sigil
updated_at: "2026-03-04"
---

# Backlog & Priorities

## Task tracking

Sigil has Volon integration (`volon.yaml` with nanite backend) but no `.volon/tasks/` directory exists yet. The `.volon/` directory contains `agent-boot.md`, `boot/` profiles, `bootstrap.md`, and `templates/`. No backlog directory found.

## Active areas (from codebase)

- 49 component types implemented across 6 categories
- Two render targets operational (go-templ, react-shadcn)
- MCP server with 9 tools, 7 resources, 2 prompts
- Live dev server with auto-reload
- Theme system with design tokens
- DataSource model for typed hooks/handlers
- Config diffing and schema migration

## Inferred priorities

1. **Additional render targets** -- current: go-templ, react-shadcn; architecture supports plugin renderers
2. **Component library expansion** -- beyond 49 types as user needs emerge
3. **Custom component schemas** -- `.sigil/components/*.schema.yaml` support exists
4. **MCP adapter maturation** -- 9 tools available; potential for expansion
5. **IDE integration** -- JSON Schema export for VS Code YAML extension

## Known gaps

- No tasks or backlog files created yet despite Volon config
- No GitHub issues or external task tracker detected
- Bootstrap exists but no iteration history

## Evidence
- Last refreshed: 2026-03-04 (mentat PCC bootstrap)
- Sources: .volon/, volon.yaml, README.md, internal/
