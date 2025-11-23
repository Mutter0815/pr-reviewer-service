package dto

import "github.com/Mutter0815/pr-reviewer-service/internal/domain"

type SetUserIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User UserDTO `json:"user"`
}

func UserDTOFromDomain(u domain.User) UserDTO {
	return UserDTO{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}
