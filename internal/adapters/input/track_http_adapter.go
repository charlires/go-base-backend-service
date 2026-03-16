package input

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/charlires/go-base-backend-service/internal/core/services"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
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
	ctx := logger.ContextWithLogger(r.Context(), slog.Default().With(
		"method", r.Method,
		"path", r.URL.Path,
		"track_id", trackID,
	))
	logger.FromCtx(ctx).Debug("TrackHTTPAdapter.GetTrackByID", "track_id", trackID)

	track, err := a.TrackService.GetTrackByID(ctx, trackID)
	if err != nil {
		logger.FromCtx(ctx).Error("TrackHTTPAdapter.GetTrackByID: failed to get track", "error", err)
		http.Error(w, "Failed to get track", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Track details: %+v", track)
}
