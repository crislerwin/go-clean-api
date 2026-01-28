include .env
export $(shell sed 's/=.*//' .env)

GO ?= go
# Ajustado para refletir a estrutura do vídeo (pasta app)
APP_PATH := ./...
GOLANGCI_LINT ?= golangci-lint
PRE_COMMIT ?= pre-commit
GOIMPORTS ?= goimports
GOFMT ?= gofmt
TEST_DATABASE_URL ?= postgres://test_user:test_pass@localhost:5433/ticket_db_test?sslmode=disable

.DEFAULT_GOAL := help

.PHONY: help all fmt imports vet lint test test-e2e test-all tidy check pre-commit-install pre-commit-run ci build run dev migration-new migration-up migration-down migration-test-up coverage coverage-html swagger

help:
	@echo "Makefile targets:"
	@echo "  pre-commit-install      Install pre-commit git hook"
	@echo "  pre-commit-run          Run pre-commit hooks against all files"
	@echo "  fmt                     Run gofmt on $(APP_PATH)"
	@echo "  vet                     Run go vet on $(APP_PATH)"
	@echo "  lint                    Run golangci-lint using .golangci.yml"
	@echo "  test-unit               Run unit tests (skips e2e)"
	@echo "  test-e2e                Run end-to-end tests (requires test database)"
	@echo "  test-all                Run all tests (unit + e2e)"
	@echo "  coverage                Run tests and show coverage"
	@echo "  coverage-html           Run tests and show coverage report in browser"
	@echo "  build                   Build the application to bin/server"
	@echo "  run                     Run the application from bin/server/app"
	@echo "  dev                     Run the application with air hot reloading"
	@echo "  migration-new           Create a new migration file"
	@echo "  migration-up            Run migrations up"
	@echo "  migration-down          Run migrations down"
	@echo "  migration-test-up       Run migrations on test database"
	@echo "  swagger                 Generate Swagger documentation"
	@echo "  ci                      Full pipeline (tidy, fmt, lint, vet, test)"

pre-commit-install:
	@command -v $(PRE_COMMIT) >/dev/null 2>&1 || { echo "pre-commit not found"; exit 1; }
	$(PRE_COMMIT) install
	$(PRE_COMMIT) install-hooks

pre-commit-run:
	$(PRE_COMMIT) run --all-files

build: swagger
	$(GO) build -o bin/server/app cmd/api/main.go

run: build
	./bin/server/app

dev:
	air -c .air.toml

swagger:
	swag init -g cmd/api/main.go --parseDependency --parseInternal

migration-new:
	@read -p "Enter migration name: " name; \
	@goose -dir sql/migrations postgres "$(DATABASE_URL)" create $$name sql

migration-up:
	@goose -dir sql/migrations postgres "$(DATABASE_URL)" up

migration-down:
	@goose -dir sql/migrations postgres "$(DATABASE_URL)" down

migration-test-up:
	@echo "Running migrations on test database..."
	@goose -dir sql/migrations postgres "$(TEST_DATABASE_URL)" up

fmt:
	$(GO) fmt $(APP_PATH)

tidy:
	$(GO) mod tidy

vet:
	$(GO) vet $(APP_PATH)

lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || { echo "golangci-lint not found"; exit 1; }
	$(GOLANGCI_LINT) run --config .golangci.yaml $(APP_PATH)

test-unit:
	$(GO) test -v $(shell go list ./... | grep -v /test/e2e)

test-e2e:
	@echo "Running e2e tests with test database..."
	TEST_DATABASE_URL="$(TEST_DATABASE_URL)" $(GO) test -v ./test/e2e/...

test-all: test-unit test-e2e

coverage:
	$(GO) test -coverprofile=coverage.out $(APP_PATH)
	$(GO) tool cover -func=coverage.out

coverage-html: coverage
	$(GO) tool cover -html=coverage.out

check: lint vet test-unit

ci: tidy fmt lint vet test-unit
