# Swagger in Kubernetes

Each backend service exposes:
- Swagger UI at /swagger/
- Swagger JSON at /swagger.json

The frontend Nginx proxies these:
- /employee-swagger/ → http://employee-service:8082/swagger/
- /employee-swagger.json → http://employee-service:8082/swagger.json
- /urgency-swagger/ → http://urgency-service:8083/swagger/
- /urgency-swagger.json → http://urgency-service:8083/swagger.json
- /activity-swagger/ → http://activity-service:8084/swagger/
- /activity-swagger.json → http://activity-service:8084/swagger.json

Make sure:
- Dockerfiles copy /docs/swagger.json (already configured)
- Services are resolvable in-cluster by DNS (Service names match above)
- Frontend Nginx config (ui/nginx.production.conf and staging) are baked into the frontend image

Troubleshooting:
- Check /docs/swagger.json exists inside the service container
- Ensure /swagger.json route is reachable (curl from another pod)
- Confirm frontend is pointing to the correct service names

