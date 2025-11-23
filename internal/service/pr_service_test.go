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
}

func (r *fakePRRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	r.created = append(r.created, pr)
	return r.createErr
}

func TestPRService_CreatePR_Success(t *testing.T) {
	ctx := context.Background()

	prRepo := &fakePRRepo{}
	svc := NewPRService(prRepo, nil, nil)

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
}

func TestPRService_CreatePR_RepoReturnsPRExists(t *testing.T) {
	ctx := context.Background()

	prRepo := &fakePRRepo{
		createErr: domain.ErrPRExists,
	}
	svc := NewPRService(prRepo, nil, nil)

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
	svc := NewPRService(prRepo, nil, nil)

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
	svc := NewPRService(prRepo, nil, nil)

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
