---
intent: pcc_global
project: sigil
updated_at: "2026-03-04"
---

# Architecture

## Tech stack

- **Language:** Go 1.25
- **Module:** `github.com/chrispian/sigil`
- **Dependencies:** cobra (CLI), yaml.v3
- **Config format:** YAML page definitions with component tree structure

## Key components

### CLI (`cmd/sigil`)
Main binary with subcommands: init, new, validate, generate, preview, serve, list, export, import, diff, migrate, schema, doctor, version, mcp.

### Internal packages (`internal/`)
- `cli` -- CLI command implementations (cobra)
- `components` -- 49 component type schemas and validation
- `config` -- project/page config loading and parsing
- `datasource` -- DataSource model (declare data shapes, generate typed hooks/handlers)
- `mcp` -- MCP server (JSON-RPC 2.0 over stdio; 9 tools, 7 resources, 2 prompts)
- `renderer` -- render target implementations
- `server` -- live dev server with auto-reload
- `theme` -- theme token parsing and CSS variable generation

## Render targets

### go-templ
Generates: `.templ` page components, shared `.templ` components, Go handler stubs per datasource, `theme.css`, Tailwind preset.

### react-shadcn
Generates: `.tsx` page components, DataTable + barrel exports, SWR hooks per datasource, TypeScript interfaces, `cn()` utility, `globals.css`, `tailwind.config.ts`.

## Data flow

1. YAML page configs loaded from `.sigil/pages/`
2. Validated against component schemas (type, props, required fields)
3. Renderer target selected (go-templ or react-shadcn)
4. Source files generated to output directory
5. Theme tokens compiled to CSS custom properties

## Project structure on disk

```
.sigil/
  sigil.yaml, pages/, themes/, datasources/, components/
```

## Evidence
- Last refreshed: 2026-03-04 (mentat PCC bootstrap)
- Sources: README.md, internal/, cmd/sigil/, docs/
