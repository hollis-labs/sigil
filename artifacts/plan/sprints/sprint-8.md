---
type: sprint-guide
sprint: 8
title: "Config Diffing + Migration + Polish"
status: todo
created_at: 2026-03-02
estimated_days: 3
depends_on: [6, 7]
---

# Sprint 8 — Config Diffing + Migration + Polish

## Objective

Add config diffing (detect changes between versions), migration tooling (upgrade
old configs to new schema versions), and overall polish. This is the "production
readiness" sprint.

## Tasks

### TASK-S8-001 | Priority A | Config diff command

**Context:** When configs change, developers need to see what changed at the
semantic level (not just text diff).

**What to do:**

Create `sigil diff <path-a> <path-b>`:
- Parse both configs
- Compare semantically: added/removed/modified components, changed props, new datasources
- Output colored diff with clear labels
- Support `sigil diff --from git:HEAD` to compare against last commit

**Acceptance criteria:**
- [ ] `sigil diff a.yaml b.yaml` shows semantic changes
- [ ] Added, removed, modified components clearly labeled
- [ ] Prop changes shown with old → new values

---

### TASK-S8-002 | Priority A | Schema migration

**Context:** As the config spec evolves, existing configs may need updating.

**What to do:**

Create `sigil migrate`:
- Detect config schema version (from `sigil:` field)
- Apply migration transforms for each version step
- Support `--dry-run` to preview changes
- For now: just the framework (version detection + transform pipeline)
- First migration: ensure all components have IDs (auto-generate from type + index)

**Acceptance criteria:**
- [ ] `sigil migrate` detects current version
- [ ] Auto-generates missing component IDs
- [ ] `--dry-run` shows what would change

---

### TASK-S8-003 | Priority B | JSON Schema generation

**Context:** IDE support (autocomplete, validation) requires JSON Schema files.

**What to do:**

Create `sigil schema export`:
- Generate `page.schema.json` from Go types + component registry
- Include component prop definitions, enum values, required fields
- Output to `schemas/` directory
- Schema should be usable with VS Code YAML extension for autocomplete

**Acceptance criteria:**
- [ ] JSON Schema generated from component registry
- [ ] Schema validates example configs correctly
- [ ] Works with VS Code YAML extension (manual test)

---

### TASK-S8-004 | Priority B | CI/release polish

**What to do:**

- `Makefile` with build, test, install targets
- Version embedding via `-ldflags`
- `sigil doctor` command: check Go version, templ installed, .sigil dir valid
- README.md with quick start, examples, architecture overview

**Acceptance criteria:**
- [ ] `make build` produces binary with correct version
- [ ] `make test` runs all tests
- [ ] `sigil doctor` checks environment
