MIGRATIONS_DIR := internal/infrastructure/db/migrations

ifneq (,$(wildcard .env))
    include .env
    export
endif

PG_DSN := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: migrate-generate
migrate-generate:
	$(GOPATH)/bin/goose -dir $(MIGRATIONS_DIR) -s create $(name) sql

.PHONY: migrate-up
migrate-up:
	docker compose run --rm migrate -dir $(MIGRATIONS_DIR) postgres "$(PG_DSN)" up

.PHONY: migrate-down
migrate-down:
	docker compose run --rm migrate -dir $(MIGRATIONS_DIR) postgres "$(PG_DSN)" down

.PHONY: migrate-status
migrate-status:
	docker compose run --rm migrate -dir $(MIGRATIONS_DIR) postgres "$(PG_DSN)" status