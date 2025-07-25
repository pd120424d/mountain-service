# Multi-Cloud Deployment Guide

This guide covers deploying the Mountain Rescue Service to both Azure VMs and AWS EC2 instances, maintaining compatibility with existing Azure deployments while adding AWS support.

## Table of Contents

- [Overview](#overview)
- [GitHub Secrets Configuration](#github-secrets-configuration)
- [Azure VM Deployment](#azure-vm-deployment)
- [AWS EC2 Deployment](#aws-ec2-deployment)
- [Local Development](#local-development)
- [Troubleshooting](#troubleshooting)

## Overview

The Mountain Rescue Service supports deployment to multiple cloud platforms:

- **Azure VM**: Original deployment target (existing)
- **AWS EC2**: New deployment target (added)

Both deployments use the same Docker images and configuration, with platform-specific networking and security configurations.

## GitHub Secrets Configuration

### Required Secrets for Both Platforms

Add these secrets to your GitHub repository (`Settings > Secrets and variables > Actions`):

#### Container Registry
```
GHCR_PAT=your_github_personal_access_token
```

#### Application Configuration
```
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ADMIN_PASSWORD=your-admin-password-change-this
SERVICE_AUTH_SECRET=your-service-auth-secret-change-this
```

#### Multi-Cloud Instance Configuration
```
# Primary cloud instance (can be Azure or AWS)
CLOUD_INSTANCE_IP=your.instance.public.ip
CLOUD_INSTANCE_USER=azureuser  # or ec2-user for AWS
CLOUD_SSH_PRIVATE_KEY=-----BEGIN OPENSSH PRIVATE KEY-----
...your private key content...
-----END OPENSSH PRIVATE KEY-----

# Deployment target selector
DEPLOYMENT_TARGET=azure  # or 'aws'
```

#### Legacy Azure Support (for backward compatibility)
```
AZURE_VM_HOST=${CLOUD_INSTANCE_IP}
AZURE_VM_USER=${CLOUD_INSTANCE_USER}
AZURE_SSH_PRIVATE_KEY=${CLOUD_SSH_PRIVATE_KEY}
```

#### AWS-Specific Secrets (when using AWS)
```
AWS_INSTANCE_IP=${CLOUD_INSTANCE_IP}
AWS_INSTANCE_USER=${CLOUD_INSTANCE_USER}
AWS_SSH_KEY_PATH=${CLOUD_SSH_PRIVATE_KEY}
```

#### CORS and API Configuration
```
CORS_ALLOWED_ORIGINS=http://your.instance.ip,https://your.instance.ip
SWAGGER_API_URL=http://your.instance.ip:8082/swagger.json
EMPLOYEE_SERVICE_URL=http://employee-service:8082
ACTIVITY_SERVICE_URL=http://activity-service:8084
```

### Environment-Specific Secrets

You can also create environment-specific secret sets:

#### For Azure Production
```
AZURE_PROD_INSTANCE_IP=your.azure.vm.ip
AZURE_PROD_INSTANCE_USER=azureuser
AZURE_PROD_SSH_KEY=your_azure_private_key
```

#### For AWS Production
```
AWS_PROD_INSTANCE_IP=your.aws.ec2.ip
AWS_PROD_INSTANCE_USER=ec2-user
AWS_PROD_SSH_KEY=your_aws_private_key
```

## Azure VM Deployment

### Prerequisites

1. Azure VM running Ubuntu 20.04+ or similar
2. SSH access configured
3. Security group allowing ports: 22, 80, 443, 8082-8084, 8090, 9082-9084

### Setup Azure VM

```bash
# Clone the repository
git clone https://github.com/pd120424d/mountain-service.git
cd mountain-service

# Copy and configure environment file
cp .env.aws .env
nano .env  # Update with Azure-specific values

# Set deployment target to Azure
echo "DEPLOYMENT_TARGET=azure" >> .env
echo "CLOUD_INSTANCE_IP=your.azure.vm.ip" >> .env
echo "CLOUD_INSTANCE_USER=azureuser" >> .env

# Setup the Azure VM
./setup-azure-vm.sh  # Use existing Azure setup script
```

### Deploy to Azure

```bash
# Full deployment
./deploy-azure.sh  # Use existing Azure deployment script

# Or use the new multi-cloud script
./deploy-aws.sh  # Works for both Azure and AWS based on DEPLOYMENT_TARGET
```

## AWS EC2 Deployment

### Prerequisites

1. AWS EC2 instance running Amazon Linux 2 or Ubuntu 20.04+
2. SSH access configured with key pair
3. Security group allowing ports: 22, 80, 443, 8082-8084, 8090, 9082-9084

### Setup AWS EC2 Instance

```bash
# Clone the repository
git clone https://github.com/pd120424d/mountain-service.git
cd mountain-service

# Copy and configure environment file
cp .env.aws .env
nano .env  # Update with AWS-specific values

# Set deployment target to AWS
echo "DEPLOYMENT_TARGET=aws" >> .env
echo "CLOUD_INSTANCE_IP=your.aws.ec2.ip" >> .env
echo "CLOUD_INSTANCE_USER=ec2-user" >> .env

# Setup the AWS instance
./setup-aws-instance.sh
```

### Deploy to AWS

```bash
# Full deployment
./deploy-aws.sh

# Quick deployment (specific services)
./quick-deploy-aws.sh --frontend --employee
```

## Local Development

### Environment Configuration

Create a local `.env` file:

```bash
# Copy the template
cp .env.aws .env

# Configure for local development
cat > .env << EOF
DEPLOYMENT_TARGET=local
CLOUD_INSTANCE_IP=localhost
EMPLOYEE_SERVICE_IMAGE=mountain-service-employee:latest
URGENCY_SERVICE_IMAGE=mountain-service-urgency:latest
ACTIVITY_SERVICE_IMAGE=mountain-service-activity:latest
VERSION_SERVICE_IMAGE=mountain-service-version:latest
FRONTEND_IMAGE=mountain-service-frontend:latest
JWT_SECRET=local-dev-secret
ADMIN_PASSWORD=admin123
SERVICE_AUTH_SECRET=local-service-secret
CORS_ALLOWED_ORIGINS=http://localhost:4200,http://localhost:80
EOF
```

### Local Deployment

```bash
# Build and run locally
docker-compose -f docker-compose.aws.yml up --build

# Or use the existing staging setup
docker-compose -f docker-compose.staging.yml up --build
```

## GitHub Actions Integration

### Workflow Configuration

Update your GitHub Actions workflows to use the new secrets:

```yaml
# .github/workflows/deploy.yml
name: Multi-Cloud Deploy

on:
  push:
    branches: [main]
    tags: ['v*']

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup environment
        run: |
          echo "DEPLOYMENT_TARGET=${{ secrets.DEPLOYMENT_TARGET }}" >> .env
          echo "CLOUD_INSTANCE_IP=${{ secrets.CLOUD_INSTANCE_IP }}" >> .env
          echo "CLOUD_INSTANCE_USER=${{ secrets.CLOUD_INSTANCE_USER }}" >> .env
          echo "JWT_SECRET=${{ secrets.JWT_SECRET }}" >> .env
          echo "ADMIN_PASSWORD=${{ secrets.ADMIN_PASSWORD }}" >> .env
          echo "SERVICE_AUTH_SECRET=${{ secrets.SERVICE_AUTH_SECRET }}" >> .env
          echo "CORS_ALLOWED_ORIGINS=${{ secrets.CORS_ALLOWED_ORIGINS }}" >> .env
      
      - name: Deploy to Cloud
        run: |
          echo "${{ secrets.CLOUD_SSH_PRIVATE_KEY }}" > /tmp/ssh_key
          chmod 600 /tmp/ssh_key
          export CLOUD_SSH_KEY_PATH=/tmp/ssh_key
          
          if [ "${{ secrets.DEPLOYMENT_TARGET }}" = "aws" ]; then
            ./deploy-aws.sh
          else
            ./deploy-azure.sh  # Use existing Azure deployment
          fi
```

## Security Group Configuration

### Azure Network Security Group

```bash
# Allow HTTP/HTTPS
az network nsg rule create --resource-group myResourceGroup \
  --nsg-name myNetworkSecurityGroup --name AllowHTTP \
  --protocol tcp --priority 1000 --destination-port-range 80

az network nsg rule create --resource-group myResourceGroup \
  --nsg-name myNetworkSecurityGroup --name AllowHTTPS \
  --protocol tcp --priority 1001 --destination-port-range 443

# Allow API ports
az network nsg rule create --resource-group myResourceGroup \
  --nsg-name myNetworkSecurityGroup --name AllowAPIs \
  --protocol tcp --priority 1002 --destination-port-range 8082-8090
```

### AWS Security Group

```bash
# Allow HTTP/HTTPS
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxxxxxx \
  --protocol tcp --port 80 --cidr 0.0.0.0/0

aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxxxxxx \
  --protocol tcp --port 443 --cidr 0.0.0.0/0

# Allow API ports
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxxxxxx \
  --protocol tcp --port 8082-8090 --cidr 0.0.0.0/0
```

## Monitoring and Maintenance

### Health Checks

Both platforms support the same health check endpoints:

```bash
# Check application health
curl http://your.instance.ip/health
curl http://your.instance.ip/api/v1/health

# Check individual services
curl http://your.instance.ip:8082/api/v1/health  # Employee service
curl http://your.instance.ip:8083/api/v1/health  # Urgency service
curl http://your.instance.ip:8084/api/v1/health  # Activity service
curl http://your.instance.ip:8090/api/v1/health  # Version service
```

### Monitoring Scripts

Both platforms include monitoring scripts:

```bash
# On the instance
~/monitor-services.sh    # Check service status
~/backup-data.sh         # Create data backup
```

## Troubleshooting

### Common Issues

1. **SSH Connection Failed**
   - Verify instance IP and user
   - Check SSH key permissions (should be 600)
   - Ensure security group allows SSH (port 22)

2. **Docker Images Not Found**
   - Check GitHub Container Registry authentication
   - Verify image names in .env file
   - Ensure images are built and pushed

3. **Services Not Starting**
   - Check Docker logs: `docker-compose logs -f`
   - Verify environment variables
   - Check database connectivity

4. **Port Access Issues**
   - Verify security group/firewall rules
   - Check if services are binding to correct ports
   - Ensure no port conflicts

### Platform-Specific Issues

#### Azure VM
- Use `azureuser` as default user
- Check Azure Network Security Group rules
- Verify VM size has sufficient resources

#### AWS EC2
- Use `ec2-user` for Amazon Linux, `ubuntu` for Ubuntu
- Check AWS Security Group rules
- Verify instance type has sufficient resources
- Ensure EBS volumes have enough space

### Getting Help

1. Check service logs: `docker-compose logs -f [service-name]`
2. Run health checks: `~/monitor-services.sh`
3. Check system resources: `htop`, `df -h`
4. Review deployment logs in GitHub Actions

## Testing Your Deployment

### Pre-Deployment Checklist

Before deploying, ensure you have:

- [ ] Cloud instance (Azure VM or AWS EC2) running and accessible
- [ ] SSH key configured with proper permissions (600)
- [ ] Security group/firewall rules configured for required ports
- [ ] GitHub secrets configured correctly
- [ ] Environment file (.env) configured with correct values

### Required Ports

Ensure your security group/firewall allows inbound traffic on:

- **22**: SSH access
- **80**: HTTP (frontend)
- **443**: HTTPS (if using SSL)
- **8082**: Employee service API
- **8083**: Urgency service API
- **8084**: Activity service API
- **8090**: Version service API
- **9082-9084**: Swagger UI interfaces (optional)

### Local Testing

Test your configuration locally before deploying:

```bash
# 1. Copy and configure environment file
cp .env.aws .env
nano .env  # Update with your values

# 2. Test Docker builds locally
cd ui && docker build -f Dockerfile.aws -t test-frontend .
cd ../api/employee && docker build -t test-employee .

# 3. Test local deployment
docker-compose -f docker-compose.aws.yml up --build
```

### Deployment Testing

#### 1. Test SSH Connection

```bash
# Test SSH access to your instance
ssh -i ~/.ssh/your-key.pem ec2-user@your.instance.ip

# Or for Azure
ssh -i ~/.ssh/your-key.pem azureuser@your.instance.ip
```

#### 2. Run Setup Script

```bash
# For AWS
./setup-aws-instance.sh

# For Azure (use existing setup)
./setup-azure-vm.sh
```

#### 3. Deploy Application

```bash
# Full deployment
./deploy-aws.sh

# Quick deployment (specific services)
./quick-deploy-aws.sh --frontend --employee
```

#### 4. Verify Deployment

```bash
# Check service health
curl http://your.instance.ip/health
curl http://your.instance.ip/api/v1/health

# Check individual services
curl http://your.instance.ip:8082/api/v1/health  # Employee
curl http://your.instance.ip:8083/api/v1/health  # Urgency
curl http://your.instance.ip:8084/api/v1/health  # Activity
curl http://your.instance.ip:8090/api/v1/health  # Version

# Check frontend
curl -I http://your.instance.ip  # Should return 200 OK
```

### GitHub Actions Testing

#### Manual Deployment Trigger

You can manually trigger deployment from GitHub:

1. Go to your repository on GitHub
2. Click "Actions" tab
3. Select "Multi-Cloud Deployment" workflow
4. Click "Run workflow"
5. Choose deployment target (azure/aws)
6. Click "Run workflow"

#### Monitoring Deployment

Monitor your deployment in GitHub Actions:

1. Check workflow logs for any errors
2. Verify all steps complete successfully
3. Check deployment summary for service URLs
4. Test the deployed application

### Post-Deployment Validation

#### Functional Testing

1. **Frontend Access**
   ```bash
   curl -I http://your.instance.ip
   # Should return: HTTP/1.1 200 OK
   ```

2. **API Endpoints**
   ```bash
   # Test employee service
   curl http://your.instance.ip:8082/api/v1/employees

   # Test urgency service
   curl http://your.instance.ip:8083/api/v1/urgencies

   # Test activity service
   curl http://your.instance.ip:8084/api/v1/activities
   ```

3. **Swagger Documentation**
   ```bash
   # Access Swagger UIs
   curl -I http://your.instance.ip/employee-swagger/
   curl -I http://your.instance.ip/urgency-swagger/
   curl -I http://your.instance.ip/activity-swagger/
   ```

#### Performance Testing

```bash
# Test response times
time curl http://your.instance.ip/api/v1/health

# Test concurrent requests
for i in {1..10}; do
  curl -s http://your.instance.ip/health &
done
wait
```

#### Security Testing

```bash
# Test HTTPS redirect (if SSL enabled)
curl -I http://your.instance.ip
# Should redirect to HTTPS if configured

# Test security headers
curl -I http://your.instance.ip
# Should include security headers like X-Frame-Options, X-Content-Type-Options
```

### Monitoring and Maintenance

#### Service Monitoring

On your instance, use the monitoring script:

```bash
# Check service status
~/monitor-services.sh

# View logs
docker-compose logs -f

# Check resource usage
htop
df -h
```

#### Backup and Recovery

```bash
# Create backup
~/backup-data.sh

# List backups
ls -la ~/backups/

# Restore from backup (if needed)
# Follow the restore procedures in the backup script
```

### Rollback Procedures

If deployment fails or issues occur:

#### Quick Rollback

```bash
# SSH to your instance
ssh -i ~/.ssh/your-key.pem user@your.instance.ip

# Stop current services
cd ~/mountain-service-deployment
docker-compose down

# Restore from backup (if available)
# Or redeploy previous version
```

#### GitHub Actions Rollback

1. Find the last successful deployment commit
2. Revert to that commit or create a new commit with fixes
3. Push to trigger new deployment
4. Monitor the deployment process

### Common Issues and Solutions

#### Issue: SSH Connection Refused
```bash
# Check instance status
aws ec2 describe-instances --instance-ids i-1234567890abcdef0
# or
az vm show --resource-group myResourceGroup --name myVM

# Check security group rules
aws ec2 describe-security-groups --group-ids sg-1234567890abcdef0
# or
az network nsg rule list --resource-group myResourceGroup --nsg-name myNSG
```

#### Issue: Docker Images Not Found
```bash
# Check if images exist in registry
docker pull ghcr.io/your-org/mountain-service-frontend:latest

# Re-authenticate to registry
echo $GHCR_PAT | docker login ghcr.io -u your-username --password-stdin
```

#### Issue: Services Not Starting
```bash
# Check logs
docker-compose logs employee-service
docker-compose logs urgency-service

# Check environment variables
docker-compose config

# Restart specific service
docker-compose restart employee-service
```

This comprehensive testing guide ensures your deployment is successful and helps troubleshoot common issues.
