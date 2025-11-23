package handlers

import (
	"net/http"

	"github.com/Mutter0815/pr-reviewer-service/internal/service"
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

func (h *PRHandler) Stub(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}
