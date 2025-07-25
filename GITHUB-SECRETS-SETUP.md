# GitHub Secrets Configuration for Multi-Cloud Deployment

This document provides the exact GitHub secrets you need to configure for deploying the Mountain Rescue Service to both Azure VMs and AWS EC2 instances, using separate secret sets for each platform.

## Approach: Separate Secrets for Each Platform

To avoid conflicts with your existing Azure deployment, we'll use:
- **Existing secrets** for Azure (no changes needed)
- **AWS_-prefixed secrets** for AWS deployment
- **Shared secrets** for common configuration

## Required GitHub Secrets

Go to your repository on GitHub → Settings → Secrets and variables → Actions → New repository secret

### 1. Container Registry Authentication (Shared)

```
Name: GHCR_PAT
Value: your_github_personal_access_token_with_packages_write_permission
Description: Used by both Azure and AWS deployments
```

### 2. AWS-Specific Secrets (New)

```
Name: AWS_DEPLOYMENT_TARGET
Value: aws
Description: Deployment target identifier for AWS
```

```
Name: AWS_INSTANCE_IP
Value: YOUR_AWS_EC2_PUBLIC_IP
Description: Public IP address of your AWS EC2 instance
```

```
Name: AWS_INSTANCE_USER
Value: ec2-user
Description: SSH username for AWS EC2 (typically ec2-user)
```

```
Name: AWS_SSH_PRIVATE_KEY
Value: -----BEGIN OPENSSH PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
...your complete AWS SSH private key content...
...
-----END OPENSSH PRIVATE KEY-----
Description: SSH private key for AWS EC2 access
```

### 3. AWS Application Security Configuration

```
Name: AWS_JWT_SECRET
Value: your-aws-specific-jwt-secret-make-it-different-from-azure
Description: JWT secret key for AWS deployment (can be different from Azure)
```

```
Name: AWS_ADMIN_PASSWORD
Value: your-aws-admin-password-123!
Description: Admin password for AWS deployment
```

```
Name: AWS_SERVICE_AUTH_SECRET
Value: your-aws-service-auth-secret-different-from-azure
Description: Service auth secret for AWS deployment
```

### 4. AWS CORS and API Configuration

```
Name: AWS_CORS_ALLOWED_ORIGINS
Value: http://YOUR_AWS_IP,https://YOUR_AWS_IP
Description: CORS origins for AWS deployment (replace YOUR_AWS_IP with actual AWS IP)
```

```
Name: AWS_EMPLOYEE_SERVICE_URL
Value: http://employee-service:8082
Description: Internal employee service URL for AWS (Docker network)
```

```
Name: AWS_ACTIVITY_SERVICE_URL
Value: http://activity-service:8084
Description: Internal activity service URL for AWS (Docker network)
```

### 5. Existing Azure Secrets (Keep Unchanged)

Your existing Azure secrets remain unchanged and continue to work:

```
AZURE_VM_HOST=your.azure.vm.ip (existing)
AZURE_VM_USER=azureuser (existing)
AZURE_SSH_PRIVATE_KEY=your_azure_ssh_key (existing)
CORS_ALLOWED_ORIGINS=http://your.azure.ip,https://your.azure.ip (existing)
JWT_SECRET=your_existing_jwt_secret (existing)
ADMIN_PASSWORD=your_existing_admin_password (existing)
SERVICE_AUTH_SECRET=your_existing_service_auth_secret (existing)
EMPLOYEE_SERVICE_URL=http://employee-service:8082 (existing)
ACTIVITY_SERVICE_URL=http://activity-service:8084 (existing)
```

**Note**: Don't change these - they're used by your current Azure deployment!

## Example Configuration for AWS Deployment

Here's a complete example of the AWS-prefixed secrets you need to add:

