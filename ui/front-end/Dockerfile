# Stage 1: Build the Angular app
FROM node:18 AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build --prod

# Stage 2: Serve the app with NGINX
FROM nginx:alpine
COPY --from=build /app/dist/front-end /usr/share/nginx/html
EXPOSE 8090
CMD ["nginx", "-g", "daemon off;"]
