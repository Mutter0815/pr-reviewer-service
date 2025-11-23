package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
)

type prRepoFake struct {
	prs       map[string]domain.PullRequest
	reviewers map[string][]string
	createErr error
}

func (r *prRepoFake) Create(ctx context.Context, pr *domain.PullRequest) error {
	if r.createErr != nil {
		return r.createErr
	}
	if r.prs == nil {
		r.prs = make(map[string]domain.PullRequest)
	}
	copy := *pr
	r.prs[pr.ID] = copy
	return nil
}

func (r *prRepoFake) AssignReviewers(ctx context.Context, prID string, reviewerIDs []string) error {
	if r.reviewers == nil {
		r.reviewers = make(map[string][]string)
	}
	r.reviewers[prID] = append([]string(nil), reviewerIDs...)

	pr := r.prs[prID]
	pr.AssignedReviewers = append([]string(nil), reviewerIDs...)
	r.prs[prID] = pr
	return nil
}

func (r *prRepoFake) GetByID(ctx context.Context, id string) (domain.PullRequest, error) {
	pr, ok := r.prs[id]
	if !ok {
		return domain.PullRequest{}, domain.ErrNotFound
	}
	return pr, nil
}

func (r *prRepoFake) ListReviewers(ctx context.Context, prID string) ([]string, error) {
	return append([]string(nil), r.reviewers[prID]...), nil
}

func (r *prRepoFake) ReassignReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error {
	list := r.reviewers[prID]
	for i, id := range list {
		if id == oldReviewerID {
			list[i] = newReviewerID
			break
		}
	}
	r.reviewers[prID] = list

	pr := r.prs[prID]
	pr.AssignedReviewers = append([]string(nil), list...)
	r.prs[prID] = pr
	return nil
}

func (r *prRepoFake) Merge(ctx context.Context, prID string) error {
	pr := r.prs[prID]
	pr.Status = domain.PullRequestStatusMerged
	now := time.Now().UTC()
	pr.MergedAt = &now
	r.prs[prID] = pr
	return nil
}

func (r *prRepoFake) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	return nil, nil
}

func (r *prRepoFake) RemoveReviewer(ctx context.Context, prID, reviewerID string) error {
	list := r.reviewers[prID]
	filtered := make([]string, 0, len(list))
	for _, id := range list {
		if id != reviewerID {
			filtered = append(filtered, id)
		}
	}
	r.reviewers[prID] = filtered

	if pr, ok := r.prs[prID]; ok {
		out := make([]string, 0, len(pr.AssignedReviewers))
		for _, id := range pr.AssignedReviewers {
			if id != reviewerID {
				out = append(out, id)
			}
		}
		pr.AssignedReviewers = out
		r.prs[prID] = pr
	}

	return nil
}

type userRepoFake struct {
	usersByID    map[string]domain.User
	activeByTeam map[string][]domain.User
}

func (r *userRepoFake) Upsert(ctx context.Context, u domain.User) error {
	if r.usersByID == nil {
		r.usersByID = make(map[string]domain.User)
	}
	r.usersByID[u.ID] = u
	return nil
}

func (r *userRepoFake) GetByID(ctx context.Context, id string) (domain.User, error) {
	u, ok := r.usersByID[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}

func (r *userRepoFake) ListActiveByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	return r.activeByTeam[teamName], nil
}

func (r *userRepoFake) SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	u, ok := r.usersByID[userID]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	u.IsActive = isActive
	r.usersByID[userID] = u
	return u, nil
}

