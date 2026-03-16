package output

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/ports"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

// Ensure implements PlaylistRepository interface
var _ ports.PlaylistRepository = (*PlaylistSQLiteAdapter)(nil)

type PlaylistSQLiteAdapter struct {
	db *sql.DB
}

func NewPlaylistSQLiteAdapter(db *sql.DB) *PlaylistSQLiteAdapter {
	return &PlaylistSQLiteAdapter{
		db: db,
	}
}

func (p *PlaylistSQLiteAdapter) AddTrackToPlaylist(ctx context.Context, playlistID string, trackID string) error {
	logger.FromCtx(ctx).Debug("PlaylistSQLiteAdapter.AddTrackToPlaylist", "playlist_id", playlistID, "track_id", trackID)
	err := p.db.QueryRowContext(ctx, "INSERT INTO playlist_tracks (playlist_id, track_id) VALUES (?, ?)", playlistID, trackID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (p *PlaylistSQLiteAdapter) CreatePlaylist(ctx context.Context, playlist domain.Playlist) (string, error) {
	logger.FromCtx(ctx).Debug("PlaylistSQLiteAdapter.CreatePlaylist", "playlist", &playlist)
	result, err := p.db.ExecContext(ctx, "INSERT INTO playlists (id, name, owner_id) VALUES (?, ?, ?)", playlist.ID, playlist.Name, playlist.OwnerID)
	if err != nil {
		return "", err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return fmt.Sprint(id), nil
}

func (p *PlaylistSQLiteAdapter) GetPlaylistByID(ctx context.Context, id string) (domain.Playlist, error) {
	logger.FromCtx(ctx).Debug("PlaylistSQLiteAdapter.GetPlaylistByID", "playlist_id", id)
	var playlist domain.Playlist
	err := p.db.QueryRowContext(ctx, "SELECT id, name, owner_id FROM playlists WHERE id = ?", id).Scan(&playlist.ID, &playlist.Name, &playlist.OwnerID)
	if err != nil {
		return domain.Playlist{}, err
	}
	// get tracks for the playlist
	rows, err := p.db.QueryContext(ctx, "SELECT t.id, t.title, t.artist FROM tracks t JOIN playlist_tracks pt ON t.id = pt.track_id WHERE pt.playlist_id = ?", id)
	if err != nil {
		return domain.Playlist{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var track domain.Track
		if err := rows.Scan(&track.ID, &track.Title, &track.Artist); err != nil {
			return domain.Playlist{}, err
		}
		playlist.Tracks = append(playlist.Tracks, track)
	}

	if err := rows.Err(); err != nil {
		return domain.Playlist{}, err
	}
	return playlist, nil
}
