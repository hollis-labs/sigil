---
type: specification
version: 1
updated_at: 2026-03-01
---

# Sigil Config Specification

## Overview

A Sigil config is a YAML document that fully describes a UI page. It declares the
layout, components, data bindings, actions, and theme — everything needed to generate
framework-specific code.

## Top-Level Schema

```yaml
sigil: "1.0"                     # Schema version (required)
kind: page                       # Config type: page | component | theme | datasource
id: <kebab-case-id>              # Unique identifier (required)
title: <string>                  # Human-readable title (required)
description: <string>            # Optional description
overlay: <overlay-type>          # Display mode (required for pages)
module: <dot.notation>           # Module grouping (optional)

# Theme (optional — inherits project default if omitted)
theme:
  extends: <theme-id>
  tokens: { ... }                # Token overrides

# DataSource declarations (optional)
datasources:
  - alias: <PascalCase>
    capabilities: [search, filter, sort, paginate, create, update, delete]

# Root layout (required for pages)
layout:
  type: <layout-type>
  props: { ... }
  children: [ ... ]

# Page-level metadata (optional)
meta:
  guards: [auth, admin]          # Access control
  tags: [dashboard, sprint]
  version: 1
```

## Overlay Types

Controls how the page is displayed:

| Value | Description |
|---|---|
| `page` | Full page in the main content area |
| `modal` | Centered dialog overlay |
| `sheet` | Slide-in panel from right |
| `drawer` | Slide-in panel from left |
| `fullscreen` | Full-screen overlay |

## Component Config

Every UI element is a component with this structure:

```yaml
id: <string>                     # Unique within the page (required)
type: <component-type>           # Registered component type (required)
props:                           # Component-specific properties (optional)
  key: value
actions:                         # Event handlers (optional)
  <event>:
    type: <action-type>
    ...
children:                        # Nested components (optional, layout types only)
  - { id, type, props, ... }
condition:                       # Conditional rendering (optional)
  field: <datasource.field>
  op: eq | neq | gt | lt | in | notin
  value: <value>
```

### Component ID Conventions

- Unique within the page
- Descriptive: `table-sprints`, `btn-new-task`, `form-login`, `search-users`
- Used for `refreshTarget` references and action targeting

## Action Config

Actions define what happens on user interaction:

```yaml
actions:
  click:                         # Event name
    type: navigate               # Action type
    page: sprint-detail          # Target page
    params:
      id: "{{row.id}}"           # Template variable from data context
```

### Action Types

| Type | Fields | Description |
|---|---|---|
| `navigate` | `page`, `params`, `url` | Navigate to a page or URL |
| `modal` | `title`, `fields`, `submit`, `refresh` | Open a modal dialog with form |
| `sheet` | `title`, `page`, `side` | Open a sheet panel |
| `http` | `url`, `method`, `body`, `refresh` | Make an HTTP request |
| `emit` | `event`, `payload` | Emit a custom event |
| `confirm` | `title`, `message`, `onConfirm` | Show confirmation dialog |
| `close` | — | Close the current overlay |

### Modal Action Detail

```yaml
actions:
  click:
    type: modal
    title: "Create Sprint"
    size: md                     # sm | md | lg | xl | full
    fields:
      - name: code
        label: Code
        type: text
        required: true
        placeholder: "SPRINT-001"
        validation:
          pattern: "^[A-Z]+-\\d+$"
          message: "Must be FORMAT-NUMBER"
      - name: title
        label: Title
        type: text
        required: true
      - name: status
        label: Status
        type: select
        options:
          - { value: planning, label: Planning }
          - { value: active, label: Active }
        default: planning
      - name: start_date
        label: Start Date
        type: date
    submit:
      datasource: Sprint
      method: POST
    refresh: table-sprints       # Component ID to refresh after submit
    onSuccess:
      type: close                # Close modal on success
```

### HTTP Action Detail

```yaml
actions:
  click:
    type: http
    url: "/api/sprints/{{row.id}}/archive"
    method: POST
    confirm:
      title: "Archive Sprint?"
      message: "This will archive sprint {{row.code}}."
    refresh: table-sprints
```

## Template Variables

Actions and props can reference data context using `{{...}}` syntax:

