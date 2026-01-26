GO ?= go
# Ajustado para refletir a estrutura do vídeo (pasta app)
APP_PATH := ./...
GOLANGCI_LINT ?= golangci-lint
PRE_COMMIT ?= pre-commit
GOIMPORTS ?= goimports
GOFMT ?= gofmt

.DEFAULT_GOAL := help

.PHONY: help all fmt imports vet lint test tidy check pre-commit-install pre-commit-run ci build run dev

help:
	@echo "Makefile targets:"
	@echo "  pre-commit-install      Install pre-commit git hook"
	@echo "  pre-commit-run          Run pre-commit hooks against all files"
	@echo "  fmt                     Run gofmt on $(APP_PATH)"
	@echo "  vet                     Run go vet on $(APP_PATH)"
	@echo "  lint                    Run golangci-lint using .golangci.yml"
	@echo "  test                    Run go test on $(APP_PATH)"
	@echo "  build                   Build the application to bin/server"
	@echo "  run                     Run the application from bin/server/app"
	@echo "  dev                     Run the application with air hot reloading"
	@echo "  ci                      Full pipeline (tidy, fmt, lint, vet, test)"

pre-commit-install:
	@command -v $(PRE_COMMIT) >/dev/null 2>&1 || { echo "pre-commit not found"; exit 1; }
	$(PRE_COMMIT) install
	$(PRE_COMMIT) install-hooks

pre-commit-run:
	$(PRE_COMMIT) run --all-files

build:
	$(GO) build -o bin/server/app cmd/api/main.go

run: build
	./bin/server/app

dev:
	air -c .air.toml

fmt:
	$(GO) fmt $(APP_PATH)

tidy:
	$(GO) mod tidy

vet:
	$(GO) vet $(APP_PATH)

lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint not found"; exit 1; }
	$(GOLANGCI_LINT) run --config .golangci.yml $(APP_PATH)

test:
	$(GO) test -v $(APP_PATH)

check: lint vet test

ci: tidy fmt lint vet test
