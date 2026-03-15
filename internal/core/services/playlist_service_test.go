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

func TestGetPlaylistByID(t *testing.T) {
	t.Parallel()

	repositoryFailure := errors.New("repository failure")

	testCases := []struct {
		testScenario                                 string
		inputPlaylistID                              string
		playlistRepoGetPlaylistByIDResp              *domain.Playlist
		playlistRepoGetPlaylistByIDErr               error
		expectedPlaylist                             *domain.Playlist
		expectedErr                                  error
		expectedPlaylistRepoGetPlaylistByIDCallCount int
		shouldCallPlaylistRepoGetPlaylistByID        bool
	}{
		{
			testScenario:    "success",
			inputPlaylistID: "playlist-001",
			playlistRepoGetPlaylistByIDResp: &domain.Playlist{
				ID:          "playlist-001",
				Name:        "Top Hits",
				Description: "Most played tracks",
				OwnerID:     "user-001",
			},
			expectedPlaylist: &domain.Playlist{
				ID:          "playlist-001",
				Name:        "Top Hits",
				Description: "Most played tracks",
				OwnerID:     "user-001",
			},
			expectedPlaylistRepoGetPlaylistByIDCallCount: 1,
			shouldCallPlaylistRepoGetPlaylistByID:        true,
		},
		{
			testScenario:                   "repository error",
			inputPlaylistID:                "playlist-999",
			playlistRepoGetPlaylistByIDErr: repositoryFailure,
			expectedPlaylist:               nil,
			expectedErr:                    repositoryFailure,
			expectedPlaylistRepoGetPlaylistByIDCallCount: 1,
			shouldCallPlaylistRepoGetPlaylistByID:        true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testScenario, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			playlistRepo := ports_mocks.NewPlaylistRepository(t)
			userRepo := ports_mocks.NewUserRepository(t)
			svc := NewPlaylistService(playlistRepo, userRepo)

			if tc.shouldCallPlaylistRepoGetPlaylistByID {
				playlistRepo.On("GetPlaylistByID", ctx, tc.inputPlaylistID).
					Return(tc.playlistRepoGetPlaylistByIDResp, tc.playlistRepoGetPlaylistByIDErr).
					Once()
			}

			actualPlaylist, err := svc.GetPlaylistByID(ctx, tc.inputPlaylistID)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expectedPlaylist, actualPlaylist)
			if tc.shouldCallPlaylistRepoGetPlaylistByID {
				playlistRepo.AssertNumberOfCalls(t, "GetPlaylistByID", tc.expectedPlaylistRepoGetPlaylistByIDCallCount)
			} else {
				playlistRepo.AssertNotCalled(t, "GetPlaylistByID", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestCreatePlaylist(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		testScenario  string
		inputPlaylist *domain.Playlist
	}{
		{
			testScenario: "panics as unimplemented",
			inputPlaylist: &domain.Playlist{
				Name:        "Fresh Finds",
				Description: "New discoveries",
				OwnerID:     "user-100",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testScenario, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			playlistRepo := ports_mocks.NewPlaylistRepository(t)
			userRepo := ports_mocks.NewUserRepository(t)
			svc := NewPlaylistService(playlistRepo, userRepo)

			assert.PanicsWithValue(t, "unimplemented", func() {
				_, _ = svc.CreatePlaylist(ctx, tc.inputPlaylist)
			})
		})
	}
}

func TestAddTrackToPlaylist(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		testScenario    string
		inputPlaylistID string
		inputTrackID    string
	}{
		{
			testScenario:    "panics as unimplemented",
			inputPlaylistID: "playlist-001",
			inputTrackID:    "track-001",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testScenario, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			playlistRepo := ports_mocks.NewPlaylistRepository(t)
			userRepo := ports_mocks.NewUserRepository(t)
			svc := NewPlaylistService(playlistRepo, userRepo)

			assert.PanicsWithValue(t, "unimplemented", func() {
				_ = svc.AddTrackToPlaylist(ctx, tc.inputPlaylistID, tc.inputTrackID)
			})
		})
	}
}
