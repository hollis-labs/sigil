# Security policy

## Supported versions

Sigil is pre-release software with no tagged versions. Security fixes are made on `main`. There are no backports.

## Report a vulnerability

Do not include an exploit, token, credential or other sensitive material in a public issue.

Use GitHub's private vulnerability-reporting flow when the repository's Security tab offers it. If it is unavailable, contact a repository maintainer privately through a contact channel published on the Hollis Labs organization or maintainer profile. Include:

- the affected commit and operating system
- which command was involved (`sigil serve`, `sigil mcp serve`, `sigil generate`, `sigil_sysop`)
- how it was started, and whether the listener was reachable beyond the local machine
- reproduction steps and the security impact
- a safe way to contact you about coordination

Maintainers will acknowledge a private report, investigate it, and coordinate disclosure; response times are best effort.

## Deployment boundary

Sigil is a build-time compiler and a local development tool, not a hosted service. Its networked pieces are for one developer on one machine:

- **`sigil serve`** runs the dev server, preview and REST API (`/api/pages`, `/api/datasources`, and so on). It has **no authentication and no TLS**, and its write endpoints (`POST`/`PUT`) change the project's YAML on disk. It binds `127.0.0.1` by default. Do not bind it to another address unless something in front of it provides authentication and TLS.
- **`sigil_sysop`** (in `sysop/`) is an admin UI over a Sigil project with the same properties: no authentication, loopback by default. Its `-addr` flag can widen that (for example `:8080` listens on all interfaces); treat doing so as exposing file-write access to anyone who can reach the port.
- **`sigil mcp serve`** speaks MCP over stdio to the agent that spawned it and can create and update page configs. It runs with the invoking user's permissions; run it only for an agent you trust with the project directory.

## Generated code

Sigil emits source; it is not a sandbox and does not vet what a YAML config asks for. Treat a config from an untrusted source as untrusted input: review generated output before building or running it. Generated code does not depend on Sigil at runtime, so its security is the security of the generated app and its own dependencies.

`sigil generate --clean` deletes files in its output directory. Point it at a directory you intend it to own.

## Data at rest

Sigil keeps no database and stores no credentials. Project configs, themes and datasource definitions are plain files; do not put secrets in datasource or page YAML that you commit or generate into a public app.

## Current security limitations

- no authentication or TLS on the dev server, REST API or Sysop UI
- the MCP server has full write access to the project it serves
- generated output is not scanned or sandboxed
- pre-release contracts

These are deployment constraints, not hidden roadmap promises. Operate within them or place Sigil behind controls that provide the missing boundary.
