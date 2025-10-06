output "ecs_task_execution_role_arn" {
  description = "ECS task execution role ARN (using LabRole)"
  value       = data.aws_iam_role.lab.arn
}

output "ecs_task_role_arn" {
  description = "ECS task role ARN (using LabRole)"
  value       = data.aws_iam_role.lab.arn
}
