# Sigil — Orchestrator Boot Prompt

Copy the text below and paste it into a new Claude Code session opened in `~/Projects-apps/sigil`.

---

## Prompt

```
Please boot into Volon:Orchestrator mode.

This is the Sigil project — a system-agnostic UI configuration and code generation tool written in Go.

Follow the Volon boot sequence:

1. Read `.volon/agent-boot.md` (ground truth and project overview)
2. Read `.volon/boot/orchestrator.md` (your role as Orchestrator)
3. Read `.volon/bootstrap.md` (current iteration state and next action)

Then:
- Initialize git if not already initialized (`git init -b main`)
- Emit the Orchestrator boot confirmation block
- Begin the canonical loop starting with Sprint 0

All architecture docs are in `docs/` and sprint guides are in `artifacts/plan/sprints/`. The sprint guides contain detailed task specifications with context, code examples, and acceptance criteria.

Start with Sprint 0: "Project Scaffold + Core Types + Config Parser". Execute tasks in priority order (A before B). Each task should produce a working commit.
```
