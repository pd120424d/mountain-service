# Kubernetes Deployment Guide

This project deploys the following components to Kubernetes:
- employee-service (Go)
- urgency-service (Go)
- activity-service (Go)
- version-service (Go)
- frontend (Angular + Nginx)
- Optional: Cloud SQL Auth Proxy sidecars for DB connectivity
- Optional: Swagger UIs exposed via frontend Nginx reverse proxy

## Prerequisites
- A Kubernetes cluster (k3s recommended for your VM)
- kubectl access configured
- Cloud SQL instances (PostgreSQL) for employee, urgency, activity
- Firestore (Native mode) if using Pub/Sub sync (step 5)
- GitHub Actions configured with secrets (see below)

## GitHub Secrets to create
Per service database (replace with your values):
- EMPLOYEE_DB_HOST, EMPLOYEE_DB_PORT, EMPLOYEE_DB_USER, EMPLOYEE_DB_PASSWORD, EMPLOYEE_DB_NAME, EMPLOYEE_DB_SSLMODE
- URGENCY_DB_HOST, URGENCY_DB_PORT, URGENCY_DB_USER, URGENCY_DB_PASSWORD, URGENCY_DB_NAME, URGENCY_DB_SSLMODE
- ACTIVITY_DB_HOST, ACTIVITY_DB_PORT, ACTIVITY_DB_USER, ACTIVITY_DB_PASSWORD, ACTIVITY_DB_NAME, ACTIVITY_DB_SSLMODE

Shared app secrets:
- JWT_SECRET
- ADMIN_PASSWORD
- SERVICE_AUTH_SECRET
- CORS_ALLOWED_ORIGINS

Cloud SQL / GCP (if using Cloud SQL Proxy or Pub/Sub):
- GCP_PROJECT_ID
- GCP_SA_KEY (base64-encoded service account key JSON)
- (Optional) CLOUDSQL_INSTANCE_EMPLOYEE, CLOUDSQL_INSTANCE_URGENCY, CLOUDSQL_INSTANCE_ACTIVITY (instance connection names)

## Workflow
1. GitHub Actions deploy job uses kubectl to create/update Kubernetes Secrets from the above repo secrets.
2. Apply manifests in k8s/ to deploy/update services.
3. Each service reads DB_ env vars from its Secret and auto-migrates on startup.

## Manifests
- k8s/namespaces.yaml
- k8s/secrets/*.yaml (templates; CI fills in actual values)
- k8s/deployments/*.yaml
- k8s/services/*.yaml
- k8s/frontend/*.yaml

## Swagger
Backend services expose Swagger at:
- /swagger/ (UI)
- /swagger.json (spec)
The frontend nginx config proxies:
- /employee-swagger/ -> employee-service/swagger/
- /employee-swagger.json -> employee-service/swagger.json
- /urgency-swagger/ -> urgency-service/swagger/
- /urgency-swagger.json -> urgency-service/swagger.json
- /activity-swagger/ -> activity-service/swagger/
- /activity-swagger.json -> activity-service/swagger.json

Ensure images include /docs/swagger.json (Dockerfiles already copy docs/ to /docs).

## Apply
- kubectl apply -f k8s/namespaces.yaml
- kubectl apply -f k8s/secrets/ (after CI populates or you create them)
- kubectl apply -f k8s/deployments/
- kubectl apply -f k8s/services/
- kubectl apply -f k8s/frontend/

See k8s/ci-deploy-example.yml for a GH Actions example.

