# AnimalSys Deployment Guide

This guide covers the deployment of the AnimalSys ERP system using Docker.

## Quick Start (One-Command Deployment)

```bash
# From the project root directory
./deployment/deploy.sh
```

This will:
1. Check prerequisites (Docker, Docker Compose)
2. Create `.env` from `.env.example` if not exists
3. Build all Docker images
4. Start all services
5. Run database migrations/seeding (if needed)
6. Display access URLs

## Prerequisites

- Docker 20.10+ installed
- Docker Compose 2.0+ installed
- At least 4GB of free RAM
- At least 10GB of free disk space

## Architecture

The deployment consists of 3 main services:

```
┌─────────────┐
│   Nginx     │ (Optional - Production only)
│  Port 80    │
└──────┬──────┘
       │
   ┌───┴────────────────┐
   │                    │
┌──▼────────┐    ┌─────▼─────┐
│ Frontend  │    │  Backend  │
│ Port 3000 │    │ Port 8080 │
└───────────┘    └─────┬─────┘
                       │
                 ┌─────▼─────┐
                 │  MongoDB  │
                 │Port 27017 │
                 └───────────┘
```

## Deployment Steps

### 1. Environment Configuration

Copy the example environment file and customize it:

```bash
cp .env.example .env
```

**Important:** Update these values in production:
- `MONGO_ROOT_PASSWORD` - MongoDB root password
- `JWT_SECRET` - Must be at least 32 characters
- `ENVIRONMENT` - Set to `production`
- `LOG_LEVEL` - Set to `warn` or `error` in production

### 2. Build and Start Services

```bash
# Development environment
docker-compose up -d

# Production environment (with Nginx)
docker-compose --profile production up -d
```

### 3. Verify Deployment

Check that all services are running:

```bash
docker-compose ps
```

All services should show status as "Up" and "healthy".

### 4. Access the Application

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **API Health Check:** http://localhost:8080/health
- **MongoDB:** localhost:27017 (from host machine)

Default admin credentials:
- Username: `admin`
- Password: `admin` (change immediately!)

## Service Management

### Start Services
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

### Restart a Service
```bash
docker-compose restart backend
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend

# Last 100 lines
docker-compose logs --tail=100 backend
```

### Update Services
```bash
# Pull latest images and restart
docker-compose pull
docker-compose up -d

# Rebuild from source
docker-compose build --no-cache
docker-compose up -d
```

## Database Operations

### Backup Database
```bash
./deployment/scripts/backup-db.sh
```

Backups are stored in `./backups/` directory with timestamp.

### Restore Database
```bash
./deployment/scripts/restore-db.sh /path/to/backup.tar.gz
```

### Seed Initial Data
```bash
docker-compose exec backend ./seed
```

## Monitoring

### Health Checks

All services have health checks configured:

```bash
# Check service health
docker-compose ps

# Check specific service health
docker inspect --format='{{.State.Health.Status}}' animalsys-backend
```

### Resource Usage

```bash
# View resource usage
docker stats

# View specific service
docker stats animalsys-backend
```

## Troubleshooting

### Service Won't Start

1. Check logs:
```bash
docker-compose logs backend
```

2. Verify environment variables:
```bash
docker-compose config
```

3. Check port conflicts:
```bash
netstat -tuln | grep -E '3000|8080|27017'
```

### Database Connection Issues

1. Verify MongoDB is healthy:
```bash
docker-compose exec mongodb mongosh --eval "db.adminCommand('ping')"
```

2. Check database logs:
```bash
docker-compose logs mongodb
```

### Frontend Can't Connect to Backend

1. Check CORS configuration in `.env`:
```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

2. Verify backend is accessible:
```bash
curl http://localhost:8080/health
```

### Out of Memory

Increase Docker memory allocation:
- Docker Desktop: Settings → Resources → Memory
- Linux: Update `/etc/docker/daemon.json`

## Security Best Practices

### Production Checklist

- [ ] Change all default passwords
- [ ] Use strong JWT_SECRET (32+ characters)
- [ ] Enable HTTPS (configure Nginx with SSL)
- [ ] Restrict MongoDB access (firewall rules)
- [ ] Set ENVIRONMENT=production
- [ ] Set LOG_LEVEL=warn or error
- [ ] Configure backup schedule
- [ ] Set up monitoring and alerts
- [ ] Review and restrict CORS origins
- [ ] Enable rate limiting
- [ ] Regular security updates

### SSL/TLS Configuration

For production with HTTPS:

1. Obtain SSL certificates (Let's Encrypt recommended)
2. Place certificates in `./deployment/nginx/ssl/`
3. Update nginx configuration
4. Deploy with production profile:

```bash
docker-compose --profile production up -d
```

## Performance Tuning

### MongoDB Optimization

```yaml
# Add to docker-compose.yml mongodb service
command: --wiredTigerCacheSizeGB 2 --maxConns 200
```

### Backend Scaling

```bash
# Run multiple backend instances
docker-compose up -d --scale backend=3
```

### Nginx Caching

Configure caching in `./deployment/nginx/nginx.conf` for static assets.

## Maintenance

### Log Rotation

Logs can grow large. Configure log rotation:

```bash
# View log sizes
docker-compose exec backend du -sh /app/logs/*

# Clear old logs (careful!)
docker-compose exec backend find /app/logs -name "*.log" -mtime +30 -delete
```

### Database Maintenance

```bash
# Compact database
docker-compose exec mongodb mongosh animalsys --eval "db.runCommand({compact: 'collection_name'})"

# Rebuild indexes
docker-compose exec mongodb mongosh animalsys --eval "db.collection.reIndex()"
```

## Upgrading

### Minor Updates

```bash
git pull origin main
docker-compose build
docker-compose up -d
```

### Major Updates

1. Backup database
2. Read CHANGELOG.md for breaking changes
3. Update environment variables if needed
4. Run migrations
5. Test in staging environment first
6. Deploy to production

## Support

For issues:
- Check logs: `docker-compose logs`
- Review this README
- Check GitHub Issues: [repository-url]
- Contact: support@animalsys.org

## License

Copyright © 2025 AnimalSys. All rights reserved.
