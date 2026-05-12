---
type: sprint-guide
sprint: 10
title: "Portfolio Frontend Standardization — react-shadcn SPA + Clockwork Migration"
status: draft
created_at: 2026-05-12
estimated_days: 10
depends_on: [9]
---

# Sprint 10 — Portfolio Frontend Standardization

## Objective

Prove the portfolio-level frontend strategy by collapsing the forked `react-clockwork`
renderer back into `react-shadcn` with an SPA output mode, then re-shipping the Clockwork
demo through the unified renderer. Establish a reusable Wails + Vite + React + shadcn
base template so future portfolio apps scaffold consistently from one source of truth.

This is a proof-of-concept sprint: success here unlocks migrating Forge and Stack
Explorer demos in a later sprint, and standardizing all future portfolio apps on the
same stack.

## Strategic context

The portfolio is converging on **Wails desktop apps** as the primary deployment target.
Wails embeds a webview with static frontend assets and has no Node runtime, so Vite
SPA is the natural fit. Next.js App Router stays available as an opt-in alt target for
hosted SaaS editions, but SPA becomes the portfolio default.

Sigil's job is to be the single source of truth for app shells, layouts, components,
and themes — fixes land once, every portfolio app inherits.

## Current state (2026-05-12)

| Demo | Renderer | Stack | Status |
|------|----------|-------|--------|
| `demo/` (Forge) | `react-shadcn` | Next.js 16 + App Router | Working, gold-standard |
| `stack-explorer-demo/` | `react-shadcn` | Next.js 16 + App Router | Working, 100% generated |
| `demo-clockwork/` | `react-clockwork` (fork) | Vite 6 + react-router-dom 7 | Working, ~2,608 lines of duplicated renderer |

The `react-clockwork` fork was the right *stack* (Vite SPA, matches Wails) for the
wrong *reason* (forking the renderer instead of extending it). About 80% of the fork
duplicates `react-shadcn` logic; only routing, layout shell, and the SSE/active-runs
provider are genuinely new.

## Target end state

- Single `react-shadcn` renderer with `target_mode: spa | app-router` flag
- Clockwork demo regenerates through `react-shadcn` SPA mode with parity vs current output
- `react-clockwork` renderer deleted; duplicate builtin schemas deleted
- Reusable `sigil-app-template/` repo: Wails + Vite + React 19 + react-router-dom 7 +
  shadcn/ui + Tailwind v4, with `.sigil/` skeleton pre-wired (shell page, theme, lens
  provider, mock API client, sample CRUD page)
