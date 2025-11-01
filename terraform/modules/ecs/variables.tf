variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment (stage/prod)"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "ecr_repository_url" {
  description = "ECR repository URL"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs for ECS tasks"
  type        = list(string)
}

variable "security_group_id" {
  description = "Security group ID for ECS tasks"
  type        = string
}

variable "execution_role_arn" {
  description = "ECS task execution role ARN"
  type        = string
}

variable "task_role_arn" {
  description = "ECS task role ARN"
  type        = string
}

variable "desired_count" {
  description = "Desired number of tasks"
  type        = number
  default     = 1
}

variable "task_cpu" {
  description = "Task CPU units"
  type        = string
  default     = "256"
}

variable "task_memory" {
  description = "Task memory in MB"
  type        = string
  default     = "512"
}

# Database configuration
variable "db_type" {
  description = "Database type: memory, mysql, or dynamo"
  type        = string
  default     = "memory"
}

# MySQL configuration
variable "mysql_host" {
  description = "MySQL host address"
  type        = string
  default     = ""
}

variable "mysql_port" {
  description = "MySQL port"
  type        = number
  default     = 3306
}

variable "mysql_database" {
  description = "MySQL database name"
  type        = string
  default     = ""
}

variable "mysql_user" {
  description = "MySQL username"
  type        = string
  default     = ""
}

variable "mysql_password" {
  description = "MySQL password"
  type        = string
  sensitive   = true
  default     = ""
}
