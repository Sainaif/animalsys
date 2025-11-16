# GEMINI.md

## Project Overview

This is a full-stack web application for an Animal Foundation CRM (Customer Relationship Management). It's built with a Go backend and a Vue.js frontend. The entire application is containerized using Docker.

**Backend:**
- **Language:** Go
- **Framework:** Gin
- **Database:** MongoDB
- **Authentication:** JWT

**Frontend:**
- **Framework:** Vue.js 3 (Composition API)
- **Build Tool:** Vite
- **UI Library:** PrimeVue
- **State Management:** Pinia

**Infrastructure:**
- **Containerization:** Docker + Docker Compose
- **Web Server:** Nginx

## Building and Running

The project uses a `Makefile` for common commands.

**To initialize the project (run this first):**
```bash
make init
```

**To start the application:**
```bash
make up
```
or for development with live-reloading:
```bash
make dev
```

**To stop the application:**
```bash
make down
```

**To view logs:**
```bash
make logs
```

**To run tests:**
```bash
# Backend
make test-backend

# Frontend
make test-frontend
```

## Development Conventions

### Backend
- The backend follows a Clean Architecture structure, separating domain logic from infrastructure concerns.
- Code is located in the `backend/internal` directory.
- API endpoints are defined in `backend/internal/delivery/http/routes`.
- Business logic is in `backend/internal/usecase`.
- Database interactions are handled in `backend/internal/infrastructure/database`.

### Frontend
- The frontend uses the Composition API.
- Components are located in `frontend/src/components`.
- Views (pages) are in `frontend/src/views`.
- State is managed with Pinia stores in `frontend/src/stores`.
- API calls are made from services in `frontend/src/services`.
