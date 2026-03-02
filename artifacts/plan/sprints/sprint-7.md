---
type: sprint-guide
sprint: 7
title: "React/shadcn Renderer"
status: todo
created_at: 2026-03-02
estimated_days: 4
depends_on: [3]
---

# Sprint 7 — React/shadcn Renderer

## Objective

Implement the second renderer target: React with shadcn/ui + Tailwind. This generates
`.tsx` files, React hooks, and a Tailwind config from Sigil configs. Developers get a
complete React project scaffold.

After this sprint, `sigil generate --target react-shadcn` produces a working React app.

## Tasks

### TASK-S7-001 | Priority A | React renderer scaffold

**Context:** The renderer interface already exists. We need a new implementation at
`internal/renderer/reactshadcn/` that registers as "react-shadcn".

**What to do:**

Create `internal/renderer/reactshadcn/renderer.go`:
- Implement `Renderer` interface (Name, Render, RenderTheme, RenderDataSourceStubs, SharedComponents)
- Register via `init()` function
- Target output: `.tsx` files using shadcn component imports

**Acceptance criteria:**
- [ ] `react-shadcn` renderer registered and discoverable
- [ ] `sigil generate --target react-shadcn` runs without error

---

### TASK-S7-002 | Priority A | Page component generation

**Context:** Each Sigil page becomes a React component file.

**What to do:**

Create `internal/renderer/reactshadcn/pages.go`:
- Map Sigil component tree → JSX/TSX tree
- Layout types (rows, columns, grid) → flex/grid divs with Tailwind classes
- Primitives → shadcn components:
  - `button` → `<Button variant={...}>`
  - `input` → `<Input placeholder={...} />`
  - `select` → `<Select>` with options
  - `badge` → `<Badge variant={...}>`
  - `data-table` → `<DataTable columns={...} data={...} />`
  - `form` → `<Form>` with fields
  - `modal` → `<Dialog>`
  - `tabs` → `<Tabs>`
  - `search-bar` → `<Input type="search" />`
  - `alert` → `<Alert variant={...}>`
- Actions → `onClick` handlers, `useRouter` for navigation
- DataSource bindings → React Query/SWR hooks

**Acceptance criteria:**
- [ ] Page components generate valid TSX
- [ ] All component types have React equivalents
- [ ] Actions map to appropriate React patterns

---

### TASK-S7-003 | Priority A | Shared components + hooks

**Context:** Shared components and data-fetching hooks should be generated once.

**What to do:**

- Generate `components/data-table.tsx` (reusable, typed DataTable component)
- Generate `hooks/use-datasource.ts` (React Query wrapper for datasource endpoints)
- Generate `lib/types.ts` (TypeScript types from datasource field definitions)
- Generate `components/ui/index.ts` (re-export shadcn components used by pages)

**Acceptance criteria:**
- [ ] DataTable component generated with column definitions
- [ ] useDatasource hook generated with proper types
- [ ] TypeScript types match datasource field definitions

---

### TASK-S7-004 | Priority B | Theme + Tailwind config

**Context:** Theme tokens should map to Tailwind config and CSS variables.

**What to do:**

- Generate `tailwind.config.ts` with sigil theme tokens
- Generate `globals.css` with `:root` CSS variables
- Map Sigil tokens → Tailwind theme extensions (colors, fonts, border-radius)

**Acceptance criteria:**
- [ ] tailwind.config.ts extends theme with Sigil tokens
- [ ] CSS variables applied to components

---

### TASK-S7-005 | Priority B | Tests + E2E

**What to do:**

- Unit tests for each page component generation
- Test shared component output
- E2E test: full project → valid TSX files that can be type-checked
- Verify generated TypeScript compiles (optional: use `tsc --noEmit` in CI)

**Acceptance criteria:**
- [ ] Component rendering tests for all types
- [ ] E2E generation test with validation
