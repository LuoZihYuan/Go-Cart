output "ecs_security_group_id" {
  description = "ECS security group ID"
  value       = aws_security_group.ecs.id
}

# Uncomment when RDS is enabled
# output "rds_security_group_id" {
#   description = "RDS security group ID"
#   value       = aws_security_group.rds.id
# }
