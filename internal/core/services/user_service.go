package services

import "context"

// Input contract — consumed by input adapters
type UserService interface {
	GetUser(ctx context.Context, query GetUserQuery) (User, error)
}

// Concrete implementation — injected by main.go
type userService struct {
	repo ports.UserRepository
}
