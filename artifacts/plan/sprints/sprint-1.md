---
type: sprint-guide
sprint: 1
title: "CLI Commands: init, new, list"
status: todo
created_at: 2026-03-01
estimated_days: 3
depends_on: [0]
---

# Sprint 1 — CLI Commands: init, new, list

## Objective

Implement the remaining CLI commands: `sigil init`, `sigil new page|datasource|theme`,
and `sigil list`. After this sprint, a developer can initialize a Sigil project, create
configs from templates, and list their inventory — all from the command line.

## Prerequisites

- [ ] Go module compiles: `go build ./cmd/sigil`
- [ ] Core types defined in `internal/config/types.go`
- [ ] Parser works: `ParseFile()` returns `Page` structs
- [ ] Validator works: `Validate()` returns structured results
- [ ] `sigil validate` command works
- [ ] Example configs exist in `examples/`

## Tasks

### TASK-S1-001 | Priority A | Implement `sigil init` command

**Context:** `sigil init` creates the `.sigil/` directory structure with a project
config, default theme, and empty directories. See `docs/04_cli-reference.md`.

**What to do:**

Create `internal/cli/init.go`:
- Create `.sigil/sigil.yaml` (project config with defaults)
- Create `.sigil/pages/` (empty directory with `.gitkeep`)
- Create `.sigil/components/` (empty)
- Create `.sigil/datasources/` (empty)
- Create `.sigil/themes/default.yaml` (default dark theme from `docs/07_theme-system.md`)
- Print success message with next steps

**Flags:** `--name <string>`, `--theme <dark|light>`

Create the project config type in `internal/config/types.go`:
```go
type ProjectConfig struct {
    Version string `yaml:"version"`
    Name    string `yaml:"name"`
    Defaults struct {
        Theme    string `yaml:"theme"`
        Renderer string `yaml:"renderer"`
        Output   string `yaml:"output"`
    } `yaml:"defaults"`
    Components struct {
        Builtin   bool   `yaml:"builtin"`
        CustomDir string `yaml:"custom_dir"`
    } `yaml:"components"`
    DataSources struct {
        Dir string `yaml:"dir"`
    } `yaml:"datasources"`
    Themes struct {
        Dir string `yaml:"dir"`
    } `yaml:"themes"`
    Generation struct {
        Clean     bool   `yaml:"clean"`
        GoPackage string `yaml:"go_package"`
        GoModule  string `yaml:"go_module"`
    } `yaml:"generation"`
}
```

**Acceptance criteria:**
- [ ] `sigil init` creates `.sigil/` directory structure
- [ ] `sigil.yaml` has sensible defaults
- [ ] `default.yaml` theme has all token categories from `docs/07_theme-system.md`
- [ ] Error if `.sigil/` already exists (with `--force` to overwrite)
- [ ] `--name` flag sets project name in sigil.yaml

---

### TASK-S1-002 | Priority A | Implement `sigil new page` command

**Context:** Creates a starter page YAML from a template. The template includes
the basic structure with placeholders that the developer fills in.

**What to do:**

Create `internal/cli/new.go` with a `new` parent command and `new page` subcommand:

```go
func NewNewCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "new",
        Short: "Create new Sigil configs",
    }
    cmd.AddCommand(NewNewPageCmd())
    cmd.AddCommand(NewNewDataSourceCmd())
    cmd.AddCommand(NewNewThemeCmd())
    return cmd
}
```

`sigil new page sprint-dashboard --title "Sprint Dashboard" --overlay page --datasource Sprint`:
- Creates `.sigil/pages/sprint-dashboard.yaml`
- Includes boilerplate layout with rows, heading, and a placeholder component
- If `--datasource` is provided, adds a datasource declaration and a data-table stub
- Print: "Created .sigil/pages/sprint-dashboard.yaml"

Template generation logic in `internal/config/templates.go`:
```go
func NewPageTemplate(id, title, overlay string, datasource string) *Page {
    page := &Page{
        Sigil:   "1.0",
        Kind:    "page",
        ID:      id,
        Title:   title,
        Overlay: overlay,
        Layout: Component{
            Type: "rows",
            Props: map[string]interface{}{"gap": 4, "padding": 6},
            Children: []Component{
                {ID: "page-title", Type: "heading", Props: map[string]interface{}{"level": 2, "text": title}},
            },
        },
        Shortcuts: []Shortcut{
            {Key: "Escape", Action: Action{Type: "close"}},
        },
    }
    if datasource != "" {
        page.DataSources = []DataSourceRef{{Alias: datasource, Capabilities: []string{"search", "filter", "sort", "paginate"}}}
        page.Layout.Children = append(page.Layout.Children, Component{
            ID:   "table-" + strings.ToLower(datasource),
            Type: "data-table",
            Props: map[string]interface{}{
                "datasource": datasource,
                "columns":    []interface{}{},
                "pagination": map[string]interface{}{"enabled": true, "pageSize": 25},
            },
        })
    }
    return page
}
```

