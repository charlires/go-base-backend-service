# Implement REST API Adapters Code Generation

## Overview

This document outlines the approach for implementing code generation for REST API adapters in the Go-based backend service project. The goal is to automate the creation of boilerplate code for HTTP handlers, HTTP request parsing, and HTTP response formatting based on an OpenAPI specification.

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

2. **Code Generation Tool**: [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen) reads `openapi.yaml` and generates Go code. It is configured via two config files at the project root — one for types, one for the server — so that the generated output is split into separate files.

3. **Generated Code Structure**: Generated code is isolated in a `gen/` sub-package to make it unambiguous what is auto-generated vs. manually maintained. The handler implementation lives in the parent package and imports from `gen/`.

```
internal/adapters/input/rest/
    gen/
        server.gen.go      <- generated: StrictServerInterface + net/http router setup
        types.gen.go       <- generated: canonical request/response types
    handlers.go            <- manually written: implements StrictServerInterface, maps to/from domain types
```

4. **`oapi-codegen` Configuration**: Two config files at the project root control what is generated:

`oapi-codegen-types.yaml` — generates only the model types into `types.gen.go`:
```yaml
package: gen
generate:
  models: true
output: internal/adapters/input/rest/gen/types.gen.go
```

`oapi-codegen.yaml` — generates the server interface and router setup into `server.gen.go` (models excluded to avoid duplication):
```yaml
package: gen
generate:
  std-http-server: true
  strict-server: true
  models: false
output: internal/adapters/input/rest/gen/server.gen.go
```

5. **Code Generation Command**: The `//go:generate` directives at the top of `internal/adapters/input/rest/handlers.go` run both configs in order (types first, then server):

```go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../oapi-codegen-types.yaml ../../../../openapi.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../oapi-codegen.yaml ../../../../openapi.yaml
```

To regenerate after any change to `openapi.yaml`, run:

```sh
go generate ./internal/adapters/input/rest/...
```

6. **Handler Implementation**: `handlers.go` implements the generated `StrictServerInterface`. Each method receives a strongly-typed `*RequestObject` and returns a strongly-typed `*ResponseObject` as defined in `server.gen.go`. The domain model types from `types.gen.go` are used for mapping between generated types and `internal/core/domain/` types inside these handlers.

7. **Wiring in `main.go`**: The `Handlers` struct is wrapped with `gen.NewStrictHandler` and passed to `gen.Handler(...)` to register all routes, replacing the manual `http.HandleFunc` registrations:

```go
gen.Handler(gen.NewStrictHandler(handlers, nil))
```
