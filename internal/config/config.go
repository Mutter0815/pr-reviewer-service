package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBHost     string `env:"DB_HOST"     envDefault:"localhost"`
	DBPort     string `env:"DB_PORT"     envDefault:"5432"`
	DBUser     string `env:"DB_USER"     envDefault:"pr_user"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"pr_pass"`
	DBName     string `env:"DB_NAME"     envDefault:"pr_db"`
}

func Load() *Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse env config: %v", err)
	}
	return &cfg
}
func (c *Config) PGURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
