#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Load environment variables
ENV_FILE="${1:-.env.aws}"
if [ -f "$ENV_FILE" ]; then
    echo -e "${GREEN}Loading environment from $ENV_FILE${NC}"
    source "$ENV_FILE"
else
    log_error "Environment file $ENV_FILE not found"
    exit 1
fi

echo -e "${GREEN}=== Mountain Service - Backend Deployment ===${NC}"
echo -e "${BLUE}Environment: $ENV_FILE${NC}"
echo -e "${BLUE}Target: $INSTANCE_USER@$INSTANCE_IP${NC}"
echo

# Validate required variables
required_vars=("INSTANCE_IP" "INSTANCE_USER" "GHCR_PAT" "GITHUB_ACTOR")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        log_error "$var is not set"
        exit 1
    fi
done

# Test SSH connection
log_info "Testing SSH connection..."
if ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no -o ConnectTimeout=10 "$INSTANCE_USER@$INSTANCE_IP" "echo 'SSH connection successful'" > /dev/null 2>&1; then
    log_success "SSH connection successful"
else
    log_error "SSH connection failed"
    exit 1
fi

# Start backend deployment
log_info "Starting backend deployment..."

# Create deployment directory and clean up
ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no "$INSTANCE_USER@$INSTANCE_IP" << 'EOF'
    mkdir -p ~/mountain-service-deployment
    cd ~/mountain-service-deployment
    
    echo "Current directory: $(pwd)"
    echo "Files in directory:"
    ls -la
    
    echo "Running containers before deployment:"
    docker ps
EOF

# Copy files to remote
log_info "Copying deployment files..."
scp -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no docker-compose.aws.yml "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/docker-compose.yml"

# Use .env.backend if it exists (created by CI/CD), otherwise use the provided env file
if [ -f ".env.backend" ]; then
    log_info "Using .env.backend file created by CI/CD with actual image names"
    scp -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no .env.backend "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/.env"
else
    log_info "Using $ENV_FILE"
    scp -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no "$ENV_FILE" "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/.env"
fi

