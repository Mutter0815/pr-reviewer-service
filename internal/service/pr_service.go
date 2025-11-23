package service

import (
	"context"
	"time"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
)

type PRService struct {
	prRepo   domain.PullRequestRepository
	userRepo domain.UserRepository
	teamRepo domain.TeamRepository
}

func NewPRService(
	prRepo domain.PullRequestRepository,
	userRepo domain.UserRepository,
	teamRepo domain.TeamRepository,
) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *PRService) CreatePR(ctx context.Context, pr *domain.PullRequest) error {
	pr.Status = domain.PullRequestStatusOpen

	if pr.CreatedAt.IsZero() {
		pr.CreatedAt = time.Now().UTC()
	}

	return s.prRepo.Create(ctx, pr)
}
