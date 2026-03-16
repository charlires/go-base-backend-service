# Sentinel Error Strategy

## Overview

Sentinel errors are predefined, exported error variables that represent specific error conditions across the application. This strategy ensures consistent error handling, easier error comparison in tests, and clean separation of concerns between layers.

## Rules

1. **Define Sentinel Errors in a Central File**: All sentinel errors must be defined in `internal/core/errors.go`, regardless of which domain they belong to. They must follow a consistent naming convention (e.g., `ErrUserNotFound`, `ErrInvalidInput`). Sentinel errors cover all error categories: domain errors, repository errors, and input validation errors.

2. **Use Sentinel Errors in Services**: When a service encounters an error condition, it must return the corresponding sentinel error instead of creating a new error instance. This allows for consistent error handling and easier error comparison in tests and higher layers.

3. **Wrapping Errors with Context**: Always use `fmt.Errorf` with the `%w` verb to wrap a sentinel error whenever there is either:
   - (a) An underlying dependency error to preserve (e.g., a repository error), or
   - (b) Additional context required to fully understand the error condition (e.g., input validation failures, business rule violations).

   The sentinel error must always be included in the wrap so that `errors.Is` comparisons work correctly at higher layers.

4. **Always Return Service-Level Sentinel Errors**: Services must never propagate raw errors from dependencies (e.g., repository errors) directly to the caller. Every error returned from a service must be a sentinel error or wrapped with one using `fmt.Errorf`, ensuring that higher layers only deal with service-level error semantics.

5. **Error Comparison in Higher Layers**: When handling errors in higher layers (e.g., HTTP adapters), use `errors.Is` to check if the error matches a specific sentinel error. This allows for consistent error handling and response generation based on the type of error.
