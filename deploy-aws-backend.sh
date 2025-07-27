#!/bin/bash

# AWS Backend Deployment Script
# This script deploys only the backend services to AWS

set -e

# Load environment variables
if [ -f .env.backend ]; then
    source .env.backend
else
    echo "Error: .env.backend file not found"
    exit 1
fi

# Validate required environment variables
required_vars=("INSTANCE_IP" "INSTANCE_USER" "EMPLOYEE_SERVICE_IMAGE" "URGENCY_SERVICE_IMAGE" "ACTIVITY_SERVICE_IMAGE" "VERSION_SERVICE_IMAGE" "GHCR_PAT" "GITHUB_ACTOR" "JWT_SECRET" "ADMIN_PASSWORD" "DB_USER" "DB_PASSWORD")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: Required environment variable $var is not set"
        exit 1
    fi
done

echo "🚀 Starting AWS Backend Deployment..."
echo "Target: $INSTANCE_USER@$INSTANCE_IP"
echo "Employee Service: $EMPLOYEE_SERVICE_IMAGE"
echo "Urgency Service: $URGENCY_SERVICE_IMAGE"
echo "Activity Service: $ACTIVITY_SERVICE_IMAGE"
echo "Version Service: $VERSION_SERVICE_IMAGE"

# Create deployment directory on remote server
ssh -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no $INSTANCE_USER@$INSTANCE_IP << 'EOF'
    mkdir -p ~/mountain-service-backend
    cd ~/mountain-service-backend
EOF

# Copy backend-specific docker-compose file to remote server
cat > docker-compose.backend.yml << 'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: mountain-rescue-postgres
    environment:
      POSTGRES_DB: mountain_rescue
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
    ports:
      - "5432:5432"
    restart: unless-stopped
    networks:
      - mountain-rescue-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d mountain_rescue"]
      interval: 30s
      timeout: 10s
      retries: 5

  employee-service:
    image: ${EMPLOYEE_SERVICE_IMAGE}
    container_name: mountain-rescue-employee-service
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: mountain_rescue
      JWT_SECRET: ${JWT_SECRET}
      ADMIN_PASSWORD: ${ADMIN_PASSWORD}
      SERVICE_AUTH_SECRET: ${SERVICE_AUTH_SECRET}
      CORS_ALLOWED_ORIGINS: ${CORS_ALLOWED_ORIGINS}
      GIN_MODE: release
    ports:
      - "8082:8082"
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - mountain-rescue-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8082/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  urgency-service:
    image: ${URGENCY_SERVICE_IMAGE}
    container_name: mountain-rescue-urgency-service
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: mountain_rescue
      JWT_SECRET: ${JWT_SECRET}
      SERVICE_AUTH_SECRET: ${SERVICE_AUTH_SECRET}
      EMPLOYEE_SERVICE_URL: ${EMPLOYEE_SERVICE_URL}
      GIN_MODE: release
    ports:
      - "8083:8083"
    depends_on:
      postgres:
        condition: service_healthy
      employee-service:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - mountain-rescue-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8083/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  activity-service:
    image: ${ACTIVITY_SERVICE_IMAGE}
    container_name: mountain-rescue-activity-service
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: mountain_rescue
      JWT_SECRET: ${JWT_SECRET}
      SERVICE_AUTH_SECRET: ${SERVICE_AUTH_SECRET}
      EMPLOYEE_SERVICE_URL: ${EMPLOYEE_SERVICE_URL}
      GIN_MODE: release
    ports:
      - "8084:8084"
    depends_on:
      postgres:
        condition: service_healthy
      employee-service:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - mountain-rescue-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8084/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  version-service:
    image: ${VERSION_SERVICE_IMAGE}
    container_name: mountain-rescue-version-service
    environment:
      GIN_MODE: release
    ports:
      - "8090:8090"
    restart: unless-stopped
    networks:
      - mountain-rescue-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8090/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

volumes:
  postgres_data:

networks:
  mountain-rescue-network:
    external: true
