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
	const query = `
		INSERT INTO pull_requests (
			pull_request_id,
			pull_request_name,
			author_id,
			status,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING;
	`

	cmd, err := r.pool.Exec(ctx, query,
		pr.ID,
		pr.Name,
		pr.AuthorID,
		pr.Status,
		pr.CreatedAt,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrPRExists
	}

	return nil
}