- Provider abstraction in `app.yaml` supports custom provider components (e.g.,
  Clockwork's `ActiveRunsProvider` + `SSEActiveRunsBridge`) instead of fork-specific
  rendering
- CI smoke test generates one canonical app in both modes and runs `vite build` +
  `next build` to catch drift early

## Phases

### Phase 1 — SPA mode in react-shadcn renderer

Extract the SPA-vs-AppRouter branches behind a `target_mode` flag in `app.yaml`.
Affects only output-shape files in `internal/renderer/reactshadcn/`:

- `layout.go` — emit `App.tsx` (SPA) or `app/(group)/layout.tsx` (AppRouter)
- new `routing.go` — emit `routes.tsx` with `<Route>` JSX for SPA mode
- `pages.go` — branch on Link/useParams/useNavigate import source
- `api.go` — branch on env var prefix (`VITE_` vs `NEXT_PUBLIC_`)

Page tree generation, props serialization, custom-component wiring, theme tokens,
modal/sheet patterns, validators all stay shared.

**Acceptance:** `sigil generate --target react-shadcn` honors `target_mode: spa` from
`app.yaml`; existing Next.js output unchanged when flag absent or `app-router`;
`go build ./...` and `go test ./...` green.

### Phase 2 — Provider abstraction for custom providers

Generalize the existing `providers:` block in `app.yaml` so apps can declare custom
provider components (not just Sigil-known types like `lens`). Clockwork's
`ActiveRunsProvider` and `SSEActiveRunsBridge` become declared providers with
`source:` fields, generated into `lib/` and wrapped around the layout.

**Acceptance:** Clockwork's two hand-written context files migrate into
`.sigil/components/` (or a new `.sigil/providers/` location, TBD in phase task);
generated layout wraps the app with declared providers in declared order;
no hand-written bridge files remain in `demo-clockwork/src/hooks/`.

### Phase 3 — Migrate Clockwork demo to react-shadcn SPA mode

Update `Makefile` `cw-demo-generate` to use `--target react-shadcn` with
`target_mode: spa`. Regenerate `demo-clockwork/` and validate parity:

- All 13 pages generate without errors
- `npm run build` (Vite) succeeds with zero TS errors
- `npm run dev` boots and pages render with mock data
- Active-runs SSE indicator works (provider-wired, not hand-written)
- Visual diff vs current output is acceptable (or improvements documented)

Switch to tmp + rsync output pattern (matches SE demo) so hand-written
`use-*.ts` mock hooks are protected with `--ignore-existing`.

**Acceptance:** `make cw-demo` produces a working app; manual smoke test on
all 13 routes passes; mock hooks survive regeneration.

### Phase 4 — Cleanup the fork

Once Clockwork demo runs through the unified renderer:

- Delete `internal/renderer/reactclockwork/` (~2,608 lines)
- Delete duplicate builtin schemas: `internal/components/builtin/{activity-panel,
  artifact-card, checkpoint-card, comment-list, filter-bar, kanban-board, run-card,
  scope-manager, status-badge, summary-cards}.yaml`
- Remove `react-clockwork` from `internal/cli/generate.go` target registration
- Remove `react-clockwork` from generate help text
- Update boot prompt and `.agentrc/agents/backend.md` to reflect single renderer

**Acceptance:** `go build ./...`, `go test ./...`, and `make cw-demo` all green
after deletions; no remaining references to `react-clockwork` outside of git history.

### Phase 5 — Base template repo

Create `sigil-app-template/` (location TBD — likely sibling to the current portfolio
repos at `~/Projects-apps/sigil-app-template`) as a forkable starter:

- Wails v2 project skeleton (Go + frontend embed)
- `frontend/` = Vite 6 + React 19 + react-router-dom 7 + shadcn/ui + Tailwind v4
- `.sigil/` skeleton: `app.yaml` (SPA mode preset), `themes/default.yaml`, one sample
  page (`sample-list.yaml`), one custom component example
- Pre-wired: lens provider, mock API client, mock data fixtures, theme toggle
- `Makefile` with `generate`, `dev`, `build`, `wails-dev`, `wails-build` targets
- README explaining the fork → rename → customize → ship flow

**Acceptance:** Fresh clone + `make generate && make wails-dev` boots a working
desktop app showing the sample page; one-line rename script swaps app name;
`make wails-build` produces a distributable binary on macOS.

### Phase 6 — Validation, CI, and documentation

Lock in the standardization with guardrails:

- CI smoke test: generate Forge through both modes (`react-shadcn` + `react-shadcn`
  with `target_mode: spa`), run `next build` on App Router output and `vite build`
  on SPA output, fail on either error
- Update `CLAUDE.md` and `.agentrc/agents/backend.md` to document `target_mode` flag
- Update `docs/04_cli-reference.md` with new flag and SPA examples
- Decision record in `docs/` capturing Wails-driven choice and the App Router fallback
  policy
- Update boot prompt (`.agentrc/boot-prompt.md`) to describe the unified renderer

**Acceptance:** CI passes both build modes; docs reflect new flow; next session
can pick up the work without re-deriving the strategy.

## Out of scope (deferred to later sprints)

- Migrating Forge demo from Next.js App Router to Vite SPA
- Migrating Stack Explorer demo from Next.js App Router to Vite SPA
- Sigil's own UI as a Wails app (built from `.sigil/pages/sigil-*.yaml`)
- Auth, multi-tenancy, server-side data fetching for hosted SaaS edition
- Tauri target as alt to Wails

## Decisions locked

- **Portfolio frontend stack:** Vite 6 + React 19 + react-router-dom 7 + shadcn/ui +
  Tailwind v4, wrapped in Wails for desktop, deployable as static SPA on the web.
  *Why:* Wails dictates SPA output; Next.js App Router is incompatible without
  full static export which forfeits its main benefits. *How to apply:* New
  portfolio apps scaffold from `sigil-app-template/`; existing Next.js demos
  (Forge, SE) remain on App Router until a separate migration sprint.
- **Renderer strategy:** Single `react-shadcn` renderer with `target_mode: spa |
  app-router` flag in `app.yaml`. *Why:* ~80% of renderer logic is target-agnostic;
  fork drift is the bigger long-term cost. *How to apply:* No new render targets
  for SPA-vs-Next.js variations — branch inside `react-shadcn` instead.
- **Sigil project repo path:** `~/Projects-apps/sigil` registered as Clockwork
  project. *Why:* Future Clockwork tasks for Sigil work need a project anchor.

## Follow-up candidates

- Migrate Forge demo to Vite SPA + Wails (prove a second app works through the
  unified renderer)
- Migrate Stack Explorer demo to Vite SPA + Wails
- Build Sigil's own UI (`sigil-page-editor`, `sigil-components`, `sigil-themes`,
  `sigil-preview`) as a Wails app of the same shape
- Add `auth` block to `app.yaml` schema once first SaaS edition starts
- Decide where custom providers live: `.sigil/components/` shared with components,
  or new `.sigil/providers/` directory (phase 2 task will decide)
- **Multi-module SPA layout** (surfaced 2026-05-12, phase 1): SPA mode currently
  emits `App.tsx` and `routes.tsx` once per module to the same root paths — the
  second module wins. Resolve in phase 3 when Clockwork demo migration forces
  the choice: single `App.tsx` rooting all modules at distinct paths vs one
  `App.tsx` per module written under `<module>/`. The decision affects how
  custom providers wrap module scopes (intersects with phase 2).

## Known limitations

- Phase 5's `sigil-app-template/` location is not yet decided; likely outside
  the Sigil repo so it can version independently. The phase task will lock this.
- CI smoke test in phase 6 assumes a CI environment; if no CI exists yet for
  Sigil, the task may produce a `make ci-smoke` target instead and defer hosted
  CI wiring.
- **Pre-existing test failure** (`TestDefaultRegistryCategories` in
  `internal/components`, sprint-9 baseline): registry test expects old
  category counts; new builtin schemas under `internal/components/builtin/`
  shifted them. Out of phase 1 scope; fix as a small cleanup before phase 4 so
  `go test ./...` returns to fully green before the renderer-fork deletion.
- **Pre-existing nondeterminism** in `internal/renderer/reactshadcn/pages.go`:
  several emissions iterate Go maps directly (`Record<string, string>` variants,
  some JSX attribute ordering), producing diff churn between regenerations.
  Worth fixing before phase 6's CI smoke test so a generated-output parity
  check is feasible.
