# Logging Strategy

## Overview

All logging is done using the `slog` structured logging library from the Go standard library. Logs are enriched with request-scoped fields and emitted at appropriate levels depending on the layer. A shared utility package handles logger propagation via `context.Context`.

## Rules

1. **Use `slog` for Structured Logging**: Use the `slog` package from the Go standard library for all logging throughout the application. This ensures consistent log formatting and easy integration with log management systems.

2. **Log Levels**: Only two log levels are used:
   - `Error` — for logging errors. Must only be emitted from the **HTTP adapter layer**, as it is the final layer responsible for handling and responding to errors.
   - `Debug` — for tracing execution flow. Must be emitted at **every function entry** in all layers (services, repositories, adapters), capturing relevant input parameters.

3. **Domain Structs Must Implement `LogValue()`**: All domain structs must implement the `LogValue() slog.Value` method to enable safe and consistent structured logging of domain entities. This method must **omit sensitive fields** (e.g., passwords, tokens, personal data) to prevent accidental exposure in logs.

4. **Avoid Logging Sensitive Information**: Sensitive information (e.g., passwords, personal data) must never appear in logs. Sensitive fields must be excluded directly in the `LogValue()` method of the relevant domain struct.

5. **Include Context in Logs**: When logging errors or important events, include relevant context (e.g., user ID, playlist ID) in the log messages to facilitate debugging and monitoring.

6. **Centralize Log Configuration**: Configure the logging output, format, and level once in `main.go` using `slog.SetDefault()` to ensure consistency across the application.

7. **Logger Injection via Context**: To propagate request-scoped fields (e.g., `request_id`, `method`, `path`) into logs emitted by all layers, use the dedicated logger utility package at `internal/pkg/logger/`. This package exposes two functions:
   - `logger.ContextWithLogger(ctx, l)` — stores a `*slog.Logger` enriched with request-scoped fields into the `context.Context`. Must be called once per request in the HTTP adapter, after enriching the logger with relevant request fields.
   - `logger.FromCtx(ctx)` — retrieves the enriched `*slog.Logger` from the context. Falls back to `slog.Default()` if no logger is present. Must be called in all layers (services, repositories) to obtain the logger for the current request.

   > Do **not** use a custom `context.Context` type or inject `*slog.Logger` directly into service/repository constructors. The former breaks compatibility with standard middleware and libraries; the latter prevents per-request field enrichment.
