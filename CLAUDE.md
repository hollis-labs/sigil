# Volon — Agent Boot Instructions

This repo runs under **Volon**, an agentic development system for Claude Code.

Volon filepath: ~/Projects-apps/volon

Store all your Volon related output in the current project folder/* for this session.

## Boot Sequence

Before any other work:

1. List the available profiles in `.volon/boot/` (e.g. `orchestrator`, `architect`, `worker`, `reviewer`)
2. Ask the user which profile to use — present the options clearly, plus "Other: Specify"
3. Wait for the user's selection before proceeding
4. Read `.volon/agent-boot.md` (core rules and ground truth sources)
5. Read `.volon/boot/<selected-profile>.md`
6. Read `.volon/bootstrap.md` (current iteration state)
7. Follow the selected profile's instructions, emit the boot confirmation block, and begin work

## Key files

| File | Purpose |
|---|---|
| `volon.yaml` | System configuration |
| `.volon/bootstrap.md` | Current iteration state — start here |
| `.volon/agent-boot.md` | Full boot reference (rules, reference map) |
| `.volon/boot/` | Role-specific boot profiles |
| `.volon/pcc/` | Project context cache |
| `docs/01_architecture.md` | System architecture |
| `docs/02_config-spec.md` | Sigil config YAML specification |
| `docs/03_component-model.md` | Component type system |
| `docs/04_cli-reference.md` | CLI commands reference |
| `docs/05_renderer-contract.md` | Renderer plugin contract |
| `docs/06_datasource-model.md` | DataSource abstraction |
| `docs/07_theme-system.md` | Theme tokens and CSS generation |
| `docs/08_mcp-integration.md` | MCP server for AI/agent integration |

## Profiles

| Profile | Use when |
|---|---|
| `orchestrator` | Driving implementation, managing tasks, writing to the repo |
| `architect` | Planning, ADRs, doc-only work — limited write scope |
| `worker` | Bounded read-only analysis (scans, audits, reports) |
| `reviewer` | Structured reviews, verification, targeted investigation |
