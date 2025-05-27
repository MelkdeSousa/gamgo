help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make [target]\033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_/-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
.PHONY:help

install: ## Install project dependencies
	@echo "Installing dependencies..."
	go mod download && go mod tidy
	go install github.com/air-verse/air@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/swaggo/swag/cmd/swag@latest
.PHONY: install

docker-up: ## Start all docker containers
	@echo "Starting docker containers..."
	docker compose up -d
.PHONY: docker-up

docker-down: ## Stop all docker containers, remove volumes and orphaned containers
	@echo "Stopping containers, removing volumes and orphans..."
	docker compose down -v --remove-orphans
.PHONY: docker-down

docker-build: ## Rebuild docker containers
	@echo "Rebuilding docker containers..."
	docker compose build
	docker compose up -d
.PHONY: docker-build

start: ## Start the application without hot-reload
	@echo "Starting application..."
	go run main.go
.PHONY: start

dev: ## Start the application with air for hot-reload
	@echo "Starting application in development mode..."
	air
.PHONY: dev

clean: ## Clean up generated files and docker volumes
	@echo "Cleaning up..."
	docker compose down -v
	rm -rf tmp/
	go clean
.PHONY: clean

db/migration-up: ## Run database migrations
	@goose up
.PHONY: db/migration-up

db/migration-down: ## Rollback database migrations
	@goose down
.PHONY: db/migration-down

db/migration-create: ## Create a new database migration
	@goose create $(name) sql
.PHONY: db/migration-create

db/seeds-up: ## Seed the database
	@goose -no-versioning -dir database/seeds up
.PHONY: db/seeds-up

db/seeds-down: ## Rollback database seeds
	@goose -no-versioning -dir database/seeds down
.PHONY: db/seeds-down

db/seeds-create: ## Create a new database seeder
	@goose -dir database/seeds create $(name) sql
.PHONY: db/seeds-create

docs/swagger-fmt: ## Format Swagger documentation
	@echo "Formatting Swagger documentation..."
	swag fmt
.PHONY: docs/swagger-fmt

docs/swagger-generate: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	swag init -g main.go -o docs/swagger
.PHONY: docs/swagger-generate