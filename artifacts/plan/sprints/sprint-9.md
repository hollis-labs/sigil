---
type: sprint-guide
sprint: 9
title: "Clockwork GUI — React Clockwork Renderer Variant"
status: in_progress
created_at: 2026-05-10
estimated_days: 14
depends_on: [8]
---

# Sprint 9 — Clockwork GUI Variant

## Objective

Port the Clockwork Manifold GUI (React 19 + shadcn/ui SPA) into Sigil's declarative
YAML system. Create a new `react-clockwork` renderer variant that generates a Vite SPA
with react-router, Clockwork's API client pattern, SSE support, and Clockwork-specific
components (kanban board, filter bar, activity panel, etc.).

## Architecture

```
.sigil/
├── sigil.yaml                  # Project config (target: react-clockwork)
├── app.yaml                    # Module definitions + API config + providers
├── themes/
│   └── clockwork-dark.yaml     # Clockwork zinc palette (extends dark)
├── pages/
│   ├── clockwork-board.yaml        # /operations — task table + filters
│   ├── clockwork-task-detail.yaml  # /tasks/:id — detail + tabs + edit
│   ├── clockwork-collections.yaml  # /collections — kanban board
│   ├── clockwork-dashboard.yaml    # /dashboard — widgets + charts
│   ├── clockwork-plans.yaml        # /plans — card grid
│   ├── clockwork-plan-detail.yaml  # /plans/:id — phases + progress
│   ├── clockwork-templates.yaml    # /templates — table + archive
│   ├── clockwork-checkpoints.yaml  # /checkpoints — split panel
│   ├── clockwork-runs.yaml         # /runs — run history
│   ├── clockwork-settings.yaml     # /settings — feature flags + about
│   ├── clockwork-projects.yaml     # /projects — CRUD listing
│   ├── clockwork-sprints.yaml      # /sprints — CRUD listing
│   └── clockwork-epics.yaml        # /epics — CRUD listing
├── datasources/
│   ├── tasks.yaml
│   ├── runs.yaml
│   ├── artifacts.yaml
│   ├── comments.yaml
│   ├── checkpoints.yaml
│   ├── collections.yaml
│   ├── plans.yaml
│   ├── templates.yaml
│   ├── projects.yaml
│   ├── sprints.yaml
│   ├── epics.yaml
│   └── models.yaml
└── components/                  # Custom component TSX sources
    ├── filter-bar.tsx
    ├── activity-panel.tsx
    ├── artifact-card.tsx
    ├── comment-list.tsx
    ├── kanban-board.tsx
    ├── status-badge.tsx
    ├── run-card.tsx
    ├── summary-cards.tsx
    ├── scope-manager.tsx
    └── checkpoint-card.tsx

internal/renderer/reactclockwork/   # New renderer variant
├── renderer.go                 # Renderer implementation (Renderer interface)
├── pages.go                    # Page → .tsx generation (react-router, not next/navigation)
├── shared.go                   # Shared components (shadcn imports)
├── theme.go                    # Theme → globals.css + tailwind.config.ts
├── api.go                      # Generate lib/clockwork-api.ts (full API client)
├── sse.go                      # Generate lib/sse.ts + hooks/use-sse.ts
├── routing.go                  # Generate routes.tsx with all page routes
├── providers.go                # Generate React Context providers
└── layout.go                   # Generate App.tsx shell (nav rail + content area)
```

## New Component Schemas (internal/components/builtin/)

| Schema | Category | Purpose |
|--------|----------|---------|
| `filter-bar.yaml` | navigation | Multi-chip filter row (status, priority, search, comboboxes) |
| `kanban-board.yaml` | data | Draggable column layout with @dnd-kit |
| `activity-panel.yaml` | overlay | SSE-driven live run event feed |
| `artifact-card.yaml` | data | File/URL artifact display with download |
| `comment-list.yaml` | data | Threaded comment UI with reply |
| `status-badge.yaml` | primitives | 9-state task status with color mapping |
| `summary-cards.yaml` | data | Horizontal stat card row (open, in-progress, review, blocked) |
| `run-card.yaml` | data | Compact run card with duration/cost/tokens |
| `scope-manager.yaml` | overlay | Project/sprint/epic assignment dialog |
| `checkpoint-card.yaml` | data | Pending checkpoint with respond/cancel |

## Tasks

### TASK-S9-001 | Priority A | React Clockwork renderer variant

