VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/chrispian/sigil/internal/cli.Version=$(VERSION)"

.PHONY: build test vet install clean

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
