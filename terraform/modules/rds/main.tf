# RDS MySQL Instance for Go-Cart E-commerce

# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${var.project_name}-${var.environment}-db-subnet"
  subnet_ids = var.subnet_ids

  tags = {
    Name        = "${var.project_name}-${var.environment}-db-subnet"
    Environment = var.environment
    Project     = var.project_name
  }
}

# DB Parameter Group
resource "aws_db_parameter_group" "main" {
  name   = "${var.project_name}-${var.environment}-mysql"
  family = "mysql8.4"

  parameter {
    name  = "character_set_server"
    value = "utf8mb4"
  }

  parameter {
    name  = "collation_server"
    value = "utf8mb4_unicode_ci"
  }

  tags = {
    Name        = "${var.project_name}-${var.environment}-mysql"
    Environment = var.environment
    Project     = var.project_name
  }
}

# RDS MySQL Instance
resource "aws_db_instance" "main" {
  identifier     = "${var.project_name}-${var.environment}"
  engine         = "mysql"
  engine_version = "8.4.6"
  instance_class = "db.t3.micro"

  allocated_storage     = 20
  max_allocated_storage = 100
  storage_type          = "gp3"
  storage_encrypted     = false

  db_name  = "gocart"
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  parameter_group_name   = aws_db_parameter_group.main.name
  vpc_security_group_ids = [var.security_group_id]

  publicly_accessible = true # Required for AWS Academy/Learner Lab
  skip_final_snapshot = true # For development/assignment purposes
  deletion_protection = false

  backup_retention_period = 0 # No backups for assignments

  tags = {
    Name        = "${var.project_name}-${var.environment}"
    Environment = var.environment
    Project     = var.project_name
  }
}

# Execute schema initialization SQL
resource "null_resource" "schema_init" {
  # Trigger re-run if init.sql changes
  triggers = {
    init_sql_hash = filemd5("${path.root}/../../scripts/mysql/init.sql")
  }

  provisioner "local-exec" {
    command = <<-EOT
      mysql -h ${aws_db_instance.main.address} \
            -P ${aws_db_instance.main.port} \
            -u ${var.db_username} \
            -p${var.db_password} \
            ${aws_db_instance.main.db_name} \
            < ${path.root}/../../scripts/mysql/init.sql
    EOT
  }

  depends_on = [aws_db_instance.main]
}
