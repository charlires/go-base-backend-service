package domain

import "log/slog"

type Track struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

func (t *Track) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", t.ID),
		slog.String("title", t.Title),
		slog.String("artist", t.Artist),
		slog.String("album", t.Album),
	)
}
