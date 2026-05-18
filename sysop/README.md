# Sigil Sysop UI

A **Sysop UI** for Sigil — a single Go binary that serves Sigil's JSON API
and a React admin SPA, same-origin. It is an admin surface over a Sigil
project: browse and inspect its pages, components, datasources, and themes.

This sub-project is the dogfood loop for the Hollis Labs UI toolchain — it
was scaffolded by `folio new sysop-ui` and built on the shared
[`@hollis-labs/sysop-ui`](https://github.com/hollis-labs/sysop-ui) kit, served
by [`go-webui`](https://github.com/hollis-labs/go-webui).

## What it shows

The SPA has five screens, each backed by a Sigil API endpoint:

| Screen | Endpoint | Content |
|---|---|---|
| Overview | `/api/project` + counts | Project config + page/component/datasource/theme counts |
| Pages | `/api/pages`, `/api/pages/{id}` | Every page config; row-click opens the parsed page |
| Components | `/api/components`, `/api/components/{type}` | The component registry (builtin + custom) |
| Datasources | `/api/datasources`, `/api/datasources/{alias}` | Datasource manifests + capabilities |
| Themes | `/api/themes`, `/api/themes/{name}` | Theme token sets |

The API is Sigil's own `internal/server.API` — the same handler `sigil serve`
mounts — so there is no API logic duplicated here.

## Layout

```
cmd/sigil_sysop/main.go   Go entrypoint — mounts Sigil's /api + the SPA
internal/webui/embed.go   //go:embed all:dist + the go-webui handler
frontend/                 Vite + React frontend (the Sysop UI)
  src/App.tsx             app shell — nav rail + screen switch
  src/pages/              one component per screen
  src/lib/                EntityListPage template + hooks
  src/api/                same-origin API client + typed context
```

The frontend builds into `internal/webui/dist/`, which the Go binary embeds —
so a single binary serves both the API and the UI.

## Prerequisites

- Go 1.26.1+
- Node.js 20+ / npm

## Build and run

```sh
make all                       # ui-build (vite → internal/webui/dist) then build
./sigil_sysop -sigil-dir PATH  # PATH is a Sigil project's .sigil directory
```

`-sigil-dir` defaults to `.sigil`, so running from a Sigil project root needs
no flag. `-addr` (default `:8080`) sets the listen address. The UI is then at
<http://localhost:8080/sysop/>; the bare root redirects there.

Before the first `make ui-build`, `go-webui` serves a "not built" placeholder
in place of the app.

## Develop

Two processes during development:

```sh
make run      # Go server on :8080 (serves /api and the last UI build)
make ui-dev   # Vite dev server with hot reload — proxies /api to :8080
```

| Command | What it does |
|---|---|
| `make ui-build` | Build the frontend into `internal/webui/dist` |
| `make ui-dev` | Run the Vite dev server (hot reload) |
| `make build` | Build the Go binary |
| `make all` | `ui-build` then `build` |
| `make run` | Build and run the server |
| `make install` | Build the embedded UI + install `sigil_sysop` to `$GOBIN` |
| `make test` / `make vet` | Go test / vet |

## Adding a screen

A list screen is column defs plus two API calls — see `src/lib/entity-page.tsx`
(`EntityListPage`) and any page under `src/pages/`. Add the page component,
register it in the `SCREENS` map in `src/App.tsx`, and add its endpoints to
`src/api/client.ts`.

## A note on the dogfood loop

The intent was for Sigil's own `react-shadcn` renderer (task FND-4 /
CW-20260515-0130) to generate these screens against the kit. At the time this
was built, FND-4 was complete on the `fnd-4-sysop-ui-integration` branch but
**not merged to `main`**, so this repo's renderer still emits the pre-kit
output. The screens here are therefore hand-composed directly against
`@hollis-labs/sysop-ui` — which is the kit's intended consumption pattern
(`OperationsTablePage`, `DataTable`, `DetailDialog`). Once FND-4 lands, these
pages can be regenerated from `.sigil/pages/*.yaml` definitions.

## Dependencies

- **`@hollis-labs/sysop-ui`** (`v0.4.0`) — the React kit + canonical theme,
  consumed as a pinned git dependency. For local kit development, link a
  working copy: `npm install file:<path-to>/libs/sysop-ui` from `frontend/`.
- **`github.com/hollis-labs/go-webui`** (`v0.1.0`) — the SPA-serving harness.
- **`github.com/chrispian/sigil`** — the parent module, wired via a `replace`
  to `../`; provides `internal/server` (the JSON API).

## License

MIT — see [LICENSE](./LICENSE).
