package rest

import (
	"context"
	"log/slog"

	"github.com/charlires/go-base-backend-service/internal/adapters/input/rest/gen"
	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

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

func domainTrackToResponse(t *domain.Track) gen.Track {
	return gen.Track{
		Id:     t.ID,
		Title:  t.Title,
		Artist: t.Artist,
		Album:  t.Album,
	}
}
