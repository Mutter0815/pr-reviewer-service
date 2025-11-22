package service

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	// TODO: реализуем позже
	return nil, nil
}
