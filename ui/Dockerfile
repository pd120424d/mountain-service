# Dockerfile for Mountain Service UI
# Supports both staging and production environments via build args

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

# Build argument to choose nginx config (default: staging)
ARG NGINX_CONFIG=staging

# Install additional tools for health checks and debugging
RUN apk add --no-cache curl wget netcat-openbsd bind-tools

# Remove default nginx configs and copy our custom config
RUN rm -f /etc/nginx/conf.d/default.conf /etc/nginx/nginx.conf
COPY nginx.${NGINX_CONFIG}.conf /etc/nginx/nginx.conf

# Copy built Angular app to NGINX
COPY --from=builder /app/dist/ui/browser /usr/share/nginx/html

# Create directories for SSL certificates and logs
# Create both the directories nginx expects during startup AND our custom log files
RUN mkdir -p /etc/nginx/ssl /var/log/nginx /var/log && \
    touch /var/log/nginx/error.log /var/log/nginx/access.log && \
    touch /var/log/frontend-access.log /var/log/frontend-error.log && \
    chown -R nginx:nginx /usr/share/nginx/html && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /etc/nginx/ssl && \
    chown nginx:nginx /var/log/frontend-*.log && \
    chmod 644 /var/log/nginx/*.log /var/log/frontend-*.log

# Expose ports (443 for potential SSL support)
EXPOSE 80 443

# Create a startup script to ensure log directories exist
RUN echo '#!/bin/sh' > /start-nginx.sh && \
    echo 'mkdir -p /var/log/nginx /var/log' >> /start-nginx.sh && \
    echo 'touch /var/log/nginx/error.log /var/log/nginx/access.log' >> /start-nginx.sh && \
    echo 'touch /var/log/frontend-access.log /var/log/frontend-error.log' >> /start-nginx.sh && \
    echo 'chown -R nginx:nginx /var/log/nginx /var/log/frontend-*.log' >> /start-nginx.sh && \
    echo 'chmod 644 /var/log/nginx/*.log /var/log/frontend-*.log' >> /start-nginx.sh && \
    echo 'exec nginx -g "daemon off;"' >> /start-nginx.sh && \
    chmod +x /start-nginx.sh

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/health || exit 1

CMD ["/start-nginx.sh"]
