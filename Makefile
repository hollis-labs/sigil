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