| Pattern | Source |
|---|---|
| `{{row.field}}` | Current row in a data-table |
| `{{form.field}}` | Current form values |
| `{{param.name}}` | URL/route parameters |
| `{{page.field}}` | Page-level variables |

## DataSource Declaration

Pages declare which data sources they need:

```yaml
datasources:
  - alias: Sprint
    capabilities: [search, filter, sort, paginate, create, update, delete]
  - alias: Task
    capabilities: [search, filter, sort, paginate]
    params:
      filter:
        sprint_id: "{{param.id}}"  # Scoped to a parent entity
```

The actual DataSource manifest (field definitions, endpoints) lives in
`.sigil/datasources/<alias>.yaml`. See `docs/06_datasource-model.md`.

## Layout Types

| Type | Description | Children |
|---|---|---|
| `rows` | Vertical stack with gap | Yes |
| `columns` | Horizontal columns with gap | Yes |
| `grid` | CSS grid with configurable columns | Yes |
| `card` | Container with optional header/footer | Yes |
| `tabs` | Tabbed content panels | Yes (tab children) |
| `split` | Resizable split panes | Yes (exactly 2) |
| `sidebar` | Sidebar + main content | Yes (exactly 2) |

### Layout Props

```yaml
layout:
  type: rows
  props:
    gap: 4                       # Gap between children (Tailwind spacing scale)
    padding: 4                   # Container padding
    align: start | center | end | stretch
    justify: start | center | end | between | around
    className: "custom-class"    # Escape hatch for custom CSS
```

### Grid Layout

```yaml
layout:
  type: grid
  props:
    columns: 3                   # Number of columns
    columns_md: 2                # Responsive: medium screens
    columns_sm: 1                # Responsive: small screens
    gap: 4
```

## Form Fields

Used in modal actions and form components:

```yaml
fields:
  - name: email
    label: Email Address
    type: email
    required: true
    placeholder: "user@example.com"
    validation:
      pattern: "^[^@]+@[^@]+$"
      message: "Enter a valid email"

  - name: role
    label: Role
    type: select
    options:
      - { value: admin, label: Administrator }
      - { value: user, label: User }
      - { value: viewer, label: Viewer }
    default: user

  - name: bio
    label: Biography
    type: textarea
    rows: 4
    maxLength: 500

  - name: active
    label: Active
    type: switch
    default: true

  - name: tags
    label: Tags
    type: multiselect
    options:
      - { value: frontend, label: Frontend }
      - { value: backend, label: Backend }
```

### Field Types

| Type | HTML equivalent | Extra props |
|---|---|---|
| `text` | `<input type="text">` | `placeholder`, `maxLength`, `pattern` |
| `email` | `<input type="email">` | `placeholder` |
| `password` | `<input type="password">` | `placeholder` |
| `number` | `<input type="number">` | `min`, `max`, `step` |
| `textarea` | `<textarea>` | `rows`, `maxLength` |
| `select` | `<select>` | `options`, `default` |
| `multiselect` | Multi-select | `options` |
| `checkbox` | `<input type="checkbox">` | `default` |
| `switch` | Toggle switch | `default` |
| `date` | Date picker | `min`, `max` |
| `datetime` | Datetime picker | `min`, `max` |
| `file` | File upload | `accept`, `maxSize` |
| `hidden` | Hidden field | `value` |

## Keyboard Shortcuts (GUI Phase)

Components that support keyboard interaction declare shortcuts in their config:

```yaml
# Page-level keyboard shortcuts
shortcuts:
  - key: Escape
    action:
      type: close              # Close current overlay
  - key: Ctrl+N
    action:
      type: modal
      title: "Create New..."
  - key: Ctrl+K
    action:
      type: emit
      event: command-palette    # Open command palette
  - key: /
    action:
      type: focus
      target: search-bar       # Focus a component by ID
    when: "!input-focused"     # Only when no input is focused

# Component-level keyboard shortcuts
- id: table-sprints
  type: data-table
  props:
    datasource: Sprint
  shortcuts:
    - key: Enter
      action:
        type: navigate
        page: sprint-detail
        params: { id: "{{selected.id}}" }
    - key: Delete
      action:
        type: confirm
        title: "Delete?"
        onConfirm:
          type: http
          url: "/api/sprints/{{selected.id}}"
          method: DELETE
      when: "has-selection"
    - key: ArrowDown
      action: { type: emit, event: select-next }
    - key: ArrowUp
      action: { type: emit, event: select-prev }
```

