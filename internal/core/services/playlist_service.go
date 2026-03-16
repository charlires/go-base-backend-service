package services

import (
	"context"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/ports"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

type PlaylistService interface {
	GetPlaylistByID(ctx context.Context, id string) (domain.Playlist, error)
	CreatePlaylist(ctx context.Context, playlist domain.Playlist) (domain.Playlist, error)
	AddTrackToPlaylist(ctx context.Context, playlistID string, trackID string) error
}

type playlistService struct {
	playlistRepo ports.PlaylistRepository
	userRepo     ports.UserRepository
}

func NewPlaylistService(playlistRepo ports.PlaylistRepository, userRepo ports.UserRepository) PlaylistService {
	return &playlistService{
		playlistRepo: playlistRepo,
		userRepo:     userRepo,
	}
}

func (p *playlistService) AddTrackToPlaylist(ctx context.Context, playlistID string, trackID string) error {
	logger.FromCtx(ctx).Debug("playlistService.AddTrackToPlaylist", "playlist_id", playlistID, "track_id", trackID)
	panic("unimplemented")
}

func (p *playlistService) CreatePlaylist(ctx context.Context, playlist domain.Playlist) (domain.Playlist, error) {
	logger.FromCtx(ctx).Debug("playlistService.CreatePlaylist", "playlist", &playlist)
	panic("unimplemented")
}

func (p *playlistService) GetPlaylistByID(ctx context.Context, id string) (domain.Playlist, error) {
	logger.FromCtx(ctx).Debug("playlistService.GetPlaylistByID", "playlist_id", id)
	playlist, err := p.playlistRepo.GetPlaylistByID(ctx, id)
	if err != nil {
		return domain.Playlist{}, err
	}
	return playlist, nil
}
