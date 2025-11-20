package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Postgres {
	return &Postgres{Pool: pool}
}
