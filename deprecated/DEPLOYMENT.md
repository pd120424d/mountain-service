# Deployment Guide (Deprecated)

This document reflects the legacy Docker Compose based deployment process. It is kept for historical reference while the project migrates to Kubernetes.

## Configuration Files (Legacy)
- docker-compose.yml
- docker-compose.aws.yml
- docker-compose.frontend.yml
- .env.* files consumed by the above

## Scripts (Legacy)
- deploy-backend.sh
- deploy-frontend.sh

## CI/CD (Legacy)
- GitHub Actions workflows in `.github/workflows` that target docker compose and EC2/VM deployment.

For the current Kubernetes-based deployment, see `k8s/README.md` and the repo root `README.md`.