EOF

# Copy the docker-compose file to remote server
scp -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no docker-compose.backend.yml $INSTANCE_USER@$INSTANCE_IP:~/mountain-service-backend/

# Create backend environment file for remote server
cat > .env.backend.remote << EOF
EMPLOYEE_SERVICE_IMAGE=$EMPLOYEE_SERVICE_IMAGE
URGENCY_SERVICE_IMAGE=$URGENCY_SERVICE_IMAGE
ACTIVITY_SERVICE_IMAGE=$ACTIVITY_SERVICE_IMAGE
VERSION_SERVICE_IMAGE=$VERSION_SERVICE_IMAGE
GHCR_PAT=$GHCR_PAT
GITHUB_ACTOR=$GITHUB_ACTOR
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
JWT_SECRET=$JWT_SECRET
ADMIN_PASSWORD=$ADMIN_PASSWORD
SERVICE_AUTH_SECRET=$SERVICE_AUTH_SECRET
CORS_ALLOWED_ORIGINS=$CORS_ALLOWED_ORIGINS
EMPLOYEE_SERVICE_URL=$EMPLOYEE_SERVICE_URL
ACTIVITY_SERVICE_URL=$ACTIVITY_SERVICE_URL
EOF

# Copy environment file to remote server
scp -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no .env.backend.remote $INSTANCE_USER@$INSTANCE_IP:~/mountain-service-backend/.env

# Copy database initialization script
scp -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no init-db.sql $INSTANCE_USER@$INSTANCE_IP:~/mountain-service-backend/

# Deploy backend on remote server
ssh -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no $INSTANCE_USER@$INSTANCE_IP << 'EOF'
    cd ~/mountain-service-backend

    echo "🔐 Logging into GitHub Container Registry..."
    echo $GHCR_PAT | docker login ghcr.io -u $GITHUB_ACTOR --password-stdin

    echo "🌐 Creating Docker network if it doesn't exist..."
    docker network create mountain-rescue-network || true

    echo "🛑 Stopping existing backend services..."
    docker-compose -f docker-compose.backend.yml down || true

    echo "🧹 Cleaning up old backend images..."
    docker image prune -f

    echo "📥 Pulling latest backend images..."
    docker-compose -f docker-compose.backend.yml pull

    echo "🚀 Starting backend services..."
    docker-compose -f docker-compose.backend.yml up -d

    echo "⏳ Waiting for backend services to be ready..."
    sleep 45

    echo "🔍 Checking backend services status..."
    docker-compose -f docker-compose.backend.yml ps

    echo "📋 Backend services logs (last 10 lines each):"
    echo "=== Employee Service ==="
    docker-compose -f docker-compose.backend.yml logs --tail=10 employee-service
    echo "=== Urgency Service ==="
    docker-compose -f docker-compose.backend.yml logs --tail=10 urgency-service
    echo "=== Activity Service ==="
    docker-compose -f docker-compose.backend.yml logs --tail=10 activity-service
    echo "=== Version Service ==="
    docker-compose -f docker-compose.backend.yml logs --tail=10 version-service
EOF

# Verify deployment
echo "🔍 Verifying backend deployment..."
sleep 15

# Test backend services
echo "Testing employee service..."
if curl -f http://$INSTANCE_IP:8082/health; then
    echo "✅ Employee service health check passed"
else
    echo "❌ Employee service health check failed"
fi

echo "Testing version service..."
if curl -f http://$INSTANCE_IP:8090/health; then
    echo "✅ Version service health check passed"
else
    echo "❌ Version service health check failed"
fi

# Cleanup local files
rm -f docker-compose.backend.yml .env.backend.remote

echo "🎉 Backend deployment completed successfully!"
echo "🔗 Employee API: http://$INSTANCE_IP:8082"
echo "🔗 Urgency API: http://$INSTANCE_IP:8083"
echo "🔗 Activity API: http://$INSTANCE_IP:8084"
echo "🔗 Version API: http://$INSTANCE_IP:8090"
