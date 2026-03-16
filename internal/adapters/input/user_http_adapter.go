package input

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/charlires/go-base-backend-service/internal/core/services"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

type UserHTTPAdapter struct {
	UserService services.UserService
}

func NewUserHTTPAdapter(userService services.UserService) *UserHTTPAdapter {
	return &UserHTTPAdapter{
		UserService: userService,
	}
}

func (a *UserHTTPAdapter) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userid")
	ctx := logger.ContextWithLogger(r.Context(), slog.Default().With(
		"method", r.Method,
		"path", r.URL.Path,
		"user_id", userID,
	))
	logger.FromCtx(ctx).Debug("UserHTTPAdapter.GetUserByID", "user_id", userID)

	user, err := a.UserService.GetUserByID(ctx, userID)
	if err != nil {
		logger.FromCtx(ctx).Error("UserHTTPAdapter.GetUserByID: failed to get user", "error", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User details: %+v", user)
}
