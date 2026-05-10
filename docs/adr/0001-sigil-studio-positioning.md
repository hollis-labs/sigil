# ADR 0001 — Sigil Studio positioning, scope, and relationship to Sigil

- **Status:** Proposed
- **Date:** 2026-04-14
- **Deciders:** Chrispian (with renderer-agent review)
- **Context pack:** `/Users/chrispian/Projects-apps/agent-workspaces/inbox/sigil_studio_agent_pack`

## Context

A context pack from an ideation session proposes **Sigil Studio**: a Wails + Go + React + tldraw desktop app that wraps Sigil and provides a turn-based, agent-first whiteboard with two modes — *Scene* (freeform tldraw canvas) and *Compose* (Sigil-native semantic authoring). The pack frames it as the "missing tool" for iterative UI design collaboration between a human and an agent.

Sigil itself is a system-agnostic UI configuration and code-generation tool. It already has: a config engine, a component registry with 49 component schemas, three render targets (react-shadcn, go-templ, preview), a theme system, a live dev server (`sigil serve`), and an MCP server. A static-HTML **catalog render target** is currently planned (see `docs/projects/catalog-renderer-prompt.md`).

We need to decide how Sigil Studio relates to Sigil, whether to pursue it, and in what order against the catalog work already planned.

## Decision

1. **Sigil Studio is a sibling product, not a feature of Sigil.** It lives in its own repository. It depends on Sigil as a Go library (and/or MCP client). Sigil must not depend on Studio.
2. **The catalog renderer ships first**, unchanged from its approved plan. It is on the critical path for Studio's Compose mode (it becomes the component palette) and it closes preview-renderer gaps regardless.
3. **Before committing to building Studio**, we will spike a lightweight alternative — **catalog + `sigil serve` + MCP edit tools** — to test whether that already covers the agent-iteration workflow. Studio is the path, not the product. If the lightweight combination satisfies the workflow, we stop there.
4. **If Studio is built**, it is decomposed into at least four sub-projects, each with its own spec:
   1. Wails app shell + document storage + command layer
   2. Scene mode (tldraw + export + handoff bundle)
   3. Sigil bridge (library consumption, loader, preview pipeline, round-trip saver)
   4. Compose mode UI (palette, inspector, preview pane)
5. **Sigil consumption from Studio is Go library for rendering/config + MCP client for agent-facing ops.** No CLI shell-out from Studio. This decision is recorded here; sub-project (3) will refine the specifics.
6. **Studio's command layer is MCP-native from day one.** Studio exposes its own commands over MCP, not only consuming external agent payloads. This makes "agent-first" a platform fact, not a framing claim.
7. **Compose-mode preview in v1 is an iframe to `sigil serve`** with Studio managing a file-watcher on the active Sigil config. We do not build a second preview pipeline.
8. **Round-trip YAML editing uses `yaml.v3` node API from day one** (not `Marshal`). Lossy round-trip is not acceptable past v0.
9. **Scene mode is the lowest-risk v1 win** and can ship as "Sigil Studio v0: visual ideation companion" before Compose exists.

## Consequences

### Positive
- Catalog work proceeds uninterrupted; nothing in the Studio pack invalidates it.
- Studio, if built, inherits catalog's variant files, visual-prop filter, and preview-gap closure for free.
- A deliberate spike on the lightweight path avoids over-building if it's unnecessary.
- MCP-native command layer means Studio is driveable by Claude Code, Cursor, or any MCP client from day one — not just the one agent a user happens to pair with.
- Scene mode is shippable independently; we can learn from real use before committing to Compose.

### Negative / risks
- Round-trip YAML is genuinely hard. Using `yaml.v3` nodes is the right tool but the code is more verbose than `Marshal`/`Unmarshal`. Accept the cost.
- Library consumption couples Studio releases to Sigil's Go API surface. Mitigate with a small adapter package in Studio and a documented Sigil public API.
- Decomposition into four sub-projects makes the project feel slower in the short term. It is not slower in total; it is just honest about scope.
- The original pack implied a single implementation plan. Rejecting that may feel like rescoping. It is.

### Neutral observations
- The Sigil Studio context pack is a strong vision doc and its principles (agent-first, two-state model, diffable intent, turn-based) are adopted. Only the execution shape changes.
- Naming collision between "Sigil" and "Sigil Studio" is acceptable given the sibling-repo structure.

## Alternatives considered

**A. Build Studio as specified in the pack, single plan, single push.**
Rejected. Scope is four projects in a trench coat. Single-plan execution of a project this size reliably produces a stuck state at ~60% completion.

**B. Fold Studio into Sigil as a new CLI subcommand (`sigil studio`) that launches a local app.**
Rejected. Conflates a config/generation tool with a desktop GUI. Complicates Sigil's dependency graph (tldraw, Wails) for users who only want the generator.

**C. Skip Studio entirely; double down on catalog + `sigil serve` + MCP.**
Not rejected. This is the hypothesis being tested by the spike in decision #3. If the spike validates it, the rest of this ADR becomes inert and we write ADR 0002 recording the simpler path.

**D. Library vs CLI vs MCP-client for Sigil consumption from Studio.**
Chose library + MCP client. CLI shell-out is brittle, slow, and harder to debug. MCP client gives Studio the same integration surface external agents use, which is desirable for consistency.

## Open questions (deferred, not decided)

- Handoff-bundle *ingest* path (agent → Studio) is unspecified in the pack. To be resolved in sub-project (1) or a follow-up ADR.
- Whether Sigil should expose a stable public Go API package specifically for Studio, or whether Studio imports internal packages directly. Sub-project (3) decides.
- Whether the catalog's static HTML output is embedded in Studio's Compose palette or whether the palette reads the same variant YAMLs directly. Either works; decide when building the palette.

## Next actions

1. Spike prompt: `docs/projects/lightweight-studio-exploration-prompt.md` — validate whether catalog + `sigil serve` + MCP covers the workflow.
2. Planning prompt: `docs/projects/sigil-studio-planning-prompt.md` — if the spike concludes we still want Studio, the planning agent takes this ADR + the original context pack and produces four sub-project specs.
3. Catalog work proceeds in parallel per `docs/projects/catalog-renderer-prompt.md`.
