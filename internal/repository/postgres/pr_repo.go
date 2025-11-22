package postgres

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PRRepo struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepo(pool *pgxpool.Pool) *PRRepo {
	return &PRRepo{pool: pool}
}

func (r *PRRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	const query = `
		INSERT INTO pull_requests (
			pull_request_id,
			pull_request_name,
			author_id,
			status,
			created_at,
			merged_at
		)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.pool.Exec(ctx, query,
		pr.ID,
		pr.Name,
		pr.AuthorID,
		string(pr.Status),
		pr.CreatedAt,
		pr.MergedAt,
	)
	return err
}
