terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  # Region auto-detected from AWS CLI config or AWS_REGION env var
}

data "aws_region" "current" {}

module "ecr" {
  source = "../modules/ecr"

  project_name = var.project_name
  environment  = var.environment
}

module "networking" {
  source = "../modules/networking"
}

module "security" {
  source = "../modules/security"

  project_name = var.project_name
  environment  = var.environment
  vpc_id       = module.networking.vpc_id
}

module "iam" {
  source = "../modules/iam"
}

module "ecs" {
  source = "../modules/ecs"

  project_name       = var.project_name
  environment        = var.environment
  aws_region         = data.aws_region.current.name
  ecr_repository_url = module.ecr.repository_url
  subnet_ids         = module.networking.public_subnet_ids
  security_group_id  = module.security.ecs_security_group_id
  execution_role_arn = module.iam.ecs_task_execution_role_arn
  task_role_arn      = module.iam.ecs_task_role_arn
  desired_count      = var.desired_count
  task_cpu           = var.task_cpu
  task_memory        = var.task_memory
}

# Uncomment when RDS and Secrets are needed
# module "secrets" {
#   source = "../modules/secrets"
#
#   project_name = var.project_name
#   environment  = var.environment
# }
#
# module "rds" {
#   source = "../modules/rds"
#
#   project_name      = var.project_name
#   environment       = var.environment
#   subnet_ids        = module.networking.public_subnet_ids
#   security_group_id = module.security.rds_security_group_id
#   master_password   = module.secrets.db_password
# }
