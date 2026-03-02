---
type: sprint-guide
sprint: 0
title: "Project Scaffold + Core Types + Config Parser"
status: todo
created_at: 2026-03-01
estimated_days: 3
depends_on: []
---

# Sprint 0 — Project Scaffold + Core Types + Config Parser

## Objective

Set up the Go project, define all core types, build the YAML config parser, and
implement schema validation. After this sprint, `sigil validate` works on real YAML
configs and reports errors.

## Context for the Orchestrator

This is a **greenfield Go project**. There is no existing code — just docs and sprint
guides. Your first job is to initialize the Go module, set up the project structure,
and implement the foundational types and parsing.

**Read these docs before starting:**
- `docs/01_architecture.md` — system overview and Go module structure
- `docs/02_config-spec.md` — the YAML config format you're parsing
- `docs/03_component-model.md` — component type system

**Tech stack decisions (already made):**
- Go 1.22+ with standard library
- Cobra for CLI (`github.com/spf13/cobra`)
- `gopkg.in/yaml.v3` for YAML parsing
- No external schema validation library — implement validation in Go against component schemas

## Deliverables

1. Go module initialized (`go.mod`, `go.sum`)
2. Core types in `internal/config/types.go`
3. YAML parser in `internal/config/parser.go`
4. Default values in `internal/config/defaults.go`
5. Validation in `internal/config/validator.go`
6. Cobra CLI skeleton with `root` and `validate` commands
7. Example configs in `examples/`
8. `go build ./cmd/sigil` produces a working binary

## Verification

```bash
go build ./cmd/sigil                                    # Compiles
go vet ./...                                            # No vet issues
go test ./...                                           # Tests pass
./sigil validate examples/sprint-dashboard.yaml          # Valid
./sigil validate examples/invalid-example.yaml           # Reports errors
```

---

## Tasks

### TASK-S0-001 | Priority A | Initialize Go module and project structure

**Context:** Starting from empty directories. Need go.mod, main.go entry point, and
the package structure defined in `docs/01_architecture.md`.

**What to do:**

1. Run `go mod init github.com/chrispian/sigil` (or appropriate module path)
2. Create `cmd/sigil/main.go` — minimal main that calls cobra root command
3. Create package directories with placeholder files:
   - `internal/config/` — `doc.go`
   - `internal/cli/` — `doc.go`
   - `internal/renderer/` — `doc.go`
   - `internal/components/` — `doc.go`
   - `internal/theme/` — `doc.go`
   - `internal/datasource/` — `doc.go`
4. Add dependencies: `go get github.com/spf13/cobra gopkg.in/yaml.v3`
5. Verify: `go build ./cmd/sigil` compiles

**Acceptance criteria:**
- [ ] `go.mod` exists with module path and Go version
- [ ] `go build ./cmd/sigil` produces a binary
- [ ] All package directories exist with at least a `doc.go`
- [ ] Cobra and yaml.v3 in `go.mod` as dependencies

---

### TASK-S0-002 | Priority A | Define core config types

**Context:** These types are the Go representation of the Sigil YAML config. Every
other package depends on them. Match the spec in `docs/02_config-spec.md` exactly.

**What to do:**

Create `internal/config/types.go`:

