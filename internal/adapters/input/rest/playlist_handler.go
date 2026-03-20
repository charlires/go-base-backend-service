package rest

import (
	"context"
	"log/slog"

	"github.com/charlires/go-base-backend-service/internal/adapters/input/rest/gen"
	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

// GetPlaylistById handles GET /playlists/{playlistId}
func (h *Handlers) GetPlaylistById(ctx context.Context, request gen.GetPlaylistByIdRequestObject) (gen.GetPlaylistByIdResponseObject, error) {
	ctx = logger.ContextWithLogger(ctx, slog.Default().With("playlist_id", request.PlaylistId))
	logger.FromCtx(ctx).Debug("Handlers.GetPlaylistById", "playlist_id", request.PlaylistId)

	playlist, err := h.PlaylistService.GetPlaylistByID(ctx, request.PlaylistId)
	if err != nil {
		logger.FromCtx(ctx).Error("Handlers.GetPlaylistById: failed to get playlist", "error", err)
		return gen.GetPlaylistById500JSONResponse{Message: "Failed to get playlist"}, nil
	}

	return gen.GetPlaylistById200JSONResponse(domainPlaylistToResponse(playlist)), nil
}

// CreatePlaylist handles POST /playlists
func (h *Handlers) CreatePlaylist(ctx context.Context, request gen.CreatePlaylistRequestObject) (gen.CreatePlaylistResponseObject, error) {
	logger.FromCtx(ctx).Debug("Handlers.CreatePlaylist")

	if request.Body == nil {
		return gen.CreatePlaylist400JSONResponse{Message: "Invalid request payload"}, nil
	}

	playlist := &domain.Playlist{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		OwnerID:     request.Body.OwnerId,
	}

	created, err := h.PlaylistService.CreatePlaylist(ctx, playlist)
	if err != nil {
		logger.FromCtx(ctx).Error("Handlers.CreatePlaylist: failed to create playlist", "error", err)
		return gen.CreatePlaylist500JSONResponse{Message: "Failed to create playlist"}, nil
	}

	return gen.CreatePlaylist201JSONResponse(domainPlaylistToResponse(created)), nil
}

// AddTrackToPlaylist handles POST /playlists/{playlistId}/tracks/{trackId}
func (h *Handlers) AddTrackToPlaylist(ctx context.Context, request gen.AddTrackToPlaylistRequestObject) (gen.AddTrackToPlaylistResponseObject, error) {
	ctx = logger.ContextWithLogger(ctx, slog.Default().With("playlist_id", request.PlaylistId, "track_id", request.TrackId))
	logger.FromCtx(ctx).Debug("Handlers.AddTrackToPlaylist", "playlist_id", request.PlaylistId, "track_id", request.TrackId)

	if err := h.PlaylistService.AddTrackToPlaylist(ctx, request.PlaylistId, request.TrackId); err != nil {
		logger.FromCtx(ctx).Error("Handlers.AddTrackToPlaylist: failed to add track to playlist", "error", err)
		return gen.AddTrackToPlaylist500JSONResponse{Message: "Failed to add track to playlist"}, nil
	}

	return gen.AddTrackToPlaylist204Response{}, nil
}

func domainPlaylistToResponse(p *domain.Playlist) gen.Playlist {
	tracks := make([]gen.Track, len(p.Tracks))
	for i, t := range p.Tracks {
		tracks[i] = domainTrackToResponse(&t)
	}
	return gen.Playlist{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerId:     p.OwnerID,
		Tracks:      &tracks,
	}
}
