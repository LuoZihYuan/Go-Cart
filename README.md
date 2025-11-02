# 🛒 Go-Cart 🛒

A RESTful e-commerce API for managing products, shopping carts, and payments, built in Go.

## 📋 Table of Contents

- [✅ Prerequisites](#-prerequisites)
- [🚀 Getting Started](#-getting-started)
  - [⚙️ Setup](#️-setup)
  - [💻 Development (Local)](#-development-local)
  - [🧪 Staging (AWS)](#-staging-aws)
  - [🏭 Production (AWS)](#-production-aws)
  - [🧹 Cleanup](#-cleanup)
- [📁 Project Structure](#-project-structure)

---

## ✅ Prerequisites

- **Docker Desktop** - For containerized development
- **Go 1.25.1+** - For local IDE support
- **AWS CLI** - For AWS deployments, configured with `aws configure`
- **Terraform 1.0+** - For infrastructure as code, run `terraform init` in `terraform/stage` and `terraform/prod`
- **VS Code or GoLand** (Optional) - For IDE features with autocomplete

---

## 🚀 Getting Started

### ⚙️ Setup

Downloads Go dependencies for IDE autocomplete and navigation features. Run this once after cloning the repository.

```bash
make setup
```

---

### 💻 Development (Local)

#### Deploy

Builds Docker image with Swagger enabled, starts container with Air hot-reload, and mounts source code as volume. When you edit any `.go` file and save, Air automatically detects changes, regenerates Swagger docs, recompiles the binary, and restarts the application (typically 2-5 seconds).

```bash
make deploy-dev
```

To test the API, open `http://localhost:8080/swagger/index.html` or use `cURL`:

```bash
# Add a product
curl -X POST http://localhost:8080/v1/products/12345/details \
  -H 'Content-Type: application/json' \
  -d '{
    "product_id": 12345,
    "sku": "ABC-123-XYZ",
    "manufacturer": "Acme Corporation",
    "category_id": 456,
    "weight": 1250,
    "some_other_id": 789
  }'

# Get a product
curl http://localhost:8080/v1/products/12345
```

#### Logs

Shows all container logs with real-time streaming. `Ctrl+C` to exit, container keeps running.

```bash
make log-dev
```

#### Shell

Opens interactive shell inside the running container for inspection and debugging. Type `exit` to close.

```bash
make shell-dev
```

#### Stop

Stops and removes the development container. Keeps Docker images for faster restart.

```bash
make stop-dev
```

#### Destroy

Removes containers, volumes, and images. Frees disk space.

```bash
make destroy-dev
```

---

### 🧪 Staging (AWS)

#### Deploy

First deployment creates ECR repository and infrastructure. Subsequent deployments build new image, push to ECR, and update ECS service. Waits for service stability and displays public IP.

```bash
make deploy-stage
```

To test the API, open `http://<STAGING_IP>:8080/swagger/index.html` or use `cURL`:

```bash
# Add a product
curl -X POST http://<STAGING_IP>:8080/v1/products/12345/details \
  -H 'Content-Type: application/json' \
  -d '{
    "product_id": 12345,
    "sku": "ABC-123-XYZ",
    "manufacturer": "Acme Corporation",
    "category_id": 456,
    "weight": 1250,
    "some_other_id": 789
  }'

# Get a product
curl http://<STAGING_IP>:8080/v1/products/12345
```

#### Logs

Streams CloudWatch logs in real-time. `Ctrl+C` to exit, service keeps running.

```bash
make log-stage
```

#### Shell

Opens interactive shell in the running ECS task on AWS. Type `exit` to close.

```bash
make shell-stage
```

#### Stop

Scales service to 0 tasks. Infrastructure remains, no compute costs.

```bash
make stop-stage
```

#### Start

Scales service back to 1 task.

```bash
make start-stage
```

#### Destroy

Deletes all AWS staging resources and removes local Docker images. Requires typing `yes` to confirm.

```bash
make destroy-stage
```

---

### 🏭 Production (AWS)

#### Deploy

Swagger disabled for security. Separate ECR repository. Waits for service stability and displays public IP.

```bash
make deploy-prod
```

To test the API, use `cURL` (Swagger is disabled in production):

```bash
# Add a product
curl -X POST http://<PRODUCTION_IP>:8080/v1/products/12345/details \
  -H 'Content-Type: application/json' \
  -d '{
    "product_id": 12345,
    "sku": "ABC-123-XYZ",
    "manufacturer": "Acme Corporation",
    "category_id": 456,
    "weight": 1250,
    "some_other_id": 789
  }'

# Get a product
curl http://<PRODUCTION_IP>:8080/v1/products/12345
```

#### Logs

Streams CloudWatch logs from production tasks. `Ctrl+C` to exit, service keeps running.

```bash
make log-prod
```

#### Shell

Opens interactive shell in one of the running production tasks. Type `exit` to close.

```bash
make shell-prod
```

#### Stop

Scales service to 0 tasks. Infrastructure remains.

```bash
make stop-prod
```

#### Start

Scales service back to 1 task.

```bash
make start-prod
```

#### Destroy

Deletes all AWS production resources and removes local Docker images. Requires typing `yes` to confirm.

```bash
make destroy-prod
```

---

### 🧹 Cleanup

Removes local build artifacts (dist/, tmp/, docs/). Does not affect Docker containers or AWS resources.

```bash
make clean
```

---

## 📁 Project Structure

```
Go-Cart/
├── cmd/                           # Application entry point
│   └── api/
│       ├── main.go               # Server initialization and Swagger configuration
│       ├── swagger.go            # Swagger setup (dev/stage builds only)
│       └── swagger_prod.go       # Empty Swagger (prod builds)
│
├── internal/                      # Application code (Go project layout standard)
│   ├── handlers/                 # HTTP request/response handling
│   │   └── product_handler.go
│   ├── models/                   # Data structures
│   │   ├── error.go
│   │   └── product.go
│   ├── repository/               # Data access layer (currently in-memory)
│   │   └── product_repository.go
│   ├── router/                   # Route registration
│   │   └── router.go
│   └── services/                 # Business logic
│       └── product_service.go
│
├── terraform/                     # Infrastructure as code
│   ├── modules/                  # Reusable Terraform modules
│   │   ├── ecr/
│   │   │   ├── main.tf           # ECR repository for Docker images
│   │   │   ├── variables.tf
│   │   │   └── outputs.tf
│   │   ├── networking/
│   │   │   ├── main.tf           # VPC and subnets
│   │   │   ├── variables.tf
│   │   │   └── outputs.tf
│   │   ├── security/
│   │   │   ├── main.tf           # Security groups
│   │   │   ├── variables.tf
│   │   │   └── outputs.tf
│   │   ├── iam/
│   │   │   ├── main.tf           # IAM roles (uses AWS Academy LabRole)
│   │   │   ├── variables.tf
│   │   │   └── outputs.tf
│   │   ├── ecs/
│   │   │   ├── main.tf           # ECS Fargate cluster and service
│   │   │   ├── variables.tf
│   │   │   └── outputs.tf
│   │   ├── rds/                  # RDS MySQL (commented, not created)
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   └── outputs.tf
│   │   └── secrets/              # Secrets Manager (commented, not created)
│   │       ├── main.tf
│   │       ├── variables.tf
│   │       └── outputs.tf
│   ├── stage/                    # Staging environment configuration
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── terraform.tfvars
│   └── prod/                     # Production environment configuration
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       └── terraform.tfvars
│
├── docs/                          # Swagger generated docs (auto-generated)
├── Makefile                       # Build and deployment commands
├── Dockerfile                     # Multi-stage Docker build
├── docker-compose.yml            # Service orchestration
├── .air.toml                      # Hot-reload configuration
└── go.mod                         # Go module definition
```
