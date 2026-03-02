---
type: specification
version: 1
updated_at: 2026-03-01
---

# Sigil Component Model

## Overview

Every UI element in Sigil is a **component** with a registered type. The component
registry maps type strings to schemas that define what props, actions, slots, and
keyboard shortcuts each component supports.

## Component Taxonomy

### Primitives (atomic elements)

| Type | Description | Key Props |
|---|---|---|
| `heading` | Section heading | `level` (1-6), `text` |
| `text` | Text paragraph | `text`, `muted`, `size` |
| `button` | Clickable button | `label`, `variant`, `icon`, `size`, `disabled` |
| `icon-button` | Icon-only button | `icon`, `variant`, `size`, `tooltip` |
| `input` | Text input | `placeholder`, `type`, `required` |
| `textarea` | Multi-line input | `rows`, `placeholder`, `maxLength` |
| `select` | Dropdown select | `options`, `placeholder`, `default` |
| `checkbox` | Checkbox input | `label`, `default` |
| `switch` | Toggle switch | `label`, `default` |
| `badge` | Status badge | `text`, `variant` (default/success/warning/danger/info) |
| `avatar` | User avatar | `src`, `fallback`, `size` |
| `separator` | Visual divider | `orientation` (horizontal/vertical) |
| `progress` | Progress bar | `value`, `max`, `variant` |
| `alert` | Alert message | `title`, `message`, `variant` |
| `label` | Form label | `text`, `for` |
| `icon` | Icon display | `name`, `size`, `color` |

### Layouts (spatial organization)

| Type | Description | Key Props |
|---|---|---|
| `rows` | Vertical stack | `gap`, `padding`, `align`, `justify` |
| `columns` | Horizontal columns | `gap`, `padding`, `align`, `justify` |
| `grid` | CSS grid | `columns`, `columns_md`, `columns_sm`, `gap` |
| `card` | Container with border | `title`, `description`, `padding`, `collapsible` |
| `tabs` | Tabbed panels | `defaultTab`, `variant` |
| `tab` | Single tab content | `label`, `icon`, `value` |
| `split` | Resizable split | `direction`, `defaultSize`, `minSize` |
| `sidebar` | Sidebar + content | `width`, `collapsible`, `defaultCollapsed` |
| `accordion` | Collapsible sections | `type` (single/multiple), `defaultOpen` |
| `accordion-item` | Accordion section | `title`, `icon` |
| `scroll-area` | Scrollable region | `height`, `direction` |
| `spacer` | Empty space | `size` |

### Navigation

| Type | Description | Key Props |
|---|---|---|
| `breadcrumb` | Navigation breadcrumbs | `items` (array of {label, page}) |
| `pagination` | Page navigation | `total`, `pageSize`, `page` |
| `nav-menu` | Navigation menu | `items` (array of {label, page, icon}) |
| `command-palette` | Searchable command list | `commands` |

### Composites (multi-element patterns)

| Type | Description | Key Props |
|---|---|---|
| `modal` | Dialog overlay | `title`, `size`, `closeable` |
| `sheet` | Side panel | `title`, `side` (left/right), `size` |
| `dropdown-menu` | Dropdown with items | `trigger`, `items` |
| `context-menu` | Right-click menu | `items` |
| `tooltip` | Hover tooltip | `content`, `side`, `delay` |
| `popover` | Click popover | `content`, `side`, `align` |
| `confirm-dialog` | Confirmation modal | `title`, `message`, `confirmLabel`, `variant` |
| `toast` | Notification toast | `title`, `message`, `variant`, `duration` |

### Data (data-bound components)

| Type | Description | Key Props |
|---|---|---|
| `data-table` | Sortable/filterable table | `datasource`, `columns`, `pagination`, `selection` |
| `list` | Simple data list | `datasource`, `itemTemplate`, `emptyMessage` |
| `detail-view` | Single record display | `datasource`, `fields`, `layout` |
| `stat-card` | Metric display | `label`, `value`, `trend`, `icon` |
| `chart` | Data chart | `type` (bar/line/pie), `datasource`, `config` |
| `timeline` | Event timeline | `datasource`, `fields` |
| `search-bar` | Search input | `placeholder`, `target`, `debounce` |

### Forms

| Type | Description | Key Props |
|---|---|---|
| `form` | Form container | `fields`, `submit`, `validation` |
| `field-group` | Grouped fields | `label`, `fields` |

## Component Schema Format

Each component type has a schema file in `.sigil/components/` or the built-in
`internal/components/builtin/`:

