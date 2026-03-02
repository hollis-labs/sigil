---
type: specification
version: 1
updated_at: 2026-03-01
---

# Sigil DataSource Model

## Overview

A DataSource is an abstract contract for backend data. Sigil configs reference data
by alias (e.g., `Sprint`, `Task`). The DataSource manifest describes what the data
looks like and what operations it supports, without coupling to any backend framework.

## DataSource Manifest

Stored in `.sigil/datasources/<alias>.yaml`:

```yaml
alias: Sprint
description: "Sprint management"

capabilities:
  - search
  - filter
  - sort
  - paginate
  - create
  - update
  - delete

fields:
  - name: id
    type: integer
    primary: true
    hidden: true
  - name: code
    type: string
    label: "Sprint Code"
    searchable: true
    sortable: true
    required: true
    validation:
      pattern: "^[A-Z]+-\\d+$"
      message: "Must be FORMAT-NUMBER"
  - name: title
    type: string
    label: "Title"
    searchable: true
    sortable: true
    required: true
  - name: status
    type: enum
    label: "Status"
    values: [planning, active, completed, archived]
    filterable: true
    default: planning
  - name: start_date
    type: date
    label: "Start Date"
    sortable: true
  - name: end_date
    type: date
    label: "End Date"
    sortable: true
  - name: created_at
    type: datetime
    label: "Created"
    sortable: true
    readonly: true
  - name: updated_at
    type: datetime
    label: "Updated"
    readonly: true
    hidden: true

relations:
  - name: tasks
    type: hasMany
    target: Task
    foreignKey: sprint_id
    label: "Tasks"

endpoints:
  list: "GET /api/sprints"
  create: "POST /api/sprints"
  read: "GET /api/sprints/{id}"
  update: "PUT /api/sprints/{id}"
  delete: "DELETE /api/sprints/{id}"

defaults:
  sort: { field: created_at, direction: desc }
  pageSize: 25
```

## Field Types

| Type | Go type | Description |
|---|---|---|
| `string` | `string` | Text value |
| `integer` | `int64` | Integer number |
| `float` | `float64` | Decimal number |
| `boolean` | `bool` | True/false |
| `date` | `time.Time` | Date only |
| `datetime` | `time.Time` | Date and time |
| `enum` | `string` | Constrained string (values list) |
| `json` | `json.RawMessage` | Arbitrary JSON |
| `text` | `string` | Long text |

## Field Properties

| Property | Type | Description |
|---|---|---|
| `name` | string | Field identifier (required) |
| `type` | string | Field type (required) |
| `label` | string | Display label |
| `primary` | boolean | Primary key |
| `required` | boolean | Required for create/update |
| `readonly` | boolean | Cannot be modified |
| `hidden` | boolean | Not shown in UI by default |
| `searchable` | boolean | Included in text search |
| `sortable` | boolean | Can be sorted |
| `filterable` | boolean | Can be filtered |
| `default` | any | Default value |
| `validation` | object | Validation rules |

## Endpoint Contract

The `endpoints` section defines the standard REST API patterns. Renderers use these
to generate handler stubs and client-side API calls.

All endpoints follow a standard request/response format:

**List (GET):**
```
GET /api/sprints?search=Q4&filter[status]=active&sort=created_at&direction=desc&page=1&pageSize=25

Response:
{
  "data": [...],
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 25,
    "page_count": 4
  }
}
```

**Create (POST):**
```
POST /api/sprints
{ "code": "SPRINT-001", "title": "Sprint 1" }

Response: 201 { "data": { ... } }
```

**Read (GET):**
```
GET /api/sprints/123

Response: 200 { "data": { ... } }
```

**Update (PUT):**
```
PUT /api/sprints/123
{ "title": "Updated Title" }

Response: 200 { "data": { ... } }
```

**Delete (DELETE):**
```
DELETE /api/sprints/123

Response: 204
```
