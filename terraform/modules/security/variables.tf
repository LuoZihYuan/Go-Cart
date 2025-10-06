variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment (stage/prod)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}
