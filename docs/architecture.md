# Architecture (Substrate port)

Sigil is a build-time UI compiler (`01_architecture.md` .. `09_sysop-ui-kit.md` specify the compiler). The `sigil`
program is a [Substrate](https://github.com/hollis-labs/substrate) app: `cmd/sigil/substrate.yaml` selects the
framework modules, the generated `cmd/sigil/substrate_wiring.go` links them, and `internal/` is Sigil's own code.
Substrate supplies mechanism (config, logging, health, HTTP host, auth guard, MCP server, CLI skeleton); Sigil owns
the compiler, its renderers and its commands. Open decisions and framework gaps are in `conversion-notes.md`.

## Rule: Substrate is the program, never the output

Substrate applies to the `sigil` program. It never appears in what Sigil emits. Generated code must stand alone
and must not import `github.com/hollis-labs/substrate` or any `libs/` module.
`TestGeneratedOutputIsFreeOfSubstrate` (`internal/cli/standalone_output_test.go`) generates the go-templ and
react-shadcn targets and fails if an emitted file mentions `hollis-labs/substrate`.

## Decisions (guide section 0)

| Question | Answer |
|---|---|
| Migrated data or fresh start? | No database and no stored accounts. State is the `.sigil/` project directory (plain files), which the port does not touch. Nothing to migrate; no startup guard needed. |
| Binary | One: `cmd/sigil`, `variant: cli`, one manifest. The compile commands (`init`, `new`, `list`, `validate`, `generate`, `preview`, `export`, `import`, `diff`, `migrate`, `schema`, `doctor`), `serve` and `mcp serve` are cobra commands on the cli module's root, which also supplies `version`, `config` and `health`. |
| Who issues identity? | Nobody. `auth.mode: off`, no `identity` module, no database. The framework refuses a non-loopback bind in this mode. |
| Public, scoped or self-authenticating routes? | Every route and tool has a scope (below). No open mounts. |
| Background work | The `serve` file watcher (`internal/server/watcher.go`, a 500 ms poll) and the reload stream. Both stop when the command's context ends, and `serve` joins the watcher before it returns. Open reload streams are closed on shutdown so they do not hold up the graceful stop. |
| UI | None in this repo. The Sysop admin app was dropped by the owner (2026-11-12); a new GUI is a separate task. |
| BOM | Nothing missing. `plugin-mcp` v0.6.0, `go-sdk` v1.8.0, `cobra`, `yaml.v3` are pinned; `libs/security/auth` and `libs/app/app-host` come in through Substrate. `pins check --app .` passes. |

## Modules

`cmd/sigil/substrate.yaml`: `config`, `logging`, `health` (profile `minimal`), `httphost`, `auth` (`mode: off`), `mcp`,
`cli`. `httphost.Start` does not listen; only `serve` calls `httphost.Serve`. There is no per-command module set, so all
of them start for every command (Gap 1); none prints anything until `serve` starts a listener.

## Layout

| Path | What |
|---|---|
| `cmd/sigil/main.go` | `boot.Start`, then `cli.Install` on the cli module's root, then execute. |
| `cmd/sigil/substrate.yaml`, `substrate_wiring.go` | The manifest; the generated wiring (`substrate generate -f cmd/sigil/substrate.yaml`). |
| `internal/cli` | The commands. `Install(root, app)` adds the compile commands, `serve` and `mcp`; `NewRootCmd()` adds the compile commands only, for tests. |
| `internal/server` | The dev server handlers (preview pages, reload stream), the JSON API (`Routes()`), the scope registration, `RequireRequestedWith`, the file watcher. It owns no listener. |
| `internal/hostguard` | The Host check, copied from Basalisk. |
| `internal/mcp` | Tools (`tools()`, one table with scopes), resources and prompts, on the `mcp` module's server. |
| `internal/scopes` | The scope names. |
| `internal/config`, `components`, `datasource`, `theme`, `renderer` | The compiler. No framework concern. |

## Request path (`sigil serve`)

```
client -> hostguard (Host check) -> httphost
                                      |- /api/*   -> RequireRequestedWith (writes) -> auth guard -> API routes
                                      |- /healthz, /mcp, /api/me (framework)
                                      |- /{$}, /pages/{id}, /assets/theme.css, /events -> auth guard -> handlers
```

## Auth model

- `auth.mode: off` in the manifest. The module reads `HOLLIS_AUTH_SIGIL` and `HOLLIS_AUTH` only when the config
  does not set a mode; the manifest sets it, so the environment cannot weaken or change it.
- `httphost.Serve` validates the bind address against the mode: wildcard and public addresses are refused
  (`auth: mode "off" refused on non-loopback, non-tailnet listener ...; use enforce`). Tailnet addresses are
  accepted by the framework (Gap 4).
- Hosts other than loopback names and addresses get `403` before any handler runs (DNS rebinding).
- Every write under `/api/` needs a non-empty `X-Requested-With` header (Gap 3).
- There is no data owner: no accounts, no per-user rows. The project directory belongs to the user running `sigil`.
- `serve` logs `auth is not enforcing` at WARN once, with the listen address, when the listener starts (Gap 1).

## Scopes

Declared in single tables that register the handler and the scope together (`server.API.Routes`,
`server.Server.Previews`, `mcp.tools`); tests fail on an entry without a scope.

| Scope | Routes and tools |
|---|---|
| `sigil.read` | every `GET /api/*`; `GET /{$}`, `/pages/{id}`, `/assets/theme.css`, `/events`; the tools `sigil_list_pages`, `sigil_get_page`, `sigil_validate`, `sigil_list_components`, `sigil_get_component_schema`, `sigil_list_datasources` |
| `sigil.write` | `POST /api/pages`, `PUT /api/pages/{id}`, `POST /api/datasources`; the tools `sigil_create_page`, `sigil_update_page`, `sigil_create_datasource` |
| `sigil.mcp` | the `mcp` module's `/mcp` endpoint and the stdio caller (enforce mode only) |

None is enforced in mode off; they are declared so the tables are correct the day the mode changes.

## Data

None stored by Sigil beyond the project's own `.sigil/` files, written with the permissions `sigil init` has always
used (directories 0755, files 0644).

## Config

The manifest sets the defaults. Any key can be overridden in a YAML file named by `SUBSTRATE_CONFIG` (one level
deep: a block's keys replace the manifest's).

| Key | Default | Meaning |
|---|---|---|
| `httphost.addr` | `127.0.0.1:3210` | `serve` listen address (`--port` changes its port). Sigil refuses any non-loopback host. |
| `auth.mode` | `off` | Leave it. Changing it needs the identity module and a data directory first. |
| `logging.level`, `logging.format` | `info`, `text` | Log output (to stderr). |
| `mcp.version` | `0.1.0` | Version the MCP server reports. |
| `mcp.stdio_token_file` | none | Needed only if `auth.mode` is changed. |

There is no `sigil:` block: the project directory is the `--sigil-dir` flag on `serve` and `mcp serve`, as it is on
every other command.

## What was dropped, and why

- Sysop (owner decision 2026-11-12), `GET /api/health`, the `-ldflags` `Version` variable and
  the old `sigil version` text, the `{"error": ...}` body. Details and replacements in `conversion-notes.md`.
