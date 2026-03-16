package output

import (
	"context"
	"database/sql"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/ports"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

// Ensure implements TrackRepository interface
var _ ports.TrackRepository = (*TrackSQLiteAdapter)(nil)

type TrackSQLiteAdapter struct {
	db *sql.DB
	// Add any dependencies (e.g., database connection) here
}

func NewTrackSQLiteAdapter(db *sql.DB) *TrackSQLiteAdapter {
	return &TrackSQLiteAdapter{
		db: db,
		// Initialize dependencies here
	}
}

func (u *TrackSQLiteAdapter) GetTrackByID(ctx context.Context, id string) (domain.Track, error) {
	logger.FromCtx(ctx).Debug("TrackSQLiteAdapter.GetTrackByID", "track_id", id)
	var track domain.Track
	err := u.db.QueryRowContext(ctx, "SELECT id, title, artist FROM tracks WHERE id = ?", id).Scan(&track.ID, &track.Title, &track.Artist)
	if err != nil {
		return domain.Track{}, err
	}
	return track, nil
}
