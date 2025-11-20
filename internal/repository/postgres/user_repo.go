package postgres

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Upsert(ctx context.Context, u domain.User) error {
	// TODO: реализем позже
	return nil
}
