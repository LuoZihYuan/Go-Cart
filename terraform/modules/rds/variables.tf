# variable "project_name" {
#   description = "Project name"
#   type        = string
# }
#
# variable "environment" {
#   description = "Environment (stage/prod)"
#   type        = string
# }
#
# variable "subnet_ids" {
#   description = "Subnet IDs for RDS"
#   type        = list(string)
# }
#
# variable "security_group_id" {
#   description = "Security group ID for RDS"
#   type        = string
# }
#
# variable "database_name" {
#   description = "Database name"
#   type        = string
#   default     = "ecommerce"
# }
#
# variable "master_username" {
#   description = "Master username"
#   type        = string
#   default     = "admin"
# }
#
# variable "master_password" {
#   description = "Master password"
#   type        = string
#   sensitive   = true
# }
#
# variable "instance_class" {
#   description = "RDS instance class"
#   type        = string
#   default     = "db.t3.micro"
# }
#
# variable "skip_final_snapshot" {
#   description = "Skip final snapshot on deletion"
#   type        = bool
#   default     = true
# }
