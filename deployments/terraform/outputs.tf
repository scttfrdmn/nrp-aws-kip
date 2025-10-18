output "vpc_id" {
  description = "VPC ID"
  value       = var.create_vpc ? aws_vpc.nrp_burst[0].id : var.existing_vpc_id
}

output "subnet_ids" {
  description = "Subnet IDs"
  value       = var.create_vpc ? aws_subnet.nrp_burst[*].id : []
}

output "security_group_id" {
  description = "Security group ID for burst instances"
  value       = aws_security_group.nrp_burst_instances.id
}

output "iam_instance_profile_name" {
  description = "IAM instance profile name"
  value       = aws_iam_instance_profile.nrp_burst.name
}

output "key_pair_name" {
  description = "SSH key pair name"
  value       = var.create_key_pair ? aws_key_pair.nrp_burst[0].key_name : ""
}

output "cloudwatch_log_group" {
  description = "CloudWatch log group name"
  value       = aws_cloudwatch_log_group.nrp_burst.name
}

output "config_snippet" {
  description = "Configuration snippet for config.yaml"
  value = <<-EOT
    aws:
      region: ${var.aws_region}
      vpcId: ${var.create_vpc ? aws_vpc.nrp_burst[0].id : var.existing_vpc_id}
      subnetIds:
        ${join("\n    ", [for subnet in (var.create_vpc ? aws_subnet.nrp_burst[*].id : []) : "- ${subnet}"])}
      securityGroupIds:
        - ${aws_security_group.nrp_burst_instances.id}
      keyName: ${var.create_key_pair ? aws_key_pair.nrp_burst[0].key_name : ""}
      iamInstanceProfile: ${aws_iam_instance_profile.nrp_burst.name}
  EOT
}
