# AWS Academy provides a LabRole that can be used for ECS tasks
data "aws_iam_role" "lab" {
  name = "LabRole"
}
