package domain

import "log/slog"

type Playlist struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Tracks      []Track `json:"tracks"` // List of track IDs
	OwnerID     string  `json:"owner_id"`
}

func (p *Playlist) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", p.ID),
		slog.String("name", p.Name),
		slog.String("owner_id", p.OwnerID),
	)
}
