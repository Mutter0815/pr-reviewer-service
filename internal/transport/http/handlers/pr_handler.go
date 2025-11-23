package handlers

import (
	"net/http"

	"github.com/Mutter0815/pr-reviewer-service/internal/service"
	"github.com/Mutter0815/pr-reviewer-service/internal/transport/http/dto"
	"github.com/Mutter0815/pr-reviewer-service/internal/transport/http/httperror"
	"github.com/gin-gonic/gin"
)

type PRHandler struct {
	prService *service.PRService
}

func NewPRHandler(prService *service.PRService) *PRHandler {
	return &PRHandler{
		prService: prService,
	}
}

func (h *PRHandler) Create(c *gin.Context) {
	var req dto.PRCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid request body",
			},
		})
		return
	}

	pr := req.ToDomain()

	if err := h.prService.CreatePR(c.Request.Context(), pr); err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.PRCreateResponse{
		PullRequest: dto.PRDTOFromDomain(*pr),
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *PRHandler) Reassign(c *gin.Context) {
	var req dto.PRReassignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid request body",
			},
		})
		return
	}

	pr, newReviewerID, err := h.prService.ReassignReviewer(
		c.Request.Context(),
		req.PullRequestID,
		req.ReviewerID,
	)
	if err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.PRReassignResponse{
		PullRequest:   dto.PRDTOFromDomain(pr),
		NewReviewerID: newReviewerID,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PRHandler) Merge(c *gin.Context) {
	var req dto.PRMergeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid request body",
			},
		})
		return
	}

	pr, err := h.prService.MergePR(c.Request.Context(), req.PullRequestID)
	if err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.PRMergeResponse{
		PR: dto.PRDTOFromDomain(pr),
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PRHandler) ListByReviewer(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "user_id query param is required",
			},
		})
		return
	}

	prs, err := h.prService.ListByReviewer(c.Request.Context(), userID)
	if err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.PRListByUserResponse{
		UserID:       userID,
		PullRequests: make([]dto.PRShortDTO, 0, len(prs)),
	}

	for _, pr := range prs {
		resp.PullRequests = append(resp.PullRequests, dto.PRShortDTOFromDomain(pr))
	}

	c.JSON(http.StatusOK, resp)
}
