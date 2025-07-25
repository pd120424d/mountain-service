#!/bin/bash

# AWS Instance Setup Script for Mountain Rescue Service
# This script prepares an AWS EC2 instance for deployment

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

# Validate required environment variables
required_vars=(
    "AWS_INSTANCE_IP"
    "AWS_INSTANCE_USER"
    "AWS_KEY_PATH"
)

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
    if [ ! -f "$AWS_KEY_PATH" ]; then
        log_error "SSH key not found at $AWS_KEY_PATH"
        exit 1
    fi
    
    # Check key permissions
    key_perms=$(stat -c "%a" "$AWS_KEY_PATH" 2>/dev/null || stat -f "%A" "$AWS_KEY_PATH" 2>/dev/null)
    if [ "$key_perms" != "600" ]; then
        log_warning "SSH key permissions are $key_perms, should be 600. Fixing..."
        chmod 600 "$AWS_KEY_PATH"
    fi
}

# Test SSH connection
test_ssh_connection() {
    log_info "Testing SSH connection to AWS instance..."
    if ssh -i "$AWS_KEY_PATH" -o ConnectTimeout=10 -o StrictHostKeyChecking=no "$AWS_INSTANCE_USER@$AWS_INSTANCE_IP" "echo 'SSH connection successful'" > /dev/null 2>&1; then
        log_success "SSH connection to AWS instance successful"
    else
        log_error "Cannot connect to AWS instance via SSH"
        log_error "Please check:"
        log_error "  - AWS_INSTANCE_IP: $AWS_INSTANCE_IP"
        log_error "  - AWS_INSTANCE_USER: $AWS_INSTANCE_USER"
        log_error "  - AWS_KEY_PATH: $AWS_KEY_PATH"
        log_error "  - Security group allows SSH (port 22)"
        log_error "  - Instance is running"
        exit 1
    fi
}

