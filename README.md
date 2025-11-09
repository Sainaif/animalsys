# Animal Foundation CRM

A comprehensive CRM system for animal foundations built with Clean Architecture principles.

## 🎯 Features

### MVP (Current Phase)
- ✅ Clean Architecture project structure
- ✅ Docker Compose development environment
- 🚧 JWT Authentication & Authorization with RBAC
- 🚧 Foundation Settings Management
- 🚧 Comprehensive Animal Management
- 🚧 Veterinary Care Tracking
- 🚧 Adoption Management
- 🚧 Dashboard & Basic Reports

### Post-MVP (Planned)
- Donor Management with Stripe Integration
- Communication Management (Email/SMS)
- Task Management
- Document Management
- Notification System
- Volunteer Management
- Partner Management
- Campaign Management
- Advanced Reporting & Analytics

## 🛠️ Technology Stack

### Backend
- **Language:** Go 1.22
- **Framework:** Gin
- **Database:** MongoDB 7.0
- **Cache:** Redis 7
- **Auth:** JWT (15min access + 7 day refresh tokens)
- **Config:** Viper
- **Logging:** Zerolog
- **Validation:** go-playground/validator

### Frontend
- **Framework:** Vue.js 3 (Composition API)
- **Build Tool:** Vite
- **UI Library:** PrimeVue
- **State:** Pinia
- **Router:** Vue Router
- **HTTP:** Axios
- **i18n:** vue-i18n (Polish + English)
- **Charts:** Chart.js

### Infrastructure
- **Containerization:** Docker + Docker Compose
- **Web Server:** Nginx
- **Deployment:** VPS-ready

## 📋 Prerequisites

- Docker (v20.10+)
- Docker Compose (v2.0+)
- Make (optional, for convenience)
- Git

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd animalsys
```

### 2. Initialize Project

```bash
make init
```

This will:
- Copy `.env.example` to `.env`
- Build all Docker images
- Start all services

### 3. Configure Environment

Edit `.env` file with your settings:

```bash
# Required changes for production
JWT_SECRET=your-super-secret-key-here

# Optional: Third-party services
EMAIL_API_KEY=your-sendgrid-api-key
SMS_ACCOUNT_SID=your-twilio-sid
SMS_AUTH_TOKEN=your-twilio-token
PAYMENT_SECRET_KEY=your-stripe-secret-key
```

### 4. Access the Application

- **Frontend:** http://localhost
- **Backend API:** http://localhost/api/v1
- **Health Check:** http://localhost/health

## 📚 Development

### Starting Services

```bash
# Start all services
make up

# Start with logs
make dev
```

### Viewing Logs

```bash
# All services
make logs

# Backend only
make logs-backend

# Frontend only
make logs-frontend
```

### Accessing Containers

```bash
# Backend shell
make backend-shell

# Frontend shell
make frontend-shell

# MongoDB shell
make db-shell

# Redis CLI
make redis-cli
```

### Stopping Services

```bash
# Stop all services
make down

# Stop and remove volumes (clean reset)
make clean
```

## 🧪 Testing

```bash
# Run backend tests
make test-backend