```
# AWS-specific secrets (add these new ones)
AWS_DEPLOYMENT_TARGET=aws
AWS_INSTANCE_IP=54.123.45.67
AWS_INSTANCE_USER=ec2-user
AWS_SSH_PRIVATE_KEY=-----BEGIN OPENSSH PRIVATE KEY-----
MIIEpAIBAAKCAQEA1234567890abcdef...
...complete AWS key content...
-----END OPENSSH PRIVATE KEY-----
AWS_JWT_SECRET=aws-specific-jwt-secret-different-from-azure
AWS_ADMIN_PASSWORD=AWSSecureAdminPass123!
AWS_SERVICE_AUTH_SECRET=aws-service-auth-secret-different-from-azure
AWS_CORS_ALLOWED_ORIGINS=http://54.123.45.67,https://54.123.45.67
AWS_EMPLOYEE_SERVICE_URL=http://employee-service:8082
AWS_ACTIVITY_SERVICE_URL=http://activity-service:8084

# Shared secret (existing)
GHCR_PAT=ghp_1234567890abcdef1234567890abcdef12345678
```

## Example Configuration for Azure Deployment

Here's a complete example for Azure deployment:

```
DEPLOYMENT_TARGET=azure
CLOUD_INSTANCE_IP=20.123.45.67
CLOUD_INSTANCE_USER=azureuser
CLOUD_SSH_PRIVATE_KEY=-----BEGIN OPENSSH PRIVATE KEY-----
MIIEpAIBAAKCAQEA9876543210fedcba...
...complete key content...
-----END OPENSSH PRIVATE KEY-----
JWT_SECRET=super-secret-jwt-key-for-production-use-random-string-here
ADMIN_PASSWORD=SecureAdminPass123!
SERVICE_AUTH_SECRET=service-auth-secret-also-random-string
CORS_ALLOWED_ORIGINS=http://20.123.45.67,https://20.123.45.67
EMPLOYEE_SERVICE_URL=http://employee-service:8082
ACTIVITY_SERVICE_URL=http://activity-service:8084
GHCR_PAT=ghp_1234567890abcdef1234567890abcdef12345678
```

## How to Generate Required Values

### 1. GitHub Personal Access Token (GHCR_PAT)

1. Go to GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Select scopes: `write:packages`, `read:packages`, `delete:packages`
4. Copy the generated token

### 2. SSH Private Key

For AWS:
```bash
# If you have the .pem file from AWS
cat ~/.ssh/your-aws-key.pem
```

For Azure:
```bash
# If you generated SSH keys for Azure
cat ~/.ssh/id_rsa
```

### 3. Strong Secrets

Generate random secrets:
```bash
# Generate JWT secret
openssl rand -base64 64

# Generate service auth secret
openssl rand -base64 32

# Generate admin password
openssl rand -base64 16
```

## Security Best Practices

1. **Use Strong Secrets**: Generate random, long secrets for JWT and service authentication
2. **Rotate Secrets**: Regularly update secrets, especially if compromised
3. **Limit Access**: Only give repository access to users who need it
4. **Monitor Usage**: Check GitHub Actions logs for any unauthorized access attempts
5. **Use Environment-Specific Secrets**: Consider separate secrets for staging/production

## Troubleshooting

### Common Issues

1. **Invalid SSH Key Format**
   - Ensure you copy the entire key including headers and footers
   - Check for extra spaces or line breaks

2. **CORS Issues**
   - Make sure CORS_ALLOWED_ORIGINS includes both HTTP and HTTPS URLs
   - Include the correct IP address

3. **Authentication Failures**
   - Verify GHCR_PAT has correct permissions
   - Check if token has expired

4. **Deployment Target Issues**
   - Ensure DEPLOYMENT_TARGET is exactly 'aws' or 'azure'
   - Check that corresponding instance variables are set

### Testing Secrets

You can test your secrets configuration by running the GitHub Actions workflow manually:

1. Go to Actions tab in your repository
2. Select "Multi-Cloud Deployment" workflow
3. Click "Run workflow"
4. Choose your deployment target
5. Monitor the workflow execution

## Migration from Azure-Only to Multi-Cloud

If you're migrating from an Azure-only setup:

1. Keep your existing Azure secrets for backward compatibility
2. Add the new multi-cloud secrets
3. Set `DEPLOYMENT_TARGET=azure` initially
4. Test the new workflow with Azure
5. When ready for AWS, update `DEPLOYMENT_TARGET=aws` and AWS-specific secrets
6. Test AWS deployment

This approach ensures zero downtime during migration.
