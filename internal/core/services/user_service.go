package services

import (
	"context"
	"errors"
	"strings"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/ports"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

var (
	ErrUserIDRequired = errors.New("user id is required")
)

// Input contract — consumed by input adapters
type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

// Implementation of the UserService interface
type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	logger.FromCtx(ctx).Debug("userService.GetUserByID", "user_id", id)
	if strings.TrimSpace(id) == "" {
		return nil, ErrUserIDRequired
	}

	return s.userRepo.GetUserByID(ctx, id)
}
