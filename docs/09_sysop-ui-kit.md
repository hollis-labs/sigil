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

## Dogfood loop — kit mode vs the hand-composed Sysop UI (CW-20260518-0035)

The Sysop UI app (`sysop/`) was hand-composed against `@hollis-labs/sysop-ui`
in CW-20260515-0133. CW-20260518-0035 ran the kit renderer over the matching
`sigil-*` page configs (`sigil-components`, `sigil-datasources`, `sigil-pages`,
`sigil-themes`) to see whether the generated output could replace those screens.

**Outcome: kept hand-composed.** The generated kit pages are valid and import
the kit, but they are not yet at parity, so swapping them in would regress —
and break — the Sysop UI:

- **API envelope mismatch (blocker).** The Sigil server returns single-key
  envelopes (`/api/components` → `{ "components": [...] }`). The generated SWR
  hooks (`fetch(url).then(r => r.json())` typed as `T[]`) and `lib/api.ts`
  (a `{ data, meta }` paginated envelope) assume a different shape, so a
  generated page hands `DataTable` an object instead of an array and breaks at
  runtime. The hand-composed app uses the kit `createApiClient` and unwraps the
  single-key envelope explicitly.
- **No detail view.** The hand-composed screens open a fetch-on-open
  `DetailDialog` per row (`onRowOpen`); the renderer does not map the modal
  `detail` page kind (see above), so generated rows are inert.
- **No `OperationsTablePage` / `MetaList` / `Metric` / `JsonViewer` mapping.**
  The hand-composed screens compose these kit pieces (summary cards, search,
  header refresh action; the Overview screen is entirely `Metric`/`MetaList`/
  `JsonViewer` with no list at all). The renderer emits a bare `DataTable`.
- **`@/components/ui/*` dependency.** Generated pages still import primitives
  the kit does not re-export (`select`, `label`, …). The Sysop UI is a thin
  kit consumer with no local `components/ui/` — adopting generated pages would
  require scaffolding shadcn primitives the kit deliberately omits.

Bug fixed during the loop (`refetch` codegen): the generated SWR list hook
exposed only SWR's `mutate`, while pages destructured a `refetch` binding — a
type error. And that binding was emitted whenever a datasource was CRUD-capable
but only *consumed* when the page YAML carried an explicit `refresh:` field,
leaving it destructured-but-unused (a `noUnusedLocals` failure). The list hook
now exposes a `refetch` alias, and a create modal always re-pulls its
datasource list after a successful create. The `<Type>Input` type import is
now gated on a create/update mutation helper actually being emitted.
