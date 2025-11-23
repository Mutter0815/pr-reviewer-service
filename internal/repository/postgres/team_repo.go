package postgres

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{pool: pool}
}

func (r *TeamRepo) Create(ctx context.Context, name string) error {
	const query = `
		INSERT INTO teams (team_name)
		VALUES ($1)
		ON CONFLICT DO NOTHING;
	`

	cmdTag, err := r.pool.Exec(ctx, query, name)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrTeamExists
	}

	return nil
}
