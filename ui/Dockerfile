# Stage 1: Build Angular App
FROM node:18 AS builder

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Generate TypeScript models from swagger files (fallback method for Docker build)
# This ensures models are available even if backend services aren't running during build
RUN echo "Generating TypeScript models..." && \
    npm run generate-employee-models-local && \
    npm run generate-urgency-models-local && \
    npm run generate-activity-models-local && \
    echo "Models generated successfully!"

# Validate generated models
RUN echo "Validating generated models..." && \
    npm run validate-models && \
    echo "Model validation passed!"

# Build the Angular application
RUN npm run build --configuration=production

# Stage 2: Serve with NGINX
FROM nginx:alpine

# Optional: Custom NGINX config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built Angular app to NGINX
COPY --from=builder /app/dist/ui/browser /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