```go
package config

// Page represents a complete Sigil page config.
type Page struct {
    Sigil       string            `yaml:"sigil" json:"sigil"`
    Kind        string            `yaml:"kind" json:"kind"`
    ID          string            `yaml:"id" json:"id"`
    Title       string            `yaml:"title" json:"title"`
    Description string            `yaml:"description,omitempty" json:"description,omitempty"`
    Overlay     string            `yaml:"overlay" json:"overlay"`
    Module      string            `yaml:"module,omitempty" json:"module,omitempty"`
    Theme       *ThemeRef         `yaml:"theme,omitempty" json:"theme,omitempty"`
    DataSources []DataSourceRef   `yaml:"datasources,omitempty" json:"datasources,omitempty"`
    Layout      Component         `yaml:"layout" json:"layout"`
    Shortcuts   []Shortcut        `yaml:"shortcuts,omitempty" json:"shortcuts,omitempty"`
    Meta        *PageMeta         `yaml:"meta,omitempty" json:"meta,omitempty"`
}

// Component represents a UI component in the config tree.
type Component struct {
    ID        string                 `yaml:"id,omitempty" json:"id,omitempty"`
    Type      string                 `yaml:"type" json:"type"`
    Props     map[string]interface{} `yaml:"props,omitempty" json:"props,omitempty"`
    Actions   map[string]Action      `yaml:"actions,omitempty" json:"actions,omitempty"`
    Children  []Component            `yaml:"children,omitempty" json:"children,omitempty"`
    Shortcuts []Shortcut             `yaml:"shortcuts,omitempty" json:"shortcuts,omitempty"`
    Condition *Condition             `yaml:"condition,omitempty" json:"condition,omitempty"`
}

// Action defines an event handler.
type Action struct {
    Type      string                 `yaml:"type" json:"type"`
    Page      string                 `yaml:"page,omitempty" json:"page,omitempty"`
    URL       string                 `yaml:"url,omitempty" json:"url,omitempty"`
    Method    string                 `yaml:"method,omitempty" json:"method,omitempty"`
    Title     string                 `yaml:"title,omitempty" json:"title,omitempty"`
    Size      string                 `yaml:"size,omitempty" json:"size,omitempty"`
    Fields    []FormField            `yaml:"fields,omitempty" json:"fields,omitempty"`
    Submit    *SubmitConfig          `yaml:"submit,omitempty" json:"submit,omitempty"`
    Refresh   string                 `yaml:"refresh,omitempty" json:"refresh,omitempty"`
    Params    map[string]string      `yaml:"params,omitempty" json:"params,omitempty"`
    Event     string                 `yaml:"event,omitempty" json:"event,omitempty"`
    Payload   map[string]interface{} `yaml:"payload,omitempty" json:"payload,omitempty"`
    Message   string                 `yaml:"message,omitempty" json:"message,omitempty"`
    OnConfirm *Action                `yaml:"onConfirm,omitempty" json:"onConfirm,omitempty"`
    OnSuccess *Action                `yaml:"onSuccess,omitempty" json:"onSuccess,omitempty"`
    Target    string                 `yaml:"target,omitempty" json:"target,omitempty"`
}

// FormField defines a field in a modal/form.
type FormField struct {
    Name        string        `yaml:"name" json:"name"`
    Label       string        `yaml:"label" json:"label"`
    Type        string        `yaml:"type" json:"type"`
    Required    bool          `yaml:"required,omitempty" json:"required,omitempty"`
    Placeholder string        `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
    Default     interface{}   `yaml:"default,omitempty" json:"default,omitempty"`
    Options     []SelectOption `yaml:"options,omitempty" json:"options,omitempty"`
    Validation  *Validation   `yaml:"validation,omitempty" json:"validation,omitempty"`
    Rows        int           `yaml:"rows,omitempty" json:"rows,omitempty"`
    MaxLength   int           `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
    Min         interface{}   `yaml:"min,omitempty" json:"min,omitempty"`
    Max         interface{}   `yaml:"max,omitempty" json:"max,omitempty"`
}

type SelectOption struct {
    Value string `yaml:"value" json:"value"`
    Label string `yaml:"label" json:"label"`
}

type Validation struct {
    Pattern string `yaml:"pattern,omitempty" json:"pattern,omitempty"`
    Message string `yaml:"message,omitempty" json:"message,omitempty"`
}

type SubmitConfig struct {
    DataSource string `yaml:"datasource" json:"datasource"`
    Method     string `yaml:"method" json:"method"`
}

// Shortcut defines a keyboard shortcut binding.
type Shortcut struct {
    Key         string `yaml:"key" json:"key"`
    Action      Action `yaml:"action" json:"action"`
    When        string `yaml:"when,omitempty" json:"when,omitempty"`
    Description string `yaml:"description,omitempty" json:"description,omitempty"`
    Global      bool   `yaml:"global,omitempty" json:"global,omitempty"`
}

// Condition for conditional rendering.
type Condition struct {
    Field string      `yaml:"field" json:"field"`
    Op    string      `yaml:"op" json:"op"`
    Value interface{} `yaml:"value" json:"value"`
}

// DataSourceRef declares a datasource dependency in a page.
type DataSourceRef struct {
    Alias        string   `yaml:"alias" json:"alias"`
    Capabilities []string `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`
    Params       map[string]interface{} `yaml:"params,omitempty" json:"params,omitempty"`
}

// ThemeRef references a theme with optional overrides.
type ThemeRef struct {
    Extends string                 `yaml:"extends,omitempty" json:"extends,omitempty"`
    Tokens  map[string]interface{} `yaml:"tokens,omitempty" json:"tokens,omitempty"`
}

// PageMeta contains page metadata.
type PageMeta struct {
    Guards  []string `yaml:"guards,omitempty" json:"guards,omitempty"`
    Tags    []string `yaml:"tags,omitempty" json:"tags,omitempty"`
    Version int      `yaml:"version,omitempty" json:"version,omitempty"`
}
```

