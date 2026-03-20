# Structured Logging Middleware for REST API Adapter

## Summary

The goal is to implement structured, request-scoped logging for the existing REST API using `log/slog` (Go standard library). Logs are written to stdout in JSON format. There is no log export or OpenTelemetry SDK involved at this stage.

## Research

The existing `internal/pkg/logger/logger.go` already provides two functions built on `log/slog`:

- `ContextWithLogger(ctx, logger)` — stores a `*slog.Logger` into the context.
- `FromCtx(ctx)` — retrieves the `*slog.Logger` from the context, falling back to `slog.Default()`.

This context-logger pattern is already in place and will be preserved. The goal is to make logging **transparent to handlers**: no handler should need to manually extract the logger, enrich it with request fields, or store it back into the context.

The generated `StrictMiddlewareFunc` type (from `oapi-codegen`) has the following signature:

```go
// type alias from strictnethttp
type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

// where StrictHandlerFunc is:
type StrictHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error)
```

The middleware receives the `operationID` string (e.g. `"GetUserById"`) at the outer call, and a fresh `context.Context` at the inner call — making it straightforward to inject an enriched logger before calling the next handler.

## Implementation Details

### 1. Logger configuration in `main.go`

A `LogConfig` const/var block will be defined in `main.go` to configure the logger before it is passed to the middleware:

```go
// Example — values are hardcoded constants for now
const (
    logLevel  = slog.LevelDebug
    logFormat = "json" // "json" | "text"
)
```

The `slog.Default()` logger is initialized in `main.go` and the same `*slog.Logger` instance is passed when wiring the middleware.

### 2. Middleware location

A new file `internal/adapters/input/rest/logging_middleware.go` will contain the `StrictMiddlewareFunc` implementation.

### 3. Fields injected by the middleware

The middleware enriches the logger with the following fields on every request before injecting it into the context:

| Field          | Source                                      |
|----------------|---------------------------------------------|
| `operation_id` | `operationID` argument of `StrictMiddlewareFunc` |
| `request_id`   | `X-Request-ID` HTTP header (generated with `uuid` if absent) |
| `method`       | `r.Method`                                  |
| `path`         | `r.URL.Path`                                |

> **Out of scope:** `path_vars` (e.g. `playlistId`, `trackId`) will not be extracted by the middleware. They remain available to handlers via the typed request object.

### 4. Context injection

After enriching, the middleware calls `logger.ContextWithLogger(ctx, enrichedLogger)` and passes the resulting context to the next `StrictHandlerFunc`. All downstream layers (services, repositories) retrieve the logger via `logger.FromCtx(ctx)`.

### 5. Error logging strategy

Handlers **never write error responses themselves**. On any service error they return `nil, err`, propagating the sentinel error up. All error-to-response mapping, logging, and writing is centralised in `ResponseErrorHandlerFunc`.

#### 5a. Handler contract

```go
func (h *Handlers) GetUserById(ctx context.Context, request gen.GetUserByIdRequestObject) (gen.GetUserByIdResponseObject, error) {
    user, err := h.UserService.GetUserByID(ctx, request.UserId)
    if err != nil {
        return nil, err  // propagate — no logging, no typed error response here
    }
    return gen.GetUserById200JSONResponse(domainUserToResponse(user)), nil
}
```

#### 5b. `ResponseErrorHandlerFunc` — central error mapper

`StrictHTTPServerOptions.ResponseErrorHandlerFunc` fires whenever a handler returns a non-nil Go `error`. It is the single place that:

1. Uses `errors.Is` against sentinel errors from `internal/core/errors.go` to decide the status code and log level.
2. Logs the error using the request-scoped logger already enriched by the middleware (`logger.FromCtx(r.Context())`).
3. Writes the **typed generated response** (e.g. `gen.GetUserById404JSONResponse`) — not a plain `http.Error` — so the response body stays consistent with the OpenAPI spec.

```
ErrNotFound      → 404 · Warn  · gen.XxxNotFoundJSONResponse
ErrInvalidInput  → 400 · Warn  · gen.XxxBadRequestJSONResponse  (if defined)
anything else    → 500 · Error · gen.Xxx500JSONResponse
```

> **Note:** because `ResponseErrorHandlerFunc` does not receive the operation name, the response type cannot be inferred generically. A helper such as `rest.ErrorResponse(w, err)` will handle the mapping and writing, keeping the func body small.

#### 5c. `RequestErrorHandlerFunc`

Fires when the request body cannot be decoded (bad JSON). Logged at `Warn` level and responds with `400`.

### 6. Prerequisites — sentinel errors in `internal/core/errors.go`

The sentinel errors are currently defined per-service (e.g. `ErrUserIDRequired` in `user_service.go`). Before implementing this spec, shared sentinels must be centralised in `internal/core/errors.go` per the rules in `003-SENTINEL-ERROR-STRATEGY.md`:

```go
var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

Service-level errors (e.g. `ErrUserIDRequired`) wrap these:

```go
var ErrUserIDRequired = fmt.Errorf("user id is required: %w", core.ErrInvalidInput)
```

This ensures `errors.Is(err, core.ErrNotFound)` works correctly in the response error handler.

### 7. Wiring in `main.go`

`gen.NewStrictHandlerWithOptions` is used instead of `NewStrictHandler`:

```go
loggingMiddleware := rest.NewLoggingMiddleware(appLogger)

gen.NewStrictHandlerWithOptions(handlers, []gen.StrictMiddlewareFunc{
    loggingMiddleware,
}, gen.StrictHTTPServerOptions{
    RequestErrorHandlerFunc:  rest.NewRequestErrorHandler(),
    ResponseErrorHandlerFunc: rest.NewResponseErrorHandler(),
})
```

`NewRequestErrorHandler` and `NewResponseErrorHandler` are defined in `logging_middleware.go`. They retrieve the logger from the context (already enriched by the middleware) and centralise all error-to-response mapping.