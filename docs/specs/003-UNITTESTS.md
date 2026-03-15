# Implement Unittests

The goal of this document is to outline the approach for implementing unit tests in the Go-based backend service project. Unit tests are essential for ensuring the correctness of individual components and functions in isolation from external dependencies.

## Clarifications (Confirmed)

1. **Initial Scope**: We will start with unit tests for `user_service` only.
2. **Mock Strategy**: `mockery` is already set up and will be used for generated mocks.
3. **Mock Library**: `github.com/stretchr/testify/mock` is approved and will be used.
4. **Coverage Goal**: Focus on critical paths and edge cases (not a fixed percentage target at this stage).
5. **Documentation Update**: Test folder structure and naming conventions are part of this spec.

## Testing Strategy

1. **Core Services**: The core services layer will be the primary focus for unit testing, as it contains the business logic. We will use mock implementations of the output ports (repositories) to test the service methods without relying on actual database interactions.
2. **Adapters**: Output adapters will be mocked in the unit tests for the core services. Input adapters will be tested using integration tests, as they involve HTTP handling and require a running instance of the service.
3. **Phased Rollout**: Phase 1 covers `user_service`; additional services (`track_service`, `playlist_service`) will be added in later phases.
4. **Test Coverage**: We will prioritize critical paths and edge cases for each service under test.

## Tools and Libraries
- **Testing Framework**: We will use Go's built-in `testing` package for writing unit tests, leveraging table-driven tests for better organization and readability.
- **Mocking**: We will use the `github.com/stretchr/testify/mock` package to create mock implementations of the output ports (repositories) for testing the core services.
- **Mock Generation**: We will use `github.com/vektra/mockery/v2` to generate mock implementations of the repository interfaces defined in the core ports.
- **Test Runner**: We will use `go test` to run the tests, and we may integrate it with a CI/CD pipeline for automated testing.

## Test Folder Structure

For the current phase (`user_service`), tests should follow this structure:

```
internal/
	core/
		services/
			user_service.go
			user_service_test.go
```

Generated mocks should live in a dedicated mocks package (example):

```
internal/
	core/
		ports/
			mocks/
				user_repository.go
```

If the repository already uses a different generated-mocks path, keep using the existing project convention consistently.

## Naming Convention

1. **Test File Names**: `<source_file>_test.go` (example: `user_service_test.go`).
2. **Test Function Names**: Prefer one method-level test function using `Test<MethodName>` (example: `TestGetUserByID`) with table-driven sub-cases.
3. **Table-Driven Cases**: Use descriptive case names (example: `"repository returns not found error"`).
4. **Mock Types**: Use clear names tied to the interface (example: `MockUserRepository`).
5. **Assertions**: Use `assert` and `require` from `testify` for clear and consistent assertions in tests.
6. **Case Field Names Must Be Explicit**: Avoid generic fields like `name`, `input`, `result`, or `error` when testing service methods that involve one or more repositories.
7. **Repo + Method Prefixing**: For test-case fields related to dependency behavior, include both repository and method in the field name (example: `repoUserGetUserByIDResp`, `repoPlaylistCreatePlaylistErr`, `expectedShouldCallUserRepoGetUserByID`).
8. **Contract-Oriented Expected Fields**: Expected output fields should include the service method context (example: `expectedUserGetUserByID`, `expectedErrCreatePlaylist`) to prevent ambiguity across multi-repo flows.

## Test Case Pattern (Apply to All Services)

1. Use **one table-driven test per service method** that includes all core and edge scenarios.
2. Each test case should define **inputs** and expected **outputs/errors** as well as **mocked dependency behavior**.
3. Dependency mock behavior can be configured per case, but assertions should focus on method contract: returned value, returned error, and whether dependency calls happened when expected.
4. Include validation scenarios (e.g., empty/blank input) in the same table to keep behavior coverage centralized.
5. For service method inputs, use descriptive field names that include `input` + parameter name (example: `inputUserID` for parameter `userID`).
6. For expected service method outputs, use descriptive field names that include the `expected` + returned type (example: `expectedUser` for the returned `domain.User` or `domain.User` object).
7. For mocked dependency method responses, use descriptive field names that include the dependency name + method name + `Resp` (example: `userRepoGetUserByIDResp` for dependency method `userRepo.GetUserByID`).
8. For mocked dependency method errors, use descriptive field names that include the dependency name + method name + `Err` (example: `userRepoGetUserByIDErr` for the mocked error of `GetUserByID` in `UserRepository`).
9. For call flags, use descriptive field names that include `shouldCall` + the dependency name + method name (example: `shouldCallUserRepoGetUserByID` for a boolean flag indicating whether `GetUserByID` in `UserRepository` should be called).
10. For methods that call multiple dependencies or multiple methods on the same dependency, use explicit per-call expectation fields (inputs, mocked responses, call flags, call counts) named with dependency + method context.

## Initial Unit Test Cases for user_service

1. **Success Path**: `GetUserByID` returns a user when repository returns data.
2. **Repository Error Path**: `GetUserByID` returns the repository error unchanged.
3. **Edge Case - Empty ID**: Service returns a validation error for empty/blank IDs and does not call the repository.

## Out of Scope (Current Phase)

1. HTTP adapter tests (covered by integration tests later).
2. Defining a strict numeric coverage threshold.G

