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
```

`lefthook` gates commits on gofmt/goimports, `golangci-lint run --new` and
`go vet ./...`, and pushes on `go test ./...`. The `make demo`, `make se-demo`
and `make cw-demo` targets regenerate output into the checked-in demo apps —
use them to see a renderer change, not to verify one.

## Boundaries

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
