# AWS Security Group Configuration for Mountain Service

## Architecture Overview

The Mountain Service uses a **reverse proxy architecture** where:
- ✅ **Frontend (nginx)** is the only service that should be externally accessible (port 80/443)
- ✅ **All API requests** are routed through nginx to backend services
- ✅ **Backend services** (ports 8082, 8083, 8084, 8090) should NOT be directly accessible externally
- ✅ **Security by design** - only the reverse proxy is exposed

## Current Status

✅ **Services are running correctly** - All containers are healthy and respond to local health checks
✅ **Docker containers are properly configured** - Services communicate internally via Docker network
✅ **Nginx reverse proxy** routes all API requests to appropriate backend services
✅ **Proper security** - Backend services are not directly exposed

## Required AWS Security Group Configuration

### Minimal Security Group Setup (Recommended)

**You only need to expose port 80 (and optionally 443 for HTTPS):**

1. **Go to AWS EC2 Console**
   - Navigate to https://console.aws.amazon.com/ec2/
   - Select your region (e.g., us-east-1)

2. **Find your EC2 instance**
   - Go to "Instances" in the left sidebar
   - Find your mountain-service instance

3. **Access Security Groups**
   - Click on your instance
   - Go to the "Security" tab
   - Click on the Security Group name (e.g., "launch-wizard-1" or similar)

4. **Required Inbound Rules**
   - Click "Edit inbound rules"
   - Ensure you have these rules (and ONLY these for web traffic):

   | Type | Protocol | Port Range | Source | Description |
   |------|----------|------------|--------|-------------|
   | SSH | TCP | 22 | Your-IP/32 | SSH access for deployment |
   | HTTP | TCP | 80 | 0.0.0.0/0 | Frontend and all APIs via nginx |
   | HTTPS | TCP | 443 | 0.0.0.0/0 | Frontend SSL (optional) |

5. **Save Rules**
   - Click "Save rules"

**Important**: Do NOT expose ports 8082, 8083, 8084, or 8090 - they should only be accessible internally!

### Using AWS CLI (Alternative)

```bash
# Get your security group ID (replace with your instance ID)
INSTANCE_ID="i-1234567890abcdef0"
SECURITY_GROUP_ID=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID --query 'Reservations[0].Instances[0].SecurityGroups[0].GroupId' --output text)

# Ensure HTTP access is allowed (should already exist)
aws ec2 authorize-security-group-ingress \
    --group-id $SECURITY_GROUP_ID \
    --protocol tcp \
    --port 80 \
    --cidr 0.0.0.0/0 \
    --description "HTTP access for frontend and APIs"

# Optional: Add HTTPS access
aws ec2 authorize-security-group-ingress \
    --group-id $SECURITY_GROUP_ID \
    --protocol tcp \
    --port 443 \
    --cidr 0.0.0.0/0 \
    --description "HTTPS access for frontend and APIs"
```

## Verification

Test the application through the proper nginx reverse proxy:

```bash
# Replace YOUR_INSTANCE_IP with your actual EC2 public IP

# Test frontend health
curl -f http://YOUR_INSTANCE_IP/health

# Test API health (routed through nginx to backend services)
curl -f http://YOUR_INSTANCE_IP/api/v1/health

# Test version endpoint
curl -f http://YOUR_INSTANCE_IP/api/v1/version

# Test employee API
curl -f http://YOUR_INSTANCE_IP/api/v1/employees

# Test urgency API
curl -f http://YOUR_INSTANCE_IP/api/v1/urgencies

# Test activity API
curl -f http://YOUR_INSTANCE_IP/api/v1/activities
```

**These should NOT work (and that's correct for security):**
```bash
# These should be blocked by security group (good!)
curl -f http://YOUR_INSTANCE_IP:8082/api/v1/health  # Should fail
curl -f http://YOUR_INSTANCE_IP:8083/api/v1/health  # Should fail
curl -f http://YOUR_INSTANCE_IP:8084/api/v1/health  # Should fail
curl -f http://YOUR_INSTANCE_IP:8090/api/v1/health  # Should fail
```

## How the Architecture Works

```
Internet → AWS Security Group (Port 80) → EC2 Instance → Nginx → Backend Services
                                                          ↓
                                                    Docker Network
                                                          ↓
                                              employee-service:8082
                                              urgency-service:8083
                                              activity-service:8084
                                              version-service:8090
```

**Benefits of this architecture:**
- ✅ **Single point of entry** - Only nginx is exposed
- ✅ **Better security** - Backend services are not directly accessible
- ✅ **Load balancing** - Nginx can distribute requests
- ✅ **SSL termination** - SSL/TLS handled at nginx level
- ✅ **Rate limiting** - Nginx provides built-in rate limiting
- ✅ **Logging** - Centralized access logs

## Next Steps

1. **Ensure your security group only exposes port 80 (and optionally 443)**
2. **Re-run the deployment** - it should now work correctly
3. **Test the application** by accessing http://YOUR_INSTANCE_IP
4. **All API calls** should go through http://YOUR_INSTANCE_IP/api/v1/...

## Security Best Practices

- ✅ **Never expose backend service ports directly** (8082, 8083, 8084, 8090)
- ✅ **Use nginx as reverse proxy** for all external access
- ✅ **Enable HTTPS** in production with proper SSL certificates
- ✅ **Use rate limiting** to prevent abuse
- ✅ **Monitor access logs** for security issues
