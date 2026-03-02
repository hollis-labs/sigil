---
type: sprint-guide
sprint: 2
title: "Component Schemas + Full Registry"
status: todo
created_at: 2026-03-01
estimated_days: 3
depends_on: [1]
---

# Sprint 2 — Component Schemas + Full Registry

## Objective

Define full schemas for all core component types (~45 components) with prop definitions,
action definitions, slots, and keyboard shortcuts. Load schemas from YAML files. The
validator uses these schemas for deep prop-level validation.

After this sprint, `sigil validate` catches prop type errors, missing required props,
and invalid values — not just unknown component types.

## Tasks

### TASK-S2-001 | Priority A | Schema loader from YAML files

**Context:** Component schemas need to be loadable from YAML files (both built-in and
custom). Built-in schemas ship embedded in the binary. Custom schemas live in
`.sigil/components/`.

**What to do:**

Use Go's `embed` package to embed built-in schemas:
```go
//go:embed builtin/*.yaml
var builtinSchemas embed.FS
```

Create `internal/components/loader.go`:
- `LoadBuiltinSchemas()` — load from embedded FS
- `LoadCustomSchemas(dir string)` — load from `.sigil/components/`
- Parse YAML → `Schema` structs → register in `Registry`

**Acceptance criteria:**
- [ ] Built-in schemas embedded in binary
- [ ] Custom schemas loaded from disk
- [ ] Custom schemas can override built-in ones (same type name)
- [ ] Parse errors reported with file path

---

### TASK-S2-002 | Priority A | Define primitive component schemas

Create YAML schema files in `internal/components/builtin/`:

**Components (16):**
heading, text, button, icon-button, input, textarea, select, checkbox,
switch, badge, avatar, separator, progress, alert, label, icon

Each schema file follows the format in `docs/03_component-model.md`.
Include all props with types, required flags, defaults, and descriptions.

**Acceptance criteria:**
- [ ] 16 primitive schema files created
- [ ] All props match `docs/03_component-model.md` tables
- [ ] Schemas load without errors
- [ ] `sigil list components` shows all primitives

---

### TASK-S2-003 | Priority A | Define layout component schemas

**Components (12):**
rows, columns, grid, card, tabs, tab, split, sidebar, accordion,
accordion-item, scroll-area, spacer

Focus on props that affect rendering: gap, padding, align, justify,
columns, direction, collapsible, etc.

**Acceptance criteria:**
- [ ] 12 layout schema files created
- [ ] Layout-specific props defined (gap, padding, columns, etc.)
- [ ] Schemas load without errors

---

### TASK-S2-004 | Priority A | Define data + form + nav + composite schemas

**Data (7):** data-table, list, detail-view, stat-card, chart, timeline, search-bar
**Forms (2):** form, field-group
**Navigation (4):** breadcrumb, pagination, nav-menu, command-palette
**Composites (8):** modal, sheet, dropdown-menu, context-menu, tooltip, popover,
confirm-dialog, toast

**Acceptance criteria:**
- [ ] All ~21 schema files created
- [ ] data-table schema is comprehensive (columns, pagination, selection, toolbar slot)
- [ ] Schemas load without errors
- [ ] `sigil list components` shows all categories

---

### TASK-S2-005 | Priority A | Deep validation using component schemas

**Context:** Sprint 0's validator checks for required fields and known types. Now
extend it to validate props against the component schema.

**What to do:**

Extend `internal/config/validator.go`:

1. For each component, look up its schema in the registry
2. Check required props are present
3. Check prop types match (string vs integer vs boolean vs array vs object)
4. Check enum values are valid
5. Warn on unknown props (not in schema)
6. Validate slot children against `accepts` lists

**Acceptance criteria:**
- [ ] Missing required props → error
- [ ] Wrong prop type → error
- [ ] Invalid enum value → error
- [ ] Unknown prop → warning (not error)
- [ ] Slot children validated against `accepts`
- [ ] Tests cover all validation paths
- [ ] Example configs still validate cleanly

---

### TASK-S2-006 | Priority B | Render component schemas (list detail)

**Context:** `sigil list components` currently shows names. Add a detail view:
`sigil list components --type data-table` shows the full schema with props.

**Acceptance criteria:**
- [ ] `sigil list components --type data-table` shows props, actions, slots, shortcuts
- [ ] Props show type, required, default, description
- [ ] Output is readable and useful as a reference
