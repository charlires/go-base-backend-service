package input

import (
	"encoding/json"
	"net/http"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/services"
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
	playlist, err := a.PlaylistService.GetPlaylistByID(r.Context(), playlistID) // Replace with actual ID extraction
	if err != nil {
		http.Error(w, "Failed to get playlist", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(playlist); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (a *PlaylistHTTPAdapter) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var playlist domain.Playlist
	if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdPlaylist, err := a.PlaylistService.CreatePlaylist(r.Context(), playlist)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdPlaylist); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (a *PlaylistHTTPAdapter) AddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID := r.PathValue("playlistid")
	trackID := r.PathValue("trackid")

	if err := a.PlaylistService.AddTrackToPlaylist(r.Context(), playlistID, trackID); err != nil {
		http.Error(w, "Failed to add track to playlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