func TestPRService_CreatePR_AssignReviewers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		users         []domain.User
		authorID      string
		wantCount     int
		wantReviewer1 string
	}{
		{
			name: "two reviewers",
			users: []domain.User{
				{ID: "u1", TeamName: "team", IsActive: true},
				{ID: "u2", TeamName: "team", IsActive: true},
				{ID: "u3", TeamName: "team", IsActive: true},
			},
			authorID:      "u1",
			wantCount:     2,
			wantReviewer1: "u2",
		},
		{
			name: "one reviewer",
			users: []domain.User{
				{ID: "a1", TeamName: "small", IsActive: true},
				{ID: "a2", TeamName: "small", IsActive: true},
			},
			authorID:      "a1",
			wantCount:     1,
			wantReviewer1: "a2",
		},
		{
			name: "no candidates",
			users: []domain.User{
				{ID: "s1", TeamName: "solo", IsActive: true},
			},
			authorID:  "s1",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &userRepoFake{
				usersByID:    make(map[string]domain.User),
				activeByTeam: make(map[string][]domain.User),
			}

			for _, u := range tt.users {
				userRepo.usersByID[u.ID] = u
				userRepo.activeByTeam[u.TeamName] = append(userRepo.activeByTeam[u.TeamName], u)
			}

			prRepo := &prRepoFake{}
			svc := NewPRService(prRepo, userRepo, nil)

			pr := &domain.PullRequest{
				ID:       "pr-" + tt.name,
				Name:     "Test PR",
				AuthorID: tt.authorID,
			}

			created, err := svc.CreatePR(ctx, pr)
			if err != nil {
				t.Fatalf("CreatePR error: %v", err)
			}

			got := prRepo.reviewers[pr.ID]
			if len(got) != tt.wantCount {
				t.Fatalf("expected %d reviewers, got %d (%v)", tt.wantCount, len(got), got)
			}

			if tt.wantCount > 0 && got[0] != tt.wantReviewer1 {
				t.Fatalf("expected first reviewer %s, got %s", tt.wantReviewer1, got[0])
			}

			if len(created.AssignedReviewers) != tt.wantCount {
				t.Fatalf("expected %d reviewers in response, got %v", tt.wantCount, created.AssignedReviewers)
			}
		})
	}
}

