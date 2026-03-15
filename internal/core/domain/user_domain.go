package domain

import "log/slog"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (u *User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", u.ID),
		slog.String("username", u.Username),
	)
}
