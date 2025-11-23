package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
)

type fakeTeamRepo struct {
	createdNames []string

	createErr   error
	getByNameFn func(ctx context.Context, name string) (domain.Team, error)
}

func (r *fakeTeamRepo) Create(ctx context.Context, name string) error {
	r.createdNames = append(r.createdNames, name)
	return r.createErr
}

func (r *fakeTeamRepo) GetByName(ctx context.Context, name string) (domain.Team, error) {
	if r.getByNameFn != nil {
		return r.getByNameFn(ctx, name)
	}
	return domain.Team{}, domain.ErrNotFound
}

type fakeUserRepo struct {
	upserted  []domain.User
	upsertErr error
}

func (r *fakeUserRepo) Upsert(ctx context.Context, u domain.User) error {
	r.upserted = append(r.upserted, u)
	return r.upsertErr
}

func TestTeamService_CreateOrUpdateTeam_Success(t *testing.T) {
	ctx := context.Background()

	teamRepo := &fakeTeamRepo{}
	userRepo := &fakeUserRepo{}

	svc := NewTeamService(teamRepo, userRepo)

	team := domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{ID: "u1", Username: "Alice", IsActive: true},
			{ID: "u2", Username: "Bob", IsActive: false},
		},
	}

	err := svc.CreateOrUpdateTeam(ctx, team)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(teamRepo.createdNames) != 1 {
		t.Fatalf("expected 1 call to Create, got %d", len(teamRepo.createdNames))
	}
	if teamRepo.createdNames[0] != "backend" {
		t.Errorf("expected team name 'backend', got %q", teamRepo.createdNames[0])
	}

	if len(userRepo.upserted) != len(team.Members) {
		t.Fatalf("expected %d calls to Upsert, got %d", len(team.Members), len(userRepo.upserted))
	}

	for i, m := range team.Members {
		u := userRepo.upserted[i]
		if u.ID != m.ID {
			t.Errorf("user[%d]: expected ID %q, got %q", i, m.ID, u.ID)
		}
		if u.Username != m.Username {
			t.Errorf("user[%d]: expected Username %q, got %q", i, m.Username, u.Username)
		}
		if u.TeamName != team.Name {
			t.Errorf("user[%d]: expected TeamName %q, got %q", i, team.Name, u.TeamName)
		}
		if u.IsActive != m.IsActive {
			t.Errorf("user[%d]: expected IsActive %v, got %v", i, m.IsActive, u.IsActive)
		}
	}
}

func TestTeamService_CreateOrUpdateTeam_TeamExists(t *testing.T) {
	ctx := context.Background()

	teamRepo := &fakeTeamRepo{
		createErr: domain.ErrTeamExists,
	}
	userRepo := &fakeUserRepo{}

	svc := NewTeamService(teamRepo, userRepo)

	team := domain.Team{
		Name: "backend",
	}

	err := svc.CreateOrUpdateTeam(ctx, team)
	if !errors.Is(err, domain.ErrTeamExists) {
		t.Fatalf("expected ErrTeamExists, got %v", err)
	}

	if len(userRepo.upserted) != 0 {
		t.Fatalf("expected 0 calls to Upsert, got %d", len(userRepo.upserted))
	}
}

func TestTeamService_CreateOrUpdateTeam_UpsertError(t *testing.T) {
	ctx := context.Background()

	teamRepo := &fakeTeamRepo{}
	upsertErr := errors.New("upsert failed")
	userRepo := &fakeUserRepo{
		upsertErr: upsertErr,
	}

	svc := NewTeamService(teamRepo, userRepo)

	team := domain.Team{
		Name: "backend",
		Members: []domain.TeamMember{
			{ID: "u1", Username: "Alice", IsActive: true},
		},
	}

	err := svc.CreateOrUpdateTeam(ctx, team)
	if !errors.Is(err, upsertErr) {
		t.Fatalf("expected upsertErr, got %v", err)
	}

	if len(teamRepo.createdNames) != 1 {
		t.Fatalf("expected 1 call to Create, got %d", len(teamRepo.createdNames))
	}
}
