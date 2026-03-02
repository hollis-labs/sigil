---
type: sprint-guide
sprint: 5
title: "Preview + Polish + Dogfood"
status: todo
created_at: 2026-03-01
estimated_days: 3
depends_on: [3, 4]
---

# Sprint 5 — Preview + Polish + Dogfood

## Objective

Add the `sigil preview` command (generate standalone HTML for quick visualization),
polish the CLI UX, and dogfood by creating Sigil configs for Sigil's own future GUI.

After this sprint, the MVP is complete: CLI + MCP + code generation + preview.

## Tasks

### TASK-S5-001 | Priority A | Implement `sigil preview` command

**Context:** Developers need to see what a config looks like without building the
full project. `sigil preview` generates a self-contained HTML file with embedded
Tailwind (CDN) and mock data, then opens it in the browser.

**What to do:**

Create `internal/renderer/preview/preview.go`:
- Generate standalone HTML using Tailwind CDN
- Mock data for datasources (generate from field definitions)
- Inline CSS from theme tokens
- Static — no HTMX or interactivity needed for preview
- Open in browser: `exec.Command("open", htmlPath)` on macOS

`sigil preview sprint-dashboard`:
- Parse the page config
- Generate HTML with Tailwind CDN classes
- Write to temp file
- Open in default browser
- Optional: `--port 3210` starts a local HTTP server instead of temp file

**Acceptance criteria:**
- [ ] `sigil preview my-page` opens browser with rendered page
- [ ] Theme tokens applied as inline styles or CSS vars
- [ ] DataSource tables show mock data (5-10 rows)
- [ ] Layout (rows, columns, grid) renders correctly
- [ ] Components render with appropriate styling

---

### TASK-S5-002 | Priority A | Implement `sigil export` and `sigil import`

**Context:** Configs should be shareable via Nanite or as JSON files.

**What to do:**

- `sigil export --format json --output configs/` — export all pages as JSON
- `sigil export --format nanite` — push to Nanite with tags (`sigil/page/<id>`)
- `sigil import --from json --input configs/dashboard.json` — import from JSON
- `sigil import --from nanite --tag "sigil/page/*"` — pull from Nanite

Nanite integration uses the `nanite` CLI (if available) or a Go client library.
For MVP, just JSON export/import. Nanite integration is stretch.

**Acceptance criteria:**
- [ ] `sigil export --format json` produces valid JSON configs
- [ ] `sigil import --from json` imports and validates configs
- [ ] Round-trip: export → import → export produces identical files

---

### TASK-S5-003 | Priority B | CLI UX polish

**Context:** Clean up output formatting, add color, improve error messages, add
`--help` descriptions.

**What to do:**

- Add color to output (green for success, red for errors, yellow for warnings)
- Add `--no-color` global flag
- Improve error messages with suggestions (e.g., "unknown type 'data-grid', did you mean 'data-table'?")
- Add `sigil version` command
- Add `sigil help` improvements (examples in help text)
- Add `--verbose` global flag for debug output

**Acceptance criteria:**
- [ ] Colored output for errors/warnings/success
- [ ] `--no-color` disables color
- [ ] Error suggestions for common mistakes
- [ ] `sigil version` prints version
- [ ] All commands have examples in `--help`

---

### TASK-S5-004 | Priority B | Dogfood: create Sigil configs for Sigil's own GUI

**Context:** When we build the GUI for Sigil itself (future sprint), we'll use Sigil
configs to define it. Start now by creating the page configs for the eventual GUI.

**What to do:**

Create `.sigil/pages/` configs for the Sigil GUI:

1. `sigil-pages.yaml` — page listing (data-table of all pages with status)
2. `sigil-page-editor.yaml` — YAML editor for a page config
3. `sigil-components.yaml` — component browser (categories + schemas)
4. `sigil-preview.yaml` — preview pane for a selected page
5. `sigil-datasources.yaml` — datasource listing and editor
6. `sigil-themes.yaml` — theme editor with live preview

These configs are real Sigil configs — they should pass `sigil validate`.
They serve as both documentation and future implementation targets.

**Acceptance criteria:**
- [ ] 6 page configs created in `.sigil/pages/`
- [ ] All pass `sigil validate`
- [ ] Configs use a variety of component types (data-table, forms, tabs, modals)
- [ ] Configs reference Sigil-specific datasources (SigilPage, SigilComponent, etc.)
