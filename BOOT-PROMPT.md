# Sigil — Boot Prompt

> Last updated: 2026-04-03

## Priority: Stack Explorer Frontend

Build the Stack Explorer frontend using Sigil. This is a flagship project — push what's possible.

**Full spec:** `docs/projects/stack-explorer-frontend.md` — read this first.

**Summary:** Stack Explorer is an AI agent tooling research platform (112 repos, 18 dimensions, 9 lenses, 47 scored repos). Build a modern analytics dashboard with full CRUD over all entities, interactive reports, lens-switchable scorecards, gap analysis charts, and pattern catalogs.

**Start with:**
1. Read the full spec at `docs/projects/stack-explorer-frontend.md`
2. Read the product vision at `/Users/chrispian/Projects-apps/stack-explorer/docs/product-vision.md`
3. Define datasource manifests (11 entities)
4. Create the dark theme
5. Build pages: Dashboard → Repos → Scorecards → Gap Analysis → Reports → Patterns → Findings → Dimensions → Settings

**Target:** `react-shadcn` renderer. Generate into the Stack Explorer project.

**Design bar:** This should look like Linear meets Datadog. Information-dense, keyboard-first, dark theme with blue accents and score heatmaps. Not a basic admin panel.

## Other Work

Previous work (design polish, migration) is documented below for context but is lower priority than the Stack Explorer frontend.

---

## Previous Context

### Active Agent

`sigil-frontend` — React/Next.js frontend development for Sigil demo and generated output.

### Design Language

- **Background:** zinc-950 (#09090b), surface zinc-900, borders zinc-800
- **Accent:** blue-500 for primary actions
- **Status colors:** green=success, amber=warning, red=danger
- **Component library:** shadcn/ui v4 + Tailwind v4
- **Target:** Next.js 16 app router

### Key Context

- **Stack:** Go CLI + React/shadcn generated output
- **Renderers:** go-templ, react-shadcn, preview (HTML)
- **57 built-in components** — see `internal/components/builtin/`
- **MCP server** for AI-assisted UI generation
- **Demo app** at `demo/` — Next.js 16 + Tailwind v4

### Files to Read

1. `CLAUDE.md` — Project conventions
2. `docs/02_config-spec.md` — YAML spec format
3. `docs/06_datasource-model.md` — DataSource contract
4. `docs/07_theme-system.md` — Theme tokens
5. `examples/sprint-dashboard.yaml` — Reference page example
