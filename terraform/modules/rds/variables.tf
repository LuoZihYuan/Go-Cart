variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "environment" {
  description = "Environment name (stage, prod)"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for DB subnet group"
  type        = list(string)
}

variable "security_group_id" {
  description = "Security group ID for RDS instance"
  type        = string
}

variable "db_username" {
  description = "Master username for database"
  type        = string
  default     = "gocart"
}

variable "db_password" {
  description = "Master password for database"
  type        = string
  sensitive   = true
  default     = "gocart-secret-password"
}
