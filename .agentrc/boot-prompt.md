# Boot Prompt — Sigil / Stack Explorer

## Context

The **Stack Explorer** is a 10-page analytics dashboard built with Sigil-generated React/shadcn pages. Frontend at `stack-explorer-demo/` (port 3334), backend Go REST API at `/Users/chrispian/Projects-apps/stack-explorer/` (port 8081). Both are fully functional with complete CRUD.

## Current state (2026-04-03)

### Frontend — complete
- 10 pages + shell layout, all reads wired to live API
- 7 modal Create forms (POST to API, refetch on success)
- 5 sheet overlays with Edit (PUT) and Delete (DELETE + AlertDialog confirm)
- Repo Detail: edit modal, delete with redirect, scoring dialog (creates scorecard + dimension scores)
- Dashboard: 5 live stat cards, 3 computed charts (repos by category, score distribution, own vs reference)
- Scorecards: live chart from scorecard data, CSV export
- Gap Analysis: computed advantage/gap charts from dimension scores
- Quick actions: layout dropdown, command palette (Cmd+K), dashboard buttons all navigate with `?action=create` to auto-open modals
- Pages using `useSearchParams` wrapped in `<Suspense>` (repos, reports, gap-analysis)
- `next build` compiles 12 routes clean

### Backend — complete
- Full CRUD for all 11 resources (repos, tags, snapshots, dimensions, lenses, scorecards, scores, patterns, findings, comparison-sets, reports)
- YAML import/export: `POST /api/repos/import`, `GET /api/repos/export`
- DB backup: `GET /api/db/backup`
- `go build`, `go vet`, `go test` all clean

### Sigil renderer — no changes needed
- `go test ./...` all pass, no renderer modifications in recent waves

## Remaining work

### Completed (2026-04-03)
- Import YAML — wired on Repos page + Settings page (file picker → POST /api/repos/import)
- Export YAML — wired on Settings page (GET /api/repos/export → browser download)
- DB Backup — wired on Settings page (GET /api/db/backup → browser download)

### Compare mode (Scorecards) — complete (2026-04-03)
- Checkbox selection on DataTable (up to 5 repos)
- "Compare (N)" button opens right-side Sheet (width scales with selection count)
- Overall scores in color-coded solid boxes
- Radar chart overlay of dimension scores (added `type="radar"` to SigilChart)
- Delta table with green/red badges (when comparing exactly 2)
- Files modified: `scorecards/page.tsx`, `sigil-chart.tsx`

### Run Scan (Dashboard) — complete (2026-04-04)
- Backend: POST/GET /api/scans endpoints, background scan worker goroutine
- Scan worker: polls pending scans every 5s, fetches GitHub API stats (stars, forks, issues, last push), creates snapshots, transitions status
- Worker file: `stack-explorer/internal/api/scan_worker.go`
- Frontend: "Run Scan" queues one scan per repo (se-repo-scan blueprint), polls for completion, scan history table with status icons
- New frontend files: `types/scan.ts`, `hooks/use-scan.ts`
- Modified: `dashboard/page.tsx`

## Key files

| File | What |
|---|---|
| `stack-explorer-demo/src/lib/api.ts` | API client — fetchList, fetchItem, postItem, updateItem, deleteItem |
| `stack-explorer-demo/src/hooks/use-*.ts` | 11 hooks with refetch support |
| `stack-explorer-demo/src/app/(explorer)/*/page.tsx` | 10 page components |
| `stack-explorer-demo/src/app/(explorer)/layout.tsx` | Shell layout, sidebar, command palette, lens selector |
| `stack-explorer-demo/src/components/data-table.tsx` | Customized DataTable (don't overwrite) |
| `stack-explorer-demo/src/components/sigil-chart.tsx` | Chart wrapper (don't overwrite) |
| `stack-explorer-demo/src/components/scorecard-heatmap.tsx` | Custom heatmap (don't overwrite) |
| `internal/renderer/reactshadcn/pages.go` | Sigil React renderer (modal actions at ~line 1864) |

## Conventions

- `make se-demo-generate` builds Sigil and syncs generated pages — hooks, data-table, sigil-chart protected by `--ignore-existing`
- Don't run `make se-demo-generate` for page tweaks — it overwrites hand-edits
- SE demo on port 3334, API on port 8081
- shadcn v4: `render={<Component />}` not `asChild`, Select `onValueChange` can be `null`, Tooltip uses `delay`
- Generated types use `unknown` for numeric fields — cast with `Number()` or `String()`
- Fix renderer gaps inline with sub-agents rather than deferring
