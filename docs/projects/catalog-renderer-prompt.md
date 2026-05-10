# Sigil Catalog Renderer — Executor Brief

You are implementing a new **catalog** render target for Sigil: a static-HTML component catalog that lets humans and agents visually browse what components exist, what variants look like, and what the active theme provides. It is NOT a Storybook integration and NOT an authoring tool.

Read this brief end-to-end, then produce a concrete implementation plan (tasks, file list, test strategy, order of work). Execute via subagents — parallelize independent work (e.g. closing preview-renderer gaps per component).

---

## Context

### Project
Sigil generates framework-specific UI code from YAML page definitions. It has three existing render targets:

- `react-shadcn` — `.tsx` + SWR hooks + TS types
- `go-templ` — `.templ` + handler stubs + HTMX
- `preview` — standalone HTML with inline CSS + mock data

Component schemas live in `internal/components/builtin/*.yaml` and are embedded in the binary. 49 components across 6 categories (primitives, layouts, navigation, data, overlay, forms). Themes live in `.sigil/themes/*.yaml` and compile to CSS custom properties (`--sigil-*`).

### Problem being solved
There is no way to browse the component catalog. You can generate and preview individual pages, but there's no index of "what components exist, what they look like, what variants are configured, what the theme provides." The catalog fills that gap for design conversations with agents.

### Reference files
| Need | File |
|---|---|
| Renderer interface | `internal/renderer/contract.go` (or see `docs/05_renderer-contract.md`) |
| Example renderer (simplest) | `internal/renderer/preview/preview.go` |
| Schemas | `internal/components/builtin/*.yaml` |
| Theme loader | `internal/theme/` |
| CLI entry | `internal/cli/` + `cmd/sigil/` |
| Agent context for this work | `.agentrc/agents/renderer.md` |
| Architecture | `docs/01_architecture.md` |

### Boot
Boot the `sigil-renderer` agent (from `.agentrc/config.yaml`) before starting. It loads renderer-specific context including the list of components with no preview renderer case.

---

## Scope (v1 / MVP)

### In scope

1. **New render target `catalog`** at `internal/renderer/catalog/`, implementing the standard `Renderer` interface. Emits static HTML + CSS + a tiny hand-written JS file.
2. **Variant authoring format**: a sibling `*.catalog.yaml` file next to each component schema in `internal/components/builtin/`. Embedded in the binary like schemas. Curated, not combinatorial — each file lists the variants that should appear on the detail page. v1 can ship with a small seed set (buttons, badges, a few layouts); the structure must support adding more without code changes.
3. **Output structure**:
   - `index.html` — gallery, grouped by category. Each tile = one component rendered in its canonical/default state, linking to its detail page.
   - `components/<type>.html` — detail page per component: rendered variants, each with its name + a short list of the visual props and their values.
   - `theme.html` — read-only theme explorer. Lists tokens from the active theme config (colors, radii, shadows, typography) with live visual samples. No editing.
   - `assets/catalog.css` — catalog chrome (layout, cards, nav). Component CSS comes from the normal theme output.
   - `assets/catalog.js` — theme toggle only. Swaps `data-theme` on `<html>`.
   - `assets/theme-<name>.css` for each configured theme (light, dark, etc.) so the toggle works without rebuilding.
4. **Rendering of components themselves reuses the `preview` renderer.** The catalog wraps preview's per-component HTML.
5. **Close preview-renderer gaps.** Per `renderer.md`, preview currently lacks cases for: `card`, `switch`, `accordion`, `accordion-item`, `combobox`, `dropdown-menu`, `tooltip`, and every component in the "schemas exist, no React render case yet" list where preview also has no case (verify by reading `internal/renderer/preview/preview.go`). Add preview rendering for each missing component so the catalog is complete. This is the forcing function — do not skip.
6. **Visual-prop filter (heuristic with override)**: props with an `enum` are "visual" (variant, size, align, etc.). Props without an enum are content/behavior and are hidden from the main detail view. Schemas may override with an explicit `style: true` or `style: false` key on a prop. Apply the override when present; fall back to the enum heuristic otherwise.
7. **CLI**: invoke via `sigil generate --target catalog --out <dir>`. Standard flags (`--theme`, `--dry-run`, `--clean`) should work. No new top-level command needed.
8. **Tests**: snapshot/golden-file tests for emitted HTML per component, index page, and theme page. Follow existing renderer test patterns (see `internal/renderer/preview/*_test.go` and `internal/renderer/reactshadcn/*_test.go`).