**Acceptance criteria:**
- [ ] All types from `docs/02_config-spec.md` represented
- [ ] YAML and JSON struct tags on all fields
- [ ] Optional fields use `omitempty`
- [ ] `go vet ./internal/config/...` passes
- [ ] Types are exported (uppercase)

---

### TASK-S0-003 | Priority A | Implement YAML config parser

**Context:** The parser reads `.yaml` files and returns typed `Page` structs. It
needs to handle both single-page files and multi-document YAML.

**What to do:**

Create `internal/config/parser.go`:

```go
package config

import (
    "fmt"
    "os"
    "gopkg.in/yaml.v3"
)

// ParseFile reads a YAML file and returns a Page config.
func ParseFile(path string) (*Page, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading %s: %w", path, err)
    }
    return Parse(data)
}

// Parse parses YAML bytes into a Page config.
func Parse(data []byte) (*Page, error) {
    var page Page
    if err := yaml.Unmarshal(data, &page); err != nil {
        return nil, fmt.Errorf("parsing YAML: %w", err)
    }
    ApplyDefaults(&page)
    return &page, nil
}
```

Also create `internal/config/defaults.go`:

```go
package config

// ApplyDefaults fills in default values for optional fields.
func ApplyDefaults(page *Page) {
    if page.Sigil == "" {
        page.Sigil = "1.0"
    }
    if page.Kind == "" {
        page.Kind = "page"
    }
    if page.Overlay == "" {
        page.Overlay = "page"
    }
    applyComponentDefaults(&page.Layout)
}

func applyComponentDefaults(c *Component) {
    // Apply defaults to component props based on type
    // Recurse into children
    for i := range c.Children {
        applyComponentDefaults(&c.Children[i])
    }
}
```

Write tests in `internal/config/parser_test.go`:
- Test: parse a valid YAML string → returns correct Page struct
- Test: parse invalid YAML → returns error
- Test: defaults are applied for omitted fields
- Test: nested components are parsed correctly
- Test: actions and shortcuts are parsed

**Acceptance criteria:**
- [ ] `ParseFile()` reads and parses YAML files
- [ ] `Parse()` parses YAML bytes
- [ ] Defaults applied for omitted optional fields
- [ ] Nested component trees parsed correctly
- [ ] Tests pass: `go test ./internal/config/...`
- [ ] Error messages include file path and context

---

### TASK-S0-004 | Priority A | Implement config validator

**Context:** Validation checks that a parsed config is structurally correct — required
fields present, component types recognized, action types valid, no duplicate IDs,
shortcuts have valid key expressions.

**What to do:**

Create `internal/config/validator.go`:

```go
package config

// ValidationResult contains all validation findings.
type ValidationResult struct {
    Valid    bool
    Errors   []ValidationError
    Warnings []ValidationWarning
}

type ValidationError struct {
    Path    string // Config path: "layout.children[0].props.columns"
    Message string
}

type ValidationWarning struct {
    Path    string
    Message string
}

// Validate checks a Page config for structural correctness.
func Validate(page *Page, registry ComponentRegistry) *ValidationResult {
    result := &ValidationResult{Valid: true}

    // 1. Required top-level fields
    validateRequired(page, result)

    // 2. Valid overlay type
    validateOverlay(page.Overlay, result)

    // 3. Walk component tree
    seen := make(map[string]bool)
    validateComponent(&page.Layout, "layout", seen, registry, result)

    // 4. Validate shortcuts
    for i, s := range page.Shortcuts {
        validateShortcut(s, fmt.Sprintf("shortcuts[%d]", i), result)
    }

    result.Valid = len(result.Errors) == 0
    return result
}
```

