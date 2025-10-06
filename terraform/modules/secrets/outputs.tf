# output "db_password" {
#   description = "Database password"
#   value       = random_password.db_password.result
#   sensitive   = true
# }
#
# output "secret_arn" {
#   description = "Secret ARN"
#   value       = aws_secretsmanager_secret.db_password.arn
# }
