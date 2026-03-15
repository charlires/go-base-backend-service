package domain

type Playlist struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Tracks      []Track `json:"tracks"` // List of track IDs
	OwnerID     string  `json:"owner_id"`
}
