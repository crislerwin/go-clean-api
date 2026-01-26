# Go Clean API

A robust Ticketing System implementation showcasing **Clean Architecture**, **Domain-Driven Design (DDD)**, and **Test-Driven Development (TDD)** in Go.

## 🚀 Overview

This project implements a ticket purchasing flow where high concurrency and data consistency are critical. It acts as a reference implementation for handling complex business rules and race conditions (e.g., preventing overselling of tickets) within a clean, maintainable codebase.

**Key Concepts:**

- **Clean Architecture**: Strict separation of concerns to ensure testability and independence from frameworks.
- **DDD (Domain-Driven Design)**: Business logic is encapsulated in rich domain models (`Order`, `Ticket`).
- **TDD (Test-Driven Development)**: Development driven by tests to ensure correctness from the start.

## 🛠️ Tech Stack

- **Language**: Go 1.25+
- **Testing**: `testify` for assertions and suites
- **Linter**: `golangci-lint`
- **Utilities**: `google/uuid`

## 📂 Project Structure

```text
.
├── internal
│   ├── domain    # Enterprise Business Rules (Entities)
│   ├── usecase   # Application Business Rules
│   └── infra     # Frameworks & Drivers (Database, etc.)
├── Makefile      # Automation commands
├── go.mod        # Dependency definitions
└── README.md     # Documentation
```

## ⚡ Getting Started

### Prerequisites

- Go installed (1.25+)
- Make

### Verified Commands

You can use the included `Makefile` to run common tasks:

```bash
# Run all tests
make test

# Run linter
make lint

# Run code formatting
make fmt
```

## 📐 Architecture Diagrams

### Domain Model (Class Diagram)

```mermaid
classDiagram
    class User {
        +UUID id
        +String name
        +String email
    }

    class Event {
        +UUID id
        +String eventName
        +Int totalCapacity
        +Decimal basePrice
    }

    class Order {
        +UUID id
        +UUID userId
        +Status status
        +Decimal totalAmount
        +DateTime expiresAt
    }

    class Ticket {
        +UUID id
        +UUID orderId
        +UUID eventId
        +String ticketNumber
        +Decimal soldPrice
        +Status status
    }

    class PaymentMethod {
        +UUID id
        +UUID userId
        +String gatewayToken
        +String lastFourDigits
        +String brand
    }

    User "1" -- "N" Order : creates
    User "1" -- "N" PaymentMethod : owns
    Order "1" -- "N" Ticket : contains
    Order "1" -- "0..1" PaymentMethod : uses
    Event "1" -- "N" Ticket : defines
    PaymentMethod "1" -- "1" Order : processes
    Ticket "1" -- "1" Event : defines
    Ticket "1" -- "1" Order : belongs_to
    Order "1" -- "1" PaymentMethod : uses
```

### Purchase Flow (Sequence Diagram)

This diagram illustrates the critical section handling for ticket purchasing to prevent race conditions.

```mermaid
sequenceDiagram
    participant U as User
    participant API as Ticket API (Go)
    participant DB as Database (Postgres)

    Note over U, API: O usuário seleciona 2 ingressos para o Evento X
    U->>API: POST /orders {eventId, qty: 2}

    rect rgb(240, 240, 240)
        Note right of API: Início da Transação (Critical Section)
        API->>DB: BEGIN TRANSACTION

        %% A query mágica que valida o estoque real
        API->>DB: SELECT count(*) FROM tickets t <br/>JOIN orders o ON t.order_id = o.id <br/>WHERE t.event_id = X <br/>AND o.status IN ('PENDING', 'PAID')
        DB-->>API: return current_sold_count (ex: 98)

        API->>API: Check: (Capacity 100 - Sold 98) >= 2?

        alt Estoque Suficiente
            API->>DB: INSERT INTO orders (status='PENDING', expires_at=now()+15m) RETURNING id
            DB-->>API: order_id

            API->>DB: INSERT INTO tickets (order_id, event_id, price...) VALUES (x2)

            API->>DB: COMMIT
            API-->>U: 201 Created {orderId, status: 'PENDING', expiration: '15min'}
        else Estoque Insuficiente
            API->>DB: ROLLBACK
            API-->>U: 409 Conflict (Sold Out)
        end
    end
```
