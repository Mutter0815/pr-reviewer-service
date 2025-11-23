package postgres

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/jackc/pgx/v5"
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

func (r *TeamRepo) GetByName(ctx context.Context, name string) (domain.Team, error) {
	const queryTeam = `
		SELECT team_name
		FROM teams
		WHERE team_name = $1;
	`

	row := r.pool.QueryRow(ctx, queryTeam, name)

	var team domain.Team
	err := row.Scan(&team.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Team{}, domain.ErrNotFound
		}
		return domain.Team{}, err
	}

	const queryMembers = `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1;
	`

	rows, err := r.pool.Query(ctx, queryMembers, name)
	if err != nil {
		return domain.Team{}, err
	}
	defer rows.Close()

	members := make([]domain.TeamMember, 0)
	for rows.Next() {
		var m domain.TeamMember
		if err := rows.Scan(&m.ID, &m.Username, &m.IsActive); err != nil {
			return domain.Team{}, err
		}
		members = append(members, m)
	}

	team.Members = members

	return team, nil
}