# Run frontend tests
make test-frontend
```

## 📖 Project Structure

```
animalsys/
├── backend/                    # Go backend
│   ├── cmd/server/            # Application entry point
│   ├── internal/
│   │   ├── domain/            # Business entities & interfaces
│   │   ├── usecase/           # Business logic
│   │   ├── infrastructure/    # Database, config, logging
│   │   └── delivery/http/     # HTTP handlers
│   └── pkg/                   # Shared packages
│
├── frontend/                   # Vue.js frontend
│   └── src/
│       ├── components/        # Vue components
│       ├── views/             # Page views
│       ├── stores/            # Pinia stores
│       ├── services/          # API services
│       └── router/            # Vue Router
│
├── nginx/                      # Nginx configuration
├── uploads/                    # File uploads (Docker volume)
├── docker-compose.yml          # Development environment
├── docker-compose.prod.yml     # Production environment
├── Makefile                    # Development commands
├── plan.md                     # Implementation plan
└── progress.md                 # Development progress
```

## 🗄️ Database Management

### Backup Database

```bash
make backup-db
```

This creates a backup file: `backup_YYYYMMDD_HHMMSS.archive`

### Restore Database

```bash
make restore-db FILE=backup_20250108_120000.archive
```

## 🔧 Configuration

### Environment Variables

Key environment variables in `.env`:

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Environment mode | `development` |
| `SERVER_PORT` | Backend port | `8080` |
| `DB_URI` | MongoDB connection string | `mongodb://mongodb:27017` |
| `DB_NAME` | Database name | `animalsys` |
| `REDIS_HOST` | Redis host | `redis` |
| `JWT_SECRET` | JWT signing key | (must change in production) |
| `STORAGE_TYPE` | Storage type (`local` or `s3`) | `local` |

### Third-Party Services

#### Email (SendGrid/Mailgun)
```env
EMAIL_PROVIDER=sendgrid
EMAIL_API_KEY=your-api-key
EMAIL_FROM=noreply@yourfoundation.org
```

#### SMS (Twilio)
```env
SMS_PROVIDER=twilio
SMS_ACCOUNT_SID=your-account-sid
SMS_AUTH_TOKEN=your-auth-token
SMS_PHONE_NUMBER=+1234567890
```

#### Payment (Stripe)
```env
PAYMENT_PROVIDER=stripe
PAYMENT_SECRET_KEY=sk_test_...
PAYMENT_PUBLISHABLE_KEY=pk_test_...
```

## 🚢 Production Deployment

### 1. Build Production Images

```bash
make prod-build
```

### 2. Update Production Environment

Create `.env.production` with production settings:
- Strong `JWT_SECRET`
- Production database credentials
- Third-party API keys
- CORS settings

### 3. Deploy

```bash
make prod-up
```

### 4. Configure SSL (with Let's Encrypt)

```bash
# Install certbot on VPS
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d yourdomain.com
```

### 5. Set up Backups

Add to crontab:
```bash
# Daily database backup at 2 AM
0 2 * * * cd /path/to/animalsys && make backup-db
```

## 📝 Development Workflow

### Phase 0: Infrastructure ✅
- [x] Project structure
- [x] Docker Compose setup
- [x] Database connection
- [x] Logging & config

### Phase 1: Authentication 🚧
- [ ] User management
- [ ] JWT authentication
- [ ] RBAC implementation
- [ ] Audit logging

### Phase 2: Foundation Settings 🔜
- [ ] Basic info management
- [ ] Animal type configuration
- [ ] Multi-language support

### Phase 3: Animal Management 🔜
- [ ] CRUD operations
- [ ] Photo management
- [ ] Advanced filtering
- [ ] Timeline tracking

See [plan.md](plan.md) for complete roadmap and [progress.md](progress.md) for current status.

## 🤝 Contributing

1. Create a feature branch from `develop`
2. Make your changes
3. Write/update tests
4. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

## 🆘 Troubleshooting

### Services Won't Start

```bash
# Check logs
make logs

# Clean and restart
make clean
make up
```

### Database Connection Issues

```bash
# Check MongoDB is running
docker-compose ps mongodb

# Check MongoDB logs
docker-compose logs mongodb

# Try connecting manually
make db-shell
```

### Port Already in Use

Edit `docker-compose.yml` and change port mappings:
```yaml
ports:
  - "8081:8080"  # Change 8080 to 8081
```

## 📞 Support

For issues and questions:
- Create an issue in the repository
- Check [plan.md](plan.md) for implementation details
- Review [progress.md](progress.md) for current status

---

**Version:** 0.1.0
**Last Updated:** 2025-01-08
**Status:** Phase 0 Complete, Phase 1 In Progress
