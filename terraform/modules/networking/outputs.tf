output "vpc_id" {
  description = "Default VPC ID"
  value       = data.aws_vpc.default.id
}

output "public_subnet_ids" {
  description = "Default subnet IDs"
  value       = data.aws_subnets.default.ids
}
