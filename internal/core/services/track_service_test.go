package services

import (
	"context"
	"testing"

	"github.com/charlires/go-base-backend-service/internal/core/ports/ports_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTrackByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		testScenario                        string
		inputTrackID                        string
		expectedTrackRepoGetTrackByIDCalled bool
	}{
		{
			testScenario:                        "panics as unimplemented",
			inputTrackID:                        "track-001",
			expectedTrackRepoGetTrackByIDCalled: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testScenario, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			trackRepo := ports_mocks.NewTrackRepository(t)
			svc := NewTrackService(trackRepo)

			assert.PanicsWithValue(t, "unimplemented", func() {
				_, _ = svc.GetTrackByID(ctx, tc.inputTrackID)
			})

			if tc.expectedTrackRepoGetTrackByIDCalled {
				trackRepo.AssertCalled(t, "GetTrackByID", ctx, tc.inputTrackID)
			} else {
				trackRepo.AssertNotCalled(t, "GetTrackByID", mock.Anything, mock.Anything)
			}
		})
	}
}
