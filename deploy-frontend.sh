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

echo -e "${GREEN}=== Mountain Service - Frontend Deployment ===${NC}"
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

# Check if backend services are running and healthy (from inside the VM)
log_info "Checking backend services health before frontend deployment..."

# Check backend health from inside the VM
backend_health_check=$(ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no "$INSTANCE_USER@$INSTANCE_IP" << 'EOF'
backend_healthy=true

echo "Checking backend services from inside VM..."

if ! curl -f http://localhost:8082/api/v1/health > /dev/null 2>&1; then
    echo "ERROR: Employee Service (port 8082) is not healthy"
    backend_healthy=false
fi

if ! curl -f http://localhost:8083/api/v1/health > /dev/null 2>&1; then
    echo "ERROR: Urgency Service (port 8083) is not healthy"
    backend_healthy=false
fi

if ! curl -f http://localhost:8084/api/v1/health > /dev/null 2>&1; then
    echo "ERROR: Activity Service (port 8084) is not healthy"
    backend_healthy=false
fi

if ! curl -f http://localhost:8090/api/v1/health > /dev/null 2>&1; then
    echo "ERROR: Version Service (port 8090) is not healthy"
    backend_healthy=false
fi

if [ "$backend_healthy" = "true" ]; then
    echo "SUCCESS: All backend services are healthy"
    exit 0
else
    echo "ERROR: Some backend services are not healthy"
    exit 1
fi
EOF
)

if [ $? -eq 0 ]; then
    log_success "All backend services are healthy. Proceeding with frontend deployment..."
else
    log_error "Backend services are not healthy. Please deploy backend first using deploy-backend.sh"
    echo "$backend_health_check"
    exit 1
fi

# Start frontend deployment
log_info "Starting frontend deployment..."

# Copy files to remote
log_info "Copying deployment files..."
scp -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no docker-compose.frontend.yml "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/docker-compose-frontend.yml"

# Use .env.frontend if it exists (created by CI/CD), otherwise use the provided env file
if [ -f ".env.frontend" ]; then
    log_info "Using .env.frontend file created by CI/CD with actual image names"
    scp -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no .env.frontend "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/.env.frontend"
else
    log_info "Using $ENV_FILE for frontend"
    scp -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no "$ENV_FILE" "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/.env.frontend"
fi

# Deploy frontend service
ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no "$INSTANCE_USER@$INSTANCE_IP" << 'EOF'
    cd ~/mountain-service-deployment
    
    echo "Current directory: $(pwd)"
    echo "Files in directory:"
    ls -la
    
    # Login to registry
    echo "$GHCR_PAT" | docker login ghcr.io -u "$GITHUB_ACTOR" --password-stdin
    
    # Stop any existing frontend container (may have been stopped by backend deployment)
    echo "Stopping existing frontend container..."
    docker stop mountain-service-deployment_frontend_1 2>/dev/null || true
    docker rm mountain-service-deployment_frontend_1 2>/dev/null || true

    # Also check for frontend containers from the frontend-only compose file
    docker-compose -f docker-compose-frontend.yml --env-file .env.frontend stop 2>/dev/null || true
    docker-compose -f docker-compose-frontend.yml --env-file .env.frontend rm -f 2>/dev/null || true
    
    # Pull frontend image
    echo "Pulling frontend image..."
    if ! docker-compose -f docker-compose-frontend.yml --env-file .env.frontend pull; then
        echo "ERROR: Failed to pull frontend Docker image from registry."
        exit 1
    fi
    
    # Ensure network exists and connect to it
    echo "Ensuring Docker network exists..."
    docker network create mountain-service-deployment_web 2>/dev/null || true

    # Deploy frontend service
    echo "Deploying frontend service..."
    docker-compose -f docker-compose-frontend.yml --env-file .env.frontend up -d --force-recreate
    
    # Wait for frontend to start
    echo "Waiting for frontend to start..."
    sleep 15
    
    # Check frontend health
    echo "Checking frontend health..."
    for i in {1..6}; do
        echo "Frontend health check attempt $i/6..."
        
        if curl -f http://localhost/health > /dev/null 2>&1; then
            echo "✓ Frontend health endpoint is responding"
            break
        elif curl -f http://localhost/ > /dev/null 2>&1; then
            echo "✓ Frontend is serving content"
            break
        else
            echo "⚠ Frontend not responding yet"
        fi
        
        if [ $i -eq 6 ]; then
            echo "❌ Frontend failed to become healthy after 3 minutes"
            echo "Frontend container status:"
            docker-compose -f docker-compose-frontend.yml ps
            echo "Frontend container logs:"
            docker-compose -f docker-compose-frontend.yml logs frontend
            exit 1
        fi
        
        echo "Waiting 30 seconds before next health check..."
        sleep 30
    done
    
    echo "Frontend container status:"
    docker-compose -f docker-compose-frontend.yml ps
    
    echo "All running containers:"
    docker ps
    
    echo "Frontend deployment completed!"
EOF

if [ $? -eq 0 ]; then
    log_success "Frontend deployment completed successfully!"
    log_info "Running final health checks..."

    # Final health check from inside the VM (since external ports may not be open)
    sleep 5
    final_health_check=$(ssh -i "$SSH_KEY_PATH" -o StrictHostKeyChecking=no "$INSTANCE_USER@$INSTANCE_IP" << 'EOF'
if curl -f http://localhost/health > /dev/null 2>&1; then
    echo "SUCCESS: Frontend health endpoint accessible"
    exit 0
elif curl -f http://localhost/ > /dev/null 2>&1; then
    echo "SUCCESS: Frontend application accessible"
    exit 0
else
    echo "WARNING: Frontend not accessible on localhost"
    exit 1
fi
EOF
)

    if [ $? -eq 0 ]; then
        log_success "Frontend is running and accessible"
    else
        log_warning "Frontend may not be fully ready yet"
        echo "$final_health_check"
    fi
    
    echo
    log_success "Frontend deployment completed successfully!"
    echo -e "${GREEN}Application URLs:${NC}"
    echo -e "${BLUE}  Frontend: http://$INSTANCE_IP${NC}"
    echo -e "${BLUE}  Employee API: http://$INSTANCE_IP:8082${NC}"
    echo -e "${BLUE}  Urgency API: http://$INSTANCE_IP:8083${NC}"
    echo -e "${BLUE}  Activity API: http://$INSTANCE_IP:8084${NC}"
    echo -e "${BLUE}  Version API: http://$INSTANCE_IP:8090${NC}"
else
    log_error "Frontend deployment failed!"
    exit 1
fi
