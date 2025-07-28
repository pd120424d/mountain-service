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

ENV_FILE="${1:-.env.aws}"

echo -e "${GREEN}=== Mountain Service - Full Deployment ===${NC}"
echo -e "${BLUE}Environment: $ENV_FILE${NC}"
echo

# Deploy backend first
log_info "Starting backend deployment..."
if ./deploy-backend.sh "$ENV_FILE"; then
    log_success "Backend deployment completed successfully!"
else
    log_error "Backend deployment failed!"
    exit 1
fi

echo
echo "=================================================="
echo

# Deploy frontend after backend is ready
log_info "Starting frontend deployment..."
if ./deploy-frontend.sh "$ENV_FILE"; then
    log_success "Frontend deployment completed successfully!"
else
    log_error "Frontend deployment failed!"
    exit 1
fi

echo
log_success "Full deployment completed successfully!"
echo -e "${GREEN}All services are now running:${NC}"

# Load environment to get INSTANCE_IP
source "$ENV_FILE"
echo -e "${BLUE}  Frontend: http://$INSTANCE_IP${NC}"
echo -e "${BLUE}  Employee API: http://$INSTANCE_IP:8082${NC}"
echo -e "${BLUE}  Urgency API: http://$INSTANCE_IP:8083${NC}"
echo -e "${BLUE}  Activity API: http://$INSTANCE_IP:8084${NC}"
echo -e "${BLUE}  Version API: http://$INSTANCE_IP:8090${NC}"
