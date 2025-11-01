variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "go-cart"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "prod"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.1.0.0/16"
}

variable "desired_count" {
  description = "Desired number of ECS tasks"
  type        = number
  default     = 1
}

variable "task_cpu" {
  description = "CPU units for ECS task"
  type        = string
  default     = "256"
}

variable "task_memory" {
  description = "Memory for ECS task"
  type        = string
  default     = "512"
}

variable "db_type" {
  description = "Database type: memory, mysql, or dynamo"
  type        = string
  default     = "memory"

  validation {
    condition     = contains(["memory", "mysql", "dynamo"], var.db_type)
    error_message = "db_type must be memory, mysql, or dynamo"
  }
}

variable "db_username" {
  description = "Database username (used for MySQL)"
  type        = string
  default     = "gocart"
}

variable "db_password" {
  description = "Database password (used for MySQL)"
  type        = string
  sensitive   = true
  default     = "gocart-secret-password"
}
