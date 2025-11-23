package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
)

type fakePRRepo struct {
	created   []*domain.PullRequest
	createErr error

	assigned map[string][]string
}

func (r *fakePRRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	r.created = append(r.created, pr)
	return r.createErr
}

func (r *fakePRRepo) AssignReviewers(ctx context.Context, prID string, reviewerIDs []string) error {
	if r.assigned == nil {
		r.assigned = make(map[string][]string)
	}
	r.assigned[prID] = append(r.assigned[prID], reviewerIDs...)
	return nil
}

func (r *fakePRRepo) GetByID(ctx context.Context, id string) (domain.PullRequest, error) {
	return domain.PullRequest{}, nil
}

func (r *fakePRRepo) ListReviewers(ctx context.Context, prID string) ([]string, error) {
	return nil, nil
}

func (r *fakePRRepo) ReassignReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error {
	return nil
}

type fakeUserRepo struct {
	usersByID    map[string]domain.User
	activeByTeam map[string][]domain.User

	upserted  []domain.User
	upsertErr error
}

func (r *fakeUserRepo) Upsert(ctx context.Context, u domain.User) error {
	r.upserted = append(r.upserted, u)
	return r.upsertErr
}

func (r *fakeUserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	u, ok := r.usersByID[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}

func (r *fakeUserRepo) ListActiveByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	return r.activeByTeam[teamName], nil
}

