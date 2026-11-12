# Conversion notes: Sigil onto Substrate

Clean break for code; no database, so no data migration. Work is uncommitted on `feat/substrate-conversion`.
Shape: one binary (`cmd/sigil`), one manifest (`cmd/sigil/substrate.yaml`): `variant: cli`, modules `httphost`,
`auth` (mode `off`), `mcp` and `cli`, plus the framework's config, logging and health. The compile commands,
`serve` and `mcp serve` are `cobra` commands on the cli module's root.

## Dropped, by owner decision (2026-11-12)

- **Sysop** (the `sysop/` tree: the `sigil_sysop` binary, its copy of the API wiring and the React/Vite admin SPA) is
  not supported and not ported. It is removed from the branch with its CI steps, Makefile targets, docs and the
  `ui-go` dependency. A new GUI for Sigil is a separate later task. The `--ui-kit sysop` renderer option (generated
  React against `@hollis-labs/sysop-ui`) is a different thing and stays.
- **Deployed unit: done.** The owner stopped and disabled the user service `sigil-sysop.service` on the agent box (the
  driver verified: inactive, disabled, nothing listening on `127.0.0.1:8094`). Nothing is left to retire.
- `GET /api/health`: the framework serves `GET /healthz` instead.
- Dev server timeouts: `httphost` sets `ReadHeaderTimeout` 10s and `IdleTimeout` 2m; the old server set none.
- `sigil serve --port` is kept. Without it the address is `httphost.addr` (default `127.0.0.1:3210`, the old default); with it, that host and the given port. A non-loopback host is refused (item 4).
- `sigil version` text and the `-ldflags` `Version` variable: see item 13.

## Gaps

1. **No per-command module set (resolved in the framework).** `boot.Start` starts every manifest module, so
   `httphost`, `auth` and `mcp` start for every command, including `sigil generate`. They hold nothing on disk and
   open no listener until `serve`. The `auth is not enforcing` warning used to be logged at module start, on stderr
   for every command; Substrate now logs it where a listener starts (`httphost.Serve`), so compile commands print
   nothing (`TestCompileCommandsPrintNothingOnStderr`) and `sigil serve` prints it once with the listen address.
2. **No shared Host check.** Basalisk carries its own `internal/hostguard`; Substrate has none. Sigil uses a copy of
   Basalisk's file and test as they are (`internal/hostguard`, package comment says so), called with no extra host
   names, outermost in the handler chain. A third variant was not written.
3. **No Substrate rule for `X-Requested-With` when auth is off.** Only the identity module enforces it, and Sigil
   has none. The dev server writes project files, and any web page can send a simple cross-origin POST to a
   loopback port, so `server.RequireRequestedWith` applies identity's rule (a non-empty header on every
   `POST`/`PUT`/`DELETE`; the server answers no preflight) to the `/api/` routes. Reads and the preview pages need
   nothing. Not applied to `/mcp`: MCP clients do not send it.
4. **The bind check accepts tailnet addresses in mode off.** `httphost.Serve` refuses a wildcard or public address
   with `auth: mode "off" refused on non-loopback, non-tailnet listener ...; use enforce`, but allows loopback,
   `localhost`, `*.ts.net` names and the tailnet IPv4 and IPv6 ranges. "Loopback only" is therefore a configuration
   convention beyond that: `sigil serve` on a tailnet address would start unauthenticated, with file-write routes.
   The owner's decision is auth off on loopback only, so `sigil serve` checks the effective address (the configured
   `httphost.addr`, or its host with `--port`) before `httphost.Serve` and refuses anything but 127.0.0.0/8, `::1` and
   `localhost` (`loopbackAddr` in `internal/cli/serve.go`; tests cover 0.0.0.0, a LAN address, a tailnet range
   address and a `.ts.net` name). The framework's acceptance of tailnet addresses in modes off and permissive is
   intended (Substrate `docs/identity.md`); Sigil's stricter rule is deliberate and stays.
