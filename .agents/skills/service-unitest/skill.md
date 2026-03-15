# Skill: Service Unit Test Writer

## Purpose

Write Go unit tests for a single service in `internal/core/services/`. Tests follow the project's table-driven pattern with mockery-generated mocks and testify assertions. After writing, verify by running `go test`.

---

## Inputs Required

Before starting, collect:

1. **Target service** — which service to test (e.g., `user_service`, `playlist_service`).
2. **Methods to test** — which methods on the service to cover. Default: all public methods.
3. **Test cases** — if the caller provides specific scenarios, include them. Otherwise derive them from the method signatures and error variables defined in the service file.

---

## Step-by-Step Execution

### Step 1 — Read the service file

Read `internal/core/services/<service_name>.go`.

Identify:
- The service interface and its method signatures.
- The concrete struct and its dependency fields injected via constructor.
- All exported error variables (e.g., `ErrUserIDRequired`).
- All validation branches that short-circuit before calling any dependency.

### Step 2 — Read the port interfaces

Read `internal/core/ports/repository.go`.

For each dependency used by the service, note:
- The interface name (e.g., `UserRepository`).
- The methods it exposes and their signatures.

### Step 3 — Check for existing mocks

Check `internal/core/ports/ports_mocks/` for a generated mock file matching each required dependency interface.

- If a mock file **exists**, proceed to Step 4.
- If a mock file **is missing**, run:

  ```bash
  go generate ./internal/core/ports/...
  ```

  Confirm the mock file was created before proceeding.

### Step 4 — Read the domain types

Read `internal/core/domain/` files for any domain structs returned or accepted by the service methods under test. Note field names and types.

### Step 5 — Determine test cases

For each method under test, define test cases covering:

1. **Happy path** — valid inputs, dependency returns expected data, service returns it unchanged.
2. **Dependency error path** — dependency returns an error, service propagates it.
3. **Validation / short-circuit paths** — invalid inputs (empty string, blank string, nil, zero value) that cause the service to return an error without calling the dependency.
4. **Multi-dependency paths** (if applicable) — when a method calls more than one dependency method, include cases where early calls succeed but later ones fail.

### Step 6 — Write the test file

Write `internal/core/services/<service_name>_test.go`.

Follow all conventions below.

### Step 7 — Run go test

Run:

```bash
go test ./internal/core/services/... -v -run Test
```

If the tests fail, diagnose and fix. Repeat until all tests pass.

---

## Naming Conventions

| Element | Convention | Example |
|---|---|---|
| Test file | `<source_file>_test.go` | `user_service_test.go` |
| Function per method | `Test<MethodName>` | `TestGetUserByID` |
| Table case field — input param | `input<ParamName>` | `inputUserID` |
| Table case field — dependency response | `<depName><MethodName>Resp` | `userRepoGetUserByIDResp` |
| Table case field — dependency error | `<depName><MethodName>Err` | `userRepoGetUserByIDErr` |
| Table case field — call flag | `shouldCall<DepName><MethodName>` | `shouldCallUserRepoGetUserByID` |
| Table case field — call count | `expected<DepName><MethodName>CallCount` | `expectedUserRepoGetUserByIDCallCount` |
| Table case field — expected return | `expected<ReturnType>` | `expectedUser` |
| Table case field — expected error | `expectedErr` | `expectedErr` |
| Table case name field | `testScenario` | `testScenario` |
| Mock type | `Mock<InterfaceName>` (generated) | `MockUserRepository` |

`<depName>` is the variable name of the dependency in the service struct (e.g., `userRepo`, `mailer`, `cache`). It is not always a repository — use the actual dependency name.

Do **not** use generic field names like `name`, `input`, `result`, or `error`.

---

## Test File Template

Use this template as the structural pattern. Replace all `<...>` placeholders with values specific to the service under test.

