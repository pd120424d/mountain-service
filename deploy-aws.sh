#!/bin/bash

# Multi-Cloud Deployment Script for Mountain Rescue Service
# This script builds Docker images and deploys them to either AWS EC2 or Azure VM
# Based on the DEPLOYMENT_TARGET environment variable

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
    echo -e "${GREEN}Loading environment variables from $ENV_FILE${NC}"
    source "$ENV_FILE"
else
    echo -e "${RED}Error: .env file not found. Please copy .env.aws to .env and configure it.${NC}"
    exit 1
fi

# Determine deployment target and set variables accordingly
DEPLOYMENT_TARGET=${DEPLOYMENT_TARGET:-aws}

if [ "$DEPLOYMENT_TARGET" = "azure" ]; then
    # Use Azure-compatible variable names (backward compatibility)
    INSTANCE_IP=${AZURE_VM_HOST:-${CLOUD_INSTANCE_IP}}
    INSTANCE_USER=${AZURE_VM_USER:-${CLOUD_INSTANCE_USER}}
    SSH_KEY_PATH=${AZURE_SSH_PRIVATE_KEY:-${CLOUD_SSH_KEY_PATH:-${SSH_KEY_PATH}}}
    COMPOSE_FILE="docker-compose.prod.yml"  # Use existing Azure compose file
else
    # Use AWS-prefixed variables
    INSTANCE_IP=${AWS_INSTANCE_IP:-${CLOUD_INSTANCE_IP}}
    INSTANCE_USER=${AWS_INSTANCE_USER:-${CLOUD_INSTANCE_USER}}
    SSH_KEY_PATH=${AWS_SSH_PRIVATE_KEY:-${CLOUD_SSH_KEY_PATH:-${SSH_KEY_PATH}}}
    COMPOSE_FILE="docker-compose.aws.yml"
fi

# Validate required environment variables (skip SSH_KEY_PATH if SSH_KEY_CONTENT is provided)
if [ -n "$SSH_KEY_CONTENT" ]; then
    required_vars=(
        "INSTANCE_IP"
        "INSTANCE_USER"
        "EMPLOYEE_SERVICE_IMAGE"
        "URGENCY_SERVICE_IMAGE"
        "ACTIVITY_SERVICE_IMAGE"
        "VERSION_SERVICE_IMAGE"
        "FRONTEND_IMAGE"
    )
else
    required_vars=(
        "INSTANCE_IP"
        "INSTANCE_USER"
        "SSH_KEY_PATH"
        "EMPLOYEE_SERVICE_IMAGE"
        "URGENCY_SERVICE_IMAGE"
        "ACTIVITY_SERVICE_IMAGE"
        "VERSION_SERVICE_IMAGE"
        "FRONTEND_IMAGE"
    )
fi

echo -e "${BLUE}Validating environment variables...${NC}"
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo -e "${RED}Error: $var is not set in .env file${NC}"
        exit 1
    fi
done

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if SSH key exists
check_ssh_key() {
    # Handle SSH key from environment variable (for CI/CD) or file path
    if [ -n "$SSH_KEY_CONTENT" ]; then
        # SSH key provided as environment variable content
        echo "$SSH_KEY_CONTENT" > /tmp/ssh_key
        SSH_KEY_PATH="/tmp/ssh_key"
        chmod 600 "$SSH_KEY_PATH"
    elif [ ! -f "$SSH_KEY_PATH" ]; then
        log_error "SSH key not found at $SSH_KEY_PATH"
        exit 1
    fi

    # Check key permissions
    key_perms=$(stat -c "%a" "$SSH_KEY_PATH" 2>/dev/null || stat -f "%A" "$SSH_KEY_PATH" 2>/dev/null)
    if [ "$key_perms" != "600" ]; then
        log_warning "SSH key permissions are $key_perms, should be 600. Fixing..."
        chmod 600 "$SSH_KEY_PATH"
    fi
}

# Test SSH connection
test_ssh_connection() {
    log_info "Testing SSH connection to $DEPLOYMENT_TARGET instance..."
    if ssh -i "$SSH_KEY_PATH" -o ConnectTimeout=10 -o StrictHostKeyChecking=no "$INSTANCE_USER@$INSTANCE_IP" "echo 'SSH connection successful'" > /dev/null 2>&1; then
        log_success "SSH connection to $DEPLOYMENT_TARGET instance successful"
    else
        log_error "Cannot connect to $DEPLOYMENT_TARGET instance via SSH"
        log_error "Please check:"
        log_error "  - INSTANCE_IP: $INSTANCE_IP"
        log_error "  - INSTANCE_USER: $INSTANCE_USER"
        log_error "  - SSH_KEY_PATH: $SSH_KEY_PATH"
        log_error "  - Security group/firewall allows SSH (port 22)"
        log_error "  - Instance is running"
        exit 1
    fi
}

