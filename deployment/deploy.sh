#!/bin/bash

# =================================================================
# AnimalSys One-Command Deployment Script
# =================================================================
# This script automates the deployment of AnimalSys ERP system
# Usage: ./deployment/deploy.sh [options]
# Options:
#   --production    Deploy with production profile (includes Nginx)
#   --rebuild       Force rebuild of Docker images
#   --no-seed       Skip database seeding
# =================================================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="AnimalSys"
BACKEND_HEALTH_URL="http://localhost:8080/health"
FRONTEND_URL="http://localhost:3000"
MAX_RETRIES=30
RETRY_INTERVAL=2

# Parse command line arguments
PRODUCTION=false
REBUILD=false
NO_SEED=false

for arg in "$@"; do
    case $arg in
        --production)
            PRODUCTION=true
            shift
            ;;
        --rebuild)
            REBUILD=true
            shift
            ;;
        --no-seed)
            NO_SEED=true
            shift
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  --production    Deploy with production profile (includes Nginx)"
            echo "  --rebuild       Force rebuild of Docker images"
            echo "  --no-seed       Skip database seeding"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $arg${NC}"
            exit 1
            ;;
    esac
done

# Function to print colored messages
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo ""
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}======================================${NC}"
    echo ""
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"

    # Check Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    print_success "Docker found: $(docker --version)"

    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    print_success "Docker Compose found"

    # Check if Docker daemon is running
    if ! docker info &> /dev/null; then
        print_error "Docker daemon is not running. Please start Docker first."
        exit 1
    fi
    print_success "Docker daemon is running"
}

# Setup environment
setup_environment() {
    print_header "Setting Up Environment"

    cd "$(dirname "$0")/.."  # Move to project root

    if [ ! -f ".env" ]; then
        print_info "Creating .env file from .env.example..."
        cp .env.example .env
        print_warning "Please review and update .env file with your configuration!"
        print_warning "Especially change: MONGO_ROOT_PASSWORD and JWT_SECRET"
    else
        print_success ".env file already exists"
    fi
}

# Build Docker images
build_images() {
    print_header "Building Docker Images"

    if [ "$REBUILD" = true ]; then
        print_info "Rebuilding images without cache..."
        docker-compose build --no-cache
    else
        print_info "Building images..."
        docker-compose build
    fi

    print_success "Docker images built successfully"
}

# Start services
start_services() {
    print_header "Starting Services"

    if [ "$PRODUCTION" = true ]; then
        print_info "Starting services with production profile..."
        docker-compose --profile production up -d
    else
        print_info "Starting services in development mode..."
        docker-compose up -d
    fi

    print_success "Services started"
}

# Wait for service to be healthy
wait_for_service() {
    local service_name=$1
    local health_check_cmd=$2
    local max_wait=$MAX_RETRIES
    local count=0

    print_info "Waiting for $service_name to be ready..."

    while [ $count -lt $max_wait ]; do
        if eval $health_check_cmd &> /dev/null; then
            print_success "$service_name is ready"
            return 0
        fi

        count=$((count + 1))
        echo -n "."
        sleep $RETRY_INTERVAL
    done

    echo ""
    print_error "$service_name failed to start within $(($max_wait * $RETRY_INTERVAL)) seconds"
    return 1
}

# Wait for all services
wait_for_services() {
    print_header "Waiting for Services to be Ready"

    # Wait for MongoDB
    wait_for_service "MongoDB" "docker-compose exec -T mongodb mongosh --eval 'db.adminCommand(\"ping\")'"

    # Wait for Backend
    wait_for_service "Backend" "curl -f $BACKEND_HEALTH_URL"

    # Wait for Frontend
    wait_for_service "Frontend" "curl -f http://localhost:3000/health"

    print_success "All services are ready!"
}

# Seed database
seed_database() {
    if [ "$NO_SEED" = true ]; then
        print_info "Skipping database seeding (--no-seed flag)"
        return 0
    fi

    print_header "Seeding Database"

    print_info "Running database seed..."
    if docker-compose exec -T backend ./seed &> /dev/null; then
        print_success "Database seeded successfully"
    else
        print_warning "Database seeding failed or seed command not found"
        print_info "You may need to run it manually: docker-compose exec backend ./seed"
    fi
}

# Display access information
display_access_info() {
    print_header "Deployment Complete!"

    echo -e "${GREEN}✓${NC} ${PROJECT_NAME} is now running!"
    echo ""
    echo -e "${BLUE}Access URLs:${NC}"
    echo -e "  Frontend:    ${GREEN}http://localhost:3000${NC}"
    echo -e "  Backend API: ${GREEN}http://localhost:8080${NC}"
    echo -e "  API Health:  ${GREEN}http://localhost:8080/health${NC}"
    echo ""

    if [ "$PRODUCTION" = true ]; then
        echo -e "${BLUE}Production Mode:${NC}"
        echo -e "  Nginx:       ${GREEN}http://localhost:80${NC}"
        echo ""
    fi

    echo -e "${BLUE}Default Admin Credentials:${NC}"
    echo -e "  Username: ${YELLOW}admin${NC}"
    echo -e "  Password: ${YELLOW}admin${NC}"
    echo -e "  ${RED}⚠  Change this immediately in production!${NC}"
    echo ""

    echo -e "${BLUE}Useful Commands:${NC}"
    echo -e "  View logs:        docker-compose logs -f"
    echo -e "  Stop services:    docker-compose down"
    echo -e "  Restart services: docker-compose restart"
    echo -e "  View status:      docker-compose ps"
    echo ""

    echo -e "${BLUE}Documentation:${NC}"
    echo -e "  Deployment Guide: ./deployment/README.md"
    echo -e "  Backend README:   ./backend/README.md"
    echo -e "  Frontend README:  ./frontend/README.md"
    echo ""
}

# Display service status
show_status() {
    print_header "Service Status"
    docker-compose ps
}

# Main deployment flow
main() {
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                                        ║${NC}"
    echo -e "${BLUE}║   ${GREEN}${PROJECT_NAME} Deployment Script${BLUE}          ║${NC}"
    echo -e "${BLUE}║                                        ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
    echo ""

    # Run deployment steps
    check_prerequisites
    setup_environment
    build_images
    start_services
    wait_for_services
    seed_database
    show_status
    display_access_info

    print_success "Deployment completed successfully! 🎉"
}

# Error handling
trap 'print_error "Deployment failed! Check the error messages above."; exit 1' ERR

# Run main function
main
