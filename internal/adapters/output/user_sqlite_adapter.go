package output

import (
	"context"
	"database/sql"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/ports"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

var _ ports.UserRepository = (*UserSQLiteAdapter)(nil)

type UserSQLiteAdapter struct {
	db *sql.DB
}

func NewUserSQLiteAdapter(db *sql.DB) *UserSQLiteAdapter {
	return &UserSQLiteAdapter{db: db}
}

func (u *UserSQLiteAdapter) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	logger.FromCtx(ctx).Debug("UserSQLiteAdapter.GetUserByID", "user_id", id)
	var user domain.User
	err := u.db.QueryRowContext(ctx, "SELECT id, username FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
