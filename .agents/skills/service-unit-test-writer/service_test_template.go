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