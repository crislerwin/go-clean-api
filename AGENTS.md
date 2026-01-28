# AGENTS Guidelines for This Repository

This repository is developed with the assistance of AI agents. These guidelines ensure consistency, quality, and a shared mental model of the system.

## Core Philosophy

We adhere to the "Theory Building" view of programming (Peter Naur). The goal is not just working code, but a clear, shared mental model ("The Theory") of the system.
**Source Code is Design.** You are not just a coder; you are a designer.

## Rules of Engagement

### 1. Protect the "Theory"

- **Explain the 'Why':** When generating complex logic, you must explicitly state the "theory" or reasoning behind the design choice in the comments or PR description.
- **Mental Model Matching:** Code structure must match the domain language. Use explicit naming conventions that reveal the intent, not just the mechanism (e.g., `calculateNetWorth()` vs `sumArray()`).

### 2. Optimize for "Cost of Change" (Constantine’s Equivalence)

- **Minimize Coupling:** Prefer code that can be understood in isolation. Avoid creating "spooky action at a distance" where changing one module breaks an unrelated one.
- **Explicit > Implicit:** Do not use "clever" one-liners if they obscure the mental model. Readability for the human reviewer is higher priority than brevity.

### 3. The "Prompt is Source" Protocol

- If I provide a high-level specification (Natural Language), treat that as the "immutable source."
- If the implementation fails, **do not just patch the code**. Analyze if the "Theory" (the prompt/spec) was ambiguous. Ask for clarification on the spec before "hacking" a fix.

### 4. Code Review Readiness

- Assume I need to reconstruct the mental model from your code immediately.
- Use **Extract Method** aggressively to label blocks of logic with domain-centric names.
- If you introduce a new pattern, explain it first.

## Development Standards

### Test-Driven Development (TDD)

Follow the **Red-Green-Refactor** cycle:

1. **Red**: Write a failing test that defines the desired behavior.
2. **Green**: Write minimal code to make the test pass.
3. **Refactor**: Improve code quality while keeping tests green.

**Testing Strategy**:

- **Unit Tests**: Test individual components/entities in isolation.
- **Integration Tests**: Test interactions between layers (UseCase + Repository, etc.).
- **E2E Tests**: Test complete user flows.
- **Table-Driven Tests**: ALWAYS prefer this format for better readability and coverage.

### Architecture & Design

- **Domain-Driven Design (DDD)**: Respect clean architecture boundaries (Domain > UseCase > Interface Adapters > Infrastructure).
- **SOLID Principles**: Apply dependencies inversion and single responsibility strictly.

### Clean Code Standards

- **Functions**: Small, focused (SRP).
- **Naming**: Intention-revealing.
- **Comments**: Explanation of _why_, not _what_.

## Workflow

1.  [ ] **Analyze**: Understand the requirement and mental model.
2.  [ ] **Plan**: Propose a plan and explain tradeoffs.
3.  [ ] **Test**: Write the test case (Table-Driven).
4.  [ ] **Implement**: Write the code to satisfy the test.
5.  [ ] **Verify**: Run tests and ensure all pass.

## Communication Templates

### When Starting a Task

```
I understand the "Theory" for [requirement] is...
Here's my approach:
...
Tradeoffs:
...
```

### When Implementing

```
I'm creating [component] to support [domain concept].
Key design decisions:
...
```

## Project Maintenance & Commands

### 1. Development

- **`make dev`**: Starts application with hot-reloading (Air).

### 2. Dependencies

- **`make tidy`**: Updates `go.mod` and `go.sum`.

### 3. Testing

- **`make test-unit`**: Fast unit tests (excludes e2e).
- **`make test-e2e`**: Full flow verification (requires DB).
- **`make ci`**: Full validation pipeline.

### 4. Database Migrations

**IMPORTANT**: Never create migration files manually. Always use goose via Makefile.

- **`make migration-new`**: Create a new migration file (will prompt for name).
- **`make migration-up`**: Apply pending migrations to development database.
- **`make migration-down`**: Rollback last migration.
- **`make migration-test-up`**: Apply migrations to test database.

### 5. Commands Recap

| Command               | Purpose                         |
| --------------------- | ------------------------------- |
| `make dev`            | Start dev server with HMR.      |
| `make test-unit`      | Run unit tests only.            |
| `make test-e2e`       | Run e2e tests.                  |
| `make ci`             | Run full pipeline.              |
| `make lint`           | Run linters.                    |
| `make migration-new`  | Create new migration (prompts). |
| `make migration-up`   | Apply pending migrations.       |
| `make migration-down` | Rollback last migration.        |
