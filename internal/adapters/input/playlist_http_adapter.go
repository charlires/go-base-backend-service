package input

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/services"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

type PlaylistHTTPAdapter struct {
	PlaylistService services.PlaylistService
}

func NewPlaylistHTTPAdapter(playlistService services.PlaylistService) *PlaylistHTTPAdapter {
	return &PlaylistHTTPAdapter{
		PlaylistService: playlistService,
	}
}

// Implement HTTP handlers for playlist-related operations here
func (a *PlaylistHTTPAdapter) GetPlaylistByID(w http.ResponseWriter, r *http.Request) {
	playlistID := r.PathValue("playlistid")
	ctx := logger.ContextWithLogger(r.Context(), slog.Default().With(
		"method", r.Method,
		"path", r.URL.Path,
		"playlist_id", playlistID,
	))
	logger.FromCtx(ctx).Debug("PlaylistHTTPAdapter.GetPlaylistByID", "playlist_id", playlistID)

	playlist, err := a.PlaylistService.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		logger.FromCtx(ctx).Error("PlaylistHTTPAdapter.GetPlaylistByID: failed to get playlist", "error", err)
		http.Error(w, "Failed to get playlist", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(playlist); err != nil {
		logger.FromCtx(ctx).Error("PlaylistHTTPAdapter.GetPlaylistByID: failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (a *PlaylistHTTPAdapter) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	ctx := logger.ContextWithLogger(r.Context(), slog.Default().With(
		"method", r.Method,
		"path", r.URL.Path,
	))
	logger.FromCtx(ctx).Debug("PlaylistHTTPAdapter.CreatePlaylist")

	var playlist domain.Playlist
	if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
		logger.FromCtx(ctx).Error("PlaylistHTTPAdapter.CreatePlaylist: invalid request payload", "error", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdPlaylist, err := a.PlaylistService.CreatePlaylist(ctx, &playlist)
	if err != nil {
		logger.FromCtx(ctx).Error("PlaylistHTTPAdapter.CreatePlaylist: failed to create playlist", "error", err)
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdPlaylist); err != nil {
		logger.FromCtx(ctx).Error("PlaylistHTTPAdapter.CreatePlaylist: failed to encode response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (a *PlaylistHTTPAdapter) AddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := r.PathValue("playlistid")
	trackID := r.PathValue("trackid")
	ctx := logger.ContextWithLogger(r.Context(), slog.Default().With(
		"method", r.Method,
		"path", r.URL.Path,
		"playlist_id", playlistID,
		"track_id", trackID,
	))
	logger.FromCtx(ctx).Debug("PlaylistHTTPAdapter.AddTrackToPlaylist", "playlist_id", playlistID, "track_id", trackID)

	if err := a.PlaylistService.AddTrackToPlaylist(ctx, playlistID, trackID); err != nil {
		logger.FromCtx(ctx).Error("PlaylistHTTPAdapter.AddTrackToPlaylist: failed to add track to playlist", "error", err)
		http.Error(w, "Failed to add track to playlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
