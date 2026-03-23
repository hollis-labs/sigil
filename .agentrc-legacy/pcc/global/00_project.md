---
intent: pcc_global
project: sigil
updated_at: "2026-03-04"
---

# Project Identity

**Sigil** turns YAML page definitions into working UI code for multiple framework targets. Define pages once as a component tree, validate against 49 component schemas, and generate framework-specific source files.

## Goals

- Declarative UI configuration via YAML page definitions
- Code generation for multiple targets: Go/Templ+HTMX (`go-templ`) and React/shadcn+Tailwind (`react-shadcn`)
- 49 built-in component types across 6 categories (layout, primitive, form, data, overlay, navigation)
- Live dev server with auto-reload on config changes
- Instant browser preview with mock data
- Deep prop validation with "did you mean?" suggestions
- MCP server for AI/agent integration (9 tools, 7 resources, 2 prompts)
- Theme system with design tokens as CSS variables
- DataSource model for typed data hooks/handlers
- Config diffing, schema migration, JSON Schema export

## Non-goals

- Not a UI framework itself; generates code for existing frameworks
- Not a runtime -- output is static source files
- No backend logic generation (handler stubs only)

## Active configuration

- `volon.yaml`: nanite storage backend, PCC enabled, worktrees enabled
- Module: `github.com/chrispian/sigil` (Go 1.25)
- Binary: `sigil`
- Project config: `.sigil/sigil.yaml`

## Evidence
- Last refreshed: 2026-03-04 (mentat PCC bootstrap)
- Sources: README.md, CLAUDE.md, volon.yaml, go.mod
