#!/bin/bash

# Quick AWS Deployment Script for Mountain Rescue Service
# This script provides a faster deployment option by only rebuilding changed services

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"

# Load environment variables
if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
else
    echo -e "${RED}Error: .env file not found. Please copy .env.aws to .env and configure it.${NC}"
    exit 1
fi

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Parse command line arguments
SERVICES_TO_DEPLOY=()
FORCE_REBUILD=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --frontend|--ui)
            SERVICES_TO_DEPLOY+=("frontend")
            shift
            ;;
        --employee)
            SERVICES_TO_DEPLOY+=("employee")
            shift
            ;;
        --urgency)
            SERVICES_TO_DEPLOY+=("urgency")
            shift
            ;;
        --activity)
            SERVICES_TO_DEPLOY+=("activity")
            shift
            ;;
        --version)
            SERVICES_TO_DEPLOY+=("version")
            shift
            ;;
        --all)
            SERVICES_TO_DEPLOY=("frontend" "employee" "urgency" "activity" "version")
            shift
            ;;
        --force)
            FORCE_REBUILD=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS] [SERVICES]"
            echo ""
            echo "Options:"
            echo "  --frontend, --ui    Deploy only frontend service"
            echo "  --employee          Deploy only employee service"
            echo "  --urgency           Deploy only urgency service"
            echo "  --activity          Deploy only activity service"
            echo "  --version           Deploy only version service"
            echo "  --all               Deploy all services (default)"
            echo "  --force             Force rebuild even if no changes detected"
            echo "  --help, -h          Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 --frontend       # Deploy only frontend"
            echo "  $0 --employee --urgency  # Deploy employee and urgency services"
            echo "  $0 --all --force    # Force rebuild and deploy all services"
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Default to all services if none specified
if [ ${#SERVICES_TO_DEPLOY[@]} -eq 0 ]; then
    SERVICES_TO_DEPLOY=("frontend" "employee" "urgency" "activity" "version")
fi

# Build and deploy specific services
deploy_services() {
    log_info "Quick deploying services: ${SERVICES_TO_DEPLOY[*]}"
    
    # Build only specified services
    for service in "${SERVICES_TO_DEPLOY[@]}"; do
        case $service in
            "frontend")
                log_info "Building frontend service..."
                cd "$SCRIPT_DIR/ui"
                docker build -f Dockerfile.aws -t "${FRONTEND_IMAGE}" .
                docker save "${FRONTEND_IMAGE}" -o /tmp/frontend.tar
                cd "$SCRIPT_DIR"
                ;;
            "employee")
                if [ -d "api/employee" ]; then
                    log_info "Building employee service..."
                    cd "api/employee"
                    docker build -t "${EMPLOYEE_SERVICE_IMAGE}" .
                    docker save "${EMPLOYEE_SERVICE_IMAGE}" -o /tmp/employee.tar
                    cd "$SCRIPT_DIR"
                fi
                ;;
            "urgency")
                if [ -d "api/urgency" ]; then
                    log_info "Building urgency service..."
                    cd "api/urgency"
                    docker build -t "${URGENCY_SERVICE_IMAGE}" .
                    docker save "${URGENCY_SERVICE_IMAGE}" -o /tmp/urgency.tar
                    cd "$SCRIPT_DIR"
                fi
                ;;
            "activity")
                if [ -d "api/activity" ]; then
                    log_info "Building activity service..."
                    cd "api/activity"
                    docker build -t "${ACTIVITY_SERVICE_IMAGE}" .
                    docker save "${ACTIVITY_SERVICE_IMAGE}" -o /tmp/activity.tar
                    cd "$SCRIPT_DIR"
                fi
                ;;
            "version")
                if [ -d "api/version-service" ]; then
                    log_info "Building version service..."
                    cd "api/version-service"
                    docker build -t "${VERSION_SERVICE_IMAGE}" .
                    docker save "${VERSION_SERVICE_IMAGE}" -o /tmp/version.tar
                    cd "$SCRIPT_DIR"
                fi
                ;;
        esac
    done
    
    # Transfer and deploy
    log_info "Transferring updated images to AWS instance..."
    
    for service in "${SERVICES_TO_DEPLOY[@]}"; do
        if [ -f "/tmp/${service}.tar" ]; then
            scp -i "$AWS_KEY_PATH" "/tmp/${service}.tar" "$AWS_INSTANCE_USER@$AWS_INSTANCE_IP:~/mountain-service-images/"
            rm "/tmp/${service}.tar"
        fi
    done
    
    # Deploy on AWS instance
    log_info "Deploying updated services on AWS instance..."
    
    ssh -i "$AWS_KEY_PATH" "$AWS_INSTANCE_USER@$AWS_INSTANCE_IP" << EOF
        set -e
        
        echo "Loading updated Docker images..."
        cd ~/mountain-service-images
        for service in ${SERVICES_TO_DEPLOY[*]}; do
            if [ -f "\${service}.tar" ]; then
                echo "Loading \${service}.tar..."
                docker load -i "\${service}.tar"
            fi
        done
        
        echo "Restarting services..."
        cd ~/mountain-service-deployment
        
        # Stop specific services
        for service in ${SERVICES_TO_DEPLOY[*]}; do
            case \$service in
                "frontend")
                    docker-compose -f docker-compose.aws.yml stop frontend || true
                    ;;
                "employee")
                    docker-compose -f docker-compose.aws.yml stop employee-service || true
                    ;;
                "urgency")
                    docker-compose -f docker-compose.aws.yml stop urgency-service || true
                    ;;
                "activity")
                    docker-compose -f docker-compose.aws.yml stop activity-service || true
                    ;;
                "version")
                    docker-compose -f docker-compose.aws.yml stop version-service || true
                    ;;
            esac
        done
        
        # Start all services (Docker Compose will only restart the stopped ones)
        docker-compose -f docker-compose.aws.yml up -d
        
        echo "Waiting for services to start..."
        sleep 15
        
        echo "Checking service health..."
        docker-compose -f docker-compose.aws.yml ps
        
        echo "Quick deployment completed successfully!"
EOF
    
    log_success "Quick deployment completed successfully"
}

# Main deployment process
main() {
    echo -e "${GREEN}=== Mountain Rescue Service - Quick AWS Deployment ===${NC}"
    echo -e "${BLUE}Target: $AWS_INSTANCE_USER@$AWS_INSTANCE_IP${NC}"
    echo -e "${BLUE}Services: ${SERVICES_TO_DEPLOY[*]}${NC}"
    echo ""
    
    # Validate SSH connection
    if ! ssh -i "$AWS_KEY_PATH" -o ConnectTimeout=5 -o StrictHostKeyChecking=no "$AWS_INSTANCE_USER@$AWS_INSTANCE_IP" "echo 'SSH OK'" > /dev/null 2>&1; then
        log_error "Cannot connect to AWS instance via SSH"
        exit 1
    fi
    
    deploy_services
    
    echo ""
    log_success "Quick deployment completed successfully!"
    echo -e "${GREEN}Updated services are now running on AWS instance.${NC}"
    echo ""
    echo -e "${BLUE}You can check the status with:${NC}"
    echo -e "${YELLOW}  ssh -i $AWS_KEY_PATH $AWS_INSTANCE_USER@$AWS_INSTANCE_IP '~/monitor-services.sh'${NC}"
}

# Run main function
main "$@"
