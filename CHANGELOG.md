# Changelog

All notable changes to Sigil are recorded here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html). Sigil is pre-release and has not cut a tagged version: there is no compatibility promise, and breaking changes can land in any commit. This file is backfilled from the git log on a good-faith basis, not exhaustively.

## [Unreleased]

### Changed

- **Sigil is a Substrate app** (`cmd/sigil/substrate.yaml`): one binary with the framework's config, logging, health, HTTP host, auth guard, MCP server and CLI skeleton. `sigil version`, `sigil config` and `sigil health` come from the framework; `sigil version` prints the module version and VCS revision from the build instead of the `-ldflags` string.
- `sigil serve` listens on `httphost.addr` (default `127.0.0.1:3210`, set in the `SUBSTRATE_CONFIG` file) (`--port` still overrides the port); Sigil refuses any address that is not loopback. A request for another `Host` gets 403, and API writes need an `X-Requested-With` header.
- `/api/*` errors use `{code, message, details?}` instead of `{"error": "..."}`; health is `GET /healthz`, replacing `GET /api/health`.
- `sigil mcp serve` runs on the Substrate `mcp` module; every route and tool declares a scope (`sigil.read`, `sigil.write`).
- Lint gate: `golangci-lint run` is clean for the whole repository.

### Removed

- **Sysop** (the `sysop/` admin app and its `sigil_sysop` binary) is dropped. A new GUI is a separate, later task.

### Security

- Go toolchain 1.26.10 (`go` directive and CI pin) for GO-2026-6809 (crypto/x509); no behavior change.

### Added

- **Sigil Sysop UI admin surface** (`sysop/`), a separate Go module with a Vite frontend over a Sigil project.
- **JSON/REST API** (`internal/server/api.go`) for browser SPAs: pages, components, datasources, themes and project.
- **`--ui-kit sysop`** for the `react-shadcn` target, generating pages and shell against `@hollis-labs/sysop-ui`.
- **`target_mode` (`spa` | `app-router`)** and a provider abstraction for custom app-shell providers in `react-shadcn`.
- Demo apps generated from the `.sigil/` project: Forge (`make demo`), Stack Explorer (`make se-demo`) and Clockwork (`make cw-demo`).
- UI compiler direction draft (`docs/ui-compiler-direction.md`).

### Changed

- Use published `libs/plugin-mcp` v0.1.1 for MCP and `libs/ui-go` v0.1.0
  for the Sysop SPA harness. Both Go modules require Go 1.26.8 and pin the
  1.26.9 toolchain; their external dependencies use released module versions.
  The Sysop module retains its existing link to the parent Sigil module in
  the same checkout. CI checks both modules, including the embedded Sysop UI
  build, without regenerating the hand-authored demo or admin screens.

- **MCP server rebuilt on `go-mcp`** and the official MCP Go SDK, targeting the 2026-07-28 spec, replacing the hand-rolled stdio JSON-RPC server.
- Module path renamed to `github.com/hollis-labs/sigil`.
- **`sigil serve` and `sigil_sysop` now bind `127.0.0.1` by default**; they previously listened on every interface with no authentication.

### Fixed

- `react-shadcn` refetch codegen corrected after dogfooding the Sysop UI.

## Pre-release history

- **2026-03-28 / 03-29.** Light theme, shadcn v4 demo app, new components; renderer fixes, data-driven grids, charting, command palette, form state and responsive behavior across the Forge demo.
- **2026-03-02.** The MVP and first extensions landed in quick sequence: MCP server (9 tools, 7 resources, 2 prompts over stdio), Go/Templ renderer, live dev server with SSE auto-reload and file watching, React/shadcn renderer, preview, export/import, config diff, schema migration, JSON Schema export and `doctor`.
- **2026-03-01.** Project scaffold: core config types, YAML parser and 14-check validator, Cobra CLI (`init`, `new`, `list`, `validate`), and the component registry (49 built-in types; 57 today).
