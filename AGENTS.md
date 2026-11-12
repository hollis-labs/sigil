# Sigil

Sigil is a UI compiler. Pages, themes and datasources are authored as YAML,
validated against embedded component schemas, and compiled to source code for
one of several render targets (Go/Templ+HTMX, React/shadcn, static HTML
preview). It is not a runtime renderer, a visual editor or a low-code platform,
and what it emits must stand alone — generated code never depends on Sigil at
runtime.

## Start Here

- `README.md` for the CLI surface, `docs/01_architecture.md` through
  `docs/09_sysop-ui-kit.md` for the specs.
- `cmd/sigil/substrate.yaml` is the manifest of the `sigil` program (Substrate app, `variant: cli`);
  `cmd/sigil/substrate_wiring.go` is generated from it (`go run github.com/hollis-labs/substrate/cmd/substrate
  generate -f cmd/sigil/substrate.yaml`). Never edit the wiring file by hand.
- `internal/renderer/renderer.go` defines the `Renderer` interface every target
  implements; `docs/05_renderer-contract.md` explains it.
- `internal/renderer/engine.go` is the only place output reaches disk. Renderers
  return `[]OutputFile`, which is what makes `--dry-run` and `--clean` work.
- `internal/components/builtin/` holds the component schemas as embedded YAML.
  The registry auto-discovers them: adding a file adds a component, and a
  component absent from the registry does not exist.
- `internal/config/types.go` owns the config shape, including which level —
  app, module or page — each field lives on.
- `internal/server/api.go` is the JSON/REST API browser SPAs consume;
  `internal/mcp/` is the agent-facing surface.
- `.sigil/` is a real Sigil project — the dogfood workspace `sigil serve` and
  `sigil generate` read.

## Commands

```bash
make build              # → bin/sigil
make test               # go test ./...
make vet
go run github.com/hollis-labs/substrate/cmd/substrate pins check --app .
python3 scripts/scan-public --tree   # history gate (--ref HEAD) returns after the program-wide history pass
```

`lefthook` gates commits on gofmt/goimports, the pinned `golangci-lint` v2.11.4 (the version CI runs) and
`go vet ./...`, and pushes on `go test ./...`. The `make demo`, `make se-demo`
and `make cw-demo` targets regenerate output into the checked-in demo apps —
use them to see a renderer change, not to verify one.

## Boundaries

Substrate applies to the `sigil` program, never to what it emits. Generated code
must not import `github.com/hollis-labs/substrate`; `TestGeneratedOutputIsFreeOfSubstrate`
holds that. Pins live in Substrate's `bom/bom.yaml`: do not edit versions in `go.mod`
by hand, run `pins check`. Sigil is one binary with one manifest
(`cmd/sigil/substrate.yaml`): the compile commands, `serve` and `mcp serve`. Every HTTP
route and every MCP tool needs a declared scope (`sigil.read` / `sigil.write`; the MCP
endpoint is `sigil.mcp`). `serve` runs with `auth.mode: off` and Sigil refuses any
non-loopback address (`loopbackAddr`, stricter than the framework's check); never change that mode without adding identity. Gaps
go in `docs/conversion-notes.md`.

Hooks under `demo/src/hooks/` and the sibling demo apps are hand-written mocks,
because SWR 2.x does not work with React 19. The demo targets protect them with
`rsync --ignore-existing`; if a generated hook lands on one, restore the mock,
because the generated version does not run.

Runtime dependencies are deliberately few: Cobra, pflag and yaml.v3, plus
`github.com/hollis-labs/libs/mcp/plugin-mcp/go-mcp` and `github.com/modelcontextprotocol/go-sdk` for
`internal/mcp/`, which targets the 2026-07-28 MCP spec. go-mcp deliberately
doesn't wrap resources/prompts, which is why the SDK is used directly. A
further dependency needs the same kind of argument, not an incidental import.

`sigil generate --clean` deletes existing files in the output directory before
generating. Never point `--output` at a directory that holds hand-written code.

A component renders in every target or its unsupported status is documented. A
schema added without a case in each renderer validates and then emits nothing.

## Contributing

Open a pull request against `main`; a maintainer will review it. `CONTRIBUTING.md` has the sequence.
