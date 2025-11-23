package service

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepository
	prRepo   domain.PullRequestRepository
}

func NewUserService(userRepo domain.UserRepository, prRepo domain.PullRequestRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	return s.userRepo.SetIsActive(ctx, userID, isActive)
}

func (s *UserService) ListReviewerPRs(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	if _, err := s.userRepo.GetByID(ctx, reviewerID); err != nil {
		return nil, err
	}

	return s.prRepo.ListByReviewer(ctx, reviewerID)
}
