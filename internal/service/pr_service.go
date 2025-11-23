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

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return err

	}
	userRepo, ok1 := s.userRepo.(interface {
		GetByID(ctx context.Context, id string) (domain.User, error)
		ListActiveByTeam(ctx context.Context, teamName string) ([]domain.User, error)
	})
	prRepo, ok2 := s.prRepo.(interface {
		AssignReviewers(ctx context.Context, prID string, reviewerIDs []string) error
	})

	if !ok1 || !ok2 {
		return nil
	}

	author, err := userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil {
		return err
	}

	candidates, err := userRepo.ListActiveByTeam(ctx, author.TeamName)
	if err != nil {
		return err
	}

	var reviewers []string
	for _, u := range candidates {
		if u.ID == pr.AuthorID {
			continue
		}
		reviewers = append(reviewers, u.ID)
	}

	if len(reviewers) == 0 {
		return domain.ErrNoCandidate
	}

	if len(reviewers) > 2 {
		reviewers = reviewers[:2]
	}

	if err := prRepo.AssignReviewers(ctx, pr.ID, reviewers); err != nil {
		return err
	}

	return nil
}
