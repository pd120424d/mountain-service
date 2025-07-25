# Quick Start Guide: AWS Deployment

This guide will help you quickly deploy the Mountain Rescue Service to AWS EC2 while maintaining compatibility with your existing Azure setup.

## Prerequisites

- AWS EC2 instance running (Amazon Linux 2 or Ubuntu 20.04+)
- SSH key pair for EC2 access
- GitHub repository with the Mountain Rescue Service code
- Docker installed locally (for testing)

## Step 1: Configure GitHub Secrets

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

### AWS-Specific Secrets (Won't Conflict with Azure)
```
AWS_DEPLOYMENT_TARGET=aws
AWS_INSTANCE_IP=YOUR_AWS_EC2_PUBLIC_IP
AWS_INSTANCE_USER=ec2-user
AWS_SSH_PRIVATE_KEY=-----BEGIN OPENSSH PRIVATE KEY-----
...your complete AWS SSH private key content...
-----END OPENSSH PRIVATE KEY-----
AWS_JWT_SECRET=your-aws-specific-jwt-secret-64-chars-long
AWS_ADMIN_PASSWORD=your-aws-secure-admin-password
AWS_SERVICE_AUTH_SECRET=your-aws-service-auth-secret
AWS_CORS_ALLOWED_ORIGINS=http://YOUR_AWS_EC2_PUBLIC_IP,https://YOUR_AWS_EC2_PUBLIC_IP
AWS_COMPOSE_ENV=production
AWS_SWAGGER_API_URL=http://YOUR_AWS_EC2_PUBLIC_IP:8082/swagger.json

# Shared secret (existing)
GHCR_PAT=your_github_personal_access_token
```

**Replace `YOUR_AWS_EC2_PUBLIC_IP` with your actual EC2 instance public IP address.**

**Note**: These AWS-prefixed secrets won't interfere with your existing Azure deployment secrets!

## Step 2: Configure AWS Security Group

Ensure your EC2 security group allows inbound traffic on these ports:

```bash
# Using AWS CLI
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxxxxxx \
  --protocol tcp --port 22 --cidr 0.0.0.0/0    # SSH

aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxxxxxx \
  --protocol tcp --port 80 --cidr 0.0.0.0/0    # HTTP

aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxxxxxx \
  --protocol tcp --port 443 --cidr 0.0.0.0/0   # HTTPS

aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxxxxxx \
  --protocol tcp --port 8082-8090 --cidr 0.0.0.0/0  # API services
```

Or use the AWS Console:
1. Go to EC2 → Security Groups
2. Select your security group
3. Add inbound rules for ports: 22, 80, 443, 8082-8090

## Step 3: Setup AWS Instance

### Option A: Automated Setup (Recommended)

1. Clone the repository locally:
```bash
git clone https://github.com/pd120424d/mountain-service.git
cd mountain-service
```

2. Configure environment:
```bash
cp .env.aws .env
# Edit .env with your AWS instance details
```

3. Run setup script:
```bash
# On Linux/Mac
./setup-aws-instance.sh

# On Windows (use Git Bash or WSL)
bash setup-aws-instance.sh
```

### Option B: Manual Setup

SSH to your EC2 instance and run:

```bash
# Update system
sudo yum update -y

# Install Docker
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create directories
mkdir -p ~/mountain-service-deployment
mkdir -p ~/mountain-service-images
mkdir -p ~/logs

# Log out and back in for Docker group membership
exit
```

## Step 4: Deploy Application

### Option A: GitHub Actions (Recommended)

1. Go to your GitHub repository
2. Click "Actions" tab
3. Select "Multi-Cloud Deployment" workflow
4. Click "Run workflow"
5. Choose "aws" as deployment target
6. Click "Run workflow"
7. Monitor the deployment progress

### Option B: Local Deployment

1. Configure local environment:
```bash
cp .env.aws .env
# Edit .env with your values
```

2. Run deployment:
```bash
# On Linux/Mac
./deploy-aws.sh

# On Windows (use Git Bash or WSL)
bash deploy-aws.sh
```

## Step 5: Verify Deployment

