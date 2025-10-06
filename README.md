# Go-Cart

A RESTful e-commerce API for managing products, shopping carts, and payments, built in Go.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Development (Local)](#development-local)
  - [Staging (AWS)](#staging-aws)
  - [Production (AWS)](#production-aws)
- [Project Structure](#project-structure)

---

## Prerequisites

### Required

- **Docker Desktop** - For containerized development
- **Go 1.25.1+** - For local IDE support
- **AWS CLI** - For AWS deployments (configured with `aws configure`)
- **Terraform 1.0+** - For infrastructure as code

### Optional

- **VS Code or GoLand** - For IDE features with autocomplete

---

## Getting Started

### Initial Setup

```bash
# Clone the repository
git clone <repository-url>
cd go-cart

# One-time setup
make setup

# Start development
make deploy-dev
```

### Development (Local)

Local development runs in Docker with hot-reload enabled.

The API is now available at http://localhost:8080 and Swagger UI at http://localhost:8080/swagger/index.html

#### Viewing Logs

```bash
make log-dev
```

Shows all container logs with real-time streaming. Ctrl+C to exit (container keeps running).

#### Shell Access

```bash
make shell-dev
```

Opens interactive shell inside the running container:

```bash
/app # ls
/app # ps aux
/app # curl http://localhost:8080/v1/products/123
/app # exit
```

#### Hot-Reload Workflow

1. Edit any `.go` file and save
2. Air detects the change automatically
3. Runs `swag init` to regenerate Swagger docs
4. Recompiles with `go build --tags dev`
5. Restarts the application
6. Check `make log-dev` to see rebuild status

Rebuild time: 2-5 seconds typically

#### Testing the API

**Using Swagger UI:**

1. Open http://localhost:8080/swagger/index.html
2. Expand "Products" section
3. Click "Try it out" on any endpoint
4. Fill in parameters
5. Click "Execute"
6. View response (body, status code, headers, cURL command)

**Using cURL:**

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

#### Stopping Development

**Stop container:**

```bash
make stop-dev
```

Stops and removes container, keeps images for faster restart.

**Complete cleanup:**

```bash
make destroy-dev
```

Removes containers, volumes, and images. Frees disk space.

#### Available Commands

- `make deploy-dev` - Start local development with hot-reload
- `make log-dev` - View logs (Ctrl+C to exit)
- `make shell-dev` - Open shell in container
- `make stop-dev` - Stop container
- `make destroy-dev` - Destroy Docker environment

---

### Staging (AWS)

Staging environment runs on AWS with Swagger enabled for testing.

#### Prerequisites

```bash
# Configure AWS CLI
aws configure

# Initialize Terraform
cd terraform/stage
terraform init
cd ../..
```

#### First Deployment

```bash
make deploy-stage
```

**What happens:**

1. Checks if ECR repository exists in AWS
2. If first time:
   - Creates ECR repository with Terraform
   - Builds Docker image for linux/amd64
   - Pushes image to ECR
   - Creates infrastructure (VPC, ECS, security groups)
3. If updating:
   - Builds new Docker image
   - Pushes to ECR
   - Updates ECS service with `force_new_deployment`
4. Waits for service to become stable
5. Displays public IP

#### Updating Staging

After code changes:

```bash
make deploy-stage
```

Rebuilds image, pushes to ECR, and triggers ECS redeployment.

#### Viewing Logs

```bash
make log-stage
```

Streams CloudWatch logs in real-time. Ctrl+C to exit, service keeps running.

#### Shell Access

```bash
make shell-stage
```

Opens interactive shell in the running ECS task on AWS.

#### Managing the Service

**Stop service (scale to 0):**

```bash
make stop-stage
```

Stops all running tasks, infrastructure remains, no compute costs.

**Restart service:**

```bash
make start-stage
```

Scales back to 1 task.

**Destroy everything:**

```bash
make destroy-stage
```

Requires typing 'yes' to confirm. Deletes ALL AWS resources and removes local Docker images.

#### Testing the API

**Using Swagger UI:**

1. Get staging public IP from deployment output
2. Open http://<STAGING_IP>:8080/swagger/index.html
3. Expand "Products" section
4. Click "Try it out" on any endpoint
5. Fill in parameters
6. Click "Execute"
7. View response

**Using cURL:**

```bash
# Replace <STAGING_IP> with your actual staging IP

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

#### Available Commands

- `make deploy-stage` - Deploy to AWS staging
- `make log-stage` - View CloudWatch logs (Ctrl+C to exit)
- `make shell-stage` - Open shell in ECS task
- `make stop-stage` - Scale to 0 tasks
- `make start-stage` - Scale back to 1 task
- `make destroy-stage` - Destroy infrastructure (DESTRUCTIVE)

---

### Production (AWS)

Production environment runs on AWS with Swagger disabled and higher resources.

#### Prerequisites

```bash
# Configure AWS CLI (if not already done)
aws configure

# Initialize Terraform
cd terraform/prod
terraform init
cd ../..
```

#### First Deployment

```bash
make deploy-prod
```

Same process as staging, but:

- Creates 1 tasks for high availability
- Larger CPU/memory allocation
- Swagger disabled in binary
- Separate ECR repository: `go-cart-prod`

#### Updating Production

```bash
make deploy-prod
```

Rebuilds, pushes, and triggers rolling deployment across 1 tasks.

#### Viewing Logs

```bash
make log-prod
```

Streams CloudWatch logs from production tasks.

#### Shell Access

```bash
make shell-prod
```

Opens shell in one of the running production tasks.

#### Managing the Service

**Stop service:**

```bash
make stop-prod
```

Scales to 0 tasks, infrastructure remains.

**Restart service:**

```bash
make start-prod
```

Scales back to 1 tasks.

**Destroy everything:**

```bash
make destroy-prod
```

Complete infrastructure teardown (requires confirmation).

#### Testing the API

Production has Swagger disabled for security. Use cURL for testing.

**Using cURL:**

```bash
# Replace <PRODUCTION_IP> with your actual production IP

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

#### Available Commands

- `make deploy-prod` - Deploy to AWS production
- `make log-prod` - View CloudWatch logs (Ctrl+C to exit)
- `make shell-prod` - Open shell in ECS task
- `make stop-prod` - Scale to 0 tasks
- `make start-prod` - Scale back to 1 tasks
- `make destroy-prod` - Destroy infrastructure (DESTRUCTIVE)

---

## Project Structure

```
project/
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
├── docker-compose.yaml            # Service orchestration
├── .air.toml                      # Hot-reload configuration
└── go.mod                         # Go module definition
```