```go
package services

import (
    "context"
    "errors"
    "testing"

    "github.com/<org>/<repo>/internal/core/domain"
    "<module>/internal/core/ports/ports_mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func Test<MethodName>(t *testing.T) {
    t.Parallel()

    // Sentinel errors for test scenarios
    dependencyFailure := errors.New("dependency failure")
    _ = dependencyFailure // remove if unused

    testCases := []struct {
        testScenario string

        // --- inputs ---
        input<ParamName> <ParamType>

        // --- mocked dependency behavior ---
        // One block per dependency method called by the service method.
        <depName><DepMethodName>Resp  <ReturnType>
        <depName><DepMethodName>Err   error
        shouldCall<DepName><DepMethodName> bool
        expected<DepName><DepMethodName>CallCount int

        // --- expected outputs ---
        expected<ReturnType> <ReturnType>
        expectedErr          error
    }{
        {
            testScenario: "success",
            // populate fields
        },
        {
            testScenario: "dependency error",
            // populate fields
        },
        {
            testScenario: "empty <param>",
            // populate fields — shouldCall* false, expectedErr = Err<Validation>
        },
        {
            testScenario: "blank <param>",
            // populate fields — same as empty but whitespace input
        },
    }

    for _, tc := range testCases {
        t.Run(tc.testScenario, func(t *testing.T) {
            t.Parallel()

            ctx := context.Background()
            <depName> := ports_mocks.New<InterfaceName>(t)
            svc := New<ServiceName>(<depName>)

            if tc.shouldCall<DepName><DepMethodName> {
                <depName>.On("<DepMethodName>", ctx, tc.input<ParamName>).
                    Return(tc.<depName><DepMethodName>Resp, tc.<depName><DepMethodName>Err).
                    Once()
            }

            actual<ReturnType>, err := svc.<MethodName>(ctx, tc.input<ParamName>)

            if tc.expectedErr != nil {
                require.Error(t, err)
                assert.ErrorIs(t, err, tc.expectedErr)
            } else {
                require.NoError(t, err)
            }

            assert.Equal(t, tc.expected<ReturnType>, actual<ReturnType>)

            if tc.shouldCall<DepName><DepMethodName> {
                <depName>.AssertNumberOfCalls(t, "<DepMethodName>", tc.expected<DepName><DepMethodName>CallCount)
            } else {
                <depName>.AssertNotCalled(t, "<DepMethodName>", mock.Anything, mock.Anything)
            }
        })
    }
}
```

Key structural rules:
- `t.Parallel()` on both the top-level test and each sub-test.
- One `Test<MethodName>` function per service method; all cases live inside it.
- Mock setup (`On(...)`) is gated by the `shouldCall*` flag.
- Use `require.Error` / `require.NoError` before `assert.*` so the test stops early on unexpected nil/non-nil errors.
- Use `AssertNumberOfCalls` when the method should be called; use `AssertNotCalled` otherwise.
- For methods with multiple dependency calls, add a separate block of fields and setup/assertion per call.

---

## Multi-Dependency Method Pattern

When a service method calls more than one dependency, or calls the same dependency multiple times, add explicit fields for **each** call:

```go
// fields
firstDep<MethodName>Resp   <Type>
firstDep<MethodName>Err    error
shouldCallFirstDep<MethodName> bool
expectedFirstDep<MethodName>CallCount int

secondDep<MethodName>Resp   <Type>
secondDep<MethodName>Err    error
shouldCallSecondDep<MethodName> bool
expectedSecondDep<MethodName>CallCount int
```

Set up and assert each independently inside the test loop.

---

## Canonical Reference

The file `internal/core/services/user_service_test.go` is the canonical example of this pattern applied to a real service in this project. When in doubt about formatting or structure, read that file.

---

## Important Constraints

- Never write tests for adapters (HTTP, SQLite). Those are integration tests.
- Never import the adapter packages from a service test.
- Never define mocks manually; always use the generated mocks from `ports_mocks/`.
- Never set a fixed coverage percentage target; focus on critical paths and edge cases.
- Keep test logic inside the table loop — avoid large helper functions unless multiple test functions share non-trivial setup.
- Do not skip the `go test` verification step.
