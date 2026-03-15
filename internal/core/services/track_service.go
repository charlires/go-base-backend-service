package services

import (
	"context"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/ports"
)

// Ensure implements TrackService interface
var _ TrackService = (*trackService)(nil)

type TrackService interface {
	GetTrackByID(ctx context.Context, id string) (domain.Track, error)
	// Define methods related to track operations here
}

type trackService struct {
	trackRepo ports.TrackRepository
}

func NewTrackService(trackRepo ports.TrackRepository) TrackService {
	return &trackService{
		trackRepo: trackRepo,
	}
}

func (s *trackService) GetTrackByID(ctx context.Context, id string) (domain.Track, error) {
	panic("unimplemented")
}
