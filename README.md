# **🛒 Go-Cart 🛒**

A RESTful e-commerce API for managing products, shopping carts, and payments, built in Go.

## **📋 Table of Contents**

- [✅ Prerequisites](#-prerequisites)
- [🚀 Getting Started](#-getting-started)
  - [⚙️ Setup](#️-setup)
  - [💻 Development (Local)](#-development-local)
  - [🧪 Staging (AWS)](#-staging-aws)
  - [🏭 Production (AWS)](#-production-aws)
  - [🧹 Cleanup](#-cleanup)
- [📁 Project Structure](#-project-structure)

## **✅ Prerequisites**

- **Docker Desktop** - For containerized development
- **Go 1.25.1+** - For local IDE support
- **AWS CLI** - For AWS deployments, configured with `aws configure`
- **Terraform 1.0+** - For infrastructure as code
- **MySQL Client** - For AWS MySQL deployments (`brew install mysql-client` on Mac)
- **VS Code or GoLand** (Optional) - For IDE features with autocomplete

## **🚀 Getting Started**

### **⚙️ Setup**

Downloads Go dependencies for IDE autocomplete and navigation features. Run this once after cloning the repository.

```bash
make setup
```

### **💻 Development (Local)**

#### **Deploy**

Builds Docker image with Swagger enabled, starts container with Air hot-reload, and mounts source code as volume. Choose your database backend: `memory` (default), `mysql` (local MySQL with persistent storage), or `dynamo` (local DynamoDB).

```bash
make deploy-dev              # In-memory storage (default)
make deploy-dev db=mysql     # Local MySQL 8.4.6 with auto-initialized schema
make deploy-dev db=dynamo    # Local DynamoDB with auto-created tables
```

When you edit any `.go` file and save, Air automatically detects changes, regenerates Swagger docs, recompiles the binary, and restarts the application (typically 2-5 seconds).

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

#### **Management**

```bash
# View all container logs with real-time streaming (Ctrl+C to exit)
make log-dev
# Open interactive shell inside the running container (type 'exit' to close)
make shell-dev
# Stop and remove all development containers
make stop-dev
# Remove containers, volumes, images, and build cache
make destroy-dev
```

### **🧪 Staging (AWS)**

#### **Deploy**

First deployment creates ECR repository and infrastructure. Subsequent deployments build new image, push to ECR, and update ECS service. Waits for service stability and displays public IP. Choose your database backend.

```bash
make deploy-stage              # In-memory storage (default)
make deploy-stage db=mysql     # AWS RDS MySQL 8.4.6 (db.t3.micro)
make deploy-stage db=dynamo    # AWS DynamoDB with PAY_PER_REQUEST billing
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

#### **Management**

```bash
# Stream CloudWatch logs in real-time (Ctrl+C to exit)
make log-stage
# Open interactive shell in the running ECS task (type 'exit' to close)
make shell-stage
# Scale service to 0 tasks (infrastructure remains, no compute costs)
make stop-stage
# Scale service back to 1 task
make start-stage
# Delete all AWS staging resources (requires typing 'yes' to confirm)
make destroy-stage
```

### **🏭 Production (AWS)**

#### **Deploy**

Swagger disabled for security. Separate ECR repository. Waits for service stability and displays public IP. Choose your database backend.

```bash
make deploy-prod              # In-memory storage (default)
make deploy-prod db=mysql     # AWS RDS MySQL 8.4.6 (db.t3.micro)
make deploy-prod db=dynamo    # AWS DynamoDB with PAY_PER_REQUEST billing
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

#### **Management**

```bash
# Stream CloudWatch logs from production tasks (Ctrl+C to exit)
make log-prod
# Open interactive shell in one of the running production tasks (type 'exit' to close)
make shell-prod
# Scale service to 0 tasks (infrastructure remains)
make stop-prod
# Scale service back to 2 tasks
make start-prod
# Delete all AWS production resources (requires typing 'yes' to confirm)
make destroy-prod
```

### **🧹 Cleanup**

Removes local build artifacts (dist/, tmp/, docs/). Does not affect Docker containers or AWS resources.

```bash
make clean
```

## **📁 Project Structure**

```
Go-Cart/
├── cmd/                           # Application entry point
│   └── api/
│       ├── main.go               # Server initialization and database switching
│       ├── swagger.go            # Swagger setup (dev/stage builds only)
│       └── swagger_prod.go       # Empty Swagger (prod builds)
│
├── internal/                      # Application code (Go project layout standard)
│   ├── handlers/                 # HTTP request/response handling
│   │   ├── cart_handler.go
│   │   └── product_handler.go
│   ├── models/                   # Data structures
│   │   ├── cart.go
│   │   ├── error.go
│   │   └── product.go
│   ├── repository/               # Data access layer
│   │   ├── interfaces.go         # Repository contracts
│   │   ├── product_memory.go     # In-memory implementation
│   │   ├── product_mysql.go      # MySQL implementation
│   │   ├── product_dynamodb.go   # DynamoDB implementation
│   │   ├── cart_memory.go
│   │   ├── cart_mysql.go
│   │   └── cart_dynamodb.go
│   ├── router/                   # Route registration
│   │   └── router.go
│   └── services/                 # Business logic
│       ├── cart_service.go
│       └── product_service.go
│
├── scripts/                       # Database initialization
│   ├── mysql/
│   │   └── init.sql              # MySQL schema (products, carts, cart_items)
│   └── dynamodb/
│       └── init-local.sh         # DynamoDB Local table creation
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
│   │   │   ├── main.tf           # Security groups (ECS and RDS)
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
│   │   ├── rds/
│   │   │   ├── main.tf           # RDS MySQL 8.4.6 (db.t3.micro)
│   │   │   ├── variables.tf
│   │   │   └── outputs.tf
│   │   └── dynamodb/
│   │       ├── main.tf           # DynamoDB tables (Products, Carts)
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
├── docker-compose.yml             # Service orchestration
├── .air.toml                      # Hot-reload configuration
└── go.mod                         # Go module definition
```
