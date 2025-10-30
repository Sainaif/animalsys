# AnimalSys Backend

Backend ERP system for animal welfare foundations built with Go, Gin, and MongoDB.

## Architecture

The backend follows **Clean Architecture** principles with clear separation of concerns:

```
backend/
├── cmd/
│   ├── server/         # Main application entry point
│   └── seed/           # Database seeding tool
├── internal/
│   ├── core/
│   │   ├── entities/   # Business entities
│   │   ├── usecases/   # Business logic
│   │   └── interfaces/ # Repository interfaces
│   ├── adapters/
│   │   ├── repository/ # MongoDB implementations
│   │   ├── http/       # HTTP handlers and routing
│   │   └── auth/       # Authentication & RBAC
│   ├── infrastructure/ # Config, Database, Logger, Middleware
│   └── pkg/            # Shared utilities (JWT, Password)
```

## Prerequisites

- Go 1.21+
- MongoDB 6.0+
- Make (optional)

## Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Key configuration variables:
- `MONGODB_URI` - MongoDB connection string
- `JWT_SECRET` - Secret for JWT token signing
- `SERVER_PORT` - HTTP server port (default: 8080)
- `CORS_ORIGINS` - Allowed CORS origins

## Installation

```bash
# Install dependencies
go mod download

# Or using Make
make install
```

## Running the Application

### 1. Start MongoDB

```bash
# Using Docker
docker run -d -p 27017:27017 --name animalsys-mongo mongo:6.0
```

### 2. Seed the Database

Create initial super admin user:

```bash
go run cmd/seed/main.go

# Output:
# Email: admin@animalsys.local
# Username: admin
# Password: Admin123!
```

### 3. Start the Server

```bash
go run cmd/server/main.go

# Or using Make
make run

# Or build and run
make build
./bin/animalsys-server
```

The server will start on `http://localhost:8080`

## API Endpoints

### Public Endpoints

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/animals/available` - List animals available for adoption
- `GET /api/v1/campaigns/active` - List active campaigns

### Authenticated Endpoints

All authenticated endpoints require `Authorization: Bearer <token>` header.

#### Users (Admin only)
- `GET /api/v1/users` - List users
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### Animals (Employee+)
- `POST /api/v1/animals` - Create animal
- `GET /api/v1/animals` - List animals
- `PUT /api/v1/animals/:id` - Update animal
- `DELETE /api/v1/animals/:id` - Delete animal
- `POST /api/v1/animals/:id/medical-records` - Add medical record
- `POST /api/v1/animals/:id/photos` - Add photo

#### Other Modules
- `/api/v1/adoptions` - Adoption management
- `/api/v1/volunteers` - Volunteer management
- `/api/v1/schedules` - Shift scheduling
- `/api/v1/documents` - Document management
- `/api/v1/finance` - Financial transactions
- `/api/v1/donors` - Donor management
- `/api/v1/donations` - Donation tracking
- `/api/v1/inventory` - Inventory management
- `/api/v1/veterinary` - Veterinary records
- `/api/v1/campaigns` - Campaign management
- `/api/v1/partners` - Partner relationships
- `/api/v1/communications` - Communication management

## Role-Based Access Control (RBAC)

The system has 6 role levels with hierarchical permissions:

1. **Super Admin** - Full system access
2. **Admin** - Administrative operations
3. **Employee** - Day-to-day operations
4. **Volunteer** - Limited access (schedules, basic info)
5. **User** - Basic user operations (adoptions)
6. **Guest** - Public access only

Higher roles inherit permissions from lower roles.

## Development

```bash
# Run tests
make test

# Run with hot reload (requires air)
make dev

# Format code
make fmt

# Lint code
make lint

# Generate mocks
make mocks
```

## Project Features

### Security
- JWT-based authentication with access & refresh tokens
- Bcrypt password hashing (cost 12)
- RBAC with 6-level role hierarchy
- Rate limiting (100 req/15min per IP)
- Security headers (XSS, Clickjacking protection)
- Audit trail for all operations

### Middleware
- CORS with whitelist support
- Request logging with zerolog
- Panic recovery with stack traces
- Request validation and sanitization

### Database
- MongoDB with connection pooling
- Automatic index creation
- Aggregation pipelines for statistics
- GridFS for file storage

## Technology Stack

- **Framework**: Gin Web Framework
- **Database**: MongoDB with official Go driver
- **Authentication**: JWT tokens
- **Logging**: Zerolog
- **Configuration**: Viper
- **Validation**: Go validator
- **Password**: Bcrypt

## License

Proprietary - All rights reserved
