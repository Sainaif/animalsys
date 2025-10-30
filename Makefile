.PHONY: help install dev build test clean deploy docker-up docker-down docker-rebuild docker-logs backup restore

# Default target
help:
	@echo "AnimalSys - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make install        - Install all dependencies"
	@echo "  make dev            - Run development servers"
	@echo "  make build          - Build for production"
	@echo "  make test           - Run all tests"
	@echo "  make test-backend   - Run backend tests"
	@echo "  make test-frontend  - Run frontend tests"
	@echo "  make lint           - Lint code"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy         - Deploy with Docker (one command)"
	@echo "  make deploy-prod    - Deploy in production mode"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-rebuild - Rebuild and restart containers"
	@echo "  make docker-logs    - View Docker logs"
	@echo ""
	@echo "Database:"
	@echo "  make backup         - Create database backup"
	@echo "  make restore        - Restore database from backup"
	@echo ""

# Install dependencies
install:
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Done!"

# Run development servers
dev:
	@echo "Starting development environment..."
	docker-compose up -d mongodb
	@echo "Waiting for MongoDB to be ready..."
	@sleep 3
	@make -j2 dev-backend dev-frontend

dev-backend:
	cd backend && air

dev-frontend:
	cd frontend && npm run dev

# Build for production
build:
	@echo "Building backend..."
	cd backend && go build -o bin/api ./cmd/api
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Build complete!"

# Run tests
test: test-backend test-frontend

test-backend:
	@echo "Running backend tests..."
	cd backend && go test ./... -v -cover -coverprofile=coverage.out

test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm run test:coverage

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin backend/tmp backend/coverage.out
	rm -rf frontend/dist frontend/coverage
	@echo "Clean complete!"

# Deployment commands
deploy:
	@chmod +x deployment/deploy.sh
	@./deployment/deploy.sh

deploy-prod:
	@chmod +x deployment/deploy.sh
	@./deployment/deploy.sh --production

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

docker-logs:
	docker-compose logs -f

# Database backup
backup:
	@echo "Creating database backup..."
	@mkdir -p backups
	docker-compose exec -T mongodb mongodump --username=$${MONGO_USERNAME} --password=$${MONGO_PASSWORD} --authenticationDatabase=admin --db=$${MONGO_DATABASE} --archive > backups/backup_$$(date +%Y%m%d_%H%M%S).archive
	@echo "Backup complete!"

# Lint code
lint:
	@echo "Linting backend..."
	cd backend && go fmt ./... && go vet ./...
	@echo "Linting frontend..."
	cd frontend && npm run lint