**Flags:** `--title`, `--overlay`, `--module`, `--layout`, `--datasource`

**Acceptance criteria:**
- [ ] `sigil new page my-page` creates `.sigil/pages/my-page.yaml`
- [ ] Generated YAML is valid (passes `sigil validate`)
- [ ] `--datasource Sprint` adds datasource ref + data-table stub
- [ ] `--title` sets the title; defaults to formatted ID if omitted
- [ ] Error if file already exists (unless `--force`)
- [ ] Output written as clean YAML (not Go struct dump)

---

### TASK-S1-003 | Priority A | Implement `sigil new datasource` command

**Context:** Creates a datasource manifest YAML. See `docs/06_datasource-model.md`.

**What to do:**

Create `internal/cli/new_datasource.go` or add to `new.go`:

`sigil new datasource Sprint --capabilities search,filter,sort,paginate,create`:
- Creates `.sigil/datasources/Sprint.yaml`
- Includes basic structure with alias, capabilities, empty fields array, default endpoints

Template:
```yaml
alias: Sprint
description: ""

capabilities:
  - search
  - filter
  - sort
  - paginate
  - create

fields:
  - name: id
    type: integer
    primary: true
    hidden: true
  # Add your fields here

relations: []

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

**Acceptance criteria:**
- [ ] `sigil new datasource Sprint` creates `.sigil/datasources/Sprint.yaml`
- [ ] Capabilities from `--capabilities` flag
- [ ] Endpoints auto-generated from alias (pluralized, lowercased)
- [ ] Includes `id` field by default

---

### TASK-S1-004 | Priority A | Implement `sigil new theme` command

Creates a theme YAML. See `docs/07_theme-system.md`.

`sigil new theme brand-dark --extends default`:
- Creates `.sigil/themes/brand-dark.yaml`
- If `--extends`, sets the extends field
- Includes empty token overrides section

**Acceptance criteria:**
- [ ] `sigil new theme my-theme` creates `.sigil/themes/my-theme.yaml`
- [ ] `--extends default` sets inheritance
- [ ] `--preset dark` starts with dark preset tokens

---

### TASK-S1-005 | Priority A | Implement `sigil list` command

**Context:** Lists all Sigil configs in the project. See `docs/04_cli-reference.md`.

**What to do:**

Create `internal/cli/list.go`:

```
sigil list pages       → scan .sigil/pages/*.yaml, print table
sigil list datasources → scan .sigil/datasources/*.yaml
sigil list themes      → scan .sigil/themes/*.yaml
sigil list components  → print from component registry
```

Output format:
```
PAGES (2 found)
  sprint-dashboard  Sprint Dashboard   page    project.sprints
  login             Login              modal   auth
```

**Acceptance criteria:**
- [ ] `sigil list pages` shows all pages with id, title, overlay, module
- [ ] `sigil list datasources` shows all datasources with alias, capabilities count
- [ ] `sigil list themes` shows all themes
- [ ] `sigil list components` shows registered component types grouped by category
- [ ] Clean tabular output

---

### TASK-S1-006 | Priority B | YAML serialization (write clean YAML)

**Context:** When `sigil new` creates configs, the output needs to be clean, readable
YAML — not a raw Go struct dump. The yaml.v3 encoder needs configuration.

**What to do:**

Create `internal/config/writer.go`:

```go
func WriteFile(path string, page *Page) error {
    data, err := MarshalYAML(page)
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

func MarshalYAML(v interface{}) ([]byte, error) {
    var buf bytes.Buffer
    enc := yaml.NewEncoder(&buf)
    enc.SetIndent(2)
    if err := enc.Encode(v); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
```

Ensure:
- 2-space indentation
- Empty/nil fields omitted (via `omitempty` tags)
- Flow style for simple inline objects (like `{ value: x, label: Y }`)
- Comments preserved if possible (stretch)

**Acceptance criteria:**
- [ ] Generated YAML is human-readable with 2-space indent
- [ ] Empty fields omitted
- [ ] Round-trip: parse → marshal → parse produces identical struct
- [ ] Generated files pass `sigil validate`
