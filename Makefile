VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/chrispian/sigil/internal/cli.Version=$(VERSION)"

.PHONY: build test vet install clean demo demo-generate demo-dev

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
