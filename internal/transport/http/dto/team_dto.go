package dto

import "github.com/Mutter0815/pr-reviewer-service/internal/domain"

type TeamMemberDTO struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type TeamRequest struct {
	TeamName string          `json:"team_name" binding:"required"`
	Members  []TeamMemberDTO `json:"members" binding:"required,dive"`
}

type TeamResponse struct {
	Team TeamDTO `json:"team"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

func (r TeamRequest) ToDomain() domain.Team {
	members := make([]domain.TeamMember, 0, len(r.Members))
	for _, m := range r.Members {
		members = append(members, domain.TeamMember{
			ID:       m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	return domain.Team{
		Name:    r.TeamName,
		Members: members,
	}
}

func TeamDTOFromDomain(team domain.Team) TeamDTO {
	members := make([]TeamMemberDTO, 0, len(team.Members))
	for _, m := range team.Members {
		members = append(members, TeamMemberDTO{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	return TeamDTO{
		TeamName: team.Name,
		Members:  members,
	}
}
