package postgres

import (
	"context"
	"errors"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Upsert(ctx context.Context, u domain.User) error {
	const query = `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			username  = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active;
	`

	_, err := r.pool.Exec(ctx, query, u.ID, u.Username, u.TeamName, u.IsActive)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	const query = `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1;
	`

	var u domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
		&u.TeamName,
		&u.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	return u, nil
}

func (r *UserRepo) ListActiveByTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	const query = `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE team_name = $1 AND is_active = TRUE;
	`

	rows, err := r.pool.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepo) SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	const query = `
		UPDATE users
		SET is_active = $2
		WHERE user_id = $1
		RETURNING user_id, username, team_name, is_active;
	`

	var u domain.User
	err := r.pool.QueryRow(ctx, query, userID, isActive).
		Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	return u, nil
}
