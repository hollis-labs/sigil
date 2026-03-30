# Boot Prompt — Sigil Frontend Agent

## Context

You are continuing work on the **Forge golden template** — a 16-page Vercel/Linear-style infra dashboard generated entirely from Sigil YAML. Boot the `sigil-frontend` agent before starting.

## What was done

### Session 3 — Mock app consolidation
All Forge pages now render inside a shared shell layout as a single working app:

**Shell & routing:**
- Created `(forge)` route group with shared `layout.tsx` (sidebar + topbar)
- Sidebar with active nav highlighting, responsive mobile hamburger (Sheet)
- Topbar with search, theme toggle (sun/moon), notification bell with badge count, "+" create dropdown
- All routes stripped of `forge-` prefix (`/overview`, `/deployments`, etc.)
- Renderer generates module-aware routes — breadcrumbs, navigate actions, nav-menu all use clean paths
- Module pages render with `<>` fragment (no `min-h-screen` wrapper) to fit inside shell layout

**Dark theme:**
- `next-themes` ThemeProvider in root layout, defaults to dark
- Sun/moon toggle in topbar

**New renderer cases (session 3):**
stat-card, icon-button, confirm-dialog, scroll-area, list, field-group, split, detail-view, pagination, timeline, chart, textarea (standalone), checkbox (standalone), date-picker, slider, context-menu, toast (action type)

**Prop audit fixes:**
button `size`, input `label`/`disabled`/`defaultValue`, alert `title`, avatar `size`, text `size`, separator `orientation`

**Loading states:**
All pages with datasources show a Loader2 spinner while data loads

**Components exercised in Forge:**
- Overview: stat-cards with icons + trends, icon-button with tooltip, toast on refresh
- Deployments: context-menu with navigate/toast actions on table rows
- Deploy-create: textarea, checkbox, labels on inputs, toast on submit
- Deployment-detail: breadcrumb, toast on redeploy
- Logs: date-pickers for start/end, scroll-area, pagination
- Monitoring-detail: sliders for threshold config
- Settings: confirm-dialog wrapping delete button
- Command-palette: Cmd+K keyboard shortcut

### Previous sessions
- Session 1: Light theme, shadcn v4 demo app, initial pages
- Session 2: 12 sharp edges fixed, 8 new renderers, 4 new pages (total 16)

## Current state

- **16 Forge pages** inside shared shell layout, 7 datasources
- **56 renderer component cases** in pages.go
- **34 installed shadcn components** in demo
- **56 component schemas** in builtin/
- Go builds, all tests pass, `make demo` works on port 3333
- `next build` compiles 27 routes clean

## What's next — prioritized

### 1. Search wiring
Topbar search and page-level search bars are static inputs. Wire them to filter datasources or navigate.

### 2. Empty states
Tables show nothing when data is empty/filtered. Add empty state messaging.

### 3. More pages to exercise remaining components
Components with renderers but not yet in Forge pages: list, field-group, split, detail-view, timeline, chart, pagination (only on logs so far).

Suggested:
- Timeline on deployment-detail logs tab
- Chart placeholder on monitoring detail
- Split layout on asset-detail

### 4. Form state management
Form inputs on deploy-create and settings are uncontrolled. Could wire up form state for a more realistic feel.

### 5. Remaining polish
- Breadcrumbs on all sub-pages (only 4 pages have them)
- Active "New Deployment" dropdown items for Asset, Job, Workspace
- Notification bell → notification panel/dropdown

## Key files

| File | What |
|---|---|
| `internal/renderer/reactshadcn/pages.go` | All React/shadcn rendering logic (56 cases) |
| `internal/config/types.go` | Action struct, Component struct |
| `internal/config/validator.go` | Action type validation (9 types) |
| `internal/components/builtin/*.yaml` | Component schemas (56) |
| `.sigil/pages/forge-*.yaml` | Forge page YAML definitions (16) |
| `.sigil/datasources/forge-*.yaml` | Forge datasource definitions (7) |
| `demo/src/app/(forge)/layout.tsx` | Shared shell layout (hand-crafted) |
| `demo/src/app/(forge)/*/page.tsx` | Generated page output (16 pages) |
| `demo/src/hooks/use-*.ts` | Mock data hooks (don't overwrite) |
| `Makefile` | `demo-generate` syncs forge pages to `(forge)/` route group |

## Important conventions

- `make demo-generate` builds sigil, generates output, syncs to demo
- Forge pages sync to `demo/src/app/(forge)/{page}/page.tsx` with `forge-` prefix stripped
- Non-forge pages sync to `demo/src/app/{page}/page.tsx` as before
- Mock hooks in `demo/src/hooks/` use static data (SWR incompatible with React 19)
- shadcn v4 uses `@base-ui/react` — no `asChild` prop on triggers
- Fix renderer gaps inline using sub-agents rather than deferring
- Pages with `module` set render with `<>` fragment, others with `min-h-screen` wrapper
- `pageRoute()` helper strips module prefix from page IDs for route generation
