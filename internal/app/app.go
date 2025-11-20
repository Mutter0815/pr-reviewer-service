package app

import (
	"context"
	"log"
	"time"

	"github.com/Mutter0815/pr-reviewer-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Cfg  *config.Config
	Pool *pgxpool.Pool
}

func New() *App {
	cfg := config.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.PGURL())
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	return &App{
		Cfg:  cfg,
		Pool: pool,
	}
}
func (a *App) Close() {
	a.Pool.Close()
}
