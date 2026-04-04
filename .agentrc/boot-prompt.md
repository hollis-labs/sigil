# Boot Prompt — Sigil / Stack Explorer

## Context

The **Stack Explorer** is a 10-page analytics dashboard built with Sigil-generated React/shadcn pages. Frontend at `stack-explorer-demo/` (port 3334), backend Go REST API at `/Users/chrispian/Projects-apps/stack-explorer/` (port 8081). Both are fully functional.

## Current state (2026-04-04)

### Sigil Renderer — fully operational
- **App-level config** (`.sigil/app.yaml`): modules, shell binding, API config, features, actions, providers
- **Layout generation**: `RenderLayout` produces `layout.tsx` from `se-shell.yaml` (sidebar, nav, command palette, mobile, theme toggle, lens selector wrapped in LensProvider)
- **API client generation**: `RenderAPIClient` produces `lib/api.ts` from app config
- **Provider generation**: `RenderProviders` produces `lib/lens-context.tsx` from providers config
- **Modal create**: Controlled form state + `postItem()` API call + SWR `mutate` refetch
- **Sheet edit/delete**: `rowClick: type: sheet` with `datasource` + `fields` → Sheet overlay with read/edit (`updateItem`) + delete (`AlertDialog` + `deleteItem`)
- **Custom component system**: `Source` field on Schema (`component` + `includes`), engine copies source files, renderer imports + renders with JSON-serialized props
- **Template references**: `{{varName}}` in custom component props resolves to JS expressions at render time; validator allows template refs for any prop type
- `go build`, `go test ./...` all pass, `next build` compiles cleanly

### Custom Components — 12 registered (2 pre-existing + 10 new)

Pre-existing:
- `sigil-chart` — Recharts wrapper (line, bar, area, pie, donut, radar)
- `scorecard-heatmap` — Self-contained heatmap with lens filtering

Extracted this session (all in `.sigil/components/`):
- `repos-by-category-chart` — Donut chart, takes `repos` prop, buckets by category
- `score-distribution-chart` — Bar chart, takes `scorecards` prop, buckets by score range
- `own-vs-reference-chart` — Grouped bars, takes `repos` + `scorecards` + `dimensionScores`
- `scan-runner` — Run Scan button + Promise.allSettled + polling + scan history table
- `scorecard-compare` — DataTable with checkboxes (max 5) + Compare sheet (radar + delta table)
- `csv-export-button` — Builds CSV from data array + columns, triggers Blob download
- `yaml-import-button` — Hidden file input, reads YAML, POSTs to configurable endpoint
- `file-download-button` — Fetches URL, parses Content-Disposition, triggers download
- `repo-scoring-dialog` — Lens selection + dimension score inputs + batch postItem
- `gap-advantage-charts` — Dual bar charts (advantages green, gaps red) from own-vs-ref deltas

### Frontend — 100% generated
- 10 pages + shell layout, all custom logic in custom components
- `make se-demo-generate` produces a fully functional app with zero hand-edits
- `next build` compiles with zero TypeScript errors
- All custom component source files auto-synced via Makefile

### Backend — complete, no changes needed
- Full CRUD for all 11 resources, YAML import/export, DB backup, scan worker
- `go build`, `go vet`, `go test` all clean

### Renderer fixes applied this session
- `mutate:` → `refetch:` in hook destructuring (hooks expose `refetch`, not `mutate`)
- Sheet state declarations rendered before output (were missing from generated code)
- Sheet imports (`deleteItem`, `updateItem`) collected before `imports.String()` call
- `as Record<string, unknown>` → `as unknown as Record<string, unknown>` for strict TS
- `pluralizeResource()` helper avoids double-s (`"lenss"` → `"lenses"`)
- Lens context provider: `!""` → `!currentId` (was always-falsy)
- Validator: `{{varName}}` template references pass type checking for any prop type

## Potential next steps

1. **Dynamic stat cards** — Dashboard and repo-detail still have hardcoded stat values; could make stat-card accept datasource references
2. **SWR compatibility** — Mock hooks work but real SWR would enable proper caching/revalidation (blocked on SWR 2.x + React 19 incompatibility)
3. **Preview renderer** — Missing render cases for many newer components
4. **Go/Templ renderer** — Hasn't kept pace with React/shadcn renderer additions
5. **Additional custom components** — Any new computed/interactive widgets follow the same extraction pattern

## Key files

| File | What |
|---|---|
| `.sigil/app.yaml` | App config — modules, shell, API, features, providers, actions |
| `.sigil/pages/se-shell.yaml` | Shell page — sidebar, nav, topbar structure |
| `.sigil/pages/se-*.yaml` | All 10 SE page definitions |
| `.sigil/components/*.yaml` | Custom component schemas (12 total) |
| `.sigil/components/*.tsx` | Custom component source files |
| `internal/renderer/reactshadcn/pages.go` | Page generation (modal create, sheet edit/delete, custom components) |
| `internal/renderer/reactshadcn/providers.go` | Context provider generation |
| `internal/renderer/engine.go` | Engine — loads app config, custom schemas, wires layout/API/provider generation |
| `internal/config/validator.go` | Validation — allows `{{ref}}` template props for any type |
| `internal/components/registry.go` | Schema struct with Source field for custom components |
| `Makefile` | `se-demo-generate` target builds, generates, syncs including custom components |
| `stack-explorer-demo/src/hooks/` | Mock hooks (use `--ignore-existing` to preserve) |

## Conventions

- `make se-demo-generate` builds Sigil and syncs all generated files to `stack-explorer-demo/`
- Custom component sources live in `.sigil/components/` alongside their schema YAMLs
- `{{varName}}` in prop values references page datasource variables
- Custom components import generated types from `@/types/` for proper TS compatibility
- Complex props (arrays, objects) are JSON-serialized in generated JSX
- SE demo on port 3334, API on port 8081
- Fix renderer gaps inline with sub-agents rather than deferring
- Backup at `stack-explorer-demo-backup-20260404-102629/` (pre-extraction snapshot)
