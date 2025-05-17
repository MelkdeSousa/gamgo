.PHONY: help install docker-up docker-down docker-build start dev clean db/migration-up db/migration-down db/migration-create db/seeds

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

install: ## Install project dependencies
	@echo "Installing dependencies..."
	go mod download
	go install github.com/air-verse/air@latest

docker-up: ## Start all docker containers
	@echo "Starting docker containers..."
	docker compose up -d

docker-down: ## Stop all docker containers, remove volumes and orphaned containers
	@echo "Stopping containers, removing volumes and orphans..."
	docker compose down -v --remove-orphans

docker-build: ## Rebuild docker containers
	@echo "Rebuilding docker containers..."
	docker compose build
	docker compose up -d

start: ## Start the application without hot-reload
	@echo "Starting application..."
	go run main.go

dev: ## Start the application with air for hot-reload
	@echo "Starting application in development mode..."
	air

clean: ## Clean up generated files and docker volumes
	@echo "Cleaning up..."
	docker compose down -v
	rm -rf tmp/
	go clean

db/migration-up: ## Run database migrations
	@goose up

db/migration-down: ## Rollback database migrations
	@goose down

db/migration-create: ## Create a new database migration
	@goose create $(name) sql

db/seeds-up: ## Seed the database
	@goose -no-versioning -dir database/seeds up

db/seeds-down: ## Rollback database seeds
	@goose -no-versioning -dir database/seeds down

db/seeds-create: ## Create a new database seeder
	@goose -dir database/seeds create $(name) sql