After deployment completes, test your application:

```bash
# Replace YOUR_IP with your EC2 public IP
curl http://YOUR_IP/health                    # Frontend health
curl http://YOUR_IP/api/v1/health            # API health
curl http://YOUR_IP:8082/api/v1/health       # Employee service
curl http://YOUR_IP:8083/api/v1/health       # Urgency service
curl http://YOUR_IP:8084/api/v1/health       # Activity service
curl http://YOUR_IP:8090/api/v1/health       # Version service
```

Open in browser:
- Frontend: `http://YOUR_IP`
- Employee Swagger: `http://YOUR_IP/employee-swagger/`
- Urgency Swagger: `http://YOUR_IP/urgency-swagger/`
- Activity Swagger: `http://YOUR_IP/activity-swagger/`

## Step 6: Monitor and Maintain

SSH to your instance to monitor:

```bash
ssh -i ~/.ssh/your-key.pem ec2-user@YOUR_IP

# Check service status
~/monitor-services.sh

# View logs
cd ~/mountain-service-deployment
docker-compose logs -f

# Create backup
~/backup-data.sh
```

## Switching Between Azure and AWS

To switch deployment targets, just update the GitHub secret:

```
DEPLOYMENT_TARGET=azure  # For Azure VM
DEPLOYMENT_TARGET=aws    # For AWS EC2
```

The same workflow and scripts work for both platforms!

## Quick Commands Reference

```bash
# Setup AWS instance
./setup-aws-instance.sh

# Full deployment
./deploy-aws.sh

# Quick deployment (specific services)
./quick-deploy-aws.sh --frontend --employee

# Monitor services
ssh -i ~/.ssh/key.pem ec2-user@YOUR_IP '~/monitor-services.sh'

# View logs
ssh -i ~/.ssh/key.pem ec2-user@YOUR_IP 'cd ~/mountain-service-deployment && docker-compose logs -f'

# Restart services
ssh -i ~/.ssh/key.pem ec2-user@YOUR_IP 'cd ~/mountain-service-deployment && docker-compose restart'
```

## Troubleshooting

### Common Issues

1. **SSH Connection Failed**
   - Check EC2 instance is running
   - Verify security group allows port 22
   - Ensure SSH key has correct permissions (600)

2. **Services Not Starting**
   - Check Docker logs: `docker-compose logs -f`
   - Verify environment variables in .env
   - Ensure sufficient disk space and memory

3. **Port Access Issues**
   - Verify security group rules
   - Check if services are binding to correct ports
   - Test with `curl` from the instance itself first

4. **GitHub Actions Failing**
   - Check all required secrets are set
   - Verify SSH key format (include headers/footers)
   - Ensure GHCR_PAT has correct permissions

### Getting Help

1. Check the comprehensive guide: `docs/MULTI-CLOUD-DEPLOYMENT.md`
2. Review GitHub secrets setup: `GITHUB-SECRETS-SETUP.md`
3. Monitor deployment logs in GitHub Actions
4. SSH to instance and run `~/monitor-services.sh`

## Cost Optimization

### AWS EC2 Instance Recommendations

- **Development**: t3.micro (1 vCPU, 1 GB RAM) - Free tier eligible
- **Staging**: t3.small (2 vCPU, 2 GB RAM)
- **Production**: t3.medium (2 vCPU, 4 GB RAM) or larger

### Storage

- Use GP3 EBS volumes for better performance/cost ratio
- Set up automated snapshots for backups
- Monitor disk usage with `df -h`

### Networking

- Use Elastic IP if you need a static IP address
- Consider Application Load Balancer for high availability
- Set up CloudWatch monitoring for metrics

## Next Steps

1. **SSL/HTTPS**: Configure SSL certificates for production
2. **Domain Name**: Set up a custom domain name
3. **Monitoring**: Implement comprehensive monitoring with CloudWatch
4. **Backup**: Set up automated database backups
5. **CI/CD**: Enhance the deployment pipeline with testing stages
6. **Scaling**: Consider auto-scaling groups for high availability

Your Mountain Rescue Service is now running on AWS! 🎉
