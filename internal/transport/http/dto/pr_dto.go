package dto

import (
	"time"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
)

type PRCreateRequest struct {
	ID       string `json:"pull_request_id"   binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorID string `json:"author_id"         binding:"required"`
}

type PRDTO struct {
	ID                string     `json:"pull_request_id"`
	Name              string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type PRCreateResponse struct {
	PR PRDTO `json:"pr"`
}

func (r PRCreateRequest) ToDomain() *domain.PullRequest {
	return &domain.PullRequest{
		ID:       r.ID,
		Name:     r.Name,
		AuthorID: r.AuthorID,
		// статус и время поставим в сервисе
	}
}

func PRDTOFromDomain(pr domain.PullRequest) PRDTO {
	return PRDTO{
		ID:                pr.ID,
		Name:              pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: append([]string(nil), pr.AssignedReviewers...),
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

type PRReassignRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_user_id"     binding:"required"`
}

type PRReassignResponse struct {
	PR         PRDTO  `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}

type PRMergeRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

type PRMergeResponse struct {
	PR PRDTO `json:"pr"`
}

type PRShortDTO struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

func PRShortDTOFromDomain(pr domain.PullRequest) PRShortDTO {
	return PRShortDTO{
		ID:       pr.ID,
		Name:     pr.Name,
		AuthorID: pr.AuthorID,
		Status:   string(pr.Status),
	}
}

type PRListByUserResponse struct {
	UserID       string       `json:"user_id"`
	PullRequests []PRShortDTO `json:"pull_requests"`
}
