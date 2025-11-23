package httperror

import (
	"errors"
	"net/http"

	"github.com/Mutter0815/pr-reviewer-service/internal/domain"
	"github.com/gin-gonic/gin"
)

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

func New(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetails{
			Code:    code,
			Message: message,
		},
	}
}

func Write(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrTeamExists):
		c.JSON(http.StatusBadRequest, New("TEAM_EXISTS", err.Error()))
	case errors.Is(err, domain.ErrPRExists):
		c.JSON(http.StatusBadRequest, New("PR_EXISTS", err.Error()))
	case errors.Is(err, domain.ErrPRMerged):
		c.JSON(http.StatusBadRequest, New("PR_MERGED", err.Error()))
	case errors.Is(err, domain.ErrNotAssigned):
		c.JSON(http.StatusBadRequest, New("NOT_ASSIGNED", err.Error()))
	case errors.Is(err, domain.ErrNoCandidate):
		c.JSON(http.StatusBadRequest, New("NO_CANDIDATE", err.Error()))
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, New("NOT_FOUND", err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, New("INTERNAL", "internal error"))
	}
}