# Deploy backend services
ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no "$INSTANCE_USER@$INSTANCE_IP" << 'EOF'
    cd ~/mountain-service-deployment
    
    # Login to registry
    echo "$GHCR_PAT" | docker login ghcr.io -u "$GITHUB_ACTOR" --password-stdin

    # Stop only backend containers to preserve frontend
    echo "Stopping existing backend containers..."
    docker stop mountain-service-deployment_employee-service_1 2>/dev/null || true
    docker stop mountain-service-deployment_urgency-service_1 2>/dev/null || true
    docker stop mountain-service-deployment_activity-service_1 2>/dev/null || true
    docker stop mountain-service-deployment_version-service_1 2>/dev/null || true
    docker stop mountain-service-deployment_employee-db_1 2>/dev/null || true
    docker stop mountain-service-deployment_urgency-db_1 2>/dev/null || true
    docker stop mountain-service-deployment_activity-db_1 2>/dev/null || true

    # Remove stopped backend containers
    docker container prune -f || true

    # Stop existing backend services with compose (but preserve frontend)
    echo "Stopping existing backend services with compose..."
    docker-compose stop employee-service urgency-service activity-service version-service employee-db urgency-db activity-db 2>/dev/null || true
    docker-compose rm -f employee-service urgency-service activity-service version-service employee-db urgency-db activity-db 2>/dev/null || true
    
    # Clear Docker image cache to prevent stale images
    echo "Clearing Docker image cache to ensure fresh images..."
    docker system prune -f --volumes || true
    docker image prune -a -f || true

    # Pull backend service images with --no-cache equivalent
    echo "Pulling backend service images (forcing fresh pull)..."

    # Remove existing images to force fresh pull
    docker rmi $(docker images -q ghcr.io/$GITHUB_ACTOR/employee-service) 2>/dev/null || true
    docker rmi $(docker images -q ghcr.io/$GITHUB_ACTOR/urgency-service) 2>/dev/null || true
    docker rmi $(docker images -q ghcr.io/$GITHUB_ACTOR/activity-service) 2>/dev/null || true
    docker rmi $(docker images -q ghcr.io/$GITHUB_ACTOR/version-service) 2>/dev/null || true

    if ! docker-compose pull employee-db urgency-db activity-db employee-service urgency-service activity-service version-service; then
        echo "ERROR: Failed to pull backend Docker images from registry."
        exit 1
    fi

    # Verify we have the correct images by showing their digests
    echo "Verifying pulled images:"
    docker images --digests | grep -E "(employee-service|urgency-service|activity-service|version-service)" || true
    
    # Create network if it doesn't exist
    echo "Creating Docker network if it doesn't exist..."
    docker network create mountain-service-deployment_web 2>/dev/null || true

    # Verify network exists
    if docker network ls | grep -q mountain-service-deployment_web; then
        echo "Docker network mountain-service-deployment_web exists"
    else
        echo "Failed to create Docker network"
        exit 1
    fi

    # Deploy backend services
    echo "Deploying backend services..."
    docker-compose up -d --force-recreate employee-db urgency-db activity-db employee-service urgency-service activity-service version-service
    
    # Wait for services to be healthy
    echo "Waiting for backend services to be healthy..."
    sleep 30
    
    # Check service health
    echo "Checking backend service health..."
    for i in {1..12}; do
        echo "Health check attempt $i/12..."
        
        # Check databases
        if docker-compose ps employee-db | grep -q "Up (healthy)"; then
            echo "✓ Employee DB is healthy"
        else
            echo "⚠ Employee DB not healthy yet"
        fi
        
        if docker-compose ps urgency-db | grep -q "Up (healthy)"; then
            echo "✓ Urgency DB is healthy"
        else
            echo "⚠ Urgency DB not healthy yet"
        fi
        
        if docker-compose ps activity-db | grep -q "Up (healthy)"; then
            echo "✓ Activity DB is healthy"
        else
            echo "⚠ Activity DB not healthy yet"
        fi
        
        # Check services
        if curl -f http://localhost:8082/api/v1/health > /dev/null 2>&1; then
            echo "✓ Employee Service is healthy"
        else
            echo "⚠ Employee Service not healthy yet"
        fi
        
        if curl -f http://localhost:8083/api/v1/health > /dev/null 2>&1; then
            echo "✓ Urgency Service is healthy"
        else
            echo "⚠ Urgency Service not healthy yet"
        fi
        
        if curl -f http://localhost:8084/api/v1/health > /dev/null 2>&1; then
            echo "✓ Activity Service is healthy"
        else
            echo "⚠ Activity Service not healthy yet"
        fi
        
        if curl -f http://localhost:8090/api/v1/health > /dev/null 2>&1; then
            echo "✓ Version Service is healthy"
        else
            echo "⚠ Version Service not healthy yet"
        fi
        
        # Check if all services are healthy
        if curl -f http://localhost:8082/api/v1/health > /dev/null 2>&1 && \
           curl -f http://localhost:8083/api/v1/health > /dev/null 2>&1 && \
           curl -f http://localhost:8084/api/v1/health > /dev/null 2>&1 && \
           curl -f http://localhost:8090/api/v1/health > /dev/null 2>&1; then
            echo "🎉 All backend services are healthy!"
            break
        fi
        
        if [ $i -eq 12 ]; then
            echo "❌ Some backend services failed to become healthy after 6 minutes"
            echo "Container status:"
            docker-compose ps
            exit 1
        fi
        
        echo "Waiting 30 seconds before next health check..."
        sleep 30
    done
    
    echo "Container status after deployment:"
    docker-compose ps
    
    echo "All running containers:"
    docker ps
    
    echo "Backend deployment completed!"
EOF

if [ $? -eq 0 ]; then
    log_success "Backend deployment completed successfully!"
    log_info "Running final health checks..."

    # Wait a bit more for services to be fully ready
    log_info "Waiting 10 seconds for services to be fully ready..."
    sleep 10

    # Health checks from inside Docker network (since ports are not exposed externally)
    log_info "Testing services from inside Docker network..."

    if docker-compose exec -T employee-service curl -f -m 5 http://localhost:8082/api/v1/health > /dev/null 2>&1; then
        log_success "Employee Service accessible"
    else
        log_error "Employee Service not accessible"
    fi

    if docker-compose exec -T urgency-service curl -f -m 5 http://localhost:8083/api/v1/health > /dev/null 2>&1; then
        log_success "Urgency Service accessible"
    else
        log_error "Urgency Service not accessible"
    fi

    if docker-compose exec -T activity-service curl -f -m 5 http://localhost:8084/api/v1/health > /dev/null 2>&1; then
        log_success "Activity Service accessible"
    else
        log_error "Activity Service not accessible"
    fi

    if docker-compose exec -T version-service curl -f -m 5 http://localhost:8090/api/v1/health > /dev/null 2>&1; then
        log_success "Version Service accessible"
    else
        log_error "Version Service not accessible"
    fi
    
    echo
    log_success "Backend deployment completed successfully!"
    echo -e "${GREEN}Backend API URLs:${NC}"
    echo -e "${BLUE}  Employee API: http://$INSTANCE_IP:8082${NC}"
    echo -e "${BLUE}  Urgency API: http://$INSTANCE_IP:8083${NC}"
    echo -e "${BLUE}  Activity API: http://$INSTANCE_IP:8084${NC}"
    echo -e "${BLUE}  Version API: http://$INSTANCE_IP:8090${NC}"
    echo
    log_info "Note: If frontend was running, it may have been stopped. Run deploy-frontend.sh to restart it."
else
    log_error "Backend deployment failed!"
    exit 1
fi