func TestPRService_CreatePR_RepoError(t *testing.T) {
	ctx := context.Background()

	repoErr := errors.New("boom")
	prRepo := &prRepoFake{createErr: repoErr}
	userRepo := &userRepoFake{
		usersByID: map[string]domain.User{
			"u1": {ID: "u1", TeamName: "team", IsActive: true},
			"u2": {ID: "u2", TeamName: "team", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"team": {
				{ID: "u1", TeamName: "team", IsActive: true},
				{ID: "u2", TeamName: "team", IsActive: true},
			},
		},
	}

	svc := NewPRService(prRepo, userRepo, nil)
	pr := &domain.PullRequest{ID: "pr-fail", Name: "fail", AuthorID: "u1"}

	if _, err := svc.CreatePR(ctx, pr); !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestPRService_ReassignReviewer(t *testing.T) {
	ctx := context.Background()

	prRepo := &prRepoFake{
		prs: map[string]domain.PullRequest{
			"pr-1": {
				ID:       "pr-1",
				Name:     "Reassign",
				AuthorID: "author",
				Status:   domain.PullRequestStatusOpen,
			},
		},
		reviewers: map[string][]string{
			"pr-1": {"u2"},
		},
	}

	userRepo := &userRepoFake{
		usersByID: map[string]domain.User{
			"author": {ID: "author", TeamName: "authors", IsActive: true},
			"u2":     {ID: "u2", TeamName: "reviewers", IsActive: true},
			"u3":     {ID: "u3", TeamName: "reviewers", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"authors": {
				{ID: "author", TeamName: "authors", IsActive: true},
			},
			"reviewers": {
				{ID: "u2", TeamName: "reviewers", IsActive: true},
				{ID: "u3", TeamName: "reviewers", IsActive: true},
			},
		},
	}

	svc := NewPRService(prRepo, userRepo, nil)

	updated, newID, err := svc.ReassignReviewer(ctx, "pr-1", "u2")
	if err != nil {
		t.Fatalf("ReassignReviewer error: %v", err)
	}

	if newID == "" {
		t.Fatalf("expected new reviewer id")
	}

	if len(updated.AssignedReviewers) != 1 {
		t.Fatalf("expected 1 reviewer after reassign, got %v", updated.AssignedReviewers)
	}

	if updated.AssignedReviewers[0] == "u2" {
		t.Fatalf("expected reviewer to change, still %s", updated.AssignedReviewers[0])
	}
}

func TestPRService_ReassignReviewer_ReusesExisting(t *testing.T) {
	ctx := context.Background()

	prRepo := &prRepoFake{
		prs: map[string]domain.PullRequest{
			"pr-small": {
				ID:       "pr-small",
				Name:     "Small team",
				AuthorID: "author",
				Status:   domain.PullRequestStatusOpen,
				AssignedReviewers: []string{
					"u2",
					"u3",
				},
			},
		},
		reviewers: map[string][]string{
			"pr-small": {"u2", "u3"},
		},
	}

	userRepo := &userRepoFake{
		usersByID: map[string]domain.User{
			"author": {ID: "author", TeamName: "backend", IsActive: true},
			"u2":     {ID: "u2", TeamName: "backend", IsActive: true},
			"u3":     {ID: "u3", TeamName: "backend", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"backend": {
				{ID: "author", TeamName: "backend", IsActive: true},
				{ID: "u2", TeamName: "backend", IsActive: true},
				{ID: "u3", TeamName: "backend", IsActive: true},
			},
		},
	}

	svc := NewPRService(prRepo, userRepo, nil)

	updated, newID, err := svc.ReassignReviewer(ctx, "pr-small", "u2")
	if err != nil {
		t.Fatalf("ReassignReviewer error: %v", err)
	}

	if newID != "u3" {
		t.Fatalf("expected new reviewer u3, got %s", newID)
	}

	if len(updated.AssignedReviewers) != 1 || updated.AssignedReviewers[0] != "u3" {
		t.Fatalf("expected reviewers [u3], got %v", updated.AssignedReviewers)
	}
}

func TestPRService_ReassignReviewer_Merged(t *testing.T) {
	ctx := context.Background()

	prRepo := &prRepoFake{
		prs: map[string]domain.PullRequest{
			"pr-merged": {
				ID:       "pr-merged",
				Name:     "Merged",
				AuthorID: "u1",
				Status:   domain.PullRequestStatusMerged,
			},
		},
		reviewers: map[string][]string{
			"pr-merged": {"u2"},
		},
	}

	userRepo := &userRepoFake{
		usersByID: map[string]domain.User{
			"u1": {ID: "u1", TeamName: "team", IsActive: true},
			"u2": {ID: "u2", TeamName: "team", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"team": {
				{ID: "u1", TeamName: "team", IsActive: true},
				{ID: "u2", TeamName: "team", IsActive: true},
			},
		},
	}

	svc := NewPRService(prRepo, userRepo, nil)

	_, _, err := svc.ReassignReviewer(ctx, "pr-merged", "u2")
	if !errors.Is(err, domain.ErrPRMerged) {
		t.Fatalf("expected ErrPRMerged, got %v", err)
	}
}

func TestPRService_Merge(t *testing.T) {
	ctx := context.Background()

	now := time.Now().UTC()
	prRepo := &prRepoFake{
		prs: map[string]domain.PullRequest{
			"pr-merge": {
				ID:        "pr-merge",
				Name:      "Merge test",
				Status:    domain.PullRequestStatusOpen,
				CreatedAt: now,
			},
			"pr-merged": {
				ID:        "pr-merged",
				Name:      "Already merged",
				Status:    domain.PullRequestStatusMerged,
				CreatedAt: now,
			},
		},
	}

	svc := NewPRService(prRepo, nil, nil)

	merged, err := svc.MergePR(ctx, "pr-merge")
	if err != nil {
		t.Fatalf("MergePR error: %v", err)
	}

	if merged.Status != domain.PullRequestStatusMerged {
		t.Fatalf("expected status MERGED, got %s", merged.Status)
	}

	if merged.MergedAt == nil {
		t.Fatalf("expected MergedAt to be set")
	}

	first, err := svc.MergePR(ctx, "pr-merged")
	if err != nil {
		t.Fatalf("first MergePR error: %v", err)
	}

	second, err := svc.MergePR(ctx, "pr-merged")
	if err != nil {
		t.Fatalf("second MergePR error: %v", err)
	}

	if first.Status != domain.PullRequestStatusMerged || second.Status != domain.PullRequestStatusMerged {
		t.Fatalf("expected MERGED status on repeated merge, got %s and %s", first.Status, second.Status)
	}
}
