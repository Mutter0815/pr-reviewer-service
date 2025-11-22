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
	// TODO: реализуем позже
	return nil
}
