# resource "random_password" "db_password" {
#   length  = 32
#   special = true
# }
#
# resource "aws_secretsmanager_secret" "db_password" {
#   name = "${var.project_name}-${var.environment}-db-password"
#
#   tags = {
#     Name        = "${var.project_name}-${var.environment}-db-password"
#     Environment = var.environment
#     Project     = var.project_name
#   }
# }
#
# resource "aws_secretsmanager_secret_version" "db_password" {
#   secret_id     = aws_secretsmanager_secret.db_password.id
#   secret_string = random_password.db_password.result
# }