**Context:** Sigil's existing `react-shadcn` renderer generates Next.js-style output
(app/ router, next/navigation imports, SWR hooks). Clockwork is a Vite SPA with
react-router and a custom API client. A new renderer variant is needed.

**What to do:**
Create `internal/renderer/reactclockwork/` as a new `Renderer` implementation:
- Register as `"react-clockwork"` via `init()` + `renderer.RegisterRenderer()`
- `Render()` generates flat `.tsx` pages with `react-router-dom` imports
- `RenderTheme()` generates `globals.css` (Tailwind + Clockwork CSS vars) + `tailwind.config.ts`
- `SharedComponents()` generates shadcn primitives + shared utilities
- `RenderLayout()` generates `App.tsx` shell (nav rail, content area, routing)
- `RenderAPIClient()` generates `lib/clockwork-api.ts` with full CRUD + SSE endpoint config
- `RenderProviders()` generates React Context providers (ApiProvider, ActiveRunsProvider)

**Acceptance criteria:**
- [ ] `sigil generate --target react-clockwork --output src/generated` works
- [ ] Generated output is a valid Vite SPA structure (not Next.js)
- [ ] Pages use `useParams`, `useNavigate` from react-router-dom
- [ ] Theme generates correct CSS variables for Clockwork's zinc palette
- [ ] All 8 Renderer interface methods implemented

---

### TASK-S9-002 | Priority A | Custom component schemas (10 new)

**Context:** Clockwork has domain-specific UI patterns not covered by Sigil's
49 existing component schemas (kanban, filter bar, activity panel, etc.).

**What to do:**
Create 10 new YAML schema files in `internal/components/builtin/`:
- `filter-bar.yaml` — props: statusOptions, priorityOptions, searchable, projectEnabled, sprintEnabled, epicEnabled
- `kanban-board.yaml` — props: columns (array), draggable, inbox (bool)
- `activity-panel.yaml` — props: maxItems, showTokens, showCost, autoScroll
- `artifact-card.yaml` — props: showPreview, showDownload, showDelete (bool)
- `comment-list.yaml` — props: editable, maxLength, showAuthor
- `status-badge.yaml` — props: status (enum of 9 states), size (sm/md/lg)
- `summary-cards.yaml` — props: cards (array of {label, count, variant, icon})
- `run-card.yaml` — props: compact, showMetrics, showTiming
- `scope-manager.yaml` — props: mode (project/sprint/epic), multiSelect
- `checkpoint-card.yaml` — props: showDetail, actions (respond/cancel)

Each schema follows the existing builtin YAML format with type, category, description,
props (type, required, default, description), and actions.

**Acceptance criteria:**
- [ ] 10 new `.yaml` files in `internal/components/builtin/`
- [ ] Each validates with existing schema loader
- [ ] Props defined with types, defaults, descriptions
- [ ] Tests pass (`go test ./internal/components/...`)

---

### TASK-S9-003 | Priority A | Clockwork-dark theme

**Context:** Clockwork uses a specific dark zinc color palette. This needs to be a
Sigil theme file that extends the existing dark theme.

**What to do:**
Create `.sigil/themes/clockwork-dark.yaml`:
- Extends: dark
- Override color tokens to match Clockwork's zinc palette:
  - background: `#09090b` (zinc-950)
  - surface/card: `#18181b` (zinc-900)
  - border: `#27272a` (zinc-800)
  - foreground: `#f4f4f5` (zinc-100)
  - muted: `#71717a` (zinc-500)
  - primary: `#6366f1` (indigo-500)
  - accent: `#8b5cf6` (violet-500)
  - danger: `#ef4444` (red-500)
  - warning: `#f59e0b` (amber-500)
  - success: `#22c55e` (green-500)
  - info: `#3b82f6` (blue-500)
- Add status color tokens for all 9 task statuses (backlog, todo, queued, doing, review, done, blocked, paused, cancelled)
- Include shadcn-compatible CSS variable aliases (`--primary`, `--secondary`, etc.)
- Set typography and radius tokens appropriate for Clockwork's compact HUD style

**Acceptance criteria:**
- [ ] `.sigil/themes/clockwork-dark.yaml` exists with full token set
- [ ] Extends existing dark theme correctly
- [ ] `sigil preview --theme clockwork-dark` works
- [ ] All 9 status colors defined

---

### TASK-S9-004 | Priority A | Page YAML definitions — Core (6 pages)

