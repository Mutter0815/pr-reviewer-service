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

func (s *PRService) CreatePR(ctx context.Context, pr *domain.PullRequest) (domain.PullRequest, error) {
	pr.Status = domain.PullRequestStatusOpen

	if pr.CreatedAt.IsZero() {
		pr.CreatedAt = time.Now().UTC()
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return domain.PullRequest{}, err
	}

	author, err := s.userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	candidates, err := s.userRepo.ListActiveByTeam(ctx, author.TeamName)
	if err != nil {
		return domain.PullRequest{}, err
	}

	var reviewers []string
	for _, u := range candidates {
		if u.ID == pr.AuthorID {
			continue
		}
		if !u.IsActive {
			continue
		}
		reviewers = append(reviewers, u.ID)
	}

	if len(reviewers) > 2 {
		reviewers = reviewers[:2]
	}

	if len(reviewers) > 0 {
		if err := s.prRepo.AssignReviewers(ctx, pr.ID, reviewers); err != nil {
			return domain.PullRequest{}, err
		}
	}

	pr.AssignedReviewers = append([]string(nil), reviewers...)

	created, err := s.prRepo.GetByID(ctx, pr.ID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	if len(created.AssignedReviewers) == 0 {
		created.AssignedReviewers = append([]string(nil), reviewers...)
	}

	return created, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (domain.PullRequest, string, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	if pr.Status == domain.PullRequestStatusMerged {
		return domain.PullRequest{}, "", domain.ErrPRMerged
	}

	reviewers, err := s.prRepo.ListReviewers(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	found := false
	for _, rID := range reviewers {
		if rID == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return domain.PullRequest{}, "", domain.ErrNotAssigned
	}

	oldReviewer, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	active, err := s.userRepo.ListActiveByTeam(ctx, oldReviewer.TeamName)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	assigned := make(map[string]struct{}, len(reviewers))
	for _, rID := range reviewers {
		assigned[rID] = struct{}{}
	}

	var newID string
	for _, u := range active {
		if u.ID == pr.AuthorID {
			continue
		}
		if u.ID == oldReviewerID {
			continue
		}
		if _, used := assigned[u.ID]; used {
			continue
		}
		newID = u.ID
		break
	}

	reuseAssigned := false
	if newID == "" {
		for _, u := range active {
			if u.ID == pr.AuthorID {
				continue
			}
			if u.ID == oldReviewerID {
				continue
			}
			if _, used := assigned[u.ID]; !used {
				continue
			}
			newID = u.ID
			reuseAssigned = true
			break
		}
	}

	if newID == "" {
		return domain.PullRequest{}, "", domain.ErrNoCandidate
	}

	if reuseAssigned {
		if err := s.prRepo.RemoveReviewer(ctx, prID, newID); err != nil {
			return domain.PullRequest{}, "", err
		}
	}

	if err := s.prRepo.ReassignReviewer(ctx, prID, oldReviewerID, newID); err != nil {
		return domain.PullRequest{}, "", err
	}

	updated, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	return updated, newID, nil
}

func (s *PRService) MergePR(ctx context.Context, prID string) (domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	if pr.Status == domain.PullRequestStatusMerged {
		return pr, nil
	}

	if err := s.prRepo.Merge(ctx, prID); err != nil {
		return domain.PullRequest{}, err
	}

	updated, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return updated, nil
}
