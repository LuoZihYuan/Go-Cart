output "products_table_name" {
  description = "Name of the Products DynamoDB table"
  value       = "Products"
}

output "products_table_arn" {
  description = "ARN of the Products DynamoDB table"
  value       = aws_dynamodb_table.products.arn
}

output "carts_table_name" {
  description = "Name of the Carts DynamoDB table"
  value       = "Carts"
}

output "carts_table_arn" {
  description = "ARN of the Carts DynamoDB table"
  value       = aws_dynamodb_table.carts.arn
}
