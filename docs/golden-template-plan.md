# Golden Template: Forge — Infra Dashboard

## Purpose

A Vercel/Linear-inspired infrastructure management dashboard built entirely with Sigil YAML definitions. Serves two goals:

1. **Dogfood Sigil** — exercise the component system, renderer, and config spec to find sharp edges
2. **Agent reference** — a styled, multi-page app shell agents can consult when building frontend

## Status: v1 Complete

All 12 pages generate, build, and render in the demo app. Run `make demo` then visit:
- `localhost:3333/forge-shell` — full app shell with sidebar nav
- `localhost:3333/forge-overview` — dashboard
- `localhost:3333/forge-{page}` — all other pages

## App Concept: "Forge"

Infrastructure management GUI: deployments, logs, monitoring, assets, scheduled jobs, workspaces, settings, users, integrations. Aimed at app/web infra (cloud, metal, AWS, Azure, Hetzner, etc). Vercel/Linear aesthetic.

## Pages (12)

| ID | Title | Key Components |
|----|-------|---------------|
| `forge-shell` | App Shell | sidebar, nav-menu, search-bar, dropdown-menu |
| `forge-overview` | Overview | card (stat), data-table, progress, badge, select |
| `forge-deployments` | Deployments | data-table, search-bar, select, badge |
| `forge-deployment-detail` | Deployment Detail | tabs, card, badge, dropdown-menu (sheet overlay) |
| `forge-logs` | Logs | data-table, search-bar, select (level + time filters) |
| `forge-assets` | Assets | data-table, search-bar, select (provider + type filters) |
| `forge-jobs` | Scheduled Jobs | data-table, search-bar, badge |
| `forge-workspaces` | Workspaces | card grid, data-table, avatar, badge |
| `forge-monitoring` | Monitoring | card (stat), progress, badge, alert icons |
| `forge-users` | Users | data-table, search-bar, select, badge |
| `forge-integrations` | Integrations | card grid, badge, button |
| `forge-settings` | Settings | tabs, input, select, switch, button (danger zone) |

## DataSources (7)

| Alias | File | Used By |
|-------|------|---------|
| Deployment | `forge-deployment.yaml` | overview, deployments, deployment-detail |
| LogEntry | `forge-log-entry.yaml` | logs |
| Asset | `forge-asset.yaml` | assets |
| ScheduledJob | `forge-scheduled-job.yaml` | jobs |
| Workspace | `forge-workspace.yaml` | workspaces |
| User | `forge-user.yaml` | users |
| Integration | `forge-integration.yaml` | integrations |

## Renderer Fixes Made (19)

All fixes in `internal/renderer/reactshadcn/pages.go` unless noted.

### Session 1 (v1 build)
1. **sidebar** — Added `case "sidebar"`. Flex layout with `<aside>` (fixed width, border-r) + `<main>` (flex-1).
2. **nav-menu** — Added `case "nav-menu"`. Vertical nav with Lucide icons, `<a>` links, hover states.
3. **tabs** — Fixed to use `child.Props["value"]` for tab value and `child.Props["label"]` for trigger text. Reads `defaultTab` from tabs props.
4. **template vars** — Added `interpolateTemplateVars()`. Converts `{{data.Alias.field}}` → `{alias?.field}`, `{{param.field}}` → `{params.field}`, `{{item.field}}` → `{item.field}`. Applied to heading, text, badge, card title/description.
5. **avatar** — Template var detection prevents broken initials fallback. Uses JSX expression for alt when template vars present.
6. **badge `info`** — Added mapping `info→outline` in `mapBadgeVariant`.
7. **dropdown separator** — Items with `separator: true` now render `<DropdownMenuSeparator />` instead of `<nil>`.
8. **card title/description** — Now run through `interpolateTemplateVars` for template var support.

### Session 2 (sharp edge sweep)
9. **detail hooks** — Pages with `params` now generate `useXById(params.id)` via `useParams()` from next/navigation. Hook files emit both list and ById exports.
10. **grid iteration** — Grids with `datasource` prop + `{{item.*}}` children wrap in `{data.map((item) => (...))}` with auto-detected key. Updated workspaces/integrations YAMLs to data-driven templates.
11. **badge semantic colors** — `mapBadgeVariant` now returns base variant + Tailwind className. success=emerald, warning=amber, info=blue, danger=destructive. No shadcn component changes needed.
12. **JSX text escaping** — Added `escapeJSXText()`. Converts `<`/`>` to `{"<"}`/`{">"}` in text content. Skips JSX expressions.
13. **dropdown trigger styling** — Trigger gets ghost-button Tailwind classes via className (base-ui Trigger already renders as `<button>`, no `asChild` in shadcn v4).
14. **dropdown item icons** — Items with `icon` prop render Lucide icon component with `mr-2 h-4 w-4` inside `<DropdownMenuItem>`.
15. **search-bar icon** — Renders with `<Search>` Lucide icon positioned absolutely inside the input (`pl-9` padding).
16. **action URL interpolation** — Added `interpolateActionURL()`. Template vars in action URLs become JS template literals: `` `/api/x/${params.id}` ``.
17. **filter/create actions** — Added `filter` and `create` action types to `renderReactActions`. Also added `Datasource`/`Field`/`Value` fields to Action struct (`internal/config/types.go`).
18. **icon sizes** — Added `mapIconSize()`. Accepts `xs`/`sm`/`md`/`lg`/`xl` enum or raw pixel values. Default `md` → `h-5 w-5`.
19. **responsive grid** — Grid now reads `columns_sm`/`columns_md` props and generates responsive Tailwind classes (`grid-cols-1 sm:grid-cols-2 md:grid-cols-3`).

## Open Sharp Edges (1)

| # | Component | Issue | Notes |
|---|-----------|-------|-------|
| 12 | **badge variants** | Sigil uses `[default, success, warning, danger, info]` but shadcn uses `[default, secondary, destructive, outline]`. Agents and users expecting shadcn names hit validation errors. | Partially addressed by fix #11 (semantic colors). Validation could accept both naming schemes. |

## What's Next

### For the Golden Template (v2+)
- **Dark theme** — Generate with the dark theme to match Vercel aesthetic. Test theme token mapping.
- **Richer interactions** — Command palette (cmd-K search), context menus, confirm dialogs.
- **Additional template styles** — GitHub-style, Stripe-style, dashboard variants.

### For Sigil Renderer
- **Component prop audit** — Systematic pass through all renderer components to ensure common use-case props are exposed (e.g., `variant`, `border`, `size`, `className`). Pattern from dogfooding: search-bar needed `variant: ghost`, card needed `border: false`. Every rendered component should support the styling knobs its shadcn counterpart offers, at minimum the top 2-3 most common variations. Audit against shadcn docs + the Forge pages as the reference for what "common" means in practice.

### For Agent Reference
- Move the generated Forge pages to a standalone reference directory (outside demo)
- Create a REFERENCE.md that maps each page to the components and patterns it demonstrates
- Use as the canonical "here's how shadcn pages should look" for frontend agents
