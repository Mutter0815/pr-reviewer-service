package handlers

import (
	"net/http"

	"github.com/Mutter0815/pr-reviewer-service/internal/service"
	"github.com/Mutter0815/pr-reviewer-service/internal/transport/http/dto"
	"github.com/Mutter0815/pr-reviewer-service/internal/transport/http/httperror"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) AddTeam(c *gin.Context) {
	var req dto.TeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid request body",
			},
		})
		return
	}

	team := req.ToDomain()

	if err := h.teamService.CreateOrUpdateTeam(c.Request.Context(), team); err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.TeamResponse{
		Team: dto.TeamDTOFromDomain(team),
	}

	c.JSON(http.StatusCreated, resp)
}
func (h *TeamHandler) GetTeamInfo(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "team_name query param is required",
			},
		})
		return
	}

	team, err := h.teamService.GetTeamInfo(c.Request.Context(), teamName)
	if err != nil {
		httperror.Write(c, err)
		return
	}

	resp := dto.TeamResponse{
		Team: dto.TeamDTOFromDomain(team),
	}

	c.JSON(http.StatusOK, resp)
}