# Setup AWS instance
setup_instance() {
    log_info "Setting up AWS instance for Mountain Rescue Service deployment..."
    
    ssh -i "$AWS_KEY_PATH" "$AWS_INSTANCE_USER@$AWS_INSTANCE_IP" << 'EOF'
        set -e
        
        echo "Updating system packages..."
        sudo yum update -y
        
        echo "Installing Docker..."
        sudo yum install -y docker
        sudo systemctl start docker
        sudo systemctl enable docker
        sudo usermod -a -G docker $USER
        
        echo "Installing Docker Compose..."
        sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        
        echo "Installing additional tools..."
        sudo yum install -y git curl wget htop
        
        echo "Creating application directories..."
        mkdir -p ~/mountain-service-deployment
        mkdir -p ~/mountain-service-images
        mkdir -p ~/logs
        
        echo "Setting up log rotation..."
        sudo tee /etc/logrotate.d/mountain-service > /dev/null << 'LOGROTATE'
/home/ec2-user/logs/*.log {
    daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    copytruncate
}
LOGROTATE
        
        echo "Setting up firewall rules..."
        # Note: AWS security groups handle most firewall rules
        # But we can set up local iptables if needed
        
        echo "Creating systemd service for auto-start..."
        sudo tee /etc/systemd/system/mountain-service.service > /dev/null << 'SERVICE'
[Unit]
Description=Mountain Rescue Service
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/ec2-user/mountain-service-deployment
ExecStart=/usr/local/bin/docker-compose -f docker-compose.aws.yml up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.aws.yml down
User=ec2-user
Group=ec2-user

[Install]
WantedBy=multi-user.target
SERVICE
        
        echo "Enabling mountain-service to start on boot..."
        sudo systemctl enable mountain-service
        
        echo "Setting up monitoring script..."
        tee ~/monitor-services.sh > /dev/null << 'MONITOR'
#!/bin/bash
# Simple monitoring script for Mountain Rescue Service

echo "=== Mountain Rescue Service Status ==="
echo "Date: $(date)"
echo ""

echo "Docker containers:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""

echo "System resources:"
echo "Memory usage:"
free -h
echo ""
echo "Disk usage:"
df -h /
echo ""

echo "Service health checks:"
services=("8082" "8083" "8084" "8090" "80")
for port in "${services[@]}"; do
    if curl -s -f "http://localhost:$port/health" > /dev/null 2>&1 || curl -s -f "http://localhost:$port/api/v1/health" > /dev/null 2>&1; then
        echo "✓ Port $port: Healthy"
    else
        echo "✗ Port $port: Not responding"
    fi
done
MONITOR
        
        chmod +x ~/monitor-services.sh
        
        echo "Setting up backup script..."
        tee ~/backup-data.sh > /dev/null << 'BACKUP'
#!/bin/bash
# Backup script for Mountain Rescue Service

BACKUP_DIR="$HOME/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"

echo "Creating backup for $DATE..."

# Backup Docker volumes
docker run --rm -v mountain-service_db_data_aws:/data -v "$BACKUP_DIR":/backup alpine tar czf "/backup/db_data_$DATE.tar.gz" -C /data .
docker run --rm -v mountain-service_urgency_db_data_aws:/data -v "$BACKUP_DIR":/backup alpine tar czf "/backup/urgency_db_data_$DATE.tar.gz" -C /data .
docker run --rm -v mountain-service_activity_db_data_aws:/data -v "$BACKUP_DIR":/backup alpine tar czf "/backup/activity_db_data_$DATE.tar.gz" -C /data .

# Backup configuration files
tar czf "$BACKUP_DIR/config_$DATE.tar.gz" -C "$HOME" mountain-service-deployment

echo "Backup completed: $BACKUP_DIR"

# Clean up old backups (keep last 7 days)
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +7 -delete
BACKUP

        chmod +x ~/backup-data.sh
        
        echo "Setting up cron job for monitoring..."
        (crontab -l 2>/dev/null; echo "*/5 * * * * $HOME/monitor-services.sh >> $HOME/logs/monitor.log 2>&1") | crontab -
        
        echo "Setting up cron job for backups..."
        (crontab -l 2>/dev/null; echo "0 2 * * * $HOME/backup-data.sh >> $HOME/logs/backup.log 2>&1") | crontab -
        
        echo "AWS instance setup completed successfully!"
        echo ""
        echo "Installed components:"
        echo "- Docker: $(docker --version)"
        echo "- Docker Compose: $(docker-compose --version)"
        echo "- Git: $(git --version)"
        echo ""
        echo "Created directories:"
        echo "- ~/mountain-service-deployment (for deployment files)"
        echo "- ~/mountain-service-images (for Docker images)"
        echo "- ~/logs (for application logs)"
        echo "- ~/backups (for data backups)"
        echo ""
        echo "Created scripts:"
        echo "- ~/monitor-services.sh (service monitoring)"
        echo "- ~/backup-data.sh (data backup)"
        echo ""
        echo "Systemd service: mountain-service (auto-start on boot)"
        echo ""
        echo "Note: You may need to log out and back in for Docker group membership to take effect."
EOF
    
    log_success "AWS instance setup completed successfully"
}

# Main setup process
main() {
    echo -e "${GREEN}=== Mountain Rescue Service - AWS Instance Setup ===${NC}"
    echo -e "${BLUE}Target: $AWS_INSTANCE_USER@$AWS_INSTANCE_IP${NC}"
    echo ""
    
    check_ssh_key
    test_ssh_connection
    setup_instance
    
    echo ""
    log_success "AWS instance setup completed successfully!"
    echo -e "${GREEN}Your AWS instance is now ready for deployment.${NC}"
    echo ""
    echo -e "${BLUE}Next steps:${NC}"
    echo -e "${YELLOW}1. Update your .env file with the correct AWS_INSTANCE_IP${NC}"
    echo -e "${YELLOW}2. Run ./deploy-aws.sh to deploy the application${NC}"
    echo ""
    echo -e "${BLUE}Useful commands on the AWS instance:${NC}"
    echo -e "${YELLOW}  ~/monitor-services.sh    - Check service status${NC}"
    echo -e "${YELLOW}  ~/backup-data.sh         - Create data backup${NC}"
    echo -e "${YELLOW}  docker-compose logs -f   - View application logs${NC}"
    echo -e "${YELLOW}  sudo systemctl status mountain-service - Check service status${NC}"
}

# Run main function
main "$@"
