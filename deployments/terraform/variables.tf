variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-west-2"
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "nrp-aws-kip"
}

variable "create_vpc" {
  description = "Whether to create a new VPC"
  type        = bool
  default     = true
}

variable "existing_vpc_id" {
  description = "Existing VPC ID (if create_vpc is false)"
  type        = string
  default     = ""
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.100.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones for subnets"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

variable "nrp_cidr_blocks" {
  description = "CIDR blocks for NRP network (for security group rules)"
  type        = list(string)
  default     = ["10.0.0.0/8"] # Update with actual NRP CIDR
}

variable "create_key_pair" {
  description = "Whether to create an SSH key pair"
  type        = bool
  default     = false
}

variable "ssh_public_key" {
  description = "SSH public key for EC2 instances"
  type        = string
  default     = ""
}

variable "create_vpn_gateway" {
  description = "Whether to create a VPN gateway for NRP connectivity"
  type        = bool
  default     = false
}

variable "log_retention_days" {
  description = "CloudWatch Logs retention period in days"
  type        = number
  default     = 7
}

variable "common_tags" {
  description = "Common tags for all resources"
  type        = map(string)
  default = {
    Project    = "NRP-AWS-KIP"
    ManagedBy  = "Terraform"
    Repository = "github.com/scttfrdmn/nrp-aws-kip"
  }
}
