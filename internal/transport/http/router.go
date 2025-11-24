package http

import (
	"github.com/Mutter0815/pr-reviewer-service/internal/service"
	"github.com/Mutter0815/pr-reviewer-service/internal/transport/http/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(services *service.Services) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	healthHandler := handlers.NewHealthHandler()
	teamHandler := handlers.NewTeamHandler(services.Team)
	userHandler := handlers.NewUserHandler(services.User)
	prHandler := handlers.NewPRHandler(services.PR)

	r.GET("/health", healthHandler.Health)
	r.POST("/team/add", teamHandler.AddTeam)
	r.GET("/team/list", teamHandler.ListTeams)

	r.GET("/team/get", teamHandler.GetTeamInfo)
	r.POST("/pullRequest/create", prHandler.Create)
	r.POST("/pullRequest/reassign", prHandler.Reassign)
	r.POST("/pullRequest/merge", prHandler.Merge)

	r.POST("/users/setIsActive", userHandler.SetIsActive)
	r.GET("/users/getReview", userHandler.GetReview)
	r.Static("/swagger", "internal/transport/http/swagger")

	return r
}
