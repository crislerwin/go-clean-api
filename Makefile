# Makefile for go-clean-api
#
# Provides common developer conveniences:
# - Formatting: `make fmt` and `make imports`
# - Static checks: `make vet`, `make lint`
# - Tests: `make test`
# - Module maintenance: `make tidy`
# - Pre-commit integration: `make pre-commit-install` and `make pre-commit-run`
# - Aggregates: `make check` (used by pre-commit hooks) and `make ci` (CI-friendly)
#
# Usage:
#   make help                # show targets
#   make pre-commit-install  # install pre-commit git hook and hooks
#   make pre-commit-run      # run pre-commit hooks against all files
#   make check               # run lint, vet, and tests
#   make ci                  # strict CI pipeline (fails on missing tools)
#
# Notes:
# - Hooks in .pre-commit-config.yaml use system-installed tools (pre-commit, goimports, golangci-lint).
# - Developers should run `make pre-commit-install` once after cloning to enable git hooks.
# - CI pipelines should run `make ci` or `make pre-commit-run` to ensure consistent checks.

GO ?= go
GOTESTPKGS := ./...
GOLANGCI_LINT ?= golangci-lint
PRE_COMMIT ?= pre-commit
GOIMPORTS ?= goimports
GOFMT ?= gofmt

.DEFAULT_GOAL := help

.PHONY: help all fmt imports vet lint test tidy check pre-commit-install pre-commit-run ci

help:
	@echo "Makefile targets:"
	@echo "  help                    Show this help"
	@echo "  pre-commit-install      Install pre-commit git hook and its hooks"
	@echo "  pre-commit-run          Run pre-commit hooks against all files"
	@echo "  fmt                     Run gofmt (formats code in-place)"
	@echo "  imports                 Run goimports (organizes imports in-place)"
	@echo "  tidy                    Run 'go mod tidy' to update go.mod and go.sum"
	@echo "  vet                     Run 'go vet' on the repository"
	@echo "  lint                    Run golangci-lint if installed"
	@echo "  test                    Run 'go test ./...' (unit tests)"
	@echo "  check                   Run lint, vet, and tests (best-effort)"
	@echo "  ci                      CI pipeline: fmt, imports, tidy, lint, vet, test"

# Install the pre-commit hook and install hooks defined in .pre-commit-config.yaml
pre-commit-install:
	@command -v $(PRE_COMMIT) >/dev/null 2>&1 || { echo "pre-commit not found; please install it (pipx/pip/homebrew)"; exit 1; }
	@echo "Installing pre-commit hooks..."
	@$(PRE_COMMIT) install -f --install-hooks
	@echo "pre-commit hooks installed."

# Run all pre-commit hooks against all files (useful in CI)
pre-commit-run:
	@command -v $(PRE_COMMIT) >/dev/null 2>&1 || { echo "pre-commit not found; please install it"; exit 1; }
	@echo "Running pre-commit hooks against all files..."
	@$(PRE_COMMIT) run --all-files

# Formatters
fmt:
	@command -v $(GOFMT) >/dev/null 2>&1 || { echo "gofmt not found in PATH; skipping fmt"; exit 1; }
	@echo "Running gofmt..."
	@$(GOFMT) -s -w .
	@echo "gofmt complete."

imports:
	@command -v $(GOIMPORTS)
