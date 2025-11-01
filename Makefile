.PHONY: help setup deploy-dev log-dev shell-dev stop-dev destroy-dev deploy-stage log-stage shell-stage stop-stage start-stage destroy-stage deploy-prod log-prod shell-prod stop-prod start-prod destroy-prod clean

help:  ## Show available commands
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

setup:  ## Setup development environment
	@echo "Checking Docker..."
	@docker --version > /dev/null || (echo "Install Docker Desktop" && exit 1)
	@docker compose version > /dev/null || (echo "Install Docker Compose" && exit 1)
	@echo "Installing Go dependencies for IDE support..."
	@go mod download
	@echo "Setup complete! Run 'make deploy-dev' to start."

# =============================================================================
# Development (Local)
# =============================================================================

deploy-dev:  ## Start local development with hot-reload (db=memory|mysql|dynamo)
	@DB_TYPE=$${db:-memory}; \
	export DB_TYPE=$$DB_TYPE; \
	docker-compose --profile dev --profile $$DB_TYPE up --build -d; \
	echo "Dev container started in background with DB_TYPE=$$DB_TYPE"; \
	echo "View logs: make log-dev"

log-dev:  ## View all running container logs (Ctrl+C to exit)
	docker-compose logs -f || true

shell-dev:  ## Open interactive shell in local API container
	docker exec -it api.gocart-dev sh

stop-dev:  ## Stop all local development containers
	@docker-compose --profile dev --profile mysql --profile dynamo down --remove-orphans
	@echo "All development containers stopped"

destroy-dev:  ## Destroy local Docker environment including volumes and build cache
	@echo "Stopping and removing containers..."
	@docker-compose --profile dev --profile mysql --profile dynamo down --volumes --remove-orphans
	@echo "Removing dev images..."
	@docker rmi api.gocart:dev mysql:8.4.6 amazon/dynamodb-local:latest amazon/aws-cli:latest 2>/dev/null || true
	@echo "Removing named volumes..."
	@docker volume rm go-cart_mysql-data 2>/dev/null || true
	@echo "Removing dangling volumes..."
	@docker volume prune -f
	@echo "Removing build history and cache..."
	@docker buildx history ls | grep gocart_dev | awk '{print $$1}' | xargs docker buildx history rm 2>/dev/null || true
	@docker system prune -f
	@echo "Development environment destroyed"

# =============================================================================
# Staging (AWS)
# =============================================================================

