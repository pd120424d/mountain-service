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
- Dockerized services with future Kubernetes compatibility

---

## CI/CD Pipeline

This project features an automated CI/CD pipeline with contract-driven model generation:

- **Automated Model Generation**: TypeScript models are automatically generated from backend Swagger specs
- **Smart Fallback**: Uses live APIs when available, falls back to local swagger files
- **Continuous Validation**: Models are validated before testing and deployment
- **Independent Deployments**: Separate AWS frontend and backend deployments for faster, safer releases
- **Path-Based Triggering**: Only deploys components that have changed

**Pipeline Flow**:
```
Code Push → Generate Models → Run Tests → Deploy (Frontend/Backend)
     ↓           ↓              ↓         ↓
   GitHub    Live APIs or   Coverage   AWS EC2
   Actions   Fallback Files  Reports   Docker
```

 **[View Detailed CI/CD Documentation](docs/CI-CD-MODEL-GENERATION.md)**

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
- **Frontend**: Angular
- **Database**: PostgreSQL
- **DevOps**: Docker, Makefiles, CI/CD ready
- **Tools**: GitHub Actions, Prometheus/Grafana (planned), Kubernetes (planned)

---

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/your-org/mountain-service.git
cd mountain-service
```

2. Start the backend and frontend locally using Docker:

```bash
docker-compose up --build
```

3. For more specific setup instructions, refer to individual README files:

- [Backend README](./api/README.md)
- [Frontend README](./ui/README.md)

---

## CI/CD and Deployment

The project uses **GitHub Actions** for CI/CD with separate AWS frontend and backend deployments.

### Test Coverage Workflows

- `frontend-test-coverage.yml` and `backend-test-coverage.yml` run on:
  - pushes to `main` or version tags (e.g. `v1.0.0`)
  - relevant changes in `ui/`, `api/`, or test script files
  - all pull requests (for coverage only, not deployment)

### AWS Deploy Workflows

- `aws-frontend-deploy.yml` and `aws-backend-deploy.yml` run **only after successful test coverage** (via `workflow_run`)
- **Path-based triggering**: Frontend deploys only on `ui/` changes, backend only on `api/` changes
- Deployments are **only triggered by push events** (e.g. merging into `main` or creating a tag)
- **Pull requests will never trigger a deployment**

### Deployment Architecture

- **Independent Deployments**: Frontend and backend deploy separately for faster, safer deployments
- Docker images are built and pushed to **GitHub Container Registry**
- Deployment runs over **SSH to AWS EC2**, loads Docker images via separate docker-compose files
- **Frontend**: Nginx reverse proxy with Angular app (port 80/443)
- **Backend**: All microservices + PostgreSQL database (ports 8082-8090)

For detailed deployment information, see [DEPLOYMENT.md](./DEPLOYMENT.md).

### GitHub Secrets

To support AWS deployments, these GitHub Secrets are required:

- `GHCR_PAT`: GitHub token for publishing to the container registry
- `AWS_SSH_PRIVATE_KEY`: SSH private key for EC2 access
- `AWS_INSTANCE_USER`: EC2 instance username (typically `ubuntu`)
- `AWS_INSTANCE_IP`: EC2 instance IP address
- `AWS_JWT_SECRET`: JWT secret for backend authentication
- `AWS_ADMIN_PASSWORD`: Admin password for backend services
- `AWS_SERVICE_AUTH_SECRET`: Service-to-service authentication secret
- `AWS_CORS_ALLOWED_ORIGINS`: Allowed CORS origins
- `AWS_DB_USER`: PostgreSQL database username
- `AWS_DB_PASSWORD`: PostgreSQL database password
- `AWS_AZURE_STORAGE_ACCOUNT_NAME`: Azure storage account name
- `AWS_AZURE_STORAGE_ACCOUNT_KEY`: Azure storage account key
- `AWS_AZURE_STORAGE_CONTAINER_NAME`: Azure storage container name

See [DEPLOYMENT.md](./DEPLOYMENT.md) for complete setup instructions.
- `SWAGGER_API_URL`: Swagger documentation URL
- `CORS_ALLOWED_ORIGINS`: Frontend URL for CORS config


#### Mimicking Secrets Locally

Create a `.env` file in the root directory with values like:

```env
GHCR_PAT=your_pat_here
AZURE_SSH_PRIVATE_KEY=path/to/private_key
AZURE_VM_USER=azureuser
AZURE_VM_HOST=your.vm.host
SWAGGER_API_URL=http://yourhost:port/swagger/index.html
CORS_ALLOWED_ORIGINS=http://localhost:4200
```

#### Running Services Locally

To spin up the services locally with production-like environment variables, you don't need GitHub Secrets.
Instead, you can use an `.env` file and the `docker-compose.staging.yml` file:

```bash
docker compose --env-file .env -f docker-compose.staging.yml up --build
```

Also, it is important to mention that api/employee/secrets/db_user and api/employee/secrets/db_password are required to be set in order to initialize and connect to the database. 

Example `.env` file:

```env
VERSION_SERVICE_IMAGE=ghcr.io/your-org/version-service:v1.2.3
SWAGGER_API_URL=http://localhost:9082/swagger.json
CORS_ALLOWED_ORIGINS=http://localhost:4200
```

This allows local simulation of the production environment variables used during CI/CD deployment.

1. Clone the repository:

```bash
git clone https://github.com/your-org/mountain-service.git
cd mountain-service
```

2. Start the backend and frontend using Docker:

```bash
docker-compose up --build
```

3. For more specific setup instructions, refer to individual README files:

- [Backend README](./api/README.md)
- [Frontend README](./ui/README.md)

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
