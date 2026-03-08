---
intent: pcc_global
project: sigil
updated_at: "2026-03-04"
---

# Workflows

## Build and install

```bash
# Build from source
make build

# Install locally
make install

# Or via Go install
go install github.com/chrispian/sigil/cmd/sigil@latest
```

## Project setup

```bash
# Initialize a new project
sigil init --name my-app

# Create a page
sigil new page dashboard --title "Dashboard"

# Validate all configs
sigil validate

# Check environment health
sigil doctor
```

## Code generation

```bash
# Generate Go/Templ code
sigil generate --target go-templ --output internal/ui

# Generate React/shadcn code
sigil generate --target react-shadcn --output src/generated
```

## Preview and dev server

```bash
# Preview a page in browser (mock data)
sigil preview dashboard

# Start live dev server (auto-reload on config changes)
sigil serve
```

## Config management

```bash
# List pages, components, or themes
sigil list pages
sigil list components
sigil list themes

# Export pages as JSON
sigil export

# Import pages from JSON
sigil import

# Compare configs semantically
sigil diff page-a.yaml page-b.yaml

# Migrate configs to latest schema
sigil migrate

# Generate JSON Schema for IDE support
sigil schema export
```

## MCP server (for AI agents)

```bash
sigil mcp serve
```

Provides 9 tools, 7 resource types, 2 prompts.

## Testing

```bash
make test    # run all tests
make vet     # go vet
```

## Evidence
- Last refreshed: 2026-03-04 (mentat PCC bootstrap)
- Sources: README.md, CLAUDE.md, Makefile, docs/04_cli-reference.md
