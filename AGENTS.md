# AGENTS.md — Sigil

## What is this and why

Sigil is a system-agnostic UI configuration and code-generation tool. You define a UI
declaratively in YAML (pages as component trees, themes as design tokens, datasources as
data shapes), Sigil validates it against 49 built-in component schemas, and generates
framework-specific source code for one of several render targets — Go/Templ+HTMX,
React/shadcn+Tailwind, or a static HTML preview. The core value proposition is **one
config language, multiple render targets**: write your pages once, generate for any
target. Codegen is the differentiator — Sigil deliberately is *not* a runtime renderer,
visual editor, or low-code platform; protect that boundary.

## Where to start

- `README.md` — feature overview, quick start, CLI table, render-target outputs.
- `cmd/sigil/` — entry point (single binary, Cobra CLI).
- `internal/cli/` — CLI command implementations.
- `internal/renderer/` — the `Renderer` interface plus concrete targets (gotempl,
  reactshadcn, preview).
- `internal/components/builtin/` — the 49 component schemas (YAML, embedded in binary).
- `docs/` — numbered architecture docs (`01_architecture.md` through
  `08_mcp-integration.md`); start with `01_architecture.md` and `02_config-spec.md`.
- `.sigil/` — a working Sigil project lives here (project config, pages, themes,
  datasources, custom components) and doubles as the dogfooding workspace.
- `.agentrc/config.yaml` — per-project agent definitions (`sigil-backend`,
  `sigil-renderer`, `sigil-frontend`, `sigil-reviewer`).

## Key domain concepts

- **Page config** — a YAML document (`sigil: "1.0"`, `kind: page`) describing a UI as a
  nested component tree with props, actions, and datasource bindings.
- **Component registry** — 49 built-in component types across 6 categories (layout,
  primitive, form, data, overlay, navigation); custom schemas are supported via
  `*.schema.yaml`.
- **Render target** — a pluggable code generator. `go-templ` emits `.templ` + HTMX +
  handler stubs; `react-shadcn` emits `.tsx` + shadcn imports + SWR hooks + TS types;
  `preview` emits static HTML with mock data. A catalog target is in active planning.
- **Theme** — design tokens (colors, radius, etc.) compiled to `--sigil-*` CSS custom
  properties and a Tailwind preset.
- **DataSource** — a declared data shape; generators emit typed React hooks or Go
  handler stubs from it.
- **MCP server** — `sigil mcp serve` exposes Sigil to AI agents over JSON-RPC 2.0 on
  stdio (9 tools, 7 resource types, 2 prompts). This is the primary agent-facing surface.
- **Live dev server** — `sigil serve` watches `.sigil/` and reloads the browser via SSE.

## Common operations

```bash
# Build the binary
make build                       # → bin/sigil

# Run tests and vet
make test
make vet

# Install locally
make install                     # go install → $GOBIN/sigil

# Initialize and author a project
sigil init --name my-app
sigil new page dashboard --title "Dashboard"
sigil validate

# Generate code
sigil generate --target go-templ     --output internal/ui
sigil generate --target react-shadcn --output src/generated

# Preview / live dev
sigil preview dashboard
sigil serve

# Agent integration
sigil mcp serve

# Run the bundled demos (Next.js / SPA)
make demo                        # demo/ on port 3333
make se-demo                     # Stack Explorer demo on port 3334
make cw-demo                     # Clockwork demo (Vite SPA)
```

The `demo*` Make targets regenerate React output into the demo apps; they use
`rsync --ignore-existing` to protect hand-written mock hooks. If hooks get clobbered,
restore the mock versions manually.

## Where to look for more

- **Architecture / specs:** `docs/01_architecture.md` … `docs/08_mcp-integration.md`,
  plus `docs/05_renderer-contract.md` for the renderer plugin contract.
- **ADRs:** `docs/adr/` — `0001-sigil-studio-positioning.md` (Proposed; positions a
  prospective Sigil Studio sibling repo). ADR 0002 (lightweight-studio spike outcome)
  is not yet written.
- **Project briefs / roadmap:** `docs/projects/` — `catalog-renderer-prompt.md`,
  `lightweight-studio-exploration-prompt.md`, `sigil-studio-planning-prompt.md`,
  `stack-explorer-api-prompt.md`, `stack-explorer-frontend.md`.
- **Portfolio knowledge:** `~/dev/agent-os/knowledge/projects/sigil.md` — cross-project
  context, sub-project status, composition points, open gaps.
- **SoT:** `.agent-ops/project.yaml` — machine-readable project source-of-truth.

## Conventions

- Three core dependencies only (Cobra, yaml.v3, pflag) — adding a dependency requires
  strong justification.
- Do not guess when uncertain — stop and ask.
- Sub-agent output stays in the sub-agent; the main context gets one-line confirmations.
- Prefer focused, minimal output; no trailing summaries.
