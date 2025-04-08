# 🏔️ Mountain Service Project

A modular, microservice-based system designed for managing mountain operations—coordinating employee shifts, tracking activities, and responding to urgencies in real time. Built with Go and Angular, this project emphasizes scalability, clarity, and role-based workflows.

---

## 📚 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Microservices](#-microservices)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
- [License](#-license)

---

## ✨ Features

- Employee shift scheduling with role constraints (Medic, Technical, Administrator)
- Activity tracking for mountain operations
- Urgency management with real-time escalation
- Role-based access and user authentication
- Angular-based frontend with localization support
- Dockerized services with future Kubernetes compatibility

---

## 🏗️ Architecture

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

## 🧩 Microservices

| Service           | Description                                | Path                        |
|-------------------|--------------------------------------------|-----------------------------|
| Employee Service  | Manages employee profiles and shifts       | `api/employee/`             |
| Activity Service  | Tracks mountain-related tasks and logs     | `api/activity/`             |
| Urgency Service   | Handles alerts and critical notifications  | `api/urgency/`              |
| API Gateway       | Entry point that routes requests           | `api/api-gateway/`          |
| Frontend (Angular)| UI for employees/admins                    | `frontend/`                 |

---

## 🛠 Tech Stack

- **Backend**: Go (Gin, GORM, Zap, Afero)
- **Frontend**: Angular
- **Database**: PostgreSQL
- **DevOps**: Docker, Makefiles, CI/CD ready
- **Tools**: GitHub Actions, Prometheus/Grafana (planned), Kubernetes (planned)

---

## 🚀 Getting Started

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

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
