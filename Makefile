ifneq ($(wildcard .env),)
	include .env
endif

PROJECT_DIR = $(shell CURDIR)
PROJECT_BIN = $(PROJECT_DIR)/bin
PROJECT_TMP = $(PROJECT_DIR)/tmp

install-deps:
	GOBIN=$(PROJECT_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(PROJECT_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


migration-status:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} status -v

migration-add:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} create $(name) sql

migration-up:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

migration-down:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v

gen-sql:
	rm -f internal/repository/product/*.go
	$(PROJECT_BIN)/sqlc generate
docker-build:
	docker build -t notes-app:0.1 .