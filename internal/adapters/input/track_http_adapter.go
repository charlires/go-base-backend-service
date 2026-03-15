package input

import (
	"fmt"
	"net/http"

	"github.com/charlires/go-base-backend-service/internal/core/services"
)

type TrackHTTPAdapter struct {
	TrackService services.TrackService
	// Add any dependencies (e.g., services) here
}

func NewTrackHTTPAdapter(trackService services.TrackService) *TrackHTTPAdapter {
	return &TrackHTTPAdapter{
		TrackService: trackService,
	}
}

// Implement HTTP handlers for track-related operations here
func (a *TrackHTTPAdapter) GetTrackByID(w http.ResponseWriter, r *http.Request) {
	trackID := r.PathValue("trackid")
	track, err := a.TrackService.GetTrackByID(r.Context(), trackID) // Replace with actual ID extraction
	if err != nil {
		http.Error(w, "Failed to get track", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Track details: %+v", track)
}