5. **`mcp` mounts `/mcp` whenever it is selected.** `sigil serve` therefore also serves the module's Streamable HTTP
   endpoint at `/mcp` (scope `sigil.mcp`). Sigil registers its tools, resources and prompts only in `mcp serve`, so
   over HTTP the server answers with no tools; nothing new is exposed. The module offers no stdio-only selection.
   `sigil mcp serve` is a local stdio process launched by the MCP client (an agent or editor) on the user's own
   machine, as the user. The owner chose local filesystem trust for a locally launched stdio MCP process, with
   enforced identity kept on HTTP surfaces. What the module requires of it: nothing in mode off; in `permissive` or `enforce` a token file
   (`mcp.stdio_token_file`) the auth module can verify, holding `sigil.mcp` in `enforce`. Tools are visible to the
   module through `Server.SDKServer()` for resources and prompts (go-mcp does not wrap them), which the module's
   API cannot express on its own: kept, as before.
6. **Index route.** `httphost.Mount` refuses the pattern `/`, so the index page is mounted as `GET /{$}`. net/http
   accepts it beside the SPA catch-all. The scope table does not know `{$}`, so `RegisterScopes` maps it to `/`.
7. **Manifest cannot carry the app's own config** (CW-20261112-0079). Sigil has no config block of its own: the
   listen address is `httphost.addr` (or `--port`) and the project directory is the `--sigil-dir` flag, so there is no
   `internal/settings` package.
8. **Error body** changed from `{"error": "..."}` to `{code, message, details?}` on `/api/*`. Client impact: the
   Sysop SPA (dropped) was the only client in this repo. The preview pages' plain-text errors are unchanged.
9. **Route and tool scopes.** `sigil.read` for every `GET` route (including the preview pages and the reload
   stream) and every read-only tool, `sigil.write` for `POST`/`PUT` and the tools that write, `sigil.mcp` for `/mcp`.
   The tables are single lists that register the handler and the scope together (`server.Routes`,
   `server.Previews`, `mcp.tools`); tests fail if an entry has no scope. None is enforced in mode off.
10. **Hygiene: tree clean; history findings wait for one planned pass after the migration (owner decision).** The repository
   is public. The working tree is clean (one mock private IPv4 address in `demo/src/hooks/use-asset.ts` was changed to
   a documentation-range address). History has 37 findings, none a credential (no `secret` class finding): path,
   host and address classes only. Owner's rule for every repo: the tree must be clean in the port PR; history is left alone and gets one planned
   rewrite pass after the migration. No history rewrite during the port, no force-push, no allow file. CI runs
   `scan-public --tree`; the history gate (`--ref HEAD`) returns after that pass. Table to execute from, from
   `python3 scripts/scan-public --ref HEAD` (74 commits on `main`; "paths still in HEAD" are paths that still exist but
   no longer contain the match):

