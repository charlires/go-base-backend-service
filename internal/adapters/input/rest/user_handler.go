package rest

import (
	"context"
	"log/slog"

	"github.com/charlires/go-base-backend-service/internal/adapters/input/rest/gen"
	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

// GetUserById handles GET /users/{userId}
func (h *Handlers) GetUserById(ctx context.Context, request gen.GetUserByIdRequestObject) (gen.GetUserByIdResponseObject, error) {
	ctx = logger.ContextWithLogger(ctx, slog.Default().With("user_id", request.UserId))
	logger.FromCtx(ctx).Debug("Handlers.GetUserById", "user_id", request.UserId)

	user, err := h.UserService.GetUserByID(ctx, request.UserId)
	if err != nil {
		logger.FromCtx(ctx).Error("Handlers.GetUserById: failed to get user", "error", err)
		return gen.GetUserById500JSONResponse{Message: "Failed to get user"}, nil
	}

	return gen.GetUserById200JSONResponse(domainUserToResponse(user)), nil
}

func domainUserToResponse(u *domain.User) gen.User {
	return gen.User{
		Id:       u.ID,
		Username: u.Username,
	}
}
