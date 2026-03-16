# Architecture

This project follows **Hexagonal Architecture** (also known as Ports and Adapters). This document explains the structure, the reasoning behind naming decisions, and the rules that must be preserved as the project evolves.

---

## Why Hexagonal Architecture

The core goal is to keep business logic completely isolated from infrastructure concerns. The application should not care whether it is being called via HTTP, gRPC, a CLI, or a message consumer. It should not care whether data is stored in Postgres, Redis, or an in-memory store. Delivery mechanisms and infrastructure are implementation details — the core logic is not.

This is achieved by defining clear boundaries:

- The **core** contains all business logic and owns the contracts (interfaces) for everything it provides and everything it depends on.
- **Adapters** live outside the core and translate between the outside world and those contracts.
- The core never imports from adapters. Adapters always depend on the core, never the other way around.

---

## Project Structure

```
internal/
├── core/
│   ├── domain/        — entities, value objects, domain events
│   ├── services/      — business logic + input contracts (interfaces)
│   └── ports/         — output contracts (interfaces the core depends on)
│
└── adapters/
    ├── input/         — inbound adapters: HTTP, gRPC, CLI, message consumers
    └── output/        — outbound adapters: Postgres, Redis, SMTP, external APIs
```

---

## Layer Responsibilities

### `core/domain/`

Pure business objects. No framework dependencies, no database tags, no serialization concerns. This is the heart of the application.

### `core/services/`

Contains two things that live together intentionally:

1. **Input contracts** — Go interfaces that define what the service can do. These are what input adapters call into.
2. **Concrete implementations** — the structs that implement those interfaces, containing the actual business logic.

They live together because the interface and its implementation are tightly coupled by design. The interface describes the service's capability; the implementation delivers it. Separating them would add indirection without benefit.

```go
// core/services/user_service.go

// Input contract — consumed by input adapters
type UserService interface {
    CreateUser(ctx context.Context, cmd CreateUserCommand) (UserID, error)
    GetUser(ctx context.Context, query GetUserQuery) (User, error)
}

// Concrete implementation — injected by main.go
type userService struct {
    repo ports.UserRepository
}
```

> **Important:** input adapters must only import and reference the `UserService` **interface**, never the concrete struct. Go's implicit interface satisfaction means there is no need for the adapter to import the implementation — only the interface type.

### `core/ports/`

Output contracts only. These are the interfaces the core depends on to do its job — repositories, external service clients, notification senders, etc. They are defined here so the core owns them, and output adapters implement them from the outside.

```go
// core/ports/user_repository.go
type UserRepository interface {
    Save(ctx context.Context, user domain.User) error
    FindByID(ctx context.Context, id domain.UserID) (domain.User, error)
}
```

### `adapters/input/`

Concrete inbound adapters. Each subdirectory is a delivery mechanism: `http/`, `grpc/`, `cli/`, `messaging/`, etc. They hold a reference to a service interface from `core/services/` and translate incoming requests into application calls.

### `adapters/output/`

Concrete outbound adapters. Each subdirectory is an infrastructure concern: `postgres/`, `redis/`, `smtp/`, etc. They implement interfaces defined in `core/ports/`.

---

## Dependency Rules

These rules are non-negotiable. Violating them breaks the architecture.

| Package | May import | Must never import |
|---|---|---|
| `core/domain` | nothing | anything |
| `core/ports` | `core/domain` | `adapters/` |
| `core/services` | `core/domain`, `core/ports` | `adapters/` |
| `adapters/input` | `core/services` (interface only) | `adapters/output`, concrete service structs |
| `adapters/output` | `core/domain`, `core/ports` | `adapters/input`, `core/services` |
| `main.go` | everything | — |

`main.go` is the only place in the codebase that is allowed to know about all three layers simultaneously. It is the composition root — it instantiates output adapters, injects them into services, and hands services to input adapters.

```go
// main.go — the only place all layers meet
userRepo    := postgres.NewUserRepository(db)
userService := services.NewUserService(userRepo)
httpServer  := http.NewServer(userService)
```

---

## Naming Rationale

Hexagonal architecture uses terms like *driving/driven*, *primary/secondary*, and *incoming/outgoing* ports. These are accurate but feel like pattern jargon when embedded in folder names.

This project uses `input/output` instead:

- **`input`** — things that trigger the application (HTTP requests, gRPC calls, consumed messages)
- **`output`** — things the application reaches out to (databases, external APIs, email)

The architecture is enforced by structure and dependency rules, not by naming. Names should be readable to any Go developer regardless of familiarity with the pattern.

---

## Testing

This structure is designed to make unit testing straightforward.

Because input contracts (interfaces) live in `core/services/` alongside their implementations, mocks can be generated targeting that package directly:

```bash
mockgen -source=internal/core/services/user_service.go \
        -destination=internal/core/services/mocks/user_service_mock.go
```

Output contracts in `core/ports/` follow the same pattern. Tests for `core/services/` inject mock repositories and never touch a real database.

```go
// core/services/user_service_test.go
func TestCreateUser(t *testing.T) {
    mockRepo := mocks.NewMockUserRepository(t)
    mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

    svc := NewUserService(mockRepo)
    _, err := svc.CreateUser(ctx, cmd)
    assert.NoError(t, err)
}
```

Keep interfaces small and focused. Large interfaces produce heavy mocks and are a signal that a service is doing too much.

---

## Summary

| Concept | Lives in | Consumed by |
|---|---|---|
| Entities & value objects | `core/domain/` | everywhere in core |
| Input contracts (interfaces) | `core/services/` | `adapters/input/` |
| Business logic (implementations) | `core/services/` | injected via `main.go` |
| Output contracts (interfaces) | `core/ports/` | `core/services/` |
| Inbound adapters | `adapters/input/` | external world → core |
| Outbound adapters | `adapters/output/` | core → external world |
| Composition root | `main.go` | wires everything together |