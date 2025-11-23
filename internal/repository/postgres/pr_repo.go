package postgres

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/jackc/pgx/v5"
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

func (r *PullRequestRepo) GetByID(ctx context.Context, id string) (domain.PullRequest, error) {
	const query = `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1;
	`

	var pr domain.PullRequest
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.PullRequest{}, domain.ErrNotFound
		}
		return domain.PullRequest{}, err
	}

	return pr, nil
}

func (r *PullRequestRepo) ListReviewers(ctx context.Context, prID string) ([]string, error) {
	const query = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		res = append(res, id)
	}

	return res, nil
}

func (r *PullRequestRepo) ReassignReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error {
	const query = `
		UPDATE pull_request_reviewers
		SET reviewer_id = $3
		WHERE pull_request_id = $1 AND reviewer_id = $2;
	`

	cmd, err := r.pool.Exec(ctx, query, prID, oldReviewerID, newReviewerID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrNotAssigned
	}

	return nil
}

func (r *PullRequestRepo) AssignReviewers(ctx context.Context, prID string, reviewerIDs []string) error {
	const query = `
		INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	for _, reviewerID := range reviewerIDs {
		if _, err := r.pool.Exec(ctx, query, prID, reviewerID); err != nil {
			return err
		}
	}

	return nil
}

func (r *PullRequestRepo) Merge(ctx context.Context, prID string) error {
	const query = `
		UPDATE pull_requests
		SET status = 'MERGED',
		    merged_at = COALESCE(merged_at, now())
		WHERE pull_request_id = $1;
	`

	cmd, err := r.pool.Exec(ctx, query, prID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}
