package dto

import "github.com/Mutter0815/pr-reviewer-service/internal/domain"

type PRCreateRequest struct {
	ID       string `json:"pull_request_id"   binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorID string `json:"author_id"         binding:"required"`
}

type PRDTO struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

type PRCreateResponse struct {
	PullRequest PRDTO `json:"pull_request"`
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
		ID:       pr.ID,
		Name:     pr.Name,
		AuthorID: pr.AuthorID,
		Status:   string(pr.Status),
	}
}
