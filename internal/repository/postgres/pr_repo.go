package postgres

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepo struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepo(pool *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{pool: pool}
}

func (r *PullRequestRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	// TODO: реализуем позже
	return nil
}
