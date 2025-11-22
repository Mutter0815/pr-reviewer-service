package postgres

import (
	"context"

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
		ON CONFLICT (team_name) DO NOTHING;
	`

	_, err := r.pool.Exec(ctx, query, name)
	return err
}
