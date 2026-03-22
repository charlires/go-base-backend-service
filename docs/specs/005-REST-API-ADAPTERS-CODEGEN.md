# Implement REST API Adapters Code Generation

## Overview

The goal of this specification is to define a a way to generate REST API adapters in the Go-based backend service project. This will automate the creation of boilerplate code for HTTP handlers, request parsing, and response formatting based on an OpenAPI specification. By using code generation, we can ensure consistency across all API endpoints, reduce manual coding errors, and speed up development when adding new endpoints or modifying existing ones.

## Decisions

- **Code generation tool**: [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen)
- **Generation mode**: `strict-server` — enforces strict request/response types via a typed `StrictServerInterface`
- **Transport**: `net/http` (Go stdlib) — `oapi-codegen` has first-class support for it; Gorilla Mux is not supported by `oapi-codegen`
- **OpenAPI spec location**: `openapi.yaml` at the project root, written manually
- **Generated types**: `types.gen.go` is the canonical source for request/response types; domain types in `internal/core/domain/` are mapped to/from these in the handler implementation
- **Existing adapters**: The manually written adapters in `internal/adapters/input/` will be replaced by the generated adapters
- **Request validation**: Out of scope for this spec

## Code Generation Strategy

1. **OpenAPI Specification**: The API contract is defined in `openapi.yaml` at the project root. This file is the single source of truth for all HTTP endpoints, request/response schemas, and status codes. It must be kept up to date as the API evolves.

2. **Code Generation Tool**: [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen) reads `openapi.yaml` and generates Go code. It is configured via an `oapi-codegen.yaml` config file at the project root.

3. **Generated Code Structure**: Generated code is isolated in a `gen/` sub-package to make it unambiguous what is auto-generated vs. manually maintained. The handler implementation lives in the parent package and imports from `gen/`.

```
internal/adapters/input/rest/
    gen/
        server.gen.go      <- generated: StrictServerInterface + net/http router setup
        types.gen.go       <- generated: canonical request/response types
    handlers.go            <- manually written: implements StrictServerInterface, maps to/from domain types
```

4. **`oapi-codegen` Configuration**: A config file `oapi-codegen.yaml` at the project root controls what is generated:

```yaml
package: gen
generate:
  strict-server: true
  models: true
output: internal/adapters/input/rest/gen/server.gen.go
```

5. **Code Generation Command**: Run the following command to regenerate the adapter code after any change to `openapi.yaml`:

```sh
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi-codegen.yaml openapi.yaml
```

This should also be added as a `//go:generate` directive at the top of `internal/adapters/input/rest/handlers.go`:

```go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../oapi-codegen.yaml ../../../openapi.yaml
```

6. **Handler Implementation**: `handlers.go` implements the generated `StrictServerInterface`. Each method receives strongly-typed request structs and returns strongly-typed response structs as defined in `types.gen.go`. Mapping between generated types and `internal/core/domain/` types is done inside these handlers.

7. **Wiring in `main.go`**: The generated router setup (e.g., `Handler(...)` or `HandlerWithOptions(...)`) is used in `main.go` to register all routes, replacing the manual `http.HandleFunc` registrations.