```yaml
# data-table.schema.yaml
type: data-table
category: data
description: "Sortable, filterable, paginated data table with row selection"
since: "1.0"

props:
  datasource:
    type: string
    required: true
    description: "DataSource alias to query"
  columns:
    type: array
    required: true
    items:
      type: object
      properties:
        key: { type: string, required: true }
        label: { type: string, required: true }
        sortable: { type: boolean, default: false }
        filterable: { type: boolean, default: false }
        render: { type: string, description: "Render component for cell value" }
        width: { type: string, description: "Column width (Tailwind or px)" }
        hidden: { type: boolean, default: false }
        align: { type: string, enum: [left, center, right], default: left }
  pagination:
    type: object
    properties:
      enabled: { type: boolean, default: true }
      pageSize: { type: integer, default: 25 }
      pageSizes: { type: array, items: { type: integer }, default: [10, 25, 50, 100] }
  selection:
    type: string
    enum: [none, single, multiple]
    default: none
  striped: { type: boolean, default: false }
  dense: { type: boolean, default: false }
  bordered: { type: boolean, default: true }
  hoverable: { type: boolean, default: true }
  emptyMessage: { type: string, default: "No data found" }
  loading: { type: boolean, default: false }

actions:
  rowClick:
    description: "Fired when a row is clicked"
    context: { row: "The clicked row data" }
  selectionChange:
    description: "Fired when row selection changes"
    context: { selected: "Array of selected rows" }
  sort:
    description: "Fired when sort changes"
    context: { field: "Sort field", direction: "asc or desc" }
  filter:
    description: "Fired when a filter changes"
    context: { field: "Filter field", value: "Filter value" }

slots:
  toolbar:
    description: "Toolbar area above the table"
    accepts: [button, icon-button, search-bar, select, dropdown-menu]
  empty:
    description: "Content shown when datasource returns no data"
    accepts: [text, icon, button]
  footer:
    description: "Footer area below the table"
    accepts: [pagination, text, button]

shortcuts:
  - key: Enter
    description: "Open selected row"
    when: has-selection
  - key: Delete
    description: "Delete selected (with confirm)"
    when: has-selection
  - key: ArrowUp
    description: "Select previous row"
  - key: ArrowDown
    description: "Select next row"
  - key: Space
    description: "Toggle row selection"
    when: "selection != none"
  - key: Ctrl+A
    description: "Select all rows"
    when: "selection == multiple"
  - key: Escape
    description: "Clear selection"
    when: has-selection

renders:
  - status-chip     # StatusChip component for enum values
  - relative-time   # Relative time display (e.g., "2 hours ago")
  - avatar          # User avatar
  - badge           # Badge display
  - link            # Clickable link
  - boolean         # Checkbox/icon for boolean values
  - currency        # Formatted currency
  - truncate        # Truncated text with tooltip
```

## Component Resolution

When Sigil encounters a component config like `{ type: "data-table", props: {...} }`,
it resolves the type through this pipeline:

1. **Registry lookup**: Find the schema for `data-table` in the component registry
2. **Prop validation**: Validate `props` against the schema's prop definitions
3. **Default application**: Fill in defaults for omitted optional props
4. **Action validation**: Validate any actions against the schema's action definitions
5. **Slot validation**: If children are provided, validate they're in accepted slot types
6. **Shortcut merging**: Merge page-level shortcuts with component-level shortcuts

## Custom Components

Projects can define custom component schemas in `.sigil/components/`:

```yaml
# .sigil/components/sprint-card.schema.yaml
type: sprint-card
category: data
description: "Sprint summary card with progress"
extends: card                    # Inherits card's base props

props:
  sprint:
    type: object
    required: true
    properties:
      code: { type: string }
      title: { type: string }
      status: { type: string }
      progress: { type: number }
      taskCount: { type: integer }
  showProgress: { type: boolean, default: true }

# Custom components need a renderer template for each target
templates:
  go-templ: "sprint-card.templ"
  react: "SprintCard.tsx"
```

Custom components reference their own renderer templates. Sigil generates code using
those templates instead of the built-in ones.

## Render Components

Cell renderers transform raw values into visual representations:

| Render | Input | Output |
|---|---|---|
| `status-chip` | String enum | Colored badge with label |
| `relative-time` | ISO datetime | "2 hours ago", "3 days ago" |
| `avatar` | { src, name } | Circular avatar image |
| `badge` | String | Styled badge |
| `link` | { text, url } | Clickable link |
| `boolean` | Boolean | Check/X icon |
| `currency` | Number | "$1,234.56" |
| `truncate` | String | First N chars with "..." tooltip |
| `progress` | Number (0-100) | Progress bar |
| `json` | Object | Formatted JSON display |
