package handlers

import (
	"net/http"

	"github.com/Mutter0815/pr-reviewer-service/internal/service"
	"github.com/Mutter0815/pr-reviewer-service/internal/transport/http/dto"
	"github.com/Mutter0815/pr-reviewer-service/internal/transport/http/httperror"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SetIsActive(c *gin.Context) {
	var req dto.SetUserIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid request body",
			},
		})
		return
	}

	user, err := h.userService.SetIsActive(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.UserResponse{
		User: dto.UserDTOFromDomain(user),
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetReview(c *gin.Context) {
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

	prs, err := h.userService.ListReviewerPRs(c.Request.Context(), userID)
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
