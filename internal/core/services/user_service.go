package services

import (
	"context"

	"github.com/charlires/go-base-backend-service/internal/core/domain"
	"github.com/charlires/go-base-backend-service/internal/core/ports"
)

// Input contract — consumed by input adapters
type UserService interface {
	GetUserByID(ctx context.Context, id string) (domain.User, error)
}

// Implementation of the UserService interface
type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}
