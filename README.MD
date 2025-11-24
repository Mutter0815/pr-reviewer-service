# Avito-Backend-trainee-assignment-autumn-2025

## **Запуск**

```bash
cp .env.example .env
```
### 1. Запуск всего проекта одной командой:
```bash
docker compose up //# запуск сервиса и e2e
```
После этого сервис доступен на[http://localhost:8080](http://localhost:8080)
**Swagger**
- Документация: internal/transport/http/swagger/openapi.yml
- UI в браузере:[http://localhost:8080/swagger](http://localhost:8080/swagger)
### 2.Запуск тестов
```go
go test ./tests/...
```
**Что я подметил и поправил по ТЗ**

- Назначение ревьюверов:
    - выбираем **случайных** активных ревьюверов из команды автора (автор не попадает в список);
    - при 0/1 кандидате PR создаётся с 0/1 ревьювером.
- Переназначение:
    - кандидаты берутся из **команды заменяемого ревьювера**, а не автора;
    - если команда маленькая (2 ревьювера), можно заменить одного другим (чтобы сценарий из ТЗ не ломался).
- /team/add:
    - первый вызов — 201;
    - повторный с тем же team_name — 400 + TEAM_EXISTS, как в OpenAPI.
- merge:
    - операция идемпотентна;
    - после MERGED любые reassign дают PR_MERGED и не меняют состав ревьюверов.

**Что сделал из доп. заданий**

- Написал интеграционный тест HTTP‑слоя (internal/transport/http/integration_test.go) — поднимается реальный роутер и гоняются запросы через httptest.
- Сделал e2e‑сценарий в Docker:
    - в docker-compose.yml есть сервис e2e, который после старта app выполняет go run ./tests/e2e и проверяет полный поток: /team/add → /pullRequest/create → /pullRequest/reassign → /pullRequest/merge → повторный reassign → /users/getReview.
### **Стек решения**
- Golang
- PostgreSQL
- Docker

### **Библиотеки, использованные в задании**
- gin (роутер и HTTP‑обработчики)
- pgx / pgxpool (подключение к PostgreSQL)
- стандартная библиотека testing + net/http/httptest (юнит‑ и интеграционные тесты)



