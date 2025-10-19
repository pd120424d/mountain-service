# Mountain Service Project

A modular, microservice-based system designed for managing mountain rescue operations—coordinating employee shifts, tracking activities, and responding to urgencies in real time. Built with Go and Angular, this project emphasizes scalability, eventual consistency, and production-ready patterns.

**Deployed at:** [http://mountain-service.duckdns.org/](http://mountain-service.duckdns.org/)

---

## Table of Contents

- [Features](#features)
- [Microservices](#microservices)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Documentation](#documentation)
- [License](#license)
- [Credits](#credits)

---

## Features

- **Employee Management**: Shift scheduling with role constraints (Medic, Technical, Administrator)
- **Activity Tracking**: Real-time activity logs with infinite scroll and cursor-based pagination
- **Urgency Management**: Critical incident tracking with real-time escalation
- **CQRS Architecture**: Separate read/write models for optimal performance
- **Role-based Access**: JWT authentication with service-to-service security
- **Localization**: Multi-language support (Serbian Cyrillic, Latin, English, Russian)
- **Admin Panel**: Safe service restarts via Kubernetes rollout
- **Observability**: Grafana dashboards with GCP Cloud Logging integration

---


## Microservices

| Service                    | Description                                      | Path                              | Tech Stack                    |
|----------------------------|--------------------------------------------------|-----------------------------------|-------------------------------|
| **Employee Service**       | Employee profiles, shifts, authentication        | `api/employee/`                   | Go, Gin, GORM, PostgreSQL     |
| **Activity Service**       | Activity tracking with CQRS read-model           | `api/activity/`                   | Go, Gin, GORM, Firestore      |
| **Urgency Service**        | Critical incident management                     | `api/urgency/`                    | Go, Gin, GORM, PostgreSQL     |
| **Activity Read-Model**    | Firestore sync via Pub/Sub events                | `api/activity-readmodel-updater/` | Go, Pub/Sub, Firestore        |
| **Docs Aggregator**        | Centralized Swagger/OpenAPI documentation        | `api/docs-aggregator/`            | Go, Swagger UI                |
| **Version Service**        | Build version and health checks                  | `api/version-service/`            | Go, Gin                       |
| **Frontend**               | Angular SPA with localization                    | `ui/`                             | Angular 18, TypeScript, Nginx |

---

## Tech Stack

### Backend
- **Language**: Go 1.22+
- **Framework**: Gin (HTTP), GORM (ORM)
- **Databases**:
  - PostgreSQL (write model, with read replicas)
  - Google Cloud Firestore (read model)
  - Redis (token blacklist)
- **Messaging**: Google Cloud Pub/Sub
- **Logging**: Uber Zap + GCP Cloud Logging
- **Testing**: gomock, go-sqlmock, testify

### Frontend
- **Framework**: Angular 18
- **Language**: TypeScript 5.5
- **UI Components**: Angular Material
- **Localization**: ngx-translate (Serbian Cyrillic/Latin, English, Russian)
- **Maps**: Leaflet with OpenStreetMap
- **Testing**: Karma, Jasmine

### Infrastructure
- **Container Orchestration**: Kubernetes (GKE)
- **Package Manager**: Helm 3
- **Ingress**: Traefik with TLS
- **CI/CD**: GitHub Actions
- **Monitoring**: Grafana + GCP Cloud Logging
- **Secrets Management**: Kubernetes Secrets via GitHub Actions

### DevOps Tools
- **Code Generation**: swagger-typescript-api, gomock, swag
- **Database Migrations**: GORM AutoMigrate
- **Containerization**: Docker multi-stage builds
- **Load Balancing**: PgBouncer (database connection pooling)

---

## Getting Started

### Prerequisites

- **Go** 1.22+ ([install](https://go.dev/doc/install))
- **Node.js** 20+ and npm ([install](https://nodejs.org/))
- **Docker** and Docker Compose ([install](https://docs.docker.com/get-docker/))
- **kubectl** ([install](https://kubernetes.io/docs/tasks/tools/))
- **Helm** 3+ ([install](https://helm.sh/docs/intro/install/))

### Local Development Setup

1. **Clone the repository:**

```bash
git clone https://github.com/pd120424d/mountain-service.git
cd mountain-service
```

2. **Set up environment variables:**

```bash
# Backend services
cd api
cp .env.example .env
# Edit .env with your database credentials, JWT secrets, etc.
```
Note: Go module root is in api/. Run Go commands from the api/ directory.


3. **Run backend services:**

```bash
# Employee service
cd api/employee
go run cmd/main.go

# Activity service
cd api/activity
go run cmd/main.go

# Urgency service
cd api/urgency
go run cmd/main.go
```

4. **Run frontend:**

```bash
cd ui
npm install
npm start
# Open http://localhost:4200
```

### Running Tests

**Backend tests with coverage:**

```bash
cd api
./backend-test-cover.sh
# Opens coverage report in browser
```

**Frontend tests with coverage:**

```bash
cd ui
npm test
# Or with coverage:
./frontend-test-cover.sh
```

### Building for Production

**Backend Docker images:**

```bash
# Build all services
cd api
docker build -t employee-service:latest -f employee/Dockerfile .
docker build -t activity-service:latest -f activity/Dockerfile .
docker build -t urgency-service:latest -f urgency/Dockerfile .
```

**Frontend Docker image:**

```bash
cd ui
docker build -t frontend:latest .
```

### Kubernetes Deployment

**Deploy with Helm:**

```bash
# Install all services
cd api/charts

helm install employee-service ./employee-service
helm install activity-service ./activity-service
helm install urgency-service ./urgency-service
helm install frontend ./frontend

# Or use the umbrella chart (if available)
helm install mountain-service ./mountain-service
```

**Check deployment status:**

```bash
kubectl get pods
kubectl get ingress
```

---


---


## Documentation

For detailed documentation on the project, including API references, deployment guides, and more, please refer to the [docs](docs/) directory.

---


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Credits

- **Google Cloud Platform** for Firestore, Pub/Sub, and Cloud Logging
- **Kubernetes** and **Helm** communities for excellent orchestration tools
- **Angular** and **Go** communities for robust frameworks
- **OpenStreetMap** for map data

---