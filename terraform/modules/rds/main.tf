# resource "aws_db_subnet_group" "main" {
#   name       = "${var.project_name}-${var.environment}"
#   subnet_ids = var.subnet_ids
#
#   tags = {
#     Name        = "${var.project_name}-${var.environment}"
#     Environment = var.environment
#     Project     = var.project_name
#   }
# }
#
# resource "aws_db_instance" "main" {
#   identifier             = "${var.project_name}-${var.environment}"
#   engine                 = "mysql"
#   engine_version         = "8.0"
#   instance_class         = var.instance_class
#   allocated_storage      = 20
#   storage_type           = "gp3"
#   db_name                = var.database_name
#   username               = var.master_username
#   password               = var.master_password
#   db_subnet_group_name   = aws_db_subnet_group.main.name
#   vpc_security_group_ids = [var.security_group_id]
#   skip_final_snapshot    = var.skip_final_snapshot
#   publicly_accessible    = false
#
#   tags = {
#     Name        = "${var.project_name}-${var.environment}"
#     Environment = var.environment
#     Project     = var.project_name
#   }
# }