# Build Docker images
build_images() {
    log_info "Building Docker images..."

    # Build frontend image
    log_info "Building frontend image..."
    if [ "$DEPLOYMENT_TARGET" = "azure" ]; then
        docker build -f ui/Dockerfile -t "${FRONTEND_IMAGE}" ui/
    else
        docker build -f ui/Dockerfile.aws -t "${FRONTEND_IMAGE}" ui/
    fi

    # Build backend services with correct context (matching working backend-deploy.yml)
    cd "$SCRIPT_DIR"

    # Employee service
    if [ -d "api/employee" ]; then
        log_info "Building employee service image..."
        docker build -f api/employee/Dockerfile -t "${EMPLOYEE_SERVICE_IMAGE}" api/
    fi

    # Urgency service
    if [ -d "api/urgency" ]; then
        log_info "Building urgency service image..."
        docker build -f api/urgency/Dockerfile -t "${URGENCY_SERVICE_IMAGE}" api/
    fi

    # Activity service
    if [ -d "api/activity" ]; then
        log_info "Building activity service image..."
        docker build -f api/activity/Dockerfile -t "${ACTIVITY_SERVICE_IMAGE}" api/
    fi

    # Version service (with build args)
    if [ -d "api/version-service" ]; then
        log_info "Building version service image..."
        # Get version and git SHA for version service
        VERSION=$(git describe --tags --abbrev=0 --match "v*" 2>/dev/null || echo dev)
        GIT_SHA=$(git rev-parse --short HEAD)
        docker build -f api/version-service/Dockerfile \
            --build-arg VERSION="$VERSION" \
            --build-arg GIT_SHA="$GIT_SHA" \
            -t "${VERSION_SERVICE_IMAGE}" api/
    fi

    log_success "All Docker images built successfully"
}

# Save and transfer Docker images
transfer_images() {
    log_info "Saving Docker images to tar files..."

    # Create temporary directory for image files
    mkdir -p /tmp/mountain-service-images

    # Save images with compression and verify
    images=(
        "${FRONTEND_IMAGE}:frontend"
        "${EMPLOYEE_SERVICE_IMAGE}:employee"
        "${URGENCY_SERVICE_IMAGE}:urgency"
        "${ACTIVITY_SERVICE_IMAGE}:activity"
        "${VERSION_SERVICE_IMAGE}:version"
    )

    for image_info in "${images[@]}"; do
        IFS=':' read -r image_name file_name <<< "$image_info"
        log_info "Saving $image_name..."

        # Save image
        docker save "$image_name" | gzip > "/tmp/mountain-service-images/${file_name}.tar.gz"

        # Verify the saved file
        if [ ! -s "/tmp/mountain-service-images/${file_name}.tar.gz" ]; then
            log_error "Failed to save $image_name - file is empty"
            exit 1
        fi

        # Test the compressed file
        if ! gzip -t "/tmp/mountain-service-images/${file_name}.tar.gz"; then
            log_error "Corrupted compressed file for $image_name"
            exit 1
        fi

        log_success "Successfully saved and verified $image_name"
    done

    log_info "Transferring Docker images to $DEPLOYMENT_TARGET instance..."

    # Create directory on remote server
    ssh -i "$SSH_KEY_PATH" "$INSTANCE_USER@$INSTANCE_IP" "mkdir -p ~/mountain-service-images"

    # Transfer image files with verification
    for image_info in "${images[@]}"; do
        IFS=':' read -r image_name file_name <<< "$image_info"
        log_info "Transferring ${file_name}.tar.gz..."

        # Transfer with retry logic
        max_retries=3
        retry_count=0

        while [ $retry_count -lt $max_retries ]; do
            if scp -i "$SSH_KEY_PATH" "/tmp/mountain-service-images/${file_name}.tar.gz" "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-images/"; then
                log_success "Successfully transferred ${file_name}.tar.gz"
                break
            else
                retry_count=$((retry_count + 1))
                log_warning "Transfer failed for ${file_name}.tar.gz, retry $retry_count/$max_retries"
                sleep 2
            fi
        done

        if [ $retry_count -eq $max_retries ]; then
            log_error "Failed to transfer ${file_name}.tar.gz after $max_retries attempts"
            exit 1
        fi
    done

    # Clean up local files
    rm -rf /tmp/mountain-service-images

    log_success "Docker images transferred successfully"
}

# Transfer deployment files
transfer_deployment_files() {
    log_info "Transferring deployment files to $DEPLOYMENT_TARGET instance..."

    # Create deployment directory on remote server
    ssh -i "$SSH_KEY_PATH" "$INSTANCE_USER@$INSTANCE_IP" "mkdir -p ~/mountain-service-deployment"

    # Transfer necessary files
    scp -i "$SSH_KEY_PATH" "$COMPOSE_FILE" "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/"
    scp -i "$SSH_KEY_PATH" .env "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/"

    # Transfer secrets directory if it exists
    if [ -d "api" ]; then
        scp -i "$SSH_KEY_PATH" -r api "$INSTANCE_USER@$INSTANCE_IP:~/mountain-service-deployment/"
    fi

    log_success "Deployment files transferred successfully"
}

