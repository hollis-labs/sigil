---
intent: pcc_global
project: sigil
updated_at: "2026-03-04"
---

# Architectural Decisions

No formal ADRs recorded yet. Key design decisions inferred from codebase:

## Inferred design choices

- **YAML as config format:** page definitions use YAML component trees with typed props, actions, and data bindings
- **Two-target renderer architecture:** go-templ and react-shadcn as pluggable render targets with a shared renderer contract (`docs/05_renderer-contract.md`)
- **49 component types:** comprehensive built-in set across 6 categories; extensible via custom schema files
- **DataSource abstraction:** declare data shapes in YAML; generators produce typed hooks (SWR for React) or handler stubs (Go)
- **Theme as CSS variables:** design tokens compile to `--sigil-*` CSS custom properties; Tailwind integration via preset/config
- **MCP server over stdio:** JSON-RPC 2.0 for AI agent integration; page CRUD, validation, component info without leaving conversation
- **Preview with mock data:** instant browser preview without backend; lowers friction for page design iteration
- **Config diffing:** semantic comparison rather than text diff; useful for review workflows
- **Schema migration:** auto-fix missing IDs and set defaults when config format evolves
- **Volon-managed development:** nanite storage backend, worktrees, quality scans; boot profiles for different agent roles

## Evidence
- Last refreshed: 2026-03-04 (mentat PCC bootstrap)
- Sources: README.md, CLAUDE.md, docs/, volon.yaml
