package service

import (
	"context"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
)

type TeamService struct {
	teamRepo domain.TeamRepository
	userRepo domain.UserRepository
}

func NewTeamService(teamRepo domain.TeamRepository, userRepo domain.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamService) CreateOrUpdateTeam(ctx context.Context, team domain.Team) error {
	if err := s.teamRepo.Create(ctx, team.Name); err != nil {
		return err
	}

	for _, m := range team.Members {
		user := domain.User{
			ID:       m.ID,
			Username: m.Username,
			TeamName: team.Name,
			IsActive: m.IsActive,
		}
		if err := s.userRepo.Upsert(ctx, user); err != nil {
			return err
		}
	}
	return nil
}
