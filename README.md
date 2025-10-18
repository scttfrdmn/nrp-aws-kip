# NRP-AWS-KIP: National Research Platform AWS Kubernetes Instance Provider

A cloud bursting solution that integrates AWS compute capacity with the National Research Platform's Kubernetes clusters using the Virtual Kubelet pattern.

## Overview

This project enables NRP Kubernetes clusters to dynamically burst workloads to AWS when on-premises capacity is exhausted. It implements a Virtual Kubelet provider that translates Kubernetes pod specifications into AWS EC2 instances.

## Architecture

```
┌─────────────────────────────────────────┐
│  NRP Kubernetes Cluster                 │
│  ┌─────────────────────────────────┐   │
│  │  NRP-AWS-KIP Virtual Node       │   │
│  │  (Virtual Kubelet)              │   │
│  └──────────────┬──────────────────┘   │
└─────────────────┼───────────────────────┘
                  │
                  │ VPN/Direct Connect
                  │
┌─────────────────▼───────────────────────┐
│  AWS Account                            │
│  ┌──────────────────────────────────┐  │
│  │  VPC / Subnets                   │  │
│  │  ┌────────┐ ┌────────┐ ┌──────┐ │  │
│  │  │ EC2    │ │ EC2    │ │ EC2  │ │  │
│  │  │ (Pod)  │ │ (Pod)  │ │ (Pod)│ │  │
│  │  └────────┘ └────────┘ └──────┘ │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

## Features

- **Dynamic Pod-to-EC2 Mapping**: Translates K8s pod specs to appropriate EC2 instance types
- **Cost Optimization**: Support for Spot instances and automatic cleanup
- **Network Integration**: Seamless connectivity between NRP and AWS
- **Resource Management**: Respects pod resource requests/limits
- **Monitoring**: Metrics and logging for burst usage

## Project Structure

```
.
├── cmd/
│   └── nrp-aws-kip/          # Main application entry point
├── pkg/
│   ├── provider/             # Virtual Kubelet provider implementation
│   ├── config/               # Configuration management
│   └── controller/           # K8s controller logic
├── internal/
│   ├── aws/                  # AWS SDK integration
│   └── metrics/              # Prometheus metrics
├── deployments/
│   ├── kubernetes/           # K8s manifests, Helm charts
│   └── terraform/            # AWS infrastructure IaC
├── docs/                     # Additional documentation
└── scripts/                  # Build and deployment scripts
```

## Prerequisites

- Go 1.21+
- Access to NRP Kubernetes cluster
- AWS account with appropriate permissions
- VPN or Direct Connect between NRP and AWS (recommended)

## Quick Start

### 1. Configure AWS Credentials

```bash
export AWS_REGION=us-west-2
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
```

### 2. Deploy Infrastructure

```bash
cd deployments/terraform
terraform init
terraform apply
```

### 3. Build and Deploy

```bash
make build
kubectl apply -f deployments/kubernetes/
```

## Configuration

See `config.yaml.example` for configuration options including:
- AWS region and instance types
- VPC/subnet configuration
- Resource limits and quotas
- Cost controls

## Development

Built by AWS and San Diego State University for the National Research Platform.

### Building

```bash
make build
```

### Testing

```bash
make test
```

## License

Apache 2.0 License - See LICENSE file

## Contributing

Contributions welcome! Please see CONTRIBUTING.md for guidelines.

## Contact

- AWS Team: [your contact]
- SDSU Team: [sdsu contact]