| Class | Rule | Findings | Example path | Earliest first-add | Latest first-add | Last commit touching an example path | Paths still in HEAD |
   |---|---|---|---|---|---|---|---|
   | host | agent-box-host | 3 | `.agent-ops/project.yaml` | c9b90c7 2026-05-15 | 0804e45 2026-09-07 | 0804e45 2026-09-07 | 1 of 3 |
   | private-path | home-path-macos | 33 | `.agentrc/boot-prompt.md` | 7561dbc 2026-03-06 | 5885287 2026-09-09 | 0804e45 2026-09-07 | 2 of 33 |
   | internal-ip | internal-ip | 1 | `demo/src/hooks/use-asset.ts` | 4820abe 2026-03-29 | 4820abe 2026-03-29 | 4820abe 2026-03-29 | 1 of 1 |
   
   Every finding (class, path or commit message, commit that first added it):
   
   - host: `message` first added in 0804e45
   - host: `.agent-ops/project.yaml` first added in c9b90c7
   - private-path: `.agentrc/boot-prompt.md` first added in 63d231c
   - private-path: `.agentrc/commands` first added in 3efb070
   - private-path: `.agentrc/roles` first added in 3efb070
   - private-path: `.agentrc/skills` first added in 3efb070
   - private-path: `.claude/skills/adr.md` first added in baab3ae
   - private-path: `.claude/skills/agentrc-install.md` first added in baab3ae
   - private-path: `.claude/skills/agentrc-manage.md` first added in baab3ae
   - private-path: `.claude/skills/blg.md` first added in baab3ae
   - private-path: `.claude/skills/blueprint-builder.md` first added in baab3ae
   - private-path: `.claude/skills/boot-prompt.md` first added in baab3ae
   - private-path: `.claude/skills/doc-note.md` first added in baab3ae
   - private-path: `.claude/skills/doc-search.md` first added in baab3ae
   - private-path: `.claude/skills/escalate.md` first added in baab3ae
   - private-path: `.claude/skills/go-build.md` first added in baab3ae
   - private-path: `.claude/skills/go-lint.md` first added in baab3ae
   - private-path: `.claude/skills/go-test.md` first added in baab3ae
   - private-path: `.claude/skills/hadron-run.md` first added in baab3ae
   - private-path: `.claude/skills/nanite` first added in baab3ae
   - private-path: `.claude/skills/qhealth.md` first added in baab3ae
   - private-path: `.claude/skills/qstatus.md` first added in baab3ae
   - private-path: `.claude/skills/shadcn-install.md` first added in baab3ae
   - private-path: `.claude/skills/shadcn-ui` first added in baab3ae
   - private-path: `.claude/skills/standup.md` first added in baab3ae
   - private-path: `.claude/skills/vault-search.md` first added in baab3ae
   - private-path: `.volon/boot/reviewer.md` first added in 7561dbc
   - private-path: `.volon/boot/worker.md` first added in 7561dbc
   - host: `AGENTS.md` first added in c9b90c7
   - private-path: `BOOT-PROMPT.md` first added in 63d231c
   - internal-ip: `demo/src/hooks/use-asset.ts` first added in 4820abe
   - private-path: `docs/adr/0001-sigil-studio-positioning.md` first added in 2031f5d
   - private-path: `docs/projects/lightweight-studio-exploration-prompt.md` first added in 2031f5d
   - private-path: `docs/projects/sigil-studio-planning-prompt.md` first added in 2031f5d
   - private-path: `docs/projects/stack-explorer-api-prompt.md` first added in 63d231c
   - private-path: `docs/projects/stack-explorer-frontend.md` first added in 63d231c
   - private-path: `sigil.cerberus.yaml` first added in 5885287
11. **Lint: `golangci-lint run` (v2.11.4, `.golangci.yml`) is clean for the whole repo (a separate change from the port).**
    `main` had 208 findings once the per-linter display caps are lifted (`--max-issues-per-linter=0
    --max-same-issues=0`); the default run showed 67. Fixed what the linter named, no refactors: `--fix` for about 90
    staticcheck, misspell and errorlint items (`Fprintf` for `WriteString(Sprintf)`, tagged `switch`, `errors.Is`;
    the misspell fix touched a comment only); errcheck (`Flush` errors returned, `MarkFlagRequired` and
    `ResponseWriter.Write` results marked ignored); six shadowed `err`; the predeclared names `max`, `min`, `cap`
    renamed; one dead assignment removed; three unused helpers deleted; unused parameters dropped from private
    functions only; four private generators lost an always-nil `error` result; in tests, file and directory modes
    tightened to 0600/0750. **Decision (accepted by the driver; do not "fix" these):** 39 gosec findings in production code (file and
    directory modes 0644/0755, reading project files by path, the browser opener) keep their behaviour and carry a
    `//nolint:gosec` with a reason, matching the pattern already in the repo, because tightening the modes would
    change the permissions of every scaffolded and generated file, a behaviour change to the compiler's output that
    is not part of this port. Generated output was compared before and after
    on the dogfood project: go-templ (105 files) and the preview HTML are byte-identical; react-shadcn (128 files) and
    the sysop-kit variant (126) differ only in line order, which is item 12. `lefthook.yml` and `AGENTS.md` now
    describe a whole-repo lint gate.
