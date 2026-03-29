# Frontend Context — Sigil

> Context for working on the Sigil demo app and React/shadcn output. Loaded by the frontend agent.
> Lives at `sigil/.agentrc/agents/frontend.md`.

## What This Agent Covers

The Next.js 16 + Tailwind v4 + shadcn/ui v4 demo app at `demo/`. Use this agent when:

- Working on the demo app (pages, components, styling)
- Installing or configuring shadcn components
- Debugging generated React output in the demo
- Updating mock hooks or data
- Tailwind/theme integration

## Demo App

- **Location:** `demo/`
- **Stack:** Next.js 16, React 19, Tailwind CSS v4, shadcn/ui v4
- **Dev server:** `make demo` (port 3333)
- **Generate + sync:** `make demo-generate`

## Key Conventions

- shadcn v4 uses `@base-ui/react` — no `asChild` prop on triggers
- Tooltip triggers need `render={<span />}` to avoid nested `<button>` hydration errors
- Accordion has no `type="single"` / `collapsible` props in v4
- SWR 2.x is incompatible with React 19 — demo uses mock hooks in `demo/src/hooks/`
- `make demo-generate` uses `--ignore-existing` to avoid overwriting mock hooks
- shadcn imports use `@/components/ui/` prefix
- Use Tailwind classes, not inline styles

## Generated Output

Sigil generates into `demo/` via `make demo-generate`:
- Pages → `demo/src/app/pages/`
- Shared components → `demo/src/components/sigil/`
- Hooks → `demo/src/hooks/` (mock versions, don't overwrite)
- Types → `demo/src/types/`
- Theme CSS → `demo/src/app/globals.css`
