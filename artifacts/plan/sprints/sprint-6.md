---
type: sprint-guide
sprint: 6
title: "Live Dev Server + Watch Mode"
status: todo
created_at: 2026-03-02
estimated_days: 3
depends_on: [5]
---

# Sprint 6 — Live Dev Server + Watch Mode

## Objective

Replace static HTML preview with a live development server. When a developer runs
`sigil serve`, it starts an HTTP server that renders pages in real-time, watches
`.sigil/` for changes, and reloads the browser automatically via SSE.

After this sprint, developers get instant visual feedback while editing configs.

## Tasks

### TASK-S6-001 | Priority A | HTTP dev server

**Context:** `sigil preview` generates a static file. Developers need a live server
that serves pages at `http://localhost:3210/pages/<id>` with navigation between pages.

**What to do:**

Create `internal/server/server.go`:
- HTTP server on configurable port (default: 3210)
- Route: `GET /` — index page listing all pages with links
- Route: `GET /pages/{id}` — render page preview HTML on-the-fly
- Route: `GET /assets/theme.css` — serve generated theme CSS
- Route: `GET /events` — SSE endpoint for live reload
- Use the existing `preview.GenerateHTML()` to render pages
- Inject a small `<script>` into each page that connects to `/events` for reload

Create `internal/cli/serve.go`:
- `sigil serve` command with `--port`, `--sigil-dir` flags
- Print URL on startup
- Graceful shutdown on SIGINT

**Acceptance criteria:**
- [ ] `sigil serve` starts HTTP server
- [ ] `localhost:3210/` shows page index
- [ ] `localhost:3210/pages/sigil-pages` renders live preview
- [ ] Navigation between pages works via links
- [ ] Theme CSS is served and applied

---

### TASK-S6-002 | Priority A | File watcher + auto-reload

**Context:** When a developer edits a page YAML, the browser should update automatically.

**What to do:**

Add file watching to the dev server:
- Watch `.sigil/pages/`, `.sigil/themes/`, `.sigil/datasources/` for changes
- On change: invalidate cached HTML, send SSE event to connected browsers
- Browser-side: `EventSource('/events')` listener triggers `location.reload()`
- Debounce: 100ms to batch rapid saves
- Use `fsnotify` package (or poll-based fallback for simplicity in MVP)

For MVP, poll-based (check mtime every 500ms) is simpler and avoids the fsnotify
dependency. Can upgrade later.

**Acceptance criteria:**
- [ ] Edit a page YAML → browser reloads within 1 second
- [ ] Edit theme YAML → styles update
- [ ] New page added → appears on index
- [ ] Deleted page → removed from index

---

### TASK-S6-003 | Priority B | Index page with navigation

**Context:** The index page should look polished and provide easy navigation.

**What to do:**

Create a styled index page:
- List all pages with ID, title, overlay type, component count
- Link each page to its live preview
- Show project name from sigil.yaml
- Show validation status (pass/warn/fail) for each page
- Include a "Validate All" button that shows results inline
- Use the project's own theme for styling

**Acceptance criteria:**
- [ ] Index page lists all pages with metadata
- [ ] Each page links to its live preview
- [ ] Validation status shown per page
- [ ] Styled using project theme tokens

---

### TASK-S6-004 | Priority B | Dev server tests

**Context:** The server needs tests for routing, SSE, and reload behavior.

**What to do:**

- Unit tests for route handlers (page render, index, theme CSS)
- Test SSE endpoint sends events
- Integration test: start server, fetch page, verify HTML content
- Test graceful shutdown

**Acceptance criteria:**
- [ ] HTTP handler tests for all routes
- [ ] SSE event delivery test
- [ ] Integration test with real HTTP server