**Context:** The primary Clockwork screens need to be defined as Sigil YAML page configs.

**What to do:**
Create 6 page YAML files in `.sigil/pages/`:

1. **`clockwork-board.yaml`** — Operations board:
   - Summary cards row (4 stat-cards: Open, In Progress, In Review, Blocked)
   - Filter bar (status chips, priority chips, search, project/sprint/epic comboboxes)
   - Data table with sortable columns (ID, title, status, priority, tags, assignee, updated)
   - SSE refresh event listener
   - URL-backed filter state

2. **`clockwork-task-detail.yaml`** — Task detail:
   - Tabs component (Details, Comments, Artifacts, Subtodos, Logs, Debug)
   - Detail view with all task fields
   - Comment list
   - Artifact cards
   - Subtodo list
   - Activity panel (SSE-driven)
   - Status badge + transition buttons
   - Blocked reason banner
   - Parent plan link

3. **`clockwork-collections.yaml`** — Collections/kanban:
   - Kanban board with columns
   - Inbox section
   - Quick-add dialog (modal)
   - Drag-and-drop between columns

4. **`clockwork-dashboard.yaml`** — Dashboard:
   - 3 tabs (Activity, Mission Control, Usage)
   - Stat cards
   - Charts (heatmap area, bar, pie)
   - Recent runs list

5. **`clockwork-settings.yaml`** — Settings:
   - Tabs (General, Features, About)
   - Feature flag toggles
   - Scheduler status card
   - About section

6. **`clockwork-models.yaml`** — Model catalog:
   - Data table with grouped rows by provider
   - Provider filter combobox
   - Cost/per-model columns

**Acceptance criteria:**
- [ ] 6 page YAML files created
- [ ] Each validates with `sigil validate`
- [ ] References existing component types correctly
- [ ] Reasonable mock layout for each page

---

### TASK-S9-005 | Priority B | Page YAML definitions — Secondary (7 pages)

**Context:** Secondary Clockwork screens for plans, templates, checkpoints, runs,
projects, sprints, and epics.

**What to do:**
Create 7 page YAML files in `.sigil/pages/`:

1. **`clockwork-plans.yaml`** — Plan listing with card grid, phase count, progress bars
2. **`clockwork-plan-detail.yaml`** — Plan detail with stepper phases, child task table
3. **`clockwork-templates.yaml`** — Template table with archive filter, instantiate action
4. **`clockwork-checkpoints.yaml`** — Split-panel: pending list (left) + detail (right), summary cards
5. **`clockwork-runs.yaml`** — Run history with filters (status, agent, date), run cards
6. **`clockwork-projects.yaml`** — Projects CRUD with data table + create/edit modals
7. **`clockwork-sprints.yaml`** — Sprints CRUD with data table + create/edit modals

**Acceptance criteria:**
- [ ] 7 page YAML files created
- [ ] Each validates with `sigil validate`
- [ ] References existing + new component types

---

### TASK-S9-006 | Priority B | Custom component TSX sources

**Context:** The new component schemas need working React TSX implementations
that integrate with Clockwork's API client and real data structures.

**What to do:**
Create 10 custom component TSX files in `.sigil/components/` paired with
`.component-schema.yaml` files:

Each component is a real React component with:
- TypeScript props interface matching the schema
- shadcn/ui imports for primitives
- Clockwork API client integration where needed
- Proper Tailwind styling using Clockwork's zinc palette

**Acceptance criteria:**
- [ ] 10 `.component.yaml` + companion `.tsx` files in `.sigil/components/`
- [ ] Components reference correct shadcn primitives
- [ ] Props match schema definitions

---

### TASK-S9-007 | Priority C | Integration + wiring

**Context:** All pieces need to work together — the renderer, theme, components,
and page definitions must produce a working Clockwork GUI.

**What to do:**
- Wire the new renderer into `internal/cli/generate.go` (add `react-clockwork` as a valid target)
- Update `sigil generate` help text with the new target
- Add `react-clockwork` to the Makefile test/build targets if applicable
- Verify `sigil init` supports the new target
- Run full test suite (`go test ./...`) to ensure no regressions
- Generate output and verify it's structurally valid (runs without TS errors)

**Acceptance criteria:**
- [ ] `sigil generate --target react-clockwork` listed in help
- [ ] Full test suite passes
- [ ] Generated output is a valid Vite SPA structure
- [ ] End-to-end: validate → generate → inspect output works
