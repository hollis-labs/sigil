# Renderer Context — Sigil

> Context for working on Sigil's code generation renderers. Loaded by the renderer agent when working on output targets.
> Lives at `sigil/.agentrc/agents/renderer.md`.

## What This Agent Covers

Sigil generates framework-specific source code from YAML page definitions. This context covers the renderer implementations, component schemas, and generated output patterns. Use this agent when:

- Adding or modifying component schemas
- Working on the React/shadcn or Go/Templ renderers
- Adding a new render target
- Fixing generated output issues
- Working on the preview system or dev server

## Render Targets

### React/shadcn (`react-shadcn`)
- **Output:** `.tsx` pages, component files, SWR hooks, TypeScript types, `globals.css`, `tailwind.config.ts`
- **Stack:** React, TypeScript, shadcn/ui, Tailwind CSS, Lucide icons, SWR for data fetching
- **Location:** `internal/renderer/reactshadcn/`
- **Files:** `renderer.go` (interface impl), `pages.go` (page → .tsx), `shared.go` (shared components), `theme.go` (theme → CSS + Tailwind config)

### Go/Templ (`go-templ`)
- **Output:** `.templ` pages, shared components, Go handler stubs, `theme.css`, `tailwind.sigil.cjs`
- **Stack:** Go/Templ, HTMX, Tailwind CSS (via preset)
- **Location:** `internal/renderer/gotempl/`
- **Files:** `renderer.go` (interface impl), `pages.go` (page → .templ), `shared.go` (shared components), `theme.go` (theme → CSS), `stubs.go` (Go handler stubs)

### Preview (`preview`)
- **Output:** Standalone HTML files with inline CSS and mock data
- **Stack:** Plain HTML/CSS, no build step, no JS framework
- **Location:** `internal/renderer/preview/`
- **Files:** `preview.go` (page → HTML with mock data)

## Component System

### 49 Built-in Components (6 categories)

**Primitives (15):** heading, text, button, icon-button, input, textarea, select, checkbox, switch, badge, avatar, separator, progress, alert, label, icon

**Layouts (12):** rows, columns, grid, card, tabs, tab, split, sidebar, accordion, accordion-item, scroll-area, spacer

**Navigation (5):** breadcrumb, pagination, nav-menu, stepper, command-palette

**Data (6):** data-table, list, detail-view, stat-card, chart, tree, timeline

**Overlay (8):** modal, sheet, drawer, popover, tooltip, dropdown-menu, context-menu, confirm-dialog

**Forms (3):** form, field-group, search-bar

### Schema Format (`internal/components/builtin/*.yaml`)

```yaml
type: button
category: primitives
description: Clickable button
props:
  label:
    type: string
    required: true
  variant:
    type: string
    default: secondary
    enum: [primary, secondary, destructive, outline, ghost, link]
  icon:
    type: string
  size:
    type: string
    default: md
    enum: [sm, md, lg]
  disabled:
    type: boolean
    default: false
actions:
  click:
    description: Fired when button is clicked
```

### Adding a New Component

1. Create `internal/components/builtin/<type>.yaml` with schema
2. Add rendering logic to each renderer:
   - `gotempl/pages.go` — Templ template generation for this component
   - `reactshadcn/pages.go` — TSX generation for this component
   - `preview/preview.go` — HTML preview rendering for this component
3. Add test coverage in the corresponding `*_test.go` files
4. The registry auto-discovers schemas from `builtin/` — no manual registration needed

## Renderer Interface

```go
type Renderer interface {
    Name() string
    Render(ctx *RenderContext) ([]OutputFile, error)
    RenderTheme(theme *ThemeConfig) ([]OutputFile, error)
    RenderDataSourceStubs(ds *DataSourceManifest) ([]OutputFile, error)
    SharedComponents(usedTypes []string) ([]OutputFile, error)
}
```

- `Render` — Generate page files from a page config
- `RenderTheme` — Generate theme CSS/config files
- `RenderDataSourceStubs` — Generate data layer stubs (hooks, handlers)
- `SharedComponents` — Generate shared component files for all used types

All methods return `[]OutputFile` (path + content + mode). The engine handles file I/O.

## Theme System

Theme files (`.sigil/themes/*.yaml`) define design tokens:

```yaml
name: dark
tokens:
  colors:
    background: "9 9 11"
    surface: "24 24 27"
    text: "244 244 245"
    accent: "79 70 229"
  radius:
    md: "0.375rem"
    lg: "0.5rem"
```

Rendered as CSS custom properties (`--sigil-*`) by each renderer's `theme.go`.

## DataSource System

DataSource manifests (`.sigil/datasources/*.yaml`) define data contracts:

```yaml
alias: Task
capabilities: [search, filter, sort, paginate, create, update, delete]
fields:
  - name: id
    type: string
    primary: true
  - name: title
    type: string
    required: true
    searchable: true
  - name: status
    type: string
    filterable: true
    values: [todo, in_progress, done]
```

Renderers use this to generate:
- **React:** SWR hooks (`useTask.ts`), TypeScript types (`task.ts`)
- **Go/Templ:** Handler stubs (`task_handler.go`)
- **Preview:** Mock data based on field types

