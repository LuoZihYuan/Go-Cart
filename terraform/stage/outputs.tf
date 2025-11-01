output "ecr_repository_url" {
  description = "ECR repository URL"
  value       = module.ecr.repository_url
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = module.ecs.cluster_name
}

output "ecs_service_name" {
  description = "ECS service name"
  value       = module.ecs.service_name
}

output "desired_count" {
  description = "Desired count of ECS tasks"
  value       = var.desired_count
}

output "rds_address" {
  description = "RDS address (only populated when db_type=mysql)"
  value       = var.db_type == "mysql" ? module.rds[0].address : "N/A"
}

output "dynamodb_tables" {
  description = "DynamoDB table names (only populated when db_type=dynamo)"
  value = var.db_type == "dynamo" ? {
    products = module.dynamodb[0].products_table_name
    carts    = module.dynamodb[0].carts_table_name
  } : null
}
