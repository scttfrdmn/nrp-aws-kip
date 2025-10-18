terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC for burst instances (optional - can use existing VPC)
resource "aws_vpc" "nrp_burst" {
  count = var.create_vpc ? 1 : 0

  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-vpc"
    }
  )
}

# Subnets
resource "aws_subnet" "nrp_burst" {
  count = var.create_vpc ? length(var.availability_zones) : 0

  vpc_id                  = aws_vpc.nrp_burst[0].id
  cidr_block              = cidrsubnet(var.vpc_cidr, 4, count.index)
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = false

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-subnet-${var.availability_zones[count.index]}"
    }
  )
}

# Security group for burst instances
resource "aws_security_group" "nrp_burst_instances" {
  name        = "${var.project_name}-instances"
  description = "Security group for NRP burst instances"
  vpc_id      = var.create_vpc ? aws_vpc.nrp_burst[0].id : var.existing_vpc_id

  # Allow all traffic from NRP network
  ingress {
    description = "All traffic from NRP"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = var.nrp_cidr_blocks
  }

  # Allow internal communication
  ingress {
    description = "Internal VPC traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
  }

  # Allow all outbound traffic
  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-instances-sg"
    }
  )
}

# IAM role for burst instances
resource "aws_iam_role" "nrp_burst_instance" {
  name = "${var.project_name}-instance-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })

  tags = var.common_tags
}

# IAM policy for instances
resource "aws_iam_role_policy" "nrp_burst_instance" {
  name = "${var.project_name}-instance-policy"
  role = aws_iam_role.nrp_burst_instance.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ec2:DescribeInstances",
          "ec2:DescribeTags",
          "cloudwatch:PutMetricData",
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogStreams"
        ]
        Resource = "*"
      }
    ]
  })
}

# Instance profile
resource "aws_iam_instance_profile" "nrp_burst" {
  name = "${var.project_name}-instance-profile"
  role = aws_iam_role.nrp_burst_instance.name

  tags = var.common_tags
}

# SSH Key pair (optional)
resource "aws_key_pair" "nrp_burst" {
  count = var.create_key_pair ? 1 : 0

  key_name   = "${var.project_name}-key"
  public_key = var.ssh_public_key

  tags = var.common_tags
}

# VPN Gateway (if needed for NRP connectivity)
resource "aws_vpn_gateway" "nrp" {
  count = var.create_vpn_gateway ? 1 : 0

  vpc_id = var.create_vpc ? aws_vpc.nrp_burst[0].id : var.existing_vpc_id

  tags = merge(
    var.common_tags,
    {
      Name = "${var.project_name}-vpn-gateway"
    }
  )
}

# CloudWatch Log Group for instance logs
resource "aws_cloudwatch_log_group" "nrp_burst" {
  name              = "/nrp-aws-kip/instances"
  retention_in_days = var.log_retention_days

  tags = var.common_tags
}