### Explicitly out of scope for v1
- Interactive prop or theme editing.
- Info-button modal with full props table, actions list, or copy-pasteable YAML snippet.
- Per-component "tokens used" lists on detail pages.
- Project-level catalog file overrides (`.sigil/catalog/*.yaml`). The escape hatch is regular Sigil pages.
- `sigil serve` integration (live catalog route). Static generation only for v1.
- Combinatorial variant generation. Variants are explicitly authored.
- Search / filter UI on the index. Category grouping is enough.

---

## Key decisions (already made — don't re-open)

| Decision | Value |
|---|---|
| Primary consumer | Designer/agent browsing what exists before authoring YAML |
| Entry format | Gallery index + detail page per component |
| Variant source | Sibling `*.catalog.yaml` in `internal/components/builtin/`, embedded. No project-level overrides in v1. |
| Delivery | Static HTML via new render target. No `serve` integration. |
| Component rendering | Reuse `preview` renderer; close its gaps as part of this work. |
| Theme handling | Emit CSS for all configured themes; JS toggle switches `data-theme`. |
| Theme explorer | Read-only page in v1. Editor is a separate future project. |
| Detail page content | Rendered variants + visual props list. Everything else is post-MVP. |
| Visual-prop identification | Enum-heuristic + optional schema `style:` override. |

---

## Constraints and patterns

- **Renderers are stateless.** All inputs via `RenderContext`; all outputs via `[]OutputFile`. Never write files directly. This preserves `--dry-run` and `--clean`.
- **Generated output must be standalone.** No Sigil runtime, no build step for the catalog. A user should be able to `open index.html` after generation.
- **Component schemas are the source of truth.** If a component isn't in the registry, it doesn't exist. Likewise, the catalog auto-discovers; no per-component registration.
- **Don't add framework logic to `engine.go`.** All catalog-specific code in `internal/renderer/catalog/`.
- **Follow existing renderer structure.** Look at `preview/` for the simplest precedent.
- **Theme CSS emission**: use or extend the existing theme renderer. Don't duplicate token-to-CSS logic. If the current theme system only emits one theme at a time, factor a small helper so the catalog can request all configured themes.
- **Embedding**: use Go's `embed` package for `*.catalog.yaml` files the same way schemas are embedded (check `internal/components/` for the existing pattern).
- **Test patterns**: golden files live next to test files. Regenerate with a `-update` flag convention if that's what the repo already uses — check first.

---

## Variant file format (proposed — finalize in your plan)

```yaml
# internal/components/builtin/button.catalog.yaml
component: button
variants:
  - name: Primary
    props:
      label: Save
      variant: primary
  - name: Destructive small
    props:
      label: Delete
      variant: destructive
      size: sm
  - name: With icon
    props:
      label: Add item
      variant: secondary
      icon: plus
```

The canonical/default rendering on the index page uses the first variant, or a `default:` block if the schema/catalog file specifies one. Decide which in the plan.

---

## Deliverables

1. **Implementation plan** — tasks in dependency order, files to add/modify, test strategy, risk callouts.
2. **Working v1** — `sigil generate --target catalog --out ./catalog-out` produces a browseable static site covering all 49 components, with curated variants for at least: button, badge, input, select, tabs, card, data-table, alert, modal. The remaining components appear on the index with at least a default rendering (no variants required for v1 seed set, but the machinery must handle them).
3. **Preview renderer gaps closed** for every component the catalog needs. Each new preview case has test coverage.
4. **Theme explorer page** rendering the active theme's tokens with visual samples.
5. **Theme toggle** working in a browser with two themes emitted (light + dark if both are configured in the project; otherwise just the one).
6. **Tests**: golden-file coverage for catalog output + unit coverage for the variant loader and visual-prop filter.
7. **Docs**: update `docs/04_cli-reference.md` with the new `--target catalog` option; add `docs/09_catalog.md` (short) describing the variant file format and how to add variants.

---

## Execution guidance

- Use subagents to parallelize: preview-gap work per component is independent and parallelizable. The catalog machinery (loader, index generator, detail generator, theme page, CLI wiring) should be done by one agent/thread to keep coherent.
- Ship preview gaps first (they unblock the catalog), then catalog infrastructure, then variant seeds, then polish.
- Run `go build ./cmd/sigil/` and `go test ./...` after each milestone. Don't defer to the end.
- If anything in this brief conflicts with what you find in the code, trust the code and note the discrepancy in your plan before proceeding.

Start by booting the `sigil-renderer` agent and producing your implementation plan. Do not begin coding until the plan is written.
