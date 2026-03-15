package ports

import (
	"context"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
)

// Output contract — implemented by output adapters

//go:generate go run github.com/vektra/mockery/v2@latest --all --output=./ports_mocks --outpkg=ports_mocks
type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

type PlaylistRepository interface {
	GetPlaylistByID(ctx context.Context, id string) (*domain.Playlist, error)
	CreatePlaylist(ctx context.Context, playlist *domain.Playlist) (string, error)
	AddTrackToPlaylist(ctx context.Context, playlistID string, trackID string) error
}

type TrackRepository interface {
	GetTrackByID(ctx context.Context, id string) (*domain.Track, error)
}