func TestPRService_CreatePR_Success(t *testing.T) {
	ctx := context.Background()

	prRepo := &fakePRRepo{}
	userRepo := &fakeUserRepo{
		usersByID: map[string]domain.User{
			"u1": {ID: "u1", TeamName: "backend", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"backend": {
				{ID: "u1", TeamName: "backend", IsActive: true},
				{ID: "u2", TeamName: "backend", IsActive: true},
			},
		},
	}

	svc := NewPRService(prRepo, userRepo, nil)

	pr := &domain.PullRequest{
		ID:       "pr-1",
		Name:     "Implement feature X",
		AuthorID: "u1",
	}

	err := svc.CreatePR(ctx, pr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(prRepo.created) != 1 {
		t.Fatalf("expected 1 call to Create, got %d", len(prRepo.created))
	}

	if prRepo.created[0] != pr {
		t.Fatalf("expected repo to receive the same PR pointer")
	}

	if pr.Status != domain.PullRequestStatusOpen {
		t.Errorf("expected status %q, got %q", domain.PullRequestStatusOpen, pr.Status)
	}

	if pr.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set, got zero value")
	}

	assigned := prRepo.assigned["pr-1"]
	if len(assigned) != 1 || assigned[0] != "u2" {
		t.Fatalf("expected one reviewer u2, got %v", assigned)
	}
}

func TestPRService_CreatePR_RepoReturnsPRExists(t *testing.T) {
	ctx := context.Background()

	prRepo := &fakePRRepo{
		createErr: domain.ErrPRExists,
	}
	userRepo := &fakeUserRepo{
		usersByID: map[string]domain.User{
			"u1": {ID: "u1", TeamName: "backend", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"backend": {
				{ID: "u1", TeamName: "backend", IsActive: true},
				{ID: "u2", TeamName: "backend", IsActive: true},
			},
		},
	}

	svc := NewPRService(prRepo, userRepo, nil)

	pr := &domain.PullRequest{
		ID:       "pr-1",
		Name:     "Implement feature X",
		AuthorID: "u1",
	}

	err := svc.CreatePR(ctx, pr)
	if !errors.Is(err, domain.ErrPRExists) {
		t.Fatalf("expected ErrPRExists, got %v", err)
	}

	if len(prRepo.created) != 1 {
		t.Fatalf("expected 1 call to Create, got %d", len(prRepo.created))
	}

	if pr.Status != domain.PullRequestStatusOpen {
		t.Errorf("expected status %q, got %q", domain.PullRequestStatusOpen, pr.Status)
	}
	if pr.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set, got zero value")
	}
}

func TestPRService_CreatePR_KeepExistingCreatedAt(t *testing.T) {
	ctx := context.Background()

	prRepo := &fakePRRepo{}
	userRepo := &fakeUserRepo{
		usersByID: map[string]domain.User{
			"u2": {ID: "u2", TeamName: "backend", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"backend": {
				{ID: "u2", TeamName: "backend", IsActive: true},
				{ID: "u3", TeamName: "backend", IsActive: true},
			},
		},
	}

	svc := NewPRService(prRepo, userRepo, nil)

	createdAt := time.Date(2025, 11, 23, 10, 0, 0, 0, time.UTC)

	pr := &domain.PullRequest{
		ID:        "pr-2",
		Name:      "Bugfix",
		AuthorID:  "u2",
		CreatedAt: createdAt,
	}

	err := svc.CreatePR(ctx, pr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !pr.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt to stay %v, got %v", createdAt, pr.CreatedAt)
	}
}

func TestPRService_CreatePR_RepoError(t *testing.T) {
	ctx := context.Background()

	repoErr := errors.New("db is down")
	prRepo := &fakePRRepo{
		createErr: repoErr,
	}
	userRepo := &fakeUserRepo{
		usersByID: map[string]domain.User{
			"u3": {ID: "u3", TeamName: "backend", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"backend": {
				{ID: "u3", TeamName: "backend", IsActive: true},
				{ID: "u4", TeamName: "backend", IsActive: true},
			},
		},
	}

	svc := NewPRService(prRepo, userRepo, nil)

	pr := &domain.PullRequest{
		ID:       "pr-3",
		Name:     "Refactor",
		AuthorID: "u3",
	}

	err := svc.CreatePR(ctx, pr)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got %v", err)
	}

	if len(prRepo.created) != 1 {
		t.Fatalf("expected 1 call to Create, got %d", len(prRepo.created))
	}
}

func TestPRService_AssignTwoReviewers(t *testing.T) {
	ctx := context.Background()

	userRepo := &fakeUserRepo{
		usersByID: map[string]domain.User{
			"u1": {ID: "u1", Username: "Author", TeamName: "backend", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"backend": {
				{ID: "u1", Username: "Author", TeamName: "backend", IsActive: true},
				{ID: "u2", Username: "Rev1", TeamName: "backend", IsActive: true},
				{ID: "u3", Username: "Rev2", TeamName: "backend", IsActive: true},
			},
		},
	}

	prRepo := &fakePRRepo{}
	svc := NewPRService(prRepo, userRepo, nil)

	pr := &domain.PullRequest{
		ID:       "pr-100",
		Name:     "Test PR",
		AuthorID: "u1",
	}

	if err := svc.CreatePR(ctx, pr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assigned := prRepo.assigned["pr-100"]
	if len(assigned) != 2 {
		t.Fatalf("expected 2 reviewers, got %d", len(assigned))
	}

	if assigned[0] != "u2" || assigned[1] != "u3" {
		t.Fatalf("expected reviewers [u2 u3], got %v", assigned)
	}
}

func TestPRService_AssignOneReviewer(t *testing.T) {
	ctx := context.Background()

	userRepo := &fakeUserRepo{
		usersByID: map[string]domain.User{
			"a1": {ID: "a1", Username: "Author", TeamName: "small", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"small": {
				{ID: "a1", Username: "Author", TeamName: "small", IsActive: true},
				{ID: "a2", Username: "OnlyRev", TeamName: "small", IsActive: true},
			},
		},
	}

	prRepo := &fakePRRepo{}
	svc := NewPRService(prRepo, userRepo, nil)

	pr := &domain.PullRequest{
		ID:       "pr-one",
		Name:     "Small team PR",
		AuthorID: "a1",
	}

	if err := svc.CreatePR(ctx, pr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assigned := prRepo.assigned["pr-one"]
	if len(assigned) != 1 {
		t.Fatalf("expected 1 reviewer, got %d", len(assigned))
	}

	if assigned[0] != "a2" {
		t.Fatalf("expected reviewer a2, got %s", assigned[0])
	}
}

func TestPRService_NoCandidates(t *testing.T) {
	ctx := context.Background()

	userRepo := &fakeUserRepo{
		usersByID: map[string]domain.User{
			"s1": {ID: "s1", Username: "Solo", TeamName: "solo", IsActive: true},
		},
		activeByTeam: map[string][]domain.User{
			"solo": {
				{ID: "s1", Username: "Solo", TeamName: "solo", IsActive: true},
			},
		},
	}

	prRepo := &fakePRRepo{}
	svc := NewPRService(prRepo, userRepo, nil)

	pr := &domain.PullRequest{
		ID:       "pr-none",
		Name:     "Solo PR",
		AuthorID: "s1",
	}

	err := svc.CreatePR(ctx, pr)
	if !errors.Is(err, domain.ErrNoCandidate) {
		t.Fatalf("expected ErrNoCandidate, got %v", err)
	}

	if len(prRepo.assigned) != 0 {
		t.Fatalf("expected no AssignReviewers calls, got %d", len(prRepo.assigned))
	}
}
