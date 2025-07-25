# AWS Deployment Secrets Setup

This document provides the exact GitHub secrets you need to add for AWS deployment **without affecting your existing Azure deployment**.

## Key Benefits

✅ **No Conflicts**: AWS secrets use `AWS_` prefix - won't interfere with Azure  
✅ **Keep Azure Working**: Your existing Azure deployment continues unchanged  
✅ **Easy Switching**: Deploy to either platform via GitHub Actions  
✅ **Separate Configs**: Different passwords, keys, and settings for each platform  

## 📋 Required GitHub Secrets

Go to your GitHub repository → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

### 1. AWS Instance Configuration

```
Name: AWS_DEPLOYMENT_TARGET
Value: aws
```

```
Name: AWS_INSTANCE_IP
Value: YOUR_AWS_EC2_PUBLIC_IP
Description: Replace with your actual AWS EC2 public IP address
```

```
Name: AWS_INSTANCE_USER
Value: ec2-user
Description: SSH username for AWS EC2 (use 'ubuntu' if using Ubuntu AMI)
```

```
Name: AWS_SSH_PRIVATE_KEY
Value: -----BEGIN OPENSSH PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
...your complete AWS SSH private key content...
...
-----END OPENSSH PRIVATE KEY-----
Description: Your AWS EC2 SSH private key content (not file path!)
```

### 2. AWS Application Security

```
Name: AWS_JWT_SECRET
Value: aws-super-secret-jwt-key-make-it-different-from-azure-64-chars-long
Description: JWT secret for AWS (should be different from Azure)
```

```
Name: AWS_ADMIN_PASSWORD
Value: AWSSecureAdminPass123!
Description: Admin password for AWS deployment (can be different from Azure)
```

```
Name: AWS_SERVICE_AUTH_SECRET
Value: aws-service-auth-secret-also-different-from-azure-32-chars
Description: Inter-service auth secret for AWS
```

### 3. AWS Network Configuration

```
Name: AWS_CORS_ALLOWED_ORIGINS
Value: http://YOUR_AWS_EC2_PUBLIC_IP,https://YOUR_AWS_EC2_PUBLIC_IP
Description: Replace YOUR_AWS_EC2_PUBLIC_IP with actual IP
```

**Note**: You don't need separate `AWS_EMPLOYEE_SERVICE_URL` and `AWS_ACTIVITY_SERVICE_URL` secrets since these are internal Docker network URLs that are the same for both platforms.

## 🔄 Your Existing Azure Secrets (Keep Unchanged!)

**Don't modify these** - they're used by your current Azure deployment:

```
AZURE_VM_HOST=your.azure.vm.ip (keep as-is)
AZURE_VM_USER=azureuser (keep as-is)
AZURE_SSH_PRIVATE_KEY=your_azure_ssh_key (keep as-is)
COMPOSE_ENV=your_existing_compose_env (keep as-is)
SWAGGER_API_URL=your_existing_swagger_url (keep as-is)
GHCR_PAT=your_github_token (shared between both platforms)
```

**Note**: Based on your existing secrets, you don't have separate JWT_SECRET, ADMIN_PASSWORD, etc. - these might be handled differently in your current setup.

## 📝 Complete Example

Here's what your GitHub secrets should look like after adding AWS support:

### New AWS Secrets (add these 8 secrets):
```
AWS_DEPLOYMENT_TARGET=aws
AWS_INSTANCE_IP=54.123.45.67
AWS_INSTANCE_USER=ec2-user
AWS_SSH_PRIVATE_KEY=-----BEGIN OPENSSH PRIVATE KEY-----
MIIEpAIBAAKCAQEA1234567890abcdef...
...complete AWS key content...
-----END OPENSSH PRIVATE KEY-----
AWS_JWT_SECRET=aws-jwt-secret-make-it-long-and-random-64-characters-minimum
AWS_ADMIN_PASSWORD=AWSSecurePass123!
AWS_SERVICE_AUTH_SECRET=aws-service-secret-32-chars-min
AWS_CORS_ALLOWED_ORIGINS=http://54.123.45.67,https://54.123.45.67
```

