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
