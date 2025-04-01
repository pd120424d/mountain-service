# Stage 1: Build Angular app
FROM node:18 AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy app source
COPY . .

# Build the app
RUN npm run build --prod

# Stage 2: Serve using Nginx
FROM nginx:alpine

# Copy custom nginx config (optional)
COPY nginx.conf /etc/nginx/nginx.conf

# Remove default content and add built app
RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/dist/ui /usr/share/nginx/html/

# Expose port 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
