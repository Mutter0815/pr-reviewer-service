package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Mutter0815/pr-reviewer-service/internal/config"
	"github.com/Mutter0815/pr-reviewer-service/internal/repository/postgres"
	"github.com/Mutter0815/pr-reviewer-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Cfg      *config.Config
	Pool     *pgxpool.Pool
	Services *service.Services
}

func New() *App {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.PGURL())
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	sql, err := os.ReadFile("migrations/0001_init.up.sql")
	if err != nil {
		log.Fatalf("failed to read migration file: %v", err)
	}

	if _, err := pool.Exec(ctx, string(sql)); err != nil {
		log.Fatalf("failed to apply migration: %v", err)
	}

	teamRepo := postgres.NewTeamRepo(pool)
	userRepo := postgres.NewUserRepo(pool)
	prRepo := postgres.NewPullRequestRepo(pool)

	teamSvc := service.NewTeamService(teamRepo, userRepo)
	userSvc := service.NewUserService(userRepo)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo)

	services := service.NewServices(teamSvc, userSvc, prSvc)

	return &App{
		Cfg:      cfg,
		Pool:     pool,
		Services: services,
	}
}

func (a *App) Close() {
	a.Pool.Close()
}
