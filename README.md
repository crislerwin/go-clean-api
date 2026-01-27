# Ticketing System API

A robust Ticketing System implementation showcasing **Clean Architecture**, **Domain-Driven Design (DDD)**, and **Test-Driven Development (TDD)** in Go.

## 🚀 Overview

This project implements a high-concurrency ticket purchasing system where data consistency is paramount. It serves as a reference for handling complex business rules (e.g., race conditions, atomic transactions) within a maintainable, modular codebase.

**Key Concepts:**

- **Clean Architecture**: Decoupled layers (Domain, UseCase, Infrastructure, Presentation) ensuring testability.
- **DDD (Domain-Driven Design)**: Rich domain models encapsulate business logic (`Order`, `Ticket`).
- **TDD (Test-Driven Development)**: Comprehensive test coverage for all layers.
- **RBAC (Role-Based Access Control)**: Secure access tailored for `User` and `Admin` roles.

## 🛠️ Tech Stack

- **Language**: Go 1.25+
- **Framework**: [Gin](https://github.com/gin-gonic/gin) (HTTP Web Framework)
- **Database**: PostgreSQL (Driver: `pgx`, Helper: `sqlx`)
- **Authentication**: JWT (`golang-jwt/jwt/v5`)
- **Migrations**: Goose
- **Testing**: Testify, SQLMock
- **Utilities**: Google UUID, Bcrypt

## 📂 Project Structure

```text
.
├── cmd/api         # Application entrypoint
├── internal
│   ├── domain      # Enterprise Business Rules (Entities)
│   ├── usecase     # Application Business Rules (Interactors)
│   └── infra       # Frameworks & Drivers
│       ├── database    # DB Connection & Transactions
│       ├── http        # Handlers, Middleware, Routes
│       └── repository  # Data Access Implementation
├── sql/migrations  # Database schema migrations
└── Makefile        # Automation commands
```

## ⚡ Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Make

### Quick Start

1.  **Start Services**:

    ```bash
    make run
    # OR
    docker-compose up -d
    ```

2.  **Run Migrations**:

    ```bash
    make migration-up
    ```

3.  **Run Tests**:
    ```bash
    make test
    ```

## 📐 Architecture Diagrams

### Domain Model (Class Diagram)

```mermaid
classDiagram
    class User {
        +UUID ID
        +String Name
        +String Email
        +String Password
        +String Role
    }

    class Event {
        +UUID ID
        +String Name
        +String Location
        +String Organization
        +String Rating
        +DateTime Date
        +String ImageURL
        +Int Capacity
        +Decimal Price
    }

    class Order {
        +UUID ID
        +UUID EventID
        +UUID UserID
        +Decimal TotalAmount
        +Int Quantity
        +String Status
        +DateTime CreatedAt
    }

    class Ticket {
        +UUID ID
        +UUID EventID
        +UUID OrderID
        +Decimal Price
        +String Status
    }

    User "1" -- "N" Order : places
    Event "1" -- "N" Ticket : defines
    Order "1" -- "N" Ticket : contains
    Ticket "1" -- "1" Event : belongs_to
```

### Purchase Flow (Sequence Diagram)

This diagram illustrates the critical section handling for ticket purchasing to prevent race conditions.

```mermaid
sequenceDiagram
    participant U as User
    participant API as Ticket API
    participant DB as Postgres

    Note over U, API: User selects 2 tickets for Event X
    U->>API: POST /orders {eventId, qty: 2}

    rect rgb(240, 240, 240)
        Note right of API: Transaction Start (Critical Section)
        API->>DB: BEGIN TRANSACTION

        API->>DB: SELECT capacity ... FOR UPDATE
        DB-->>API: Row Lock (Event Info)

        API->>API: Calculate Available Capacity

        alt Sufficient Stock
            API->>DB: INSERT INTO orders (PENDING)
            API->>DB: INSERT INTO tickets (x2)
            API->>DB: COMMIT
            API-->>U: 201 Created
        else Insufficient Stock
            API->>DB: ROLLBACK
            API-->>U: 409 Conflict (Sold Out)
        end
    end
```

## 🔌 API Routes

**Swagger Documentation**: [http://localhost:8080/swagger](http://localhost:8080/swagger)

### Public

- `POST /api/v1/signup`: Register a new user.
  - **Body**: `{"name": "...", "email": "...", "password": "..."}`
  - **Default Role**: `user`
- `POST /api/v1/login`: various Authenticate and receive JWT.
  - **Body**: `{"email": "...", "password": "..."}`

### Protected (Authenticated)

- `POST /api/v1/orders`: Purchase tickets.
  - **Requires**: Role `user` or `admin`.
  - **Body**: `{"event_id": "...", "quantity": 2}`

### Admin Only

- `POST /api/v1/events`: Create a new event.
  - **Requires**: Role `admin`.
  - **Body**: `{"name": "...", "capacity": 100, "price": 50.0, ...}`

## TODO

- [ ] Add e2e tests
- [ ] Add edit user endpoint
- [ ] Add edit event endpoint
- [ ] Add event description field
- [ ] Add soft delete event endpoint
- [ ] Add change order status webhook endpoint for payment gateway
- [ ] Add event image upload endpoint // using cloud storage
- [ ] Add event check-in endpoint // using ticket id