### Shortcut Schema

```yaml
shortcuts:
  - key: <key-expression>        # Key or combo: Escape, Ctrl+N, Shift+Enter
    action: <action-config>      # Standard action config
    when: <condition>            # Optional: when to enable
    description: <string>        # For keyboard shortcut help display
    global: false                # If true, fires even when inputs are focused
```

### Key Expression Format

| Expression | Meaning |
|---|---|
| `Escape` | Escape key |
| `Enter` | Enter key |
| `Ctrl+N` | Control + N |
| `Ctrl+Shift+K` | Control + Shift + K |
| `ArrowDown` | Down arrow |
| `/` | Forward slash |
| `?` | Question mark (shows shortcut help) |

### Standard Shortcuts (recommended defaults)

| Key | Action | Context |
|---|---|---|
| `Escape` | Close current overlay / deselect | Global |
| `Ctrl+K` | Open command palette | Global |
| `Ctrl+N` | Create new (context-dependent) | Global |
| `/` | Focus search | When no input focused |
| `?` | Show keyboard shortcuts | When no input focused |
| `Enter` | Open selected item | Data tables |
| `Delete` | Delete selected (with confirm) | Data tables |
| `ArrowUp/Down` | Navigate list | Data tables, lists |

## Full Example

```yaml
sigil: "1.0"
kind: page
id: sprint-dashboard
title: Sprint Dashboard
description: "Overview of all sprints with task counts"
overlay: page
module: project.sprints

theme:
  extends: default

datasources:
  - alias: Sprint
    capabilities: [search, filter, sort, paginate, create, update, delete]

shortcuts:
  - key: Escape
    action: { type: close }
  - key: Ctrl+N
    action:
      type: modal
      title: "New Sprint"
      fields:
        - { name: code, label: Code, type: text, required: true }
        - { name: title, label: Title, type: text, required: true }
      submit: { datasource: Sprint, method: POST }
      refresh: table-sprints
  - key: /
    action: { type: focus, target: search-sprints }
    when: "!input-focused"

layout:
  type: rows
  props:
    gap: 4
    padding: 6
  children:
    - type: columns
      props:
        justify: between
        align: center
      children:
        - id: page-title
          type: heading
          props:
            level: 2
            text: Sprint Dashboard
        - id: btn-new-sprint
          type: button
          props:
            label: New Sprint
            variant: primary
            icon: plus
          actions:
            click:
              type: modal
              title: "Create Sprint"
              fields:
                - { name: code, label: Code, type: text, required: true }
                - { name: title, label: Title, type: text, required: true }
                - { name: status, label: Status, type: select, options: [{value: planning, label: Planning}], default: planning }
              submit: { datasource: Sprint, method: POST }
              refresh: table-sprints

    - id: search-sprints
      type: search-bar
      props:
        placeholder: "Search sprints..."
        target: table-sprints

    - id: table-sprints
      type: data-table
      props:
        datasource: Sprint
        columns:
          - { key: code, label: Code, sortable: true }
          - { key: title, label: Title, sortable: true }
          - { key: status, label: Status, filterable: true, render: status-chip }
          - { key: created_at, label: Created, sortable: true, render: relative-time }
        pagination: { enabled: true, pageSize: 25 }
        selection: single
        striped: true
        dense: false
      actions:
        rowClick:
          type: navigate
          page: sprint-detail
          params: { id: "{{row.id}}" }
      shortcuts:
        - key: Enter
          action: { type: navigate, page: sprint-detail, params: { id: "{{selected.id}}" } }
        - key: Delete
          action:
            type: confirm
            title: "Delete Sprint?"
            message: "Are you sure you want to delete {{selected.code}}?"
            onConfirm:
              type: http
              url: "/api/sprints/{{selected.id}}"
              method: DELETE
              refresh: table-sprints
          when: has-selection

meta:
  guards: [auth]
  tags: [dashboard, sprint, project]
  version: 1
```