## Generated Output Examples

### React/shadcn page (`pages/task-board.tsx`)
```tsx
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

export default function TaskBoard() {
  const { data: task, isLoading } = useTask();
  return (
    <div className="flex flex-col gap-6 p-6">
      <h1>Task Board</h1>
      <Button variant="default"><Plus className="mr-2 h-4 w-4" />New Task</Button>
    </div>
  );
}
```

### Go/Templ page (`pages/task_board.templ`)
```templ
templ TaskBoard() {
    <div class="flex flex-col gap-6 p-6">
        <h1>Task Board</h1>
        <button class="sigil-btn sigil-btn-primary"
            hx-get="/ui/modals/create-task"
            hx-target="#modal-container">
            <span class="icon icon-plus"></span> New Task
        </button>
    </div>
}
```

## Patterns to Follow

- **Renderers are stateless.** All context flows through `RenderContext`. No global state.
- **All output via `OutputFile`.** Never write files directly — return `[]OutputFile` so `--dry-run` and `--clean` work.
- **Component schemas are the source of truth.** If a component isn't in the registry, it doesn't exist.
- **Each renderer handles all 49 components.** When adding a component, all three renderers must support it.
- **Props map directly.** Component props from YAML map to framework props (React) or template attrs (Templ). Don't invent intermediate representations.
- **Actions map to framework idioms.** `navigate` → React Router / HTMX hx-get. `http` → fetch / HTMX. `modal` → dialog open.

## Anti-Patterns to Avoid

- **Don't add framework logic to `engine.go`.** Framework-specific code belongs in the renderer implementation.
- **Don't generate code that requires runtime Sigil.** Generated code must be standalone — no sigil runtime library.
- **Don't skip renderers.** A component must render in all targets or explicitly document unsupported status.
- **Don't hardcode shadcn import paths.** Use `@/components/ui/` prefix consistently — it's the shadcn convention.
- **Don't generate inline styles in React output.** Use Tailwind classes. Inline styles are only for preview HTML.

## Live Dev Server

- `sigil serve` starts HTTP on port 3000
- Routes: `/` (index), `/pages/{id}` (preview), `/assets/theme.css`, `/events` (SSE)
- File watcher monitors `.sigil/` for changes and pushes SSE events for auto-reload
- Implementation: `internal/server/` (server.go, watcher.go, index.go)

## Demo App

Next.js 16 + Tailwind v4 + shadcn/ui v4 at `demo/`. Real rendering of generated output.

- `make demo-generate` — build, generate react-shadcn, sync to demo
- `make demo` — generate + start dev server (port 3333)
- Demo hooks use mock data (SWR incompatible with React 19)
- shadcn v4 uses `@base-ui/react` — no `asChild` prop on triggers
- Tooltip triggers need `render={<span />}` to avoid nested `<button>` hydration errors
- Accordion has no `type="single"` / `collapsible` props in v4

## Current Status (Sprint 9)

### React renderer components (24 of 50 schemas have render cases)

**Fully rendered:** button, badge, input, label, select, search-bar, tabs, tab, modal, sheet, separator, progress, avatar, alert, icon, text, heading, spacer, form, data-table, card, switch, dropdown-menu, accordion, accordion-item, combobox, tooltip

**Schemas exist, no React render case yet:** breadcrumb, chart, checkbox (inline only), command-palette, confirm-dialog, context-menu, detail-view, dropdown-menu separators/labels, field-group, grid responsive, icon-button, list, nav-menu, pagination, popover (standalone), scroll-area, sidebar, split, stat-card, timeline, toast

### Shared generated components
- `DataTable` — TanStack React Table wrapper
- `Combobox` — stateful Command + Popover wrapper with search/select

### Theme system
- Light base theme with shadcn-compatible tokens
- Dark theme extending light
- CSS output emits both `--sigil-*` and shadcn standard vars (`--background`, `--primary`, etc.)
- Tailwind config includes shadcn semantic colors + sigil namespace
- Shadow tokens prefixed `--sigil-shadow-*` to avoid radius key collision

### Known issues
- SWR 2.x incompatible with React 19 / Next 16 — demo uses mock hooks
- Preview renderer doesn't have cases for card, switch, accordion, combobox, dropdown-menu, tooltip
- Go/Templ renderer doesn't have cases for any of the new components
- `make demo-generate` uses `--ignore-existing` for hooks but is fragile

## Reference Files

| Need | Reference |
|------|-----------|
| New component schema | `internal/components/builtin/button.yaml` |
| React rendering | `internal/renderer/reactshadcn/pages.go` |
| Shared components | `internal/renderer/reactshadcn/shared.go` |
| Templ rendering | `internal/renderer/gotempl/pages.go` |
| Preview rendering | `internal/renderer/preview/preview.go` |
| Theme generation (React) | `internal/renderer/reactshadcn/theme.go` |
| Theme generation (Templ) | `internal/renderer/gotempl/theme.go` |
| DataSource stubs | `internal/renderer/gotempl/stubs.go` |
| Demo app | `demo/` |
| Component showcase | `.sigil/pages/component-showcase.yaml` |
| Complete example | `examples/task-board-demo/` |
