# Deployment Guide

This document describes the deployment configuration for the Mountain Service.

## Configuration Files

### Docker Compose
- **`docker-compose.yml`** - Single unified compose file for all environments
- **`.env.staging`** - Staging environment configuration
- **`.env.aws`** - AWS production environment configuration

### Deployment Scripts
- **`deploy-simple.sh`** - Unified deployment script for all environments
- **`backend-test-cover.sh`** - Backend test coverage script
- **`frontend-test-cover.sh`** - Frontend test coverage script

## Usage

### Local Development (Staging)
```bash
# Start all services
docker-compose --env-file .env.staging up -d

# Start with Swagger UI
COMPOSE_PROFILES=swagger docker-compose --env-file .env.staging up -d

# Stop services
docker-compose down
```

### AWS Production Deployment
```bash
# Full deployment (backend + frontend)
./deploy-simple.sh full .env.aws

# Backend only
./deploy-simple.sh backend .env.aws

# Frontend only
./deploy-simple.sh frontend .env.aws
```

### Testing
```bash
# Backend tests
./backend-test-cover.sh

# Frontend tests
./frontend-test-cover.sh

# All tests (via Makefile)
make coverage
```

## Environment Variables

### Staging (.env.staging)
- Uses local build context
- Exposes services on high ports (10001-10003)
- Includes Swagger UI
- Default credentials for development

### AWS (.env.aws)
- Uses pre-built Docker images from registry
- Standard ports (5432, 5433, 5434)
- Production security settings
- Environment variables from CI/CD secrets

## CI/CD Integration

GitHub Actions workflows use the following commands:
- Backend deployment: `./deploy-simple.sh backend .env.aws`
- Frontend deployment: `./deploy-simple.sh frontend .env.aws`
- Path triggers updated to include new configuration files

## Quick Reference

| Task | Command |
|------|---------|
| Local development | `docker-compose --env-file .env.staging up -d` |
| AWS full deploy | `./deploy-simple.sh full .env.aws` |
| AWS backend only | `./deploy-simple.sh backend .env.aws` |
| AWS frontend only | `./deploy-simple.sh frontend .env.aws` |
| Run tests | `make coverage` |
| Stop services | `docker-compose down` |
| View logs | `docker-compose logs -f` |
| Help | `./deploy-simple.sh --help` |
