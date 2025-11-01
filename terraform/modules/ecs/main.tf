resource "aws_ecs_cluster" "main" {
  name = "${var.project_name}-${var.environment}"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name        = "${var.project_name}-${var.environment}"
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/${var.project_name}-${var.environment}"
  retention_in_days = 7

  tags = {
    Name        = "${var.project_name}-${var.environment}"
    Environment = var.environment
    Project     = var.project_name
  }
}

# Build environment variables based on db_type
locals {
  base_environment = [
    {
      name  = "GIN_MODE"
      value = "release"
    },
    {
      name  = "DB_TYPE"
      value = var.db_type
    }
  ]

  mysql_environment = var.db_type == "mysql" ? [
    {
      name  = "MYSQL_HOST"
      value = var.mysql_host
    },
    {
      name  = "MYSQL_PORT"
      value = tostring(var.mysql_port)
    },
    {
      name  = "MYSQL_DATABASE"
      value = var.mysql_database
    },
    {
      name  = "MYSQL_USER"
      value = var.mysql_user
    },
    {
      name  = "MYSQL_PASSWORD"
      value = var.mysql_password
    },
    {
      name  = "MYSQL_MAX_CONNECTIONS"
      value = "20"
    },
    {
      name  = "MYSQL_MAX_IDLE_CONNECTIONS"
      value = "5"
    }
  ] : []

  dynamo_environment = var.db_type == "dynamo" ? [
    {
      name  = "DYNAMODB_REGION"
      value = var.aws_region
    },
    {
      name  = "DYNAMODB_TABLE_PREFIX"
      value = var.dynamodb_table_prefix
    }
  ] : []

  # Combine all environments
  environment = concat(
    local.base_environment,
    local.mysql_environment,
    local.dynamo_environment
  )
}

resource "aws_ecs_task_definition" "app" {
  family                   = "${var.project_name}-${var.environment}"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = var.execution_role_arn
  task_role_arn            = var.task_role_arn

  container_definitions = jsonencode([
    {
      name      = "api"
      image     = "${var.ecr_repository_url}:${var.environment}"
      essential = true

      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        }
      ]

      environment = local.environment

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  tags = {
    Name        = "${var.project_name}-${var.environment}"
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_ecs_service" "app" {
  name                   = "${var.project_name}-${var.environment}"
  cluster                = aws_ecs_cluster.main.id
  task_definition        = aws_ecs_task_definition.app.arn
  desired_count          = var.desired_count
  launch_type            = "FARGATE"
  enable_execute_command = true
  force_new_deployment   = true

  network_configuration {
    subnets          = var.subnet_ids
    security_groups  = [var.security_group_id]
    assign_public_ip = true
  }

  tags = {
    Name        = "${var.project_name}-${var.environment}"
    Environment = var.environment
    Project     = var.project_name
  }
}