**Validation checks to implement:**

| # | Check | Severity |
|---|---|---|
| 1 | `sigil`, `id`, `title`, `overlay` present | Error |
| 2 | `overlay` is valid enum value | Error |
| 3 | `layout.type` is present | Error |
| 4 | Component `type` exists in registry (if registry provided) | Error |
| 5 | Component `id` is unique within the page | Error |
| 6 | Action `type` is valid enum | Error |
| 7 | Modal action has `title` and `fields` | Error |
| 8 | Navigate action has `page` or `url` | Error |
| 9 | HTTP action has `url` and `method` | Error |
| 10 | DataSource refs have `alias` | Error |
| 11 | Shortcut `key` is non-empty | Error |
| 12 | Form field has `name`, `label`, `type` | Error |
| 13 | Select field has `options` | Warning |
| 14 | Components without `id` | Warning |

Write tests in `internal/config/validator_test.go`.

**Acceptance criteria:**
- [ ] `Validate()` returns structured results with errors and warnings
- [ ] All 14 checks implemented
- [ ] Error paths indicate where in the config the issue is
- [ ] Tests cover valid config, missing fields, bad types, duplicate IDs
- [ ] `go test ./internal/config/...` passes

---

### TASK-S0-005 | Priority A | Implement Cobra CLI skeleton with validate command

**Context:** The CLI is the primary interface. Start with just `root` and `validate`
commands. More commands added in Sprint 1.

**What to do:**

Create `internal/cli/root.go`:
```go
package cli

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "sigil",
        Short: "UI configuration and code generation tool",
        Long:  "Sigil turns declarative YAML configs into framework-specific UI code.",
    }
    cmd.AddCommand(NewValidateCmd())
    return cmd
}
```

Create `internal/cli/validate.go`:
```go
func NewValidateCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "validate [path]",
        Short: "Validate Sigil configs",
        Long:  "Validate Sigil YAML configs against component schemas.",
        Args:  cobra.MaximumNArgs(1),
        RunE:  runValidate,
    }
    cmd.Flags().Bool("strict", false, "Fail on warnings too")
    return cmd
}

func runValidate(cmd *cobra.Command, args []string) error {
    // If no args, validate all .sigil/pages/*.yaml
    // If arg, validate that specific file
    // Parse → Validate → Print results
    // Exit code 0 = valid, 1 = errors
}
```

Update `cmd/sigil/main.go`:
```go
package main

import (
    "os"
    "github.com/chrispian/sigil/internal/cli"
)

func main() {
    if err := cli.NewRootCmd().Execute(); err != nil {
        os.Exit(1)
    }
}
```

**Output format for validate:**
```
✓ examples/sprint-dashboard.yaml — valid
✗ examples/invalid.yaml — 2 errors, 1 warning
  ERROR [layout.children[0]]: unknown component type "data-grid"
  ERROR [layout.children[0].props]: required prop "columns" missing for "data-table"
  WARN  [datasources]: datasource "User" referenced but not declared
```

**Acceptance criteria:**
- [ ] `sigil` prints help
- [ ] `sigil validate examples/sprint-dashboard.yaml` validates and prints result
- [ ] `sigil validate` (no args) scans `.sigil/pages/`
- [ ] `--strict` flag causes warnings to be treated as errors
- [ ] Exit code 0 for valid, 1 for errors
- [ ] `go build ./cmd/sigil` produces working binary

---

### TASK-S0-006 | Priority A | Create example configs

**Context:** Examples serve as test fixtures and documentation. They should cover
common patterns and exercise the parser/validator.

**What to do:**

