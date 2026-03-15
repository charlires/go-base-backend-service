package ports

import (
	"context"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
)

// Output contract — implemented by output adapters

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (domain.User, error)
}

type PlaylistRepository interface {
	GetPlaylistByID(ctx context.Context, id string) (domain.Playlist, error)
	CreatePlaylist(ctx context.Context, playlist domain.Playlist) (string, error)
	AddTrackToPlaylist(ctx context.Context, playlistID string, trackID string) error
}

type TrackRepository interface {
	GetTrackByID(ctx context.Context, id string) (domain.Track, error)
}
