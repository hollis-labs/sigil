# Lightweight Studio Alternative — Exploration Brief

You are a local research + small-spike agent. Your job is to answer one question with evidence:

> **Does the combination of the Sigil catalog renderer, `sigil serve` live preview, and an MCP "edit YAML section" tool cover enough of the agent-iteration workflow that we don't need to build Sigil Studio as a Wails desktop app?**

You are not building Studio. You are testing whether Studio is necessary.

## Why this brief exists

The Sigil Studio context pack (`/Users/chrispian/Projects-apps/agent-workspaces/inbox/sigil_studio_agent_pack`) proposes a full Wails + tldraw desktop app. ADR 0001 (`docs/adr/0001-sigil-studio-positioning.md`) identifies that much of Studio's pitched value — agent-first command layer, iterative preview, diffable config edits, component palette — may already be reachable by composing three things Sigil already has or is about to have:

1. **Catalog renderer** (planned, see `docs/projects/catalog-renderer-prompt.md`) — static component browser.
2. **`sigil serve`** — existing live preview server with SSE reload on file watch.
3. **Sigil's MCP server** (`internal/mcp/`) — could expose an "edit this section of this config" tool and a "create component at this path" tool.

If this combination covers the use cases, we save a multi-month Wails build. If it doesn't, we have clear evidence of the gap and Studio proceeds.

## Your work

### Phase 1 — Read and map the terrain
1. Read ADR 0001 (`docs/adr/0001-sigil-studio-positioning.md`).
2. Read the Studio context pack in full (`/Users/chrispian/Projects-apps/agent-workspaces/inbox/sigil_studio_agent_pack/`).
3. Read the catalog plan (`docs/projects/catalog-renderer-prompt.md`).
4. Skim `internal/server/` (serve + watcher), `internal/mcp/` (existing tools), `docs/08_mcp-integration.md`.
5. Inventory the existing Sigil MCP tools. What's already exposed? What would need adding to reach the "lightweight studio" hypothesis?

### Phase 2 — Map user flows against the lightweight stack
For each user flow from the context pack (`01_product/02_user_flows.md`: Flow A "Agent sketches a UI", Flow B "Agent proposes a real Sigil page", Flow C "Agent explains visually"), answer:

- Can this flow be executed today or with minimal additions using only catalog + `sigil serve` + MCP + a chat interface (Claude Code, Cursor, etc.)?
- If yes, what's the rough UX? Write a 3-5 step walkthrough.
- If no, what specifically fails or is missing?
- For things that are missing, what's the smallest MCP tool or Sigil addition that would close the gap (e.g. "an MCP tool that accepts a tldraw scene JSON and returns a PNG")?

Be specific. "It's less integrated" is not a valid gap. "The agent has no way to produce an image artifact the user can screenshot for further iteration" is.

### Phase 3 — Small, time-boxed spike (optional but preferred)
If time permits (~half-day max), spike the weakest gap you identified. Examples:
- Add an MCP tool `sigil_config_patch` that applies a JSON-patch-like edit to a Sigil config file while preserving YAML comments/formatting via `yaml.v3` nodes.
- Add an MCP tool `sigil_component_add` that appends a component to a named container in a page config.
- Prototype the workflow end-to-end: `catalog` open in a browser tab, `sigil serve` in another, an agent driving both via MCP, human reviewing in the browser.

Do not polish the spike. This is to validate the hypothesis, not to ship.

### Phase 4 — Write a recommendation
Produce `docs/adr/0002-studio-vs-lightweight-path.md` (proposed status) with:

- **Decision**: one of
  - *Accept lightweight path*: catalog + `sigil serve` + MCP additions are sufficient. Close out Studio. Enumerate the MCP tools to add.
  - *Reject lightweight path*: build Studio as planned. Enumerate the specific workflows that cannot be served without Studio.
  - *Partial*: ship lightweight path first; revisit Studio after 4-8 weeks of real use. Enumerate the triggers that would cause us to start Studio.
- **Evidence**: your phase 2 walkthroughs and phase 3 spike results.
- **MCP additions required** (for any decision): a concrete list of tool signatures Sigil should expose, with brief rationale.
- **Cost comparison**: rough effort bands for the two paths, not hour estimates.

## What to avoid

- **Don't rebuild Studio in miniature.** If your answer is "we need a new desktop UI," the lightweight path has failed and that's the honest conclusion — don't start sneaking in Wails shells to make it work.
- **Don't underestimate the catalog.** The catalog (planned, not yet built) is a substantial piece of the hypothesis. Assume it exists at the level described in its brief; don't discount the hypothesis on the basis that the catalog isn't shipped yet.
- **Don't trust the context pack's framing as gospel.** The pack pitches Studio because that's what it's about. Your job is to evaluate the pitch, not propagate it.
- **Don't propose a UI library or framework.** If the answer involves "we need to build X," X must be additions to existing Sigil surfaces (MCP tools, CLI commands, maybe a new render target), not a new application.

## Success criteria

- `docs/adr/0002-studio-vs-lightweight-path.md` committed with a clear Accept / Reject / Partial decision and concrete evidence.
- If Accept or Partial: a prioritized list of MCP tools to add.
- If Reject: a short list of workflows that demonstrably cannot be served without Studio, plus a green-light for the planning-agent brief (`sigil-studio-planning-prompt.md`) to proceed.

## Working rules

- Time-box the spike. If it's taking more than a half-day, you've already answered the question — the gap is real enough to matter.
- One pass through the user flows is enough. Do not go deep on one flow at the expense of the others.
- Flag to the human if you find a third path neither Studio nor the lightweight stack covers well (e.g. "actually what we need is a browser extension" or "the real blocker is that Sigil's config can't express X"). Don't force a binary answer if the evidence says the question is wrong.
