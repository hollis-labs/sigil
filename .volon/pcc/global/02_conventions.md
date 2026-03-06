---
intent: pcc_global
project: sigil
updated_at: "2026-03-04"
---

# Conventions

## Coding standards

- Go standard library style
- CLI built with cobra framework
- Internal packages under `internal/`
- Component schemas define validation rules for all 49 types
- Comprehensive documentation in `docs/` (8 numbered spec files)

## Branching and git

- Worktrees enabled (`volon.yaml`: `git.use_worktrees: true`, root: `.worktrees/`)
- Branch prefix: `volon/`
- PR mode: optional
- Commit policy: per-iteration

## Naming conventions

- Page configs: `.sigil/pages/<page-id>.yaml`
- Theme files: `.sigil/themes/<theme-name>.yaml`
- DataSource manifests: `.sigil/datasources/<name>.yaml`
- Custom component schemas: `.sigil/components/<name>.schema.yaml`
- Docs: numbered prefix (e.g., `01_architecture.md`, `08_mcp-integration.md`)
- CSS variables: `--sigil-*` (theme tokens)

## Test commands

```bash
make test      # go test ./...
make vet       # go vet
make build     # build sigil binary
make install   # install locally
```

## Documentation structure

8 spec documents covering architecture, config spec, component model, CLI reference, renderer contract, datasource model, theme system, and MCP integration.

## Volon integration

- Boot profiles: orchestrator, architect, worker, reviewer (in `.volon/boot/`)
- Nanite storage backend for task sync
- Quality scans enabled (dead_code, security, correctness, perf_smells)

## Evidence
- Last refreshed: 2026-03-04 (mentat PCC bootstrap)
- Sources: CLAUDE.md, volon.yaml, Makefile, README.md, docs/