# Deploy on cloud instance
deploy_on_cloud() {
    log_info "Deploying application on $DEPLOYMENT_TARGET instance..."

    # Execute deployment commands on remote server
    ssh -i "$SSH_KEY_PATH" -T "$INSTANCE_USER@$INSTANCE_IP" << 'EOF'
        set -e

        echo "Loading Docker images..."
        cd ~/mountain-service-images

        # Show available images before loading
        echo "Available tar files:"
        ls -la *.tar 2>/dev/null || echo "No tar files found"

        # Load all available image files (both .tar.gz and .tar)
        loaded_count=0

        # Load compressed images first
        for image in *.tar.gz; do
            if [ -f "$image" ]; then
                echo "Loading compressed $image..."
                if gzip -t "$image" && gzip -dc "$image" | docker load; then
                    echo "✓ Successfully loaded $image"
                    loaded_count=$((loaded_count + 1))
                else
                    echo "✗ Failed to load $image"
                    exit 1
                fi
            fi
        done

        # Load uncompressed tar files
        for image in *.tar; do
            if [ -f "$image" ]; then
                echo "Loading $image..."
                if docker load -i "$image"; then
                    echo "✓ Successfully loaded $image"
                    loaded_count=$((loaded_count + 1))
                else
                    echo "✗ Failed to load $image"
                    exit 1
                fi
            fi
        done

        echo "Successfully loaded $loaded_count Docker images"

        echo "Stopping existing containers..."
        cd ~/mountain-service-deployment

        # Determine which compose file to use based on deployment target
        if [ -f "docker-compose.aws.yml" ]; then
            COMPOSE_FILE="docker-compose.aws.yml"
        elif [ -f "docker-compose.prod.yml" ]; then
            COMPOSE_FILE="docker-compose.prod.yml"
        else
            echo "Error: No compose file found"
            exit 1
        fi

        echo "Using compose file: $COMPOSE_FILE"

        docker compose -f "$COMPOSE_FILE" down || true

        echo "Starting new containers..."
        docker compose -f "$COMPOSE_FILE" up -d

        echo "Waiting for services to start..."
        sleep 30

        echo "Checking service health..."
        docker compose -f "$COMPOSE_FILE" ps

        echo "Deployment completed successfully!"
EOF

    log_success "Application deployed successfully on $DEPLOYMENT_TARGET instance"
}

# Alternative: Pull images from registry (like working backend-deploy.yml)
pull_images_from_registry() {
    log_info "Pulling Docker images from registry on $DEPLOYMENT_TARGET instance..."

    # Execute on remote server
    ssh -i "$SSH_KEY_PATH" "$INSTANCE_USER@$INSTANCE_IP" << EOF
        set -e

        echo "Logging into GitHub Container Registry..."
        echo "$GHCR_PAT" | docker login ghcr.io -u "$GITHUB_ACTOR" --password-stdin

        echo "Pulling images from registry..."
        docker pull "${EMPLOYEE_SERVICE_IMAGE}" || echo "Warning: Failed to pull employee service image"
        docker pull "${URGENCY_SERVICE_IMAGE}" || echo "Warning: Failed to pull urgency service image"
        docker pull "${ACTIVITY_SERVICE_IMAGE}" || echo "Warning: Failed to pull activity service image"
        docker pull "${VERSION_SERVICE_IMAGE}" || echo "Warning: Failed to pull version service image"
        docker pull "${FRONTEND_IMAGE}" || echo "Warning: Failed to pull frontend image"

        echo "Successfully pulled images from registry"
EOF

    log_success "Images pulled from registry successfully"
}

# Main deployment process
main() {
    echo -e "${GREEN}=== Mountain Rescue Service - Multi-Cloud Deployment ===${NC}"
    echo -e "${BLUE}Target Platform: $DEPLOYMENT_TARGET${NC}"
    echo -e "${BLUE}Target Instance: $INSTANCE_USER@$INSTANCE_IP${NC}"
    echo ""

    check_ssh_key
    test_ssh_connection

    # Choose deployment method based on environment
    if [ -n "$GHCR_PAT" ] && [ -n "$GITHUB_ACTOR" ]; then
        log_info "Using registry-based deployment (recommended)"
        pull_images_from_registry
    else
        log_info "Using local build and transfer method"
        build_images
        transfer_images
    fi

    transfer_deployment_files
    deploy_on_cloud

    echo ""
    log_success "Deployment completed successfully!"
    echo -e "${GREEN}Your application should now be available at:${NC}"
    echo -e "${BLUE}  Frontend: http://$INSTANCE_IP${NC}"
    echo -e "${BLUE}  Employee API: http://$INSTANCE_IP:8082${NC}"
    echo -e "${BLUE}  Urgency API: http://$INSTANCE_IP:8083${NC}"
    echo -e "${BLUE}  Activity API: http://$INSTANCE_IP:8084${NC}"
    echo -e "${BLUE}  Version API: http://$INSTANCE_IP:8090${NC}"
    echo ""
    if [ "$DEPLOYMENT_TARGET" = "aws" ]; then
        echo -e "${YELLOW}Note: Make sure your AWS security group allows inbound traffic on these ports.${NC}"
    else
        echo -e "${YELLOW}Note: Make sure your Azure network security group allows inbound traffic on these ports.${NC}"
    fi
}

# Run main function
main "$@"
