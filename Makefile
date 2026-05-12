VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/chrispian/sigil/internal/cli.Version=$(VERSION)"

.PHONY: build test vet install clean demo demo-generate demo-dev se-demo-generate se-demo-dev se-demo

build:
	go build $(LDFLAGS) -o bin/sigil ./cmd/sigil

test:
	go test ./...

vet:
	go vet ./...

install:
	go install $(LDFLAGS) ./cmd/sigil

clean:
	rm -rf bin/

# ── Demo ────────────────────────────────────────────────
demo-generate: build
	bin/sigil generate --target react-shadcn --output /tmp/sigil-demo-gen
	rsync -a /tmp/sigil-demo-gen/types/ demo/src/types/
	rsync -a --ignore-existing /tmp/sigil-demo-gen/hooks/ demo/src/hooks/
	rsync -a /tmp/sigil-demo-gen/components/data-table.tsx demo/src/components/data-table.tsx
	@for f in /tmp/sigil-demo-gen/pages/*.tsx; do \
		page=$$(basename "$$f" .tsx); \
		case "$$page" in \
			forge-*) \
				route=$${page#forge-}; \
				mkdir -p "demo/src/app/(forge)/$${route}"; \
				rsync -a "$$f" "demo/src/app/(forge)/$${route}/page.tsx"; \
				;; \
			*) \
				mkdir -p "demo/src/app/$${page}"; \
				rsync -a "$$f" "demo/src/app/$${page}/page.tsx"; \
				;; \
		esac; \
	done
	@echo "✓ Demo pages synced"

demo-dev: demo-generate
	cd demo && npm run dev

demo: demo-dev

# ── Stack Explorer Demo ────────────────────────────────────
SE_ENTITIES := repo tag snapshot dimension lens scorecard dimensionscore pattern finding comparisonset reportconfig

se-demo-generate: build
	bin/sigil generate --target react-shadcn --theme stack-explorer \
		--pages se-dashboard,se-repos,se-repo-detail,se-scorecards,se-gap-analysis,se-reports,se-patterns,se-findings,se-dimensions,se-settings \
		--output /tmp/sigil-se-gen
	@for e in $(SE_ENTITIES); do \
		cp /tmp/sigil-se-gen/types/$$e.ts stack-explorer-demo/src/types/$$e.ts 2>/dev/null || true; \
	done
	@for e in $(SE_ENTITIES); do \
		[ -f "stack-explorer-demo/src/hooks/use-$$e.ts" ] || cp /tmp/sigil-se-gen/hooks/use-$$e.ts stack-explorer-demo/src/hooks/use-$$e.ts 2>/dev/null || true; \
	done
	rsync -a --ignore-existing /tmp/sigil-se-gen/components/data-table.tsx stack-explorer-demo/src/components/data-table.tsx
	@# Sync all custom component source files (always overwrite)
	@for f in /tmp/sigil-se-gen/components/*.tsx; do \
		base=$$(basename "$$f"); \
		[ "$$base" = "data-table.tsx" ] || [ "$$base" = "index.ts" ] || \
		cp "$$f" "stack-explorer-demo/src/components/$$base"; \
	done
	@# Sync generated layout and API client
	@if [ -f /tmp/sigil-se-gen/app/\(explorer\)/layout.tsx ]; then \
		cp /tmp/sigil-se-gen/app/\(explorer\)/layout.tsx stack-explorer-demo/src/app/\(explorer\)/layout.tsx; \
	fi
	@if [ -f /tmp/sigil-se-gen/lib/api.ts ]; then \
		mkdir -p stack-explorer-demo/src/lib; \
		cp /tmp/sigil-se-gen/lib/api.ts stack-explorer-demo/src/lib/api.ts; \
	fi
	@for f in /tmp/sigil-se-gen/lib/*-context.tsx; do \
		[ -f "$$f" ] && cp "$$f" stack-explorer-demo/src/lib/ || true; \
	done
	@for f in /tmp/sigil-se-gen/pages/*.tsx; do \
		page=$$(basename "$$f" .tsx); \
		route=$${page#se-}; \
		case "$$route" in \
			repo-detail) \
				mkdir -p "stack-explorer-demo/src/app/(explorer)/repo-detail/[id]"; \
				rsync -a "$$f" "stack-explorer-demo/src/app/(explorer)/repo-detail/[id]/page.tsx"; \
				sed -i '' 's/params\.filter as string/params.id as string/g' "stack-explorer-demo/src/app/(explorer)/repo-detail/[id]/page.tsx"; \
				;; \
			*) \
				mkdir -p "stack-explorer-demo/src/app/(explorer)/$${route}"; \
				rsync -a "$$f" "stack-explorer-demo/src/app/(explorer)/$${route}/page.tsx"; \
				;; \
		esac; \
	done
	@echo "✓ Stack Explorer pages synced"

se-demo-dev: se-demo-generate
	cd stack-explorer-demo && npm run dev -- -p 3334

se-demo: se-demo-dev

# ── Clockwork Demo ────────────────────────────────────────
# Sprint 10 phase 3: regenerated through `react-shadcn` with `--target-mode spa`
# (Vite SPA output). The fork-mode (`react-clockwork`) is deleted in phase 4.
# Uses tmp + rsync (matches se-demo-generate) so hand-written mock hooks under
# demo-clockwork/src/hooks/ are protected by --ignore-existing.
cw-demo-generate: build
	bin/sigil generate --target react-shadcn --target-mode spa --theme clockwork-dark \
		--pages clockwork-board,clockwork-task-detail,clockwork-collections,clockwork-dashboard,clockwork-templates,clockwork-checkpoints,clockwork-runs,clockwork-models,clockwork-projects,clockwork-sprints,clockwork-settings,clockwork-plans,clockwork-plan-detail \
		--output /tmp/sigil-cw-gen
	@# Stale fork-output cleanup: the active-runs provider migrated to a custom
	@# provider declaration in Phase 2, so the renderer now copies the source
	@# into lib/. Remove the pre-Phase-3 copies under hooks/ so imports resolve
	@# to the generated lib/ versions instead.
	@rm -f demo-clockwork/src/hooks/active-runs-context.tsx \
	       demo-clockwork/src/hooks/sse-active-runs-bridge.tsx
	@# Note: demo-clockwork/src/lib/api-context.tsx and clockwork-api.ts are
	@# hand-written companion files for the existing mock hooks; they survive
	@# regeneration. The newer lib/api.ts (from Sigil) sits alongside them.
	@# Types: regenerate (overwrite).
	@# Types: regenerate (overwrite). Tolerate missing dir — the clockwork
	@# demo's mock hooks declare their own types inline, and Sigil's
	@# phase-3.5 datasource filter skips unreferenced manifests so no
	@# unrelated types are emitted here.
	@if [ -d /tmp/sigil-cw-gen/types ]; then \
		rsync -a /tmp/sigil-cw-gen/types/ demo-clockwork/src/types/; \
	fi
	@# Hooks: preserve hand-written mock hooks (use-*.ts) and the generic SSE
	@# primitive (use-sse.ts). --ignore-existing means generated hooks only fill
	@# gaps; hand-edited files are never clobbered. Tolerate missing dir.
	@if [ -d /tmp/sigil-cw-gen/hooks ]; then \
		rsync -a --ignore-existing /tmp/sigil-cw-gen/hooks/ demo-clockwork/src/hooks/; \
	fi
	@# Lib: regenerate (overwrite) — includes api.ts, utils.ts, and any provider
	@# files copied from .sigil/providers/ (active-runs-context.tsx,
	@# sse-active-runs-bridge.tsx).
	rsync -a /tmp/sigil-cw-gen/lib/ demo-clockwork/src/lib/
	@# Components: regenerate all (overwrite) — Sigil owns these.
	rsync -a /tmp/sigil-cw-gen/components/ demo-clockwork/src/components/
	@# App.tsx + routes.tsx: regenerate (overwrite).
	rsync -a /tmp/sigil-cw-gen/App.tsx demo-clockwork/src/App.tsx
	rsync -a /tmp/sigil-cw-gen/routes.tsx demo-clockwork/src/routes.tsx
	@# Theme assets at project root.
	rsync -a /tmp/sigil-cw-gen/globals.css demo-clockwork/src/globals.css
	rsync -a /tmp/sigil-cw-gen/tailwind.config.ts demo-clockwork/tailwind.config.ts
	@# Pages: regenerate (overwrite). All clockwork-* pages land directly under src/pages/.
	rsync -a /tmp/sigil-cw-gen/pages/ demo-clockwork/src/pages/
	@echo "✓ Clockwork pages generated (SPA via react-shadcn)"

cw-demo-install:
	cd demo-clockwork && npm install

cw-demo-dev: cw-demo-generate
	cd demo-clockwork && npm run dev

cw-demo: cw-demo-generate cw-demo-dev
