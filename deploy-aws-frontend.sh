#!/bin/bash

# AWS Frontend Deployment Script
# This script deploys only the frontend container to AWS

set -e

# Load environment variables
if [ -f .env.frontend ]; then
    source .env.frontend
else
    echo "Error: .env.frontend file not found"
    exit 1
fi

# Validate required environment variables
required_vars=("INSTANCE_IP" "INSTANCE_USER" "FRONTEND_IMAGE" "GHCR_PAT" "GITHUB_ACTOR")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: Required environment variable $var is not set"
        exit 1
    fi
done

echo "🚀 Starting AWS Frontend Deployment..."
echo "Target: $INSTANCE_USER@$INSTANCE_IP"
echo "Frontend Image: $FRONTEND_IMAGE"

# Create deployment directory on remote server
ssh -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no $INSTANCE_USER@$INSTANCE_IP << 'EOF'
    mkdir -p ~/mountain-service-frontend
    cd ~/mountain-service-frontend
EOF

# Copy frontend-specific docker-compose file to remote server
cat > docker-compose.frontend.yml << 'EOF'
version: '3.8'

services:
  frontend:
    image: ${FRONTEND_IMAGE}
    container_name: mountain-rescue-frontend
    ports:
      - "80:80"
      - "443:443"
    restart: unless-stopped
    networks:
      - mountain-rescue-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  mountain-rescue-network:
    external: true
EOF

# Copy the docker-compose file to remote server
scp -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no docker-compose.frontend.yml $INSTANCE_USER@$INSTANCE_IP:~/mountain-service-frontend/

# Create frontend environment file for remote server
cat > .env.frontend.remote << EOF
FRONTEND_IMAGE=$FRONTEND_IMAGE
GHCR_PAT=$GHCR_PAT
GITHUB_ACTOR=$GITHUB_ACTOR
EOF

# Copy environment file to remote server
scp -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no .env.frontend.remote $INSTANCE_USER@$INSTANCE_IP:~/mountain-service-frontend/.env

# Deploy frontend on remote server
ssh -i ~/.ssh/deploy_key -o StrictHostKeyChecking=no $INSTANCE_USER@$INSTANCE_IP << 'EOF'
    cd ~/mountain-service-frontend
    
    echo "🔐 Logging into GitHub Container Registry..."
    echo $GHCR_PAT | docker login ghcr.io -u $GITHUB_ACTOR --password-stdin
    
    echo "🛑 Stopping existing frontend container..."
    docker-compose -f docker-compose.frontend.yml down || true
    
    echo "🧹 Cleaning up old frontend images..."
    docker image prune -f
    
    echo "📥 Pulling latest frontend image..."
    docker-compose -f docker-compose.frontend.yml pull
    
    echo "🚀 Starting frontend container..."
    docker-compose -f docker-compose.frontend.yml up -d
    
    echo "⏳ Waiting for frontend to be ready..."
    sleep 15
    
    echo "🔍 Checking frontend container status..."
    docker-compose -f docker-compose.frontend.yml ps
    
    echo "📋 Frontend container logs (last 20 lines):"
    docker-compose -f docker-compose.frontend.yml logs --tail=20 frontend
EOF

# Verify deployment
echo "🔍 Verifying frontend deployment..."
sleep 10

# Test frontend health endpoint
if curl -f http://$INSTANCE_IP/health; then
    echo "✅ Frontend health check passed"
else
    echo "❌ Frontend health check failed"
    exit 1
fi

# Test frontend application
if curl -f http://$INSTANCE_IP/ > /dev/null; then
    echo "✅ Frontend application is accessible"
else
    echo "❌ Frontend application is not accessible"
    exit 1
fi

# Cleanup local files
rm -f docker-compose.frontend.yml .env.frontend.remote

echo "🎉 Frontend deployment completed successfully!"
echo "🌐 Application URL: http://$INSTANCE_IP"
