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

# Optional: Custom NGINX config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built Angular app to NGINX
COPY --from=builder /app/dist/ui/browser /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