Create `examples/sprint-dashboard.yaml` — the full example from `docs/02_config-spec.md`
(the sprint dashboard with data-table, search, create modal, keyboard shortcuts).

Create `examples/login-form.yaml`:
```yaml
sigil: "1.0"
kind: page
id: login
title: Login
overlay: modal

layout:
  type: rows
  props:
    gap: 6
    padding: 8
    align: center
  children:
    - id: login-heading
      type: heading
      props:
        level: 2
        text: Sign In
    - id: login-form
      type: form
      props:
        fields:
          - { name: email, label: Email, type: email, required: true, placeholder: "you@example.com" }
          - { name: password, label: Password, type: password, required: true }
          - { name: remember, label: "Remember me", type: checkbox }
        submit:
          datasource: Auth
          method: POST
          label: Sign In
      actions:
        submit:
          type: http
          url: "/api/auth/login"
          method: POST
          onSuccess:
            type: navigate
            page: dashboard

shortcuts:
  - key: Escape
    action: { type: close }
  - key: Enter
    action: { type: emit, event: submit-form }
    when: "!input-focused"
```

Create `examples/invalid-example.yaml` — intentionally invalid config for testing:
```yaml
sigil: "1.0"
kind: page
# Missing: id, title
overlay: invalid-type

layout:
  type: unknown-layout
  children:
    - type: data-table
      # Missing: id, props.datasource, props.columns
    - id: duplicate-id
      type: button
    - id: duplicate-id
      type: button
```

**Acceptance criteria:**
- [ ] `examples/sprint-dashboard.yaml` — valid, exercises data-table + modal + shortcuts
- [ ] `examples/login-form.yaml` — valid, exercises form + auth flow
- [ ] `examples/invalid-example.yaml` — has known errors for validator testing
- [ ] All valid examples parse without error
- [ ] Invalid example triggers multiple validation errors

---

### TASK-S0-007 | Priority B | Component registry foundation

**Context:** The component registry maps type strings (e.g., "data-table") to schemas
that define what props/actions/slots each type supports. For Sprint 0, we need a basic
registry that can be populated and queried. The actual schemas are added in Sprint 2.

**What to do:**

Create `internal/components/registry.go`:

```go
package components

// Schema defines a component type's props, actions, and slots.
type Schema struct {
    Type        string
    Category    string // primitives, layouts, data, forms, navigation, composites
    Description string
    Props       map[string]PropDef
    Actions     map[string]ActionDef
    Slots       map[string]SlotDef
    Shortcuts   []ShortcutDef
}

type PropDef struct {
    Type        string      // string, integer, boolean, array, object, enum
    Required    bool
    Default     interface{}
    Description string
    Enum        []string    // For enum type
}

type ActionDef struct {
    Description string
}

type SlotDef struct {
    Description string
    Accepts     []string // Accepted component types
}

type ShortcutDef struct {
    Key         string
    Description string
    When        string
}

// Registry holds all known component schemas.
type Registry struct {
    schemas map[string]*Schema
}

func NewRegistry() *Registry {
    return &Registry{schemas: make(map[string]*Schema)}
}

func (r *Registry) Register(schema *Schema) {
    r.schemas[schema.Type] = schema
}

func (r *Registry) Get(componentType string) (*Schema, bool) {
    s, ok := r.schemas[componentType]
    return s, ok
}

func (r *Registry) Has(componentType string) bool {
    _, ok := r.schemas[componentType]
    return ok
}

func (r *Registry) Types() []string { ... }
func (r *Registry) ByCategory(category string) []*Schema { ... }
```

Register a minimal set of types for Sprint 0 (just so the validator can check types):
- All layout types: rows, columns, grid, card, tabs, tab, split, sidebar
- Key primitives: heading, text, button, icon-button, input, select, badge
- Key data: data-table, search-bar, form

Full schemas with prop definitions come in Sprint 2.

**Acceptance criteria:**
- [ ] `Registry` type with Register/Get/Has/Types methods
- [ ] `Schema` type with props/actions/slots definitions
- [ ] ~20 component types registered (type name + category, minimal props)
- [ ] Validator can check `registry.Has(componentType)` during validation
- [ ] `go test ./internal/components/...` passes
