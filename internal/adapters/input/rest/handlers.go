package rest

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../../oapi-codegen.yaml ../../../../openapi.yaml

import (
	"context"
	"log/slog"

	"github.com/charlires/go-base-backend-service/internal/adapters/input/rest/gen"
	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/services"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

// Handlers implements gen.StrictServerInterface and wires all service dependencies.
type Handlers struct {
	UserService     services.UserService
	TrackService    services.TrackService
	PlaylistService services.PlaylistService
}

// Compile-time check that Handlers satisfies the StrictServerInterface.
var _ gen.StrictServerInterface = (*Handlers)(nil)

// NewHandlers creates a new Handlers instance.
func NewHandlers(
	userService services.UserService,
	trackService services.TrackService,
	playlistService services.PlaylistService,
) *Handlers {
	return &Handlers{
		UserService:     userService,
		TrackService:    trackService,
		PlaylistService: playlistService,
	}
}

// GetUserById handles GET /users/{userId}
func (h *Handlers) GetUserById(ctx context.Context, request gen.GetUserByIdRequestObject) (gen.GetUserByIdResponseObject, error) {
	ctx = logger.ContextWithLogger(ctx, slog.Default().With("user_id", request.UserId))
	logger.FromCtx(ctx).Debug("Handlers.GetUserById", "user_id", request.UserId)

	user, err := h.UserService.GetUserByID(ctx, request.UserId)
	if err != nil {
		logger.FromCtx(ctx).Error("Handlers.GetUserById: failed to get user", "error", err)
		return gen.GetUserById500JSONResponse{Message: "Failed to get user"}, nil
	}

	return gen.GetUserById200JSONResponse(domainUserToResponse(user)), nil
}

// GetTrackById handles GET /tracks/{trackId}
func (h *Handlers) GetTrackById(ctx context.Context, request gen.GetTrackByIdRequestObject) (gen.GetTrackByIdResponseObject, error) {
	ctx = logger.ContextWithLogger(ctx, slog.Default().With("track_id", request.TrackId))
	logger.FromCtx(ctx).Debug("Handlers.GetTrackById", "track_id", request.TrackId)

	track, err := h.TrackService.GetTrackByID(ctx, request.TrackId)
	if err != nil {
		logger.FromCtx(ctx).Error("Handlers.GetTrackById: failed to get track", "error", err)
		return gen.GetTrackById500JSONResponse{Message: "Failed to get track"}, nil
	}

	return gen.GetTrackById200JSONResponse(domainTrackToResponse(&track)), nil
}

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

// --- Domain mapping helpers ---

func domainUserToResponse(u *domain.User) gen.User {
	return gen.User{
		Id:       u.ID,
		Username: u.Username,
	}
}

func domainTrackToResponse(t *domain.Track) gen.Track {
	return gen.Track{
		Id:     t.ID,
		Title:  t.Title,
		Artist: t.Artist,
		Album:  t.Album,
	}
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
