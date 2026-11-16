GO ?= $(HOME)/.local/go/bin/go
export PATH := $(HOME)/.local/go/bin:$(PATH)

.PHONY: run build test vet fmt lint contract contract-express

run:
	WEB_DIR=../dist $(GO) run ./cmd/server

build:
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...
	$(GO) vet -tags=contract ./contract/

fmt:
	gofmt -s -w .

lint: vet
	@test -z "$$(gofmt -s -l .)" || (echo "gofmt needed:"; gofmt -s -l .; exit 1)

contract:
	./scripts/contract.sh go

contract-express:
	./scripts/contract.sh express
