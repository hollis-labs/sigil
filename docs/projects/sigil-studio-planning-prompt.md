# Sigil Studio — Planning Agent Brief

You are a local planning agent. Your job is to read the context pack and the ADR, respect the decisions already made, and produce **four sub-project specs + a sequencing plan**. You are NOT writing code. You are producing specs that a follow-on implementation agent will execute.

## Preconditions

**Only proceed with this brief if the lightweight spike (see `lightweight-studio-exploration-prompt.md`) has been run and its conclusion is "we still want Studio."** If the spike has not been run, stop and say so. If the spike concluded the lightweight path is sufficient, stop and write ADR 0002 closing out Studio instead.

## Inputs (read in this order)

1. `docs/adr/0001-sigil-studio-positioning.md` — decisions already made. Do not re-open them.
2. `/Users/chrispian/Projects-apps/agent-workspaces/inbox/sigil_studio_agent_pack/` — the original context pack. Read every file.
3. `docs/projects/catalog-renderer-prompt.md` — the catalog plan. Studio's Compose palette consumes this.
4. Sigil docs: `docs/01_architecture.md`, `docs/02_config-spec.md`, `docs/03_component-model.md`, `docs/05_renderer-contract.md`, `docs/07_theme-system.md`, `docs/08_mcp-integration.md`.
5. Sigil code: skim `internal/components/`, `internal/config/`, `internal/renderer/`, `internal/theme/`, `internal/server/`, `internal/mcp/`. Identify what is safely importable as a library vs. what is internal-only.

## Decisions already locked (from ADR 0001)

- Sibling repository. Studio depends on Sigil; Sigil does not depend on Studio.
- Catalog ships first and is not blocked by this work.
- Sigil consumed as Go library (rendering, config) + MCP client (agent ops). No CLI shell-out.
- Command layer is MCP-native from day one — Studio exposes its commands over MCP.
- Compose-mode v1 preview is an iframe to `sigil serve` with Studio managing a file-watcher.
- Round-trip YAML uses `yaml.v3` node API, not `Marshal`.
- Scene mode is the first shippable release (v0). Compose follows.
- Four sub-projects, each with its own spec.

Do not revisit these. If you find a reason to, flag it as an addendum to the ADR — do not silently override.

## Your deliverables

Produce the following, all in `docs/projects/studio/`:

### 1. `00-sequencing.md`
A single sequencing doc covering:
- Dependency graph between the four sub-projects.
- Which sub-project is the v0 release (Scene mode on Wails shell + command layer + MCP surface, no Sigil bridge yet).
- Which sub-projects can run in parallel once v0 ships.
- Risk callouts per sub-project with mitigation notes.
- Rough effort bands (S/M/L/XL), not hour estimates.

### 2. `01-app-shell-spec.md` — Wails shell + documents + command layer + MCP surface
Scope:
- Wails bootstrap, React frontend scaffold, Tailwind + shadcn/ui baseline.
- Local document storage (scene docs + compose docs) with save/load.
- Command layer: every meaningful user or agent action is a named command with a schema. UI wires through commands, not direct state mutations.
- MCP server surface exposing those commands. Studio is driveable by any MCP client.
- Document browser / recent files shell.
- Bundle export/import plumbing (schemas only; Scene mode fills in the artifact content).

Out of scope for this sub-project: tldraw, Sigil bridge, Compose UI.

### 3. `02-scene-mode-spec.md` — tldraw + export + handoff bundle
Scope:
- Embed tldraw in the React frontend.
- Scene document schema (close to tldraw-native per the pack).
- PNG export (required), SVG export (if low cost).
- Handoff bundle export: image + notes + intent + history metadata.
- Handoff bundle *ingest* (agent → Studio). This was an open question in the ADR; resolve it here.
- Image paste/upload for screenshots and reference material.

Out of scope: Compose mode, semantic conversion of scenes.

### 4. `03-sigil-bridge-spec.md` — library consumption + round-trip + preview
Scope:
- Which Sigil packages Studio imports. Recommend changes to Sigil (e.g. moving internal → pkg/) where necessary; enumerate them.
- Sigil config loader/saver using `yaml.v3` node API with formatting and comments preserved.
- Preview pipeline: Studio launches or attaches to `sigil serve`, iframes it, watches the config file for changes.
- Catalog consumption: Compose palette reads either the catalog's static HTML (embedded) or the variant YAMLs directly. Pick one with reasoning.
- Theme/token introspection so Compose's inspector can display and edit tokens.

Out of scope: Compose UI itself (that's sub-project 4). This sub-project provides the API; the next one consumes it.

### 5. `04-compose-mode-spec.md` — palette + inspector + preview pane
Scope:
- Compose-mode layout: palette (left), canvas (center showing the `sigil serve` iframe or a custom view), inspector (right).
- Component palette driven by the catalog.
- Prop inspector — visual props first (using the catalog's visual-prop filter), full props behind an advanced toggle.
- Token/theme inspector.
- Drag/drop or click-to-add from palette into the active Sigil config.
- Save writes back to the config via the bridge's node-API saver.

Out of scope: anything that requires extending Sigil's renderer or schema system beyond what exists when this sub-project starts.

## Format expectations for each spec

- Problem statement (2-3 sentences).
- In scope / out of scope as explicit lists.
- Architecture sketch (packages/modules, interfaces, data flow).
- Open questions with proposed answers.
- Test strategy.
- Dependencies on other sub-projects.
- Deliverable checklist.

Specs should be terse and concrete. Avoid restating the ADR or the context pack — link to them.

## Working rules

- One question per clarification round if you need to ask the human. Do not batch.
- If a decision is underspecified by the ADR or pack, propose the answer in your spec with reasoning. Do not wait for permission on minor choices.
- If you find a reason to contradict the ADR, stop, write up the contradiction with evidence, and ask. Do not proceed.
- Sub-specs should not exceed ~400 lines each. If one wants to grow larger, decompose further.

## Done criteria

Five files in `docs/projects/studio/` (`00-sequencing.md` through `04-compose-mode-spec.md`), each committed. A short summary message naming the v0 release boundary and the highest-risk item.
