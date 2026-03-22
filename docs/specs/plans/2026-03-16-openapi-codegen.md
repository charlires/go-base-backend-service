# Research Document: OpenAPI Code Generation for Go REST API Adapters

**Date:** 2026-03-16
**Status:** Draft
**Lead:** @charlires

---

## 1. Objective
We need to integrate a library into our existing Go project that automatically generates HTTP routing handlers and data models from an existing OpenAPI v3 YAML specification file. The goal is to enforce contract-first development and reduce or eliminate manual boilerplate.

## 2. Context & Constraints
* **Language/Framework:** Go (v1.21+)
* **Input:** OpenAPI v3+ YAML file
* **Output Required:** Go structs (models), interfaces for handlers, and ideally integration with standard `net/http`.
* **Constraints:** The generated code should not force a heavy web framework dependency. It must support strict typing.

## 3. Options Evaluated
*List the libraries or approaches discovered during research.*

### Option A: `oapi-codegen/oapi-codegen`
* **Description:** The most popular Go code generator for OpenAPI 3.
* **Pros:** Highly maintained, supports multiple routers (Chi, Echo, Gin, Fiber, strict `net/http`), generates both client and server code, strict typing support, team is already familiar with it.
* **Cons:** Can generate very large files; complex configuration options can be overwhelming at first.

### Option B: `go-swagger/go-swagger`
* **Description:** A complete toolkit for OpenAPI.
* **Pros:** Extremely mature, heavily tested.
* **Cons:** Primarily built for Swagger 2.0 (OpenAPI 2). Support for OpenAPI 3 is limited/bolted-on. Imposes its own heavy ecosystem.

### Option C: `ogen-go/ogen`
* **Description:** A newer, high-performance OpenAPI v3 code generator.
* **Pros:** Excellent performance, very strict type safety, zero reflection used in generated code.
* **Cons:** Steeper learning curve, less community adoption compared to `oapi-codegen`, custom internal routing engine instead of standard ones.

## 4. Recommendation & Rationale
*State the chosen path and why it wins.*
**Decision:** Proceed with **`oapi-codegen/oapi-codegen`**. 
**Why:** It offers the best balance of community support, OpenAPI v3 compatibility, and router flexibility. Unlike `go-swagger`, it natively handles v3, and unlike `ogen`, it easily plugs into standard Go routers like Chi or `net/http` without enforcing its own routing paradigm. 

## 5. Proof of Concept (PoC) / Implementation Notes
*Provide a minimal snippet or the CLI command required to prove it works. This feeds directly into the "Plan" phase.*
**Tool Installation:**
`go install github.com/oapi-codegen/oapi-codegen/cmd/oapi-codegen@latest`

**Generation Command (Example):**
`oapi-codegen -generate types,server,spec -package api openapi.yaml > api/generated.go`

## 6. Open Questions / Risks
*What do we still need to figure out during the Plan phase?*
* How will we handle authentication middleware with the generated handlers?
* Should we commit the generated `generated.go` file to version control, or generate it dynamically in the CI/CD pipeline?