12. **Pre-existing: react-shadcn output is not deterministic.** Two runs of the unmodified `main` binary differ (48
    changed lines across `pages/*.tsx`): Go map iteration order leaks into prop order and `variants` maps. The
    ported build differs from `main` only in that same ordering (checked: no file differs beyond line order). Not
    fixed here; it affects any byte-for-byte comparison of generated React.
13. **Version output (client impact: none found).** The cli module prints build information, replacing the `-ldflags`
    `Version` variable and the old `sigil version` command ("sigil dev", or the `git describe` string from `make`).
    Measured: a clean build at a tagged commit prints `sigil v0.9.9 <revision>`; a plain `go build` in a checkout
    prints `sigil v0.0.0-<date>-<revision>` (with `+dirty` when the tree is modified); a build without VCS
    information (`-buildvcs=false`, a source tarball) prints `sigil (devel)`. Searched for consumers of the old
    string: the repo has no deploy.yaml or install script, `sigil.cerberus.yaml` only runs `make build` and
    `./bin/sigil serve`, and the MCP server reports its own version (`mcp.version`, `0.1.0`), which never read
    `Version`. The deployed sigil-sysop unit is being retired by the owner, so nothing downstream of it matters.


## Known limits

- **The preview index does not escape project values.** `internal/server/index.go` writes page id, title, overlay and
  module into the index HTML as they are (it did so before the port). With the Host check and the write guard, a
  page on another site cannot plant content, so the risk is low, but a project file with markup in a title runs in
  the preview. Not changed here; the new Sigil GUI task (CW-20261112-0108) must not inherit it.

## Old to new

### Commands and flags

| Old | New |
|---|---|
| `sigil version` ("sigil dev") | `sigil version` from the cli module (build info); `-ldflags ...cli.Version` removed |
| (none) | `sigil config`, `sigil health` from the cli module |
| every compile command (init, new, list, validate, generate, preview, export, import, diff, migrate, schema, doctor) | unchanged, installed on the cli module's root by `cli.Install` |
| `sigil serve --port N` (default 3210) | unchanged flag (default: the port of `httphost.addr`, 3210); the host is `httphost.addr`'s and must be loopback |
| `sigil serve --sigil-dir` | unchanged |
| `sigil mcp serve` | unchanged command; the server is the `mcp` module's (`mcp.version` replaces the hard-coded `0.1.0` in the manifest) |
| `sigil_sysop -addr -sigil-dir` | removed (Sysop dropped) |

### HTTP

| Old | New |
|---|---|
| `GET /api/health` | `GET /healthz` (framework) |
| `GET/POST/PUT /api/pages`, `/api/components`, `/api/datasources`, `/api/themes`, `/api/project` | same paths and bodies; writes need `X-Requested-With`; errors are `{code, message, details?}` |
| `GET /`, `/pages/{id}`, `/assets/theme.css`, `/events` | same, mounted with scope `sigil.read`; a foreign `Host` gets 403 |
| (new) | `GET /api/me`, `GET /healthz`, `/mcp` (framework modules) |

### Environment

| Old | New |
|---|---|
| `NO_COLOR`, `--no-color`, `--verbose` | unchanged (presentation, not settings) |
| (none) | `SUBSTRATE_CONFIG`: YAML file overriding any manifest key (`httphost.addr`, `logging`, `auth`, `mcp`); the merge is one level deep |

### Operator

- `sigil.cerberus.yaml`: the `--port 3210` argument is gone (the default address is the same); nothing else changes.
- `sigil serve --port N` works as before; `httphost.addr` in a `SUBSTRATE_CONFIG` file changes the default.
