# DynamoDB Tables for Go-Cart E-commerce

# Products table
resource "aws_dynamodb_table" "products" {
  name         = "${var.project_name}-${var.environment}-Products"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "product_id"

  attribute {
    name = "product_id"
    type = "N"
  }

  tags = {
    Name        = "${var.project_name}-${var.environment}-Products"
    Environment = var.environment
    Project     = var.project_name
  }
}

# Carts table
resource "aws_dynamodb_table" "carts" {
  name         = "${var.project_name}-${var.environment}-Carts"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "cart_id"

  attribute {
    name = "cart_id"
    type = "N"
  }

  tags = {
    Name        = "${var.project_name}-${var.environment}-Carts"
    Environment = var.environment
    Project     = var.project_name
  }
}