deploy-stage:  ## Deploy to AWS staging (requires Terraform) (db=memory|mysql|dynamo)
	@if [ ! -d "terraform/stage" ]; then \
		echo "Terraform not configured yet."; \
		echo "deploy-stage is for AWS deployment only."; \
		echo "Use 'make deploy-dev' for local development."; \
		exit 1; \
	fi
	@DB_TYPE=$${db:-memory}; \
	echo "Getting AWS configuration..."; \
	AWS_REGION=$$(aws configure get region) || (echo "AWS region not configured. Run: aws configure" && exit 1); \
	\
	echo "Checking if ECR repository exists..."; \
	if ! aws ecr describe-repositories --repository-names go-cart-stage > /dev/null 2>&1; then \
		echo "First-time deployment detected."; \
		echo "Step 1: Initializing Terraform and creating ECR repository..."; \
		cd terraform/stage && terraform init > /dev/null 2>&1 && terraform apply -target=module.ecr -var="db_type=$$DB_TYPE" -auto-approve && cd ../..; \
		ECR_REPO_URL=$$(aws ecr describe-repositories \
			--repository-names go-cart-stage \
			--query 'repositories[0].repositoryUri' \
			--output text); \
		echo "ECR repository created: $$ECR_REPO_URL"; \
		\
		echo "Step 2: Building and pushing Docker image..."; \
		docker-compose --profile stage build && \
		aws ecr get-login-password --region $$AWS_REGION | docker login --username AWS --password-stdin $$(echo $$ECR_REPO_URL | cut -d'/' -f1) && \
		docker tag api.gocart:stage $$ECR_REPO_URL:stage && \
		docker push $$ECR_REPO_URL:stage && \
		\
		echo "Step 3: Creating remaining infrastructure with DB_TYPE=$$DB_TYPE..." && \
		cd terraform/stage && terraform apply -var="db_type=$$DB_TYPE" -auto-approve && cd ../..; \
	else \
		echo "Infrastructure exists. Updating deployment..."; \
		ECR_REPO_URL=$$(aws ecr describe-repositories \
			--repository-names go-cart-stage \
			--query 'repositories[0].repositoryUri' \
			--output text); \
		echo "Using ECR repository: $$ECR_REPO_URL"; \
		\
		echo "1/2 Building and pushing Docker image..." && \
		docker-compose --profile stage build && \
		aws ecr get-login-password --region $$AWS_REGION | docker login --username AWS --password-stdin $$(echo $$ECR_REPO_URL | cut -d'/' -f1) && \
		docker tag api.gocart:stage $$ECR_REPO_URL:stage && \
		docker push $$ECR_REPO_URL:stage && \
		\
		echo "2/2 Updating ECS service with DB_TYPE=$$DB_TYPE..." && \
		cd terraform/stage && terraform apply -var="db_type=$$DB_TYPE" -auto-approve && cd ../..; \
	fi && \
	echo "" && \
	echo "Staging deployed to AWS with DB_TYPE=$$DB_TYPE!" && \
	echo "Waiting for service to become stable..." && \
	aws ecs wait services-stable --cluster go-cart-stage --services go-cart-stage && \
	echo "Service is stable. Getting public IP..." && \
	TASK_ARN=$$(aws ecs list-tasks --cluster go-cart-stage --service-name go-cart-stage --query 'taskArns[0]' --output text 2>/dev/null) && \
	if [ -n "$$TASK_ARN" ] && [ "$$TASK_ARN" != "None" ]; then \
		ENI_ID=$$(aws ecs describe-tasks --cluster go-cart-stage --tasks $$TASK_ARN --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' --output text 2>/dev/null) && \
		PUBLIC_IP=$$(aws ec2 describe-network-interfaces --network-interface-ids $$ENI_ID --query 'NetworkInterfaces[0].Association.PublicIp' --output text 2>/dev/null) && \
		echo "API available at: http://$$PUBLIC_IP:8080"; \
	else \
		echo "Task not running yet. View status: make log-stage"; \
	fi

log-stage:  ## View AWS staging logs (Ctrl+C to exit)
	aws logs tail /ecs/go-cart-stage --follow --format short || true

shell-stage:  ## Open interactive shell in AWS staging container
	@echo "Finding running task..."
	@CLUSTER=$$(cd terraform/stage && terraform output -raw ecs_cluster_name 2>/dev/null && cd ../..) && \
	SERVICE=$$(cd terraform/stage && terraform output -raw ecs_service_name 2>/dev/null && cd ../..) && \
	TASK_ID=$$(aws ecs list-tasks --cluster $$CLUSTER --service-name $$SERVICE --query 'taskArns[0]' --output text 2>/dev/null | cut -d'/' -f3) && \
	if [ -z "$$TASK_ID" ] || [ "$$TASK_ID" = "None" ]; then \
		echo "No running tasks found in staging."; \
		exit 1; \
	fi && \
	echo "Connecting to task: $$TASK_ID" && \
	aws ecs execute-command \
		--cluster $$CLUSTER \
		--task $$TASK_ID \
		--container api \
		--interactive \
		--command "/bin/sh"

stop-stage:  ## Stop AWS staging service (scale to 0)
	@echo "Stopping staging service..."
	@CLUSTER=$$(cd terraform/stage && terraform output -raw ecs_cluster_name 2>/dev/null && cd ../..) && \
	SERVICE=$$(cd terraform/stage && terraform output -raw ecs_service_name 2>/dev/null && cd ../..) && \
	aws ecs update-service --cluster $$CLUSTER --service $$SERVICE --desired-count 0 > /dev/null
	@echo "Staging service stopped (infrastructure still exists)"

start-stage:  ## Start AWS staging service
	@echo "Starting staging service..."
	@CLUSTER=$$(cd terraform/stage && terraform output -raw ecs_cluster_name 2>/dev/null && cd ../..) && \
	SERVICE=$$(cd terraform/stage && terraform output -raw ecs_service_name 2>/dev/null && cd ../..) && \
	DESIRED_COUNT=$$(cd terraform/stage && terraform output -raw desired_count 2>/dev/null && cd ../..) || DESIRED_COUNT=1 && \
	aws ecs update-service --cluster $$CLUSTER --service $$SERVICE --desired-count $$DESIRED_COUNT > /dev/null
	@echo "Staging service started"

destroy-stage:  ## Destroy AWS staging infrastructure (DESTRUCTIVE)
	@echo "WARNING: This will DESTROY all staging AWS resources"
	@read -p "Type 'yes' to confirm: " confirm && [ "$$confirm" = "yes" ] || (echo "Aborted" && exit 1)
	@echo "Destroying Terraform infrastructure..."
	@cd terraform/stage && terraform destroy -auto-approve
	@echo "Cleaning local Docker images..."
	@docker rmi api.gocart:stage 2>/dev/null || true
	@docker images | grep "go-cart-stage" | awk '{print $$3}' | xargs docker rmi -f 2>/dev/null || true
	@echo ""
	@read -p "Clean up unused Docker images and containers? (y/n): " cleanup && \
	if [ "$$cleanup" = "y" ] || [ "$$cleanup" = "Y" ]; then \
		echo "Cleaning up unused Docker resources..." && \
		docker system prune -f; \
	else \
		echo "Skipping Docker cleanup"; \
	fi
	@echo "Cleanup complete"

# =============================================================================
# Production (AWS)
# =============================================================================

deploy-prod:  ## Deploy to AWS production (requires Terraform) (db=memory|mysql|dynamo)
	@if [ ! -d "terraform/prod" ]; then \
		echo "Terraform not configured yet."; \
		echo "deploy-prod is for AWS deployment only."; \
		echo "Use 'make deploy-dev' for local development."; \
		exit 1; \
	fi
	@DB_TYPE=$${db:-memory}; \
	echo "Getting AWS configuration..."; \
	AWS_REGION=$$(aws configure get region) || (echo "AWS region not configured. Run: aws configure" && exit 1); \
	\
	echo "Checking if ECR repository exists..."; \
	if ! aws ecr describe-repositories --repository-names go-cart-prod > /dev/null 2>&1; then \
		echo "First-time deployment detected."; \
		echo "Step 1: Initializing Terraform and creating ECR repository..."; \
		cd terraform/prod && terraform init > /dev/null 2>&1 && terraform apply -target=module.ecr -var="db_type=$$DB_TYPE" -auto-approve && cd ../..; \
		ECR_REPO_URL=$$(aws ecr describe-repositories \
			--repository-names go-cart-prod \
			--query 'repositories[0].repositoryUri' \
			--output text); \
		echo "ECR repository created: $$ECR_REPO_URL"; \
		\
		echo "Step 2: Building and pushing Docker image..."; \
		docker-compose --profile prod build && \
		aws ecr get-login-password --region $$AWS_REGION | docker login --username AWS --password-stdin $$(echo $$ECR_REPO_URL | cut -d'/' -f1) && \
		docker tag api.gocart:prod $$ECR_REPO_URL:prod && \
		docker push $$ECR_REPO_URL:prod && \
		\
		echo "Step 3: Creating remaining infrastructure with DB_TYPE=$$DB_TYPE..." && \
		cd terraform/prod && terraform apply -var="db_type=$$DB_TYPE" -auto-approve && cd ../..; \
	else \
		echo "Infrastructure exists. Updating deployment..."; \
		ECR_REPO_URL=$$(aws ecr describe-repositories \
			--repository-names go-cart-prod \
			--query 'repositories[0].repositoryUri' \
			--output text); \
		echo "Using ECR repository: $$ECR_REPO_URL"; \
		\
		echo "1/2 Building and pushing Docker image..." && \
		docker-compose --profile prod build && \
		aws ecr get-login-password --region $$AWS_REGION | docker login --username AWS --password-stdin $$(echo $$ECR_REPO_URL | cut -d'/' -f1) && \
		docker tag api.gocart:prod $$ECR_REPO_URL:prod && \
		docker push $$ECR_REPO_URL:prod && \
		\
		echo "2/2 Updating ECS service with DB_TYPE=$$DB_TYPE..." && \
		cd terraform/prod && terraform apply -var="db_type=$$DB_TYPE" -auto-approve && cd ../..; \
	fi && \
	echo "" && \
	echo "Production deployed to AWS with DB_TYPE=$$DB_TYPE!" && \
	echo "Waiting for service to become stable..." && \
	aws ecs wait services-stable --cluster go-cart-prod --services go-cart-prod && \
	echo "Service is stable. Getting public IP..." && \
	TASK_ARN=$$(aws ecs list-tasks --cluster go-cart-prod --service-name go-cart-prod --query 'taskArns[0]' --output text 2>/dev/null) && \
	if [ -n "$$TASK_ARN" ] && [ "$$TASK_ARN" != "None" ]; then \
		ENI_ID=$$(aws ecs describe-tasks --cluster go-cart-prod --tasks $$TASK_ARN --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' --output text 2>/dev/null) && \
		PUBLIC_IP=$$(aws ec2 describe-network-interfaces --network-interface-ids $$ENI_ID --query 'NetworkInterfaces[0].Association.PublicIp' --output text 2>/dev/null) && \
		echo "API available at: http://$$PUBLIC_IP:8080"; \
	else \
		echo "Task not running yet. View status: make log-prod"; \
	fi

log-prod:  ## View AWS production logs (Ctrl+C to exit)
	aws logs tail /ecs/go-cart-prod --follow --format short || true

shell-prod:  ## Open interactive shell in AWS production container
	@echo "Finding running task..."
	@CLUSTER=$$(cd terraform/prod && terraform output -raw ecs_cluster_name 2>/dev/null && cd ../..) && \
	SERVICE=$$(cd terraform/prod && terraform output -raw ecs_service_name 2>/dev/null && cd ../..) && \
	TASK_ID=$$(aws ecs list-tasks --cluster $$CLUSTER --service-name $$SERVICE --query 'taskArns[0]' --output text 2>/dev/null | cut -d'/' -f3) && \
	if [ -z "$$TASK_ID" ] || [ "$$TASK_ID" = "None" ]; then \
		echo "No running tasks found in production."; \
		exit 1; \
	fi && \
	echo "Connecting to task: $$TASK_ID" && \
	aws ecs execute-command \
		--cluster $$CLUSTER \
		--task $$TASK_ID \
		--container api \
		--interactive \
		--command "/bin/sh"

stop-prod:  ## Stop AWS production service (scale to 0)
	@echo "Stopping production service..."
	@CLUSTER=$$(cd terraform/prod && terraform output -raw ecs_cluster_name 2>/dev/null && cd ../..) && \
	SERVICE=$$(cd terraform/prod && terraform output -raw ecs_service_name 2>/dev/null && cd ../..) && \
	aws ecs update-service --cluster $$CLUSTER --service $$SERVICE --desired-count 0 > /dev/null
	@echo "Production service stopped (infrastructure still exists)"

start-prod:  ## Start AWS production service
	@echo "Starting production service..."
	@CLUSTER=$$(cd terraform/prod && terraform output -raw ecs_cluster_name 2>/dev/null && cd ../..) && \
	SERVICE=$$(cd terraform/prod && terraform output -raw ecs_service_name 2>/dev/null && cd ../..) && \
	DESIRED_COUNT=$$(cd terraform/prod && terraform output -raw desired_count 2>/dev/null && cd ../..) || DESIRED_COUNT=2 && \
	aws ecs update-service --cluster $$CLUSTER --service $$SERVICE --desired-count $$DESIRED_COUNT > /dev/null
	@echo "Production service started"

destroy-prod:  ## Destroy AWS production infrastructure (DESTRUCTIVE)
	@echo "WARNING: This will DESTROY all production AWS resources"
	@read -p "Type 'yes' to confirm: " confirm && [ "$$confirm" = "yes" ] || (echo "Aborted" && exit 1)
	@echo "Destroying Terraform infrastructure..."
	@cd terraform/prod && terraform destroy -auto-approve
	@echo "Cleaning local Docker images..."
	@docker rmi api.gocart:prod 2>/dev/null || true
	@docker images | grep "go-cart-prod" | awk '{print $$3}' | xargs docker rmi -f 2>/dev/null || true
	@echo ""
	@read -p "Clean up unused Docker images and containers? (y/n): " cleanup && \
	if [ "$$cleanup" = "y" ] || [ "$$cleanup" = "Y" ]; then \
		echo "Cleaning up unused Docker resources..." && \
		docker system prune -f; \
	else \
		echo "Skipping Docker cleanup"; \
	fi
	@echo "Cleanup complete"

# =============================================================================
# Utilities
# =============================================================================

clean:  ## Remove local build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf dist/ tmp/ docs/
	@echo "Cleanup complete!"