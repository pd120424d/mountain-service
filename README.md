# 🏔️ Mountain Service Project

A modular, microservice-based system designed for managing mountain operations—coordinating employee shifts, tracking activities, and responding to urgencies in real time. Built with Go and Angular, this project emphasizes scalability, clarity, and role-based workflows.

---

## Table of Contents

- [Features](#-features)
- [CI/CD Pipeline](#-cicd-pipeline)
- [Architecture](#-architecture)
- [Microservices](#-microservices)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
- [License](#-license)

---

## Features

- Employee shift scheduling with role constraints (Medic, Technical, Administrator)
- Activity tracking for mountain operations
- Urgency management with real-time escalation
- Role-based access and user authentication
- Angular-based frontend with localization support
- Admin panel: safe service restarts (Kubernetes deployment rollout restart)
- Observability: Grafana dashboards at /grafana (admin-only)
- Kubernetes deployment (Helm charts in api/charts/) with Secrets via GitHub Actions

---

## CI/CD Pipeline

This project uses GitHub Actions for:
- TypeScript model generation/validation
- Backend/frontend test coverage
- Kubernetes deployment with kubectl (see .github/workflows/k8s-deploy.yml)

See docs/CI-CD-MODEL-GENERATION.md for model generation details.

---

## Architecture

The system follows a microservice architecture:

```
                       +-------------------+
                       |   Angular Frontend |
                       +-------------------+
                                |
                                v
                     +---------------------+
                     |    API Gateway       |
                     +---------------------+
                                |
        -------------------------------------------------
        |                     |                      |
        v                     v                      v
+----------------+   +------------------+   +------------------+
|  Employee Svc  |   |  Activity Svc     |   |  Urgency Svc     |
+----------------+   +------------------+   +------------------+

```

Each service is independently deployable and communicates via HTTP/REST APIs.

---

## Microservices

| Service           | Description                                | Path                        |
|-------------------|--------------------------------------------|-----------------------------|
| Employee Service  | Manages employee profiles and shifts       | `api/employee/`             |
| Activity Service  | Tracks mountain-related tasks and logs     | `api/activity/`             |
| Urgency Service   | Handles alerts and critical notifications  | `api/urgency/`              |
| API Gateway       | Entry point that routes requests           | `api/api-gateway/`          |
| Frontend (Angular)| UI for employees/admins                    | `frontend/`                 |

---

## Tech Stack

- **Backend**: Go (Gin, GORM, Zap, Afero)
- **Frontend**: Angular, Typescript
- **Database**: PostgreSQL, Redis, Firebase
- **DevOps**: Docker, Makefiles, CI/CD ready
- **Tools**: GitHub Actions, Prometheus/Grafana (planned), Kubernetes

---

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/your-org/mountain-service.git
cd mountain-service
```

2. Kubernetes quick start:

```bash
kubectl apply -f k8s/namespaces.yaml
kubectl -n mountain apply -f k8s/deployments/
kubectl -n mountain apply -f k8s/services/
kubectl -n mountain apply -f k8s/frontend/
```

3. GitHub Secrets required: see docs/README.k8s-variables.md

4. For legacy Docker instructions, see deprecated/README.md

- [Backend README](./api/README.md)
- [Frontend README](./ui/README.md)

---

## CI/CD and Deployment

- Test coverage workflows are kept: backend-test-coverage.yml, frontend-test-coverage.yml
- Deployment to Kubernetes is handled by .github/workflows/k8s-deploy.yml
- Legacy AWS/docker-compose workflows and scripts are in deprecated/

Secrets required for K8s: see docs/README.k8s-variables.md

3. For more specific setup instructions, refer to individual README files:

- [Backend README](./api/README.md)

---

## Observability and Admin operations

- Grafana: Exposed at https://mountain-service.duckdns.org/grafana (link available in Admin panel). See docs/OBSERVABILITY-GRAFANA.md for Helm install and Ingress example.
- Admin service restarts: Admin panel includes a dropdown and confirmation modal to safely trigger a rollout restart of selected services.
  - Backend endpoint: POST /api/v1/admin/k8s/restart with body {"deployment":"<name>"}
  - Secured by admin JWT and Kubernetes RBAC; only allowlisted Deployments can be restarted.
  - Enabled via Helm in the employee-service chart:

    values.yaml
    - serviceAccount.create: true
    - rbac.create: true
    - rbac.allowedDeployments: [employee-service, urgency-service, activity-service, version-service, docs-aggregator, docs-ui]

  - Install/upgrade example:
    helm upgrade --install employee-service api/charts/employee-service -n mountain-service -f your-values.yaml

For details:
- docs/OBSERVABILITY-GRAFANA.md — Grafana setup and exposure at /grafana
- docs/ADMIN-RESTART.md — How the restart feature works, Helm values, and testing the endpoint

- [Frontend README](./ui/README.md)

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
