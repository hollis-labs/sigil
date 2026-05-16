---
type: specification
version: 1
updated_at: 2026-05-16
---

# Sysop UI kit integration (`--ui-kit sysop`)

The `react-shadcn` renderer can generate pages that import shared components
from **`@hollis-labs/sysop-ui`** — the Sysop UI kit — instead of per-app
hand-written component `.tsx`. This is opt-in via the `--ui-kit` flag; the
default (`--ui-kit` omitted) keeps the legacy per-app shadcn output untouched.

```sh
sigil generate --target react-shadcn --target-mode spa --ui-kit sysop ...
```

## What changes in kit mode

| Concern            | Legacy (`--ui-kit` empty)            | Kit mode (`--ui-kit sysop`)                          |
|--------------------|--------------------------------------|------------------------------------------------------|
| Data tables        | app-local `@/components/data-table` (tanstack `accessorKey`) | kit `DataTable` (`items` / `columns` with function `cell` / `getRowId` / `onRowOpen` / `emptyState`) |
| Status badges      | hand-rolled variant-map `<Badge>`    | kit `<StatusBadge>` (theme-aware tone lookup)        |
| Empty states       | `emptyMessage` string                | kit `<EmptyState variant=... />`                     |
| Filter bars        | hand-written custom `filter-bar.tsx` | kit `<FilterBar>` + `<FilterCycleToggle>` chips      |
| shadcn primitives  | `@/components/ui/*`                  | kit barrel for the 13 it re-exports; the rest stay `@/components/ui/*` |
| App shell (SPA)    | hand-emitted sidebar/topbar + `next-themes` + `@/components/ui/*` | thin `App.tsx` composed from kit `NavRail` + `PageHeader` + `ThemeSwitcher` — matches the folio `sysop-ui` preset shape |
| Theme              | generated `globals.css` + `tailwind.config.ts` (Tailwind v3) | thin `globals.css` importing `@hollis-labs/sysop-ui/theme.css` (Tailwind v4, `[data-theme]` palettes) |

The kit only re-exports a subset of shadcn primitives (badge, button, command,
dialog, input, input-group, popover, scroll-area, skeleton, sonner, table,
textarea, tooltip). Primitives it does **not** ship — `select`, `card`,
`label`, `tabs`, `separator`, … — still resolve to app-local `@/components/ui/*`.

## How filter state is wired

Filter controls (header buttons with `emit`/`filter` actions, `select`
filters, and the kit `FilterBar`'s `FilterCycleToggle` chips) all converge on
one shared `useState` per `datasource`+`field` pair via `registerFilter`. The
`data-table` reads those bindings and applies `.filter()` chains to its
`items`, so a chip click "actually feeds" the table. Search is registered the
same way via `registerSearch`. An empty filter value (`value: ""`) normalizes
to `"all"`, which the table treats as unfiltered.

## Implementation

- Flag: `--ui-kit` → `GenerateConfig.UIKit` → `RenderContext.UIKit` /
  `LayoutContext.UIKit`, plus `renderer.UIKitAware.SetUIKit` for the
  contextless interface methods.
- Kit logic lives in `internal/renderer/reactshadcn/kit.go`, gated on
  `importTracker.kitMode()`. Legacy paths are unchanged.

## App shell (SPA)

In kit mode the SPA `App.tsx` is emitted by `renderKitLayoutSPA` as a thin
composition of the kit shell — the same shape the folio `sysop-ui` preset
scaffolds:

```tsx
<div className="flex h-screen bg-bg text-text">
  <NavRail items={nav} logo={...} logoLabel="..." />
  <div className="flex min-w-0 flex-1 flex-col">
    <PageHeader title={title}><ThemeSwitcher /></PageHeader>
    <main>{/* routed page */}</main>
  </div>
</div>
```

react-router routing and per-module provider wrapping are retained (Sigil apps
are multi-page/multi-module). `applyTheme(getInitialTheme())` runs at module
load. The `PageHeader` title is route-derived from the nav table.

## Known limitations / follow-ups

- `summary-cards` and the modal `detail` page kind (kit `DetailDialog`) are not
  yet kit-mapped — they fall back to legacy emission.
- The kit shell drops the legacy command palette, the topbar actions dropdown,
  and the topbar provider entity-selector (e.g. the lens switcher) — providers
  still wrap and stay functional, but their switcher UI is not yet rendered.
- App-router (non-SPA) kit-mode still emits the legacy shell; only SPA is
  migrated.
- Generated pages still render their own in-page `heading`; combined with the
  shell `PageHeader` this can double up. A `heading` → `PageHeader` mapping is
  a follow-up.
- In kit mode the engine still copies hand-written custom-component `.tsx`
  (e.g. `filter-bar.tsx`) into `components/` even when the page uses the kit
  equivalent. The copies are unused; a project-wide `tsc` may still check them.