**Note**: `COMPOSE_ENV` and `SWAGGER_API_URL` are automatically generated from your AWS_INSTANCE_IP, so you don't need separate secrets for them.

### Existing Secrets (don't change):
```
AZURE_VM_HOST=20.123.45.67
AZURE_VM_USER=azureuser
AZURE_SSH_PRIVATE_KEY=-----BEGIN OPENSSH PRIVATE KEY-----...
COMPOSE_ENV=production
SWAGGER_API_URL=http://20.123.45.67:8082/swagger.json
GHCR_PAT=ghp_1234567890abcdef1234567890abcdef12345678
```

## 🚀 How to Deploy

### Deploy to AWS:
1. Go to **Actions** tab in your GitHub repository
2. Select **Multi-Cloud Deployment** workflow
3. Click **Run workflow**
4. Choose **aws** as deployment target
5. Click **Run workflow**

### Deploy to Azure (unchanged):
1. Go to **Actions** tab in your GitHub repository
2. Select **Multi-Cloud Deployment** workflow
3. Click **Run workflow**
4. Choose **azure** as deployment target
5. Click **Run workflow**

## 🔧 Generate Strong Secrets

Use these commands to generate secure secrets:

```bash
# Generate AWS JWT secret (64 characters)
openssl rand -base64 48

# Generate AWS service auth secret (32 characters)
openssl rand -base64 24

# Generate AWS admin password (16 characters)
openssl rand -base64 12
```

## ✅ Verification Checklist

Before deploying, ensure you have:

- [ ] Added all 10 AWS-prefixed secrets to GitHub
- [ ] Replaced `YOUR_AWS_EC2_PUBLIC_IP` with actual IP in 2 secrets
- [ ] Used your actual AWS SSH private key content (not file path)
- [ ] Generated strong, unique secrets (different from Azure)
- [ ] Kept all existing Azure secrets unchanged
- [ ] AWS EC2 instance is running and accessible
- [ ] AWS Security Group allows ports: 22, 80, 443, 8082-8090

## 🔍 Testing Your Setup

Test your secrets configuration:

1. **Manual Deployment Test**:
   - Go to Actions → Multi-Cloud Deployment → Run workflow
   - Choose "aws" target
   - Monitor workflow execution

2. **SSH Connection Test**:
   ```bash
   ssh -i ~/.ssh/your-aws-key.pem ec2-user@YOUR_AWS_IP
   ```

3. **Service Health Test** (after deployment):
   ```bash
   curl http://YOUR_AWS_IP/health
   curl http://YOUR_AWS_IP/api/v1/health
   ```

## 🆘 Troubleshooting

### Common Issues:

1. **SSH Key Format Error**
   - Ensure you copy the entire key including headers/footers
   - No extra spaces or line breaks

2. **Workflow Fails with "Secret not found"**
   - Check secret names match exactly (case-sensitive)
   - Ensure all AWS_ prefixed secrets are added

3. **Services Don't Start**
   - Check AWS Security Group allows required ports
   - Verify EC2 instance has sufficient resources (t3.small minimum)

4. **CORS Errors**
   - Ensure AWS_CORS_ALLOWED_ORIGINS has correct IP
   - Include both HTTP and HTTPS URLs

## 🎉 Success!

After successful deployment, your applications will be available at:

- **Azure** (unchanged): `http://your.azure.ip`
- **AWS** (new): `http://your.aws.ip`

Both deployments run independently with separate configurations!

## 📞 Need Help?

1. Check the comprehensive guide: `docs/MULTI-CLOUD-DEPLOYMENT.md`
2. Review the quick start: `QUICK-START-AWS.md`
3. Monitor GitHub Actions logs for detailed error messages
4. Test SSH connection to your AWS instance first
