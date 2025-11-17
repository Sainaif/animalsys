.PHONY: help build up down logs clean restart ps backend-shell frontend-shell db-shell test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d
	@echo "Services started!"
	@echo "Frontend: http://localhost"
	@echo "Backend: http://localhost/api"
	@echo "MongoDB: localhost:27017"
	@echo "Redis: localhost:6379"

down: ## Stop all services
	docker-compose down

logs: ## Show logs from all services
	docker-compose logs -f

logs-backend: ## Show logs from backend only
	docker-compose logs -f backend

logs-frontend: ## Show logs from frontend only
	docker-compose logs -f frontend

clean: ## Stop services and remove volumes
	docker-compose down -v
	rm -rf backend/tmp
	@echo "Cleaned up!"

restart: down up ## Restart all services

ps: ## Show running services
	docker-compose ps

backend-shell: ## Open shell in backend container
	docker-compose exec backend sh

frontend-shell: ## Open shell in frontend container
	docker-compose exec frontend sh

db-shell: ## Open MongoDB shell
	docker-compose exec mongodb mongosh animalsys

redis-cli: ## Open Redis CLI
	docker-compose exec redis redis-cli

test-backend: ## Run backend tests
	cd backend && go test -v ./...

test-frontend: ## Run frontend tests
	cd frontend && npm test

dev: ## Start development environment
	@echo "Starting development environment..."
	docker-compose up

prod-build: ## Build production images
	docker-compose -f docker-compose.prod.yml build

prod-up: ## Start production environment
	docker-compose -f docker-compose.prod.yml up -d

backup-db: ## Backup MongoDB database
	@echo "Creating database backup..."
	docker-compose exec -T mongodb mongodump --db=animalsys --archive > backup_$(shell date +%Y%m%d_%H%M%S).archive
	@echo "Backup created!"

restore-db: ## Restore MongoDB database (usage: make restore-db FILE=backup.archive)
	@echo "Restoring database from $(FILE)..."
	docker-compose exec -T mongodb mongorestore --db=animalsys --archive < $(FILE)
	@echo "Database restored!"

init: ## Initialize project (copy env, build, start)
	cp .env.example .env
	@echo "Created .env file. Please update it with your settings."
	make build
	make up
	@echo "Project initialized! Visit http://localhost"

seed: ## Seed database with initial admin user
	@echo "Seeding database..."
	docker-compose exec backend go run ./cmd/seed
	@echo "Database seeded!"

reseed: ## Drop the Mongo database and run the seed script
	./scripts/reseed.sh
