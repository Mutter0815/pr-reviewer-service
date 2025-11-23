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

	created, err := h.prService.CreatePR(c.Request.Context(), pr)
	if err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.PRCreateResponse{
		PR: dto.PRDTOFromDomain(created),
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
		req.OldUserID,
	)
	if err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.PRReassignResponse{
		PR:         dto.PRDTOFromDomain(pr),
		ReplacedBy: newReviewerID,
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
