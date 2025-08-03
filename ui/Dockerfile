# Staging/Development Dockerfile for Mountain Service UI
# Simple configuration for local development and staging environments

# Stage 1: Build Angular App
FROM node:18 AS builder

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Build the Angular application
RUN npm run build --configuration=production

# Stage 2: Serve with NGINX
FROM nginx:alpine

# Install additional tools for health checks and debugging
RUN apk add --no-cache curl wget

# Copy staging/development NGINX config
COPY nginx.staging.conf /etc/nginx/nginx.conf

# Copy built Angular app to NGINX
COPY --from=builder /app/dist/ui/browser /usr/share/nginx/html

# Create log directories and ensure proper permissions
RUN mkdir -p /var/log/nginx /var/log && \
    touch /var/log/frontend-access.log /var/log/frontend-error.log && \
    chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/log/nginx && \
    chown nginx:nginx /var/log/frontend-*.log && \
    chmod 644 /var/log/frontend-*.log

EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/health || exit 1

CMD ["nginx", "-g", "daemon off;"]
