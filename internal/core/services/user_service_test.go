package services

import (
	"context"
	"errors"
	"testing"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/ports/ports_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	repositoryFailure := errors.New("repository failure")

	testCases := []struct {
		testScenario                         string
		inputUserID                          string
		userRepoGetUserByIDResp              *domain.User
		userRepoGetUserByIDErr               error
		expectedUser                         *domain.User
		expectedErr                          error
		expectedUserRepoGetUserByIDCallCount int
		shouldCallUserRepoGetUserByID        bool
	}{
		{
			testScenario:                         "success",
			inputUserID:                          "user-001",
			userRepoGetUserByIDResp:              &domain.User{ID: "user-001", Username: "carlos"},
			expectedUser:                         &domain.User{ID: "user-001", Username: "carlos"},
			expectedUserRepoGetUserByIDCallCount: 1,
			shouldCallUserRepoGetUserByID:        true,
		},
		{
			testScenario:                         "repository error",
			inputUserID:                          "user-999",
			userRepoGetUserByIDErr:               repositoryFailure,
			expectedUser:                         nil,
			expectedErr:                          repositoryFailure,
			expectedUserRepoGetUserByIDCallCount: 1,
			shouldCallUserRepoGetUserByID:        true,
		},
		{
			testScenario:                         "empty id",
			inputUserID:                          "",
			expectedUser:                         nil,
			expectedErr:                          ErrUserIDRequired,
			expectedUserRepoGetUserByIDCallCount: 0,
			shouldCallUserRepoGetUserByID:        false,
		},
		{
			testScenario:                         "blank id",
			inputUserID:                          "   ",
			expectedUser:                         nil,
			expectedErr:                          ErrUserIDRequired,
			expectedUserRepoGetUserByIDCallCount: 0,
			shouldCallUserRepoGetUserByID:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testScenario, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			userRepo := ports_mocks.NewUserRepository(t)
			svc := NewUserService(userRepo)

			if tc.shouldCallUserRepoGetUserByID {
				userRepo.On("GetUserByID", ctx, tc.inputUserID).Return(tc.userRepoGetUserByIDResp, tc.userRepoGetUserByIDErr).Once()
			}

			actualUser, err := svc.GetUserByID(ctx, tc.inputUserID)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expectedUser, actualUser)
			if tc.shouldCallUserRepoGetUserByID {
				userRepo.AssertNumberOfCalls(t, "GetUserByID", tc.expectedUserRepoGetUserByIDCallCount)
			} else {
				userRepo.AssertNotCalled(t, "GetUserByID", mock.Anything, mock.Anything)
			}
		})
	}
}
