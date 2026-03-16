---
mode: 'agent'
description: 'Implement structured logging across all layers of the service following the logging strategy spec.'
---

Apply every rule below to the file(s) in scope. If a whole layer is in scope, apply to all files in that layer.

## Rules

### 1. Logger utility package
Ensure `internal/pkg/logger/logger.go` exists and exposes exactly these two functions:
- `ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context` — stores the logger in the context.
- `FromCtx(ctx context.Context) *slog.Logger` — retrieves it, falling back to `slog.Default()`.

Do **not** create this file if it already exists.

### 2. `LogValue()` on domain structs
Every domain struct in `internal/core/domain/` must implement `LogValue() slog.Value` using `slog.GroupValue(...)`.  
- Include all non-sensitive identifying fields (e.g. `id`, `name`, `title`, `artist`, `album`, `owner_id`).  
- **Omit** any sensitive fields (passwords, tokens, personal data).  
- Add `import "log/slog"` to the file if not already present.

### 3. HTTP adapter layer (`internal/adapters/input/`)
At the **top** of every handler function:
1. Build an enriched logger with `slog.Default().With("method", r.Method, "path", r.URL.Path, <entity_id_key>, <entity_id_value>)`.
2. Store it in the context with `logger.ContextWithLogger(r.Context(), ...)` and use that `ctx` for all downstream calls.
3. Emit a `Debug` log: `logger.FromCtx(ctx).Debug("<AdapterType>.<MethodName>", <input args as key-value pairs>)`.
4. On every error path, emit an `Error` log **before** writing the HTTP error response:  
   `logger.FromCtx(ctx).Error("<AdapterType>.<MethodName>: <description>", "error", err)`.

> `Error` logs must **only** be emitted in the HTTP adapter layer.

### 4. Service layer (`internal/core/services/`)
At the **top** of every method body (before any other logic), emit:
```go
logger.FromCtx(ctx).Debug("<ServiceType>.<MethodName>", <input args as key-value pairs>)
```
Use `logger` imported from `github.com/charlires/go-base-backend-service/internal/pkg/logger`.  
Do **not** add `Error` logs here.

### 5. Repository / output adapter layer (`internal/adapters/output/`)
At the **top** of every method body (before any other logic), emit:
```go
logger.FromCtx(ctx).Debug("<AdapterType>.<MethodName>", <input args as key-value pairs>)
```
Do **not** add `Error` logs here.

## Key constraints
- Never log sensitive fields (passwords, tokens, personal data).
- Never emit `Error` logs outside the HTTP adapter layer.
- Never inject `*slog.Logger` directly into service or repository constructors.
- Always use `logger.FromCtx(ctx)` in services and repositories to obtain the logger.
- Domain struct pointer receivers must be used for `LogValue()` (i.e. `func (x *Foo) LogValue() slog.Value`).
- When passing a domain struct to a log call, pass a pointer so `LogValue()` is picked up automatically by `slog`.
