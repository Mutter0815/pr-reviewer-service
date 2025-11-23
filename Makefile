# Имя бинарника
APP_NAME ?= pr-reviewer-service

# Путь к main-файлу (поменяй, если у тебя другой)
CMD_DIR ?= ./cmd/app

# docker compose (если у тебя старая версия, поменяй на "docker-compose")
COMPOSE ?= docker compose

# Имена сервисов из docker-compose.yml
APP_SERVICE ?= app      # сервис с нашим Go-приложением
DB_SERVICE  ?= db       # сервис с Postgres

.PHONY: build run up down stop restart logs logs-app db-shell db-logs clean

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

logs-app:
	$(COMPOSE) logs -f $(APP_SERVICE)

db-logs:
	$(COMPOSE) logs -f $(DB_SERVICE)

db-shell:
	$(COMPOSE) exec $(DB_SERVICE) psql -U pr_user -d pr_db

