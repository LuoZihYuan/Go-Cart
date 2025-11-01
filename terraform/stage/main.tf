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

# Conditionally create RDS module
module "rds" {
  source = "../modules/rds"
  count  = var.db_type == "mysql" ? 1 : 0

  project_name      = var.project_name
  environment       = var.environment
  subnet_ids        = module.networking.public_subnet_ids
  security_group_id = module.security.rds_security_group_id
  db_username       = var.db_username
  db_password       = var.db_password
}

# Conditionally create DynamoDB module
module "dynamodb" {
  source = "../modules/dynamodb"
  count  = var.db_type == "dynamo" ? 1 : 0

  project_name = var.project_name
  environment  = var.environment
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

  # Database configuration
  db_type = var.db_type

  # MySQL configuration (only used when db_type=mysql)
  mysql_host     = var.db_type == "mysql" ? module.rds[0].address : ""
  mysql_port     = var.db_type == "mysql" ? module.rds[0].port : 3306
  mysql_database = var.db_type == "mysql" ? module.rds[0].database_name : ""
  mysql_user     = var.db_type == "mysql" ? module.rds[0].username : ""
  mysql_password = var.db_type == "mysql" ? module.rds[0].password : ""

  # DynamoDB configuration (only used when db_type=dynamo)
  dynamodb_table_prefix = var.db_type == "dynamo" ? module.dynamodb[0].table_prefix : ""
}
